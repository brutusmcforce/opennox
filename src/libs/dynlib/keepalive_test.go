package dynlib_test

import (
	"os"
	"strings"
	"testing"
)

// The native bindings pass Go pointers to native functions as uintptr through a call wrapper of their own.
// purego.SyscallN keeps such arguments alive only for its own call site, so each wrapper has to repeat the
// pragma; without it the Go pointers are free to be collected (or their stack to move) during the native call,
// which corrupts memory only under GC pressure. Guard the pragma, it is easy to drop by accident.
func TestCallWrappersKeepPointersAlive(t *testing.T) {
	for _, name := range []string{"../sdl/sdl.go", "../gl/gl.go"} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Error(err)
			continue
		}
		if !strings.Contains(string(data), "//go:uintptrescapes\nfunc call(") {
			t.Errorf("%s: func call must be marked //go:uintptrescapes", name)
		}
	}
}
