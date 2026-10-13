package sdl

import (
	"image"
	"strings"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/brutusmcforce/opennox/libs/client/seat"
	"github.com/brutusmcforce/opennox/libs/types"
)

func (win *Window) ReplaceInputs(cfg seat.InputConfig) seat.InputConfig {
	oldCfg := win.onInput
	win.onInput = cfg
	return oldCfg
}

func (win *Window) OnInput(fnc func(ev seat.InputEvent)) {
	win.onInput = append(win.onInput, fnc)
}

func (win *Window) SetTextInput(enable bool) {
	if win.textInp == enable {
		return
	}
	win.textInp = enable
	if enable {
		sdl.StartTextInput()
	} else {
		sdl.StopTextInput()
	}
}

func (win *Window) InputTick() {
	for {
		switch ev := sdl.PollEvent().(type) {
		case nil:
			// no more events
			return
		case sdl.TextEditingEvent:
			win.processTextEditingEvent(&ev)
		case sdl.TextInputEvent:
			win.processTextInputEvent(&ev)
		case sdl.KeyboardEvent:
			win.processKeyboardEvent(&ev)
		case sdl.MouseButtonEvent:
			win.processMouseButtonEvent(&ev)
		case sdl.MouseMotionEvent:
			win.processMotionEvent(&ev)
		case sdl.MouseWheelEvent:
			win.processWheelEvent(&ev)
		case sdl.ControllerAxisEvent:
			if debugGpad {
				win.log.Debug("SDL_CONTROLLERAXISMOTION",
					"ev", ev.GetType(), "dev", ev.Which, "axis", ev.Axis, "val", ev.Value)
			}
			win.processGamepadAxisEvent(&ev)
		case sdl.ControllerButtonEvent:
			if debugGpad {
				win.log.Debug("SDL_CONTROLLERBUTTON",
					"ev", ev.GetType(), "dev", ev.Which, "btn", ev.Button, "state", ev.State)
			}
			win.processGamepadButtonEvent(&ev)
		case sdl.ControllerDeviceEvent:
			switch ev.GetType() {
			case sdl.CONTROLLERDEVICEADDED:
				if debugGpad {
					win.log.Debug("SDL_CONTROLLERDEVICEADDED", "ev", ev.GetType(), "dev", ev.Which)
				}
				win.processGamepadDeviceEvent(&ev)
			case sdl.CONTROLLERDEVICEREMOVED:
				if debugGpad {
					win.log.Debug("SDL_CONTROLLERDEVICEREMOVED", "ev", ev.GetType(), "dev", ev.Which)
				}
				win.processGamepadDeviceEvent(&ev)
			case sdl.CONTROLLERDEVICEREMAPPED:
				if debugGpad {
					win.log.Debug("SDL_CONTROLLERDEVICEREMAPPED", "ev", ev.GetType())
				}
			}
		case sdl.WindowEvent:
			win.processWindowEvent(&ev)
		case sdl.QuitEvent:
			win.processQuitEvent(&ev)
		}
		// TODO: touch events for WASM
	}
}

func (win *Window) inputEvent(ev seat.InputEvent) {
	for _, fnc := range win.onInput {
		fnc(ev)
	}
}

func (win *Window) processQuitEvent(ev *sdl.QuitEvent) {
	win.inputEvent(seat.WindowClosed)
}

func (win *Window) processWindowEvent(ev *sdl.WindowEvent) {
	switch ev.Event {
	case sdl.WINDOWEVENT_FOCUS_LOST:
		win.inputEvent(seat.WindowUnfocused)
	case sdl.WINDOWEVENT_FOCUS_GAINED:
		win.inputEvent(seat.WindowFocused)
	}
}

func (win *Window) processTextEditingEvent(ev *sdl.TextEditingEvent) {
	win.inputEvent(&seat.TextEditEvent{
		Text: ev.GetText(),
	})
}

func (win *Window) processTextInputEvent(ev *sdl.TextInputEvent) {
	text := ev.GetText()
	if sdl.GetModState()&sdl.KMOD_CTRL != 0 && len(text) == 1 && strings.ToLower(text) == "v" {
		return // ignore "V" from Ctrl-V
	}
	win.inputEvent(&seat.TextInputEvent{
		Text: text,
	})
}

func (win *Window) processKeyboardEvent(ev *sdl.KeyboardEvent) {
	if win.textInp && ev.State == sdl.PRESSED && sdl.GetModState()&sdl.KMOD_CTRL != 0 && ev.Keysym.Scancode == sdl.SCANCODE_V {
		text, err := sdl.GetClipboardText()
		if err != nil {
			win.log.Error("cannot get clipboard text", "err", err)
			return
		}
		win.inputEvent(&seat.TextInputEvent{
			Text: text,
		})
		return
	}
	key := scanCodeToKeyNum[ev.Keysym.Scancode]
	win.inputEvent(&seat.KeyboardEvent{
		Key:     key,
		Pressed: ev.State == sdl.PRESSED,
	})
}

func (win *Window) processMouseButtonEvent(ev *sdl.MouseButtonEvent) {
	pressed := ev.State == sdl.PRESSED
	// TODO: handle focus, or move to other place
	//if pressed {
	//	h.iface.WindowEvent(WindowFocus)
	//}

	var button seat.MouseButton
	switch ev.Button {
	case sdl.BUTTON_LEFT:
		button = seat.MouseButtonLeft
	case sdl.BUTTON_RIGHT:
		button = seat.MouseButtonRight
	case sdl.BUTTON_MIDDLE:
		button = seat.MouseButtonMiddle
	default:
		return
	}
	win.inputEvent(&seat.MouseButtonEvent{
		Button:  button,
		Pressed: pressed,
	})
}

func (win *Window) processMotionEvent(ev *sdl.MouseMotionEvent) {
	pos := image.Point{X: int(ev.X), Y: int(ev.Y)}
	// mouse events are in window coordinates, which differ from pixels on high DPI screens
	if ww, wh := win.win.GetSize(); ww > 0 && wh > 0 {
		sz := win.ScreenSize()
		pos.X = pos.X * sz.X / int(ww)
		pos.Y = pos.Y * sz.Y / int(wh)
	}
	win.inputEvent(&seat.MouseMoveEvent{
		Pos: pos,
		Rel: types.Pointf{X: float32(ev.XRel), Y: float32(ev.YRel)},
	})
}

func (win *Window) processWheelEvent(ev *sdl.MouseWheelEvent) {
	win.inputEvent(&seat.MouseWheelEvent{
		Wheel: int(ev.Y),
	})
}

var gamepadButtons = map[sdl.GameControllerButton]seat.GamepadButton{
	sdl.CONTROLLER_BUTTON_A:             seat.GamepadA,
	sdl.CONTROLLER_BUTTON_B:             seat.GamepadB,
	sdl.CONTROLLER_BUTTON_X:             seat.GamepadX,
	sdl.CONTROLLER_BUTTON_Y:             seat.GamepadY,
	sdl.CONTROLLER_BUTTON_BACK:          seat.GamepadBack,
	sdl.CONTROLLER_BUTTON_GUIDE:         seat.GamepadGuide,
	sdl.CONTROLLER_BUTTON_START:         seat.GamepadStart,
	sdl.CONTROLLER_BUTTON_LEFTSTICK:     seat.GamepadLeftStick,
	sdl.CONTROLLER_BUTTON_RIGHTSTICK:    seat.GamepadRightStick,
	sdl.CONTROLLER_BUTTON_LEFTSHOULDER:  seat.GamepadLeftShoulder,
	sdl.CONTROLLER_BUTTON_RIGHTSHOULDER: seat.GamepadRightShoulder,
	sdl.CONTROLLER_BUTTON_DPAD_UP:       seat.GamepadDpadUp,
	sdl.CONTROLLER_BUTTON_DPAD_DOWN:     seat.GamepadDpadDown,
	sdl.CONTROLLER_BUTTON_DPAD_LEFT:     seat.GamepadDpadLeft,
	sdl.CONTROLLER_BUTTON_DPAD_RIGHT:    seat.GamepadDpadRight,
}

var gamepadAxes = map[sdl.GameControllerAxis]seat.GamepadAxis{
	sdl.CONTROLLER_AXIS_LEFTX:        seat.GamepadLeftX,
	sdl.CONTROLLER_AXIS_LEFTY:        seat.GamepadLeftY,
	sdl.CONTROLLER_AXIS_RIGHTX:       seat.GamepadRightX,
	sdl.CONTROLLER_AXIS_RIGHTY:       seat.GamepadRightY,
	sdl.CONTROLLER_AXIS_TRIGGERLEFT:  seat.GamepadTriggerLeft,
	sdl.CONTROLLER_AXIS_TRIGGERRIGHT: seat.GamepadTriggerRight,
}

func (win *Window) processGamepadButtonEvent(ev *sdl.ControllerButtonEvent) {
	b, ok := gamepadButtons[sdl.GameControllerButton(ev.Button)]
	if !ok {
		return
	}
	win.inputEvent(&seat.GamepadButtonEvent{
		Button:  b,
		Pressed: ev.State == sdl.PRESSED,
	})
}

func (win *Window) processGamepadAxisEvent(ev *sdl.ControllerAxisEvent) {
	a, ok := gamepadAxes[sdl.GameControllerAxis(ev.Axis)]
	if !ok {
		return
	}
	v := float32(ev.Value) / 32767
	v = max(-1, min(1, v))
	win.inputEvent(&seat.GamepadAxisEvent{Axis: a, Value: v})
}

func (win *Window) processGamepadDeviceEvent(ev *sdl.ControllerDeviceEvent) {
	switch ev.Type {
	case sdl.CONTROLLERDEVICEADDED:
		// Which is the device index here
		ind := int(ev.Which)
		if !sdl.IsGameController(ind) {
			return
		}
		ctrl := sdl.GameControllerOpen(ind)
		if ctrl == nil {
			win.log.Warn("cannot open game controller", "index", ind, "err", sdl.GetError())
			return
		}
		id := ctrl.Joystick().InstanceID()
		if _, ok := win.pads[id]; ok {
			ctrl.Close()
			return
		}
		if win.pads == nil {
			win.pads = make(map[sdl.JoystickID]*sdl.GameController)
		}
		win.pads[id] = ctrl
		name := ctrl.Name()
		win.log.Info("game controller connected", "name", name)
		win.inputEvent(&seat.GamepadDeviceEvent{Name: name, Connected: true})
	case sdl.CONTROLLERDEVICEREMOVED:
		// Which is the instance ID here
		ctrl, ok := win.pads[ev.Which]
		if !ok {
			return
		}
		name := ctrl.Name()
		ctrl.Close()
		delete(win.pads, ev.Which)
		win.log.Info("game controller disconnected", "name", name)
		win.inputEvent(&seat.GamepadDeviceEvent{Name: name, Connected: false})
	}
}

// closeGamepads closes all open game controllers.
func (win *Window) closeGamepads() {
	for id, ctrl := range win.pads {
		ctrl.Close()
		delete(win.pads, id)
	}
}
