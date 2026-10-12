package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func spellEffectPortal(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	return spellEffectPortalSlot(id, b, c, -1)
}
func spellEffectNamedPortal(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	return spellEffectPortalSlot(id, b, c, id-46)
}
func spellEffectPortalSlot(id int32, source, origin *server.Object, slot int32) int32 {
	typ := spellEffectGlyphType()
	pos := source.PosVec
	if origin != nil && uint32(origin.TypeInd) == typ {
		pos = origin.PosVec
	}
	if source.ObjClass&4 == 0 {
		return 1
	}
	ud := source.UpdateData
	if slot < 0 {
		slot = 0
		for slot < 4 && spellEffectObject(ud, 116+int(slot)*4) != nil {
			slot++
		}
		if slot == 4 {
			oldest := GetServer().S().Frame()
			slot = int32(memmap.PtrWord(source.CObj()))
			for i := int32(0); i < 4; i++ {
				u := spellEffectObject(ud, 116+int(i)*4)
				if u.Field34 < oldest {
					oldest = u.Field34
					slot = i
				}
			}
		}
	}
	off := 116 + int(slot)*4
	if u := spellEffectObject(ud, off); u != nil {
		Nox_xxx_unitMove_4E7010(u, pos)
		u.Field34 = GetServer().S().Frame()
	} else {
		names := [4]string{"TeleportGlyph1", "TeleportGlyph2", "TeleportGlyph3", "TeleportGlyph4"}
		u := GetServer().S().NewObjectByTypeID(names[slot])
		*controlPtr(ud, off) = u.CObj()
		if u == nil {
			spellEffectAudio(id, 0, source)
			return 1
		}
		spellEffectCreate(u, source, pos)
	}
	*controlByte(ud, 156+int(slot)) = 3
	spellEffectAudio(id, 0, source)
	return 1
}
func spellEffectDetonateGlyph(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	var nearest *server.Object
	distance := float32(1e8)
	for u := GetServer().S().Objs.List; u != nil; u = u.ObjNext {
		if u.HasOwner(b) && GetServer().S().Types.ByInd(int(u.TypeInd)).ID() == "Glyph" {
			dx := float64(b.PosVec.X) - float64(u.PosVec.X)
			dy := float64(b.PosVec.Y) - float64(u.PosVec.Y)
			d := dy*dy + dx*dx
			if d < float64(distance) {
				distance = float32(d)
				nearest = u
			}
		}
	}
	if nearest == nil {
		return 0
	}
	spellEffectAudio(id, 0, b)
	Nox_xxx_dieGlyph_54DF30(nearest)
	return 1
}
