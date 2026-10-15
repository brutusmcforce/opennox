package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	obj4   [6]ns.ObjectID
	ivar5  int
	obj6   ns.ObjectID
	obj7   ns.ObjectID
	obj8   ns.ObjectID
	obj9   ns.ObjectID
	obj10  ns.ObjectID
	obj11  ns.ObjectID
	obj12  ns.ObjectID
	obj13  ns.ObjectID
	obj14  ns.ObjectID
	obj15  ns.ObjectID
	obj16  ns.ObjectID
	obj17  ns.ObjectID
	gvar18 ns.ObjectGroupID
	gvar19 ns.ObjectGroupID
	gvar20 ns.ObjectGroupID
	gvar21 ns.ObjectGroupID
	gvar22 ns.ObjectGroupID
	gvar23 ns.WallGroupID
	wp24   ns.WaypointID
	wp25   ns.WaypointID
	wp26   ns.WaypointID
	wp27   ns.WaypointID
	wp28   ns.WaypointID
	wp29   ns.WaypointID
	wp30   ns.WaypointID
	wp31   ns.WaypointID
	wp32   ns.WaypointID
	wp33   ns.WaypointID
	wp34   ns.WaypointID
	wp35   ns.WaypointID
	wp36   ns.WaypointID
	wp37   ns.WaypointID
)

func MapInitialize() {
	obj6 = ns.Object("SpikeBlock")
	obj7 = ns.Object("BlockTrigger")
	obj8 = ns.Object("DownstairsLeftBookcase")
	obj9 = ns.Object("DownstairsRightBookcase")
	obj10 = ns.Object("UpstairsLeftBookcase")
	obj11 = ns.Object("UpstairsRightBookcase")
	obj12 = ns.Object("OfficeBookcase")
	obj13 = ns.Object("LibraryBookcase")
	obj14 = ns.Object("DownstairsOrbMover")
	obj15 = ns.Object("UpstairsOrbMover")
	obj16 = ns.Object("OfficeOrbMover")
	obj17 = ns.Object("LibraryOrbMover")
	gvar23 = ns.WallGroup("BlockWalls")
	gvar22 = ns.ObjectGroup("BlockSwitches")
	gvar20 = ns.ObjectGroup("OfficeElevatorSwitches")
	gvar21 = ns.ObjectGroup("LibraryElevatorSwitches")
	gvar18 = ns.ObjectGroup("downstairsElevatorSwitches")
	gvar19 = ns.ObjectGroup("upstairsElevatorSwitches")
	wp37 = ns.Waypoint("Waypoint1")
	wp36 = ns.Waypoint("Waypoint3")
	wp25 = ns.Waypoint("Waypoint5")
	wp24 = ns.Waypoint("Waypoint6")
	wp26 = ns.Waypoint("Waypoint7")
	wp27 = ns.Waypoint("Waypoint8")
	wp28 = ns.Waypoint("Waypoint10")
	wp29 = ns.Waypoint("Waypoint11")
	wp31 = ns.Waypoint("Waypoint12")
	wp30 = ns.Waypoint("Waypoint13")
	wp33 = ns.Waypoint("Waypoint15")
	wp32 = ns.Waypoint("Waypoint17")
	wp35 = ns.Waypoint("Waypoint20")
	wp34 = ns.Waypoint("Waypoint22")
	ns.ObjectOn(obj14)
	ns.ObjectOn(obj15)
	ns.ObjectOn(obj16)
	ns.ObjectOn(obj17)
	obj4[0] = ns.Object("Door1")
	obj4[1] = ns.Object("Door2")
	obj4[2] = ns.Object("Door3")
	obj4[3] = ns.Object("Door4")
	obj4[4] = ns.Object("Door5")
	obj4[5] = ns.Object("Door6")
	ivar5 = 0
	for {
		if !(ivar5 <= 5) {
			goto LABEL1
		}
		ns.LockDoor(obj4[ivar5])
		ivar5 += 1
	}
LABEL1:
	return
}
func CloseSecretElevatorDownstairs() {
	ns.Move(obj8, wp24)
	ns.Move(obj9, wp25)
	ns.WallClose(ns.Wall(149, 143))
	ns.ObjectGroupOn(gvar18)
}
func OpenSecretElevatorDownstairs() {
	ns.ObjectGroupOff(gvar18)
	ns.Move(obj8, wp26)
	ns.Move(obj9, wp27)
	ns.WallOpen(ns.Wall(149, 143))
	ns.FrameTimer(90, CloseSecretElevatorDownstairs)
}
func CloseSecretElevatorUpstairs() {
	ns.ObjectGroupOn(gvar19)
	ns.Move(obj10, wp28)
	ns.Move(obj11, wp29)
	ns.WallClose(ns.Wall(202, 222))
}
func OpenSecretElevatorUpstairs() {
	ns.ObjectGroupOff(gvar19)
	ns.Move(obj10, wp30)
	ns.Move(obj11, wp31)
	ns.WallOpen(ns.Wall(202, 222))
	ns.FrameTimer(90, CloseSecretElevatorUpstairs)
}
func CloseOfficeWall() {
	ns.WallClose(ns.Wall(135, 125))
	ns.ObjectGroupOn(gvar20)
}
func MoveOfficeBookcaseBack() {
	ns.Move(obj12, wp32)
	ns.FrameTimer(20, CloseOfficeWall)
}
func OpenOfficeElevator() {
	ns.ObjectGroupOff(gvar20)
	ns.WallOpen(ns.Wall(135, 125))
	ns.Move(obj12, wp33)
	ns.FrameTimer(90, MoveOfficeBookcaseBack)
}
func CloseLibraryWall() {
	ns.WallClose(ns.Wall(204, 190))
	ns.ObjectGroupOn(gvar21)
}
func MoveLibraryBookcaseBack() {
	ns.Move(obj13, wp34)
	ns.FrameTimer(20, CloseLibraryWall)
}
func OpenLibraryElevator() {
	ns.ObjectGroupOff(gvar21)
	ns.WallOpen(ns.Wall(204, 190))
	ns.Move(obj13, wp35)
	ns.FrameTimer(90, MoveLibraryBookcaseBack)
}
func ToggleCellGates() {
	ivar5 = 0
	if !ns.IsLocked(obj4[ivar5]) {
		goto LABEL1
	}
	ivar5 = 0
	for {
		if !(ivar5 <= 5) {
			goto LABEL2
		}
		ns.UnlockDoor(obj4[ivar5])
		ns.Chat(ns.GetCaller(), "Fortress:GatesUnlocked")
		ivar5 += 1
	}
LABEL2:
	goto LABEL4
LABEL1:
	ivar5 = 0
	for {
		if !(ivar5 <= 5) {
			goto LABEL4
		}
		ns.LockDoor(obj4[ivar5])
		ns.Chat(ns.GetCaller(), "Fortress:GatesLocked")
		ivar5 += 1
	}
LABEL4:
	return
}
func BlockReturn() {
	ns.Move(obj6, wp36)
	ns.ObjectOn(obj7)
	ns.ObjectGroupOn(gvar22)
}
func SmushSniper() {
	ns.ObjectGroupOff(gvar22)
	ns.WallGroupOpen(gvar23)
	ns.Move(obj6, wp37)
	ns.FrameTimer(80, BlockReturn)
}
func CloseBlockWalls() {
	ns.WallGroupClose(gvar23)
	ns.ObjectOff(ns.GetTrigger())
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
