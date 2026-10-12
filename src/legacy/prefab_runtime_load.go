package legacy

import (
	"io"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func prefabLoad(index int32) uint32 {
	if index < 0 || index >= int32((*prefabGlobal(prefabCount))) {
		return 0
	}
	GetServer().Nox_xxx_free503F40()
	*memmap.PtrUint32(0x5D4594, 1599572) = math.MaxUint32
	(*prefabGlobal(prefabScript)) = 0
	if prefabOpen((*prefabGlobal(prefabPath))) == 0 {
		return 0
	}
	defer prefabClose()
	if prefabSeek(index) == 0 {
		return 0
	}
	handle := (*prefabGlobal(prefabFile))
	prefabRawWord(handle)
	nameLen := int(prefabRawByte(handle))
	if nameLen > 63 {
		return 0
	}
	if prefabReadRaw(handle, make([]byte, nameLen)) != nameLen {
		return 0
	}
	prefabRawByte(handle)
	version := prefabRawByte(handle)
	prefabRawWord(handle)
	prefabRawWord(handle)
	if version > 1 {
		n := int32(prefabRawWord(handle))
		if n < 0 {
			return 0
		}
		if _, err := fileByHandle((*FILE)(mapRoomPointer(handle))).Seek(int64(n), io.SeekCurrent); err != nil {
			return 0
		}
	}
	if prefabRawWord(handle) != 0xcafedead {
		return 0
	}
	*memmap.PtrUint32(0x5D4594, 739980) = prefabRawWord(handle)
	*memmap.PtrUint32(0x5D4594, 739984) = prefabRawWord(handle)
	var bounds [8]uint32
	for _, i := range []int{0, 1, 6, 7, 2, 3, 4, 5} {
		bounds[i] = prefabRawWord(handle)
	}
	prefabSetBounds(&bounds)
	copy(unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, 1599500)), 8), bounds[:])
	var rect [4]int32
	geometryWallBounds(&bounds, &rect)
	cf := cryptfile.Global()
	cf.SetXOR(true)
	stream := prefabWire{}
	for {
		n := int(stream.byte(0))
		if stream.bad {
			return 0
		}
		if n == 0 {
			break
		}
		name := make([]byte, n)
		stream.data(name)
		stream.word(0)
		if stream.bad {
			return 0
		}
		id := prefabByteString(name)
		recognized, err := Nox_xxx_mapReadSection(cf, unsafe.Pointer(&bounds), id)
		if err != nil {
			mapLog.Println(err)
			return 0
		}
		if !recognized {
			obj := GetServer().S().NewObjectByTypeID(id)
			if obj == nil {
				return 0
			}
			if err := obj.CallXfer(unsafe.Pointer(&rect)); err != nil {
				GetServer().S().Objs.FreeObject(obj)
				return 0
			}
			objectXferPlace(obj, nil, unsafe.Pointer(&rect))
		}
	}
	cf.SetXOR(false)
	(*prefabGlobal(prefabLoaded)) = uint32(index)
	(*prefabGlobal(prefabPlaced)) = 0
	(*prefabGlobal(prefabSelected)) = uint32(index)
	return 1
}
func prefabInstantiate(point *types.Pointf) uint32 {
	var world types.Pointf
	if mapPaintTransform(point, &world) == 0 {
		return 0
	}
	if (*prefabGlobal(prefabLoaded)) != (*prefabGlobal(prefabSelected)) || (*prefabGlobal(prefabLoaded)) == math.MaxUint32 || (*prefabGlobal(prefabPlaced)) == 1 {
		if prefabLoad(int32((*prefabGlobal(prefabSelected)))) == 0 {
			return 0
		}
	}
	row := (*prefabGlobal(prefabMetadata)) + 76*(*prefabGlobal(prefabSelected))
	width, height := *prefabFloat(row, 64), *prefabFloat(row, 68)
	var bounds [8]uint32
	bounds[2], bounds[3] = uint32(int64(world.X)), uint32(int64(world.Y))
	corner := types.Ptf(width+point.X, height+point.Y)
	var transformed types.Pointf
	mapPaintTransform(&corner, &transformed)
	bounds[4], bounds[5] = uint32(int64(transformed.X)), uint32(int64(transformed.Y))
	corner = types.Ptf(width+point.X, point.Y)
	mapPaintTransform(&corner, &transformed)
	bounds[0], bounds[1] = uint32(int64(transformed.X)), uint32(int64(transformed.Y))
	corner = types.Ptf(point.X, height+point.Y)
	mapPaintTransform(&corner, &transformed)
	bounds[6], bounds[7] = uint32(int64(transformed.X)), uint32(int64(transformed.Y))
	prefabSetBounds(&bounds)
	var rect [4]int32
	geometryWallBounds(&bounds, &rect)
	gx, gy := memmap.Uint32(0x5D4594, 739980), memmap.Uint32(0x5D4594, 739984)
	*memmap.PtrUint32(0x5D4594, 1599484) = gx
	*memmap.PtrUint32(0x5D4594, 1599488) = gy
	*memmap.PtrFloat32(0x5D4594, 1599492) = float32(int32(gx * 23))
	*memmap.PtrFloat32(0x5D4594, 1599496) = float32(int32(gy * 23))
	dx := int32(int64(float64(world.X) - float64(memmap.Int32(0x5D4594, 1599508))))
	dy := int32(int64(float64(world.Y) - float64(memmap.Int32(0x5D4594, 1599512))))
	if prefabPlaceTiles(dx, dy) == 0 || prefabPlaceWalls(dx, dy) == 0 || prefabPlaceWaypoints(dx, dy) == 0 || prefabPlaceObjects(dx, dy) == 0 {
		return 0
	}
	s := GetServer().S()
	s.WPs.Sub_579D20()
	for wp := s.WPs.Pending; wp != nil; wp = wp.WpNext {
		wp.Flags |= 0x80000000
	}
	*mapPaintGlobal(paintObjectCounter) = prefabInteresting(&rect, *mapPaintGlobal(paintObjectCounter))
	if GetServer().Sub504720(uint32(dx), uint32(dy)) == 0 {
		return 0
	}
	for wp := s.WPs.Pending; wp != nil; wp = wp.WpNext {
		wp.Field1 = 0
	}
	for obj := s.Objs.Pending; obj != nil; obj = obj.ObjNext {
		obj.ScriptIDVal = 0
	}
	s.Nox_xxx_waypoint_5799C0()
	s.Objs.ObjectsClearPending()
	(*prefabGlobal(prefabPlaced)) = 1
	if (*prefabGlobal(prefabScript)) != 0 {
		*memmap.PtrUint32(0x973F18, 35880)++
		prefabAdjustScript((*prefabGlobal(prefabInstance)), dx, dy)
		a := alloc.GoString((*byte)(memmap.PtrOff(0x973F18, 36008)))
		b := alloc.GoString((*byte)(memmap.PtrOff(0x973F18, 38056)))
		c := alloc.GoString((*byte)(memmap.PtrOff(0x973F18, 30760)))
		if memmap.Uint32(0x5D4594, 1599580) != 0 {
			ifs.Remove(a)
			ifs.Rename(b, a)
			prefabCombineScripts(a, c, b)
		} else {
			*memmap.PtrUint32(0x5D4594, 1599580) = 1
			ifs.Rename(c, b)
		}
	}
	(*prefabGlobal(prefabInstance))++
	return 1
}
