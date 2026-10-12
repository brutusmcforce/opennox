package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"math"
	"unsafe"
)

var presentationChantTree *server.PhonemeLeaf

func presentationPhonemeMark(index int) {
	*effectMapped(1096596 + 4*uintptr(index)) = GetServer().S().Frame()
}
func presentationPhonemeInit() *gui.Window {
	for i := 0; i < 8; i++ {
		name := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 151272+4*uintptr(i))))
		img := Nox_xxx_gLoadImg(name)
		*memmap.PtrPtr(0x5D4594, 1096564+4*uintptr(i)) = unsafe.Pointer(img.C())
		if img == nil {
			return nil
		}
	}
	w := GetClient().Cli().GUI.NewWindowRaw(nil, 64, (int(nox_win_width)-100)/2, (int(nox_win_height)-100)/2, 1, 1, nil)
	w.SetDraw(func(*gui.Window, *gui.WindowData) int { return presentationPhonemeDraw() })
	return w
}
func presentationPhonemeDraw() int {
	for i := 0; i < 8; i++ {
		stamp := effectMapped(1096596 + 4*uintptr(i))
		if *stamp == 0 {
			continue
		}
		img := GetClient().R2().GetBag().AsImage(noxrender.ImageHandle(*memmap.PtrPtr(0x5D4594, 1096564+4*uintptr(i))))
		pos := image.Pt(int(nox_win_width)/2+int(memmap.Int32(0x587000, 151208+8*uintptr(i)))-16, int(nox_win_height)/2+int(memmap.Int32(0x587000, 151212+8*uintptr(i)))-41)
		GetClient().R2().DrawImageAt(img, pos)
		if GetServer().S().Frame()-*stamp > 3 {
			*stamp = 0
		}
	}
	return 1
}
func presentationChantStart(id byte) {
	*memmap.PtrUint8(0x5D4594, 1303504) = id
	*memmap.PtrUint8(0x5D4594, 1303512) = 0
	*effectMapped(1303516) = GetServer().S().Frame()
	presentationChantTree = GetServer().S().Spells.PhonemeTree()
}
func presentationChantClear() { *memmap.PtrUint8(0x5D4594, 1303504) = 0 }
func presentationChantTick() {
	id := memmap.Uint8(0x5D4594, 1303504)
	if id == 0 {
		return
	}
	s := GetServer().S()
	index := memmap.PtrUint8(0x5D4594, 1303512)
	ph := s.Spells.Phoneme(spell.ID(id), int(*index))
	if s.Frame() >= *effectMapped(1303516) {
		audioID := spellLifePhoneme(int32(nox_player_netCode_85319C), int8(ph))
		Nox_xxx_clientPlaySoundSpecial_452D80(sound.ID(audioID), 100)
		presentationPhonemeMark(int(memmap.Uint32(0x587000, 163576+4*uintptr(ph))))
		*effectMapped(1303516) = s.Frame() + 3
		presentationChantTree = presentationChantTree.Next(ph)
		*index++
	}
	if presentationChantTree.Ind == int32(memmap.Uint8(0x5D4594, 1303504)) {
		presentationChantClear()
	}
}
func presentationShields() []*client.Drawable {
	return unsafe.Slice((**client.Drawable)(memmap.PtrOff(0x5D4594, 1217468)), 9)
}
func presentationShieldInit() bool {
	for i, name := range []string{"ReflectiveShieldNW", "ReflectiveShieldN", "ReflectiveShieldNE", "ReflectiveShieldW", "", "ReflectiveShieldE", "ReflectiveShieldSW", "ReflectiveShieldS", "ReflectiveShieldSE"} {
		if name == "" {
			presentationShields()[i] = nil
			continue
		}
		presentationShields()[i] = GetClient().Nox_new_drawable_for_thing(int(effectType(name)))
	}
	for i, dr := range presentationShields() {
		if i == 4 {
			continue
		}
		if dr == nil {
			return false
		}
		dr.ObjFlags |= 0x1000000
	}
	*effectMapped(1217504) = 0
	return true
}
func presentationShieldDestroy() {
	for i, dr := range presentationShields() {
		if dr != nil {
			GetClient().Nox_xxx_spriteDelete_45A4B0(dr)
		}
		presentationShields()[i] = nil
	}
	*effectMapped(1217504) = 0
}
func presentationShieldDraw(vp *noxrender.Viewport, dr *client.Drawable) {
	i := uintptr(dr.AnimDir)
	shield := presentationShields()[i]
	shield.PosVec = dr.PosVec.Add(image.Pt(int(memmap.Int32(0x587000, 161776+8*i)), int(int16(dr.ZVal))+int(memmap.Int32(0x587000, 161780+8*i))))
	client.CallDrawableDrawDiscard(shield.DrawFuncPtr, vp, shield)
}
func presentationTurnUndead(pos *[2]int16) {
	typ := effectMapped(1217508)
	if *typ == 0 {
		*typ = effectType("UndeadKiller")
	}
	for angle := 0; angle < 256; angle += 6 {
		dr := effectSpawn(int(*typ), image.Pt(int(pos[0]), int(pos[1])))
		if dr == nil {
			continue
		}
		*(*uint16)(unsafe.Add(dr.C(), 508)) = uint16(angle)
		dr.Field_117 = math.Float32bits(*memmap.PtrFloat32(0x587000, 194136+8*uintptr(angle)) * 4)
		dr.Field_118 = math.Float32bits(*memmap.PtrFloat32(0x587000, 194140+8*uintptr(angle)) * 4)
		dr.Field_119 = 0
		dr.AnimStart = GetServer().S().Frame()
		dr.Field_81, dr.Field_82 = uint32(int32(pos[0])), uint32(int32(pos[1]))
		dr.Field_115 = drawableUpdateIdentity(updateID_sprite)
		GetClient().Cli().Objs.List5Add(dr)
		GetClient().Cli().Objs.List6Add(dr)
	}
}
func presentationBubble(typ int, pos image.Point, z int16, a, b, c, d, e byte, lifetime int) {
	if *effectMapped(1217512) == 0 {
		for i, name := range []string{"RedBubbleParticle", "WhiteBubbleParticle", "LightBlueBubbleParticle", "OrangeBubbleParticle", "GreenBubbleParticle", "VioletBubbleParticle", "LightVioletBubbleParticle", "YellowBubbleParticle"} {
			*effectMapped(1217512 + 4*uintptr(i)) = effectType(name)
		}
	}
	dr := effectSpawn(typ, pos)
	if dr == nil {
		return
	}
	dr.ZVal = uint16(z)
	*effectWord(dr, 432) = uint32(noxcolor.RGB5551Color(byte(dr.LightColor.R), byte(dr.LightColor.G), byte(dr.LightColor.B)).Color32())
	rgb := [3]byte{200, 200, 255}
	for i, value := range [][3]byte{{255, 128, 128}, {255, 255, 255}, {200, 200, 255}, {255, 100, 50}, {64, 255, 64}, {255, 100, 255}, {255, 200, 255}, {255, 255, 200}} {
		if i != 2 && uint32(typ) == *effectMapped(1217512 + 4*uintptr(i)) {
			rgb = value
			break
		}
	}
	*effectWord(dr, 436) = uint32(noxcolor.RGB5551Color(rgb[0], rgb[1], rgb[2]).Color32())
	copy(unsafe.Slice(effectByte(dr, 440), 7), []byte{a, 1, b, b, d, e, c})
	cl := GetClient().Cli()
	cl.Objs.List6Add(dr)
	cl.Objs.TransparentDecay(dr, lifetime)
	cl.Objs.List34Add(dr)
}
func presentationEquip(op byte, id, mask uint32, mods *[4]byte) {
	s := GetServer().S()
	npc := s.NPCs.ByID(int(id))
	if npc == nil {
		return
	}
	slots := npc.Armor[:]
	aggregate := &npc.ArmorEquip
	if op == 80 || op == 81 {
		slots = npc.Weapon[:]
		aggregate = &npc.WeaponEquip
	}
	for i := range slots {
		if slots[i].Field0 != 0 {
			continue
		}
		slots[i].Field0 = mask
		*aggregate |= mask
		for j, m := range mods {
			slots[i].Field4[j] = s.Modif.Nox_xxx_modifGetDescById413330(int(m)).C()
		}
		return
	}
}
func presentationBookReward(kind, id, auto int) {
	seq := uint32(GetClient().GetInputSeq())
	if noxflags.HasGame(noxflags.GameModeCoop) && seq-*effectMapped(1217504) < 2 {
		return
	}
	*effectMapped(1217504) = seq
	effect := 2
	if kind == 2 {
		effect = 0
	} else if kind == 3 {
		effect = 3
	}
	pos := image.Pt(5, int(nox_win_height)/3)
	bookSetForward(uintptr(kind), id, pos)
	for _, p := range [][6]int{{0, 0, 271, 166, 1, 1}, {0, 0, 135, 166, 2, 1}, {0, 166, 135, 166, 2, 1}, {271, 0, 271, 166, 1, 2}, {135, 0, 135, 166, 2, 2}, {135, 166, 135, 166, 2, 2}} {
		effectScreenParticles(effect, pos.X+p[0], pos.Y+p[1], p[2], p[3], p[4], p[5])
	}
	if kind != 4 && auto == 1 {
		bookAdd(kind, id)
	}
}
