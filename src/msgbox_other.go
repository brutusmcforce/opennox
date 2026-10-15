//go:build !windows && !darwin && !server

package opennox

import (
	"fmt"
	"os"
)

// nativeErrorMessage prints an error without SDL. Everything that can show a dialog here (zenity, kdialog...)
// comes with SDL, so there is nothing left to try.
func nativeErrorMessage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}
