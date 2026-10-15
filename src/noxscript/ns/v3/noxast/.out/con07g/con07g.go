package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

func MapInitialize() {
}
func MapEntry() {
	ns.Music(28, 100)
}
func PlayerDeath() {
	ns.DeathScreen(7)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	case "MapEntry":
		MapEntry()
	case "PlayerDeath":
		PlayerDeath()
	}
}
