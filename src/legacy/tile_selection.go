package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func tileDefinitionsAll() *[176]server.TileDef {
	return (*[176]server.TileDef)(unsafe.Pointer(&worldTileDefinitions[0]))
}

func selectTileName(name *byte) bool {
	selected := memmap.PtrUint32(0x973F18, 35912)
	found := false
	// Scan physical storage, including entries beyond nox_tile_def_cnt. Keep
	// the shared C comparator until string handling is ported; it uses the C
	// locale's byte folding, rather than Unicode case folding.
	tiles := tileDefinitionsAll()
	for i := range tiles {
		p := &tiles[i]
		if textCompareNarrow((*byte)(unsafe.Pointer(&p.NameBuf[0])), (*byte)(unsafe.Pointer(name))) == 0 {
			*selected = uint32(i)
			found = true
		}
	}
	if textCompareNarrow((*byte)(unsafe.Pointer(name)), (*byte)(unsafe.Pointer(alloc.InternCString("NONE")))) == 0 {
		*selected = 255
		return true
	}
	if !found {
		*selected = 0
	}
	return found
}

func selectTileImage(index int32) bool {
	selected := memmap.PtrUint32(0x973F18, 35912)
	if index < 0 || index >= 176 {
		*selected = 0
		return false
	}
	*selected = uint32(index)
	return true
}

func selectTileVariation(variation int32) bool {
	// Existing callers select/reset a valid tile first. NONE (255) is not a
	// valid tile for this operation in the original C either.
	p := &tileDefinitionsAll()[int32(memmap.Uint32(0x973F18, 35912))]
	if variation <= int32(p.Field52)*int32(p.Field53)-1 {
		dword_5d4594_3835348 = uint32(variation)
		return true
	}
	dword_5d4594_3835348 = 0
	return false
}

func setTileFlag(value int32) bool {
	if value != 0 && value != 1 {
		return false
	}
	*memmap.PtrUint32(0x973F18, 35916) = uint32(value)
	return true
}

func nox_xxx_tileGetDefByName_51D4D0(name *int8) int32 {
	return int32(bool2int(selectTileName((*byte)(unsafe.Pointer(name)))))
}
