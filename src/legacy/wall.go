package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Sub_526CA0                               func(a1 string) int
	Nox_xxx_mapSetWallInGlobalDir0pr1_5004D0 func()
	Nox_xxx_map_5004F0                       func()
	Sub_4FF990                               func(a1 uint32)
	Sub_5000B0                               func(a1 *server.Object) int
)

func nox_xxx_wallFlags(ind int) uint32 {
	return GetServer().S().Walls.DefByInd(ind).Flags32
}

func nox_xxx_mapDamageToWalls_534FC0(a1 *[4]int32, a2 unsafe.Pointer, a3 float32, a4, a5 int, a6 unsafe.Pointer) bool {
	rect := image.Rect(int(a1[0]), int(a1[1]), int(a1[2]), int(a1[3]))
	return bool(GetServer().Nox_xxx_mapDamageToWalls_534FC0(rect, *(*types.Pointf)(a2), float32(a3), a4, object.DamageType(a5), AsObjectP(a6)))
}

func Sub_5071C0() bool {
	return voteHead != nil
}

func Nox_xxx_math_509ED0(pos types.Pointf) int {
	cpos, free := alloc.New(types.Pointf{})
	defer free()
	*cpos = pos
	return int(int32(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(cpos))))))
}

func Nox_xxx_math_509EA0(a1 int) int {
	return int(int32(geometryDirection4Index(int32(a1))))
}
