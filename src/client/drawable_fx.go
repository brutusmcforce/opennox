package client

import (
	"image"
	"unsafe"
)

// DrawableFX is an effect attached to a drawable (see Drawable.Field_114 and legacy combatFX).
type DrawableFX struct {
	Kind, Field4                               uint32         // 0, 4
	History                                    [6]image.Point // 8: last positions, newest first
	Count                                      byte           // 56
	_                                          [3]byte
	Owner                                      *Drawable
	Next, Previous, GlobalNext, GlobalPrevious *DrawableFX // Next is the link of Field_114
}

func (fx *DrawableFX) C() unsafe.Pointer {
	return unsafe.Pointer(fx)
}
