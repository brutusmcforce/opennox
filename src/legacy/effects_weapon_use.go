package legacy

import (
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func effectsRecharge(it *server.Object, delta int32) int {
	if it.UseData.Ptr == nil {
		return 0
	}
	data := it.UseDataWand()
	percent := &data.Percent
	old := int32(*percent)
	if old >= 100 {
		return 0
	}
	next := old + delta
	if next >= 100 {
		next = 100
	}
	*percent = uint32(next)
	charge := *percent * uint32(data.MaxCharges) / 100
	if charge == uint32(data.Charges) {
		return 0
	}
	data.Charges = byte(charge)
	return 1
}
func effectsWandShot(u *server.Object, typ int, pos types.Pointf, dir uint32) uint32 {
	core := GetServer().S()
	shot := core.Objs.NewObject(core.Types.ByInd(typ))
	if shot == nil {
		return 0
	}
	GetServer().CreateObjectAt(shot, u, pos)
	shot.Direction1, shot.Direction2 = server.Dir16(dir), server.Dir16(dir)
	dx, dy := movementDirectionVector(int32(dir))
	// The first product spills across the second direction-table lookup;
	// the Y product stays in the x87 register until owner velocity is added.
	x := float32(float64(dx) * float64(shot.SpeedCur))
	shot.VelVec = types.Pointf{X: float32(float64(u.VelVec.X) + float64(x)), Y: float32(float64(u.VelVec.Y) + float64(dy)*float64(shot.SpeedCur))}
	return dir
}
func effectsLesserFireball(u, it *server.Object) int {
	if it.UseData.Ptr == nil {
		return 1
	}
	data := it.UseDataWand()
	index := &data.MissileInd
	if *index == 0 {
		*index = uint32(GetServer().S().Types.IndByID(data.Missile()))
	}
	if data.MaxCharges != 0 && data.Charges == 0 {
		inventorySound(222, u, 0, 0)
		return 0
	}
	core := GetServer().S()
	if core.Frame()-data.LastUse < data.Delay-uint32(effectsReadiness(it)) {
		return 0
	}
	radius := float64(u.Shape.Circle.R) + 4
	dir := int32(int16(u.Direction1))
	dx, dy := movementDirectionVector(dir)
	pos := u.PosVec
	x := float32(radius*float64(dx) + float64(pos.X))
	end := types.Pointf{X: float32(float64(x) + float64(u.VelVec.X)), Y: float32(radius*float64(dy) + float64(pos.Y) + float64(u.VelVec.Y))}
	if !core.MapTraceRayAt(pos, end, nil, nil, 5) {
		end = pos
	}
	effectsWandShot(u, int(*index), end, uint32(dir))
	if data.Flags&1 != 0 {
		effectsWandShot(u, int(*index), end, uint32((dir+8)&255))
		effectsWandShot(u, int(*index), end, uint32((dir-8)&255))
	}
	inventorySound(int(data.Sound), u, 0, 0)
	if data.MaxCharges != 0 {
		maximum := data.MaxCharges
		data.Charges--
		data.Percent = 100 * uint32(data.Charges) / uint32(maximum)
		if u.ObjClass&4 != 0 {
			equipmentReportWandCharges(u, it)
		}
	}
	data.LastUse = core.Frame()
	return 1
}
func effectsWandCast(u, it *server.Object) int {
	data := it.UseDataWand()
	cache := memmap.PtrUint32(0x5d4594, 2488732)
	if *cache == 0 {
		*cache = uint32(GetServer().S().Types.IndByID("ForceWand"))
	}
	if data.MaxCharges != 0 && data.Charges == 0 {
		inventorySound(222, u, 0, 0)
		return 0
	}
	core := GetServer().S()
	if core.Frame()-data.LastUse < data.Delay-2*uint32(effectsReadiness(it)) {
		return 0
	}
	arg, free := alloc.New(server.SpellAcceptArg{})
	defer free()
	arg.Obj = u
	if u.ObjClass&4 != 0 {
		pl := equipmentPlayer(u)
		arg.Pos = types.Pointf{X: float32(int32(pl.CursorVec.X)), Y: float32(int32(pl.CursorVec.Y))}
		if target := layout.As[server.PlayerUpdateData](u.UpdateData).CursorObj; target != nil {
			arg.Obj = target
		}
	} else {
		arg.Pos = u.PosVec
	}
	id := spell.ID(data.Spell)
	data.Flags |= 4
	if GetServer().Nox_xxx_spellAccept4FD400(id, u, u, u, arg, 4) {
		data.LastUse = core.Frame()
		if it.ObjClass&0x1000 == 0 || it.ObjSubClass&0x4040000 == 0 {
			maximum := data.MaxCharges
			if maximum != 0 {
				data.Charges--
				data.Percent = 100 * uint32(data.Charges) / uint32(maximum)
				if u.ObjClass&4 != 0 {
					equipmentReportWandCharges(u, it)
				}
			}
		}
	}
	return 1
}
func effectsFireWand(u, it *server.Object) int {
	dx, dy := movementDirectionVector(int32(int16(u.Direction1)))
	radius := float64(u.Shape.Circle.R)
	x, y := float32(float64(dx)*radius*1.5), float32(float64(dy)*radius*1.5)
	rng := GetServer().S().Rand.Logic
	speed := float32(rng.FloatClamp(12, 25))
	vx := float32(rng.FloatClamp(-2, 2) + float64(speed)*float64(dx))
	vy := float32(rng.FloatClamp(-2, 2) + float64(speed)*float64(dy))
	pos := types.Pointf{X: float32(float64(x) + float64(u.PosVec.X)), Y: float32(float64(y) + float64(u.PosVec.Y))}
	temporarySpark(pos, types.Ptf(vx, vy), 1, 20, 0, nil)
	frame := GetServer().S().Frame()
	timestamp := &it.Field34
	if frame-*timestamp > uint32(GetServer().S().TickRate()) {
		inventorySound(9, u, 0, 0)
		*timestamp = GetServer().S().Frame()
	}
	return 0
}
func effectsUse(u, it *server.Object) int {
	if it == nil || it.Use.Ptr == nil {
		return 1
	}
	if result := int(sub_419E60(asObjectC(u))); result == 1 {
		return result
	}
	return int(it.Use.CallResult(u, it))
}
