package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/server"
)

func effectsInventory(u *server.Object, mask byte) int {
	if u == nil || mask == 0 {
		return 0
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjFlags&0x100 == 0 || it.ObjClass&0x13001000 == 0 {
			continue
		}
		for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4) {
			if m == nil {
				continue
			}
			for off := uintptr(200160); off < 200280; off += 20 {
				if m.Engage112 == memmap.LoadPtr(0x587000, off) && mask == *memmap.PtrUint8(0x587000, off+4) {
					return 1
				}
			}
		}
	}
	return 0
}
func effectsEngageFlag(u *server.Object, mask byte, sound int) {
	*(*byte)(unsafe.Pointer(&u.Field110)) |= mask
	inventorySound(sound, u, 0, 0)
}
func effectsDisengageFlag(u *server.Object, mask byte, sound int) {
	if effectsInventory(u, mask) == 0 {
		*(*byte)(unsafe.Pointer(&u.Field110)) &^= mask
	}
	inventorySound(sound, u, 0, 0)
}
func effectsSpeed(m *server.ModifierEff, u *server.Object, engage bool) {
	if u == nil || u.ObjClass&4 == 0 {
		return
	}
	sound := 60
	if engage {
		*(*byte)(unsafe.Pointer(&u.Field110)) |= 16
		u.SpeedBonus = float32(float64(m.EngageFloat120) + float64(u.SpeedBonus))
		sound = 59
	} else {
		if effectsInventory(u, 16) == 0 {
			*(*byte)(unsafe.Pointer(&u.Field110)) &^= 16
		}
		u.SpeedBonus = float32(float64(u.SpeedBonus) - float64(m.EngageFloat120))
	}
	gameplayReportSpeed(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), u, 0, uint32(math.Float32bits(u.SpeedBonus)))
	inventorySound(sound, u, 0, 0)
}
func effectsProtection(u *server.Object, fn unsafe.Pointer, buff int, key string, cap1, cap2 float64) float64 {
	if u == nil {
		return 0
	}
	var sum float64
	add := func(it *server.Object) {
		for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4) {
			if m != nil && m.Engage112 == fn {
				sum += float64(m.EngageFloat120)
			}
		}
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjFlags&0x100 != 0 && it.ObjClass&0x13001000 != 0 {
			add(it)
		}
	}
	if u.ObjClass&0x13001000 != 0 {
		add(u)
	}
	value := float32(sum)
	if sum > cap1 {
		value = float32(cap1)
	}
	result := float64(value)
	if u.Buffs&(1<<buff) != 0 {
		result = GetServer().S().Balance.FloatInd(key, int(u.BuffsPower[buff])-1) + float64(value)
	}
	if result > cap2 {
		result = cap2
	}
	return result
}
func effectsRegeneration(m *server.ModifierEff, it *server.Object) {
	if it == nil {
		return
	}
	u := it.InvHolder
	if u == nil || u.HealthData == nil {
		return
	}
	core := GetServer().S()
	frame, fps := core.Frame(), uint32(core.TickRate())
	if frame-u.Frame134 < fps || u.ObjFlags&0x8020 != 0 {
		return
	}
	maxHP := uint16(resourceGetMaxHP(u))
	if uint16(resourceGetHP(u)) >= maxHP {
		return
	}
	period := uint32(m.Update100.Val)
	if it.ObjClass&0x2000000 != 0 && equipmentArmorBits(it)&0x4000 != 0 {
		period /= 3
	}
	if frame%(period*fps/uint32(uint16(resourceGetMaxHP(u)))) == 0 {
		resourceAdjustHP(u, 1)
	}
}
func effectsReplenish(m *server.ModifierEff, it *server.Object) {
	if it == nil || GetServer().S().Frame()%uint32(m.Update100.Val) != 0 {
		return
	}
	if it.ObjClass&0x1000 == 0 {
		return
	}
	data := it.UseDataWand()
	old, max := data.Charges, data.MaxCharges
	data.Charges++
	if data.Charges > max {
		data.Charges = max
	}
	if old != data.Charges && it.InvHolder != nil && it.InvHolder.ObjClass&4 != 0 {
		equipmentReportWandCharges(it.InvHolder, it)
	}
}
func effectsGrip(m *server.ModifierEff, value *int32, invert bool) int {
	v := m.DefendCollide88.Val == 0
	if invert {
		v = !v
	}
	*value = int32(bool2int(v))
	return int(*value)
}
func effectsGripSearch(a, b, it, target *server.Object) int {
	if it == nil || it.ObjClass&0x3001000 == 0 {
		return 0
	}
	for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4)[2:] {
		if m != nil && m.DefendCollide88.Fnc == modifierKey(modifierIDGripEffect) {
			var v int32
			effectsGrip(m, &v, false)
			if v == 0 {
				return 1
			}
		}
	}
	return 0
}
func effectsStatus(m *server.ModifierEff, u, target *server.Object, stun bool) {
	if target.ObjClass&6 == 0 {
		return
	}
	arg, free := alloc.New(server.SpellAcceptArg{})
	defer free()
	*arg = server.SpellAcceptArg{Obj: target, Pos: target.PosVec}
	code := uint32(1)
	if stun {
		Nox_xxx_castStun_52C2C0(spell.ID(74), u, u, u, arg, int(int8(m.AttackPreHit52.Val)))
		code = 0
	} else {
		Nox_xxx_castConfuse_52C1E0(spell.ID(12), u, u, u, arg, int(int8(m.AttackPreHit52.Val)))
	}
	if target.ObjClass&4 != 0 {
		gameplayTextInformation(int(uint8(target.UpdateDataPlayer().Player.PlayerInd)), 13, unsafe.Pointer(&code))
	}
}
func effectsRecoil(m *server.ModifierEff, it, target *server.Object) {
	if it != nil && target != nil {
		GetServer().ApplyForce(target, it.PosVec, float64(m.AttackPreHit52.Valf))
	}
}

// The C effects truncate through signed 64-bit x87 conversion and then keep the
// low 32 bits. Invalid conversions produce the integer-indefinite low word zero.
func effectsTruncWord(v float64) int32 {
	if math.IsNaN(v) || v >= 9223372036854775808.0 || v < -9223372036854775808.0 {
		return 0
	}
	return int32(int64(v))
}
func effectsLightning(m *server.ModifierEff, it, u, target *server.Object) {
	if target == nil {
		return
	}
	_ = target.CallDamageValue(u, it, effectsTruncWord(float64(m.AttackPreHit52.Valf)), 9)
	pos := target.PosVec
	visibilityFXPoint(129, pos)
	inventorySound(225, target, 0, 0)
}
func effectsDrainMana(m *server.ModifierEff, u, target *server.Object, damage int32) {
	if u == nil || target == nil || target.ObjClass&6 == 0 {
		return
	}
	v := effectsTruncWord(float64(m.AttackPreDmg64.Valf) * float64(damage))
	if v == 0 {
		v = 1
	}
	mana := uint16(resourceGetMana(target))
	if int32(mana) < v {
		resourceSubMana(target, int32(mana))
		resourceAddMana(u, int16(mana))
	} else {
		resourceSubMana(target, v)
		resourceAddMana(u, int16(v))
	}
}
func effectsVampirism(m *server.ModifierEff, u, target *server.Object, damage int32) {
	if u == nil || target == nil || target.ObjClass&6 == 0 || target.ObjClass&2 != 0 && target.ObjSubClass&0x40 != 0 {
		return
	}
	v := effectsTruncWord(float64(damage) * float64(m.AttackPreDmg64.Valf))
	if v == 0 {
		v = 1
	}
	hp := int32(uint16(resourceGetHP(target)))
	if hp < v {
		v = hp
	}
	resourceAdjustHP(u, v)
}
func effectsPoison(m *server.ModifierEff, u, target *server.Object) {
	if target.ObjClass&4 != 0 && *(*byte)(unsafe.Add(target.UpdateData, 88)) == 16 && stateFront(&target.PosVec, int32(int16(target.Direction1)), &u.PosVec)&1 != 0 {
		return
	}
	if target.ObjClass&6 != 0 && resourcePoison(target, 1, m.AttackPreDmg64.Val) && target.ObjClass&4 != 0 {
		v := int32(2)
		gameplayTextInformation(int(uint8(target.UpdateDataPlayer().Player.PlayerInd)), 13, unsafe.Pointer(&v))
	}
}
func effectsSympathy(m *server.ModifierEff, u, target *server.Object, damage int32) {
	if u == nil || target == nil || target.ObjClass&6 == 0 {
		return
	}
	hp := int32(uint16(resourceGetHP(target)))
	if hp < damage {
		damage = hp
	}
	resourceDamage(u, effectsTruncWord(float64(damage)*float64(m.AttackPreDmg64.Valf)))
}
func effectsReadiness(it *server.Object) int32 {
	if it == nil || it.ObjClass&0x13001000 == 0 {
		return 0
	}
	for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4)[2:] {
		if m != nil && m.Attack40.Fnc == modifierKey(modifierIDReadinessEffect) {
			return m.Attack40.Val
		}
	}
	return 0
}
func effectsRechargeRate(it *server.Object) int32 {
	if it == nil {
		return 0
	}
	if it.ObjClass&0x1000 != 0 && it.ObjSubClass&0x4000000 != 0 {
		return floatToInt32(float32(GetServer().S().Balance.Float("OblivionStaffRechargeRate")))
	}
	for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4)[2:] {
		if m != nil && m.Attack40.Fnc == modifierKey(modifierIDReplenishmentEffect) {
			return m.Attack40.Val
		}
	}
	return 0
}
