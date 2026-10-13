package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

func monsterAutoSpells(u *server.Object) int16 {
	core := GetServer().S()
	names := [9]string{"UrchinShaman", "Wizard", "WizardWhite", "Beholder", "Lich", "LichLord", "Demon", "WizardGreen", "WillOWisp"}
	if *memmap.PtrUint32(0x5D4594, 2491640) == 0 {
		for i, name := range names {
			*memmap.PtrUint32(0x5D4594, 2491624+uintptr(4*i)) = uint32(core.Types.IndByID(name))
		}
	}
	var ids [9]uint32
	for i := range ids {
		ids[i] = *memmap.PtrUint32(0x5D4594, 2491624+uintptr(4*i))
	}
	ud := (*server.MonsterUpdateData)(u.UpdateData)
	put := func(off uintptr, value uint32) { // off is a 386 offset
		*(*uint32)(unsafe.Add(u.UpdateData, layoutServerMonsterUpdateData.Off(off))) = value
	}
	ret := int16(u.TypeInd)
	switch uint32(u.TypeInd) {
	case ids[0]:
		put(2040, 2)
		put(1640, 0x8000000)
		ud.Field362_0 = 0
		ud.Field362_2 = uint16(core.TickRate() >> 1)
		put(1536, 0x20000000)
		ud.Field366_0 = 3 * uint16(core.TickRate())
		ud.Field366_2 = 5 * uint16(core.TickRate())
		ud.Field368_0 = 0
		put(1720, 0x40000000)
		put(1696, 0x40000000)
		put(1752, 0x80000000)
		ud.Field368_2 = 3 * uint16(core.TickRate())
		ud.Field370_0 = uint16(core.TickRate())
		ud.Field370_2 = 3 * uint16(core.TickRate())
		ret = int16(core.TickRate())
	case ids[1], ids[2]:
		put(2040, 3)
		put(1540, 0x8000000)
		put(1640, 0x8000000)
		put(1644, 0x10000000)
		put(1776, 0x20000000)
		for _, off := range []uintptr{1596, 1688, 1660, 1584} {
			put(off, 0x40000000)
		}
		put(1504, 0x80000000)
		ret = 0
	case ids[3]:
		put(2040, 3)
		put(1584, 0x40000000)
		ret = int16(bool2int(noxflags.HasGame(4096)))
		if ret == 0 {
			put(1504, 0x80000000)
		}
	case ids[4]:
		put(2040, 3)
		put(1540, 0x8000000)
		put(1640, 0x8000000)
		put(1772, 0x10000000)
		put(1776, 0x20000000)
		put(1620, 0x20000000)
		put(1644, 0x80000000)
		for _, off := range []uintptr{1596, 1660, 1584} {
			put(off, 0x40000000)
		}
		ret = 0
	case ids[5]:
		put(2040, 3)
		put(1540, 0x8000000)
		put(1640, 0x8000000)
		put(1772, 0x10000000)
		ud.Field368_0 = 3 * uint16(core.TickRate())
		put(1612, 0x40000000)
		put(1644, 0x80000000)
		ud.Field368_2 = 5 * uint16(core.TickRate())
		ret = int16(core.TickRate())
	case ids[6]:
		put(2040, 3)
		put(1540, 0x8000000)
		put(1784, 0x20000000)
		put(1596, 0x40000000)
		put(1528, 0x40000000)
		put(1644, 0x80000000)
		ret = 0
	case ids[7]:
		put(2040, 3)
		put(1640, 0x8000000)
		put(1536, 0x20000000)
		for _, off := range []uintptr{1720, 1728, 1688} {
			put(off, 0x40000000)
		}
		put(1504, 0x80000000)
		ret = 0
	case ids[8]:
		put(2040, 3)
		put(1660, 0x40000000)
	}
	*(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](u.UpdateData).Field509)) = 1
	return ret
}

// The C ABI passes each product as float32 to nox_float2int. Conversion to
// int32 on the supported 386/SSE2 target also preserves its indefinite result.
func creationScaleDurability(u *server.Object, d *server.Modifier) int32 {
	u.HealthData.Cur = uint16(d.Durability52)
	u.HealthData.Max = uint16(d.Durability52)
	if !noxflags.HasGame(4096) {
		return 0
	}
	scale := float64(float32(GetServer().S().Balance.Float("QuestDurabilityMultiplier")))
	u.HealthData.Cur = uint16(int32(float32(float64(u.HealthData.Cur) * scale)))
	result := int32(float32(float64(u.HealthData.Max) * scale))
	u.HealthData.Max = uint16(result)
	return result
}
func createWeapon(u *server.Object) int32 {
	core := GetServer().S()
	init := u.InitData
	d := core.Modif.Nox_xxx_getProjectileClassById413250(int(u.TypeInd))
	if *memmap.PtrUint32(0x5D4594, 2491660) == 0 {
		for i, name := range []string{"OblivionHeart", "OblivionWierdling", "OblivionOrb"} {
			*memmap.PtrUint32(0x5D4594, 2491660+uintptr(4*i)) = uint32(core.Types.IndByID(name))
		}
	}
	if d != nil && u.HealthData != nil {
		creationScaleDurability(u, d)
	}
	setMod := func(slot int, name string) {
		// Both the initialization buffer and modifier descriptors are C-owned.
		// Store address bits without a Go write barrier: the old slot may still
		// contain uninitialized bytes, which must not be scanned as a Go pointer.
		mods := &layout.As[server.ModifierInitData](init).Modifiers
		*(*uintptr)(unsafe.Pointer(&mods[slot])) = uintptr(unsafe.Pointer(core.Modif.Nox_xxx_modifGetDescById413330(core.Modif.Nox_xxx_modifGetIdByName413290(name))))
	}
	if uint32(u.TypeInd) == *memmap.PtrUint32(0x5D4594, 2491660) {
		setMod(2, "Lightning4")
	} else if uint32(u.TypeInd) == *memmap.PtrUint32(0x5D4594, 2491664) {
		setMod(2, "Vampirism2")
		setMod(3, "Lightning3")
	}
	if u.Class()&0x1000000 != 0 {
		if u.SubClass()&0x82 != 0 {
			data := u.UseDataAmmo()
			name := "DefaultAmmoAmount"
			if noxflags.HasGame(4096) {
				name = "DefaultAmmoAmountQuest"
			}
			n := byte(int32(float32(core.Balance.Float(name))))
			data.Charges = n
			data.Infinite = 0
			data.MaxCharges = n
		} else if u.SubClass()&0xc != 0 {
			u.UseDataBow().Reload = 0
		}
	}
	if !noxflags.HasGame(4096) {
		return 0
	}
	result := int32(u.Class())
	if u.Class()&0x1000 != 0 && u.SubClass()&0x47f0000 != 0 {
		data := u.UseDataWand()
		scale := float64(float32(core.Balance.Float("QuestStaffChargeMultiplier")))
		// Doubling stays in x87's double-precision register until the product spill.
		if u.SubClass()&0x40000 != 0 {
			scale += scale
		}
		data.MaxCharges = byte(int32(float32(float64(data.MaxCharges) * scale)))
		result = int32(float32(float64(data.Charges) * scale))
		data.Charges = byte(result)
	}
	return result
}

// This legacy return mixes a definition address and integer conversion bits.
// Keep it an integer across the Go stack, where small values are not pointers.
func createArmor(u *server.Object) uintptr {
	d := GetServer().S().Modif.Nox_xxx_equipClothFindDefByTT413270(int(u.TypeInd))
	if d == nil || u.HealthData == nil {
		return uintptr(unsafe.Pointer(d))
	}
	return uintptr(uint32(creationScaleDurability(u, d)))
}
func diePolyp(u *server.Object) {
	core := GetServer().S()
	id := monsterCache(2491672, "ToxicCloud")
	t := core.NewObjectByTypeInd(int(id))
	if t != nil {
		data := &t.UpdateDataToxicCloud().Life
		GetServer().CreateObjectAt(t, nil, u.PosVec)
		*data = uint32(int32(float32(core.Balance.Float("ToxicCloudLifetime") * float64(int32(core.TickRate())))))
	}
	core.Audio.EventObj(284, u, 0, 0)
	GetServer().DelayedDelete(u)
}

func createObelisk(u *server.Object) int32 {
	*(*uint32)(u.UpdateData) = 50
	u.NeedSync()
	return 0
}

func createAnim(u *server.Object) { u.SetXStatus(2) }

func createTrigger(u *server.Object) *byte {
	u.UpdateDataTrigger().Colors = [6]byte{90, 90, 90, 10, 10, 10}
	return (*byte)(u.UpdateData)
}

func createMonsterGenerator(u *server.Object) *uint32 {
	ud := u.UpdateDataMonsterGen()
	ud.Flags92 = 2
	ud.FuncInd52, ud.FuncInd60, ud.FuncInd76, ud.FuncInd68 = 0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff
	return (*uint32)(u.UpdateData)
}

func createRewardMarker(u *server.Object) *uint32 {
	data := u.InitDataRewardMarker()
	data.Categories = 255
	data.Chance = 0
	return &data.Categories
}
