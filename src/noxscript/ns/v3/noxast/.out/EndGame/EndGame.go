package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	obj4 ns.ObjectID
	obj5 ns.ObjectID
	obj6 ns.ObjectID
	obj7 ns.ObjectID
)

func StartUp() {
	ns.ObjectOn(obj4)
	ns.ObjectOn(obj5)
	ns.ObjectOn(obj6)
	ns.ObjectOn(obj7)
}
func MapInitialize() {
	obj4 = ns.Object("Mover1")
	obj5 = ns.Object("Mover2")
	obj6 = ns.Object("Mover3")
	obj7 = ns.Object("Mover4")
	ns.FrameTimer(5, StartUp)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
