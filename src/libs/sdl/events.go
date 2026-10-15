package sdl

import (
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/dynlib"
)

// Event types (SDL_EventType).
const (
	QUIT                     = 0x100
	WINDOWEVENT_FOCUS_GAINED = 0x20E
	WINDOWEVENT_FOCUS_LOST   = 0x20F
	KEYDOWN                  = 0x300
	KEYUP                    = 0x301
	TEXTEDITING              = 0x302
	TEXTINPUT                = 0x303
	MOUSEMOTION              = 0x400
	MOUSEBUTTONDOWN          = 0x401
	MOUSEBUTTONUP            = 0x402
	MOUSEWHEEL               = 0x403
	CONTROLLERAXISMOTION     = 0x650
	CONTROLLERBUTTONDOWN     = 0x651
	CONTROLLERBUTTONUP       = 0x652
	CONTROLLERDEVICEADDED    = 0x653
	CONTROLLERDEVICEREMOVED  = 0x654
	CONTROLLERDEVICEREMAPPED = 0x655
)

// Event is one of the event types below; other SDL events are not decoded.
type Event interface {
	GetType() uint32
}

type CommonEvent struct {
	Type      uint32
	Timestamp uint64 // in nanoseconds
}

func (e CommonEvent) GetType() uint32 { return e.Type }

type QuitEvent struct {
	CommonEvent
}

// WindowEvent is a window focus event: SDL3 has a separate event type for each window event,
// Event is the same as Type (WINDOWEVENT_FOCUS_GAINED or WINDOWEVENT_FOCUS_LOST).
type WindowEvent struct {
	CommonEvent
	WindowID uint32
	Event    uint32
}

type Keysym struct {
	Scancode Scancode
	Sym      uint32
	Mod      Keymod
}

type KeyboardEvent struct {
	CommonEvent
	WindowID uint32
	State    uint8 // PRESSED or RELEASED
	Repeat   bool
	Keysym   Keysym
}

type TextEditingEvent struct {
	CommonEvent
	WindowID uint32
	Text     string
	Start    int32
	Length   int32
}

func (e *TextEditingEvent) GetText() string { return e.Text }

type TextInputEvent struct {
	CommonEvent
	WindowID uint32
	Text     string
}

func (e *TextInputEvent) GetText() string { return e.Text }

// MouseMotionEvent has positions in window coordinates (see Window.GetSize).
type MouseMotionEvent struct {
	CommonEvent
	WindowID   uint32
	Which      uint32
	State      uint32
	X, Y       float32
	XRel, YRel float32
}

type MouseButtonEvent struct {
	CommonEvent
	WindowID uint32
	Which    uint32
	Button   uint8
	State    uint8 // PRESSED or RELEASED
	Clicks   uint8
	X, Y     float32
}

// MouseWheelEvent has fractional amounts for precise devices (touchpads); a mouse wheel step is 1.
type MouseWheelEvent struct {
	CommonEvent
	WindowID  uint32
	Which     uint32
	X, Y      float32
	Direction uint32
}

type ControllerAxisEvent struct {
	CommonEvent
	Which JoystickID
	Axis  uint8
	Value int16
}

type ControllerButtonEvent struct {
	CommonEvent
	Which  JoystickID
	Button uint8
	State  uint8 // PRESSED or RELEASED
}

type ControllerDeviceEvent struct {
	CommonEvent
	Which JoystickID
}

// eventBuf receives events from SDL_PollEvent (SDL_Event is a 128 byte union). SDL events are only polled from the
// main thread, so one buffer is enough, and it doesn't allocate on every call.
var eventBuf [128]byte

// PollEvent returns the next pending event, or nil if there's none. Events of other types are skipped.
func PollEvent() Event {
	for {
		if call(fn.PollEvent, uintptr(unsafe.Pointer(&eventBuf)))&0xff == 0 {
			return nil
		}
		if ev := decodeEvent(eventBuf[:]); ev != nil {
			return ev
		}
	}
}

func state(down byte) uint8 {
	if down != 0 {
		return PRESSED
	}
	return RELEASED
}

// decodeEvent decodes an SDL_Event (offsets of the SDL3 structures on 64-bit platforms).
// Text is copied: SDL only keeps it until the next poll.
func decodeEvent(b []byte) Event {
	le := binary.LittleEndian
	u32 := func(off int) uint32 { return le.Uint32(b[off:]) }
	i32 := func(off int) int32 { return int32(le.Uint32(b[off:])) }
	f32 := func(off int) float32 { return math.Float32frombits(le.Uint32(b[off:])) }
	text := func(off int) string { return dynlib.GoString(uintptr(le.Uint64(b[off:]))) }
	com := CommonEvent{Type: u32(0), Timestamp: le.Uint64(b[8:])}
	switch com.Type {
	case QUIT:
		return QuitEvent{CommonEvent: com}
	case WINDOWEVENT_FOCUS_GAINED, WINDOWEVENT_FOCUS_LOST:
		return WindowEvent{CommonEvent: com, WindowID: u32(16), Event: com.Type}
	case KEYDOWN, KEYUP:
		return KeyboardEvent{CommonEvent: com, WindowID: u32(16), State: state(b[36]), Repeat: b[37] != 0,
			Keysym: Keysym{Scancode: Scancode(u32(24)), Sym: u32(28), Mod: Keymod(le.Uint16(b[32:]))}}
	case TEXTEDITING:
		return TextEditingEvent{CommonEvent: com, WindowID: u32(16), Text: text(24), Start: i32(32), Length: i32(36)}
	case TEXTINPUT:
		return TextInputEvent{CommonEvent: com, WindowID: u32(16), Text: text(24)}
	case MOUSEMOTION:
		return MouseMotionEvent{CommonEvent: com, WindowID: u32(16), Which: u32(20), State: u32(24),
			X: f32(28), Y: f32(32), XRel: f32(36), YRel: f32(40)}
	case MOUSEBUTTONDOWN, MOUSEBUTTONUP:
		return MouseButtonEvent{CommonEvent: com, WindowID: u32(16), Which: u32(20), Button: b[24], State: state(b[25]),
			Clicks: b[26], X: f32(28), Y: f32(32)}
	case MOUSEWHEEL:
		return MouseWheelEvent{CommonEvent: com, WindowID: u32(16), Which: u32(20), X: f32(24), Y: f32(28),
			Direction: u32(32)}
	case CONTROLLERAXISMOTION:
		return ControllerAxisEvent{CommonEvent: com, Which: JoystickID(u32(16)), Axis: b[20], Value: int16(le.Uint16(b[24:]))}
	case CONTROLLERBUTTONDOWN, CONTROLLERBUTTONUP:
		return ControllerButtonEvent{CommonEvent: com, Which: JoystickID(u32(16)), Button: b[20], State: state(b[21])}
	case CONTROLLERDEVICEADDED, CONTROLLERDEVICEREMOVED, CONTROLLERDEVICEREMAPPED:
		return ControllerDeviceEvent{CommonEvent: com, Which: JoystickID(u32(16))}
	}
	return nil
}
