package legacy

import (
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func projectileDamage(t, owner, u *server.Object, damage, kind int32) int32 {
	return t.CallDamageValue(owner, u, damage, kind)
}
func projectileWall(u *server.Object, x, y, damage, kind int32) {
	GetServer().Nox_xxx_damageToMap_534BC0(int(x), int(y), int(damage), object.DamageType(kind), u)
}
func projectileHasContact() bool { return dword_5d4594_2488620 != 0 }
func projectileContact(u *server.Object, damage, kind int32) {
	if dword_5d4594_2488620 != 0 {
		projectileWall(u, int32(memmap.Uint32(0x5d4594, 2488612)), int32(memmap.Uint32(0x5d4594, 2488616)), damage, kind)
	}
}
func projectileGrid(p types.Pointf) (int32, int32) {
	return floatToInt32(float32(float64(p.X) * 0.043478262)), floatToInt32(float32(float64(p.Y) * 0.043478262))
}
func projectileFX(code byte, u *server.Object) {
	visibilityFXPoint(byte(code), u.PosVec)
}
func projectileFront(t, u *server.Object) bool {
	return stateFront(&t.PosVec, int32(int16(t.Direction1)), &u.PosVec)&1 != 0
}
func projectileGeneric(u, t *server.Object) {
	stone, imp := memmap.PtrUint32(0x5d4594, 1567948), memmap.PtrUint32(0x5d4594, 1567952)
	if *stone == 0 {
		*stone = uint32(GetServer().S().Types.IndByID("ThrowingStone"))
		*imp = uint32(GetServer().S().Types.IndByID("ImpShot"))
	}
	var damage int32
	name := ""
	if uint32(u.TypeInd) == *stone {
		name = "UrchinStoneDamage"
	} else if uint32(u.TypeInd) == *imp {
		name = "ImpShotDamage"
	}
	if name != "" {
		damage = floatToInt32(float32(nox_xxx_gamedataGetFloat_419D40(internCStr(name))))
	} else {
		damage = int32(*equipmentWord(u.CollideData, 0))
	}
	if t != nil {
		if byte(projectileDamage(t, u.FindOwnerChainPlayer(), u, damage, 11)) == 0 {
			return
		}
	} else {
		projectileContact(u, damage, 11)
	}
	GetServer().DelayedDelete(u)
}
func projectileGenericSpark(u, t *server.Object) {
	damage := int32(*equipmentWord(u.CollideData, 0))
	if t != nil {
		if byte(projectileDamage(t, u.FindOwnerChainPlayer(), u, damage, 11)) == 0 {
			return
		}
	} else {
		x, y := projectileGrid(u.NewPos)
		projectileWall(u, x, y, damage, 11)
	}
	GetServer().DelayedDelete(u)
}
func projectileDamageField(u, t *server.Object) {
	if t == nil || t.HealthData == nil {
		return
	}
	damage := int32(*(*byte)(u.CollideData)) >> 1
	if *(*byte)(u.CollideData) == 1 && GetServer().S().Frame()&1 != 0 {
		damage = 1
	}
	projectileDamage(t, u.FindOwnerChainPlayer(), u, damage, int32(*equipmentWord(u.CollideData, 4)))
}
func projectileManaDrain(u, t *server.Object) {
	if t == nil || t.ObjClass&4 == 0 || *(*uint16)(unsafe.Add(t.UpdateData, 4)) == 0 {
		return
	}
	resourceSubMana(t, int32(*(*byte)(u.CollideData)))
	stamp := (*uint16)(unsafe.Add(u.CObj(), int(layoutServerObject.Off(542)) /* Field542 */))
	core := GetServer().S()
	if core.Frame()-uint32(int32(int16(*stamp))) > uint32(int32(core.TickRate())>>1) {
		inventorySound(228, u, 0, 0)
		*stamp = uint16(core.Frame())
	}
}
func projectileBomb(u, t *server.Object) {
	if noxflags.HasGame(2048) && GetServer().S().Players.FirstUnit().ObjFlags&2 != 0 {
		return
	}
	nox_xxx_scriptCallByEventBlock_502490(unsafe.Add(u.UpdateData, 1272), t.CObj(), u.CObj(), 21)
	if t != nil && t.ObjClass&6 != 0 && !Nox_xxx_unitsHaveSameTeam_4EC520(u, t) && t.ObjFlags&0x8000 == 0 {
		*temporaryRefWord(u.UpdateData, 2176) = t
		resourceDamage(u, 999)
	}
}
func projectileDie(u, t *server.Object) {
	if t == nil || Nox_xxx_unitsHaveSameTeam_4EC520(u, t) || t.ObjClass&6 == 0 {
		return
	}
	u.ObjFlags |= 0x8000
	if u.Death != nil {
		u.CallDeath()
	} else {
		GetServer().DelayedDelete(u)
	}
}
func projectileTrapEligible(u, t *server.Object) bool {
	if noxflags.HasGame(2048) && GetServer().S().Players.FirstUnit().ObjFlags&2 != 0 {
		return false
	}
	ok := true
	if t != nil {
		if Nox_xxx_unitsHaveSameTeam_4EC520(u, t) {
			ok = false
		} else if noxflags.HasGame(512) {
			a, b := u.FindOwnerChainPlayer(), t.FindOwnerChainPlayer()
			if a.ObjClass&4 != 0 && b.ObjClass&4 != 0 {
				ok = false
			}
		}
	}
	if nox_common_playerIsAbilityActive_4FC250(asObjectC(t), 4) != 0 {
		ok = false
	}
	return ok
}
func projectileSulphur(u, t *server.Object, n *types.Pointf) {
	if t != nil {
		if Nox_xxx_unitsHaveSameTeam_4EC520(u, t) {
			return
		}
		damage := int32(*equipmentWord(u.CollideData, 0))
		if noxflags.HasGame(4096) && u.Collide == collisionKey(collisionIdentityYellowStarShot) {
			damage *= 3
		}
		if projectileDamage(t, u.FindOwnerChainPlayer(), u, damage, 11) != 0 {
			GetServer().DelayedDelete(u)
		}
	} else if n != nil {
		collisionReflect(n, &u.VelVec)
		x, y := projectileGrid(u.NewPos)
		projectileWall(u, x, y, int32(*equipmentWord(u.CollideData, 0)), 11)
	}
}
func projectileSulphurTrail(u, t *server.Object, n *types.Pointf) {
	if u != nil && !noxflags.HasGame(4) {
		projectileFX(136, u)
	}
	projectileSulphur(u, t, n)
}
func projectileDeathFragment(u, t *server.Object, n *types.Pointf) {
	if t != nil {
		projectileDamage(t, u.FindOwnerChainPlayer(), u, 20, 2)
	} else if n != nil {
		collisionReflect(n, &u.VelVec)
		inventorySound(37, u, 0, 0)
		x, y := projectileGrid(u.NewPos)
		projectileWall(u, x, y, 20, 2)
		return
	}
	GetServer().DelayedDelete(u)
}
func projectileWallSpark(u, t *server.Object, n *types.Pointf) {
	if t != nil {
		if projectileDamage(t, u.FindOwnerChainPlayer(), u, int32(*equipmentWord(u.CollideData, 0)), 11) != 0 {
			GetServer().DelayedDelete(u)
		}
	} else if n != nil {
		collisionReflect(n, &u.VelVec)
		x, y := projectileGrid(u.NewPos)
		projectileWall(u, x, y, int32(*equipmentWord(u.CollideData, 0)), 11)
	}
}
func projectileSparkOwner(u, t *server.Object) {
	if t != nil && t.ObjClass&4 != 0 && u.ObjOwner == nil {
		u.Field34 = GetServer().S().Frame()
		GetServer().S().ObjSetOwner(t, u)
	}
}
func projectileSpark(u, t *server.Object, n *types.Pointf) {
	switch *equipmentWord(u.UpdateData, 12) {
	case 4:
		return
	case 5:
		if t != nil {
			inventorySound(351, u, 0, 0)
			GetServer().DelayedDelete(u)
			*(*byte)(unsafe.Add(t.CObj(), int(layoutServerObject.Off(541)) /* Field541 */))++
			*(*uint16)(unsafe.Add(t.CObj(), int(layoutServerObject.Off(542)) /* Field542 */)) = 1000
			if t.ObjClass&4 != 0 {
				gameplayTextPrivate(t, (*byte)(unsafe.Pointer(internCStr("objcoll.c:WebbingSlow"))), 0)
			}
		}
	default:
		projectileSulphur(u, t, n)
	}
}
func projectileWeb(u, t *server.Object) {
	if t == nil {
		return
	}
	inventorySound(351, u, 0, 0)
	GetServer().DelayedDelete(u)
	if projectileDamage(t, u.FindOwnerChainPlayer(), u, 0, 2) != 0 {
		if t.ObjClass&6 != 0 {
			spellLifeApplyBuff(t, 4, int16(uint16(GetServer().S().TickRate())*4), 3)
		}
		if t.ObjClass&4 != 0 {
			gameplayTextPrivate(t, (*byte)(unsafe.Pointer(internCStr("objcoll.c:WebbingSlow"))), 0)
		}
	}
}
func projectileFist(u, t *server.Object) {
	if u.ZVal <= 0 && t != nil {
		projectileDamage(t, u.FindOwnerChainPlayer(), u, int32(*equipmentWord(u.UpdateData, 0)), 2)
	}
}
func projectileTeleportWake(u, t *server.Object) {
	if t == nil || t.HasEnchant(14) || (noxflags.HasGame(4096) && u.ObjOwner != nil && u.ObjOwner.ObjClass&4 == 0) || t.ObjClass&0x3001016 == 0 {
		return
	}
	if !t.HasEnchant(0) {
		projectileFX(138, t)
	}
	inventorySound(147, t, 0, 0)
	stateTeleport(t, (*types.Pointf)(u.CollideData))
	if !t.HasEnchant(0) {
		projectileFX(137, t)
	}
	inventorySound(147, t, 0, 0)
}
func projectileTrap(u, t *server.Object, gas bool) {
	if t == nil || !projectileTrapEligible(u, t) {
		return
	}
	name := "ClosedBearTrap"
	if gas {
		name = "ToxicCloud"
	}
	p := GetServer().S().NewObjectByTypeID(name)
	if p == nil {
		return
	}
	var owner server.Obj
	if u.ObjOwner != nil {
		owner = u.ObjOwner
	}
	GetServer().CreateObjectAt(p, owner, u.PosVec)
	if gas {
		*equipmentWord(p.UpdateData, 0) = uint32(floatToInt32(float32(float64(nox_xxx_gamedataGetFloat_419D40(internCStr("ToxicCloudLifetime"))) * float64(int32(GetServer().S().TickRate())))))
		inventorySound(847, u, 0, 0)
		GetServer().DelayedDelete(u)
	} else {
		GetServer().DelayedDelete(u)
		spellLifeApplyBuff(t, 5, 90, 5)
		spellLifeApplyBuff(t, 14, 90, 5)
		inventorySound(846, u, 0, 0)
	}
}
