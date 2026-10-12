package legacy

import (
	"encoding/binary"
	"image"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"
	"github.com/brutusmcforce/opennox/client"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
)

var clientFadeObjects uint32 = 1

func clientGameLocalDrawable() *client.Drawable {
	return LocalPlayerDrawable
}
func clientGameKeepDrawable(dr *client.Drawable) bool {
	return dr == clientGameLocalDrawable() || dr.DrawFuncPtr == drawableDrawKey(drawKey_nox_thing_animate_draw) && dr.DrawData != nil && *(*uint32)(unsafe.Add(dr.DrawData, 12)) == 1
}
func clientGameRemoveDrawable(dr *client.Drawable, code uint16) {
	if code&0x8000 != 0 {
		GetClient().Cli().Nox_xxx_cliDestroyObj_45A9A0(dr)
	} else {
		GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
	}
}
func clientGameRestoreLight(dr *client.Drawable) {
	typ := GetClient().Cli().Things.TypeByInd(int(dr.TypeIDVal))
	particleLightIntensity(unsafe.Add(dr.C(), 136), typ.LightIntensity, true)
}

// clientGameState handles complete case groups; unselected messages retain the C fallback.
func clientGameState(op netmsg.Op, data []byte) (int, bool) {
	size := 0
	switch int(op) {
	case 47:
		size = 9
	case 48, 81, 82:
		size = 11
	case 50, 51, 52, 53, 55, 56, 59, 60, 62, 67:
		size = 3
	case 54:
		size = 1
	case 57, 65, 94, 95:
		size = 4
	case 61, 92:
		size = 6
	case 66, 69, 73, 74, 100, 102:
		size = 5
	case 68, 79, 80, 83, 84, 90, 93, 107, 221, 222:
		size = 7
	case 71, 91:
		size = 2
	case 104:
		size = 8
	case 105:
		size = 21
	case 169:
		return clientGameNotice(data), true
	default:
		return 0, false
	}
	word := func(off int) uint16 { return binary.LittleEndian.Uint16(data[off:]) }
	dword := func(off int) uint32 { return binary.LittleEndian.Uint32(data[off:]) }
	if int(op) == 66 {
		combatHealthAdd(uint32(word(1)), int16(word(3)))
		return size, true
	}
	if !Nox_client_isConnected() {
		return size, true
	}
	drawable := func() *client.Drawable { return GetClient().Cli().Objs.ByNetCode(word(1)) }
	local := func() bool { return int(word(1)&0x7fff) == ClientPlayerNetCode() }
	switch int(op) {
	case 47, 48:
		code, typ, x, y := word(1), word(3), int(word(5)), int(word(7))
		if int(op) == 48 && code == 0 && typ == 0 {
			Nox_xxx_cliUpdateCameraPos_435600(x, y)
			InputSetKeyTimeoutLegacy(9)
			break
		}
		if dr := GetClient().Nox_xxx_spriteCreate_48E970(int(typ), code, x, y); dr != nil {
			dr.Field_72 = GetServer().S().Frame()
			if int(op) == 48 {
				dr.SetFrameMB(int(data[10]))
				drawableStreamDirection(dr, data[9])
				data[9] &= 15
				drawableStreamAnimation(dr, uint32(data[9]))
			}
		}
		if int(op) == 48 && int(code) == ClientPlayerNetCode() && Sub_416120(9) {
			Nox_xxx_cliUpdateCameraPos_435600(x, y)
		}
	case 50, 51:
		if dr := drawable(); dr != nil {
			if int(op) == 51 {
				dr.Field_120, dr.Field_121, dr.Field_122 = 1, 1, 1
				if clientFadeObjects != 0 {
					if dr != clientGameLocalDrawable() {
						GetClient().Cli().Objs.TransparentDecay(dr, int(GetServer().S().TickRate()))
					}
					break
				}
			}
			if !clientGameKeepDrawable(dr) {
				clientGameRemoveDrawable(dr, word(1))
			}
		}
	case 52:
		combatFriendAdd(uint32(word(1) & 0x7fff))
	case 53:
		combatFriendRemove(uint32(word(1) & 0x7fff))
	case 54:
		combatFriendClear()
	case 55, 56:
		if dr := drawable(); dr != nil {
			if int(op) == 55 {
				dr.ObjFlags |= 0x1000000
			} else {
				if dr.ObjClass&0x40000 != 0 {
					dr.DrawFuncPtr = nil
				}
				dr.ObjFlags &^= 0x1000000
			}
		}
	case 57, 107:
		if dr := drawable(); dr != nil {
			if int(op) == 57 {
				dr.SetFrameMB(int(data[3]))
			} else {
				dr.SetFrameMB(int(dword(3)))
			}
		}
	case 59, 60:
		if !noxflags.HasGame(noxflags.GameHost) {
			if w := worldSecretFind(int16(word(1))); w != nil && w.Flags4&4 != 0 {
				if int(op) == 59 {
					*(*byte)(unsafe.Add(w.Data, 22)) = 23
					*(*byte)(unsafe.Add(w.Data, 21)) = 3
				} else if w.Data != nil {
					*(*byte)(unsafe.Add(w.Data, 22)) = 0
					*(*byte)(unsafe.Add(w.Data, 21)) = 1
				}
			}
		}
	case 61:
		walls := &GetServer().S().Walls
		pos := image.Pt(int(data[4]), int(data[5]))
		w := walls.GetWallAtGrid(pos)
		if w == nil {
			w = walls.CreateAtGrid(pos)
		}
		if w != nil {
			w.Tile1, w.Dir0, w.Field2 = data[1], data[2], data[3]
		}
	case 62:
		walls := &GetServer().S().Walls
		if w := walls.GetWallAtGrid(image.Pt(int(data[1]), int(data[2]))); w != nil {
			walls.DeleteAtGrid(image.Pt(int(w.X5), int(w.Y6)))
		}
	case 65:
		if combatAllyLookup(uint32(word(1))) != nil {
			combatAllyFirst(uint32(word(1)), 2*uint16(data[3]))
		}
	case 67:
		sub_470CB0(int(int16(word(1))))
	case 68:
		uiInventoryItemHealth(int(word(1)&0x7fff), int16(word(3)), int16(word(5)))
	case 69:
		if local() {
			nox_xxx_cliSetMana_470D10(int(word(3)))
		}
	case 71:
		sub_470D20(int(data[1]), int(memmap.Int32(0x587000, 157092)))
	case 73:
		sub_467450(int(dword(1)))
	case 74:
		sub_467490(int(dword(1)))
	case 79, 80, 81, 82:
		wire := word(1)
		code := uint32(wire & 0x7fff)
		binary.LittleEndian.PutUint16(data[1:], uint16(code))
		implicit := [4]byte{255, 255, 255, 255}
		mods := &implicit
		if int(op) == 81 || int(op) == 82 {
			mods = (*[4]byte)(data[7:11])
		}
		if wire&0x8000 != 0 {
			playerStateEquip(data[0], code, dword(3), mods)
		} else {
			presentationEquip(data[0], code, dword(3), mods)
		}
	case 83, 84:
		playerStateUnequip(data[0], uint32(word(1)), dword(3))
	case 90:
		if dr := drawable(); dr != nil {
			before := dr.Buffs&0x8000 != 0
			dr.Buffs = dword(3)
			isLocal := dr == clientGameLocalDrawable()
			if isLocal {
				sub_467410(int(dr.Buffs))
			}
			if before && dr.Buffs&0x8000 == 0 && !(isLocal && memmap.Uint8(0x5D4594, 1062536)&8 != 0) {
				clientGameRestoreLight(dr)
			}
		}
	case 91:
		p := memmap.PtrUint8(0x5D4594, 1062536)
		old := *p
		*p = data[1]
		if dr := clientGameLocalDrawable(); old&8 != 0 && *p&8 == 0 && dr != nil && dr.Buffs&0x8000 == 0 {
			clientGameRestoreLight(dr)
		}
	case 92:
		if dr := drawable(); dr != nil {
			particleLightColor(unsafe.Add(dr.C(), 136), int(data[3]), int(data[4]), int(data[5]))
		}
	case 93:
		if dr := drawable(); dr != nil {
			particleLightIntensity(unsafe.Add(dr.C(), 136), math.Float32frombits(dword(3)), true)
		}
	case 94, 95:
		if dr := drawable(); dr != nil {
			dr.ZVal = uint16(data[3])
			if int(op) == 95 {
				dr.ZVal = -dr.ZVal
			}
		}
	case 100:
		sub_467930(int(word(1)&0x7fff), int(data[3]), int(data[4]))
	case 102:
		if dr := clientGameLocalDrawable(); dr != nil {
			*(*uint32)(unsafe.Add(dr.C(), 120)) = dword(1)
		}
	case 104:
		if local() {
			sub_467470(int(data[7]), math.Float32frombits(dword(3)))
		}
	case 105:
		wire := word(1)
		code := int(wire & 0x7fff)
		binary.LittleEndian.PutUint16(data[1:], uint16(code))
		npcs := &GetServer().S().NPCs
		p := npcs.ByID(code)
		if p != nil {
			npcs.Set(p, code)
		} else {
			p = npcs.New(code)
		}
		if p != nil {
			for i := range p.Color8 {
				off := 3 + 3*i
				p.Color8[i] = particleRGB(int(data[off]), int(data[off+1]), int(data[off+2]))
			}
			p.Field1312 = uint32(wire >> 15)
		}
	case 221:
		if local() {
			uiMeterSetTotal(0, 2247, int(word(3)), int(word(5)))
		} else if combatAllyLookup(uint32(word(1))) != nil {
			combatAllyPair(uint32(word(1)), word(3), word(5))
		}
	case 222:
		if local() {
			uiMeterSetTotal(1, 2243, int(word(3)), int(word(5)))
		}
	}
	return size, true
}
