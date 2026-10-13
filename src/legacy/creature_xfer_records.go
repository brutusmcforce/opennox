package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/server"
)

func creatureXferMerchant(r objectXferStream, u *server.Object, v int) {
	if v < 43 || u.ObjSubClass&8 == 0 {
		return
	}
	p := u.InitData
	if v >= 50 {
		r.raw(unsafe.Pointer(&layout.As[server.ShopkeeperInitData](p).Field1716), 4)
	}
	if v >= 61 {
		r.raw(unsafe.Pointer(&layout.As[server.ShopkeeperInitData](p).Field1720), 4)
	}
	if v >= 48 {
		creatureXferText(r, unsafe.Pointer(&layout.As[server.ShopkeeperInitData](p).ShopText))
	}
	count := byte(0)
	if p != nil {
		count = *(*byte)(p)
	}
	count = r.byte(count)
	if p == nil {
		return
	}
	if r.read() {
		*(*byte)(p) = count
	}
	stock := layout.As[server.ShopkeeperInitData](p)
	for i := 0; i < int(stock.StockCount); i++ {
		entry := &stock.Stock[i]
		if r.read() {
			Nox_xxx_XFer_ReadShopItem_52A840(entry, v)
		} else {
			Nox_xxx_XFer_WriteShopItem_52A5F0(entry)
		}
	}
}
func creatureXferColors(r objectXferStream, u *server.Object) {
	for i := 0; i < 6; i++ {
		p := unsafe.Add(u.UpdateData, layoutServerMonsterUpdateData.Off(uintptr(2076+3*i)))
		if r.read() {
			var color server.Color3
			r.raw(unsafe.Pointer(&color), 3)
			u.Nox_xxx_setNPCColor_4E4A90(byte(i), &color)
		} else {
			r.raw(p, 3)
		}
	}
}
func creatureXferMonster(u *server.Object) int {
	cache := memmap.PtrUint32(0x5d4594, 2487692)
	if *cache == 0 {
		*cache = uint32(GetServer().S().Types.IndByID("Glyph"))
	}
	r, v, saved, ok := objectXferStart(u, 64)
	if !ok {
		return 0
	}
	p := u.UpdateData
	at := func(off, n int) { creatureXferRaw(r, layoutServerMonsterUpdateData, p, off, n) } // 386 offsets, as saved
	creatureXferHeader(r, u, v, false, u.EditorData())
	creatureXferStats(r, u, v, false)
	if v >= 41 && creatureXferAction(u) == 0 {
		return 0
	}
	if v >= 42 {
		at(1445, 1)
	}
	creatureXferMerchant(r, u, v)
	if v >= 44 {
		at(0, 4)
	}
	if v >= 45 {
		u.ObjSubClass |= object.SubClass(r.word(uint32(u.ObjSubClass) & 0x180))
	}
	if v >= 49 {
		r.raw(unsafe.Pointer(u.HealthData), 2)
	}
	if v >= 51 {
		at(1348, 1)
		at(1340, 1)
		at(1444, 1)
		at(2036, 1)
	}
	if u.ObjSubClass&0x20 != 0 && v >= 54 {
		creatureXferColors(r, u)
	}
	if u.ObjSubClass&0x20 != 0 && v >= 55 {
		creatureXferVoice(u)
	}
	if v >= 62 && creatureXferBuffs(u) == 0 {
		return 0
	}
	if v >= 63 && u.ObjSubClass&0x80000 != 0 {
		creatureXferVoice(u)
	}
	if v >= 64 {
		poison := r.byte(u.Poison540)
		if r.read() && poison != 0 {
			resourceSetPoison(u, int32(poison))
		}
	}
	if r.read() {
		if *(*byte)(unsafe.Add(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).Field361), 1)) != 0 && noxflags.HasGame(1) {
			u.HealthData.Max = 0
			u.HealthData.Cur = 0
		}
		if u.ObjFlags&0x8000 != 0 && !monsterIsZombie(u) {
			u.ObjFlags |= 0x40
		}
	}
	if u.Field34 != 0 {
		if !r.read() {
			u.Field34 = saved
			return 1
		}
		if objectXferInventory(uint16(v), u, int32(u.Field34)) == 0 {
			return 0
		}
	}
	if r.read() {
		if noxflags.HasGame(0x200000) || nox_xxx_gameIsSwitchToSolo_4DB240() == 0 {
			creatureXferDefaults(u)
		}
		if u.ObjClass&2 != 0 && u.ObjSubClass&0x2000 != 0 {
			found := false
			for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
				if uint32(it.TypeInd) == *cache {
					found = true
				}
			}
			spells := unsafe.Slice(&layout.As[server.MonsterUpdateData](p).Field511, 3)
			count := 0
			for _, id := range spells {
				if id != 0 {
					count++
				}
			}
			if !found && count != 0 {
				glyph := GetServer().S().NewObjectByTypeID("Glyph")
				if glyph != nil {
					data := glyph.InitDataGlyph()
					// Match the existing prefix copy even if the three source slots have gaps.
					copy(data.Spells[:count], spells[:count])
					*(*byte)(unsafe.Pointer(&data.SpellsCnt)) = byte(count)
					data.SpellArg.Obj = nil
					data.SpellArg.Pos = u.PosVec
				}
				inventoryInsert(u, glyph, 1)
			}
		}
	}
	u.Field34 = saved
	return 1
}

func creatureXferNPC(u *server.Object) int {
	r := objectXferStream{cryptfile.Global()}
	p, scripts, saved := u.UpdateData, u.EditorData(), u.Field34
	v := int(int16(r.short(62)))
	if v > 62 {
		return 0
	}
	if r.read() {
		def := Nox_xxx_monsterDefByTT_517560(int(u.TypeInd))
		layout.As[server.MonsterUpdateData](p).MonsterDef = def
		if def != nil {
			*(*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).StatusFlags)) = uint32(def.StatusFlags92)
		}
	}
	if objectXferCommon(u, v) == 0 {
		return 0
	}
	at := func(off, n int) { creatureXferRaw(r, layoutServerMonsterUpdateData, p, off, n) } // 386 offsets, as saved
	creatureXferHeader(r, u, v, true, scripts)
	creatureXferColors(r, u)
	if r.read() && v == 31 {
		r.short(0)
		count := r.byte(0)
		for i := 0; i < int(count); i++ {
			var name [256]byte
			n := r.byte(0)
			r.raw(unsafe.Pointer(&name[0]), int(n))
			r.byte(0)
			r.byte(0)
		}
	}
	creatureXferStats(r, u, v, true)
	if v >= 41 && creatureXferAction(u) == 0 {
		return 0
	}
	if v >= 42 {
		at(1445, 1)
	}
	if v >= 44 {
		at(0, 4)
	}
	if v >= 45 {
		health := uint32(0)
		if u.HealthData != nil {
			health = uint32(u.HealthData.Max)
		}
		health = r.word(health)
		if u.HealthData != nil {
			u.HealthData.Max = uint16(health)
		}
	}
	if v >= 46 {
		u.ObjSubClass |= object.SubClass(r.word(uint32(u.ObjSubClass) & 0x180))
	}
	if v >= 48 {
		r.raw(unsafe.Pointer(u.HealthData), 2)
	}
	if v >= 51 {
		r.raw(unsafe.Pointer(&u.Experience), 4)
	}
	if v >= 52 {
		creatureXferVoice(u)
	}
	if v >= 61 && creatureXferBuffs(u) == 0 {
		return 0
	}
	if v >= 62 {
		poison := r.byte(u.Poison540)
		if r.read() && poison != 0 {
			resourceSetPoison(u, int32(poison))
		}
	}
	if r.read() {
		if *(*byte)(unsafe.Add(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).Field361), 1)) != 0 && noxflags.HasGame(1) && u.HealthData != nil {
			u.HealthData.Max = 0
			u.HealthData.Cur = 0
		}
		if u.ObjFlags&0x8000 != 0 {
			u.ObjFlags |= 0x40
		}
	}
	if u.Field34 != 0 && r.read() && objectXferInventory(uint16(v), u, int32(u.Field34)) == 0 {
		return 0
	}
	creatureXferEquipment(u)
	if r.read() && noxflags.HasGame(1) {
		for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
			if it.ObjFlags&0x100 == 0 {
				continue
			}
			it.ObjFlags &^= 0x100
			if it.ObjClass&0x1001000 != 0 {
				equipmentNPCEquipWeapon(u, it)
			} else {
				equipmentNPCEquipArmor(u, it)
			}
		}
	}
	u.Field34 = saved
	return 1
}
