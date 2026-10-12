package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// Copy only through the terminator, retaining the rest of the fixed record.
func worldCopyNarrow(dst, src unsafe.Pointer) {
	for i := 0; ; i++ {
		v := *controlByte(src, i)
		*controlByte(dst, i) = v
		if v == 0 {
			return
		}
	}
}
func worldCopyWide(dst, src unsafe.Pointer) {
	for i := 0; ; i += 2 {
		v := *controlHalf(src, i)
		*controlHalf(dst, i) = v
		if v == 0 {
			return
		}
	}
}
func worldWideEqual(a, b unsafe.Pointer) bool {
	for i := 0; ; i += 2 {
		v := *controlHalf(a, i)
		if v != *controlHalf(b, i) {
			return false
		}
		if v == 0 {
			return true
		}
	}
}
func worldAnkhRemember(pl unsafe.Pointer, a *server.Object) {
	for i := 0; i < 5; i++ {
		p := controlPtr(pl, 4796+4*i)
		if *p == nil {
			*p = a.CObj()
			return
		}
	}
}
func worldAnkhAlreadyAwarded(b *server.Object) {
	if worldCollideTicks()-*worldCollideClock() <= 1500 {
		return
	}
	worldCollideMessage(b, "objcoll.c:ExtraLifeAlreadyAwarded")
	worldCollideSound(925, b)
	*worldCollideClock() = worldCollideTicks()
}
func worldCollideAnkh(a, b *server.Object) {
	data := a.InitData
	if b == nil || b.ObjClass&4 == 0 {
		return
	}
	bd := b.UpdateData
	pl := controlPlayer(b)
	seen := false
	for i := 0; i < 5; i++ {
		if *controlPtr(pl, 4796+4*i) == a.CObj() {
			seen = true
			break
		}
	}
	core := GetServer().S()
	for i := 0; i < 64; i++ {
		record := unsafe.Add(data, 80*i)
		if core.Frame()-*equipmentWord(record, 76) > 240*uint32(core.TickRate()) {
			worldCopyWide(record, memmap.PtrOff(0x5D4594, 1568012))
			*controlByte(record, 51) = memmap.Uint8(0x5D4594, 1568016)
			*controlByte(record, 50) = 0
			*equipmentWord(record, 76) = 0
		}
		if *controlByte(record, 50) == *controlByte(pl, 2251) && worldWideEqual(record, unsafe.Add(pl, 2185)) && alloc.GoString(controlByte(record, 51)) == alloc.GoString(controlByte(pl, 2112)) {
			worldAnkhRemember(pl, a)
			worldAnkhAlreadyAwarded(b)
			return
		}
	}
	if seen {
		worldAnkhAlreadyAwarded(b)
		return
	}
	limit := uint32(floatToInt32(float32(core.Balance.Float("MaxExtraLives"))))
	if *equipmentWord(bd, 320) < limit {
		if item := core.NewObjectByTypeID("AnkhTradable"); item != nil {
			item.Pickup.Get()(b, item, 1, 0)
		}
		a.Field34 = core.Frame()
		worldCollideSound(1004, a)
		visibilityFXPoint(130, a.PosVec)
		worldCollideMessage(b, "objcoll.c:AwardExtraLife")
		// Re-read the player reference from the original update record after pickup.
		pl = *controlPtr(bd, 276)
		worldAnkhRemember(pl, a)
		index := controlByte(data, 5120)
		record := unsafe.Add(data, 80*int(*index))
		worldCopyWide(record, unsafe.Add(pl, 2185))
		*controlByte(record, 50) = *controlByte(pl, 2251)
		worldCopyNarrow(unsafe.Add(record, 51), unsafe.Add(pl, 2112))
		*equipmentWord(record, 76) = core.Frame()
		*index += 1
		if *index >= 64 {
			*index = 0
		}
	} else if worldCollideTicks()-*worldCollideClock() > 1500 {
		worldCollideMessage(b, "pickup.c:MaxTradableAnkhsReached")
		worldCollideSound(925, b)
		*worldCollideClock() = worldCollideTicks()
	}
}
