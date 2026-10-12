package gui

import (
	"unsafe"
)

var animationCallbacksGo = make(map[unsafe.Pointer]func() int)

// RegisterAnimationCallbackGo binds a stable identity during initialization.
func RegisterAnimationCallbackGo(key unsafe.Pointer, fn func() int) {
	if key == nil || fn == nil {
		panic("invalid animation callback")
	}
	if _, ok := animationCallbacksGo[key]; ok {
		panic("animation callback already registered")
	}
	animationCallbacksGo[key] = fn
}

// CallAnimationCallback invokes a registered native animation identity.
func CallAnimationCallback(key unsafe.Pointer) int {
	if fn := animationCallbacksGo[key]; fn != nil {
		return fn()
	}
	panic("unregistered animation callback")
}
