package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

var onlineConnectionFailed, onlineSessionStatus uint32
var onlineRetryPending, onlineRetryActive, onlineRetryDeadline, onlineRetryOrigin uint32

func onlineSessionReset() uint32 {
	onlineRetryPending = 0
	onlineRetryActive = 0
	onlineRetryDeadline = 0
	onlineRetryOrigin = 0
	return 0
}
func onlineSessionRetry() uint32 {
	NetworkLogPrint(fmt.Sprintf("RECON: TryReconnectAgain called on frame (%d)", int32(gameFrame())))
	onlineRetryActive = 0
	result := gameFrame() + 120*gameFPS()
	onlineRetryDeadline = gameFrame() + 120*gameFPS()
	return result
}
func onlineSessionStart() {
	if onlineRetryPending == 1 || onlineRetryActive == 1 || onlineRetryDeadline != 0 || onlineRetryOrigin != 0 {
		return
	}
	NetworkLogPrint(fmt.Sprintf("RECON: Starting reconnection process frame (%d)", int32(gameFrame())))
	onlineRetryPending = 1
	onlineRetryActive = 0
	onlineRetryOrigin = gameFrame()
	onlineRetryDeadline = gameFrame() + 120*gameFPS()
}
func onlineSessionLogin() uint32 {
	clear(unsafe.Slice((*byte)(memmap.PtrOff(0x85B3FC, 10308)), 108))
	// The former account selector has no writer and remains at its initial -1.
	return 0
}
func onlineSessionAttempt() uint32 {
	if gameFrame()-onlineRetryOrigin > 3600*gameFPS() {
		onlineSessionReset()
		return onlineSessionMarker(1)
	}
	if onlineRetryPending == 0 {
		return 0
	}
	if onlineRetryActive != 0 {
		return onlineRetryActive
	}
	NetworkLogPrint("RECON: Attempting to re-login")
	onlineConnectionFailed = 0
	onlineSessionLogin()
	return onlineSessionRetry()
}
func onlineSessionMarker(v uint32) uint32   { *memmap.PtrUint32(0x5D4594, 528268) = v; return v }
func onlineSessionBriefing(v uint32) uint32 { *memmap.PtrUint32(0x5D4594, 527720) = v; return v }
func onlineSessionMap() uint32              { *memmap.PtrUint32(0x5D4594, 371700) = 1; return 0 }
func onlineSessionEvent() uint32 {
	// All twelve shipped legacy queue pointers remain nil; their allocator is gone.
	// The live disconnect path still reports the unavailable queue through this hook.
	return uint32(Sub_41E300(11))
}
func onlineSessionListCleanup() uint32 {
	// The channel/user lists have no population path. Their cleanup notification
	// remains observable when the current status is 7.
	if onlineSessionStatus == 7 {
		return onlineSessionEvent()
	}
	return onlineSessionStatus
}
