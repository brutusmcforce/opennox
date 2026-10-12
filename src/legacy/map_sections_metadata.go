package legacy

import (
	"github.com/brutusmcforce/opennox/libs/wall"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"unsafe"
)

func mapSectionMetadata(kind int, region *[8]uint32) uint32 {
	// kind: 0 windows, 1 breakable walls, 2 secret walls.
	maxVersion := uint16(2)
	countOffset := uintptr(741336)
	flag := byte(64)
	if kind == 1 {
		maxVersion = 1
		countOffset = 741340
		flag = 8
	}
	if kind == 2 {
		countOffset = 741348
		flag = 4
	}
	version := int16(mapSectionWord(maxVersion))
	if version > int16(maxVersion) {
		return 0
	}
	count := memmap.PtrUint16(0x5D4594, countOffset)
	if !cryptfile.Global().ReadOnly() {
		each := func(fn func(*server.Wall)) {
			GetServer().S().Walls.EachWallXxx(func(w *server.Wall) bool {
				if byte(w.Flags4)&flag == 0 {
					return true
				}
				if region != nil {
					pt := [2]int32{23 * int32(w.X5), 23 * int32(w.Y6)}
					if geometryWallPoint(&pt, (*[8]int32)(unsafe.Pointer(region))) == 0 {
						return true
					}
				}
				fn(w)
				return true
			})
		}
		*count = 0
		each(func(*server.Wall) { *count++ })
		mapSectionIO(unsafe.Pointer(count), 2)
		each(func(w *server.Wall) {
			// Preserve the single coordinate operation: the legacy checksum depends on IO boundaries.
			point := [2]uint32{uint32(w.X5), uint32(w.Y6)}
			mapSectionIO(unsafe.Pointer(&point), 8)
			if kind == 2 {
				p := w.Data
				mapSectionIO(unsafe.Add(p, 16), 4)
				mapSectionIO(unsafe.Add(p, 20), 1)
				mapSectionIO(unsafe.Add(p, 21), 1)
				mapSectionIO(unsafe.Add(p, 22), 1)
				mapSectionIO(unsafe.Add(p, 24), 4)
				mapSectionIO(unsafe.Add(p, 28), 4)
			}
		})
		return 1
	}
	mapSectionIO(unsafe.Pointer(count), 2)
	if *count == 0 {
		return 1
	}
	generation := noxflags.HasGame(noxflags.GameFlag(0x400000))
	for i := int32(0); ; i++ {
		var data *[8]uint32
		var x, y int32
		if kind == 2 {
			data = (*[8]uint32)(mapRoomCalloc(1, 32))
			mapSectionIO(unsafe.Pointer(&data[1]), 8)
			mapSectionIO(unsafe.Pointer(&data[4]), 4)
			mapSectionIO(unsafe.Pointer(&data[5]), 1)
			if version >= 2 {
				mapSectionIO(unsafe.Add(unsafe.Pointer(data), 21), 1)
				mapSectionIO(unsafe.Add(unsafe.Pointer(data), 22), 1)
				mapSectionIO(unsafe.Pointer(&data[6]), 4)
				mapSectionIO(unsafe.Pointer(&data[7]), 4)
			}
			x, y = int32(data[1]), int32(data[2])
		} else {
			var point [2]uint32
			mapSectionIO(unsafe.Pointer(&point), 8)
			x, y = int32(point[0]), int32(point[1])
		}
		if region != nil {
			b := mapSectionRegionBounds(region)
			x += b[0] - memmap.Int32(0x5D4594, 739980)
			y += b[1] - memmap.Int32(0x5D4594, 739984)
		}
		if data != nil {
			data[1], data[2] = uint32(x), uint32(y)
		}
		var w *server.Wall
		if generation {
			if node := prefabWallFind(x, y); node != 0 {
				w = (*server.Wall)(mapRoomPointer(*prefabWord(node, 0)))
			}
		} else {
			w = GetServer().S().Walls.GetWallAtGrid(image.Pt(int(x), int(y)))
		}
		if w != nil {
			w.Flags4 |= wall.Flags(flag)
			switch kind {
			case 0:
				if version < 2 {
					w.Field2 = 0
				}
			case 1:
				index := memmap.PtrUint32(0x5D4594, 741344)
				w.Field10 = uint16(*index)
				*index++
				if !generation {
					GetServer().S().Walls.AddBreakable(w)
				}
			case 2:
				index := memmap.PtrUint32(0x5D4594, 741352)
				w.Field10 = uint16(*index)
				*index++
				w.Data = unsafe.Pointer(data)
				data[3] = mapRoomRaw(w.C())
				if byte(data[5]>>8) == 0 {
					if data[5]&8 != 0 {
						data[7] = 0xffffffff
						data[5] = data[5]&0xff0000ff | 3<<8 | 23<<16
					} else {
						data[7] = 0
						data[5] = data[5]&0xff0000ff | 1<<8
					}
				}
				if !generation {
					worldSecretInsert(w.Data)
				}
			}
		} else if data != nil {
			mapRoomRelease(unsafe.Pointer(data))
		}
		// The C format uses an unsigned initial nonempty check and signed loop bound.
		if i+1 >= int32(int16(*count)) {
			break
		}
	}
	return 1
}
