package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/legacy/timer"
	"strings"
	"unsafe"
)

func audioAssetSlot(id int32) unsafe.Pointer {
	if dword_5d4594_1045432 == 0 || id < 0 || id >= 1023 {
		return nil
	}
	return memmap.PtrOff(0x5D4594, 840628+200*uintptr(id))
}
func audioAssetDelay(p unsafe.Pointer) int32 {
	if dword_5d4594_1045432 == 0 {
		return 0
	}
	return *(*int32)(unsafe.Add(p, 64))
}

func audioAssetSample(catalog unsafe.Pointer, key *byte) int32 {
	base := *(*unsafe.Pointer)(catalog)
	lo, hi := uint32(0), *(*uint32)(unsafe.Add(catalog, 4))
	for lo < hi {
		mid := (lo + hi) / 2
		p := unsafe.Add(base, 36*uintptr(mid))
		cmp := int(textCompareNarrow(key, (*byte)(p)))
		if cmp < 0 {
			hi = mid
		} else if cmp > 0 {
			lo = mid + 1
		} else {
			return int32(mid)
		}
	}
	return -1
}
func audioAssetReadName(f *binfile.MemFile, scratch []byte, signed bool) string {
	n := int(f.ReadU8())
	if signed {
		n = int(int8(n))
	}
	f.Read(scratch[:n])
	scratch[n] = 0
	return alloc.GoStringS(scratch[:n+1])
}
func audioAssetStore(slot unsafe.Pointer, off uintptr, v uint32) {
	*(*uint32)(unsafe.Add(slot, off)) = v
}
func audioAssetCatalog() unsafe.Pointer { return unsafe.Pointer(uintptr(dword_5d4594_1045420)) }

func audioAssetEvent(f *binfile.MemFile, scratch []byte) int {
	name := audioAssetReadName(f, scratch, false)
	id := sound.ByName(name)
	slot := audioAssetSlot(int32(id))
	if id == 0 || slot == nil {
		thingSkipAVNTInner(f)
		return 1
	}
	for {
		switch f.ReadU8() {
		case 0:
			audioAssetStore(slot, 0, 1)
			return 1
		case 1:
			audioAssetStore(slot, 48, uint32(f.ReadU8()))
		case 2:
			audioAssetStore(slot, 4, uint32(f.ReadU8()))
		case 3:
			(*timer.Timer)(unsafe.Add(slot, 16)).Init(163 * int32(f.ReadU8()))
		case 4:
			audioAssetStore(slot, 56, uint32(f.ReadU8()))
		case 5:
			audioAssetStore(slot, 60, uint32(f.ReadU8()))
		case 6:
			audioAssetStore(slot, 76, uint32(int32(f.ReadI8())))
			audioAssetStore(slot, 80, uint32(int32(f.ReadI8())))
		case 7:
			count := 0
			for {
				n := int(f.ReadU8())
				if n == 0 {
					break
				}
				f.Read(scratch[:n])
				scratch[n] = 0
				if count < 32 {
					index := int16(audioAssetSample(audioAssetCatalog(), &scratch[0]))
					*(*uint16)(unsafe.Add(slot, 128+2*count)) = uint16(index)
					if index != -1 {
						count++
					}
				}
			}
			audioAssetStore(slot, 192, uint32(count))
		case 8:
			audioAssetStore(slot, 68, f.ReadU32())
			audioAssetStore(slot, 72, f.ReadU32())
		case 9:
			if v := f.ReadI16(); v > 0 {
				audioAssetStore(slot, 64, uint32(15*int32(v)))
			}
		case 10:
			audioAssetStore(slot, 8, uint32(int32(f.ReadI16())))
		default:
			return 0
		}
	}
}
func audioAssetRecord(f *binfile.MemFile, scratch []byte) int {
	name := audioAssetReadName(f, scratch, true)
	id := sound.ByName(name)
	slot := audioAssetSlot(int32(id))
	if id == 0 || slot == nil {
		f.Skip(9)
		for {
			n := int(f.ReadI8())
			if n == 0 {
				break
			}
			f.Skip(n)
		}
		return 1
	}
	priority := f.ReadI16()
	audioAssetStore(slot, 4, 2)
	audioAssetStore(slot, 8, uint32(int32(priority)))
	(*timer.Timer)(unsafe.Add(slot, 16)).Init(163 * int32(f.ReadU8()))
	if v := f.ReadI16(); v > 0 {
		audioAssetStore(slot, 64, uint32(15*int32(v)))
	}
	audioAssetStore(slot, 56, uint32(int32(f.ReadI8())))
	audioAssetStore(slot, 76, uint32(int32(f.ReadI8())))
	audioAssetStore(slot, 80, uint32(int32(f.ReadI8())))
	mode := int32(f.ReadI8())
	audioAssetStore(slot, 48, uint32(mode))
	if mode >= 3 {
		return 0
	}
	count := 0
	for {
		n := int(f.ReadI8())
		if n == 0 {
			break
		}
		f.Read(scratch[:n])
		scratch[n] = 0
		if dot := strings.LastIndexByte(alloc.GoStringS(scratch[:n+1]), '.'); dot >= 0 {
			scratch[dot] = 0
		}
		index := int16(audioAssetSample(audioAssetCatalog(), &scratch[0]))
		*(*uint16)(unsafe.Add(slot, 128+2*count)) = uint16(index)
		if index != -1 {
			count++
		}
	}
	audioAssetStore(slot, 0, 1)
	audioAssetStore(slot, 192, uint32(count))
	return 1
}
func audioAssetDefinitions(f *binfile.MemFile, scratch []byte) int {
	count := f.ReadI32()
	for i := int32(0); i < count; i++ {
		if audioAssetRecord(f, scratch) == 0 {
			return 0
		}
	}
	return 1
}
