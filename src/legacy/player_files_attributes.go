package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/server"
)

func playerFileAttributes(u *server.Object, info unsafe.Pointer) int {
	if info == nil {
		return 0
	}
	var data unsafe.Pointer
	var p *server.Player
	if u != nil {
		data = u.UpdateData
		if data != nil {
			p = u.UpdateDataPlayer().Player
		}
	}
	r := playerFileStream()
	version := int16(r.short(5))
	if version > 5 {
		return 0
	}
	if version >= 5 {
		mode := uint32(2)
		if noxflags.HasGame(2048) {
			mode = 1
		}
		if noxflags.HasGame(4096) {
			mode = 4
		}
		mode = r.word(mode)
		if p != nil && p.PlayerInd != 31 && ((mode == 4) != noxflags.HasGame(4096)) {
			Nox_xxx_playerCallDisconnect_4DEAB0(ntype.PlayerInd(p.PlayerInd), 1)
			return 0
		}
	}
	size := r.byte(byte(playerFileWideLen(info)))
	if size >= 25 {
		return 0
	}
	r.raw(info, 2*int(size))
	*(*uint16)(unsafe.Add(info, 2*int(size))) = 0
	if data != nil {
		bytes := 2 * (playerFileWideLen(info) + 1)
		copy(unsafe.Slice((*byte)(unsafe.Pointer(&p.NameFinal[0])), bytes), unsafe.Slice((*byte)(info), bytes))
	}
	if r.read() && data != nil {
		crc := protectionStringChecksum(info, uint32(2*playerFileWideLen(info)))
		setProtectionRecord(int32(p.ProtPlayerOrigName), uint32(crc))
	}
	r.raw(unsafe.Add(info, 50), 4)
	r.raw(unsafe.Add(info, 54), 4)
	if r.read() && data != nil {
		setProtectionRecord(int32(p.ProtPlayerField2239), *(*uint32)(unsafe.Add(info, 54)))
		setProtectionRecord(int32(p.ProtPlayerField2235), *(*uint32)(unsafe.Add(info, 50)))
	}
	r.raw(unsafe.Add(info, 58), 4)
	r.raw(unsafe.Add(info, 62), 4)
	r.raw(unsafe.Add(info, 66), 1)
	if r.read() && data != nil {
		setProtectionRecord(int32(p.ProtPlayerClass), uint32(*(*byte)(unsafe.Add(info, 66))))
	}
	r.raw(unsafe.Add(info, 67), 1)
	for off := 68; off <= 80; off += 3 {
		r.raw(unsafe.Add(info, off), 3)
	}
	if version >= 2 {
		for off := 83; off <= 87; off++ {
			r.raw(unsafe.Add(info, off), 1)
		}
	}
	if r.read() && u != nil {
		if pl := GetServer().S().Players.ByID(int(u.NetCode)); pl != nil {
			clientPlayerColors(pl)
		}
	}
	r.raw(unsafe.Add(info, 88), 1)
	if version >= 3 {
		lives := uint32(0)
		if data != nil {
			lives = layout.As[server.PlayerUpdateData](data).Field80
		}
		lives = r.word(lives)
		if r.read() && data != nil {
			layout.As[server.PlayerUpdateData](data).Field80 = lives
		}
		max := uint32(floatToInt32(float32(GetServer().S().Balance.Float("MaxExtraLives"))))
		if data != nil && layout.As[server.PlayerUpdateData](data).Field80 > max {
			return 0
		}
		if version == 3 {
			for i := 0; i < 9; i++ {
				r.word(0)
			}
		}
	}
	if r.read() {
		questRuntimeReset(u)
	}
	if version >= 4 {
		stage := uint32(0)
		if p != nil {
			stage = p.Field4696
		}
		stage = r.word(stage)
		if r.read() {
			if p != nil {
				p.Field4696 = stage
			}
			if data != nil {
				questRuntimeHighestMessage(int(p.PlayerInd), uint16(p.Field4696))
			}
		}
	}
	return 1
}
func playerFileStatus(u *server.Object) int {
	data := u.UpdateData
	if GetServer().S().Players.ByID(int(u.NetCode)) == nil {
		return 0
	}
	r := playerFileStream()
	version := int16(r.short(2))
	if version > 2 {
		return 0
	}
	if !playerFilePresent(r, !noxflags.HasGame(2048)) {
		return 1
	}
	if !noxflags.HasGame(2048) {
		return 0
	}
	maxHP := r.short(uint16(resourceGetMaxHP(u)))
	if r.read() {
		resourceSetMaxHP(u, maxHP)
		resourceSetHP(u, maxHP)
	}
	maxMana := r.short(uint16(resourceGetMaxMana(u)))
	if r.read() {
		resourceSetMaxMana(u, maxMana)
		resourceRefreshMana(u)
	}
	hp := memmap.PtrUint32(0x5D4594, 527696)
	*hp = uint32(u.HealthData.Cur)
	r.raw(unsafe.Pointer(hp), 2)
	mana := memmap.PtrUint32(0x5D4594, 527700)
	*mana = uint32(layout.As[server.PlayerUpdateData](data).ManaCur)
	r.raw(unsafe.Pointer(mana), 2)
	poison := r.byte(u.Poison540)
	if r.read() {
		resourceSetPoison(u, int32(poison))
	}
	r.raw(unsafe.Pointer(&u.Field541), 1)
	r.raw(unsafe.Pointer(&u.Field542), 2)
	r.raw(unsafe.Pointer(&u.Experience), 4)
	if r.read() {
		p := u.UpdateDataPlayer().Player
		updateProtectionFloat(int32(p.ProtUnitExperience), u.Experience, false)
		gameplayReportExperience(u)
	}
	if version >= 2 {
		r.raw(unsafe.Pointer(&u.Direction1), 2)
		if r.read() {
			*(*uint16)(unsafe.Pointer(&u.Direction2)) = *(*uint16)(unsafe.Pointer(&u.Direction1))
		}
	}
	return 1
}
