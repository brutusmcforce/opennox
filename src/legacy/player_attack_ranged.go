package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func attackWeaponBits(it *server.Object) uint32 {
	return uint32(GetServer().S().Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(it))
}
func attackReload(u *server.Object, kind uint32) int {
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if attackWeaponBits(it) == kind {
			return equipmentEquipWeapon(u, it, 1, 1)
		}
	}
	return 0
}
func attackDirection(u *server.Object) types.Pointf {
	off := uintptr(194136 + 8*int(int16(u.Direction1)))
	return types.Pointf{X: memmap.Float32(0x587000, off), Y: memmap.Float32(0x587000, off+4)}
}
func attackMuzzle(u *server.Object) types.Pointf {
	d := attackDirection(u)
	r := float64(u.Shape.Circle.R) + 4
	return types.Pointf{X: float32(r*float64(d.X) + float64(u.PosVec.X)), Y: float32(r*float64(d.Y) + float64(u.PosVec.Y))}
}
func attackRay(a, b types.Pointf, flags byte) int {
	return bool2int(GetServer().S().MapTraceRayAt(a, b, nil, nil, server.MapTraceFlags(int8(flags))))
}
func attackProjectileVelocity(u, p *server.Object) {
	d := attackDirection(u)
	p.VelVec = types.Pointf{X: float32(float64(d.X) * float64(p.SpeedCur)), Y: float32(float64(d.Y) * float64(p.SpeedCur))}
	p.Direction1 = u.Direction1
	p.Direction2 = u.Direction1
}
func attackShotEffects(u, it, p *server.Object) int {
	for _, m := range itemModifiers(it)[2:] {
		if m == nil {
			continue
		}
		if m.AttackPreHit52.Fnc == modifierKey(modifierIDRecoilEffect) {
			p.InitDataModifier().Modifiers[3] = m
		} else if m.Attack40.Fnc == modifierKey(modifierIDProjectileSpeedEffect) {
			p.SpeedCur = float32(float64(m.Attack40.Valf) * float64(p.SpeedCur))
		}
	}
	return 0
}
func attackReportAmmo(u, it *server.Object) {
	b := it.UseDataAmmo()
	player := u.UpdateDataPlayer().Player
	gameplayReportCharges(int(player.PlayerInd), it, b.Charges, b.MaxCharges)
}
func attackShoot(u, ammo, it *server.Object, kind uint32) int {
	pos := attackMuzzle(u)
	result := attackRay(u.PosVec, pos, 5)
	if result == 0 {
		return result
	}
	name := "WeakArcherArrow"
	if ammo != nil {
		name = "ArcherBolt"
		if kind == 4 {
			name = "ArcherArrow"
		}
	}
	p := GetServer().S().NewObjectByTypeID(name)
	if p != nil {
		p.CollideDataProjectile().SetOwnerObj(u)
		GetServer().CreateObjectAt(p, u, pos)
		if ammo != nil {
			Nox_xxx_modifSetItemAttrs_4E4990(p, ammo.InitDataModifier())
		}
		attackShotEffects(u, it, p)
		attackProjectileVelocity(u, p)
	}
	if ammo != nil {
		b := ammo.UseDataAmmo()
		if b.Infinite == 0 && u.ObjClass&4 != 0 {
			b.Charges--
			attackReportAmmo(u, ammo)
			if b.Charges == 0 {
				GetServer().DelayedDelete(ammo)
			}
		}
	}
	sound := 886
	if kind == 4 {
		sound = 885
	}
	inventorySound(sound, u, 0, 0)
	return result
}
func attackBow(u, it *server.Object) int {
	kind := attackWeaponBits(it)
	b := &it.UseDataBow().Reload
	quest := noxflags.HasGame(4096)
	msg := func(s string) { nox_xxx_netPriMsgToPlayer_4DA2C0(asObjectC(u), (*int8)(internCStr(s)), 0) }
	emptySound := func() {
		if kind == 4 {
			if !quest {
				inventorySound(887, u, 0, 0)
			}
		} else if kind == 8 {
			inventorySound(888, u, 0, 0)
		}
	}
	if *b != 0 {
		msg("pattack.c:ReloadingQuiver")
		emptySound()
		*b--
		return 0
	}
	for ammo := u.InvFirstItem; ammo != nil; ammo = ammo.InvNextItem {
		if ammo.ObjFlags&0x100 != 0 && attackWeaponBits(ammo) == 2 {
			a := ammo.UseDataAmmo()
			if a.Charges != 0 || a.Infinite == 1 {
				attackShoot(u, ammo, it, kind)
				return 1
			}
		}
	}
	if quest && kind == 4 {
		attackShoot(u, nil, it, 4)
	}
	if u.ObjClass&4 == 0 {
		return 0
	}
	if attackReload(u, 2) == 1 {
		msg("pattack.c:ReloadQuiver")
		*b = 0
	} else if !(quest && kind == 4) {
		msg("pattack.c:NoQuiver")
	}
	emptySound()
	Nox_xxx_playerSetState_4FA020(u, 13)
	return 0
}
