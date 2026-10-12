package legacy

import (
	"runtime"
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

// These non-zero-size global bytes give each callback table entry one stable,
// distinct process-lifetime identity that can be stored in the foreign blob.
var monsterCallbackIdentityBytes [26]byte

const (
	monsterCallbackID_0  = 0  // nox_xxx_strikeOgre_549220
	monsterCallbackID_1  = 1  // nox_xxx_strikeScorpion_5495B0
	monsterCallbackID_2  = 2  // nox_xxx_strikeVileZombie_549700
	monsterCallbackID_3  = 3  // nox_xxx_strikeStoneGolem_5497E0
	monsterCallbackID_4  = 4  // nox_xxx_strikeMechGolem_549960
	monsterCallbackID_5  = 5  // nox_xxx_strikeWasp_549980
	monsterCallbackID_6  = 6  // nox_xxx_strikeSpider_549BC0
	monsterCallbackID_7  = 7  // nox_xxx_strikeSpittingSpider_549CA0
	monsterCallbackID_8  = 8  // nox_xxx_strikeGhost_549A60
	monsterCallbackID_9  = 9  // nox_xxx_strikeBomber_549BB0
	monsterCallbackID_10 = 10 // nox_xxx_strikeMonsterDefault_549380
	monsterCallbackID_11 = 11 // sub_549D80
	monsterCallbackID_12 = 12 // sub_549E00
	monsterCallbackID_13 = 13 // sub_549E70
	monsterCallbackID_14 = 14 // sub_549E90
	monsterCallbackID_15 = 15 // sub_549FA0
	monsterCallbackID_16 = 16 // nox_bomberDead_54A150
	monsterCallbackID_17 = 17 // sub_54A250
	monsterCallbackID_18 = 18 // sub_54A310
	monsterCallbackID_19 = 19 // sub_54A750
	monsterCallbackID_20 = 20 // nox_xxx_monsterDeadTroll_54A270
	monsterCallbackID_21 = 21 // sub_54A890
	monsterCallbackID_22 = 22 // sub_54A900
	monsterCallbackID_23 = 23 // sub_54A7D0
	monsterCallbackID_24 = 24 // sub_54A850
	monsterCallbackID_25 = 25 // sub_54A950
)

func monsterCallbackKey(id int) unsafe.Pointer {
	return unsafe.Pointer(&monsterCallbackIdentityBytes[id])
}

var monsterCallbackHandlers = make(map[unsafe.Pointer]func(*server.Object) int32)

func init() {
	for key, fn := range map[unsafe.Pointer]func(*server.Object) int32{
		monsterCallbackKey(monsterCallbackID_0):  func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 0))) },  // nox_xxx_strikeOgre_549220
		monsterCallbackKey(monsterCallbackID_1):  func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 1))) },  // nox_xxx_strikeScorpion_5495B0
		monsterCallbackKey(monsterCallbackID_2):  func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 2))) },  // nox_xxx_strikeVileZombie_549700
		monsterCallbackKey(monsterCallbackID_3):  func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 3))) },  // nox_xxx_strikeStoneGolem_5497E0
		monsterCallbackKey(monsterCallbackID_4):  func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 4))) },  // nox_xxx_strikeMechGolem_549960
		monsterCallbackKey(monsterCallbackID_5):  func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 5))) },  // nox_xxx_strikeWasp_549980
		monsterCallbackKey(monsterCallbackID_6):  func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 6))) },  // nox_xxx_strikeSpider_549BC0
		monsterCallbackKey(monsterCallbackID_7):  func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 7))) },  // nox_xxx_strikeSpittingSpider_549CA0
		monsterCallbackKey(monsterCallbackID_8):  func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 8))) },  // nox_xxx_strikeGhost_549A60
		monsterCallbackKey(monsterCallbackID_9):  func(*server.Object) int32 { return 1 },                                       // nox_xxx_strikeBomber_549BB0
		monsterCallbackKey(monsterCallbackID_10): func(u *server.Object) int32 { return int32(bool2int(monsterStrike(u, 10))) }, // nox_xxx_strikeMonsterDefault_549380
		monsterCallbackKey(monsterCallbackID_11): func(u *server.Object) int32 { monsterDeathExplosion(u, false); return 1 },    // sub_549D80
		monsterCallbackKey(monsterCallbackID_12): func(u *server.Object) int32 { monsterDeathExplosion(u, true); return 1 },     // sub_549E00
		monsterCallbackKey(monsterCallbackID_13): func(u *server.Object) int32 { monsterPointFX(u, 129); return 1 },             // sub_549E70
		monsterCallbackKey(monsterCallbackID_14): func(u *server.Object) int32 { monsterDeathDebris(u); return 1 },              // sub_549E90
		monsterCallbackKey(monsterCallbackID_15): func(u *server.Object) int32 { monsterDeathChunks(u); return 1 },              // sub_549FA0
		monsterCallbackKey(monsterCallbackID_16): func(u *server.Object) int32 { return int32(Nox_bomberDead_54A150(u)) },       // nox_bomberDead_54A150
		monsterCallbackKey(monsterCallbackID_17): func(u *server.Object) int32 { monsterPointFX(u, 129); return 1 },             // sub_54A250
		monsterCallbackKey(monsterCallbackID_18): func(u *server.Object) int32 { monsterDeathLoot(u, 17); return 1 },            // sub_54A310
		monsterCallbackKey(monsterCallbackID_19): func(u *server.Object) int32 { monsterDeathLoot(u, 18); return 1 },            // sub_54A750
		monsterCallbackKey(monsterCallbackID_20): func(u *server.Object) int32 { monsterDeathTroll(u); return 1 },               // nox_xxx_monsterDeadTroll_54A270
		monsterCallbackKey(monsterCallbackID_21): func(u *server.Object) int32 { monsterDeathLoot(u, 20); return 1 },            // sub_54A890
		monsterCallbackKey(monsterCallbackID_22): func(u *server.Object) int32 { monsterDeathLoot(u, 21); return 1 },            // sub_54A900
		monsterCallbackKey(monsterCallbackID_23): func(u *server.Object) int32 { monsterDeathLoot(u, 22); return 1 },            // sub_54A7D0
		monsterCallbackKey(monsterCallbackID_24): func(u *server.Object) int32 { monsterDeathLoot(u, 23); return 1 },            // sub_54A850
		monsterCallbackKey(monsterCallbackID_25): func(u *server.Object) int32 { monsterDeathLoot(u, 24); return 1 },            // sub_54A950
	} {
		monsterCallbackHandlers[key] = fn
	}
}

func monsterCallbackResult(key unsafe.Pointer, u *server.Object) int32 {
	if fn := monsterCallbackHandlers[key]; fn != nil {
		result := fn(u)
		runtime.KeepAlive(u)
		return result
	}
	panic("unregistered monster callback")
}

func monsterCallbackDiscard(key unsafe.Pointer, u *server.Object) {
	if fn := monsterCallbackHandlers[key]; fn != nil {
		fn(u)
		runtime.KeepAlive(u)
		return
	}
	panic("unregistered monster callback")
}
