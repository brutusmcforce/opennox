package legacy

import (
	"github.com/brutusmcforce/opennox/internal/binfile"
)

var (
	Sub_42BFB0             func()
	Nox_xxx_objectTOCgetTT func(a1 uint16) int
	Sub_42C310             func(a1 int, a2 uint16)
	Sub_42C2E0             func(a1 int) uint16
	Sub_42C300             func() uint16
	Sub_42BFE0             func()
	Sub_4E3AD0             func(ind int) int
)

func Sub_4F0640() {
	resourceLinkCatalogs()
}
func Sub_485CF0() {
	floorAssetFree()
}
func Sub_485F30() {
	edgeAssetFree()
}

func Nox_thing_skip_spells_415100(f *binfile.MemFile) {
	thingSkipSpells(f)
}

func Nox_thing_read_ability_415320(f *binfile.MemFile) {
	thingSkipAbilities(f)
}

func Nox_thing_skip_AUD_414D40(f *binfile.MemFile) {
	thingSkipAUD(f)
}

func Nox_thing_skip_AVNT_452B00(f *binfile.MemFile) {
	thingSkipAVNT(f)
}

func Nox_thing_read_image_415240(f *binfile.MemFile) {
	thingSkipImages(f)
}

func Nox_thing_read_floor_485B30(f *binfile.MemFile, buf []byte) int {
	if cap(buf) < 256*1024 {
		panic(cap(buf))
	}
	return floorAssetBind(f, thingReaderScratch(buf))
}

func Nox_thing_read_edge_485D40(f *binfile.MemFile, buf []byte) int {
	if cap(buf) < 256*1024 {
		panic(cap(buf))
	}
	return edgeAssetBind(f, thingReaderScratch(buf))
}

func Nox_thing_read_WALL_414F60(f *binfile.MemFile, buf []byte) int {
	if cap(buf) < 256*1024 {
		panic(cap(buf))
	}
	return thingSkipWall(f, thingReaderScratch(buf))
}

func Nox_thing_read_FLOR_414DB0(f *binfile.MemFile) int {
	return floorAssetSkip(f)
}

func Nox_thing_read_EDGE_414E70(f *binfile.MemFile, buf []byte) int {
	if cap(buf) < 256*1024 {
		panic(cap(buf))
	}
	return edgeAssetSkip(f, thingReaderScratch(buf))
}

func Nox_thing_read_audio_415660(f *binfile.MemFile, buf []byte) int {
	if cap(buf) < 256*1024 {
		panic(cap(buf))
	}
	return audioAssetDefinitions(f, thingReaderScratch(buf))
}

func Nox_thing_read_AVNT_452890(f *binfile.MemFile, buf []byte) int {
	if cap(buf) < 256*1024 {
		panic(cap(buf))
	}
	return audioAssetEvent(f, thingReaderScratch(buf))
}

func Nox_thing_read_FLOR_411540(f *binfile.MemFile, buf []byte) int {
	if cap(buf) < 256*1024 {
		panic(cap(buf))
	}
	return floorAssetDefinition(f, thingReaderScratch(buf))
}

func Nox_thing_read_EDGE_411850(f *binfile.MemFile, buf []byte) int {
	if cap(buf) < 256*1024 {
		panic(cap(buf))
	}
	return edgeAssetDefinition(f, thingReaderScratch(buf))
}

func LoadAllBinFileSectionsResetCounters() {
	worldTileDefinitionCount = 0
	dword_5d4594_251572 = 0
}

// The old C bridge accepted nonempty short slices backed by the checked capacity.
func thingReaderScratch(buf []byte) []byte {
	_ = buf[0] // Preserve the original empty-slice rejection before any reader runs.
	return buf[:cap(buf)]
}
