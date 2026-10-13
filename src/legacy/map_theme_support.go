package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"math"
	"os"
	"time"
	"unsafe"
)

// Theme records remain owned by the shared C allocator while the rest of the
// map generator consumes them. Numeric conversion preserves the libc boundaries.
// They are indexed with pointer word arithmetic (base + 60*i), so they are word regions (see mapRoomCalloc).
func mapThemeAlloc(n, size uint32) uint32 { return mapRoomRaw(mapRoomCalloc(n, uintptr(size))) }
func mapThemeFree(p uint32)               { mapRoomRelease(mapRoomPointer(p)) }
func mapThemeTemplate() *uint32           { return (*uint32)(unsafe.Pointer(&dword_5d4594_2487524)) }

var mapThemeClock = func() uint32 { return uint32(time.Now().Unix()) }

func mapThemeInt(s string) int32 { return resourceAtoi(s) }
func mapThemeFloat(s string) float64 {
	v, _, ok := resourceParseFloat(s, 64, false)
	if !ok {
		for len(s) != 0 && resourceSpace(s[0]) {
			s = s[1:]
		}
		// strtod consumes the signed zero before an incomplete hexadecimal prefix.
		// Keep this compatibility correction local to the former atof boundary.
		if len(s) >= 3 && s[0] == '-' && s[1] == '0' && (s[2] == 'x' || s[2] == 'X') {
			return math.Copysign(0, -1)
		}
	}
	return v
}
func mapThemeTable(off uintptr, s string) int {
	for i := 0; ; i++ {
		p := *memmap.PtrUint32(0x587000, off+uintptr(4*i))
		if p == 0 {
			return -1
		}
		if mapThemeLower(populationString(p)) == mapThemeLower(s) {
			return i
		}
	}
}
func mapThemeReadByte(f uint32) (byte, bool) {
	file := fileByHandle((*FILE)(mapRoomPointer(f)))
	var b [1]byte
	_, err := file.Bin.Read(b[:])
	if err != nil {
		file.Err = err
	}
	return b[0], file.Err == nil
}
func mapThemeFile(cfg, name uint32) uint32 {
	clear(unsafe.Slice(mapThemeByte(cfg, 0), 1116))
	defaults := map[int]uint32{4: 5, 8: 2, 12: 3, 16: 1, 20: 3, 24: 40, 28: 20, 32: 10, 36: 5, 40: 20, 48: 100, 52: 100, 56: 1, 64: 1157234688, 72: 100, 76: mapThemeClock()}
	for off, v := range defaults {
		*populationWord(cfg, off) = v
	}
	*populationBlob(2487520) = 0
	file, err := binfile.BinfileOpen("mapgen/"+populationString(name)+".thm", binfile.ReadOnly)
	if err != nil {
		if !os.IsNotExist(err) {
			binfile.Log.Println(err)
		}
		return 0
	}
	f := mapRoomRaw(unsafe.Pointer(NewFileHandle(file.File)))
	if err = file.SetKey(1); err != nil {
		binfile.Log.Println(err)
		return 0
	}
	result := uint32(1)
	for mapThemeNext(f) {
		switch mapThemeLower(mapThemeText()) {
		case "algorithm_data":
			result = mapThemeAlgorithm(cfg, f)
		case "exit":
			result = mapThemeExit(cfg, f)
		case "prefab":
			result = mapThemePrefabs(cfg, f)
		case "spell_set":
			result = mapThemeSpells(cfg, f)
		case "weapon_set":
			result = mapThemeEquipment(cfg, f, false)
		case "armor_set":
			result = mapThemeEquipment(cfg, f, true)
		case "decor":
			result = mapThemeDecor(cfg, f)
		case "ambient_light":
			for _, off := range []int{536, 540, 544} {
				if !mapThemeNext(f) {
					result = 0
					break
				}
				*populationWord(cfg, off) = uint32(mapThemeInt(mapThemeText()))
			}
		default:
			result = 0
		}
		if result == 0 {
			break
		}
	}
	file.Close()
	if result != 1 {
		return result
	}
	if mapThemeValidate(cfg+88) == 0 || mapThemeValidate(cfg+120) == 0 {
		return 0
	}
	if *populationWord(cfg, 184) != 0 && mapThemeValidate(cfg+184) == 0 {
		return 0
	}
	return 1
}
