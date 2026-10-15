package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	flag4  bool
	wp5    [4]ns.WaypointID
	wp6    [20]ns.WaypointID
	wp7    ns.WaypointID
	gvar8  int
	gvar9  int
	obj10  ns.ObjectID
	obj11  ns.ObjectID
	obj12  ns.ObjectID
	obj13  ns.ObjectID
	gvar14 int
	gvar15 ns.ObjectGroupID
	flag16 bool
	gvar17 ns.ObjectGroupID
	gvar18 ns.ObjectGroupID
	gvar19 ns.ObjectGroupID
	gvar20 ns.WallGroupID
	wp21   ns.WaypointID
	wp22   ns.WaypointID
	wp23   ns.WaypointID
	wp24   ns.WaypointID
	wp25   ns.WaypointID
	wp26   ns.WaypointID
	obj27  ns.ObjectID
	ivar28 int
	obj29  [50]ns.ObjectID
	ivar30 int
	gvar31 int
	gvar32 ns.ObjectGroupID
	gvar33 ns.ObjectGroupID
	gvar34 ns.ObjectGroupID
	gvar35 ns.ObjectGroupID
	gvar36 ns.ObjectGroupID
	wp37   ns.WaypointID
	wp38   ns.WaypointID
	wp39   ns.WaypointID
	wp40   ns.WaypointID
	wp41   ns.WaypointID
	wp42   [16]ns.WaypointID
	wp43   [6]ns.WaypointID
	wp44   [4]ns.WaypointID
	wp45   [8]ns.WaypointID
	wp46   [4]ns.WaypointID
	wp47   [4]ns.WaypointID
	wp48   [4]ns.WaypointID
	wp49   [4]ns.WaypointID
	fvar50 float32
	fvar51 float32
	fvar52 float32
	fvar53 float32
	fvar54 float32
	gvar55 int
	ivar56 int
	ivar57 int
	ivar58 int
	gvar59 ns.TimerID
	flag60 bool
	flag61 bool
	flag62 bool
	flag63 bool
	flag64 bool
	flag65 bool
	flag66 bool
	obj67  ns.ObjectID
)

func init() {
	flag4 = false
	gvar8 = 0
	gvar9 = 1
	gvar14 = gvar8
	flag16 = false
	ivar28 = 50
	ivar30 = 0
	fvar50 = 7
	fvar51 = 75
	fvar52 = 150
	fvar53 = 30
	fvar54 = 70
	gvar55 = 5
	ivar56 = 600
	ivar57 = 5
	ivar58 = 2
	gvar59 = 0
	flag60 = false
	flag61 = true
	flag62 = false
	flag63 = false
	flag64 = false
	flag65 = false
	flag66 = false
}
func BarrelCaveIn() {
	var (
		v0 ns.ObjectGroupID
		v1 float32
		v2 float32
	)
	v1 = ns.GetObjectX(ns.GetHost())
	v2 = ns.GetObjectY(ns.GetHost())
	v0 = ns.ObjectGroup("BarrelExplosionTriggers")
	ns.Effect(ns.JIGGLE, v1, v2, 10, 0)
	ns.AudioEvent(ns.WallDestroyed, wp7)
	ns.AudioEvent(ns.FireballExplode, wp7)
	ns.CreateObject("CaveRocksLarge", wp5[0])
	ns.CreateObject("CaveRocksSmall", wp5[1])
	ns.CreateObject("CaveRocksSmall", wp5[2])
	ns.CreateObject("CaveRocksPebbles", wp5[3])
	ns.ObjectGroupOff(v0)
	ns.FrameTimer(20, BarrelsExplode)
}
func BarrelsExplode() {
	var v0 ns.ObjectID
	v0 = ns.Object("PowderBarrel")
	ns.Damage(v0, 0, 30, 1)
	ns.CreateObject("Flame", wp6[0])
	ns.CreateObject("SmallFlame", wp6[1])
	ns.CreateObject("MediumFlame", wp6[2])
	ns.CreateObject("Flame", wp6[3])
	ns.CreateObject("MediumFlame", wp6[4])
	ns.CreateObject("Flame", wp6[5])
	ns.CreateObject("MediumFlame", wp6[6])
	ns.CreateObject("Flame", wp6[7])
	ns.CreateObject("MediumFlame", wp6[8])
	ns.CreateObject("Flame", wp6[9])
	ns.CreateObject("SmallFlame", wp6[10])
	ns.CreateObject("MediumFlame", wp6[11])
	ns.CreateObject("Flame", wp6[12])
	ns.CreateObject("SmallFlame", wp6[13])
	ns.CreateObject("MediumFlame", wp6[14])
	ns.CreateObject("Flame", wp6[15])
	ns.CreateObject("SmallFlame", wp6[16])
	ns.CreateObject("MediumFlame", wp6[17])
	ns.CreateObject("Flame", wp6[18])
	ns.CreateObject("SmallFlame", wp6[19])
	ns.FrameTimer(30, EndExplodingBarrelSetPiece)
}
func EndExplodingBarrelSetPiece() {
	ns.WideScreen(false)
	ns.Frozen(ns.GetHost(), false)
}
func BeginExplodingBarrelSetPiece() {
	var (
		v0 float32
		v1 float32
	)
	if !(!flag4 && ns.IsCaller(ns.GetHost())) {
		goto LABEL1
	}
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.Effect(ns.JIGGLE, v0, v1, 10, 0)
	ns.AudioEvent(ns.EarthRumbleMajor, wp7)
	flag4 = true
	ns.Frozen(ns.GetHost(), true)
	ns.WideScreen(true)
	ns.FrameTimer(30, BarrelCaveIn)
LABEL1:
	return
}
func CaptainDialogStart() {
	var v0 int
	ns.LookAtObject(obj10, ns.GetHost())
	switch v0 = gvar14; v0 {
	case gvar8:
		goto LABEL1
	case gvar9:
		goto LABEL2
	default:
		goto LABEL3
	}
LABEL1:
	ns.TellStory(ns.SwordsmanHurt, "Con01a:CaptainGreet")
	goto LABEL3
LABEL2:
	ns.TellStory(ns.SwordsmanHurt, "Con01a:CaptainProd")
	goto LABEL3
LABEL3:
	return
}
func CaptainDialogEnd() {
	var v0 int
	switch v0 = gvar14; v0 {
	case gvar8:
		goto LABEL1
	case gvar9:
		goto LABEL2
	default:
		goto LABEL3
	}
LABEL1:
	ns.JournalEntry(ns.GetHost(), "LocateAldwyn", 2)
	ns.PrintToAll("Con01a:NewJournalEntry")
	gvar14 = gvar9
	goto LABEL3
LABEL2:
	goto LABEL3
LABEL3:
	return
}
func Gatekeeper2DialogStart() {
	ns.LookAtObject(ns.GetTrigger(), ns.GetHost())
	ns.TellStory(ns.SwordsmanHurt, "Con02a:Gatekeeper2Greet")
}
func Gatekeeper2DialogEnd() {
}
func Gatekeeper3DialogStart() {
	ns.LookAtObject(ns.GetTrigger(), ns.GetHost())
	ns.TellStory(ns.SwordsmanHurt, "Con02a:Gatekeeper3Greet")
}
func Gatekeeper3DialogEnd() {
}
func InitializeDialogs() {
	ns.StoryPic(obj10, "AirshipCaptainPic")
	ns.StoryPic(obj11, "IxGuard2Pic")
	ns.StoryPic(obj12, "IxGuard1Pic")
	ns.SetDialog(obj11, ns.NORMAL, Gatekeeper2DialogStart, Gatekeeper2DialogEnd)
	ns.SetDialog(obj12, ns.NORMAL, Gatekeeper3DialogStart, Gatekeeper3DialogEnd)
}
func StartCaptainConversation() {
	ns.SetDialog(obj10, ns.NORMAL, CaptainDialogStart, CaptainDialogEnd)
	ns.StartDialog(obj10, ns.GetHost())
}
func PlayerDeath() {
	ns.DeathScreen(1)
}
func MovePlayerToExit() {
	var (
		v0 float32
		v1 float32
		v2 ns.WaypointID
	)
	v2 = ns.Waypoint("Exit")
	v0 = ns.GetWaypointX(v2)
	v1 = ns.GetWaypointY(v2)
	ns.Frozen(ns.GetHost(), false)
	ns.MoveObject(ns.GetHost(), v0, v1)
}
func MapSwitch() {
	if !(!flag16 && ns.IsCaller(ns.GetHost())) {
		goto LABEL1
	}
	ns.Frozen(ns.GetHost(), true)
	flag16 = true
	ns.ObjectGroupOff(gvar15)
	ns.Blind()
	ns.FrameTimer(60, MovePlayerToExit)
LABEL1:
	return
}
func MineShake4() {
	var (
		v0 float32
		v1 float32
		v2 ns.ObjectGroupID
		v3 ns.WaypointID
	)
	v3 = ns.Waypoint("MineShake4AudioOrigin")
	v2 = ns.ObjectGroup("MineShake4Triggers")
	ns.AudioEvent(ns.EarthRumbleMajor, v3)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.ObjectGroupOff(v2)
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.CreateObject("CaveRocksTiny", wp48[0])
	ns.CreateObject("CaveRocksTiny", wp48[3])
	ns.FrameTimer(7, MineShake4SEG2)
}
func PlayEntranceLiftSounds() {
	ns.AudioEvent(ns.ChangeSpellbar, wp21)
	ns.AudioEvent(ns.ChangeSpellbar, wp22)
	ns.FrameTimer(7, PlayEntranceLiftSounds)
}
func PlayExitLiftSounds() {
	ns.AudioEvent(ns.ChangeSpellbar, wp24)
	ns.AudioEvent(ns.ChangeSpellbar, wp25)
	ns.AudioEvent(ns.ChangeSpellbar, wp26)
	ns.FrameTimer(7, PlayExitLiftSounds)
}
func ActivateEntranceElevator() {
	var v0 ns.ObjectID
	ns.GiveXp(ns.GetHost(), 100)
	v0 = ns.Object("EntranceSwitch")
	ns.ObjectOff(v0)
	ns.AudioEvent(ns.CreatureCageAppears, wp21)
	ns.AudioEvent(ns.ChangeSpellbar, wp21)
	ns.ObjectGroupOn(gvar17)
}
func ActivateExitElevator() {
	var v0 ns.ObjectID
	ns.GiveXp(ns.GetHost(), 100)
	v0 = ns.Object("ExitSwitch")
	ns.ObjectOff(v0)
	ns.AudioEvent(ns.CreatureCageAppears, wp24)
	ns.AudioEvent(ns.ChangeSpellbar, wp24)
	ns.ObjectGroupOn(gvar18)
	MineShake4()
}
func ActivateSecretElevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.CreatureCageAppears, wp23)
	ns.AudioEvent(ns.ChangeSpellbar, wp23)
	ns.WallGroupOpen(gvar20)
	ns.ObjectGroupOn(gvar19)
}
func PlayOutdoorMusic() {
	ns.Music(21, 100)
}
func DisableEntranceLift() {
	ns.ObjectGroupOff(gvar17)
}
func OwnObjects() {
	ns.SetOwner(ns.GetHost(), obj10)
	ns.SetOwner(ns.GetHost(), obj11)
	ns.SetOwner(ns.GetHost(), obj12)
}
func MapInitialize() {
	var v0 int
	v0 = 0
	gvar32 = ns.ObjectGroup("TunnelEntranceTriggers")
	gvar33 = ns.ObjectGroup("TunnelTriggers")
	gvar34 = ns.ObjectGroup("EntranceShakeTriggers")
	gvar17 = ns.ObjectGroup("EntranceLiftAndGears")
	gvar18 = ns.ObjectGroup("ExitLiftAndGears")
	gvar35 = ns.ObjectGroup("EntranceExitTriggers")
	gvar36 = ns.ObjectGroup("NPCs")
	gvar15 = ns.ObjectGroup("MapSwitchTriggers")
	gvar19 = ns.ObjectGroup("SecretElevatorGroup")
	gvar20 = ns.WallGroup("SecretElevatorWalls")
	obj11 = ns.Object("Ryan")
	obj12 = ns.Object("Ed")
	obj27 = ns.Object("PlayerStart")
	obj10 = ns.Object("Airship_Captain")
	obj13 = ns.Object("Staff")
	wp37 = ns.Waypoint("QuakeAudioOrigin")
	wp40 = ns.Waypoint("CreationWP")
	wp21 = ns.Waypoint("EntranceLiftAudioOrigin1")
	wp22 = ns.Waypoint("EntranceLiftAudioOrigin2")
	wp24 = ns.Waypoint("ExitLiftAudioOrigin1")
	wp25 = ns.Waypoint("ExitLiftAudioOrigin2")
	wp26 = ns.Waypoint("ExitLiftAudioOrigin3")
	wp23 = ns.Waypoint("SecretElevatorAudioOrigin")
	wp38 = ns.Waypoint("TunnelAudioOrigin")
	wp39 = ns.Waypoint("ExitAudioOrigin")
	wp7 = ns.Waypoint("BarrelAudioOrigin")
	wp41 = ns.Waypoint("SecretAudioWP")
	wp42[0] = ns.Waypoint("TunnelPitWP0")
	wp42[1] = ns.Waypoint("TunnelPitWP1")
	wp42[2] = ns.Waypoint("TunnelPitWP2")
	wp42[3] = ns.Waypoint("TunnelPitWP3")
	wp42[4] = ns.Waypoint("TunnelPitWP4")
	wp42[5] = ns.Waypoint("TunnelPitWP5")
	wp42[6] = ns.Waypoint("TunnelPitWP6")
	wp42[7] = ns.Waypoint("TunnelPitWP7")
	wp42[8] = ns.Waypoint("TunnelPitWP8")
	wp42[9] = ns.Waypoint("TunnelPitWP9")
	wp42[10] = ns.Waypoint("TunnelPitWP10")
	wp42[11] = ns.Waypoint("TunnelPitWP11")
	wp42[12] = ns.Waypoint("TunnelPitWP12")
	wp42[13] = ns.Waypoint("TunnelPitWP13")
	wp42[14] = ns.Waypoint("TunnelPitWP14")
	wp42[15] = ns.Waypoint("TunnelPitWP15")
	wp43[0] = ns.Waypoint("ExitPitWP0")
	wp43[1] = ns.Waypoint("ExitPitWP1")
	wp43[2] = ns.Waypoint("ExitPitWP2")
	wp43[3] = ns.Waypoint("ExitPitWP3")
	wp43[4] = ns.Waypoint("ExitPitWP4")
	wp43[5] = ns.Waypoint("ExitPitWP5")
	wp44[0] = ns.Waypoint("EntranceShakePit0")
	wp44[1] = ns.Waypoint("EntranceShakePit1")
	wp44[2] = ns.Waypoint("EntranceShakePit2")
	wp44[3] = ns.Waypoint("EntranceShakePit3")
	wp45[0] = ns.Waypoint("MineShake1Pit1")
	wp45[1] = ns.Waypoint("MineShake1Pit2")
	wp45[2] = ns.Waypoint("MineShake1Pit3")
	wp45[3] = ns.Waypoint("MineShake1Pit4")
	wp45[4] = ns.Waypoint("MineShake1Pit5")
	wp45[5] = ns.Waypoint("MineShake1Pit6")
	wp45[6] = ns.Waypoint("MineShake1Pit7")
	wp45[7] = ns.Waypoint("MineShake1Pit8")
	wp46[0] = ns.Waypoint("MineShake2Pit1")
	wp46[1] = ns.Waypoint("MineShake2Pit2")
	wp46[2] = ns.Waypoint("MineShake2Pit3")
	wp46[3] = ns.Waypoint("MineShake2Pit4")
	wp47[0] = ns.Waypoint("MineShake3Pit1")
	wp47[1] = ns.Waypoint("MineShake3Pit2")
	wp47[2] = ns.Waypoint("MineShake3Pit3")
	wp47[3] = ns.Waypoint("MineShake3Pit4")
	wp48[0] = ns.Waypoint("MineShake4Pit1")
	wp48[1] = ns.Waypoint("MineShake4Pit2")
	wp48[2] = ns.Waypoint("MineShake4Pit3")
	wp48[3] = ns.Waypoint("MineShake4Pit4")
	wp49[0] = ns.Waypoint("MineShake5Pit1")
	wp49[1] = ns.Waypoint("MineShake5Pit2")
	wp49[2] = ns.Waypoint("MineShake5Pit3")
	wp49[3] = ns.Waypoint("MineShake5Pit4")
	wp5[0] = ns.Waypoint("CavePit1")
	wp5[1] = ns.Waypoint("CavePit2")
	wp5[2] = ns.Waypoint("CavePit3")
	wp5[3] = ns.Waypoint("CavePit4")
	wp6[0] = ns.Waypoint("FlameWP1")
	wp6[1] = ns.Waypoint("FlameWP2")
	wp6[2] = ns.Waypoint("FlameWP3")
	wp6[3] = ns.Waypoint("FlameWP4")
	wp6[4] = ns.Waypoint("FlameWP5")
	wp6[5] = ns.Waypoint("FlameWP6")
	wp6[6] = ns.Waypoint("FlameWP7")
	wp6[7] = ns.Waypoint("FlameWP8")
	wp6[8] = ns.Waypoint("FlameWP9")
	wp6[9] = ns.Waypoint("FlameWP10")
	wp6[10] = ns.Waypoint("FlameWP11")
	wp6[11] = ns.Waypoint("FlameWP12")
	wp6[12] = ns.Waypoint("FlameWP13")
	wp6[13] = ns.Waypoint("FlameWP14")
	wp6[14] = ns.Waypoint("FlameWP15")
	wp6[15] = ns.Waypoint("FlameWP16")
	wp6[16] = ns.Waypoint("FlameWP17")
	wp6[17] = ns.Waypoint("FlameWP18")
	wp6[18] = ns.Waypoint("FlameWP19")
	wp6[19] = ns.Waypoint("FlameWP20")
	v0 = 0
	for {
		if !(v0 < 50) {
			goto LABEL1
		}
		obj29[v0] = 0
		v0 += 1
	}
LABEL1:
	ns.Pickup(obj10, obj13)
	InitializeDialogs()
	ns.StartupScreen(1)
	ns.FrameTimer(35, DisableEntranceLift)
	ns.FrameTimer(2, OwnObjects)
	ns.FrameTimer(5, PlayOutdoorMusic)
	ns.FrameTimer(5, StartCaptainConversation)
}
func MapEntry() {
	ns.ObjectGroupOn(gvar15)
}
func EnableTeleporter() {
	ns.ObjectOn(obj27)
}
func NoMonsters() {
	if !ns.IsAttackedBy(ns.GetCaller(), ns.GetHost()) {
		goto LABEL1
	}
	ns.GoBackHome(ns.GetCaller())
LABEL1:
	return
}
func PlayUndergroundMusic() {
	ns.Music(18, 100)
}
func Quake() {
	var (
		v0  float32
		v1  float32
		v2  float32
		v3  float32
		v4  float32
		v5  int
		v6  int
		v7  int
		v8  int
		v9  ns.ObjectID
		v10 int
	)
	if !flag60 {
		goto LABEL1
	}
	ns.PrintToAll("Quake initiated")
LABEL1:
	if flag61 {
		goto LABEL2
	}
	return
LABEL2:
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	v7 = ns.Random(ivar58, ivar57)
	ns.MoveWaypoint(wp37, v0, v1)
	ns.AudioEvent(ns.EarthRumbleMinor, wp37)
	ns.AudioEvent(ns.BoulderMove, wp37)
	ns.Effect(ns.JIGGLE, v0, v1, fvar50, 0)
	v8 = 0
	for {
		if !(v8 <= v7) {
			goto LABEL3
		}
		if !flag60 {
			goto LABEL4
		}
		ns.PrintToAll("''For'' loop")
	LABEL4:
		v2 = ns.RandomFloat(fvar53, fvar54)
		v3 = ns.RandomFloat(fvar53, fvar54)
		v5 = ns.Random(0, 3)
		if !flag60 {
			goto LABEL5
		}
		ns.PrintToAll("Attempting switch")
	LABEL5:
		switch v10 = v5; v10 {
		case 0:
			goto LABEL6
		case 1:
			goto LABEL7
		case 2:
			goto LABEL8
		case 3:
			goto LABEL9
		default:
			goto LABEL10
		}
	LABEL6:
		ns.MoveWaypoint(wp40, v0+v2, v1+v3)
		goto LABEL11
	LABEL7:
		ns.MoveWaypoint(wp40, v0+v2, v1-v3)
		goto LABEL11
	LABEL8:
		ns.MoveWaypoint(wp40, v0-v2, v1+v3)
		goto LABEL11
	LABEL9:
		ns.MoveWaypoint(wp40, v0-v2, v1-v3)
		goto LABEL11
	LABEL10:
		goto LABEL11
	LABEL11:
		ns.Delete(obj29[ivar30])
		v9 = ns.CreateObject("CaveRocksTiny", wp40)
		obj29[ivar30] = v9
		ivar30 += 1
		if !(ivar30 >= ivar28) {
			goto LABEL12
		}
		ivar30 = 0
	LABEL12:
		v4 = ns.RandomFloat(fvar51, fvar52)
		ns.Raise(v9, v4)
		v8 += 1
	}
LABEL3:
	v6 = ns.Random(5, ivar56)
	gvar59 = ns.FrameTimer(v6, Quake)
}
func DebugOn() {
	flag60 = true
}
func DebugOff() {
	flag60 = false
}
func TunnelCollapse() {
	var (
		v0 float32
		v1 float32
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.Effect(ns.JIGGLE, v0, v1, 10, 0)
	ns.AudioEvent(ns.GolemMove, wp38)
	ns.FrameTimer(4, FallingRocksSEG1)
}
func FallingRocksSEG1() {
	var (
		v0 float32
		v1 float32
	)
	_ = v1
	_ = v0
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.AudioEvent(ns.EarthRumbleMajor, wp38)
	ns.CreateObject("CaveRocksSmall", wp42[0])
	ns.CreateObject("CaveRocksSmall", wp42[7])
	ns.CreateObject("CaveRocksMedium", wp42[9])
	ns.CreateObject("CaveRocksSmall", wp42[14])
	ns.FrameTimer(10, FallingRocksSEG2)
}
func SetpieceShake1() {
	var (
		v0 float32
		v1 float32
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.Effect(ns.JIGGLE, v0, v1, 10, 0)
	ns.AudioEvent(ns.EarthRumbleMinor, wp38)
}
func NullCallback() {
}
func AwardExitExperience() {
	ns.MoveWaypoint(wp39, ns.GetObjectX(ns.GetHost()), ns.GetObjectY(ns.GetHost()))
	ns.AudioEvent(ns.FlagDrop, wp39)
	ns.GiveXp(ns.GetHost(), 200)
}
func CaveExitCollapseSEG2() {
	var (
		v0 ns.WaypointID
		v1 float32
		v2 float32
	)
	v1 = ns.GetObjectX(ns.GetHost())
	v2 = ns.GetObjectY(ns.GetHost())
	v0 = ns.Waypoint("ExitBoulderWP")
	ns.NoWallSound(false)
	ns.Effect(ns.JIGGLE, v1, v2, 10, 0)
	ns.AudioEvent(ns.FistHit, wp39)
	ns.CreateObject("BoulderCave", v0)
}
func CaveExitCollapse() {
	var (
		v0 ns.WallGroupID
		v1 ns.ObjectGroupID
		v2 float32
		v3 float32
	)
	v0 = ns.WallGroup("CaveExitWalls")
	v1 = ns.ObjectGroup("CaveExitTriggers")
	v2 = ns.GetObjectX(ns.GetHost())
	v3 = ns.GetObjectY(ns.GetHost())
	ns.ObjectGroupOff(v1)
	flag63 = true
	ns.Effect(ns.JIGGLE, v2, v3, 10, 0)
	ns.AudioEvent(ns.EarthRumbleMajor, wp39)
	ns.NoWallSound(true)
	ns.WallGroupClose(v0)
	ns.CreateObject("CaveRocksHuge", wp43[0])
	ns.CreateObject("CaveBoulders", wp43[1])
	ns.CreateObject("CaveRocksHuge", wp43[2])
	ns.CreateObject("CaveRocksLarge", wp43[3])
	ns.CreateObject("CaveBoulders", wp43[4])
	ns.CreateObject("CaveRocksHuge", wp43[5])
	flag61 = false
	ns.FrameTimer(20, CaveExitCollapseSEG2)
}
func EnterExitPolygon() {
	if !(!flag63 && ns.IsCaller(ns.GetHost())) {
		goto LABEL1
	}
	CaveExitCollapse()
	PlayOutdoorMusic()
LABEL1:
	return
}
func BeginCollapseSetPiece() {
	if !(!flag62 && ns.IsCaller(ns.GetHost())) {
		goto LABEL1
	}
	PlayUndergroundMusic()
	flag62 = true
	ns.WideScreen(true)
	ns.Frozen(ns.GetHost(), true)
	ns.FrameTimer(15, SetpieceShake1)
	ns.FrameTimer(60, TunnelCollapse)
	ns.Effect(ns.JIGGLE, ns.GetObjectX(ns.GetHost()), ns.GetObjectY(ns.GetHost()), 10, 0)
	ns.AudioEvent(ns.EarthRumbleMajor, wp38)
LABEL1:
	return
}
func SquashSpider() {
	var (
		v0 float32
		v1 float32
	)
	v0 = ns.GetObjectX(obj67)
	v1 = ns.GetObjectY(obj67)
	ns.MoveObject(obj67, v0, v1+15)
	ns.WallGroupOpen(ns.WallGroup("BearCageWalls"))
	ns.Frozen(ns.GetHost(), false)
	ns.WideScreen(false)
	ns.NoWallSound(false)
}
func KillBear() {
	ns.AudioEvent(ns.BearHurt, wp38)
	ns.Damage(ns.Object("BigSpider"), 0, 200, 0)
	ns.FrameTimer(10, SquashSpider)
}
func DropCrushingRock() {
	obj67 = ns.CreateObject("CaveBoulders", ns.Waypoint("SpiderWP"))
	ns.Raise(obj67, 100)
	ns.FrameTimer(5, KillBear)
}
func SquashSpiderReport() {
	var v0 ns.ObjectID
	v0 = ns.Object("BigSpider")
	ns.SetCallback(v0, 11, NullCallback)
	ns.LookAtObject(v0, ns.Object("BearFaceLocation"))
	ns.HitLocation(v0, ns.GetObjectX(v0)-1, ns.GetObjectY(v0)+1)
	ns.AudioEvent(ns.BearAttackInit, wp38)
	ns.WallGroupClose(ns.WallGroup("BearCageWalls"))
	ns.FrameTimer(20, DropCrushingRock)
}
func BearDie() {
	ns.AudioEvent(ns.FistHit, wp38)
	ns.Effect(ns.JIGGLE, ns.GetObjectX(ns.GetHost()), ns.GetObjectY(ns.GetHost()), 10, 0)
}
func ReleaseStun() {
	var (
		v0 float32
		v1 float32
		v2 ns.WaypointID
		v3 ns.WaypointID
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	v2 = ns.Waypoint("BoulderWP2")
	v3 = ns.Waypoint("BoulderWP3")
	ns.AudioEvent(ns.FistHit, wp38)
	ns.Effect(ns.JIGGLE, v0, v1, 10, 0)
	ns.CreateObject("BoulderCave", v2)
	ns.CreateObject("BoulderCave", v3)
}
func GoBigSpider() {
	var (
		v0 ns.ObjectID
		v1 ns.WaypointID
	)
	v0 = ns.Object("BigSpider")
	v1 = ns.Waypoint("SpiderWP")
	ns.AudioEvent(ns.BearRecognize, wp38)
	ns.SetCallback(v0, 11, SquashSpiderReport)
	ns.Move(v0, v1)
}
func FallingRocksSEG5() {
	var (
		v0 ns.ObjectID
		v1 ns.WaypointID
		v2 ns.WaypointID
		v3 ns.WaypointID
		v4 ns.WallGroupID
		v5 ns.WallGroupID
		v6 ns.WallGroupID
		v7 float32
		v8 float32
	)
	v7 = ns.GetObjectX(ns.GetHost())
	v8 = ns.GetObjectY(ns.GetHost())
	ns.Effect(ns.JIGGLE, v7, v8, 10, 0)
	v0 = ns.Object("Treasure")
	v4 = ns.WallGroup("TunnelWalls")
	v5 = ns.WallGroup("TunnelDestructableWalls2")
	v6 = ns.WallGroup("TunnelDestructableWalls3")
	v1 = ns.Waypoint("BoulderWP1")
	v2 = ns.Waypoint("BoulderWP4")
	v3 = ns.Waypoint("BoulderWP5")
	ns.WallGroupOpen(v5)
	ns.WallGroupBreak(v6)
	ns.NoWallSound(false)
	ns.AudioEvent(ns.EarthRumbleMajor, wp38)
	ns.NoWallSound(true)
	ns.WallGroupClose(v4)
	ns.Delete(v0)
	ns.CreateObject("BoulderCave", v1)
	ns.CreateObject("BoulderCave", v2)
	ns.CreateObject("BoulderCave", v3)
	ns.CreateObject("CaveBoulders", wp42[0])
	ns.CreateObject("CaveBoulders", wp42[2])
	ns.CreateObject("CaveBoulders", wp42[4])
	ns.CreateObject("CaveBoulders", wp42[6])
	ns.CreateObject("CaveBoulders", wp42[8])
	ns.CreateObject("CaveBoulders", wp42[10])
	ns.CreateObject("CaveBoulders", wp42[12])
	ns.CreateObject("CaveBoulders", wp42[14])
	ns.FrameTimer(5, GoBigSpider)
	ns.FrameTimer(20, ReleaseStun)
}
func FallingRocksSEG4() {
	var (
		v0 float32
		v1 float32
		v2 ns.WallGroupID
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	v2 = ns.WallGroup("TunnelDestructableWalls1")
	ns.Effect(ns.JIGGLE, v0, v1, 15, 0)
	ns.NoWallSound(true)
	ns.WallGroupOpen(v2)
	ns.CreateObject("CaveRocksLarge", wp42[2])
	ns.CreateObject("CaveRocksLarge", wp42[5])
	ns.CreateObject("CaveRocksLarge", wp42[11])
	ns.CreateObject("CaveRocksLarge", wp42[12])
	ns.AudioEvent(ns.FistHit, wp38)
	ns.FrameTimer(10, FallingRocksSEG5)
}
func FallingRocksSEG3() {
	var (
		v0 float32
		v1 float32
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.Effect(ns.JIGGLE, v0, v1, 10, 0)
	ns.CreateObject("CaveRocksMedium", wp42[1])
	ns.CreateObject("CaveRocksLarge", wp42[6])
	ns.CreateObject("CaveRocksMedium", wp42[8])
	ns.CreateObject("CaveRocksLarge", wp42[15])
	ns.AudioEvent(ns.EarthRumbleMajor, wp38)
	ns.FrameTimer(7, FallingRocksSEG4)
}
func FallingRocksSEG2() {
	var (
		v0 float32
		v1 float32
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.Effect(ns.JIGGLE, v0, v1, 15, 0)
	ns.AudioEvent(ns.GolemMove, wp38)
	ns.CreateObject("CaveRocks", wp42[3])
	ns.CreateObject("CaveRocks", wp42[4])
	ns.CreateObject("CaveRocks", wp42[10])
	ns.CreateObject("CaveRocks", wp42[13])
	ns.FrameTimer(7, FallingRocksSEG3)
}
func MineShake1SEG2() {
	var (
		v0 float32
		v1 float32
		v2 ns.WaypointID
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	v2 = ns.Waypoint("MineShake1AudioOrigin")
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.AudioEvent(ns.EarthRumbleMinor, v2)
	ns.CreateObject("CaveRocksTiny", wp45[1])
	ns.CreateObject("CaveRocksPebbles", wp45[2])
	ns.CreateObject("CaveRocksTiny", wp45[4])
	ns.CreateObject("CaveRocksTiny", wp45[6])
	ivar56 = 500
	gvar55 = 200
}
func MineShake1() {
	var (
		v0 float32
		v1 float32
		v2 ns.ObjectGroupID
		v3 ns.WaypointID
	)
	if !(!flag64 && ns.IsCaller(ns.GetHost())) {
		goto LABEL1
	}
	flag64 = true
	v3 = ns.Waypoint("MineShake1AudioOrigin")
	v2 = ns.ObjectGroup("MineShake1Triggers")
	ns.AudioEvent(ns.EarthRumbleMajor, v3)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.ObjectGroupOff(v2)
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.CreateObject("CaveRocksTiny", wp45[0])
	ns.CreateObject("CaveRocksTiny", wp45[3])
	ns.CreateObject("CaveRocksPebbles", wp45[5])
	ns.CreateObject("CaveRocksTiny", wp45[7])
	ns.FrameTimer(7, MineShake1SEG2)
LABEL1:
	return
}
func MineShake2SEG2() {
	var (
		v0 float32
		v1 float32
		v2 ns.WaypointID
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	v2 = ns.Waypoint("MineShake2AudioOrigin")
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.AudioEvent(ns.EarthRumbleMinor, v2)
	ns.CreateObject("CaveRocksTiny", wp46[1])
	ns.CreateObject("CaveRocksTiny", wp46[2])
	fvar50 = 8
	ivar56 = 120
}
func MineShake2() {
	var (
		v0 float32
		v1 float32
		v2 ns.ObjectGroupID
		v3 ns.WaypointID
	)
	if !(!flag65 && ns.IsCaller(ns.GetHost())) {
		goto LABEL1
	}
	flag65 = true
	v3 = ns.Waypoint("MineShake2AudioOrigin")
	v2 = ns.ObjectGroup("MineShake2Triggers")
	ns.AudioEvent(ns.EarthRumbleMajor, v3)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.ObjectGroupOff(v2)
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.CreateObject("CaveRocksTiny", wp46[0])
	ns.CreateObject("CaveRocksTiny", wp46[3])
	ns.FrameTimer(7, MineShake2SEG2)
LABEL1:
	return
}
func MineShake3SEG2() {
	var (
		v0 float32
		v1 float32
		v2 ns.WaypointID
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	v2 = ns.Waypoint("MineShake3AudioOrigin")
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.AudioEvent(ns.EarthRumbleMinor, v2)
	ns.CreateObject("CaveRocksTiny", wp47[1])
	ns.CreateObject("CaveRocksTiny", wp47[2])
}
func MineShake3() {
	var (
		v0 float32
		v1 float32
		v2 ns.ObjectGroupID
		v3 ns.WaypointID
	)
	if !(!flag66 && ns.IsCaller(ns.GetHost())) {
		goto LABEL1
	}
	flag66 = true
	v3 = ns.Waypoint("MineShake3AudioOrigin")
	v2 = ns.ObjectGroup("MineShake3Triggers")
	ns.AudioEvent(ns.EarthRumbleMajor, v3)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.ObjectGroupOff(v2)
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.CreateObject("CaveRocksTiny", wp47[0])
	ns.CreateObject("CaveRocksTiny", wp47[3])
	ns.FrameTimer(7, MineShake3SEG2)
LABEL1:
	return
}
func MineShake4SEG2() {
	var (
		v0 float32
		v1 float32
		v2 ns.WaypointID
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	v2 = ns.Waypoint("MineShake4AudioOrigin")
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.AudioEvent(ns.EarthRumbleMinor, v2)
	ns.CreateObject("CaveRocksTiny", wp48[1])
	ns.CreateObject("CaveRocksTiny", wp48[2])
	fvar50 = 8
	ivar56 = 120
}
func MineShake5SEG2() {
	var (
		v0 float32
		v1 float32
		v2 ns.WaypointID
	)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	v2 = ns.Waypoint("MineShake5AudioOrigin")
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.AudioEvent(ns.EarthRumbleMinor, v2)
	ns.CreateObject("CaveRocksTiny", wp49[1])
	ns.CreateObject("CaveRocksTiny", wp49[2])
	fvar50 = 9
	ivar57 = 2
	ivar58 = 1
	ivar56 = 20
}
func MineShake5() {
	var (
		v0 float32
		v1 float32
		v2 ns.ObjectGroupID
		v3 ns.WaypointID
	)
	v3 = ns.Waypoint("MineShake5AudioOrigin")
	v2 = ns.ObjectGroup("MineShake5Triggers")
	ns.AudioEvent(ns.EarthRumbleMajor, v3)
	v0 = ns.GetObjectX(ns.GetHost())
	v1 = ns.GetObjectY(ns.GetHost())
	ns.ObjectGroupOff(v2)
	ns.Effect(ns.JIGGLE, v0, v1, 5, 0)
	ns.CreateObject("CaveRocksTiny", wp49[0])
	ns.CreateObject("CaveRocksTiny", wp49[3])
	ns.FrameTimer(7, MineShake5SEG2)
}
func EntranceShake() {
	ns.ObjectGroupOff(gvar34)
	ns.ObjectGroupOn(gvar35)
	flag61 = true
	ivar56 = 600
	Quake()
}
func StopShakes() {
	ns.ObjectGroupOff(gvar35)
	flag61 = false
	ns.CancelTimer(gvar59)
	ns.ObjectGroupOn(gvar34)
}
func LittleSpiderRecognize() {
	ns.CreatureIdle(ns.Object("LittleSpider"))
	ns.AggressionLevel(ns.Object("LittleSpider"), 0.83)
	ns.SetCallback(ns.Object("LittleSpider"), 3, NullCallback)
}
func Secret01Found() {
	ns.ObjectGroupOff(ns.ObjectGroup("Secret01Triggers"))
	ns.MoveWaypoint(wp41, ns.GetObjectX(ns.GetHost()), ns.GetObjectY(ns.GetHost()))
	ns.AudioEvent(ns.SecretFound, wp41)
	ns.PrintToAll("GeneralPrint:SecretFound")
	ns.GiveXp(ns.GetHost(), 50)
}
func Secret02Found() {
	ns.ObjectGroupOff(ns.ObjectGroup("Secret02Triggers"))
	ns.MoveWaypoint(wp41, ns.GetObjectX(ns.GetHost()), ns.GetObjectY(ns.GetHost()))
	ns.AudioEvent(ns.SecretFound, wp41)
	ns.PrintToAll("GeneralPrint:SecretFound")
	ns.GiveXp(ns.GetHost(), 50)
}
func SecretElevatorFound() {
	ns.ObjectOff(ns.GetTrigger())
	ns.MoveWaypoint(wp41, ns.GetObjectX(ns.GetHost()), ns.GetObjectY(ns.GetHost()))
	ns.AudioEvent(ns.SecretFound, wp41)
	ns.PrintToAll("GeneralPrint:SecretFound")
	ns.GiveXp(ns.GetHost(), 50)
}
func Secret03Found() {
	ns.ObjectOff(ns.GetTrigger())
	ns.MoveWaypoint(wp41, ns.GetObjectX(ns.GetHost()), ns.GetObjectY(ns.GetHost()))
	ns.AudioEvent(ns.SecretFound, wp41)
	ns.PrintToAll("GeneralPrint:SecretFound")
	ns.GiveXp(ns.GetHost(), 50)
}
func OnEvent(typ string) {
	switch typ {
	case "PlayerDeath":
		PlayerDeath()
	case "MapInitialize":
		MapInitialize()
	case "MapEntry":
		MapEntry()
	}
}
