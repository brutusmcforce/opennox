package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// MapPrefabFreeNode releases a prefab cache wrapper or raw payload with the same allocator
// used by its constructors and removal paths. Group references have their
// own tracked Go allocator and must not use this function.
func MapPrefabFreeNode(p unsafe.Pointer) { mapRoomRelease(p) }

// MapPrefabFreeCaches frees the prefab caches (objects, waypoints, tiles, walls). Their nodes link with 32-bit pointer
// words: data at 0, next at 4 (tiles: 16). Objects and waypoints are freed only when the prefab was not placed.
func MapPrefabFreeCaches(freeObject func(*server.Object)) {
	placed := *prefabGlobal(prefabPlaced) != 0
	for n, next := *prefabGlobal(prefabObjects), uint32(0); n != 0; n = next {
		next = *prefabWord(n, 4)
		if !placed {
			freeObject(populationObject(*prefabWord(n, 0)))
		}
		mapRoomRelease(mapRoomPointer(n))
	}
	for n, next := *prefabGlobal(prefabWaypoints), uint32(0); n != 0; n = next {
		next = *prefabWord(n, 4)
		if !placed {
			alloc.Free(prefabWaypoint(*prefabWord(n, 0)))
			*prefabWord(n, 0) = 0
		}
		mapRoomRelease(mapRoomPointer(n))
	}
	for n, next := *prefabGlobal(prefabTiles), uint32(0); n != 0; n = next {
		next = *prefabWord(n, 16)
		data := *prefabWord(n, 0)
		for s, next2 := *prefabWord(data, 16), uint32(0); s != 0; s = next2 {
			next2 = *prefabWord(s, 16)
			mapPaintSubtileFree((*[5]uint32)(mapRoomPointer(s)))
		}
		mapRoomRelease(mapRoomPointer(data))
		*prefabWord(n, 0) = 0
		mapRoomRelease(mapRoomPointer(n))
	}
	for n, next := *prefabGlobal(prefabWalls), uint32(0); n != 0; n = next {
		next = *prefabWord(n, 4)
		w := (*server.Wall)(mapRoomPointer(*prefabWord(n, 0)))
		// Placed secret data is moved out of the cache, even on partial failure.
		if w.Flags4&4 != 0 && w.Data != nil {
			mapRoomRelease(w.Data)
		}
		mapRoomRelease(unsafe.Pointer(w))
		*prefabWord(n, 0) = 0
		mapRoomRelease(mapRoomPointer(n))
	}
}
