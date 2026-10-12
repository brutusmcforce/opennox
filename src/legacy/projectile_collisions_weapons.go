package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func projectileBoltDamage(strength int32, d *server.Modifier) float64 {
	return float64(controlBoltDamage(int32(strength), unsafe.Pointer(d)))
}
func projectileArrow(u, t *server.Object) {
	core := GetServer().S()
	def := core.Modif.Nox_xxx_getProjectileClassById413250(int(u.TypeInd))
	if def == nil {
		return
	}
	cd := u.CollideData
	if u.ObjOwner == t && t != nil {
		return
	}
	strength := int32(30)
	if u.ObjOwner != nil {
		strength = int32(equipmentStrength(u.ObjOwner))
	}
	if noxflags.HasGame(4096) {
		owner := u.FindOwnerChainPlayer()
		if owner != nil && t != nil && owner.ObjClass&4 != 0 && t.ObjClass&4 != 0 && !core.IsEnemyTo(owner, t) {
			return
		}
	}
	if t == nil {
		if projectileHasContact() {
			projectileContact(u, effectsTruncWord(projectileBoltDamage(strength, def)), 11)
		}
		GetServer().DelayedDelete(u)
		return
	}
	equipmentStrength(u.ObjOwner)
	bolt := memmap.PtrUint32(0x5d4594, 1568000)
	if *bolt == 0 {
		*bolt = uint32(core.Types.IndByID("ArcherBolt"))
	}
	if t.ObjFlags&0x8000 != 0 {
		return
	}
	r := attackRecord{Weapon: u, Owner: *temporaryRefWord(cd, 4), Pos: u.PosVec, Type: 11, Radius: u.Shape.Circle.R, Damage: float32(projectileBoltDamage(strength, def))}
	attackItemEffects(u, r.Owner, &r)
	if u.ObjOwner != t {
		attackPreEffects(t, *temporaryRefWord(cd, 4), u, &r)
	}
	damage := effectsTruncWord(float64(r.Damage) + math.Float64frombits(uint64(qword_581450_9544)))
	accepted := byte(projectileDamage(t, u.FindOwnerChainPlayer(), u, damage, int32(r.Type)))
	if uint32(u.TypeInd) == *bolt {
		if t.HealthData != nil && t.HealthData.Cur == 0 {
			if t.HealthData.Max == 0 {
				GetServer().DelayedDelete(u)
			}
			return
		}
	} else if accepted == 0 {
		return
	}
	GetServer().DelayedDelete(u)
}
func projectileMonsterArrow(u, t *server.Object) {
	off := 4
	if noxflags.HasGame(2048) {
		off = 0
	}
	damage := int32(*equipmentWord(u.CollideData, off))
	if t != nil {
		if t.ObjFlags&0x8000 != 0 {
			return
		}
		projectileDamage(t, u.FindOwnerChainPlayer(), u, damage, 3)
	} else {
		projectileContact(u, damage, 11)
	}
	GetServer().DelayedDelete(u)
}
func projectileChakramCandidate(t *server.Object, pos *types.Pointf) {
	// These two decompiled expressions numerically convert the float
	// interpretation of the flag words; they are not bit casts.
	if uint32(effectsTruncWord(float64(math.Float32frombits(uint32(t.ObjClass)))))&0x20006 == 0 {
		return
	}
	if uint32(effectsTruncWord(float64(math.Float32frombits(uint32(t.ObjFlags)))))&0x8020 != 0 || t.HasEnchant(0) {
		return
	}
	if t == projectileChakramOwner || uintptr(t.CObj()) == uintptr(memmap.Uint32(0x5d4594, 1567932)) {
		return
	}
	u := projectileChakramObj
	if !GetServer().S().MapTraceVision(t, u) {
		return
	}
	dx, dy := float64(pos.X)-float64(t.PosVec.X), float64(pos.Y)-float64(t.PosVec.Y)
	d := dy*dy + dx*dx
	best := memmap.PtrFloat32(0x5d4594, 1567836)
	if d <= 160000 && d < float64(*best) {
		*best = float32(d)
		dword_5d4594_1567928 = memmap.PtrWord(t.CObj())
	}
}
func projectileChakramSelect(u *server.Object) *server.Object {
	dword_5d4594_1567928 = 0
	*memmap.PtrUint32(0x5d4594, 1567932) = *equipmentWord(u.UpdateData, 12)
	projectileChakramOwner = u.ObjOwner
	projectileChakramObj = u
	*memmap.PtrUint32(0x5d4594, 1567836) = 1259902592
	rect := types.Rectf{Min: types.Pointf{X: float32(float64(u.PosVec.X) - 400), Y: float32(float64(u.PosVec.Y) - 400)}, Max: types.Pointf{X: float32(float64(u.PosVec.X) + 400), Y: float32(float64(u.PosVec.Y) + 400)}}
	GetServer().S().Map.EachObjInRect(rect, func(t *server.Object) bool { projectileChakramCandidate(t, &u.PosVec); return true })
	t := objectFromWord(uint32(int32(dword_5d4594_1567928)))
	if t != nil {
		*(*byte)(unsafe.Add(u.UpdateData, 24)) = 2
		dx, dy := float64(t.PosVec.X)-float64(u.PosVec.X), float64(t.PosVec.Y)-float64(u.PosVec.Y)
		y := float32(dy)
		length := float32(math.Sqrt(dy*float64(y)+dx*dx) + math.Float64frombits(uint64(qword_581450_10176)))
		u.VelVec.X = float32(dx * float64(u.SpeedCur) / float64(length))
		u.VelVec.Y = float32(float64(y) * float64(u.SpeedCur) / float64(length))
	}
	return t
}
func projectileChakramFallback(u *server.Object) {
	length := float32(math.Sqrt(float64(u.VelVec.X)*float64(u.VelVec.X) + float64(u.VelVec.Y)*float64(u.VelVec.Y)))
	dir := int32(int32(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&u.VelVec))))))
	dir = (int32(nox_common_randomInt_415FA0(-64, 64)) + dir + 128) & 255
	off := uintptr(194136 + 8*dir)
	u.VelVec.X = float32(float64(length) * float64(memmap.Float32(0x587000, off)))
	vy := float64(length) * float64(memmap.Float32(0x587000, off+4))
	u.NewPos = u.PrevPos
	u.PosVec.X = u.PrevPos.X
	u.VelVec.Y = float32(vy)
	u.PosVec.Y = u.PrevPos.Y
	Nox_xxx_moveUpdateSpecial_517970(u)
}
func projectileChakram(u, t *server.Object, n *types.Pointf) {
	ud := u.UpdateData
	it := u.InvFirstItem
	if it == nil || it.ObjFlags&0x20 != 0 {
		GetServer().DelayedDelete(u)
		return
	}
	if t == nil || t.Material&0x30 != 0 {
		projectileFX(150, u)
	}
	count := (*byte)(unsafe.Add(ud, 4))
	mode := (*byte)(unsafe.Add(ud, 24))
	if t == nil {
		if n == nil {
			return
		}
		if *count != 0 {
			collisionReflect(n, &u.VelVec)
		} else {
			projectileChakramFallback(u)
		}
		advance := true
		if *count != 0 {
			*count--
		} else if *mode == 0 {
			*mode = 2
			advance = false
		}
		projectileContact(u, 1, 0)
		if *mode == 1 {
			inventoryDrop(u, it, &u.PosVec)
			GetServer().DelayedDelete(u)
			return
		}
		if *count == 0 && advance {
			*mode = 0
			*temporaryRefWord(ud, 8) = u.ObjOwner
			return
		}
		projectileChakramSelect(u)
		return
	}
	if t == u.ObjOwner {
		inventoryRemove(u, it)
		inventoryInsert(u.ObjOwner, it, 1)
		if owner := u.ObjOwner; owner != nil && owner.ObjClass&4 != 0 && *temporaryRefWord(owner.UpdateData, 104) == nil {
			equipmentEquipWeapon(owner, it, 1, 1)
		}
		inventorySound(892, u, 0, 0)
		GetServer().DelayedDelete(u)
		return
	}
	if Nox_xxx_unitsHaveSameTeam_4EC520(u, t) {
		return
	}
	owner := u.ObjOwner
	def := GetServer().S().Modif.Nox_xxx_getProjectileClassById413250(int(u.TypeInd))
	if owner == nil || owner.ObjFlags&0x8020 != 0 {
		inventoryDrop(u, it, &u.PosVec)
		GetServer().DelayedDelete(u)
		return
	}
	if t.ObjFlags&0x8000 != 0 || def == nil {
		return
	}
	strength := int32(equipmentStrength(owner))
	r := attackRecord{Pos: u.PosVec, Weapon: u, Owner: owner, Damage: float32(projectileBoltDamage(strength, def)), Radius: float32(float64(u.Shape.Circle.R) + 30)}
	attackItemEffects(u, owner, &r)
	attackPreEffects(t, owner, u, &r)
	projectileDamage(t, owner, u, effectsTruncWord(float64(r.Damage)+math.Float64frombits(uint64(qword_581450_9544))), 0)
	if t.ObjFlags&0x8020 == 0 {
		*temporaryRefWord(ud, 12) = t
	}
	if *count != 0 {
		damageReflect(u, t)
	} else {
		projectileChakramFallback(u)
	}
	advance := true
	if *count == 0 && *mode == 0 {
		*mode = 2
		advance = false
	}
	if *count != 0 {
		*count--
	}
	if *mode == 1 {
		it = u.InvFirstItem
		inventorySound(893, u, 0, 0)
		inventoryRemove(u, it)
		GetServer().CreateObjectAt(it, nil, u.PosVec)
		GetServer().DelayedDelete(u)
		return
	}
	if *count == 0 && advance {
		*temporaryRefWord(ud, 8) = owner
		*mode = 0
		return
	}
	projectileChakramSelect(u)
}
