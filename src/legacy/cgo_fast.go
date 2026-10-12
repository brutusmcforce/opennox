//go:build !safe

package legacy

import "os"

var cgoSafe = os.Getenv("NOX_SAFE") == "true"
