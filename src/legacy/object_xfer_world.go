package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func objectXferInventory(version uint16, u *server.Object, count int32) int {
	r := objectXferStream{cryptfile.Global()}
	for i := int32(0); i < count; i++ {
		var typ uint16
		if version < 60 {
			n := int(r.byte(0))
			var name [256]byte
			r.raw(unsafe.Pointer(&name[0]), n)
			typ = uint16(GetServer().S().Types.IndByID(alloc.GoStringS(name[:])))
		} else {
			typ = uint16(Nox_xxx_objectTOCgetTT(r.short(0)))
		}
		if typ == 0 {
			return 0
		}
		r.cf.ReadAlignedU32() // Record length is consumed, not used to skip a rejected child.
		child := GetServer().S().NewObjectByTypeInd(int(typ))
		if child == nil {
			return 0
		}
		if child.CallXfer(nil) != nil {
			GetServer().S().Objs.FreeObject(child)
			return 0
		}
		child.Field125 = nil
		child.InvNextItem = u.InvFirstItem
		if u.InvFirstItem != nil {
			u.InvFirstItem.Field125 = child
		}
		u.InvFirstItem = child
		child.InvHolder = u
	}
	return 1
}
func objectXferPlace(u *server.Object, owner, offset unsafe.Pointer) int {
	s := GetServer().S()
	admitted := noxflags.HasGame(0x400000) || s.Types.ByInd(int(u.TypeInd)).Allowed()
	if admitted {
		if offset != nil {
			p := (*[2]int32)(offset)
			x := int32(uint32(memmap.Int32(0x5D4594, 739980)) * 23)
			y := int32(uint32(memmap.Int32(0x5D4594, 739984)) * 23)
			u.PosVec.X = float32(float64(u.PosVec.X) - float64(x) + float64(p[0]) - 11)
			u.PosVec.Y = float32(float64(u.PosVec.Y) - float64(y) + float64(p[1]) - 11)
		}
		if noxflags.HasGame(0x400000) {
			nox_xxx_unitAddToList_5048A0(int32(memmap.PtrWord(u.CObj())))
			return 1
		}
		if noxflags.HasGame(0x200000) || Sub_4E3AD0(int(u.TypeInd)) != 0 {
			var own server.Obj
			if owner != nil {
				own = (*server.Object)(owner)
			}
			GetServer().CreateObjectAt(u, own, u.PosVec)
			return 1
		}
	}
	for it := u.InvFirstItem; it != nil; {
		next := it.InvNextItem
		s.Objs.FreeObject(it)
		it = next
	}
	// The baseline correction clears the inventory before the parent destructor.
	u.InvFirstItem = nil
	s.Objs.FreeObject(u)
	return 0
}
func objectXferLight(u *server.Object) int {
	r, v, saved, ok := objectXferStart(u, 60)
	if !ok {
		return 0
	}
	var light [140]byte
	if !r.read() {
		objs := &GetClient().Cli().Objs
		var d *client.Drawable
		if objectXferEditor() {
			for it := objs.FirstList1(); it != nil; it = it.NextPtr {
				if it.NetCode32 == u.Extent {
					d = it
					break
				}
			}
		} else {
			if uint32(u.ObjClass)&0x20400000 != 0 {
				d = objs.ByNetCodeStatic(int(u.Extent))
			} else {
				d = objs.ByNetCodeDynamic(int(u.NetCode))
			}
			// The original aborted without a drawable, which lights created at runtime (e.g. by scripts) don't
			// have: saving their map crashed the game. Their record is saved without light values instead.
		}
		if d != nil {
			// the saved record keeps the 386 layout of the drawable from LightFlags on
			layoutClientDrawable.Runs(136, uintptr(len(light)), func(o386, o, n uintptr) {
				copy(light[o386-136:], unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(d), o)), n))
			})
		}
	}
	p := unsafe.Pointer(&light[0])
	// Keep each field operation in the original stream order.
	for _, f := range [][2]int{{0, 4}, {4, 4}, {8, 4}, {12, 4}, {16, 12}, {28, 2}, {30, 2}, {32, 4}} {
		r.raw(unsafe.Add(p, f[0]), f[1])
	}
	if v >= 2 {
		for _, f := range [][2]int{{40, 2}, {42, 48}, {90, 16}, {106, 16}, {122, 2}, {124, 2}, {126, 2}, {128, 4}, {134, 2}, {136, 2}, {138, 1}} {
			r.raw(unsafe.Add(p, f[0]), f[1])
		}
		if v > 40 {
			if v >= 42 {
				r.raw(unsafe.Add(p, 36), 4)
			} else {
				*objectXferWord(p, 36) = uint32(r.byte(0))
			}
		} else if r.read() {
			*objectXferWord(p, 36) = 0
		}
	} else {
		for _, off := range []int{40, 122, 124, 126, 134} {
			*(*uint16)(unsafe.Add(p, off)) = 0
		}
		*objectXferWord(p, 128) = 0
		light[138] = 128
		if r.read() && (*(*float32)(unsafe.Add(p, 4)) > 63 || float64(int32(*objectXferWord(p, 12)))*memmap.Float64(0x581450, 9752) > memmap.Float64(0x581450, 9744)) {
			particleLightIntensityRaw(p, 63)
		}
	}
	if r.read() && objectXferEditor() {
		copy(u.EditorData().Light[:], light[:])
	}
	return objectXferFinish(r, u, v, saved)
}
