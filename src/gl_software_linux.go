//go:build !server

package opennox

import (
	"os"
	"syscall"
)

// glSoftwareEnv makes Mesa render OpenGL on the CPU (llvmpipe).
const glSoftwareEnv = "LIBGL_ALWAYS_SOFTWARE"

// restartWithSoftwareGL replaces the process with a new one with Mesa's software renderer, unless it's already used.
// It only returns if it didn't restart.
func restartWithSoftwareGL() {
	if os.Getenv(glSoftwareEnv) != "" {
		return // already tried
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	// Mesa picks the driver when the display is opened, so it can't be changed in this process
	_ = syscall.Exec(exe, os.Args, append(os.Environ(), glSoftwareEnv+"=1"))
}
