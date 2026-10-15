package opennox

import (
	"fmt"
	"image/color"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/noxfont"
	"github.com/brutusmcforce/opennox/libs/strman"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy"
)

type videoOpt struct {
	ID     uint
	Flag   noxflags.EngineFlag
	Bool   *bool
	CFlag  *uint32
	TextID strman.ID
	Text   string
	Get    func() bool
	Toggle func()
	// Label returns a dynamic text of the option, e.g. with its current value.
	Label  func() string
	Def    bool
	Hidden bool
	// Slider shows the option as a label with the value and a slider below it, instead of a check box.
	Slider *videoSlider
}

// videoSlider is an option that selects one of a list of values with a slider.
type videoSlider struct {
	Name    string
	Choices []int
	Get     func() int
	Set     func(v int)
	// Format returns the text of a value; by default, the number itself.
	Format func(v int) string
}

// index returns the index of the choice closest to v.
func (s *videoSlider) index(v int) int {
	best := 0
	for i, c := range s.Choices {
		if abs(c-v) < abs(s.Choices[best]-v) {
			best = i
		}
	}
	return best
}

func (s *videoSlider) label() string {
	v := s.Get()
	if s.Format != nil {
		return s.Name + ": " + s.Format(v)
	}
	return fmt.Sprintf("%s: %d", s.Name, v)
}

// The slider runs over twice the number of steps, so that a click selects the closest choice
// (the slider truncates) and the arrow keys (which move by 2) select the next one.
func sliderPos(i int) int   { return 2 * i }
func sliderIndex(v int) int { return (v + 1) / 2 }

// advOptsSliderLabelID is added to a slider option's ID to get the ID of its label.
const advOptsSliderLabelID = 100

type guiAdvOptions struct {
	cli                          *Client
	nox_win_advVideoOpts_1522600 *gui.Window
	noxVideoAdvList              []*videoOpt
	noxVideoAdvOpts              map[uint]*videoOpt

	// options are split into pages when they don't fit the panel
	pages    [][]*gui.Window
	page     int // kept while the game runs
	pageNext *gui.Window
}

// advOptsNextPageID is the ID of the button that switches option pages.
const advOptsNextPageID = 2097

func (c *guiAdvOptions) Init(cli *Client) {
	c.cli = cli
	c.noxVideoAdvOpts = make(map[uint]*videoOpt)
	c.noxVideoAdvList = []*videoOpt{
		{ID: 2051, Get: c.cli.GetFiltering, Toggle: c.cli.ToggleFiltering, Text: "Smooth image", TextID: "AdVidOpt.wnd:Filtering"},
		{ID: 2050, Get: c.cli.GetStretch, Toggle: c.cli.ToggleStretch, Text: "Stretch image", TextID: "AdVidOpt.wnd:Stretched"},
		{ID: 2010, Bool: &legacyUnlockSurface, TextID: "AdVidOpt.wnd:ClipWalls", Def: true, Hidden: true},
		{ID: 2012, Flag: noxflags.EngineSoftShadowEdge, TextID: "AdVidOpt.wnd:GouradShading"},
		{ID: 2014, Bool: &guiCon.translucent, TextID: "AdVidOpt.wnd:TranslucentConsole"},
		{ID: 2015, Bool: &c.cli.r.Part.RenderGlow, TextID: "AdVidOpt.wnd:RenderGlow"},
		{ID: 2016, CFlag: legacy.Get_nox_client_fadeObjects_80836_ptr(), TextID: "AdVidOpt.wnd:FadeObjects"},
		{ID: 2017, Bool: &nox_client_showTooltips_80840, TextID: "AdVidOpt.wnd:ShowTooltips"},
		{ID: 2020, Bool: &nox_client_drawFrontWalls_80812, TextID: "AdVidOpt.wnd:DrawFrontWalls"},
		{ID: 2021, CFlag: legacy.Get_nox_client_translucentFrontWalls_805844_ptr(), TextID: "AdVidOpt.wnd:TranslucentFrontWalls"},
		{ID: 2031, CFlag: legacy.Get_nox_client_highResFloors_154952_ptr(), TextID: "AdVidOpt.wnd:InterlacedFloors"},
		{ID: 2032, Bool: &nox_client_lockHighResFloors_1193152, TextID: "AdVidOpt.wnd:LockHiResFloors"},
		{ID: 2033, Bool: &nox_client_texturedFloors2_154960, TextID: "AdVidOpt.wnd:FlatShadedFloors"},
		{ID: 2052, Flag: noxflags.EngineNoSoftLights, Text: "Disable Soft Light", TextID: "AdVidOpt.wnd:NoSoftLights"},
		//{ID: 2041, CFlag: &C.nox_client_renderBubbles_80844, TextID: "AdVidOpt.wnd:RenderBubbles"},
		{ID: 2040, Bool: &nox_client_renderGUI_80828, TextID: "AdVidOpt.wnd:RenderGUI"},
	}
	c.noxVideoAdvList = append(c.noxVideoAdvList, openNoxVideoOpts(c.cli)...)
	for _, opt := range c.noxVideoAdvList {
		if opt.ID == 0 {
			panic("ID must be set")
		}
		c.noxVideoAdvOpts[opt.ID] = opt
	}
}

func (c *Client) NewGUIAdvOptsOn(par *gui.Window) {
	c.guiAdv.NewOn(par)
}

func (c *Client) GUIAdvVideoOptsLoad() {
	c.guiAdv.nox_client_advVideoOptsLoad()
}

func (c *guiAdvOptions) NewOn(par *gui.Window) {
	c.nox_win_advVideoOpts_1522600 = c.newAdvVideoOpts()
	c.nox_win_advVideoOpts_1522600.SetParent(par)
	c.nox_client_advVideoOptsLoad()
}

func (c *guiAdvOptions) nox_client_advVideoOptsLoad() {
	for id, opt := range c.noxVideoAdvOpts {
		if opt.Slider != nil {
			c.syncSlider(opt)
			continue
		}
		if w := c.nox_win_advVideoOpts_1522600.ChildByID(id); w != nil {
			v := opt.Def
			if opt.Get != nil {
				v = opt.Get()
			} else if opt.Flag != 0 {
				v = noxflags.HasEngine(opt.Flag)
			} else if opt.Bool != nil {
				v = *opt.Bool
			} else if opt.CFlag != nil {
				v = *opt.CFlag != 0
			}
			w.DrawData().Field0Set(0x4, v)
			if opt.Label != nil {
				w.DrawData().SetText(opt.Label())
			}
		}
	}
}

// syncSlider moves the slider of an option to its current value and updates its label.
func (c *guiAdvOptions) syncSlider(opt *videoOpt) {
	root := c.nox_win_advVideoOpts_1522600
	s := opt.Slider
	if w := root.ChildByID(opt.ID); w != nil {
		w.Func94(gui.AsWindowEvent(0x400A, uintptr(sliderPos(s.index(s.Get()))), 0))
	}
	if w := root.ChildByID(opt.ID + advOptsSliderLabelID); w != nil {
		w.Func94(&gui.StaticTextSetText{Str: s.label(), Val: -1})
	}
}

func (c *guiAdvOptions) nox_client_advVideoOptsProc_4CB5D0(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	switch ev := ev.(type) {
	case *WindowEvent0x4009:
		// a slider was moved (sent when the mouse button is released, or by the arrow keys)
		opt, ok := c.noxVideoAdvOpts[ev.Win.ID()]
		if !ok || opt.Slider == nil {
			return nil
		}
		s := opt.Slider
		i := min(max(sliderIndex(ev.Val), 0), len(s.Choices)-1)
		if v := s.Choices[i]; v != s.Get() {
			s.Set(v)
			clientPlaySoundSpecial(sound.SoundButtonPress, 100)
		}
		c.syncSlider(opt)
		return nil
	case *WindowEvent0x4005:
		clientPlaySoundSpecial(sound.SoundShellSelect, 100)
		return gui.RawEventResp(1)
	case *WindowEvent0x4007:
		clientPlaySoundSpecial(sound.SoundButtonPress, 100)
		id := ev.Win.ID()
		switch id {
		case 2010:
			// always enabled
			legacyUnlockSurface = true
			win.DrawData().Field0Set(0x4, true)
			return nil
		case 2033:
			v := nox_client_texturedFloors2_154960
			nox_client_texturedFloors2_154960 = !v
			nox_client_texturedFloors_154956 = !v
			nox_xxx_tileSetDrawFn_481420()
			legacy.Set_dword_5d4594_1193156(0)
			return nil
		case advOptsNextPageID:
			c.setPage(c.page + 1)
			return nil
		case 2099:
			detectBestVideoSettings()
			c.nox_client_advVideoOptsLoad()
			if par := win.Parent(); par != nil {
				par.Func94(ev)
			}
			return nil
		}
		if opt, ok := c.noxVideoAdvOpts[id]; ok {
			if opt.Toggle != nil {
				opt.Toggle()
			} else if opt.Flag != 0 {
				noxflags.ToggleEngine(opt.Flag)
			} else if opt.Bool != nil {
				v := *opt.Bool
				*opt.Bool = !v
			} else if opt.CFlag != nil {
				v := *opt.CFlag != 0
				*opt.CFlag = uint32(bool2int(!v))
			}
		}
		switch id {
		case 2020:
			legacy.Sub_49B3C0()
		case 2040:
			nox_xxx_xxxRenderGUI_587000_80832 = nox_client_renderGUI_80828
		}
		if opt, ok := c.noxVideoAdvOpts[id]; ok && opt.Label != nil {
			// multi-value options update the label and keep the check box state in sync with the value
			c.nox_client_advVideoOptsLoad()
			// the check box flips its state after this event is handled, so set the opposite state here
			ev.Win.DrawData().Field0Set(0x4, !opt.Get())
		}
		return nil
	}
	return nil
}

func (c *guiAdvOptions) newAdvVideoOpts() *gui.Window {
	draw, drawFree := gui.NewWindowData()
	defer drawFree()

	draw.Status = gui.StatusEnabled | gui.StatusNoFocus

	draw.SetHighlightColor(noxcolor.RGB5551Color(192, 128, 128))
	draw.SetTextColor(noxcolor.RGB5551Color(240, 180, 42))
	draw.SetEnabledColor(color.Transparent)
	draw.SetDisabledColor(color.Transparent)
	draw.SetSelectedColor(color.Transparent)
	draw.SetBackgroundColor(color.Transparent)
	draw.SetFont(c.cli.r.Fonts.FontPtrByName(noxfont.LargeName))

	const (
		width  = 184
		height = 315
	)

	root := c.cli.GUI.NewUserWindow(nil, 2000, draw.Status, 280, 38, width, height, draw, c.nox_client_advVideoOptsProc_4CB5D0)

	const (
		pad     = 3
		cheight = 15
		crow    = 20
	)
	// a slider takes two rows: its label and the slider
	rows := func(opt *videoOpt) int {
		if opt.Slider != nil {
			return 2
		}
		return 1
	}
	total := 0
	for _, opt := range c.noxVideoAdvList {
		if !opt.Hidden {
			total += rows(opt)
		}
	}
	// options that don't fit go to more pages, with a button in the last row to switch between them
	perPage := (height - 2*pad) / crow
	if total > perPage {
		perPage--
	}
	c.pages = nil
	c.pageNext = nil

	row := perPage
	for _, opt := range c.noxVideoAdvList {
		if opt.Hidden {
			continue
		}
		if row+rows(opt) > perPage {
			c.pages = append(c.pages, nil)
			row = 0
		}
		y := pad + crow*row
		row += rows(opt)
		page := &c.pages[len(c.pages)-1]
		if s := opt.Slider; s != nil {
			label := c.cli.GUI.NewStaticText(root, opt.ID+advOptsSliderLabelID, pad, y, width-2*pad, cheight, false, false, s.label())
			slider := NewHorizontalSlider(root, opt.ID, pad+12, y+crow+1, width-2*pad-24, 14, 0, sliderPos(len(s.Choices)-1))
			slider.Func94(gui.AsWindowEvent(0x400A, uintptr(sliderPos(s.index(s.Get()))), 0))
			*page = append(*page, label, slider)
			continue
		}
		str := opt.Text
		if opt.Label != nil {
			str = opt.Label()
		} else if opt.TextID != "" {
			if e, ok := c.cli.Strings().GetInFile(opt.TextID, "psscript.c"); ok {
				str = e.Value().Str
			}
		}
		win := NewCheckBox(c.cli.GUI, root, opt.ID, pad, y, width, cheight, str)
		*page = append(*page, win)
	}
	if len(c.pages) > 1 {
		c.pageNext = NewTextButton(c.cli.GUI, root, advOptsNextPageID, pad, pad+crow*perPage, width-2*pad, cheight, "")
	}
	c.setPage(c.page)
	return root
}

// setPage shows the given page of options (wrapping around) and hides the others.
func (c *guiAdvOptions) setPage(page int) {
	if len(c.pages) == 0 {
		return
	}
	c.page = page % len(c.pages)
	for i, list := range c.pages {
		for _, win := range list {
			win.SetHidden(i != c.page)
		}
	}
	if c.pageNext != nil {
		c.pageNext.DrawData().SetText(fmt.Sprintf("Next Page (%d/%d)", c.page+1, len(c.pages)))
	}
}
