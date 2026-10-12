package legacy

import (
	"encoding/binary"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"io"
	"unsafe"
)

func playerFileServerSection(id uint32, u *server.Object, info unsafe.Pointer) (int, bool) {
	switch id {
	case 2:
		return playerFileAttributes(u, info), true
	case 3:
		return playerFileStatus(u), true
	case 4:
		return playerFileInventory(u), true
	case 8:
		return playerFileGuides(u), true
	case 5:
		return playerFileSpells(u), true
	case 6:
		return playerFileEnchantment(u), true
	case 9:
		return playerFileJournal(u), true
	case 10:
		return playerFileGame(u), true
	case 11:
		return cryptfile.Global().ReadWriteAlign(), true
	}
	return 0, false
}
func playerFileServerLoad(path string, ind int) int {
	core := GetServer().S()
	p := core.Players.ByInd(ntype.PlayerInd(ind))
	if p == nil || p.PlayerUnit == nil {
		return 0
	}
	u := p.PlayerUnit
	if cryptfile.OpenGlobal(path, cryptfile.ReadOnly, 27) != nil {
		return 0
	}
	if noxflags.HasGame(2048) {
		ServerCheatGod(false)
	}
	*memmap.PtrUint32(0x5D4594, 527696) = uint32(uint16(resourceGetHP(u)))
	*memmap.PtrUint32(0x5D4594, 527700) = uint32(uint16(resourceGetMana(u)))
	controlResetPlayer(u)
	core.Players.SetXxx(u, 1)
	r := playerFileStream()
	for {
		id := r.word(0)
		if id == 0 {
			cryptfile.Close()
			controlLevelFromXP(u)
			resourceDamage(u, int32(uint32(uint16(resourceGetMaxHP(u)))-*memmap.PtrUint32(0x5D4594, 527696)))
			resourceSubMana(u, int32(uint32(uint16(resourceGetMaxMana(u)))-*memmap.PtrUint32(0x5D4594, 527700)))
			resourceHPHistory(u)
			core.Players.SetXxx(u, 0)
			return 1
		}
		length, _ := r.cf.ReadAlignedU32()
		ret, found := playerFileServerSection(id, u, unsafe.Add(unsafe.Pointer(p), 2185))
		if !found {
			r.cf.Seek(int64(int32(length)), io.SeekCurrent)
			continue
		}
		if ret == 0 {
			for it := u.InvFirstItem; it != nil; {
				next := it.InvNextItem
				GetServer().DelayedDelete(it)
				it = next
			}
			core.Players.SetXxx(u, 0)
			cryptfile.Close()
			return 0
		}
	}
}
func playerFileClientLoad(path string) int {
	if cryptfile.OpenGlobal(path, cryptfile.ReadOnly, 27) != nil {
		return 0
	}
	quickbarClearSlots()
	r := playerFileStream()
	for {
		id := r.word(0)
		if id == 0 {
			break
		}
		length, _ := r.cf.ReadAlignedU32()
		found := false
		for off := uintptr(55936); *memmap.PtrUint32(0x587000, off) != 0; off += 12 {
			if *memmap.PtrUint32(0x587000, off+4) != id {
				continue
			}
			found = true
			if callPlayerFileSection(*memmap.PtrPtr(0x587000, off+8)) == 0 {
				cryptfile.Close()
				return 0
			}
			break
		}
		if !found {
			r.cf.Seek(int64(int32(length)), io.SeekCurrent)
		}
	}
	cryptfile.Close()
	playerFileCopyString(memmap.PtrOff(0x85B3FC, 10984), path)
	return 1
}
func playerFileClientWrite(info unsafe.Pointer, all int) int {
	base := memmap.PtrOff(0x85B3FC, 10980)
	copy(unsafe.Slice((*byte)(base), int(unsafe.Sizeof(server.SaveGameInfo{}))), unsafe.Slice((*byte)(info), int(unsafe.Sizeof(server.SaveGameInfo{}))))
	if cryptfile.OpenGlobal(alloc.GoString((*byte)(unsafe.Add(base, 4))), cryptfile.WriteOnly, 27) != nil {
		return 0
	}
	r := playerFileStream()
	for off := uintptr(55936); *memmap.PtrUint32(0x587000, off) != 0; off += 12 {
		id := *memmap.PtrUint32(0x587000, off+4)
		if all == 0 && id != 1 {
			continue
		}
		r.raw(memmap.PtrOff(0x587000, off+4), 4)
		r.cf.SectionStart()
		ret := callPlayerFileSection(*memmap.PtrPtr(0x587000, off+8))
		r.cf.SectionEnd()
		if ret == 0 {
			cryptfile.Close()
			return 0
		}
	}
	r.word(0)
	cryptfile.Close()
	return 1
}
func playerFileExtract(path string, out unsafe.Pointer) {
	f, err := binfile.BinfileOpen(path, binfile.ReadOnly)
	if err != nil {
		return
	}
	if f.SetKey(27) != nil {
		return
	}
	var buf [4]byte
	for {
		clear(buf[:])
		f.Read(buf[:])
		id := binary.LittleEndian.Uint32(buf[:])
		if id == 0 {
			*(*uint32)(out) = 0
			f.Close()
			return
		}
		before := int32(f.Written())
		clear(buf[:])
		f.ReadAligned(buf[:])
		length := binary.LittleEndian.Uint32(buf[:])
		after := int32(f.Written())
		switch id {
		case 2, 3, 4, 8, 5, 6, 9, 10, 11:
		default:
			f.FileSeek(int64(int32(length)), io.SeekCurrent)
			continue
		}
		delta := after - before
		*(*uint32)(out) = id
		lengthAt := unsafe.Add(out, int(delta)-4)
		*(*uint32)(lengthAt) = length
		out = unsafe.Add(lengthAt, 8)
		var value [1]byte
		for n := int32(length); n > 0; n-- {
			f.Read(value[:])
			*(*byte)(out) = value[0]
			out = unsafe.Add(out, 1)
		}
	}
}
