package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/server"
)

func unitDialogueCallback(index int32, caller, trigger *server.Object) {
	if err := GetServer().S().NoxScriptVM.CallByIndex(int(index), caller, trigger); err != nil {
		scriptLog.Println(err)
	}
}
func unitForceDialogue(player, monster *server.Object) {
	if player == nil || monster == nil || player.ObjClass&4 == 0 || monster.ObjClass&2 == 0 || player.ObjFlags&0x8020 != 0 || monster.Field5&0x10 == 0 {
		return
	}
	a, b := int32(*(*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](monster.UpdateData).DialogStartFunc))), int32(*(*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](monster.UpdateData).DialogEndFunc)))
	if a != -1 && b != -1 {
		unitDialogueCallback(a, player, monster)
	}
}
func unitFinishDialogue(player *server.Object, response byte) {
	stateUnfreeze(player, 0)
	partner := &layout.As[server.PlayerUpdateData](player.UpdateData).DialogWith
	monster := *partner
	if monster == nil {
		return
	}
	a, b := int32(*(*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](monster.UpdateData).DialogStartFunc))), int32(*(*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](monster.UpdateData).DialogEndFunc)))
	if a == -1 || b == -1 {
		return
	}
	reliableEnqueue(int(player.UpdateDataPlayer().Player.PlayerInd), []byte{0xd0, 4}, nil, 1, 0)
	*partner = nil
	if layout.As[server.MonsterUpdateData](monster.UpdateData).DialogFlags != 1 {
		response = 0
	}
	layout.As[server.MonsterUpdateData](monster.UpdateData).DialogResult = response
	unitDialogueCallback(b, player, monster)
}
