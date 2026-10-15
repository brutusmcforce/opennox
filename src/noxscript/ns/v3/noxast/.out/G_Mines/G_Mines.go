package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

var (
	obj4    ns.ObjectID
	obj5    ns.ObjectID
	obj6    ns.ObjectID
	obj7    ns.ObjectID
	obj8    ns.ObjectID
	wp9     ns.WaypointID
	wp10    ns.WaypointID
	wp11    ns.WaypointID
	wp12    ns.WaypointID
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
	obj25   ns.ObjectID
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
	gvar63  ns.ObjectGroupID
	obj64   ns.ObjectID
	obj65   ns.ObjectID
	obj66   ns.ObjectID
	obj67   ns.ObjectID
	obj68   ns.ObjectID
	obj69   ns.ObjectID
	gvar70  ns.WallGroupID
	gvar71  ns.WallGroupID
	gvar72  ns.WallGroupID
	gvar73  ns.WallGroupID
	gvar74  ns.WallGroupID
	gvar75  ns.WallGroupID
	gvar76  ns.WallGroupID
	gvar77  ns.WallGroupID
	gvar78  ns.WallGroupID
	gvar79  ns.WallID
	gvar80  int
	gvar81  int
	gvar82  ns.WallID
	gvar83  ns.WallID
	gvar84  ns.WallID
	gvar85  ns.WallID
	gvar86  ns.WallID
	gvar87  ns.WallID
	gvar88  ns.WallID
	gvar89  ns.WallID
	obj90   ns.ObjectID
	obj91   ns.ObjectID
	obj92   ns.ObjectID
	obj93   ns.ObjectID
	obj94   ns.ObjectID
	obj95   ns.ObjectID
	ivar96  int
	ivar97  int
	ivar98  int
	ivar99  int
	ivar100 int
	ivar101 int
	obj102  ns.ObjectID
	obj103  ns.ObjectID
	obj104  ns.ObjectID
	obj105  ns.ObjectID
	gvar106 ns.ObjectGroupID
	gvar107 ns.WallGroupID
	wp108   ns.WaypointID
	gvar109 ns.ObjectGroupID
	gvar110 ns.ObjectGroupID
	gvar111 ns.ObjectGroupID
	flag112 bool
	gvar113 ns.WallGroupID
	gvar114 ns.WallGroupID
	gvar115 ns.WallGroupID
	gvar116 ns.WallGroupID
	obj117  ns.ObjectID
	obj118  ns.ObjectID
	obj119  ns.ObjectID
	wp120   ns.WaypointID
	gvar121 ns.ObjectGroupID
	gvar122 ns.WallGroupID
	gvar123 ns.WallID
	gvar124 ns.WallID
	gvar125 ns.WallID
	gvar126 ns.WallGroupID
	gvar127 ns.WallID
	gvar128 ns.WallID
	gvar129 ns.WallID
	gvar130 ns.WallGroupID
	gvar131 ns.WallID
	gvar132 ns.WallID
	gvar133 ns.WallID
	gvar134 ns.WallGroupID
	flag135 bool
	flag136 bool
	flag137 bool
	flag138 bool
	flag139 bool
	obj140  ns.ObjectID
	obj141  ns.ObjectID
	obj142  ns.ObjectID
	obj143  ns.ObjectID
	gvar144 ns.ObjectGroupID
	gvar145 ns.ObjectGroupID
	gvar146 ns.ObjectGroupID
	gvar147 ns.ObjectGroupID
	gvar148 int
	gvar149 ns.WallGroupID
	gvar150 ns.WallGroupID
	gvar151 ns.WallID
	gvar152 ns.WallID
	gvar153 ns.WallID
	gvar154 ns.WallID
	gvar155 ns.WallID
	gvar156 ns.WallID
	gvar157 ns.WallID
	gvar158 ns.WallID
	gvar159 ns.WallID
	gvar160 ns.WallID
	gvar161 ns.WallID
	gvar162 ns.WallID
	gvar163 ns.WallID
	gvar164 ns.WallID
	gvar165 ns.WallID
	gvar166 ns.WallID
	gvar167 ns.WallID
	gvar168 ns.WallGroupID
	gvar169 ns.WallID
	gvar170 ns.WallID
	gvar171 ns.WallID
	wp172   ns.WaypointID
	wp173   ns.WaypointID
	wp174   ns.WaypointID
	wp175   ns.WaypointID
	wp176   ns.WaypointID
	ivar177 int
	flag178 bool
	gvar179 ns.WallGroupID
	obj180  ns.ObjectID
	gvar181 ns.WallGroupID
	gvar182 ns.WallID
	gvar183 ns.WallID
	gvar184 ns.WallID
	gvar185 ns.WallID
	gvar186 ns.WallID
	gvar187 ns.WallID
	gvar188 ns.WallID
	gvar189 ns.WallID
	gvar190 ns.WallID
	gvar191 ns.WallID
	gvar192 ns.WallID
	gvar193 ns.WallID
	gvar194 ns.WallID
	gvar195 ns.WallID
	gvar196 ns.WallGroupID
	gvar197 ns.WallID
	gvar198 ns.WallID
	gvar199 ns.WallID
	gvar200 ns.WallID
	gvar201 ns.WallID
	gvar202 ns.WallID
	gvar203 ns.WallID
	gvar204 ns.WallID
	gvar205 ns.WallID
	gvar206 ns.WallID
	gvar207 ns.WallID
	gvar208 ns.WallID
	gvar209 ns.WallID
	gvar210 ns.WallID
	flag211 bool
	wp212   ns.WaypointID
	obj213  ns.ObjectID
	obj214  ns.ObjectID
	wp215   ns.WaypointID
	obj216  ns.ObjectID
	obj217  ns.ObjectID
	obj218  ns.ObjectID
	obj219  ns.ObjectID
	obj220  ns.ObjectID
	obj221  ns.ObjectID
	obj222  ns.ObjectID
	obj223  ns.ObjectID
	flag224 bool
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
	obj241  ns.ObjectID
	obj242  ns.ObjectID
	obj243  ns.ObjectID
	obj244  ns.ObjectID
	obj245  ns.ObjectID
	obj246  ns.ObjectID
	obj247  ns.ObjectID
	obj248  ns.ObjectID
	obj249  ns.ObjectID
	obj250  ns.ObjectID
	obj251  ns.ObjectID
	obj252  ns.ObjectID
	obj253  ns.ObjectID
	obj254  ns.ObjectID
	gvar255 ns.WallGroupID
	gvar256 ns.WallID
	gvar257 ns.WallID
	gvar258 ns.WallID
	gvar259 ns.WallID
	gvar260 ns.WallID
	gvar261 ns.WallGroupID
	gvar262 ns.WallID
	gvar263 ns.WallID
	gvar264 ns.WallID
	gvar265 ns.WallID
	gvar266 ns.WallID
	gvar267 ns.WallID
	gvar268 ns.WallID
	gvar269 ns.WallID
	gvar270 ns.WallID
	gvar271 ns.WallID
	gvar272 ns.WallID
	gvar273 ns.WallID
	gvar274 ns.WallID
	gvar275 ns.WallID
	gvar276 ns.WallID
	gvar277 ns.WallID
	gvar278 ns.WallID
	gvar279 ns.WallID
	gvar280 ns.WallID
	gvar281 ns.WallID
	gvar282 ns.WallID
	obj283  ns.ObjectID
	obj284  ns.ObjectID
	obj285  ns.ObjectID
	obj286  ns.ObjectID
	obj287  ns.ObjectID
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
	wp301   ns.WaypointID
	wp302   ns.WaypointID
	wp303   ns.WaypointID
	wp304   ns.WaypointID
	obj305  ns.ObjectID
	obj306  ns.ObjectID
	obj307  ns.ObjectID
	wp308   ns.WaypointID
	wp309   ns.WaypointID
	wp310   ns.WaypointID
	wp311   ns.WaypointID
	obj312  ns.ObjectID
	obj313  ns.ObjectID
	gvar314 ns.WallGroupID
	gvar315 ns.WallID
	gvar316 ns.WallID
	gvar317 ns.WallID
	obj318  ns.ObjectID
	obj319  ns.ObjectID
	gvar320 ns.WallGroupID
	gvar321 ns.WallID
	gvar322 ns.WallID
	gvar323 ns.WallID
	gvar324 ns.WallGroupID
	gvar325 ns.WallID
	gvar326 ns.WallID
	gvar327 ns.WallID
	gvar328 ns.WallID
	gvar329 ns.WallID
	gvar330 ns.WallID
	obj331  ns.ObjectID
	obj332  ns.ObjectID
	obj333  ns.ObjectID
	obj334  ns.ObjectID
	gvar335 ns.ObjectGroupID
	wp336   ns.WaypointID
	wp337   ns.WaypointID
	wp338   ns.WaypointID
	wp339   ns.WaypointID
	wp340   ns.WaypointID
	gvar341 ns.WallID
	gvar342 ns.WallID
	gvar343 ns.WallID
	gvar344 ns.WallID
	gvar345 ns.WallID
	gvar346 ns.WallID
	gvar347 ns.WallID
	gvar348 ns.WallID
	gvar349 ns.WallID
	gvar350 ns.WallID
	gvar351 ns.WallID
	gvar352 ns.WallID
	gvar353 ns.WallID
	gvar354 ns.WallID
	gvar355 ns.WallID
	gvar356 ns.WallID
	gvar357 ns.WallID
	gvar358 ns.WallID
	gvar359 ns.WallGroupID
	gvar360 ns.WallGroupID
	gvar361 ns.TimerID
	ivar362 int
	obj363  ns.ObjectID
	wp364   ns.WaypointID
	obj365  ns.ObjectID
	obj366  ns.ObjectID
	obj367  ns.ObjectID
	obj368  ns.ObjectID
	obj369  ns.ObjectID
	obj370  ns.ObjectID
	obj371  ns.ObjectID
	obj372  ns.ObjectID
	obj373  ns.ObjectID
	obj374  ns.ObjectID
	flag375 bool
	flag376 bool
	flag377 bool
	flag378 bool
	gvar379 ns.WallID
	gvar380 ns.WallID
	gvar381 ns.WallID
	gvar382 ns.WallID
	gvar383 ns.WallID
	gvar384 ns.WallID
	gvar385 ns.WallID
	gvar386 ns.WallID
	gvar387 ns.WallID
	gvar388 ns.WallID
	gvar389 ns.WallID
	gvar390 ns.WallID
	gvar391 ns.WallID
	gvar392 ns.WallID
	gvar393 ns.WallID
	gvar394 ns.WallID
	gvar395 ns.WallID
	gvar396 ns.WallID
	gvar397 ns.WallID
	gvar398 ns.WallID
	gvar399 ns.WallID
	gvar400 ns.WallID
	gvar401 ns.WallID
	gvar402 ns.WallID
	gvar403 ns.WallID
	gvar404 ns.WallID
	gvar405 ns.WallID
	gvar406 ns.WallID
	gvar407 ns.WallID
	gvar408 ns.WallID
	gvar409 ns.WallID
	gvar410 ns.WallID
	gvar411 ns.WallID
	gvar412 ns.WallID
	gvar413 ns.WallID
	gvar414 ns.WallID
	gvar415 ns.WallID
	gvar416 ns.WallID
	gvar417 ns.WallID
	gvar418 ns.WallID
	gvar419 ns.WallID
	gvar420 ns.WallID
	gvar421 ns.WallID
	gvar422 ns.WallGroupID
	gvar423 ns.WallGroupID
	gvar424 ns.WallGroupID
	gvar425 ns.WallGroupID
	gvar426 ns.WallGroupID
	gvar427 ns.WallGroupID
	ivar428 int
	ivar429 int
	ivar430 int
	gvar431 ns.WallID
	gvar432 ns.WallID
	gvar433 ns.WallID
	gvar434 ns.WallID
	gvar435 ns.WallID
	gvar436 ns.WallID
	gvar437 ns.WallID
	gvar438 ns.WallID
	gvar439 ns.WallID
	gvar440 ns.WallID
	gvar441 ns.WallID
	gvar442 ns.WallID
	gvar443 ns.WallID
	gvar444 ns.WallID
	gvar445 ns.WallID
	gvar446 ns.WallID
	gvar447 ns.WallID
	gvar448 ns.WallID
	gvar449 ns.WallID
	gvar450 ns.WallID
	gvar451 ns.WallID
	gvar452 ns.WallID
	gvar453 ns.WallID
	gvar454 ns.WallID
	gvar455 ns.WallID
	obj456  ns.ObjectID
	wp457   ns.WaypointID
	obj458  ns.ObjectID
	obj459  ns.ObjectID
	obj460  ns.ObjectID
	obj461  ns.ObjectID
	obj462  ns.ObjectID
	obj463  ns.ObjectID
	obj464  ns.ObjectID
	obj465  ns.ObjectID
	gvar466 ns.ObjectGroupID
	wp467   ns.WaypointID
	gvar468 ns.WallGroupID
	gvar469 ns.WallGroupID
	ivar470 int
	ivar471 int
	gvar472 ns.TimerID
	flag473 bool
	flag474 bool
	wp475   ns.WaypointID
	fvar476 float32
	fvar477 float32
	wp478   ns.WaypointID
	wp479   ns.WaypointID
	wp480   ns.WaypointID
	wp481   ns.WaypointID
	wp482   ns.WaypointID
	wp483   ns.WaypointID
)

func init() {
	flag112 = true
	flag135 = false
	flag136 = false
	flag137 = false
	flag138 = false
	flag139 = false
	gvar148 = 0
	ivar177 = 0
	flag178 = false
	flag211 = true
	flag224 = false
	gvar361 = 0
	ivar362 = 0
	flag375 = false
	flag376 = false
	flag377 = false
	flag378 = false
	ivar428 = 0
	ivar429 = 0
	ivar430 = 0
	ivar470 = 0
	ivar471 = 0
	gvar472 = 0
	flag473 = false
	flag474 = false
	fvar476 = 0
	fvar477 = 0
}
func StartAltBlocks() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj4)
	ns.ObjectOn(obj5)
	ns.ObjectOn(obj6)
}
func StartAltBlocks2() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.SpikeBlockMove, wp9)
	ns.ObjectOn(obj7)
	ns.ObjectOn(obj8)
}
func Alt1Grind() {
	ns.AudioEvent(ns.SpikeBlockMove, wp10)
}
func Alt2Grind() {
	ns.AudioEvent(ns.SpikeBlockMove, wp11)
}
func Alt3Grind() {
	ns.AudioEvent(ns.SpikeBlockMove, wp12)
}
func AltBlocksInit() {
	obj4 = ns.Object("AltBlockMover1")
	obj5 = ns.Object("AltBlockMover2")
	obj6 = ns.Object("AltBlockMover3")
	obj7 = ns.Object("AltBlockMover4")
	obj8 = ns.Object("AltBlockMover5")
	wp9 = ns.Waypoint("Grind2WP")
	wp10 = ns.Waypoint("Alt1WP")
	wp11 = ns.Waypoint("Alt2WP")
	wp12 = ns.Waypoint("Alt3WP")
}
func APitsRelease1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj13)
}
func APitsRelease2() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj14)
}
func APitsRelease3() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj15)
}
func APitsRelease4() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj16)
}
func APitsRelease5() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj17)
}
func APitsRelease6() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj18)
}
func APitsRelease7() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj19)
}
func APitsRelease8() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj20)
}
func APitsRelease9() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj21)
}
func APitsRelease10() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj22)
}
func APitsRelease11() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj23)
}
func APitsRelease12() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj24)
}
func APitsRelease13() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj25)
}
func APitsRelease14() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj26)
}
func APitsRelease15() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj27)
}
func APitsRelease16() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj28)
}
func APitsRelease17() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj29)
}
func APitsRelease18() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj30)
}
func APitsRelease19() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj31)
}
func APitsRelease20() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj32)
}
func APitsRelease21() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj33)
}
func APitsRelease22() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj34)
}
func APitsRelease23() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj35)
}
func APitsRelease24() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj36)
}
func APitsRelease25() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj37)
}
func APitsRelease26() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj38)
}
func APitsRelease27() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj39)
}
func APitsRelease28() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj40)
}
func APitsRelease29() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj41)
}
func APitsRelease30() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj42)
}
func APitsRelease31() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj43)
}
func APitsRelease32() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj44)
}
func APitsRelease33() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj45)
}
func APitsRelease34() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj46)
}
func APitsRelease35() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj47)
}
func APitsRelease36() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj48)
}
func APitsRelease37() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj49)
}
func APitsRelease38() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj50)
}
func APitsRelease39() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj51)
}
func APitsRelease40() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj52)
}
func APitsRelease41() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj53)
}
func APitsRelease42() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj54)
}
func APitsRelease43() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj55)
}
func APitsRelease44() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj56)
}
func APitsRelease45() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj57)
}
func APitsRelease46() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj58)
}
func APitsRelease47() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj59)
}
func APitsRelease48() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj60)
}
func APitsRelease49() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj61)
}
func APitsRelease50() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj62)
}
func SealAPits() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectGroupOff(gvar63)
}
func APitsInit() {
	obj13 = ns.Object("APit1")
	obj14 = ns.Object("APit2")
	obj15 = ns.Object("APit3")
	obj16 = ns.Object("APit4")
	obj17 = ns.Object("APit5")
	obj18 = ns.Object("APit6")
	obj19 = ns.Object("APit7")
	obj20 = ns.Object("APit8")
	obj21 = ns.Object("APit9")
	obj22 = ns.Object("APit10")
	obj23 = ns.Object("APit11")
	obj24 = ns.Object("APit12")
	obj25 = ns.Object("APit13")
	obj26 = ns.Object("APit14")
	obj27 = ns.Object("APit15")
	obj28 = ns.Object("APit16")
	obj29 = ns.Object("APit17")
	obj30 = ns.Object("APit18")
	obj31 = ns.Object("APit19")
	obj32 = ns.Object("APit20")
	obj33 = ns.Object("APit21")
	obj34 = ns.Object("APit22")
	obj35 = ns.Object("APit23")
	obj36 = ns.Object("APit24")
	obj37 = ns.Object("APit25")
	obj38 = ns.Object("APit26")
	obj39 = ns.Object("APit27")
	obj40 = ns.Object("APit28")
	obj41 = ns.Object("APit29")
	obj42 = ns.Object("APit30")
	obj43 = ns.Object("APit31")
	obj44 = ns.Object("APit32")
	obj45 = ns.Object("APit33")
	obj46 = ns.Object("APit34")
	obj47 = ns.Object("APit35")
	obj48 = ns.Object("APit36")
	obj49 = ns.Object("APit37")
	obj50 = ns.Object("APit38")
	obj51 = ns.Object("APit39")
	obj52 = ns.Object("APit40")
	obj53 = ns.Object("APit41")
	obj54 = ns.Object("APit42")
	obj55 = ns.Object("APit43")
	obj56 = ns.Object("APit44")
	obj57 = ns.Object("APit45")
	obj58 = ns.Object("APit46")
	obj59 = ns.Object("APit47")
	obj60 = ns.Object("APit48")
	obj61 = ns.Object("APit49")
	obj62 = ns.Object("APit50")
	gvar63 = ns.ObjectGroup("APits")
}
func StartBreakSpikes() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj67)
	ns.ObjectOn(obj68)
	ns.ObjectOn(obj69)
}
func Smash1() {
	if !ns.IsCaller(obj65) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar70)
LABEL1:
	return
}
func Smash2() {
	if !ns.IsCaller(obj65) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar71)
LABEL1:
	return
}
func Smash3() {
	if !ns.IsCaller(obj65) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar72)
LABEL1:
	return
}
func Smash4() {
	if !ns.IsCaller(obj65) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar73)
LABEL1:
	return
}
func Smash5() {
	if !ns.IsCaller(obj65) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar74)
LABEL1:
	return
}
func Smash6() {
	if !ns.IsCaller(obj65) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar75)
LABEL1:
	return
}
func Smash7() {
	if !ns.IsCaller(obj65) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar76)
LABEL1:
	return
}
func Smash8() {
	if !ns.IsCaller(obj65) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar77)
LABEL1:
	return
}
func Smash9() {
	if !ns.IsCaller(obj65) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar78)
LABEL1:
	return
}
func KillGen() {
	if !(ns.IsCaller(obj64) || ns.IsCaller(obj65) || ns.IsCaller(obj66)) {
		goto LABEL1
	}
	ns.Damage(ns.GetTrigger(), obj65, 1000, 0)
LABEL1:
	return
}
func BuzzSawInit() {
	obj64 = ns.Object("BreakSpike1")
	obj65 = ns.Object("BreakSpike2")
	obj66 = ns.Object("BreakSpike3")
	obj67 = ns.Object("BreakMover1")
	obj68 = ns.Object("BreakMover2")
	obj69 = ns.Object("BreakMover3")
	gvar70 = ns.WallGroup("Break1")
	gvar71 = ns.WallGroup("Break2")
	gvar72 = ns.WallGroup("Break3")
	gvar73 = ns.WallGroup("Break4")
	gvar74 = ns.WallGroup("Break5")
	gvar75 = ns.WallGroup("Break6")
	gvar76 = ns.WallGroup("Break7")
	gvar77 = ns.WallGroup("Break8")
	gvar78 = ns.WallGroup("Break9")
	gvar79 = ns.Wall(43, 73)
	gvar82 = ns.Wall(44, 74)
	gvar86 = ns.Wall(45, 75)
	gvar82 = ns.Wall(44, 70)
	gvar83 = ns.Wall(45, 71)
	gvar84 = ns.Wall(46, 72)
	gvar85 = ns.Wall(47, 73)
	gvar86 = ns.Wall(46, 68)
	gvar87 = ns.Wall(47, 69)
	gvar88 = ns.Wall(48, 70)
	gvar89 = ns.Wall(49, 71)
	obj90 = ns.Object("SpikeGen1")
	obj91 = ns.Object("SpikeGen2")
	obj92 = ns.Object("SpikeGen3")
	obj93 = ns.Object("SpikeGen4")
	obj94 = ns.Object("SpikeGen5")
	obj95 = ns.Object("SpikeGen6")
	ivar96 = ns.CurrentHealth(obj90)
	ivar97 = ns.CurrentHealth(obj91)
	ivar98 = ns.CurrentHealth(obj92)
	ivar99 = ns.CurrentHealth(obj93)
	ivar100 = ns.CurrentHealth(obj94)
	ivar101 = ns.CurrentHealth(obj95)
}
func StartSouthBlocks() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.SpikeBlockMove, wp108)
	ns.ObjectOn(obj102)
	ns.ObjectOn(obj103)
}
func StartEastBlocks() {
	ns.ObjectGroupOff(gvar106)
	ns.ObjectOn(obj104)
	ns.ObjectOn(obj105)
}
func SealSouthWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupClose(gvar107)
}
func CrisCrosInit() {
	obj102 = ns.Object("SpikeSouthMover1")
	obj103 = ns.Object("SpikeSouthMover2")
	obj104 = ns.Object("SpikeEastMover1")
	obj105 = ns.Object("SpikeEastMover2")
	gvar106 = ns.ObjectGroup("EastBlockTriggers")
	gvar107 = ns.WallGroup("SpikeWallsSouth")
	wp108 = ns.Waypoint("GrindWP")
}
func StopFirstElevator() {
	ns.ObjectOff(obj117)
}
func StartFirstElevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj118)
	ns.ObjectOn(obj119)
	ns.ObjectOn(obj117)
	ns.AudioEvent(ns.CreatureCageAppears, wp120)
	ns.AudioEvent(ns.ChangeSpellbar, wp120)
	ns.Print("GeneralPrint:MsgElevOn")
}
func GoodyPitSwitch() {
	if !flag112 {
		goto LABEL1
	}
	ns.ObjectGroupOff(gvar109)
	flag112 = false
	ns.WallGroupOpen(gvar114)
	goto LABEL2
LABEL1:
	ns.ObjectGroupOn(gvar109)
	flag112 = true
	ns.WallGroupClose(gvar114)
LABEL2:
	return
}
func GoodySwitch() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar113)
}
func GoodyPitInit() {
	gvar109 = ns.ObjectGroup("GoodiePits")
	gvar110 = ns.ObjectGroup("HiddenPits")
	gvar111 = ns.ObjectGroup("GoodiePitTriggers")
	gvar113 = ns.WallGroup("WrongWalls")
	gvar114 = ns.WallGroup("RightWalls")
	gvar115 = ns.WallGroup("SortaWalls")
	gvar116 = ns.WallGroup("GoodyPitWalls")
	obj117 = ns.Object("FirstElevator")
	obj118 = ns.Object("FirstGearA")
	obj119 = ns.Object("FirstGearB")
	wp120 = ns.Waypoint("GearSoundA")
	ns.ObjectOff(obj118)
	ns.ObjectOff(obj119)
	ns.FrameTimer(36, StopFirstElevator)
}
func SmashGoodyPitWalls() {
	ns.ObjectGroupOff(gvar111)
	ns.WallGroupBreak(gvar116)
}
func OpenSortaWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar115)
	ns.WallGroupOpen(gvar114)
}
func StartPuzzle() {
	if !(flag136 == true && flag137 == true) {
		goto LABEL1
	}
	flag135 = true
	ns.WallClose(gvar123)
	ns.WallClose(gvar124)
	ns.WallClose(gvar125)
	ns.WallClose(gvar127)
	ns.WallClose(gvar128)
	ns.WallClose(gvar129)
	ns.WallOpen(gvar131)
	ns.WallOpen(gvar132)
	ns.WallOpen(gvar133)
LABEL1:
	return
}
func PlateAActivate() {
	flag136 = true
	if !flag135 {
		goto LABEL1
	}
	ns.WallOpen(gvar123)
	ns.WallOpen(gvar124)
	ns.WallOpen(gvar125)
LABEL1:
	if flag135 {
		goto LABEL2
	}
	StartPuzzle()
LABEL2:
	return
}
func PlateBActivate() {
	flag137 = true
	if !flag135 {
		goto LABEL1
	}
	ns.WallOpen(gvar123)
	ns.WallOpen(gvar124)
	ns.WallOpen(gvar125)
LABEL1:
	if flag135 {
		goto LABEL2
	}
	StartPuzzle()
LABEL2:
	return
}
func PlateARelease() {
	flag136 = false
	if !(flag135 == true && flag137 != true && flag138 != true && flag139 != true) {
		goto LABEL1
	}
	ns.WallClose(gvar123)
	ns.WallClose(gvar124)
	ns.WallClose(gvar125)
LABEL1:
	return
}
func PlateBRelease() {
	flag137 = false
	if !(flag135 == true && flag136 != true && flag138 != true && flag139 != true) {
		goto LABEL1
	}
	ns.WallClose(gvar123)
	ns.WallClose(gvar124)
	ns.WallClose(gvar125)
LABEL1:
	return
}
func PlateCActivate() {
	flag138 = true
	if !flag135 {
		goto LABEL1
	}
	ns.WallOpen(gvar123)
	ns.WallOpen(gvar124)
	ns.WallOpen(gvar125)
LABEL1:
	if !flag135 {
		goto LABEL2
	}
	if !(flag138 == true && flag139 == true) {
		goto LABEL2
	}
	ns.ObjectGroupOff(gvar121)
	ns.WallOpen(gvar123)
	ns.WallOpen(gvar124)
	ns.WallOpen(gvar125)
	ns.WallOpen(gvar127)
	ns.WallOpen(gvar128)
	ns.WallOpen(gvar129)
LABEL2:
	return
}
func PlateDActivate() {
	flag139 = true
	if !flag135 {
		goto LABEL1
	}
	ns.WallOpen(gvar123)
	ns.WallOpen(gvar124)
	ns.WallOpen(gvar125)
LABEL1:
	if !flag135 {
		goto LABEL2
	}
	if !(flag138 == true && flag139 == true) {
		goto LABEL2
	}
	ns.ObjectGroupOff(gvar121)
	ns.WallOpen(gvar123)
	ns.WallOpen(gvar124)
	ns.WallOpen(gvar125)
	ns.WallOpen(gvar127)
	ns.WallOpen(gvar128)
	ns.WallOpen(gvar129)
LABEL2:
	return
}
func PlateCRelease() {
	flag138 = false
	if !(flag135 == true && flag136 != true && flag137 != true && flag139 != true) {
		goto LABEL1
	}
	ns.WallClose(gvar123)
	ns.WallClose(gvar124)
	ns.WallClose(gvar125)
LABEL1:
	return
}
func PlateDRelease() {
	flag139 = false
	if !(flag135 == true && flag136 != true && flag137 != true && flag138 != true) {
		goto LABEL1
	}
	ns.WallOpen(gvar123)
	ns.WallOpen(gvar124)
	ns.WallOpen(gvar125)
LABEL1:
	return
}
func OpenBooWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupBreak(gvar134)
}
func HallWorkInit() {
	gvar121 = ns.ObjectGroup("PuzzleTriggers")
	gvar122 = ns.WallGroup("PuzzleWallA")
	gvar123 = ns.Wall(34, 120)
	gvar124 = ns.Wall(35, 121)
	gvar125 = ns.Wall(36, 122)
	gvar126 = ns.WallGroup("PuzzleWallB")
	gvar127 = ns.Wall(26, 128)
	gvar128 = ns.Wall(27, 129)
	gvar129 = ns.Wall(28, 130)
	gvar130 = ns.WallGroup("PuzzleWallBoo")
	gvar131 = ns.Wall(36, 116)
	gvar132 = ns.Wall(37, 115)
	gvar133 = ns.Wall(38, 114)
	gvar134 = ns.WallGroup("BooWalls")
}
func OpenJail1() {
	if !ns.IsLocked(obj140) {
		goto LABEL1
	}
	ns.UnlockDoor(obj140)
LABEL1:
	return
}
func OpenJail2() {
	if !ns.IsLocked(obj141) {
		goto LABEL1
	}
	ns.UnlockDoor(obj141)
LABEL1:
	return
}
func OpenJail3() {
	if !ns.IsLocked(obj142) {
		goto LABEL1
	}
	ns.UnlockDoor(obj142)
LABEL1:
	return
}
func OpenJail4() {
	if !ns.IsLocked(obj143) {
		goto LABEL1
	}
	ns.UnlockDoor(obj143)
LABEL1:
	return
}
func JailRoomInit() {
	obj140 = ns.Object("Jail1")
	obj141 = ns.Object("Jail2")
	obj142 = ns.Object("Jail3")
	obj143 = ns.Object("Jail4")
	ns.LockDoor(obj140)
	ns.LockDoor(obj141)
	ns.LockDoor(obj142)
	ns.LockDoor(obj143)
}
func JumpPitsInit() {
	gvar144 = ns.ObjectGroup("JumpPits1")
	gvar145 = ns.ObjectGroup("JumpPits2")
	gvar146 = ns.ObjectGroup("JumpPits3")
	gvar147 = ns.ObjectGroup("JumpPits4")
}
func StartTheMovers() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj216)
	ns.ObjectOn(obj217)
	ns.ObjectOn(obj218)
	ns.ObjectOn(obj219)
	ns.ObjectOn(obj220)
	ns.ObjectOn(obj221)
}
func KillMover1() {
	ns.ObjectOff(obj220)
}
func KillMover2() {
	ns.ObjectOff(obj221)
}
func OpenSouthMagicWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar149)
}
func OpenNorthMagicWall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar151)
	ns.WallOpen(gvar152)
	ns.WallOpen(gvar153)
	ns.WallOpen(gvar154)
	ns.WallOpen(gvar155)
	ns.WallOpen(gvar156)
	ns.WallOpen(gvar157)
	ns.WallOpen(gvar158)
	ns.WallOpen(gvar159)
	ns.WallOpen(gvar160)
	ns.WallOpen(gvar161)
	ns.WallOpen(gvar162)
	ns.WallOpen(gvar163)
	ns.WallOpen(gvar164)
	ns.WallOpen(gvar165)
	ns.WallOpen(gvar166)
	ns.WallOpen(gvar167)
}
func StartAnkhRoom() {
	ns.ObjectOff(ns.GetTrigger())
	AnkhCountDown()
}
func AnkhCountDown() {
	if !(ivar177 < 10 && flag178 == false) {
		goto LABEL1
	}
	ns.AudioEvent(ns.GenerateTick, wp172)
	ns.AudioEvent(ns.GenerateTick, wp173)
	ns.AudioEvent(ns.GenerateTick, wp174)
	ns.AudioEvent(ns.GenerateTick, wp175)
	ns.AudioEvent(ns.GenerateTick, wp176)
	ivar177 += 1
	ns.SecondTimer(1, StartAnkhRoom)
	goto LABEL2
LABEL1:
	if flag178 {
		goto LABEL2
	}
	ns.WallClose(gvar169)
	ns.WallClose(gvar170)
	ns.WallClose(gvar171)
	ns.Print("GeneralPrint:MsgWallRaised")
LABEL2:
	return
}
func TouchAnkh() {
	ns.ObjectOff(ns.GetTrigger())
	flag178 = true
}
func BoulderTrigger() {
	obj223 = ns.GetCaller()
	flag224 = true
}
func BoulderRelease() {
	flag224 = false
}
func DragonShoot() {
	if !flag224 {
		goto LABEL1
	}
	ns.CastSpellObjectObject(ns.SPELL_FIREBALL, obj222, obj223)
LABEL1:
	return
}
func OpenReturnWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar179)
}
func StartRockElevator() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj180)
}
func OpenGoodWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar182)
	ns.WallOpen(gvar183)
	ns.WallOpen(gvar184)
	ns.WallOpen(gvar185)
	ns.WallOpen(gvar186)
	ns.WallOpen(gvar187)
	ns.WallOpen(gvar188)
	ns.WallOpen(gvar189)
	ns.WallOpen(gvar190)
	ns.WallOpen(gvar191)
	ns.WallOpen(gvar192)
	ns.WallOpen(gvar193)
	ns.WallOpen(gvar194)
	ns.WallOpen(gvar195)
}
func OpenBadWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar197)
	ns.WallOpen(gvar198)
	ns.WallOpen(gvar199)
	ns.WallOpen(gvar200)
	ns.WallOpen(gvar201)
	ns.WallOpen(gvar202)
	ns.WallOpen(gvar203)
	ns.WallOpen(gvar204)
	ns.WallOpen(gvar205)
	ns.WallOpen(gvar206)
	ns.WallOpen(gvar207)
	ns.WallOpen(gvar208)
	ns.WallOpen(gvar209)
	ns.WallOpen(gvar210)
}
func MakeAWish() {
	var v0 int
	if !flag211 {
		goto LABEL1
	}
	v0 = ns.MaxHealth(ns.GetCaller())
	ns.RestoreHealth(ns.GetCaller(), v0-ns.CurrentHealth(ns.GetCaller()))
	ns.AudioEvent(ns.RestoreHealthName, wp212)
	ns.Print("GeneralPrint:WellSignRefresh")
	flag211 = true
LABEL1:
	return
}
func ResetWell() {
	flag211 = true
}
func LastBlockStart() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj214)
}
func LastBlockStop() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOff(obj214)
	ns.Move(obj213, wp215)
}
func PitExitInit() {
	obj225 = ns.Object("FinalPit1")
	obj226 = ns.Object("FinalPit2")
	obj227 = ns.Object("FinalPit3")
	obj228 = ns.Object("FinalPit4")
	obj229 = ns.Object("FinalPit5")
	obj230 = ns.Object("FinalPit6")
	obj231 = ns.Object("FinalPit7")
	obj232 = ns.Object("FinalPit8")
	obj233 = ns.Object("FinalPit9")
	obj234 = ns.Object("FinalPit10")
	obj235 = ns.Object("FinalPit11")
	obj236 = ns.Object("FinalPit12")
	obj237 = ns.Object("FinalPit13")
	obj238 = ns.Object("FinalPit14")
	obj239 = ns.Object("FinalPit15")
	obj240 = ns.Object("FinalPit16")
	obj241 = ns.Object("FinalPit17")
	obj242 = ns.Object("FinalPit18")
	obj243 = ns.Object("FinalPit19")
	obj244 = ns.Object("FinalPit20")
	obj245 = ns.Object("FinalPit21")
	obj246 = ns.Object("FinalPit22")
	obj247 = ns.Object("FinalPit23")
	obj248 = ns.Object("FinalPit24")
	obj249 = ns.Object("FinalPit25")
	obj250 = ns.Object("FinalPit26")
	obj251 = ns.Object("FinalPit27")
	obj252 = ns.Object("FinalPit28")
	obj253 = ns.Object("FinalPit29")
	obj254 = ns.Object("FinalPit30")
	gvar255 = ns.WallGroup("SurfaceWalls")
	gvar256 = ns.Wall(54, 204)
	gvar257 = ns.Wall(55, 203)
	gvar258 = ns.Wall(43, 209)
	gvar259 = ns.Wall(44, 210)
	gvar260 = ns.Wall(45, 211)
	gvar261 = ns.WallGroup("GoodyWalls")
	gvar262 = ns.Wall(62, 212)
	gvar263 = ns.Wall(63, 211)
	gvar264 = ns.Wall(55, 221)
	gvar265 = ns.Wall(56, 222)
	gvar266 = ns.Wall(57, 223)
	gvar267 = ns.Wall(45, 229)
	gvar268 = ns.Wall(46, 228)
	gvar269 = ns.Wall(47, 227)
	gvar270 = ns.Wall(48, 226)
	gvar271 = ns.Wall(49, 225)
	gvar272 = ns.Wall(50, 224)
	gvar273 = ns.Wall(51, 223)
	gvar274 = ns.Wall(37, 221)
	gvar275 = ns.Wall(38, 220)
	gvar276 = ns.Wall(39, 219)
	gvar277 = ns.Wall(40, 218)
	gvar278 = ns.Wall(41, 217)
	gvar279 = ns.Wall(42, 216)
	gvar280 = ns.Wall(43, 215)
	gvar281 = ns.Wall(44, 214)
	gvar282 = ns.Wall(45, 213)
}
func PusherInit() {
	obj283 = ns.Object("PushSpikeMover1")
	obj284 = ns.Object("PushSpikeMover2")
	obj285 = ns.Object("PushSpikeMover3")
	obj286 = ns.Object("PushSpikeMover4")
	obj287 = ns.Object("PushSpikeMover5")
	obj288 = ns.Object("PushSpikeMover6")
	obj289 = ns.Object("PushSpikeMover7")
	obj290 = ns.Object("PushSpikeMover8")
	obj291 = ns.Object("PushBlock1")
	obj292 = ns.Object("PushBlock2")
	obj293 = ns.Object("PushBlock12Trigger")
	wp294 = ns.Waypoint("PushBlockHome1")
	wp295 = ns.Waypoint("PushBlockHome2")
	wp296 = ns.Waypoint("PushBlockStart1")
	wp297 = ns.Waypoint("PushBlockStart2")
	obj298 = ns.Object("PushBlock3")
	obj299 = ns.Object("PushBlock4")
	obj300 = ns.Object("PushBlock34Trigger")
	wp301 = ns.Waypoint("PushBlockHome3")
	wp302 = ns.Waypoint("PushBlockHome4")
	wp303 = ns.Waypoint("PushBlockStart3")
	wp304 = ns.Waypoint("PushBlockStart4")
	obj305 = ns.Object("PushBlock5")
	obj306 = ns.Object("PushBlock6")
	obj307 = ns.Object("PushBlock56Trigger")
	wp308 = ns.Waypoint("PushBlockHome5")
	wp309 = ns.Waypoint("PushBlockHome6")
	wp310 = ns.Waypoint("PushBlockStart5")
	wp311 = ns.Waypoint("PushBlockStart6")
	obj312 = ns.Object("MovingRay")
	obj313 = ns.Object("MovingRayMover")
	gvar314 = ns.WallGroup("MovingRayWalls")
	gvar315 = ns.Wall(166, 192)
	gvar316 = ns.Wall(167, 191)
	gvar317 = ns.Wall(168, 190)
	obj318 = ns.Object("SpinSpikeMover1")
	obj319 = ns.Object("SpinSpikeMover2")
	gvar320 = ns.WallGroup("PusherExitWalls1")
	gvar321 = ns.Wall(148, 174)
	gvar322 = ns.Wall(149, 173)
	gvar323 = ns.Wall(150, 172)
	gvar324 = ns.WallGroup("PusherExitWalls2")
	gvar325 = ns.Wall(152, 178)
	gvar326 = ns.Wall(153, 177)
	gvar327 = ns.Wall(154, 176)
	gvar328 = ns.Wall(124, 198)
	gvar328 = ns.Wall(125, 197)
	gvar328 = ns.Wall(126, 196)
}
func RayTrapsInit() {
	gvar335 = ns.ObjectGroup("SentryTrapTriggers")
	obj331 = ns.Object("SentryRay1")
	obj332 = ns.Object("SentryRay2")
	obj333 = ns.Object("SentryBlockMover1")
	obj334 = ns.Object("SentryBlockMover2")
	wp336 = ns.Waypoint("SentryBlockHome1")
	wp337 = ns.Waypoint("SentryBlockHome2")
	wp338 = ns.Waypoint("SentryBlockStart1")
	wp339 = ns.Waypoint("SentryBlockStart2")
	wp340 = ns.Waypoint("RayTrapWP")
	gvar341 = ns.Wall(119, 105)
	gvar342 = ns.Wall(120, 106)
	gvar343 = ns.Wall(121, 107)
	gvar344 = ns.Wall(120, 104)
	gvar345 = ns.Wall(121, 105)
	gvar346 = ns.Wall(122, 106)
	gvar347 = ns.Wall(121, 103)
	gvar348 = ns.Wall(122, 104)
	gvar349 = ns.Wall(123, 105)
	gvar350 = ns.Wall(122, 102)
	gvar351 = ns.Wall(123, 103)
	gvar352 = ns.Wall(124, 104)
	gvar353 = ns.Wall(123, 101)
	gvar354 = ns.Wall(124, 102)
	gvar355 = ns.Wall(125, 103)
	gvar356 = ns.Wall(124, 100)
	gvar357 = ns.Wall(125, 101)
	gvar358 = ns.Wall(126, 102)
	gvar359 = ns.WallGroup("SentryTrapWalls")
	gvar360 = ns.WallGroup("StopRayWalls")
	obj363 = ns.Object("RayElevator")
	wp364 = ns.Waypoint("RayWP")
	ns.FrameTimer(36, StopRayElevator)
}
func Secret1Init() {
	obj365 = ns.Object("Door1")
	obj366 = ns.Object("Door2")
	obj367 = ns.Object("Door3")
	obj368 = ns.Object("Door4")
	obj369 = ns.Object("LeftGen1")
	obj370 = ns.Object("LeftGen2")
	obj371 = ns.Object("LeftGen3")
	obj372 = ns.Object("RightGen1")
	obj373 = ns.Object("RightGen2")
	obj374 = ns.Object("RightGen3")
	gvar379 = ns.Wall(67, 75)
	gvar380 = ns.Wall(68, 74)
	gvar381 = ns.Wall(69, 73)
	gvar382 = ns.Wall(67, 77)
	gvar383 = ns.Wall(68, 78)
	gvar384 = ns.Wall(69, 77)
	gvar385 = ns.Wall(70, 76)
	gvar386 = ns.Wall(71, 77)
	gvar387 = ns.Wall(71, 75)
	gvar388 = ns.Wall(72, 74)
	gvar389 = ns.Wall(71, 73)
	gvar390 = ns.Wall(59, 91)
	gvar391 = ns.Wall(60, 90)
	gvar392 = ns.Wall(61, 89)
	gvar393 = ns.Wall(62, 88)
	gvar394 = ns.Wall(63, 87)
	gvar395 = ns.Wall(73, 77)
	gvar396 = ns.Wall(74, 76)
	gvar397 = ns.Wall(75, 75)
	gvar398 = ns.Wall(35, 71)
	gvar399 = ns.Wall(36, 70)
	gvar400 = ns.Wall(37, 69)
	gvar401 = ns.Wall(30, 66)
	gvar402 = ns.Wall(29, 65)
	gvar403 = ns.Wall(28, 64)
	gvar404 = ns.Wall(27, 63)
	gvar405 = ns.Wall(28, 62)
	gvar406 = ns.Wall(29, 61)
	gvar407 = ns.Wall(30, 60)
	gvar408 = ns.Wall(31, 59)
	gvar409 = ns.Wall(32, 60)
	gvar410 = ns.Wall(33, 61)
	gvar411 = ns.Wall(34, 62)
	gvar412 = ns.Wall(35, 63)
	gvar413 = ns.Wall(34, 64)
	gvar414 = ns.Wall(33, 65)
	gvar415 = ns.Wall(32, 66)
	gvar416 = ns.Wall(31, 67)
	gvar417 = ns.Wall(63, 95)
	gvar418 = ns.Wall(64, 94)
	gvar419 = ns.Wall(65, 93)
	gvar420 = ns.Wall(66, 92)
	gvar421 = ns.Wall(67, 91)
	gvar422 = ns.WallGroup("SecretEndWalls")
	gvar423 = ns.WallGroup("DoorBackWalls1")
	gvar424 = ns.WallGroup("DoorBackWalls2")
	gvar425 = ns.WallGroup("Door2Walls")
	gvar426 = ns.WallGroup("Door3Walls")
	gvar427 = ns.WallGroup("Door4Walls")
	gvar431 = ns.Wall(38, 44)
	gvar432 = ns.Wall(39, 43)
	gvar433 = ns.Wall(40, 42)
	gvar434 = ns.Wall(41, 41)
	gvar435 = ns.Wall(42, 40)
	gvar436 = ns.Wall(43, 39)
	gvar437 = ns.Wall(44, 38)
	gvar438 = ns.Wall(45, 39)
	gvar439 = ns.Wall(46, 40)
	gvar440 = ns.Wall(47, 41)
	gvar441 = ns.Wall(48, 42)
	gvar442 = ns.Wall(49, 43)
	gvar443 = ns.Wall(50, 44)
	gvar444 = ns.Wall(49, 45)
	gvar445 = ns.Wall(48, 46)
	gvar446 = ns.Wall(47, 47)
	gvar447 = ns.Wall(46, 48)
	gvar448 = ns.Wall(45, 49)
	gvar449 = ns.Wall(44, 50)
	gvar450 = ns.Wall(43, 49)
	gvar451 = ns.Wall(42, 48)
	gvar452 = ns.Wall(41, 47)
	gvar453 = ns.Wall(40, 46)
	gvar454 = ns.Wall(39, 45)
	gvar455 = ns.Wall(38, 44)
	obj456 = ns.Object("NewGear")
	wp457 = ns.Waypoint("NewGearWP")
	ns.LockDoor(obj365)
	ns.LockDoor(obj366)
	ns.LockDoor(obj367)
	ns.LockDoor(obj368)
	ns.ObjectOff(obj369)
	ns.ObjectOff(obj370)
	ns.ObjectOff(obj371)
	ns.ObjectOff(obj372)
	ns.ObjectOff(obj373)
	ns.ObjectOff(obj374)
}
func SkullGunInit() {
	obj458 = ns.Object("SkullGun1")
	obj459 = ns.Object("SkullGun2")
	obj460 = ns.Object("SkullWall1Trigger")
	obj462 = ns.Object("Plate1")
	obj463 = ns.Object("Plate2")
	obj464 = ns.Object("Plate3")
	obj465 = ns.Object("Plate4")
	gvar466 = ns.ObjectGroup("SkullGuns")
	obj461 = ns.Object("StartTrapTrigger")
	wp467 = ns.Waypoint("SkullWP")
	gvar468 = ns.WallGroup("SkullWalls1")
	gvar469 = ns.WallGroup("SkullWalls2")
}
func TheForceInit() {
	wp475 = ns.Waypoint("FirstForceWP")
	fvar476 = ns.GetWaypointX(wp475)
	fvar477 = ns.GetWaypointY(wp475)
	wp478 = ns.Waypoint("FF1")
	wp479 = ns.Waypoint("FF2")
	wp480 = ns.Waypoint("FF3")
	wp481 = ns.Waypoint("FF4")
	wp482 = ns.Waypoint("FF5")
	wp483 = ns.Waypoint("FF6")
}
func MapInitialize() {
	gvar149 = ns.WallGroup("SouthRockWalls")
	gvar150 = ns.WallGroup("NorthRockWalls")
	gvar151 = ns.Wall(31, 137)
	gvar152 = ns.Wall(32, 136)
	gvar153 = ns.Wall(33, 135)
	gvar154 = ns.Wall(34, 134)
	gvar155 = ns.Wall(35, 133)
	gvar156 = ns.Wall(36, 132)
	gvar157 = ns.Wall(37, 131)
	gvar158 = ns.Wall(38, 130)
	gvar159 = ns.Wall(39, 129)
	gvar160 = ns.Wall(40, 128)
	gvar161 = ns.Wall(41, 127)
	gvar162 = ns.Wall(42, 126)
	gvar163 = ns.Wall(43, 125)
	gvar164 = ns.Wall(44, 124)
	gvar165 = ns.Wall(45, 123)
	gvar166 = ns.Wall(46, 122)
	gvar167 = ns.Wall(47, 121)
	gvar168 = ns.WallGroup("AnkhWalls")
	gvar169 = ns.Wall(138, 228)
	gvar170 = ns.Wall(139, 227)
	gvar171 = ns.Wall(140, 226)
	wp172 = ns.Waypoint("TickWP1")
	wp173 = ns.Waypoint("TickWP2")
	wp174 = ns.Waypoint("TickWP3")
	wp175 = ns.Waypoint("TickWP4")
	wp176 = ns.Waypoint("TickWP5")
	gvar179 = ns.WallGroup("ReturnWalls")
	obj180 = ns.Object("RockElevator")
	gvar181 = ns.WallGroup("GoodWalls")
	gvar182 = ns.Wall(201, 79)
	gvar183 = ns.Wall(202, 78)
	gvar184 = ns.Wall(203, 77)
	gvar185 = ns.Wall(204, 76)
	gvar186 = ns.Wall(205, 75)
	gvar187 = ns.Wall(206, 74)
	gvar188 = ns.Wall(207, 73)
	gvar189 = ns.Wall(208, 72)
	gvar190 = ns.Wall(209, 71)
	gvar191 = ns.Wall(207, 81)
	gvar192 = ns.Wall(208, 80)
	gvar193 = ns.Wall(209, 79)
	gvar194 = ns.Wall(210, 78)
	gvar195 = ns.Wall(211, 77)
	gvar196 = ns.WallGroup("BadWalls")
	gvar197 = ns.Wall(197, 75)
	gvar198 = ns.Wall(198, 74)
	gvar199 = ns.Wall(199, 73)
	gvar200 = ns.Wall(200, 72)
	gvar201 = ns.Wall(201, 71)
	gvar202 = ns.Wall(202, 70)
	gvar203 = ns.Wall(203, 69)
	gvar204 = ns.Wall(204, 68)
	gvar205 = ns.Wall(205, 67)
	gvar206 = ns.Wall(195, 69)
	gvar207 = ns.Wall(196, 68)
	gvar208 = ns.Wall(197, 67)
	gvar209 = ns.Wall(198, 66)
	gvar210 = ns.Wall(199, 65)
	wp212 = ns.Waypoint("WellWP")
	obj213 = ns.Object("LastBlock")
	obj214 = ns.Object("LastBlockMover")
	wp215 = ns.Waypoint("LastBlockEnd")
	obj216 = ns.Object("FrontBlockMover1")
	obj217 = ns.Object("FrontBlockMover2")
	obj218 = ns.Object("RearBlockMover1")
	obj219 = ns.Object("RearBlockMover2")
	obj220 = ns.Object("GeneratorMover1")
	obj221 = ns.Object("GeneratorMover2")
	obj222 = ns.Object("FireDragon")
	AltBlocksInit()
	APitsInit()
	BuzzSawInit()
	CrisCrosInit()
	GoodyPitInit()
	HallWorkInit()
	JailRoomInit()
	PitExitInit()
	PusherInit()
	RayTrapsInit()
	Secret1Init()
	SkullGunInit()
	TheForceInit()
}
func OpenPit1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj225)
}
func OpenPit2() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj226)
}
func OpenPit3() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj227)
}
func OpenPit4() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj228)
}
func OpenPit5() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj229)
}
func OpenPit6() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj230)
}
func OpenPit7() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj231)
}
func OpenPit8() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj232)
}
func OpenPit9() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj233)
}
func OpenPit10() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj234)
}
func OpenPit11() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj235)
}
func OpenPit12() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj236)
}
func OpenPit13() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj237)
}
func OpenPit14() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj238)
}
func OpenPit15() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj239)
}
func OpenPit16() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj240)
}
func OpenPit17() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj241)
}
func OpenPit18() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj242)
}
func OpenPit19() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj243)
}
func OpenPit20() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj244)
}
func OpenPit21() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj245)
}
func OpenPit22() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj246)
}
func OpenPit23() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj247)
}
func OpenPit24() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj248)
}
func OpenPit25() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj249)
}
func OpenPit26() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj250)
}
func OpenPit27() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj251)
}
func OpenPit28() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj252)
}
func OpenPit29() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj253)
}
func OpenPit30() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj254)
}
func OpenSurfaceWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar256)
	ns.WallOpen(gvar257)
	ns.WallOpen(gvar258)
	ns.WallOpen(gvar259)
	ns.WallOpen(gvar260)
}
func OpenGoodyWalls() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar262)
	ns.WallOpen(gvar263)
	ns.WallOpen(gvar264)
	ns.WallOpen(gvar265)
	ns.WallOpen(gvar266)
}
func OpenPitSouth() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar267)
	ns.WallOpen(gvar268)
	ns.WallOpen(gvar269)
	ns.WallOpen(gvar270)
	ns.WallOpen(gvar271)
	ns.WallOpen(gvar272)
	ns.WallOpen(gvar273)
}
func OpenPitNorth() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar274)
	ns.WallOpen(gvar275)
	ns.WallOpen(gvar276)
	ns.WallOpen(gvar277)
	ns.WallOpen(gvar278)
	ns.WallOpen(gvar279)
	ns.WallOpen(gvar280)
	ns.WallOpen(gvar281)
	ns.WallOpen(gvar282)
}
func StartPusher() {
	ns.ObjectOn(obj283)
	ns.ObjectOn(obj284)
	ns.ObjectOn(obj285)
	ns.ObjectOn(obj286)
	ns.ObjectOn(obj287)
	ns.ObjectOn(obj288)
	ns.ObjectOn(obj289)
	ns.ObjectOn(obj290)
}
func StartPushBlock12() {
	ns.ObjectOff(obj293)
	ns.Move(obj291, wp296)
	ns.Move(obj292, wp297)
	ns.SecondTimer(5, ReturnPushBlock12)
}
func ReturnPushBlock12() {
	ns.Move(obj291, wp294)
	ns.Move(obj292, wp295)
	ns.SecondTimer(3, ResetPushBlock12)
}
func ResetPushBlock12() {
	ns.ObjectOn(obj293)
}
func StartPushBlock34() {
	ns.ObjectOff(obj300)
	ns.Move(obj298, wp303)
	ns.Move(obj299, wp304)
	ns.SecondTimer(5, ReturnPushBlock34)
}
func ReturnPushBlock34() {
	ns.Move(obj298, wp301)
	ns.Move(obj299, wp302)
	ns.SecondTimer(3, ResetPushBlock34)
}
func ResetPushBlock34() {
	ns.ObjectOn(obj300)
}
func StartPushBlock56() {
	ns.ObjectOff(obj307)
	ns.Move(obj305, wp310)
	ns.Move(obj306, wp311)
	ns.SecondTimer(5, ReturnPushBlock56)
}
func ReturnPushBlock56() {
	ns.Move(obj305, wp308)
	ns.Move(obj306, wp309)
	ns.SecondTimer(3, ResetPushBlock56)
}
func ResetPushBlock56() {
	ns.ObjectOn(obj307)
}
func StartMovingRay() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj312)
	ns.ObjectOn(obj313)
	ns.WallOpen(gvar315)
	ns.WallOpen(gvar316)
	ns.WallOpen(gvar317)
}
func StartSpinSpike() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOn(obj318)
	ns.ObjectOn(obj319)
}
func OpenPusherExit1() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar321)
	ns.WallOpen(gvar322)
	ns.WallOpen(gvar323)
}
func OpenPusherExit2() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar325)
	ns.WallOpen(gvar326)
	ns.WallOpen(gvar327)
}
func OpenStuckExit() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar328)
	ns.WallOpen(gvar329)
	ns.WallOpen(gvar330)
}
func StopRayElevator() {
	ns.ObjectOff(obj363)
}
func StopSecondTime() {
	ns.ObjectOff(obj363)
	ns.SecondTimer(1, StartShooting)
}
func StartShooting() {
	ns.ObjectOn(obj332)
	ns.SecondTimer(2, OpenSetOne)
}
func MoveTheRay() {
	ns.MoveObject(obj332, ns.GetObjectX(obj363), ns.GetObjectY(obj363))
}
func StartSentryBlocks() {
	ns.ObjectGroupOff(gvar335)
	ns.AudioEvent(ns.SpikeBlockMover, wp340)
	ns.Move(obj332, wp339)
	ns.WallGroupClose(gvar360)
	StartRayTrapTimer()
}
func StartRayTrapTimer() {
	if !(ivar362 < 16) {
		goto LABEL1
	}
	ns.AudioEvent(ns.GenerateTick, wp340)
	ivar362 += 1
	gvar361 = ns.SecondTimer(1, StartRayTrapTimer)
	goto LABEL2
LABEL1:
	ns.CancelTimer(gvar361)
	ivar362 = 0
	ns.ObjectOn(obj332)
	ns.SecondTimer(2, OpenSetOne)
LABEL2:
	return
}
func OpenSetOne() {
	ns.NoWallSound(true)
	ns.WallOpen(gvar341)
	ns.WallOpen(gvar342)
	ns.WallOpen(gvar343)
	ns.NoWallSound(false)
	ns.FrameTimer(5, BreakSetOne)
}
func BreakSetOne() {
	ns.WallBreak(gvar344)
	ns.WallBreak(gvar345)
	ns.WallBreak(gvar346)
	ns.SecondTimer(2, OpenSetTwo)
}
func OpenSetTwo() {
	ns.NoWallSound(true)
	ns.WallOpen(gvar347)
	ns.WallOpen(gvar348)
	ns.WallOpen(gvar349)
	ns.NoWallSound(false)
	ns.FrameTimer(5, BreakSetTwo)
}
func BreakSetTwo() {
	ns.WallBreak(gvar350)
	ns.WallBreak(gvar351)
	ns.WallBreak(gvar352)
	ns.SecondTimer(2, OpenSetThree)
}
func OpenSetThree() {
	ns.NoWallSound(true)
	ns.WallOpen(gvar353)
	ns.WallOpen(gvar354)
	ns.WallOpen(gvar355)
	ns.NoWallSound(false)
	ns.FrameTimer(5, BreakSetThree)
}
func BreakSetThree() {
	ns.WallBreak(gvar356)
	ns.WallBreak(gvar357)
	ns.WallBreak(gvar358)
	ns.SecondTimer(3, StopTrap)
}
func OpenSentryWalls() {
	ns.WallGroupOpen(gvar359)
}
func StopTrap() {
	ns.ObjectOff(obj331)
	ns.ObjectOff(obj332)
	ns.ObjectOff(obj333)
	ns.ObjectOff(obj334)
	ns.Move(obj332, wp337)
	ns.SecondTimer(2, ResetTrap)
}
func ResetTrap() {
	ns.WallGroupOpen(gvar360)
}
func CloseSentryWalls() {
	ns.WallGroupClose(gvar359)
}
func CloseSecretStart() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallClose(gvar401)
	ns.WallClose(gvar402)
	ns.WallClose(gvar403)
	ns.WallClose(gvar404)
	ns.WallClose(gvar405)
	ns.WallClose(gvar406)
	ns.WallClose(gvar407)
	ns.WallClose(gvar408)
	ns.WallClose(gvar409)
	ns.WallClose(gvar410)
	ns.WallClose(gvar411)
	ns.WallClose(gvar412)
	ns.WallClose(gvar413)
	ns.WallClose(gvar414)
	ns.WallClose(gvar415)
	ns.WallClose(gvar416)
}
func OpenSecretEnd() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar431)
	ns.WallOpen(gvar432)
	ns.WallOpen(gvar433)
	ns.WallOpen(gvar434)
	ns.WallOpen(gvar435)
	ns.WallOpen(gvar436)
	ns.WallOpen(gvar437)
	ns.WallOpen(gvar438)
	ns.WallOpen(gvar439)
	ns.WallOpen(gvar440)
	ns.WallOpen(gvar441)
	ns.WallOpen(gvar442)
	ns.WallOpen(gvar443)
	ns.WallOpen(gvar444)
	ns.WallOpen(gvar445)
	ns.WallOpen(gvar446)
	ns.WallOpen(gvar447)
	ns.WallOpen(gvar448)
	ns.WallOpen(gvar449)
	ns.WallOpen(gvar450)
	ns.WallOpen(gvar451)
	ns.WallOpen(gvar452)
	ns.WallOpen(gvar453)
	ns.WallOpen(gvar454)
	ns.WallOpen(gvar455)
}
func NewFenceWallRoutine() {
	ns.ObjectOff(ns.GetTrigger())
	ns.AudioEvent(ns.CreatureCageAppears, wp457)
	ns.AudioEvent(ns.ChangeSpellbar, wp457)
	ns.ObjectOn(obj456)
	ns.WallOpen(gvar398)
	ns.WallOpen(gvar399)
	ns.WallOpen(gvar400)
}
func OpenDoor1() {
	if !(ns.IsLocked(obj365) && flag375 == false) {
		goto LABEL1
	}
	ns.UnlockDoor(obj365)
LABEL1:
	return
}
func OpenDoor2() {
	if !(ns.IsLocked(obj366) && flag376 == false) {
		goto LABEL1
	}
	ns.UnlockDoor(obj366)
LABEL1:
	return
}
func OpenDoor3() {
	if !(ns.IsLocked(obj367) && flag377 == false) {
		goto LABEL1
	}
	ns.UnlockDoor(obj367)
LABEL1:
	return
}
func OpenDoor4() {
	if !(ns.IsLocked(obj368) && flag378 == false) {
		goto LABEL1
	}
	ns.UnlockDoor(obj368)
LABEL1:
	return
}
func CloseDoor1() {
	if !(!ns.IsLocked(obj365) && flag375 == false) {
		goto LABEL1
	}
	ns.LockDoor(obj365)
LABEL1:
	return
}
func CloseDoor2() {
	if !(!ns.IsLocked(obj366) && flag376 == false) {
		goto LABEL1
	}
	ns.LockDoor(obj366)
LABEL1:
	return
}
func CloseDoor3() {
	if !(!ns.IsLocked(obj367) && flag377 == false) {
		goto LABEL1
	}
	ns.LockDoor(obj367)
LABEL1:
	return
}
func CloseDoor4() {
	if !(!ns.IsLocked(obj368) && flag378 == false) {
		goto LABEL1
	}
	ns.LockDoor(obj368)
LABEL1:
	return
}
func OpenDoor1Wall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar417)
	ns.WallOpen(gvar418)
	ns.WallOpen(gvar419)
	ns.WallOpen(gvar420)
	ns.WallOpen(gvar421)
}
func OpenDoor2Wall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar425)
}
func OpenDoor3Wall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar426)
}
func OpenDoor4Wall() {
	ns.ObjectOff(ns.GetTrigger())
	ns.WallGroupOpen(gvar427)
}
func OpenBackWall1() {
	ivar429 += 1
	if !(ivar429 >= 3) {
		goto LABEL1
	}
	ns.WallOpen(gvar390)
	ns.WallOpen(gvar391)
	ns.WallOpen(gvar392)
	ns.WallOpen(gvar393)
	ns.WallOpen(gvar394)
LABEL1:
	return
}
func OpenBackWall2() {
	ivar430 += 1
	if !(ivar430 >= 3) {
		goto LABEL1
	}
	ns.WallOpen(gvar395)
	ns.WallOpen(gvar396)
	ns.WallOpen(gvar397)
LABEL1:
	return
}
func OpenFence() {
	if ivar428 != 0 {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar379)
	ns.WallOpen(gvar380)
	ns.WallOpen(gvar381)
	ivar428 += 1
	goto LABEL2
LABEL1:
	ns.ObjectOff(ns.GetTrigger())
	ns.WallOpen(gvar382)
	ns.WallOpen(gvar383)
	ns.WallOpen(gvar384)
	ns.WallOpen(gvar385)
	ns.WallOpen(gvar386)
	ns.WallOpen(gvar387)
	ns.WallOpen(gvar388)
	ns.WallOpen(gvar389)
	ivar428 += 1
LABEL2:
	ns.Print("GeneralPrint:MsgWallLowered")
}
func StartLG1() {
	ns.ObjectOn(obj369)
}
func StartLG2() {
	ns.ObjectOn(obj370)
}
func StartLG3() {
	ns.ObjectOn(obj371)
}
func LeftGenGone() {
	ns.ObjectOn(obj370)
	OpenBackWall1()
}
func CenterGenGone() {
	ns.ObjectOn(obj369)
	ns.ObjectOn(obj371)
	OpenBackWall1()
}
func RightGenGone() {
	ns.ObjectOn(obj370)
	OpenBackWall1()
}
func StartRG1() {
	ns.ObjectOn(obj372)
}
func StartRG2() {
	ns.ObjectOn(obj373)
}
func StartRG3() {
	ns.ObjectOn(obj374)
}
func LeftRGGone() {
	ns.ObjectOn(obj373)
	OpenBackWall2()
}
func CenterRGGone() {
	ns.ObjectOn(obj372)
	ns.ObjectOn(obj374)
	OpenBackWall2()
}
func RightRGGone() {
	ns.ObjectOn(obj373)
	OpenBackWall2()
}
func Door1Chain() {
	if !ns.IsLocked(obj365) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.UnlockDoor(obj365)
	flag375 = true
LABEL1:
	return
}
func Door2Chain() {
	if !ns.IsLocked(obj366) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.UnlockDoor(obj366)
	flag376 = true
LABEL1:
	return
}
func Door3Chain() {
	if !ns.IsLocked(obj367) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.UnlockDoor(obj367)
	flag377 = true
LABEL1:
	return
}
func Door4Chain() {
	if !ns.IsLocked(obj368) {
		goto LABEL1
	}
	ns.ObjectOff(ns.GetTrigger())
	ns.UnlockDoor(obj368)
	flag378 = true
LABEL1:
	return
}
func ShutOffGuns() {
	ns.ObjectOff(ns.GetTrigger())
	ns.ObjectOff(obj458)
	ns.ObjectOff(obj459)
}
func OpenSkullWalls1() {
	ns.WallGroupOpen(gvar468)
}
func CloseSkullWalls1() {
	if flag474 {
		goto LABEL1
	}
	ns.WallGroupClose(gvar468)
LABEL1:
	return
}
func PuzzleSolved() {
	ns.CancelTimer(gvar472)
	ns.WallGroupOpen(gvar468)
	ns.WallGroupOpen(gvar469)
	ns.ObjectGroupOff(gvar466)
	ns.ObjectOn(obj461)
	ivar470 = 0
	flag473 = true
	flag474 = true
}
func StartSkullGuns() {
	ns.ObjectGroupOn(gvar466)
	ns.SecondTimer(4, StopSkullGuns)
}
func StopSkullGuns() {
	ns.ObjectGroupOff(gvar466)
}
func StartCountDown() {
	if flag474 {
		goto LABEL1
	}
	ns.WallGroupClose(gvar468)
	ns.WallGroupClose(gvar469)
	ns.ObjectOff(obj461)
	Countdown()
LABEL1:
	return
}
func Countdown() {
	if !(ivar470 < 26) {
		goto LABEL1
	}
	ns.AudioEvent(ns.GenerateTick, wp467)
	ivar470 += 1
	gvar472 = ns.SecondTimer(1, Countdown)
	goto LABEL2
LABEL1:
	StartSkullGuns()
LABEL2:
	return
}
func PlateHit() {
	ivar471 += 1
	if !(ivar471 >= 4) {
		goto LABEL1
	}
	PuzzleSolved()
LABEL1:
	return
}
func PlateRelease() {
	if !flag474 {
		goto LABEL1
	}
	flag474 = false
LABEL1:
	if flag473 {
		goto LABEL2
	}
	ivar471 -= 1
LABEL2:
	return
}
func FirstForceSeries() {
	ns.PrintToAll("I work")
	ns.PushObjectTo(ns.GetTrigger(), fvar476, fvar477)
}
func FFHit1() {
	ns.PushObject(ns.GetCaller(), 25, ns.GetObjectX(ns.GetCaller())-10, ns.GetObjectY(ns.GetCaller())+10)
}
func FFHit3() {
	ns.PushObject(ns.GetCaller(), 25, ns.GetObjectX(ns.GetCaller())-10, ns.GetObjectY(ns.GetCaller())-10)
}
func FFHit4() {
	ns.PushObject(ns.GetCaller(), 25, ns.GetObjectX(ns.GetCaller())+10, ns.GetObjectY(ns.GetCaller())+10)
}
func FFHit5() {
	ns.PushObject(ns.GetCaller(), 25, ns.GetObjectX(ns.GetCaller())-10, ns.GetObjectY(ns.GetCaller())+10)
}
func FFHit6() {
	ns.PushObject(ns.GetCaller(), 1000, ns.GetWaypointX(wp483), ns.GetWaypointY(wp483))
}
func OnEvent(typ string) {
	switch typ {
	case "MapInitialize":
		MapInitialize()
	}
}
