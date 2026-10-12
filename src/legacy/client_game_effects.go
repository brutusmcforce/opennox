package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"math"
	"unsafe"
)

var clientGameBlueSpark, clientGameVioletSpark uint32

func clientGameEffectType(offset uintptr, name string) uint32 {
	p := effectMapped(offset)
	if *p == 0 {
		*p = effectType(name)
	}
	return *p
}
func clientGameEffectLink(typ uint32, pos image.Point) {
	if dr := effectSpawn(int(typ), pos); dr != nil {
		effectLink(dr)
	}
}

func clientGameEffects(player int, op netmsg.Op, data []byte) (int, bool) {
	size := 0
	switch int(op) {
	case 125, 140, 141, 142, 143, 144, 145, 149:
		size = 9
	case 126:
		size = 12
	case 127, 155, 156, 157:
		size = 3
	case 128:
		size = 4
	case 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 150, 154, 160, 163:
		size = 5
	case 147, 159, 161:
		size = 6
	case 152, 162:
		size = 11
	case 158:
		if data[1] < 1 || data[1] > 14 {
			return -1, true
		}
		size = 7
	case 164:
		var coords [2]int32
		offset := 1 + int(drawableStreamFirst(unsafe.Pointer(&data[1]), player, &coords))
		if offset >= len(data) {
			return 0, true
		}
		for {
			n := int(drawableStreamNext(unsafe.Pointer(&data[offset]), player, &coords))
			if n <= 0 {
				return offset - n, true
			}
			offset += n
			if offset >= len(data) {
				return 0, true
			}
		}
	default:
		return 0, false
	}
	if int(op) == 147 && *effectMapped(1200852) == 0 {
		*effectMapped(1200852) = effectType("Spark")
		*effectMapped(1200856) = effectType("MediumFireBoom")
		*effectMapped(1197380) = effectType("FireBoom")
	}
	if int(op) == 152 {
		clientGameEffectType(1200844, "GreenZap")
	}
	if !Nox_client_isConnected() {
		return size, true
	}
	word := func(off int) uint16 { return binary.LittleEndian.Uint16(data[off:]) }
	signedPos := func() image.Point { return image.Pt(int(int16(word(1))), int(int16(word(3)))) }
	switch int(op) {
	case 125, 140, 141, 142, 143, 144, 145:
		effectDispatchRay((*[9]byte)(data))
		if int(op) != 143 && int(op) != 144 && int(op) != 145 {
			if data[0] == 140 || data[0] == 142 {
				effectLightningParticles(int(clientGameBlueSpark), image.Pt(int(word(1)), int(word(3))), image.Pt(int(word(5)), int(word(7))))
			} else if data[0] == 125 {
				effectCreateEnergySparks(int(word(5)), int(word(7)), 10, int(clientGameBlueSpark))
			}
		}
	case 126:
		drawableSummonStart(int16(word(5)), [2]uint16{word(1), word(3)}, word(7), data[9], int16(word(10)))
	case 127:
		drawableSummonStop(int16(word(1)))
	case 128:
		drawableShield(uint32(word(1)), int(data[3]))
	case 129, 130, 131, 132:
		typ, count, speed, ttl := clientGameBlueSpark, 50, 1000, 30
		if int(op) != 129 {
			count, speed, ttl = 25, 500, 25
			switch int(op) {
			case 130:
				typ = *effectMapped(1200780)
			case 131:
				typ = *effectMapped(1200784)
			case 132:
				typ = clientGameVioletSpark
			}
		}
		pos := signedPos()
		effectCreatePointSparks(int(typ), count, speed, ttl, pos.X, pos.Y)
	case 133, 134, 135, 136, 137, 139:
		offsets := map[int]uintptr{133: 1200872, 134: 1200876, 135: 1200880, 136: 1200884, 137: 1200888, 139: 1200892}
		names := map[int]string{133: "FireBoom", 134: "MediumFireBoom", 135: "CounterspellBoom", 136: "ThinFireBoom", 137: "TeleportPoof", 139: "DamagePoof"}
		pos := signedPos()
		if int(op) == 137 || int(op) == 139 {
			pos.Y += 2
		}
		clientGameEffectLink(clientGameEffectType(offsets[int(op)], names[int(op)]), pos)
	case 138:
		if *effectMapped(1200900) == 0 {
			*effectMapped(1200900) = effectType("Smoke")
			*effectMapped(1200896) = effectType("Puff")
		}
		pos := signedPos()
		if dr := effectSpawn(int(*effectMapped(1200900)), pos); dr != nil {
			dr.ZVal = 20
			effectLink(dr)
		}
		for i := 0; i < 6; i++ {
			y := pos.Y + effectRand(-15, 15)
			x := pos.X + effectRand(-15, 15)
			if dr := effectSpawn(int(*effectMapped(1200896)), image.Pt(x, y)); dr != nil {
				dr.ZVal = uint16(effectRand(5, 25))
				effectLink(dr)
			}
		}
	case 147:
		pos := signedPos()
		effectSparkBurst(pos, int(*effectMapped(1200852)), data[5])
		typ := *effectMapped(1197380)
		if data[5] <= 170 {
			typ = *effectMapped(1200856)
		}
		clientGameEffectLink(typ, pos)
	case 149:
		objectRenderBeamAppend(unsafe.Pointer(&data[0]))
		if effectRand(0, 100) < 25 {
			if local := clientGameLocalDrawable(); local != nil {
				dx, dy := uint32(word(5))-uint32(local.PosVec.X), uint32(word(7))-uint32(local.PosVec.Y)
				distance := int(math.Sqrt(float64(dx*dx + dy*dy)))
				if distance < 600 {
					audioEventPlay(297, int32(100*(600-distance)/600), 0, 0)
				}
			}
		}
		if !noxflags.HasGame(noxflags.GamePause) {
			dx, dy := int(word(5))-int(word(1)), int(word(7))-int(word(3))
			squared := int32(dx*dx + dy*dy)
			distance := math.MinInt32
			if squared >= 0 {
				distance = int(math.Sqrt(float64(squared)))
			}
			if distance == 0 {
				distance = 1
			}
			pos := image.Pt(int(word(5))-4*dx/distance, int(word(7))-4*dy/distance)
			if dr := effectSpawn(int(clientGameVioletSpark), pos); dr != nil {
				effectInitSpark(dr, pos, 1500, 5, 20)
				dr.ZVal = 22
				dr.VelZ = int8(effectRand(-4, 4))
				effectLink(dr)
			}
		}
	case 150:
		typ := clientGameEffectType(1200860, "BlueSpark")
		pos := signedPos()
		for i := 0; i < 5; i++ {
			if dr := effectSpawn(int(typ), pos); dr != nil {
				*effectWord(dr, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint32(dr.PosVec.X) << 12
				*effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(dr.PosVec.Y) << 12
				dr.Field_74_4 = byte(effectRand(0, 255))
				*effectWord(dr, int(layoutClientDrawable.Off(440)) /* Union._ */) = uint32(effectRand(1333, 4000))
				*effectWord(dr, int(layoutClientDrawable.Off(448)) /* Union._ */) = gameFrame() + uint32(effectRand(5, 20))
				*effectWord(dr, int(layoutClientDrawable.Off(444)) /* Union._ */) = gameFrame()
				dr.ZVal = 20
				dr.VelZ = int8(effectRand(-5, 5))
				effectLink(dr)
			}
		}
	case 152:
		x, y := int(word(1)), int(word(3))
		pos := image.Pt(x+(int(word(5))-x)/2, y+(int(word(7))-y)/2)
		if dr := effectSpawn(int(*effectMapped(1200844)), pos); dr != nil {
			*effectByte(dr, int(layoutClientDrawable.Off(432)) /* Union._ */) = 0
			*effectWord(dr, int(layoutClientDrawable.Off(433)) /* Union._ */) = uint32(word(9))
			*effectWord(dr, int(layoutClientDrawable.Off(437)) /* Union._ */) = binary.LittleEndian.Uint32(data[1:])
			*effectWord(dr, int(layoutClientDrawable.Off(441)) /* Union._ */) = binary.LittleEndian.Uint32(data[5:])
		}
	case 154:
		*effectMapped(1096520) = 1
	case 155, 156, 157:
		clientMapProgress(data)
	case 158:
		if data[1] <= 7 {
			presentationRayAdd((*[7]byte)(data))
		} else {
			presentationRayRemove((*[7]byte)(data))
		}
	case 159:
		if dr := GetClient().Cli().Objs.ByNetCode(word(1)); dr != nil {
			*effectWord(dr, int(layoutClientDrawable.Off(432)) /* Union._ */) = gameFrame()
			*effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */) = math.Float32bits(float32(data[3]))
			*effectWord(dr, int(layoutClientDrawable.Off(440)) /* Union._ */) = math.Float32bits(float32(int8(data[4])))
			*effectWord(dr, int(layoutClientDrawable.Off(444)) /* Union._ */) = math.Float32bits(float32(data[5]))
		}
	case 160:
		pos := [2]int16{int16(word(1)), int16(word(3))}
		presentationTurnUndead(&pos)
	case 161:
		if *effectMapped(1200864) == 0 {
			*effectMapped(1200864) = effectType("ArrowTrap1Smoke")
			*effectMapped(1200868) = effectType("ArrowTrap2Smoke")
		}
		pos := signedPos()
		typ := *effectMapped(1200868)
		if data[5] == 1 {
			pos.X += 15
			typ = *effectMapped(1200864)
		} else {
			pos.X -= 3
		}
		clientGameEffectLink(typ, pos)
	case 162:
		typ := clientGameEffectType(1200848, "HealOrb")
		coords := [4]uint16{word(1), word(3), word(5), word(7)}
		count := min(int(word(9)>>2), 7) + 1
		for i := 0; i < count; i++ {
			speed := byte(effectRand(6, 12))
			y := effectRand(-20, 20)
			x := effectRand(-20, 20)
			effectCreateOrb(int(typ), &coords, x, y, speed, 0)
		}
	case 163:
		radius := uint32(floatToInt32(float32(GetServer().S().Balance.Float("ManaBombOutRadius"))))
		pos := signedPos()
		for i := 0; i < 150; i++ {
			r := int(radius>>2) + effectRand(0, int(radius))
			if r > int(radius) {
				r = int(radius)
			}
			angle := effectRand(0, 255)
			where := image.Pt(pos.X+r*int(memmap.Int32(0x587000, 192088+8*uintptr(angle)))/16, pos.Y+r*int(memmap.Int32(0x587000, 192092+8*uintptr(angle)))/16)
			if dr := effectSpawn(int(*effectMapped(1200784)), where); dr != nil {
				*effectWord(dr, int(layoutClientDrawable.Off(432)) /* Union._ */) = uint32(dr.PosVec.X) << 12
				*effectWord(dr, int(layoutClientDrawable.Off(436)) /* Union._ */) = uint32(dr.PosVec.Y) << 12
				dr.Field_74_4 = 0
				*effectWord(dr, int(layoutClientDrawable.Off(440)) /* Union._ */) = 0
				*effectWord(dr, int(layoutClientDrawable.Off(448)) /* Union._ */) = gameFrame() + uint32(effectRand(30, 40))
				*effectWord(dr, int(layoutClientDrawable.Off(444)) /* Union._ */) = gameFrame()
				dr.ZVal = 0
				dr.VelZ = int8(effectRand(4, 10))
				effectLink(dr)
			}
		}
	}
	return size, true
}
