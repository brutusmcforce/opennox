package legacy

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

type briefingSlide struct {
	Image    uint32
	Text     *uint16
	Voice    *byte
	Duration uint32
}
type briefingChapter struct{ Begin, Loss briefingSlide }
type briefingScore struct {
	Player                            *server.Player
	Kills, Generators, Secrets, Found uint16
	Total                             uint32
}

var _ = [1]struct{}{}[32-unsafe.Sizeof(briefingChapter{})]
var _ = [1]struct{}{}[16-unsafe.Sizeof(briefingScore{})]

func briefingString(id string) string {
	s := GetServer().S().Strings().GetStringInFile(strman.ID(id), "GUIBrief.c")
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	return s
}

// briefingFormat formats localized strings with the original optional arguments.
func briefingFormat(id string, args ...any) string {
	return fmt.Sprintf(briefingString(id), args...)
}
func briefingLoadImage(name string) uint32 {
	im := Nox_xxx_gLoadImg(name)
	if im == nil {
		return 0
	}
	return uint32(uintptr(im.C()))
}
func briefingLoadChapters() *uint16 {
	table := (*[33]briefingChapter)(memmap.PtrOff(0x5D4594, 831300))
	for ch := 0; ch < 11; ch++ {
		for class := 0; class < 3; class++ {
			name := alloc.GoString(*(**byte)(memmap.PtrOff(0x587000, 122944+uintptr(4*class))))
			duration := uint32(ch + 3)
			if ch == 0 {
				duration = 1
			} else if ch == 1 {
				duration = []uint32{4, 2, 3}[class]
			}
			row := &table[class*11+ch]
			for k, slide := range []*briefingSlide{&row.Begin, &row.Loss} {
				kind := "Begin"
				if k == 1 {
					kind = "Loss"
				}
				image := fmt.Sprintf("%sChapter%s%d", name, kind, ch+1)
				slide.Image = briefingLoadImage(image)
				v, _ := GetServer().S().Strings().GetVariantInFile(strman.ID("Briefing:"+image), "GUIBrief.c")
				slide.Text = alloc.InternCString16(v.Str)
				slide.Voice = alloc.InternCString(v.Str2)
				slide.Duration = duration
			}
		}
	}
	*memmap.PtrUint32(0x5D4594, 831264) = briefingLoadImage("CreditsImage")
	v, _ := GetServer().S().Strings().GetVariantInFile("Nox:Credits", "GUIBrief.c")
	text := alloc.InternCString16(v.Str)
	*memmap.PtrPtr(0x5D4594, 831268) = unsafe.Pointer(text)
	*memmap.PtrPtr(0x5D4594, 831272) = unsafe.Pointer(alloc.InternCString(v.Str2))
	return text
}
func briefingSprites() []*uint32 {
	return []*uint32{
		(*uint32)(&dword_5d4594_832496),
		(*uint32)(&dword_5d4594_832492),
		(*uint32)(&dword_5d4594_832500),
		(*uint32)(&dword_5d4594_832504),
		(*uint32)(&dword_5d4594_832508),
		(*uint32)(&dword_5d4594_832512),
		(*uint32)(&dword_5d4594_832516),
		(*uint32)(&dword_5d4594_832520),
		(*uint32)(&dword_5d4594_832524),
		(*uint32)(&dword_5d4594_832528),
		(*uint32)(&dword_5d4594_832532),
		(*uint32)(&dword_5d4594_832536),
	}
}
func briefingInitSprites() *client.Drawable {
	if dword_5d4594_832484 == 0 {
		dword_5d4594_832484 = uint32(uintptr(GetClient().R2().GetFonts().FontPtrByName("default")))
	}
	names := []string{"GauntletExitB", "BeholderGenerator", "Ankh", "SoulGate", "SilverKey", "GoldKey", "QuestGoldChest", "QuestGoldPile", "DunMirChest4", "WarHammer", "HastePotion", "ConjurerSpellBook"}
	var dr *client.Drawable
	for i, p := range briefingSprites() {
		dr = (*client.Drawable)(unsafe.Pointer(uintptr(*p)))
		if dr == nil {
			dr = GetClient().Nox_new_drawable_for_thing(GetClient().Cli().Things.IndByID(names[i]))
			*p = uint32(uintptr(dr.C()))
		}
		*effectWord(dr, 120) |= 0x1000000
	}
	return dr
}
func briefingCompare(a, b uint32) int {
	if a == b {
		return 0
	}
	if a < b {
		return 1
	}
	return -1
}
func briefingSetImage(image uint32) uint32 {
	if image == 0 {
		image = briefingLoadImage("WarriorChapterBegin8")
	}
	*memmap.PtrUint32(0x5D4594, 832460) = image
	return image
}
func briefingSetCaption(text uint32) uint32 { *memmap.PtrUint32(0x5D4594, 832464) = text; return text }
func briefingSetStage(stage uint32)         { *memmap.PtrUint32(0x5D4594, 832468) = stage }
func briefingStage() uint32                 { return memmap.Uint32(0x5D4594, 832468) }
func briefingWinReport(data unsafe.Pointer) int {
	b := unsafe.Slice((*byte)(data), 90)
	rows := (*[6]briefingScore)(memmap.PtrOff(0x5D4594, 832364))
	clear(rows[:])
	*memmap.PtrUint32(0x5D4594, 832356) = uint32(binary.LittleEndian.Uint16(b[2:]))
	*memmap.PtrUint32(0x5D4594, 831228) = uint32(binary.LittleEndian.Uint16(b[4:]))
	n := 0
	for i := range rows {
		r := b[6+i*14:]
		id := binary.LittleEndian.Uint16(r)
		if id == 0 {
			continue
		}
		rows[i] = briefingScore{Player: GetServer().S().Players.ByID(int(id)), Kills: binary.LittleEndian.Uint16(r[8:]), Generators: binary.LittleEndian.Uint16(r[2:]), Secrets: binary.LittleEndian.Uint16(r[4:]), Found: binary.LittleEndian.Uint16(r[6:]), Total: binary.LittleEndian.Uint32(r[10:])}
		n++
	}
	// The sender packs participating records. Sorting only the original count also
	// preserves the captured bounded sparse-input behavior and stable equal scores.
	sort.SliceStable(rows[:n], func(i, j int) bool { return briefingCompare(rows[i].Total, rows[j].Total) < 0 })
	if dword_5d4594_832476 == 0 {
		w := (*gui.Window)(unsafe.Pointer(uintptr(nox_wnd_briefing_831232))).ChildByID(1010)
		r := GetClient().R2()
		font := r.GetFonts().AsFont(w.DrawData().FontPtr)
		width := 0
		for _, id := range []string{"GUIBrief.c:GeneratorsDestroyed", "GUIBrief.c:Kills", "GUIBrief.c:numSecretsFound", "GUIBrief.c:TotalScore"} {
			if v := r.GetStringSizeWrapped(font, briefingString(id), 0).X; v > width {
				width = v
			}
		}
		if width > 85 {
			width = 85
		}
		dword_5d4594_832476 = uint32(width)
	}
	return briefingShow(254, 1, 1)
}
func briefingSelection(data unsafe.Pointer, show int, instructions bool) int {
	dword_5d4594_832480 = 0
	Nox_client_resetScreenParticles_431510()
	Nox_xxx_bookHideMB_45ACA0(1)
	Sub_446780()
	b := unsafe.Slice((*byte)(data), 69)
	briefingSetImage(briefingLoadImage(alloc.GoString(&b[5])))
	if name := alloc.GoString(&b[37]); name != "" {
		briefingSetCaption(uint32(uintptr(unsafe.Pointer(alloc.InternCString16(briefingString(name))))))
	} else {
		off := uintptr(832544)
		if instructions {
			off = 832548
		}
		briefingSetCaption(uint32(uintptr(memmap.PtrOff(0x5D4594, off))))
	}
	briefingSetStage(uint32(binary.LittleEndian.Uint16(b[2:])))
	mode := byte(2)
	if instructions {
		mode = 4
	} else if b[4]&2 != 0 {
		dword_5d4594_832480 = 1
	}
	if show != 0 {
		return briefingShow(254, 1, mode)
	}
	return show
}
