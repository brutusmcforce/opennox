package legacy

import "github.com/brutusmcforce/opennox/server"

// noModifiers stands in for the modifiers of an object without init data. It is never written.
var noModifiers [4]*server.ModifierEff

// itemModifiers returns the four modifier slots of a weapon, armor or wand (server.ModifierInitData).
// An object of such a class that was created without init data (e.g. spawned bare by a script) has none:
// it gets four empty slots instead of a crash.
func itemModifiers(u *server.Object) []*server.ModifierEff {
	if u == nil || u.InitData == nil {
		return noModifiers[:]
	}
	return (*server.ModifierInitData)(u.InitData).Modifiers[:]
}
