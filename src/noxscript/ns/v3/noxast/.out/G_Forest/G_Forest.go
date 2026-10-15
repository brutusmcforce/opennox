package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	obj4    ns.ObjectID
	obj5    ns.ObjectID
	obj6    ns.ObjectID
	obj7    ns.ObjectID
	obj8    ns.ObjectID
	obj9    ns.ObjectID
	gvar10  ns.ObjectGroupID
	gvar11  ns.ObjectGroupID
	wp12    ns.WaypointID
	wp13    ns.WaypointID
	fvar14  float32
	fvar15  float32
	fvar16  float32
	fvar17  float32
	flag18  bool
	flag19  bool
	wp20    ns.WaypointID
	obj21   ns.ObjectID
	gvar22  ns.ObjectGroupID
	gvar23  ns.ObjectGroupID
	gvar24  ns.WallGroupID
	gvar25  ns.WallGroupID
	obj26   ns.ObjectID
	obj27   ns.ObjectID
	obj28   ns.ObjectID
	obj29   ns.ObjectID
	obj30   ns.ObjectID
	obj31   ns.ObjectID
	obj32   ns.ObjectID
	obj33   ns.ObjectID
	obj34   ns.ObjectID
	obj35   ns.ObjectID
	obj36   ns.ObjectID
	obj37   ns.ObjectID
	obj38   ns.ObjectID
	obj39   ns.ObjectID
	gvar40  ns.ObjectGroupID
	wp41    ns.WaypointID
	wp42    ns.WaypointID
	wp43    ns.WaypointID
	wp44    ns.WaypointID
	wp45    ns.WaypointID
	wp46    ns.WaypointID
	wp47    ns.WaypointID
	wp48    ns.WaypointID
	wp49    ns.WaypointID
	wp50    ns.WaypointID
	wp51    ns.WaypointID
	wp52    ns.WaypointID
	wp53    ns.WaypointID
	wp54    ns.WaypointID
	gvar55  ns.WallGroupID
	gvar56  ns.ObjectGroupID
	gvar57  ns.ObjectGroupID
	gvar58  ns.ObjectGroupID
	gvar59  ns.ObjectGroupID
	obj60   ns.ObjectID
	obj61   ns.ObjectID
	wp62    ns.WaypointID
	wp63    ns.WaypointID
	wp64    ns.WaypointID
	wp65    ns.WaypointID
	ivar66  int
	ivar67  int
	obj68   ns.ObjectID
	obj69   ns.ObjectID
	obj70   ns.ObjectID
	obj71   ns.ObjectID
	obj72   ns.ObjectID
	obj73   ns.ObjectID
	obj74   ns.ObjectID
	obj75   ns.ObjectID
	obj76   ns.ObjectID
	gvar77  ns.ObjectGroupID
	gvar78  ns.ObjectGroupID
	wp79    ns.WaypointID
	wp80    ns.WaypointID
	wp81    ns.WaypointID
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
	wp93    ns.WaypointID
	wp94    ns.WaypointID
	gvar95  ns.WallGroupID
	gvar96  ns.WallGroupID
	gvar97  ns.WallGroupID
	gvar98  ns.WallGroupID
	ivar99  int
	obj100  ns.ObjectID
	obj101  ns.ObjectID
	obj102  ns.ObjectID
	obj103  ns.ObjectID
	obj104  ns.ObjectID
	obj105  ns.ObjectID
	wp106   ns.WaypointID
	wp107   ns.WaypointID
	wp108   ns.WaypointID
	wp109   ns.WaypointID
	wp110   ns.WaypointID
	wp111   ns.WaypointID
	wp112   ns.WaypointID
	wp113   ns.WaypointID
	wp114   ns.WaypointID
	wp115   ns.WaypointID
	gvar116 ns.ObjectGroupID
	gvar117 ns.ObjectGroupID
	gvar118 ns.WallGroupID
	obj119  ns.ObjectID
	obj120  ns.ObjectID
	obj121  ns.ObjectID
	obj122  ns.ObjectID
	obj123  ns.ObjectID
	obj124  ns.ObjectID
)

func init() {
	fvar14 = 0
	fvar15 = 3
	fvar16 = 10
	fvar17 = 15
	flag18 = false
	flag19 = false
	ivar66 = 0
	ivar67 = 0
	ivar99 = 90
}
func InitializeElevators() {
	obj4 = ns.Object("IsleElevator01")
	obj5 = ns.Object("IsleElevator02")
	obj6 = ns.Object("IsleElevator03")
	obj7 = ns.Object("IsleElevator04")
	obj8 = ns.Object("BanditElevator01")
	obj9 = ns.Object("BanditElevator02")
	gvar10 = ns.ObjectGroup("GearGroup01")
	gvar11 = ns.ObjectGroup("GearGroup02")
	wp12 = ns.Waypoint("Elevator01AudioOrigin")
	wp13 = ns.Waypoint("Elevator02AudioOrigin")
}
func ActivateBanditElevator01() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectGroupOn(gvar10)
	ns.ObjectOn(obj8)
	ns.AudioEvent(ns.CreatureCageAppears, wp12)
	ns.AudioEvent(ns.ChangeSpellbar, wp12)
	ns.Print("GeneralPrint:MsgElevOn")
}
func ActivateBanditElevator02() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectGroupOn(gvar11)
	ns.ObjectOn(obj9)
	ns.AudioEvent(ns.CreatureCageAppears, wp13)
	ns.AudioEvent(ns.ChangeSpellbar, wp13)
}
func ActivateIsleElevator01() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj4)
}
func ActivateIsleElevator02() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj5)
}
func ActivateIsleElevator03() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj6)
}
func ActivateIsleElevator04() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj7)
}
func OpenSecretShortcut() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(ns.Wall(37, 131))
}
func FireballTrapTickerLoop() {
	if !flag18 {
		goto LABEL1
	}
	fvar14 += 0.5
	ns.AudioEvent(ns.GenerateTick, wp20)
	if fvar14 != fvar15 {
		goto LABEL2
	}
	ns.WallGroupOpen(gvar25)
LABEL2:
	if !(fvar14 >= fvar16) {
		goto LABEL3
	}
	ns.ObjectGroupOn(gvar22)
LABEL3:
	if !(fvar14 >= fvar17) {
		goto LABEL4
	}
	ResetFireballTrap()
	goto LABEL1
LABEL4:
	ns.FrameTimer(15, FireballTrapTickerLoop)
LABEL1:
	return
}
func ResetFireballTrap() {
	flag18 = false
	fvar14 = 0
	ns.ObjectGroupOff(gvar22)
	ns.WallGroupOpen(gvar24)
	ns.ObjectGroupOn(gvar23)
	ns.ObjectOff(obj21)
	ns.Print("GeneralPrint:MsgTrapOff")
}
func InitializeFireballTrap() {
	wp20 = ns.Waypoint("FireballTrapAudioOrigin")
	obj21 = ns.Object("TickerGear")
	gvar22 = ns.ObjectGroup("FireballTrapSkullGuns")
	gvar23 = ns.ObjectGroup("FireballTrapPressurePlates")
	gvar24 = ns.WallGroup("FireballTrapEntranceWalls")
	gvar25 = ns.WallGroup("FireballTrapSurroundingWalls")
}
func ActivateFireballTrap() {
	ns.ObjectGroupOff(gvar23)
	flag18 = true
	ns.WallGroupClose(gvar24)
	ns.ObjectOn(obj21)
	FireballTrapTickerLoop()
}
func RollingBoulderAudioLoop() {
	ns.AudioEvent(ns.BoulderMove, wp44)
	ns.FrameTimer(30, RollingBoulderAudioLoop)
}
func InitializeForest() {
	obj29 = ns.Object("CavePit")
	obj26 = ns.Object("Boulder01")
	obj28 = ns.Object("BanditCampExitSkullGun")
	obj27 = ns.Object("BanditCampExitBlock")
	obj30 = ns.Object("RotatingSpike2")
	obj31 = ns.Object("RotatingSpike3")
	obj32 = ns.Object("RotatingSpike6")
	obj33 = ns.Object("RotatingSpike7")
	obj35 = ns.Object("LeftSentry")
	obj34 = ns.Object("RightSentry")
	obj36 = ns.Object("SecretBookcase01")
	obj37 = ns.Object("SecretBookcase02")
	obj38 = ns.Object("Candleabra01Mover")
	obj39 = ns.Object("Candleabra02Mover")
	gvar40 = ns.ObjectGroup("GearGroup03")
	wp41 = ns.Waypoint("WellWP")
	wp42 = ns.Waypoint("BlockGateAudioOrigin")
	wp43 = ns.Waypoint("CavePitAudioOrigin")
	wp44 = ns.Waypoint("RollingBoulderAudioOrigin")
	wp49 = ns.Waypoint("Boulder01WP")
	wp50 = ns.Waypoint("BanditCampExitBlockWP")
	wp45 = ns.Waypoint("RotatingSpikeWP2")
	wp46 = ns.Waypoint("RotatingSpikeWP3")
	wp47 = ns.Waypoint("RotatingSpikeWP6")
	wp48 = ns.Waypoint("RotatingSpikeWP7")
	wp51 = ns.Waypoint("LeftSentryWP")
	wp52 = ns.Waypoint("RightSentryWP")
	wp53 = ns.Waypoint("SecretBookcase01WP")
	wp54 = ns.Waypoint("SecretBookcase02WP")
	gvar55 = ns.WallGroup("BanditCampExitWalls")
	ns.Move(obj26, wp49)
	ns.Move(obj30, wp45)
	ns.Move(obj31, wp46)
	ns.Move(obj32, wp47)
	ns.Move(obj33, wp48)
	ns.Move(obj35, wp51)
	ns.Move(obj34, wp52)
	ns.ObjectOn(obj38)
	ns.ObjectOn(obj39)
	RollingBoulderAudioLoop()
}
func WellOfRestoration() {
	ns.RestoreHealth(ns.GetCaller(), ns.MaxHealth(ns.GetCaller())-ns.CurrentHealth(ns.GetCaller()))
	ns.AudioEvent(ns.RestoreHealthName, wp41)
	ns.Print("GeneralPrint:WellSignRefresh")
}
func OpenBanditCampExitWalls() {
	ns.WallGroupOpen(gvar55)
}
func CloseBanditCampExitWalls() {
	ns.WallGroupClose(gvar55)
}
func OpenBanditCampExitBlock() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj27, wp50)
	ns.ObjectGroupOn(gvar40)
	ns.AudioEvent(ns.SpikeBlockMove, wp42)
}
func StartSkullGun() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj28)
}
func DisableBanditCampExitSkullGun() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOff(obj28)
}
func CoopBridgeLeverPress() {
	ns.WallOpen(ns.Wall(121, 133))
}
func CoopBridgeLeverRelease() {
	ns.WallClose(ns.Wall(121, 133))
}
func OpenCavePit() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.ChangeSpellbar, wp43)
	ns.ObjectOn(obj29)
}
func OpenSecretBookcase01() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(ns.Wall(117, 69))
	ns.Move(obj36, wp53)
}
func OpenSecretBookcase02() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(ns.Wall(125, 219))
	ns.Move(obj37, wp54)
}
func InitializeTreeHaus() {
	obj60 = ns.Object("SpikeBlock01")
	obj61 = ns.Object("SpikeBlock02")
	wp62 = ns.Waypoint("SpikeBlock01Dest")
	wp63 = ns.Waypoint("SpikeBlock02Dest")
	wp64 = ns.Waypoint("SpikeBlock01Home")
	wp65 = ns.Waypoint("SpikeBlock02Home")
	gvar56 = ns.ObjectGroup("ArrowTrapGroup01")
	gvar57 = ns.ObjectGroup("ArrowTrapGroup02")
	gvar58 = ns.ObjectGroup("ArrowTrapGroup03")
	gvar59 = ns.ObjectGroup("BlockTrapPressurePlates")
}
func InitializeUnderGround() {
	obj69 = ns.Object("RotatingSpike11")
	obj70 = ns.Object("RotatingSpike12")
	obj71 = ns.Object("RotatingSpike13")
	obj72 = ns.Object("RotatingSpike14")
	obj73 = ns.Object("RotatingSpike15")
	obj74 = ns.Object("RotatingSpike16")
	obj75 = ns.Object("RotatingSpike17")
	obj76 = ns.Object("RotatingSpike18")
	obj68 = ns.Object("BoulderHallSkullGun")
	gvar77 = ns.ObjectGroup("RotatingSpikeTriggers")
	gvar78 = ns.ObjectGroup("CaveSkullGuns")
	wp79 = ns.Waypoint("RotatingSpikeHome11")
	wp80 = ns.Waypoint("RotatingSpikeHome12")
	wp81 = ns.Waypoint("RotatingSpikeHome13")
	wp82 = ns.Waypoint("RotatingSpikeHome14")
	wp83 = ns.Waypoint("RotatingSpikeHome15")
	wp84 = ns.Waypoint("RotatingSpikeHome16")
	wp85 = ns.Waypoint("RotatingSpikeHome17")
	wp86 = ns.Waypoint("RotatingSpikeHome18")
	wp87 = ns.Waypoint("RotatingSpikeDest11")
	wp88 = ns.Waypoint("RotatingSpikeDest12")
	wp89 = ns.Waypoint("RotatingSpikeDest13")
	wp90 = ns.Waypoint("RotatingSpikeDest14")
	wp91 = ns.Waypoint("RotatingSpikeDest15")
	wp92 = ns.Waypoint("RotatingSpikeDest16")
	wp93 = ns.Waypoint("RotatingSpikeDest17")
	wp94 = ns.Waypoint("RotatingSpikeDest18")
	gvar95 = ns.WallGroup("ArcherRunWalls")
	gvar96 = ns.WallGroup("ArcherRunExitWalls")
	gvar97 = ns.WallGroup("BoulderHallWalls")
	gvar98 = ns.WallGroup("BanditTreasureRoomWalls")
	ns.WallGroupClose(gvar98)
}
func InitializeGolem() {
	obj100 = ns.Object("Golem")
	obj101 = ns.Object("GolemUrchin")
	obj102 = ns.Object("SpikeBlock03")
	obj103 = ns.Object("SpikeBlock04")
	obj104 = ns.Object("SpikeBlock05")
	obj105 = ns.Object("SpikeBlock06")
	wp106 = ns.Waypoint("SpikeBlock03Home")
	wp107 = ns.Waypoint("SpikeBlock04Home")
	wp108 = ns.Waypoint("SpikeBlock05Home")
	wp109 = ns.Waypoint("SpikeBlock06Home")
	wp110 = ns.Waypoint("SpikeBlock03Dest")
	wp111 = ns.Waypoint("SpikeBlock04Dest")
	wp112 = ns.Waypoint("SpikeBlock05Dest")
	wp113 = ns.Waypoint("SpikeBlock06Dest")
	wp114 = ns.Waypoint("GolemButton")
	wp115 = ns.Waypoint("UrchinAudioOrigin")
	gvar116 = ns.ObjectGroup("BlockPusherPlates")
	gvar117 = ns.ObjectGroup("BlockPusherTriggers")
	gvar118 = ns.WallGroup("GolemWalls")
}
func InitializeGeneratorGrids() {
	obj119 = ns.Object("GenGrid01Generator04")
	obj120 = ns.Object("GenGrid01Generator05")
	obj121 = ns.Object("GenGrid01Generator06")
	obj122 = ns.Object("GenGrid01Generator07")
	obj123 = ns.Object("GenGrid01Generator08")
	obj124 = ns.Object("GenGrid01Generator09")
	ns.ObjectOff(obj119)
	ns.ObjectOff(obj120)
	ns.ObjectOff(obj121)
	ns.ObjectOff(obj122)
	ns.ObjectOff(obj123)
	ns.ObjectOff(obj124)
}
func MapInitialize() {
	InitializeElevators()
	InitializeForest()
	InitializeFireballTrap()
	InitializeTreeHaus()
	InitializeUnderGround()
	InitializeGolem()
	InitializeGeneratorGrids()
}
func DisableArrowTrap01() {
	ns.ObjectGroupOff(gvar56)
}
func ArrowTrap01() {
	ns.ObjectGroupOn(gvar56)
	ns.FrameTimer(1, DisableArrowTrap01)
}
func DisableArrowTrap02() {
	ns.ObjectGroupOff(gvar57)
}
func ArrowTrap02() {
	ns.ObjectGroupOn(gvar57)
	ns.FrameTimer(1, DisableArrowTrap02)
}
func DisableArrowTrap03() {
	ns.ObjectGroupOff(gvar58)
}
func ArrowTrap03() {
	ns.ObjectGroupOn(gvar58)
	ns.FrameTimer(1, DisableArrowTrap03)
}
func ResetBlockTrap() {
	ns.ObjectGroupOn(gvar59)
}
func SpikeBlocksReturn() {
	ns.Move(obj60, wp64)
	ns.Move(obj61, wp65)
	ns.FrameTimer(80, ResetBlockTrap)
}
func ActivateSpikeBlockTrap() {
	ns.ObjectGroupOff(gvar59)
	ns.Move(obj60, wp62)
	ns.Move(obj61, wp63)
	ns.FrameTimer(110, SpikeBlocksReturn)
}
func ActivateBoulderHallSkullGun() {
	ns.ObjectOn(obj68)
}
func DisableBoulderHallSkullGun() {
	ns.ObjectOff(obj68)
}
func OpenBoulderHallWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar97)
}
func ArcherRunButtonPress() {
	ns.ObjectOff(ns.GetTrigger())
	ivar66 += 1
	if !(ivar66 >= 2) {
		goto LABEL1
	}
	ns.WallGroupOpen(gvar95)
	ns.Print("GeneralPrint:MsgSeWallOpen")
LABEL1:
	return
}
func OpenArcherRunExitWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar96)
}
func CheckForBothButtons() {
	if !(ivar67 >= 2) {
		goto LABEL1
	}
	ns.WallGroupOpen(gvar98)
LABEL1:
	return
}
func BanditTreasureRoomPress() {
	ivar67 += 1
	ns.FrameTimer(1, CheckForBothButtons)
}
func BanditTreasureRoomRelease() {
	ivar67 -= 1
	ns.WallGroupClose(gvar98)
}
func EnableRotatingSpikeTriggers() {
	ns.ObjectGroupOn(gvar77)
}
func ResetRotatingSpikes() {
	ns.Move(obj69, wp79)
	ns.Move(obj70, wp80)
	ns.Move(obj71, wp81)
	ns.Move(obj72, wp82)
	ns.Move(obj73, wp83)
	ns.Move(obj74, wp84)
	ns.Move(obj75, wp85)
	ns.Move(obj76, wp86)
	ns.FrameTimer(60, EnableRotatingSpikeTriggers)
}
func ActivateRotatingSpikes() {
	ns.ObjectGroupOff(gvar77)
	ns.Move(obj69, wp87)
	ns.Move(obj70, wp88)
	ns.Move(obj71, wp89)
	ns.Move(obj72, wp90)
	ns.Move(obj73, wp91)
	ns.Move(obj74, wp92)
	ns.Move(obj75, wp93)
	ns.Move(obj76, wp94)
	ns.FrameTimer(60, ResetRotatingSpikes)
}
func DisableCaveSkullGuns() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectGroupOff(gvar78)
	ns.Print("GeneralPrint:MsgTrapOff")
}
func ReleaseGolem() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar118)
	ns.ObjectOn(obj100)
	if !ns.IsCaller(obj101) {
		goto LABEL1
	}
	ns.AudioEvent(ns.UrchinTaunt, wp115)
LABEL1:
	if !(ns.CurrentHealth(obj101) > 0) {
		goto LABEL2
	}
	ns.CreatureIdle(obj101)
	ns.AggressionLevel(obj101, 0.83)
LABEL2:
	return
}
func GolemUrchinEnter() {
	ns.Walk(obj101, ns.GetWaypointX(wp114)+7, ns.GetWaypointY(wp114)-7)
}
func EnterGolemRoom() {
	ns.ObjectOff(ns.GetTrigger())
	ns.FrameTimer(15, GolemUrchinEnter)
}
func EnableBlockPusher() {
	ns.ObjectGroupOn(gvar116)
	ns.ObjectGroupOn(gvar117)
}
func ResetBlockPusher() {
	if !(ns.IsCaller(obj102) || ns.IsCaller(obj103)) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj102, wp106)
	ns.Move(obj103, wp107)
	ns.Move(obj104, wp108)
	ns.Move(obj105, wp109)
	ns.AudioEvent(ns.SpikeBlockMove, wp110)
	ns.AudioEvent(ns.SpikeBlockMove, wp112)
	ns.FrameTimer(ivar99, EnableBlockPusher)
LABEL1:
	return
}
func ActivateBlockPusherSeg2() {
	if !(ns.IsCaller(obj104) || ns.IsCaller(obj105)) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj102, wp110)
	ns.Move(obj103, wp111)
	ns.AudioEvent(ns.SpikeBlockMove, wp106)
LABEL1:
	return
}
func ActivateBlockPusher() {
	ns.ObjectGroupOff(gvar116)
	ns.Move(obj104, wp112)
	ns.Move(obj105, wp113)
	ns.AudioEvent(ns.SpikeBlockMove, wp108)
}
func GenGrid01Gen01Destroy() {
	ns.ObjectOn(obj119)
}
func GenGrid01Gen02Destroy() {
	ns.ObjectOn(obj120)
}
func GenGrid01Gen03Destroy() {
	ns.ObjectOn(obj121)
}
func GenGrid01Gen04Destroy() {
	ns.ObjectOn(obj120)
	ns.ObjectOn(obj122)
}
func GenGrid01Gen05Destroy() {
	ns.ObjectOn(obj119)
	ns.ObjectOn(obj121)
	ns.ObjectOn(obj123)
}
func GenGrid01Gen06Destroy() {
	ns.ObjectOn(obj120)
	ns.ObjectOn(obj124)
}
func GenGrid01Gen07Destroy() {
	ns.ObjectOn(obj119)
	ns.ObjectOn(obj123)
}
func GenGrid01Gen08Destroy() {
	ns.ObjectOn(obj120)
	ns.ObjectOn(obj124)
}
func GenGrid01Gen09Destroy() {
	ns.ObjectOn(obj121)
	ns.ObjectOn(obj123)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
