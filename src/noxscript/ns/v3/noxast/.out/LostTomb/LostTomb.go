package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	wp4   ns.WaypointID
	gvar5 ns.ObjectGroupID
	gvar6 ns.ObjectGroupID
	obj7  ns.ObjectID
	obj8  ns.ObjectID
	obj9  ns.ObjectID
	obj10 ns.ObjectID
	obj11 ns.ObjectID
	obj12 ns.ObjectID
)

func ResetFistTrap01() {
	ns.ObjectGroupOn(gvar6)
	ns.ObjectGroupOn(gvar5)
	ns.AudioEvent(ns.PressurePlateRelease, wp4)
}
func ActivateFistTrap01() {
	ns.ObjectGroupOff(gvar6)
	ns.ObjectGroupOff(gvar5)
	ns.CastSpellLocationLocation(ns.SPELL_FIST, ns.GetWaypointX(wp4), ns.GetWaypointY(wp4), ns.GetWaypointX(wp4), ns.GetWaypointY(wp4))
	ns.FrameTimer(300, ResetFistTrap01)
}
func CloseRightGates() {
	ns.ObjectOn(obj11)
	ns.LockDoor(obj7)
	ns.LockDoor(obj8)
}
func OpenRightGates() {
	ns.ObjectOff(obj11)
	ns.UnlockDoor(obj7)
	ns.UnlockDoor(obj8)
	ns.FrameTimer(90, CloseRightGates)
}
func CloseLeftGates() {
	ns.ObjectOn(obj12)
	ns.LockDoor(obj9)
	ns.LockDoor(obj10)
}
func OpenLeftGates() {
	ns.ObjectOff(obj12)
	ns.UnlockDoor(obj9)
	ns.UnlockDoor(obj10)
	ns.FrameTimer(90, CloseLeftGates)
}
func MapInitialize() {
	obj7 = ns.Object("RightGate1")
	obj8 = ns.Object("RightGate2")
	obj9 = ns.Object("LeftGate1")
	obj10 = ns.Object("LeftGate2")
	obj11 = ns.Object("RightGateSwitch")
	obj12 = ns.Object("LeftGateSwitch")
	wp4 = ns.Waypoint("FistTarget01")
	gvar5 = ns.ObjectGroup("FistTriggers01")
	gvar6 = ns.ObjectGroup("FistLights")
	ns.LockDoor(obj7)
	ns.LockDoor(obj8)
	ns.LockDoor(obj9)
	ns.LockDoor(obj10)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
