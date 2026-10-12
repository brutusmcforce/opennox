package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/libs/wall"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

const worldTileGridCapacity = 128

// worldTileCell retains the legacy 44-byte cell as eleven raw words.
// Words 5 and 10 hold unmanaged subtile-list addresses.
type worldTileCell [11]uint32




var worldTileGrid **worldTileCell
var worldTileDefinitions [176]server.TileDef
var worldTileDefinitionCount uint32
var worldSecretHead unsafe.Pointer

func worldGridAllocate() int {
	worldTileGrid = (**worldTileCell)(legacyCalloc(worldTileGridCapacity, 4))
	if worldTileGrid == nil {
		return 0
	}
	rows := unsafe.Slice(worldTileGrid, worldTileGridCapacity)
	for i := range rows {
		rows[i] = (*worldTileCell)(legacyCalloc(worldTileGridCapacity, 44))
		if rows[i] == nil {
			return 0
		}
	}
	return 1
}
func worldGridFreeRows() {
	for _, row := range unsafe.Slice(worldTileGrid, worldTileGridCapacity) {
		if row != nil {
			legacyFree(unsafe.Pointer(row))
		}
	}
	// The caller owns the outer pointer table; preserve that allocation and value.
}
func worldTileWater(p types.Pointf) bool {
	offsets := [5]uintptr{26516, 26520, 26524, 26528, 26532}
	if memmap.Int32(0x587000, 26520) == -1 {
		names := [5]string{"WaterNoTeleport", "WaterDeepNoTeleport", "WaterShallowNoTeleport", "WaterSwampDeepNoTeleport", "WaterSwampShallowNoTeleport"}
		for i := range worldTileDefinitions {
			name := worldTileDefinitions[i].Name()
			for j, want := range names {
				if name == want {
					*memmap.PtrUint32(0x587000, offsets[j]) = uint32(i)
					break
				}
			}
		}
	}
	tile := uint32(tileAtPoint(p))
	for _, off := range offsets {
		if tile == memmap.Uint32(0x587000, off) {
			return true
		}
	}
	return false
}
func worldDoorAttach(value unsafe.Pointer, x, y int) *server.Wall {
	w := GetServer().S().Walls.CreateAtGrid(image.Pt(x, y))
	if w != nil {
		w.Data = value
		w.Flags4 |= wall.Flags(0x10)
	}
	return w
}
func worldDrawableAttach(d *client.Drawable, x, y int) unsafe.Pointer {
	walls := &GetServer().S().Walls
	w := walls.GetWallAtGrid2(image.Pt(x, y))
	if w == nil {
		w = walls.CreateAtGrid(image.Pt(x, y))
		if w == nil {
			return nil
		}
	}
	w.Field32 = memmap.PtrWord(d.C())
	w.Flags4 |= wall.Flags(0x10)
	point := [2]int32{int32(d.PosVec.X), int32(d.PosVec.Y)}
	p := mapPolygonFind(&point, 0, false)
	if p == nil {
		p = mapPolygonFindEdge(&point, 0, 10)
	}
	value := byte(1)
	if p != nil {
		value = *(*byte)(unsafe.Add(unsafe.Pointer(p), int(layoutLegacyMapPolygon.Off(130)) /* Level */))
	}
	w.Field8 = (w.Field8 & 0xff00) | uint16(value)
	return unsafe.Pointer(p)
}
func worldSecretInsert(p unsafe.Pointer) unsafe.Pointer {
	*(*unsafe.Pointer)(p) = worldSecretHead // the record starts with its next pointer
	worldSecretHead = p
	return p
}
func worldSecretNext(p unsafe.Pointer) unsafe.Pointer {
	if p == nil {
		return nil
	}
	return *(*unsafe.Pointer)(p)
}
func worldSecretFind(id int16) *server.Wall {
	for p := worldSecretHead; p != nil; p = worldSecretNext(p) {
		w := *(**server.Wall)(unsafe.Add(p, 12))
		if int32(w.Field10) == int32(id) {
			return w
		}
	}
	return nil
}
func worldSecretRemove(p unsafe.Pointer) unsafe.Pointer {
	var prev unsafe.Pointer
	for cur := worldSecretHead; cur != nil; cur = worldSecretNext(cur) {
		if cur == p {
			next := worldSecretNext(cur)
			if prev == nil {
				worldSecretHead = next
			} else {
				*(*unsafe.Pointer)(prev) = next
			}
			legacyFree(p)
			return cur
		}
		prev = cur
	}
	return nil
}
func worldSecretClear() unsafe.Pointer {
	for p := worldSecretHead; p != nil; {
		next := worldSecretNext(p)
		legacyFree(p)
		p = next
	}
	worldSecretHead = nil
	return nil
}
