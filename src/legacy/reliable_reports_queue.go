package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

// C-owned pool storage remains shared with the engine lifecycle. Keep this
// layout until all consumers of the legacy allocation state have moved.
type reliableMessage struct {
	Frame                                    uint32
	LastSent                                 [32]uint32
	Countdown                                [32]byte
	Retries                                  byte
	_                                        [3]byte
	Acknowledged, Sent, Recipients, Priority uint32
	Ordered                                  byte
	_                                        byte
	Sequence                                 [32]uint16
	To                                       byte
	Data                                     [150]byte
	Size                                     byte
	_                                        [2]byte
	Related                                  *server.Object
	Next, Previous                           *reliableMessage
}

type reliableRateState struct {
	Budget, Delay, Rate, Reserved byte
	High, Low                     uint32
}

func reliableHead() **reliableMessage {
	return &dword_5d4594_1565512
}
func reliableTail() **reliableMessage {
	return &dword_5d4594_1565516
}
func reliableCapacity() *uint32      { return (*uint32)(unsafe.Pointer(&dword_5d4594_1565520)) }
func reliablePool() *unsafe.Pointer  { return &reliablePoolPtr }
func reliableSequences() *[32]uint16 { return (*[32]uint16)(memmap.PtrOff(0x5D4594, 1565524)) }
func reliableRates() *[32]reliableRateState {
	return (*[32]reliableRateState)(memmap.PtrOff(0x5D4594, 1565124))
}
func reliableMask() uint32 { return uint32(dword_5d4594_2649712) }
func reliableRate() uint32 { return *memmap.PtrUint32(0x587000, 4728) }

func reliableRecalculate(to int) int {
	p := &reliableRates()[to]
	divisor := uint32(1)
	if Get_dword_5d4594_2650652() == 1 {
		divisor = reliableRate()
	}
	frames := GetServer().S().TickRate() / divisor
	p.Low = 0
	if p.Delay > 2 {
		p.Low = uint32(p.Budget) * uint32(p.Delay-1) * frames
	}
	p.High = uint32(p.Delay) * uint32(p.Budget) * frames
	return int(p.High)
}
func reliableResetRate(to int) int {
	p := &reliableRates()[to]
	p.Budget = 1
	p.Delay = 2
	p.Rate = byte(reliableRate())
	return reliableRecalculate(to)
}
func reliableResetRates() int {
	clear(reliableRates()[:])
	var rv int
	for i := 0; i < 32; i++ {
		rv = reliableResetRate(i)
	}
	return rv
}
func reliableResetSequences() int                  { clear(reliableSequences()[:]); return 0 }
func ReliableResetSequence(to ntype.PlayerInd) int { reliableSequences()[to] = 0; return int(to) }
func reliableInit() int {
	if p := *reliablePool(); p != nil {
		alloc.AsClass(p).Free()
	}
	*reliablePool() = nil
	reliableResetSequences()
	*reliableHead() = nil
	*reliableTail() = nil
	return reliableResetRates()
}
func reliableUnlink(n *reliableMessage) {
	if n.Previous != nil {
		n.Previous.Next = n.Next
	} else {
		*reliableHead() = n.Next
	}
	if n.Next != nil {
		n.Next.Previous = n.Previous
	} else {
		*reliableTail() = n.Previous
	}
	alloc.AsClass(*reliablePool()).FreeObjectFirst(unsafe.Pointer(n))
}
func reliableTrim() int {
	for n := *reliableHead(); n != nil; {
		next := n.Next
		if n.Data[0] >= 49 && n.Data[0] <= 51 {
			reliableUnlink(n)
		}
		n = next
	}
	return 0
}
func reliableEnqueue(to int, data []byte, related *server.Object, priority int, ordered byte) int {
	mask := reliableMask()
	if to != 255 && to&128 != 0 && mask&^(uint32(1)<<uint(to&127)) == 0 {
		return 1
	}
	if len(data) > 150 {
		return 0
	}
	pool := alloc.AsClass(*reliablePool())
	if pool == nil {
		capacity := 256
		dynamic := noxflags.HasGame(noxflags.GameFlag(2048))
		if dynamic {
			capacity = 512
		} else if noxflags.HasGame(noxflags.GameFlag(1)) {
			capacity = 3072
		}
		*reliableCapacity() = uint32(capacity)
		if dynamic {
			pool = alloc.NewDynamicClass("importantClass", 416, capacity)
		} else {
			pool = alloc.NewClass("importantClass", 416, capacity)
		}
		*reliablePool() = pool.UPtr()
	}
	n := (*reliableMessage)(pool.NewObject())
	if n == nil {
		if reliablePressure() != 1 {
			return 0
		}
		n = (*reliableMessage)(pool.NewObject())
		if n == nil {
			return 0
		}
	}
	copy(n.Data[:], data)
	n.Size = byte(len(data))
	n.Related = related
	n.To = byte(to)
	n.Priority = uint32(priority)
	n.Ordered = ordered
	n.Recipients = mask
	n.Frame = GetServer().S().Frame()
	if ordered != 0 {
		seq := reliableSequences()
		for i := 0; i < 32; i++ {
			include := to == 255 && mask&(uint32(1)<<uint(i)) != 0 || to&128 == 0 && i == to || to != 255 && to&128 != 0 && i != to&127 && mask&(uint32(1)<<uint(i)) != 0
			if include {
				n.Sequence[i] = seq[i]
				seq[i]++
			}
		}
	}
	n.Next = *reliableHead()
	if n.Next != nil {
		n.Next.Previous = n
	} else {
		*reliableTail() = n
	}
	*reliableHead() = n
	return 1
}
func reliableAcknowledge(mask uint32, n *reliableMessage, to int) {
	if n.Related != nil && n.Data[0] >= 49 && n.Data[0] <= 51 {
		n.Related.Field37 &^= mask
	}
	audience := reliableMask() & n.Recipients
	if n.To < 128 {
		if int(n.To) == to {
			reliableUnlink(n)
		}
		return
	}
	n.Acknowledged |= mask
	if n.To != 255 {
		audience &^= uint32(1) << uint(n.To&127)
	}
	if audience&n.Acknowledged == audience {
		reliableUnlink(n)
	}
}
func reliableRemoveRecipient(to byte) int {
	for n := *reliableHead(); n != nil; {
		next := n.Next
		reliableAcknowledge(uint32(1)<<uint(to), n, int(to))
		n = next
	}
	return 0
}
func reliableACK(to int, frame uint32) int {
	for n := *reliableHead(); n != nil; {
		next := n.Next
		if n.LastSent[to] == frame {
			reliableAcknowledge(uint32(1)<<uint(to), n, to)
		}
		n = next
	}
	return 0
}
func reliableRemoveSlowPlayer(to int) int {
	pl := GetServer().S().Players.ByInd(ntype.PlayerInd(to))
	if pl == nil {
		return 0
	}
	reliableRemoveRecipient(byte(to))
	Nox_xxx_netNeedTimestampStatus_4174F0(pl, 128)
	return bool2int(noxflags.HasGame(noxflags.GameFlag(1)))
}
func reliablePressure() int {
	var oldest *reliableMessage
	threshold := uint32(999999999)
	counts := [32]uint16{}
	slow := -1
	var maximum uint16
	for n := *reliableHead(); n != nil; n = n.Next {
		if n.To < 128 && n.To != 31 {
			counts[n.To]++
			if counts[n.To] > maximum {
				slow = int(n.To)
				maximum = counts[n.To]
			}
		}
		if n.Frame < threshold {
			threshold = n.Frame
			oldest = n
		}
	}
	if slow != -1 {
		reliableRemoveSlowPlayer(slow)
	}
	if oldest == nil {
		return 0
	}
	// Removing the slow player may already have released the selected node.
	for n := *reliableHead(); n != nil; n = n.Next {
		if n == oldest {
			reliableUnlink(n)
			break
		}
	}
	return 1
}
func reliableCoalesce(to int, data []byte, related *server.Object, priority int) int {
	opcode := data[0]
	for n := *reliableHead(); n != nil; {
		next := n.Next
		if n.Data[0] == opcode {
			if to == 255 || to&128 != 0 {
				reliableUnlink(n)
			} else {
				reliableAcknowledge(uint32(1)<<uint(to), n, to)
			}
		}
		n = next
	}
	return reliableEnqueue(to, data, related, priority, 0)
}
func reliableAdapt(to byte) int {
	GetServer().S().Players.ByInd(ntype.PlayerInd(to))
	count := uint32(0)
	mask := uint32(1) << uint(to)
	for n := *reliableTail(); n != nil; n = n.Previous {
		if n.Acknowledged&mask == 0 && (n.To == 255 || n.To < 128 && n.To == to || n.To >= 128 && n.To != 255 && to != n.To&127) {
			count++
		}
	}
	p := &reliableRates()[to]
	p.Rate = byte(reliableRate())
	if count > p.High {
		p.Delay++
		if p.Delay > 5 {
			if p.Budget == 2 {
				reliableRemoveSlowPlayer(int(to))
			}
			p.Delay = 5
			p.Budget = 2
		}
		p.Rate = byte(reliableRate())
	} else if int32(p.Low) > 0 && int32(count) < int32(p.Low) {
		if p.Budget == 2 {
			p.Budget = 1
			return reliableRecalculate(int(to))
		}
		p.Delay--
		if p.Delay < 2 {
			p.Delay = 2
		}
	}
	return reliableRecalculate(int(to))
}
