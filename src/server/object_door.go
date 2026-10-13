package server

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/object"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
)

// DoorUpdateData is the update data of doors and gates (DoorUpdate).
type DoorUpdateData struct {
	Field0      byte    // 0, 0
	KeyType     byte    // 0, 1: required key (1 silver, 2 gold, 3 ruby, 4 sapphire), 5 means locked by a mechanism; 0 if unlocked
	Field2      [2]byte // 0, 2
	BaseDir     uint32  // 1, 4: closed direction (0-31)
	PrevDir     uint32  // 2, 8: direction at the last sync
	Dir         uint32  // 3, 12: current direction (0-31)
	GridX       uint32  // 4, 16: wall grid X of the door
	GridY       uint32  // 5, 20: wall grid Y of the door
	Field24     uint32  // 6, 24
	AngleQueued uint32  // 7, 28: set while the door is linked into the angle queue
	Torque      float32 // 8, 32: accumulated angular force
	AngleNext   uint32  // 9, 36: next door data in the angle queue (as a memmap word)
	Angle       int16   // 10, 40: current angle (0-255)
	Field42     uint16  // 10, 42
	HitFrame    uint32  // 11, 44: frame of the last collision
	Notified    byte    // 12, 48: set when the door was opened in quest (sent to joining players)
	Field49     [3]byte // 12, 49
}

var _ = [1]struct{}{}[52-unsafe.Sizeof(DoorUpdateData{})]

type serverDoors struct {
	flagXxx bool

	keyHolder    *Object
	keyHolderSet bool
}

func (s *serverDoors) Sub_4D72C0() bool {
	return s.flagXxx
}

func (s *serverDoors) Sub_4D72B0(a1 bool) {
	s.flagXxx = a1
}

func (s *Server) PlayersHaveSilverKey() *Object {
	var found *Object
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		if !pl.IsActive() {
			continue
		}
		u := pl.PlayerUnit
		if u == nil {
			continue
		}
		cnt := 0
		for it := u.FirstItem(); it != nil; it = it.NextItem() {
			if int(it.TypeInd) == s.Types.SilverKeyID() {
				cnt++
			}
		}
		if cnt > 0 {
			found = u
		}
	}
	if found == nil {
		return nil
	}
	for it := found.FirstItem(); it != nil; it = it.NextItem() {
		if int(it.TypeInd) == s.Types.SilverKeyID() {
			return it
		}
	}
	return nil
}

func (s *Server) DoorCheckKey(u, door *Object) *Object {
	ud2 := door.UpdateDataDoor()
	if ud2.KeyType == 5 {
		return nil
	}
	if door.ObjOwner != nil {
		return nil
	}
	var found *Object
	for it := u.FirstItem(); it != nil; it = it.NextItem() {
		if !it.Class().Has(object.ClassKey) {
			continue
		}
		tname := s.Types.ByInd(int(it.TypeInd)).ID()
		exp := ""
		switch ud2.KeyType {
		case 1:
			exp = "SilverKey"
		case 2:
			exp = "GoldKey"
		case 3:
			exp = "RubyKey"
		case 4:
			exp = "SapphireKey"
		}
		if exp != "" && tname == exp {
			found = it
			break
		}
	}
	if found == nil && u.Class().Has(object.ClassPlayer) && noxflags.HasGame(noxflags.GameModeQuest) &&
		s.Doors.Sub_4D72C0() && ud2.KeyType == 1 {
		return s.PlayersHaveSilverKey()
	}
	return found
}
