package legacy

import (
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/strman"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func doubleToInt32(v float64) int32 {
	bits := math.Float64bits(v)
	exp := int((bits>>52)&2047) - 1023
	if exp < 0 {
		return 0
	}
	if exp >= 31 {
		return -2147483648
	}
	magnitude := (uint64(1)<<52 | bits&0xfffffffffffff) >> uint(52-exp)
	if bits>>63 != 0 {
		return -int32(magnitude)
	}
	return int32(magnitude)
}
func runtimeObserverCount() int {
	if !noxflags.HasGame(0x8000) {
		return 0
	}
	count := 0
	players := &GetServer().S().Players
	for p := players.First(); p != nil; p = players.Next(p) {
		if p.Field3680&0x21 == 1 && p.PlayerInd != 31 {
			count++
		}
	}
	return count
}
func runtimeMeterWave() {
	scale := math.Float64frombits(uint64(qword_581450_9552))
	for i := 0; i < 320; i++ {
		angle := float64(i+192) * memmap.Float64(0x581450, 9768) * memmap.Float64(0x581450, 9760)
		*memmap.PtrUint32(0x5D4594, 1309840+4*uintptr(i)) = uint32(effectsTruncWord(math.Sin(angle) * scale))
	}
}
func runtimePauseActor() *server.Object {
	return (*server.Object)(unsafe.Pointer(uintptr(dword_5d4594_2523780)))
}
func runtimePauseEffect() *server.Object {
	return (*server.Object)(unsafe.Pointer(uintptr(dword_5d4594_2523776)))
}
func runtimePauseStart(actor *server.Object, kind int32) {
	if dword_5d4594_2523804 == 1 || noxflags.HasGame(noxflags.GamePause) {
		return
	}
	if actor != nil {
		dword_5d4594_2523780 = uint32(uintptr(actor.CObj()))
	} else {
		actor = runtimePauseActor()
	}
	effect := runtimePauseEffect()
	if kind == 0 || kind == 1 {
		name := "LevelUp"
		if kind == 1 {
			name = "OblivionUp"
		}
		effect = GetServer().S().NewObjectByTypeID(name)
		actor = runtimePauseActor()
		dword_5d4594_2523776 = uint32(uintptr(effect.CObj()))
	}
	if effect != nil {
		if actor != nil {
			GetServer().CreateObjectAt(effect, nil, actor.PosVec)
		} else {
			GetServer().S().Objs.FreeObject(effect)
			dword_5d4594_2523776 = 0
		}
		actor = runtimePauseActor()
	}
	if (kind == 0 || kind == 1) && actor != nil {
		visibilityFXPoint(154, actor.PosVec)
		actor = runtimePauseActor()
	}
	if kind == 0 && actor != nil {
		inventorySound(902, actor, 2, int(*(*uint32)(unsafe.Add(actor.CObj(), 36))))
		text := GetServer().S().Strings().GetStringInFile(strman.ID("expLevel.c:LevelUP"), `C:\NoxPost\src\common\GameMech\PauseFX.c`)
		Nox_xxx_netSendLineMessage_4D9EB0(runtimePauseActor(), text)
		actor = runtimePauseActor()
	}
	if actor != nil {
		ud := actor.UpdateData
		if Nox_xxx_playerSetState_4FA020(actor, 30) {
			*(*byte)(unsafe.Add(ud, 236)) = 4
		}
	}
	delay := uint32(0)
	if kind == 0 || kind == 1 {
		delay = 5000
	}
	*memmap.PtrUint32(0x5D4594, 2523796) = delay
	*memmap.PtrUint32(0x5D4594, 2523800) = 0
	*memmap.PtrUint32(0x5D4594, 2523772) = uint32(kind)
	dword_5d4594_2523804 = 1
	Sub_413A00(1)
	*memmap.PtrUint64(0x5D4594, 2523788) = uint64(uint32(PlatformTicks()))
}
func runtimePauseStop() {
	if dword_5d4594_2523804 == 0 {
		return
	}
	actor := runtimePauseActor()
	kind := memmap.Uint32(0x5D4594, 2523772)
	if actor != nil && (kind == 0 || kind == 1) {
		visibilityFXPoint(154, actor.PosVec)
		actor = runtimePauseActor()
	}
	if effect := runtimePauseEffect(); effect != nil {
		GetServer().DelayedDelete(effect)
		actor = runtimePauseActor()
	}
	dword_5d4594_2523776 = 0
	if actor != nil {
		Nox_xxx_playerSetState_4FA020(actor, 13)
	}
	dword_5d4594_2523780 = 0
	if Sub_45D9B0() == 0 {
		Sub_413A00(0)
	}
	dword_5d4594_2523804 = 0
}
func runtimeHostOwnership() {
	s := GetServer().S()
	host := s.Players.ByInd(31)
	if host == nil || host.PlayerUnit == nil {
		return
	}
	cache := memmap.PtrUint32(0x5D4594, 1563124)
	if *cache == 0 {
		*cache = uint32(s.Types.IndByID("SaveGameLocation"))
	}
	for marker := s.Objs.List; marker != nil; marker = marker.ObjNext {
		if uint32(marker.TypeInd) != *cache {
			continue
		}
		for child := marker.Field129; child != nil; {
			next := child.Field128
			s.ObjSetOwner(host.PlayerUnit, child)
			if child.Class().Has(object.ClassMonster) && *(*byte)(unsafe.Add(child.UpdateData, 1440))&0x80 != 0 {
				child.ObjSubClass |= 0x80
				gameplayReportAcquireCreature(int(host.PlayerInd), child)
				s.Players.Nox_xxx_netMarkMinimapObject_417190(host.PlayerIndex(), child, 1)
			}
			child = next
		}
		*(*uint32)(unsafe.Add(marker.CObj(), 44)) = 0
		GetServer().DelayedDelete(marker)
		return
	}
}
