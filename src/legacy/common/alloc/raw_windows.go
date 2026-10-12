package alloc

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

func rawMap(size int) []byte {
	addr, err := windows.VirtualAlloc(0, uintptr(size), windows.MEM_RESERVE|windows.MEM_COMMIT, windows.PAGE_READWRITE)
	if err != nil || addr == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(addr)), size)
}
func rawUnmap(data []byte) error {
	return windows.VirtualFree(uintptr(unsafe.Pointer(&data[0])), 0, windows.MEM_RELEASE)
}
