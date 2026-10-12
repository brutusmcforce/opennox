//go:build linux

package netstr

import (
	"syscall"
	"unsafe"
)

func netCanRead(fd uintptr) (uint32, syscall.Errno) {
	var n uint32
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCINQ), uintptr(unsafe.Pointer(&n)))
	return n, err
}
