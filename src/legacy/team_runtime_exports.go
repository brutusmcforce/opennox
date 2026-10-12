package legacy

import ()

func teamRuntimeBool(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
