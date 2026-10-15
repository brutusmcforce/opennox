//go:build windows && !server

package opennox

import "golang.org/x/sys/windows"

// nativeErrorMessage shows an error without SDL: the Windows app is built without a console to print to.
func nativeErrorMessage(msg string) {
	title, err := windows.UTF16PtrFromString("OpenNox")
	if err != nil {
		return
	}
	text, err := windows.UTF16PtrFromString(msg)
	if err != nil {
		return
	}
	_, _ = windows.MessageBox(0, text, title, windows.MB_OK|windows.MB_ICONERROR)
}
