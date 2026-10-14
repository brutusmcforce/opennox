//go:build !server

package ail

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/legacy/common/alloc/handles"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/noxtest"
)

func testMusicFile(t testing.TB) string {
	path := filepath.Join(noxtest.DataPath(t, "music"), "title.wav")
	if _, err := os.Stat(ifs.Normalize(path)); err != nil {
		t.Skip("no music file:", err)
	}
	return path
}

// The music opens a new stream for every track, after every save load and voice line,
// so closing a stream must free its OpenAL source.
func TestStreamCloseFrees(t *testing.T) {
	path := testMusicFile(t)
	t.Setenv("ALSOFT_DRIVERS", "null") // no sound
	handles.Init()
	dig := WaveOutOpen()
	if dig == 0 {
		t.Skip("no OpenAL device")
	}
	d := dig.get()
	for i := 0; i < 600; i++ { // more than OpenAL Soft's 256 sources
		s := dig.OpenStream(path)
		must.NonZero(t, s, must.Sprintf("stream %d", i))
		s.Start()
		d.doWork()
		must.EqOp(t, 4, s.Status())
		s.Pause(true)
		must.EqOp(t, 8, s.Status())
		d.doWork()
		s.Pause(false)
		must.EqOp(t, 4, s.Status())
		must.NoError(t, s.Close())
		must.EqOp(t, 2, s.Status())
	}
	d.mu.Lock()
	head := d.streamHead
	d.mu.Unlock()
	must.Nil(t, head)
}

// Positions saved with the music state must resume at the same sample.
func TestADPCMSeek(t *testing.T) {
	path := testMusicFile(t)
	decode := func(r *wavReader, n int) []int16 {
		var out []int16
		buf := make([]int16, 16*1024)
		for len(out) < n {
			k, ok := r.Decode(buf)
			out = append(out, buf[:k]...)
			if !ok {
				break
			}
		}
		return out[:min(n, len(out))]
	}
	r, err := openWav(path)
	must.NoError(t, err)
	all := decode(r, 400000)
	r.Close()
	for _, pos := range []int{0, 1, 2040, 2041, 2042, 100000, 150000} {
		r, err := openWav(path)
		must.NoError(t, err)
		r.Seek(pos)
		start := r.Position() // the start of the block holding pos
		ch := r.Channels()
		got := decode(r, 4096)
		r.Close()
		d := r.dec.(*adpcmDecoder)
		spb := (d.blockSize/ch-4)*2 + 1
		must.True(t, start <= pos && pos < start+spb, must.Sprintf("pos %d: start %d", pos, start))
		must.True(t, slices.Equal(all[start*ch:start*ch+len(got)], got), must.Sprintf("pos %d: start %d", pos, start))
	}
}
