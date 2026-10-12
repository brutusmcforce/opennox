package legacy

import (
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
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
	a, b := int32(*equipmentWord(monster.UpdateData, 2096)), int32(*equipmentWord(monster.UpdateData, 2100))
	if a != -1 && b != -1 {
		unitDialogueCallback(a, player, monster)
	}
}
func unitFinishDialogue(player *server.Object, response byte) {
	stateUnfreeze(player, 0)
	partner := (**server.Object)(unsafe.Add(player.UpdateData, 284))
	monster := *partner
	if monster == nil {
		return
	}
	a, b := int32(*equipmentWord(monster.UpdateData, 2096)), int32(*equipmentWord(monster.UpdateData, 2100))
	if a == -1 || b == -1 {
		return
	}
	reliableEnqueue(int(player.UpdateDataPlayer().Player.PlayerInd), []byte{0xd0, 4}, nil, 1, 0)
	*partner = nil
	if *controlByte(monster.UpdateData, 2104) != 1 {
		response = 0
	}
	*controlByte(monster.UpdateData, 2105) = response
	unitDialogueCallback(b, player, monster)
}
