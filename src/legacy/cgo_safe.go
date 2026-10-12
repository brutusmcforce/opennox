//go:build safe

package legacy

import "github.com/brutusmcforce/opennox/common/memmap"

const cgoSafe = true

func init() {
	memmap.SetRuntimeChecks(true)
}
