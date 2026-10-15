package sdl

import (
	"errors"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/dynlib"
)

// Window is SDL_Window.
type Window struct {
	p uintptr
}

// CreateWindow creates a window. SDL3 creates windows without a position, so it's set afterwards,
// unless it's WINDOWPOS_UNDEFINED.
func CreateWindow(title string, x, y, w, h int32, flags WindowFlags) (*Window, error) {
	if err := Load(); err != nil {
		return nil, err
	}
	t := dynlib.CString(title)
	p := call(fn.CreateWindow, uintptr(unsafe.Pointer(t)), uintptr(w), uintptr(h), uintptr(flags))
	if p == 0 {
		return nil, GetError()
	}
	win := &Window{p: p}
	if x != WINDOWPOS_UNDEFINED || y != WINDOWPOS_UNDEFINED {
		win.SetPosition(x, y)
	}
	return win, nil
}

func (w *Window) Destroy() error {
	if w == nil || w.p == 0 {
		return errors.New("invalid window")
	}
	call(fn.DestroyWindow, w.p)
	w.p = 0
	return nil
}

func (w *Window) GLCreateContext() (GLContext, error) {
	p := call(fn.GLCreateContext, w.p)
	if p == 0 {
		return nil, GetError()
	}
	return GLContext(unsafe.Pointer(p)), nil //nolint:govet // native handle
}

func (w *Window) GLMakeCurrent(ctx GLContext) error {
	return errorFrom(call(fn.GLMakeCurrent, w.p, uintptr(ctx)))
}

func (w *Window) GLSwap() {
	call(fn.GLSwapWindow, w.p)
}

// GLGetDrawableSize returns the size of the window in pixels (it differs from GetSize on high density displays).
func (w *Window) GLGetDrawableSize() (width, height int32) {
	call(fn.GetWindowSizeInPixels, w.p, uintptr(unsafe.Pointer(&width)), uintptr(unsafe.Pointer(&height)))
	return
}

func (w *Window) GetSize() (width, height int32) {
	call(fn.GetWindowSize, w.p, uintptr(unsafe.Pointer(&width)), uintptr(unsafe.Pointer(&height)))
	return
}

func (w *Window) SetSize(width, height int32) {
	call(fn.SetWindowSize, w.p, uintptr(width), uintptr(height))
}

func (w *Window) GetPosition() (x, y int32) {
	call(fn.GetWindowPosition, w.p, uintptr(unsafe.Pointer(&x)), uintptr(unsafe.Pointer(&y)))
	return
}

func (w *Window) SetPosition(x, y int32) {
	call(fn.SetWindowPosition, w.p, uintptr(x), uintptr(y))
}

// GetDisplayIndex returns the ID of the display the window is on (see GetDisplayBounds).
func (w *Window) GetDisplayIndex() (int, error) {
	id := uint32(call(fn.GetDisplayForWindow, w.p))
	if id == 0 {
		return 0, GetError()
	}
	return int(id), nil
}

// SetFullscreen switches the window to desktop fullscreen (any flags) or back to a window (0).
func (w *Window) SetFullscreen(flags uint32) error {
	return errorFrom(call(fn.SetWindowFullscreen, w.p, cbool(flags != 0)))
}

func (w *Window) SetResizable(resizable bool) {
	call(fn.SetWindowResizable, w.p, cbool(resizable))
}

func (w *Window) SetBordered(bordered bool) {
	call(fn.SetWindowBordered, w.p, cbool(bordered))
}

// SetGrab confines the mouse to the window.
func (w *Window) SetGrab(grabbed bool) {
	call(fn.SetWindowMouseGrab, w.p, cbool(grabbed))
}

// StartTextInput enables text input events (and the on-screen keyboard or IME, if any) for the window.
func (w *Window) StartTextInput() {
	call(fn.StartTextInput, w.p)
}

func (w *Window) StopTextInput() {
	call(fn.StopTextInput, w.p)
}

// Surface is SDL_Surface.
type Surface struct {
	p     uintptr
	W, H  int32
	Pitch int32
}

// SDL_Surface fields on 64-bit platforms
const (
	surfaceW      = 8
	surfaceH      = 12
	surfacePitch  = 16
	surfacePixels = 24
)

func CreateSurface(width, height int32, format uint32) (*Surface, error) {
	p := call(fn.CreateSurface, uintptr(width), uintptr(height), uintptr(format))
	if p == 0 {
		return nil, GetError()
	}
	s := &Surface{p: p}
	s.W = *(*int32)(unsafe.Pointer(p + surfaceW))         //nolint:govet // native memory
	s.H = *(*int32)(unsafe.Pointer(p + surfaceH))         //nolint:govet // native memory
	s.Pitch = *(*int32)(unsafe.Pointer(p + surfacePitch)) //nolint:govet // native memory
	return s, nil
}

// Pixels returns the pixel memory of the surface (Pitch*H bytes).
func (s *Surface) Pixels() []byte {
	pix := *(*unsafe.Pointer)(unsafe.Pointer(s.p + surfacePixels)) //nolint:govet // native memory
	return unsafe.Slice((*byte)(pix), int(s.Pitch)*int(s.H))
}

func (s *Surface) Free() {
	if s.p != 0 {
		call(fn.DestroySurface, s.p)
		s.p = 0
	}
}

// Cursor is SDL_Cursor.
type Cursor struct {
	p uintptr
}

func CreateColorCursor(s *Surface, hotX, hotY int32) *Cursor {
	p := call(fn.CreateColorCursor, s.p, uintptr(hotX), uintptr(hotY))
	if p == 0 {
		return nil
	}
	return &Cursor{p: p}
}

func SetCursor(c *Cursor) {
	var p uintptr
	if c != nil {
		p = c.p
	}
	call(fn.SetCursor, p)
}

func GetDefaultCursor() *Cursor {
	p := call(fn.GetDefaultCursor)
	if p == 0 {
		return nil
	}
	return &Cursor{p: p}
}

func FreeCursor(c *Cursor) {
	if c != nil && c.p != 0 {
		call(fn.DestroyCursor, c.p)
		c.p = 0
	}
}

// ShowCursor shows or hides the mouse cursor.
func ShowCursor(show bool) error {
	if show {
		return errorFrom(call(fn.ShowCursor))
	}
	return errorFrom(call(fn.HideCursor))
}

// GameController is SDL_Gamepad.
type GameController struct {
	p uintptr
}

// Joystick is the joystick of a gamepad; only its instance ID is used.
type Joystick struct {
	id JoystickID
}

// IsGameController checks if a joystick (by instance ID, as in the device events) is a gamepad.
func IsGameController(id JoystickID) bool {
	return call(fn.IsGamepad, uintptr(id))&0xff != 0
}

// GameControllerOpen opens a gamepad by its joystick instance ID.
func GameControllerOpen(id JoystickID) *GameController {
	p := call(fn.OpenGamepad, uintptr(id))
	if p == 0 {
		return nil
	}
	return &GameController{p: p}
}

func (c *GameController) Close() {
	if c.p != 0 {
		call(fn.CloseGamepad, c.p)
		c.p = 0
	}
}

func (c *GameController) Name() string {
	return dynlib.GoString(call(fn.GetGamepadName, c.p))
}

func (c *GameController) Joystick() *Joystick {
	return &Joystick{id: JoystickID(call(fn.GetGamepadID, c.p))}
}

func (j *Joystick) InstanceID() JoystickID {
	return j.id
}
