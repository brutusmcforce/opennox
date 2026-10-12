package legacy

import (
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func mapPolygonCall(cb *server.ScriptCallback, u *server.Object, event int) {
	if cb.Func != -1 {
		GetServer().NoxScriptC().ScriptCallback(cb, u, nil, server.ScriptEventType(event))
	}
}
func mapPolygonPlayerLeave(u *server.Object, pl unsafe.Pointer) {
	cache := equipmentWord(pl, 3664)
	if *cache != 0 && *cache != mapPolygonUnset {
		mapPolygonCall(&mapPolygonAt(*cache).Leave, u, 27)
		*cache = 0
		*(*byte)(unsafe.Add(pl, 3668)) = 1
	}
}
func mapPolygonPlayer(u *server.Object) {
	if u == nil || uint32(u.ObjClass)&4 == 0 {
		return
	}
	player := u.UpdateDataPlayer().Player
	pl := unsafe.Pointer(player)
	cache := equipmentWord(pl, int(layoutServerPlayer.Off(3664)) /* field3664 */)
	var p *mapPolygon
	if player.PlayerInd == 31 {
		index := *equipmentWord(pl, int(layoutServerPlayer.Off(3660)) /* field3660 */)
		if index == 0 {
			*cache = 0
			mapPolygonPlayerLeave(u, pl)
			return
		}
		if index == mapPolygonUnset {
			mapPolygonPlayerLeave(u, pl)
			return
		}
		p = mapPolygonAt(index)
	} else {
		if *cache != mapPolygonUnset && u.PosVec == u.PrevPos {
			return
		}
		point := [2]int32{floatToInt32(u.PosVec.X), floatToInt32(u.PosVec.Y)}
		p = mapPolygonFind(&point, *cache, false)
		if p == nil {
			if *cache == 0 || *cache == mapPolygonUnset {
				mapPolygonPlayerLeave(u, pl)
				return
			}
			p = mapPolygonAt(*cache)
		}
	}
	if *cache == p.ID {
		return
	}
	if *cache != mapPolygonUnset {
		if *cache != 0 {
			mapPolygonCall(&mapPolygonAt(*cache).Leave, u, 29)
		}
		bit := uint32(1) << uint(player.PlayerInd&31)
		if p.Visited&bit == 0 && p.Flags&1 != 0 && noxflags.HasGame(4096) {
			questRuntimeIncrement(u, 4672, 16)
			gameplayTextPrivate(u, alloc.InternCString("GeneralPrint:SecretFound"), 0)
			GetServer().S().Audio.EventObj(904, u, 0, 0)
			players := &GetServer().S().Players
			for other := players.FirstUnit(); other != nil; other = players.NextUnit(other) {
				if other != u {
					gameplayTextInformation(int(other.UpdateDataPlayer().Player.PlayerInd), 20, unsafe.Pointer(&u.NetCode))
				}
			}
			p.Flags &^= 1
		}
		p.Visited |= bit
		mapPolygonCall(&p.Enter, u, 28)
	}
	*cache = p.ID
	*(*byte)(unsafe.Add(pl, int(layoutServerPlayer.Off(3668)) /* field3668 */)) = p.Level
}
func mapPolygonMonster(u *server.Object) {
	if u == nil || uint32(u.ObjClass)&2 == 0 {
		return
	}
	cache := (*uint32)(u.UpdateData)
	if *cache != mapPolygonUnset && u.PosVec == u.PrevPos {
		return
	}
	point := [2]int32{floatToInt32(u.PosVec.X), floatToInt32(u.PosVec.Y)}
	var p *mapPolygon
	if *cache == mapPolygonUnset {
		p = mapPolygonFind(&point, 0, false)
	} else {
		p = mapPolygonFind(&point, *cache, true)
	}
	if p != nil {
		if *cache != p.ID {
			if *cache != mapPolygonUnset {
				if *cache != 0 {
					mapPolygonCall(&mapPolygonAt(*cache).Leave, u, 26)
				}
				mapPolygonCall(&p.Enter, u, 25)
			}
			*cache = p.ID
		}
	} else if *cache != 0 && *cache != mapPolygonUnset {
		old := mapPolygonAt(*cache)
		if old.Active != 0 {
			mapPolygonCall(&old.Leave, u, 24)
		}
		*cache = 0
	}
}
func mapPolygonActorInit(p unsafe.Pointer) unsafe.Pointer {
	if p != nil {
		*equipmentWord(p, 3660) = mapPolygonUnset
		*equipmentWord(p, 3664) = mapPolygonUnset
	}
	return p
}
func mapPolygonColor() {
	dr := minimapLocalDrawable()
	if dr == nil {
		return
	}
	point := GetClient().Viewport().World.Max
	player := GetServer().S().Players.ByID(int(dr.NetCode32))
	if player == nil {
		return
	}
	pl := unsafe.Pointer(player)
	cache := equipmentWord(pl, int(layoutServerPlayer.Off(3660)) /* field3660 */)
	if mapPolygonNext > 1 {
		if *cache != mapPolygonUnset && int32(point.X) == memmap.Int32(0x5D4594, 811364) && int32(point.Y) == memmap.Int32(0x5D4594, 811368) {
			return
		}
		pt := [2]int32{int32(point.X), int32(point.Y)}
		p := mapPolygonFind(&pt, *cache, false)
		if p != nil {
			if *cache != p.ID {
				*(*byte)(unsafe.Add(pl, int(layoutServerPlayer.Off(3668)) /* field3668 */)) = p.Level
				*cache = p.ID
				GetClient().R2().Data().SetLightColor(noxrender.RGB{R: int(p.Color[0]), G: int(p.Color[1]), B: int(p.Color[2])})
			}
			return
		}
	}
	*(*byte)(unsafe.Add(pl, int(layoutServerPlayer.Off(3668)) /* field3668 */)) = 1
	*cache = 0
	GetClient().R2().Data().SetLightColor(noxrender.RGB{R: int(memmap.Uint32(0x587000, 142296)), G: int(memmap.Uint32(0x587000, 142300)), B: int(memmap.Uint32(0x587000, 142304))})
}
