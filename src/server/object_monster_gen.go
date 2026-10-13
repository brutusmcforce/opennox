package server

import "unsafe"

type MonsterGenUpdateData struct {
	Field0     [12]*Object // 0, 0
	Field48    uint32      // 12, 48
	FuncInd52  uint32      // 13, 52
	Field56    uint32      // 14, 56
	FuncInd60  uint32      // 15, 60
	Field64    uint32      // 16, 64
	FuncInd68  uint32      // 17, 68
	Field72    uint32      // 18, 72
	FuncInd76  uint32      // 19, 76
	Rates      [3]byte     // 20, 80: spawn rate by level
	Field83    [3]byte     // 20, 83
	Current    byte        // 21, 86: monsters alive
	Max        byte        // 21, 87
	Frame88    uint32      // 22, 88
	Flags92    uint32      // 23, 92
	HealthPrev [32]uint16  // 24, 96: per player index, the health last reported to that player
	Field160   uint32      // 40, 160
}

// ScriptAt returns the script callback at a 386 offset of the data (48: damaged, 56: dead, 64: spawn, 72).
func (d *MonsterGenUpdateData) ScriptAt(off int) *ScriptCallback {
	switch off {
	case 48:
		return (*ScriptCallback)(unsafe.Pointer(&d.Field48))
	case 56:
		return (*ScriptCallback)(unsafe.Pointer(&d.Field56))
	case 64:
		return (*ScriptCallback)(unsafe.Pointer(&d.Field64))
	case 72:
		return (*ScriptCallback)(unsafe.Pointer(&d.Field72))
	}
	panic("invalid generator script offset")
}
