package legacy

import (
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

// Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch handles player actions.
// The caller checks the consumed length afterward, so reject incomplete actions
// before reading their fields. The explicit player and update
// arguments remain distinct from the unit's own update data.
func Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(_ ntype.PlayerInd, data []byte, player *server.Player, unit *server.Object, update unsafe.Pointer) int {
	if n := serverActionLength(data); n == 0 || len(data) < n {
		return -1
	}
	ud := (*server.PlayerUpdateData)(update)
	code := func(off int) uint32 { return networkDynamicUnitCode(uint32(binary.LittleEndian.Uint16(data[off:]))) }
	switch data[0] {
	case 64:
		controlSetWaypoint(unit, math.Float32bits(float32(binary.LittleEndian.Uint16(data[3:]))), math.Float32bits(float32(binary.LittleEndian.Uint16(data[5:]))))
		return 7
	case 114:
		id := code(1)
		if ud.Player.Field3680&3 == 0 && ud.Trade70 == nil && ud.DialogWith == nil && unit.ObjFlags&2 == 0 {
			if item := controlEquippedByCode(unit, id); item != nil {
				pos := types.Pointf{X: float32(binary.LittleEndian.Uint16(data[3:])), Y: float32(binary.LittleEndian.Uint16(data[5:]))}
				inventoryTargetDrop(unit, item, &pos)
			}
		}
		return 7
	case 115:
		id := code(1)
		if !Nox_xxx_gameGet_4DB1B0() && ud.Player.Field3680&3 == 0 && ud.Trade70 == nil && ud.DialogWith == nil && unit.ObjFlags&2 == 0 {
			if item := objectLookupByNetCode(id); item != nil {
				var weight int32
				for it := unit.InvFirstItem; it != nil; it = it.InvNextItem {
					weight += int32(it.Weight)
				}
				if weight+int32(item.Weight) <= int32(unit.CarryCapacity) {
					serverActionPickup(unit, item)
				} else {
					gameplayTextPrivate(unit, alloc.InternCString("pickup.c:CarryingTooMuch"), 0)
				}
			}
		}
		return 3
	case 116, 117:
		id := code(1)
		if !Nox_xxx_gameGet_4DB1B0() && ud.Player.Field3680&3 == 0 {
			if item := controlEquippedByCode(unit, id); item != nil {
				if data[0] == 116 {
					effectsUse(unit, item)
				} else {
					equipmentTryEquip(unit, item)
				}
			}
		}
		return 3
	case 118:
		id := code(1)
		if ud.Player.Field3680&3 == 0 {
			if item := controlEquippedByCode(unit, id); item != nil && !(ud.State == 1 && item.ObjClass&0x1000000 != 0 && item.ObjSubClass&8 != 0) {
				equipmentTryDequip(unit, item)
			}
		}
		return 3
	case 120:
		id := code(1)
		if ud.Player.Field3680&1 == 0 {
			if binary.LittleEndian.Uint16(data[1:]) == 0 {
				monsterOrder(unit, nil, int(data[3]))
			} else if target := objectLookupByNetCode(id); target != nil {
				monsterOrder(unit, target, int(data[3]))
			}
		}
		return 4
	case 121:
		serverActionSpell(data, unit, ud)
		return 22
	case 123:
		id := code(1)
		if ud.Player.Field3680&3 == 0 && ud.Trade70 == nil && ud.DialogWith == nil {
			if target := objectLookupByNetCode(id); target != nil && target.Collide != nil {
				target.CallCollide(int(uintptr(unit.CObj())), 0)
			}
		}
		return 3
	case 165:
		if data[1] >= 255 {
			// the player only has 255 alias slots; index 255 would overwrite the player unit pointer
			return -1
		}
		// Preserve the shared layout and raw byte writes of the original alias table.
		dst := unsafe.Slice((*byte)(unsafe.Add(player.C(), 16+8*int(data[1]))), 8)
		copy(dst, data[2:10])
		return 10
	case 224:
		id := code(1)
		var target *server.Object
		if binary.LittleEndian.Uint16(data[1:]) != 0 {
			target = objectLookupByNetCode(id)
		}
		equipmentSecondary(unit, target)
		return 3
	case 226:
		id := code(1)
		item := controlEquippedByCode(unit, id)
		if item == nil {
			if s := unit.UpdateDataPlayer().Trade70; s != nil {
				if n := shopFind(s.Stock, id); n != nil {
					item = n.Object
				}
			}
		}
		if item == nil {
			item = objectLookupByNetCode(id)
		}
		if item != nil {
			// only readable objects (books) have the expected use data; the net code comes from the client
			if item.ObjClass&0x40000 == 0 || item.UseData.Ptr == nil {
				return 4
			}
			out := []byte{226, 0, 0, 0}
			binary.LittleEndian.PutUint16(out[1:], uint16(GetServer().S().GetUnitNetCode(item)))
			if data[3] == 2 {
				out[3] = byte(bookGuideID(alloc.GoString((*byte)(item.UseData.Ptr))))
			} else {
				out[3] = *(*byte)(item.UseData.Ptr)
			}
			reliableEnqueue(int(ud.Player.PlayerInd), out, nil, 1, 0)
		}
		return 4
	case 238:
		switch data[1] {
		case 0, 1, 2, 3:
			kind := uint32(data[1] & 1)
			if kind == 0 && noxflags.HasGame(4096) {
				kind = 3
			}
			name := (*uint16)(unsafe.Pointer(&data[2]))
			if data[1] < 2 {
				voteCast(kind, unit, name)
			} else {
				voteWithdraw(kind, unit, name)
			}
			return 52
		case 4:
			voteCast(2, unit, nil)
			return 2
		case 5:
			voteWithdraw(2, unit, nil)
			return 2
		default:
			return -1
		}
	case 201:
		return serverActionTrade(data, unit, ud)
	case 240:
		switch data[1] {
		case 3:
			if target := ud.Player.PlayerUnit; target != nil && target.ObjFlags&0x8000 != 0 {
				ud.Field137 = 0
				controlRespawn(target)
			}
		case 27:
			Sub_4DD0B0(unit)
		default:
			return -1
		}
		return 2
	case 241:
		if item := controlEquippedByCode(unit, uint32(binary.LittleEndian.Uint16(data[1:]))); item != nil {
			inventoryDrop(unit, item, &unit.PosVec)
			gameplayTextPrivate(unit, alloc.InternCString("pickup.c:CarryingTooMuch"), 0)
		}
		return 3
	default:
		return -1
	}
}

func serverActionLength(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	switch data[0] {
	case 64, 114:
		return 7
	case 115, 116, 117, 118, 123, 224, 241:
		return 3
	case 120, 226:
		return 4
	case 121:
		return 22
	case 165:
		return 10
	case 240:
		return 2
	case 238:
		if len(data) < 2 {
			return 2
		}
		if data[1] < 4 {
			return 52
		}
		return 2
	case 201:
		if len(data) < 2 {
			return 2
		}
		switch data[1] {
		case 15, 16, 21, 22, 24, 26, 28, 30:
			return 4
		case 23, 25:
			return 5
		default:
			return 2
		}
	default:
		return 0
	}
}

func serverActionPickup(unit, item *server.Object) {
	class := byte(unit.UpdateDataPlayer().Player.PlayerClass())
	typ := byte(item.TypeInd)
	if byte(uint32(item.ObjClass)>>16) == 17 && (typ == 106 && class != 1 || (typ == 107 || typ == 109) && class != 0) {
		gameplayTextPrivate(unit, alloc.InternCString("pickup.c:ObjectEquipClassFail"), 0)
		inventorySound(925, unit, 2, int(unit.NetCode))
		return
	}
	Nox_xxx_inventoryServPlace_4F36F0(unit, item, 1, 1)
}

func serverActionSpell(data []byte, unit *server.Object, ud *server.PlayerUpdateData) {
	allowed := true
	if ud.Player.Field3680&1 != 0 {
		gameplayTextPrivate(unit, alloc.InternCString("GeneralPrint:NoSpellWarningGeneral"), 0)
		allowed = false
	}
	if ud.Player.Field3680&2 != 0 {
		gameplayTextPrivate(unit, alloc.InternCString("GeneralPrint:ConjureNoSpellWarning1"), 0)
		allowed = false
	}
	if !noxflags.HasGame(2048) && unit.ObjFlags&0x4000 != 0 {
		allowed = false
	}
	if noxflags.HasGame(128) || !allowed {
		return
	}
	var ids [5]int32
	var count int32
	for i := range ids {
		ids[i] = int32(binary.LittleEndian.Uint32(data[1+4*i:]))
		if ids[i] != 0 {
			count++
		}
	}
	s := GetServer().S()
	if count == 1 && s.Spells.HasFlags(spell.ID(ids[0]), 32) && ud.CursorObj != nil && !s.IsEnemyTo(unit, ud.CursorObj) && !noxflags.HasGame(4096) {
		return
	}
	if spellLifeInsertBook(unit, unsafe.Pointer(&ids[0]), count, 3, int32(data[21])) == 0 && count == 1 {
		for _, id := range ids {
			if id != 0 {
				gameplayReportSpellStat(int(ud.Player.PlayerInd), uint32(id), 0)
			}
		}
	}
}

func serverActionTrade(data []byte, unit *server.Object, ud *server.PlayerUpdateData) int {
	switch data[1] {
	case 14:
		if ud.Trade70 != nil {
			shopCancelTrade(ud.Trade70)
		}
		return 2
	case 15:
		id := networkDynamicUnitCode(uint32(binary.LittleEndian.Uint16(data[2:])))
		if item := controlEquippedByCode(unit, id); item != nil && ud.Trade70 != nil {
			if tradeAddOffer(ud.Trade70, unit, item) == 1 {
				inventoryRemove(unit, item)
			}
		}
		return 4
	case 16:
		id := networkDynamicUnitCode(uint32(binary.LittleEndian.Uint16(data[2:])))
		if ud.Trade70 != nil {
			shopWithdraw(ud.Trade70, id)
		}
		return 4
	case 17:
		if ud.Trade70 != nil {
			shopAccept(ud.Trade70, unit)
		}
		return 2
	case 18:
		if ud.Trade70 != nil {
			shopExit(ud.Trade70)
		}
		return 2
	case 21:
		if !Nox_xxx_gameGet_4DB1B0() && ud.Player.Field3680&3 == 0 {
			id := networkDynamicUnitCode(uint32(binary.LittleEndian.Uint16(data[2:])))
			if target := objectLookupByNetCode(id); target != nil && target.ObjSubClass&8 != 0 {
				tradeStart(unit, target)
			}
		}
		return 4
	case 22:
		if ud.Trade70 != nil {
			tradeBuy(unit, ud.Trade70, uint32(binary.LittleEndian.Uint16(data[2:])))
		}
		return 4
	case 23:
		if ud.Trade70 != nil {
			tradeBuyMany(unit, ud.Trade70, int32(binary.LittleEndian.Uint16(data[2:])), uint32(data[4]))
		}
		return 5
	case 24:
		if ud.Trade70 != nil {
			tradeSell(unit, ud.Trade70, uint32(binary.LittleEndian.Uint16(data[2:])))
		}
		return 4
	case 25:
		if ud.Trade70 != nil {
			shopSell(unit, ud.Trade70, int32(binary.LittleEndian.Uint16(data[2:])), uint32(data[4]))
		}
		return 5
	case 26:
		if ud.Trade70 != nil {
			shopRepair(unit, ud.Trade70, uint32(binary.LittleEndian.Uint16(data[2:])))
		}
		return 4
	case 28:
		if ud.Trade70 != nil {
			tradeSaleQuote(unit, ud.Trade70, uint32(binary.LittleEndian.Uint16(data[2:])))
		}
		return 4
	case 30:
		if ud.Trade70 != nil {
			shopRepairQuote(unit, ud.Trade70, uint32(binary.LittleEndian.Uint16(data[2:])))
		}
		return 4
	default:
		return -1
	}
}
