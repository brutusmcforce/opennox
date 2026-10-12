package client

import (
	"runtime"
	"unsafe"

	"github.com/brutusmcforce/opennox/client/noxrender"
)

var drawableDrawCallbacksGo = make(map[unsafe.Pointer]func(*noxrender.Viewport, *Drawable) int32)

// RegisterDrawableDrawCallbackGo binds a stable identity during initialization.
// Drawable and object-type records retain their original pointer-sized fields.
func RegisterDrawableDrawCallbackGo(key unsafe.Pointer, fn func(*noxrender.Viewport, *Drawable) int32) {
	if key == nil || fn == nil {
		panic("invalid drawable callback")
	}
	if _, ok := drawableDrawCallbacksGo[key]; ok {
		panic("drawable callback already registered")
	}
	drawableDrawCallbacksGo[key] = fn
}

// CallDrawableDrawResult preserves the complete signed C-int result. Callers
// require a registered callback.
func CallDrawableDrawResult(key unsafe.Pointer, vp *noxrender.Viewport, dr *Drawable) int32 {
	var result int32
	if fn := drawableDrawCallbacksGo[key]; fn != nil {
		result = fn(vp, dr)
	} else {
		panic("unregistered drawable draw callback")
	}
	runtime.KeepAlive(vp)
	runtime.KeepAlive(dr)
	return result
}

// CallDrawableDrawDiscard preserves the void convention used for shield draws.
func CallDrawableDrawDiscard(key unsafe.Pointer, vp *noxrender.Viewport, dr *Drawable) {
	if fn := drawableDrawCallbacksGo[key]; fn != nil {
		fn(vp, dr)
	} else {
		panic("unregistered drawable draw callback")
	}
	runtime.KeepAlive(vp)
	runtime.KeepAlive(dr)
}
