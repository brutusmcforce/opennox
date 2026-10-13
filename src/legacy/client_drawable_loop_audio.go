package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
)

func clientDrawableLoopAudio(source, listener *client.Drawable) {
	metadata := audioAssetSlot(int32(source.AudioLoop))
	vp := drawViewport()
	if metadata == nil || vp == nil {
		return
	}
	volume, pan := int32(0), int32(0)
	if source.ObjFlags&0x1000000 != 0 && source.Flags70Val&12 == 0 {
		dx, dy := int32(listener.PosVec.X-source.PosVec.X), int32(listener.PosVec.Y-source.PosVec.Y)
		limit := audioAssetDelay(metadata)
		if dx < limit && dy < limit && limit > 0 {
			sum := dx*dx + dy*dy + 1
			// The 386 C code converts sqrt through signed 64-bit before narrowing.
			distance := int32(int64(math.Sqrt(float64(sum))))
			if distance < limit {
				volume = 100 * (limit - distance) / limit
				if volume > 100 {
					volume = 100
				}
				if volume < 0 {
					volume = 0
				}
				half := int32(nox_win_width) / 2
				if half != 0 {
					pan = 50 * int32(source.PosVec.X-vp.World.Max.X-vp.Screen.Min.X) / half
				}
			}
		}
	}
	handle := (*audioEventHandle)(unsafe.Pointer(&source.Audio))
	current := audioEventHandleGet(handle)
	if volume != 0 {
		if current != nil {
			audioEventFadePan(current, pan)
			audioEventFadeVolume(current, volume)
		} else {
			if event := audioEventNew((*audioEventMetadata)(metadata)); event != nil {
				audioEventSetVolume(event, volume)
				audioEventSetPan(event, pan)
				audioEventHandleSet(handle, event)
			}
		}
	} else if current != nil {
		audioEventDelete(current)
	}
}
