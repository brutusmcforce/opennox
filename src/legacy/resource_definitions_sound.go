package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"os"
	"unsafe"
)

var resourceSoundHead unsafe.Pointer

func resourceToken(f *binfile.Binfile, out []byte) bool {
	n := 0
	leading := true
	var previous byte
	for {
		c, ok := monsterDefinitionReadByte(f)
		if !ok {
			return false
		}
		if resourceSpace(c) {
			if !leading {
				out[n] = 0
				return true
			}
		} else {
			leading = false
			if c == '/' && previous == '/' {
				// Preserve the existing sound-set SkipLine boundary (first non-newline).
				if err := f.SkipLine(); err != nil {
					binfile.Log.Println(err)
				}
				n = 0
				leading = true
			} else {
				if n >= len(out)-1 {
					return false
				}
				out[n] = c
				n++
			}
		}
		previous = c
	}
}
func resourceSoundSetField(row unsafe.Pointer, name string, id uint32) bool {
	for off := uintptr(64704); ; off += 8 {
		p := (*byte)(memmap.LoadPtr(0x587000, off))
		if p == nil {
			return false
		}
		if alloc.GoString(p) == name {
			*resourceWord(row, uintptr(memmap.Uint32(0x587000, off+4))) = id
			return true
		}
	}
}
func resourceSoundLoad(path string) int {
	f, err := binfile.BinfileOpen(path, binfile.ReadOnly)
	if err != nil {
		if !os.IsNotExist(err) {
			binfile.Log.Println(err)
		}
		return 0
	}
	// Keep partial list results on errors, but release this reader on every exit.
	defer f.Close()
	if err = f.SetKey(5); err != nil {
		binfile.Log.Println(err)
		return 0
	}
	var a, b [256]byte
	for resourceToken(f, a[:]) {
		row, _ := alloc.Calloc(1, 84)
		*(*unsafe.Pointer)(unsafe.Add(row, 76)) = resourceSoundHead
		resourceSoundHead = row
		name, _ := alloc.CString(alloc.GoStringS(a[:]))
		*(*unsafe.Pointer)(row) = unsafe.Pointer(name)
		for resourceToken(f, a[:]) && alloc.GoStringS(a[:]) != "END" && resourceToken(f, b[:]) {
			if !resourceSoundSetField(row, alloc.GoStringS(a[:]), uint32(sound.ByName(alloc.GoStringS(b[:])))) {
				return 0
			}
		}
	}
	return 1
}
func resourceSoundByName(name string) unsafe.Pointer {
	for p := resourceSoundHead; p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 76)) {
		if alloc.GoString(*(**byte)(p)) == name {
			return p
		}
	}
	return nil
}
func resourceMonsterSound(u *server.Object) unsafe.Pointer {
	if u == nil || uint32(u.ObjClass)&2 == 0 {
		return nil
	}
	return *(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 488))
}
