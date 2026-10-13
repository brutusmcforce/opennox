package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func sustainedTeleportStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	delay := uint32(1)
	if sustainedGame(2048) {
		delay = uint32(sustainedTableInt("TeleportDelay", d.Level))
	}
	d.Frame68 = sustainedFrame() + delay
	return 0
}
func sustainedBlinkStart(p unsafe.Pointer) uint32 {
	if sustainedGame(4096) && (*server.DurSpell)(p).Flag20 == 1 {
		return 1
	}
	return sustainedTeleportStart(p)
}
func sustainedTeleportWake(u *server.Object, from, to *types.Pointf) uint32 {
	id := stateType(2487916, "TeleportWake")
	wake := spellEffectNew(id)
	if wake == nil {
		return 0
	}
	wake.CollideDataTeleportWake().Dest = *to
	spellEffectCreate(wake, u, *from)
	ret := sustainedFrame() + sustainedFPS()
	wake.Field34 = ret
	return ret
}
func sustainedTeleportSound(d *server.DurSpell, phase int32, u *server.Object, hidden bool) {
	mode, code := int32(0), int32(0)
	if hidden {
		if u.ObjClass&4 == 0 {
			return
		}
		mode = 2
		code = int32(u.NetCode)
	}
	id := int32(nox_xxx_spellGetAud44_424800(int(int32(d.Spell)), int(int32(phase))))
	nox_xxx_aud_501960(int32(id), asObjectC(u), int(mode), int32(code))
}
func sustainedBlinkTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Target48
	if u == nil || u.ObjFlags&0x8020 != 0 {
		return 1
	}
	if d.Frame68-1 != sustainedFrame() {
		return 0
	}
	if spellLifeHasBuff(u, 14) {
		sustainedAudio(231, u)
		return 1
	}
	var dest types.Pointf
	if sustainedGame(4096) && u.ObjClass&4 != 0 {
		if start := layout.As[server.PlayerUpdateData](u.UpdateData).Field77; start != nil {
			inventoryRandomPlacement(60, &start.PosVec, &dest)
		} else {
			controlFindStart(&dest, u)
		}
	} else if GetServer().S().Nox_xxx_waypoint_579F00(&dest, u) == 0 {
		controlFindStart(&dest, u)
	}
	sustainedTeleportWake(u, &u.PosVec, &dest)
	sustainedFX(137, u.PosVec)
	spellEffectAudio(int32(d.Spell), 0, u)
	if !spellLifeHasBuff(u, 0) {
		sustainedFX(137, u.PosVec)
		sustainedFX(137, dest)
	}
	stateTeleport(u, &dest)
	hidden := spellLifeHasBuff(u, 0)
	if !hidden {
		sustainedFX(137, u.PosVec)
	}
	sustainedTeleportSound(d, 0, u, hidden)
	spellEffectAlert(d.Caster16, u)
	return 1
}
func sustainedGlyphStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u := d.Obj12
	if u.ObjClass&4 != 0 && (u.UpdateData == nil || u.UpdateDataPlayer().Field29[d.Spell-122] == nil) {
		return 1
	}
	return sustainedTeleportStart(p)
}
func sustainedGlyphUse(ud unsafe.Pointer, index int, mark bool) {
	pud := layout.As[server.PlayerUpdateData](ud)
	glyph := pud.Field29[index]
	if mark {
		glyph.Field34 = sustainedFrame()
	}
	charge := &(*[4]byte)(unsafe.Pointer(&pud.Field39))[index]
	*charge--
	if *charge == 0 {
		sustainedFX(129, glyph.PosVec)
		sustainedDelete(glyph)
		pud.Field29[index] = nil
	}
}
func sustainedGlyphTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	owner, u := d.Obj12, d.Target48
	if owner == nil || owner.ObjFlags&0x20 != 0 || u == nil || u.ObjFlags&0x8020 != 0 || u != owner && d.Flag20 == 0 {
		return 1
	}
	if d.Frame68-1 != sustainedFrame() {
		return 0
	}
	if spellLifeHasBuff(u, 14) {
		sustainedAudio(231, u)
		return 1
	}
	if owner.ObjClass&4 != 0 {
		ud := owner.UpdateData
		index := int(d.Spell) - 122
		pud := layout.As[server.PlayerUpdateData](ud)
		glyph := pud.Field29[index]
		if glyph == nil {
			return 1
		}
		old := u.PosVec
		sustainedTeleportWake(u, &u.PosVec, &glyph.PosVec)
		spellEffectAudio(int32(d.Spell), 1, u)
		stateTeleport(u, &glyph.PosVec)
		hidden := spellLifeHasBuff(u, 0)
		if !hidden {
			sustainedFX(137, old)
			sustainedFX(137, u.PosVec)
		}
		sustainedTeleportSound(d, 1, u, hidden)
		sustainedGlyphUse(ud, index, true)
	}
	if d.Caster16 != nil && d.Caster16.ObjFlags&0x20 == 0 {
		spellEffectAlert(d.Caster16, u)
	}
	return 1
}
func sustainedRandomGlyphTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	owner, u := d.Caster16, d.Target48
	if owner == nil || owner.ObjFlags&0x8020 != 0 || u == nil || u.ObjFlags&0x8020 != 0 {
		return 1
	}
	if d.Frame68-1 != sustainedFrame() {
		return 0
	}
	if spellLifeHasBuff(u, 14) {
		sustainedAudio(231, u)
		return 1
	}
	if owner.ObjClass&4 != 0 {
		ud := owner.UpdateData
		count := 0
		for i := 0; i < 4; i++ {
			if layout.As[server.PlayerUpdateData](ud).Field29[i] != nil {
				count++
			}
		}
		if count == 0 {
			return 1
		}
		old := u.PosVec
		var glyph *server.Object
		index := 0
		for glyph == nil {
			index = GetServer().S().Rand.Logic.IntClamp(0, 3)
			glyph = layout.As[server.PlayerUpdateData](ud).Field29[index]
		}
		sustainedTeleportWake(u, &u.PosVec, &glyph.PosVec)
		stateTeleport(u, &glyph.PosVec)
		if !spellLifeHasBuff(u, 0) {
			sustainedFX(137, old)
			sustainedFX(137, u.PosVec)
		}
		spellEffectAudio(int32(d.Spell), 0, u)
		sustainedGlyphUse(ud, index, false)
	}
	spellEffectAlert(owner, u)
	return 1
}
func sustainedTeleportToPointTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u, owner := d.Target48, d.Caster16
	if u == nil || owner == nil || u.ObjFlags&0x8020 != 0 || owner.ObjFlags&0x8020 != 0 {
		return 1
	}
	if d.Frame68-1 != sustainedFrame() {
		return 0
	}
	if spellLifeHasBuff(u, 14) {
		sustainedAudio(231, u)
		return 1
	}
	sustainedFX(137, u.PosVec)
	spellEffectAudio(int32(d.Spell), 0, u)
	sustainedTeleportWake(u, &u.PosVec, &d.Pos2)
	stateTeleport(u, &d.Pos2)
	hidden := spellLifeHasBuff(u, 0)
	if !hidden {
		sustainedFX(137, u.PosVec)
	}
	sustainedTeleportSound(d, 0, u, hidden)
	spellEffectAlert(owner, u)
	return 1
}
func sustainedSwapStart(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u, t := d.Caster16, d.Target48
	if u == nil || d.Flags88&0x20 != 0 || t == nil || t.ObjFlags&0x8020 != 0 || u == t || t.ObjClass&2 != 0 && t.ObjSubClass&0x4000 != 0 {
		return 1
	}
	return sustainedTeleportStart(p)
}
func sustainedSwapTick(p unsafe.Pointer) uint32 {
	d := (*server.DurSpell)(p)
	u, t := d.Caster16, d.Target48
	if t == nil || u == nil || t.ObjFlags&0x8020 != 0 || u.ObjFlags&0x8020 != 0 {
		return 1
	}
	if d.Frame68-1 != sustainedFrame() {
		return 0
	}
	if spellLifeHasBuff(t, 14) || spellLifeHasBuff(u, 14) {
		sustainedAudio(231, t)
		sustainedAudio(231, u)
		return 1
	}
	if d.Flag20 == 0 {
		if !sustainedInteract(u, t) {
			resourcePriority(u, "ExecDur.c:NeedClearLOSForSwap")
			return 1
		}
		if u.ObjClass&4 != 0 {
			pl := controlPlayer(u)
			x, y := float64(pl.Field10), float64(pl.Field12)
			rect := [4]float32{float32(float64(u.PosVec.X) - x), float32(float64(u.PosVec.Y) - y), float32(x + float64(u.PosVec.X)), float32(y + float64(u.PosVec.Y))}
			if int32(geometryRectFloat((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&t.PosVec))), (*[4]float32)(unsafe.Pointer(unsafe.Pointer(&rect))))) == 0 {
				resourcePriority(u, "ExecDur.c:NeedClearLOSForSwap")
				return 1
			}
		}
	}
	old := t.PosVec
	stateTeleport(t, &u.PosVec)
	stateTeleport(u, &old)
	if !spellLifeHasBuff(t, 0) && !spellLifeHasBuff(u, 0) {
		sustainedFX(137, u.PosVec)
		sustainedFX(137, t.PosVec)
	}
	spellEffectAudio(int32(d.Spell), 0, t)
	spellEffectAudio(int32(d.Spell), 0, u)
	spellEffectAlert(u, t)
	return 1
}
