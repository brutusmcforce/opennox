package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/noxscript/ns/asm"

	"github.com/brutusmcforce/opennox/server/noxscript"
)

var (
	Nox_script_shouldReadMoreXxx     func(fi asm.Builtin) bool
	Nox_script_shouldReadEvenMoreXxx func(fi asm.Builtin) bool
)

func CallScriptBuiltin(fi asm.Builtin) (int, bool) {
	if fi < 0 || int(fi) >= len(noxScriptBuiltins) {
		return 0, false
	}
	fnc := noxScriptBuiltins[fi]
	if fnc == nil {
		return 0, false
	}
	res := fnc(GetServer().NoxScriptC())
	return res, true
}

func Nox_script_StartupScreen_516600_A() {
	scriptInventoryStartup()
}

func Sub_512E80(str string) int {
	cstr, _ := CWString(str)
	return scriptBindingIntern(unsafe.Pointer(cstr))
}

var noxScriptBuiltins = [asm.BuiltinGetScore + 1]noxscript.Builtin{
	asm.BuiltinSetQuestStatus:      func(vm noxscript.VM) int { name := vm.PopString(); questProgressSet(name, vm.PopU32(), 0); return 0 },
	asm.BuiltinSetQuestStatusFloat: func(vm noxscript.VM) int { name := vm.PopString(); questProgressSet(name, vm.PopU32(), 1); return 0 },
	asm.BuiltinGetQuestStatus:      func(vm noxscript.VM) int { vm.PushU32(questProgressInt(vm.PopString())); return 0 },
	asm.BuiltinGetQuestStatusFloat: func(vm noxscript.VM) int { vm.PushF32(float32(questProgressFloat(vm.PopString()))); return 0 },
	asm.BuiltinResetQuestStatus:    func(vm noxscript.VM) int { questProgressReset(vm.PopString()); return 0 },
	asm.BuiltinSetRoamFlag:         func(vm noxscript.VM) int { return scriptBindingRoamByte(vm, false) },
	asm.BuiltinGroupSetRoamFlag:    func(vm noxscript.VM) int { return scriptBindingRoamByte(vm, true) },
	asm.BuiltinJournalDelete:       func(vm noxscript.VM) int { return scriptBindingJournal(vm, false) },
	asm.BuiltinJournalEdit:         func(vm noxscript.VM) int { return scriptBindingJournal(vm, true) },
	asm.BuiltinGiveXp:              func(vm noxscript.VM) int { return scriptBindingGiveXP(vm) },
	asm.BuiltinIsTalking:           func(vm noxscript.VM) int { return scriptBindingHostState(vm, false) },
	asm.BuiltinMakeFriendly:        func(vm noxscript.VM) int { return scriptBindingOwnership(vm, 0) },
	asm.BuiltinMakeEnemy:           func(vm noxscript.VM) int { return scriptBindingOwnership(vm, 1) },
	asm.BuiltinBecomePet:           func(vm noxscript.VM) int { return scriptBindingOwnership(vm, 2) },
	asm.BuiltinBecomeEnemy:         func(vm noxscript.VM) int { return scriptBindingOwnership(vm, 3) },
	asm.BuiltinUnknownb8:           func(vm noxscript.VM) int { return scriptBindingSubclass(vm, 0x100) },
	asm.BuiltinUnknownb9:           func(vm noxscript.VM) int { return scriptBindingSubclass(vm, 0x80) },
	asm.BuiltinSetHalberd:          scriptInventoryHalberd,
	asm.BuiltinIsTrading:           func(vm noxscript.VM) int { return scriptBindingHostState(vm, true) },
}
