package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

func playerFileAbilities(u *server.Object, r objectXferStream) {
	core := GetServer().S()
	first := r.byte(byte(bool2int(core.Abils.IsActive(u, 1))))
	if r.read() && first == 1 {
		Sub_4FC670(1)
	}
	fourth := r.byte(byte(bool2int(core.Abils.IsActive(u, 4))))
	remaining := r.word(uint32(core.Abils.Sub4FC030(u, 4)))
	if r.read() && fourth == 1 {
		Nox_xxx_playerExecuteAbil_4FBB70(u, 4)
		core.Abils.Sub4FC070(u, 4, int(int32(remaining)))
	}
	start := 1
	if first == 1 {
		start = 2
	}
	for i := start; i < 6; i++ {
		abil := server.Ability(i)
		cd := r.word(uint32(core.Abils.GetCooldown(u, abil)))
		if r.read() {
			core.Abils.SetCooldown(u, abil, int(int32(cd)))
			if cd != 0 {
				Nox_xxx_netAbilRepotState_4D8100(u, abil, 0)
			}
		}
	}
}
func playerFileEnchantment(u *server.Object) int {
	p := u.UpdateDataPlayer().Player
	r := playerFileStream()
	version := int16(r.short(5))
	if version > 5 {
		return 0
	}
	if r.read() && noxflags.HasGame(2048) {
		Nox_xxx_spellCastByPlayer_4FEEF0()
		GetServer().S().Spells.Dur.Sub4FE8A0(0)
	}
	present := playerFilePresent(r, noxflags.HasGame(8192))
	if present {
		if !noxflags.HasGame(2048) {
			return 0
		}
		if r.read() {
			count := r.byte(0)
			for i := 0; i < int(count); i++ {
				name := playerFileName(r, "")
				enchant, ok := server.ParseEnchant(name)
				id := int32(enchant)
				if !ok {
					id = -1
				}
				timer := r.short(0)
				power := byte(2)
				if version >= 2 {
					power = r.byte(power)
				}
				arg := server.SpellAcceptArg{Obj: u, Pos: u.PosVec}
				GetServer().Nox_xxx_spellAccept4FD400(server.EnchantID(id).Spell(), u, u, u, &arg, int(power))
				if timer == 0 {
					timer = uint16(GetServer().S().TickRate())
				}
				u.BuffsDur[int(id)] = timer
				if id == 26 && version >= 3 {
					duration := int32(r.word(0))
					if rec := spellLifeFindDuration(51, u); rec != nil {
						rec.Field72 = duration
					}
				}
			}
		} else {
			r.byte(byte(bookEnchantCountActive(u)))
			for id := int32(bookEnchantFirst()); id != -1; id = int32(bookEnchantNext(int32(id))) {
				if !spellLifeHasBuff(u, int32(int8(id))) {
					continue
				}
				playerFileName(r, server.EnchantID(id).String())
				r.raw(unsafe.Pointer(&u.BuffsDur[int(id)]), 2)
				r.byte(byte(spellLifeBuffPower(u, id)))
				if id == 26 {
					duration := uint32(100)
					if rec := spellLifeFindDuration(51, u); rec != nil {
						duration = uint32(rec.Field72)
					}
					r.word(duration)
				}
			}
		}
		if version >= 5 && p.InfoBytes[66] == 0 {
			playerFileAbilities(u, r)
		}
	}
	// Version four's tail is outside the optional enchantment body.
	if version == 4 && p.InfoBytes[66] == 0 {
		playerFileAbilities(u, r)
	}
	return 1
}
