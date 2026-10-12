package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"math/bits"
	"time"
	"unsafe"
)

var statisticsPlayers, statisticsEvents, statisticsArrayCount uint32

// The session initializer is disabled. Preserve the elapsed-time origin used by
// the live flush path, including explicit fixture ownership of that state.
var statisticsStart uint32
var statisticsTime = func() uint32 { return uint32(time.Now().Unix()) }

func statisticsI32(p unsafe.Pointer, n int) *int32  { return (*int32)(unsafe.Add(p, n)) }
func statisticsU32(p unsafe.Pointer, n int) *uint32 { return (*uint32)(unsafe.Add(p, n)) }
func statisticsString(p unsafe.Pointer) string      { return alloc.GoString((*byte)(p)) }
func statisticsTag(off uintptr) string              { return statisticsString(memmap.PtrOff(0x587000, off)) }
func statisticsCopyString(p unsafe.Pointer, s string) {
	copy(unsafe.Slice((*byte)(p), len(s)+1), s)
	*controlByte(p, len(s)) = 0
}
func statisticsAllocate(n int) unsafe.Pointer {
	if n == 0 {
		n = 1
	}
	p, _ := alloc.Malloc(uintptr(n))
	return p
}
func statisticsFree(p unsafe.Pointer) {
	if p != nil {
		alloc.FreePtr(p)
	}
}

func statisticsRow(pl unsafe.Pointer) uint32 {
	index := statisticsI32(pl, 4648)
	if *index != -1 {
		return uint32(*index)
	}
	row := statisticsPlayers
	statisticsPlayers++
	statisticsCopyString(memmap.PtrOff(0x5D4594, 600124+32*uintptr(row)), statisticsString(unsafe.Add(pl, 2096)))
	ind := int(*controlByte(pl, 2064)) + 1
	if ind == 32 {
		ind = 0
	}
	*memmap.PtrUint32(0x5D4594, 600136+32*uintptr(row)) = bits.ReverseBytes32(nox_xxx_net_getIP_554200(ind))
	*memmap.PtrUint32(0x5D4594, 600140+32*uintptr(row)) = *statisticsU32(pl, 2068)
	*memmap.PtrUint8(0x5D4594, 600144+32*uintptr(row)) = *controlByte(pl, 2251)
	*index = int32(row)
	return row
}
func statisticsEvent(actor, target unsafe.Pointer) uint32 {
	if !noxflags.HasGame(0x2000) {
		return 0
	}
	if noxflags.HasGame(4096) {
		return 1
	}
	if actor == nil {
		return 0
	}
	a := statisticsRow(actor)
	b := uint32(255)
	if target != nil {
		b = statisticsRow(target)
	}
	*memmap.PtrUint8(0x5D4594, 608320+2*uintptr(statisticsEvents)) = byte(a)
	*memmap.PtrUint8(0x5D4594, 608321+2*uintptr(statisticsEvents)) = byte(b)
	statisticsEvents++
	if statisticsPlayers >= 255 {
		return statisticsFlush()
	}
	return statisticsPlayers
}
func statisticsParticipation(pl unsafe.Pointer, value byte) uint32 {
	if !noxflags.HasGame(0x2000) {
		return 0
	}
	if noxflags.HasGame(4096) {
		return 1
	}
	index := *statisticsI32(pl, 4648)
	if index == -1 {
		return ^uint32(0)
	}
	off := uint32(index) * 32
	*memmap.PtrUint8(0x5D4594, 600152+uintptr(off)) = value
	return off
}
func statisticsRegister(pl unsafe.Pointer) {
	if !noxflags.HasGame(0x2000) || noxflags.HasGame(4096) || pl == nil || *statisticsI32(pl, 4648) != -1 {
		return
	}
	row := statisticsRow(pl)
	off := uintptr(row) * 32
	flags := *statisticsU32(pl, 3680)
	active := byte(1)
	if flags&1 != 0 && flags&0x20 == 0 || *controlByte(pl, 2064) == 31 && noxflags.HasEngine(noxflags.EngineNoRendering) {
		active = 0
	}
	*memmap.PtrUint8(0x5D4594, 600152+off) = active
	*memmap.PtrUint32(0x5D4594, 600148+off) = statisticsTime()
	*memmap.PtrUint8(0x5D4594, 600145+off) = 1
}
func statisticsMatchArrays(report, rows unsafe.Pointer, count int) {
	if names := *controlPtr(report, 608); names != nil {
		for i := uint32(0); i < statisticsArrayCount; i++ {
			statisticsFree(*controlPtr(names, int(i)*4))
		}
	}
	for off := 608; off <= 632; off += 4 {
		statisticsFree(*controlPtr(report, off))
		width := 4
		if off == 620 || off == 624 || off == 632 {
			width = 1
		}
		*controlPtr(report, off) = statisticsAllocate(width * count)
	}
	for i := 0; i < count; i++ {
		row := unsafe.Add(rows, 32*i)
		name := statisticsString(row)
		// C reserved ten bytes for a name; allocate enough for the actual string.
		p := statisticsAllocate(max(10, len(name)+1))
		statisticsCopyString(p, name)
		*controlPtr(*controlPtr(report, 608), 4*i) = p
		for _, pair := range [][2]int{{612, 12}, {616, 16}} {
			*statisticsU32(*controlPtr(report, pair[0]), 4*i) = *statisticsU32(row, pair[1])
		}
		for _, pair := range [][2]int{{620, 20}, {624, 21}, {632, 28}} {
			*controlByte(*controlPtr(report, pair[0]), i) = *controlByte(row, pair[1])
		}
		elapsed := statisticsTime() - *statisticsU32(row, 24)
		*statisticsU32(row, 24) = elapsed
		*statisticsU32(*controlPtr(report, 628), 4*i) = elapsed
	}
	*controlHalf(report, 6) = uint16(count)
	statisticsArrayCount = uint32(count)
}
func statisticsEventArray(report, events unsafe.Pointer, count int) uint32 {
	statisticsFree(*controlPtr(report, 636))
	p := statisticsAllocate(2 * count)
	*controlPtr(report, 636) = p
	copy(unsafe.Slice((*byte)(p), 2*count), unsafe.Slice((*byte)(events), 2*count))
	*memmap.PtrUint32(0x5D4594, 741308) = uint32(count)
	return uint32(2 * count)
}
func statisticsSend(report unsafe.Pointer, mode int) uint32 {
	data := statisticsReport(report, mode, int16(memmap.Uint32(0x5D4594, 741308)))
	_, n := statisticsEncode(data)
	*memmap.PtrInt32(0x5D4594, 741312) = n
	// The legacy service directory has no population path.
	return 1
}
func statisticsFlush() uint32 {
	if !noxflags.HasGame(0x2000) {
		return 0
	}
	if noxflags.HasGame(4096) {
		return 1
	}
	report := memmap.PtrOff(0x5D4594, 599476)
	statisticsMatchArrays(report, memmap.PtrOff(0x5D4594, 600124), int(statisticsPlayers))
	statisticsEventArray(report, memmap.PtrOff(0x5D4594, 608320), int(statisticsEvents))
	*memmap.PtrUint32(0x5D4594, 599504) = statisticsTime() - statisticsStart
	statisticsSend(report, 1)
	clear(unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 600124)), 0x2000))
	clear(unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 608320)), 0x20000))
	statisticsPlayers = 0
	statisticsEvents = 0
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		pl.Field4648 = -1
	}
	if pl := s.Players.ByInd(31); pl != nil {
		statisticsRegister(pl.C())
	}
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		if pl.PlayerInd != 31 {
			statisticsRegister(pl.C())
		}
	}
	return 0
}
