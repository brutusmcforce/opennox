//go:build !server

package ail

import (
	"math"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

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
// so closing a stream must free its SDL stream and logical device.
func TestStreamCloseFrees(t *testing.T) {
	path := testMusicFile(t)
	t.Setenv("SDL_AUDIODRIVER", "dummy") // no sound
	handles.Init()
	dig := WaveOutOpen()
	if dig == 0 {
		t.Skip("no audio device")
	}
	d := dig.get()
	for i := 0; i < 600; i++ { // many: leaked streams and devices would pile up
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

// Voices are released when the audio device closes (e.g. sound turned off in the options).
func TestSampleReleaseFrees(t *testing.T) {
	t.Setenv("SDL_AUDIODRIVER", "dummy") // no sound
	handles.Init()
	dig := WaveOutOpen()
	if dig == 0 {
		t.Skip("no audio device")
	}
	d := dig.get()
	for i := 0; i < 600; i++ { // many: leaked streams and devices would pile up
		s := dig.AllocateSample()
		must.NonZero(t, s, must.Sprintf("sample %d", i))
		d.doWork()
		s.Release()
		must.Nil(t, s.get())
	}
	d.mu.Lock()
	head := d.sampleHead
	d.mu.Unlock()
	must.Nil(t, head)
}

// Corrupt music files must not crash the game (opennox#681: index 253 in the step table).
func TestADPCMCorruptBlock(t *testing.T) {
	block := make([]byte, 1024)
	for i := range block {
		block[i] = byte(i*37 + 253)
	}
	block[2], block[6] = 253, 255
	must.NotEq(t, 0, len(decodeADPCMStereo(nil, block)))
	must.NotEq(t, 0, len(decodeADPCMMono(nil, block)))
}

// Sound effects (ADPCM blocks) are decoded, panned and queued, and play on the device.
func TestSamplePlays(t *testing.T) {
	t.Setenv("SDL_AUDIODRIVER", "dummy") // no sound
	handles.Init()
	dig := WaveOutOpen()
	if dig == 0 {
		t.Skip("no audio device")
	}
	d := dig.get()
	h := dig.AllocateSample()
	must.NonZero(t, h)
	defer h.Release()
	const block = 512
	h.SetType(4, 0) // mono ADPCM
	h.SetADPCMBlockSize(block)
	h.SetPlaybackRate(22050)
	h.Init()
	h.SetPan(0) // all left
	data := make([]byte, 4*block)
	for i := range data {
		data[i] = byte(i * 37)
	}
	n := h.BufferReady() // the game asks which buffer to fill first
	must.EqOp(t, 0, n)
	h.LoadBuffer(uint32(n), data)
	d.doWork()
	s := h.get()
	must.True(t, s.stream.Queued() > 0, must.Sprint("nothing queued"))
	must.False(t, s.stream.Paused(), must.Sprint("the stream must play"))
	// mono data is turned into stereo with the pan: pan 0 (-63/64, nearly all left) leaves the right channel quiet
	var left, right float64
	for i := 0; i+1 < len(s.stereoBuf); i += 2 {
		left += math.Abs(float64(s.stereoBuf[i]))
		right += math.Abs(float64(s.stereoBuf[i+1]))
	}
	must.True(t, left > 0, must.Sprint("no sound"))
	must.True(t, right < left/20, must.Sprintf("left %.0f, right %.0f", left, right))
}

// The movie player's stream counts the frames it played.
func TestPCMStreamPlayed(t *testing.T) {
	t.Setenv("SDL_AUDIODRIVER", "dummy") // no sound
	handles.Init()
	dig := WaveOutOpen()
	if dig == 0 {
		t.Skip("no audio device")
	}
	p := dig.OpenPCMStream(22050, 2)
	must.NotNil(t, p)
	defer p.Close()
	p.QueueStereo16(make([][2]int16, 22050)) // 1 s
	must.EqOp(t, 0, p.Played())
	must.EqOp(t, 22050, p.Queued())
	p.Start()
	deadline := time.Now().Add(3 * time.Second)
	for p.Played() == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	must.True(t, p.Played() > 0, must.Sprint("nothing played"))
	must.EqOp(t, 22050, p.Played()+p.Queued())
}
