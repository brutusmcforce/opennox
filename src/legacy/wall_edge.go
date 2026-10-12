package legacy

import (
	"image"

	"github.com/brutusmcforce/opennox/client/noxrender"
)

// wallEdgeDraw draws a lit wall sprite. The software renderer is gone, so it's drawn by the GPU renderer
// (see Nox_client_gpu_edgeDraw); the clip region is the intersection of the viewport and the given limits.
func wallEdgeDraw(handle noxrender.ImageHandle, pos image.Point, first, second *[3]uint32, bottom, left, right, crop int) {
	if handle == nil || Nox_client_gpu_edgeDraw == nil {
		return
	}
	clip := image.Rect(
		max(int(int32(dword_5d4594_3807140)), left), int(int32(dword_5d4594_3807136)),
		min(int(int32(dword_5d4594_3807116)), right), min(int(int32(dword_5d4594_3807152)), bottom),
	)
	img := GetClient().R2().GetBag().AsImage(handle)
	Nox_client_gpu_edgeDraw(img, pos,
		[3]int{int(first[0]), int(first[1]), int(first[2])},
		[3]int{int(second[0]), int(second[1]), int(second[2])},
		clip, crop)
}
