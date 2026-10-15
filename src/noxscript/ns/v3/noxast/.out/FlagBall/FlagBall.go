package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var obj4 ns.ObjectID

func MapInitialize() {
	obj4 = ns.Object("Ball")
}
func GoneOut() {
	if !ns.IsCaller(obj4) {
		goto LABEL1
	}
	ns.GoBackHome(obj4)
LABEL1:
	return
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
