package legacy

import (
	"io"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// Word-sized addresses remain required by the live 386 map-file ABI. Payloads
// use their subsystem allocator; Go owns the cache state and algorithms.
const (
	prefabMetadata = iota
	prefabCount
	prefabLoaded
	prefabPlaced
	prefabSelected
	prefabObjects
	prefabWalls
	prefabTiles
	prefabWaypoints
	prefabPath
	prefabAlternate
	prefabIntro
	prefabScript
	prefabInstance
	prefabLastWaypoint
	prefabFile
)

var prefabState [16]uint32

func prefabWord(p uint32, off int) *uint32     { return (*uint32)(unsafe.Add(mapRoomPointer(p), off)) }
func prefabFloat(p uint32, off int) *float32   { return (*float32)(unsafe.Add(mapRoomPointer(p), off)) }
func prefabPoint(p uint32) *types.Pointf       { return (*types.Pointf)(mapRoomPointer(p)) }
func prefabWaypoint(p uint32) *server.Waypoint { return (*server.Waypoint)(mapRoomPointer(p)) }
func prefabMetadataAt(index int32) uint32 {
	if index < 0 || index > int32((*prefabGlobal(prefabCount))) {
		return 0
	}
	return (*prefabGlobal(prefabMetadata)) + uint32(index)*76
}
func prefabDimension(index int32, off int) float64 {
	if index < 0 || index >= int32((*prefabGlobal(prefabCount))) {
		return -1
	}
	return float64(*prefabFloat((*prefabGlobal(prefabMetadata))+uint32(index)*76, off))
}
func prefabFindName(name uint32) uint32 {
	for i := int32(0); i < int32((*prefabGlobal(prefabCount))); i++ {
		if prefabCompareNames(name, (*prefabGlobal(prefabMetadata))+uint32(i)*76) == 0 {
			return uint32(i)
		}
	}
	return math.MaxUint32
}
func prefabSetPath(source uint32, alternate bool) uint32 {
	slot, sentinel := prefabPath, uintptr(1599608)
	if alternate {
		slot, sentinel = prefabAlternate, 1599612
	} else {
		prefabClose()
	}
	dst := unsafe.Slice((*byte)(mapRoomPointer((*prefabGlobal(slot)))), 2048)
	if source == 0 {
		dst[0] = memmap.Uint8(0x5D4594, sentinel)
		return 0
	}
	// C strncpy pads the first 2047 bytes and preserves the last byte.
	ended := false
	for i := 0; i < 2047; i++ {
		var c byte
		if !ended {
			c = *(*byte)(unsafe.Add(mapRoomPointer(source), i))
			ended = c == 0
		}
		dst[i] = c
	}
	return 1
}
func prefabOpen(path uint32) uint32 {
	if (*prefabGlobal(prefabFile)) == 0 {
		if err := cryptfile.OpenGlobal(alloc.GoString((*byte)(mapRoomPointer(path))), cryptfile.ReadOnly, -1); err != nil {
			return 0
		}
		(*prefabGlobal(prefabFile)) = mapRoomRaw(unsafe.Pointer(NewFileHandle(cryptfile.Global().File.File)))
	}
	_, _ = fileByHandle((*FILE)(mapRoomPointer((*prefabGlobal(prefabFile))))).Seek(0, io.SeekStart)
	return 1
}
func prefabClose() uint32 {
	old := (*prefabGlobal(prefabFile))
	if old != 0 {
		cryptfile.Close()
		(*prefabGlobal(prefabFile)) = 0
	}
	return old
}
func prefabSeek(index int32) uint32 {
	handle := (*prefabGlobal(prefabFile))
	if handle == 0 || index < 0 || index >= int32((*prefabGlobal(prefabCount))) {
		return 0
	}
	off := int32(*prefabWord((*prefabGlobal(prefabMetadata))+uint32(index)*76, 72))
	_, _ = fileByHandle((*FILE)(mapRoomPointer(handle))).Seek(int64(off), io.SeekStart)
	return handle
}
func prefabSelect(index int32) uint32 {
	if index < 0 || index >= int32((*prefabGlobal(prefabCount))) {
		return 0
	}
	(*prefabGlobal(prefabSelected)) = uint32(index)
	return uint32(bool2int(prefabLoad(index) == 0))
}
func prefabResetWaypoint()                   { (*prefabGlobal(prefabLastWaypoint)) = 0 }
func prefabSetWaypointKind(kind byte) uint32 { *memmap.PtrUint8(0x973F18, 35972) = kind; return 1 }
func prefabCreateWaypoint(point *types.Pointf) uint32 {
	var pos types.Pointf
	if mapPaintTransform(point, &pos) == 0 {
		return 0
	}
	wp := GetServer().S().NewWaypoint(pos)
	if wp == nil {
		return 0
	}
	if previous := (*prefabGlobal(prefabLastWaypoint)); previous != 0 && memmap.Uint32(0x973F18, 35976) == 1 {
		kind := int8(memmap.Uint8(0x973F18, 35972))
		appendWaypointLink(prefabWaypoint(previous), wp, kind)
		appendWaypointLink(wp, prefabWaypoint(previous), kind)
	}
	raw := mapRoomRaw(unsafe.Pointer(wp))
	(*prefabGlobal(prefabLastWaypoint)) = raw
	return raw
}
func prefabFindWaypoint(point *types.Pointf) uint32 {
	var pos types.Pointf
	if mapPaintTransform(point, &pos) == 0 {
		return 0
	}
	return mapRoomRaw(unsafe.Pointer(GetServer().S().WPs.Sub_579AD0(pos)))
}
func prefabConnectWaypoint(a, b *types.Pointf) uint32 {
	if a == nil || b == nil {
		return 0
	}
	x := prefabFindWaypoint(a)
	if x == 0 {
		return 0
	}
	y := prefabFindWaypoint(b)
	if y == 0 {
		return 0
	}
	return uint32(bool2int(appendWaypointLink(prefabWaypoint(x), prefabWaypoint(y), int8(memmap.Uint8(0x973F18, 35972)))))
}
func prefabRelativePosition(raw uint32, out *types.Pointf) uint32 {
	if (*prefabGlobal(prefabLoaded)) != (*prefabGlobal(prefabSelected)) || (*prefabGlobal(prefabLoaded)) == math.MaxUint32 || (*prefabGlobal(prefabPlaced)) == 1 {
		return 0
	}
	origin := types.Ptf(float32(memmap.Int32(0x5D4594, 1599508)), float32(memmap.Int32(0x5D4594, 1599512)))
	var base, pos types.Pointf
	geometryMapCoordinates(&origin, &base)
	geometryMapCoordinates(&populationObject(raw).PosVec, &pos)
	out.X, out.Y = pos.X-base.X, pos.Y-base.Y
	return 1
}
