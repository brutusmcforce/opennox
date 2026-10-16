//go:build !server

package opennox

import (
	"errors"
	"fmt"
	"image"
	"unicode/utf16"

	"github.com/spf13/viper"

	libsdl "github.com/brutusmcforce/opennox/libs/sdl"

	"github.com/brutusmcforce/opennox/libs/client/seat/sdl"
	"github.com/brutusmcforce/opennox/libs/env"
	"github.com/brutusmcforce/opennox/libs/noximage"

	"github.com/brutusmcforce/opennox/client/gpu"
	"github.com/brutusmcforce/opennox/client/input"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/client/render"
	"github.com/brutusmcforce/opennox/internal/version"
	"github.com/brutusmcforce/opennox/legacy"
)

func init() {
	viper.SetDefault(configVideoFiltering, true)
	viper.SetDefault(configVideoStretch, false)
}

// initGPU initializes the OpenGL renderer. All drawing is done on the GPU.
func (c *Client) initGPU() error {
	g, err := gpu.New(c.Log)
	if err != nil {
		return fmt.Errorf("cannot initialize OpenGL renderer (OpenGL 3.3 is required): %w", err)
	}
	c.gpu = g
	c.r.SetGPU(g)
	legacy.Nox_client_gpu_edgeDraw = func(img *noxrender.Image, pos image.Point, left, right [3]int, clip image.Rectangle, cut int) bool {
		// wall light values are 8.8 fixed point multipliers
		conv := func(v [3]int) [3]float32 {
			return [3]float32{float32(v[0]) / 256, float32(v[1]) / 256, float32(v[2]) / 256}
		}
		return c.r.DrawImageGradientGPU(img, pos, cut, conv(left), conv(right), clip)
	}
	c.Win.SetPresentFunc(func(img *noximage.Image16, view image.Rectangle, win image.Point) {
		g.Gamma = getGamma()
		g.Filtering = c.Win.GetFiltering()
		g.Present(img, view, win)
		c.updateSysCursor()
		switch {
		case img == movieFrame:
			frameStats.presented(frameMovie)
		case nox_client_gui_flag_815132 != 0:
			frameStats.presented(frameMenu)
		default:
			frameStats.presented(frameGame)
		}
		if img != nil && img != movieFrame && uiScaleState.world == nil {
			// Menus are drawn at a low resolution and stretched to the window.
			// Draw them at the window resolution instead, so that the text stays smooth.
			// It's not done in game, because the world drawing expects the target to match the image.
			sz := image.Point{}
			if nox_client_gui_flag_815132 != 0 && view.Dx() > img.Rect.Dx() && view.Dy() > img.Rect.Dy() {
				sz = view.Size()
			}
			g.SetTargetSize(img, sz)
		}
	})
	return nil
}

func (c *Client) initSeat(sz image.Point) error {
	sst, err := sdl.New(c.Log, "OpenNox "+version.ClientVersion(), sz)
	if errors.Is(err, sdl.ErrOpenGL) {
		// opennox#692: old GPUs (e.g. Intel before Sandy Bridge) only have OpenGL 2.1; Mesa's software renderer
		// has 3.3, and it's fast enough for the game. It only returns if the restart can't be done.
		c.Log.Warn("no OpenGL 3.3, restarting with software rendering", "err", err)
		restartWithSoftwareGL()
		errorMessage("OpenNox needs OpenGL 3.3, which your graphics card or its driver doesn't support. Updating the graphics driver may help.\n\n%s", err.Error())
		return err
	}
	if err != nil {
		// e.g. SDL is missing: the Windows and macOS apps have no console to print it to
		errorMessage("%s", err.Error())
		return err
	}
	c.Seat = sst
	if env.IsE2E() {
		c.Seat = e2eWrapSeat(c.Seat)
	}
	c.Win, err = render.New(c.Seat)
	if err != nil {
		_ = c.Seat.Close()
		return err
	}

	c.initSysCursor(c.Seat)
	// game controllers act as mouse and keyboard (after the system cursor, which needs the real seat)
	c.gamepad = newGamepadSeat(c.Seat, c)
	c.Seat = c.gamepad

	inp := input.New(c.Log, c.Seat, false, c.Strings().Lang())
	c.Inp = inp

	inp.OnQuit(mainloopStop)
	inp.OnToggleFullScreen(c.Win.ToggleWindowMode)
	inp.OnKeyPress(gameexOnKeyboardPress)
	inp.OnMouseWheel(func(dv int) {
		// mix event handler is triggered only for wheel events
		call_OnLibraryNotice_265(dv)
	})
	inp.OnInputString(func(str string) {
		for _, c := range utf16.Encode([]rune(str)) {
			legacy.NoxInputOnChar(c)
		}
	})

	c.Win.OnViewResize(inp.SetWinSize)
	OnPixBufferResize(inp.SetDrawWinSize)

	c.Win.SetFiltering(viper.GetBool(configVideoFiltering))
	c.Win.SetStretched(viper.GetBool(configVideoStretch))
	if err := c.initGPU(); err != nil {
		return err
	}
	applyVSync()
	sst.SetGamma(getGamma())
	return nil
}

func (c *Client) freeSeat() {
	frameStats.finish()
	if c.gpu != nil {
		legacy.Nox_client_gpu_edgeDraw = nil
		c.r.SetGPU(nil)
		c.gpu.Close()
		c.gpu = nil
	}
	if c.Seat != nil {
		c.freeSysCursor()
		c.Seat.Close()
		c.Seat = nil
	}
}

// errorMessage shows an error in a message box. SDL does it natively on every platform we ship; the fallback is
// for the errors that SDL itself cannot be loaded for.
func errorMessage(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if err := libsdl.ShowSimpleMessageBox(libsdl.MESSAGEBOX_ERROR, "OpenNox", msg); err == nil {
		return
	}
	nativeErrorMessage(msg)
}
