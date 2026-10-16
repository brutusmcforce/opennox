//go:build !server

package music

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"
	"github.com/brutusmcforce/opennox/legacy/common/alloc/handles"
	"github.com/brutusmcforce/opennox/legacy/timer"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/noxtest"
)

const musicTitle = 14 // title.wav

type testMusic struct {
	*Module
	ticks uint64
	state uint32 // 5d4594_816348
}

func newTestMusic(t *testing.T, dir string) *testMusic {
	t.Setenv("SDL_AUDIODRIVER", "dummy") // no sound
	handles.Init()
	tm := &testMusic{}
	timer.PlatformTicks = func() uint64 { return tm.ticks }
	t.Cleanup(func() { timer.PlatformTicks = nil })
	dig := ail.WaveOutOpen()
	if dig == 0 {
		t.Skip("no audio device")
	}
	var (
		v816368, v816372, v93156, v93160, v816340 uint32
		drv                                       ail.Driver
		counter                                   timer.TimerGroup
		master                                    = new(timer.TimerGroup)
	)
	counter.Init()
	master.Init()
	tm.Module = NewModule(dir, &v816368, &v816372, &drv, &v93156, &v93160, &v816340, &counter, &master, &tm.state,
		func() ail.Driver { return dig },
		func() bool { return false },
		func() string { return "" },
		func() uint64 { return tm.ticks },
	)
	tm.Init()
	return tm
}

func (tm *testMusic) update(ms uint64) {
	tm.ticks += ms
	ail.Serve() // as the game loop
	tm.Sub_43D2D0()
	tm.Update()
}

func TestMusicPauseResume(t *testing.T) {
	dir := noxtest.DataPath(t, "music")
	if _, err := os.Stat(ifs.Normalize(filepath.Join(dir, "title.wav"))); err != nil {
		t.Skip("no music file:", err)
	}
	tm := newTestMusic(t, dir)
	tm.SetNextMusic(MusicState{MusicIdx: musicTitle, Volume: 100, Position: 50000})
	tm.update(10)
	must.EqOp(t, 1, tm.state)
	s := tm.playingStream
	must.NonZero(t, s)
	must.EqOp(t, 4, s.Status())
	// the start position is used once, as in the original
	must.EqOp(t, 0, tm.blocks[tm.block_idx].Inner.Position)

	// a voice line pauses the music after fading it out
	tm.Sub_43DBD0()
	must.EqOp(t, 4, tm.state)
	tm.update(100)
	must.EqOp(t, 4, tm.state) // still fading out
	must.EqOp(t, 4, s.Status())
	for i := 0; i < 20 && tm.state == 4; i++ {
		tm.update(100)
	}
	must.EqOp(t, 3, tm.state)
	must.EqOp(t, 8, s.Status())

	// and resumes the same stream afterwards
	tm.Sub_43DBE0()
	must.EqOp(t, 1, tm.state)
	tm.update(10)
	must.EqOp(t, 1, tm.state)
	must.EqOp(t, s, tm.playingStream)
	must.EqOp(t, 4, s.Status())
}

func TestMusicMissingFile(t *testing.T) {
	tm := newTestMusic(t, t.TempDir())
	tm.SetNextMusic(MusicState{MusicIdx: musicTitle, Volume: 100})
	tm.update(10)
	// a track that cannot be opened is given up, not retried every frame
	must.EqOp(t, 0, tm.state)
	must.EqOp(t, 0, tm.blocks[tm.block_idx].Inner.MusicIdx)
}

// A position from a save is used once: when the track ends, it starts again from the beginning.
func TestMusicRestartFromStart(t *testing.T) {
	dir := noxtest.DataPath(t, "music")
	if _, err := os.Stat(ifs.Normalize(filepath.Join(dir, "title.wav"))); err != nil {
		t.Skip("no music file:", err)
	}
	tm := newTestMusic(t, dir)
	const samples = 1622015 // title.wav
	// 0.1 s before the end
	tm.SetNextMusic(MusicState{MusicIdx: musicTitle, Volume: 100, Position: samples - 22050/10})
	tm.update(10)
	first := tm.playingStream
	must.NonZero(t, first)
	deadline := time.Now().Add(5 * time.Second)
	for tm.playingStream == first && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond) // the audio device plays the stream in real time
		tm.update(10)
	}
	for (tm.state != 1 || tm.playingStream == first) && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
		tm.update(10)
	}
	must.EqOp(t, 1, tm.state)
	must.NotEqOp(t, first, tm.playingStream)
	must.Less(t, 22050, tm.playingStream.Position()) // within the first second
}
