package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func spellStartTeleport(d *server.DurSpell) uint32 {
	if d.Target48 == nil {
		d.Target48 = d.Caster16
	}
	caster, target := d.Caster16, d.Target48
	if worldTileWater(d.Pos2) {
		text := GetServer().S().Strings().GetStringInFile(strman.ID("UnseenTarget"), "C:\\NoxPost\\src\\Server\\Magic\\Spell\\ExecDur.c")
		Nox_xxx_netSendLineMessage_4D9EB0(target, text)
		sustainedAudio(231, target)
		return 1
	}
	if !spellEffectPlacement(target.PosVec, d.Pos2) {
		spellEffectInform(caster)
		return 1
	}
	return sustainedTeleportStart(d.C())
}

func spellStartPixies(id spell.ID, owner, origin *server.Object, level int) int {
	core := GetServer().S()
	cache := memmap.PtrUint32(0x5D4594, 2489140)
	radius := float32(origin.Shape.Circle.R + 4)
	if *cache == 0 {
		*cache = uint32(core.Types.IndByID("Pixie"))
	}
	count := int32(owner.CountSubOfType(int(*cache)))
	// Unlike teleport delay, the original count truncates the double directly.
	limit := int32(int64(core.Balance.FloatInd("PixieCount", level-1)))
	if count >= limit {
		return 1
	}
	for left := int32(uint32(limit) - uint32(count)); left > 0; left-- {
		dir := core.Rand.Logic.IntClamp(0, 255)
		// C keeps each product and coordinate addition wide until the point store.
		x := float64(radius)*float64(*memmap.PtrFloat32(0x587000, 194136+8*uintptr(dir))) + float64(origin.PosVec.X)
		y := float64(radius)*float64(*memmap.PtrFloat32(0x587000, 194140+8*uintptr(dir))) + float64(origin.PosVec.Y)
		pos := types.Ptf(float32(x), float32(y))
		if !core.MapTraceRay(origin.PosVec, pos, 5) {
			continue
		}
		u := core.NewObjectByTypeID("Pixie")
		if u == nil {
			continue
		}
		ud := u.UpdateDataPixie()
		var objOwner server.Obj
		if owner != nil {
			objOwner = owner
		}
		GetServer().CreateObjectAt(u, objOwner, pos)
		u.Direction1, u.Direction2 = server.Dir16(dir), server.Dir16(dir)
		u.VelVec = types.Pointf{}
		ud.Target = core.Nox_xxx_spellFlySearchTarget(nil, u, 32, 600, 0, owner)
		ud.Owner = owner
		ud.SpellID = uint32(id)
		u.Pos39 = origin.PosVec
		ud.Deadline = core.Frame() + uint32(core.TickRate())*uint32(core.Rand.Logic.IntClamp(30, 90))
		ud.SeenOwner = core.Frame()
	}
	spellEffectAudio(int32(id), 0, origin)
	return 1
}
