package legacy

import (
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/datapath"
	"github.com/brutusmcforce/opennox/libs/ifs"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func mapCatalogFirst() *Nox_map_list_item {
	return (*Nox_map_list_item)(unsafe.Pointer(listNext((*legacyListNode)(unsafe.Pointer(&legacyGlobals.nox_common_maplist)))))
}
func mapCatalogNext(p *Nox_map_list_item) *Nox_map_list_item {
	if p == nil {
		return nil
	}
	return (*Nox_map_list_item)(unsafe.Pointer(listNext((*legacyListNode)(unsafe.Pointer(p)))))
}
func mapCatalogAdd(p *Nox_map_list_item) {
	name := alloc.GoString(&p.Name[0])
	for it := mapCatalogFirst(); it != nil; it = mapCatalogNext(it) {
		if name <= alloc.GoString(&it.Name[0]) {
			listAppend((*legacyListNode)(unsafe.Pointer(it)), (*legacyListNode)(unsafe.Pointer(p)))
			return
		}
	}
	listAppend((*legacyListNode)(unsafe.Pointer(&legacyGlobals.nox_common_maplist)), (*legacyListNode)(unsafe.Pointer(p)))
}
func mapCatalogFree() {
	for it := mapCatalogFirst(); it != nil; {
		next := mapCatalogNext(it)
		listRemove((*legacyListNode)(unsafe.Pointer(it)))
		legacyFree(unsafe.Pointer(it))
		it = next
	}
}
func mapASCIIEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		x, y := a[i], b[i]
		if x >= 'A' && x <= 'Z' {
			x += 'a' - 'A'
		}
		if y >= 'A' && y <= 'Z' {
			y += 'a' - 'A'
		}
		if x != y {
			return false
		}
	}
	return true
}
func mapCycleMask(group int) uint32 { return memmap.Uint32(0x587000, 191836+8*uintptr(group)) }
func mapCycleGroup(name string) int {
	for i := 0; i < 6; i++ {
		if mapASCIIEqual(alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 191832+8*uintptr(i)))), name) {
			return i
		}
	}
	return -1
}
func mapCycleFlagGroup(flags uint32) int {
	for i := 0; i < 6; i++ {
		if flags&mapCycleMask(i) != 0 {
			return i
		}
	}
	return 0
}
func mapCycleCounts() []uint32  { return unsafe.Slice(memmap.PtrUint32(0x5D4594, 1548428), 6) }
func mapCycleIndices() []uint32 { return unsafe.Slice(memmap.PtrUint32(0x5D4594, 1548452), 6) }
func mapCycleName(group, index int) *byte {
	return memmap.PtrUint8(0x5D4594, 1529228+128*uintptr(index+25*group))
}
func mapCycleEnabled() int {
	if memmap.Uint32(0x5D4594, 1548484) != 0 || noxflags.HasEngine(noxflags.EngineNoRendering) {
		return 1
	}
	return 0
}
func mapCycleSetEnabled(v uint32) uint32 { *memmap.PtrUint32(0x5D4594, 1548484) = v; return v }
func mapCycleReset()                     { clear(mapCycleCounts()); clear(mapCycleIndices()) }
func mapCycleSetIndex(flags uint32, index uint32) int {
	group := mapCycleFlagGroup(flags)
	mapCycleIndices()[group] = index
	return group
}
func mapCycleGetIndex(flags uint32) uint32 { return mapCycleIndices()[mapCycleFlagGroup(flags)] }
func mapCycleNext() *byte {
	group := mapCycleFlagGroup(uint32(noxflags.GetGame()))
	count := mapCycleCounts()[group]
	if count == 0 {
		return nil
	}
	indices := mapCycleIndices()
	index := indices[group]
	// Preserve the original strict greater-than comparison, including index==count.
	if int32(index) > int32(count) {
		index = 0
		indices[group] = 0
	}
	indices[group] = (indices[group] + 1) % count
	return mapCycleName(group, int(index))
}
func mapCycleStrip(p *byte) {
	if p == nil {
		return
	}
	raw := unsafe.Slice(p, len(alloc.GoString(p)))
	// CR is replaced first; an LF after it is no longer visible to the second scan.
	if i := strings.IndexByte(string(raw), '\r'); i >= 0 {
		raw[i] = 0
		raw = raw[:i]
	}
	if i := strings.IndexByte(string(raw), '\n'); i >= 0 {
		raw[i] = 0
	}
}
func mapCycleLoad() {
	counts := mapCycleCounts()
	clear(counts)
	path := datapath.Data() + "\\mapcycle.txt"
	alloc.StrCopyZero(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1524108), 1024), path)
	file, err := ifs.Open(path)
	if err != nil {
		*memmap.PtrUint32(0x5D4594, 1548484) = 0
		return
	}
	f := binfile.NewTextFile(file)
	defer f.Close()
	// Use the actual text owner: ReadString consumes a complete line, while the
	// old fgets facade copies at most 126 bytes and rejects an EOF-terminated tail.
	read := func() (string, bool) {
		line, err := f.ReadString()
		if err != nil {
			return "", false
		}
		if len(line) > 126 {
			line = line[:126]
		}
		if i := strings.IndexByte(string(line), 0); i >= 0 {
			line = line[:i]
		}
		if i := strings.IndexAny(string(line), "\r\n"); i >= 0 {
			line = line[:i]
		}
		return string(line), true
	}
	group := 1
	if line, ok := read(); ok {
		if g := mapCycleGroup(line); g >= 0 {
			group = g
		} else {
			alloc.StrCopyZero(unsafe.Slice(mapCycleName(1, int(counts[1])), 128), line)
			counts[1]++
		}
	}
	for {
		line, ok := read()
		if !ok {
			break
		}
		if g := mapCycleGroup(line); g >= 0 {
			group = g
			continue
		}
		if counts[group] >= 25 || line == "" {
			continue
		}
		// strtok discards leading delimiters and stops at the first following one.
		name := strings.TrimLeft(line, ".\n")
		if i := strings.IndexAny(name, ".\n"); i >= 0 {
			name = name[:i]
		}
		if err := Nox_common_checkMapFile(name); err != nil {
			gameLog.Println("check map file:", err)
			continue
		}
		if uint32(Nox_mapToGameFlags(int(memmap.Uint32(0x973F18, 3800))))&mapCycleMask(group) == 0 {
			continue
		}
		alloc.StrCopyZero(unsafe.Slice(mapCycleName(group, int(counts[group])), 128), line)
		counts[group]++
	}
}

type mapQuestRow struct {
	Group uint32
	Name  [20]byte
	Used  uint32
	Last  uint32
}

var _ = [1]struct{}{}[32-unsafe.Sizeof(mapQuestRow{})]

func mapQuestRows() []mapQuestRow {
	return unsafe.Slice((*mapQuestRow)(memmap.PtrOff(0x5D4594, 1525132)), 128)
}
func mapQuestBuild() int {
	rows := mapQuestRows()
	count := 0
	dword_5d4594_1548476 = 0
	for p := mapCatalogFirst(); p != nil; p = mapCatalogNext(p) {
		if p.Field_6 == 0 || Nox_mapToGameFlags(int(int32(p.Field_7)))&noxflags.GameModeQuest == 0 || count >= 128 {
			continue
		}
		row := &rows[count]
		alloc.StrCopyZero(row.Name[:], alloc.GoString(&p.Name[0])+".map")
		row.Group = 0
		count++
		dword_5d4594_1548476 = uint32(count)
	}
	dword_5d4594_1548476 = uint32(count)
	group := uint32(1)
	for i := 0; i < count; i++ {
		if rows[i].Group != 0 {
			continue
		}
		rows[i].Group = group
		group++
		for j := i + 1; j < count; j++ {
			a, b := alloc.GoString(&rows[i].Name[0]), alloc.GoString(&rows[j].Name[0])
			a = a[:min(len(a), 6)]
			b = b[:min(len(b), 6)]
			if mapASCIIEqual(a, b) {
				rows[j].Group = rows[i].Group
			}
		}
	}
	return count
}
func mapQuestReset() {
	dword_5d4594_1548480 = 1000
	for i := 0; i < int(dword_5d4594_1548476); i++ {
		mapQuestRows()[i].Used = 0
		mapQuestRows()[i].Last = 0
	}
}
func mapQuestChoose() *byte {
	count := int(dword_5d4594_1548476)
	if count == 0 {
		return nil
	}
	rows := mapQuestRows()
	if count == 1 {
		return &rows[0].Name[0]
	}
	maximum := uint32(0)
	// The C maximum comparison promotes its signed local to unsigned.
	for i := 0; i < count; i++ {
		if rows[i].Used > maximum {
			maximum = rows[i].Used
		}
	}
	if maximum == 0 {
		return &rows[GetServer().S().Rand.Logic.IntClamp(0, count-1)].Name[0]
	}
	equal := true
	for i := 1; i < count; i++ {
		if rows[i].Used != rows[0].Used {
			equal = false
			break
		}
	}
	if equal {
		maximum++
	}
	last := int(memmap.Uint32(0x587000, 191880))
	clock := uint32(dword_5d4594_1548480)
	eligible := func(i int) bool {
		return rows[i].Used < maximum && i != last && rows[i].Group != rows[last].Group && clock-rows[i].Last > 4
	}
	candidates := 0
	for i := 0; i < count; i++ {
		if eligible(i) {
			candidates++
		}
	}
	// The baseline correction falls back when family/history exclusions exhaust the catalog.
	if candidates == 0 {
		return &rows[GetServer().S().Rand.Logic.IntClamp(0, count-1)].Name[0]
	}
	choice := GetServer().S().Rand.Logic.IntClamp(0, candidates-1)
	index := 0
	for i := 0; i < count; i++ {
		if eligible(i) {
			if index == choice {
				return &rows[i].Name[0]
			}
			index++
		}
	}
	return &rows[choice].Name[0]
}
func mapQuestPlayed(name *byte) {
	if name == nil {
		return
	}
	rows := mapQuestRows()
	count := int(dword_5d4594_1548476)
	for i := 0; i < count; i++ {
		if !mapASCIIEqual(alloc.GoString(&rows[i].Name[0]), alloc.GoString(name)) {
			continue
		}
		*memmap.PtrUint32(0x587000, 191880) = uint32(i)
		clock := uint32(dword_5d4594_1548480)
		for j := 0; j < count; j++ {
			if rows[j].Group == rows[i].Group {
				rows[j].Used++
				rows[j].Last = clock
			}
		}
		dword_5d4594_1548480 = uint32(clock + 1)
		return
	}
}
