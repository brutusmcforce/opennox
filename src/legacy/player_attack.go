package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

// attackPlayer keeps the original byte-sized animation counter: elapsed frame
// arithmetic wraps at uint32, then the quotient is narrowed before hit gates.
func attackPlayer(u *server.Object) int {
	ud := u.UpdateData
	var it *server.Object
	var bits uint32
	var anim, prior byte
	var priorOffset int
	player := u.ObjClass&4 != 0
	if player {
		pd := *(*unsafe.Pointer)(unsafe.Add(ud, 276))
		it = *temporaryRefWord(ud, 104)
		bits = *(*uint32)(unsafe.Add(pd, 4))
		anim = *(*byte)(unsafe.Add(pd, 8))
		priorOffset = 236
	} else {
		if u.ObjClass&2 == 0 || u.ObjSubClass&0x10 == 0 {
			return 0
		}
		bits = *(*uint32)(unsafe.Add(ud, 2056))
		it = *temporaryRefWord(ud, 2064)
		anim = *(*byte)(unsafe.Add(ud, 2068))
		priorOffset = 481
	}
	prior = *(*byte)(unsafe.Add(ud, priorOffset))
	var def *server.Modifier
	if it != nil {
		def = GetServer().S().Modif.Nox_xxx_getProjectileClassById413250(int(it.TypeInd))
		if def == nil {
			return 0
		}
		it.PosVec = u.PosVec
		it.PrevPos = u.PosVec
	} else if anim == 0 {
		anim = byte(nox_common_randomInt_415FA0(23, 24))
		if player {
			pd := *(*unsafe.Pointer)(unsafe.Add(ud, 276))
			if *(*byte)(unsafe.Add(pd, 2251)) == 0 && nox_common_randomInt_415FA0(0, 100) >= 75 {
				anim = 25
			}
			*(*byte)(unsafe.Add(pd, 8)) = anim
		} else {
			*(*byte)(unsafe.Add(ud, 2068)) = anim
		}
	}
	strength := equipmentStrength(u)
	core := GetServer().S()
	frame := core.Frame()
	var frames, delay int32
	var current byte
	animation := func(id int) { nox_xxx_animPlayerGetFrameRange_4F9F90(int(int32(id)), &frames, &delay) }
	elapsed := func() byte { return byte((frame - u.Field34) / uint32(delay+1)) }
	start := func(readiness int, immediate bool) {
		if player && *(*uint32)(ud) == 0 {
			end := frame
			if !immediate {
				end += uint32(frames * (delay + 1))
			}
			*(*uint32)(ud) = end - uint32(readiness)
		}
	}
	finish := func() int {
		stored := current
		if int(current) >= int(frames) {
			stored = byte(frames - 1)
		}
		*(*byte)(unsafe.Add(ud, priorOffset)) = stored
		return bool2int(int(current) < int(frames))
	}
	half := func() bool { return int(current) == int(frames)/2 && current > prior }
	melee := func(damageType byte, static uint32, area bool, sound int) {
		r := attackRecord{Owner: u, Weapon: it, Pos: u.PosVec, Type: damageType, HitStatic: static, Front: 1}
		r.Damage = float32(controlBoltDamage(int32(strength), unsafe.Pointer(def)))
		r.Radius = float32(float64(u.Shape.Circle.R) + float64(def.Range68))
		if area {
			d := attackDirection(u)
			r.Pos = types.Pointf{X: float32(float64(d.X)*35 + float64(u.PosVec.X)), Y: float32(float64(d.Y)*35 + float64(u.PosVec.Y))}
			r.Radius = def.Range68
		}
		attackItemEffects(it, u, &r)
		hit := attackTrace(u, &r)
		if area {
			gameplayReportEarthquake(&u.PosVec, int(floatToInt32(float32(float64(strength)*0.1))))
			inventorySound(882, u, 0, 0)
		} else if hit == 0 {
			inventorySound(sound, u, 0, 0)
		}
	}
	if nox_common_playerIsAbilityActive_4FC250(asObjectC(u), 2) != 0 && nox_xxx_probablyWarcryCheck_4FC3E0(asObjectC(u), 2) != 0 {
		animation(46)
		current = elapsed()
		if uint16(current)<<8|uint16(prior) == 770 {
			gameplayReportEarthquake(&u.PosVec, 15)
			core.Map.EachObjInCircle(u.PosVec, 300, func(t *server.Object) bool { attackWarcry(t, u); return true })
			nox_xxx_castCounterSpell_52BBB0(13, asObjectC(u), asObjectC(u), asObjectC(u))
		}
		if int(current) >= int(frames) {
			sub_4FC440(asObjectC(u), 2)
		}
		return finish()
	}
	if nox_common_playerIsAbilityActive_4FC250(asObjectC(u), 1) != 0 {
		if u.HasEnchant(25) || u.HasEnchant(5) {
			return 0
		}
		animation(45)
		current = elapsed()
		speed := float64(u.SpeedBase) * 6
		u.SpeedCur = float32(speed)
		d := attackDirection(u)
		u.ForceVec.X = float32(speed*float64(d.X) + float64(u.ForceVec.X))
		u.ForceVec.Y = float32(speed*float64(d.Y) + float64(u.ForceVec.Y))
		if int(current) >= int(frames)-1 {
			current = 0
		}
		return finish()
	}
	if it == nil {
		animation(int(anim))
		start(0, false)
		current = elapsed()
		if int(current) < int(frames) {
			return finish()
		}
		unarmed, free := alloc.New(server.Modifier{})
		defer free()
		if anim == 23 || anim == 24 {
			unarmed.DamageCoeffOrArmor64 = float32(0.039999999)
			unarmed.DamageMin72 = 5
		} else if anim == 25 {
			unarmed.DamageCoeffOrArmor64 = float32(0.039999999)
			unarmed.DamageMin72 = 10
		}
		r := attackRecord{Owner: u, Pos: u.PosVec, Type: 10, Front: 1, Radius: float32(float64(u.Shape.Circle.R) + 20)}
		r.Damage = float32(controlBoltDamage(int32(strength), unsafe.Pointer(unarmed)))
		if attackTrace(u, &r) == 0 {
			inventorySound(879, u, 0, 0)
		}
		return finish()
	}
	if bits&0x47f8000 != 0 {
		kind := 31
		if bits&0x8000 != 0 || *(*byte)(unsafe.Add(it.UseData.Ptr, 96))&2 != 0 {
			kind = 29
		}
		animation(kind)
		start(0, kind != 29)
		current = elapsed()
		if int(current) >= int(frames) {
			if bits&0x47f0000 != 0 {
				*(*uint32)(unsafe.Add(it.UseData.Ptr, 96)) &^= 2
			}
			return finish()
		}
		if kind == 29 && half() {
			melee(0, 0, false, 879)
		} else if u.ObjClass&2 != 0 && kind != 29 && current == 1 && prior == 0 {
			data := it.UseData.Ptr
			effectsUse(u, it)
			if *(*byte)(unsafe.Add(data, 108)) == 0 && *(*byte)(unsafe.Add(data, 109)) != 0 {
				equipmentNPCDequipWeapon(u, it)
			}
		}
		return finish()
	}
	if bits&0x7800000 != 0 {
		kind := 31
		if bits&0x3800000 != 0 {
			kind = 32
		}
		animation(kind)
		start(0, false)
		current = elapsed()
		if kind == 32 && half() {
			melee(0, 0, false, 879)
		}
		return finish()
	}
	if bits&0xc0 != 0 {
		animation(44)
		start(0, false)
		current = elapsed()
		if !half() {
			return finish()
		}
		pos := attackMuzzle(u)
		round := bits&0x40 != 0
		traceFlags := byte(4)
		name := "FanChakramInMotion"
		if round {
			traceFlags = 5
			name = "RoundChakramInMotion"
		}
		if attackRay(u.PosVec, pos, traceFlags) == 0 {
			inventorySound(323, u, 0, 0)
			return finish()
		}
		p := core.NewObjectByTypeID(name)
		if p == nil {
			return 0
		}
		if round {
			data := p.UpdateData
			inventoryRemove(u, it)
			GetServer().CreateObjectAt(p, u, pos)
			inventoryInsert(p, it, 1)
			Nox_xxx_modifSetItemAttrs_4E4990(p, it.InitDataModifier())
			attackProjectileVelocity(u, p)
			*(*byte)(unsafe.Add(data, 4)) = 4
			*(*types.Pointf)(unsafe.Add(data, 16)) = u.PosVec
			*(*byte)(unsafe.Add(data, 24)) = 2
			inventorySound(891, u, 0, 0)
		} else {
			b := unsafe.Slice((*byte)(it.UseData.Ptr), 3)
			*temporaryRefWord(p.CollideData, 4) = u
			GetServer().CreateObjectAt(p, u, pos)
			Nox_xxx_modifSetItemAttrs_4E4990(p, it.InitDataModifier())
			attackProjectileVelocity(u, p)
			inventorySound(891, u, 0, 0)
			if b[2] == 0 {
				b[1]--
				if b[1] != 0 {
					if player {
						attackReportAmmo(u, it)
					}
				} else {
					inventoryRemove(u, it)
					GetServer().DelayedDelete(it)
					attackReload(u, 128)
				}
			}
		}
		return finish()
	}
	// Preserve priority for masks that contain more than one weapon class.
	kind, sound := 0, 0
	damageType := byte(0)
	static := uint32(0)
	area := false
	switch {
	case bits&0x200 != 0:
		kind = 28
		sound = 880
	case bits&0x100 != 0:
		kind = 27
		sound = 881
	case bits&0x400 != 0:
		kind = 37
		sound = 881
	case bits&0x4000 != 0:
		kind = 39
		damageType = 2
		static = 1
		area = true
	case bits&0x800 != 0:
		kind = 26
		sound = 884
		damageType = 2
		static = 1
	case bits&0x3000 != 0:
		kind = 35
		sound = 883
		static = 1
	case bits&4 != 0:
		animation(33)
		if player && *(*uint32)(ud) == 0 {
			start(int(effectsReadiness(it)), false)
		}
		current = elapsed() + byte(effectsReadiness(it))
		if int(current) >= int(frames)-1 && current > prior {
			attackBow(u, it)
			current = byte(frames)
		}
		return finish()
	case bits&8 != 0:
		animation(34)
		if player && *(*uint32)(ud) == 0 {
			start(int(effectsReadiness(it)), false)
		}
		current = elapsed()
		if current == 1 && prior == 0 {
			attackBow(u, it)
		}
		if current >= 1 {
			current += byte(effectsReadiness(it))
		}
		return finish()
	default:
		return finish()
	}
	animation(kind)
	start(0, false)
	current = elapsed()
	if half() {
		melee(damageType, static, area, sound)
	}
	return finish()
}
