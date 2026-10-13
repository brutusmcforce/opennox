package legacy

import (
	"strings"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func unitActionName(id int32, restricted bool) *byte {
	if restricted {
		allowed := false
		for _, v := range unsafe.Slice(memmap.PtrInt32(0x587000, 262056), 4) {
			if id == v && id >= 0 && id < 39 {
				allowed = true
				break
			}
		}
		if !allowed {
			id = 38
		}
	}
	if id < 0 || id >= 72 {
		return nil
	}
	return (*byte)(memmap.LoadPtr(0x587000, 261768+4*uintptr(id)))
}
func unitActionIndex(name string, restricted bool) int32 {
	name, _, _ = strings.Cut(name, "\x00")
	limit, fallback := int32(72), int32(0)
	if restricted {
		limit, fallback = 39, 38
	}
	for i := int32(0); i < limit; i++ {
		if alloc.GoString(unitActionName(i, false)) == name {
			return i
		}
	}
	return fallback
}
func unitLocalOrder(player int, order int32) int {
	pl := GetServer().S().Players.ByInd(ntype.PlayerInd(player))
	pl.SummonOrderAll = uint32(order)
	return gameplayReportCreature(player, byte(order))
}
func unitHurtSound(u *server.Object) {
	if u.ObjClass&2 == 0 {
		return
	}
	core := GetServer().S()
	next := &layout.As[server.MonsterUpdateData](u.UpdateData).Field133
	if core.Frame() >= *next {
		fps := uint32(core.TickRate())
		*next = core.Frame() + uint32(core.Rand.Logic.IntClamp(int(int32(2*fps)), int(int32(4*fps))))
		if set := resourceMonsterSound(u); set != nil {
			core.Audio.EventObj(sound.ID(set.Sound(2)), u, 0, 0)
		}
	}
}
func unitDamageTimer(u *server.Object) {
	stamp := &layout.As[server.MonsterUpdateData](u.UpdateData).Field130
	if *stamp == 0 {
		*stamp = GetServer().S().Frame()
	}
}
func unitDebug(kind int, frame, code uint32, name string) {
	if !noxflags.HasEngine(noxflags.EngineShowAI) {
		return
	}
	name, _, _ = strings.Cut(name, "\x00")
	if kind == 0 {
		ai.Log.Printf("%d: Lost sight of %s(#%d)\n", int32(frame), name, int32(code))
	} else {
		ai.Log.Printf("%d: %s(#%d) FRUSTRATED\n", int32(frame), name, int32(code))
	}
}
