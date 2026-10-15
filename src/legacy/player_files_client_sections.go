package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func playerFileCopyString(dst unsafe.Pointer, s string) {
	copy(unsafe.Slice((*byte)(dst), len(s)+1), s+"\x00")
}
func playerFileWideLen(p unsafe.Pointer) int {
	n := 0
	for *(*uint16)(unsafe.Add(p, 2*n)) != 0 {
		n++
	}
	return n
}
func playerFileGUI() int {
	r := playerFileStream()
	version := int16(r.short(3))
	if version > 3 {
		return 0
	}
	if ret := int(quickbarSave()); ret == 0 {
		return ret
	}
	if version >= 2 {
		row := r.byte(byte(quickbarMain().Selected))
		if r.read() {
			quickbarSelectRow(int(row))
		}
		trap := r.byte(quickbarRecords[0].Selected)
		if r.read() {
			quickbarTrapSelect(int(trap))
		}
	}
	if version >= 3 {
		p := GetServer().S().Players.ByID(int(nox_player_netCode_85319C))
		order := byte(4)
		if p != nil {
			order = *(*byte)(unsafe.Pointer(&p.SummonOrderAll))
		}
		order = r.byte(order)
		if r.read() && noxflags.HasGame(2048) {
			unitLocalOrder(int(uint8(p.PlayerInd)), int32(order))
		}
	}
	return 1
}

// playerFileSaveName transfers the map name of a save (Player.SaveName). The historical format transfers twice
// the length of the name, while the player only has room for 32 bytes (up to Field4792, the quest state), so it
// goes through a buffer: long map names didn't fit otherwise. The file format doesn't change.
func playerFileSaveName(r objectXferStream, p *server.Player, name string) {
	const saveNameRoom = unsafe.Offsetof(p.Field4792) - unsafe.Offsetof(p.SaveNameBuf)
	saveName := unsafe.Slice((*byte)(unsafe.Pointer(&p.SaveNameBuf)), saveNameRoom)
	var buf []byte
	if !r.read() {
		buf = make([]byte, 2*len(name))
		copy(buf, name)
	}
	n := int(r.short(uint16(len(buf) / 2)))
	if r.read() {
		buf = make([]byte, 2*n)
	}
	if n > 0 {
		r.raw(unsafe.Pointer(&buf[0]), 2*n)
	}
	// writing also sets it, as the original
	n = min(n, int(saveNameRoom)-1)
	clear(saveName)
	copy(saveName[:n], buf)
}

func playerFileGame(u *server.Object) int {
	p := u.UpdateDataPlayer().Player
	if noxflags.HasGame(8192) {
		return 1
	}
	r := playerFileStream()
	version := int16(r.short(5))
	if version > 5 {
		return 0
	}
	if version >= 5 {
		if r.read() {
			r.word(0)
		} else {
			id := r.word(uint32(GetServer().S().Objs.LastObjectScriptID()))
			GetServer().S().Objs.SetLastObjectScriptID(server.ObjectScriptID(id))
		}
	}
	playerFileSaveName(r, p, alloc.GoString(sessionMapName()))
	if version >= 2 {
		ret := 0
		if r.read() {
			ret = questProgressRead()
		} else {
			ret = questProgressWrite()
		}
		if ret == 0 {
			return ret
		}
	}
	if version >= 3 {
		if ret := Sub_5000B0(u); ret == 0 {
			return ret
		}
	}
	if version >= 4 {
		audioEventSetByte(int8(r.byte(audioEventByte())))
	} else {
		audioEventSetByte(0)
	}
	return 1
}
func playerFileMetadata() int {
	r := playerFileStream()
	version := int16(r.short(12))
	if version > 12 {
		return 0
	}
	base := memmap.PtrOff(0x85B3FC, 10980)
	field := func(off int) unsafe.Pointer { return unsafe.Add(base, off) }
	flags := (*uint32)(base)
	if noxflags.HasGame(8192) {
		*flags &^= 1
		if noxflags.HasGame(4096) || questRuntimeWord(1556160) != 0 || questRuntimeWord(1556164) != 0 {
			*flags |= 4
		} else {
			*flags |= 2
		}
	} else {
		*flags = *flags&^6 | 1
	}
	r.raw(base, 4)
	if r.read() {
		n := int(int16(r.short(0)))
		r.raw(field(4), n)
		*(*byte)(field(4 + n)) = 0
		n = int(r.byte(0))
		r.raw(field(1028), n)
		*(*byte)(field(1028 + n)) = 0
	} else {
		n := r.short(uint16(len(alloc.GoString((*byte)(field(4))))))
		r.raw(field(4), int(int16(n)))
		n8 := r.byte(byte(len(alloc.GoString((*byte)(field(1028))))))
		r.raw(field(1028), int(n8))
	}
	noxGetLocalTime((*noxSystemTime)(field(1188)))
	for off := 1188; off < 1204; off += 2 {
		r.raw(field(off), 2)
	}
	for _, off := range []int{1207, 1204, 1210, 1213, 1216} {
		r.raw(field(off), 3)
	}
	for off := 1219; off < 1224; off++ {
		r.raw(field(off), 1)
	}
	if r.read() {
		n := int(r.byte(0))
		r.raw(field(1224), 2*n)
		*(*uint16)(field(1224 + 2*n)) = 0
	} else {
		n := r.byte(byte(playerFileWideLen(field(1224))))
		r.raw(field(1224), 2*int(n))
	}
	r.raw(field(1274), 1)
	*(*byte)(field(1275)) = 0
	r.raw(field(1275), 1)
	if !r.read() {
		*(*byte)(field(1276)) = byte(sub_467590())
	}
	r.raw(field(1276), 1)
	if version >= 11 {
		// SaveGameInfo.MapNameBuf has 32 bytes
		name := alloc.GoString(sessionMapName())
		playerFileCopyString(field(1156), name[:min(len(name), 31)])
		n := int(r.byte(byte(len(alloc.GoString((*byte)(field(1156)))))))
		if n > 31 {
			buf := make([]byte, n) // a damaged file: keep the start of the name
			r.raw(unsafe.Pointer(&buf[0]), n)
			copy(unsafe.Slice((*byte)(field(1156)), 31), buf)
			n = 31
		} else {
			r.raw(field(1156), n)
		}
		*(*byte)(field(1156 + n)) = 0
	}
	if version < 12 {
		*(*byte)(field(1277)) = 0
	} else {
		r.raw(field(1277), 1)
	}
	return 1
}
func playerFileSaveRequest() int {
	code := uint32(nox_player_netCode_85319C)
	data := [3]byte{0xc1, byte(code), byte(code >> 8)}
	Nox_xxx_netClientSend2_4E53C0(31, unsafe.Pointer(&data[0]), 3, 0, 1)
	return 1
}
