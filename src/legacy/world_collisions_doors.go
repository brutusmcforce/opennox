package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

// Preserve the legacy platform adapter's uint32 result before the uint64 clock math.
func worldCollideTicks() uint64  { return uint64(uint32(PlatformTicks())) }
func worldCollideClock() *uint64 { return (*uint64)(unsafe.Pointer(&qword_5d4594_1567940)) }
func worldCollideDoor(a, b *server.Object) {
	d := a.UpdateData
	if b == nil || *equipmentWord(d, 12) != *equipmentWord(d, 4) {
		return
	}
	if owner := a.ObjOwner; owner != nil {
		if a.Field34 <= GetServer().S().Frame() {
			a.ObjOwner = nil
		} else if owner != b {
			if worldCollideTicks()-*worldCollideClock() > 1500 {
				if a.ObjSubClass&4 != 0 {
					worldCollideSound(244, a)
					worldCollideMessage(b, "objcoll.c:GateLockedMagic")
				} else {
					worldCollideSound(240, a)
					worldCollideMessage(b, "objcoll.c:DoorLockedMagic")
				}
				*worldCollideClock() = worldCollideTicks()
			}
			return
		}
	}
	keyType := *controlByte(d, 1)
	if keyType == 0 {
		return
	}
	if keyType == 5 {
		if worldCollideTicks()-*worldCollideClock() > 1500 {
			if a.ObjSubClass&4 != 0 {
				worldCollideSound(244, a)
				worldCollideMessage(b, "objcoll.c:GateLockedMechanism")
			} else {
				worldCollideSound(240, a)
				worldCollideMessage(b, "objcoll.c:DoorLockedMechanism")
			}
			*worldCollideClock() = worldCollideTicks()
		}
		return
	}
	key := GetServer().S().DoorCheckKey(b, a)
	if key == nil {
		if worldCollideTicks()-*worldCollideClock() > 1500 {
			if a.ObjSubClass&4 != 0 {
				worldCollideSound(244, a)
				controlLockedDoor(b, internCStr("objcoll.c:GateLockedKey"), keyType)
			} else {
				worldCollideSound(240, a)
				controlLockedDoor(b, internCStr("objcoll.c:DoorLockedKey"), keyType)
			}
			*worldCollideClock() = worldCollideTicks()
		}
		return
	}
	x, y := int32(*equipmentWord(d, 16)), int32(*equipmentWord(d, 20))
	*controlByte(d, 1) = 0
	cx, cy := float64(23*x), float64(23*y)
	rect := types.Rectf{Min: types.Ptf(float32(cx-34), float32(cy-34)), Max: types.Ptf(float32(cx+34), float32(cy+34))}
	var adjacent [2]int32
	switch *equipmentWord(d, 4) {
	case 0:
		adjacent = [2]int32{x - 1, y - 1}
	case 8:
		adjacent = [2]int32{x + 1, y - 1}
	case 16:
		adjacent = [2]int32{x + 1, y + 1}
	case 24:
		adjacent = [2]int32{x - 1, y + 1}
	}
	if noxflags.HasGame(noxflags.GameModeQuest) {
		stateDoorNotify(a)
		questRuntimeSetSoulFrame(GetServer().S().Frame())
	}
	GetServer().S().Map.EachObjInRect(rect, func(u *server.Object) bool { stateCloseDoor(u, unsafe.Pointer(&adjacent)); return true })
	worldCollideSound(234, a)
	if holder := key.InvHolder; holder != nil && holder != b && holder.ObjClass&4 != 0 && noxflags.HasGame(noxflags.GameModeQuest) && GetServer().S().Doors.Sub_4D72C0() {
		worldCollideMessage(holder, "GeneralPrint:KeyShared1")
	}
	GetServer().DelayedDelete(key)
}
func worldCollideChest(a, b *server.Object) {
	if b == nil || b.ObjClass&4 == 0 || a.ObjFlags&0x8000 != 0 {
		return
	}
	if noxflags.HasGame(noxflags.GameModeQuest) && a.ObjSubClass&0xf00 != 0 {
		found := false
		for it := b.InvFirstItem; it != nil; it = it.InvNextItem {
			if it.ObjClass&0x40 != 0 && GetServer().S().Types.ByInd(int(it.TypeInd)).ID() == "SilverKey" {
				GetServer().DelayedDelete(it)
				worldCollideSound(234, a)
				found = true
				break
			}
		}
		if !found {
			if worldCollideTicks()-*worldCollideClock() > 1500 {
				worldCollideSound(1012, a)
				worldCollideMessage(b, "objcoll.c:ChestLockedSilver")
				*worldCollideClock() = worldCollideTicks()
			}
			return
		}
	}
	if a.Death != nil {
		a.CallDeath()
	}
	inventoryChest(a, b)
	inventoryDropAll(a)
}
