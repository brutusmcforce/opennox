package client

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client/noxrender"
)

// Draw data of sprites, as loaded from the thing database (see legacy client_sprite_parse.go).

// SpriteStaticData is the draw data of static sprites.
type SpriteStaticData struct {
	Size  uint32 // 0, 0
	Image noxrender.ImageHandle
}

// SpriteStaticRandomData is the draw data of sprites with a random static frame.
type SpriteStaticRandomData struct {
	Size   uint32                 // 0, 0
	Frames *noxrender.ImageHandle // 1, 4: [Count]
	Count  byte                   // 2, 8
	_      [3]byte
}

// SpriteAnimationData is the draw data of animated sprites.
type SpriteAnimationData struct {
	Size         uint32                 // 0, 0
	Frames       *noxrender.ImageHandle // 1, 4: [Count]
	Count, Delay byte                   // 2, 8
	Reserved     uint16
	Kind         AnimKind // 3, 12
}

// SpriteStateData is the draw data of sprites animated by state.
type SpriteStateData struct {
	Size   uint32 // 0, 0
	States [3]AnimationVector
}

// SpriteFrame returns frame i of a frame array.
func SpriteFrame(frames *noxrender.ImageHandle, i int) noxrender.ImageHandle {
	return *(*noxrender.ImageHandle)(unsafe.Add(unsafe.Pointer(frames), uintptr(i)*unsafe.Sizeof(noxrender.ImageHandle(nil))))
}

// SpriteFrames returns the frame array of sprite draw data that starts with a size and a frame array
// (SpriteStaticRandomData, SpriteAnimationData).
func SpriteFrames(dd unsafe.Pointer) *noxrender.ImageHandle {
	return (*SpriteAnimationData)(dd).Frames
}
