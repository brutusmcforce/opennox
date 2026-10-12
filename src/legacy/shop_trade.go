package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

func shopSend(u *server.Object, b []byte, ordered int) uint32 {
	return uint32(Nox_xxx_netSendPacket_4E5030(int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), b, 0, 1, ordered))
}
func shopSendShort(u *server.Object, code uint16, ordered int) uint32 {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], code)
	return shopSend(u, b[:], ordered)
}
func shopSendCode(u, item *server.Object, code uint16) uint32 {
	var b [4]byte
	binary.LittleEndian.PutUint16(b[:], code)
	binary.LittleEndian.PutUint16(b[2:], uint16(item.NetCode))
	return shopSend(u, b[:], 1)
}
func shopSendItem(u *server.Object, n *shopItem) uint32 {
	it := n.Object
	var b [18]byte
	b[0], b[1] = 0xc9, 8
	binary.LittleEndian.PutUint16(b[2:], it.TypeInd)
	binary.LittleEndian.PutUint16(b[4:], uint16(it.NetCode))
	binary.LittleEndian.PutUint32(b[6:], n.Value)
	if hp := it.HealthData; hp != nil {
		binary.LittleEndian.PutUint32(b[10:], uint32(hp.Cur))
	}
	for i := 14; i < 18; i++ {
		b[i] = 0xff
	}
	if uint32(it.ObjClass)&0x13001000 != 0 {
		for i, p := range unsafe.Slice((*unsafe.Pointer)(it.InitData), 4) {
			if p != nil {
				b[14+i] = *(*byte)(unsafe.Add(p, 4))
			}
		}
	}
	return shopSend(u, b[:], 1)
}
func shopSendAcceptance(u *server.Object, s *shopSession) uint32 {
	b := [3]byte{0xc9, 3, 0}
	for i := 0; i < 2; i++ {
		if s.Accepted[i] != 0 {
			if s.Units[i] == u {
				b[2] |= 1
			} else {
				b[2] |= 2
			}
		}
	}
	return shopSend(u, b[:], 1)
}
func shopSendGold(u *server.Object, s *shopSession) uint32 {
	var b [14]byte
	b[0], b[1] = 0xc9, 6
	side := -1
	if s.Units[0] == u {
		side = 0
	} else if s.Units[1] == u {
		side = 1
	}
	if side >= 0 {
		other := 1 - side
		binary.LittleEndian.PutUint32(b[2:], *shopGold(s.Gold[side]))
		if s.Units[other].ObjClass&4 == 0 {
			binary.LittleEndian.PutUint32(b[6:], s.Totals[other]-s.Totals[side])
		}
		binary.LittleEndian.PutUint32(b[10:], *shopGold(s.Gold[other]))
	}
	return shopSend(u, b[:], 1)
}
func shopTotal(s *shopSession, u *server.Object) uint32 {
	side := 1
	if s.Units[0] == u {
		side = 0
	}
	value := *shopGold(s.Gold[side])
	for n := s.Offers[side]; n != nil; n = n.Next {
		value += n.Value
	}
	return value
}
func shopBalance(s *shopSession) uint32 {
	result := uint32(s.Units[0].ObjClass & 4)
	if result == 4 && s.Units[1].ObjClass&4 == 4 {
		return result
	}
	for i := 0; i < 2; i++ {
		gold := shopGold(s.Gold[i])
		if s.Units[i].ObjClass&4 != 0 {
			shopAddGold(s.Units[i], *gold)
		}
		*gold = 0
	}
	s.Totals[0] = shopTotal(s, s.Units[0])
	s.Totals[1] = shopTotal(s, s.Units[1])
	result = s.Totals[1]
	if s.Totals[0] == s.Totals[1] {
		return result
	}
	side := 0
	if s.Totals[1] < s.Totals[0] {
		side = 1
	}
	amount := s.Totals[1-side] - s.Totals[side]
	gold := shopGold(s.Gold[side])
	if s.Units[side].ObjClass&4 != 0 {
		paid := min(amount, shopGetGold(s.Units[side]))
		s.Totals[side] += paid
		*gold += paid
		return shopSubGold(s.Units[side], paid)
	}
	s.Totals[side] = s.Totals[1-side]
	*gold += amount
	return memmap.PtrWord(unsafe.Pointer(gold))
}
func shopTransfer(s *shopSession, u *server.Object, n *shopItem) {
	for ; n != nil; n = n.Next {
		if u.ObjClass&4 != 0 {
			shopPut(u, n.Object)
		} else if added := shopAdd(s, n.Object); added != nil {
			player := s.Units[0]
			if player.ObjClass&4 == 0 {
				player = s.Units[1]
			}
			shopSendItem(player, added)
		}
	}
}
func shopTransferGold(u, gold *server.Object) uint32 {
	if u.ObjClass&4 != 0 {
		return shopAddGold(u, *shopGold(gold))
	}
	return memmap.PtrWord(u.CObj())
}
func shopCancelTrade(s *shopSession) {
	for i := 0; i < 2; i++ {
		for n := s.Offers[i]; n != nil; n = n.Next {
			shopPut(s.Units[i], n.Object)
		}
		shopAddGold(s.Units[i], *shopGold(s.Gold[i]))
		shopSendShort(s.Units[i], 457, 0)
	}
	shopDetach(s, s.Units[0])
	shopDetach(s, s.Units[1])
	shopDestroy(s)
}
func shopAccept(s *shopSession, u *server.Object) {
	if s.Units[0] == u {
		s.Accepted[0] = 1
	} else if s.Units[1] == u {
		s.Accepted[1] = 1
	}
	for _, u := range s.Units {
		if u.ObjClass&4 != 0 {
			shopSendAcceptance(u, s)
		}
	}
	if s.Accepted[0] != 1 || s.Accepted[1] != 1 {
		return
	}
	shopTransfer(s, s.Units[1], s.Offers[0])
	shopTransfer(s, s.Units[0], s.Offers[1])
	shopTransferGold(s.Units[1], s.Gold[0])
	shopTransferGold(s.Units[0], s.Gold[1])
	shopFreeList(s.Offers[0])
	next := s.Offers[1]
	s.Offers[0] = nil
	shopFreeList(next)
	s.Offers[1] = nil
	*shopGold(s.Gold[0]), *shopGold(s.Gold[1]) = 0, 0
	s.Totals = [2]uint32{}
	if s.Kind == 1 {
		player := s.Units[0]
		if player.ObjClass&4 == 0 {
			player = s.Units[1]
		}
		shopSendShort(player, 1993, 1)
	} else {
		shopCancelTrade(s)
	}
}
func shopWithdraw(s *shopSession, code uint32) uint32 {
	side := 0
	n := shopFind(s.Offers[0], code)
	if n == nil {
		side = 1
		n = shopFind(s.Offers[1], code)
	}
	if n == nil {
		return 0
	}
	u := s.Units[side]
	if u.ObjClass&4 != 0 {
		shopPut(u, n.Object)
	} else if added := shopAdd(s, n.Object); added != nil {
		player := s.Units[0]
		if player.ObjClass&4 == 0 {
			player = s.Units[1]
		}
		shopSendItem(player, added)
	}
	if n.Next != nil {
		n.Next.Prev = n.Prev
	}
	if n.Prev != nil {
		n.Prev.Next = n.Next
	}
	index := 1
	if s.Units[0] == u {
		index = 0
	}
	if s.Offers[index] == n {
		s.Offers[index] = n.Next
	}
	shopBalance(s)
	s.Accepted = [2]uint32{}
	if s.Units[0].ObjClass&4 != 0 {
		if s.Units[1].ObjClass&4 == 0 && s.Totals[1] <= s.Totals[0] {
			s.Accepted[1] = 1
		}
	} else if s.Totals[0] <= s.Totals[1] {
		s.Accepted[0] = 1
	}
	for _, u := range s.Units {
		if u.ObjClass&4 != 0 {
			shopSendGold(u, s)
			shopSendCode(u, n.Object, 1481)
			shopSendAcceptance(u, s)
		}
	}
	shopFreeItem(n)
	return 1
}
