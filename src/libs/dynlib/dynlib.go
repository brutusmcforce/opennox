// Package dynlib loads native libraries (SDL3, OpenAL) at runtime, without cgo.
package dynlib

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Lib is a loaded native library.
type Lib struct {
	Name string
	h    uintptr
}

// candidates returns where to look for a library: next to the executable (and in the Frameworks folder of a macOS
// app bundle) first, so that the libraries shipped with OpenNox win, then the system's search path.
// Windows already searches the executable's folder first.
func candidates(names []string) []string {
	if runtime.GOOS == "windows" {
		return names
	}
	var out []string
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		for _, name := range names {
			if filepath.IsAbs(name) {
				continue
			}
			out = append(out, filepath.Join(dir, name))
			if runtime.GOOS == "darwin" {
				out = append(out, filepath.Join(dir, "..", "Frameworks", name))
			}
		}
	}
	return append(out, names...)
}

// Open loads the first library from the list that can be loaded.
func Open(names ...string) (*Lib, error) {
	var errs []string
	for _, name := range candidates(names) {
		if filepath.IsAbs(name) {
			if _, err := os.Stat(name); err != nil {
				continue // not there: only report the system search below
			}
		}
		h, err := open(name)
		if err == nil {
			return &Lib{Name: name, h: h}, nil
		}
		errs = append(errs, err.Error())
	}
	return nil, fmt.Errorf("cannot load %s (is it installed?): %s", strings.Join(names, " or "), strings.Join(errs, "; "))
}

// Sym returns the address of a function in the library, or 0 if there's no such function.
func (l *Lib) Sym(name string) uintptr {
	if l == nil || l.h == 0 {
		return 0
	}
	p, err := sym(l.h, name)
	if err != nil {
		return 0
	}
	return p
}

// Syms looks up functions by name and stores their addresses. It returns an error listing the missing ones.
func (l *Lib) Syms(funcs map[string]*uintptr) error {
	var missing []string
	for name, p := range funcs {
		if *p = l.Sym(name); *p == 0 {
			missing = append(missing, name)
		}
	}
	if len(missing) != 0 {
		return fmt.Errorf("%s: missing functions: %s", l.Name, strings.Join(missing, ", "))
	}
	return nil
}

// CString returns a NUL-terminated copy of a string, to pass to native functions.
func CString(s string) *byte {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return &b[0]
}
