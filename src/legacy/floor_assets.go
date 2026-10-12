package legacy

import (
	"bytes"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func floorAssetName(f *binfile.MemFile) string {
	n := int(f.ReadU8())
	b := make([]byte, n)
	f.Read(b)
	if i := bytes.IndexByte(b, 0); i >= 0 {
		b = b[:i]
	}
	return string(b)
}
func floorAssetFacade(tile *server.TileDef) bool {
	for i := uintptr(0); ; i++ {
		p := *memmap.PtrPtr(0x587000, 26488+4*i)
		if p == nil {
			return false
		}
		if alloc.GoString((*byte)(p)) == tile.Name() {
			return true
		}
	}
}
func floorAssetDefinition(f *binfile.MemFile, scratch []byte) int {
	if uint32(worldTileDefinitionCount) >= 176 {
		return 0
	}
	tile := &tileDefinitionsAll()[uint32(worldTileDefinitionCount)]
	f.Skip(4)
	name := floorAssetName(f)
	clear(tile.NameBuf[:31])
	copy(tile.NameBuf[:31], name)
	r, g, b := f.ReadU8(), f.ReadU8(), f.ReadU8()
	tile.Color48 = noxcolor.RGB5551Color(r, g, b).Color32()
	if r == 255 && g == 0 && b == 255 {
		tile.Color48 = 0x80000000
	}
	tile.Field58 = 0
	if floorAssetFacade(tile) {
		tile.Field58 = 1
	}
	tile.Field36 = f.ReadI32()
	tile.Field40 = f.ReadI32()
	tile.Field57 = f.ReadU8()
	tile.Field53 = f.ReadU8()
	tile.Field52 = f.ReadU8()
	tile.Field54 = f.ReadU8()
	tile.Field46 = 0
	tile.Field55 = f.ReadU8()
	tile.Field56 = tile.Field55
	tile.Field44 = uint16(tile.Field52) * uint16(tile.Field53)
	tile.Data32 = nil
	count := int(tile.Field52) * int(tile.Field53) * int(tile.Field54)
	for i := 0; i < count; i++ {
		id := f.ReadI32()
		scratch[0] = memmap.Uint8(0x5D4594, 251576)
		if id == -1 {
			f.Skip(1)
			n := int(f.ReadU8())
			f.Read(scratch[:n])
			scratch[n] = 0
		}
	}
	worldTileDefinitionCount++
	return 1
}
func edgeAssetDefinition(f *binfile.MemFile, scratch []byte) int {
	index := uint32(dword_5d4594_251572)
	if index >= 64 {
		return 0
	}
	base := uintptr(28644 + 60*index)
	f.Skip(4)
	name := floorAssetName(f)
	dst := unsafe.Slice(memmap.PtrUint8(0x85B3FC, base), 32)
	copy(dst, name)
	dst[len(name)] = 0
	*memmap.PtrInt32(0x85B3FC, base+36) = f.ReadI32()
	*memmap.PtrInt32(0x85B3FC, base+40) = f.ReadI32()
	*memmap.PtrUint8(0x85B3FC, base+57) = f.ReadU8()
	a := f.ReadU8()
	*memmap.PtrUint8(0x85B3FC, base+54) = a
	*memmap.PtrUint16(0x85B3FC, base+46) = 0
	v := f.ReadU8()
	*memmap.PtrUint8(0x85B3FC, base+55) = v
	*memmap.PtrUint8(0x85B3FC, base+56) = v
	if f.ReadU8() == 1 {
		return 0
	}
	b, c := f.ReadU8(), f.ReadU8()
	*memmap.PtrUint8(0x85B3FC, base+53) = b
	*memmap.PtrUint8(0x85B3FC, base+52) = c
	*memmap.PtrUint16(0x85B3FC, base+44) = 2 * (uint16(b) + uint16(c))
	*memmap.PtrPtr(0x85B3FC, base+32) = nil
	count := 2 * int(a) * (int(b) + int(c))
	for i := 0; i < count; i++ {
		id := f.ReadI32()
		scratch[0] = memmap.Uint8(0x5D4594, 251580)
		if id == -1 {
			f.Skip(1)
			n := int(f.ReadU8())
			f.Read(scratch[:n])
			scratch[n] = 0
		}
	}
	if f.ReadU32() != 0x454e4420 {
		return 0
	}
	dword_5d4594_251572++
	return 1
}
func floorAssetSkipImages(f *binfile.MemFile, count int) {
	for i := 0; i < count; i++ {
		if f.ReadI32() == -1 {
			f.Skip(1)
			n := int(f.ReadU8())
			f.Skip(n)
		}
	}
}
func floorAssetSkip(f *binfile.MemFile) int {
	f.Skip(4)
	n := int(f.ReadU8())
	f.Skip(n + 12)
	a, b, c := f.ReadU8(), f.ReadU8(), f.ReadU8()
	f.Skip(1)
	floorAssetSkipImages(f, int(a)*int(b)*int(c))
	if f.ReadU32() == 0x454e4420 {
		return 1
	}
	return 0
}
func edgeAssetSkip(f *binfile.MemFile, scratch []byte) int {
	f.Skip(4)
	n := int(f.ReadU8())
	f.Read(scratch[:n])
	f.Skip(9)
	a := f.ReadU8()
	f.Skip(1)
	if f.ReadU8() == 1 {
		return 0
	}
	b, c := f.ReadU8(), f.ReadU8()
	floorAssetSkipImages(f, 2*int(a)*(int(b)+int(c)))
	if f.ReadU32() == 0x454e4420 {
		return 1
	}
	return 0
}
func floorAssetImage(f *binfile.MemFile, scratch []byte, typ *byte, initial byte) unsafe.Pointer {
	index := f.ReadI32()
	scratch[0] = initial
	if index == -1 {
		*typ = f.ReadU8()
		n := int(f.ReadU8())
		f.Read(scratch[:n])
		scratch[n] = 0
	}
	name := alloc.GoString((*byte)(unsafe.Pointer(&scratch[0])))
	return unsafe.Pointer(GetClient().R2().GetBag().ImageRef(int(index), *typ, name).C())
}
func floorAssetBind(f *binfile.MemFile, scratch []byte) int {
	f.Skip(4)
	name := floorAssetName(f)
	index := -1
	for i := 0; i < int(worldTileDefinitionCount); i++ {
		if tileDefinitionsAll()[i].Name() == name {
			index = i
			break
		}
	}
	// The C baseline must first make the empty-definition case a defined failure.
	if index < 0 {
		return 0
	}
	f.Skip(12)
	a, b, c := f.ReadU8(), f.ReadU8(), f.ReadU8()
	f.Skip(1)
	count := int(a) * int(b) * int(c)
	ptr := legacyCalloc(uintptr(count), 4)
	tileDefinitionsAll()[index].Data32 = ptr
	words := unsafe.Slice((*unsafe.Pointer)(ptr), count)
	typ := a
	for i := range words {
		words[i] = floorAssetImage(f, scratch, &typ, memmap.Uint8(0x5D4594, 1193192))
	}
	return 1
}
func edgeAssetBind(f *binfile.MemFile, scratch []byte) int {
	f.Skip(4)
	name := floorAssetName(f)
	index := -1
	for i := 0; i < int(int32(dword_5d4594_251572)); i++ {
		if alloc.GoString(memmap.PtrUint8(0x85B3FC, 28644+60*uintptr(i))) == name {
			index = i
			break
		}
	}
	if index < 0 {
		return 0
	}
	f.Skip(9)
	a := f.ReadU8()
	f.Skip(1)
	if f.ReadU8() == 1 {
		return 0
	}
	b, c := f.ReadU8(), f.ReadU8()
	count := 2 * int(a) * (int(b) + int(c))
	// Keep the original allocation extent: handles use four bytes of each five.
	ptr := legacyCalloc(uintptr(count), 5)
	*memmap.PtrPtr(0x85B3FC, 28676+60*uintptr(index)) = ptr
	if ptr == nil {
		return 0
	}
	words := unsafe.Slice((*unsafe.Pointer)(ptr), count)
	typ := a
	for i := range words {
		words[i] = floorAssetImage(f, scratch, &typ, memmap.Uint8(0x5D4594, 1193196))
	}
	if f.ReadU32() == 0x454e4420 {
		return 1
	}
	return 0
}
func floorAssetFree() {
	for i := 0; i < int(worldTileDefinitionCount); i++ {
		p := &tileDefinitionsAll()[i].Data32
		if *p != nil {
			legacyFree(*p)
			*p = nil
		}
	}
}
func edgeAssetFree() {
	for i := int32(0); i < int32(dword_5d4594_251572); i++ {
		p := memmap.PtrPtr(0x85B3FC, 28676+60*uintptr(i))
		if *p != nil {
			legacyFree(*p)
			*p = nil
		}
	}
}
