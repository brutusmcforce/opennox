package sdl

import "unsafe"

type (
	GLattr               int
	GLContext            unsafe.Pointer
	WindowFlags          uint64
	Keymod               uint16
	Scancode             uint32
	JoystickID           uint32
	GameControllerButton int
	GameControllerAxis   int
)

// Rect is SDL_Rect.
type Rect struct {
	X, Y, W, H int32
}

// Init flags (SDL_InitFlags).
const (
	INIT_AUDIO   = 0x00000010
	INIT_VIDEO   = 0x00000020
	INIT_GAMEPAD = 0x00002000
)

// OpenGL attributes (SDL_GLAttr). Note that SDL3 renumbered the context attributes.
const (
	GL_DOUBLEBUFFER          GLattr = 0x05
	GL_CONTEXT_MAJOR_VERSION GLattr = 0x11
	GL_CONTEXT_MINOR_VERSION GLattr = 0x12
	GL_CONTEXT_FLAGS         GLattr = 0x13
	GL_CONTEXT_PROFILE_MASK  GLattr = 0x14

	GL_CONTEXT_PROFILE_CORE            = 0x0001
	GL_CONTEXT_FORWARD_COMPATIBLE_FLAG = 0x0002
)

// Window flags (SDL_WindowFlags).
const (
	WINDOW_FULLSCREEN         WindowFlags = 0x00000001
	WINDOW_OPENGL             WindowFlags = 0x00000002
	WINDOW_HIDDEN             WindowFlags = 0x00000008
	WINDOW_BORDERLESS         WindowFlags = 0x00000010
	WINDOW_RESIZABLE          WindowFlags = 0x00000020
	WINDOW_HIGH_PIXEL_DENSITY WindowFlags = 0x00002000

	WINDOWPOS_UNDEFINED = 0x1FFF0000
	WINDOWPOS_CENTERED  = 0x2FFF0000
)

const (
	RELEASED = 0
	PRESSED  = 1
)

const (
	KMOD_LCTRL Keymod = 0x0040
	KMOD_RCTRL Keymod = 0x0080
	KMOD_CTRL         = KMOD_LCTRL | KMOD_RCTRL
)

const SCANCODE_V Scancode = 25

const (
	BUTTON_LEFT   = 1
	BUTTON_MIDDLE = 2
	BUTTON_RIGHT  = 3
)

// Gamepad buttons and axes (SDL_GamepadButton, SDL_GamepadAxis), with the Xbox names of SDL2.
const (
	CONTROLLER_BUTTON_A             GameControllerButton = 0 // SOUTH
	CONTROLLER_BUTTON_B             GameControllerButton = 1 // EAST
	CONTROLLER_BUTTON_X             GameControllerButton = 2 // WEST
	CONTROLLER_BUTTON_Y             GameControllerButton = 3 // NORTH
	CONTROLLER_BUTTON_BACK          GameControllerButton = 4
	CONTROLLER_BUTTON_GUIDE         GameControllerButton = 5
	CONTROLLER_BUTTON_START         GameControllerButton = 6
	CONTROLLER_BUTTON_LEFTSTICK     GameControllerButton = 7
	CONTROLLER_BUTTON_RIGHTSTICK    GameControllerButton = 8
	CONTROLLER_BUTTON_LEFTSHOULDER  GameControllerButton = 9
	CONTROLLER_BUTTON_RIGHTSHOULDER GameControllerButton = 10
	CONTROLLER_BUTTON_DPAD_UP       GameControllerButton = 11
	CONTROLLER_BUTTON_DPAD_DOWN     GameControllerButton = 12
	CONTROLLER_BUTTON_DPAD_LEFT     GameControllerButton = 13
	CONTROLLER_BUTTON_DPAD_RIGHT    GameControllerButton = 14

	CONTROLLER_AXIS_LEFTX        GameControllerAxis = 0
	CONTROLLER_AXIS_LEFTY        GameControllerAxis = 1
	CONTROLLER_AXIS_RIGHTX       GameControllerAxis = 2
	CONTROLLER_AXIS_RIGHTY       GameControllerAxis = 3
	CONTROLLER_AXIS_TRIGGERLEFT  GameControllerAxis = 4
	CONTROLLER_AXIS_TRIGGERRIGHT GameControllerAxis = 5
)

const PIXELFORMAT_ABGR8888 = 0x16762004
