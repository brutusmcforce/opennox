package legacy

import (
	"math"

	"github.com/brutusmcforce/opennox/common/memmap"
)

func durabilityBand(current, maximum uint16) int {
	if maximum == 0 {
		return 4
	}
	if current == maximum {
		return 0
	}
	// The thresholds remain shared with C. Preserve their ordered comparisons,
	// including equality at half/quarter and nonfinite threshold values.
	if float64(current) >= float64(maximum)*math.Float64frombits(uint64(qword_581450_9544)) {
		return 1
	}
	if float64(current) < float64(maximum)*memmap.Float64(0x581450, 9608) {
		return 3
	}
	return 2
}
