//go:build !server && !linux

package opennox

// restartWithSoftwareGL does nothing: the software OpenGL fallback is for Mesa on Linux.
func restartWithSoftwareGL() {}
