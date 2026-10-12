package server

import (
	"runtime"
	"unsafe"
)

// Modifier callback registrations are initialization-only. Keys have stable
// addresses because the original records store a single pointer per callback.
type ModifierEffect3Func func(*ModifierEff, *Object, *Object) int32

var modifierEffects3 = make(map[unsafe.Pointer]ModifierEffect3Func)

func RegisterModifierEffect3(key unsafe.Pointer, fn ModifierEffect3Func) { modifierEffects3[key] = fn }

type ModifierEffect5Func func(*ModifierEff, *Object, *Object, *Object, unsafe.Pointer)

var modifierEffects5 = make(map[unsafe.Pointer]ModifierEffect5Func)

func RegisterModifierEffect5(key unsafe.Pointer, fn ModifierEffect5Func) { modifierEffects5[key] = fn }

type ModifierEffect6Func func(*ModifierEff, *Object, *Object, *Object, *Object, unsafe.Pointer)

var modifierEffects6 = make(map[unsafe.Pointer]ModifierEffect6Func)

func RegisterModifierEffect6(key unsafe.Pointer, fn ModifierEffect6Func) { modifierEffects6[key] = fn }

func CallModifierEffect3Result(key unsafe.Pointer, mod *ModifierEff, a, b *Object) int32 {
	var result int32
	if fn := modifierEffects3[key]; fn != nil {
		result = fn(mod, a, b)
	} else {
		panic("unregistered modifier callback")
	}
	runtime.KeepAlive(mod)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	return result
}

func CallModifierEffect3Discard(key unsafe.Pointer, mod *ModifierEff, a, b *Object) {
	if fn := modifierEffects3[key]; fn != nil {
		fn(mod, a, b)
	} else {
		panic("unregistered modifier callback")
	}
	runtime.KeepAlive(mod)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
}

func CallModifierEffect5(key unsafe.Pointer, mod *ModifierEff, a, b, c *Object, data unsafe.Pointer) {
	if fn := modifierEffects5[key]; fn != nil {
		fn(mod, a, b, c, data)
	} else {
		panic("unregistered modifier callback")
	}
	runtime.KeepAlive(mod)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	runtime.KeepAlive(c)
	runtime.KeepAlive(data)
}

func CallModifierEffect6(key unsafe.Pointer, mod *ModifierEff, a, b, c, d *Object, data unsafe.Pointer) {
	if fn := modifierEffects6[key]; fn != nil {
		fn(mod, a, b, c, d, data)
	} else {
		panic("unregistered modifier callback")
	}
	runtime.KeepAlive(mod)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	runtime.KeepAlive(c)
	runtime.KeepAlive(d)
	runtime.KeepAlive(data)
}
