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
	obj14   ns.ObjectID
	obj15   ns.ObjectID
	obj16   ns.ObjectID
	obj17   ns.ObjectID
	obj18   ns.ObjectID
	obj19   ns.ObjectID
	obj20   ns.ObjectID
	obj21   ns.ObjectID
	obj22   ns.ObjectID
	obj23   ns.ObjectID
	obj24   ns.ObjectID
	ivar25  int
	ivar26  int
	ivar27  int
	ivar28  int
	ivar29  int
	ivar30  int
	ivar31  int
	ivar32  int
	obj33   ns.ObjectID
	obj34   ns.ObjectID
	obj35   ns.ObjectID
	obj36   ns.ObjectID
	gvar37  int
	gvar38  int
	gvar39  int
	obj40   ns.ObjectID
	obj41   ns.ObjectID
	obj42   ns.ObjectID
	obj43   ns.ObjectID
	obj44   ns.ObjectID
	obj45   ns.ObjectID
	obj46   ns.ObjectID
	obj47   ns.ObjectID
	gvar48  ns.WallGroupID
	flag49  bool
	obj50   ns.ObjectID
	obj51   ns.ObjectID
	gvar52  ns.WallID
	gvar53  ns.WallID
	gvar54  ns.WallGroupID
	gvar55  ns.WallGroupID
	flag56  bool
	obj57   ns.ObjectID
	fvar58  float32
	fvar59  float32
	flag60  bool
	gvar61  ns.WallID
	gvar62  ns.WallID
	gvar63  ns.WallID
	gvar64  ns.WallID
	gvar65  ns.WallID
	gvar66  ns.WallID
	obj67   ns.ObjectID
	obj68   ns.ObjectID
	obj69   ns.ObjectID
	obj70   ns.ObjectID
	ivar71  int
	obj72   ns.ObjectID
	obj73   ns.ObjectID
	obj74   ns.ObjectID
	obj75   ns.ObjectID
	obj76   ns.ObjectID
	obj77   ns.ObjectID
	ivar78  int
	gvar79  ns.WallID
	gvar80  ns.WallID
	gvar81  ns.WallID
	gvar82  ns.WallGroupID
	gvar83  ns.WallGroupID
	gvar84  ns.WallGroupID
	gvar85  ns.WallGroupID
	ivar86  int
	obj87   ns.ObjectID
	obj88   ns.ObjectID
	obj89   ns.ObjectID
	obj90   ns.ObjectID
	obj91   ns.ObjectID
	obj92   ns.ObjectID
	obj93   ns.ObjectID
	obj94   ns.ObjectID
	obj95   ns.ObjectID
	obj96   ns.ObjectID
	obj97   ns.ObjectID
	obj98   ns.ObjectID
	obj99   ns.ObjectID
	obj100  ns.ObjectID
	obj101  ns.ObjectID
	obj102  ns.ObjectID
	obj103  ns.ObjectID
	obj104  ns.ObjectID
	obj105  ns.ObjectID
	obj106  ns.ObjectID
	obj107  ns.ObjectID
	obj108  ns.ObjectID
	ivar109 int
	flag110 bool
	gvar111 ns.WallGroupID
	obj112  ns.ObjectID
	obj113  ns.ObjectID
	obj114  ns.ObjectID
	obj115  ns.ObjectID
	obj116  ns.ObjectID
	gvar117 ns.WallID
	gvar118 ns.WallID
	gvar119 ns.WallID
	gvar120 ns.WallID
	gvar121 ns.WallID
	gvar122 ns.WallID
	gvar123 ns.WallID
	gvar124 ns.WallID
	gvar125 ns.WallID
	gvar126 ns.WallID
	gvar127 ns.WallID
	gvar128 ns.WallID
	gvar129 ns.WallID
	gvar130 ns.WallID
	gvar131 ns.WallID
	gvar132 ns.WallID
	gvar133 ns.WallGroupID
	gvar134 ns.WallID
	gvar135 ns.WallID
	gvar136 ns.WallID
	gvar137 ns.WallID
	gvar138 ns.WallID
	gvar139 ns.WallID
	gvar140 ns.WallID
	gvar141 ns.WallID
	gvar142 ns.WallID
	gvar143 ns.WallGroupID
	ivar144 int
	gvar145 ns.WallID
	gvar146 ns.ObjectGroupID
	wp147   ns.WaypointID
	wp148   ns.WaypointID
	obj149  ns.ObjectID
	obj150  ns.ObjectID
	obj151  ns.ObjectID
	gvar152 ns.WallID
	gvar153 ns.WallID
	gvar154 ns.WallID
	gvar155 ns.WallID
	gvar156 ns.WallID
	gvar157 ns.WallID
	gvar158 ns.WallGroupID
	gvar159 ns.WallGroupID
	gvar160 ns.WallGroupID
	ivar161 int
	obj162  ns.ObjectID
	obj163  ns.ObjectID
	obj164  ns.ObjectID
	obj165  ns.ObjectID
	obj166  ns.ObjectID
	obj167  ns.ObjectID
	obj168  ns.ObjectID
	obj169  ns.ObjectID
	obj170  ns.ObjectID
	obj171  ns.ObjectID
	obj172  ns.ObjectID
	obj173  ns.ObjectID
	obj174  ns.ObjectID
	obj175  ns.ObjectID
	obj176  ns.ObjectID
	obj177  ns.ObjectID
	obj178  ns.ObjectID
	obj179  ns.ObjectID
	obj180  ns.ObjectID
	obj181  ns.ObjectID
	obj182  ns.ObjectID
	obj183  ns.ObjectID
	gvar184 ns.WallID
	gvar185 ns.WallID
	gvar186 ns.WallID
	gvar187 ns.WallID
	gvar188 ns.WallID
	obj189  ns.ObjectID
	obj190  ns.ObjectID
	obj191  ns.ObjectID
	obj192  ns.ObjectID
	obj193  ns.ObjectID
	obj194  ns.ObjectID
	obj195  ns.ObjectID
	obj196  ns.ObjectID
	obj197  ns.ObjectID
	obj198  ns.ObjectID
	obj199  ns.ObjectID
	obj200  ns.ObjectID
	gvar201 ns.ObjectGroupID
	wp202   ns.WaypointID
	wp203   ns.WaypointID
	wp204   ns.WaypointID
	wp205   ns.WaypointID
	wp206   ns.WaypointID
	wp207   ns.WaypointID
	wp208   ns.WaypointID
	wp209   ns.WaypointID
	obj210  ns.ObjectID
	obj211  ns.ObjectID
	obj212  ns.ObjectID
	gvar213 ns.WallID
	gvar214 ns.WallID
	gvar215 ns.WallID
	gvar216 ns.WallID
	gvar217 ns.WallID
	gvar218 ns.WallID
	gvar219 ns.WallID
	gvar220 ns.WallID
	wp221   ns.WaypointID
	flag222 bool
	obj223  ns.ObjectID
	obj224  ns.ObjectID
	obj225  ns.ObjectID
	obj226  ns.ObjectID
	obj227  ns.ObjectID
	obj228  ns.ObjectID
	obj229  ns.ObjectID
	obj230  ns.ObjectID
	obj231  ns.ObjectID
	obj232  ns.ObjectID
	obj233  ns.ObjectID
	obj234  ns.ObjectID
	obj235  ns.ObjectID
	obj236  ns.ObjectID
	obj237  ns.ObjectID
	obj238  ns.ObjectID
	obj239  ns.ObjectID
	obj240  ns.ObjectID
	gvar241 ns.WallID
	gvar242 ns.WallID
	gvar243 ns.WallID
	gvar244 ns.WallGroupID
	obj245  ns.ObjectID
	obj246  ns.ObjectID
	obj247  ns.ObjectID
	obj248  ns.ObjectID
	obj249  ns.ObjectID
	obj250  ns.ObjectID
	obj251  ns.ObjectID
	obj252  ns.ObjectID
	obj253  ns.ObjectID
	gvar254 ns.ObjectGroupID
	gvar255 ns.WallGroupID
	wp256   ns.WaypointID
	ivar257 int
	gvar258 ns.TimerID
	obj259  ns.ObjectID
	obj260  ns.ObjectID
	obj261  ns.ObjectID
	obj262  ns.ObjectID
	gvar263 ns.WallID
	gvar264 ns.WallID
	gvar265 ns.WallID
	gvar266 ns.WallID
	gvar267 ns.WallID
	gvar268 ns.WallGroupID
	gvar269 ns.WallGroupID
	gvar270 ns.WallGroupID
	obj271  ns.ObjectID
	obj272  ns.ObjectID
	gvar273 ns.ObjectGroupID
	wp274   ns.WaypointID
	flag275 bool
	flag276 bool
	obj277  ns.ObjectID
	wp278   ns.WaypointID
	wp279   ns.WaypointID
	gvar280 ns.TimerID
	gvar281 ns.TimerID
	gvar282 ns.ObjectGroupID
	gvar283 ns.WallID
	gvar284 ns.WallID
	gvar285 ns.WallID
	gvar286 ns.WallID
	gvar287 ns.WallID
	gvar288 ns.WallID
	gvar289 ns.WallID
	gvar290 ns.WallID
	wp291   ns.WaypointID
	gvar292 ns.TimerID
	ivar293 int
	obj294  ns.ObjectID
	obj295  ns.ObjectID
	gvar296 ns.WallID
	gvar297 ns.WallID
	gvar298 ns.WallID
	gvar299 ns.WallID
	gvar300 ns.WallID
	gvar301 ns.WallID
	obj302  ns.ObjectID
	obj303  ns.ObjectID
	obj304  ns.ObjectID
	obj305  ns.ObjectID
	obj306  ns.ObjectID
	obj307  ns.ObjectID
)

func init() {
	ivar25 = 0
	ivar26 = 0
	ivar27 = 0
	ivar28 = 0
	ivar29 = 0
	ivar30 = 0
	ivar31 = 0
	ivar32 = 0
	flag49 = false
	flag56 = false
	flag60 = false
	ivar71 = 1
	ivar78 = 1
	ivar86 = 0
	ivar109 = 0
	flag110 = false
	ivar144 = 0
	ivar161 = 0
	flag222 = false
	ivar257 = 0
	gvar258 = 0
	flag275 = false
	flag276 = false
	gvar280 = 0
	gvar281 = 0
	gvar292 = 0
	ivar293 = 0
}
func TombPlateHit1() {
	ivar25 += 1
	CheckRows(1)
}
func TombPlateRelease1() {
	ivar25 -= 1
	CheckRows(1)
}
func TombPlateHit2() {
	ivar26 += 1
	CheckRows(2)
}
func TombPlateRelease2() {
	ivar26 -= 1
	CheckRows(2)
}
func TombPlateHit3() {
	ivar27 += 1
	CheckRows(3)
}
func TombPlateRelease3() {
	ivar27 -= 1
	CheckRows(3)
}
func TombPlateHit4() {
	ivar28 += 1
	CheckRows(4)
}
func TombPlateRelease4() {
	ivar28 -= 1
	CheckRows(4)
}
func SetHit1() {
	ivar29 += 1
	WhatRow(1)
}
func SetRelease1() {
	ivar29 -= 1
	WhatRow(1)
}
func SetHit2() {
	ivar30 += 2
	WhatRow(2)
}
func SetRelease2() {
	ivar30 -= 1
	WhatRow(2)
}
func SetHit3() {
	ivar31 += 1
	WhatRow(3)
}
func SetRelease3() {
	ivar31 -= 1
	WhatRow(3)
}
func SetHit4() {
	ivar32 += 1
	WhatRow(4)
}
func SetRelease4() {
	ivar32 -= 1
	WhatRow(4)
}
func ArrowsInit() {
	obj4 = ns.Object("TombArrow1")
	obj5 = ns.Object("TombArrow2")
	obj6 = ns.Object("TombArrow3")
	obj7 = ns.Object("TombArrow4")
	obj8 = ns.Object("TombArrow5")
	obj9 = ns.Object("TombArrow6")
	obj10 = ns.Object("TombArrow7")
	obj11 = ns.Object("TombArrow8")
	obj12 = ns.Object("TombArrow9")
	obj13 = ns.Object("TombArrow10")
	obj14 = ns.Object("TombArrow11")
	obj15 = ns.Object("TombArrow12")
	obj16 = ns.Object("Arrow1")
	obj17 = ns.Object("Arrow2")
	obj18 = ns.Object("Arrow3")
	obj19 = ns.Object("Arrow4")
	obj20 = ns.Object("Arrow5")
	obj21 = ns.Object("Arrow6")
	obj22 = ns.Object("Arrow7")
	obj23 = ns.Object("Arrow8")
	obj24 = ns.Object("Arrow9")
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
	default:
		goto LABEL5
	}
LABEL1:
	if !(ivar25 > 0) {
		goto LABEL6
	}
	ns.ObjectOn(obj4)
	ns.ObjectOn(obj5)
	ns.ObjectOn(obj6)
	goto LABEL7
LABEL6:
	ns.ObjectOff(obj4)
	ns.ObjectOff(obj5)
	ns.ObjectOff(obj6)
LABEL7:
	goto LABEL5
LABEL2:
	if !(ivar26 > 0) {
		goto LABEL8
	}
	ns.ObjectOn(obj7)
	ns.ObjectOn(obj8)
	ns.ObjectOn(obj9)
	goto LABEL9
LABEL8:
	ns.ObjectOff(obj7)
	ns.ObjectOff(obj8)
	ns.ObjectOff(obj9)
LABEL9:
	goto LABEL5
LABEL3:
	if !(ivar27 > 0) {
		goto LABEL10
	}
	ns.ObjectOn(obj10)
	ns.ObjectOn(obj11)
	ns.ObjectOn(obj12)
	goto LABEL11
LABEL10:
	ns.ObjectOff(obj10)
	ns.ObjectOff(obj11)
	ns.ObjectOff(obj12)
LABEL11:
	goto LABEL5
LABEL4:
	if !(ivar28 > 0) {
		goto LABEL12
	}
	ns.ObjectOn(obj13)
	ns.ObjectOn(obj14)
	ns.ObjectOn(obj15)
	goto LABEL13
LABEL12:
	ns.ObjectOff(obj13)
	ns.ObjectOff(obj14)
	ns.ObjectOff(obj15)
LABEL13:
	goto LABEL5
LABEL5:
	return
}
func WhatRow(a1 int) {
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
	default:
		goto LABEL5
	}
LABEL1:
	if !(ivar29 > 0) {
		goto LABEL6
	}
	ns.ObjectOn(obj16)
	ns.ObjectOn(obj17)
	ns.ObjectOn(obj18)
	goto LABEL7
LABEL6:
	ns.ObjectOff(obj16)
	ns.ObjectOff(obj17)
	ns.ObjectOff(obj18)
LABEL7:
	goto LABEL5
LABEL2:
	if !(ivar30 > 0) {
		goto LABEL8
	}
	ns.ObjectOn(obj19)
	ns.ObjectOn(obj20)
	ns.ObjectOn(obj21)
	goto LABEL9
LABEL8:
	ns.ObjectOff(obj19)
	ns.ObjectOff(obj20)
	ns.ObjectOff(obj21)
LABEL9:
	goto LABEL5
LABEL3:
	if !(ivar31 > 0) {
		goto LABEL10
	}
	ns.ObjectOn(obj22)
	ns.ObjectOn(obj23)
	ns.ObjectOn(obj24)
	goto LABEL11
LABEL10:
	ns.ObjectOff(obj22)
	ns.ObjectOff(obj23)
	ns.ObjectOff(obj24)
LABEL11:
	goto LABEL5
LABEL4:
	if !(ivar32 > 0) {
		goto LABEL12
	}
	ns.ObjectOn(obj16)
	ns.ObjectOn(obj17)
	ns.ObjectOn(obj18)
	ns.ObjectOn(obj19)
	ns.ObjectOn(obj20)
	ns.ObjectOn(obj21)
	ns.ObjectOn(obj22)
	ns.ObjectOn(obj23)
	ns.ObjectOn(obj24)
LABEL12:
	goto LABEL5
LABEL5:
	return
}
func StartBarrels() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj40)
	ns.ObjectOn(obj41)
	ns.ObjectOn(obj42)
	ns.ObjectOn(obj43)
}
func BarrelHit1() {
	if !ns.IsCaller(obj33) {
		goto LABEL1
	}
	if flag49 {
		goto LABEL1
	}
	flag49 = true
	ns.WallGroupClose(gvar48)
	ns.ObjectOn(obj47)
	ns.ObjectOff(obj40)
	ns.ObjectOff(obj41)
	ns.ObjectOff(obj42)
	ns.ObjectOff(obj43)
LABEL1:
	return
}
func BarrelHit2() {
	if !ns.IsCaller(obj34) {
		goto LABEL1
	}
	if flag49 {
		goto LABEL1
	}
	flag49 = true
	ns.WallGroupClose(gvar48)
	ns.ObjectOff(obj40)
	ns.ObjectOff(obj41)
	ns.ObjectOff(obj42)
	ns.ObjectOff(obj43)
LABEL1:
	return
}
func BarrelHit3() {
	if !ns.IsCaller(obj35) {
		goto LABEL1
	}
	if flag49 {
		goto LABEL1
	}
	flag49 = true
	ns.WallGroupClose(gvar48)
	ns.ObjectOff(obj40)
	ns.ObjectOff(obj41)
	ns.ObjectOff(obj42)
	ns.ObjectOff(obj43)
LABEL1:
	return
}
func BarrelHit4() {
	if !ns.IsCaller(obj36) {
		goto LABEL1
	}
	if flag49 {
		goto LABEL1
	}
	flag49 = true
	ns.WallGroupClose(gvar48)
	ns.ObjectOff(obj40)
	ns.ObjectOff(obj41)
	ns.ObjectOff(obj42)
	ns.ObjectOff(obj43)
LABEL1:
	return
}
func BarrelsInit() {
	obj33 = ns.Object("Barrel1")
	obj34 = ns.Object("Barrel2")
	obj35 = ns.Object("Barrel3")
	obj36 = ns.Object("Barrel4")
	obj36 = ns.Object("Barrel5")
	obj36 = ns.Object("Barrel6")
	obj36 = ns.Object("Barrel7")
	obj40 = ns.Object("BarrelMover1")
	obj41 = ns.Object("BarrelMover2")
	obj42 = ns.Object("BarrelMover3")
	obj43 = ns.Object("BarrelMover4")
	obj44 = ns.Object("BarrelMover5")
	obj45 = ns.Object("BarrelMover6")
	obj46 = ns.Object("BarrelMover7")
	gvar48 = ns.WallGroup("BarrelWalls")
	obj47 = ns.Object("TempTeleporter")
	ns.ObjectOff(obj47)
}
func StartBlocker() {
	if flag56 {
		goto LABEL1
	}
	ns.ObjectOn(obj50)
	ns.ObjectOn(obj51)
	flag56 = true
LABEL1:
	return
}
func RaiseBlockerWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupClose(gvar54)
	ns.WallGroupOpen(gvar55)
	ns.WallOpen(gvar52)
	ns.WallOpen(gvar53)
}
func BlockerInit() {
	obj50 = ns.Object("SentryBlock")
	obj51 = ns.Object("SentryBlockMover")
	gvar52 = ns.Wall(173, 169)
	gvar53 = ns.Wall(177, 165)
	gvar54 = ns.WallGroup("BlockerWalls")
	gvar55 = ns.WallGroup("BlockerExit")
}
func FireDoomStatue() {
	if flag60 {
		goto LABEL1
	}
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj57, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
LABEL1:
	return
}
func DoomStatueOff() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar64)
	ns.WallOpen(gvar65)
	ns.WallOpen(gvar66)
}
func TakeShot() {
	var v0 int
	switch v0 = ivar71; v0 {
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
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj67, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar71 += 1
	goto LABEL5
LABEL2:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj68, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar71 += 1
	goto LABEL5
LABEL3:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj69, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar71 += 1
	goto LABEL5
LABEL4:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj70, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar71 = 1
	goto LABEL5
LABEL5:
	return
}
func KillEm() {
	var v0 int
	switch v0 = ivar78; v0 {
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
	default:
		goto LABEL7
	}
LABEL1:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj72, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar78 += 1
	goto LABEL7
LABEL2:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj73, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar78 += 1
	goto LABEL7
LABEL3:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj74, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar78 += 1
	goto LABEL7
LABEL4:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj75, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar78 += 1
	goto LABEL7
LABEL5:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj76, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar78 += 1
	goto LABEL7
LABEL6:
	ns.CastSpellObjectLocation(ns.SPELL_DEATH_RAY, obj77, ns.GetObjectX(ns.GetCaller()), ns.GetObjectY(ns.GetCaller()))
	ivar78 = 1
	goto LABEL7
LABEL7:
	return
}
func DeathRayInit() {
	obj57 = ns.Object("DoomStatue")
	gvar61 = ns.Wall(80, 14)
	gvar62 = ns.Wall(81, 13)
	gvar63 = ns.Wall(82, 12)
	gvar64 = ns.Wall(76, 10)
	gvar65 = ns.Wall(77, 9)
	gvar66 = ns.Wall(78, 8)
	fvar58 = ns.GetObjectX(obj57)
	fvar59 = ns.GetObjectY(obj57)
	obj67 = ns.Object("Death1")
	obj68 = ns.Object("Death2")
	obj69 = ns.Object("Death3")
	obj70 = ns.Object("Death4")
	obj72 = ns.Object("RayStatue1")
	obj73 = ns.Object("RayStatue2")
	obj74 = ns.Object("RayStatue3")
	obj75 = ns.Object("RayStatue4")
	obj76 = ns.Object("RayStatue5")
	obj77 = ns.Object("RayStatue6")
}
func OpenMonsterWalls() {
	var v0 int
	if !(ivar86 < 5) {
		goto LABEL1
	}
	ivar86 += 1
LABEL1:
	switch v0 = ivar86; v0 {
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
	ns.WallGroupOpen(gvar82)
	goto LABEL6
LABEL3:
	ns.WallGroupOpen(gvar85)
	goto LABEL6
LABEL4:
	ns.WallGroupOpen(gvar84)
	goto LABEL6
LABEL5:
	ns.WallGroupOpen(gvar83)
	goto LABEL6
LABEL6:
	return
}
func SpecialPlate() {
	var v0 int
	if ivar86 == 0 {
		goto LABEL1
	}
	if !(ivar86 < 5) {
		goto LABEL2
	}
	ivar86 += 1
LABEL2:
	switch v0 = ivar86; v0 {
	case 1:
		goto LABEL3
	case 2:
		goto LABEL4
	case 3:
		goto LABEL5
	case 4:
		goto LABEL6
	default:
		goto LABEL1
	}
LABEL3:
	ns.WallGroupOpen(gvar82)
	goto LABEL1
LABEL4:
	ns.WallGroupOpen(gvar85)
	goto LABEL1
LABEL5:
	ns.WallGroupOpen(gvar84)
	goto LABEL1
LABEL6:
	ns.WallGroupOpen(gvar83)
	goto LABEL1
LABEL1:
	if ivar86 != 0 {
		goto LABEL7
	}
	ns.WallOpen(gvar79)
	ns.WallOpen(gvar80)
	ns.WallOpen(gvar81)
LABEL7:
	return
}
func FourWallInit() {
	gvar79 = ns.Wall(98, 138)
	gvar80 = ns.Wall(99, 139)
	gvar81 = ns.Wall(100, 140)
	gvar82 = ns.WallGroup("MonsterWallsUL")
	gvar83 = ns.WallGroup("MonsterWallsUR")
	gvar84 = ns.WallGroup("MonsterWallsLL")
	gvar85 = ns.WallGroup("MonsterWallsLR")
}
func Infrawel() {
	ns.Enchant(ns.GetCaller(), ns.ENCHANT_INFRAVISION, 30)
	ns.Print("GeneralPrint:WishingWellPrintEnchantment")
}
func OpenInvisoWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar111)
	flag110 = true
	ns.ObjectOn(obj87)
	ns.ObjectOn(obj88)
}
func InvisoHurt() {
	ns.ObjectOn(ns.GetTrigger())
}
func InvisoMonster() {
	var v0 int
	if !(ivar109 < 5) {
		goto LABEL1
	}
	ns.Enchant(ns.GetCaller(), ns.ENCHANT_INVISIBLE, 0)
	ivar109 += 1
	switch v0 = ivar109; v0 {
	case 1:
		goto LABEL2
	case 2:
		goto LABEL3
	case 3:
		goto LABEL4
	case 4:
		goto LABEL5
	default:
		goto LABEL1
	}
LABEL2:
	obj89 = ns.GetCaller()
	goto LABEL1
LABEL3:
	obj90 = ns.GetCaller()
	goto LABEL1
LABEL4:
	obj91 = ns.GetCaller()
	goto LABEL1
LABEL5:
	obj92 = ns.GetCaller()
	goto LABEL1
LABEL1:
	return
}
func StartInviso() {
	ns.ObjectOff(ns.GetTrigger())
	InvisAll()
}
func GuardGenInviso() {
	ns.Enchant(ns.GetCaller(), ns.ENCHANT_INVISIBLE, 0)
}
func InvisAll() {
	if ns.HasEnchant(obj89, ns.ENCHANT_INVISIBLE) {
		goto LABEL1
	}
	ns.Enchant(obj89, ns.ENCHANT_INVISIBLE, 0)
LABEL1:
	if ns.HasEnchant(obj90, ns.ENCHANT_INVISIBLE) {
		goto LABEL2
	}
	ns.Enchant(obj90, ns.ENCHANT_INVISIBLE, 0)
LABEL2:
	if ns.HasEnchant(obj91, ns.ENCHANT_INVISIBLE) {
		goto LABEL3
	}
	ns.Enchant(obj91, ns.ENCHANT_INVISIBLE, 0)
LABEL3:
	if ns.HasEnchant(obj92, ns.ENCHANT_INVISIBLE) {
		goto LABEL4
	}
	ns.Enchant(obj92, ns.ENCHANT_INVISIBLE, 0)
LABEL4:
	if ns.HasEnchant(obj87, ns.ENCHANT_INVISIBLE) {
		goto LABEL5
	}
	ns.Enchant(obj87, ns.ENCHANT_INVISIBLE, 0)
LABEL5:
	if ns.HasEnchant(obj88, ns.ENCHANT_INVISIBLE) {
		goto LABEL6
	}
	ns.Enchant(obj88, ns.ENCHANT_INVISIBLE, 0)
LABEL6:
	return
}
func InvisGenInit() {
	obj87 = ns.Object("InvisoGen1")
	obj88 = ns.Object("InvisoGen2")
	obj99 = ns.Object("GuardGen1")
	obj100 = ns.Object("GuardGen2")
	obj93 = ns.Object("LastInviso1")
	obj94 = ns.Object("LastInviso2")
	obj101 = ns.Object("DeathGens1")
	obj102 = ns.Object("DeathGens2")
	obj103 = ns.Object("DeathGens3")
	obj104 = ns.Object("DeathGens4")
	obj105 = ns.Object("BooGen1")
	obj106 = ns.Object("BooGen2")
	obj107 = ns.Object("InvisoAlt1")
	obj108 = ns.Object("InvisoAlt2")
	gvar111 = ns.WallGroup("InvisoWalls")
	ns.Enchant(obj87, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj88, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj99, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj100, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj93, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj94, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj95, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj96, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj97, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj98, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj101, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj102, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj103, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj104, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj105, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj106, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj107, ns.ENCHANT_INVISIBLE, 0)
	ns.Enchant(obj108, ns.ENCHANT_INVISIBLE, 0)
	ns.ObjectOff(obj87)
	ns.ObjectOff(obj88)
	ns.ObjectOff(obj105)
	ns.ObjectOff(obj106)
}
func InvisoOthers() {
	if ns.IsCaller(ns.GetHost()) {
		goto LABEL1
	}
	ns.Enchant(ns.GetCaller(), ns.ENCHANT_INVISIBLE, 0)
LABEL1:
	return
}
func RevealGens() {
	ns.ObjectOn(obj105)
	ns.ObjectOn(obj106)
}
func StartCaveElevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.CreatureCageAppears, wp147)
	ns.AudioEvent(ns.ChangeSpellbar, wp147)
	ns.Print("GeneralPrint:MsgElevOn")
	ns.ObjectGroupOn(gvar146)
	ns.ObjectOn(obj112)
}
func StopKeyElevator() {
	ns.ObjectOff(obj113)
}
func StartKeyElevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.MoveObject(obj114, 2334, 5048)
	ns.AudioEvent(ns.CreatureCageAppears, wp148)
	ns.AudioEvent(ns.ChangeSpellbar, wp148)
	ns.ObjectOn(obj113)
	ns.FrameTimer(36, StopKeyElevator)
}
func OpenTombElevatorWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar117)
	ns.WallOpen(gvar118)
	ns.WallOpen(gvar119)
}
func StartAltBlocks() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj115)
	ns.ObjectOn(obj116)
}
func OpenLastTombWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar120)
	ns.WallOpen(gvar121)
	ns.WallOpen(gvar122)
	ns.WallOpen(gvar123)
	ns.WallOpen(gvar124)
}
func OpenFireWall1() {
	ns.WallOpen(gvar125)
}
func OpenFireWall2() {
	ns.WallOpen(gvar126)
}
func OpenFireWall3() {
	ns.WallOpen(gvar127)
}
func OpenEastWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar128)
	ns.WallOpen(gvar129)
	ns.WallOpen(gvar130)
}
func OpenSurpriseWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar131)
	ns.WallOpen(gvar132)
}
func OpenLotsaWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar134)
	ns.WallOpen(gvar135)
	ns.WallOpen(gvar136)
	ns.WallOpen(gvar137)
	ns.WallOpen(gvar138)
	ns.WallOpen(gvar139)
	ns.WallOpen(gvar140)
	ns.WallOpen(gvar141)
	ns.WallOpen(gvar142)
}
func OpenGolemWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ivar144 += 1
	if !(ivar144 >= 2) {
		goto LABEL1
	}
	ns.WallGroupOpen(gvar143)
	ns.Print("GeneralPrint:MsgWallLowered")
LABEL1:
	return
}
func LaunchSurprise() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar145)
}
func XGen1Die() {
	ns.ObjectOn(obj150)
}
func XGen2Die() {
	ns.ObjectOn(obj151)
}
func OneByOneInit() {
	gvar152 = ns.Wall(13, 177)
	gvar153 = ns.Wall(20, 170)
	gvar154 = ns.Wall(23, 187)
	gvar155 = ns.Wall(30, 180)
}
func PickOneInit() {
	gvar156 = ns.Wall(79, 171)
	gvar157 = ns.Wall(80, 170)
	gvar158 = ns.WallGroup("LeftWalls")
	gvar159 = ns.WallGroup("RightWalls")
	gvar160 = ns.WallGroup("CenterWalls")
}
func PitLineInit() {
	obj162 = ns.Object("PitLine1")
	obj163 = ns.Object("PitLine2")
	obj164 = ns.Object("PitLine3")
	obj165 = ns.Object("PitLine4")
	obj166 = ns.Object("PitLine5")
	obj167 = ns.Object("PitLine6")
	obj168 = ns.Object("PitLine7")
	obj169 = ns.Object("PitLine8")
	obj170 = ns.Object("PitLine9")
	obj171 = ns.Object("PitLine10")
	obj172 = ns.Object("PitLine11")
	obj173 = ns.Object("PitLine12")
	obj174 = ns.Object("PitLine13")
	obj175 = ns.Object("PitLine14")
	obj176 = ns.Object("PitLine15")
	obj177 = ns.Object("PitLine16")
	obj178 = ns.Object("PitLine17")
	obj179 = ns.Object("PitLine18")
	obj180 = ns.Object("PitLine19")
	obj181 = ns.Object("PitLine20")
	obj182 = ns.Object("PitLine21")
	obj183 = ns.Object("PitLine22")
	gvar185 = ns.Wall(172, 148)
	gvar186 = ns.Wall(171, 149)
	gvar187 = ns.Wall(172, 150)
	gvar188 = ns.Wall(173, 151)
	gvar184 = ns.Wall(189, 131)
	obj189 = ns.Object("HallGen1")
	obj190 = ns.Object("HallGen2")
	obj191 = ns.Object("HallGen3")
	obj192 = ns.Object("HallGen4")
	ns.ObjectOff(obj189)
	ns.ObjectOff(obj190)
	ns.ObjectOff(obj191)
}
func PusherInit() {
	obj193 = ns.Object("Pusher1")
	obj194 = ns.Object("Pusher2")
	obj195 = ns.Object("Pusher3")
	obj196 = ns.Object("Pusher4")
	obj197 = ns.Object("PusherMover1")
	obj198 = ns.Object("PusherMover2")
	obj199 = ns.Object("PusherMover3")
	obj200 = ns.Object("PusherMover4")
	gvar201 = ns.ObjectGroup("PusherPlates")
	wp202 = ns.Waypoint("PusherStart1")
	wp203 = ns.Waypoint("PusherStart2")
	wp204 = ns.Waypoint("PusherStart3")
	wp205 = ns.Waypoint("PusherStart4")
	wp206 = ns.Waypoint("PusherHome1")
	wp207 = ns.Waypoint("PusherHome2")
	wp208 = ns.Waypoint("PusherHome3")
	wp209 = ns.Waypoint("PusherHome4")
}
func SentMechInit() {
	obj210 = ns.Object("Fluffy")
	obj211 = ns.Object("MechSentry")
	obj212 = ns.Object("MechSentry2")
	gvar213 = ns.Wall(145, 61)
	gvar214 = ns.Wall(146, 60)
	gvar215 = ns.Wall(147, 59)
	gvar216 = ns.Wall(143, 59)
	gvar217 = ns.Wall(144, 58)
	gvar218 = ns.Wall(145, 57)
	gvar219 = ns.Wall(170, 60)
	gvar220 = ns.Wall(171, 59)
	wp221 = ns.Waypoint("MechSentrySound")
}
func ThreeSpkInit() {
	obj227 = ns.Object("RotSpike1")
	obj228 = ns.Object("RotSpike2")
	obj229 = ns.Object("RotSpike3")
	obj230 = ns.Object("RotSpike4")
	obj223 = ns.Object("RotSpikeMover1")
	obj224 = ns.Object("RotSpikeMover2")
	obj225 = ns.Object("RotSpikeMover3")
	obj226 = ns.Object("RotSpikeMover4")
	obj236 = ns.Object("KillSpikeMover1")
	obj237 = ns.Object("KillSpikeMover2")
	obj238 = ns.Object("KillSpikeMover3")
	obj239 = ns.Object("KillSpikeMover4")
	obj240 = ns.Object("KillSpikeMover5")
	obj231 = ns.Object("KillSpike1")
	obj232 = ns.Object("KillSpike2")
	obj233 = ns.Object("KillSpike3")
	obj234 = ns.Object("KillSpike4")
	obj235 = ns.Object("KillSpike5")
	gvar241 = ns.Wall(89, 117)
	gvar242 = ns.Wall(90, 118)
	gvar243 = ns.Wall(91, 119)
	gvar244 = ns.WallGroup("ReleaseWalls")
	obj245 = ns.Object("FRow1")
	obj246 = ns.Object("FRow2")
	obj247 = ns.Object("FRow3")
	obj248 = ns.Object("MRow1")
	obj249 = ns.Object("MRow2")
	obj250 = ns.Object("MRow3")
	obj251 = ns.Object("BRow1")
	obj252 = ns.Object("BRow2")
	obj253 = ns.Object("BRow3")
	ns.ObjectOff(obj248)
	ns.ObjectOff(obj249)
	ns.ObjectOff(obj250)
	ns.ObjectOff(obj251)
	ns.ObjectOff(obj252)
	ns.ObjectOff(obj253)
}
func TombGunsInit() {
	gvar254 = ns.ObjectGroup("TombGuns")
	gvar255 = ns.WallGroup("TombTrapWalls")
	wp256 = ns.Waypoint("GunWP")
}
func TombRayInit() {
	obj259 = ns.Object("TombRay")
	obj260 = ns.Object("TombRayMover")
	obj277 = ns.Object("PushingBlock")
	wp278 = ns.Waypoint("PushingWP")
	wp279 = ns.Waypoint("PushingHome")
	gvar263 = ns.Wall(120, 194)
	gvar264 = ns.Wall(121, 195)
	gvar265 = ns.Wall(117, 187)
	gvar266 = ns.Wall(118, 188)
	gvar267 = ns.Wall(119, 189)
	gvar268 = ns.WallGroup("TombTrapSentryWalls")
	gvar269 = ns.WallGroup("TombTrapFenceWalls")
	gvar270 = ns.WallGroup("TombTrapGenWalls")
	obj261 = ns.Object("StoneBlockSentry")
	obj262 = ns.Object("StoneBlockSentryMover")
	gvar273 = ns.ObjectGroup("SendTriggers")
	wp274 = ns.Waypoint("DestWP")
}
func WallUpInit() {
	gvar282 = ns.ObjectGroup("WallUpTriggers")
	gvar283 = ns.Wall(104, 40)
	gvar284 = ns.Wall(105, 41)
	gvar285 = ns.Wall(106, 42)
	gvar286 = ns.Wall(107, 43)
	gvar287 = ns.Wall(108, 44)
	gvar288 = ns.Wall(109, 45)
	gvar289 = ns.Wall(110, 46)
	gvar290 = ns.Wall(111, 47)
	wp291 = ns.Waypoint("WallUpAudioWP")
}
func XtraAreaInit() {
	obj294 = ns.Object("LeftSkull")
	obj295 = ns.Object("RightSkull")
	gvar296 = ns.Wall(234, 172)
	gvar297 = ns.Wall(235, 173)
	gvar298 = ns.Wall(236, 174)
	gvar299 = ns.Wall(8, 238)
	gvar300 = ns.Wall(9, 237)
	gvar301 = ns.Wall(10, 236)
	obj302 = ns.Object("FGen1")
	obj303 = ns.Object("FGen2")
	obj304 = ns.Object("FGen3")
	obj305 = ns.Object("FGenA1")
	obj306 = ns.Object("FGenA2")
	obj307 = ns.Object("FGenA3")
	ns.ObjectOff(obj303)
	ns.ObjectOff(obj304)
	ns.ObjectOff(obj306)
	ns.ObjectOff(obj307)
}
func MapInitialize() {
	obj112 = ns.Object("CaveElevator")
	gvar146 = ns.ObjectGroup("CaveGears")
	wp147 = ns.Waypoint("CaveGearWP")
	obj113 = ns.Object("KeyElevator")
	wp148 = ns.Waypoint("KeyElevatorWP")
	obj114 = ns.Object("EyeCandyKey")
	obj115 = ns.Object("AltBlockMover1")
	obj116 = ns.Object("AltBlockMover2")
	gvar117 = ns.Wall(47, 129)
	gvar118 = ns.Wall(48, 130)
	gvar119 = ns.Wall(47, 131)
	gvar120 = ns.Wall(225, 145)
	gvar121 = ns.Wall(226, 144)
	gvar122 = ns.Wall(227, 143)
	gvar123 = ns.Wall(228, 142)
	gvar124 = ns.Wall(229, 141)
	gvar125 = ns.Wall(75, 19)
	gvar126 = ns.Wall(69, 25)
	gvar127 = ns.Wall(63, 31)
	gvar128 = ns.Wall(35, 57)
	gvar129 = ns.Wall(36, 58)
	gvar130 = ns.Wall(37, 59)
	gvar131 = ns.Wall(22, 56)
	gvar132 = ns.Wall(23, 55)
	gvar133 = ns.WallGroup("LotsaWalls")
	gvar134 = ns.Wall(21, 51)
	gvar135 = ns.Wall(22, 52)
	gvar136 = ns.Wall(23, 53)
	gvar137 = ns.Wall(24, 54)
	gvar138 = ns.Wall(25, 55)
	gvar139 = ns.Wall(26, 56)
	gvar140 = ns.Wall(25, 53)
	gvar141 = ns.Wall(26, 52)
	gvar142 = ns.Wall(27, 51)
	gvar143 = ns.WallGroup("GolemWalls")
	gvar145 = ns.Wall(9, 191)
	obj149 = ns.Object("XGen1")
	obj150 = ns.Object("XGen2")
	obj151 = ns.Object("XGen3")
	ArrowsInit()
	BarrelsInit()
	BlockerInit()
	DeathRayInit()
	FourWallInit()
	InvisGenInit()
	OneByOneInit()
	PickOneInit()
	PitLineInit()
	PusherInit()
	SentMechInit()
	ThreeSpkInit()
	TombGunsInit()
	TombRayInit()
	WallUpInit()
	XtraAreaInit()
	ns.ObjectOff(obj150)
	ns.ObjectOff(obj151)
	ns.FrameTimer(36, StopKeyElevator)
}
func OpenULWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar152)
}
func OpenURWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar153)
}
func OpenLLWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar154)
}
func OpenLRWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar155)
}
func OpenLeft() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar158)
	ivar161 += 1
}
func OpenRight() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar159)
	ivar161 += 1
}
func OpenCenter() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar160)
	if ivar161 != 2 {
		goto LABEL1
	}
	ns.WallOpen(gvar156)
	ns.WallOpen(gvar157)
LABEL1:
	return
}
func PitHit1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj162)
}
func PitHit2() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj163)
}
func PitHit3() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj164)
}
func PitHit4() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj165)
}
func PitHit5() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj166)
}
func PitHit6() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj167)
}
func PitHit7() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj168)
}
func PitHit8() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj169)
}
func PitHit9() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj170)
}
func PitHit10() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj171)
}
func PitHit11() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj172)
}
func PitHit12() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj173)
}
func PitHit13() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj174)
}
func PitHit14() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj175)
}
func PitHit15() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj176)
}
func PitHit16() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj177)
}
func PitHit17() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj178)
}
func PitHit18() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj179)
}
func PitHit19() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj180)
}
func PitHit20() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj181)
}
func PitHit21() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj182)
}
func PitHit22() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj183)
}
func OpenPitLinePass() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar184)
	ns.WallOpen(gvar185)
	ns.WallOpen(gvar186)
	ns.WallOpen(gvar187)
	ns.WallOpen(gvar188)
}
func HallGen4Die() {
	ns.ObjectOn(obj191)
}
func HallGen3Die() {
	ns.ObjectOn(obj190)
}
func HallGen2Die() {
	ns.ObjectOn(obj189)
}
func StartPusher() {
	ns.ObjectGroupOff(gvar201)
	ns.Move(obj193, wp202)
	ns.Move(obj194, wp203)
	ns.Move(obj195, wp204)
	ns.Move(obj196, wp205)
	ns.FrameTimer(45, ReturnPusher)
}
func ReturnPusher() {
	ns.Move(obj193, wp206)
	ns.Move(obj194, wp207)
	ns.Move(obj195, wp208)
	ns.Move(obj196, wp209)
	ns.FrameTimer(45, ResetPusher)
}
func ResetPusher() {
	ns.ObjectGroupOn(gvar201)
}
func StartMechSentry() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj211)
	ns.ObjectOn(obj212)
}
func StopMechSentry() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOff(obj211)
	ns.ObjectOff(obj212)
	ns.AudioEvent(ns.Clank1, wp221)
	ns.Print("GeneralPrint:MsgTrapOff")
}
func OpenFluffyWall() {
	if flag222 {
		goto LABEL1
	}
	ns.WallOpen(gvar213)
	ns.WallOpen(gvar214)
	ns.WallOpen(gvar215)
	flag222 = true
LABEL1:
	if !ns.IsCaller(ns.GetHost()) {
		goto LABEL2
	}
	if ns.IsAttackedBy(obj210, ns.GetCaller()) {
		goto LABEL2
	}
	OpenFluffyExit()
LABEL2:
	return
}
func OpenOffWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar219)
	ns.WallOpen(gvar220)
}
func OpenFluffyExit() {
	ns.WallOpen(gvar216)
	ns.WallOpen(gvar217)
	ns.WallOpen(gvar218)
}
func StartThreeSpk() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj223)
	ns.ObjectOn(obj224)
	ns.ObjectOn(obj225)
	ns.ObjectOn(obj226)
}
func ButtonActivate() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar241)
	ns.WallOpen(gvar242)
	ns.WallOpen(gvar243)
}
func PlateActivate() {
	ns.WallGroupOpen(gvar244)
}
func StartKillSpike() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj236)
	ns.ObjectOn(obj237)
	ns.ObjectOn(obj238)
	ns.ObjectOn(obj239)
	ns.ObjectOn(obj240)
}
func FRow1Die() {
	ns.ObjectOn(obj248)
}
func FRow2Die() {
	ns.ObjectOn(obj249)
}
func FRow3Die() {
	ns.ObjectOn(obj250)
}
func MRow1Die() {
	ns.ObjectOn(obj251)
	ns.ObjectOn(obj249)
}
func MRow2Die() {
	ns.ObjectOn(obj252)
	ns.ObjectOn(obj248)
	ns.ObjectOn(obj250)
}
func MRow3Die() {
	ns.ObjectOn(obj253)
	ns.ObjectOn(obj249)
}
func BRow1Die() {
	ns.ObjectOn(obj252)
}
func BRow2Die() {
	ns.ObjectOn(obj251)
	ns.ObjectOn(obj253)
}
func BRow3Die() {
	ns.ObjectOn(obj252)
}
func KillGen() {
	if !(ns.IsCaller(obj231) || ns.IsCaller(obj232) || ns.IsCaller(obj233) || ns.IsCaller(obj234) || ns.IsCaller(obj235)) {
		goto LABEL1
	}
	ns.Damage(ns.GetTrigger(), obj231, 1000, 0)
LABEL1:
	return
}
func CountDown() {
	if !(ivar257 < 11) {
		goto LABEL1
	}
	ns.WallGroupClose(gvar255)
	ns.AudioEvent(ns.GenerateTick, wp256)
	ivar257 += 1
	gvar258 = ns.SecondTimer(1, CountDown)
	goto LABEL2
LABEL1:
	StartTombGuns()
LABEL2:
	return
}
func StopCountDown() {
	ns.CancelTimer(gvar258)
	ns.WallGroupOpen(gvar255)
	ivar257 = 0
}
func StartTombGuns() {
	ns.ObjectGroupOn(gvar254)
	ns.SecondTimer(5, StopTombGuns)
}
func StopTombGuns() {
	ns.ObjectGroupOff(gvar254)
	ns.WallGroupOpen(gvar255)
}
func StartTombRayTrap() {
	ns.WallGroupClose(gvar269)
	ns.WallGroupOpen(gvar268)
	ns.WallGroupOpen(gvar270)
	ns.LockDoor(obj271)
	ns.LockDoor(obj272)
	ns.ObjectOff(ns.GetTrigger())
	ns.FrameTimer(15, StartTheRay)
}
func StopTombRayTrap() {
	flag276 = true
	ns.ObjectOff(obj260)
	ns.ObjectOff(obj259)
	ns.WallOpen(gvar263)
	ns.WallOpen(gvar264)
	ns.WallOpen(gvar265)
	ns.WallOpen(gvar266)
	ns.WallOpen(gvar267)
	ns.UnlockDoor(obj271)
	ns.UnlockDoor(obj272)
}
func StartTheRay() {
	if flag276 {
		goto LABEL1
	}
	ns.ObjectOn(obj260)
	ns.ObjectOn(obj259)
LABEL1:
	return
}
func SendPlayers() {
	if !flag275 {
		goto LABEL1
	}
	ns.MoveObject(ns.GetCaller(), ns.GetWaypointX(wp274), ns.GetWaypointY(wp274))
	ns.ObjectGroupOff(gvar273)
LABEL1:
	return
}
func StartStoneBlockSentry() {
	ns.Move(obj277, wp278)
	ns.SecondTimer(1, FireSentryRay)
}
func FireSentryRay() {
	ns.ObjectOn(obj261)
	ns.ObjectOn(obj262)
}
func StopStoneBlockSentry() {
	ns.CancelTimer(gvar280)
	ns.CancelTimer(gvar281)
	ns.ObjectOff(obj261)
	ns.ObjectOff(obj262)
	ns.Move(obj277, wp279)
}
func SentryWait() {
	if !ns.IsCaller(obj261) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	gvar281 = ns.SecondTimer(4, StopStoneBlockSentry)
LABEL1:
	return
}
func StartWallUp() {
	ns.ObjectGroupOff(gvar282)
	Counting()
}
func Counting() {
	ns.AudioEvent(ns.GenerateTick, wp291)
	ivar293 += 1
	RaiseWall()
}
func RaiseWall() {
	var v0 int
	switch v0 = ivar293; v0 {
	case 2:
		goto LABEL1
	case 4:
		goto LABEL2
	case 6:
		goto LABEL3
	case 8:
		goto LABEL4
	case 10:
		goto LABEL5
	case 12:
		goto LABEL6
	case 14:
		goto LABEL7
	case 16:
		goto LABEL8
	default:
		goto LABEL9
	}
LABEL1:
	ns.WallClose(gvar283)
	goto LABEL9
LABEL2:
	ns.WallClose(gvar284)
	goto LABEL9
LABEL3:
	ns.WallClose(gvar285)
	goto LABEL9
LABEL4:
	ns.WallClose(gvar286)
	goto LABEL9
LABEL5:
	ns.WallClose(gvar287)
	goto LABEL9
LABEL6:
	ns.WallClose(gvar288)
	goto LABEL9
LABEL7:
	ns.WallClose(gvar289)
	goto LABEL9
LABEL8:
	ns.WallClose(gvar290)
	goto LABEL9
LABEL9:
	if !(ivar293 < 17) {
		goto LABEL10
	}
	gvar292 = ns.SecondTimer(1, Counting)
LABEL10:
	return
}
func StopWallUp() {
	ns.ObjectOff(ns.GetTrigger())
	ns.CancelTimer(gvar292)
	ns.WallOpen(gvar283)
	ns.WallOpen(gvar284)
	ns.WallOpen(gvar285)
	ns.WallOpen(gvar286)
	ns.WallOpen(gvar287)
	ns.WallOpen(gvar288)
	ns.WallOpen(gvar289)
	ns.WallOpen(gvar290)
}
func StartLeftSkull() {
	ns.ObjectOn(obj294)
}
func StopLeftSkull() {
	ns.ObjectOff(obj294)
}
func StartRightSkull() {
	ns.ObjectOn(obj295)
}
func StopRightSkull() {
	ns.ObjectOff(obj295)
	ns.ClearOwner(obj295)
}
func OpenLichWall() {
	ns.WallOpen(gvar296)
	ns.WallOpen(gvar297)
	ns.WallOpen(gvar298)
}
func OpenGoodyWall() {
	ns.WallOpen(gvar299)
	ns.WallOpen(gvar300)
	ns.WallOpen(gvar301)
}
func FGen1Die() {
	ns.ObjectOn(obj303)
}
func FGen2Die() {
	ns.ObjectOn(obj304)
}
func FGenA1Die() {
	ns.ObjectOn(obj306)
}
func FGenA2Die() {
	ns.ObjectOn(obj307)
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
