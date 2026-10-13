package seat

import (
	"image"

	"github.com/brutusmcforce/opennox/libs/client/keybind"
	"github.com/brutusmcforce/opennox/libs/types"
)

type InputConfig []func(ev InputEvent)

type Input interface {
	// InputTick runs input event processing, triggering all relevant callback functions.
	InputTick()
	// ReplaceInputs replaces all inputs by a provided configuration, and returns the previous configuration.
	ReplaceInputs(InputConfig) InputConfig
	// OnInput adds a handler function that will be called on each input event.
	OnInput(fnc func(ev InputEvent))
	// SetTextInput enables or disables text input mode.
	SetTextInput(enable bool)
}

type InputEvent interface {
	isInputEvent()
}

type MouseButton int

const (
	MouseButtonLeft = MouseButton(iota)
	MouseButtonRight
	MouseButtonMiddle
)

var _ InputEvent = &MouseMoveEvent{}

type MouseMoveEvent struct {
	Relative bool
	Pos      image.Point
	Rel      types.Pointf
}

func (e *MouseMoveEvent) isInputEvent() {}

var _ InputEvent = &MouseButtonEvent{}

type MouseButtonEvent struct {
	Button  MouseButton
	Pressed bool
}

func (e *MouseButtonEvent) isInputEvent() {}

var _ InputEvent = &MouseWheelEvent{}

type MouseWheelEvent struct {
	Wheel int
}

func (e *MouseWheelEvent) isInputEvent() {}

const (
	WindowClosed = WindowEvent(iota)
	WindowFocused
	WindowUnfocused
)

var _ InputEvent = WindowEvent(0)

type WindowEvent int

func (e WindowEvent) isInputEvent() {}

var _ InputEvent = &TextEditEvent{}

type TextEditEvent struct {
	Text string
}

func (e *TextEditEvent) isInputEvent() {}

var _ InputEvent = &TextInputEvent{}

type TextInputEvent struct {
	Text string
}

func (e *TextInputEvent) isInputEvent() {}

var _ InputEvent = &KeyboardEvent{}

type KeyboardEvent struct {
	Key     keybind.Key
	Pressed bool
}

func (e *KeyboardEvent) isInputEvent() {}

// GamepadButton is a game controller button, in the layout of an Xbox controller (A is the bottom face button).
type GamepadButton int

const (
	GamepadA = GamepadButton(iota)
	GamepadB
	GamepadX
	GamepadY
	GamepadBack
	GamepadGuide
	GamepadStart
	GamepadLeftStick
	GamepadRightStick
	GamepadLeftShoulder
	GamepadRightShoulder
	GamepadDpadUp
	GamepadDpadDown
	GamepadDpadLeft
	GamepadDpadRight
)

// GamepadAxis is a game controller stick axis or trigger.
type GamepadAxis int

const (
	GamepadLeftX = GamepadAxis(iota)
	GamepadLeftY
	GamepadRightX
	GamepadRightY
	GamepadTriggerLeft
	GamepadTriggerRight
)

var _ InputEvent = &GamepadButtonEvent{}

type GamepadButtonEvent struct {
	Button  GamepadButton
	Pressed bool
}

func (e *GamepadButtonEvent) isInputEvent() {}

var _ InputEvent = &GamepadAxisEvent{}

// GamepadAxisEvent is a new position of a stick axis (-1 to 1; down and right are positive) or a trigger (0 to 1).
type GamepadAxisEvent struct {
	Axis  GamepadAxis
	Value float32
}

func (e *GamepadAxisEvent) isInputEvent() {}

var _ InputEvent = &GamepadDeviceEvent{}

// GamepadDeviceEvent is sent when a game controller is connected or disconnected.
type GamepadDeviceEvent struct {
	Name      string
	Connected bool
}

func (e *GamepadDeviceEvent) isInputEvent() {}
