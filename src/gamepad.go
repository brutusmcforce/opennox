package opennox

import (
	"image"
	"math"
	"time"

	"github.com/brutusmcforce/opennox/libs/client/keybind"
	"github.com/brutusmcforce/opennox/libs/client/seat"
	"github.com/brutusmcforce/opennox/libs/types"
)

// Game controller support.
//
// The controller is turned into the mouse and keyboard input the game knows, so menus, key bindings and the network
// code work unchanged. In the game:
//   - the left stick moves: the cursor is put around the player in the stick's direction (near walks, far runs) and
//     the MoveForward key is held;
//   - the right stick moves the cursor (aiming, inventory, dialogs);
//   - buttons press the key bound to an action (see gamepadGameButtons), so keyboard rebinding applies too.
//
// In menus both sticks move the cursor, A (or the right trigger) clicks and B is Escape.
//
// While the controller moves the cursor, the game draws the cursor itself (the system cursor stays where the mouse
// is); moving the mouse switches back.

// gamepadGameButtons are the actions of the controller buttons in the game.
var gamepadGameButtons = map[seat.GamepadButton]keybind.Event{
	seat.GamepadA:             keybind.EventAction,
	seat.GamepadB:             keybind.EventJump,
	seat.GamepadX:             keybind.EventQuickHealth,
	seat.GamepadY:             keybind.EventQuickMana,
	seat.GamepadDpadLeft:      keybind.EventInvokeSlot1,
	seat.GamepadDpadUp:        keybind.EventInvokeSlot2,
	seat.GamepadDpadRight:     keybind.EventInvokeSlot3,
	seat.GamepadDpadDown:      keybind.EventInvokeSlot4,
	seat.GamepadLeftShoulder:  keybind.EventInvertSpellTarget,
	seat.GamepadRightShoulder: keybind.EventNextSpellSet,
	seat.GamepadBack:          keybind.EventToggleInventory,
	seat.GamepadStart:         keybind.EventToggleQuitMenu,
	seat.GamepadLeftStick:     keybind.EventToggleMap,
	seat.GamepadRightStick:    keybind.EventQuickCurePoison,
}

// gamepadGameTriggers are the actions of the triggers in the game.
var gamepadGameTriggers = map[seat.GamepadAxis]keybind.Event{
	seat.GamepadTriggerLeft:  keybind.EventInvokeSlot5,
	seat.GamepadTriggerRight: keybind.EventAction,
}

// gamepadMenuButtons are the keys of the controller buttons in menus.
var gamepadMenuButtons = map[seat.GamepadButton]keybind.Key{
	seat.GamepadA:         keybind.KeyLeftMouse,
	seat.GamepadB:         keybind.KeyEsc,
	seat.GamepadStart:     keybind.KeyEsc,
	seat.GamepadDpadUp:    keybind.KeyUp,
	seat.GamepadDpadDown:  keybind.KeyDown,
	seat.GamepadDpadLeft:  keybind.KeyLeft,
	seat.GamepadDpadRight: keybind.KeyRight,
}

const (
	gamepadDeadZone   = 0.2  // stick values below are ignored
	gamepadMoveStart  = 0.35 // left stick: start moving
	gamepadMoveStop   = 0.25 // left stick: stop moving
	gamepadTriggerOn  = 0.5
	gamepadTriggerOff = 0.3
	gamepadWalkDist   = 0.08 // cursor distance from the player when moving, as part of the screen height: walking
	gamepadRunDist    = 0.38 // and running at full tilt
	gamepadCursorRate = 1.1  // cursor speed at full tilt: screen heights per second
)

// gamepadGame is what the controller needs from the game.
type gamepadGame interface {
	// GamepadInGame checks if a player character is controlled (otherwise the controller drives menus).
	GamepadInGame() bool
	// GamepadKey returns the key bound to an action, or 0.
	GamepadKey(ev keybind.Event) keybind.Key
}

// gamepadSource is a button or trigger that holds a key down.
type gamepadSource int

const (
	gamepadSourceTrigger = gamepadSource(100) // + axis
	gamepadSourceMove    = gamepadSource(200)
)

// gamepadSeat wraps a seat, turning controller events into mouse and keyboard events.
type gamepadSeat struct {
	seat.Seat
	game      gamepadGame
	onInput   []func(seat.InputEvent)
	textInput bool

	axes   [seat.GamepadTriggerRight + 1]float32
	cursor types.Pointf // cursor position in screen pixels
	sent   image.Point  // the last cursor position sent
	active bool         // the controller moved the cursor since the mouse did
	moving bool         // the move key is held for the left stick
	held   map[gamepadSource]keybind.Key
	last   time.Time
	now    func() time.Time
}

func newGamepadSeat(s seat.Seat, game gamepadGame) *gamepadSeat {
	g := &gamepadSeat{
		Seat: s, game: game,
		held: make(map[gamepadSource]keybind.Key),
		now:  time.Now,
	}
	s.OnInput(g.event)
	return g
}

func (g *gamepadSeat) OnInput(fnc func(ev seat.InputEvent)) {
	g.onInput = append(g.onInput, fnc)
}

func (g *gamepadSeat) ReplaceInputs(cfg seat.InputConfig) seat.InputConfig {
	old := g.onInput
	g.onInput = cfg
	return old
}

func (g *gamepadSeat) SetTextInput(enable bool) {
	g.textInput = enable
	g.Seat.SetTextInput(enable)
}

func (g *gamepadSeat) InputTick() {
	g.Seat.InputTick()
	g.tick()
}

// CursorActive checks if the controller drives the cursor (the game should draw it at the cursor position).
func (g *gamepadSeat) CursorActive() bool {
	return g.active
}

func (g *gamepadSeat) emit(ev seat.InputEvent) {
	for _, fnc := range g.onInput {
		fnc(ev)
	}
}

func (g *gamepadSeat) event(ev seat.InputEvent) {
	switch ev := ev.(type) {
	case *seat.GamepadButtonEvent:
		g.button(ev.Button, ev.Pressed)
	case *seat.GamepadAxisEvent:
		g.axis(ev.Axis, ev.Value)
	case *seat.GamepadDeviceEvent:
		if !ev.Connected {
			g.releaseAll()
			g.axes = [len(g.axes)]float32{}
		}
	case *seat.MouseMoveEvent:
		if !ev.Relative {
			g.cursor = types.Pointf{X: float32(ev.Pos.X), Y: float32(ev.Pos.Y)}
			g.sent = ev.Pos
			g.active = false
		}
		g.emit(ev)
	default:
		g.emit(ev)
	}
}

func (g *gamepadSeat) button(b seat.GamepadButton, pressed bool) {
	src := gamepadSource(b)
	if !pressed {
		g.release(src)
		return
	}
	var key keybind.Key
	if g.game.GamepadInGame() {
		if ev, ok := gamepadGameButtons[b]; ok {
			key = g.game.GamepadKey(ev)
		}
	} else {
		key = gamepadMenuButtons[b]
	}
	if key.IsKeyboard() && g.textInput && key != keybind.KeyEsc {
		return // typing (e.g. chat): only Escape
	}
	g.press(src, key)
}

func (g *gamepadSeat) axis(a seat.GamepadAxis, v float32) {
	if int(a) < 0 || int(a) >= len(g.axes) {
		return
	}
	g.axes[a] = v
	if a != seat.GamepadTriggerLeft && a != seat.GamepadTriggerRight {
		return
	}
	src := gamepadSourceTrigger + gamepadSource(a)
	_, down := g.held[src]
	if down && v < gamepadTriggerOff {
		g.release(src)
	} else if !down && v > gamepadTriggerOn {
		var key keybind.Key
		if g.game.GamepadInGame() {
			key = g.game.GamepadKey(gamepadGameTriggers[a])
		} else if a == seat.GamepadTriggerRight {
			key = keybind.KeyLeftMouse
		}
		g.press(src, key)
	}
}

func (g *gamepadSeat) press(src gamepadSource, key keybind.Key) {
	if key == 0 {
		return
	}
	if _, ok := g.held[src]; ok {
		return
	}
	g.held[src] = key
	g.emitKey(key, true)
}

func (g *gamepadSeat) release(src gamepadSource) {
	key, ok := g.held[src]
	if !ok {
		return
	}
	delete(g.held, src)
	g.emitKey(key, false)
}

func (g *gamepadSeat) releaseAll() {
	for src := range g.held {
		g.release(src)
	}
	g.moving = false
}

func (g *gamepadSeat) emitKey(key keybind.Key, pressed bool) {
	switch key {
	case keybind.KeyLeftMouse:
		g.emit(&seat.MouseButtonEvent{Button: seat.MouseButtonLeft, Pressed: pressed})
	case keybind.KeyRightMouse:
		g.emit(&seat.MouseButtonEvent{Button: seat.MouseButtonRight, Pressed: pressed})
	case keybind.KeyMiddleMouse:
		g.emit(&seat.MouseButtonEvent{Button: seat.MouseButtonMiddle, Pressed: pressed})
	case keybind.KeyMouseWheelUp:
		if pressed {
			g.emit(&seat.MouseWheelEvent{Wheel: 1})
		}
	case keybind.KeyMouseWheelDown:
		if pressed {
			g.emit(&seat.MouseWheelEvent{Wheel: -1})
		}
	default:
		if key.IsKeyboard() {
			g.emit(&seat.KeyboardEvent{Key: key, Pressed: pressed})
		}
	}
}

// deadZone maps a stick value to 0 inside the dead zone, and the rest to 0..1 (keeping the sign).
func deadZone(v float32) float32 {
	a := float32(math.Abs(float64(v)))
	if a < gamepadDeadZone {
		return 0
	}
	return float32(math.Copysign(float64((a-gamepadDeadZone)/(1-gamepadDeadZone)), float64(v)))
}

func (g *gamepadSeat) tick() {
	now := g.now()
	dt := float32(now.Sub(g.last).Seconds())
	g.last = now
	if dt > 0.1 || dt < 0 {
		dt = 0.1
	}
	sz := g.ScreenSize()
	if sz.X <= 0 || sz.Y <= 0 {
		return
	}
	lx, ly := g.axes[seat.GamepadLeftX], g.axes[seat.GamepadLeftY]
	rx, ry := deadZone(g.axes[seat.GamepadRightX]), deadZone(g.axes[seat.GamepadRightY])
	if !g.game.GamepadInGame() {
		if g.moving {
			g.release(gamepadSourceMove)
			g.moving = false
		}
		// menus: either stick moves the cursor
		g.moveCursor(sz, max(-1, min(1, rx+deadZone(lx))), max(-1, min(1, ry+deadZone(ly))), dt)
		return
	}
	m := float32(math.Hypot(float64(lx), float64(ly)))
	if g.moving && m < gamepadMoveStop {
		g.release(gamepadSourceMove)
		g.moving = false
	} else if !g.moving && m > gamepadMoveStart {
		g.moving = true
	}
	if g.moving {
		// the player is in the middle of the screen; the game moves towards the cursor
		h := float32(min(sz.X, sz.Y))
		t := max(0, min(1, (m-gamepadMoveStart)/(1-gamepadMoveStart)))
		r := h * (gamepadWalkDist + t*(gamepadRunDist-gamepadWalkDist))
		g.setCursor(sz, types.Pointf{
			X: float32(sz.X)/2 + lx/m*r,
			Y: float32(sz.Y)/2 + ly/m*r,
		})
		g.press(gamepadSourceMove, g.game.GamepadKey(keybind.EventMoveForward))
		return
	}
	g.moveCursor(sz, rx, ry, dt)
}

// moveCursor moves the cursor with a stick: faster the more it is tilted.
func (g *gamepadSeat) moveCursor(sz image.Point, x, y, dt float32) {
	if x == 0 && y == 0 {
		return
	}
	speed := gamepadCursorRate * float32(sz.Y)
	p := g.cursor
	p.X += x * float32(math.Abs(float64(x))) * speed * dt
	p.Y += y * float32(math.Abs(float64(y))) * speed * dt
	g.setCursor(sz, p)
}

func (g *gamepadSeat) setCursor(sz image.Point, p types.Pointf) {
	p.X = max(0, min(float32(sz.X-1), p.X))
	p.Y = max(0, min(float32(sz.Y-1), p.Y))
	g.cursor = p
	g.active = true
	pos := image.Point{X: int(p.X), Y: int(p.Y)}
	if pos == g.sent {
		return
	}
	rel := pos.Sub(g.sent)
	g.sent = pos
	g.emit(&seat.MouseMoveEvent{Pos: pos, Rel: types.Pointf{X: float32(rel.X), Y: float32(rel.Y)}})
}

// GamepadInGame implements gamepadGame.
func (c *Client) GamepadInGame() bool {
	return gameGetPlayState() == 3 && c.ClientPlayerUnit() != nil
}

// GamepadKey implements gamepadGame: the first key (of the keyboard or mouse) bound alone to the action.
func (c *Client) GamepadKey(ev keybind.Event) keybind.Key {
	if c.ctrl == nil {
		return 0
	}
	for _, b := range c.ctrl.listBindings() {
		if len(b.keys) != 1 || b.keys[0].IsJoystick() {
			continue
		}
		for _, e := range b.events {
			if e == ev {
				return b.keys[0]
			}
		}
	}
	return 0
}

// gamepadCursor checks if the controller drives the cursor, so that the game draws it instead of the system.
func (c *Client) gamepadCursor() bool {
	return c.gamepad != nil && c.gamepad.CursorActive()
}
