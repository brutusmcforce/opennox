package gui

import "unsafe"

type TooltipCallbackGo func(*Window, *WindowData, uintptr)

var tooltipCallbacksGo = make(map[unsafe.Pointer]TooltipCallbackGo)

// RegisterTooltipCallbackGo binds a stable identity to a native tooltip callback.
// Registration is intended for initialization; foreign function pointers remain
// supported by Window.TooltipFunc's raw C fallback.
func RegisterTooltipCallbackGo(key unsafe.Pointer, fn TooltipCallbackGo) {
	if key == nil || fn == nil {
		panic("invalid tooltip callback")
	}
	if _, ok := tooltipCallbacksGo[key]; ok {
		panic("tooltip callback already registered")
	}
	tooltipCallbacksGo[key] = fn
}
