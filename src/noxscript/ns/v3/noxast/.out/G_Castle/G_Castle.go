package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	obj4    ns.ObjectID
	obj5    ns.ObjectID
	obj6    ns.ObjectID
	ivar7   int
	ivar8   int
	ivar9   int
	obj10   ns.ObjectID
	obj11   ns.ObjectID
	obj12   ns.ObjectID
	ivar13  int
	ivar14  int
	ivar15  int
	obj16   ns.ObjectID
	obj17   ns.ObjectID
	obj18   ns.ObjectID
	obj19   ns.ObjectID
	obj20   ns.ObjectID
	obj21   ns.ObjectID
	obj22   ns.ObjectID
	obj23   ns.ObjectID
	obj24   ns.ObjectID
	obj25   ns.ObjectID
	ivar26  int
	ivar27  int
	ivar28  int
	ivar29  int
	ivar30  int
	ivar31  int
	ivar32  int
	ivar33  int
	ivar34  int
	ivar35  int
	obj36   ns.ObjectID
	obj37   ns.ObjectID
	obj38   ns.ObjectID
	obj39   ns.ObjectID
	obj40   ns.ObjectID
	obj41   ns.ObjectID
	ivar42  int
	ivar43  int
	ivar44  int
	ivar45  int
	ivar46  int
	ivar47  int
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
	gvar59  ns.ObjectGroupID
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
	gvar91  ns.WallID
	gvar92  ns.WallID
	gvar93  ns.WallID
	gvar94  ns.WallID
	gvar95  ns.WallID
	gvar96  ns.TimerID
	gvar97  ns.TimerID
	fvar98  float32
	fvar99  float32
	flag100 bool
	flag101 bool
	flag102 bool
	wp103   ns.WaypointID
	flag104 bool
	obj105  ns.ObjectID
	obj106  ns.ObjectID
	obj107  ns.ObjectID
	gvar108 ns.WallGroupID
	wp109   ns.WaypointID
	wp110   ns.WaypointID
	gvar111 ns.WallID
	gvar112 ns.WallGroupID
	gvar113 ns.WallGroupID
	gvar114 ns.WallGroupID
	gvar115 ns.WallID
	gvar116 ns.WallID
	gvar117 ns.WallGroupID
	flag118 bool
	obj119  ns.ObjectID
	obj120  ns.ObjectID
	wp121   ns.WaypointID
	wp122   ns.WaypointID
	obj123  ns.ObjectID
	obj124  ns.ObjectID
	wp125   ns.WaypointID
	obj126  ns.ObjectID
	obj127  ns.ObjectID
	wp128   ns.WaypointID
	wp129   ns.WaypointID
	obj130  ns.ObjectID
	obj131  ns.ObjectID
	wp132   ns.WaypointID
	obj133  ns.ObjectID
	obj134  ns.ObjectID
	wp135   ns.WaypointID
	wp136   ns.WaypointID
	obj137  ns.ObjectID
	obj138  ns.ObjectID
	wp139   ns.WaypointID
	obj140  ns.ObjectID
	obj141  ns.ObjectID
	wp142   ns.WaypointID
	wp143   ns.WaypointID
	obj144  ns.ObjectID
	obj145  ns.ObjectID
	wp146   ns.WaypointID
	obj147  ns.ObjectID
	obj148  ns.ObjectID
	obj149  ns.ObjectID
	obj150  ns.ObjectID
	gvar151 ns.ObjectGroupID
	gvar152 ns.ObjectGroupID
	gvar153 ns.WallGroupID
	gvar154 ns.WallGroupID
	gvar155 ns.WallGroupID
	gvar156 ns.WallGroupID
	gvar157 ns.WallGroupID
	flag158 bool
	obj159  ns.ObjectID
	obj160  ns.ObjectID
	gvar161 ns.ObjectGroupID
	wp162   ns.WaypointID
	wp163   ns.WaypointID
	gvar164 ns.ObjectGroupID
	gvar165 ns.WallID
	gvar166 ns.WallID
	gvar167 ns.WallID
	gvar168 ns.WallID
	gvar169 ns.WallID
	gvar170 ns.WallID
	gvar171 ns.WallID
	gvar172 ns.WallID
	gvar173 ns.WallID
	gvar174 ns.WallID
	gvar175 ns.WallID
	gvar176 ns.WallID
	gvar177 ns.WallID
	gvar178 ns.WallID
	gvar179 ns.WallID
	gvar180 ns.WallID
	gvar181 ns.WallID
	gvar182 ns.WallID
	gvar183 ns.WallID
	gvar184 ns.WallID
	gvar185 ns.WallID
	gvar186 ns.WallID
	wp187   ns.WaypointID
	obj188  ns.ObjectID
	obj189  ns.ObjectID
	obj190  ns.ObjectID
	obj191  ns.ObjectID
	obj192  ns.ObjectID
	obj193  ns.ObjectID
	obj194  ns.ObjectID
	obj195  ns.ObjectID
	gvar196 int
	gvar197 int
	gvar198 int
	gvar199 int
	gvar200 int
	gvar201 int
	flag202 bool
	obj203  ns.ObjectID
	obj204  ns.ObjectID
	obj205  ns.ObjectID
	obj206  ns.ObjectID
	obj207  ns.ObjectID
	obj208  ns.ObjectID
	obj209  ns.ObjectID
	obj210  ns.ObjectID
	flag211 bool
	ivar212 int
	obj213  ns.ObjectID
	obj214  ns.ObjectID
	obj215  ns.ObjectID
	obj216  ns.ObjectID
	gvar217 ns.WallGroupID
	gvar218 ns.WallID
	gvar219 ns.WallGroupID
	gvar220 ns.WallID
	wp221   ns.WaypointID
	obj222  ns.ObjectID
	obj223  ns.ObjectID
	obj224  ns.ObjectID
	obj225  ns.ObjectID
	gvar226 ns.WallGroupID
	gvar227 ns.WallGroupID
	gvar228 ns.WallID
	gvar229 ns.WallID
	wp230   ns.WaypointID
	obj231  ns.ObjectID
	gvar232 ns.WallGroupID
	flag233 bool
	obj234  ns.ObjectID
	obj235  ns.ObjectID
	obj236  ns.ObjectID
	obj237  ns.ObjectID
	obj238  ns.ObjectID
	obj239  ns.ObjectID
	obj240  ns.ObjectID
	obj241  ns.ObjectID
	obj242  ns.ObjectID
	ivar243 int
	ivar244 int
	ivar245 int
	gvar246 ns.WallGroupID
	gvar247 ns.WallGroupID
	flag248 bool
	obj249  ns.ObjectID
	gvar250 ns.WallGroupID
	gvar251 ns.WallGroupID
	gvar252 ns.WallGroupID
	gvar253 ns.WallGroupID
	obj254  ns.ObjectID
	obj255  ns.ObjectID
	obj256  ns.ObjectID
	obj257  ns.ObjectID
	obj258  ns.ObjectID
	ivar259 int
	flag260 bool
	obj261  ns.ObjectID
	gvar262 ns.ObjectGroupID
	gvar263 ns.WallGroupID
	gvar264 ns.WallGroupID
	obj265  ns.ObjectID
	obj266  ns.ObjectID
	obj267  ns.ObjectID
	obj268  ns.ObjectID
	obj269  ns.ObjectID
	obj270  ns.ObjectID
	gvar271 ns.WallID
	gvar272 ns.WallID
	gvar273 ns.WallID
	gvar274 ns.WallID
	gvar275 ns.WallID
)

func init() {
	ivar7 = 0
	ivar8 = 0
	ivar9 = 0
	ivar13 = 0
	ivar14 = 0
	ivar15 = 0
	ivar26 = 0
	ivar27 = 0
	ivar28 = 0
	ivar29 = 0
	ivar30 = 0
	ivar31 = 0
	ivar32 = 0
	ivar33 = 0
	ivar34 = 0
	ivar35 = 0
	ivar42 = 0
	ivar43 = 0
	ivar44 = 0
	ivar45 = 0
	ivar46 = 0
	ivar47 = 0
	gvar96 = 0
	gvar97 = 0
	fvar98 = 5
	fvar99 = 5
	flag100 = true
	flag101 = false
	flag102 = false
	flag104 = true
	flag118 = false
	flag158 = false
	gvar196 = 0
	gvar197 = 1
	gvar198 = 2
	gvar199 = 3
	gvar200 = 4
	gvar201 = gvar197
	flag202 = false
	flag211 = false
	ivar212 = 1
	flag233 = false
	ivar243 = 0
	ivar244 = 0
	ivar245 = 0
	flag248 = false
	ivar259 = 0
	flag260 = false
}
func GetArrows() {
	obj4 = ns.Object("NorthArrow1")
	obj5 = ns.Object("NorthArrow2")
	obj6 = ns.Object("NorthArrow3")
	obj10 = ns.Object("SouthArrow1")
	obj11 = ns.Object("SouthArrow2")
	obj12 = ns.Object("SouthArrow3")
	obj16 = ns.Object("EastArrow1")
	obj17 = ns.Object("EastArrow2")
	obj18 = ns.Object("EastArrow3")
	obj19 = ns.Object("EastArrow4")
	obj20 = ns.Object("EastArrow5")
	obj21 = ns.Object("EastArrow6")
	obj22 = ns.Object("EastArrow7")
	obj23 = ns.Object("EastArrow8")
	obj24 = ns.Object("EastArrow9")
	obj25 = ns.Object("EastArrow10")
	obj36 = ns.Object("WestArrow1")
	obj37 = ns.Object("WestArrow2")
	obj38 = ns.Object("WestArrow3")
	obj39 = ns.Object("WestArrow4")
	obj40 = ns.Object("WestArrow5")
	obj41 = ns.Object("WestArrow6")
}
func NorthRow1Hit() {
	ivar7 += 1
	CheckRows(4)
}
func NorthRow1Release() {
	ivar7 -= 1
	CheckRows(4)
}
func NorthRow2Hit() {
	ivar8 += 1
	CheckRows(5)
}
func NorthRow2Release() {
	ivar8 -= 1
	CheckRows(5)
}
func NorthRow3Hit() {
	ivar9 += 1
	CheckRows(6)
}
func NorthRow3Release() {
	ivar9 -= 1
	CheckRows(6)
}
func SouthRow1Hit() {
	ivar13 += 1
	CheckRows(1)
}
func SouthRow1Release() {
	ivar13 -= 1
	CheckRows(1)
}
func SouthRow2Hit() {
	ivar14 += 1
	CheckRows(2)
}
func SouthRow2Release() {
	ivar14 -= 1
	CheckRows(2)
}
func SouthRow3Hit() {
	ivar15 += 1
	CheckRows(3)
}
func SouthRow3Release() {
	ivar15 -= 1
	CheckRows(3)
}
func EastRow1Hit() {
	ivar26 += 1
	CheckRows(13)
}
func EastRow1Release() {
	ivar26 -= 1
	CheckRows(13)
}
func EastRow2Hit() {
	ivar27 += 1
	CheckRows(14)
}
func EastRow2Release() {
	ivar27 -= 1
	CheckRows(14)
}
func EastRow3Hit() {
	ivar28 += 1
	CheckRows(15)
}
func EastRow3Release() {
	ivar28 -= 1
	CheckRows(15)
}
func EastRow4Hit() {
	ivar29 += 1
	CheckRows(16)
}
func EastRow4Release() {
	ivar29 -= 1
	CheckRows(16)
}
func EastRow5Hit() {
	ivar30 += 1
	CheckRows(17)
}
func EastRow5Release() {
	ivar30 -= 1
	CheckRows(17)
}
func EastRow6Hit() {
	ivar31 += 1
	CheckRows(18)
}
func EastRow6Release() {
	ivar31 -= 1
	CheckRows(18)
}
func EastRow7Hit() {
	ivar32 += 1
	CheckRows(19)
}
func EastRow7Release() {
	ivar32 -= 1
	CheckRows(19)
}
func EastRow8Hit() {
	ivar33 += 1
	CheckRows(20)
}
func EastRow8Release() {
	ivar33 -= 1
	CheckRows(20)
}
func EastRow9Hit() {
	ivar34 += 1
	CheckRows(21)
}
func EastRow9Release() {
	ivar34 -= 1
	CheckRows(21)
}
func EastRow10Hit() {
	ivar35 += 1
	CheckRows(22)
}
func EastRow10Release() {
	ivar35 -= 1
	CheckRows(22)
}
func WestRow1Hit() {
	ivar42 += 1
	CheckRows(7)
}
func WestRow1Release() {
	ivar42 -= 1
	CheckRows(7)
}
func WestRow2Hit() {
	ivar43 += 1
	CheckRows(8)
}
func WestRow2Release() {
	ivar43 -= 1
	CheckRows(8)
}
func WestRow3Hit() {
	ivar44 += 1
	CheckRows(9)
}
func WestRow3Release() {
	ivar44 -= 1
	CheckRows(9)
}
func WestRow4Hit() {
	ivar45 += 1
	CheckRows(10)
}
func WestRow4Release() {
	ivar45 -= 1
	CheckRows(10)
}
func WestRow5Hit() {
	ivar46 += 1
	CheckRows(11)
}
func WestRow5Release() {
	ivar46 -= 1
	CheckRows(11)
}
func WestRow6Hit() {
	ivar47 += 1
	CheckRows(12)
}
func WestRow6Release() {
	ivar47 -= 1
	CheckRows(12)
}
func CheckRows(a1 int) {
	var v0 int
	switch v0 = a1; v0 {
	case 1:
		goto LABEL1
	case 2:
		goto LABEL2
	case 3:
		goto LABEL3
	case 4:
		goto LABEL4
	case 5:
		goto LABEL5
	case 6:
		goto LABEL6
	case 7:
		goto LABEL7
	case 8:
		goto LABEL8
	case 9:
		goto LABEL9
	case 10:
		goto LABEL10
	case 11:
		goto LABEL11
	case 12:
		goto LABEL12
	case 13:
		goto LABEL13
	case 14:
		goto LABEL14
	case 15:
		goto LABEL15
	case 16:
		goto LABEL16
	case 17:
		goto LABEL17
	case 18:
		goto LABEL18
	case 19:
		goto LABEL19
	case 20:
		goto LABEL20
	case 21:
		goto LABEL21
	case 22:
		goto LABEL22
	default:
		goto LABEL23
	}
LABEL1:
	if !(ivar13 > 0) {
		goto LABEL24
	}
	ns.ObjectOn(obj10)
	goto LABEL25
LABEL24:
	ns.ObjectOff(obj10)
LABEL25:
	goto LABEL23
LABEL2:
	if !(ivar14 > 0) {
		goto LABEL26
	}
	ns.ObjectOn(obj11)
	goto LABEL27
LABEL26:
	ns.ObjectOff(obj11)
LABEL27:
	goto LABEL23
LABEL3:
	if !(ivar15 > 0) {
		goto LABEL28
	}
	ns.ObjectOn(obj12)
	goto LABEL29
LABEL28:
	ns.ObjectOff(obj12)
LABEL29:
	goto LABEL23
LABEL4:
	if !(ivar7 > 0) {
		goto LABEL30
	}
	ns.ObjectOn(obj4)
	goto LABEL31
LABEL30:
	ns.ObjectOff(obj4)
LABEL31:
	goto LABEL23
LABEL5:
	if !(ivar8 > 0) {
		goto LABEL32
	}
	ns.ObjectOn(obj5)
	goto LABEL33
LABEL32:
	ns.ObjectOff(obj5)
LABEL33:
	goto LABEL23
LABEL6:
	if !(ivar9 > 0) {
		goto LABEL34
	}
	ns.ObjectOn(obj6)
	goto LABEL35
LABEL34:
	ns.ObjectOff(obj6)
LABEL35:
	goto LABEL23
LABEL7:
	if !(ivar42 > 0) {
		goto LABEL36
	}
	ns.ObjectOn(obj36)
	goto LABEL37
LABEL36:
	ns.ObjectOff(obj36)
LABEL37:
	goto LABEL23
LABEL8:
	if !(ivar43 > 0) {
		goto LABEL38
	}
	ns.ObjectOn(obj37)
	goto LABEL39
LABEL38:
	ns.ObjectOff(obj37)
LABEL39:
	goto LABEL23
LABEL9:
	if !(ivar44 > 0) {
		goto LABEL40
	}
	ns.ObjectOn(obj38)
	goto LABEL41
LABEL40:
	ns.ObjectOff(obj38)
LABEL41:
	goto LABEL23
LABEL10:
	if !(ivar45 > 0) {
		goto LABEL42
	}
	ns.ObjectOn(obj39)
	goto LABEL43
LABEL42:
	ns.ObjectOff(obj39)
LABEL43:
	goto LABEL23
LABEL11:
	if !(ivar46 > 0) {
		goto LABEL44
	}
	ns.ObjectOn(obj40)
	goto LABEL45
LABEL44:
	ns.ObjectOff(obj40)
LABEL45:
	goto LABEL23
LABEL12:
	if !(ivar47 > 0) {
		goto LABEL46
	}
	ns.ObjectOn(obj41)
	goto LABEL47
LABEL46:
	ns.ObjectOff(obj41)
LABEL47:
	goto LABEL23
LABEL13:
	if !(ivar26 > 0) {
		goto LABEL48
	}
	ns.ObjectOn(obj16)
	goto LABEL49
LABEL48:
	ns.ObjectOff(obj16)
LABEL49:
	goto LABEL23
LABEL14:
	if !(ivar27 > 0) {
		goto LABEL50
	}
	ns.ObjectOn(obj17)
	goto LABEL51
LABEL50:
	ns.ObjectOff(obj17)
LABEL51:
	goto LABEL23
LABEL15:
	if !(ivar28 > 0) {
		goto LABEL52
	}
	ns.ObjectOn(obj18)
	goto LABEL53
LABEL52:
	ns.ObjectOff(obj18)
LABEL53:
	goto LABEL23
LABEL16:
	if !(ivar29 > 0) {
		goto LABEL54
	}
	ns.ObjectOn(obj19)
	goto LABEL55
LABEL54:
	ns.ObjectOff(obj19)
LABEL55:
	goto LABEL23
LABEL17:
	if !(ivar30 > 0) {
		goto LABEL56
	}
	ns.ObjectOn(obj20)
	goto LABEL57
LABEL56:
	ns.ObjectOff(obj20)
LABEL57:
	goto LABEL23
LABEL18:
	if !(ivar31 > 0) {
		goto LABEL58
	}
	ns.ObjectOn(obj21)
	goto LABEL59
LABEL58:
	ns.ObjectOff(obj21)
LABEL59:
	goto LABEL23
LABEL19:
	if !(ivar32 > 0) {
		goto LABEL60
	}
	ns.ObjectOn(obj22)
	goto LABEL61
LABEL60:
	ns.ObjectOff(obj22)
LABEL61:
	goto LABEL23
LABEL20:
	if !(ivar33 > 0) {
		goto LABEL62
	}
	ns.ObjectOn(obj23)
	goto LABEL63
LABEL62:
	ns.ObjectOff(obj23)
LABEL63:
	goto LABEL23
LABEL21:
	if !(ivar34 > 0) {
		goto LABEL64
	}
	ns.ObjectOn(obj24)
	goto LABEL65
LABEL64:
	ns.ObjectOff(obj24)
LABEL65:
	goto LABEL23
LABEL22:
	if !(ivar35 > 0) {
		goto LABEL66
	}
	ns.ObjectOn(obj25)
	goto LABEL67
LABEL66:
	ns.ObjectOff(obj25)
LABEL67:
	goto LABEL23
LABEL23:
	return
}
func CrusherGets() {
	obj48 = ns.Object("CrusherDoor1")
	obj49 = ns.Object("CrusherDoor2")
	obj50 = ns.Object("CrusherDoor3")
	obj51 = ns.Object("CrusherDoor4")
	obj52 = ns.Object("CrushTrigger1")
	obj53 = ns.Object("CrushTrigger2")
	obj54 = ns.Object("CrushTrigger3")
	obj55 = ns.Object("CrushTrigger4")
	obj56 = ns.Object("CrushTrigger5")
	obj57 = ns.Object("CrushTrigger6")
	obj58 = ns.Object("CrushTrigger7")
	gvar59 = ns.ObjectGroup("CrusherTriggers")
	obj60 = ns.Object("CBlock1")
	obj61 = ns.Object("CBlock2")
	obj62 = ns.Object("CBlock3")
	obj63 = ns.Object("CBlock4")
	obj64 = ns.Object("CBlock5")
	obj65 = ns.Object("CBlock6")
	obj66 = ns.Object("CBlock7")
	obj67 = ns.Object("CBlockMover1")
	obj68 = ns.Object("CBlockMover2")
	obj69 = ns.Object("CBlockMover3")
	obj70 = ns.Object("CBlockMover4")
	obj71 = ns.Object("CBlockMover5")
	obj72 = ns.Object("CBlockMover6")
	obj73 = ns.Object("CBlockMover7")
	wp74 = ns.Waypoint("CBlockWP1")
	wp75 = ns.Waypoint("CBlockWP2")
	wp76 = ns.Waypoint("CBlockWP3")
	wp77 = ns.Waypoint("CBlockWP4")
	wp78 = ns.Waypoint("CBlockWP5")
	wp79 = ns.Waypoint("CBlockWP6")
	wp80 = ns.Waypoint("CBlockWP7")
	wp81 = ns.Waypoint("CBlockHome1")
	wp82 = ns.Waypoint("CBlockHome2")
	wp83 = ns.Waypoint("CBlockHome3")
	wp84 = ns.Waypoint("CBlockHome4")
	wp85 = ns.Waypoint("CBlockHome5")
	wp86 = ns.Waypoint("CBlockHome6")
	wp87 = ns.Waypoint("CBlockHome7")
	wp88 = ns.Waypoint("CrushSound1")
	wp89 = ns.Waypoint("CrushSound2")
	wp90 = ns.Waypoint("CrushSound3")
	gvar91 = ns.Wall(242, 142)
	gvar92 = ns.Wall(243, 143)
	gvar93 = ns.Wall(214, 132)
	gvar94 = ns.Wall(215, 131)
	gvar95 = ns.Wall(216, 130)
}
func SetMe() {
	if !(ns.Distance(ns.GetWaypointX(wp81), ns.GetWaypointY(wp81), ns.GetObjectX(obj67), ns.GetObjectY(obj67)) <= 10) {
		goto LABEL1
	}
	ns.MoveObject(obj67, ns.GetWaypointX(wp81), ns.GetWaypointY(wp81))
	ns.MoveObject(obj68, ns.GetWaypointX(wp82), ns.GetWaypointY(wp82))
	ns.MoveObject(obj69, ns.GetWaypointX(wp83), ns.GetWaypointY(wp83))
	ns.MoveObject(obj70, ns.GetWaypointX(wp84), ns.GetWaypointY(wp84))
	ns.MoveObject(obj71, ns.GetWaypointX(wp85), ns.GetWaypointY(wp85))
	ns.MoveObject(obj72, ns.GetWaypointX(wp86), ns.GetWaypointY(wp86))
	ns.MoveObject(obj73, ns.GetWaypointX(wp87), ns.GetWaypointY(wp87))
	ns.UnlockDoor(obj48)
	ns.UnlockDoor(obj49)
	ns.UnlockDoor(obj50)
	ns.UnlockDoor(obj51)
	goto LABEL2
LABEL1:
	ns.MoveObject(obj67, ns.GetObjectX(obj67)-fvar98, ns.GetObjectY(obj67)-fvar99)
	ns.MoveObject(obj68, ns.GetObjectX(obj68)-fvar98, ns.GetObjectY(obj68)-fvar99)
	ns.MoveObject(obj69, ns.GetObjectX(obj69)-fvar98, ns.GetObjectY(obj69)-fvar99)
	ns.MoveObject(obj70, ns.GetObjectX(obj70)-fvar98, ns.GetObjectY(obj70)-fvar99)
	ns.MoveObject(obj71, ns.GetObjectX(obj71)-fvar98, ns.GetObjectY(obj71)-fvar99)
	ns.MoveObject(obj72, ns.GetObjectX(obj72)-fvar98, ns.GetObjectY(obj72)-fvar99)
	ns.MoveObject(obj73, ns.GetObjectX(obj73)-fvar98, ns.GetObjectY(obj73)-fvar99)
	gvar97 = ns.FrameTimer(1, SetMe)
LABEL2:
	return
}
func StartCrusher() {
	if flag101 {
		goto LABEL1
	}
	ns.CancelTimer(gvar97)
	ns.AudioEvent(ns.SpikeBlockMove, wp88)
	ns.AudioEvent(ns.SpikeBlockMove, wp89)
	ns.AudioEvent(ns.SpikeBlockMove, wp90)
	if flag102 {
		goto LABEL2
	}
	ns.LockDoor(obj48)
	ns.LockDoor(obj49)
	ns.LockDoor(obj50)
	ns.LockDoor(obj51)
LABEL2:
	ns.Move(obj60, wp74)
	ns.Move(obj61, wp75)
	ns.Move(obj62, wp76)
	ns.Move(obj63, wp77)
	ns.Move(obj64, wp78)
	ns.Move(obj65, wp79)
	ns.Move(obj66, wp80)
	flag101 = true
	flag100 = false
LABEL1:
	return
}
func CrusherWait() {
	gvar96 = ns.SecondTimer(15, ResetCrusher)
}
func ResetCrusher() {
	if flag100 {
		goto LABEL1
	}
	flag101 = false
	flag100 = true
	ns.MoveObject(obj67, ns.GetObjectX(obj67)-fvar98, ns.GetObjectY(obj67)-fvar99)
	ns.MoveObject(obj68, ns.GetObjectX(obj68)-fvar98, ns.GetObjectY(obj68)-fvar99)
	ns.MoveObject(obj69, ns.GetObjectX(obj69)-fvar98, ns.GetObjectY(obj69)-fvar99)
	ns.MoveObject(obj70, ns.GetObjectX(obj70)-fvar98, ns.GetObjectY(obj70)-fvar99)
	ns.MoveObject(obj71, ns.GetObjectX(obj71)-fvar98, ns.GetObjectY(obj71)-fvar99)
	ns.MoveObject(obj72, ns.GetObjectX(obj72)-fvar98, ns.GetObjectY(obj72)-fvar99)
	ns.MoveObject(obj73, ns.GetObjectX(obj73)-fvar98, ns.GetObjectY(obj73)-fvar99)
	ns.Move(obj60, wp81)
	ns.Move(obj61, wp82)
	ns.Move(obj62, wp83)
	ns.Move(obj63, wp84)
	ns.Move(obj64, wp85)
	ns.Move(obj65, wp86)
	ns.Move(obj66, wp87)
	SetMe()
LABEL1:
	return
}
func CrusherReady() {
	if !ns.IsCaller(obj63) {
		goto LABEL1
	}
	ns.UnlockDoor(obj48)
	ns.UnlockDoor(obj49)
	ns.UnlockDoor(obj50)
	ns.UnlockDoor(obj51)
	ns.CancelTimer(gvar96)
LABEL1:
	return
}
func CrusherSolved() {
	flag102 = true
	if !ns.IsLocked(obj50) {
		goto LABEL1
	}
	ns.UnlockDoor(obj48)
	ns.UnlockDoor(obj49)
	ns.UnlockDoor(obj50)
	ns.UnlockDoor(obj51)
LABEL1:
	return
}
func CrusherNotSolved() {
	flag102 = false
	if ns.IsLocked(obj50) {
		goto LABEL1
	}
	if flag100 {
		goto LABEL1
	}
	ns.LockDoor(obj48)
	ns.LockDoor(obj49)
	ns.LockDoor(obj50)
	ns.LockDoor(obj51)
LABEL1:
	return
}
func OpenCrusherSecret() {
	ns.WallOpen(gvar93)
	ns.WallOpen(gvar94)
	ns.WallOpen(gvar95)
}
func MakeAWish() {
	var v0 int
	if !flag104 {
		goto LABEL1
	}
	v0 = ns.MaxHealth(ns.GetCaller())
	ns.RestoreHealth(ns.GetCaller(), v0-ns.CurrentHealth(ns.GetCaller()))
	ns.AudioEvent(ns.RestoreHealthName, wp103)
	ns.Print("GeneralPrint:WellSignRefresh")
LABEL1:
	return
}
func ResetWell() {
	flag104 = true
}
func DungeonGets() {
	wp103 = ns.Waypoint("WellWP")
	obj105 = ns.Object("ScaryBlock")
	obj106 = ns.Object("ScaryBlockMover")
	obj107 = ns.Object("ScaryBlock2Mover")
	gvar115 = ns.Wall(195, 111)
	gvar116 = ns.Wall(193, 113)
	gvar108 = ns.WallGroup("SpikeBlockWalls")
	wp109 = ns.Waypoint("BlockSound1WP")
	wp110 = ns.Waypoint("BlockSound2WP")
	gvar111 = ns.Wall(148, 178)
	gvar112 = ns.WallGroup("WallGroup1")
	gvar113 = ns.WallGroup("WallGroup2")
	gvar114 = ns.WallGroup("NorthFenceGroup")
	gvar117 = ns.WallGroup("EastFenceWalls")
	obj119 = ns.Object("BlockMover1")
	obj120 = ns.Object("BlockMover2")
	wp121 = ns.Waypoint("BlockWP1")
	wp122 = ns.Waypoint("BlockWP2")
	obj123 = ns.Object("BlockRow1Trigger")
	obj124 = ns.Object("BlockRow1Trigger2")
	wp125 = ns.Waypoint("BlockRow1WP")
	obj126 = ns.Object("BlockMover3")
	obj127 = ns.Object("BlockMover4")
	wp128 = ns.Waypoint("BlockWP3")
	wp129 = ns.Waypoint("BlockWP4")
	obj130 = ns.Object("BlockRow2Trigger")
	obj131 = ns.Object("BlockRow2Trigger2")
	wp132 = ns.Waypoint("BlockRow2WP")
	obj133 = ns.Object("BlockMover5")
	obj134 = ns.Object("BlockMover6")
	wp135 = ns.Waypoint("BlockWP5")
	wp136 = ns.Waypoint("BlockWP6")
	obj137 = ns.Object("BlockRow3Trigger")
	obj138 = ns.Object("BlockRow3Trigger2")
	wp139 = ns.Waypoint("BlockRow3WP")
	obj140 = ns.Object("BlockMover7")
	obj141 = ns.Object("BlockMover8")
	wp142 = ns.Waypoint("BlockWP7")
	wp143 = ns.Waypoint("BlockWP8")
	obj144 = ns.Object("BlockRow4Trigger")
	obj145 = ns.Object("BlockRow4Trigger2")
	wp146 = ns.Waypoint("BlockRow4WP")
}
func OpenWall1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar111)
}
func OpenWallGroup1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar112)
}
func OpenWallGroup2() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar113)
}
func OpenNorthFenceGroup() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar114)
}
func StartScaryBlock() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj106)
	ns.WallGroupOpen(gvar108)
	ns.AudioEvent(ns.SpikeBlockMove, wp109)
}
func OpenScaryBlockWall1() {
	if !ns.IsCaller(obj105) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar115)
LABEL1:
	return
}
func OpenScaryBlockWall2() {
	if !ns.IsCaller(obj105) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj107)
	ns.AudioEvent(ns.SpikeBlockMove, wp110)
	ns.WallOpen(gvar116)
LABEL1:
	return
}
func OpenEastFenceWalls() {
	if flag118 {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	flag118 = true
	ns.WallGroupOpen(gvar117)
	goto LABEL2
LABEL1:
	flag118 = false
	ns.WallGroupClose(gvar117)
LABEL2:
	return
}
func StartBlockRow1() {
	ns.ObjectOff(obj123)
	ns.ObjectOff(obj124)
	ns.Move(obj119, wp121)
	ns.Move(obj120, wp122)
	ns.AudioEvent(ns.SpikeBlockMove, wp125)
	ns.FrameTimer(15, BlockRow1Boom)
}
func BlockRow1Boom() {
	ns.AudioEvent(ns.HammerMissing, wp125)
	ns.Effect(ns.JIGGLE, ns.GetWaypointX(wp125), ns.GetWaypointY(wp125), 10, 0)
	ns.FrameTimer(20, ResetBlockRow1)
}
func ResetBlockRow1() {
	ns.ObjectOn(obj123)
	ns.ObjectOn(obj124)
}
func StartBlockRow2() {
	ns.ObjectOff(obj130)
	ns.ObjectOff(obj131)
	ns.Move(obj126, wp128)
	ns.Move(obj127, wp129)
	ns.AudioEvent(ns.SpikeBlockMove, wp132)
	ns.FrameTimer(15, BlockRow2Boom)
}
func BlockRow2Boom() {
	ns.AudioEvent(ns.HammerMissing, wp132)
	ns.Effect(ns.JIGGLE, ns.GetWaypointX(wp132), ns.GetWaypointY(wp132), 10, 0)
	ns.FrameTimer(20, ResetBlockRow2)
}
func ResetBlockRow2() {
	ns.ObjectOn(obj130)
	ns.ObjectOn(obj131)
}
func StartBlockRow3() {
	ns.ObjectOff(obj137)
	ns.ObjectOff(obj138)
	ns.Move(obj133, wp135)
	ns.Move(obj134, wp136)
	ns.AudioEvent(ns.SpikeBlockMove, wp139)
	ns.FrameTimer(15, BlockRow3Boom)
}
func BlockRow3Boom() {
	ns.AudioEvent(ns.HammerMissing, wp139)
	ns.Effect(ns.JIGGLE, ns.GetWaypointX(wp139), ns.GetWaypointY(wp139), 10, 0)
	ns.FrameTimer(20, ResetBlockRow3)
}
func ResetBlockRow3() {
	ns.ObjectOn(obj137)
	ns.ObjectOn(obj138)
}
func StartBlockRow4() {
	ns.ObjectOff(obj144)
	ns.ObjectOff(obj145)
	ns.Move(obj140, wp142)
	ns.Move(obj141, wp143)
	ns.AudioEvent(ns.SpikeBlockMove, wp146)
	ns.FrameTimer(15, BlockRow4Boom)
}
func BlockRow4Boom() {
	ns.AudioEvent(ns.HammerMissing, wp146)
	ns.Effect(ns.JIGGLE, ns.GetWaypointX(wp146), ns.GetWaypointY(wp146), 10, 0)
	ns.FrameTimer(20, ResetBlockRow4)
}
func ResetBlockRow4() {
	ns.ObjectOn(obj144)
	ns.ObjectOn(obj145)
}
func TrapGeneratorOn() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar153)
	ns.ObjectGroupOn(gvar152)
}
func TrapGeneratorOff() {
	ns.WallGroupClose(gvar153)
	ns.ObjectGroupOff(gvar152)
}
func SwitchSpikeGroup1() {
	if ns.IsObjectOn(obj148) {
		goto LABEL1
	}
	ns.ObjectGroupOn(gvar151)
	ns.AudioEvent(ns.FloorSpikesUp, wp187)
	goto LABEL2
LABEL1:
	ns.ObjectGroupOff(gvar151)
	ns.AudioEvent(ns.FloorSpikesDown, wp187)
LABEL2:
	return
}
func SwitchGun1() {
	if ns.IsObjectOn(obj147) {
		goto LABEL1
	}
	ns.ObjectOn(obj147)
	goto LABEL2
LABEL1:
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOff(obj147)
	ns.Print("GeneralPrint:MsgTrapOff")
LABEL2:
	return
}
func OpenSpellWallsA() {
	ns.WallGroupOpen(gvar156)
	ns.ObjectOff(ns.GetTrigger())
}
func OpenFenceWalls() {
	if flag158 {
		goto LABEL1
	}
	ns.WallGroupOpen(gvar157)
	flag158 = true
	goto LABEL2
LABEL1:
	ns.WallGroupClose(gvar157)
	flag158 = false
LABEL2:
	return
}
func OpenTreasureWalls() {
	ns.WallGroupOpen(gvar155)
	ns.ObjectOff(ns.GetTrigger())
}
func Infrawel() {
	ns.Enchant(ns.GetCaller(), ns.ENCHANT_INFRAVISION, 30)
	ns.Print("GeneralPrint:WishingWellPrintEnchantment")
}
func StartFourTrap() {
	ns.ObjectGroupOff(gvar164)
	ns.WallOpen(gvar165)
	ns.WallOpen(gvar166)
	ns.WallOpen(gvar167)
	ns.WallOpen(gvar168)
	ns.WallOpen(gvar169)
	ns.WallOpen(gvar170)
	ns.WallOpen(gvar171)
	ns.WallOpen(gvar172)
	ns.WallOpen(gvar173)
	ns.WallOpen(gvar174)
	ns.WallOpen(gvar175)
	ns.WallOpen(gvar176)
	ns.WallOpen(gvar177)
	ns.WallOpen(gvar178)
	ns.WallOpen(gvar179)
	ns.WallOpen(gvar180)
	ns.WallOpen(gvar181)
	ns.WallOpen(gvar182)
	ns.WallOpen(gvar183)
	ns.WallOpen(gvar184)
	ns.WallOpen(gvar185)
	ns.WallOpen(gvar186)
}
func StopMechavator() {
	ns.ObjectOff(obj160)
}
func StartMechavator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.CreatureCageAppears, wp163)
	ns.AudioEvent(ns.ChangeSpellbar, wp163)
	ns.ObjectOn(obj160)
	ns.ObjectGroupOn(gvar161)
}
func SetMech() {
	ns.MoveObject(obj159, ns.GetWaypointX(wp162), ns.GetWaypointY(wp162))
}
func GetSentryRays() {
	obj189 = ns.Object("SentryRay1")
	obj190 = ns.Object("SentryRay2")
	obj191 = ns.Object("SentryRay3")
	obj192 = ns.Object("SentryRay4")
	obj193 = ns.Object("SentryRay5")
	obj194 = ns.Object("RayTrigger1")
	obj195 = ns.Object("RayTrigger2")
}
func TeleportGets() {
	obj203 = ns.Object("DoomRay1")
	obj204 = ns.Object("DoomRay2")
	obj205 = ns.Object("DoomRayMover1")
	obj206 = ns.Object("DoomRayMover2")
	obj207 = ns.Object("DoomStatue1")
	obj208 = ns.Object("DoomStatue2")
	obj209 = ns.Object("DoomStatue3")
	obj210 = ns.Object("DoomStatue4")
}
func WizardsGets() {
	obj213 = ns.Object("SouthBlock")
	obj214 = ns.Object("SouthBlock2")
	obj215 = ns.Object("SouthBlockMover")
	obj216 = ns.Object("SouthBlock2Mover")
	gvar217 = ns.WallGroup("SouthWalls2")
	gvar219 = ns.WallGroup("SouthTWalls")
	gvar218 = ns.Wall(58, 202)
	gvar220 = ns.Wall(47, 221)
	wp221 = ns.Waypoint("RightSoundWP")
	obj222 = ns.Object("NorthBlock")
	obj223 = ns.Object("NorthBlock2")
	obj224 = ns.Object("NorthBlockMover")
	obj225 = ns.Object("NorthBlock2Mover")
	gvar226 = ns.WallGroup("NorthWalls1")
	gvar227 = ns.WallGroup("NorthTWalls")
	gvar228 = ns.Wall(52, 196)
	gvar229 = ns.Wall(33, 207)
	wp230 = ns.Waypoint("LeftSoundWP")
	obj231 = ns.Object("ArenaDoor")
	gvar232 = ns.WallGroup("ArenaWalls")
	obj234 = ns.Object("FenceDoor1")
	obj235 = ns.Object("FenceDoor2")
	obj236 = ns.Object("WizArrow1")
	obj237 = ns.Object("WizArrow2")
	obj238 = ns.Object("WizArrow3")
	obj239 = ns.Object("WizArrow4")
	obj240 = ns.Object("WizArrow5")
	obj241 = ns.Object("WizArrow6")
	obj242 = ns.Object("WizArrow7")
	gvar246 = ns.WallGroup("WizExitWalls1")
	ns.LockDoor(obj231)
	ns.LockDoor(obj234)
	ns.LockDoor(obj235)
}
func WizPart1Gets() {
	gvar247 = ns.WallGroup("MagicWalls1")
	obj249 = ns.Object("SecretChairMover")
	gvar250 = ns.WallGroup("ChairWalls")
	gvar251 = ns.WallGroup("WallRing1")
	gvar252 = ns.WallGroup("WallRing2")
	gvar253 = ns.WallGroup("WallRing3")
	obj254 = ns.Object("BonusRay1")
	obj255 = ns.Object("BonusRay2")
	obj256 = ns.Object("BonusRayChain")
	obj257 = ns.Object("GuardRay")
	obj258 = ns.Object("FourDoor")
	obj261 = ns.Object("BreakDropDoor")
	gvar262 = ns.ObjectGroup("BreakDropTriggers")
	gvar263 = ns.WallGroup("BreakingWalls")
	gvar264 = ns.WallGroup("DroppingWalls")
	obj265 = ns.Object("InvisoGen1")
	obj266 = ns.Object("InvisoGen2")
	obj267 = ns.Object("InvisoGen3")
	obj268 = ns.Object("InvisoChest1")
	obj269 = ns.Object("InvisoChest2")
	obj270 = ns.Object("InvisoChest3")
	gvar271 = ns.Wall(70, 122)
	gvar272 = ns.Wall(71, 123)
	gvar273 = ns.Wall(72, 124)
	gvar274 = ns.Wall(70, 116)
	gvar275 = ns.Wall(71, 115)
	ns.ObjectOn(obj249)
	ns.LockDoor(obj261)
	ns.Enchant(obj265, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj266, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj267, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj268, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj269, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj270, ns.ENCHANT_INVISIBLE, 0)
}
func MapInitialize() {
	obj149 = ns.Object("TrapGenerator")
	obj148 = ns.Object("TestSpike")
	obj150 = ns.Object("TestGenerator")
	gvar151 = ns.ObjectGroup("SpikeGroup1")
	gvar152 = ns.ObjectGroup("WizGroup1")
	obj147 = ns.Object("SkullGun1")
	gvar157 = ns.WallGroup("FenceWalls")
	gvar156 = ns.WallGroup("SpellWallsA")
	gvar155 = ns.WallGroup("TreasureWalls")
	gvar153 = ns.WallGroup("TrapGeneratorWalls")
	gvar154 = ns.WallGroup("TrapGeneratorExitWalls")
	obj160 = ns.Object("MechElevator")
	obj159 = ns.Object("Mech")
	wp162 = ns.Waypoint("MechWP")
	gvar161 = ns.ObjectGroup("MechGears")
	wp163 = ns.Waypoint("MechavatorSoundWP")
	gvar164 = ns.ObjectGroup("FourTrapTrigs")
	gvar165 = ns.Wall(89, 115)
	gvar166 = ns.Wall(90, 116)
	gvar167 = ns.Wall(91, 117)
	gvar168 = ns.Wall(92, 118)
	gvar169 = ns.Wall(93, 119)
	gvar170 = ns.Wall(94, 120)
	gvar171 = ns.Wall(95, 119)
	gvar172 = ns.Wall(96, 118)
	gvar173 = ns.Wall(97, 117)
	gvar174 = ns.Wall(98, 116)
	gvar175 = ns.Wall(99, 115)
	gvar176 = ns.Wall(100, 114)
	gvar177 = ns.Wall(99, 113)
	gvar178 = ns.Wall(98, 112)
	gvar179 = ns.Wall(97, 111)
	gvar180 = ns.Wall(96, 110)
	gvar181 = ns.Wall(95, 109)
	gvar182 = ns.Wall(94, 110)
	gvar183 = ns.Wall(93, 111)
	gvar184 = ns.Wall(92, 112)
	gvar185 = ns.Wall(91, 113)
	gvar186 = ns.Wall(90, 114)
	wp187 = ns.Waypoint("FloorSpikeSound")
	obj188 = ns.Object("TestChest")
	GetArrows()
	GetSentryRays()
	DungeonGets()
	CrusherGets()
	TeleportGets()
	WizardsGets()
	WizPart1Gets()
	ns.ObjectGroupOff(gvar161)
	ns.FrameTimer(36, StopMechavator)
}
func StartRays() {
	SentryRun()
	ns.ObjectOff(obj194)
	ns.ObjectOff(obj195)
}
func SentryRun() {
	var v0 int
	switch v0 = gvar201; v0 {
	case gvar196:
		goto LABEL1
	case gvar197:
		goto LABEL2
	case gvar198:
		goto LABEL3
	case gvar199:
		goto LABEL4
	case gvar200:
		goto LABEL5
	default:
		goto LABEL6
	}
LABEL1:
	ns.ObjectOn(obj191)
	ns.ObjectOn(obj192)
	ns.ObjectOff(obj193)
	gvar201 = gvar197
	goto LABEL6
LABEL2:
	ns.ObjectOn(obj191)
	ns.ObjectOff(obj192)
	ns.ObjectOn(obj193)
	gvar201 = gvar198
	goto LABEL6
LABEL3:
	ns.ObjectOff(obj191)
	ns.ObjectOn(obj192)
	ns.ObjectOn(obj193)
	gvar201 = gvar196
	goto LABEL6
LABEL4:
	ns.ObjectOn(obj191)
	ns.ObjectOn(obj192)
	ns.ObjectOn(obj193)
	gvar201 = gvar196
	goto LABEL6
LABEL5:
	ns.ObjectOff(obj189)
	ns.ObjectOff(obj190)
	ns.ObjectOn(obj191)
	ns.ObjectOn(obj192)
	ns.ObjectOn(obj193)
	gvar201 = gvar196
	goto LABEL6
LABEL6:
	if flag202 {
		goto LABEL7
	}
	ns.FrameTimer(90, SentryRun)
	goto LABEL8
LABEL7:
	ns.ObjectOff(obj189)
	ns.ObjectOff(obj190)
	ns.ObjectOff(obj191)
	ns.ObjectOff(obj192)
	ns.ObjectOff(obj193)
LABEL8:
	return
}
func SentryRaySwitch() {
	if flag202 {
		goto LABEL1
	}
	flag202 = true
	ns.ObjectOff(ns.GetTrigger())
	ns.Print("GeneralPrint:MsgTrapOff")
	goto LABEL2
LABEL1:
	flag202 = false
	SentryRun()
LABEL2:
	return
}
func TakeShot() {
	var v0 int
	switch v0 = ivar212; v0 {
	case 1:
		goto LABEL1
	case 2:
		goto LABEL2
	case 3:
		goto LABEL3
	case 4:
		goto LABEL4
	default:
		goto LABEL5
	}
LABEL1:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj207, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar212 += 1
	goto LABEL5
LABEL2:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj208, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar212 += 1
	goto LABEL5
LABEL3:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj209, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar212 += 1
	goto LABEL5
LABEL4:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj210, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar212 = 1
	goto LABEL5
LABEL5:
	return
}
func StartNewRays() {
	if flag211 {
		goto LABEL1
	}
	ns.ObjectOn(obj205)
	ns.ObjectOn(obj206)
	flag211 = true
LABEL1:
	return
}
func StopNewRays() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOff(obj205)
	ns.ObjectOff(obj206)
	ns.ObjectOff(obj203)
	ns.ObjectOff(obj204)
}
func StartSouthBlock() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar217)
	ns.ObjectOn(obj215)
	ns.ObjectOn(obj216)
	ns.AudioEvent(ns.SpikeBlockMove, wp221)
}
func OpenSouthWall3() {
	if !ns.IsCaller(obj213) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar219)
	ns.WallOpen(gvar220)
LABEL1:
	return
}
func StartNorthBlock() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar226)
	ns.ObjectOn(obj224)
	ns.ObjectOn(obj225)
	ns.AudioEvent(ns.SpikeBlockMove, wp230)
}
func OpenNorthWall2() {
	if !ns.IsCaller(obj222) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar228)
	ns.WallGroupOpen(gvar227)
	ns.WallOpen(gvar229)
LABEL1:
	return
}
func OpenArenaWalls() {
	if flag233 {
		goto LABEL1
	}
	ns.WallGroupOpen(gvar232)
	flag233 = true
	goto LABEL2
LABEL1:
	ns.WallGroupClose(gvar232)
	flag233 = false
LABEL2:
	return
}
func OpenArenaDoor() {
	ns.ObjectOff(ns.GetTrigger())
	ns.UnlockDoor(obj231)
}
func OpenFenceDoor1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.UnlockDoor(obj234)
}
func OpenFenceDoor2() {
	ns.ObjectOff(ns.GetTrigger())
	ns.UnlockDoor(obj235)
}
func WizHitRow1() {
	ivar243 += 1
	CheckWizRows(1)
}
func WizReleaseRow1() {
	ivar243 -= 1
	CheckWizRows(1)
}
func WizHitRow2() {
	ivar244 += 1
	CheckWizRows(2)
}
func WizReleaseRow2() {
	ivar244 -= 1
	CheckWizRows(2)
}
func WizHitRow3() {
	ivar245 = 1
	CheckWizRows(3)
}
func WizReleaseRow3() {
	ivar245 -= 1
	CheckWizRows(3)
}
func OpenWizExitWalls1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar246)
}
func CheckWizRows(a1 int) {
	var v0 int
	switch v0 = a1; v0 {
	case 1:
		goto LABEL1
	case 2:
		goto LABEL2
	case 3:
		goto LABEL3
	default:
		goto LABEL4
	}
LABEL1:
	if !(ivar243 > 0) {
		goto LABEL5
	}
	ns.ObjectOn(obj236)
	goto LABEL6
LABEL5:
	ns.ObjectOff(obj236)
LABEL6:
	goto LABEL4
LABEL2:
	if !(ivar244 > 0) {
		goto LABEL7
	}
	ns.ObjectOn(obj237)
	goto LABEL8
LABEL7:
	ns.ObjectOff(obj237)
LABEL8:
	goto LABEL4
LABEL3:
	if !(ivar245 > 0) {
		goto LABEL9
	}
	if ns.IsObjectOn(obj238) {
		goto LABEL10
	}
	ns.ObjectOn(obj238)
	ns.ObjectOn(obj239)
	ns.ObjectOn(obj240)
	ns.ObjectOn(obj241)
	ns.ObjectOn(obj242)
LABEL10:
	goto LABEL11
LABEL9:
	ns.ObjectOff(obj238)
	ns.ObjectOff(obj239)
	ns.ObjectOff(obj240)
	ns.ObjectOff(obj241)
	ns.ObjectOff(obj242)
LABEL11:
	goto LABEL4
LABEL4:
	return
}
func OpenMagicWalls1() {
	ns.ObjectOff(ns.GetTrigger())
	if flag248 {
		goto LABEL1
	}
	ns.WallGroupOpen(gvar247)
	flag248 = true
	goto LABEL2
LABEL1:
	ns.WallGroupClose(gvar247)
	flag248 = false
LABEL2:
	return
}
func OpenChairWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar250)
}
func OpenWallRing1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar251)
	ns.Print("GeneralPrint:MsgWallLowered")
}
func OpenWallRing2() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar252)
	ns.Print("GeneralPrint:MsgWallLowered")
}
func OpenWallRing3() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar253)
	ns.Print("GeneralPrint:MsgWallLowered")
}
func StartBonusRay() {
	if ns.IsObjectOn(obj254) {
		goto LABEL1
	}
	ns.ObjectOff(obj256)
	ns.ObjectOn(obj254)
	ns.ObjectOn(obj255)
	ns.SecondTimer(1, StopBonusRay)
LABEL1:
	return
}
func StopBonusRay() {
	if !ns.IsObjectOn(obj254) {
		goto LABEL1
	}
	ns.ObjectOn(obj256)
	ns.ObjectOff(obj254)
	ns.ObjectOff(obj255)
LABEL1:
	return
}
func StartGuardRay() {
	ns.ObjectOn(obj257)
}
func StopGuardRay() {
}
func SecretDown() {
	ivar259 += 1
	if !(ivar259 >= 4) {
		goto LABEL1
	}
	if !ns.IsLocked(obj258) {
		goto LABEL1
	}
	ns.UnlockDoor(obj258)
LABEL1:
	return
}
func SecretUp() {
	ivar259 -= 1
}
func FourDoorKill() {
	if flag260 {
		goto LABEL1
	}
	ns.LockDoor(obj258)
	flag260 = true
LABEL1:
	return
}
func OpenFourDoor() {
	if !ns.IsLocked(obj258) {
		goto LABEL1
	}
	ns.UnlockDoor(obj258)
LABEL1:
	return
}
func BreakWalls() {
	ns.WallGroupBreak(gvar263)
}
func DropWalls() {
	ns.ObjectGroupOff(gvar262)
	ns.NoWallSound(true)
	ns.WallGroupOpen(gvar264)
	ns.NoWallSound(false)
	ns.FrameTimer(1, BreakWalls)
}
func OpenBreakDropDoor() {
	ns.ObjectOff(ns.GetTrigger())
	ns.Print("GeneralPrint:MsgDoorUnlock")
	ns.UnlockDoor(obj261)
}
func InvisoSpawn() {
	ns.Enchant(ns.GetCaller(), ns.ENCHANT_INVISIBLE, 0)
}
func CloseInvisoWall() {
	ns.WallClose(gvar271)
	ns.WallClose(gvar272)
	ns.WallClose(gvar273)
	ns.WallOpen(gvar274)
	ns.WallOpen(gvar275)
}
func OpenInvisoWall() {
	ns.WallOpen(gvar271)
	ns.WallOpen(gvar272)
	ns.WallOpen(gvar273)
	ns.WallClose(gvar274)
	ns.WallClose(gvar275)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
