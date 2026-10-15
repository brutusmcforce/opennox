// Package sdl is a minimal SDL3 binding without cgo: the library is loaded at runtime (see dynlib).
//
// It covers what OpenNox uses. The API follows the go-sdl2 names it replaced, adapted to SDL3 where SDL3 differs.
// The functions must be called from the main thread, like SDL itself requires.
package sdl

import (
	"errors"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"

	"github.com/brutusmcforce/opennox/libs/dynlib"
)

func init() {
	// SDL and OpenGL must be used from the main thread: keep the main goroutine on it (go-sdl2 did the same).
	// Package initialization runs on the main goroutine.
	runtime.LockOSThread()
}

func libNames() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"SDL3.dll"}
	case "darwin":
		// Homebrew isn't on the default search path (Apple Silicon: /opt/homebrew, Intel: /usr/local)
		return []string{"libSDL3.0.dylib", "libSDL3.dylib", "SDL3.framework/SDL3",
			"/opt/homebrew/lib/libSDL3.0.dylib", "/usr/local/lib/libSDL3.0.dylib"}
	}
	return []string{"libSDL3.so.0", "libSDL3.so"}
}

var fn struct {
	SetHint               uintptr
	Init                  uintptr
	InitSubSystem         uintptr
	Quit                  uintptr
	GetError              uintptr
	Free                  uintptr
	GLSetAttribute        uintptr
	GLSetSwapInterval     uintptr
	GLGetProcAddress      uintptr
	GLCreateContext       uintptr
	GLMakeCurrent         uintptr
	GLSwapWindow          uintptr
	GetWindowSizeInPixels uintptr
	CreateWindow          uintptr
	DestroyWindow         uintptr
	GetWindowSize         uintptr
	SetWindowSize         uintptr
	GetWindowPosition     uintptr
	SetWindowPosition     uintptr
	GetDisplayForWindow   uintptr
	SetWindowFullscreen   uintptr
	SetWindowResizable    uintptr
	SetWindowBordered     uintptr
	SetWindowMouseGrab    uintptr
	GetDisplayBounds      uintptr
	PollEvent             uintptr
	StartTextInput        uintptr
	StopTextInput         uintptr
	GetModState           uintptr
	GetClipboardText      uintptr
	IsGamepad             uintptr
	OpenGamepad           uintptr
	CloseGamepad          uintptr
	GetGamepadName        uintptr
	GetGamepadID          uintptr
	CreateSurface         uintptr
	DestroySurface        uintptr
	CreateColorCursor     uintptr
	SetCursor             uintptr
	GetDefaultCursor      uintptr
	DestroyCursor         uintptr
	ShowCursor            uintptr
	HideCursor            uintptr
}

var (
	loadOnce sync.Once
	loadErr  error
)

// Load loads the SDL3 library. Other functions call it as needed.
func Load() error {
	loadOnce.Do(func() {
		lib, err := dynlib.Open(libNames()...)
		if err != nil {
			loadErr = err
			return
		}
		loadErr = lib.Syms(map[string]*uintptr{
			"SDL_SetHint":               &fn.SetHint,
			"SDL_Init":                  &fn.Init,
			"SDL_InitSubSystem":         &fn.InitSubSystem,
			"SDL_Quit":                  &fn.Quit,
			"SDL_GetError":              &fn.GetError,
			"SDL_free":                  &fn.Free,
			"SDL_GL_SetAttribute":       &fn.GLSetAttribute,
			"SDL_GL_SetSwapInterval":    &fn.GLSetSwapInterval,
			"SDL_GL_GetProcAddress":     &fn.GLGetProcAddress,
			"SDL_GL_CreateContext":      &fn.GLCreateContext,
			"SDL_GL_MakeCurrent":        &fn.GLMakeCurrent,
			"SDL_GL_SwapWindow":         &fn.GLSwapWindow,
			"SDL_GetWindowSizeInPixels": &fn.GetWindowSizeInPixels,
			"SDL_CreateWindow":          &fn.CreateWindow,
			"SDL_DestroyWindow":         &fn.DestroyWindow,
			"SDL_GetWindowSize":         &fn.GetWindowSize,
			"SDL_SetWindowSize":         &fn.SetWindowSize,
			"SDL_GetWindowPosition":     &fn.GetWindowPosition,
			"SDL_SetWindowPosition":     &fn.SetWindowPosition,
			"SDL_GetDisplayForWindow":   &fn.GetDisplayForWindow,
			"SDL_SetWindowFullscreen":   &fn.SetWindowFullscreen,
			"SDL_SetWindowResizable":    &fn.SetWindowResizable,
			"SDL_SetWindowBordered":     &fn.SetWindowBordered,
			"SDL_SetWindowMouseGrab":    &fn.SetWindowMouseGrab,
			"SDL_GetDisplayBounds":      &fn.GetDisplayBounds,
			"SDL_PollEvent":             &fn.PollEvent,
			"SDL_StartTextInput":        &fn.StartTextInput,
			"SDL_StopTextInput":         &fn.StopTextInput,
			"SDL_GetModState":           &fn.GetModState,
			"SDL_GetClipboardText":      &fn.GetClipboardText,
			"SDL_IsGamepad":             &fn.IsGamepad,
			"SDL_OpenGamepad":           &fn.OpenGamepad,
			"SDL_CloseGamepad":          &fn.CloseGamepad,
			"SDL_GetGamepadName":        &fn.GetGamepadName,
			"SDL_GetGamepadID":          &fn.GetGamepadID,
			"SDL_CreateSurface":         &fn.CreateSurface,
			"SDL_DestroySurface":        &fn.DestroySurface,
			"SDL_CreateColorCursor":     &fn.CreateColorCursor,
			"SDL_SetCursor":             &fn.SetCursor,
			"SDL_GetDefaultCursor":      &fn.GetDefaultCursor,
			"SDL_DestroyCursor":         &fn.DestroyCursor,
			"SDL_ShowCursor":            &fn.ShowCursor,
			"SDL_HideCursor":            &fn.HideCursor,
		})
	})
	return loadErr
}

func call(f uintptr, args ...uintptr) uintptr {
	r, _, _ := purego.SyscallN(f, args...)
	return r
}

func cbool(v bool) uintptr {
	if v {
		return 1
	}
	return 0
}

// errorFrom returns SDL's error if a call returned false (SDL3 functions return a C bool).
func errorFrom(r uintptr) error {
	if r&0xff == 0 {
		return GetError()
	}
	return nil
}

// GetError returns the last SDL error, or nil if there's none.
func GetError() error {
	if Load() != nil {
		return loadErr
	}
	if s := dynlib.GoString(call(fn.GetError)); s != "" {
		return errors.New(s)
	}
	return nil
}

// SetHint sets an SDL hint. It does nothing if SDL can't be loaded (Init reports that).
func SetHint(name, value string) bool {
	if Load() != nil {
		return false
	}
	n, v := dynlib.CString(name), dynlib.CString(value)
	return call(fn.SetHint, uintptr(unsafe.Pointer(n)), uintptr(unsafe.Pointer(v)))&0xff != 0
}

func Init(flags uint32) error {
	if err := Load(); err != nil {
		return err
	}
	return errorFrom(call(fn.Init, uintptr(flags)))
}

func InitSubSystem(flags uint32) error {
	if err := Load(); err != nil {
		return err
	}
	return errorFrom(call(fn.InitSubSystem, uintptr(flags)))
}

func Quit() {
	if Load() != nil {
		return
	}
	call(fn.Quit)
}

func GLSetAttribute(attr GLattr, value int) error {
	if err := Load(); err != nil {
		return err
	}
	return errorFrom(call(fn.GLSetAttribute, uintptr(attr), uintptr(value)))
}

func GLSetSwapInterval(interval int) error {
	if err := Load(); err != nil {
		return err
	}
	return errorFrom(call(fn.GLSetSwapInterval, uintptr(interval)))
}

// GLGetProcAddress returns the address of an OpenGL function of the current context, or nil.
func GLGetProcAddress(name string) unsafe.Pointer {
	if Load() != nil {
		return nil
	}
	n := dynlib.CString(name)
	return unsafe.Pointer(call(fn.GLGetProcAddress, uintptr(unsafe.Pointer(n)))) //nolint:govet // native function address
}

func GetModState() Keymod { return Keymod(call(fn.GetModState)) }

func GetClipboardText() (string, error) {
	p := call(fn.GetClipboardText)
	if p == 0 {
		return "", GetError()
	}
	s := dynlib.GoString(p)
	call(fn.Free, p)
	return s, nil
}

// GetDisplayBounds returns the bounds of a display, given its ID (see Window.GetDisplayIndex).
func GetDisplayBounds(display int) (Rect, error) {
	var r Rect
	err := errorFrom(call(fn.GetDisplayBounds, uintptr(display), uintptr(unsafe.Pointer(&r))))
	return r, err
}
