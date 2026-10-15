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
	obj99   ns.ObjectID
	obj100  ns.ObjectID
	obj101  ns.ObjectID
	obj102  ns.ObjectID
	obj103  ns.ObjectID
	obj104  ns.ObjectID
	ivar105 int
	obj106  ns.ObjectID
	obj107  ns.ObjectID
	obj108  ns.ObjectID
	obj109  ns.ObjectID
	obj110  ns.ObjectID
	obj111  ns.ObjectID
	wp112   ns.WaypointID
	wp113   ns.WaypointID
	wp114   ns.WaypointID
	wp115   ns.WaypointID
	wp116   ns.WaypointID
	wp117   ns.WaypointID
	wp118   ns.WaypointID
	wp119   ns.WaypointID
	wp120   ns.WaypointID
	wp121   ns.WaypointID
	gvar122 ns.ObjectGroupID
	gvar123 ns.ObjectGroupID
	gvar124 ns.WallGroupID
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
	ivar105 = 90
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
	obj106 = ns.Object("Golem")
	obj107 = ns.Object("GolemUrchin")
	obj108 = ns.Object("SpikeBlock03")
	obj109 = ns.Object("SpikeBlock04")
	obj110 = ns.Object("SpikeBlock05")
	obj111 = ns.Object("SpikeBlock06")
	wp112 = ns.Waypoint("SpikeBlock03Home")
	wp113 = ns.Waypoint("SpikeBlock04Home")
	wp114 = ns.Waypoint("SpikeBlock05Home")
	wp115 = ns.Waypoint("SpikeBlock06Home")
	wp116 = ns.Waypoint("SpikeBlock03Dest")
	wp117 = ns.Waypoint("SpikeBlock04Dest")
	wp118 = ns.Waypoint("SpikeBlock05Dest")
	wp119 = ns.Waypoint("SpikeBlock06Dest")
	wp120 = ns.Waypoint("GolemButton")
	wp121 = ns.Waypoint("UrchinAudioOrigin")
	gvar122 = ns.ObjectGroup("BlockPusherPlates")
	gvar123 = ns.ObjectGroup("BlockPusherTriggers")
	gvar124 = ns.WallGroup("GolemWalls")
}
func InitializeGeneratorGrids() {
	obj99 = ns.Object("GenGrid01Generator04")
	obj100 = ns.Object("GenGrid01Generator05")
	obj101 = ns.Object("GenGrid01Generator06")
	obj102 = ns.Object("GenGrid01Generator07")
	obj103 = ns.Object("GenGrid01Generator08")
	obj104 = ns.Object("GenGrid01Generator09")
	ns.ObjectOff(obj99)
	ns.ObjectOff(obj100)
	ns.ObjectOff(obj101)
	ns.ObjectOff(obj102)
	ns.ObjectOff(obj103)
	ns.ObjectOff(obj104)
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
func GenGrid01Gen01Destroy() {
	ns.ObjectOn(obj99)
}
func GenGrid01Gen02Destroy() {
	ns.ObjectOn(obj100)
}
func GenGrid01Gen03Destroy() {
	ns.ObjectOn(obj101)
}
func GenGrid01Gen04Destroy() {
	ns.ObjectOn(obj100)
	ns.ObjectOn(obj102)
}
func GenGrid01Gen05Destroy() {
	ns.ObjectOn(obj99)
	ns.ObjectOn(obj101)
	ns.ObjectOn(obj103)
}
func GenGrid01Gen06Destroy() {
	ns.ObjectOn(obj100)
	ns.ObjectOn(obj104)
}
func GenGrid01Gen07Destroy() {
	ns.ObjectOn(obj99)
	ns.ObjectOn(obj103)
}
func GenGrid01Gen08Destroy() {
	ns.ObjectOn(obj100)
	ns.ObjectOn(obj104)
}
func GenGrid01Gen09Destroy() {
	ns.ObjectOn(obj101)
	ns.ObjectOn(obj103)
}
func ReleaseGolem() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar124)
	ns.ObjectOn(obj106)
	if !ns.IsCaller(obj107) {
		goto LABEL1
	}
	ns.AudioEvent(ns.UrchinTaunt, wp121)
LABEL1:
	if !(ns.CurrentHealth(obj107) > 0) {
		goto LABEL2
	}
	ns.CreatureIdle(obj107)
	ns.AggressionLevel(obj107, 0.83)
LABEL2:
	return
}
func GolemUrchinEnter() {
	ns.Walk(obj107, ns.GetWaypointX(wp120)+7, ns.GetWaypointY(wp120)-7)
}
func EnterGolemRoom() {
	ns.ObjectOff(ns.GetTrigger())
	ns.FrameTimer(15, GolemUrchinEnter)
}
func EnableBlockPusher() {
	ns.ObjectGroupOn(gvar122)
	ns.ObjectGroupOn(gvar123)
}
func ResetBlockPusher() {
	if !(ns.IsCaller(obj108) || ns.IsCaller(obj109)) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj108, wp112)
	ns.Move(obj109, wp113)
	ns.Move(obj110, wp114)
	ns.Move(obj111, wp115)
	ns.AudioEvent(ns.SpikeBlockMove, wp116)
	ns.AudioEvent(ns.SpikeBlockMove, wp118)
	ns.FrameTimer(ivar105, EnableBlockPusher)
LABEL1:
	return
}
func ActivateBlockPusherSeg2() {
	if !(ns.IsCaller(obj110) || ns.IsCaller(obj111)) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj108, wp116)
	ns.Move(obj109, wp117)
	ns.AudioEvent(ns.SpikeBlockMove, wp112)
LABEL1:
	return
}
func ActivateBlockPusher() {
	ns.ObjectGroupOff(gvar122)
	ns.Move(obj110, wp118)
	ns.Move(obj111, wp119)
	ns.AudioEvent(ns.SpikeBlockMove, wp114)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
