package legacy

import (
	"encoding/binary"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/console"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"unsafe"
)

var combatFeedRows, combatFeedWrite, combatFeedRead uint32

func combatFeedRing() *[100][6]uint32 { return (*[100][6]uint32)(memmap.PtrOff(0x5D4594, 1201428)) }
func combatFeedInit() {
	combatFeedRead, combatFeedWrite = 0, 0
	for i, name := range []string{"ArcherBolt", "ArcherArrow", "Bow", "CrossBow"} {
		p := memmap.PtrUint32(0x5D4594, 1203844+uintptr(i)*4)
		if *p == 0 {
			*p = uint32(GetClient().Cli().Things.IndByID(name))
		}
	}
	combatFeedDefaultIcon = Nox_xxx_spellIcon_424A90(15)
}
func combatFeedAdd(raw unsafe.Pointer) {
	p := unsafe.Slice((*byte)(raw), 11)
	index := combatFeedWrite
	if (combatFeedWrite+1)%100 == combatFeedRead {
		combatFeedRead = (combatFeedRead + 1) % 100
	}
	cause := binary.LittleEndian.Uint16(p[8:])
	if p[10] == 1 {
		if uint32(cause) == memmap.Uint32(0x5D4594, 1203844) {
			cause = memmap.Uint16(0x5D4594, 1203856)
		} else if uint32(cause) == memmap.Uint32(0x5D4594, 1203848) {
			cause = memmap.Uint16(0x5D4594, 1203852)
		}
		binary.LittleEndian.PutUint16(p[8:], cause)
	}
	combatFeedRing()[index] = [6]uint32{uint32(binary.LittleEndian.Uint16(p[2:])), uint32(binary.LittleEndian.Uint16(p[4:])), uint32(binary.LittleEndian.Uint16(p[6:])), uint32(cause), uint32(p[10]), gameFrame()}
	combatFeedWrite = (index + 1) % 100
	combatFeedConsole(p)
}
func combatFeedConsole(p []byte) {
	text := func(id string) string { return GetServer().S().Strings().GetStringInFile(strman.ID(id), "deathmsg.c") }
	attacker := ""
	killer, assist, victim := binary.LittleEndian.Uint16(p[2:]), binary.LittleEndian.Uint16(p[4:]), binary.LittleEndian.Uint16(p[6:])
	players := &GetServer().S().Players
	if pl := players.ByID(int(killer)); killer != 0 && pl != nil {
		attacker = consoleCommandFormat(text("die.c:LocalizeAttacker"), pl.Name())
		if a := players.ByID(int(assist)); assist != 0 && a != nil {
			attacker += " + " + a.Name()
		}
	} else {
		attacker = consoleCommandFormat(text("die.c:LocalizeAttacker"), text("die.c:AttackerNasty"))
	}
	who := ""
	if pl := players.ByID(int(victim)); victim != 0 && pl != nil {
		who = consoleCommandFormat(text("die.c:LocalizeVictim"), pl.Name())
	}
	GetConsole().Print(console.ColorWhite, consoleCommandFormat(alloc.GoString16(memmap.PtrUint16(0x587000, 161668)), who, attacker))
}
func combatFeedDraw() {
	combatFeedRows = 0
	r := GetClient().R2()
	r.Data().SetAlphaEnabled(false)
	r.Data().SetMultiply14(0)
	r.Data().SetColorize17(0)
	limit := int(nox_win_height) / 4 / 36
	for index := combatFeedRead; index != combatFeedWrite; index = (index + 1) % 100 {
		if int(combatFeedRows) > limit {
			break
		}
		if gameFrame()-combatFeedRing()[combatFeedRead][5] <= 90 {
			combatFeedRow(&combatFeedRing()[index])
			combatFeedRows++
		} else {
			combatFeedRead = (combatFeedRead + 1) % 100
		}
	}
}
func combatFeedRow(row *[6]uint32) {
	r := GetClient().R2()
	face := r.GetFonts().AsFont(r.GetFonts().FontPtrByName("large"))
	var names [3]string
	for i, code := range row[:3] {
		if code != 0 {
			if pl := GetServer().S().Players.ByID(int(code)); pl != nil {
				names[i] = pl.Name()
			}
		}
	}
	if row[0] != 0 && GetServer().S().Players.ByID(int(row[0])) != nil && row[1] != 0 && GetServer().S().Players.ByID(int(row[1])) != nil {
		names[1] = "+" + names[1]
	}
	var handle unsafe.Pointer
	if row[4] == 1 {
		if typ := GetClient().Cli().Things.TypeByInd(int(row[3])); typ != nil {
			if typ.ObjClass&0x1001000 != 0 {
				objectBaseMaterials(int(row[3]))
			}
			handle = memmap.WordPtr(typ.PrettyImage)
		}
	} else if row[4] == 2 {
		id := 0
		switch row[3] {
		case 1, 12:
			id = 5
		case 2:
			handle = unsafe.Pointer(Nox_xxx_spellGetAbilityIcon_425310(1, 0))
		case 4:
			id = 130
		case 5:
			id = 60
		case 9, 17:
			id = 43
		case 15:
			id = 56
		case 16:
			id = 16
		}
		if id != 0 {
			handle = Nox_xxx_spellIcon_424A90(id)
		}
	}
	if handle == nil {
		handle = combatFeedDefaultIcon
	}
	var off, sz image.Point
	var img *noxrender.Image
	if handle != nil {
		img = r.GetBag().AsImage(noxrender.ImageHandle(handle))
		off, sz, _ = img.Meta()
	}
	widths := 0
	height := 0
	for _, name := range names {
		size := r.GetStringSizeWrapped(face, name, 0)
		widths += size.X
		height = size.Y
	}
	x := (int(nox_win_width) - (widths + sz.X + 10)) / 2
	y := int(36 * combatFeedRows)
	r.DrawRectFilledAlpha(x-5, y, widths+sz.X+20, 36)
	textY := y + (36-height)/2
	draw := func(i int) {
		if row[i] == 0 {
			return
		}
		color := noxcolor.RGBA5551(memmap.Uint32(0x5D4594, 2597996))
		if row[i] == uint32(ClientPlayerNetCode()) {
			color = noxcolor.RGBA5551(dword_8531A0_2572)
		}
		r.Data().SetTextColor(color)
		x = r.DrawString(face, names[i], image.Pt(x, textY))
	}
	draw(0)
	draw(1)
	x += 5
	if img != nil {
		r.DrawImageAt(img, image.Pt(x-off.X, y+(36-sz.Y)/2-off.Y))
	}
	x += sz.X + 5
	draw(2)
}
