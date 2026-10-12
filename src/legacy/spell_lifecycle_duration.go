package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/server"
	"runtime"
)

func spellLifeCancelDurations(mode int32) int32 {
	d := &GetServer().S().Spells.Dur
	for p := d.List; p != nil; {
		next := p.Next
		if mode != 1 || p.Target48 == nil || p.Target48.ObjClass&4 != 4 {
			d.CancelSpell(p)
		}
		p = next
	}
	return 0
}
func spellLifeCancelPlayer(u *server.Object) int32 {
	d := &GetServer().S().Spells.Dur
	for p := d.List; p != nil; {
		next := p.Next
		if p.Caster16 == u {
			d.CancelSpell(p)
		}
		p = next
	}
	return 0
}
func spellLifeCancelWand(u, item *server.Object) {
	if item.ObjClass&0x1000 == 0 {
		return
	}
	d := &GetServer().S().Spells.Dur
	if item.ObjSubClass&0x40000 != 0 {
		d.CancelFor(spell.ID(43), u)
	}
	if item.ObjSubClass&0x4000000 != 0 {
		d.CancelFor(spell.ID(59), u)
	}
}
func spellLifeCancelSelected(u *server.Object) {
	d := &GetServer().S().Spells.Dur
	for p := d.List; p != nil; {
		next := p.Next
		if p.Caster16 == u {
			switch p.Spell {
			case 24, 43, 35, 8, 22, 59, 67:
				d.CancelSpell(p)
			}
		}
		p = next
	}
}
func spellLifeFindDuration(id int32, u *server.Object) *server.DurSpell {
	for p := GetServer().S().Spells.Dur.List; p != nil; p = p.Next {
		if p.Flags88&1 == 0 && p.Spell == uint32(id) && p.Target48 != nil && p.Target48 == u {
			return p
		}
	}
	return nil
}
func spellLifeRayMessage(p *server.DurSpell) uint32 {
	result := p.Spell - 7
	var buf [7]byte
	buf[0] = 158
	buf[2] = byte(p.Level)
	switch p.Spell {
	case 7:
		buf[1] = 3
	case 9:
		buf[1] = 2
	case 22:
		buf[1] = 5
	case 24:
		buf[1] = 4
	case 35:
		result = controlRaw(p.Target48)
		if p.Caster16 == p.Target48 {
			return result
		}
		buf[1] = 6
		binary.LittleEndian.PutUint16(buf[3:], uint16(nox_xxx_netGetUnitCodeServ_578AC0(asObjectC(p.Target48))))
		binary.LittleEndian.PutUint16(buf[5:], uint16(nox_xxx_netGetUnitCodeServ_578AC0(asObjectC(p.Caster16))))
		return spellLifeSendRay(p, buf)
	case 43:
		for sub := p.Sub108; sub != nil; sub = sub.Next {
			result = spellLifeRayMessage(sub)
		}
		return result
	case 59:
		buf[1] = 1
		buf[2] = byte(p.Caster16.Direction1)
	default:
		return result
	}
	result = controlRaw(p.Target48)
	if p.Target48 == nil {
		return result
	}
	binary.LittleEndian.PutUint16(buf[5:], uint16(nox_xxx_netGetUnitCodeServ_578AC0(asObjectC(p.Target48))))
	binary.LittleEndian.PutUint16(buf[3:], uint16(nox_xxx_netGetUnitCodeServ_578AC0(asObjectC(p.Caster16))))
	return spellLifeSendRay(p, buf)
}
func spellLifeSendRay(p *server.DurSpell, buf [7]byte) uint32 {
	reliableEnqueue(255, buf[:], nil, 1, 1)
	runtime.KeepAlive(buf)
	playerStateMark(p.Caster16, 2)
	playerStateMark(p.Target48, 2)
	return 0
}
