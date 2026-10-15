package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	obj4    ns.ObjectID
	obj5    ns.ObjectID
	gvar6   ns.WallGroupID
	gvar7   ns.WallGroupID
	gvar8   ns.WallGroupID
	obj9    ns.ObjectID
	obj10   ns.ObjectID
	obj11   ns.ObjectID
	obj12   ns.ObjectID
	obj13   ns.ObjectID
	obj14   ns.ObjectID
	obj15   ns.ObjectID
	obj16   ns.ObjectID
	obj17   ns.ObjectID
	obj18   ns.ObjectID
	obj19   ns.ObjectID
	obj20   ns.ObjectID
	gvar21  ns.ObjectGroupID
	gvar22  ns.ObjectGroupID
	gvar23  ns.ObjectGroupID
	gvar24  ns.ObjectGroupID
	gvar25  ns.ObjectGroupID
	gvar26  ns.ObjectGroupID
	gvar27  ns.ObjectGroupID
	gvar28  ns.ObjectGroupID
	gvar29  ns.ObjectGroupID
	gvar30  ns.ObjectGroupID
	gvar31  ns.ObjectGroupID
	gvar32  ns.ObjectGroupID
	gvar33  ns.ObjectGroupID
	gvar34  int
	gvar35  int
	gvar36  int
	gvar37  int
	gvar38  int
	ivar39  int
	obj40   ns.ObjectID
	obj41   ns.ObjectID
	obj42   ns.ObjectID
	obj43   ns.ObjectID
	obj44   ns.ObjectID
	obj45   ns.ObjectID
	obj46   ns.ObjectID
	obj47   ns.ObjectID
	obj48   ns.ObjectID
	obj49   ns.ObjectID
	obj50   ns.ObjectID
	obj51   ns.ObjectID
	gvar52  ns.ObjectGroupID
	wp53    ns.WaypointID
	wp54    ns.WaypointID
	wp55    ns.WaypointID
	wp56    ns.WaypointID
	wp57    ns.WaypointID
	wp58    ns.WaypointID
	wp59    ns.WaypointID
	wp60    ns.WaypointID
	wp61    ns.WaypointID
	wp62    ns.WaypointID
	wp63    ns.WaypointID
	wp64    ns.WaypointID
	wp65    ns.WaypointID
	wp66    ns.WaypointID
	wp67    ns.WaypointID
	wp68    ns.WaypointID
	wp69    ns.WaypointID
	wp70    ns.WaypointID
	wp71    ns.WaypointID
	wp72    ns.WaypointID
	wp73    ns.WaypointID
	obj74   ns.ObjectID
	obj75   ns.ObjectID
	obj76   ns.ObjectID
	obj77   ns.ObjectID
	obj78   ns.ObjectID
	obj79   ns.ObjectID
	obj80   ns.ObjectID
	obj81   ns.ObjectID
	wp82    ns.WaypointID
	wp83    ns.WaypointID
	wp84    ns.WaypointID
	wp85    ns.WaypointID
	wp86    ns.WaypointID
	wp87    ns.WaypointID
	wp88    ns.WaypointID
	wp89    ns.WaypointID
	wp90    ns.WaypointID
	wp91    ns.WaypointID
	wp92    ns.WaypointID
	flag93  bool
	gvar94  ns.WallGroupID
	gvar95  ns.WallGroupID
	gvar96  ns.WallGroupID
	gvar97  ns.WallGroupID
	gvar98  ns.WallGroupID
	ivar99  int
	ivar100 int
	obj101  ns.ObjectID
	obj102  ns.ObjectID
	obj103  ns.ObjectID
	obj104  ns.ObjectID
	obj105  ns.ObjectID
	obj106  ns.ObjectID
	wp107   ns.WaypointID
	wp108   ns.WaypointID
	wp109   ns.WaypointID
	wp110   ns.WaypointID
	wp111   ns.WaypointID
	wp112   ns.WaypointID
	gvar113 ns.WallGroupID
	gvar114 ns.WallGroupID
	obj115  [13]ns.ObjectID
	obj116  [13]ns.ObjectID
	ivar117 int
	wp118   ns.WaypointID
	wp119   ns.WaypointID
	obj120  ns.ObjectID
	obj121  ns.ObjectID
	wp122   ns.WaypointID
	wp123   ns.WaypointID
	wp124   ns.WaypointID
	gvar125 ns.ObjectGroupID
	gvar126 ns.WallGroupID
	gvar127 int
	gvar128 int
	gvar129 int
	gvar130 ns.WallGroupID
	gvar131 ns.WallGroupID
	gvar132 ns.WallGroupID
	gvar133 ns.WallGroupID
	gvar134 ns.WallGroupID
	gvar135 ns.WallGroupID
)

func init() {
	gvar34 = 0
	gvar35 = 1
	gvar36 = gvar35
	gvar37 = gvar35
	gvar38 = gvar35
	ivar39 = 55
	flag93 = false
	ivar99 = 0
	ivar100 = 90
	ivar117 = 120
	gvar127 = 0
	gvar128 = 1
	gvar129 = 0
}
func OpenAirlock() {
	ns.ObjectOn(obj5)
	ns.WallGroupOpen(gvar6)
}
func InitializeAirlock() {
	obj4 = ns.Object("AirlockButton")
	obj5 = ns.Object("AirlockResetButton")
	gvar6 = ns.WallGroup("AirlockFrontWalls")
	gvar7 = ns.WallGroup("AirlockBackWalls")
	ns.ObjectOff(obj5)
}
func CloseAirlock() {
	ns.ObjectOff(obj4)
	ns.WallGroupClose(gvar7)
	ns.FrameTimer(5, OpenAirlock)
}
func ResetAirlock() {
	ns.ObjectOff(obj5)
	ns.WallGroupClose(gvar6)
	ns.WallGroupOpen(gvar7)
	ns.ObjectOn(obj4)
}
func InitializeAmbush() {
	gvar8 = ns.WallGroup("AmbushWalls")
}
func Ambush() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar8)
}
func InitializeArrowTraps() {
	obj9 = ns.Object("ArrowTrap01")
	obj10 = ns.Object("ArrowTrap02")
	obj11 = ns.Object("ArrowTrap03")
	obj12 = ns.Object("ArrowTrap04")
	obj13 = ns.Object("ArrowTrap05")
	obj14 = ns.Object("ArrowTrap06")
	obj15 = ns.Object("ArrowTrap07")
	obj16 = ns.Object("ArrowTrap08")
	obj17 = ns.Object("ArrowTrap09")
	obj18 = ns.Object("ArrowTrap10")
	obj19 = ns.Object("ArrowTrap11")
	obj20 = ns.Object("ArrowTrap12")
	gvar21 = ns.ObjectGroup("ArrowTrapGroup01")
	gvar22 = ns.ObjectGroup("ArrowTrapGroup02")
	gvar23 = ns.ObjectGroup("ArrowTrapGroup03")
	gvar24 = ns.ObjectGroup("ArrowTrapGroup04")
	gvar25 = ns.ObjectGroup("ArrowTrapGroup05")
	gvar26 = ns.ObjectGroup("ArrowTrapGroup06")
	gvar27 = ns.ObjectGroup("ArrowTrapGroup07")
	gvar28 = ns.ObjectGroup("ArrowTrapGroup08")
	gvar29 = ns.ObjectGroup("ArrowTrapGroup09")
	gvar30 = ns.ObjectGroup("ArrowTrapGroup10")
	gvar31 = ns.ObjectGroup("ArrowTrapGroup11")
	gvar32 = ns.ObjectGroup("ArrowTrapGroup12")
	gvar33 = ns.ObjectGroup("ArrowTrapGroup13")
}
func DisableArrowTrap02Row01() {
	ns.ObjectOff(obj9)
	ns.ObjectOff(obj10)
}
func ArrowTrap02Row01() {
	ns.ObjectOn(obj9)
	ns.ObjectOn(obj10)
	ns.FrameTimer(1, DisableArrowTrap02Row01)
}
func DisableArrowTrap02Row02() {
	ns.ObjectOff(obj9)
	ns.ObjectOff(obj10)
	ns.ObjectOff(obj11)
}
func ArrowTrap02Row02() {
	ns.ObjectOn(obj9)
	ns.ObjectOn(obj10)
	ns.ObjectOn(obj11)
	ns.FrameTimer(1, DisableArrowTrap02Row02)
}
func DisableArrowTrap02Row03() {
	ns.ObjectOff(obj10)
	ns.ObjectOff(obj11)
	ns.ObjectOff(obj12)
}
func ArrowTrap02Row03() {
	ns.ObjectOn(obj10)
	ns.ObjectOn(obj11)
	ns.ObjectOn(obj12)
	ns.FrameTimer(1, DisableArrowTrap02Row03)
}
func DisableArrowTrap02Row04() {
	ns.ObjectOff(obj11)
	ns.ObjectOff(obj12)
	ns.ObjectOff(obj13)
}
func ArrowTrap02Row04() {
	ns.ObjectOn(obj11)
	ns.ObjectOn(obj12)
	ns.ObjectOn(obj13)
	ns.FrameTimer(1, DisableArrowTrap02Row04)
}
func DisableArrowTrap02Row05() {
	ns.ObjectOff(obj12)
	ns.ObjectOff(obj13)
	ns.ObjectOff(obj14)
}
func ArrowTrap02Row05() {
	ns.ObjectOn(obj12)
	ns.ObjectOn(obj13)
	ns.ObjectOn(obj14)
	ns.FrameTimer(1, DisableArrowTrap02Row05)
}
func DisableArrowTrap02Row06() {
	ns.ObjectOff(obj13)
	ns.ObjectOff(obj14)
	ns.ObjectOff(obj15)
}
func ArrowTrap02Row06() {
	ns.ObjectOn(obj13)
	ns.ObjectOn(obj14)
	ns.ObjectOn(obj15)
	ns.FrameTimer(1, DisableArrowTrap02Row06)
}
func DisableArrowTrap02Row07() {
	ns.ObjectOff(obj14)
	ns.ObjectOff(obj15)
	ns.ObjectOff(obj16)
}
func ArrowTrap02Row07() {
	ns.ObjectOn(obj14)
	ns.ObjectOn(obj15)
	ns.ObjectOn(obj16)
	ns.FrameTimer(1, DisableArrowTrap02Row07)
}
func DisableArrowTrap02Row08() {
	ns.ObjectOff(obj15)
	ns.ObjectOff(obj16)
	ns.ObjectOff(obj17)
}
func ArrowTrap02Row08() {
	ns.ObjectOn(obj15)
	ns.ObjectOn(obj16)
	ns.ObjectOn(obj17)
	ns.FrameTimer(1, DisableArrowTrap02Row08)
}
func DisableArrowTrap02Row09() {
	ns.ObjectOff(obj16)
	ns.ObjectOff(obj17)
	ns.ObjectOff(obj18)
}
func ArrowTrap02Row09() {
	ns.ObjectOn(obj16)
	ns.ObjectOn(obj17)
	ns.ObjectOn(obj18)
	ns.FrameTimer(1, DisableArrowTrap02Row09)
}
func DisableArrowTrap02Row10() {
	ns.ObjectOff(obj17)
	ns.ObjectOff(obj18)
}
func ArrowTrap02Row10() {
	ns.ObjectOn(obj17)
	ns.ObjectOn(obj18)
	ns.FrameTimer(1, DisableArrowTrap02Row10)
}
func DisableArrowTrap03Row01() {
	ns.ObjectGroupOff(gvar22)
}
func ArrowTrap03Row01() {
	ns.ObjectGroupOn(gvar22)
	ns.FrameTimer(1, DisableArrowTrap03Row01)
}
func DisableArrowTrap03Row02() {
	ns.ObjectGroupOff(gvar23)
}
func ArrowTrap03Row02() {
	ns.ObjectGroupOn(gvar23)
	ns.FrameTimer(1, DisableArrowTrap03Row02)
}
func DisableArrowTrap03Row03() {
	ns.ObjectGroupOff(gvar24)
}
func ArrowTrap03Row03() {
	ns.ObjectGroupOn(gvar24)
	ns.FrameTimer(1, DisableArrowTrap03Row03)
}
func DisableArrowTrap03Row04() {
	ns.ObjectGroupOff(gvar25)
}
func ArrowTrap03Row04() {
	ns.ObjectGroupOn(gvar25)
	ns.FrameTimer(1, DisableArrowTrap03Row04)
}
func DisableArrowTrap03Row05() {
	ns.ObjectGroupOff(gvar26)
}
func ArrowTrap03Row05() {
	ns.ObjectGroupOn(gvar26)
	ns.FrameTimer(1, DisableArrowTrap03Row05)
}
func DisableArrowTrap03Row06() {
	ns.ObjectGroupOff(gvar27)
}
func ArrowTrap03Row06() {
	ns.ObjectGroupOn(gvar27)
	ns.FrameTimer(1, DisableArrowTrap03Row06)
}
func DisableArrowTrap04Row01() {
	ns.ObjectOff(obj19)
}
func ArrowTrap04Row01() {
	ns.ObjectOn(obj19)
	ns.FrameTimer(1, DisableArrowTrap04Row01)
}
func DisableArrowTrap04Row02() {
	ns.ObjectOff(obj20)
}
func ArrowTrap04Row02() {
	ns.ObjectOn(obj20)
	ns.FrameTimer(1, DisableArrowTrap04Row02)
}
func DisableArrowTrap05Row01() {
	ns.ObjectGroupOff(gvar28)
}
func ArrowTrap05Row01() {
	ns.ObjectGroupOn(gvar28)
	ns.FrameTimer(1, DisableArrowTrap05Row01)
}
func DisableArrowTrap05Row02() {
	ns.ObjectGroupOff(gvar29)
}
func ArrowTrap05Row02() {
	ns.ObjectGroupOn(gvar29)
	ns.FrameTimer(1, DisableArrowTrap05Row02)
}
func DisableArrowTrap05Row03() {
	ns.ObjectGroupOff(gvar30)
}
func ArrowTrap05Row03() {
	ns.ObjectGroupOn(gvar30)
	ns.FrameTimer(1, DisableArrowTrap05Row03)
}
func DisableArrowTrap05Row04() {
	ns.ObjectGroupOff(gvar31)
}
func ArrowTrap05Row04() {
	ns.ObjectGroupOn(gvar31)
	ns.FrameTimer(1, DisableArrowTrap05Row04)
}
func DisableArrowTrap05Row05() {
	ns.ObjectGroupOff(gvar32)
}
func ArrowTrap05Row05() {
	ns.ObjectGroupOn(gvar32)
	ns.FrameTimer(1, DisableArrowTrap05Row05)
}
func DisableArrowTrap05Row06() {
	ns.ObjectGroupOff(gvar33)
}
func ArrowTrap05Row06() {
	ns.ObjectGroupOn(gvar33)
	ns.FrameTimer(1, DisableArrowTrap05Row06)
}
func BlockCrusherRow1Loop() {
	var v0 int
	switch v0 = gvar36; v0 {
	case gvar34:
		goto LABEL1
	case gvar35:
		goto LABEL2
	default:
		goto LABEL3
	}
LABEL1:
	ns.Move(obj40, wp59)
	ns.Move(obj41, wp60)
	gvar36 = gvar35
	goto LABEL3
LABEL2:
	ns.Move(obj40, wp53)
	ns.Move(obj41, wp54)
	gvar36 = gvar34
	goto LABEL3
LABEL3:
	ns.FrameTimer(ivar39, BlockCrusherRow1Loop)
}
func BlockCrusherRow2Loop() {
	var v0 int
	switch v0 = gvar37; v0 {
	case gvar34:
		goto LABEL1
	case gvar35:
		goto LABEL2
	default:
		goto LABEL3
	}
LABEL1:
	ns.Move(obj42, wp61)
	ns.Move(obj43, wp62)
	gvar37 = gvar35
	goto LABEL3
LABEL2:
	ns.Move(obj42, wp55)
	ns.Move(obj43, wp56)
	gvar37 = gvar34
	goto LABEL3
LABEL3:
	ns.FrameTimer(ivar39, BlockCrusherRow2Loop)
}
func BlockCrusherRow3Loop() {
	var v0 int
	switch v0 = gvar38; v0 {
	case gvar34:
		goto LABEL1
	case gvar35:
		goto LABEL2
	default:
		goto LABEL3
	}
LABEL1:
	ns.Move(obj44, wp63)
	ns.Move(obj45, wp64)
	gvar38 = gvar35
	goto LABEL3
LABEL2:
	ns.Move(obj44, wp57)
	ns.Move(obj45, wp58)
	gvar38 = gvar34
	goto LABEL3
LABEL3:
	ns.FrameTimer(ivar39, BlockCrusherRow3Loop)
}
func InitializeBlockCrusher() {
	obj40 = ns.Object("SpikeBlock06")
	obj41 = ns.Object("SpikeBlock07")
	obj42 = ns.Object("SpikeBlock08")
	obj43 = ns.Object("SpikeBlock09")
	obj44 = ns.Object("SpikeBlock10")
	obj45 = ns.Object("SpikeBlock11")
	obj46 = ns.Object("SpikeBlock14")
	obj47 = ns.Object("SpikeBlock15")
	obj48 = ns.Object("SpikeBlock16")
	obj49 = ns.Object("SpikeBlock17")
	obj50 = ns.Object("Block12Mover")
	obj51 = ns.Object("Block13Mover")
	gvar52 = ns.ObjectGroup("DuckAroundPressurePlates")
	wp53 = ns.Waypoint("Block06Dest")
	wp54 = ns.Waypoint("Block07Dest")
	wp55 = ns.Waypoint("Block08Dest")
	wp56 = ns.Waypoint("Block09Dest")
	wp57 = ns.Waypoint("Block10Dest")
	wp58 = ns.Waypoint("Block11Dest")
	wp59 = ns.Waypoint("Block06Home")
	wp60 = ns.Waypoint("Block07Home")
	wp61 = ns.Waypoint("Block08Home")
	wp62 = ns.Waypoint("Block09Home")
	wp63 = ns.Waypoint("Block10Home")
	wp64 = ns.Waypoint("Block11Home")
	wp65 = ns.Waypoint("Block14WP")
	wp66 = ns.Waypoint("Block15WP")
	wp67 = ns.Waypoint("Block16WP")
	wp68 = ns.Waypoint("Block17WP")
	wp69 = ns.Waypoint("CrushSFX1WP")
	wp70 = ns.Waypoint("CrushSFX2WP")
	wp71 = ns.Waypoint("CrushSFX3WP")
	wp72 = ns.Waypoint("DuckAroundAudioOrigin")
	wp73 = ns.Waypoint("BlockAudioOrigin01")
}
func EnableBlockCrusher() {
	ns.ObjectOff(ns.GetTrigger())
	BlockCrusherRow1Loop()
	ns.FrameTimer(20, BlockCrusherRow2Loop)
	ns.FrameTimer(40, BlockCrusherRow3Loop)
}
func StartBlocks() {
	ns.ObjectOn(obj50)
	ns.ObjectOn(obj51)
	ns.AudioEvent(ns.SpikeBlockMove, wp73)
}
func StopBlocks() {
	ns.ObjectOff(obj50)
	ns.ObjectOff(obj51)
}
func EnableDuckAroundBlocks() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj46, wp65)
	ns.Move(obj47, wp66)
	ns.Move(obj48, wp67)
	ns.Move(obj49, wp68)
	ns.AudioEvent(ns.SpikeBlockMove, wp72)
}
func CrushSFX1() {
	ns.AudioEvent(ns.HammerMissing, wp69)
}
func CrushSFX2() {
	ns.AudioEvent(ns.HammerMissing, wp70)
}
func CrushSFX3() {
	ns.AudioEvent(ns.HammerMissing, wp71)
}
func EnableBlockTrap() {
	flag93 = false
	ns.ObjectOn(obj79)
	ns.ObjectOn(obj80)
	ns.AudioEvent(ns.SmallGong, wp92)
}
func BlockTrapResetComplete() {
	if !ns.IsCaller(obj76) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar94)
	ns.FrameTimer(300, EnableBlockTrap)
LABEL1:
	return
}
func InitializeBlockTrap() {
	obj74 = ns.Object("SpikeBlock01")
	obj75 = ns.Object("SpikeBlock02")
	obj76 = ns.Object("SpikeBlock03")
	obj77 = ns.Object("SpikeBlock04")
	obj78 = ns.Object("SpikeBlock05")
	obj79 = ns.Object("BlockTrapResetButton")
	obj80 = ns.Object("BlockTrapEndPlate")
	obj81 = ns.Object("BlockTrapStartPlate")
	wp82 = ns.Waypoint("Block01Home")
	wp83 = ns.Waypoint("Block02Home")
	wp84 = ns.Waypoint("Block03Home")
	wp85 = ns.Waypoint("Block04Home")
	wp86 = ns.Waypoint("Block05Home")
	wp87 = ns.Waypoint("Block01Dest")
	wp88 = ns.Waypoint("Block02Dest")
	wp89 = ns.Waypoint("Block03Dest")
	wp90 = ns.Waypoint("Block04Dest")
	wp91 = ns.Waypoint("Block05Dest")
	wp92 = ns.Waypoint("BlockTrapAudioOrigin")
	gvar94 = ns.WallGroup("BlockTrapPathWalls")
	gvar95 = ns.WallGroup("BlockTrapButtonWall")
	gvar96 = ns.WallGroup("BlockTrapMonsterWall")
	gvar97 = ns.WallGroup("BlockTrapSecretWall")
	gvar98 = ns.WallGroup("BlockTrapExitWalls")
}
func TriggerBlockTrap() {
	if !(!flag93 && ns.HasClass(ns.GetCaller(), ns.PLAYER)) {
		goto LABEL1
	}
	flag93 = true
	ns.WallGroupClose(gvar94)
	ns.WallGroupOpen(gvar95)
	ns.Move(obj74, wp87)
	ns.Move(obj75, wp88)
	ns.Move(obj76, wp89)
	ns.Move(obj77, wp90)
	ns.Move(obj78, wp91)
	ns.Print("GeneralPrint:MsgTrapOn")
	ns.AudioEvent(ns.SpikeBlockMove, wp92)
	ns.AudioEvent(ns.BigGong, wp92)
LABEL1:
	return
}
func ResetBlockTrap() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj74, wp82)
	ns.Move(obj75, wp83)
	ns.Move(obj76, wp84)
	ns.Move(obj77, wp85)
	ns.Move(obj78, wp86)
	ns.WallGroupOpen(gvar94)
	ns.WallGroupOpen(gvar96)
	ns.AudioEvent(ns.SpikeBlockMove, wp89)
	ns.Print("GeneralPrint:MsgTrapOff")
	ns.ObjectOn(obj81)
}
func OpenBlockPathWalls() {
	ns.WallGroupOpen(gvar94)
	EnableBlockTrap()
}
func BlockReturn() {
	ns.Move(obj74, wp82)
	ns.Move(obj75, wp83)
	ns.Move(obj76, wp84)
	ns.Move(obj77, wp85)
	ns.Move(obj78, wp86)
	ns.AudioEvent(ns.SpikeBlockMove, wp89)
	ns.ObjectOn(obj81)
}
func BlockTrapEnd() {
	if !ns.IsCaller(obj76) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOff(obj79)
	ns.FrameTimer(90, BlockReturn)
LABEL1:
	return
}
func OpenBlockTrapExit() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar98)
}
func BlockTrapSecretFound() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar97)
}
func RotatingSpikeLoop() {
	ns.Move(obj101, wp107)
	ns.Move(obj102, wp108)
	ns.Move(obj103, wp109)
	ns.Move(obj104, wp110)
	ns.Move(obj105, wp111)
	ns.Move(obj106, wp112)
	ns.FrameTimer(ivar100, RotatingSpikeLoop)
}
func InitializeButtonPuzzle() {
	gvar113 = ns.WallGroup("ButtonPuzzleWall")
	gvar114 = ns.WallGroup("ButtonPuzzleSecretWalls")
	obj101 = ns.Object("RotatingSpike1")
	obj102 = ns.Object("RotatingSpike2")
	obj103 = ns.Object("RotatingSpike3")
	obj104 = ns.Object("RotatingSpike4")
	obj105 = ns.Object("RotatingSpike5")
	obj106 = ns.Object("RotatingSpike6")
	wp107 = ns.Waypoint("RotatingSpikeWP1")
	wp108 = ns.Waypoint("RotatingSpikeWP2")
	wp109 = ns.Waypoint("RotatingSpikeWP3")
	wp110 = ns.Waypoint("RotatingSpikeWP4")
	wp111 = ns.Waypoint("RotatingSpikeWP5")
	wp112 = ns.Waypoint("RotatingSpikeWP6")
}
func ButtonPuzzleButtonPress() {
	ns.ObjectOff(ns.GetTrigger())
	ivar99 += 1
	if ivar99 != 2 {
		goto LABEL1
	}
	ns.Print("GeneralPrint:MsgSeWallOpen")
	ns.WallGroupOpen(gvar113)
LABEL1:
	return
}
func EnableRotatingSpikeSet1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj101, wp107)
	ns.Move(obj102, wp108)
	ns.Move(obj103, wp109)
	ns.Move(obj104, wp110)
	ns.Move(obj105, wp111)
	ns.Move(obj106, wp112)
	ns.FrameTimer(ivar100, RotatingSpikeLoop)
}
func OpenButtonPuzzleSecretWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar114)
}
func InitializeFistTraps() {
	obj116[0] = ns.Object("FistTrap01")
	obj116[1] = ns.Object("FistTrap02")
	obj116[2] = ns.Object("FistTrap03")
	obj116[3] = ns.Object("FistTrap04")
	obj116[4] = ns.Object("FistTrap05")
	obj116[5] = ns.Object("FistTrap06")
	obj116[6] = ns.Object("FistTrap07")
	obj116[7] = ns.Object("FistTrap08")
	obj116[8] = ns.Object("FistTrap09")
	obj116[9] = ns.Object("FistTrap10")
	obj116[10] = ns.Object("FistTrap11")
	obj116[11] = ns.Object("FistTrap12")
	obj116[12] = ns.Object("FistTrap13")
	obj115[0] = ns.Object("FistTrapLight01")
	obj115[1] = ns.Object("FistTrapLight02")
	obj115[2] = ns.Object("FistTrapLight03")
	obj115[3] = ns.Object("FistTrapLight04")
	obj115[4] = ns.Object("FistTrapLight05")
	obj115[5] = ns.Object("FistTrapLight06")
	obj115[6] = ns.Object("FistTrapLight07")
	obj115[7] = ns.Object("FistTrapLight08")
	obj115[8] = ns.Object("FistTrapLight09")
	obj115[9] = ns.Object("FistTrapLight10")
	obj115[10] = ns.Object("FistTrapLight11")
	obj115[11] = ns.Object("FistTrapLight12")
	obj115[12] = ns.Object("FistTrapLight13")
}
func ResetFistTrap(a1 int) {
	ns.ObjectOn(obj115[a1])
	ns.ObjectOn(obj116[a1])
	ns.MoveWaypoint(wp118, ns.GetObjectX(obj116[a1]), ns.GetObjectY(obj116[a1]))
	ns.AudioEvent(ns.TriggerReleased, wp118)
}
func DropFist(a1 int) {
	ns.CastSpellObjectLocation(ns.SPELL_FIST, obj116[a1], ns.GetObjectX(obj116[a1]), ns.GetObjectY(obj116[a1]))
	ns.ObjectOff(obj116[a1])
	ns.ObjectOff(obj115[a1])
	ns.FrameTimerWithArg(ivar117, a1, ResetFistTrap)
}
func ActivateFistTrap01() {
	DropFist(0)
}
func ActivateFistTrap02() {
	DropFist(1)
}
func ActivateFistTrap03() {
	DropFist(2)
}
func ActivateFistTrap04() {
	DropFist(3)
}
func ActivateFistTrap05() {
	DropFist(4)
}
func ActivateFistTrap06() {
	DropFist(5)
}
func ActivateFistTrap07() {
	DropFist(6)
}
func ActivateFistTrap08() {
	DropFist(7)
}
func ActivateFistTrap09() {
	DropFist(8)
}
func ActivateFistTrap10() {
	DropFist(9)
}
func ActivateFistTrap11() {
	DropFist(10)
}
func ActivateFistTrap12() {
	DropFist(11)
}
func ActivateFistTrap13() {
	DropFist(12)
}
func InitializePeekABooTrap() {
	gvar126 = ns.WallGroup("Peek-A-BooWalls")
}
func InitializeWallPuzzle() {
	gvar132 = ns.WallGroup("WallPuzzleInnerWalls")
	gvar133 = ns.WallGroup("WallPuzzleOuterWalls")
	gvar134 = ns.WallGroup("WallPuzzleRightWalls")
	gvar135 = ns.WallGroup("WallPuzzleLeftWalls")
}
func InitializeTogglePuzzle() {
	var v0 int
	gvar130 = ns.WallGroup("TogglePuzzleRightWalls")
	gvar131 = ns.WallGroup("TogglePuzzleLeftWalls")
	v0 = ns.Random(0, 1000)
	if !(v0 < 250 || v0 > 750) {
		goto LABEL1
	}
	gvar129 = gvar127
	goto LABEL2
LABEL1:
	gvar129 = gvar128
LABEL2:
	return
}
func InitializeMiscTraps() {
	gvar125 = ns.ObjectGroup("SkullGunGroup01")
	obj120 = ns.Object("RotatingSpike7")
	obj121 = ns.Object("RotatingSpike8")
	wp122 = ns.Waypoint("RotatingSpike7WP")
	wp123 = ns.Waypoint("RotatingSpike8WP")
	wp124 = ns.Waypoint("RotSpikeSet2AudioOrigin")
}
func MapInitialize() {
	InitializeButtonPuzzle()
	InitializeAirlock()
	InitializeArrowTraps()
	InitializePeekABooTrap()
	InitializeAmbush()
	InitializeBlockTrap()
	InitializeWallPuzzle()
	InitializeBlockCrusher()
	InitializeTogglePuzzle()
	InitializeFistTraps()
	InitializeMiscTraps()
	wp119 = ns.Waypoint("WellWP")
}
func WellOfRestoration() {
	ns.RestoreHealth(ns.GetCaller(), ns.MaxHealth(ns.GetCaller())-ns.CurrentHealth(ns.GetCaller()))
	ns.AudioEvent(ns.RestoreHealthName, wp119)
	ns.Print("GeneralPrint:WellSignRefresh")
}
func DisableSkullGunGroup01() {
	ns.ObjectGroupOff(gvar125)
}
func EnableSkullGunGroup01() {
	ns.ObjectGroupOn(gvar125)
}
func OpenLordTrapWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(ns.Wall(109, 53))
}
func OpenRubyRoomExit() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(ns.WallGroup("RubyRoomWalls"))
}
func OpenSilverKeyShortcutWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(ns.WallGroup("SilverKeyShortcutWalls"))
}
func OpenCoffinRoomWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(ns.WallGroup("CoffinRoomWalls"))
}
func OpenSecretShortcutWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(ns.WallGroup("SecretShortcutWalls"))
}
func EnableRotatingSpikeSet2() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.Gear3, wp124)
	ns.Move(obj120, wp122)
	ns.Move(obj121, wp123)
}
func OpenPeekABooWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar126)
}
func ClosePeekABooWalls() {
	ns.WallGroupClose(gvar126)
}
func ToggleWalls() {
	if gvar129 != gvar128 {
		goto LABEL1
	}
	gvar129 = gvar127
	ns.WallGroupOpen(gvar130)
	ns.WallGroupClose(gvar131)
	goto LABEL2
LABEL1:
	gvar129 = gvar128
	ns.WallGroupOpen(gvar131)
	ns.WallGroupClose(gvar130)
LABEL2:
	return
}
func LowerInnerWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar132)
	ns.WallGroupOpen(gvar134)
}
func LowerOuterWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar133)
	ns.WallGroupOpen(gvar135)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
