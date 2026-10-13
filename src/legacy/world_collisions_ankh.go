package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

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
func worldAnkhRemember(pl *server.Player, a *server.Object) {
	for i := 0; i < 5; i++ {
		p := &pl.Ankhs[i]
		if *p == nil {
			*p = a
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
	data := a.InitDataAnkh()
	if b == nil || b.ObjClass&4 == 0 {
		return
	}
	bd := b.UpdateData
	pl := controlPlayer(b)
	seen := false
	for i := 0; i < 5; i++ {
		if pl.Ankhs[i] == a {
			seen = true
			break
		}
	}
	core := GetServer().S()
	for i := 0; i < 64; i++ {
		record := &data.Records[i]
		if core.Frame()-record.Frame > 240*uint32(core.TickRate()) {
			worldCopyWide(unsafe.Pointer(&record.Name[0]), memmap.PtrOff(0x5D4594, 1568012))
			record.Serial[0] = memmap.Uint8(0x5D4594, 1568016)
			record.Class = 0
			record.Frame = 0
		}
		if record.Class == pl.InfoBytes[66] && worldWideEqual(unsafe.Pointer(&record.Name[0]), unsafe.Pointer(&pl.InfoBytes)) && alloc.GoString(&record.Serial[0]) == alloc.GoString(&pl.SerialBuf[0]) {
			worldAnkhRemember((*server.Player)(unsafe.Pointer(pl)), a)
			worldAnkhAlreadyAwarded(b)
			return
		}
	}
	if seen {
		worldAnkhAlreadyAwarded(b)
		return
	}
	limit := uint32(floatToInt32(float32(core.Balance.Float("MaxExtraLives"))))
	if layout.As[server.PlayerUpdateData](bd).Field80 < limit {
		if item := core.NewObjectByTypeID("AnkhTradable"); item != nil {
			item.Pickup.Get()(b, item, 1, 0)
		}
		a.Field34 = core.Frame()
		worldCollideSound(1004, a)
		visibilityFXPoint(130, a.PosVec)
		worldCollideMessage(b, "objcoll.c:AwardExtraLife")
		// Re-read the player reference from the original update record after pickup.
		pl = layout.As[server.PlayerUpdateData](bd).Player
		worldAnkhRemember((*server.Player)(unsafe.Pointer(pl)), a)
		index := &data.Next
		record := &data.Records[*index]
		worldCopyWide(unsafe.Pointer(&record.Name[0]), unsafe.Pointer(&pl.InfoBytes))
		record.Class = pl.InfoBytes[66]
		worldCopyNarrow(unsafe.Pointer(&record.Serial[0]), unsafe.Pointer(&pl.SerialBuf))
		record.Frame = core.Frame()
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
