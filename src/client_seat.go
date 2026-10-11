//go:build !server

package opennox

import (
	"image"
	"unicode/utf16"

	"github.com/spf13/viper"
	"github.com/tawesoft/golib/v2/dialog"

	"github.com/opennox/libs/client/seat/sdl"
	"github.com/opennox/libs/env"
	"github.com/opennox/libs/noximage"

	"github.com/opennox/opennox/v1/client/gpu"
	"github.com/opennox/opennox/v1/client/input"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/client/render"
	"github.com/opennox/opennox/v1/internal/version"
	"github.com/opennox/opennox/v1/legacy"
)

const (
	configVideoRenderer = "video.renderer"

	rendererSoftware = "software"
	rendererGL       = "gl"
)

func init() {
	viper.SetDefault(configVideoFiltering, true)
	viper.SetDefault(configVideoStretch, false)
	viper.SetDefault(configVideoRenderer, rendererGL)
}

// initGPU enables GPU rendering, if it's selected in the config.
func (c *Client) initGPU() {
	if env.IsE2E() || viper.GetString(configVideoRenderer) != rendererGL {
		return
	}
	g, err := gpu.New(c.Log)
	if err != nil {
		c.Log.Error("cannot initialize GPU renderer, using software rendering", "err", err)
		return
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
		if img != nil && uiScaleState.world == nil {
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
}

func (c *Client) initSeat(sz image.Point) error {
	sst, err := sdl.New(c.Log, "OpenNox "+version.ClientVersion(), sz)
	if err != nil {
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
	c.initGPU()
	if err != nil {
		return err
	}
	sst.SetGamma(getGamma())
	return nil
}

func (c *Client) freeSeat() {
	if c.gpu != nil {
		legacy.Nox_client_gpu_edgeDraw = nil
		c.r.SetGPU(nil)
		c.gpu.Close()
		c.gpu = nil
	}
	if c.Seat != nil {
		c.Seat.Close()
		c.Seat = nil
	}
}

func errorMessage(format string, args ...any) {
	dialog.Error(format, args...)
}
