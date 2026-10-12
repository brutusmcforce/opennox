package legacy

import (
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"strings"
	"unsafe"
)

func monsterLoadCallback(def unsafe.Pointer, name string, table, field uintptr) bool {
	if strings.EqualFold(name, "NULL") {
		*(*unsafe.Pointer)(unsafe.Add(def, field)) = nil
		return true
	}
	for off := table; ; off += 8 {
		p := *memmap.PtrPtr(0x587000, off)
		if p == nil {
			return false
		}
		if alloc.GoString((*byte)(p)) == name {
			*(*unsafe.Pointer)(unsafe.Add(def, field)) = *memmap.PtrPtr(0x587000, off+4)
			return true
		}
	}
}

// The original callback numerically converts the float interpretation of these
// words to a byte. Its compiled x87 conversion uses a truncating signed word.
func monsterCallbackByte(word uint32) byte {
	v := math.Float32frombits(word)
	if math.IsNaN(float64(v)) || v >= 32768 || v <= -32769 {
		return 0
	}
	return byte(int32(v))
}
func monsterMeleeCandidate(t, u *server.Object) {
	cl := monsterCallbackByte(uint32(t.Class()))
	if t == u || (monsterCallbackByte(uint32(t.Flags()))&0x11 != 0 && cl&6 == 0) {
		return
	}
	if cl&6 == 0 && (t.HealthData == nil || t.HealthData.Max == 0) {
		return
	}
	core := GetServer().S()
	if !core.IsEnemyTo(u, t) && *memmap.PtrUint32(0x5D4594, 2491568) == 0 {
		return
	}
	dx := float64(t.PosVec.X) - float64(u.PosVec.X)
	dy := float32(float64(t.PosVec.Y) - float64(u.PosVec.Y))
	dir := unsafe.Slice(memmap.PtrFloat32(0x587000, uintptr(194136+8*int(int16(u.Direction1)))), 2)
	dist := float32(math.Sqrt(float64(dy)*float64(dy)+dx*dx) + .001)
	dot := float64(dy)/float64(dist)*float64(dir[1]) + dx/float64(dist)*float64(dir[0])
	if !(dot > .5) {
		return
	}
	gap := float32(float64(dist) - (float64(t.Shape.Circle.R) + float64(u.Shape.Circle.R)))
	nearest := memmap.PtrFloat32(0x5D4594, 2491572)
	if gap < *nearest {
		*nearest = gap
		*memmap.PtrPtr(0x5D4594, 2491564) = t.CObj()
	}
}
func monsterPickMeleeTarget(u *server.Object, all uint32) *server.Object {
	d := u.UpdateDataMonster().MonsterDef
	radius := float64(d.MeleeAttackRange112) + float64(u.Shape.Circle.R) + float64(*memmap.PtrFloat32(0x587000, 287328))
	rect := types.Rectf{Min: types.Pointf{X: float32(float64(u.PosVec.X) - radius), Y: float32(float64(u.PosVec.Y) - radius)}, Max: types.Pointf{X: float32(float64(u.PosVec.X) + radius), Y: float32(float64(u.PosVec.Y) + radius)}}
	*memmap.PtrUint32(0x5D4594, 2491568) = all
	*memmap.PtrPtr(0x5D4594, 2491564) = nil
	*memmap.PtrFloat32(0x5D4594, 2491572) = d.MeleeAttackRange112
	GetServer().S().Map.EachObjInRect(rect, func(t *server.Object) bool { monsterMeleeCandidate(t, u); return true })
	return (*server.Object)(*memmap.PtrPtr(0x5D4594, 2491564))
}
func monsterMeleePoison(u, t *server.Object) bool {
	d := u.UpdateDataMonster().MonsterDef
	if d.MeleeAttackPoisonChange136 == 0 {
		return false
	}
	if GetServer().S().Rand.Logic.IntClamp(1, 100) > int(int32(d.MeleeAttackPoisonChange136)) {
		return false
	}
	return resourcePoison(t, int32(d.MeleeAttackPoisonStrength140), int32(d.MeleeAttackPoisonMax144))
}
func monsterPoisonMessage(u, t *server.Object, kind int) {
	if !monsterMeleePoison(u, t) {
		return
	}
	name := "aifunc.c:Poisoned"
	switch kind {
	case 1:
		name = "aifunc.c:PoisonedByScorpion"
	case 2:
		name = "aifunc.c:PoisonedByZombie"
	case 5:
		name = "aifunc.c:PoisonedByWasp"
	}
	p, free := alloc.CString(name)
	defer free()
	gameplayTextPrivate(t, p, 0)
}

// kind follows the shipped strike table. Wasp applies poison before force;
// scorpion, vile zombie and spiders apply it afterwards.
func monsterStrike(u *server.Object, kind int) bool {
	switch kind {
	case 0:
		return monsterOgreStrike(u)
	case 3:
		*memmap.PtrUint32(0x5D4594, 2491560) = 0
		return monsterGolemStrike(u)
	case 4:
		*memmap.PtrUint32(0x5D4594, 2491560) = 1
		return monsterGolemStrike(u)
	case 9:
		return true
	}
	ud := u.UpdateDataMonster()
	t := monsterPickMeleeTarget(u, 0)
	if t == nil {
		return kind != 5 && kind != 6 && kind != 7
	}
	core := GetServer().S()
	if !core.MapTraceRay(u.PosVec, t.PosVec, 5) {
		return false
	}
	d := ud.MonsterDef
	t.CallDamage(u, u, int(d.MeleeAttackDamage116), object.DamageType(d.MeleeAttackDamageType124))
	if kind == 5 {
		monsterPoisonMessage(u, t, kind)
	}
	// Damage and poison callbacks can change the definition/impact.
	force := ud.MonsterDef.MeleeAttackImpact120
	if force > 0 {
		GetServer().ApplyForce(t, u.PosVec, float64(force))
	}
	switch kind {
	case 1, 2, 6, 7:
		monsterPoisonMessage(u, t, kind)
	case 8:
		Nox_xxx_buffApplyTo_4FF380(t, 5, int(2*uint16(core.TickRate())), 3)
		if st := u.MonsterPushAction(ai.ActionType(25)); st != nil {
			st.Args[0] = uintptr(math.Float32bits(t.PosVec.X))
			st.Args[1] = uintptr(math.Float32bits(t.PosVec.Y))
		}
		if st := u.MonsterPushAction(ai.ActionType(41)); st != nil {
			st.Args[0] = uintptr(core.Frame() + uint32(core.Rand.Logic.IntClamp(int(2*core.TickRate()), int(4*core.TickRate()))))
		}
		if st := u.MonsterPushAction(ai.ActionType(24)); st != nil {
			st.Args[0] = uintptr(math.Float32bits(t.PosVec.X))
			st.Args[1] = uintptr(math.Float32bits(t.PosVec.Y))
			st.Args[2] = 0
		}
	}
	return true
}
func monsterOgreCandidate(t, u *server.Object) {
	ud := u.UpdateDataMonster()
	if t == u {
		return
	}
	dx := float64(t.PosVec.X) - float64(u.PosVec.X)
	dy := float64(t.PosVec.Y) - float64(u.PosVec.Y)
	y := float32(dy)
	dist := float32(math.Sqrt(dy*dy+dx*dx) + .0099999998)
	if !(float64(dist)-(float64(t.Shape.Circle.R)+float64(u.Shape.Circle.R)) <= float64(ud.MonsterDef.MeleeAttackRange112)) {
		return
	}
	dir := unsafe.Slice(memmap.PtrFloat32(0x587000, uintptr(194136+8*int(int16(u.Direction1)))), 2)
	if !(float64(y)/float64(dist)*float64(dir[1])+dx/float64(dist)*float64(dir[0]) > .40000001) {
		return
	}
	if !GetServer().S().MapTraceRay(u.PosVec, t.PosVec, 5) {
		return
	}
	d := ud.MonsterDef
	t.CallDamage(u, u, int(d.MeleeAttackDamage116), object.DamageType(d.MeleeAttackDamageType124))
	GetServer().ApplyForce(t, u.PosVec, float64(ud.MonsterDef.MeleeAttackImpact120))
	*memmap.PtrUint32(0x5D4594, 2491556) = 1
}
func monsterOgreStrike(u *server.Object) bool {
	r := float32(float64(u.UpdateDataMonster().MonsterDef.MeleeAttackRange112) + float64(u.Shape.Circle.R) + float64(*memmap.PtrFloat32(0x587000, 287328)))
	*memmap.PtrUint32(0x5D4594, 2491556) = 0
	GetServer().S().Map.EachObjInCircle(u.PosVec, r, func(t *server.Object) bool { monsterOgreCandidate(t, u); return true })
	return *memmap.PtrUint32(0x5D4594, 2491556) != 0
}
func monsterAreaCandidate(t, u *server.Object) {
	ud := u.UpdateDataMonster()
	if t == u {
		return
	}
	if stateFront(&u.PosVec, int32(int16(u.Direction1)), &t.PosVec)&1 == 0 {
		return
	}
	if !(float64(stateDistance(u, t)) <= float64(ud.MonsterDef.MeleeAttackRange112)) {
		return
	}
	if !GetServer().S().MapTraceRay(u.PosVec, t.PosVec, 5) {
		return
	}
	d := ud.MonsterDef
	t.CallDamage(u, u, int(d.MeleeAttackDamage116), object.DamageType(d.MeleeAttackDamageType124))
	if t.Class()&6 != 0 {
		*memmap.PtrUint32(0x5D4594, 2491576) = 1
	}
	if force := ud.MonsterDef.MeleeAttackImpact120; force > 0 {
		GetServer().ApplyForce(t, u.PosVec, float64(force))
	}
}
func monsterGolemStrike(u *server.Object) bool {
	r := float32(float64(u.UpdateDataMonster().MonsterDef.MeleeAttackRange112) + float64(u.Shape.Circle.R) + float64(*memmap.PtrFloat32(0x587000, 287328)))
	*memmap.PtrUint32(0x5D4594, 2491576) = 0
	GetServer().S().Map.EachObjInCircle(u.PosVec, r, func(t *server.Object) bool { monsterAreaCandidate(t, u); return true })
	gameplayReportEarthquake(&u.PosVec, 30)
	return *memmap.PtrUint32(0x5D4594, 2491576) != 0
}
func monsterPointFX(u *server.Object, code byte) {
	visibilityFXPoint(byte(code), u.PosVec)
}
func monsterDeathExplosion(u *server.Object, big bool) {
	radius, force, damage, size := float32(96), float32(100), 96, byte(128)
	if noxflags.HasGame(2048) {
		damage = 30
	}
	if big {
		radius, force, damage, size = 150, 150, 148, 255
	}
	spellEffectPushAround(u.PosVec, radius, 10, force, u, nil, 0)
	mapDamageUnitsAround(u.PosVec, float32(radius), 10, int32(damage), 7, u, nil)
	visibilityFXSpark(u.PosVec, byte(size))
	GetServer().S().Audio.EventObj(42, u, 0, 0)
	GetServer().DelayedDelete(u)
}
func monsterDebrisPlace(u, t *server.Object, radius float32) {
	p, free := alloc.New(types.Pointf{})
	defer free()
	inventoryRandomPlacement(radius, &u.PosVec, p)
	GetServer().CreateObjectAt(t, nil, *p)
}
func monsterDebrisRaise(t *server.Object, lo, hi float64, field29 float32) {
	rng := GetServer().S().Rand.Logic
	Nox_xxx_unitRaise_4E46F0(t, float32(rng.FloatClamp(lo, hi)))
	t.Field27 = float32(rng.FloatClamp(-2, 0))
	t.Field29 = math.Float32bits(field29)
	t.ObjFlags |= 0x800000
}
func monsterDebrisDecay(t *server.Object, lo, hi int) {
	core := GetServer().S()
	delay := core.Rand.Logic.IntClamp(lo, hi)
	Nox_xxx_unitSetDecayTime_511660(t, int(core.TickRate()*uint32(delay)))
}
func monsterDeathDebris(u *server.Object) {
	core := GetServer().S()
	core.Audio.EventObj(494, u, 0, 0)
	monsterPointFX(u, 138)
	for off := uintptr(287976); ; off += 4 {
		name := *memmap.PtrPtr(0x587000, off)
		if name == nil {
			break
		}
		t := core.NewObjectByTypeID(alloc.GoString((*byte)(name)))
		if t == nil {
			break
		}
		monsterDebrisPlace(u, t, 30)
		monsterDebrisRaise(t, 10, 70, 2)
		monsterDebrisDecay(t, 10, 20)
	}
}
func monsterDeathChunks(u *server.Object) {
	core := GetServer().S()
	core.Audio.EventObj(487, u, 0, 0)
	monsterPointFX(u, 138)
	n := 6
	if noxflags.HasGame(2048) {
		n = core.Rand.Logic.IntClamp(20, 30)
	}
	index := uint32(dword_5d4594_2491580)
	for i := 0; i < n; i++ {
		name := *memmap.PtrPtr(0x587000, 288240+4*uintptr(index))
		t := core.NewObjectByTypeID(alloc.GoString((*byte)(name)))
		if t == nil {
			break
		}
		monsterDebrisPlace(u, t, 30)
		Nox_xxx_unitRaise_4E46F0(t, float32(core.Rand.Logic.FloatClamp(10, 70)))
		t.Field27 = float32(core.Rand.Logic.FloatClamp(-2, 0))
		t.ObjFlags |= 0x800000
		t.Field29 = math.Float32bits(float32(*(*byte)(memmap.PtrOff(0x587000, 287332+uintptr(dword_5d4594_2491580)))))
		GetServer().ApplyForce(t, u.PosVec, float64(float32(core.Rand.Logic.FloatClamp(5, 20))))
		if noxflags.HasGame(2048) {
			monsterDebrisDecay(t, 10, 20)
		} else {
			monsterDebrisDecay(t, 5, 10)
		}
		index = (uint32(dword_5d4594_2491580) + 1) % *memmap.PtrUint32(0x587000, 287344)
		dword_5d4594_2491580 = uint32(index)
	}
}
func monsterDeathTroll(u *server.Object) {
	core := GetServer().S()
	id := monsterCache(2491584, "SmallToxicCloud")
	t := core.NewObjectByTypeInd(int(id))
	if t == nil {
		return
	}
	data := (*uint32)(t.UpdateData)
	GetServer().CreateObjectAt(t, u, u.PosVec)
	core.Audio.EventObj(644, u, 0, 0)
	lifetime := float32(core.Balance.Float("SmallToxicCloudLifetime") * float64(int32(core.TickRate())))
	*data = uint32(int32(lifetime))
}
func monsterDropLoot(u *server.Object, name string, mods [4]string, ammo int) {
	if !noxflags.HasGame(2048) {
		return
	}
	core := GetServer().S()
	t := core.NewObjectByTypeID(name)
	if t == nil {
		return
	}
	monsterDebrisPlace(u, t, 50)
	if t.Class()&0x13001000 != 0 && (mods[0] != "" || mods[1] != "" || mods[2] != "" || mods[3] != "") {
		data, free := alloc.New(server.ModifierInitData{})
		defer free()
		*data = server.ModifierInitData{} // C's fifth source word was undefined reserved padding.
		for i, name := range mods {
			data.Modifiers[i] = core.Modif.Nox_xxx_modifGetDescById413330(core.Modif.Nox_xxx_modifGetIdByName413290(name))
		}
		Nox_xxx_modifSetItemAttrs_4E4990(t, unsafe.Pointer(data))
	}
	if ammo != 0 && t.Class()&0x1000000 != 0 && t.SubClass()&0x82 != 0 {
		*(*byte)(unsafe.Add(t.UseData.Ptr, 1)) = byte(ammo)
	}
	if t.Class()&0x2000000 != 0 {
		equipmentArmorDropSound(t)
	} else if t.Class()&0x1001000 != 0 {
		equipmentDropSound(t)
	}
}
func monsterDeathSkull(u *server.Object) {
	core := GetServer().S()
	core.Audio.EventObj(356, u, 0, 0)
	monsterPointFX(u, 138)
	skull := core.NewObjectByTypeID("Skull")
	if skull == nil {
		return
	}
	monsterDebrisPlace(u, skull, 20)
	Nox_xxx_unitRaise_4E46F0(skull, 40)
	skull.Field27 = float32(core.Rand.Logic.FloatClamp(-2, 0))
	skull.Field29 = math.Float32bits(4)
	skull.ObjFlags |= 0x800000
	GetServer().ApplyForce(skull, u.PosVec, float64(float32(core.Rand.Logic.FloatClamp(5, 25))))
	var n int
	if noxflags.HasGame(2048) {
		monsterDebrisDecay(skull, 10, 20)
		n = core.Rand.Logic.IntClamp(10, 20)
	} else {
		monsterDebrisDecay(skull, 2, 5)
		n = core.Rand.Logic.IntClamp(5, 10)
	}
	index := uint32(dword_5d4594_2491588)
	for i := 0; i < n; i++ {
		name := *memmap.PtrPtr(0x587000, 288868+4*uintptr(index))
		t := core.NewObjectByTypeID(alloc.GoString((*byte)(name)))
		if t == nil {
			break
		}
		monsterDebrisPlace(u, t, 20)
		monsterDebrisRaise(t, 10, 35, 4)
		GetServer().ApplyForce(t, u.PosVec, float64(float32(core.Rand.Logic.FloatClamp(5, 25))))
		if noxflags.HasGame(2048) {
			monsterDebrisDecay(t, 10, 20)
		} else {
			monsterDebrisDecay(t, 2, 5)
		}
		index = (uint32(dword_5d4594_2491588) + 1) % *memmap.PtrUint32(0x587000, 287348)
		dword_5d4594_2491588 = uint32(index)
	}
}
func monsterDeathLoot(u *server.Object, kind int) {
	roll := GetServer().S().Rand.Logic.IntClamp(0, 100)
	if kind == 17 || kind == 18 {
		monsterDeathSkull(u)
	}
	switch kind {
	case 17, 18, 22:
		if roll <= 20 {
			return
		}
		material := "Material2"
		if kind == 22 {
			material = "Material1"
		}
		if roll <= 50 {
			monsterDropLoot(u, "Sword", [4]string{"WeaponPower1", material}, 0)
		} else {
			shield := "WoodenShield"
			if kind == 18 {
				shield = "SteelShield"
			}
			monsterDropLoot(u, shield, [4]string{"", material}, 0)
		}
	case 20:
		if roll > 20 {
			name := "Bow"
			if roll > 50 {
				name = "Quiver"
			}
			monsterDropLoot(u, name, [4]string{}, 0)
		}
	case 21:
		if roll > 25 {
			monsterDropLoot(u, "OgreAxe", [4]string{"WeaponPower1", "Material2"}, 0)
		}
	case 23:
		if roll > 25 {
			monsterDropLoot(u, "StaffWooden", [4]string{"WeaponPower1"}, 0)
		}
	case 24:
		if roll > 25 {
			monsterDropLoot(u, "FanChakram", [4]string{}, 5)
		}
	}
}
