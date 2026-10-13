package server

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

// WandUseData is the use data of wands and staffs (WandUse, WandCastUse).
type WandUseData struct {
	Kind       uint32   // 0, 0: 0 for WandUse, 1 for WandCastUse
	MissileBuf [80]byte // 1, 4: missile type ID (WandUse)
	MissileInd uint32   // 21, 84: cached missile type index (WandUse)
	Sound      uint32   // 22, 88: shot sound (WandUse)
	Spell      uint32   // 23, 92: spell ID (WandCastUse)
	Flags      uint32   // 24, 96: 1 multi shot, 2 melee swing, 4 sustained spell is active
	Delay      uint32   // 25, 100: frames between uses
	LastUse    uint32   // 26, 104: frame of the last use
	Charges    byte     // 27, 108
	MaxCharges byte     // 27, 109: 0 if the wand has no charges
	Field110   [2]byte  // 27, 110
	Percent    uint32   // 28, 112: charge percent
}

var _ = [1]struct{}{}[116-unsafe.Sizeof(WandUseData{})]

func (d *WandUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

func (d *WandUseData) Missile() string {
	return alloc.GoStringS(d.MissileBuf[:])
}

func (p UseDataPtr) AsWand() *WandUseData {
	return useDataPtrAs[WandUseData](p)
}

// UseDataWand returns the use data of wands and staffs (WandUse, WandCastUse).
func (obj *Object) UseDataWand() *WandUseData {
	// TODO: verify this conversion by checking ObjectType
	return obj.UseData.AsWand()
}
