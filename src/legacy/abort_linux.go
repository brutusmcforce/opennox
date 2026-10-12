package legacy

import (
	"runtime"
	"syscall"
)

func platformAbort() {
	runtime.LockOSThread()
	_ = syscall.Tgkill(syscall.Getpid(), syscall.Gettid(), syscall.SIGABRT)
}
