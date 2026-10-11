//go:build !server

package opennox

import (
	"context"
	"strings"

	"github.com/spf13/viper"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/opennox/libs/console"
)

// Vertical synchronization.
//
// With vsync enabled, presented frames are synchronized with the display refresh, which removes tearing.
// The frame rate is then also limited by the refresh rate of the display.
// Adaptive vsync is the same, but frames are not delayed if they are late, which avoids stutter.

const configVideoVSync = "video.vsync"

const (
	vsyncOff      = "off"
	vsyncOn       = "on"
	vsyncAdaptive = "adaptive"
)

var vsyncModes = []string{vsyncOff, vsyncOn, vsyncAdaptive}

func init() {
	viper.SetDefault(configVideoVSync, vsyncOff)
	noxConsole.Register(&console.Command{
		Token: "vsync",
		Help:  "vsync [off|on|adaptive]: get or set vertical synchronization",
		Flags: console.ClientServer,
		Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			if len(tokens) > 1 {
				return false
			}
			if len(tokens) == 1 {
				mode := strings.ToLower(tokens[0])
				if !isVSyncMode(mode) {
					return false
				}
				setVSyncMode(mode)
			}
			c.Printf(console.ColorLightYellow, "vsync = %s", viper.GetString(configVideoVSync))
			return true
		},
	})
}

func isVSyncMode(mode string) bool {
	for _, m := range vsyncModes {
		if m == mode {
			return true
		}
	}
	return false
}

// setVSyncMode changes and saves the vsync mode. It's applied on the main loop, which owns the OpenGL context.
//
// It may be called from any goroutine, including the main loop itself (e.g. from the GUI),
// so the change is queued asynchronously instead of waiting for the main loop.
func setVSyncMode(mode string) {
	viper.Set(configVideoVSync, mode)
	writeConfigLater()
	go noxServer.QueueInLoop(context.Background(), applyVSync)
}

// applyVSync sets the swap interval of the current OpenGL context according to the config.
func applyVSync() {
	mode := strings.ToLower(viper.GetString(configVideoVSync))
	interval := 0
	switch mode {
	case vsyncOn:
		interval = 1
	case vsyncAdaptive:
		interval = -1
	}
	if err := sdl.GLSetSwapInterval(interval); err != nil {
		if interval != -1 {
			videoLog.Printf("cannot set vsync to %q: %v", mode, err)
			return
		}
		videoLog.Printf("adaptive vsync is not supported, using regular vsync: %v", err)
		if err := sdl.GLSetSwapInterval(1); err != nil {
			videoLog.Printf("cannot enable vsync: %v", err)
			return
		}
	}
	videoLog.Printf("vsync: %s", mode)
}
