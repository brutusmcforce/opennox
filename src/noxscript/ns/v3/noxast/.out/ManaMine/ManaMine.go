package script

import (
	"strconv"
	"github.com/brutusmcforce/opennox/noxscript/ns/v3"
)

var (
	gvar4  ns.ObjectGroupID
	gvar5  ns.ObjectGroupID
	gvar6  ns.ObjectGroupID
	obj7   ns.ObjectID
	gvar8  int
	gvar9  int
	gvar10 int
	gvar11 int
	gvar12 int
	gvar13 int
	gvar14 int
	gvar15 int
	gvar16 int
	gvar17 int
	obj18  ns.ObjectID
	obj19  ns.ObjectID
	obj20  ns.ObjectID
	obj21  ns.ObjectID
	wp22   [20]ns.WaypointID
	wp23   ns.WaypointID
	wp24   ns.WaypointID
	wp25   ns.WaypointID
	wp26   ns.WaypointID
)

func MapInitialize() {
	gvar4 = ns.ObjectGroup("Flames")
	gvar5 = ns.ObjectGroup("BlueFlames")
	gvar6 = ns.ObjectGroup("BlockSwitches")
	obj7 = ns.Object("FlameSwitch")
	obj21 = ns.Object("BlueFlameSwitch")
	obj20 = ns.Object("BlockDoor2")
	obj19 = ns.Object("BlockDoor1")
	obj18 = ns.Object("SpikeBlock")
	wp22[0] = ns.Waypoint("Waypoint1")
	wp22[1] = ns.Waypoint("Waypoint2")
	wp22[2] = ns.Waypoint("Waypoint3")
	wp22[3] = ns.Waypoint("Waypoint4")
	wp22[4] = ns.Waypoint("Waypoint5")
	wp22[5] = ns.Waypoint("Waypoint6")
	wp22[6] = ns.Waypoint("Waypoint7")
	wp22[7] = ns.Waypoint("Waypoint8")
	wp22[8] = ns.Waypoint("Waypoint9")
	wp22[9] = ns.Waypoint("Waypoint10")
	wp22[10] = ns.Waypoint("Waypoint11")
	wp22[11] = ns.Waypoint("Waypoint12")
	wp22[12] = ns.Waypoint("Waypoint13")
	wp22[13] = ns.Waypoint("Waypoint14")
	wp22[14] = ns.Waypoint("Waypoint15")
	wp22[15] = ns.Waypoint("Waypoint16")
	wp22[16] = ns.Waypoint("Waypoint17")
	wp22[17] = ns.Waypoint("Waypoint18")
	wp22[18] = ns.Waypoint("Waypoint19")
	wp22[19] = ns.Waypoint("Waypoint20")
	wp23 = ns.Waypoint("Waypoint21")
	wp24 = ns.Waypoint("Waypoint22")
	wp25 = ns.Waypoint("Waypoint23")
	wp26 = ns.Waypoint("Waypoint24")
	ns.ObjectGroupOff(gvar5)
	ns.ObjectGroupOff(gvar4)
}
func FlameReturnSEG2() {
	ns.ObjectOn(obj7)
}
func FlameReturn() {
	var v0 int
	ns.ObjectGroupOff(gvar4)
	v0 = 1
	for {
		if !(v0 <= 10) {
			goto LABEL1
		}
		ns.Move(ns.Object("Flame"+strconv.Itoa(v0)), wp22[v0-1])
		v0 += 1
	}
LABEL1:
	ns.FrameTimer(38, FlameReturnSEG2)
}
func EnableFlameTrap() {
	var v0 int
	ns.ObjectGroupOn(gvar4)
	ns.ObjectOff(obj7)
	ns.AudioEvent(ns.DemonBreath, wp23)
	v0 = 1
	for {
		if !(v0 <= 10) {
			goto LABEL1
		}
		ns.Move(ns.Object("Flame"+strconv.Itoa(v0)), wp22[v0+9])
		v0 += 1
	}
LABEL1:
	ns.FrameTimer(38, FlameReturn)
}
func DisableSpikeBlock() {
	ns.Move(obj18, wp24)
	ns.UnlockDoor(obj19)
	ns.UnlockDoor(obj20)
	ns.ObjectGroupOn(gvar6)
}
func EnableSpikeBlock() {
	ns.ObjectGroupOff(gvar6)
	ns.LockDoor(obj19)
	ns.LockDoor(obj20)
	ns.Move(obj18, wp25)
	ns.SecondTimer(4, DisableSpikeBlock)
}
func DisableBlueFlames() {
	ns.ObjectGroupOff(gvar5)
	ns.ObjectOn(obj21)
}
func EnableBlueFlames() {
	ns.ObjectOff(obj21)
	ns.AudioEvent(ns.DemonBreath, wp26)
	ns.ObjectGroupOn(gvar5)
	ns.SecondTimer(1, DisableBlueFlames)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
