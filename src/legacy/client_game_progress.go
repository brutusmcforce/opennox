package legacy

import (
	"github.com/brutusmcforce/opennox/server"
	"encoding/binary"
	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"
	"github.com/brutusmcforce/opennox/libs/strman"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

var clientGameMapFrame uint32

func clientGameProgressString(key string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(key), "cdecode.c")
}
func clientGameProgressText(key string, args ...any) string {
	return consoleCommandFormat(clientGameProgressString(key), args...)
}

func clientGameProgress(op netmsg.Op, data []byte) (int, bool) {
	size := 0
	switch int(op) {
	case 72:
		size = 14
	case 75, 77, 110:
		size = 5
	case 76:
		size = 9
	case 78:
		size = 11
	case 85, 101, 103, 106:
		size = 7
	case 86, 87, 88, 89:
		return clientGameWinner(byte(op), data), true
	case 96, 97, 109:
		size = 3
	case 108:
		size = 5
	case 111:
		size = 4
	default:
		return 0, false
	}
	word := func(off int) uint16 { return binary.LittleEndian.Uint16(data[off:]) }
	dword := func(off int) uint32 { return binary.LittleEndian.Uint32(data[off:]) }
	if int(op) == 106 {
		if p := GetServer().S().Players.ByID(int(word(1))); p != nil {
			if !noxflags.HasGame(noxflags.GameHost) {
				playerStateRemoveStatus(p, 0x423)
				playerStateAddStatus(p, dword(3)&0x423)
			}
			if !noxflags.HasEngine(noxflags.EngineNoRendering) && int(word(1)) == ClientPlayerNetCode() {
				Nox_client_onClientStatusA(int(byte(p.Field3680)))
				sub_470C40(int(p.Field3680 >> 10 & 1))
			}
		}
		return size, true
	}
	if !Nox_client_isConnected() {
		return size, true
	}
	code := word(1) & 0x7fff
	switch int(op) {
	case 72:
		if int(code) == ClientPlayerNetCode() {
			if p := Get_dword_8531A0_2576(); p != nil && !noxflags.HasGame(noxflags.GameHost) {
				for _, pair := range [][2]int{{2247, 3}, {2243, 5}, {2235, 9}, {2239, 11}} {
					*(*uint32)(unsafe.Add(p.C(), pair[0])) = uint32(word(pair[1]))
				}
				*(*uint16)(unsafe.Add(p.C(), int(layoutServerPlayer.Off(3652)) /* field3652 */)) = word(7)
				*(*byte)(unsafe.Add(p.C(), int(layoutServerPlayer.Off(3684)) /* Level */)) = data[13]
			}
			nox_xxx_inventoryNameSignInit_4671E0()
		}
	case 75, 76:
		mods := server.ModifierInitData{Field16: 0xffffffff}
		if int(op) == 76 {
			for i := 0; i < 4; i++ {
				mods.Modifiers[i] = (*server.ModifierEff)(nox_xxx_modifGetDescById_413330(int32(data[5+i])))
			}
		}
		if nox_xxx_spritePickup_461660(int32(code), int32(word(3)), &mods) == 0 {
			nox_xxx_send2ServInvenFail_461630(int16(code))
		}
	case 77:
		sub_461A80(int32(code))
	case 78:
		if p := GetServer().S().Players.ByID(int(code)); p != nil {
			if !noxflags.HasGame(noxflags.GameHost) {
				p.Lessons = int32(dword(3))
				p.Field2140 = dword(7)
				p.Field2144 = gameFrame()
			}
			if noxflags.HasGame(1024) && int32(dword(7)) >= int32(memmap.Uint16(0x5D4594, 371434)) {
				audioEventPlay(312, 100, 0, 0)
				Nox_xxx_printCentered_445490(clientGameProgressText("Eliminated", p.Name()))
			}
		}
	case 85:
		if p := GetServer().S().Players.ByID(int(code)); p != nil {
			fields := unsafe.Slice((*uint32)(unsafe.Add(p.C(), int(layoutServerPlayer.Off(2152)) /* Field2152 */)), 3)
			if !noxflags.HasGame(noxflags.GameHost) && gameFrame() > fields[2] {
				fields[0], fields[1], fields[2] = uint32(word(3)), uint32(word(5)), gameFrame()
			}
			if fields[0] == fields[1]-1 {
				Nox_xxx_printCentered_445490(clientGameProgressText("SH_NearVictory", p.Name()))
			}
		}
	case 96:
		sub_462040(int32(code))
	case 97:
		sub_4624D0(int32(code))
	case 101:
		if dr := GetClient().Cli().Objs.ByNetCode(word(1)); dr != nil {
			old := *effectWord(dr, int(layoutClientDrawable.Off(280)) /* Flags70Val */)
			status := dword(3)
			*effectWord(dr, int(layoutClientDrawable.Off(280)) /* Flags70Val */) = status
			if dr.Class()&0x20000 != 0 {
				if old&0x400 == 0 && status&0x400 != 0 {
					objectGeneratorCountdown(dr)
				}
				if status&0x800 != 0 {
					dr.ObjClass &^= 0x80000
					dr.ObjFlags &^= 0x20000000
				}
			}
		}
	case 103:
		if dr := GetClient().Cli().Objs.ByNetCode(word(1)); dr != nil {
			for i := 0; i < 4; i++ {
				*drawableModifier(dr, i) = nox_xxx_modifGetDescById_413330(int32(data[3+i]))
			}
			*effectWord(dr, int(layoutClientDrawable.Off(448)) /* Union._ */) = 0xffffffff
		}
	case 108:
		typ := word(3)
		binary.LittleEndian.PutUint16(data[3:], typ&0x7fff)
		summonAdd(uint32(word(1)), uint32(typ&0x7fff), typ&0x8000 != 0)
		c := GetClient()
		dr := c.Cli().Objs.ByNetCode(word(1))
		if dr == nil {
			dr = c.Nox_xxx_spriteCreate_48E970(int(typ&0x7fff), code, 0, 0)
		}
		if dr != nil {
			c.Cli().Objs.MinimapAdd(dr, 1)
		}
		combatAllyAdd(uint32(word(1)), 0, 0)
	case 109:
		original := word(1)
		binary.LittleEndian.PutUint16(data[1:], code)
		summonRemove(uint32(code), original&0x8000 != 0)
		combatAllyRemove(uint32(code))
		// The legacy dispatcher performs lookup after clearing the high bit in place.
		if dr := GetClient().Cli().Objs.ByNetCodeDynamic(int(code)); dr != nil {
			GetClient().Cli().Objs.RemoveHealthBar(dr, 1)
		}
	case 110:
		sub_467440(int(dword(1)))
	case 111:
		bookSpellReward(int(data[1]), int(data[2]), int(data[3]&127), int(data[3]>>7))
	}
	return size, true
}
