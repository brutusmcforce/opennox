package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func creatureXferAdjust(p *uint32, delta uint32) uint32 {
	sum := *p + delta
	*p = sum
	if int32(sum) < 1 {
		*p = 1
	}
	return sum
}

func creatureXferDefaults(u *server.Object) {
	if u == nil {
		return
	}
	p := u.UpdateData
	def := Nox_xxx_monsterDefByTT_517560(int(u.TypeInd))
	if def == nil {
		return
	}
	if *(*byte)(unsafe.Add(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).Field361), 1)) == 0 {
		health := uint16(def.Health68)
		u.HealthData.Cur = health
		u.HealthData.Max = health
		u.HealthData.Field2 = health
	}
	status := (*uint32)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).StatusFlags))
	if *(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).Field361)) == 1 {
		*status = uint32(def.StatusFlags92)
	} else {
		const replace = uint32(0x3fffff &^ 0x19c40)
		*status = (*status &^ replace) | (uint32(def.StatusFlags92) & replace)
	}
	if *status&0x20 == 0 {
		*status &^= 0x1800
	}
	if *(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).Field335)) == 1 {
		layout.As[server.MonsterUpdateData](p).RetreatLevel = def.RetreatRatio80
	}
	if *(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).Field337)) == 1 {
		layout.As[server.MonsterUpdateData](p).ResumeLevel = def.ResumeRatio84
	}
	if *(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](p).Field509)) == 1 {
		monsterAutoSpells(u)
	}
}

func creatureXferEquipment(u *server.Object) int {
	if u == nil {
		return 0
	}
	weapon, shield := false, false
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjFlags&0x100 == 0 {
			continue
		}
		switch {
		case it.ObjClass&0x1001000 != 0 && it.ObjSubClass&0x7ffe40c != 0:
			if shield {
				it.ObjFlags &^= 0x100
			} else {
				weapon = true
			}
		case it.ObjClass&0x2000000 != 0 && it.ObjSubClass&2 != 0:
			if weapon {
				it.ObjFlags &^= 0x100
			} else {
				shield = true
			}
		}
	}
	return 0
}

func creatureXferPostload(u *server.Object) uint32 {
	p := u.UpdateData
	md := layout.As[server.MonsterUpdateData](p)
	for i := 0; i <= int(md.AIStackInd); i++ {
		st := &md.AIStack[i]
		id := st.Action
		count := int32(*memmap.PtrUint32(0x587000, 255604+uintptr(16*id)))
		for j := int32(0); j < count; j++ {
			arg := &st.Args[2*j] // loaded as the low 32 bits (see creatureXferArgument)
			kind := *memmap.PtrUint32(0x587000, 255608+uintptr(16*id+4*uint32(j)))
			switch kind {
			case 1:
				if *arg != 0 {
					*arg = uintptr(unsafe.Pointer(objectLookupByScriptID(uint32(*arg))))
				}
			case 2:
				if *arg != 0 {
					wp := GetServer().S().WPs.ByInd(int(*arg))
					*arg = uintptr(unsafe.Pointer(wp))
				}
			}
		}
	}
	return memmap.PtrWord(p)
}

func creatureXferVoice(u *server.Object) uint32 {
	r := objectXferStream{cryptfile.Global()}
	if r.read() {
		var name [256]byte
		n := r.byte(0)
		r.raw(unsafe.Pointer(&name[0]), int(n))
		set := Nox_xxx_getDefaultSoundSet_424350(alloc.GoString(&name[0]))
		if u == nil || u.ObjClass&2 == 0 {
			return 0
		}
		layout.As[server.MonsterUpdateData](u.UpdateData).SoundSet122 = set
		return 1
	}
	set := Nox_xxx_monsterGetSoundSet_424300(u)
	if set == nil {
		b := byte(0)
		if _, err := r.cf.ReadWrite(unsafe.Slice(&b, 1)); err != nil {
			return 0
		}
		return 1
	}
	name := set.Name
	n := r.byte(byte(len(alloc.GoString(name))))
	if _, err := r.cf.ReadWrite(unsafe.Slice(name, int(n))); err != nil {
		return 0
	}
	return 1
}
