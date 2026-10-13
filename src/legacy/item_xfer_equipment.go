package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func itemXferModifiers(r objectXferStream, u *server.Object) {
	if r.read() {
		attrs := server.ModifierInitData{Field16: 0xffffffff}
		for i := 0; i < 4; i++ {
			n := int(r.byte(0))
			var name [256]byte
			r.raw(unsafe.Pointer(&name[0]), n)
			mods := GetServer().S().Modif
			attrs.Modifiers[i] = mods.Nox_xxx_modifGetDescById413330(mods.Nox_xxx_modifGetIdByName413290(alloc.GoString(&name[0])))
		}
		stateAttributes(u, &attrs)
	} else {
		for _, m := range (*[4]*server.ModifierEff)(u.InitData) {
			if m == nil {
				r.byte(0)
				continue
			}
			name := m.Name()
			n := int(r.byte(byte(len(name))))
			r.cf.ReadWrite([]byte(name)[:n])
		}
	}
}
func itemXferHealth(r objectXferStream, u *server.Object, armor bool) {
	hp := r.short(uint16(resourceGetHP(u)))
	if hp > u.HealthData.Max {
		hp = u.HealthData.Max
	}
	if !r.read() {
		return
	}
	if Nox_xxx_gameIsSwitchToSolo_4DB240() || Nox_xxx_gameIsNotMultiplayer_4DB250() || noxflags.HasGame(4096) && GetServer().S().Players.AnyXxx() {
		resourceSetHP(u, hp)
		return
	}
	var mod *server.Modifier
	if armor {
		mod = GetServer().S().Modif.Nox_xxx_equipClothFindDefByTT413270(int(uint16(u.TypeInd)))
	} else {
		mod = GetServer().S().Modif.Nox_xxx_getProjectileClassById413250(int(uint16(u.TypeInd)))
	}
	if mod != nil {
		hp = uint16(mod.Durability52)
		u.HealthData.Max = hp
		u.HealthData.Field2 = hp
		resourceSetHP(u, hp)
	}
}
func itemXferWeapon(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 64)
	if !ok {
		return 0
	}
	if v < 11 && r.read() {
		attrs := server.ModifierInitData{Field16: 0xffffffff}
		stateAttributes(u, &attrs)
		return 1
	}
	itemXferModifiers(r, u)
	if v >= 41 && uint32(u.Class())&0x1000 != 0 && uint32(u.SubClass())&0x47f0000 != 0 && (v >= 62 || uint32(u.SubClass())&0x4000000 == 0) {
		p := u.UseDataWand()
		count := &p.Charges
		max := &p.MaxCharges
		amount := &p.Percent
		oldMax, oldAmount := *max, int32(*amount)
		a, b := r.byte(*count), r.byte(*max)
		c := *amount
		if v >= 61 {
			c = r.word(c)
		}
		if !noxflags.HasGame(4096) || a <= oldMax && int32(c) >= 0 && int32(c) <= oldAmount && oldMax == b {
			*count, *max, *amount = a, b, c
		} else {
			*count, *max, *amount = 0, oldMax, 0
		}
	}
	if v >= 42 {
		itemXferHealth(r, u, false)
	}
	if v == 63 {
		r.byte(0)
	}
	if v >= 64 {
		r.raw(unsafe.Pointer(&u.UpdateDataWeaponArmor().Flags), 4)
	}
	return objectXferFinish(r, u, v, saved)
}
func itemXferArmor(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 62)
	if !ok {
		return 0
	}
	if v < 11 && r.read() {
		var attrs server.ModifierInitData
		stateAttributes(u, &attrs)
		return 1
	}
	itemXferModifiers(r, u)
	if v >= 41 {
		itemXferHealth(r, u, true)
	}
	if v == 61 {
		r.byte(0)
	}
	if v >= 62 {
		r.raw(unsafe.Pointer(&u.UpdateDataWeaponArmor().Flags), 4)
	}
	return objectXferFinish(r, u, v, saved)
}
func itemXferAmmo(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	itemXferModifiers(r, u)
	p := u.UseDataAmmo()
	if r.read() {
		a, b := r.byte(0), r.byte(0)
		if noxflags.HasGame(4096) {
			p.Infinite = 0
		}
		p.Charges, p.MaxCharges = a, b
	} else {
		r.raw(unsafe.Pointer(&p.Charges), 1)
		r.raw(unsafe.Pointer(&p.MaxCharges), 1)
	}
	return objectXferFinish(r, u, v, saved)
}
func itemXferTeam(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	itemXferModifiers(r, u)
	if r.read() && uint32(u.Class())&0x10000000 != 0 {
		u.UpdateDataFlag().Home = u.PosVec
	}
	return objectXferFinish(r, u, v, saved)
}
