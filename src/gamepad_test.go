package opennox

import (
	"image"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/brutusmcforce/opennox/libs/client/keybind"
	"github.com/brutusmcforce/opennox/libs/client/seat"
)

// fakeSeat delivers events given to send, and reports a fixed screen size.
type fakeSeat struct {
	seat.Seat
	onInput []func(seat.InputEvent)
	text    bool
}

func (s *fakeSeat) ScreenSize() image.Point           { return image.Point{X: 1000, Y: 800} }
func (s *fakeSeat) OnInput(fnc func(seat.InputEvent)) { s.onInput = append(s.onInput, fnc) }
func (s *fakeSeat) InputTick()                        {}
func (s *fakeSeat) SetTextInput(enable bool)          { s.text = enable }
func (s *fakeSeat) send(evs ...seat.InputEvent) {
	for _, ev := range evs {
		for _, fnc := range s.onInput {
			fnc(ev)
		}
	}
}

type fakeGamepadGame struct {
	inGame bool
}

func (g *fakeGamepadGame) GamepadInGame() bool { return g.inGame }

// the default bindings: MOUSE_BUTTON_RIGHT = MoveForward, MOUSE_BUTTON_LEFT = Action, SPACE = Jump, ...
func (g *fakeGamepadGame) GamepadKey(ev keybind.Event) keybind.Key {
	switch ev {
	case keybind.EventMoveForward:
		return keybind.KeyRightMouse
	case keybind.EventAction:
		return keybind.KeyLeftMouse
	case keybind.EventJump:
		return keybind.KeySpace
	case keybind.EventInvokeSlot5:
		return keybind.KeyG
	}
	return 0
}

func newTestGamepad(inGame bool) (*fakeSeat, *gamepadSeat, *[]seat.InputEvent, *time.Time) {
	fs := &fakeSeat{}
	g := newGamepadSeat(fs, &fakeGamepadGame{inGame: inGame})
	var got []seat.InputEvent
	g.OnInput(func(ev seat.InputEvent) { got = append(got, ev) })
	now := time.Unix(1000, 0)
	g.now = func() time.Time { return now }
	g.InputTick() // start the clock
	return fs, g, &got, &now
}

func TestGamepadMove(t *testing.T) {
	fs, g, got, _ := newTestGamepad(true)
	// full tilt right: the cursor goes far right of the player (running) and MoveForward is held
	fs.send(&seat.GamepadAxisEvent{Axis: seat.GamepadLeftX, Value: 1})
	g.InputTick()
	require.Equal(t, []seat.InputEvent{
		&seat.MouseMoveEvent{Pos: image.Point{X: 500 + 304, Y: 400}, Rel: (*got)[0].(*seat.MouseMoveEvent).Rel},
		&seat.MouseButtonEvent{Button: seat.MouseButtonRight, Pressed: true},
	}, *got)
	require.True(t, g.CursorActive())
	*got = nil
	// a little up: walking distance, still moving
	fs.send(&seat.GamepadAxisEvent{Axis: seat.GamepadLeftX, Value: 0}, &seat.GamepadAxisEvent{Axis: seat.GamepadLeftY, Value: -0.36})
	g.InputTick()
	require.Len(t, *got, 1)
	mv := (*got)[0].(*seat.MouseMoveEvent)
	require.Equal(t, 500, mv.Pos.X)
	require.InDelta(t, 400-800*(0.08+0.30*0.01/0.65), mv.Pos.Y, 1)
	*got = nil
	// released: the move key goes up
	fs.send(&seat.GamepadAxisEvent{Axis: seat.GamepadLeftY, Value: 0})
	g.InputTick()
	require.Equal(t, []seat.InputEvent{&seat.MouseButtonEvent{Button: seat.MouseButtonRight, Pressed: false}}, *got)
}

func TestGamepadButtons(t *testing.T) {
	fs, g, got, _ := newTestGamepad(true)
	fs.send(
		&seat.GamepadButtonEvent{Button: seat.GamepadB, Pressed: true},
		&seat.GamepadButtonEvent{Button: seat.GamepadB, Pressed: false},
		&seat.GamepadAxisEvent{Axis: seat.GamepadTriggerRight, Value: 0.9},
		&seat.GamepadAxisEvent{Axis: seat.GamepadTriggerRight, Value: 0.4}, // between the thresholds: still held
		&seat.GamepadAxisEvent{Axis: seat.GamepadTriggerRight, Value: 0.1},
		&seat.GamepadAxisEvent{Axis: seat.GamepadTriggerLeft, Value: 1},
		&seat.GamepadButtonEvent{Button: seat.GamepadY, Pressed: true}, // QuickMana: not bound here
	)
	require.Equal(t, []seat.InputEvent{
		&seat.KeyboardEvent{Key: keybind.KeySpace, Pressed: true},
		&seat.KeyboardEvent{Key: keybind.KeySpace, Pressed: false},
		&seat.MouseButtonEvent{Button: seat.MouseButtonLeft, Pressed: true},
		&seat.MouseButtonEvent{Button: seat.MouseButtonLeft, Pressed: false},
		&seat.KeyboardEvent{Key: keybind.KeyG, Pressed: true},
	}, *got)
	*got = nil
	// a disconnected controller releases what it holds
	fs.send(&seat.GamepadDeviceEvent{Name: "pad", Connected: false})
	require.Equal(t, []seat.InputEvent{&seat.KeyboardEvent{Key: keybind.KeyG, Pressed: false}}, *got)
	*got = nil
	// while typing, buttons do not type
	g.SetTextInput(true)
	require.True(t, fs.text)
	fs.send(&seat.GamepadButtonEvent{Button: seat.GamepadB, Pressed: true})
	require.Empty(t, *got)
}

func TestGamepadMenuCursor(t *testing.T) {
	fs, g, got, now := newTestGamepad(false)
	// the mouse puts the cursor somewhere; other events pass through
	fs.send(&seat.MouseMoveEvent{Pos: image.Point{X: 100, Y: 100}})
	require.Len(t, *got, 1)
	require.False(t, g.CursorActive())
	*got = nil
	// full tilt right for 0.1s: 1.1 screen heights per second
	fs.send(&seat.GamepadAxisEvent{Axis: seat.GamepadRightX, Value: 1})
	*now = now.Add(100 * time.Millisecond)
	g.InputTick()
	require.Len(t, *got, 1)
	require.Equal(t, image.Point{X: 100 + 88, Y: 100}, (*got)[0].(*seat.MouseMoveEvent).Pos)
	require.True(t, g.CursorActive())
	*got = nil
	// A clicks, B is Escape
	fs.send(
		&seat.GamepadAxisEvent{Axis: seat.GamepadRightX, Value: 0},
		&seat.GamepadButtonEvent{Button: seat.GamepadA, Pressed: true},
		&seat.GamepadButtonEvent{Button: seat.GamepadA, Pressed: false},
		&seat.GamepadButtonEvent{Button: seat.GamepadB, Pressed: true},
	)
	require.Equal(t, []seat.InputEvent{
		&seat.MouseButtonEvent{Button: seat.MouseButtonLeft, Pressed: true},
		&seat.MouseButtonEvent{Button: seat.MouseButtonLeft, Pressed: false},
		&seat.KeyboardEvent{Key: keybind.KeyEsc, Pressed: true},
	}, *got)
	// the mouse takes the cursor back
	fs.send(&seat.MouseMoveEvent{Pos: image.Point{X: 5, Y: 5}})
	require.False(t, g.CursorActive())
}
