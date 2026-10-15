package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	flag4  bool
	obj5   ns.ObjectID
	obj6   ns.ObjectID
	obj7   ns.ObjectID
	obj8   ns.ObjectID
	obj9   ns.ObjectID
	obj10  ns.ObjectID
	obj11  ns.ObjectID
	obj12  ns.ObjectID
	obj13  ns.ObjectID
	wp14   ns.WaypointID
	wp15   ns.WaypointID
	wp16   ns.WaypointID
	wp17   ns.WaypointID
	wp18   ns.WaypointID
	wp19   ns.WaypointID
	wp20   ns.WaypointID
	wp21   ns.WaypointID
	wp22   ns.WaypointID
	wp23   ns.WaypointID
	wp24   ns.WaypointID
	wp25   ns.WaypointID
	wp26   ns.WaypointID
	wp27   ns.WaypointID
	wp28   ns.WaypointID
	wp29   ns.WaypointID
	wp30   ns.WaypointID
	gvar31 ns.WallGroupID
	gvar32 ns.WallGroupID
	gvar33 ns.ObjectGroupID
	gvar34 ns.ObjectGroupID
	gvar35 ns.ObjectGroupID
	obj36  ns.ObjectID
	obj37  ns.ObjectID
	gvar38 ns.ObjectGroupID
	gvar39 ns.ObjectGroupID
)

func init() {
	flag4 = false
}
func MapInitialize() {
	gvar31 = ns.WallGroup("GlobeWalls")
	gvar32 = ns.WallGroup("LowerPassage")
	gvar33 = ns.ObjectGroup("globes")
	gvar38 = ns.ObjectGroup("GlobeAnchors")
	gvar34 = ns.ObjectGroup("SecretPassageSwitches")
	gvar35 = ns.ObjectGroup("LowerPassageSwitches")
	gvar39 = ns.ObjectGroup("OrbAnchors")
	wp14 = ns.Waypoint("Waypoint3")
	wp15 = ns.Waypoint("Waypoint4")
	wp16 = ns.Waypoint("Waypoint1")
	wp17 = ns.Waypoint("Waypoint2")
	wp18 = ns.Waypoint("Waypoint6")
	wp19 = ns.Waypoint("Waypoint5")
	wp20 = ns.Waypoint("Waypoint7")
	wp21 = ns.Waypoint("Waypoint8")
	wp22 = ns.Waypoint("Waypoint27")
	wp23 = ns.Waypoint("Waypoint28")
	wp24 = ns.Waypoint("Waypoint31")
	wp25 = ns.Waypoint("Waypoint32")
	wp26 = ns.Waypoint("Waypoint29")
	wp27 = ns.Waypoint("Waypoint30")
	wp28 = ns.Waypoint("Waypoint11")
	wp29 = ns.Waypoint("Waypoint13")
	wp30 = ns.Waypoint("Waypoint26")
	obj5 = ns.Object("URBookcase")
	obj6 = ns.Object("LRBookcase")
	obj7 = ns.Object("ULBookcase")
	obj8 = ns.Object("LLBookcase")
	obj10 = ns.Object("OLBookcase")
	obj11 = ns.Object("ORBookcase")
	obj12 = ns.Object("ILBookcase")
	obj13 = ns.Object("IRBookcase")
	obj9 = ns.Object("SecretBookcase")
	obj36 = ns.Object("TrapSwitch")
	obj37 = ns.Object("PurpleLight")
	ns.ObjectGroupOn(gvar38)
	ns.ObjectGroupOn(gvar39)
}
func ReleaseManaTrap() {
	ns.Move(obj5, wp14)
	ns.Move(obj6, wp15)
	ns.Move(obj7, wp16)
	ns.Move(obj8, wp17)
	ns.WallGroupClose(gvar31)
	ns.ObjectOn(obj36)
}
func DisableGlobes() {
	ns.ObjectGroupOff(gvar33)
	ns.ObjectOff(obj37)
	ns.FrameTimer(5, ReleaseManaTrap)
}
func EnableGlobes() {
	ns.ObjectGroupOn(gvar33)
	ns.FrameTimer(30, DisableGlobes)
}
func ActivateManaTrap() {
	ns.ObjectOff(obj36)
	ns.Move(obj5, wp18)
	ns.Move(obj6, wp19)
	ns.Move(obj7, wp20)
	ns.Move(obj8, wp21)
	ns.WallGroupOpen(gvar31)
	ns.ObjectOn(obj37)
	ns.FrameTimer(70, EnableGlobes)
}
func OpenExit() {
	ns.WallOpen(ns.Wall(181, 165))
}
func CloseSecretPassageSEG2() {
	ns.WallClose(ns.Wall(220, 180))
	ns.ObjectGroupOn(gvar34)
}
func CloseSecretPassage() {
	ns.Move(obj9, wp29)
	ns.FrameTimer(20, CloseSecretPassageSEG2)
}
func OpenSecretPassage() {
	ns.ObjectGroupOff(gvar34)
	ns.WallOpen(ns.Wall(220, 180))
	ns.Move(obj9, wp28)
	ns.FrameTimer(120, CloseSecretPassage)
}
func CloseLowerPassageSEG3() {
	ns.ObjectGroupOn(gvar35)
	ns.WallGroupClose(gvar32)
}
func CloseLowerPassageSEG2() {
	ns.Move(obj10, wp22)
	ns.Move(obj11, wp23)
	ns.FrameTimer(20, CloseLowerPassageSEG3)
}
func CloseLowerPassage() {
	ns.Move(obj12, wp24)
	ns.Move(obj13, wp25)
	ns.FrameTimer(20, CloseLowerPassageSEG2)
}
func OpenLowerPassageSEG2() {
	ns.WallGroupOpen(gvar32)
	ns.Move(obj12, wp22)
	ns.Move(obj13, wp23)
	ns.FrameTimer(120, CloseLowerPassage)
}
func OpenLowerPassage() {
	ns.ObjectGroupOff(gvar35)
	ns.Move(obj10, wp26)
	ns.Move(obj11, wp27)
	ns.FrameTimer(8, OpenLowerPassageSEG2)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
