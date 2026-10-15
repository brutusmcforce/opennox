package server

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/noxscript/ns/asm"
)

// A map script that loops forever is stopped instead of freezing the server.
func TestScriptEndlessLoop(t *testing.T) {
	var vm NoxScriptVM
	loop := ScriptFunc{FuncDef: asm.FuncDef{Name: "loop", Code: []uint32{uint32(asm.OpJump), 0}}}
	vm.vm.funcs = []ScriptFunc{{}, {}, loop}
	err := vm.callByFunc(&vm.vm.funcs[2], nil, nil, nil)
	must.ErrorContains(t, err, "script runs too long")
	must.EqOp(t, 0, len(vm.vm.stack))
}
