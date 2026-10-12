package legacy

import "os"

// legacyAbort never returns or unwinds Go cleanup. Linux first raises the same
// fatal signal as the former C boundary; exit also covers a blocked/handled signal.
func legacyAbort() {
	platformAbort()
	os.Exit(2)
}
