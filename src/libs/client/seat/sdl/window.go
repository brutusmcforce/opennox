package sdl

import (
	"errors"
	"fmt"
	"image"
	"log/slog"
	"os"

	"github.com/brutusmcforce/opennox/libs/client/seat"
	"github.com/brutusmcforce/opennox/libs/client/seat/opengl"
	"github.com/brutusmcforce/opennox/libs/env"
	noxlog "github.com/brutusmcforce/opennox/libs/log"
	"github.com/brutusmcforce/opennox/libs/sdl"
)

var (
	debugGpad = os.Getenv("NOX_DEBUG_GPAD") == "true"
)

var _ seat.Seat = &Window{}

// ErrOpenGL is wrapped by the errors of New when the window was created, but not an OpenGL 3.3 core context for it
// (e.g. the GPU or its driver only supports an older version).
var ErrOpenGL = errors.New("OpenGL 3.3 is not available")

// New creates a new SDL window which implements a Seat interface.
func New(log *slog.Logger, title string, sz image.Point) (*Window, error) {
	log = noxlog.WithSystem(log, "sdl")
	// SDL3 is per-monitor DPI aware on Windows by default.
	// TODO: if we ever decide to use multiple windows, this will need to be moved elsewhere; same for sdl.Quit
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, fmt.Errorf("SDL Initialization failed: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL version: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL version: %w", err)
	}
	// the attribute ids are GL_CONTEXT_PROFILE_MASK and GL_CONTEXT_FLAGS (opennox#749: GL_CONTEXT_PROFILE_CORE
	// and GL_CONTEXT_FORWARD_COMPATIBLE_FLAG are values, used as ids they set the green and blue sizes);
	// macOS (and Wine on it) only has 3.2+ contexts as forward-compatible core profiles
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL core: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL forward: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL attribute: %w", err)
	}

	flags := sdl.WINDOW_RESIZABLE | sdl.WINDOW_OPENGL | sdl.WINDOW_HIGH_PIXEL_DENSITY
	if env.IsE2E() {
		// E2E tests compare screenshots, so the size in pixels must not depend on the display scale
		// and the window must not be resized by the window manager
		flags &^= sdl.WINDOW_HIGH_PIXEL_DENSITY | sdl.WINDOW_RESIZABLE
	}
	win, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(sz.X), int32(sz.Y), flags)
	if err != nil {
		sdl.Quit()
		return nil, fmt.Errorf("SDL Window creation failed: %w", err)
	}
	h := &Window{
		log: log,
		win: win, prevSz: sz,
	}
	h.SetScreenMode(seat.Windowed)
	if err := h.initGL(); err != nil {
		_ = win.Destroy()
		sdl.Quit()
		return nil, fmt.Errorf("%w: %w", ErrOpenGL, err)
	}
	// game controllers: SDL reports the connected ones as added devices; E2E tests inject their own events
	if !env.IsE2E() {
		if err := sdl.InitSubSystem(sdl.INIT_GAMEPAD); err != nil {
			log.Warn("game controllers are not available", "err", err)
		}
	}
	return h, nil
}

type Window struct {
	log      *slog.Logger
	win      *sdl.Window
	gl       opengl.Window
	prevPos  image.Point
	prevSz   image.Point
	textInp  bool
	mode     seat.ScreenMode
	cursor   *cursor // system cursor set by the game, if any
	mpos     image.Point
	onResize []func(sz image.Point)
	onInput  []func(ev seat.InputEvent)
	pads     map[sdl.JoystickID]*sdl.GameController
	wheel    float32 // fraction of a wheel step not reported yet (touchpads scroll in fractions)
}

func (win *Window) Close() error {
	if win.win == nil {
		return nil
	}
	win.closeGamepads()
	win.gl.Close()
	err := win.win.Destroy()
	win.win = nil
	win.onResize = nil
	win.onInput = nil
	sdl.Quit()
	return err
}

func (win *Window) NewSurface(sz image.Point, filter bool) seat.Surface {
	return win.gl.NewSurface(sz, filter)
}

func (win *Window) Clear() {
	win.gl.Clear()
}

func (win *Window) ScreenSize() image.Point {
	w, h := win.win.GLGetDrawableSize()
	return image.Point{
		X: int(w), Y: int(h),
	}
}

func (win *Window) screenPos() image.Point {
	x, y := win.win.GetPosition()
	return image.Point{
		X: int(x), Y: int(y),
	}
}

func (win *Window) displayRect() sdl.Rect {
	disp, err := win.win.GetDisplayIndex()
	if err != nil {
		win.log.Warn("can't get display index", "err", err)
		return sdl.Rect{}
	}
	rect, err := sdl.GetDisplayBounds(disp)
	if err != nil {
		win.log.Warn("can't get display bounds", "err", err)
		return sdl.Rect{}
	}
	return rect
}

func (win *Window) ScreenMaxSize() image.Point {
	rect := win.displayRect()
	return image.Point{
		X: int(rect.W), Y: int(rect.H),
	}
}

func (win *Window) setSize(sz image.Point) {
	win.log.Info("window", "size", sz)
	win.win.SetSize(int32(sz.X), int32(sz.Y))
}

func (win *Window) center() {
	win.win.SetPosition(sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED)
}

func (win *Window) ResizeScreen(sz image.Point) {
	if win.mode != seat.Windowed {
		return
	}
	win.setSize(sz)
	win.prevSz = sz
}

func (win *Window) SetScreenMode(mode seat.ScreenMode) {
	if env.IsE2E() {
		// E2E tests need an exact window size, see New
		mode = seat.Windowed
	}
	if win.mode == mode {
		return
	}
	if win.mode == seat.Windowed {
		// preserve size and pos, so we can restore them later
		win.prevSz = win.ScreenSize()
		win.prevPos = win.screenPos()
	}
	switch mode {
	case seat.Windowed:
		win.win.SetFullscreen(0)
		// E2E tests need an exact window size; window managers don't resize fixed-size windows
		win.win.SetResizable(!env.IsE2E())
		win.win.SetBordered(true)
		win.setSize(win.prevSz)
		if win.prevPos != (image.Point{}) {
			win.win.SetPosition(int32(win.prevPos.X), int32(win.prevPos.Y))
		} else {
			win.center()
		}
		win.win.SetGrab(false)
	case seat.Fullscreen:
		win.win.SetResizable(false)
		win.win.SetBordered(false)
		win.setSize(win.ScreenMaxSize())
		if err := win.win.SetFullscreen(1); err != nil { // desktop fullscreen
			win.log.Warn("cannot switch to fullscreen", "err", err)
		}
		// keep the cursor in the window, e.g. when there are multiple monitors
		win.win.SetGrab(true)
	case seat.Borderless:
		win.win.SetFullscreen(0)
		win.win.SetResizable(false)
		win.win.SetBordered(true)
		win.setSize(win.ScreenMaxSize())
		win.center()
		win.win.SetGrab(false)
	}
	win.mode = mode
	win.updateCursor()
}

// SetGamma sets screen gamma parameter.
func (win *Window) SetGamma(v float32) {
	win.gl.SetGamma(v)
}

func (win *Window) OnScreenResize(fnc func(sz image.Point)) {
	win.onResize = append(win.onResize, fnc)
}

func (win *Window) initGL() error {
	gtx, err := win.win.GLCreateContext()
	if err != nil {
		return fmt.Errorf("OpenGL creation failed: %w", err)
	}
	err = win.win.GLMakeCurrent(gtx)
	if err != nil {
		return fmt.Errorf("OpenGL bind failed: %w", err)
	}
	sdl.GLSetSwapInterval(0)
	if err := win.gl.Init(win.log); err != nil {
		return err
	}
	return nil
}

func (win *Window) Present() {
	win.win.GLSwap()
}
