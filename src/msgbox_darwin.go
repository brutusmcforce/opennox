//go:build darwin && !server

package opennox

import (
	"os/exec"
	"strings"
)

// nativeErrorMessage shows an error without SDL, using AppleScript: the app bundle has no console to print to.
func nativeErrorMessage(msg string) {
	quote := strings.NewReplacer(`\`, `\\`, `"`, `\"`)
	_ = exec.Command("osascript", "-e", `display alert "OpenNox" message "`+quote.Replace(msg)+`" as critical`).Run()
}
