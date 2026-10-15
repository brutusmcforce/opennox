//go:build !server

package opennox

import (
	"fmt"
	"image"

	"github.com/spf13/viper"

	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/legacy"
)

// OpenNox-specific options, shown in the advanced video options list (see gui_video_adv.go).
//
// Options with a range of values are sliders over a list of choices.
// Changes are saved to the config right away.

var (
	uiScaleChoices = []int{100, 125, 150, 175, 200, 225, 250, 275, 300, 325, 350, 375, 400}
	maxFPSChoices  = []int{30, 60, 75, 90, 100, 120, 144, 165, 180, 200, 240, 360}
)

func openNoxVideoOpts(c *Client) []*videoOpt {
	return []*videoOpt{
		{
			ID: 2060,
			Slider: &videoSlider{
				Name:    "UI Scale",
				Choices: uiScaleChoices,
				Get:     videoGetUIScale,
				Set: func(v int) {
					videoSetUIScaleOpts(v)
					_ = writeConfig()
					c.applyUIScaleInGame()
				},
				Format: func(v int) string { return fmt.Sprintf("%d%%", v) },
			},
		},
		{
			ID: 2061,
			// 30 (the tick rate) draws no frames between the ticks
			Slider: &videoSlider{
				Name:    "Frame Rate Limit",
				Choices: maxFPSChoices,
				Get:     func() int { return c.interp.maxFPS },
				Set: func(v int) {
					c.interp.setMaxFPS(v)
					viper.Set(configVideoMaxFPS, v)
					_ = writeConfig()
				},
			},
		},
		{
			ID: 2062,
			Label: func() string {
				return "Vertical Sync: " + vsyncMode()
			},
			Get: func() bool { return vsyncMode() != vsyncOff },
			Toggle: func() {
				cur := vsyncMode()
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
		{
			ID:     2065,
			Text:   "Auto-Equip Pickups",
			Get:    func() bool { return gameAutoEquip },
			Toggle: func() { setGameAutoEquip(!gameAutoEquip) },
		},
		{
			ID: 2066,
			Slider: &videoSlider{
				Name:    "Map Opacity",
				Choices: largeMapOpacityChoices,
				Get:     func() int { return legacy.MinimapLargeOpacity },
				Set: func(v int) {
					setLargeMapOpacity(v)
					viper.Set(configLargeMapOpacity, legacy.MinimapLargeOpacity)
					_ = writeConfig()
				},
				Format: func(v int) string { return fmt.Sprintf("%d%%", v) },
			},
		},
	}
}

// applyUIScaleInGame applies the UI scale right away, if the game is running.
// In menus, the UI scale is applied when the game starts.
//
// Only the GUI buffer is resized: the world is still drawn at the same resolution, so nothing else is affected.
// The GUI windows were placed for the old GUI size, so they are moved to keep their place on the screen.
func (c *Client) applyUIScaleInGame() {
	if nox_client_gui_flag_815132 != 0 || uiScaleState.inWorld || noxPixBuffer.img == nil || c.ClientPlayerUnit() == nil {
		return
	}
	world := c.videoGetGameMode()
	if max := videoGetMaxSize(); world.X > max.X || world.Y > max.Y {
		world = image.Point{X: min(world.X, max.X), Y: min(world.Y, max.Y)}
	}
	old := videoGetWindowSize()
	mode := videoUIBufferSize(world)
	if mode == old {
		return
	}
	videoLog.Printf("ui scale: resizing GUI from %v to %v", old, mode)
	if mode == world {
		world = image.Point{}
	}
	videoResizeView(mode)
	videoSetWorldSize(world)
	nox_video_initPixbuffer_486090(mode)
	sub_49F610(mode)
	relayoutGUI(c, old, mode)
}

// relayoutShift returns how far to move a window along one axis when the GUI size changes from old to cur.
// The window keeps its distance to what it's closest to: the start, the end or the center of the screen.
func relayoutShift(off, sz, old, cur int) int {
	toStart := abs(off)
	toEnd := abs(old - (off + sz))
	toCenter := abs(off + sz/2 - old/2)
	switch {
	case toCenter <= toStart && toCenter <= toEnd:
		return (cur - old) / 2
	case toStart <= toEnd:
		return 0
	}
	return cur - old
}

// relayoutGUI moves top-level GUI windows after the GUI size changes from old to cur.
//
// The game places windows only once, when they are created, so they must be moved to where they would be placed
// for the new size. There is no record of how each window was placed, so it is guessed from its position:
//   - windows that are still where their window file put them are not moved;
//   - other windows keep their distance to the closest of: the left/top edge, the right/bottom edge, the center.
//
// Windows of the same size placed in a row or a column (e.g. spell bar slots) are moved together, as one window.
// Empty placeholder windows move together with a window at the same position, if any.
func relayoutGUI(c *Client, old, cur image.Point) {
	if cur == old {
		return
	}
	type item struct {
		win  *gui.Window
		rect image.Rectangle // old placement
		pos  image.Point     // new position
		grp  int
	}
	isPlaceholder := func(it *item) bool {
		sz := it.rect.Size()
		return sz.X <= 1 || sz.Y <= 1
	}
	var list []*item
	for win := c.GUI.Last(); win != nil; win = win.Prev() {
		off := win.Offs()
		list = append(list, &item{win: win, rect: image.Rectangle{Min: off, Max: off.Add(win.Size())}})
	}
	// group windows of the same size that share a row or a column
	for i, it := range list {
		it.grp = i
	}
	var find func(i int) int
	find = func(i int) int {
		if list[i].grp != i {
			list[i].grp = find(list[i].grp)
		}
		return list[i].grp
	}
	for i, a := range list {
		if isPlaceholder(a) {
			continue
		}
		for j := i + 1; j < len(list); j++ {
			b := list[j]
			if b.rect.Size() != a.rect.Size() || (b.rect.Min.X != a.rect.Min.X && b.rect.Min.Y != a.rect.Min.Y) {
				continue
			}
			list[find(j)].grp = find(i)
		}
	}
	bounds := make(map[int]image.Rectangle)
	for i, it := range list {
		g := find(i)
		if r, ok := bounds[g]; ok {
			bounds[g] = r.Union(it.rect)
		} else {
			bounds[g] = it.rect
		}
	}
	// new positions
	bookList := legacy.SpellBookListWindow()
	for i, it := range list {
		win := it.win
		it.pos = it.rect.Min
		switch {
		case win == guiCon.root:
			continue // recreated below
		case it.rect.Size() == old:
			// full screen window
			it.pos = image.Point{}
			continue
		case win.ID() == 9200:
			// chat, see sub_46A730
			it.pos = image.Pt(cur.X/2, 2*cur.Y/3)
			continue
		case win.ID() == 9000:
			// quit menu
			it.pos = quitMenuPos(win, cur)
			continue
		case win == bookList:
			// see nox_xxx_bookInit_45B9D0
			it.pos = image.Pt(5, cur.Y-323)
			continue
		}
		file, fromFile := c.GUI.FileRect(win)
		fromFile = fromFile && file.Size() == it.rect.Size()
		b := bounds[find(i)]
		if !fromFile || file.Min.X != it.rect.Min.X {
			it.pos.X += relayoutShift(b.Min.X, b.Dx(), old.X, cur.X)
		}
		if !fromFile || file.Min.Y != it.rect.Min.Y {
			it.pos.Y += relayoutShift(b.Min.Y, b.Dy(), old.Y, cur.Y)
		}
	}
	// placeholder windows usually hold parts of another window at the same position, so they must move together
	for _, it := range list {
		if !isPlaceholder(it) {
			continue
		}
		for _, it2 := range list {
			if it2.rect.Min == it.rect.Min && !isPlaceholder(it2) && it2.win != guiCon.root {
				it.pos = it.rect.Min.Add(it2.pos.Sub(it2.rect.Min))
				break
			}
		}
	}
	for _, it := range list {
		if it.win == guiCon.root {
			continue
		}
		if it.rect.Size() == old {
			it.win.SizeVal = cur
		}
		it.win.SetPos(it.pos)
	}
	legacy.RelayoutQuickBar(old, cur)
	guiCon.Resize(cur)
}
