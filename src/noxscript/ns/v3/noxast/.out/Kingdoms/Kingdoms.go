package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	obj4   ns.ObjectID
	obj5   ns.ObjectID
	obj6   ns.ObjectID
	obj7   ns.ObjectID
	obj8   ns.ObjectID
	obj9   ns.ObjectID
	wp10   ns.WaypointID
	wp11   ns.WaypointID
	wp12   ns.WaypointID
	wp13   ns.WaypointID
	wp14   ns.WaypointID
	wp15   ns.WaypointID
	wp16   ns.WaypointID
	wp17   ns.WaypointID
	wp18   ns.WaypointID
	wp19   ns.WaypointID
	wp20   ns.WaypointID
	wp21   ns.WaypointID
	gvar22 ns.WallGroupID
	gvar23 ns.WallGroupID
)

func BlueOpen2() {
	ns.Move(obj6, wp18)
	ns.Move(obj7, wp21)
	ns.FrameTimer(150, BlueClose)
}
func BlueClose() {
	ns.Move(obj6, wp17)
	ns.Move(obj7, wp20)
	ns.FrameTimer(23, BlueClose2)
}
func BlueClose2() {
	ns.Move(obj6, wp16)
	ns.Move(obj7, wp19)
	ns.FrameTimer(1, EnableBPlate)
}
func EnableBPlate() {
	ns.WallGroupClose(gvar22)
	ns.ObjectOn(obj9)
}
func RedOpen2() {
	ns.Move(obj4, wp12)
	ns.Move(obj5, wp15)
	ns.FrameTimer(150, RedClose)
}
func RedClose() {
	ns.Move(obj4, wp11)
	ns.Move(obj5, wp14)
	ns.FrameTimer(23, RedClose2)
}
func RedClose2() {
	ns.Move(obj4, wp10)
	ns.Move(obj5, wp13)
	ns.FrameTimer(1, EnableRPlate)
}
func EnableRPlate() {
	ns.WallGroupClose(gvar23)
	ns.ObjectOn(obj8)
}
func MapInitialize() {
	obj4 = ns.Object("RCase1")
	obj5 = ns.Object("RCase2")
	obj8 = ns.Object("RPlate")
	obj6 = ns.Object("BCase1")
	obj7 = ns.Object("BCase2")
	obj9 = ns.Object("BPlate")
	wp10 = ns.Waypoint("R1a")
	wp11 = ns.Waypoint("R1b")
	wp12 = ns.Waypoint("R1c")
	wp13 = ns.Waypoint("R2a")
	wp14 = ns.Waypoint("R2b")
	wp15 = ns.Waypoint("R2c")
	wp16 = ns.Waypoint("B1a")
	wp17 = ns.Waypoint("B1b")
	wp18 = ns.Waypoint("B1c")
	wp19 = ns.Waypoint("B2a")
	wp20 = ns.Waypoint("B2b")
	wp21 = ns.Waypoint("B2c")
	gvar22 = ns.WallGroup("BWalls")
	gvar23 = ns.WallGroup("RWalls")
}
func BlueOpen() {
	ns.ObjectOff(obj9)
	ns.WallGroupOpen(gvar22)
	ns.Move(obj6, wp17)
	ns.Move(obj7, wp20)
	ns.FrameTimer(13, BlueOpen2)
}
func RedOpen() {
	ns.ObjectOff(obj8)
	ns.WallGroupOpen(gvar23)
	ns.Move(obj4, wp11)
	ns.Move(obj5, wp14)
	ns.FrameTimer(13, RedOpen2)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
