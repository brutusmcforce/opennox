package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	obj4    ns.ObjectID
	obj5    ns.ObjectID
	obj6    ns.ObjectID
	obj7    ns.ObjectID
	obj8    ns.ObjectID
	obj9    ns.ObjectID
	obj10   ns.ObjectID
	obj11   ns.ObjectID
	obj12   ns.ObjectID
	obj13   ns.ObjectID
	gvar14  ns.WallID
	gvar15  ns.WallID
	gvar16  ns.WallID
	gvar17  ns.WallID
	gvar18  ns.WallID
	gvar19  ns.WallID
	gvar20  ns.WallID
	gvar21  ns.WallID
	obj22   ns.ObjectID
	obj23   ns.ObjectID
	obj24   ns.ObjectID
	obj25   ns.ObjectID
	gvar26  ns.WallGroupID
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
	wp71    ns.WaypointID
	wp72    ns.WaypointID
	wp73    ns.WaypointID
	wp74    ns.WaypointID
	wp75    ns.WaypointID
	wp76    ns.WaypointID
	wp77    ns.WaypointID
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
	wp114   ns.WaypointID
	wp115   ns.WaypointID
	wp116   ns.WaypointID
	wp117   ns.WaypointID
	wp118   ns.WaypointID
	wp119   ns.WaypointID
	wp120   ns.WaypointID
	wp121   ns.WaypointID
	wp122   ns.WaypointID
	wp123   ns.WaypointID
	wp124   ns.WaypointID
	wp125   ns.WaypointID
	flag126 bool
	wp127   ns.WaypointID
	gvar128 ns.TimerID
	gvar129 ns.ObjectGroupID
	gvar130 ns.WallGroupID
	gvar131 ns.WallGroupID
	gvar132 ns.WallGroupID
	gvar133 ns.WallGroupID
	gvar134 ns.TimerID
	ivar135 int
	flag136 bool
	gvar137 ns.WallID
	gvar138 ns.WallID
	gvar139 ns.WallID
	gvar140 ns.WallID
	gvar141 ns.WallID
	obj142  ns.ObjectID
	obj143  ns.ObjectID
	obj144  ns.ObjectID
	wp145   ns.WaypointID
	gvar146 ns.WallID
	gvar147 ns.WallID
	gvar148 ns.WallID
	gvar149 ns.WallID
	gvar150 ns.WallID
	gvar151 ns.WallID
	flag152 bool
	gvar153 ns.WallGroupID
	gvar154 ns.WallGroupID
	gvar155 ns.WallGroupID
	gvar156 ns.WallGroupID
	gvar157 ns.WallID
	gvar158 ns.WallID
	gvar159 ns.WallID
	gvar160 ns.WallID
	gvar161 ns.WallGroupID
	ivar162 int
	gvar163 ns.ObjectGroupID
	gvar164 ns.WallGroupID
	gvar165 ns.WallID
	gvar166 ns.WallID
	ivar167 int
	obj168  ns.ObjectID
	gvar169 ns.ObjectGroupID
	wp170   ns.WaypointID
	gvar171 ns.WallGroupID
	gvar172 ns.WallGroupID
	gvar173 ns.WallGroupID
	gvar174 ns.WallGroupID
	gvar175 ns.WallID
	gvar176 ns.WallID
	gvar177 ns.WallID
	gvar178 ns.WallID
	gvar179 ns.WallID
	gvar180 ns.WallID
	gvar181 ns.WallID
	gvar182 ns.WallID
	gvar183 ns.WallID
	obj184  ns.ObjectID
	obj185  ns.ObjectID
	obj186  ns.ObjectID
	obj187  ns.ObjectID
	obj188  ns.ObjectID
	obj189  ns.ObjectID
	obj190  ns.ObjectID
	obj191  ns.ObjectID
	obj192  ns.ObjectID
	obj193  ns.ObjectID
	obj194  ns.ObjectID
	wp195   ns.WaypointID
	obj196  ns.ObjectID
	obj197  ns.ObjectID
	obj198  ns.ObjectID
	obj199  ns.ObjectID
	flag200 bool
	obj201  ns.ObjectID
	obj202  ns.ObjectID
	obj203  ns.ObjectID
	obj204  ns.ObjectID
	gvar205 ns.WallGroupID
	gvar206 ns.WallGroupID
	gvar207 ns.WallGroupID
	gvar208 ns.WallGroupID
	gvar209 ns.WallGroupID
	gvar210 ns.WallGroupID
	obj211  ns.ObjectID
	obj212  ns.ObjectID
	obj213  ns.ObjectID
	obj214  ns.ObjectID
	obj215  ns.ObjectID
	obj216  ns.ObjectID
	ivar217 int
	ivar218 int
	ivar219 int
	obj220  ns.ObjectID
	obj221  ns.ObjectID
	obj222  ns.ObjectID
	gvar223 ns.ObjectGroupID
	gvar224 ns.WallGroupID
	gvar225 ns.WallGroupID
	gvar226 ns.WallGroupID
	gvar227 ns.WallGroupID
	gvar228 ns.WallGroupID
	gvar229 ns.WallGroupID
	gvar230 ns.WallGroupID
	ivar231 int
	gvar232 ns.WallID
	obj233  ns.ObjectID
	obj234  ns.ObjectID
	obj235  ns.ObjectID
	obj236  ns.ObjectID
	obj237  ns.ObjectID
	obj238  ns.ObjectID
	gvar239 ns.ObjectGroupID
	obj240  ns.ObjectID
	obj241  ns.ObjectID
	obj242  ns.ObjectID
	obj243  ns.ObjectID
	gvar244 ns.ObjectGroupID
	wp245   ns.WaypointID
	wp246   ns.WaypointID
	obj247  ns.ObjectID
	obj248  ns.ObjectID
	obj249  ns.ObjectID
	obj250  ns.ObjectID
	gvar251 ns.ObjectGroupID
	wp252   ns.WaypointID
	wp253   ns.WaypointID
	obj254  ns.ObjectID
	obj255  ns.ObjectID
	obj256  ns.ObjectID
	obj257  ns.ObjectID
	obj258  ns.ObjectID
	gvar259 ns.ObjectGroupID
	gvar260 ns.WallID
	gvar261 ns.WallGroupID
	gvar262 ns.WallGroupID
	gvar263 ns.WallGroupID
	gvar264 ns.WallGroupID
	gvar265 ns.WallGroupID
	gvar266 ns.ObjectGroupID
	gvar267 ns.ObjectGroupID
	gvar268 ns.ObjectGroupID
	gvar269 ns.TimerID
	obj270  ns.ObjectID
	obj271  ns.ObjectID
	obj272  ns.ObjectID
	obj273  ns.ObjectID
	gvar274 ns.WallID
	gvar275 ns.WallID
	gvar276 ns.WallID
	gvar277 ns.WallID
	gvar278 ns.WallID
	obj279  ns.ObjectID
	obj280  ns.ObjectID
	obj281  ns.ObjectID
	obj282  ns.ObjectID
	obj283  ns.ObjectID
	obj284  ns.ObjectID
	obj285  ns.ObjectID
	obj286  ns.ObjectID
	gvar287 ns.ObjectGroupID
	obj288  ns.ObjectID
	obj289  ns.ObjectID
	obj290  ns.ObjectID
	obj291  ns.ObjectID
	obj292  ns.ObjectID
	obj293  ns.ObjectID
	wp294   ns.WaypointID
	wp295   ns.WaypointID
	wp296   ns.WaypointID
	wp297   ns.WaypointID
	obj298  ns.ObjectID
	obj299  ns.ObjectID
	obj300  ns.ObjectID
	obj301  ns.ObjectID
	wp302   ns.WaypointID
	wp303   ns.WaypointID
	wp304   ns.WaypointID
	wp305   ns.WaypointID
	obj306  ns.ObjectID
	obj307  ns.ObjectID
	obj308  ns.ObjectID
	obj309  ns.ObjectID
	wp310   ns.WaypointID
	wp311   ns.WaypointID
	wp312   ns.WaypointID
	wp313   ns.WaypointID
	wp314   ns.WaypointID
	wp315   ns.WaypointID
	wp316   ns.WaypointID
	gvar317 ns.WallGroupID
	gvar318 ns.WallGroupID
	gvar319 ns.WallGroupID
	gvar320 ns.WallGroupID
	ivar321 int
	ivar322 int
	ivar323 int
	flag324 bool
	flag325 bool
	flag326 bool
	obj327  ns.ObjectID
	obj328  ns.ObjectID
	obj329  ns.ObjectID
	obj330  ns.ObjectID
	obj331  ns.ObjectID
	obj332  ns.ObjectID
	obj333  ns.ObjectID
	obj334  ns.ObjectID
	gvar335 ns.WallID
	gvar336 ns.WallID
	gvar337 ns.WallID
	gvar338 ns.WallID
	gvar339 ns.WallID
	ivar340 int
	wp341   ns.WaypointID
	obj342  ns.ObjectID
	gvar343 ns.WallGroupID
	wp344   ns.WaypointID
	wp345   ns.WaypointID
	wp346   ns.WaypointID
	wp347   ns.WaypointID
	wp348   ns.WaypointID
	gvar349 ns.TimerID
	ivar350 int
	ivar351 int
	flag352 bool
)

func init() {
	flag126 = false
	gvar128 = 0
	gvar134 = 0
	ivar135 = 0
	flag136 = false
	flag152 = false
	ivar162 = 0
	ivar167 = 0
	flag200 = false
	ivar217 = 0
	ivar218 = 0
	ivar219 = 0
	ivar231 = 0
	gvar269 = 0
	ivar321 = 2
	ivar322 = 2
	ivar323 = 1
	flag324 = false
	flag325 = false
	flag326 = false
	ivar340 = 0
	gvar349 = 0
	ivar350 = 0
	ivar351 = 0
	flag352 = false
}
func BreakBuzzWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.NoWallSound(true)
	ns.WallOpen(gvar14)
	ns.WallOpen(gvar15)
	ns.WallOpen(gvar16)
	ns.WallOpen(gvar17)
	ns.NoWallSound(false)
	ns.FrameTimer(5, BreakBreakWalls)
}
func BreakBreakWalls() {
	ns.WallBreak(gvar18)
	ns.WallBreak(gvar19)
	ns.WallBreak(gvar20)
	ns.WallBreak(gvar21)
}
func BuzzSawGets() {
	obj4 = ns.Object("RotSpikeMover1")
	obj5 = ns.Object("RotSpikeMover2")
	obj6 = ns.Object("RotSpikeMover3")
	obj7 = ns.Object("RotSpikeMover4")
	obj8 = ns.Object("RotSpikeMover5")
	obj9 = ns.Object("RotSpikeMover6")
	obj10 = ns.Object("BuzzSawTrapTrigger")
	obj11 = ns.Object("BuzzSawSecretTrigger")
	obj12 = ns.Object("SilverDoor1")
	obj13 = ns.Object("SilverDoor2")
	gvar14 = ns.Wall(55, 137)
	gvar15 = ns.Wall(56, 136)
	gvar16 = ns.Wall(57, 135)
	gvar17 = ns.Wall(58, 134)
	gvar18 = ns.Wall(56, 138)
	gvar19 = ns.Wall(57, 137)
	gvar20 = ns.Wall(58, 136)
	gvar21 = ns.Wall(59, 135)
}
func StartBuzzSaw() {
	if !ns.IsLocked(obj12) {
		goto LABEL1
	}
	ns.UnlockDoor(obj12)
	ns.UnlockDoor(obj13)
LABEL1:
	ns.ObjectOff(obj10)
	ns.ObjectOff(obj11)
	ns.ObjectOn(obj6)
	ns.ObjectOn(obj7)
}
func CircleGets() {
	obj22 = ns.Object("CircleBlockMover1")
	obj23 = ns.Object("CircleBlockMover2")
	obj24 = ns.Object("CircleBlockMover3")
	obj25 = ns.Object("CircleBlockMover4")
	gvar26 = ns.WallGroup("CircleWalls")
	obj27 = ns.Object("WackyBlock1")
	obj28 = ns.Object("WackyBlock2")
	obj29 = ns.Object("WackyBlock3")
	obj30 = ns.Object("WackyBlock4")
	obj31 = ns.Object("WackyBlock5")
	obj32 = ns.Object("WackyBlock6")
	obj33 = ns.Object("WackyBlock7")
	obj34 = ns.Object("WackyBlock8")
	obj35 = ns.Object("WackyBlock9")
	obj36 = ns.Object("WackyBlock10")
	obj37 = ns.Object("WackyBlock11")
	obj38 = ns.Object("WackyBlock12")
	obj39 = ns.Object("WackyBlock13")
	obj40 = ns.Object("WackyBlock14")
	obj41 = ns.Object("WackyBlock15")
	obj42 = ns.Object("WackyBlock16")
	obj43 = ns.Object("WackyBlock17")
	obj44 = ns.Object("WackyBlock18")
	obj45 = ns.Object("WackyBlock19")
	obj46 = ns.Object("WackyBlock20")
	obj47 = ns.Object("WackyBlock21")
	obj48 = ns.Object("WackyBlock22")
	obj49 = ns.Object("WackyMover1")
	obj50 = ns.Object("WackyMover2")
	obj51 = ns.Object("WackyMover3")
	obj52 = ns.Object("WackyMover4")
	obj53 = ns.Object("WackyMover5")
	obj54 = ns.Object("WackyMover6")
	obj55 = ns.Object("WackyMover7")
	obj56 = ns.Object("WackyMover8")
	obj57 = ns.Object("WackyMover9")
	obj58 = ns.Object("WackyMover10")
	obj59 = ns.Object("WackyMover11")
	obj60 = ns.Object("WackyMover12")
	obj61 = ns.Object("WackyMover13")
	obj62 = ns.Object("WackyMover14")
	obj63 = ns.Object("WackyMover15")
	obj64 = ns.Object("WackyMover16")
	obj65 = ns.Object("WackyMover17")
	obj66 = ns.Object("WackyMover18")
	obj67 = ns.Object("WackyMover19")
	obj68 = ns.Object("WackyMover20")
	obj69 = ns.Object("WackyMover21")
	obj70 = ns.Object("WackyMover22")
	wp71 = ns.Waypoint("WackyStart1")
	wp72 = ns.Waypoint("WackyStart2")
	wp73 = ns.Waypoint("WackyStart3")
	wp74 = ns.Waypoint("WackyStart4")
	wp75 = ns.Waypoint("WackyStart5")
	wp76 = ns.Waypoint("WackyStart6")
	wp77 = ns.Waypoint("WackyStart7")
	wp78 = ns.Waypoint("WackyStart8")
	wp79 = ns.Waypoint("WackyStart9")
	wp80 = ns.Waypoint("WackyStart10")
	wp81 = ns.Waypoint("WackyStart11")
	wp82 = ns.Waypoint("WackyStart12")
	wp83 = ns.Waypoint("WackyStart13")
	wp84 = ns.Waypoint("WackyStart14")
	wp85 = ns.Waypoint("WackyStart15")
	wp86 = ns.Waypoint("WackyStart16")
	wp87 = ns.Waypoint("WackyStart17")
	wp88 = ns.Waypoint("WackyStart18")
	wp89 = ns.Waypoint("WackyStart19")
	wp90 = ns.Waypoint("WackyStart20")
	wp91 = ns.Waypoint("WackyStart21")
	wp92 = ns.Waypoint("WackyStart22")
	wp93 = ns.Waypoint("WHome1")
	wp94 = ns.Waypoint("WHome2")
	wp95 = ns.Waypoint("WHome3")
	wp96 = ns.Waypoint("WHome4")
	wp97 = ns.Waypoint("WHome5")
	wp98 = ns.Waypoint("WHome6")
	wp99 = ns.Waypoint("WHome7")
	wp100 = ns.Waypoint("WHome8")
	wp101 = ns.Waypoint("WHome9")
	wp102 = ns.Waypoint("WHome10")
	wp103 = ns.Waypoint("WHome11")
	wp104 = ns.Waypoint("WHome12")
	wp105 = ns.Waypoint("WHome13")
	wp106 = ns.Waypoint("WHome14")
	wp107 = ns.Waypoint("WHome15")
	wp108 = ns.Waypoint("WHome16")
	wp109 = ns.Waypoint("WHome17")
	wp110 = ns.Waypoint("WHome18")
	wp111 = ns.Waypoint("WHome19")
	wp112 = ns.Waypoint("WHome20")
	wp113 = ns.Waypoint("WHome21")
	wp114 = ns.Waypoint("WHome22")
	wp115 = ns.Waypoint("WackSound1")
	wp116 = ns.Waypoint("WackSound2")
	wp117 = ns.Waypoint("WackSound3")
	wp118 = ns.Waypoint("WackSound4")
	wp119 = ns.Waypoint("WackSound5")
	wp120 = ns.Waypoint("WackSound6")
	wp121 = ns.Waypoint("WackSound7")
	wp122 = ns.Waypoint("WackSound8")
	wp123 = ns.Waypoint("WackSound9")
	wp124 = ns.Waypoint("WackSound10")
	wp125 = ns.Waypoint("WackSound11")
	wp127 = ns.Waypoint("WackySoundWP")
}
func StartCircle() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj22)
	ns.ObjectOn(obj23)
	ns.ObjectOn(obj24)
	ns.ObjectOn(obj25)
}
func EndCircle() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar26)
}
func StartWackiness() {
	if flag126 {
		goto LABEL1
	}
	ns.MoveWaypoint(wp127, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ns.AudioEvent(ns.Clank3, wp127)
	ns.Print("GeneralPrint:MsgTrapOn")
	flag126 = true
	ns.Move(obj27, wp71)
	ns.Move(obj28, wp72)
	ns.AudioEvent(ns.SpikeBlockMove, wp115)
	gvar128 = ns.SecondTimer(30, WackyGroup11)
LABEL1:
	return
}
func WackyGroup1() {
	if !(ns.IsCaller(obj27) || ns.IsCaller(obj28)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp116)
	ns.Move(obj29, wp73)
	ns.Move(obj30, wp74)
LABEL1:
	return
}
func WackyGroup2() {
	if !(ns.IsCaller(obj29) || ns.IsCaller(obj30)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp117)
	ns.Move(obj31, wp75)
	ns.Move(obj32, wp76)
LABEL1:
	return
}
func WackyGroup3() {
	if !(ns.IsCaller(obj31) || ns.IsCaller(obj32)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp118)
	ns.Move(obj33, wp77)
	ns.Move(obj34, wp78)
LABEL1:
	return
}
func WackyGroup4() {
	if !(ns.IsCaller(obj33) || ns.IsCaller(obj34)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp119)
	ns.Move(obj35, wp79)
	ns.Move(obj36, wp80)
LABEL1:
	return
}
func WackyGroup5() {
	if !(ns.IsCaller(obj35) || ns.IsCaller(obj36)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp120)
	ns.Move(obj37, wp81)
	ns.Move(obj38, wp82)
LABEL1:
	return
}
func WackyGroup6() {
	if !(ns.IsCaller(obj37) || ns.IsCaller(obj38)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp121)
	ns.Move(obj39, wp83)
	ns.Move(obj40, wp84)
LABEL1:
	return
}
func WackyGroup7() {
	if !(ns.IsCaller(obj39) || ns.IsCaller(obj40)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp122)
	ns.Move(obj41, wp85)
	ns.Move(obj42, wp86)
LABEL1:
	return
}
func WackyGroup8() {
	if !(ns.IsCaller(obj41) || ns.IsCaller(obj42)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp123)
	ns.Move(obj43, wp87)
	ns.Move(obj44, wp88)
LABEL1:
	return
}
func WackyGroup9() {
	if !(ns.IsCaller(obj43) || ns.IsCaller(obj44)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp124)
	ns.Move(obj45, wp89)
	ns.Move(obj46, wp90)
LABEL1:
	return
}
func WackyGroup10() {
	if !(ns.IsCaller(obj45) || ns.IsCaller(obj46)) {
		goto LABEL1
	}
	ns.AudioEvent(ns.SpikeBlockMove, wp125)
	ns.Move(obj47, wp91)
	ns.Move(obj48, wp92)
LABEL1:
	return
}
func WackyGroup11() {
	ns.CancelTimer(gvar128)
	ns.AudioEvent(ns.SpikeBlockMove, wp115)
	ns.AudioEvent(ns.SpikeBlockMove, wp116)
	ns.AudioEvent(ns.SpikeBlockMove, wp117)
	ns.AudioEvent(ns.SpikeBlockMove, wp118)
	ns.AudioEvent(ns.SpikeBlockMove, wp119)
	ns.AudioEvent(ns.SpikeBlockMove, wp120)
	ns.AudioEvent(ns.SpikeBlockMove, wp121)
	ns.AudioEvent(ns.SpikeBlockMove, wp122)
	ns.AudioEvent(ns.SpikeBlockMove, wp123)
	ns.AudioEvent(ns.SpikeBlockMove, wp124)
	ns.AudioEvent(ns.SpikeBlockMove, wp125)
	ns.Move(obj27, wp93)
	ns.Move(obj28, wp94)
	ns.Move(obj29, wp95)
	ns.Move(obj30, wp96)
	ns.Move(obj31, wp97)
	ns.Move(obj32, wp98)
	ns.Move(obj33, wp99)
	ns.Move(obj34, wp100)
	ns.Move(obj35, wp101)
	ns.Move(obj36, wp102)
	ns.Move(obj37, wp103)
	ns.Move(obj38, wp104)
	ns.Move(obj39, wp105)
	ns.Move(obj40, wp106)
	ns.Move(obj41, wp107)
	ns.Move(obj42, wp108)
	ns.Move(obj43, wp109)
	ns.Move(obj44, wp110)
	ns.Move(obj45, wp111)
	ns.Move(obj46, wp112)
	ns.Move(obj47, wp113)
	ns.Move(obj48, wp114)
	flag126 = false
}
func FinaleGets() {
	gvar129 = ns.ObjectGroup("FinaleTriggers")
	gvar130 = ns.WallGroup("FinaleWalls")
	gvar131 = ns.WallGroup("NorthFinale")
	gvar132 = ns.WallGroup("MiddleFinale")
	gvar133 = ns.WallGroup("SouthFinale")
	gvar137 = ns.Wall(16, 208)
	gvar138 = ns.Wall(17, 207)
	gvar139 = ns.Wall(18, 206)
	gvar140 = ns.Wall(19, 205)
	gvar141 = ns.Wall(20, 204)
}
func StartFinale() {
	ns.ObjectGroupOff(gvar129)
}
func SwapFinaleWalls() {
	if flag136 {
		goto LABEL1
	}
	ns.WallGroupOpen(gvar130)
	flag136 = true
	goto LABEL2
LABEL1:
	ns.WallGroupClose(gvar130)
	flag136 = false
LABEL2:
	gvar134 = ns.SecondTimer(5, SwapFinaleWalls)
}
func NorthFinaleButton() {
	ns.ObjectOff(ns.GetTrigger())
	ivar135 += 1
	ns.WallGroupOpen(gvar131)
	if !(ivar135 >= 3) {
		goto LABEL1
	}
	ns.CancelTimer(gvar134)
	ns.WallGroupOpen(gvar130)
LABEL1:
	return
}
func MiddleFinaleButton() {
	ns.ObjectOff(ns.GetTrigger())
	ivar135 += 1
	ns.WallGroupOpen(gvar132)
	if !(ivar135 >= 3) {
		goto LABEL1
	}
	ns.CancelTimer(gvar134)
	ns.WallGroupOpen(gvar130)
LABEL1:
	return
}
func SouthFinaleButton() {
	ns.ObjectOff(ns.GetTrigger())
	ivar135 += 1
	ns.WallGroupOpen(gvar133)
	if !(ivar135 >= 3) {
		goto LABEL1
	}
	ns.CancelTimer(gvar134)
	ns.WallGroupOpen(gvar130)
LABEL1:
	return
}
func OpenExit3() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar137)
	ns.WallOpen(gvar138)
	ns.WallOpen(gvar139)
	ns.WallOpen(gvar140)
	ns.WallOpen(gvar141)
}
func GiveFireProtect() {
	ns.Enchant(ns.GetCaller(), ns.ENCHANT_PROTECT_FROM_FIRE, 60)
	ns.Print("GeneralPrint:WishingWellPrintEnchantment")
}
func RevealFirstSkullGun() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar146)
}
func StartPusherBlock() {
	ns.ObjectOn(obj143)
}
func StopPusherBlock() {
	ns.ObjectOff(obj143)
}
func PutAwayPusherBlock() {
	if flag200 {
		goto LABEL1
	}
	ns.ObjectOff(obj144)
	ns.ObjectOff(obj143)
	ns.Move(obj142, wp145)
	flag200 = true
LABEL1:
	if !ns.IsLocked(obj196) {
		goto LABEL2
	}
	ns.UnlockDoor(obj196)
	ns.UnlockDoor(obj197)
	ns.UnlockDoor(obj198)
	ns.UnlockDoor(obj199)
LABEL2:
	return
}
func TriggerMagicWalls1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar147)
	ns.WallOpen(gvar148)
	ns.WallOpen(gvar149)
	ns.WallOpen(gvar150)
	ns.WallOpen(gvar151)
	ns.Print("GeneralPrint:MsgWallLowered")
}
func OpenMassiveSecret() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar153)
}
func OpenFrontFenceWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar154)
	ns.Print("GeneralPrint:MsgWallLowered")
}
func OpenBackFenceWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar155)
	ns.WallGroupOpen(gvar226)
	ns.Print("GeneralPrint:MsgWallLowered")
}
func OpenSoulGateTrap() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar156)
	ns.WallOpen(gvar157)
	ns.WallOpen(gvar158)
	ns.WallOpen(gvar159)
	ns.WallOpen(gvar160)
}
func BarrelRelease() {
	ivar162 += 1
	if !(ivar162 >= 3) {
		goto LABEL1
	}
	ns.WallGroupOpen(gvar161)
LABEL1:
	return
}
func OpenGoodiePits() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.ChangeSpellbar, wp195)
	ns.ObjectGroupOn(gvar163)
}
func OpenDemonWalls() {
	ivar167 += 1
	if !(ivar167 >= 2) {
		goto LABEL1
	}
	ns.WallGroupBreak(gvar164)
LABEL1:
	return
}
func OpenColumnWalls() {
	ns.ObjectGroupOff(gvar169)
	ns.WallOpen(gvar175)
	ns.WallOpen(gvar176)
	ns.WallOpen(gvar177)
	ns.WallOpen(gvar178)
	ns.WallOpen(gvar179)
	ns.WallOpen(gvar180)
	ns.WallOpen(gvar181)
	ns.WallOpen(gvar182)
	ns.WallOpen(gvar183)
}
func OpenDemonPit() {
	ns.ObjectOn(obj184)
}
func OpenDemonSecret() {
	ns.ObjectOff(ns.GetTrigger())
	if !ns.IsObjectOn(obj184) {
		goto LABEL1
	}
	ns.WallOpen(gvar165)
	ns.WallOpen(gvar166)
	ns.Print("GeneralPrint:MsgWallLowered")
	goto LABEL2
LABEL1:
	ns.ObjectOn(obj184)
LABEL2:
	return
}
func StartGenRow2A() {
	ns.ObjectOn(obj188)
}
func StartGenRow2B() {
	ns.ObjectOn(obj189)
}
func StartGenRow2C() {
	ns.ObjectOn(obj190)
}
func StartGenRow3A() {
	ns.ObjectOn(obj191)
	ns.ObjectOn(obj189)
}
func StartGenRow3B() {
	ns.ObjectOn(obj192)
	ns.ObjectOn(obj188)
	ns.ObjectOn(obj190)
}
func StartGenRow3C() {
	ns.ObjectOn(obj193)
	ns.ObjectOn(obj189)
}
func StopSkullFix() {
	ns.ObjectOff(obj194)
}
func StartClearingSpike() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj201)
	ns.ObjectOn(obj202)
}
func BreakClearWalls1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar205)
	ns.SecondTimer(2, BreakClearWalls2)
}
func BreakClearWalls2() {
	ns.WallGroupBreak(gvar206)
}
func BreakClearWalls3() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar207)
	ns.SecondTimer(2, BreakClearWalls4)
}
func BreakClearWalls4() {
	ns.WallGroupBreak(gvar208)
}
func BreakClearWalls5() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar209)
	ns.SecondTimer(2, BreakClearWalls6)
}
func BreakClearWalls6() {
	ns.WallGroupBreak(gvar210)
}
func ArrowRow1() {
	if ns.IsObjectOn(obj211) {
		goto LABEL1
	}
	ns.ObjectOn(obj211)
	ns.ObjectOn(obj212)
	ns.ObjectOn(obj213)
	ns.ObjectOn(obj214)
	ns.ObjectOn(obj215)
	ns.ObjectOn(obj216)
LABEL1:
	return
}
func ArrowRow1Release() {
	ivar217 = 1
	if !(ivar217 <= 0) {
		goto LABEL1
	}
	ns.ObjectOff(obj211)
	ns.ObjectOff(obj212)
LABEL1:
	return
}
func ArrowRow2() {
	ivar218 = 1
	if !(ivar218 > 0) {
		goto LABEL1
	}
	ns.ObjectOn(obj213)
	ns.ObjectOn(obj214)
LABEL1:
	return
}
func ArrowRow2Release() {
	ivar218 = 1
	if !(ivar218 <= 0) {
		goto LABEL1
	}
	ns.ObjectOff(obj213)
	ns.ObjectOff(obj214)
LABEL1:
	return
}
func ArrowRow3() {
	ivar219 = 1
	if !(ivar219 > 0) {
		goto LABEL1
	}
	ns.ObjectOn(obj215)
	ns.ObjectOn(obj216)
LABEL1:
	return
}
func ArrowRow3Release() {
	ivar219 = 1
	if !(ivar219 <= 0) {
		goto LABEL1
	}
	ns.ObjectOff(obj215)
	ns.ObjectOff(obj216)
LABEL1:
	return
}
func StartLavaBlocks() {
	ns.ObjectGroupOff(gvar223)
	ns.ObjectOn(obj220)
	ns.ObjectOn(obj221)
	ns.ObjectOn(obj222)
}
func ConWalls1Plate() {
	ns.WallGroupOpen(gvar224)
}
func ConWalls2Plate() {
	ns.WallGroupOpen(gvar225)
}
func MonstersGets() {
	gvar227 = ns.WallGroup("MonsterWalls1")
	gvar228 = ns.WallGroup("MonsterWalls2")
	gvar229 = ns.WallGroup("MonsterWalls3")
	gvar230 = ns.WallGroup("MonsterWalls4")
	gvar232 = ns.Wall(216, 144)
}
func ShooterGets() {
	obj233 = ns.Object("ShootingBlock1")
	obj234 = ns.Object("ShootingBlock2")
	obj235 = ns.Object("ShootingBlock3")
	obj236 = ns.Object("ShootingBlockMover1")
	obj237 = ns.Object("ShootingBlockMover2")
	obj238 = ns.Object("ShootingBlockMover3")
	gvar239 = ns.ObjectGroup("ShootingBlockTriggers")
}
func SlammerGets() {
	obj240 = ns.Object("Slammer1")
	obj241 = ns.Object("Slammer2")
	obj242 = ns.Object("SlammerMover1")
	obj243 = ns.Object("SlammerMover2")
	gvar244 = ns.ObjectGroup("SlammerGroup1")
	wp245 = ns.Waypoint("Slammer1Home")
	wp246 = ns.Waypoint("Slammer2Home")
	obj247 = ns.Object("Slammer3")
	obj248 = ns.Object("Slammer4")
	obj249 = ns.Object("SlammerMover3")
	obj250 = ns.Object("SlammerMover4")
	gvar251 = ns.ObjectGroup("SlammerGroup2")
	wp252 = ns.Waypoint("Slammer3Home")
	wp253 = ns.Waypoint("Slammer4Home")
}
func SouthWayGets() {
	obj254 = ns.Object("SouthBlock1")
	obj255 = ns.Object("SouthBlock2")
	obj256 = ns.Object("SouthBlockMover1")
	obj257 = ns.Object("SouthBlockMover2")
	obj258 = ns.Object("SouthHallDoor")
	gvar259 = ns.ObjectGroup("SouthTriggers")
	gvar260 = ns.Wall(103, 125)
	gvar226 = ns.WallGroup("SouthFenceWalls")
	gvar261 = ns.WallGroup("SouthTeleportWalls")
	gvar262 = ns.WallGroup("SwapWallsStart")
	gvar263 = ns.WallGroup("SwapWallsEnd")
	gvar264 = ns.WallGroup("SwapWallsExit")
	gvar265 = ns.WallGroup("SouthFenceSet")
	ns.LockDoor(obj258)
}
func SpikesGets() {
	gvar266 = ns.ObjectGroup("SpikeGroup1")
	gvar267 = ns.ObjectGroup("SpikeGroup2")
	gvar268 = ns.ObjectGroup("SpikeTriggers")
}
func StarCastGets() {
	obj270 = ns.Object("StarCastleMover1")
	obj271 = ns.Object("StarCastleMover2")
	obj272 = ns.Object("StarCastleMover3")
	obj273 = ns.Object("StarCastleMover4")
	gvar274 = ns.Wall(66, 118)
	gvar275 = ns.Wall(67, 119)
	gvar276 = ns.Wall(68, 120)
	gvar277 = ns.Wall(69, 121)
	gvar278 = ns.Wall(70, 122)
	obj279 = ns.Object("BFBMover1")
	obj280 = ns.Object("BFBMover2")
	obj281 = ns.Object("BFBMover3")
	obj282 = ns.Object("BFBMover4")
	obj283 = ns.Object("BFBMover5")
	obj284 = ns.Object("BFBMover6")
	obj285 = ns.Object("BFBTrigger")
	obj286 = ns.Object("GeneratorMover")
	gvar287 = ns.ObjectGroup("MovingGenTriggers")
	obj288 = ns.Object("GenMover2")
	obj289 = ns.Object("GenMover3")
}
func StatueGets() {
	obj290 = ns.Object("North1")
	obj292 = ns.Object("South1")
	obj291 = ns.Object("West1")
	obj293 = ns.Object("East1")
	wp294 = ns.Waypoint("North1Home")
	wp295 = ns.Waypoint("South1Home")
	wp297 = ns.Waypoint("East1Home")
	wp296 = ns.Waypoint("West1Home")
	obj298 = ns.Object("North2")
	obj300 = ns.Object("South2")
	obj299 = ns.Object("West2")
	obj301 = ns.Object("East2")
	wp302 = ns.Waypoint("North2Home")
	wp303 = ns.Waypoint("South2Home")
	wp305 = ns.Waypoint("East2Home")
	wp304 = ns.Waypoint("West2Home")
	obj306 = ns.Object("North3")
	obj308 = ns.Object("South3")
	obj307 = ns.Object("West3")
	obj309 = ns.Object("East3")
	wp310 = ns.Waypoint("North3Home")
	wp311 = ns.Waypoint("South3Home")
	wp313 = ns.Waypoint("East3Home")
	wp312 = ns.Waypoint("West3Home")
	wp314 = ns.Waypoint("Group1WP")
	wp315 = ns.Waypoint("Group2WP")
	wp316 = ns.Waypoint("Group3WP")
	gvar317 = ns.WallGroup("EastOopsWalls")
	gvar319 = ns.WallGroup("WestOopsWalls")
	gvar318 = ns.WallGroup("SouthOopsWalls")
	gvar320 = ns.WallGroup("TadaWalls")
	obj330 = ns.Object("Group1Light")
	obj331 = ns.Object("Group2Light")
	obj332 = ns.Object("Group3Light")
	obj327 = ns.Object("Switch1")
	obj328 = ns.Object("Switch2")
	obj329 = ns.Object("Switch3")
	obj333 = ns.Object("StudGen1")
	obj334 = ns.Object("StudGen2")
	gvar335 = ns.Wall(19, 209)
	gvar336 = ns.Wall(20, 208)
	gvar337 = ns.Wall(21, 207)
	gvar338 = ns.Wall(22, 206)
	gvar339 = ns.Wall(23, 205)
	wp341 = ns.Waypoint("StatueSoundWP")
	ns.ObjectOff(obj331)
	ns.ObjectOff(obj332)
}
func TimeRoomGets() {
	obj196 = ns.Object("TimerDoor1")
	obj197 = ns.Object("TimerDoor2")
	obj198 = ns.Object("TimerDoor3")
	obj199 = ns.Object("TimerDoor4")
	obj342 = ns.Object("TimerPlate")
	wp344 = ns.Waypoint("TickWP1")
	wp345 = ns.Waypoint("TickWP2")
	wp346 = ns.Waypoint("TickWP3")
	wp347 = ns.Waypoint("TickWP4")
	wp348 = ns.Waypoint("TickWP5")
	gvar343 = ns.WallGroup("TimerRoomWalls")
	ns.LockDoor(obj196)
	ns.LockDoor(obj197)
	ns.LockDoor(obj198)
	ns.LockDoor(obj199)
}
func MapInitialize() {
	obj142 = ns.Object("PusherBlock")
	obj143 = ns.Object("PusherBlockMover")
	obj144 = ns.Object("PusherBlockTrigger")
	wp145 = ns.Waypoint("PusherBlockHomeWP")
	gvar146 = ns.Wall(86, 116)
	gvar147 = ns.Wall(50, 190)
	gvar148 = ns.Wall(51, 189)
	gvar149 = ns.Wall(52, 188)
	gvar150 = ns.Wall(53, 187)
	gvar151 = ns.Wall(54, 186)
	gvar153 = ns.WallGroup("MassiveSecretWalls")
	gvar154 = ns.WallGroup("FrontFenceWalls")
	gvar155 = ns.WallGroup("BackFenceWalls")
	gvar156 = ns.WallGroup("SoulGateTrapWalls")
	gvar157 = ns.Wall(153, 129)
	gvar158 = ns.Wall(154, 130)
	gvar159 = ns.Wall(155, 131)
	gvar160 = ns.Wall(156, 132)
	gvar161 = ns.WallGroup("BarrelWalls")
	gvar163 = ns.ObjectGroup("GoodiePits")
	gvar164 = ns.WallGroup("DemonWalls1")
	gvar165 = ns.Wall(177, 11)
	gvar166 = ns.Wall(178, 12)
	obj168 = ns.Object("ColumnBlockMover")
	gvar169 = ns.ObjectGroup("ColumnTriggers")
	wp170 = ns.Waypoint("ColumnBlockHome")
	gvar171 = ns.WallGroup("ColumnWalls1")
	gvar172 = ns.WallGroup("ColumnWalls2")
	gvar173 = ns.WallGroup("ColumnWalls3")
	gvar174 = ns.WallGroup("ColumnWalls4")
	gvar175 = ns.Wall(26, 100)
	gvar176 = ns.Wall(27, 99)
	gvar177 = ns.Wall(28, 98)
	gvar178 = ns.Wall(38, 104)
	gvar179 = ns.Wall(39, 103)
	gvar180 = ns.Wall(40, 102)
	gvar181 = ns.Wall(34, 92)
	gvar182 = ns.Wall(35, 91)
	gvar183 = ns.Wall(36, 90)
	obj184 = ns.Object("DemonPit")
	obj194 = ns.Object("SkullFix")
	wp195 = ns.Waypoint("PitSoundWP")
	ns.ObjectOn(obj168)
	obj201 = ns.Object("ClearingSpikeMover1")
	obj202 = ns.Object("ClearingSpikeMover2")
	obj203 = ns.Object("ClearingSpike1")
	obj204 = ns.Object("ClearingSpike2")
	gvar205 = ns.WallGroup("ClearWalls1")
	gvar206 = ns.WallGroup("ClearWalls2")
	gvar207 = ns.WallGroup("ClearWalls3")
	gvar208 = ns.WallGroup("ClearWalls4")
	gvar209 = ns.WallGroup("ClearWalls5")
	gvar210 = ns.WallGroup("ClearWalls6")
	obj211 = ns.Object("Arrow1")
	obj212 = ns.Object("Arrow2")
	obj213 = ns.Object("Arrow3")
	obj214 = ns.Object("Arrow4")
	obj215 = ns.Object("Arrow5")
	obj216 = ns.Object("Arrow6")
	obj220 = ns.Object("LavaBlockMover1")
	obj221 = ns.Object("LavaBlockMover2")
	obj222 = ns.Object("LavaBlockMover3")
	gvar223 = ns.ObjectGroup("LavaBlockTriggers")
	gvar224 = ns.WallGroup("ConWalls1")
	gvar225 = ns.WallGroup("ConWalls2")
	obj185 = ns.Object("GenRow1A")
	obj186 = ns.Object("GenRow1B")
	obj187 = ns.Object("GenRow1C")
	obj188 = ns.Object("GenRow2A")
	obj189 = ns.Object("GenRow2B")
	obj190 = ns.Object("GenRow2C")
	obj191 = ns.Object("GenRow3A")
	obj192 = ns.Object("GenRow3B")
	obj193 = ns.Object("GenRow3C")
	ns.ObjectOff(obj188)
	ns.ObjectOff(obj189)
	ns.ObjectOff(obj190)
	ns.ObjectOff(obj191)
	ns.ObjectOff(obj192)
	ns.ObjectOff(obj193)
	BuzzSawGets()
	CircleGets()
	FinaleGets()
	MonstersGets()
	ShooterGets()
	SlammerGets()
	SouthWayGets()
	SpikesGets()
	StarCastGets()
	StatueGets()
	TimeRoomGets()
}
func MonsterRelease() {
	var v0 int
	if !(ivar231 <= 4) {
		goto LABEL1
	}
	ivar231 += 1
LABEL1:
	switch v0 = ivar231; v0 {
	case 1:
		goto LABEL2
	case 2:
		goto LABEL3
	case 3:
		goto LABEL4
	case 4:
		goto LABEL5
	default:
		goto LABEL6
	}
LABEL2:
	ns.WallGroupOpen(gvar227)
	goto LABEL6
LABEL3:
	ns.WallGroupOpen(gvar228)
	goto LABEL6
LABEL4:
	ns.WallGroupOpen(gvar229)
	goto LABEL6
LABEL5:
	ns.WallGroupOpen(gvar230)
	goto LABEL6
LABEL6:
	return
}
func OpenMonsterExit() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar232)
}
func StartShootingBlocks() {
	ns.ObjectGroupOff(gvar239)
	ns.ObjectOn(obj236)
	ns.SecondTimer(3, ShootBlock2)
}
func ShootBlock1() {
	ns.ObjectOn(obj236)
	ns.SecondTimer(3, ShootBlock2)
}
func ShootBlock2() {
	ns.ObjectOn(obj237)
	ns.SecondTimer(3, ShootBlock3)
}
func ShootBlock3() {
	ns.ObjectOn(obj238)
	ns.SecondTimer(3, ShootBlock1)
}
func NoShootBlock1() {
	ns.ObjectOff(obj236)
}
func NoShootBlock2() {
	ns.ObjectOff(obj237)
}
func NoShootBlock3() {
	ns.ObjectOff(obj238)
}
func StartSlammerGroup1() {
	ns.ObjectOn(obj242)
	ns.ObjectOn(obj243)
}
func SlammerGroup1Wait() {
	ns.ObjectGroupOff(gvar244)
	ns.SecondTimer(3, SlammerGroup1Retreat)
}
func SlammerGroup1Retreat() {
	ns.ObjectOff(obj242)
	ns.ObjectOff(obj243)
	ns.Move(obj240, wp245)
	ns.Move(obj241, wp246)
}
func ResetSlammerGroup1() {
	ns.ObjectOff(obj242)
	ns.ObjectOff(obj243)
	ns.ObjectGroupOn(gvar244)
}
func StartSlammerGroup2() {
	ns.ObjectOn(obj249)
	ns.ObjectOn(obj250)
}
func SlammerGroup2Wait() {
	ns.ObjectGroupOff(gvar251)
	ns.SecondTimer(3, SlammerGroup2Retreat)
}
func SlammerGroup2Retreat() {
	ns.ObjectOff(obj249)
	ns.ObjectOff(obj250)
	ns.Move(obj247, wp252)
	ns.Move(obj248, wp253)
}
func ResetSlammerGroup2() {
	ns.ObjectOff(obj249)
	ns.ObjectOff(obj250)
	ns.ObjectGroupOn(gvar251)
}
func StartSouthBlocks() {
	ns.ObjectGroupOff(gvar259)
	ns.WallOpen(gvar260)
	ns.WallGroupOpen(gvar261)
	ns.WallGroupClose(gvar226)
	ns.ObjectOn(obj256)
	ns.ObjectOn(obj257)
}
func StopSouthBlocks() {
	if !ns.IsCaller(obj255) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOff(obj256)
	ns.ObjectOff(obj257)
LABEL1:
	return
}
func StartSwapRoom() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar262)
	ns.WallGroupClose(gvar263)
}
func LowerDividerWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar263)
}
func LowerExitWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar264)
}
func OpenSouthFenceSet() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar265)
}
func OpenSouthHallDoor() {
	ns.ObjectOff(ns.GetTrigger())
	ns.UnlockDoor(obj258)
	ns.Print("GeneralPrint:MsgDoorUnlock")
}
func StartSpikes() {
	ns.ObjectGroupOff(gvar268)
	gvar269 = ns.SecondTimer(3, SpikeGroup2Down)
}
func SpikeGroup1Up() {
	ns.ObjectGroupOn(gvar266)
	gvar269 = ns.SecondTimer(3, SpikeGroup1Down)
}
func SpikeGroup1Down() {
	ns.ObjectGroupOff(gvar266)
	gvar269 = ns.SecondTimer(2, SpikeGroup2Up)
}
func SpikeGroup2Up() {
	ns.ObjectGroupOn(gvar267)
	gvar269 = ns.SecondTimer(3, SpikeGroup2Down)
}
func SpikeGroup2Down() {
	ns.ObjectGroupOff(gvar267)
	gvar269 = ns.SecondTimer(2, SpikeGroup1Up)
}
func StartStarCastle() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar274)
	ns.WallOpen(gvar275)
	ns.WallOpen(gvar276)
	ns.WallOpen(gvar277)
	ns.WallOpen(gvar278)
	ns.ObjectOn(obj270)
	ns.ObjectOn(obj271)
	ns.ObjectOn(obj272)
	ns.ObjectOn(obj273)
}
func EndStarCastle() {
	ns.ObjectOff(obj270)
	ns.ObjectOff(obj271)
	ns.ObjectOff(obj272)
	ns.ObjectOff(obj273)
}
func StartBFB() {
	ns.ObjectOff(obj285)
	ns.ObjectOn(obj279)
	ns.ObjectOn(obj280)
	ns.ObjectOn(obj281)
	ns.ObjectOn(obj282)
	ns.ObjectOn(obj283)
	ns.ObjectOn(obj284)
}
func StartMovingGen() {
	ns.ObjectGroupOff(gvar287)
	ns.ObjectOn(obj286)
	ns.ObjectOn(obj288)
	ns.ObjectOn(obj289)
}
func StopMovingGen() {
	ns.ObjectOff(obj286)
}
func StopGenMover2() {
	ns.ObjectOff(obj288)
}
func StopGenMover3() {
	ns.ObjectOff(obj289)
}
func GenGone() {
	ivar340 += 1
	if ivar340 != 2 {
		goto LABEL1
	}
	ns.WallOpen(gvar335)
	ns.WallOpen(gvar336)
	ns.WallOpen(gvar337)
	ns.WallOpen(gvar338)
	ns.WallOpen(gvar339)
LABEL1:
	return
}
func Group1Lever() {
	var v0 int
	ivar321 += 1
	if !(ivar321 > 4) {
		goto LABEL1
	}
	ivar321 = 1
LABEL1:
	if flag324 {
		goto LABEL2
	}
	flag324 = true
LABEL2:
	switch v0 = ivar321; v0 {
	case 1:
		goto LABEL3
	case 2:
		goto LABEL4
	case 3:
		goto LABEL5
	case 4:
		goto LABEL6
	default:
		goto LABEL7
	}
LABEL3:
	ns.MoveObject(obj290, ns.GetWaypointX(wp314), ns.GetWaypointY(wp314))
	ns.MoveObject(obj293, ns.GetWaypointX(wp297), ns.GetWaypointY(wp297))
	ns.MoveObject(obj292, ns.GetWaypointX(wp295), ns.GetWaypointY(wp295))
	ns.MoveObject(obj291, ns.GetWaypointX(wp296), ns.GetWaypointY(wp296))
	ns.ObjectOff(obj330)
	goto LABEL7
LABEL4:
	ns.MoveObject(obj293, ns.GetWaypointX(wp314), ns.GetWaypointY(wp314))
	ns.MoveObject(obj290, ns.GetWaypointX(wp294), ns.GetWaypointY(wp294))
	ns.MoveObject(obj292, ns.GetWaypointX(wp295), ns.GetWaypointY(wp295))
	ns.MoveObject(obj291, ns.GetWaypointX(wp296), ns.GetWaypointY(wp296))
	ns.ObjectOn(obj330)
	goto LABEL7
LABEL5:
	ns.MoveObject(obj292, ns.GetWaypointX(wp314), ns.GetWaypointY(wp314))
	ns.MoveObject(obj293, ns.GetWaypointX(wp297), ns.GetWaypointY(wp297))
	ns.MoveObject(obj290, ns.GetWaypointX(wp294), ns.GetWaypointY(wp294))
	ns.MoveObject(obj291, ns.GetWaypointX(wp296), ns.GetWaypointY(wp296))
	ns.ObjectOff(obj330)
	goto LABEL7
LABEL6:
	ns.MoveObject(obj291, ns.GetWaypointX(wp314), ns.GetWaypointY(wp314))
	ns.MoveObject(obj293, ns.GetWaypointX(wp297), ns.GetWaypointY(wp297))
	ns.MoveObject(obj292, ns.GetWaypointX(wp295), ns.GetWaypointY(wp295))
	ns.MoveObject(obj290, ns.GetWaypointX(wp294), ns.GetWaypointY(wp294))
	ns.ObjectOff(obj330)
	goto LABEL7
LABEL7:
	OpenOopsWalls()
}
func Group2Lever() {
	var v0 int
	ivar322 += 1
	if !(ivar322 > 4) {
		goto LABEL1
	}
	ivar322 = 1
LABEL1:
	if flag325 {
		goto LABEL2
	}
	flag325 = true
LABEL2:
	switch v0 = ivar322; v0 {
	case 1:
		goto LABEL3
	case 2:
		goto LABEL4
	case 3:
		goto LABEL5
	case 4:
		goto LABEL6
	default:
		goto LABEL7
	}
LABEL3:
	ns.MoveObject(obj298, ns.GetWaypointX(wp315), ns.GetWaypointY(wp315))
	ns.MoveObject(obj301, ns.GetWaypointX(wp305), ns.GetWaypointY(wp305))
	ns.MoveObject(obj300, ns.GetWaypointX(wp303), ns.GetWaypointY(wp303))
	ns.MoveObject(obj299, ns.GetWaypointX(wp304), ns.GetWaypointY(wp304))
	ns.ObjectOn(obj331)
	goto LABEL7
LABEL4:
	ns.MoveObject(obj301, ns.GetWaypointX(wp315), ns.GetWaypointY(wp315))
	ns.MoveObject(obj298, ns.GetWaypointX(wp302), ns.GetWaypointY(wp302))
	ns.MoveObject(obj300, ns.GetWaypointX(wp303), ns.GetWaypointY(wp303))
	ns.MoveObject(obj299, ns.GetWaypointX(wp304), ns.GetWaypointY(wp304))
	ns.ObjectOff(obj331)
	goto LABEL7
LABEL5:
	ns.MoveObject(obj300, ns.GetWaypointX(wp315), ns.GetWaypointY(wp315))
	ns.MoveObject(obj301, ns.GetWaypointX(wp305), ns.GetWaypointY(wp305))
	ns.MoveObject(obj298, ns.GetWaypointX(wp302), ns.GetWaypointY(wp302))
	ns.MoveObject(obj299, ns.GetWaypointX(wp304), ns.GetWaypointY(wp304))
	ns.ObjectOff(obj331)
	goto LABEL7
LABEL6:
	ns.MoveObject(obj299, ns.GetWaypointX(wp315), ns.GetWaypointY(wp315))
	ns.MoveObject(obj301, ns.GetWaypointX(wp305), ns.GetWaypointY(wp305))
	ns.MoveObject(obj300, ns.GetWaypointX(wp303), ns.GetWaypointY(wp303))
	ns.MoveObject(obj298, ns.GetWaypointX(wp302), ns.GetWaypointY(wp302))
	ns.ObjectOff(obj331)
	goto LABEL7
LABEL7:
	OpenOopsWalls()
}
func Group3Lever() {
	var v0 int
	ivar323 += 1
	if !(ivar323 > 4) {
		goto LABEL1
	}
	ivar323 = 1
LABEL1:
	if flag326 {
		goto LABEL2
	}
	flag326 = true
LABEL2:
	switch v0 = ivar323; v0 {
	case 1:
		goto LABEL3
	case 2:
		goto LABEL4
	case 3:
		goto LABEL5
	case 4:
		goto LABEL6
	default:
		goto LABEL7
	}
LABEL3:
	ns.MoveObject(obj306, ns.GetWaypointX(wp316), ns.GetWaypointY(wp316))
	ns.MoveObject(obj309, ns.GetWaypointX(wp313), ns.GetWaypointY(wp313))
	ns.MoveObject(obj308, ns.GetWaypointX(wp311), ns.GetWaypointY(wp311))
	ns.MoveObject(obj307, ns.GetWaypointX(wp312), ns.GetWaypointY(wp312))
	ns.ObjectOff(obj332)
	goto LABEL7
LABEL4:
	ns.MoveObject(obj309, ns.GetWaypointX(wp316), ns.GetWaypointY(wp316))
	ns.MoveObject(obj306, ns.GetWaypointX(wp310), ns.GetWaypointY(wp310))
	ns.MoveObject(obj308, ns.GetWaypointX(wp311), ns.GetWaypointY(wp311))
	ns.MoveObject(obj307, ns.GetWaypointX(wp312), ns.GetWaypointY(wp312))
	ns.ObjectOff(obj332)
	goto LABEL7
LABEL5:
	ns.MoveObject(obj308, ns.GetWaypointX(wp316), ns.GetWaypointY(wp316))
	ns.MoveObject(obj309, ns.GetWaypointX(wp313), ns.GetWaypointY(wp313))
	ns.MoveObject(obj306, ns.GetWaypointX(wp310), ns.GetWaypointY(wp310))
	ns.MoveObject(obj307, ns.GetWaypointX(wp312), ns.GetWaypointY(wp312))
	ns.ObjectOff(obj332)
	goto LABEL7
LABEL6:
	ns.MoveObject(obj307, ns.GetWaypointX(wp316), ns.GetWaypointY(wp316))
	ns.MoveObject(obj309, ns.GetWaypointX(wp313), ns.GetWaypointY(wp313))
	ns.MoveObject(obj308, ns.GetWaypointX(wp311), ns.GetWaypointY(wp311))
	ns.MoveObject(obj306, ns.GetWaypointX(wp310), ns.GetWaypointY(wp310))
	ns.ObjectOn(obj332)
	goto LABEL7
LABEL7:
	OpenOopsWalls()
}
func OpenOopsWalls() {
	if !(ivar321 == 2 && ivar322 == 2 && ivar323 == 2) {
		goto LABEL1
	}
	ns.WallGroupOpen(gvar317)
LABEL1:
	if !(ivar321 == 3 && ivar322 == 3 && ivar323 == 3) {
		goto LABEL2
	}
	ns.WallGroupOpen(gvar318)
LABEL2:
	if !(ivar321 == 4 && ivar322 == 4 && ivar323 == 4) {
		goto LABEL3
	}
	ns.WallGroupOpen(gvar319)
LABEL3:
	if !(ivar321 == 1 && ivar322 == 1 && ivar323 == 1) {
		goto LABEL4
	}
	ns.WallGroupOpen(gvar317)
	ns.WallGroupOpen(gvar319)
	ns.WallGroupOpen(gvar318)
LABEL4:
	if !(ivar321 == 2 && ivar322 == 1 && ivar323 == 4) {
		goto LABEL5
	}
	ns.ObjectOff(obj327)
	ns.ObjectOff(obj328)
	ns.ObjectOff(obj329)
	ns.AudioEvent(ns.Clank1, wp341)
	ns.WallGroupOpen(gvar320)
LABEL5:
	return
}
func StartTimer() {
	if ns.IsLocked(obj196) {
		goto LABEL1
	}
	if flag352 {
		goto LABEL1
	}
	ns.LockDoor(obj196)
	ns.LockDoor(obj197)
	ns.LockDoor(obj198)
	ns.LockDoor(obj199)
LABEL1:
	if ivar350 == 10 {
		goto LABEL2
	}
	if flag352 {
		goto LABEL3
	}
	ns.AudioEvent(ns.GenerateTick, wp344)
	ns.AudioEvent(ns.GenerateTick, wp344)
	ns.AudioEvent(ns.GenerateTick, wp344)
	ns.AudioEvent(ns.GenerateTick, wp344)
	ivar350 += 1
	gvar349 = ns.SecondTimer(1, StartTimer)
LABEL3:
	goto LABEL4
LABEL2:
	OpenTimerRoomWalls()
	ns.CancelTimer(gvar349)
LABEL4:
	return
}
func StopTimer() {
	ns.CancelTimer(gvar349)
	ivar350 = 0
	if !ns.IsLocked(obj196) {
		goto LABEL1
	}
	ns.UnlockDoor(obj196)
	ns.UnlockDoor(obj197)
	ns.UnlockDoor(obj198)
	ns.UnlockDoor(obj199)
LABEL1:
	return
}
func OpenTimerRoomWalls() {
	flag352 = true
	ns.WallGroupOpen(gvar343)
}
func TimerGenGone() {
	ivar351 += 1
	if !(ivar351 >= 3) {
		goto LABEL1
	}
	ns.ObjectOff(obj342)
	if !ns.IsLocked(obj196) {
		goto LABEL2
	}
	ns.UnlockDoor(obj196)
	ns.UnlockDoor(obj197)
	ns.UnlockDoor(obj198)
	ns.UnlockDoor(obj199)
LABEL2:
	OpenRoomDoors()
LABEL1:
	return
}
func OpenRoomDoors() {
	if !ns.IsLocked(obj196) {
		goto LABEL1
	}
	ns.UnlockDoor(obj196)
	ns.UnlockDoor(obj197)
	ns.UnlockDoor(obj198)
	ns.UnlockDoor(obj199)
LABEL1:
	return
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
