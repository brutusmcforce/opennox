package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	gvar4 ns.ObjectGroupID
	gvar5 ns.ObjectGroupID
	gvar6 ns.WallGroupID
	gvar7 ns.WallGroupID
	wp8   ns.WaypointID
	wp9   ns.WaypointID
	wp10  ns.WaypointID
	wp11  ns.WaypointID
	wp12  ns.WaypointID
)

func EnableLightgroup1() {
	ns.ObjectGroupOff(gvar5)
	ns.ObjectGroupOn(gvar4)
	ns.FrameTimer(10, EnableLightgroup2)
}
func MapInitialize() {
	wp8 = ns.Waypoint("Waypoint1")
	wp9 = ns.Waypoint("Waypoint2")
	wp10 = ns.Waypoint("Waypoint3")
	wp11 = ns.Waypoint("Waypoint4")
	wp12 = ns.Waypoint("Waypoint5")
	gvar4 = ns.ObjectGroup("LanternGroup1")
	gvar5 = ns.ObjectGroup("LanternGroup2")
	gvar6 = ns.WallGroup("LeftSecretWalls")
	gvar7 = ns.WallGroup("RightSecretWalls")
}
func EnableLightgroup2() {
	ns.ObjectGroupOff(gvar4)
	ns.ObjectGroupOn(gvar5)
	ns.FrameTimer(10, EnableLightgroup1)
}
func LeftWallsOpen() {
	ns.WallGroupOpen(gvar6)
}
func LeftElevatorChime() {
	ns.AudioEvent(ns.Chime, wp8)
}
func RightWallsOpen() {
	ns.WallGroupOpen(gvar7)
}
func RightElevatorChime() {
	ns.AudioEvent(ns.Chime, wp9)
}
func OpenElevatorWalls() {
	ns.WallOpen(ns.Wall(126, 174))
}
func LeftTeleporterChime() {
	ns.AudioEvent(ns.MetallicBong, wp10)
}
func RightTeleporterChime() {
	ns.AudioEvent(ns.MetallicBong, wp11)
}
func PitChime() {
	ns.AudioEvent(ns.SmallGong, wp12)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
