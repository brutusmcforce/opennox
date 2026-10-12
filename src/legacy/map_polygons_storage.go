package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"strconv"
	"unsafe"
)

var mapPolygonNext uint32 = 1
var mapPolygonVertexNext uint32 = 1
var mapPolygonRemapHead *mapPolygonRemap

type mapPolygonVertex struct {
	ID     uint32
	X, Y   float32
	Active uint32
}
type mapPolygonRemap struct {
	Index, Old uint32
	Next       *mapPolygonRemap
}
type mapPolygon struct {
	Metadata       unsafe.Pointer
	Name           [76]byte
	ID, Active     uint32
	Bounds         [4]int32
	Color          [3]byte
	_              byte
	Vertices       *uint32
	Enter, Leave   server.ScriptCallback
	Count          uint16
	Level          byte
	_              byte
	Flags, Visited uint32
}

var _ [140 - unsafe.Sizeof(mapPolygon{})]byte
var _ [unsafe.Sizeof(mapPolygon{}) - 140]byte
var _ [108 - unsafe.Offsetof(mapPolygon{}.Vertices)]byte
var _ [unsafe.Offsetof(mapPolygon{}.Vertices) - 108]byte
var _ [16 - unsafe.Sizeof(mapPolygonVertex{})]byte
var _ [unsafe.Sizeof(mapPolygonVertex{}) - 16]byte
var _ [12 - unsafe.Sizeof(mapPolygonRemap{})]byte
var _ [unsafe.Sizeof(mapPolygonRemap{}) - 12]byte

const mapPolygonUnset = uint32(0xdeadface)

func mapPolygonAt(i uint32) *mapPolygon {
	return (*mapPolygon)(memmap.PtrOff(0x5D4594, 552228+uintptr(i)*140))
}
func mapPolygonVertexAt(i uint32) *mapPolygonVertex {
	return (*mapPolygonVertex)(memmap.PtrOff(0x5D4594, 535844+uintptr(i)*16))
}
func mapPolygonGet(i uint32) *mapPolygon {
	if i == mapPolygonUnset {
		return nil
	}
	return mapPolygonAt(i)
}
func mapPolygonIDs(p *mapPolygon) []uint32 { return unsafe.Slice(p.Vertices, int(p.Count)) }
func mapPolygonRemapAdd(index, old uint32) *mapPolygonRemap {
	p := (*mapPolygonRemap)(legacyCalloc(1, 12))
	if p != nil {
		p.Index = index
		p.Old = old
		p.Next = mapPolygonRemapHead
		mapPolygonRemapHead = p
	}
	return p
}
func mapPolygonRemapClear() unsafe.Pointer {
	for p := mapPolygonRemapHead; p != nil; {
		next := p.Next
		legacyFree(unsafe.Pointer(p))
		p = next
	}
	mapPolygonRemapHead = nil
	return nil
}
func mapPolygonRemapIDs(p *mapPolygon) {
	for i, id := range mapPolygonIDs(p) {
		for node := mapPolygonRemapHead; node != nil; node = node.Next {
			if node.Old == id {
				mapPolygonIDs(p)[i] = node.Index
				break
			}
		}
	}
}
func mapPolygonVertexFirst() *mapPolygonVertex { return mapPolygonVertexAfter(0) }
func mapPolygonVertexAfter(id uint32) *mapPolygonVertex {
	for i := id + 1; i < mapPolygonVertexNext; i++ {
		p := mapPolygonVertexAt(i)
		if p.Active != 0 {
			return p
		}
	}
	return nil
}
func mapPolygonFirst() *mapPolygon { return mapPolygonAfter(0) }
func mapPolygonAfter(id uint32) *mapPolygon {
	for i := id + 1; i < mapPolygonNext; i++ {
		p := mapPolygonAt(i)
		if p.Active != 0 {
			return p
		}
	}
	return nil
}
func mapPolygonFreeVertex() uint32 {
	for i := uint32(1); i < mapPolygonVertexNext; i++ {
		if mapPolygonVertexAt(i).Active == 0 {
			return i
		}
	}
	i := mapPolygonVertexNext
	mapPolygonVertexNext++
	return i
}
func mapPolygonFreeIndex() uint32 {
	for i := uint32(1); i < mapPolygonNext; i++ {
		if mapPolygonAt(i).Active == 0 {
			return i
		}
	}
	i := mapPolygonNext
	mapPolygonNext++
	return i
}
func mapPolygonVertexSet(x, y float32, index, old uint32) *mapPolygonVertex {
	if old != 0 {
		mapPolygonRemapAdd(index, old)
	}
	p := mapPolygonVertexAt(index)
	p.ID = index
	if index >= mapPolygonVertexNext {
		mapPolygonVertexNext = index + 1
	}
	p.X = x
	p.Y = y
	p.Active = 1
	return p
}
func mapPolygonResetRecords() unsafe.Pointer {
	var result unsafe.Pointer
	for i := uint32(1); i < 256; i++ {
		p := mapPolygonAt(i)
		if p.Metadata != nil {
			if memmap.Uint32(0x5D4594, 588076) != 0 {
				legacyFree(p.Metadata)
			}
			p.Metadata = nil
		}
		result = unsafe.Pointer(p.Vertices)
		if p.Vertices != nil {
			if memmap.Uint32(0x5D4594, 588076) != 0 {
				legacyFree(unsafe.Pointer(p.Vertices))
			}
			p.Vertices = nil
		}
		p.Count = 0
		p.ID = 0
		p.Enter.Func = -1
		p.Leave.Func = -1
		p.Active = 0
	}
	mapPolygonNext = 1
	return result
}
func mapPolygonResetVertices() unsafe.Pointer {
	for i := uint32(1); i < 1024; i++ {
		mapPolygonVertexAt(i).Active = 0
	}
	mapPolygonVertexNext = 1
	return memmap.PtrOff(0x5D4594, 552240)
}
func mapPolygonReset() unsafe.Pointer {
	mapPolygonResetRecords()
	mapPolygonResetVertices()
	mapPolygonRemapClear()
	*memmap.PtrUint32(0x5D4594, 588076) = 1
	*memmap.PtrUint16(0x5D4594, 588072) = 0
	return nil
}
func mapPolygonDefaultsRead(p *mapPolygon) int {
	mapPolygonStringCopy(&p.Name[0], memmap.PtrUint8(0x587000, 60364))
	for i := range p.Color {
		p.Color[i] = memmap.Uint8(0x587000, 60464+uintptr(i))
	}
	p.Level = memmap.Uint8(0x587000, 60490)
	p.Flags = 0
	p.Visited = 0
	return 0
}
func mapPolygonDefaultsWrite(p *mapPolygon) int {
	mapPolygonStringCopy(memmap.PtrUint8(0x587000, 60364), &p.Name[0])
	for i, v := range p.Color {
		*memmap.PtrUint8(0x587000, 60464+uintptr(i)) = v
	}
	*memmap.PtrUint8(0x587000, 60490) = p.Level
	p.Flags = 0
	p.Visited = 0
	return 0
}
func mapPolygonNew() *mapPolygon {
	i := mapPolygonFreeIndex()
	p := mapPolygonAt(i)
	p.Metadata = nil
	if noxflags.HasGame(0x200000) {
		p.Metadata = legacyCalloc(1, 256)
		if p.Metadata == nil {
			return nil
		}
	}
	p.ID = i
	alloc.StrCopy(p.Name[:], strconv.FormatInt(int64(int32(i)), 10))
	mapPolygonDefaultsRead(p)
	p.Active = 1
	return p
}
func mapPolygonBounds(p *mapPolygon) {
	ids := mapPolygonIDs(p)
	v := mapPolygonVertexAt(ids[0])
	x, y := floatToInt32(v.X), floatToInt32(v.Y)
	p.Bounds = [4]int32{x, y, x, y}
	for _, id := range ids[1:] {
		v := mapPolygonVertexAt(id)
		if !(float64(v.X) >= float64(p.Bounds[0])) {
			p.Bounds[0] = floatToInt32(v.X)
		} else if float64(v.X) > float64(p.Bounds[2]) {
			p.Bounds[2] = floatToInt32(v.X)
		}
		if !(float64(v.Y) >= float64(p.Bounds[1])) {
			p.Bounds[1] = floatToInt32(v.Y)
		} else if float64(v.Y) > float64(p.Bounds[3]) {
			p.Bounds[3] = floatToInt32(v.Y)
		}
	}
}
func mapPolygonConstruct() {
	count := memmap.Uint16(0x5D4594, 588072)
	if count < 3 {
		return
	}
	p := mapPolygonNew()
	if p == nil {
		return
	}
	p.Vertices = (*uint32)(legacyCalloc(uintptr(count), 4))
	p.Count = count
	copy(mapPolygonIDs(p), unsafe.Slice(memmap.PtrUint32(0x5D4594, 534820), int(count)))
	mapPolygonBounds(p)
	*memmap.PtrUint16(0x5D4594, 588072) = 0
}

func mapPolygonStringCopy(dst, src *byte) {
	s := alloc.GoString(src)
	b := unsafe.Slice(dst, len(s)+1)
	copy(b, s)
	b[len(s)] = 0
}
