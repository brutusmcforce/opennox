package legacy

import (
	"github.com/brutusmcforce/opennox/internal/binfile"
	"io"
)

func thingSkipSized(f *binfile.MemFile) { n := int(f.ReadU8()); f.Skip(n) }
func thingSkipImageRef(f *binfile.MemFile) {
	if f.ReadI32() == -1 {
		f.Skip(1)
		thingSkipSized(f)
	}
}
func thingSkipAUD(f *binfile.MemFile) int {
	count := f.ReadI32()
	for i := int32(0); i < count; i++ {
		n := int(f.ReadU8())
		f.Skip(n + 9)
		for {
			n = int(f.ReadU8())
			if n == 0 {
				break
			}
			f.Skip(n)
		}
	}
	return 1
}
func thingSkipSpells(f *binfile.MemFile) int {
	count := f.ReadI32()
	for i := int32(0); i < count; i++ {
		n := int(f.ReadU8())
		f.Skip(n + 3)
		thingSkipSized(f)
		thingSkipImageRef(f)
		thingSkipImageRef(f)
		f.Skip(4)
		thingSkipSized(f)
		n = int(f.ReadI16())
		f.Skip(n)
		for j := 0; j < 3; j++ {
			thingSkipSized(f)
		}
	}
	return 1
}
func thingSkipAbilities(f *binfile.MemFile) int {
	count := f.ReadI32()
	for i := int32(0); i < count; i++ {
		n := int(f.ReadU8())
		f.Skip(n + 1)
		for j := 0; j < 3; j++ {
			thingSkipImageRef(f)
		}
		thingSkipSized(f)
		n = int(f.ReadI16())
		f.Skip(n)
		for j := 0; j < 3; j++ {
			thingSkipSized(f)
		}
	}
	return 1
}
func thingSkipImages(f *binfile.MemFile) int {
	count := f.ReadI32()
	for i := int32(0); i < count; i++ {
		thingSkipSized(f)
		kind := f.ReadU8()
		frames := 1
		if kind == 2 {
			frames = int(f.ReadU8())
			f.Skip(1)
			thingSkipSized(f)
		}
		for j := 0; j < frames; j++ {
			thingSkipImageRef(f)
		}
	}
	return 1
}
func thingSkipAVNTInner(f *binfile.MemFile) int {
	for {
		switch f.ReadU8() {
		case 0:
			return 1
		case 1, 2, 3, 4, 5:
			f.Skip(1)
		case 6, 9, 10:
			f.Skip(2)
		case 7:
			for {
				n := int(f.ReadU8())
				if n == 0 {
					break
				}
				f.Skip(n)
			}
		case 8:
			f.Skip(8)
		default:
			return 0
		}
	}
}
func thingSkipAVNT(f *binfile.MemFile) int { thingSkipSized(f); return thingSkipAVNTInner(f) }

func thingReadAlignedByte(f *binfile.MemFile) byte {
	off, _ := f.Seek(0, io.SeekCurrent)
	if extra := off % 8; extra != 0 {
		f.Skip(int(8 - extra))
	}
	val := f.ReadU8()
	f.Skip(7)
	return val
}
func thingSkipWall(f *binfile.MemFile, scratch []byte) int {
	readName := func() { n := int(f.ReadU8()); f.Read(scratch[:n]); scratch[n] = 0 }
	f.Skip(4)
	readName()
	f.Skip(14)
	materials := int(thingReadAlignedByte(f))
	for i := 0; i < materials; i++ {
		if i == 8 {
			return 0
		}
		readName()
	}
	for i := 0; i < 3; i++ {
		readName()
	}
	f.Skip(1)
	for side := 0; side < 15; side++ {
		count := int(thingReadAlignedByte(f))
		for i := 0; i < count; i++ {
			for face := 0; face < 4; face++ {
				f.Skip(8)
				thingSkipImageRef(f)
			}
		}
	}
	if f.ReadU32() == 0x454e4420 {
		return 1
	}
	return 0
}
