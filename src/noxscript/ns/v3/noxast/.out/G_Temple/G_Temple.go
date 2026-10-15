package script

import (
	"strconv"
	"github.com/brutusmcforce/opennox/noxscript/ns/v3"
)

var (
	gvar4   ns.ObjectGroupID
	gvar5   ns.ObjectGroupID
	gvar6   ns.ObjectGroupID
	gvar7   ns.ObjectGroupID
	gvar8   ns.ObjectGroupID
	gvar9   ns.ObjectGroupID
	gvar10  ns.ObjectGroupID
	gvar11  ns.ObjectGroupID
	gvar12  ns.ObjectGroupID
	gvar13  ns.ObjectGroupID
	gvar14  ns.ObjectGroupID
	gvar15  ns.ObjectGroupID
	gvar16  ns.ObjectGroupID
	gvar17  ns.ObjectGroupID
	gvar18  ns.ObjectGroupID
	gvar19  ns.ObjectGroupID
	gvar20  ns.ObjectGroupID
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
	gvar34  ns.ObjectGroupID
	gvar35  ns.ObjectGroupID
	gvar36  ns.ObjectGroupID
	gvar37  ns.ObjectGroupID
	gvar38  ns.ObjectGroupID
	gvar39  ns.ObjectGroupID
	gvar40  int
	gvar41  int
	gvar42  int
	ivar43  int
	flag44  bool
	obj45   ns.ObjectID
	obj46   ns.ObjectID
	obj47   ns.ObjectID
	obj48   ns.ObjectID
	obj49   ns.ObjectID
	obj50   ns.ObjectID
	obj51   ns.ObjectID
	obj52   ns.ObjectID
	obj53   ns.ObjectID
	obj54   ns.ObjectID
	obj55   ns.ObjectID
	obj56   ns.ObjectID
	obj57   ns.ObjectID
	obj58   ns.ObjectID
	obj59   ns.ObjectID
	obj60   ns.ObjectID
	obj61   ns.ObjectID
	obj62   ns.ObjectID
	obj63   ns.ObjectID
	obj64   ns.ObjectID
	obj65   ns.ObjectID
	obj66   ns.ObjectID
	obj67   ns.ObjectID
	obj68   ns.ObjectID
	obj69   ns.ObjectID
	obj70   ns.ObjectID
	obj71   ns.ObjectID
	obj72   ns.ObjectID
	obj73   ns.ObjectID
	obj74   ns.ObjectID
	obj75   ns.ObjectID
	gvar76  ns.ObjectGroupID
	gvar77  ns.ObjectGroupID
	wp78    ns.WaypointID
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
	wp95    ns.WaypointID
	wp96    ns.WaypointID
	wp97    ns.WaypointID
	wp98    ns.WaypointID
	wp99    ns.WaypointID
	wp100   ns.WaypointID
	wp101   ns.WaypointID
	wp102   ns.WaypointID
	wp103   ns.WaypointID
	wp104   ns.WaypointID
	wp105   ns.WaypointID
	wp106   ns.WaypointID
	wp107   ns.WaypointID
	wp108   ns.WaypointID
	wp109   ns.WaypointID
	wp110   ns.WaypointID
	wp111   ns.WaypointID
	wp112   ns.WaypointID
	wp113   ns.WaypointID
	gvar114 ns.WallGroupID
	gvar115 ns.WallGroupID
	gvar116 ns.WallGroupID
	gvar117 ns.WallGroupID
	gvar118 ns.WallGroupID
	gvar119 ns.WallGroupID
	gvar120 ns.WallGroupID
	flag121 bool
	obj122  ns.ObjectID
	obj123  ns.ObjectID
	obj124  ns.ObjectID
	wp125   ns.WaypointID
	wp126   ns.WaypointID
	wp127   ns.WaypointID
	wp128   ns.WaypointID
	wp129   ns.WaypointID
	obj130  ns.ObjectID
	obj131  ns.ObjectID
	obj132  ns.ObjectID
	wp133   ns.WaypointID
	wp134   ns.WaypointID
	wp135   ns.WaypointID
	gvar136 ns.WallGroupID
	gvar137 ns.WallGroupID
	gvar138 ns.WallGroupID
	obj139  [32]ns.ObjectID
	gvar140 ns.WallGroupID
	gvar141 ns.WallGroupID
	gvar142 ns.WallGroupID
	gvar143 ns.WallGroupID
	gvar144 ns.WallGroupID
	obj145  ns.ObjectID
	obj146  ns.ObjectID
	obj147  ns.ObjectID
	obj148  ns.ObjectID
	obj149  ns.ObjectID
	obj150  ns.ObjectID
	obj151  ns.ObjectID
	obj152  ns.ObjectID
	obj153  ns.ObjectID
	obj154  ns.ObjectID
	obj155  ns.ObjectID
	obj156  ns.ObjectID
	obj157  ns.ObjectID
	obj158  ns.ObjectID
)

func init() {
	gvar40 = 0
	gvar41 = 1
	gvar42 = gvar40
	ivar43 = 150
	flag44 = false
	flag121 = false
}
func InitializeArrowTraps() {
	gvar4 = ns.ObjectGroup("ArrowTrap01Row01")
	gvar5 = ns.ObjectGroup("ArrowTrap01Row02")
	gvar6 = ns.ObjectGroup("ArrowTrap01Row03")
	gvar7 = ns.ObjectGroup("ArrowTrap02Row01")
	gvar8 = ns.ObjectGroup("ArrowTrap03Row01")
	gvar9 = ns.ObjectGroup("ArrowTrap04Row01")
	gvar10 = ns.ObjectGroup("ArrowTrap05Row01")
	gvar11 = ns.ObjectGroup("ArrowTrap05Row02")
	gvar12 = ns.ObjectGroup("ArrowTrap05Row03")
	gvar13 = ns.ObjectGroup("ArrowTrap06Row01")
	gvar14 = ns.ObjectGroup("ArrowTrap07Row01")
	gvar15 = ns.ObjectGroup("ArrowTrap08Row01")
	gvar16 = ns.ObjectGroup("ArrowTrap08Row02")
	gvar17 = ns.ObjectGroup("ArrowTrap09Row01")
	gvar18 = ns.ObjectGroup("ArrowTrap09Row02")
	gvar19 = ns.ObjectGroup("ArrowTrap09Row03")
	gvar20 = ns.ObjectGroup("ArrowTrap09Row04")
	gvar21 = ns.ObjectGroup("ArrowTrap09Row05")
	gvar22 = ns.ObjectGroup("ArrowTrap09Row06")
	gvar23 = ns.ObjectGroup("ArrowTrap09Row07")
	gvar24 = ns.ObjectGroup("ArrowTrap09Row08")
	gvar25 = ns.ObjectGroup("ArrowTrap09Row09")
	gvar26 = ns.ObjectGroup("ArrowTrap09Row10")
	gvar27 = ns.ObjectGroup("ArrowTrap09Row11")
	gvar28 = ns.ObjectGroup("ArrowTrap09Row12")
	gvar29 = ns.ObjectGroup("ArrowTrap09Row13")
	gvar30 = ns.ObjectGroup("ArrowTrap09Row14")
	gvar31 = ns.ObjectGroup("ArrowTrap09Row15")
	gvar32 = ns.ObjectGroup("ArrowTrap09Row16")
	gvar33 = ns.ObjectGroup("ArrowTrap09Row17")
	gvar34 = ns.ObjectGroup("ArrowTrap09Row18")
	gvar35 = ns.ObjectGroup("ArrowTrap09Row19")
	gvar36 = ns.ObjectGroup("ArrowTrap09Row20")
	gvar37 = ns.ObjectGroup("ArrowTrap09Row21")
	gvar38 = ns.ObjectGroup("ArrowTrap09Row22")
	gvar39 = ns.ObjectGroup("ArrowTrap10Row01")
}
func DisableArrowTrap01Row01() {
	ns.ObjectGroupOff(gvar4)
}
func ArrowTrap01Row01() {
	ns.ObjectGroupOn(gvar4)
	ns.FrameTimer(1, DisableArrowTrap01Row01)
}
func DisableArrowTrap01Row02() {
	ns.ObjectGroupOff(gvar5)
}
func ArrowTrap01Row02() {
	ns.ObjectGroupOn(gvar5)
	ns.FrameTimer(1, DisableArrowTrap01Row02)
}
func DisableArrowTrap01Row03() {
	ns.ObjectGroupOff(gvar6)
}
func ArrowTrap01Row03() {
	ns.ObjectGroupOn(gvar6)
	ns.FrameTimer(1, DisableArrowTrap01Row03)
}
func DisableArrowTrap02Row01() {
	ns.ObjectGroupOff(gvar7)
}
func ArrowTrap02Row01() {
	ns.ObjectGroupOn(gvar7)
	ns.FrameTimer(1, DisableArrowTrap02Row01)
}
func DisableArrowTrap03Row01() {
	ns.ObjectGroupOff(gvar8)
}
func ArrowTrap03Row01() {
	ns.ObjectGroupOn(gvar8)
}
func DisableArrowTrap04Row01() {
	ns.ObjectGroupOff(gvar9)
}
func ArrowTrap04Row01() {
	ns.ObjectGroupOn(gvar9)
}
func DisableArrowTrap05Row01() {
	ns.ObjectGroupOff(gvar10)
}
func ArrowTrap05Row01() {
	ns.ObjectGroupOn(gvar10)
	ns.FrameTimer(1, DisableArrowTrap05Row01)
}
func DisableArrowTrap05Row02() {
	ns.ObjectGroupOff(gvar11)
}
func ArrowTrap05Row02() {
	ns.ObjectGroupOn(gvar11)
	ns.FrameTimer(1, DisableArrowTrap05Row02)
}
func DisableArrowTrap05Row03() {
	ns.ObjectGroupOff(gvar12)
}
func ArrowTrap05Row03() {
	ns.ObjectGroupOn(gvar12)
	ns.FrameTimer(1, DisableArrowTrap05Row03)
}
func DisableArrowTrap06Row01() {
	ns.ObjectGroupOff(gvar13)
}
func ArrowTrap06Row01() {
	ns.ObjectGroupOn(gvar13)
	ns.FrameTimer(1, DisableArrowTrap06Row01)
}
func DisableArrowTrap07Row01() {
	ns.ObjectGroupOff(gvar14)
}
func ArrowTrap07Row01() {
	ns.ObjectGroupOn(gvar14)
	ns.FrameTimer(1, DisableArrowTrap07Row01)
}
func DisableArrowTrap08Row01() {
	ns.ObjectGroupOff(gvar15)
}
func ArrowTrap08Row01() {
	ns.ObjectGroupOn(gvar15)
	ns.FrameTimer(1, DisableArrowTrap08Row01)
}
func DisableArrowTrap08Row02() {
	ns.ObjectGroupOff(gvar16)
}
func ArrowTrap08Row02() {
	ns.ObjectGroupOn(gvar16)
	ns.FrameTimer(1, DisableArrowTrap08Row02)
}
func DisableArrowTrap09Row01() {
	ns.ObjectGroupOff(gvar17)
}
func ArrowTrap09Row01() {
	ns.ObjectGroupOn(gvar17)
	ns.FrameTimer(1, DisableArrowTrap09Row01)
}
func DisableArrowTrap09Row02() {
	ns.ObjectGroupOff(gvar18)
}
func ArrowTrap09Row02() {
	ns.ObjectGroupOn(gvar18)
	ns.FrameTimer(1, DisableArrowTrap09Row02)
}
func DisableArrowTrap09Row03() {
	ns.ObjectGroupOff(gvar19)
}
func ArrowTrap09Row03() {
	ns.ObjectGroupOn(gvar19)
	ns.FrameTimer(1, DisableArrowTrap09Row03)
}
func DisableArrowTrap09Row04() {
	ns.ObjectGroupOff(gvar20)
}
func ArrowTrap09Row04() {
	ns.ObjectGroupOn(gvar20)
	ns.FrameTimer(1, DisableArrowTrap09Row04)
}
func DisableArrowTrap09Row05() {
	ns.ObjectGroupOff(gvar21)
}
func ArrowTrap09Row05() {
	ns.ObjectGroupOn(gvar21)
	ns.FrameTimer(1, DisableArrowTrap09Row05)
}
func DisableArrowTrap09Row06() {
	ns.ObjectGroupOff(gvar22)
}
func ArrowTrap09Row06() {
	ns.ObjectGroupOn(gvar22)
	ns.FrameTimer(1, DisableArrowTrap09Row06)
}
func DisableArrowTrap09Row07() {
	ns.ObjectGroupOff(gvar23)
}
func ArrowTrap09Row07() {
	ns.ObjectGroupOn(gvar23)
	ns.FrameTimer(1, DisableArrowTrap09Row07)
}
func DisableArrowTrap09Row08() {
	ns.ObjectGroupOff(gvar24)
}
func ArrowTrap09Row08() {
	ns.ObjectGroupOn(gvar24)
	ns.FrameTimer(1, DisableArrowTrap09Row08)
}
func DisableArrowTrap09Row09() {
	ns.ObjectGroupOff(gvar25)
}
func ArrowTrap09Row09() {
	ns.ObjectGroupOn(gvar25)
	ns.FrameTimer(1, DisableArrowTrap09Row09)
}
func DisableArrowTrap09Row10() {
	ns.ObjectGroupOff(gvar26)
}
func ArrowTrap09Row10() {
	ns.ObjectGroupOn(gvar26)
	ns.FrameTimer(1, DisableArrowTrap09Row10)
}
func DisableArrowTrap09Row11() {
	ns.ObjectGroupOff(gvar27)
}
func ArrowTrap09Row11() {
	ns.ObjectGroupOn(gvar27)
	ns.FrameTimer(1, DisableArrowTrap09Row11)
}
func DisableArrowTrap09Row12() {
	ns.ObjectGroupOff(gvar28)
}
func ArrowTrap09Row12() {
	ns.ObjectGroupOn(gvar28)
	ns.FrameTimer(1, DisableArrowTrap09Row12)
}
func DisableArrowTrap09Row13() {
	ns.ObjectGroupOff(gvar29)
}
func ArrowTrap09Row13() {
	ns.ObjectGroupOn(gvar29)
	ns.FrameTimer(1, DisableArrowTrap09Row13)
}
func DisableArrowTrap09Row14() {
	ns.ObjectGroupOff(gvar30)
}
func ArrowTrap09Row14() {
	ns.ObjectGroupOn(gvar30)
	ns.FrameTimer(1, DisableArrowTrap09Row14)
}
func DisableArrowTrap09Row15() {
	ns.ObjectGroupOff(gvar31)
}
func ArrowTrap09Row15() {
	ns.ObjectGroupOn(gvar31)
	ns.FrameTimer(1, DisableArrowTrap09Row15)
}
func DisableArrowTrap09Row16() {
	ns.ObjectGroupOff(gvar32)
}
func ArrowTrap09Row16() {
	ns.ObjectGroupOn(gvar32)
	ns.FrameTimer(1, DisableArrowTrap09Row16)
}
func DisableArrowTrap09Row17() {
	ns.ObjectGroupOff(gvar33)
}
func ArrowTrap09Row17() {
	ns.ObjectGroupOn(gvar33)
	ns.FrameTimer(1, DisableArrowTrap09Row17)
}
func DisableArrowTrap09Row18() {
	ns.ObjectGroupOff(gvar34)
}
func ArrowTrap09Row18() {
	ns.ObjectGroupOn(gvar34)
	ns.FrameTimer(1, DisableArrowTrap09Row18)
}
func DisableArrowTrap09Row19() {
	ns.ObjectGroupOff(gvar35)
}
func ArrowTrap09Row19() {
	ns.ObjectGroupOn(gvar35)
	ns.FrameTimer(1, DisableArrowTrap09Row19)
}
func DisableArrowTrap09Row20() {
	ns.ObjectGroupOff(gvar36)
}
func ArrowTrap09Row20() {
	ns.ObjectGroupOn(gvar36)
	ns.FrameTimer(1, DisableArrowTrap09Row20)
}
func DisableArrowTrap09Row21() {
	ns.ObjectGroupOff(gvar37)
}
func ArrowTrap09Row21() {
	ns.ObjectGroupOn(gvar37)
	ns.FrameTimer(1, DisableArrowTrap09Row21)
}
func DisableArrowTrap09Row22() {
	ns.ObjectGroupOff(gvar38)
}
func ArrowTrap09Row22() {
	ns.ObjectGroupOn(gvar38)
	ns.FrameTimer(1, DisableArrowTrap09Row22)
}
func DisableArrowTrap10Row01() {
	ns.ObjectGroupOff(gvar39)
}
func ArrowTrap10Row01() {
	ns.ObjectGroupOn(gvar39)
	ns.FrameTimer(1, DisableArrowTrap10Row01)
}
func SpikeBlock01Loop() {
	var v0 int
	switch v0 = gvar42; v0 {
	case gvar40:
		goto LABEL1
	case gvar41:
		goto LABEL2
	default:
		goto LABEL3
	}
LABEL1:
	ns.Move(obj45, wp79)
	ns.AudioEvent(ns.SpikeBlockMove, wp78)
	gvar42 = gvar41
	goto LABEL3
LABEL2:
	ns.Move(obj45, wp78)
	ns.AudioEvent(ns.SpikeBlockMove, wp79)
	gvar42 = gvar40
	goto LABEL3
LABEL3:
	ns.FrameTimer(ivar43, SpikeBlock01Loop)
}
func ActivateSpikeBlockGroup01() {
	ns.Move(obj46, wp80)
	ns.Move(obj47, wp81)
	ns.Move(obj48, wp82)
}
func InitializeBlocks() {
	obj45 = ns.Object("SpikeBlock01")
	obj46 = ns.Object("SpikeBlock02")
	obj47 = ns.Object("SpikeBlock03")
	obj48 = ns.Object("SpikeBlock04")
	obj49 = ns.Object("TempleEntrance")
	obj50 = ns.Object("StoneBlock01")
	obj51 = ns.Object("StoneBlock02")
	obj52 = ns.Object("StoneBlock03")
	obj53 = ns.Object("StoneBlock04")
	obj54 = ns.Object("StoneBlock05")
	obj55 = ns.Object("StoneBlock06")
	obj56 = ns.Object("StoneBlock07")
	obj57 = ns.Object("StoneBlock08")
	obj58 = ns.Object("StoneBlock09")
	obj59 = ns.Object("StoneBlock10")
	obj60 = ns.Object("StoneBlock11")
	obj61 = ns.Object("StoneBlock12")
	obj62 = ns.Object("StoneBlock13")
	obj63 = ns.Object("StoneBlock14")
	obj64 = ns.Object("StoneBlock15")
	obj65 = ns.Object("StoneBlock16")
	obj66 = ns.Object("StoneBlock17")
	obj67 = ns.Object("StoneBlock18")
	obj68 = ns.Object("GeneratorPusher01")
	obj69 = ns.Object("GeneratorPusher02")
	obj70 = ns.Object("GeneratorPusher03")
	obj71 = ns.Object("GeneratorPusher04")
	obj72 = ns.Object("GeneratorPusherMover01")
	obj73 = ns.Object("GeneratorPusherMover02")
	obj74 = ns.Object("GeneratorPusherMover03")
	obj75 = ns.Object("GeneratorPusherMover04")
	gvar76 = ns.ObjectGroup("BlockGroup04Phase2Triggers")
	gvar77 = ns.ObjectGroup("PusherTriggers")
	wp78 = ns.Waypoint("SpikeBlock01Home")
	wp79 = ns.Waypoint("SpikeBlock01Dest")
	wp80 = ns.Waypoint("SpikeBlock02WP")
	wp81 = ns.Waypoint("SpikeBlock03WP")
	wp82 = ns.Waypoint("SpikeBlock04WP")
	wp83 = ns.Waypoint("StoneBlock01WP")
	wp84 = ns.Waypoint("StoneBlock02WP")
	wp85 = ns.Waypoint("StoneBlock03WP")
	wp86 = ns.Waypoint("StoneBlock04WP")
	wp87 = ns.Waypoint("StoneBlock05WP")
	wp88 = ns.Waypoint("StoneBlock06WP")
	wp89 = ns.Waypoint("StoneBlock07WP")
	wp90 = ns.Waypoint("StoneBlock08WP")
	wp91 = ns.Waypoint("StoneBlock09WP")
	wp92 = ns.Waypoint("StoneBlock10WP")
	wp93 = ns.Waypoint("StoneBlock11WP")
	wp94 = ns.Waypoint("StoneBlock12WP")
	wp95 = ns.Waypoint("StoneBlock13WP")
	wp96 = ns.Waypoint("StoneBlock14WP")
	wp97 = ns.Waypoint("StoneBlock15WP")
	wp98 = ns.Waypoint("StoneBlock16WP")
	wp99 = ns.Waypoint("StoneBlock15WP2")
	wp100 = ns.Waypoint("StoneBlock16WP2")
	wp101 = ns.Waypoint("StoneBlock17WP")
	wp102 = ns.Waypoint("StoneBlock18WP")
	wp103 = ns.Waypoint("GeneratorPusher01WP")
	wp104 = ns.Waypoint("GeneratorPusher02WP")
	wp105 = ns.Waypoint("GeneratorPusher03WP")
	wp106 = ns.Waypoint("GeneratorPusher04WP")
	wp107 = ns.Waypoint("StoneBlockAudioOrigin01")
	wp108 = ns.Waypoint("StoneBlockAudioOrigin02")
	wp109 = ns.Waypoint("StoneBlockAudioOrigin03")
	wp110 = ns.Waypoint("StoneBlockAudioOrigin04")
	wp111 = ns.Waypoint("StoneBlockAudioOrigin05")
	wp112 = ns.Waypoint("StoneBlockAudioOrigin06")
	wp113 = ns.Waypoint("TempleEntranceAudioOrigin")
	gvar114 = ns.WallGroup("BlockGroup01Walls")
	gvar115 = ns.WallGroup("BlockGroup02Walls")
	gvar116 = ns.WallGroup("ElevatorWallgroup01")
	gvar117 = ns.WallGroup("ElevatorWallgroup02")
	gvar118 = ns.WallGroup("ElevatorWallgroup03")
	gvar119 = ns.WallGroup("ElevatorWallgroup04")
	gvar120 = ns.WallGroup("GeneratorPusherWalls")
	ns.Move(obj70, wp105)
	ns.Move(obj71, wp106)
}
func ActivateBlockGroup01() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj50, wp83)
	ns.Move(obj51, wp84)
	ns.WallGroupOpen(gvar114)
	ns.AudioEvent(ns.SpikeBlockMove, wp107)
}
func ActivateBlockGroup02() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Print("GeneralPrint:MsgSeWallOpen")
	ns.Move(obj52, wp85)
	ns.Move(obj53, wp86)
	ns.Move(obj54, wp87)
	ns.WallGroupOpen(gvar115)
	ns.AudioEvent(ns.SpikeBlockMove, wp108)
}
func OpenElevator01Walls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar116)
}
func OpenPitSecretWall() {
	if !ns.IsCaller(obj53) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(ns.Wall(164, 150))
	ns.WallOpen(ns.Wall(163, 151))
	ns.WallOpen(ns.Wall(162, 152))
LABEL1:
	return
}
func ActivateStoneBlock06() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj55, wp88)
	ns.AudioEvent(ns.SpikeBlockMove, wp109)
}
func ActivateStoneBlock07() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(ns.Wall(157, 139))
	ns.WallOpen(ns.Wall(143, 151))
	ns.Move(obj56, wp89)
	ns.AudioEvent(ns.SpikeBlockMove, wp88)
}
func ActivateStoneBlockGroup03() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj57, wp90)
	ns.Move(obj58, wp91)
	ns.AudioEvent(ns.SpikeBlockMove, wp90)
}
func ActivateSpikeBlock01() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj45, wp79)
	gvar42 = gvar41
	ns.AudioEvent(ns.SpikeBlockMove, wp78)
	ns.FrameTimer(ivar43, SpikeBlock01Loop)
}
func OpenElevatorWallgroup02() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar117)
}
func ActivateGeneratorHallPushers() {
	if flag44 {
		goto LABEL1
	}
	ns.ObjectGroupOff(gvar77)
	flag44 = true
	ns.WallGroupOpen(gvar120)
	ns.Move(obj68, wp103)
	ns.Move(obj69, wp104)
LABEL1:
	return
}
func OpenElevatorWallgroup03() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar118)
}
func GeneratorPusher01Die() {
	ns.ObjectOff(obj72)
}
func GeneratorPusher02Die() {
	ns.ObjectOff(obj73)
}
func GeneratorPusher03Die() {
	ns.ObjectOff(obj74)
}
func GeneratorPusher04Die() {
	ns.ObjectOff(obj75)
}
func OpenTempleEntrance() {
	ns.Move(obj63, wp96)
	ns.AudioEvent(ns.SpikeBlockMove, wp113)
	ns.AudioEvent(ns.BigGong, wp113)
}
func EnableTempleEntrance() {
	if !ns.IsCaller(obj63) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj49)
LABEL1:
	return
}
func ActivateBlockGroup04Phase1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj64, wp97)
	ns.Move(obj65, wp98)
	ns.AudioEvent(ns.SpikeBlockMove, wp112)
	ActivateSpikeBlockGroup01()
}
func EnableBlockGroup04Phase2() {
	ns.ObjectGroupOn(gvar76)
}
func ActivateBlockGroup04Phase2() {
	ns.ObjectGroupOff(gvar76)
	ns.Move(obj64, wp99)
	ns.Move(obj65, wp100)
	ns.AudioEvent(ns.SpikeBlockMove, wp99)
	ns.FrameTimer(150, EnableBlockGroup04Phase2)
}
func OpenElevatorWallgroup04() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar119)
}
func ActivateStoneBlockGroup05() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj66, wp101)
	ns.Move(obj67, wp102)
	ns.AudioEvent(ns.SpikeBlockMove, wp111)
}
func FONTrap01Loop() {
	if !flag121 {
		goto LABEL1
	}
	ns.CastSpellObjectLocation(ns.SPELL_FORCE_OF_NATURE, obj122, ns.GetWaypointX(wp125), ns.GetWaypointY(wp125))
	ns.SecondTimer(5, FONTrap01Loop)
LABEL1:
	return
}
func FireFONTrap02() {
	ns.CastSpellObjectLocation(ns.SPELL_FORCE_OF_NATURE, obj123, ns.GetWaypointX(wp126), ns.GetWaypointY(wp126))
	ns.CastSpellObjectLocation(ns.SPELL_FORCE_OF_NATURE, obj124, ns.GetWaypointX(wp127), ns.GetWaypointY(wp127))
	ns.SecondTimer(5, FireFONTrap02)
}
func InitializeFONtraps() {
	obj122 = ns.Object("FON_Origin01")
	obj123 = ns.Object("FON_Origin02")
	obj124 = ns.Object("FON_Origin03")
	wp125 = ns.Waypoint("FON_Target01")
	wp126 = ns.Waypoint("FON_Target02")
	wp127 = ns.Waypoint("FON_Target03")
}
func FireFONTrap01() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Print("GeneralPrint:MsgTrapOn")
	flag121 = true
	ns.CastSpellObjectLocation(ns.SPELL_FORCE_OF_NATURE, obj122, ns.GetWaypointX(wp125), ns.GetWaypointY(wp125))
	ns.SecondTimer(5, FONTrap01Loop)
}
func DisableFONTrap01() {
	flag121 = false
}
func OpenGoldKeyExitWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(ns.Wall(94, 152))
}
func ActivateFONtrap() {
	ns.ObjectOff(ns.GetTrigger())
	FireFONTrap02()
}
func InitializeMaze() {
	obj130 = ns.Object("MazeSpike01")
	obj131 = ns.Object("MazeSpike02")
	obj132 = ns.Object("MazeSpike03")
	wp133 = ns.Waypoint("MazeSpike01WP")
	wp134 = ns.Waypoint("MazeSpike02WP")
	wp135 = ns.Waypoint("MazeSpike03WP")
}
func InitializeMiscRooms() {
	gvar136 = ns.WallGroup("SnipeRoomWalls")
	gvar138 = ns.WallGroup("SnipeRoomExitWalls")
	gvar137 = ns.WallGroup("PusherRoomExitWalls")
}
func InitializeTemple() {
	gvar142 = ns.WallGroup("ElevatorWallgroup05")
	gvar143 = ns.WallGroup("BossRoomSecretWalls")
	gvar144 = ns.WallGroup("FireRoomWalls")
}
func InitializeGoodiePits() {
	var v0 int
	v0 = 0
	v0 = 0
	for {
		if !(v0 < 32) {
			goto LABEL1
		}
		obj139[v0] = ns.Object("GoodiePit" + strconv.Itoa(v0+1))
		v0 += 1
	}
LABEL1:
	gvar140 = ns.WallGroup("GoodiePitExitWalls")
	gvar141 = ns.WallGroup("ElevatorWallgroup06")
}
func InitializeGeneratorGrids() {
	obj145 = ns.Object("GenGrid01Generator01")
	obj146 = ns.Object("GenGrid01Generator02")
	obj147 = ns.Object("GenGrid01Generator03")
	obj148 = ns.Object("GenGrid02Generator01")
	obj149 = ns.Object("GenGrid02Generator02")
	obj150 = ns.Object("GenGrid02Generator03")
	obj151 = ns.Object("GenGrid02Generator04")
	obj152 = ns.Object("GenGrid02Generator05")
	obj153 = ns.Object("GenGrid02Generator06")
	obj154 = ns.Object("GenGrid02Generator07")
	obj155 = ns.Object("GenGrid03Generator01")
	obj156 = ns.Object("GenGrid03Generator02")
	obj157 = ns.Object("GenGrid03Generator03")
	obj158 = ns.Object("GenGrid02Golem")
	ns.ObjectOff(obj145)
	ns.ObjectOff(obj146)
	ns.ObjectOff(obj147)
	ns.ObjectOff(obj149)
	ns.ObjectOff(obj150)
	ns.ObjectOff(obj151)
	ns.ObjectOff(obj152)
	ns.ObjectOff(obj153)
	ns.ObjectOff(obj154)
	ns.ObjectOff(obj156)
	ns.ObjectOff(obj157)
}
func MapInitialize() {
	InitializeArrowTraps()
	InitializeBlocks()
	InitializeMaze()
	InitializeMiscRooms()
	InitializeFONtraps()
	InitializeTemple()
	InitializeGoodiePits()
	InitializeGeneratorGrids()
	wp128 = ns.Waypoint("WellWP")
	wp129 = ns.Waypoint("WellWP2")
}
func WellOfRestoration() {
	ns.RestoreHealth(ns.GetCaller(), ns.MaxHealth(ns.GetCaller())-ns.CurrentHealth(ns.GetCaller()))
	ns.AudioEvent(ns.RestoreHealthName, wp128)
	ns.Print("GeneralPrint:WellSignRefresh")
}
func WellOfRestoration2() {
	ns.RestoreHealth(ns.GetCaller(), ns.MaxHealth(ns.GetCaller())-ns.CurrentHealth(ns.GetCaller()))
	ns.AudioEvent(ns.RestoreHealthName, wp129)
	ns.Print("GeneralPrint:WellSignRefresh")
}
func OpenSlidingGeneratorWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(ns.Wall(137, 111))
}
func ActivateMazePath01() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj130, wp133)
}
func ActivateMazePath02() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj131, wp134)
}
func ActivateMazePath03() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Move(obj132, wp135)
}
func OpenSnipeRoomWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar136)
	ns.Print("GeneralPrint:MsgSeWallOpen")
}
func OpenSnipeRoomExitWalls() {
	ns.WallGroupOpen(gvar138)
	ns.WallGroupClose(gvar137)
}
func OpenPusherRoomExitWalls() {
	ns.WallGroupClose(gvar138)
	ns.WallGroupOpen(gvar137)
}
func OpenGoodiePitExitWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar140)
}
func OpenElevatorWallgroup06() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar141)
}
func OpenGoodiePit(a1 int) {
	ns.ObjectOn(obj139[a1-1])
}
func GoodiePit1() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(1)
}
func GoodiePit2() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(2)
}
func GoodiePit3() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(3)
}
func GoodiePit4() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(4)
}
func GoodiePit5() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(5)
}
func GoodiePit6() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(6)
}
func GoodiePit7() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(7)
}
func GoodiePit8() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(8)
}
func GoodiePit9() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(9)
}
func GoodiePit10() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(10)
}
func GoodiePit11() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(11)
}
func GoodiePit12() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(12)
}
func GoodiePit13() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(13)
}
func GoodiePit14() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(14)
}
func GoodiePit15() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(15)
}
func GoodiePit16() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(16)
}
func GoodiePit17() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(17)
}
func GoodiePit18() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(18)
}
func GoodiePit19() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(19)
}
func GoodiePit20() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(20)
}
func GoodiePit21() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(21)
}
func GoodiePit22() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(22)
}
func GoodiePit23() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(23)
}
func GoodiePit24() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(24)
}
func GoodiePit25() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(25)
}
func GoodiePit26() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(26)
}
func GoodiePit27() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(27)
}
func GoodiePit28() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(28)
}
func GoodiePit29() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(29)
}
func GoodiePit30() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(30)
}
func GoodiePit31() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(31)
}
func GoodiePit32() {
	ns.ObjectOff(ns.GetTrigger())
	OpenGoodiePit(32)
}
func OpenElevatorWallgroup05() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar142)
}
func OpenBossRoomWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar143)
}
func OpenFireRoomWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar144)
}
func GenGrid01Front01Destroy() {
	ns.ObjectOn(obj145)
}
func GenGrid01Front02Destroy() {
	ns.ObjectOn(obj146)
}
func GenGrid01Front03Destroy() {
	ns.ObjectOn(obj147)
}
func GenGrid01Back01Destroy() {
	ns.ObjectOn(obj146)
}
func GenGrid01Back02Destroy() {
	ns.ObjectOn(obj145)
	ns.ObjectOn(obj147)
}
func GenGrid01Back03Destroy() {
	ns.ObjectOn(obj146)
}
func GenGrid02Gen01Destroy() {
	ns.ObjectOn(obj149)
}
func GenGrid02Gen02Destroy() {
	ns.ObjectOn(obj150)
}
func GenGrid02Gen03Destroy() {
	ns.ObjectOn(obj151)
}
func GenGrid02Gen04Destroy() {
	ns.ObjectOn(obj152)
}
func GenGrid02Gen05Destroy() {
	ns.ObjectOn(obj153)
}
func GenGrid02Gen06Destroy() {
	ns.ObjectOn(obj154)
}
func GenGrid02Gen07Destroy() {
	ns.AggressionLevel(obj158, 0.83)
}
func GenGrid03Gen01Destroy() {
	ns.ObjectOn(obj156)
}
func GenGrid03Gen02Destroy() {
	ns.ObjectOn(obj155)
	ns.ObjectOn(obj157)
}
func GenGrid03Gen03Destroy() {
	ns.ObjectOn(obj156)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
