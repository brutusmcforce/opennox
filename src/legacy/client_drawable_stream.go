package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func drawableStreamKeys(p unsafe.Pointer, offset *int, player int, permanent bool) (uint16, uint16) {
	table := (*[255]server.PlayerNetData)(memmap.PtrOff(0x5D4594, 1198020))
	slot := *(*byte)(unsafe.Add(p, *offset))
	*offset++
	if slot != 255 {
		r := table[slot]
		return r.Field0, r.Field2
	}
	code := *(*uint16)(unsafe.Add(p, *offset))
	typ := *(*uint16)(unsafe.Add(p, *offset+2))
	*offset += 4
	frame := GetServer().S().Frame()
	slot = selectNetworkAlias(table, int32(code), int32(typ), frame)
	if slot != 255 {
		expiry := frame + 60
		if permanent {
			expiry = 0xffffffff
		}
		table[slot] = server.PlayerNetData{Field0: code, Field2: typ, Frame4: expiry}
		var packet [10]byte
		packet[0], packet[1] = 0xa5, slot
		binary.LittleEndian.PutUint16(packet[2:], code)
		binary.LittleEndian.PutUint16(packet[4:], typ)
		binary.LittleEndian.PutUint32(packet[6:], expiry)
		GetServer().S().NetList.AddToMsgListCli(ntype.PlayerInd(player), netlist.Kind0, packet[:])
	}
	return code, typ
}
func drawableStreamDirection(dr *client.Drawable, status byte) {
	dr.AnimDir = (status >> 4) & 7
	if dr.AnimDir > 3 {
		dr.AnimDir++
	}
}
func drawableStreamAnimation(dr *client.Drawable, anim uint32) {
	if dr.AnimInd != anim {
		dr.AnimInd = anim
		dr.AnimStart = GetServer().S().Frame()
	}
}
func drawableStreamFirst(p unsafe.Pointer, player int, pos *[2]int32) int32 {
	n := 0
	code, typ := drawableStreamKeys(p, &n, player, true)
	x, y := *(*uint16)(unsafe.Add(p, n)), *(*uint16)(unsafe.Add(p, n+2))
	n += 4
	status := *(*byte)(unsafe.Add(p, n))
	n++
	frame := byte(0)
	if status&0x80 != 0 {
		frame = *(*byte)(unsafe.Add(p, n))
		n++
	}
	anim := *(*byte)(unsafe.Add(p, n))
	n++
	if code != 0 || typ != 0 {
		if dr := GetClient().Nox_xxx_spriteCreate_48E970(int(typ), code, int(x), int(y)); dr != nil {
			dr.Field_72 = GetServer().S().Frame()
			drawableStreamDirection(dr, status)
			drawableStreamAnimation(dr, uint32(anim))
			dr.SetFrameMB(int(frame))
		}
	}
	pos[0], pos[1] = int32(x), int32(y)
	Nox_xxx_cliUpdateCameraPos_435600(int(x), int(y))
	return int32(n)
}
func drawableStreamNext(p unsafe.Pointer, player int, pos *[2]int32) int32 {
	n := 0
	absolute := *(*byte)(p) == 0
	if absolute {
		if *(*byte)(unsafe.Add(p, 1)) == 0 && *(*byte)(unsafe.Add(p, 2)) == 0 {
			return -3
		}
		n++
	}
	code, typ := drawableStreamKeys(p, &n, player, false)
	if absolute {
		pos[0] = int32(*(*int16)(unsafe.Add(p, n)))
		pos[1] = int32(*(*int16)(unsafe.Add(p, n+2)))
		n += 4
	} else {
		pos[0] += int32(*(*int8)(unsafe.Add(p, n)))
		pos[1] += int32(*(*int8)(unsafe.Add(p, n+1)))
		n += 2
	}
	if pos[0] < 0 || pos[0] > 6000 || pos[1] < 0 || pos[1] > 6000 {
		return -int32(n)
	}
	dr := GetClient().Nox_xxx_spriteCreate_48E970(int(typ), code, int(pos[0]), int(pos[1]))
	if dr == nil {
		return -int32(n)
	}
	dr.Field_72 = GetServer().S().Frame()
	if dr.ObjClass&0x200000 != 0 {
		status := *(*byte)(unsafe.Add(p, n))
		n++
		drawableStreamDirection(dr, status)
		if status&0x80 != 0 {
			dr.SetFrameMB(int(*(*byte)(unsafe.Add(p, n))))
			n++
		}
		anim := status & 15
		if dr.ObjClass&4 != 0 {
			anim = *(*byte)(unsafe.Add(p, n))
			n++
		}
		drawableStreamAnimation(dr, uint32(anim))
	} else {
		GetClient().Cli().Objs.AddIndex2D(dr)
	}
	pos[0], pos[1] = int32(dr.PosVec.X), int32(dr.PosVec.Y)
	return int32(n)
}

func nox_xxx_netCliProcUpdateStream_494A60(p *uint8, player int32, pos *uint32) int32 {
	return int32(drawableStreamFirst(unsafe.Pointer(p), int(player), (*[2]int32)(unsafe.Pointer(pos))))
}
