package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/server"
)

func objectReportInit() int {
	for i, name := range []string{"TeleportPentagram", "PressurePlate", "Spike", "PeriodicSpike"} {
		id := uint32(GetServer().S().Types.IndByID(name))
		*memmap.PtrUint32(0x5D4594, 2386972+uintptr(4*i)) = id
		if id == 0 {
			return 0
		}
	}
	return 1
}
func objectReportSprite(viewer *server.Object, to int, u *server.Object) int {
	code := uint16(GetServer().S().GetUnitNetCode(u))
	b := []byte{57, byte(code), byte(code >> 8), 0}
	cls := uint32(u.ObjClass)
	typ := uint32(u.TypeInd)
	switch {
	case cls&0x400000 != 0 && uint32(u.ObjSubClass)&0x18 != 0:
		b[0], b[3] = 179, *(*byte)(u.UpdateData)
		return gameplayReportSend(to, b, true, 1)
	case typ == memmap.Uint32(0x5D4594, 2386972): // TeleportPentagram
		b[3] = u.UpdateDataTransporter().Anim
	case typ == memmap.Uint32(0x5D4594, 2386980):
		b[3] = ^byte(uint32(u.ObjFlags)>>24) & 1
	case typ == memmap.Uint32(0x5D4594, 2386976):
		b[0], b[3] = 180, *(*byte)(u.UpdateData)&1
	case cls&0x4000 != 0:
		b[3] = byte(u.UpdateDataElevator().Field_4) >> 2
	case cls&0x8000 != 0:
		if u.UpdateData != nil {
			if link := u.UpdateDataElevatorShaft().Partner; link != nil {
				b[3] = byte(link.UpdateDataElevator().Field_4) >> 2
			}
		}
	case cls&0x80 != 0:
		b[0], b[3] = 178, byte(u.UpdateDataDoor().Dir)
	default:
		visibilitySpecialUpdate(viewer, u)
		return 0
	}
	return bool2int(GetServer().S().NetList.AddToMsgListCli(ntype.PlayerInd(to), netlist.Kind1, b))
}

// The original entry returns only the signed low byte of its branch result.
func objectReportRecipient(viewer, u *server.Object) int {
	pl := viewer.UpdateDataPlayer().Player
	to := int(pl.PlayerInd)
	bit := uint32(1) << uint(to)
	ret := uint32(byte(u.ObjFlags))
	if uint32(u.ObjFlags)&0x20 != 0 || uint32(u.ObjClass)&0x40000000 != 0 {
		return int(int8(ret))
	}
	if u != viewer && u.Field35&bit != 0 {
		gameplayReportFriend(to, u, bool2int(u.Field36&bit != 0))
		ret = ^bit
		u.Field35 &^= bit
	}
	if u.Field38&bit == 0 && uint32(u.ObjClass)&1 == 0 {
		return int(int8(ret))
	}
	if uint32(u.ObjFlags)&0x800 != 0 && !u.HasOwner(viewer) && !controlFlags(512) {
		return 0
	}
	ret = uint32(u.ObjClass)
	visible := ret&0x20400000 != 0
	if !visible {
		ret = uint32(Nox_xxx_playerMapTracksObj_4173D0(to, u))
		visible = ret != 0
		if !visible {
			visible = GetServer().S().MapTraceRay(pl.Pos3632(), u.PosVec, server.MapTraceFlags(69))
			ret = uint32(bool2int(visible))
		}
	}
	if visible && u != viewer && u.IsHidden() {
		visible = false
	}
	if visible {
		if u.Field37&bit != 0 {
			if u.Field5&0x20 != 0 {
				return int(int8(ret))
			}
		} else if u.Field140[to]&0xfff != 0 {
			u.Field140[to] |= (u.Field140[to] & 0xfff) << 16
		}
		sent := 0
		cls := uint32(u.ObjClass)
		switch {
		case cls&0x400000 != 0:
			sent = objectReportSprite(viewer, to, u)
		case cls&0x200000 != 0:
			switch {
			case cls&2 != 0:
				sent = objectReportMonster(to, u, 1)
			case cls&4 != 0:
				sent = objectReportPlayer(viewer, u, 1, 1)
			default:
				sent = objectReportPhantom(to, u)
			}
		case cls&0x100000 != 0:
			sent = objectReportSimple(to, u)
		default:
			u.Field38 = 0
		}
		ret = uint32(visibilitySpecialUpdate(viewer, u))
		if sent != 0 {
			u.Field38 &^= bit
			u.Field37 |= bit
			ret = u.Field37
		}
	} else if u.Field37&bit != 0 {
		if uint32(u.ObjClass)&6 != 0 {
			visibilityOutOfSight(to, u)
		} else {
			visibilityInShadows(to, u)
		}
		u.Field37 &^= bit
		u.Field38 |= bit
		ret = u.Field37
	}
	return int(int8(ret))
}

// The count helper was private to this scheduling caller in C.
func objectReportMinimapCount(to int) int {
	if to < 0 || to >= 32 {
		return 0
	}
	head := GetServer().S().Players.ByIndRaw(ntype.PlayerInd(to)).Field4580
	if head == nil {
		return 0
	}
	count := 1
	for p := head.Field8; p != head; p = p.Field8 {
		count++
	}
	return count
}
func objectReportSchedule(data *server.PlayerUpdateData) int {
	count := objectReportMinimapCount(int(data.Player.PlayerInd))
	if count == 0 {
		return 0
	}
	return bool2int(count > 60 || GetServer().S().Frame()-data.Field67 > uint32(60/count))
}
