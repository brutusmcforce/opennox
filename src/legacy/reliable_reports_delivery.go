package legacy

import (
	"encoding/binary"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func reliableClientSend(to int, data []byte, related *server.Object, priority int) int {
	if noxflags.HasGame(noxflags.GameFlag(1)) {
		GetServer().S().NetList.AddToMsgListCli(ntype.PlayerInd(to), netlist.Kind0, data)
		return 1
	}
	if to == 255 || to&128 != 0 {
		return 0
	}
	return reliableEnqueue(to, data, related, priority, 0)
}
func reliableDeliver(to byte, kind int) {
	s := GetServer().S()
	pl := s.Players.ByInd(ntype.PlayerInd(to))
	hosting := noxflags.HasGame(noxflags.GameFlag(1))
	if pl != nil && hosting && *(*byte)(unsafe.Add(pl.C(), int(layoutServerPlayer.Off(3680)) /* Field3680 */))&16 == 0 {
		return
	}
	send := func(data []byte) bool {
		if to == 31 {
			return s.NetList.AddToMsgListCli(ntype.PlayerInd(to), netlist.Kind(kind), data)
		}
		return s.NetList.ClientSend0(ntype.PlayerInd(to), netlist.Kind(kind), data, GetNetPlayerBufSize())
	}
	bit := uint32(1) << uint(to)
	timestamp := true
	retried := 0
	p := &reliableRates()[to]
	for n := *reliableTail(); n != nil; {
		next := n.Previous
		if n.Related != nil && uint32(n.Related.ObjFlags)&32 != 0 {
			n.Related = nil
		}
		if n.Acknowledged&bit != 0 || n.To != 255 && (n.To < 128 && n.To != to || n.To >= 128 && to == n.To&127) {
			n = next
			continue
		}
		if n.Related != nil && n.Related.Field37&bit == 0 || n.Priority != 0 && reliableMask()&bit == 0 {
			reliableAcknowledge(bit, n, int(to))
			return
		}
		if n.Sent&bit != 0 {
			if n.Countdown[to] != 0 {
				n.Countdown[to]--
				n = next
				continue
			}
			if retried >= int(p.Budget) {
				n = next
				continue
			}
			n.Sent &^= bit
			n.Retries++
			retried++
		}
		if timestamp {
			var stamp [5]byte
			stamp[0] = 170
			binary.LittleEndian.PutUint32(stamp[1:], s.Frame())
			data := stamp[:]
			if kind != 0 && to != 31 {
				data = data[:1]
			}
			if !send(data) {
				return
			}
			timestamp = false
		}
		data := n.Data[:int(n.Size)]
		if n.Ordered != 0 {
			scratch := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1564964)), 160)
			scratch[0] = 204
			binary.LittleEndian.PutUint16(scratch[1:], n.Sequence[to])
			scratch[3] = n.Size
			copy(scratch[4:], data)
			data = scratch[:int(n.Size)+4]
		}
		if send(data) {
			n.Sent |= bit
			n.Countdown[to] = byte(s.TickRate() * uint32(p.Delay) / reliableRate())
			n.LastSent[to] = s.Frame()
			if noxflags.HasEngine(noxflags.EngineReplayRead) {
				reliableAcknowledge(bit, n, int(to))
			}
		}
		n = next
	}
	if hosting && s.Frame()%(s.TickRate()*uint32(p.Delay)) == 0 {
		reliableAdapt(to)
	}
}
