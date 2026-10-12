package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// Preserve the qualified fixed-pool record layout and borrowed object references.
type voteRecord struct {
	kind       uint32
	count      byte
	_          [3]byte
	voters     uint32
	minimum    byte
	_          [3]byte
	team       *server.ObjectTeam
	teamOnly   uint32
	frame      uint32
	target     *server.Object
	_          [3]uint32
	next, prev *voteRecord
}

var votePool alloc.ClassT[voteRecord]
var voteHead *voteRecord
var voteQuestResetSetting uint32 = 6
var voteQuestKickSetting uint32 = 5

func voteInit() int {
	votePool = alloc.NewClassT("VoteClass", voteRecord{}, 64)
	voteHead = nil
	return 1
}
func voteClose()                     { votePool.Free(); votePool.Class = nil; voteHead = nil }
func voteStatus(to int, on byte) int { return int(reliableEnqueue(to, []byte{238, 6, on}, nil, 1, 1)) }
func voteResetNotify(to int) int     { return int(reliableEnqueue(to, []byte{238, 7}, nil, 1, 1)) }
func voteCreate(kind uint32, u *server.Object) *voteRecord {
	if u == nil || u.ObjClass&4 == 0 {
		return nil
	}
	first := voteHead == nil
	p := votePool.NewObject()
	if p == nil {
		return nil
	}
	*p = voteRecord{}
	p.kind = kind
	p.frame = GetServer().S().Frame()
	p.team = u.TeamPtr()
	switch kind {
	case 0, 1:
		p.minimum = memmap.Uint8(0x587000, 229980)
	case 2, 3:
		p.minimum = 6
	default:
		p.minimum = memmap.Uint8(0x587000, 229984)
	}
	p.next = voteHead
	if voteHead != nil {
		voteHead.prev = p
	}
	voteHead = p
	if first {
		voteStatus(255, 1)
	}
	return p
}
func voteDelete(p *voteRecord) {
	if p == nil {
		return
	}
	if p.kind == 2 {
		for i := 0; i < 32; i++ {
			if p.voters&(uint32(1)<<i) != 0 {
				voteResetNotify(i)
			}
		}
	}
	if p.next != nil {
		p.next.prev = p.prev
	}
	if p.prev != nil {
		p.prev.next = p.next
	} else {
		voteHead = p.next
	}
	votePool.FreeObjectFirst(p)
	if voteHead == nil {
		voteStatus(255, 0)
	}
}
func votePlayerBit(u *server.Object) uint32 {
	return uint32(1) << u.UpdateDataPlayer().Player.PlayerInd
}
func voteRemovePlayer(u *server.Object) {
	if u == nil || u.ObjClass&4 == 0 {
		return
	}
	bit := votePlayerBit(u)
	for p := voteHead; p != nil; {
		next := p.next
		if p.voters&bit != 0 {
			p.voters &^= bit
			p.count--
		}
		if p.count == 0 {
			voteDelete(p)
		}
		p = next
	}
}
func voteNamesEqual(a, b *uint16) bool {
	for a != nil && b != nil {
		if *a != *b {
			return false
		}
		if *a == 0 {
			return true
		}
		a = (*uint16)(unsafe.Add(unsafe.Pointer(a), 2))
		b = (*uint16)(unsafe.Add(unsafe.Pointer(b), 2))
	}
	return a == b
}
func voteTarget(name *uint16) *server.Player {
	if name == nil {
		return nil
	}
	players := &GetServer().S().Players
	for p := players.First(); p != nil; p = players.Next(p) {
		if p.Active == 1 && voteNamesEqual(&p.NameFinal[0], name) {
			return p
		}
	}
	return nil
}
func voteFind(kind uint32, target *server.Object) *voteRecord {
	for p := voteHead; p != nil; p = p.next {
		if p.kind == kind && (kind == 2 || p.target == target) {
			return p
		}
	}
	return nil
}
func voteCast(kind uint32, u *server.Object, name *uint16) {
	if u == nil || u.ObjClass&4 == 0 {
		return
	}
	var target *server.Object
	switch kind {
	case 0, 1:
		threshold := memmap.Uint32(0x587000, 229980)
		if threshold == 0 || threshold > 32 {
			return
		}
	case 2:
		if voteQuestResetSetting == 0 || u.UpdateDataPlayer().Player.Field4792 == 0 {
			return
		}
	case 3:
		if voteQuestKickSetting == 0 || u.UpdateDataPlayer().Player.Field4792 == 0 {
			return
		}
	default:
		return
	}
	if kind != 2 {
		player := voteTarget(name)
		if player == nil || player.PlayerInd == 31 {
			return
		}
		target = player.PlayerUnit
		if target == nil || target == u {
			return
		}
		if kind == 3 {
			if player.Field4792 == 0 {
				return
			}
		} else if noxflags.HasGamePlay(4) && !u.TeamPtr().SameAs(target.TeamPtr()) {
			return
		}
	}
	p := voteFind(kind, target)
	if p == nil {
		p = voteCreate(kind, u)
		if p == nil {
			return
		}
		p.target = target
		if kind < 2 && noxflags.HasGamePlay(4) {
			p.teamOnly = 1
		}
	}
	bit := votePlayerBit(u)
	if p.voters&bit == 0 {
		p.voters |= bit
		p.count++
	}
}
func voteWithdraw(kind uint32, u *server.Object, name *uint16) {
	if u == nil || u.ObjClass&4 == 0 || kind > 3 {
		return
	}
	var target *server.Object
	if kind != 2 {
		player := voteTarget(name)
		if player == nil || player.PlayerInd == 31 || player.PlayerUnit == nil {
			return
		}
		target = player.PlayerUnit
	}
	bit := votePlayerBit(u)
	for p := voteHead; p != nil; p = p.next {
		if p.kind != kind || kind != 2 && p.target != target {
			continue
		}
		if p.voters&bit != 0 {
			p.voters &^= bit
			p.count--
			if p.count == 0 {
				voteDelete(p)
			}
			return
		}
		// Reset votes are unique by kind; other kinds search for the voter's record.
		if kind == 2 {
			return
		}
	}
}
func voteThreshold(p *voteRecord) int {
	if p.count >= p.minimum {
		return 1
	}
	count := uint32(0)
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		if p.teamOnly != 1 || p.team.SameAs(u.TeamPtr()) {
			count++
		}
	}
	if uint32(p.count) >= count-1 && p.count >= 2 {
		return 1
	}
	return 0
}
func voteBlockPlayer(pl *server.Player) {
	serverConfigBlockedAdd(15, &pl.NameFinal[0], (*byte)(unsafe.Add(pl.C(), int(layoutServerPlayer.Off(2112)) /* SerialBuf */)))
}
func voteUpdateNormal(p *voteRecord) {
	u := p.target
	if u.ObjFlags&0x20 != 0 {
		voteDelete(p)
		return
	}
	p.team = u.TeamPtr()
	if voteThreshold(p) == 1 {
		data := u.UpdateDataPlayer()
		Nox_xxx_playerCallDisconnect_4DEAB0(ntype.PlayerInd(data.Player.PlayerInd), 4)
		voteBlockPlayer(data.Player)
		voteDelete(p)
	}
}
func voteUpdateReset(p *voteRecord) {
	if int(p.count) < questRuntimeCount() {
		return
	}
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		if u.UpdateDataPlayer().Player.Field4792 == 1 {
			controlRespawn(u)
		}
	}
	questRuntimeSetStage(0)
	GetServer().SwitchMap(alloc.GoString(mapQuestChoose()))
	voteDelete(p)
}
func voteUpdateQuest(p *voteRecord) {
	u := p.target
	if u == nil || u.ObjFlags&0x20 != 0 || u.UpdateDataPlayer().Player.Field4792 == 0 {
		voteDelete(p)
		return
	}
	if p.count < p.minimum {
		count := questRuntimeCount()
		if count <= 1 {
			voteDelete(p)
			return
		}
		if int(p.count) < count-1 || p.count < 2 {
			return
		}
	}
	data := u.UpdateDataPlayer()
	Sub_4DCFB0(u)
	voteBlockPlayer(data.Player)
	voteDelete(p)
}
func voteTick() {
	for p := voteHead; p != nil; {
		next := p.next
		switch p.kind {
		case 0, 1:
			voteUpdateNormal(p)
		case 2:
			voteUpdateReset(p)
		case 3:
			voteUpdateQuest(p)
		}
		p = next
	}
}
