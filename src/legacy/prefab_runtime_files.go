package legacy

import (
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
)

func prefabReadRaw(handle uint32, dst []byte) int {
	n, _ := fileByHandle((*FILE)(mapRoomPointer(handle))).Read(dst)
	return n
}
func prefabRawWord(handle uint32) uint32 {
	var b [4]byte
	prefabReadRaw(handle, b[:])
	return binary.LittleEndian.Uint32(b[:])
}
func prefabRawByte(handle uint32) byte { var b [1]byte; prefabReadRaw(handle, b[:]); return b[0] }
func prefabLibrary() uint32 {
	(*prefabGlobal(prefabCount)) = 0
	for _, entry := range [][2]int{{prefabPath, 2048}, {prefabAlternate, 2048}, {prefabMetadata, 2048 * 76}} {
		if (*prefabGlobal(entry[0])) == 0 {
			(*prefabGlobal(entry[0])) = mapRoomRaw(mapRoomCalloc(1, uintptr(entry[1])))
		}
		if (*prefabGlobal(entry[0])) == 0 {
			return 0
		}
	}
	path := (*prefabGlobal(prefabPath))
	if *(*byte)(mapRoomPointer(path)) == 0 || prefabOpen(path) == 0 {
		return 0
	}
	defer prefabClose()
	handle := (*prefabGlobal(prefabFile))
	if prefabRawWord(handle) != 0xcafedead {
		return 0
	}
	f := fileByHandle((*FILE)(mapRoomPointer(handle)))
	for {
		size := prefabRawWord(handle)
		if size == 0 {
			return 1
		}
		if (*prefabGlobal(prefabCount)) >= 2048 {
			return 0
		}
		pos, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return 0
		}
		row := (*prefabGlobal(prefabMetadata)) + (*prefabGlobal(prefabCount))*76
		*prefabWord(row, 72) = uint32(pos - 4)
		length := int(prefabRawByte(handle))
		// A C metadata name has 64 bytes including the terminator.
		if length >= 64 {
			return 0
		}
		name := make([]byte, length)
		if prefabReadRaw(handle, name) != length {
			return 0
		}
		// strcpy only overwrites the used prefix, retaining the rest of a reused row.
		for i, c := range name {
			*(*byte)(unsafe.Add(mapRoomPointer(row), i)) = c
			if c == 0 {
				length = i
				break
			}
		}
		*(*byte)(unsafe.Add(mapRoomPointer(row), length)) = 0
		prefabRawByte(handle)
		prefabRawByte(handle)
		*prefabWord(row, 64) = prefabRawWord(handle)
		*prefabWord(row, 68) = prefabRawWord(handle)
		// The record byte count includes name-length, name and the ten fixed bytes.
		rest := int32(size) - 1 - int32(len(name)) - 10
		if rest < 0 {
			return 0
		}
		if _, err := f.Seek(int64(rest), io.SeekCurrent); err != nil {
			return 0
		}
		(*prefabGlobal(prefabCount))++
	}
}
func prefabScriptScan(handle, output uint32) uint32 {
	file := fileByHandle((*FILE)(mapRoomPointer(handle)))
	read := func(b []byte) bool {
		n, err := file.Bin.Read(b)
		if err != nil {
			file.Err = err
		}
		return n == len(b)
	}
	word := func() (uint32, bool) { var b [4]byte; ok := read(b[:]); return binary.LittleEndian.Uint32(b[:]), ok }
	seek := func(n int32) bool { return file.Bin.FileSeek(int64(n), io.SeekCurrent) == nil }
	length, ok := word()
	if !ok || length > 1024 {
		return 0
	}
	if !read(make([]byte, int(length))) {
		return 0
	}
	if !read(unsafe.Slice((*byte)(mapRoomPointer(output)), 4)) {
		return 0
	}
	count, ok := word()
	if !ok {
		return 0
	}
	for i := int32(0); i < int32(count); i++ {
		var op [1]byte
		if !read(op[:]) || !seek(1) {
			return 0
		}
		if op[0] > 36 {
			return 0
		}
		offset := uintptr(op[0]) * 268
		args := int(memmap.Uint8(0x587000, 218640+offset))
		for j := 0; j < args; j++ {
			kind := memmap.Uint32(0x587000, 218648+offset+uintptr(j)*8)
			switch kind {
			case 0, 3, 4, 5, 6:
				if !seek(4) {
					return 0
				}
			case 1:
				if !seek(8) {
					return 0
				}
			case 2, 7:
				var n [1]byte
				if !read(n[:]) || !seek(int32(n[0])) {
					return 0
				}
			}
		}
	}
	return count
}
