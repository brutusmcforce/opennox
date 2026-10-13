package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/server"
)

func itemXferObelisk(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 61)
	if !ok {
		return 0
	}
	if v >= 61 {
		r.raw(u.UpdateData, 4)
		// The former mana-derived callback is empty and has no observable effect.
		present := byte(0)
		if noxflags.HasGame(2048) {
			c := GetClient().Cli()
			dr := c.Objs.ByNetCodeStatic(int(u.Extent))
			if dr != nil {
				for it := c.Objs.FirstMinimapList(); it != nil; it = it.Field_102 {
					if it == dr {
						present = 1
						break
					}
				}
			}
		}
		r.byte(present)
	}
	return objectXferFinish(r, u, v, saved)
}
func itemXferGenerator(u *server.Object) int {
	p := u.UpdateDataMonsterGen()
	r, v, saved, ok := itemXferPositiveStart(u, 63)
	if !ok {
		return 0
	}
	for i, n := 0, int(r.byte(3)); i < n; i++ {
		r.raw(unsafe.Pointer(&p.Rates[i]), 1)
	}
	r.raw(unsafe.Pointer(&p.Current), 1)
	r.raw(unsafe.Pointer(&p.Max), 1)
	r.raw(unsafe.Pointer(&p.Frame88), 4)
	for _, pair := range [][2]int{{48, server.ScriptNameGen52}, {56, server.ScriptNameGen60}, {72, server.ScriptNameGen76}, {64, server.ScriptNameGen68}} {
		objectXferScript(unsafe.Pointer(p.ScriptAt(pair[0])), objectXferScriptName(u, pair[1]))
	}
	if r.read() {
		for row, rows := 0, int(r.byte(0)); row < rows; row++ {
			for col, n := 0, int(r.byte(0)); col < n; col++ {
				length := int(r.byte(0))
				var name [256]byte
				r.raw(unsafe.Pointer(&name[0]), length)
				ch := GetServer().S().NewObjectByTypeID(alloc.GoString(&name[0]))
				if ch == nil {
					return 0
				}
				r.short(0)
				r.cf.ReadAlignedU32()
				if ch.CallXfer(nil) != nil {
					GetServer().S().Objs.FreeObject(ch)
					return 0
				}
				p.Field0[row*4+col] = ch
			}
		}
	} else {
		r.byte(3)
		for row := 0; row < 3; row++ {
			children := p.Field0[4*row : 4*row+4]
			count := byte(0)
			for _, ch := range children {
				if ch != nil {
					count++
				}
			}
			r.byte(count)
			for _, ch := range children {
				if ch == nil {
					continue
				}
				name := GetServer().S().Types.ByInd(int(ch.TypeInd)).ID()
				n := int(r.byte(byte(len(name))))
				r.cf.ReadWrite([]byte(name)[:n])
				Nox_xxx_xfer_saveObj51DF90(r.cf, ch)
			}
		}
	}
	if v >= 62 {
		for i, n := 0, int(r.byte(3)); i < n; i++ {
			r.raw(unsafe.Pointer(&p.Field83[i]), 1)
		}
	}
	if v >= 63 {
		r.raw(unsafe.Pointer(&p.Flags92), 4)
	}
	return objectXferFinish(r, u, v, saved)
}
func itemXferRewardMarker(u *server.Object) int {
	p := u.InitDataRewardMarker()
	r, v, saved, ok := itemXferPositiveStart(u, 63)
	if !ok {
		return 0
	}
	r.raw(unsafe.Pointer(&p.Categories), 4)
	r.raw(unsafe.Pointer(&p.Explicit), 4)
	for group, mask := range [][]byte{p.Spells[:], p.Abilities[:], p.Guides[:]} {
		count := uint16(0)
		for _, b := range mask {
			if b == 1 {
				count++
			}
		}
		count = r.short(count)
		if r.read() {
			for i := 0; i < int(count); i++ {
				n := int(r.byte(0))
				var name [256]byte
				r.raw(unsafe.Pointer(&name[0]), n)
				id := 0
				switch group {
				case 0:
					id = int(spell.ParseID(alloc.GoString(&name[0])))
				case 1:
					id = int(bookAbilityID(alloc.GoString(&name[0])))
				case 2:
					id = int(bookGuideID(alloc.GoString(&name[0])))
				}
				if id == 0 {
					return 0
				}
				mask[id] = 1
			}
		} else {
			// Historical counts include only byte==1, while names include all nonzero bytes.
			for id, b := range mask {
				if b == 0 {
					continue
				}
				name := ""
				switch group {
				case 0:
					name = spell.ID(id).String()
				case 1:
					name = server.Ability(id).String()
				case 2:
					name = alloc.GoString(bookGuideName(int32(id)))
				}
				n := int(r.byte(byte(len(name))))
				r.cf.ReadWrite([]byte(name)[:n])
			}
		}
	}
	for _, f := range []*uint32{&p.Field196, &p.Field192, &p.Field200, &p.Field204, &p.Field208} {
		r.raw(unsafe.Pointer(f), 4)
	}
	if v >= 62 {
		r.raw(unsafe.Pointer(&p.Chance), 4)
	}
	if v >= 63 {
		r.raw(unsafe.Pointer(&p.Flags), 1)
	}
	return objectXferFinish(r, u, v, saved)
}
