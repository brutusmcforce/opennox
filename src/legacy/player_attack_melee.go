package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

type attackRecord struct {
	Damage    float32
	Type      byte
	_         [3]byte
	Radius    float32
	Owner     *server.Object
	Pos       types.Pointf
	HitStatic uint32
	Weapon    *server.Object
	Front     byte
	_         [3]byte
}

func attackPreEffects(t, u, it *server.Object, r *attackRecord) int {
	if it == nil {
		return 0
	}
	if !noxflags.HasGamePlay(1) && u != nil && u.ObjClass&6 != 0 && !GetServer().S().IsEnemyTo(u, t) {
		return 0
	}
	if t.HasEnchant(23) || t.HasEnchant(27) {
		return 1
	}
	for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4)[2:] {
		if m != nil && m.AttackPreHit52.Fnc != nil {
			server.CallModifierEffect5(m.AttackPreHit52.Fnc, m, it, u, t, unsafe.Pointer(r))
		}
	}
	return 0
}
func attackItemEffects(it, u *server.Object, r *attackRecord) int {
	for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4) {
		if m != nil && m.Attack40.Fnc != nil {
			server.CallModifierEffect5(m.Attack40.Fnc, m, it, u, nil, unsafe.Pointer(r))
		}
	}
	return 0
}
func attackHit(t *server.Object, r *attackRecord) {
	if t == nil || r == nil || r.Owner == nil || t == r.Owner {
		return
	}
	u := r.Owner
	if byte(stateFront(&u.PosVec, int32(int16(u.Direction1)), &t.PosVec))&r.Front == 0 {
		return
	}
	if t.ObjFlags&0x8040 != 0 || (r.HitStatic == 0 && t.ObjFlags&8 != 0) {
		return
	}
	if t.Material != 0x4000 {
		dword_5d4594_2488656 = 1
	}
	if attackRay(u.PosVec, t.PosVec, 5) == 0 {
		return
	}
	attackPreEffects(t, u, r.Weapon, r)
	_ = t.CallDamageValue(r.Owner, r.Weapon, effectsTruncWord(float64(r.Damage)+0.5), int32(r.Type))
	if noxflags.HasGame(2048) && r.Owner.ObjClass&4 != 0 && t.ObjClass&2 == 0 && t.HealthData != nil && t.HealthData.Max != 0 && t.ObjFlags&0x8020 == 0 {
		visibilityFXPoint(139, t.PosVec)
	}
	if r.Weapon == nil {
		u = r.Owner
		var anim byte
		if u.ObjClass&4 != 0 {
			anim = *(*byte)(unsafe.Add(*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 276)), 8))
		} else {
			anim = *(*byte)(unsafe.Add(u.UpdateData, 2068))
		}
		if anim == 25 {
			GetServer().ApplyForce(t, u.PosVec, 20)
		}
	}
}
func attackNearest(t, u *server.Object) {
	if u == nil || t == nil || u == t {
		return
	}
	f := t.ObjFlags
	if f&0x8049 != 0 || (t.ObjClass&6 == 0 && f&0x10 != 0 && f&0x80 == 0) || (t.ObjClass&6 != 0 && !GetServer().S().IsEnemyTo(u, t)) || !GetServer().S().CanInteract(u, t, 1) {
		return
	}
	limit := math.Float32frombits(uint32(dword_5d4594_2488652))
	if !(limit > 0) {
		return
	}
	dx := float32(float64(t.PosVec.X) - float64(u.PosVec.X))
	dy64 := float64(t.PosVec.Y) - float64(u.PosVec.Y)
	dy := float32(dy64)
	distance := math.Sqrt(dy64*float64(dy) + float64(dx)*float64(dx))
	if distance == 0 {
		return
	}
	dir := attackDirection(u)
	if !(float64(dy)/distance*float64(dir.Y)+float64(dx)/distance*float64(dir.X) > 0.5) {
		return
	}
	if t.Shape.Kind == 2 {
		distance -= float64(t.Shape.Circle.R)
	} else if t.Shape.Kind == 3 {
		delta := types.Pointf{X: dx, Y: dy}
		n := collisionBoxDistance(&u.PosVec, limit, t, &delta)
		if n < 0 {
			return
		}
		distance = float64(limit) - n
	}
	if distance < 0 {
		distance = 0
	}
	old := objectFromInt(int32(dword_5d4594_2488660))
	if (distance < float64(limit) || old != nil && old.ObjClass&2 == 0 && t.ObjClass&2 != 0) && (old == nil || old.ObjClass&2 == 0) {
		dword_5d4594_2488652 = uint32(math.Float32bits(float32(distance)))
		dword_5d4594_2488660 = memmap.PtrWord(t.CObj())
	}
}
func attackTrace(u *server.Object, r *attackRecord) int {
	if u == nil || r == nil {
		return 0
	}
	dword_5d4594_2488656 = 0
	dword_5d4594_2488660 = 0
	// This dependency includes shape extents; EachObjInCircle uses centers only.
	extra := float64(0)
	if r.Weapon != nil && r.Weapon.ObjSubClass&0x4000 != 0 {
		motionRadial(&r.Pos, r.Radius, func(candidate *server.Object) { attackHit(candidate, r) })
		extra = 25
	} else {
		dword_5d4594_2488652 = uint32(math.Float32bits(r.Radius))
		motionRadial(&u.PosVec, r.Radius, func(candidate *server.Object) { attackNearest(candidate, u) })
		if dword_5d4594_2488660 != 0 {
			attackHit(objectFromInt(int32(dword_5d4594_2488660)), r)
		}
	}
	bounds := [4]int32{
		floatToInt32(float32(float64(r.Pos.X)-float64(r.Radius)-extra)) / 23,
		floatToInt32(float32(float64(r.Pos.Y)-float64(r.Radius)-extra)) / 23,
		floatToInt32(float32(float64(r.Radius)+float64(r.Pos.X)+extra)) / 23,
		floatToInt32(float32(float64(r.Pos.Y)+float64(r.Radius)+extra)) / 23,
	}
	it := r.Weapon
	if it == nil {
		it = u
	}
	nox_xxx_mapDamageToWalls_534FC0((*[4]int32)(unsafe.Pointer(&bounds)), unsafe.Pointer(&r.Pos), float32(float32(extra+float64(r.Radius))), int(int32(effectsTruncWord(float64(r.Damage)+0.5))), int(int32(r.Type)), it.CObj())
	if r.Weapon != nil && dword_5d4594_2488656 != 0 {
		damage := float32(float64(nox_xxx_gamedataGetFloat_419D40(internCStr("ItemDamagePercentage"))) * float64(r.Damage))
		target := int32(dword_5d4594_2488660)
		damageDurability(r.Weapon, r.Owner, objectFromInt(int32(target)), objectFromInt(int32(target)), damage, int32(r.Type), true)
	}
	return int(dword_5d4594_2488656)
}
func attackWarcry(t, u *server.Object) int16 {
	if u == nil {
		return 0
	}
	if t == nil {
		return 0
	}
	if t.ObjClass&2 != 0 && t.ObjSubClass&0x20000 != 0 && t.ObjFlags&0x8020 == 0 {
		spellLifeApplyBuff(t, 5, 90, 3)
	}
	return int16(memmap.PtrWord(t.CObj()))
}
