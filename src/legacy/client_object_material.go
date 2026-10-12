package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"unsafe"
)

func objectTeamColor(arg unsafe.Pointer) int {
	if arg == nil {
		return 0
	}
	name := alloc.GoString(*(**byte)(arg))
	for off := uintptr(177488); ; off += 8 {
		p := (*byte)(memmap.LoadPtr(0x587000, off))
		if p == nil {
			return 0
		}
		if name == alloc.GoString(p) {
			return int(memmap.Uint32(0x587000, off+4))
		}
	}
}
func objectDrawableTeamColor(dr *client.Drawable) int {
	if *effectWord(dr, int(layoutClientDrawable.Off(112)) /* ObjClass */)&0x10000000 == 0 {
		return 0
	}
	return objectTeamColor(*(*unsafe.Pointer)(unsafe.Pointer(effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */))))
}
func objectTeamByColor(color int) *server.Team {
	teams := &GetServer().S().Teams
	for t := teams.First(); t != nil; t = teams.Next(t) {
		if int(t.ColorInd) == color {
			return t
		}
	}
	return nil
}
func objectBaseMaterials(typ int) uint32 {
	def := GetServer().S().Modif.Nox_xxx_getProjectileClassById413250(typ)
	if def == nil {
		return 0
	}
	objectSetBaseMaterials(def)
	// Preserve the retained C/UI entry point's pointer bits with its low byte
	// replaced by the final base material's red component.
	return memmap.PtrWord(def.C())&0xffffff00 | uint32(def.Colors12[6].R)
}
func objectSetBaseMaterials(def *server.Modifier) {
	d := GetClient().R2().Data()
	for i := 1; i <= 6; i++ {
		cl := def.Colors12[i]
		d.SetMaterialRGB(i, int(cl.R), int(cl.G), int(cl.B))
	}
}
func objectEquipmentMaterials(dr *client.Drawable, armor bool) {
	modif := &GetServer().S().Modif
	var def *server.Modifier
	if armor {
		def = modif.Nox_xxx_equipClothFindDefByTT413270(int(dr.TypeIDVal))
	} else {
		def = modif.Nox_xxx_getProjectileClassById413250(int(dr.TypeIDVal))
	}
	if def == nil {
		return
	}
	objectSetBaseMaterials(def)
	for i, slot := range def.ColorIndexes() {
		p := *drawableModifier(dr, i)
		if p == nil {
			continue
		}
		cl := (*server.ModifierEff)(p).Color24
		GetClient().R2().Data().SetMaterialRGB(int(slot), int(cl.R), int(cl.G), int(cl.B))
	}
}
func objectEquipmentDraw(vp *noxrender.Viewport, dr *client.Drawable, armor, animated bool) int {
	objectEquipmentMaterials(dr, armor)
	if animated {
		return spriteAnimateDraw(vp, dr)
	}
	return spriteStaticDraw(vp, dr)
}
func objectFlagDraw(vp *noxrender.Viewport, dr *client.Drawable) int {
	objectEquipmentDraw(vp, dr, false, true)
	if !noxflags.HasGame(noxflags.GameFlag(128)) || *effectWord(dr, int(layoutClientDrawable.Off(120)) /* ObjFlags */)&0x1000000 == 0 {
		return 1
	}
	team := objectTeamByColor(objectDrawableTeamColor(dr))
	if team == nil {
		return 1
	}
	p := vp.ToScreenPos(dr.PosVec)
	p.Y -= int(int16(dr.ZVal)) + int(int64(*(*float32)(unsafe.Add(dr.C(), int(layoutClientDrawable.Off(100)) /* ZSizeMax */))))
	r := GetClient().R2()
	r.Data().SetTextColor(noxcolor.RGBA5551(dword_8531A0_2572))
	name := team.Name()
	size := r.GetStringSizeWrapped(nil, name, 0)
	r.DrawString(nil, name, image.Pt(p.X-size.X/2, p.Y))
	return 1
}
