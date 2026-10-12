package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"strings"
	"unsafe"
)

func audioAssetSlot(id int32) *audioEventMetadata {
	if dword_5d4594_1045432 == 0 || id < 0 || id >= 1023 {
		return nil
	}
	return &audioEventMetadataTable[id]
}
func audioAssetDelay(p *audioEventMetadata) int32 {
	if dword_5d4594_1045432 == 0 {
		return 0
	}
	return int32(p.Distance)
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
func audioAssetCatalog() unsafe.Pointer { return memmap.WordPtr(dword_5d4594_1045420) }

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
			slot.Enabled = 1
			return 1
		case 1:
			slot.Priority = int32(uint32(f.ReadU8()))
		case 2:
			slot.Flags = uint32(f.ReadU8())
		case 3:
			slot.Volume.Init(163 * int32(f.ReadU8()))
		case 4:
			slot.Limit = uint32(f.ReadU8())
		case 5:
			slot.Loops = uint32(f.ReadU8())
		case 6:
			slot.PitchMin = int32(uint32(int32(f.ReadI8())))
			slot.PitchMax = int32(uint32(int32(f.ReadI8())))
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
					slot.Samples[count] = index
					if index != -1 {
						count++
					}
				}
			}
			slot.SampleCount = uint32(count)
		case 8:
			slot.DelayMin = f.ReadU32()
			slot.DelayMax = f.ReadU32()
		case 9:
			if v := f.ReadI16(); v > 0 {
				slot.Distance = uint32(15 * int32(v))
			}
		case 10:
			slot.Field8 = uint32(int32(f.ReadI16()))
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
	slot.Flags = 2
	slot.Field8 = uint32(int32(priority))
	slot.Volume.Init(163 * int32(f.ReadU8()))
	if v := f.ReadI16(); v > 0 {
		slot.Distance = uint32(15 * int32(v))
	}
	slot.Limit = uint32(int32(f.ReadI8()))
	slot.PitchMin = int32(uint32(int32(f.ReadI8())))
	slot.PitchMax = int32(uint32(int32(f.ReadI8())))
	mode := int32(f.ReadI8())
	slot.Priority = int32(uint32(mode))
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
		slot.Samples[count] = index
		if index != -1 {
			count++
		}
	}
	slot.Enabled = 1
	slot.SampleCount = uint32(count)
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
