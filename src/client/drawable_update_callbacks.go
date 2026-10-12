package client

import (
	"runtime"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/noxrender"
)

var drawableUpdateCallbacksGo = make(map[unsafe.Pointer]func(*noxrender.Viewport, *Drawable) int32)

// RegisterDrawableUpdateCallbackGo binds a stable identity during initialization.
// Drawable and object-type records retain their original pointer-sized fields.
func RegisterDrawableUpdateCallbackGo(key unsafe.Pointer, fn func(*noxrender.Viewport, *Drawable) int32) {
	if key == nil || fn == nil {
		panic("invalid drawable update callback")
	}
	if _, ok := drawableUpdateCallbacksGo[key]; ok {
		panic("drawable update callback already registered")
	}
	drawableUpdateCallbacksGo[key] = fn
}

// CallDrawableUpdateResult preserves the complete signed C-int result. Callers
// require a registered callback.
func CallDrawableUpdateResult(key unsafe.Pointer, vp *noxrender.Viewport, dr *Drawable) int32 {
	var result int32
	if fn := drawableUpdateCallbacksGo[key]; fn != nil {
		result = fn(vp, dr)
	} else {
		panic("unregistered drawable update callback")
	}
	runtime.KeepAlive(vp)
	runtime.KeepAlive(dr)
	return result
}

// CallDrawableUpdateDiscard preserves the secondary update slot’s void convention.
func CallDrawableUpdateDiscard(key unsafe.Pointer, vp *noxrender.Viewport, dr *Drawable) {
	if fn := drawableUpdateCallbacksGo[key]; fn != nil {
		fn(vp, dr)
	} else {
		panic("unregistered drawable update callback")
	}
	runtime.KeepAlive(vp)
	runtime.KeepAlive(dr)
}
