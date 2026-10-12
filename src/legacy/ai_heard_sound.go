package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

// heardSoundAction is sub_545DA0. It turns a still-fresh heard-sound record
// into the C action sequence and clears the record only after those pushes.
func heardSoundAction(u *server.Object) int {
	ud := u.UpdateDataMonster()
	if ud.Field97 == 0 || GetServer().S().Frame()-ud.Field101 >= 3*GetServer().S().TickRate() {
		return 0
	}

	u.MonsterPushAction(ai.DEPENDENCY_NO_INTERESTING_SOUND)
	u.MonsterPushAction(ai.DEPENDENCY_NO_VISIBLE_ENEMY)
	if st := u.MonsterPushAction(ai.ACTION_WAIT); st != nil {
		st.SetArgs(GetServer().S().Frame() + uint32(nox_common_randomInt_415FA0(int(int32(GetServer().S().TickRate())), int(int32(2*GetServer().S().TickRate())))))
	}
	if st := u.MonsterPushAction(ai.ACTION_FACE_LOCATION); st != nil {
		st.SetArgs(types.Pointf{X: ud.Field99X, Y: ud.Field99Y})
	}
	ud.Field97 = 0
	return 1
}

// investigateHeardSound is sub_5466F0. A fresh record is consumed even when
// its tile, precheck, or ray path prevents it from scheduling actions.
func investigateHeardSound(u *server.Object) int {
	ud := u.UpdateDataMonster()
	if ud.Field97 == 0 || GetServer().S().Frame()-ud.Field101 >= 3*GetServer().S().TickRate() {
		return 0
	}

	heard := (*types.Pointf)(unsafe.Pointer(&ud.Field99X))
	if uint32(u.ObjSubClass)&0x400 != 0 || tileAtPoint(*heard) != 6 {
		if GetServer().Sub_50B810(u, heard) {
			if GetServer().S().MapTraceRayAt(u.PosVec, *heard, nil, nil, server.MapTraceFlag1|server.MapTraceFlag4) {
				heardSoundAction(u)
			} else {
				u.MonsterPushAction(ai.DEPENDENCY_NO_INTERESTING_SOUND)
				u.MonsterPushAction(ai.DEPENDENCY_NO_VISIBLE_ENEMY)
				u.MonsterPushAction(ai.DEPENDENCY_NOT_FRUSTRATED)
				if st := u.MonsterPushAction(ai.DEPENDENCY_LOCATION_IS_SAFE); st != nil {
					st.SetArgs(*heard)
				}
				if st := u.MonsterPushAction(ai.ACTION_MOVE_TO); st != nil {
					st.SetArgs(*heard, uint32(0))
				}
			}
		}
	}
	ud.Field97 = 0
	return 1
}
