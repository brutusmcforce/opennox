//go:build unix && !linux

package netstr

import (
	"errors"
	"syscall"

	"golang.org/x/sys/unix"
)

// fionread is FIONREAD on macOS and the BSDs (_IOR('f', 127, int)); x/sys/unix doesn't define it for all of them.
const fionread = 0x4004667f

// netCanRead returns the number of bytes waiting on the socket (macOS and the BSDs have no TIOCINQ).
func netCanRead(fd uintptr) (uint32, syscall.Errno) {
	n, err := unix.IoctlGetInt(int(fd), fionread)
	if err != nil {
		var errno syscall.Errno
		if errors.As(err, &errno) {
			return 0, errno
		}
		return 0, syscall.EINVAL
	}
	return uint32(n), 0
}
