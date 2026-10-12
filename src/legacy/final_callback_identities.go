package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

const (
	playerSectionGUI = iota
	playerSectionMetadata
	playerSectionMusic
)

var playerSectionKeys [3]byte

func playerSectionKey(id int) unsafe.Pointer { return unsafe.Pointer(&playerSectionKeys[id]) }

var playerFileCallbacks = make(map[unsafe.Pointer]func(unsafe.Pointer) int)
var screenParticleCallbacks = make(map[unsafe.Pointer]func(*noxrender.Viewport, *Nox_screenParticle) int)

func callPlayerFileCallback(key, arg unsafe.Pointer) int {
	if fn := playerFileCallbacks[key]; fn != nil {
		return fn(arg)
	}
	panic("unregistered player file callback")
}

func callPlayerFileSection(key unsafe.Pointer) int  { return callPlayerFileCallback(key, nil) }
func callPlayerFileMetadata(key unsafe.Pointer) int { return callPlayerFileCallback(key, nil) }

const (
	tooltipConversation = iota
	tooltipTeamAssign
	tooltipTeamDamage
)

var finalTooltipKeys [3]byte

func finalTooltipKey(id int) unsafe.Pointer { return unsafe.Pointer(&finalTooltipKeys[id]) }

var screenParticleKey, flameCleanseKey byte

func screenParticleCallbackKey() unsafe.Pointer { return unsafe.Pointer(&screenParticleKey) }
func flameCleanseCallbackKey() unsafe.Pointer   { return unsafe.Pointer(&flameCleanseKey) }

func init() {
	playerFileCallbacks[playerSectionKey(playerSectionGUI)] = func(unsafe.Pointer) int { return playerFileGUI() }
	playerFileCallbacks[playerSectionKey(playerSectionMetadata)] = func(unsafe.Pointer) int { return playerFileMetadata() }
	playerFileCallbacks[playerSectionKey(playerSectionMusic)] = func(unsafe.Pointer) int { return playerFileMusic() }
	screenParticleCallbacks[screenParticleCallbackKey()] = screenParticleDraw

	gui.RegisterTooltipCallbackGo(finalTooltipKey(tooltipConversation), func(*gui.Window, *gui.WindowData, uintptr) {})
	gui.RegisterTooltipCallbackGo(finalTooltipKey(tooltipTeamAssign), func(_ *gui.Window, data *gui.WindowData, _ uintptr) {
		serverOptionsTooltip(false, *(*byte)(unsafe.Pointer(data)))
	})
	gui.RegisterTooltipCallbackGo(finalTooltipKey(tooltipTeamDamage), func(_ *gui.Window, data *gui.WindowData, _ uintptr) {
		serverOptionsTooltip(true, *(*byte)(unsafe.Pointer(data)))
	})
	// FlameCleanse is assigned dynamically, not a named object-type update.
	server.RegisterObjectUpdateCallbackGo(flameCleanseCallbackKey(), temporaryFlameCleanse)
}

func callScreenParticleDraw(key unsafe.Pointer, vp *noxrender.Viewport, p *Nox_screenParticle) int {
	if fn := screenParticleCallbacks[key]; fn != nil {
		return fn(vp, p)
	}
	panic("unregistered screen particle callback")
}
