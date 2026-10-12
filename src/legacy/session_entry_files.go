package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/libs/datapath"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

func sessionFileExists(path string) bool {
	_, err := ifs.Stat(path)
	// The existing access adapter returns -2 for other errors; callers distinguish only -1.
	return !os.IsNotExist(err)
}
func sessionFileOpen(path string, key int) bool {
	if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, key); err != nil {
		if !os.IsNotExist(err) {
			binfile.Log.Println(err)
		}
		return false
	}
	return true
}
func sessionMapValidate(name *byte, checksum uint32) int32 {
	if checksum == 0 {
		return 6
	}
	if name == nil {
		return 0
	}
	path := alloc.GoString(name)
	if !strings.ContainsRune(path, '\\') {
		part := "maps\\" + path
		part = part[:len(part)-4]
		sep := alloc.GoString(memmap.PtrUint8(0x587000, 191672))
		path = part + sep + path
	}
	if !sessionFileExists(path) {
		return 0
	}
	var result int32
	if !sessionFileExists(path) {
		result = 1
	}
	if sessionFileOpen(path, 19) {
		result |= 2
		magic, _ := cryptfile.Global().ReadU32()
		if int32(magic) == -86050098 {
			crc, _ := cryptfile.Global().ReadAlignedU32()
			if crc == checksum {
				result |= 4
			}
		}
		cryptfile.Close()
	}
	return result
}
func sessionSaveSlots() int32 {
	root := saveRoot()
	_ = ifs.Mkdir(root + "\\Save\\")
	var count int32
	if sessionFileExists(root + "\\Save\\AUTOSAVE\\Player.plr") {
		count++
	}
	for slot := 1; slot < 14; slot++ {
		if sessionFileExists(fmt.Sprintf("%s\\Save\\SAVE%04d\\Player.plr", root, slot)) {
			count++
		}
	}
	return count
}
func sessionSaveMetadata(path string, info *server.SaveGameInfo) int32 {
	if i := strings.IndexByte(path, 0); i >= 0 {
		path = path[:i]
	}
	out := unsafe.Slice((*byte)(unsafe.Pointer(info)), 1278)
	out[1028] = memmap.Uint8(0x5D4594, 527728)
	if !sessionFileOpen(path, 27) {
		return 0
	}
	f := cryptfile.Global()
	for {
		id, _ := f.ReadU32()
		if id == 0 {
			break
		}
		size, _ := f.ReadAlignedU32()
		if id != 1 {
			_ = f.Seek(int64(int32(size)), io.SeekCurrent)
			continue
		}
		current := unsafe.Slice(memmap.PtrUint8(0x85B3FC, 10980), 1278)
		saved := [1278]byte{}
		copy(saved[:], current)
		if memmap.Uint32(0x587000, 55936) != 0 && memmap.Uint32(0x587000, 55948) != 0 {
			if callPlayerFileMetadata(*memmap.PtrPtr(0x587000, 55956)) == 0 {
				cryptfile.Close()
				return 0
			}
		}
		copy(out, current)
		copy(current, saved[:])
	}
	cryptfile.Close()
	copy(out[4:], path)
	out[4+len(path)] = 0
	return 1
}
func sessionCharacterCount() int32 {
	root := saveRoot()
	dir := root + "\\Save\\"
	_ = ifs.Mkdir(dir)
	_ = ifs.Chdir(dir)
	paths, _ := filepath.Glob(ifs.Normalize("*.plr"))
	var count int32
	for _, path := range paths {
		st, err := ifs.Stat(path)
		if err != nil || st.IsDir() {
			continue
		}
		var info server.SaveGameInfo
		if sessionSaveMetadata(dir+filepath.Base(path), &info) != 0 && info.Flags&2 != 0 {
			count++
		}
	}
	// restore the working directory to the game directory
	_ = ifs.Chdir(datapath.Data())
	return count
}
