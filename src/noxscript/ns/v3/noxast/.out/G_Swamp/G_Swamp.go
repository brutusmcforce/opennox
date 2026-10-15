package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	obj4   ns.ObjectID
	obj5   ns.ObjectID
	wp6    ns.WaypointID
	wp7    ns.WaypointID
	wp8    ns.WaypointID
	wp9    ns.WaypointID
	wp10   ns.WaypointID
	wp11   ns.WaypointID
	wp12   ns.WaypointID
	gvar13 ns.WallGroupID
	gvar14 ns.WallGroupID
	gvar15 ns.WallGroupID
	obj16  ns.ObjectID
	obj17  ns.ObjectID
	obj18  ns.ObjectID
	obj19  ns.ObjectID
	obj20  ns.ObjectID
	gvar21 ns.WallGroupID
	gvar22 ns.WallGroupID
	gvar23 ns.WallGroupID
	gvar24 ns.WallGroupID
	wp25   ns.WaypointID
	gvar26 int
	gvar27 int
	flag28 bool
	ivar29 int
	ivar30 int
	ivar31 int
	gvar32 int
	ivar33 int
	obj34  ns.ObjectID
	obj35  ns.ObjectID
	obj36  ns.ObjectID
	obj37  ns.ObjectID
	obj38  ns.ObjectID
	obj39  ns.ObjectID
	obj40  ns.ObjectID
	obj41  ns.ObjectID
	obj42  ns.ObjectID
	obj43  ns.ObjectID
	obj44  ns.ObjectID
	obj45  ns.ObjectID
	obj46  ns.ObjectID
	obj47  ns.ObjectID
	obj48  ns.ObjectID
	obj49  ns.ObjectID
	obj50  ns.ObjectID
	obj51  ns.ObjectID
	obj52  ns.ObjectID
	obj53  ns.ObjectID
	obj54  ns.ObjectID
	obj55  ns.ObjectID
	obj56  ns.ObjectID
	obj57  ns.ObjectID
	obj58  ns.ObjectID
	gvar59 ns.ObjectGroupID
	gvar60 ns.ObjectGroupID
	gvar61 ns.ObjectGroupID
	wp62   ns.WaypointID
	wp63   ns.WaypointID
	wp64   ns.WaypointID
	wp65   ns.WaypointID
	wp66   ns.WaypointID
	wp67   ns.WaypointID
	wp68   ns.WaypointID
	wp69   ns.WaypointID
	wp70   ns.WaypointID
	wp71   ns.WaypointID
	wp72   ns.WaypointID
	wp73   ns.WaypointID
	wp74   ns.WaypointID
	wp75   ns.WaypointID
	wp76   ns.WaypointID
	wp77   ns.WaypointID
	wp78   ns.WaypointID
	wp79   ns.WaypointID
	wp80   ns.WaypointID
	wp81   ns.WaypointID
	wp82   ns.WaypointID
	wp83   ns.WaypointID
	wp84   ns.WaypointID
	wp85   ns.WaypointID
	wp86   ns.WaypointID
	wp87   ns.WaypointID
	wp88   ns.WaypointID
	gvar89 ns.WallGroupID
	obj90  ns.ObjectID
	obj91  ns.ObjectID
	wp92   ns.WaypointID
	wp93   ns.WaypointID
	wp94   ns.WaypointID
	wp95   ns.WaypointID
	gvar96 ns.WallGroupID
	gvar97 ns.WallGroupID
)

func init() {
	gvar26 = 0
	gvar27 = 1
	flag28 = true
	ivar29 = 0
	ivar30 = 300
	ivar31 = 15
	gvar32 = gvar27
	ivar33 = 90
}
func InitializeBlockPuzzle() {
	obj4 = ns.Object("SpikeBlock01")
	obj5 = ns.Object("SpikeBlock02")
	wp6 = ns.Waypoint("SpikeBlock01WP01")
	wp7 = ns.Waypoint("SpikeBlock01WP02")
	wp8 = ns.Waypoint("SpikeBlock01WP03")
	wp9 = ns.Waypoint("SpikeBlock02WP01")
	wp10 = ns.Waypoint("SpikeBlock02WP02")
	wp11 = ns.Waypoint("SpikeBlock02WP03")
	wp12 = ns.Waypoint("BlockPuzzleAudioOrigin")
	gvar13 = ns.WallGroup("BlockPuzzleSecretWalls01")
	gvar14 = ns.WallGroup("BlockPuzzleSecretWalls02")
	gvar15 = ns.WallGroup("BlockPuzzleSecretWalls03")
}
func BlockPuzzleStep01() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.SpikeBlockMove, wp12)
	ns.Move(obj4, wp6)
	ns.Move(obj5, wp9)
	ns.WallGroupOpen(gvar13)
}
func BlockPuzzleStep02() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.SpikeBlockMove, wp6)
	ns.Move(obj4, wp7)
	ns.Move(obj5, wp10)
	ns.WallGroupOpen(gvar14)
	ns.WallGroupOpen(gvar15)
}
func BlockPuzzleStep03() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.SpikeBlockMove, wp7)
	ns.AudioEvent(ns.SpikeBlockMove, wp10)
	ns.Move(obj4, wp8)
	ns.Move(obj5, wp11)
}
func OpenBlockPuzzleSecretExit() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(ns.Wall(41, 73))
}
func InitializeKeyElevator() {
	obj16 = ns.Object("SilverKey")
	obj17 = ns.Object("KeyElevator")
	obj18 = ns.Object("SilverKeySpot")
	obj19 = ns.Object("RaiseKeySwitch")
	obj20 = ns.Object("KeyElevatorSwitch")
	gvar21 = ns.WallGroup("KeyElevatorWallgroup01")
	gvar22 = ns.WallGroup("KeyElevatorWallgroup02")
	gvar23 = ns.WallGroup("KeyElevatorWallgroup03")
	gvar24 = ns.WallGroup("KeyElevatorWallgroup04")
}
func LowerKeyElevatorWallgroup01() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar21)
}
func LowerKeyElevatorWallgroup02() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar22)
}
func LowerKeyElevatorWallgroup04() {
	ns.WallGroupOpen(gvar24)
}
func LowerKeyElevatorWallgroup03() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar24)
}
func DisableKeyElevator() {
	ns.ObjectOff(obj17)
	ns.MoveObject(obj16, ns.GetObjectX(obj18), ns.GetObjectY(obj18))
}
func RaiseSilverKeyElevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj17)
	ns.FrameTimer(50, DisableKeyElevator)
}
func EnableKeyElevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOff(obj19)
	ns.ObjectOn(obj17)
}
func InitializeSwamp() {
	obj90 = ns.Object("StoneBlock01")
	obj91 = ns.Object("OgreHut02Elevator")
	wp92 = ns.Waypoint("StoneBlock01WP")
	wp93 = ns.Waypoint("StoneBlock01AudioOrigin")
	wp94 = ns.Waypoint("OgreWarlordWP")
	wp95 = ns.Waypoint("OgreHut02ElevatorAudioOrigin")
	gvar96 = ns.WallGroup("WaspPathSecretWalls")
	gvar97 = ns.WallGroup("OgreHutExitWalls")
}
func InitializeOgreCamp() {
	obj34 = ns.Object("OgreHut01Elevator")
	obj35 = ns.Object("OgreCampExitElevator")
	obj36 = ns.Object("RotatingSpike01")
	obj37 = ns.Object("RotatingSpike04")
	obj38 = ns.Object("RotatingSpikeMover01")
	obj39 = ns.Object("RotatingSpikeMover02")
	obj40 = ns.Object("RotatingSpikeMover03")
	obj41 = ns.Object("RotatingSpikeMover04")
	obj42 = ns.Object("SilverKeyBlockSwitch")
	obj43 = ns.Object("SilverKeyBlock01")
	obj44 = ns.Object("SilverKeyBlock02")
	obj45 = ns.Object("CageDoor01")
	obj46 = ns.Object("CageDoor02")
	obj47 = ns.Object("SpikeBlock03")
	obj48 = ns.Object("SpikeBlock04")
	obj49 = ns.Object("SpikeBlock05")
	obj50 = ns.Object("SpikeBlock06")
	obj51 = ns.Object("SpikeBlock07")
	obj52 = ns.Object("SpikeBlock08")
	obj53 = ns.Object("IronBlock01")
	obj54 = ns.Object("IronBlock02")
	obj55 = ns.Object("IronBlock03")
	obj56 = ns.Object("IronBlock04")
	obj57 = ns.Object("BlockGateGenerator01")
	obj58 = ns.Object("BlockGateGenerator02")
	gvar59 = ns.ObjectGroup("ArrowTrapGroup01")
	gvar60 = ns.ObjectGroup("ArrowTrapGroup02")
	gvar61 = ns.ObjectGroup("ArrowTrapGroup03")
	wp62 = ns.Waypoint("SilverKeyBlock01Home")
	wp63 = ns.Waypoint("SilverKeyBlock02Home")
	wp64 = ns.Waypoint("SilverKeyBlock01Dest")
	wp65 = ns.Waypoint("SilverKeyBlock02Dest")
	wp66 = ns.Waypoint("SilverKeyPuzzleAudioOrigin")
	wp67 = ns.Waypoint("BlockPusherWP1")
	wp68 = ns.Waypoint("BlockPusherWP2")
	wp69 = ns.Waypoint("BlockPusherWP3")
	wp70 = ns.Waypoint("BlockPusherWP4")
	wp71 = ns.Waypoint("BlockPusherWP5")
	wp72 = ns.Waypoint("BlockPusherWP6")
	wp73 = ns.Waypoint("BlockPusherWP7")
	wp74 = ns.Waypoint("BlockPusherWP8")
	wp75 = ns.Waypoint("BlockPusherWP9")
	wp76 = ns.Waypoint("BlockPusherWP10")
	wp77 = ns.Waypoint("BlockPusherWP11")
	wp78 = ns.Waypoint("BlockPusherWP12")
	wp83 = ns.Waypoint("IronBlockWP01")
	wp84 = ns.Waypoint("IronBlockWP02")
	wp85 = ns.Waypoint("IronBlockWP03")
	wp86 = ns.Waypoint("IronBlockWP04")
	wp82 = ns.Waypoint("BlockGateAudioOrigin")
	wp79 = ns.Waypoint("BlockPusherRow1AudioOrigin")
	wp80 = ns.Waypoint("BlockPusherRow2AudioOrigin")
	wp81 = ns.Waypoint("BlockPusherRow3AudioOrigin")
	wp87 = ns.Waypoint("OgreElevator01AudioOrigin")
	wp88 = ns.Waypoint("OgreElevator02AudioOrigin")
	gvar89 = ns.WallGroup("OgreHut01GeneratorWalls")
	ns.LockDoor(obj45)
	ns.LockDoor(obj46)
	BlockPusherLoop()
	ns.FrameTimer(40, OgreCampSetup)
}
func MapInitialize() {
	wp25 = ns.Waypoint("WellWP")
	InitializeSwamp()
	InitializeBlockPuzzle()
	InitializeOgreCamp()
}
func WellOfRestoration() {
	ns.RestoreHealth(ns.GetCaller(), ns.MaxHealth(ns.GetCaller())-ns.CurrentHealth(ns.GetCaller()))
	ns.AudioEvent(ns.RestoreHealthName, wp25)
	ns.Print("GeneralPrint:WellSignRefresh")
}
func OgreCampSetup() {
	ns.ObjectOff(obj34)
	ns.ObjectOff(obj35)
}
func SilverKeyBlockTimer() {
	if !flag28 {
		goto LABEL1
	}
	if !(ivar29 >= ivar30) {
		goto LABEL2
	}
	ivar29 = 0
	CloseSilverKeyBlocks()
	goto LABEL1
LABEL2:
	ivar29 += ivar31
	ns.AudioEvent(ns.GenerateTick, wp66)
	ns.FrameTimer(ivar31, SilverKeyBlockTimer)
LABEL1:
	return
}
func BlockPusherLoop() {
	var v0 int
	switch v0 = gvar32; v0 {
	case gvar27:
		goto LABEL1
	case gvar26:
		goto LABEL2
	default:
		goto LABEL3
	}
LABEL1:
	ns.Move(obj47, wp67)
	ns.Move(obj48, wp70)
	ns.Move(obj49, wp72)
	ns.Move(obj50, wp73)
	ns.Move(obj51, wp75)
	ns.Move(obj52, wp78)
	gvar32 = gvar26
	goto LABEL3
LABEL2:
	ns.Move(obj47, wp68)
	ns.Move(obj48, wp69)
	ns.Move(obj49, wp71)
	ns.Move(obj50, wp74)
	ns.Move(obj51, wp76)
	ns.Move(obj52, wp77)
	gvar32 = gvar27
	goto LABEL3
LABEL3:
	ns.AudioEvent(ns.SpikeBlockMove, wp79)
	ns.AudioEvent(ns.SpikeBlockMove, wp80)
	ns.AudioEvent(ns.SpikeBlockMove, wp81)
	ns.FrameTimer(ivar33, BlockPusherLoop)
}
func EnableOgreHut01Elevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar89)
	ns.ObjectOn(obj34)
	ns.Print("GeneralPrint:MsgElevOn")
	ns.AudioEvent(ns.CreatureCageAppears, wp87)
	ns.AudioEvent(ns.ChangeSpellbar, wp87)
}
func ToggleRotatingSpikes() {
	if !ns.IsObjectOn(obj38) {
		goto LABEL1
	}
	ns.ObjectOff(obj38)
	ns.ObjectOff(obj39)
	ns.ObjectOff(obj40)
	ns.ObjectOff(obj41)
	goto LABEL2
LABEL1:
	ns.ObjectOn(obj38)
	ns.ObjectOn(obj39)
	ns.ObjectOn(obj40)
	ns.ObjectOn(obj41)
LABEL2:
	return
}
func BreakRotatingSpikeWall01() {
	if !ns.IsCaller(obj36) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallBreak(ns.Wall(57, 23))
LABEL1:
	return
}
func BreakRotatingSpikeWall02() {
	if !ns.IsCaller(obj36) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallBreak(ns.Wall(53, 19))
LABEL1:
	return
}
func BreakRotatingSpikeWall03() {
	if !ns.IsCaller(obj37) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallBreak(ns.Wall(55, 29))
LABEL1:
	return
}
func BreakRotatingSpikeWall04() {
	if !ns.IsCaller(obj37) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallBreak(ns.Wall(51, 25))
LABEL1:
	return
}
func BreakRotatingSpikeWall05() {
	if !ns.IsCaller(obj37) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallBreak(ns.Wall(47, 21))
LABEL1:
	return
}
func DisableSilverKeyBlocks() {
	ns.ObjectOff(ns.GetTrigger())
	flag28 = false
	ns.ObjectOff(obj42)
	ns.Move(obj43, wp64)
	ns.Move(obj44, wp65)
}
func OpenSilverKeyBlocks() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj43, wp64)
	ns.Move(obj44, wp65)
	ns.AudioEvent(ns.SpikeBlockMove, wp62)
	SilverKeyBlockTimer()
}
func CloseSilverKeyBlocks() {
	if !flag28 {
		goto LABEL1
	}
	ns.ObjectOn(obj42)
	ns.Move(obj43, wp62)
	ns.Move(obj44, wp63)
	ns.AudioEvent(ns.SpikeBlockMove, wp62)
LABEL1:
	return
}
func DisableArrowTrap01() {
	ns.ObjectGroupOff(gvar59)
}
func ArrowTrap01() {
	ns.ObjectGroupOn(gvar59)
	ns.FrameTimer(1, DisableArrowTrap01)
}
func DisableArrowTrap02() {
	ns.ObjectGroupOff(gvar60)
}
func ArrowTrap02() {
	ns.ObjectGroupOn(gvar60)
	ns.FrameTimer(1, DisableArrowTrap02)
}
func DisableArrowTrap03() {
	ns.ObjectGroupOff(gvar61)
}
func ArrowTrap03() {
	ns.ObjectGroupOn(gvar61)
	ns.FrameTimer(1, DisableArrowTrap03)
}
func UnlockOgreCage() {
	ns.ObjectOff(ns.GetTrigger())
	ns.UnlockDoor(obj45)
	ns.UnlockDoor(obj46)
	ns.Print("GeneralPrint:MsgDoorUnlock")
}
func EnableOgreCampExitElevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj35)
	ns.Print("GeneralPrint:MsgElevOn")
	ns.AudioEvent(ns.CreatureCageAppears, wp88)
	ns.AudioEvent(ns.ChangeSpellbar, wp88)
}
func BlockPusherRow1Boom() {
	if !(ns.IsCaller(obj47) || ns.IsCaller(obj48)) {
		goto LABEL1
	}
	ns.Effect(ns.JIGGLE, ns.GetWaypointX(wp79), ns.GetWaypointY(wp79), 4, 0)
	ns.AudioEvent(ns.HammerMissing, wp79)
LABEL1:
	return
}
func BlockPusherRow2Boom() {
	if !(ns.IsCaller(obj49) || ns.IsCaller(obj50)) {
		goto LABEL1
	}
	ns.Effect(ns.JIGGLE, ns.GetWaypointX(wp80), ns.GetWaypointY(wp80), 4, 0)
	ns.AudioEvent(ns.HammerMissing, wp80)
LABEL1:
	return
}
func BlockPusherRow3Boom() {
	if !(ns.IsCaller(obj51) || ns.IsCaller(obj52)) {
		goto LABEL1
	}
	ns.Effect(ns.JIGGLE, ns.GetWaypointX(wp81), ns.GetWaypointY(wp81), 4, 0)
	ns.AudioEvent(ns.HammerMissing, wp81)
LABEL1:
	return
}
func OpenOgreBlockGate() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj57)
	ns.ObjectOn(obj58)
	ns.Move(obj53, wp83)
	ns.Move(obj54, wp84)
	ns.Move(obj55, wp85)
	ns.Move(obj56, wp86)
	ns.WallGroupOpen(ns.WallGroup("OgreCampExitBlockWalls"))
	ns.AudioEvent(ns.SpikeBlockMove, wp82)
}
func OpenWaspPathExitWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar96)
}
func OpenOgreHutExitWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar97)
}
func EnableOgreHut02Elevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj91)
	ns.AudioEvent(ns.CreatureCageAppears, wp95)
	ns.AudioEvent(ns.ChangeSpellbar, wp95)
}
func GenerateOgreWarlord() {
	ns.Move(ns.GetCaller(), wp94)
}
func OpenOgreVillageBridge() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.SpikeBlockMove, wp93)
	ns.Move(obj90, wp92)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
