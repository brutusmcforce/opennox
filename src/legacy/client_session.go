package legacy

import (
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/console"
	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"
	"github.com/brutusmcforce/opennox/libs/strman"
)

func clientSessionDrawable(code uint16) *client.Drawable {
	if code&0x8000 != 0 {
		return GetClient().Cli().Objs.ByNetCodeStatic(int(code & 0x7fff))
	}
	return GetClient().Cli().Objs.ByNetCodeDynamic(int(code & 0x7fff))
}

func clientGameSession(ind int, op netmsg.Op, data []byte) int {
	size := 0
	switch int(op) {
	case 196:
		return clientSessionTeam(data)
	case 201:
		return clientSessionTrade(data)
	case 240:
		return clientSessionQuest(data)
	case 166, 167, 178, 179, 180, 217, 218, 224, 226:
		size = 4
	case 171, 215, 219, 234:
		size = 5
	case 174, 202, 205, 206, 207, 209, 214, 220, 225, 229, 230, 231, 232:
		size = 3
	case 175:
		size = 20
	case 176:
		size = 49
	case 177:
		size = 60
	case 181:
		size = 14
	case 189, 235, 237:
		size = 2
	case 195:
		size = 12
	case 197, 198, 203:
		size = 1
	case 204:
		if len(data) < 4 {
			return -1
		}
		size = 4 + int(data[3])
	case 210:
		size = 7
	case 211:
		size = 13
	case 213:
		if len(data) < 2 || data[1] < 1 || data[1] > 3 {
			return -1
		}
		size = 68
	case 216, 223:
		size = 6
	case 233:
		size = 9
	case 238:
		if len(data) < 2 {
			return -1
		}
		switch data[1] {
		case 6:
			size = 3
		case 7:
			size = 2
		default:
			return -1
		}
	default:
		return -1
	}
	if len(data) < size {
		return -1
	}
	word := func(off int) uint16 { return binary.LittleEndian.Uint16(data[off:]) }
	dword := func(off int) uint32 { return binary.LittleEndian.Uint32(data[off:]) }
	// These updates intentionally apply regardless of the connection gate.
	switch int(op) {
	case 174:
		return size
	case 175, 176, 177:
		clientSessionSettings(int(op), data)
		return size
	case 195:
		code := word(1)
		if word(3) != 0 || code != 0 {
			if dr := GetClient().Nox_xxx_spriteCreate_48E970(int(word(3)), code, int(word(5)), int(word(7))); dr != nil {
				dr.Field_72 = GetServer().S().Frame()
				dr.SetFrameMB(int(data[10]))
				dir := (data[9] >> 4) & 7
				if dir > 3 {
					dir++
				}
				dr.AnimDir = dir
				if dr.AnimInd != uint32(data[11]) {
					dr.AnimStart = GetServer().S().Frame()
					dr.AnimInd = uint32(data[11])
				}
			}
			if int(code) == ClientPlayerNetCode() && Sub_416120(8) {
				Nox_xxx_cliUpdateCameraPos_435600(int(word(5)), int(word(7)))
			}
		} else {
			Nox_xxx_cliUpdateCameraPos_435600(int(word(5)), int(word(7)))
			InputSetKeyTimeoutLegacy(8)
		}
		return size
	case 197:
		browserNotice(false)
		return size
	case 198:
		browserNotice(true)
		return size
	case 215:
		if p := GetServer().S().Players.ByID(int(word(1))); p != nil {
			*(*uint16)(unsafe.Pointer(&p.Field2148)) = word(3)
		}
		return size
	case 223:
		quickbarSetFlash(int(dword(1)), data[5])
		return size
	case 224:
		sub_467750(int32(word(1)&0x7fff), int8(data[3]))
		return size
	case 225:
		dword_5d4594_1062488 = uint32(word(1) & 0x7fff)
		return size
	case 226:
		code := word(1)
		dr := uiShopDrawable(uint32(code & 0x7fff))
		if dr == nil {
			dr = uiInventoryItem(uint32(code & 0x7fff))
		}
		if dr == nil {
			dr = clientSessionDrawable(code)
		}
		if dr != nil {
			(*client.DrawableUnionBook)(unsafe.Pointer(&dr.Union)).Index = uint32(data[3])
		}
		return size
	case 234:
		if noxflags.HasGame(0x2000) {
			GetClient().Sub_45A670(dword(1))
		}
		return size
	case 238:
		if data[1] == 7 {
			voteGUIReset()
		} else {
			*memmap.PtrUint32(0x5D4594, 1197304) = uint32(data[2])
			hidden := 1
			if data[2] == 1 {
				hidden = 0
			}
			interactionVoteHide(hidden)
		}
		return size
	}
	if !Nox_client_isConnected() {
		return size
	}
	switch int(op) {
	case 166, 167:
		kind := 1
		if data[0] == 167 {
			kind = 2
		}
		audioEventPlay(int32(word(2)&0x3ff), int32((word(2)>>9)&0x7e), int32(int8(data[1])), kind)
	case 171:
		reliableACK(ind, dword(1))
	case 178:
		if dr := clientSessionDrawable(word(1)); dr != nil {
			dr.Field_74_4 = data[3]
		}
	case 179:
		if dr := clientSessionDrawable(word(1)); dr != nil {
			dr.SetActive()
			particleLightIntensity(dr, float32(16*int(data[3])/10), true)
			dr.SetFrameMB(8 * int(data[3]) / 50)
			if dr.AnimFrameSlave == 8 {
				dr.AnimFrameSlave = 7
			}
		}
	case 180:
		if dr := clientSessionDrawable(word(1)); dr != nil {
			flags := (*uint32)(unsafe.Pointer(&dr.ObjClass))
			value := float32(0)
			if data[3] != 0 {
				*flags |= 0x80000
				value = 41.958
			} else {
				*flags &^= 0x80000
			}
			particleLightIntensity(dr, value, true)
			dr.SetFrameMB(int(data[3]))
			dr.Field_72 = GetServer().S().Frame()
		}
	case 181:
		if dr := GetClient().Nox_xxx_spriteCreate_48E970(int(word(3)), word(1)&0x7fff, int(word(5)), int(word(7))); dr != nil {
			*(*uint16)(unsafe.Pointer(&dr.Field_127)) = word(9)
			dr.Field_117 = math.Float32bits(float32(int8(data[12])) * 0.0625)
			dr.Field_118 = math.Float32bits(float32(int8(data[13])) * 0.0625)
			dr.Field_119 = math.Float32bits(float32(int8(data[11])) * 0.0625)
			dr.AnimStart = GetServer().S().Frame()
			dr.Field_81 = uint32(word(5))
			dr.Field_82 = uint32(word(7))
			dr.Field_115 = drawableUpdateIdentity(updateID_sprite)
			GetClient().Cli().Objs.List5Add(dr)
		}
	case 189:
		key := "invalidPass"
		if data[1] == 1 {
			key = "sysopAccess"
		}
		GetConsole().Print(console.ColorRed, clientGameProgressString(key))
	case 202:
		if word(1) == 57005 {
			chatBubbleClear()
		} else {
			chatBubbleRemove(uint32(word(1)))
		}
	case 203:
		interactionMessagesClear()
	case 204:
		clientSequenceEnqueue(data)
	case 205:
		quickbarAbilityReward(int(data[1]), int(data[2]&0x7f), uintptr(data[2]>>7))
	case 206:
		quickbarAbilityState(uint32(data[1]), uint32(data[2]))
	case 207:
		quickbarAbilityFlags(uint32(data[1]), uint32(data[2]))
	case 209:
		bookGuideReward(int(data[1]), int(data[2]&0x7f))
	case 210:
		code := word(1)
		dr := clientSessionDrawable(code)
		if data[5] == 1 {
			if dr == nil {
				dr = GetClient().Nox_xxx_spriteCreate_48E970(int(word(3)), code&0x7fff, 0, 0)
			}
			if dr != nil {
				GetClient().Cli().Objs.MinimapAdd(dr, data[6])
			}
		} else if dr != nil {
			GetClient().Cli().Objs.RemoveHealthBar(dr, data[6])
		}
	case 211:
		if GetServer().S().Frame()-dword(9) >= 30 {
			reliableClientSend(31, []byte{212}, nil, 1)
		} else if dword(1) != 0 {
			serverConfigTimerSet(1)
			serverConfigTimerReset(int32(dword(5)))
		} else {
			serverConfigTimerSet(0)
		}
	case 213:
		end := 2
		for end < 66 && data[end] != 0 {
			end++
		}
		name := string(data[2:end])
		p := Get_dword_8531A0_2576()
		switch data[1] {
		case 1:
			if p != nil {
				journalAdd(p, name, word(66))
			}
			journalMeasure()
		case 2:
			if p != nil {
				journalRemove(p, name)
			}
			journalMeasure()
		case 3:
			if p != nil {
				journalUpdate(p, name, word(66))
			}
		}
	case 214:
		briefingShow(int(data[1]), int(data[2]), 0)
	case 216:
		teamUICTFTooltip(data[2], data[1])
		scoreboardSetFlag(data[1], data[3], word(4))
	case 217:
		clientSessionBall(data[1])
		scoreboardSetBall(data[1], word(2))
	case 218:
		code := uint32(word(1) & 0x7fff)
		if combatAllyLookup(code) != nil {
			combatAllyFlag(code, data[3])
		}
		GetServer().S().NPCs.Set328(int(code), int(data[3]))
	case 219:
		code := word(1)
		dr := clientSessionDrawable(code)
		if dr == nil {
			dr = GetClient().Nox_xxx_spriteCreate_48E970(int(word(3)), code&0x7fff, 0, 0)
		}
		if dr != nil {
			GetClient().Cli().Objs.MinimapAdd(dr, 1)
		}
		combatAllyAdd(uint32(code&0x7fff), 0, 0)
	case 220:
		combatAllyRemove(uint32(word(1)))
		if dr := clientSessionDrawable(word(1)); dr != nil {
			GetClient().Cli().Objs.RemoveHealthBar(dr, 1)
		}
	case 229:
		Sub_43D9B0(uint32(data[1]), uint32(data[2]))
	case 230:
		audioEventMusicSave()
	case 231:
		audioEventMusicRestore()
	case 232:
		clientSessionDeath(word(1))
	case 233:
		if noxflags.HasGame(0x2000) && int(word(1)) == ClientPlayerNetCode() {
			GetClient().Sub_45A670(dword(3))
		}
		if data[8]&1 != 0 {
			if p := GetServer().S().Players.ByID(int(word(1))); p != nil {
				playerStateRespawn(p, data[7])
			}
		}
	case 235:
		quickbarResetAbility(data[1])
	case 237:
		summonSetCommand(uint32(data[1]))
	}
	return size
}

func clientSessionDeath(code uint16) {
	p := GetServer().S().Players.ByID(int(code))
	if p == nil {
		return
	}
	if p == Get_dword_8531A0_2576() {
		uiAmountCancel()
		Sub_478000()
	}
	if dr := GetClient().Cli().Objs.ByNetCodeDynamic(int(code)); dr != nil {
		GetClient().Cli().Objs.RemoveHealthBar(dr, 3)
	}
	if noxflags.HasGame(4096) {
		return
	}
	word := func(off int) *uint32 { return (*uint32)(unsafe.Add(p.C(), off)) }
	*word(4) = 0
	// Preserve the original four-byte offset into the weapon records.
	for i := 0; i < 27; i++ {
		base := 2328 + 24*i
		for _, off := range []int{0, 4, 8, 12, 20} {
			*word(base + off) = 0
		}
	}
	for i := 0; i < 26; i++ {
		base := 2972 + 24*i
		mask := *word(base)
		if mask&0xc0d == 0 {
			*word(0) &^= mask
			for off := 0; off < 20; off += 4 {
				*word(base + off) = 0
			}
		}
	}
}

func clientSessionBall(state byte) {
	*memmap.PtrUint8(0x5D4594, 1045644) = state
	var imageName, key string
	switch state {
	case 0:
		imageName, key = "BallAtHome", "BallHomeTT"
	case 1:
		imageName, key = "BallAway", "BallAwayTT"
	case 2:
		imageName, key = "BallRed", "BallRedTT"
	case 4:
		imageName, key = "BallBlue", "BallBlueTT"
	default:
		return
	}
	if w := teamUIWindow(1045636); w != nil {
		w.DrawDataVal.BgImageHnd = noxrender.ImageHandle(Nox_xxx_gLoadImg(imageName).C())
		sm := GetServer().S().Strings()
		w.DrawData().SetTooltip(GetClient().Cli().Strings(), sm.GetStringInFile(strman.ID(key), "guifb.c"))
	}
}
