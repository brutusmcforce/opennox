package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

type itemRespawnRecord struct {
	typ            uint32
	obj            *server.Object
	pos            types.Pointf
	direction      uint16
	_              uint16
	frame, pending uint32
	attrs          server.ModifierInitData
	ammo           [2]byte
	_              [2]byte
	next, prev     *itemRespawnRecord
}

// Keep the qualified 386 pool layout, including modifier words and ammunition bytes.





var itemRespawnPool alloc.ClassT[itemRespawnRecord]
var itemRespawnHead *itemRespawnRecord
var itemRespawnEnabled uint32 = 1
var itemRespawnCrown uint32

func itemRespawnInit() int {
	itemRespawnPool = alloc.NewClassT("Respawn", itemRespawnRecord{}, 384)
	return bool2int(itemRespawnPool.Class != nil)
}
func itemRespawnFree() { itemRespawnPool.Free(); itemRespawnPool.Class = nil }
func itemRespawnReset() {
	itemRespawnHead = nil
	itemRespawnPool.FreeAllObjects()
	itemRespawnEnabled = 1
}
func itemRespawnAdd(u *server.Object) uintptr {
	if itemRespawnEnabled == 0 {
		return 0
	}
	p := itemRespawnPool.NewObject()
	if p == nil {
		return 0
	}
	*p = itemRespawnRecord{typ: uint32(u.TypeInd), obj: u, pos: u.PosVec, direction: uint16(u.Direction1)}
	if u.ObjClass&0x13001000 != 0 {
		p.attrs = *u.InitDataModifier()
	}
	if u.ObjClass&0x1000000 != 0 && GetServer().S().Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(u)&0x82 != 0 {
		d := (*[2]byte)(u.UseData.Ptr)
		p.ammo = [2]byte{d[1], d[0]}
	}
	old := itemRespawnHead
	p.next = old
	if old != nil {
		old.prev = p
	}
	itemRespawnHead = p
	return uintptr(unsafe.Pointer(old))
}
func itemRespawnRemove(u *server.Object) {
	for p := itemRespawnHead; p != nil; p = p.next {
		if p.obj != u {
			continue
		}
		if p.prev != nil {
			p.prev.next = p.next
		} else {
			itemRespawnHead = p.next
		}
		if p.next != nil {
			p.next.prev = p.prev
		}
		itemRespawnPool.FreeObjectFirst(p)
		return
	}
}
func itemOwnerSameTeam(a, b *server.Object) bool {
	for x := a; x != nil; x = x.ObjOwner {
		for y := b; y != nil; y = y.ObjOwner {
			if x == y || x.TeamVal.SameAs(&y.TeamVal) {
				return true
			}
		}
	}
	return false
}
func itemDropCrowns(owner *server.Object, stamp uint32) {
	cache := memmap.PtrUint32(0x5D4594, 1568248)
	if *cache == 0 {
		*cache = uint32(GetServer().S().Types.IndByID("Crown"))
	}
	for u := owner.Field129; u != nil; u = u.Field128 {
		if uint32(u.TypeInd) == *cache {
			data := u.UpdateData
			inventoryCrownDrop(owner, u, &owner.PosVec)
			*(*uint32)(unsafe.Add(data, 4)) = stamp
		}
	}
}
func itemRespawnTick() {
	s := GetServer().S()
	if itemRespawnCrown == 0 {
		itemRespawnCrown = uint32(s.Types.IndByID("Crown"))
	}
	if noxflags.HasGame(4608) {
		return
	}
	itemRespawnEnabled = 0
	schedule := func(p *itemRespawnRecord) { p.pending = 1; p.frame = s.Frame() + 30*s.TickRate() }
	allowed := func(typ uint32) bool { return s.Types.ByInd(int(typ)).Allowed() }
	for p := itemRespawnHead; p != nil; p = p.next {
		if p.pending == 0 {
			u := p.obj
			if u == nil {
				schedule(p)
			} else if u.ObjClass&2 != 0 {
				if u.ObjFlags&0x20 != 0 {
					p.obj = nil
					schedule(p)
				} else if u.ObjFlags&0x8000 != 0 {
					schedule(p)
				}
			} else if u.ObjFlags&0x20 != 0 {
				p.obj = nil
				if allowed(uint32(u.TypeInd)) {
					schedule(p)
				}
			} else if u.ObjClass&0x3001000 != 0 || uint32(u.TypeInd) == itemRespawnCrown {
				if u.InvHolder != nil || !allowed(uint32(u.TypeInd)) {
					if u.InvHolder != nil && allowed(uint32(u.TypeInd)) && uint32(u.TypeInd) != itemRespawnCrown && serverConfigFlagsQuery(2) != 0 {
						schedule(p)
					}
				} else if s.Frame() > 5*s.TickRate()+u.Field32 {
					dx, dy := float64(p.pos.X)-float64(u.PosVec.X), float64(p.pos.Y)-float64(u.PosVec.Y)
					if dx*dx+dy*dy > 2500 {
						visibilityFXPoint(129, u.PosVec)
						s.Audio.EventPos(sound.ID(283), u.PosVec, 0, 0)
						Nox_xxx_unitMove_4E7010(u, p.pos)
						if u.ObjClass&0x1000 != 0 {
							effectsRecharge(u, 100)
						} else if u.ObjClass&0x1000000 != 0 && s.Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(u)&0x82 != 0 {
							d := (*[2]byte)(u.UseData.Ptr)
							d[1], d[0] = p.ammo[0], p.ammo[1]
						}
						if u.HealthData != nil {
							Nox_xxx_unitSetHP_4E4560(u, u.HealthData.Max)
						}
						visibilityFXPoint(129, p.pos)
						s.Audio.EventPos(sound.ID(283), p.pos, 0, 0)
					}
				}
			} else if u.InvHolder != nil {
				schedule(p)
			}
			if p.pending == 0 {
				continue
			}
		}
		if s.Frame() < p.frame || !allowed(p.typ) {
			continue
		}
		next := s.NewObjectByTypeInd(int(p.typ))
		if next != nil {
			GetServer().CreateObjectAt(next, nil, p.pos)
			visibilityFXPoint(129, p.pos)
			next.Direction1, next.Direction2 = server.Dir16(p.direction), server.Dir16(p.direction)
			if next.ObjClass&0x13001000 != 0 {
				stateAttributes(next, &p.attrs)
			}
			if next.ObjClass&0x1000000 != 0 && s.Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(next)&0x82 != 0 {
				d := (*[2]byte)(next.UseData.Ptr)
				d[1], d[0] = p.ammo[0], p.ammo[1]
			}
			s.Audio.EventObj(sound.ID(283), next, 0, 0)
		}
		if old := p.obj; old != nil && old.ObjClass&2 != 0 {
			GetServer().DelayedDelete(old)
		}
		p.pending = 0
		p.obj = next
	}
}
