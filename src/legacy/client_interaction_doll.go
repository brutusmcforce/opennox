package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"unsafe"
)

func interactionDollImage(handle uint32, pos image.Point) {
	r := GetClient().R2()
	r.DrawImageAt(r.GetBag().AsImage(noxrender.ImageHandle(unsafe.Pointer(uintptr(handle)))), pos)
}
func interactionDollLoad() int {
	load := func(nameOffset uintptr) uint32 {
		name := alloc.GoString(*(**byte)(memmap.PtrOff(0x587000, nameOffset)))
		return uint32(uintptr(Nox_xxx_gLoadImg(name).C()))
	}
	for gender := uintptr(0); gender < 2; gender++ {
		*memmap.PtrUint32(0x973A20, 16+4*gender) = load(180960 + 4*gender)
		*memmap.PtrUint32(0x973A20, 24+4*gender) = load(180968 + 4*gender)
		for i := uintptr(0); i < 26; i++ {
			off := 104*gender + 4*i
			*memmap.PtrUint32(0x973A20, 32+off) = load(180976 + off)
		}
		for i := uintptr(0); i < 27; i++ {
			off := 108*gender + 4*i
			*memmap.PtrUint32(0x973A20, 256+off) = load(181184 + off)
		}
	}
	*memmap.PtrUint32(0x5D4594, 1319052) = uint32(uintptr(Nox_xxx_gLoadImg("MaleMedievalCloakTop").C()))
	return 1
}
func interactionDollLayer(typ uint32, pos image.Point, table *uint32, index uint32, overlay bool) int16 {
	dr := uiInventoryEquippedType(typ)
	if dr == nil {
		return 0
	}
	mods := &GetServer().S().Modif
	var def *server.Modifier
	if uint32(dr.ObjClass)&0x2000000 != 0 {
		def = mods.Nox_xxx_equipClothFindDefByTT413270(int(dr.TypeIDVal))
	} else {
		def = mods.Nox_xxx_getProjectileClassById413250(int(dr.TypeIDVal))
	}
	if def != nil {
		data := GetClient().R2().Data()
		for i := 1; i < 7; i++ {
			c := def.Colors12[i]
			data.SetMaterialRGB(i, int(c.R), int(c.G), int(c.B))
		}
		indices := def.ColorIndexes()
		for i := 0; i < 4; i++ {
			if mod := uiInventoryItemModifier(dr, i); mod != nil {
				c := mod.Color24
				data.SetMaterialRGB(int(indices[i]), int(c.R), int(c.G), int(c.B))
			}
		}
	}
	handle := memmap.Uint32(0x5D4594, 1319052)
	if !overlay {
		handle = *(*uint32)(unsafe.Add(unsafe.Pointer(table), uintptr(index)*4))
	}
	interactionDollImage(handle, pos)
	return int16(uintptr(dr.C()))
}
func interactionDollDraw(origin image.Point) int16 {
	pos := origin.Add(image.Pt(11, 15))
	uiMeterSetColor(uint32(nox_color_black_2650656))
	nox_client_drawRectFilledOpaque_49CE30(pos.X, pos.Y, 200, 200)
	active := memmap.Uint32(0x852978, 8)
	result := int16(active)
	if active == 0 {
		return result
	}
	p := Get_dword_8531A0_2576()
	if p == nil {
		return result
	}
	data := GetClient().R2().Data()
	for i, c := range []uint32{p.Colors.Hair, p.Colors.Goatee, p.Colors.UnkColor, p.Colors.Beard, p.Colors.Mustache, p.Colors.Skin} {
		data.SetMaterial(i+1, noxcolor.RGBA5551(c))
	}
	gender := uintptr(*(*byte)(unsafe.Add(unsafe.Pointer(p), 2252)))
	base := uintptr(16)
	if p.Colors.Skin == p.Colors.Hair {
		base = 24
	}
	interactionDollImage(memmap.Uint32(0x973A20, base+4*gender), pos)
	armor := memmap.PtrUint32(0x973A20, 32+104*gender)
	drawArmor := func(i uint32, overlay bool) {
		typ := GetServer().S().Armor.Sub_415CD0(1 << i)
		interactionDollLayer(uint32(typ), pos, armor, i, overlay)
	}
	for i := uint32(0); i < 24; i++ {
		if p.ArmorEquip&(1<<i) != 0 {
			drawArmor(i, false)
		}
	}
	if p.ArmorEquip&2 != 0 {
		drawArmor(1, true)
	}
	for i := uint32(24); i < 26; i++ {
		if p.ArmorEquip&(1<<i) != 0 {
			drawArmor(i, false)
		}
	}
	weapons := memmap.PtrUint32(0x973A20, 256+108*gender)
	for i := uint32(0); i < 27; i++ {
		result = int16(p.WeaponEquip)
		if p.WeaponEquip&(1<<i) != 0 {
			typ := GetServer().S().Weapons.Sub_415840(1 << i)
			result = interactionDollLayer(uint32(typ), pos, weapons, i, false)
		}
	}
	return result
}

func sub_4BF7E0(point *uint32) int16 {
	p := (*[2]int32)(unsafe.Pointer(point))
	return int16(interactionDollDraw(image.Pt(int(p[0]), int(p[1]))))
}

func sub_4BF9F0(mask, typ, x, y, table, index, overlay int32) int16 {
	return int16(interactionDollLayer(uint32(typ), image.Pt(int(x), int(y)), (*uint32)(unsafe.Pointer(uintptr(uint32(table)))), uint32(index), overlay != 0))
}

func sub_4BFAD0() int32 { return int32(interactionDollLoad()) }
