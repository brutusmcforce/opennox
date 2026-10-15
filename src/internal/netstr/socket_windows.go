//go:build windows

package netstr

import (
	"syscall"
	"unsafe"
)

var procIoctlSocket = syscall.NewLazyDLL("ws2_32.dll").NewProc("ioctlsocket")

const fionread = 0x4004667F // FIONREAD in winsock2.h

func netCanRead(fd uintptr) (uint32, syscall.Errno) {
	var n uint32 // u_long is 32 bits on Windows
	r, _, err := procIoctlSocket.Call(fd, fionread, uintptr(unsafe.Pointer(&n)))
	if r != 0 {
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			return n, errno
		}
		return n, 12345 // TODO
	}
	return n, 0
}
