package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/server"
)

type spellLifeBook struct {
	Unused        uint32
	Owner         *server.Object
	Spells        [5]int32
	Index, Queued byte
	Padding30     [2]byte
	Tree          *server.PhonemeLeaf
	Phoneme       byte
	Padding37     [3]byte
	Frame, Delay  uint32
	Self          int32
	Next, Prev    *spellLifeBook
}

// NewSpellLifeBookPool allocates the pool behind spellLifeBookPool.
func NewSpellLifeBookPool() unsafe.Pointer {
	return alloc.NewClass("magicEntityClass", unsafe.Sizeof(spellLifeBook{}), 64).UPtr()
}
func spellLifeBookPool() alloc.ClassT[spellLifeBook] {
	return alloc.AsClassT[spellLifeBook](legacyGlobals.nox_alloc_magicEnt_1569668)
}
func spellLifeBookHead() *spellLifeBook {
	return (*spellLifeBook)(memmap.WordPtr(dword_5d4594_1569672))
}
func spellLifeSetBookHead(p *spellLifeBook) {
	dword_5d4594_1569672 = memmap.PtrWord(unsafe.Pointer(p))
}
func spellLifeUnlinkBook(p *spellLifeBook) *spellLifeBook {
	next := p.Next
	if next != nil {
		next.Prev = p.Prev
	}
	if p.Prev != nil {
		p.Prev.Next = next
	} else {
		spellLifeSetBookHead(next)
	}
	spellLifeBookPool().FreeObjectFirst(p)
	return next
}
func spellLifeReset(mode, imaginary int32) int32 {
	s := GetServer().S()
	s.Spells.Dur.Sub4FE8A0(int(mode))
	spellLifeBookPool().FreeAllObjects()
	spellLifeSetBookHead(nil)
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		p := u.UpdateData
		layout.As[server.PlayerUpdateData](p).Field47_0 = 0
		layout.As[server.PlayerUpdateData](p).SpellCastStart = 0
		layout.As[server.PlayerUpdateData](p).TrapSpells = [5]uint32{}
		*(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](p).TrapSpellsCnt)) = 0
	}
	if imaginary != 0 && Nox_setImaginaryCaster() == 0 {
		return 0
	}
	return 1
}
func spellLifeInform(u *server.Object, code int32) {
	gameplayTextInformation(int(int32(controlPlayer(u).PlayerInd)), 0, unsafe.Pointer(&code))
}
func spellLifeBookError(u *server.Object, code int32, audio int) int32 {
	spellLifeInform(u, code)
	nox_xxx_aud_501960(int32(int32(audio)), asObjectC(u), 0, 0)
	return 0
}
func spellLifeInsertBook(u *server.Object, list unsafe.Pointer, n, delay, self int32) int32 {
	if u.ObjFlags&0x8022 != 0 {
		return 0
	}
	words := unsafe.Slice((*int32)(list), 5)
	for _, v := range words {
		if v < 0 || v >= 137 {
			return 0
		}
	}
	if u.ObjClass&4 == 0 {
		return 0
	}
	ud := u.UpdateData
	pl := controlPlayer(u)
	if layout.As[server.PlayerUpdateData](ud).Trade70 != nil {
		return 0
	}
	for _, v := range words {
		if v != 0 && pl.SpellLvl[int(v)] == 0 {
			return 0
		}
	}
	if layout.As[server.PlayerUpdateData](ud).SpellCastStart != 0 {
		return 0
	}
	trap := false
	for i := int32(0); i < n; i++ {
		if *spellLifeWord(list, int(i)*4) == 34 {
			trap = true
		}
	}
	if trap {
		if spellLifeCheckMana(u, list, n) == 0 {
			return spellLifeBookError(u, 12, 232)
		}
		if pl.InfoBytes[66] == 2 {
			if !bool(nox_xxx_checkSummonedCreaturesLimit_500D70(asObjectC(u), 5)) {
				return spellLifeBookError(u, 4, 231)
			}
			count := stateCount(u, 2, 0x2000)
			if int32(count) >= int32(int64(nox_xxx_gamedataGetFloat_419D40(internCStr("MaxBomberCount")))) {
				return spellLifeBookError(u, 5, 231)
			}
		} else if int32(*(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).CurTraps))) >= int32(int64(nox_xxx_gamedataGetFloat_419D40(internCStr("MaxTrapCount")))) {
			return spellLifeBookError(u, 5, 231)
		}
		for i := int32(0); i < n; i++ {
			id := int32(*spellLifeWord(list, int(i)*4))
			if e := spellLifeCheckClass(u, id); e != 0 {
				return spellLifeBookError(u, e, 231)
			}
			if e := spellLifeCantCast(u, id, 1); e != 0 {
				return spellLifeBookError(u, e, 231)
			}
		}
	} else {
		if e := spellLifeCheckClass(u, words[0]); e != 0 {
			return spellLifeBookError(u, e, 231)
		}
		if e := spellLifeCantCast(u, words[0], 0); e != 0 {
			return spellLifeBookError(u, e, 231)
		}
	}
	Nox_xxx_playerSetState_4FA020(u, 2)
	layout.As[server.PlayerUpdateData](ud).Field47_0 = 1
	layout.As[server.PlayerUpdateData](ud).SpellCastStart = GetServer().S().Frame()
	p := spellLifeBookPool().NewObject()
	if p == nil {
		return 0
	}
	*p = spellLifeBook{}
	p.Owner = u
	p.Self = self
	p.Frame = GetServer().S().Frame()
	p.Delay = uint32(delay)
	p.Tree = GetServer().S().Spells.PhonemeTree()
	for i := int32(0); i < 5; i++ {
		if i < n {
			p.Spells[i] = words[i]
			if words[i] == 34 {
				p.Queued = 1
			}
		}
	}
	p.Next = spellLifeBookHead()
	if p.Next != nil {
		p.Next.Prev = p
	}
	spellLifeSetBookHead(p)
	return 1
}
func spellLifeAdvanceBook(p *spellLifeBook) {
	p.Phoneme = 0
	p.Tree = GetServer().S().Spells.PhonemeTree()
	p.Frame = GetServer().S().Frame() + p.Delay
	p.Index++
}
func spellLifeCastBooks() {
	s := GetServer().S()
	for p := spellLifeBookHead(); p != nil; {
		u := p.Owner
		if u.ObjFlags&0x8020 != 0 {
			p = spellLifeUnlinkBook(p)
			continue
		}
		if s.Frame() < p.Frame {
			p = p.Next
			continue
		}
		var ud unsafe.Pointer
		if u.ObjClass&4 != 0 {
			ud = u.UpdateData
		}
		id := int32(*(*uint32)(unsafe.Pointer(&p.Spells[int(p.Index)])))
		if p.Phoneme == 0 {
			msg := [2]byte{112, byte(id)}
			nox_netlist_addToMsgListCli_40EBC0(int(int32(layout.As[server.Player](unsafe.Pointer(layout.As[server.PlayerUpdateData](ud).Player)).PlayerInd)), 1, (*uint8)(unsafe.Pointer(&msg[0])), 2)
		}
		if p.Tree.Ind != id {
			settings := unsafe.Pointer((unsafe.Pointer)(unsafe.Pointer(serverConfigSettings())))
			ph := s.Spells.Phoneme(spell.ID(id), int(p.Phoneme))
			if dword_5d4594_2650652 == 0 || *spellLifeWord(settings, 62) != 0 {
				spellLifeBroadcastPhoneme(u, int8(ph))
			}
			p.Tree = p.Tree.Next(ph)
			if u.ObjClass&4 != 0 {
				*(*unsafe.Pointer)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).SpellPhonemeLeaf)) = p.Tree.C()
			}
			p.Phoneme++
			p.Frame = s.Frame() + p.Delay
			p = p.Next
			continue
		}
		// The source reads one word beyond Spells for index 4: the queue/index bytes.
		nextID := int32(*(*uint32)(unsafe.Pointer(&p.Spells[1+int(p.Index)])))
		if p.Queued != 0 && id != 34 && u.ObjClass&4 != 0 {
			allowed := true
			for i := 0; i < int(*(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).TrapSpellsCnt))); i++ {
				if int32(layout.As[server.PlayerUpdateData](ud).TrapSpells[i]) == id {
					spellLifeInform(u, 6)
					allowed = false
				}
			}
			if allowed {
				if spellLifeSpendMana(u, id, 2) < 0 {
					spellLifeBookError(u, 11, 232)
				} else {
					i := *(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).TrapSpellsCnt))
					*(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).TrapSpellsCnt)) = i + 1
					layout.As[server.PlayerUpdateData](ud).TrapSpells[int(i)] = uint32(id)
				}
			}
		}
		if id != 34 && nextID != 0 {
			spellLifeAdvanceBook(p)
			p = p.Next
			continue
		}
		if u.ObjClass&4 != 0 {
			pl := controlPlayer(u)
			layout.As[server.PlayerUpdateData](ud).Field55 = pl.CursorVec.X
			layout.As[server.PlayerUpdateData](ud).Field56 = pl.CursorVec.Y
			if p.Self != 0 {
				*(*unsafe.Pointer)(unsafe.Pointer(&pl.Obj3640)) = u.CObj()
			} else {
				*(*unsafe.Pointer)(unsafe.Pointer(&pl.Obj3640)) = *(*unsafe.Pointer)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).CursorObj))
			}
			GetServer().PlayerSpell(u)
			layout.As[server.PlayerUpdateData](ud).SpellCastStart = 0
			layout.As[server.PlayerUpdateData](ud).Field47_0 = 0
			*(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](ud).TrapSpellsCnt)) = 0
		} else {
			Nox_xxx_castSpellByUser_4FDD20(int(id), u, nil)
		}
		p = spellLifeUnlinkBook(p)
	}
}
func spellLifeCounterBooks(u *server.Object, radius float32) {
	for p := spellLifeBookHead(); p != nil; {
		owner := p.Owner
		same := owner.ObjClass&4 != 0 && u.TeamVal.SameAs(&owner.TeamVal)
		dx := float64(owner.PosVec.X) - float64(u.PosVec.X)
		dy := float64(owner.PosVec.Y) - float64(u.PosVec.Y)
		if !same && math.Sqrt(dx*dx+dy*dy)+0.1 < float64(radius) && nox_xxx_mapCheck_537110(asObjectC(u), asObjectC(owner)) != 0 {
			if owner.ObjClass&4 != 0 {
				layout.As[server.PlayerUpdateData](owner.UpdateData).SpellCastStart = 0
				layout.As[server.PlayerUpdateData](owner.UpdateData).Field47_0 = 0
				spellLifeInform(owner, 15)
				GetServer().S().Audio.EventObj(231, owner, 0, 0)
				Nox_xxx_playerSetState_4FA020(owner, 13)
			}
			p = spellLifeUnlinkBook(p)
		} else {
			p = p.Next
		}
	}
}
