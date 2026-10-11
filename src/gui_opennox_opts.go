//go:build !server

package opennox

import (
	"fmt"

	"github.com/spf13/viper"
)

// OpenNox-specific video options, shown in the advanced video options list (see gui_video_adv.go).
//
// Options with multiple values are shown as a check box with the current value in the label;
// clicking it switches to the next value. Changes are saved to the config right away.

var (
	uiScaleChoices = []int{100, 125, 150, 200, 250, 300}
	maxFPSChoices  = []int{30, 60, 120, 144, 165, 240}
)

// nextChoice returns the value that follows the current one in the list, wrapping around.
func nextChoice(list []int, cur int) int {
	for _, v := range list {
		if v > cur {
			return v
		}
	}
	return list[0]
}

func openNoxVideoOpts(c *Client) []*videoOpt {
	return []*videoOpt{
		{
			ID: 2060,
			Label: func() string {
				return fmt.Sprintf("UI Scale: %d%%", videoGetUIScale())
			},
			Get: func() bool { return videoGetUIScale() > 100 },
			Toggle: func() {
				videoSetUIScaleOpts(nextChoice(uiScaleChoices, videoGetUIScale()))
				_ = writeConfig()
			},
		},
		{
			ID: 2061,
			Label: func() string {
				if c.interp.maxFPS <= int(noxServer.TickRate()) {
					return "Frame Rate Limit: 30 (no smoothing)"
				}
				return fmt.Sprintf("Frame Rate Limit: %d", c.interp.maxFPS)
			},
			Get: func() bool { return c.interp.enabled() },
			Toggle: func() {
				v := nextChoice(maxFPSChoices, c.interp.maxFPS)
				c.interp.setMaxFPS(v)
				viper.Set(configVideoMaxFPS, v)
				_ = writeConfig()
			},
		},
		{
			ID: 2062,
			Label: func() string {
				return "Vertical Sync: " + viper.GetString(configVideoVSync)
			},
			Get: func() bool { return viper.GetString(configVideoVSync) != vsyncOff },
			Toggle: func() {
				cur := viper.GetString(configVideoVSync)
				next := vsyncModes[0]
				for i, m := range vsyncModes {
					if m == cur {
						next = vsyncModes[(i+1)%len(vsyncModes)]
					}
				}
				setVSyncMode(next)
				_ = writeConfig()
			},
		},
		{
			ID:   2063,
			Text: "Show FPS",
			Get:  func() bool { return c.fps.enabled },
			Toggle: func() {
				c.fps.enabled = !c.fps.enabled
				viper.Set(configVideoShowFPS, c.fps.enabled)
				_ = writeConfig()
			},
		},
		{
			ID:   2064,
			Text: "Smooth Fonts (restart)",
			Get:  func() bool { return viper.GetBool(configVideoSmoothFonts) },
			Toggle: func() {
				viper.Set(configVideoSmoothFonts, !viper.GetBool(configVideoSmoothFonts))
				_ = writeConfig()
			},
		},
	}
}
