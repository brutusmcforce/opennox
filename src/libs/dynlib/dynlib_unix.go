//go:build !windows

package dynlib

import "github.com/ebitengine/purego"

func open(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}

func sym(h uintptr, name string) (uintptr, error) {
	return purego.Dlsym(h, name)
}
