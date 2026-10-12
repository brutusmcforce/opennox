//go:build 386

package client

// Layout assertions: these structs keep the layout of the original 32-bit binary.
// On other architectures pointers are larger and the layout differs.

import (
	"unsafe"
)

var _ = [1]struct{}{}[20-unsafe.Sizeof(DrawableUnion{})]
var _ = [1]struct{}{}[unsafe.Sizeof(DrawableUnionMonster{})-unsafe.Sizeof(DrawableUnion{})]
var _ = [1]struct{}{}[unsafe.Sizeof(DrawableUnionItem{})-unsafe.Sizeof(DrawableUnion{})]
var _ = [1]struct{}{}[unsafe.Sizeof(DrawableUnionEffect{})-unsafe.Sizeof(DrawableUnion{})]
var _ = [1]struct{}{}[unsafe.Sizeof(DrawableUnionDoor{})-unsafe.Sizeof(DrawableUnion{})]
var _ = [1]struct{}{}[64-unsafe.Offsetof(DrawableFX{}.Next)]
var _ = [1]struct{}{}[772-unsafe.Sizeof(MonsterDrawData{})]
var _ = [1]struct{}{}[40-unsafe.Sizeof(PlayerEquipAnimation{})]
var _ = [1]struct{}{}[264-unsafe.Sizeof(PlayerAnimation{})]
var _ = [1]struct{}{}[14524-unsafe.Sizeof(PlayerDrawData{})]
var _ = [1]struct{}{}[48-unsafe.Sizeof(AnimationVector{})]
