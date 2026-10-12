package server

import (
	"runtime"
	"unsafe"
)

// DurSpellCallbackFunc preserves the 32-bit result of a duration callback.
// Create and Update use this result; Destroy discards it.
type DurSpellCallbackFunc func(*DurSpell) int32

var durSpellCallbacks = make(map[unsafe.Pointer]DurSpellCallbackFunc)

// RegisterDurSpellCallback binds a stable callback key during initialization.
func RegisterDurSpellCallback(key unsafe.Pointer, fn DurSpellCallbackFunc) {
	durSpellCallbacks[key] = fn
}

func CallDurSpellResult(key unsafe.Pointer, sp *DurSpell) int32 {
	var result int32
	if fn := durSpellCallbacks[key]; fn != nil {
		result = fn(sp)
	} else {
		panic("unregistered duration spell callback")
	}
	runtime.KeepAlive(sp)
	return result
}

func CallDurSpellDiscard(key unsafe.Pointer, sp *DurSpell) {
	if fn := durSpellCallbacks[key]; fn != nil {
		fn(sp)
	} else {
		panic("unregistered duration spell callback")
	}
	runtime.KeepAlive(sp)
}
