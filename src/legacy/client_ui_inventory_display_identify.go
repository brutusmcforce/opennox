package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/server"
)

func uiInventoryWideAt(offset uintptr) string {
	return alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, offset)))
}
func uiInventoryListText(w *gui.Window, text string) {
	p := alloc.InternCString16(text)
	w.Func94(&gui.RawEvent{Event: 16397, Arg1: uintptr(unsafe.Pointer(p)), Arg2: uintptr(uint32(0xffffffff))})
}
func uiInventoryIdentify(pos image.Point) uint32 {
	mouse := GetClient().GetMousePos()
	r := GetClient().R2()
	r.Data().SetTextColor(noxcolor.RGBA5551(nox_color_white_2523948))
	uiMeterSetColor(uint32(nox_color_black_2650656))
	nox_client_drawRectFilledOpaque_49CE30(pos.X+11, pos.Y+15, 200, 200)
	rel := mouse.Sub(uiWindowPosition(uiInventoryMainWindow()))
	cursor := 7
	if uiInventoryHitRect(rel, 136352) || uiInventoryHitRect(rel, 136368) {
		if !uiInventoryHitRect(rel, 136368) || (rel.Y-13)/50 == 1 {
			cursor = 6
		}
	} else if uiInventoryHitRect(rel, 136336) {
		cursor = 0
	} else if uiShopActive() != 0 && uiShopMode() == 2 {
		if uiShopInside(rel) {
			cursor = 6
		}
	}
	nox_client_setCursorType_477610(cursor)
	dr := uiInventorySelectedItem()
	panel := uiInventoryIdentifyWindow()
	header := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1063124)), 256)
	if dr == nil {
		if dword_5d4594_1063120 != 0 {
			dword_5d4594_1063120 = 0
			alloc.StrCopyZero16(header, uiInventoryText("thing.db:IdentifyDescription"))
			return uint32(gui.EventRespInt(panel.ChildByID(9156).Func94(&gui.RawEvent{Event: 16399})))
		}
		return 0
	}
	if dword_5d4594_1063120 == dword_5d4594_1063116 {
		return uiInventoryPointer(dr.C())
	}
	dword_5d4594_1063120 = dword_5d4594_1063116
	// Localization entries can have random variants. Preserve the original
	// lookup order and install the prefix before asking for the item name.
	prefix := uiInventoryText("IdentifyItem") + " "
	alloc.StrCopyZero16(header, prefix)
	name := alloc.GoString16(uiItemTooltip(dr))
	if name == uiInventoryWideAt(1063652) {
		dword_5d4594_1063120 = 0
	}
	alloc.StrCopyZero16(header, prefix+name)
	panel.ChildByID(9151).Func94(&gui.StaticTextSetText{Str: alloc.GoString16(&header[0])})
	list := panel.ChildByID(9156)
	list.Func94(&gui.RawEvent{Event: 16399})
	add := func(text string) { uiInventoryListText(list, text) }
	current, maximum := dr.Field_73_1, dr.Field_73_2
	var text string
	if noxflags.HasGame(2048) {
		if maximum != 0 {
			var a, b float32
			uiInventoryScaledDurability(dr, &a, &b)
			text = uiInventoryFormatInts("IdentifyDurability", int(a), int(b))
		} else {
			text = uiInventoryText("IdentifyDurabilityIndestructable")
		}
	} else {
		keys := [5]string{"IdentifyDurabilityNoDamage", "IdentifyDurabilitySlight", "IdentifyDurabilityModerate", "IdentifyDurabilitySevere", "IdentifyDurabilityIndestructable"}
		text = uiInventoryFormatLiteral(uiInventoryText(keys[durabilityBand(current, maximum)]))
	}
	add(text)
	add(uiInventoryWideAt(1063656))
	add(uiInventoryFormatInts("IdentifyWeight", int(dr.Weight), 0))
	add(uiInventoryWideAt(1063660))
	class, sub := uint32(dr.Class()), uint32(dr.SubClass())
	modifiers := &GetServer().S().Modif
	if class&0x2000000 != 0 {
		def := modifiers.Nox_xxx_equipClothFindDefByTT413270(int(dr.TypeIDVal))
		scale := float64(1)
		if m := uiInventoryItemModifier(dr, 0); m != nil && m.Defend76.Fnc == modifierKey(modifierIDArmorMultiplierEffect) {
			scale = float64(m.Defend76.Valf)
		}
		value := int(float32(scale*float64(def.DamageCoeffOrArmor64)*1000.0 + 0.5))
		if sub&2 != 0 {
			text = uiInventoryFormatLiteral(uiInventoryText("ArmorValueLabelNA"))
		} else {
			text = uiInventoryFormatInts("ArmorValueLabel", value, 0)
		}
		add(text)
		add(uiInventoryWideAt(1063664))
	} else if class&0x1001000 != 0 {
		p := uiMeterPlayer()
		arrow := memmap.PtrUint32(0x5D4594, 1063644)
		bolt := memmap.PtrUint32(0x5D4594, 1063648)
		if *arrow == 0 {
			*arrow = uint32(GetClient().Cli().Things.IndByID("ArcherArrow"))
			*bolt = uint32(GetClient().Cli().Things.IndByID("ArcherBolt"))
		}
		selected := false
		var def *server.Modifier
		if p != nil && sub&2 != 0 {
			mask := p.WeaponEquip
			if mask&4 != 0 {
				selected = true
				def = modifiers.Nox_xxx_getProjectileClassById413250(int(*arrow))
			} else if mask&8 != 0 {
				selected = true
				def = modifiers.Nox_xxx_getProjectileClassById413250(int(*bolt))
			}
		}
		if def == nil {
			def = modifiers.Nox_xxx_getProjectileClassById413250(int(dr.TypeIDVal))
		}
		electric, fire := float32(uiInventoryElementValue(dr, false)), float32(uiInventoryElementValue(dr, true))
		scale := float32(1)
		if m := uiInventoryItemModifier(dr, 0); m != nil && m.Attack40.Fnc == modifierKey(modifierIDDamageMultiplierEffect) {
			scale = m.Attack40.Valf
		}
		baseDamage := controlBoltDamage(*(*int32)(unsafe.Pointer(&p.InfoBytes[54])), def.C())
		total := float32(baseDamage*float64(scale) + float64(electric) + float64(fire))
		base := float64(def.DamageMin72)
		if def.TypeInd == *bolt && noxflags.HasGame(2048) {
			base = GetServer().S().Balance.Float("BoltSoloDamageMin")
		}
		baseValue := float32(base * float64(scale))
		strength := float32(float64(total) - float64(baseValue) - float64(fire) - float64(electric))
		if strength < 0 {
			total = float32(float64(total) - float64(strength))
			strength = 0
		}
		if sub&12 != 0 {
			text = uiInventoryFormatLiteral(uiInventoryText("WeaponDamageLabelNA"))
		} else if sub&2 == 0 || selected {
			text = uiInventoryFormatFloat("WeaponDamageLabel", float64(total))
		} else {
			text = uiInventoryFormatFloat("WeaponDamageLabelUnknownPlus", float64(total))
		}
		add(text)
		add("  " + uiInventoryFormatFloat("BaseDamageLabel", float64(baseValue)))
		add("  " + uiInventoryFormatFloat("StrengthDamageLabel", float64(strength)))
		if fire > 0 {
			add("  " + uiInventoryFormatFloat("FireDamageLabel", float64(fire)))
		}
		if electric > 0 {
			add("  " + uiInventoryFormatFloat("ElectricalDamageLabel", float64(electric)))
		}
		add(uiInventoryWideAt(1063668))
	}
	if class&0x13001000 != 0 {
		var def *server.Modifier
		if class&0x11001000 != 0 {
			def = modifiers.Nox_xxx_getProjectileClassById413250(int(dr.TypeIDVal))
		} else {
			def = modifiers.Nox_xxx_equipClothFindDefByTT413270(int(dr.TypeIDVal))
		}
		if def == nil {
			add(uiInventoryFormatLiteral(uiInventoryText("IdentifySpecialAttributes")))
			add(uiInventoryFormatLiteral(uiInventoryText("IdentifyUnknown")))
		} else if class&0x10000000 == 0 {
			first := true
			for i := 2; i < 4; i++ {
				if m := uiInventoryItemModifier(dr, i); m != nil {
					desc := m.IdentDesc16
					if desc != nil {
						if first {
							add(uiInventoryFormatLiteral(uiInventoryText("IdentifySpecialAttributes")))
							first = false
						}
						add("  " + uiInventoryFormatLiteral(alloc.GoString16(desc)))
					}
				}
			}
			if !first {
				add(uiInventoryWideAt(1063672))
			}
		}
	}
	typ := GetClient().Cli().Things.TypeByInd(int(dr.TypeIDVal))
	if typ != nil && typ.Desc != nil {
		add(alloc.GoString16(typ.Desc))
	}
	icon := panel.ChildByID(9155)
	if icon == nil {
		return uint32(0xfffffffe)
	}
	handle := uint32(0)
	if typ != nil {
		handle = typ.PrettyImage
	}
	uiMeterSetIcon(icon, handle)
	return 0
}
