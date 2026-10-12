package legacy

import (
	"bytes"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func mapMetadataInfo() uint32 {
	ver := mapSectionWord(3)
	if ver > 3 {
		return 0
	}
	read := cryptfile.Global().ReadOnly()
	if ver >= 1 {
		// These field boundaries are also checksum boundaries.
		off := uintptr(2408)
		for _, n := range []int{64, 512, 16, 64, 64, 128, 128, 256, 128, 32, 4} {
			mapSectionIO(memmap.PtrOff(0x973F18, off), n)
			off += uintptr(n)
		}
		if ver == 2 {
			mapSectionIO(memmap.PtrOff(0x973F18, 3804), 1)
			mapSectionIO(memmap.PtrOff(0x973F18, 3805), 1)
		} else if read {
			*memmap.PtrUint8(0x973F18, 3804) = 2
			*memmap.PtrUint8(0x973F18, 3805) = 16
		}
	}
	if ver < 3 {
		*memmap.PtrUint8(0x973F18, 3806) = memmap.Uint8(0x5D4594, 741376)
		*memmap.PtrUint8(0x973F18, 3838) = memmap.Uint8(0x5D4594, 741380)
	} else {
		for _, off := range []uintptr{3806, 3838} {
			p := memmap.PtrOff(0x973F18, off)
			n := int(mapSectionByte(byte(alloc.Strlen(p))))
			mapSectionIO(p, n)
			*(*byte)(unsafe.Add(p, n)) = 0
		}
	}
	return 1
}

func mapMetadataSetAmbient(v [3]uint32) {
	for i, x := range v {
		*memmap.PtrUint32(0x587000, 142296+uintptr(4*i)) = x
	}
}
func mapMetadataAmbient() uint32 {
	if int16(mapSectionWord(1)) < 1 {
		return 0
	}
	if cryptfile.Global().ReadOnly() {
		values := [3]uint32{mapSectionDword(0), mapSectionDword(0), mapSectionDword(0)}
		mapMetadataSetAmbient(values)
		if noxflags.HasGame(0x200002) {
			GetClient().R2().Data().SetLightColor(noxrender.RGB{R: int(values[0]), G: int(values[1]), B: int(values[2])})
		}
	} else {
		for _, off := range []uintptr{142296, 142300, 142304} {
			mapSectionIO(memmap.PtrOff(0x587000, off), 4)
		}
	}
	return 1
}

func mapMetadataTOC() uint32 {
	if int16(mapSectionWord(1)) > 1 {
		return 0
	}
	Sub_42BFB0()
	if cryptfile.Global().ReadOnly() {
		count := mapSectionWord(0)
		for i := 0; i < int(count); i++ {
			code := mapSectionWord(0)
			n := mapSectionByte(0)
			var buf [256]byte
			mapSectionIO(unsafe.Pointer(&buf[0]), int(n))
			name := buf[:int(n)]
			if end := bytes.IndexByte(name, 0); end >= 0 {
				name = name[:end]
			}
			var id int
			if !noxflags.HasGame(noxflags.GameClient) || noxflags.HasGame(noxflags.GameHost) {
				id = GetServer().S().Types.IndByID(string(name))
			} else {
				id = GetClient().Cli().Things.IndByID(string(name))
			}
			Sub_42C310(int(uint16(id)), code)
		}
	} else {
		Sub_42BFE0()
		mapSectionWord(Sub_42C300())
		s := GetServer().S()
		for id := uint16(0); int(id) < s.Types.Count(); id++ {
			code := Sub_42C2E0(int(id))
			if code == 0 {
				continue
			}
			mapSectionWord(code)
			name := s.Types.ByInd(int(id)).ID()
			n := mapSectionByte(byte(len(name)))
			data := append([]byte(name), 0)
			mapSectionIO(unsafe.Pointer(&data[0]), int(n))
		}
	}
	return 1
}
