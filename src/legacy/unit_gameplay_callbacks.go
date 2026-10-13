package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func unitUndeadUpdate(u *server.Object) {
	// The collide data holds the turn undead spell as a pointer word (see sustainedTurnUndeadStart).
	spell := (*server.DurSpell)(u.CollideDataUndeadKiller().SpellPtr())
	if spell != nil && spell.Flags88&1 != 0 || GetServer().S().Frame()-u.Field34 > 70 {
		GetServer().DelayedDelete(u)
	}
}
func unitRead(u, it *server.Object, warp bool) bool {
	if u.ObjClass&4 == 0 {
		return true
	}
	data := it.UseDataRead()
	stamp := &data.LastRead
	core := GetServer().S()
	if *stamp != 0 && core.Frame()-*stamp <= 3*uint32(core.TickRate()) {
		return true
	}
	if !core.MapTraceVision(u, it) {
		return true
	}
	if !warp {
		gameplayTextPrivate(u, &data.TextBuf[0], 1)
	} else if questRuntimeWord(1556120) != 0 {
		threshold := uint32(Nox_server_questNextStageThreshold_4D74F0(int(int32(questRuntimeStage()))))
		gameplayTextInformation(int(u.UpdateDataPlayer().Player.PlayerInd), 21, unsafe.Pointer(&threshold))
	} else {
		gameplayTextPrivate(u, alloc.InternCString("GeneralPrint:WarpClosed"), 1)
	}
	*stamp = core.Frame()
	return true
}
