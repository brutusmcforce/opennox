package opennox

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/spf13/viper"

	"github.com/brutusmcforce/opennox/libs/console"
	"github.com/brutusmcforce/opennox/libs/platform"
)

const configVideoShowFPS = "video.show_fps"

func init() {
	viper.SetDefault(configVideoShowFPS, false)
	registerOnConfigRead(func() {
		noxClient.fps.enabled = viper.GetBool(configVideoShowFPS)
	})
	noxConsole.Register(&console.Command{
		Token: "fps",
		Help:  "fps: toggle frame rate display in the upper left corner",
		Flags: console.ClientServer | console.NoDedicated,
		Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			if len(tokens) != 0 {
				return false
			}
			noxClient.fps.enabled = !noxClient.fps.enabled
			viper.Set(configVideoShowFPS, noxClient.fps.enabled)
			writeConfigLater()
			return true
		},
	})
}

// fpsOverlay shows the frame rate in the upper left corner.
type fpsOverlay struct {
	enabled bool
	frames  int
	since   time.Duration
	fps     float64
}

// countFrame must be called for each presented frame.
func (c *Client) countFrame() {
	f := &c.fps
	now := platform.Ticks()
	if f.since == 0 {
		f.since = now
	}
	f.frames++
	if dt := now - f.since; dt >= 500*time.Millisecond {
		f.fps = float64(f.frames) / dt.Seconds()
		f.frames = 0
		f.since = now
	}
}

// drawFPSOverlay draws the frame rate, if enabled.
func (c *Client) drawFPSOverlay() {
	if !c.fps.enabled {
		return
	}
	str := fmt.Sprintf("FPS: %.0f", c.fps.fps)
	r := c.r
	sz := r.GetStringSizeWrapped(nil, str, 0)
	const pad = 3
	d := r.Data()
	prevClip := d.Clip()
	d.SetClip(false)
	r.DrawRectFilledAlpha(0, 0, sz.X+2*pad, sz.Y+2*pad)
	d.SetTextColor(nox_color_white_2523948)
	r.DrawString(nil, str, image.Pt(pad, pad))
	d.SetClip(prevClip)
}
