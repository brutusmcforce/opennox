package opennox

import (
	"context"
	"image"
	"strconv"
	"strings"

	"github.com/opennox/libs/client/seat"
	"github.com/opennox/libs/console"
	"github.com/opennox/libs/log"
	"github.com/spf13/viper"
)

const (
	noxDefaultWidth  = seat.DefaultWidth
	noxDefaultHeight = seat.DefaultHeight
	noxDefaultDepth  = seat.DefaultDepth
)

const (
	gammaMin = 0.1
	gammaDef = 1.0
	gammaMax = 3.0
	gammaInc = 0.1
)

const (
	configVideoUIScale = "video.ui_scale"

	videoUIScaleMin = 100
	videoUIScaleDef = 100
	videoUIScaleMax = 400
)

var (
	videoLog         = log.New("video")
	noxVideoModeMenu = image.Point{
		X: noxDefaultWidth,
		Y: noxDefaultHeight,
	}
	noxVideoMax = image.Point{
		X: noxMaxWidth,
		Y: noxMaxHeight,
	}
	noxVideoModes = []image.Point{
		{640, 480},
		{800, 600},
		{1024, 768},
		{1280, 720},
		{1920, 1080},
		{2560, 1440},
		{3840, 2160},
	}
	noxVideoGamma   = float32(gammaDef)
	noxVideoUIScale = videoUIScaleDef

	nox_video_cutSize = 100

	g_scaled_cfg     = 0
	g_fullscreen_cfg = 0
)

func init() {
	viper.SetDefault(configVideoUIScale, videoUIScaleDef)
	registerOnConfigRead(func() {
		videoSetUIScale(viper.GetInt(configVideoUIScale))
	})
	noxConsole.Register(&console.Command{
		Token: "ui_scale",
		Help:  "ui_scale [percent]: get or set UI scale in percent (100-400), applied on the next game start",
		Flags: console.ClientServer,
		Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			if len(tokens) == 0 {
				c.Printf(console.ColorLightYellow, "ui_scale = %d%%", videoGetUIScale())
				return true
			} else if len(tokens) > 1 {
				return false
			}
			v, err := strconv.Atoi(strings.TrimSuffix(tokens[0], "%"))
			if err != nil {
				return false
			}
			videoSetUIScaleOpts(v)
			c.Printf(console.ColorLightYellow, "ui_scale = %d%% (applied on the next game start)", videoGetUIScale())
			return true
		},
	})
	noxConsole.Register(&console.Command{
		Token:  "gamma",
		HelpID: "gammahelp",
		Flags:  console.ClientServer,
		Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			if len(tokens) == 0 {
				c.Printf(console.ColorLightYellow, "gamma = %d", getGammaSlider())
				return true
			} else if len(tokens) > 1 {
				return false
			}
			s := tokens[0]
			if len(s) == 0 {
				return false
			}
			add := 0
			if s[0] == '+' {
				add = +1
				s = s[1:]
			} else if s[0] == '-' {
				add = -1
				s = s[1:]
			}
			v, err := strconv.Atoi(s)
			if err != nil {
				return false
			}
			gamma := v
			if add != 0 {
				gamma = getGammaSlider() + add*gamma
			}
			setGammaSlider(gamma)
			return true
		},
	})
}

func videoModeID() int {
	mode := noxClient.videoGetGameMode()
	for i, res := range noxVideoModes {
		if mode == res {
			return i
		}
	}
	for i, res := range noxVideoModes {
		if mode.X == res.X {
			return i
		}
	}
	for i, res := range noxVideoModes {
		if mode.Y == res.Y {
			return i
		}
	}
	area := mode.X * mode.Y
	min := 0
	ind := -1
	for i, res := range noxVideoModes {
		area2 := res.X * res.Y
		if ind == -1 || abs(area-area2) < min {
			ind = i
			min = abs(area - area2)
		}
	}
	return ind
}

func videoModeByID(code byte) (image.Point, bool) {
	if int(code) >= len(noxVideoModes) {
		return image.Point{}, false
	}
	return noxVideoModes[code], true
}

func videoGetMenuMode() image.Point {
	return noxVideoModeMenu
}

func videoSetMenuMode(mode image.Point) {
	noxVideoModeMenu = mode
}

func videoSetMaxSize(sz image.Point) {
	noxVideoMax = sz
}

func videoGetMaxSize() image.Point {
	return noxVideoMax
}

func videoResizeView(mode image.Point) {
	max := videoGetMaxSize()
	if mode.X > max.X {
		mode.X = max.X
	}
	if mode.Y > max.Y {
		mode.Y = max.Y
	}
	videoSetWindowSize(mode)
}

func changeWindowedOrFullscreen() {
	cfgUpdateFullScreen()
}

func nox_video_setCutSize(v int) {
	if v < 65 {
		v = 65
	}
	if v > 100 {
		v = 100
	}
	nox_video_cutSize = v
}

func nox_video_getCutSize() int {
	return nox_video_cutSize
}

func getGammaSlider() int {
	return int(100 * (getGamma() - gammaMin) / gammaMax)
}

func setGammaSlider(v int) {
	fv := float32(v) / 100
	setGamma(gammaMin + gammaMax*fv)
}

func setGammaSliderOpts(v int) {
	setGammaSlider(v)
	viper.Set(configVideoGamma, getGamma())
	writeConfigLater()
}

func updateGamma(value int) {
	var dv float32
	if value > 0 {
		dv = +gammaInc
	} else if value < 0 {
		dv = -gammaInc
	}
	setGamma(getGamma() + dv)
}

func getGamma() float32 {
	return noxVideoGamma
}

func setGamma(v float32) {
	if v < gammaMin {
		v = gammaMin
	} else if v > gammaMax {
		v = gammaMax
	}
	if noxVideoGamma == v {
		return
	}
	videoLog.Printf("gamma: %v", v)
	noxVideoGamma = v
	noxClient.SetScreenGamma(v)
}

// videoGetUIScale returns UI scale in percent.
func videoGetUIScale() int {
	return noxVideoUIScale
}

// videoSetUIScale sets UI scale in percent. It will be applied on the next video mode reset.
func videoSetUIScale(perc int) {
	if perc < videoUIScaleMin {
		perc = videoUIScaleMin
	} else if perc > videoUIScaleMax {
		perc = videoUIScaleMax
	}
	if noxVideoUIScale == perc {
		return
	}
	videoLog.Printf("ui scale: %d%%", perc)
	noxVideoUIScale = perc
}

func videoSetUIScaleOpts(perc int) {
	videoSetUIScale(perc)
	viper.Set(configVideoUIScale, videoGetUIScale())
	writeConfigLater()
}

// videoUIBufferSize returns the size of the GUI buffer for a given window size and current UI scale.
// The GUI is rendered at this (lower) resolution and then upscaled to the window size,
// while the world is still rendered at full resolution (see video_uiscale.go).
func videoUIBufferSize(mode image.Point) image.Point {
	scale := videoGetUIScale()
	if scale <= 100 || mode.X <= 0 || mode.Y <= 0 {
		return mode
	}
	ratio := float64(mode.Y) / float64(mode.X)
	// Nox requires buffer width to be aligned to 32 pixels.
	w := (mode.X * 100 / scale) &^ 31
	if w < noxDefaultWidth {
		w = noxDefaultWidth
	}
	h := int(float64(w)*ratio + 0.5)
	if h < noxDefaultHeight {
		// keep the minimal height required by the GUI, adjust width to preserve aspect ratio
		h = noxDefaultHeight
		w = (int(float64(h)/ratio+31) &^ 31)
	}
	if w > mode.X || h > mode.Y {
		// the window is too small to scale it down
		return mode
	}
	return image.Point{X: w, Y: h}
}
