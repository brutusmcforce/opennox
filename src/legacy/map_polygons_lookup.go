package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func mapPolygonNearest(x, y, distance float32) *mapPolygonVertex {
	var found *mapPolygonVertex
	for i := uint32(1); i < mapPolygonVertexNext; i++ {
		p := mapPolygonVertexAt(i)
		if p.Active == 0 {
			continue
		}
		dy := float64(float32(p.Y - y))
		dx := float32(p.X - x)
		// The Y square and sum use double precision; the X square uses float.
		// A newly selected distance is narrowed back to the float threshold.
		squared := dy*dy + float64(float32(dx*dx))
		if squared < float64(distance) {
			distance = float32(squared)
			found = p
		}
	}
	return found
}
func mapPolygonSelected(p *mapPolygonVertex) bool {
	for _, id := range unsafe.Slice(memmap.PtrUint32(0x5D4594, 534820), int(memmap.Uint16(0x5D4594, 588072))) {
		if id == p.ID {
			return true
		}
	}
	return false
}
func mapPolygonSelect(x, y float32) *mapPolygonVertex {
	point := [2]int32{floatToInt32(x), floatToInt32(y)}
	count := memmap.PtrUint16(0x5D4594, 588072)
	if noxflags.HasGame(0x200000) {
		if p := mapPolygonNearest(x, y, 900); p != nil {
			if !mapPolygonSelected(p) {
				*memmap.PtrUint32(0x5D4594, 534820+uintptr(*count)*4) = p.ID
				*count++
			}
			return nil
		}
		if mapPolygonFind(&point, 0, false) != nil {
			return nil
		}
	}
	p := mapPolygonVertexSet(x, y, mapPolygonFreeVertex(), 0)
	*memmap.PtrUint32(0x5D4594, 534820+uintptr(*count)*4) = p.ID
	*count++
	return p
}
func mapPolygonSegmentStart(p *mapPolygon) [4]int32 {
	v := mapPolygonVertexAt(*p.Vertices)
	return [4]int32{floatToInt32(v.X), floatToInt32(v.Y), 0, 0}
}
func mapPolygonSegmentNext(p *mapPolygon, segment *[4]int32, i uint32) {
	id := mapPolygonIDs(p)[i%uint32(p.Count)]
	v := mapPolygonVertexAt(id)
	off := 0
	if i&1 != 0 {
		off = 2
	}
	segment[off] = floatToInt32(v.X)
	segment[off+1] = floatToInt32(v.Y)
}
func mapPolygonContains(p *mapPolygon, point *[2]int32) bool {
	if p == nil {
		return false
	}
	ray := [4]int32{point[0], point[1], 0, 0}
	phase := memmap.PtrUint32(0x5D4594, 588080)
	if *phase&2 != 0 {
		ray[2] = 5888
	}
	*phase++
	if *phase&2 != 0 {
		ray[3] = 5888
	}
	segment := mapPolygonSegmentStart(p)
	var crossings uint8
	for i := uint16(1); i <= p.Count; i++ {
		mapPolygonSegmentNext(p, &segment, uint32(i))
		if int32(geometrySegments((*[4]int32)(unsafe.Pointer(unsafe.Pointer(&ray))), (*[4]int32)(unsafe.Pointer(unsafe.Pointer(&segment))))) != 0 {
			crossings++
		}
	}
	return crossings&1 != 0
}
func mapPolygonAdmits(p *mapPolygon, point *[2]int32, scripts bool) bool {
	if p.Active == 0 || scripts && p.Enter.Func == -1 && p.Leave.Func == -1 {
		return false
	}
	return int32(geometryRectInt((*[2]int32)(unsafe.Pointer(unsafe.Pointer(point))), (*[4]int32)(unsafe.Pointer(unsafe.Pointer(&p.Bounds))))) != 0 && mapPolygonContains(p, point)
}
func mapPolygonFind(point *[2]int32, cached uint32, scripts bool) *mapPolygon {
	if cached != 0 && cached != mapPolygonUnset {
		p := mapPolygonAt(cached)
		if mapPolygonAdmits(p, point, scripts) {
			return p
		}
	}
	for i := uint32(1); i < mapPolygonNext; i++ {
		p := mapPolygonAt(i)
		if p.ID != cached && mapPolygonAdmits(p, point, scripts) {
			return p
		}
	}
	return nil
}
func mapPolygonIndex(point *[2]int32, cached uint32) uint32 {
	if p := mapPolygonFind(point, cached, false); p != nil {
		return p.ID
	}
	return 0
}
func mapPolygonEdge(p *mapPolygon, point *[2]int32, distance float32) bool {
	segment := mapPolygonSegmentStart(p)
	for i := uint32(1); i <= uint32(p.Count); i++ {
		mapPolygonSegmentNext(p, &segment, i)
		if int32(geometryProjectEdge((*[2]int32)(unsafe.Pointer(point)), (*[4]int32)(unsafe.Pointer(unsafe.Pointer(&segment))), float32(distance))) != 0 {
			return true
		}
	}
	return false
}
func mapPolygonFindEdge(point *[2]int32, cached uint32, distance float32) *mapPolygon {
	if cached != 0 && cached != mapPolygonUnset {
		p := mapPolygonAt(cached)
		if p.Active != 0 && mapPolygonEdge(p, point, distance) {
			return p
		}
	}
	for i := uint32(1); i < mapPolygonNext; i++ {
		p := mapPolygonAt(i)
		if p.Active != 0 && p.ID != cached && mapPolygonEdge(p, point, distance) {
			return p
		}
	}
	return nil
}
