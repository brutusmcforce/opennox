//go:build !server

package noxmovie

import (
	"errors"
	"image"
	"image/draw"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brutusmcforce/opennox/libs/client/seat"
	"github.com/brutusmcforce/opennox/libs/noximage"

	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"

	"github.com/brutusmcforce/opennox/libs/env"
	"github.com/brutusmcforce/opennox/vqa-decode/movies"
)

type Frame struct {
	Image *noximage.Image16
	Audio [][2]int16
}

// DisplayFunc shows a movie frame on the screen.
type DisplayFunc func(img *noximage.Image16)

type MoviePlayer struct {
	queue     chan *Frame
	movie     *movies.VqaFile
	file      io.Closer
	seat      seat.Seat
	oldInputs seat.InputConfig
	display   DisplayFunc
	audioDrv  ail.Driver
	audio     *ail.PCMStream // nil in E2E tests: the movie is skipped
	audioGain float32
	playing   atomic.Bool // Play owns the audio stream and frees it
	stop      chan struct{}
}

func (player *MoviePlayer) Start() {
	go player.decode()
}

func (player *MoviePlayer) decode() {
	for {
		select {
		case <-player.stop:
			return
		default:
		}
		frame, samples, err := player.movie.DecodeNextFrame()
		if err != nil {
			close(player.queue)
			return
		}
		img16 := noximage.NewImage16(frame.Rect)
		draw.Draw(img16, img16.Rect, frame, image.Pt(0, 0), draw.Src)
		select {
		case <-player.stop:
			return
		case player.queue <- &Frame{
			Image: img16,
			Audio: samples,
		}:
		}
	}
}

func NewPlayerWithHandle(file io.ReadSeekCloser, mvSeat seat.Seat, audioDrv ail.Driver, display DisplayFunc) (plr *MoviePlayer, err error) {
	vqa, err := movies.OpenMovieWithHandle(file)
	if err != nil {
		return nil, err
	}
	queue := make(chan *Frame, vqa.Header.Fps*2)
	stop := make(chan struct{})

	audio := audioDrv.OpenPCMStream(int(vqa.Header.SampleRate), 2)
	if audio == nil && !env.IsE2E() {
		return nil, errors.New("cannot play movie: no audio")
	}

	oldInputs := mvSeat.ReplaceInputs(nil)

	var player = &MoviePlayer{
		queue:     queue,
		movie:     vqa,
		file:      file,
		seat:      mvSeat,
		display:   display,
		audioDrv:  audioDrv,
		audio:     audio,
		stop:      stop,
		oldInputs: oldInputs,
		audioGain: 1,
	}

	// TODO: actually replace the preexisting events
	mvSeat.OnInput(func(ev seat.InputEvent) {
		switch evt := ev.(type) {
		case seat.WindowEvent:
			switch evt {
			case seat.WindowClosed:
				// If user tries to close the window, stop the loop.
				player.Close()
			}
		case *seat.KeyboardEvent:
			if evt.Pressed {
				// Stop the loop on keyboard key press as well.
				player.Close()
			}
		case *seat.MouseButtonEvent:
			if evt.Pressed {
				// Stop the loop on mouse key press as well.
				player.Close()
			}
		}
	})

	return player, nil
}

func NewPlayer(fname string, mvSeat seat.Seat, audioDrv ail.Driver, display DisplayFunc) (plr *MoviePlayer, err error) {
	vqaFile, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	return NewPlayerWithHandle(vqaFile, mvSeat, audioDrv, display)
}

func (player *MoviePlayer) Close() {
	select {
	case <-player.stop:
		return
	default:
	}
	close(player.stop)
	player.file.Close()
	if player.audio != nil {
		// silence now; the stream is freed once the audio goroutine of Play is done with it
		player.audio.Pause()
		if !player.playing.Load() {
			player.audio.Close()
		}
	}
	player.seat.ReplaceInputs(player.oldInputs)
}

func (player *MoviePlayer) SetAudioGain(v float32) {
	player.audioGain = v
}

func (player *MoviePlayer) Play() {
	// TODO: if we got no audio, we could probably still keep displaying frames...
	// Note that the frame delays are all calculated based on sound position
	audio := player.audio
	if audio == nil {
		player.Close()
		return
	}
	player.playing.Store(true)
	defer audio.Close() // runs after wg.Wait below: the audio goroutine is done with the stream
	audio.SetGain(player.audioGain)

	framesChan := make(chan *Frame, player.movie.Header.Fps)

	select {
	case <-player.stop:
		return
	case <-time.After(time.Second / time.Duration(player.movie.Header.Fps/5)):
	}
	// This forces to wait for the audio thread to exit.
	// Not closing it properly may lead to crashes later during the game.
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	// Audio thread: queue the audio of each frame, then pass the frame on to be displayed
	go func() {
		defer wg.Done()
		started := false
		for {
			select {
			case <-player.stop:
				close(framesChan)
				return
			case newFrame, ok := <-player.queue:
				if !ok {
					// the whole file is decoded: let the main loop finish the playback
					close(framesChan)
					return
				}
				audio.QueueStereo16(newFrame.Audio)
				if !started {
					audio.Start()
					started = true
				}
				select {
				case <-player.stop:
					return
				case framesChan <- newFrame:
				}
			}
		}
	}()

	var (
		currentFrame              = 0
		samplesWithCommitedFrames = 0
		sampleDuration            = time.Second / time.Duration(player.movie.Header.SampleRate)
	)

loop:
	for {
		// Process input events.
		player.seat.InputTick()
		// A key may have closed the player. Check it before the select below:
		// it picks randomly when a frame is ready too.
		select {
		case <-player.stop:
			break loop
		default:
		}

		var nextFrame *Frame
		// Run until movie player tells us to stop.
		select {
		case <-player.stop:
			break loop
		case newFrame, ok := <-framesChan:
			if !ok {
				// The whole file is displayed: wait until the queued audio has played
				timeToSleepBeforeAudioEnd := time.Duration(audio.Queued()) * sampleDuration
				select {
				case <-player.stop:
				case <-time.After(timeToSleepBeforeAudioEnd):
				}
				break loop
			}
			nextFrame = newFrame
		}

		// Run audio callbacks.
		ail.Serve()

		// The frames follow the audio: wait until the audio played up to this frame
		currentPosition := time.Duration(audio.Played()) * sampleDuration
		expectedPosition := time.Duration(samplesWithCommitedFrames) * sampleDuration
		if sleep := expectedPosition - currentPosition; sleep > 0 {
			select {
			case <-player.stop:
				break loop
			case <-time.After(sleep):
			}
		}

		player.display(nextFrame.Image)
		samplesWithCommitedFrames += len(nextFrame.Audio)
		currentFrame++
	}
}
