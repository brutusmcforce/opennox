//go:build windows

package dynlib

import "syscall"

func open(name string) (uintptr, error) {
	h, err := syscall.LoadLibrary(name)
	return uintptr(h), err
}

func sym(h uintptr, name string) (uintptr, error) {
	return syscall.GetProcAddress(syscall.Handle(h), name)
}
