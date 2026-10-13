package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/server"
)

func objectDeathImpEgg(u *server.Object) uint32 {
	GetServer().S().Audio.EventObj(764, u, 0, 0)
	u.ObjFlags |= 0x40
	return uint32(u.ObjFlags)
}

func objectDeathPotion(u *server.Object) {
	GetServer().S().Audio.EventObj(753, u, 0, 0)
	GetServer().DelayedDelete(u)
}

func objectDeathBarrel(u *server.Object) {
	core := GetServer().S()
	cache := memmap.PtrUint32(0x5D4594, 2491696)
	if *cache == 0 {
		*cache = uint32(core.Types.IndByID("BarrelBreaking"))
	}
	if t := core.NewObjectByTypeInd(int(*cache)); t != nil {
		GetServer().CreateObjectAt(t, nil, u.PosVec)
	}
	core.Audio.EventObj(286, u, 0, 0)
	stateLoot(u, &u.PosVec)
	GetServer().DelayedDelete(u)
}
func objectDeathCreate(u *server.Object, spawn bool) int16 {
	core := GetServer().S()
	data := u.DeathDataCreateObject()
	if t := core.NewObjectByTypeID(alloc.GoString(&data.Name[0])); t != nil {
		GetServer().CreateObjectAt(t, nil, u.PosVec)
	}
	if id := data.Sound; id != 0 {
		core.Audio.EventObj(sound.ID(id), u, 0, 0)
	}
	if spawn {
		u.ObjFlags |= 0x8000
		return int16(u.ObjFlags)
	}
	GetServer().DelayedDelete(u)
	return 0
}
func objectDeathMarker(u *server.Object) {
	if owner := u.FindOwnerChainPlayer(); owner != nil {
		// The retained owner-chain service can return a non-player terminal owner;
		// preserve the C callback's raw four-slot layout access (the teleport markers of players).
		marks := &layout.As[server.PlayerUpdateData](owner.UpdateData).Field29
		for i := range marks {
			p := &marks[i]
			if *p == u {
				*p = nil
				break
			}
		}
	}
	monsterPointFX(u, 138)
	GetServer().DelayedDelete(u)
}
func objectDeathBoulder(u *server.Object) {
	core := GetServer().S()
	core.Audio.EventObj(757, u, 0, 0)
	monsterPointFX(u, 138)
	n := core.Rand.Logic.IntClamp(20, 30)
	index := uint32(dword_5d4594_2491704)
	for i := 0; i < n; i++ {
		name := memmap.LoadPtr(0x587000, 291512+4*uintptr(index))
		t := core.NewObjectByTypeID(alloc.GoString((*byte)(name)))
		if t == nil {
			return
		} // Original early failure leaves the source alive.
		monsterDebrisPlace(u, t, 30)
		Nox_xxx_unitRaise_4E46F0(t, float32(core.Rand.Logic.FloatClamp(10, 70)))
		t.Field27 = float32(core.Rand.Logic.FloatClamp(-2, 0))
		t.ObjFlags |= 0x800000
		t.Field29 = math.Float32bits(float32(*(*byte)(memmap.PtrOff(0x587000, 290328+uintptr(dword_5d4594_2491704)))))
		GetServer().ApplyForce(t, u.PosVec, float64(float32(core.Rand.Logic.FloatClamp(5, 20))))
		monsterDebrisDecay(t, 45, 75)
		index = (uint32(dword_5d4594_2491704) + 1) % *memmap.PtrUint32(0x587000, 290340)
		dword_5d4594_2491704 = uint32(index)
	}
	GetServer().DelayedDelete(u)
}
func objectDeathArmor(u *server.Object) {
	core := GetServer().S()
	plural := false
	if lang := core.Strings().Lang(); lang == 0 || lang == 1 {
		if def := core.Modif.Nox_xxx_equipClothFindDefByTT413270(int(u.TypeInd)); def != nil {
			name := alloc.GoString16(def.Desc8)
			plural = len(name) != 0 && (name[len(name)-1] == 's' || name[len(name)-1] == 'S')
		}
	}
	holder := u.InvHolder
	pos := &u.PosVec
	if holder != nil {
		pos = &holder.PosVec
	}
	key := "ArmorDieGeneric"
	id := sound.ID(memmap.PtrWord(unsafe.Pointer(pos)))
	switch {
	case u.Material&16 != 0:
		key = "ArmorDieMetal"
		id = 806
	case u.Material&8 != 0:
		key = "ArmorDieWood"
		id = 812
	case u.Material&4 != 0:
		key = "ArmorDieHide"
		id = 809
	case u.Material&2 != 0:
		key = "ArmorDieCloth"
		id = 815
	}
	if plural && key != "ArmorDieGeneric" {
		key += "Plural"
	}
	format := alloc.InternCString16(core.Strings().GetStringInFile(strman.ID(key), "Die.c"))
	name := unitItemName(u)
	textFormatLine(holder, format, textFormatPointer(unsafe.Pointer(name)))
	core.Audio.EventPos(id, *pos, 0, 0)
	GetServer().DelayedDelete(u)
}
func objectDeathWeapon(u *server.Object) {
	core := GetServer().S()
	holder := u.InvHolder
	pos := &u.PosVec
	if holder != nil {
		pos = &holder.PosVec
	}
	key := "WeaponDieGeneric"
	var id sound.ID
	var name *uint16
	switch {
	case u.Material&16 != 0:
		name = unitItemName(u)
		key = "WeaponDieMetal"
		id = 818
	case u.Material&8 != 0:
		name = unitItemName(u)
		key = "WeaponDieWood"
		id = 819
	default:
		name = alloc.InternCString16(core.Armor.Sub_415B60(u))
	}
	format := alloc.InternCString16(core.Strings().GetStringInFile(strman.ID(key), "Die.c"))
	textFormatLine(holder, format, textFormatPointer(unsafe.Pointer(name)))
	if id != 0 {
		core.Audio.EventPos(id, *pos, 0, 0)
	}
	GetServer().DelayedDelete(u)
}
