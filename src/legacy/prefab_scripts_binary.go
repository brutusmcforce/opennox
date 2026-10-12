package legacy

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/brutusmcforce/opennox/noxscript/ns/asm"
	"github.com/brutusmcforce/opennox/internal/binfile"
)

var prefabScriptOffsets [3]uint32 // globals, strings, functions in the first input

func prefabScriptReadInt(f *binfile.File) uint32 {
	var b [4]byte
	_, _ = f.Read(b[:])
	return binary.LittleEndian.Uint32(b[:])
}
func prefabScriptReadFloat(f *binfile.File) float64 {
	return float64(math.Float32frombits(prefabScriptReadInt(f)))
}
func prefabScriptWriteInt(f *binfile.File, v uint32) uint32 {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	n, _ := f.Write(b[:])
	return uint32(n / 4)
}
func prefabScriptWriteFloat(f *binfile.File, v float32) uint32 {
	return prefabScriptWriteInt(f, math.Float32bits(v))
}
func prefabScriptInstructions(in, out *binfile.File, relocate bool) uint32 {
	for {
		op := prefabScriptReadInt(in)
		// An incomplete stream has no terminating instruction; stop instead of looping.
		if in.Err != nil {
			return 0
		}
		if op > 73 {
			return op
		}
		switch op {
		case 0, 1, 2:
			scope, v := prefabScriptReadInt(in), prefabScriptReadInt(in)
			if scope != 0 && relocate && int32(v) >= 4 {
				v += prefabScriptOffsets[0] - 4
			}
			prefabScriptWriteInt(out, op)
			prefabScriptWriteInt(out, scope)
			prefabScriptWriteInt(out, v)
		case 3:
			scope, v := prefabScriptReadInt(in), prefabScriptReadInt(in)
			if scope != 0 && relocate {
				v += prefabScriptOffsets[1]
			}
			prefabScriptWriteInt(out, op)
			prefabScriptWriteInt(out, scope)
			prefabScriptWriteInt(out, v)
		case 4, 6, 19, 20, 21, 70:
			v := prefabScriptReadInt(in)
			if relocate {
				if op == 6 {
					v += prefabScriptOffsets[1]
				} else if op == 70 {
					v += prefabScriptOffsets[2] - 2
				}
			}
			prefabScriptWriteInt(out, op)
			prefabScriptWriteInt(out, v)
		case 5:
			v := float32(prefabScriptReadFloat(in))
			prefabScriptWriteInt(out, op)
			prefabScriptWriteFloat(out, v)
		case 69:
			builtin := prefabScriptReadInt(in)
			if relocate {
				if Nox_script_shouldReadMoreXxx(asm.Builtin(builtin)) {
					out.Seek(-4, io.SeekCurrent)
					v := prefabScriptReadInt(out) + prefabScriptOffsets[2] - 2
					out.Seek(-4, io.SeekCurrent)
					prefabScriptWriteInt(out, v)
					if Nox_script_shouldReadEvenMoreXxx(asm.Builtin(builtin)) {
						out.Seek(-12, io.SeekCurrent)
						v = prefabScriptReadInt(out) + prefabScriptOffsets[2] - 2
						out.Seek(-4, io.SeekCurrent)
						prefabScriptWriteInt(out, v)
						out.Seek(8, io.SeekCurrent)
					}
				}
			}
			prefabScriptWriteInt(out, op)
			prefabScriptWriteInt(out, builtin)
		case 72:
			return prefabScriptWriteInt(out, op)
		default:
			prefabScriptWriteInt(out, op)
		}
		if in.Err != nil || out.Err != nil {
			return 0
		}
	}
}

func prefabScriptCopyBytes(in, out *binfile.File, n uint32) bool {
	// C's largest temporary has 4096 bytes; larger records were undefined.
	if n > 4096 {
		return false
	}
	data := make([]byte, int(n))
	nr, err := in.Read(data)
	if err != nil || nr != len(data) {
		return false
	}
	nw, err := out.Write(data)
	return err == nil && nw == len(data)
}
func prefabScriptCopyWord(in, out *binfile.File) uint32 {
	v := prefabScriptReadInt(in)
	prefabScriptWriteInt(out, v)
	return v
}
func prefabScriptCopyString(in, out *binfile.File) bool {
	n := prefabScriptCopyWord(in, out)
	return prefabScriptCopyBytes(in, out, n)
}
func prefabScriptStrings(a, b, out *binfile.File) uint32 {
	prefabScriptReadInt(a)
	prefabScriptCopyWord(b, out)
	na, nb := prefabScriptReadInt(a), prefabScriptReadInt(b)
	prefabScriptOffsets[1] = na
	prefabScriptWriteInt(out, na+nb)
	for _, v := range []struct {
		f *binfile.File
		n uint32
	}{{a, na}, {b, nb}} {
		for i := int32(0); i < int32(v.n); i++ {
			if !prefabScriptCopyString(v.f, out) {
				return 0
			}
		}
	}
	return nb
}
func prefabScriptLocals(in, out *binfile.File) uint32 {
	prefabScriptCopyWord(in, out)
	n := prefabScriptCopyWord(in, out)
	prefabScriptCopyWord(in, out)
	for i := int32(0); i < int32(n); i++ {
		prefabScriptCopyWord(in, out)
	}
	return n
}
func prefabScriptReserved(a, b, out *binfile.File) uint32 {
	both := func() uint32 { prefabScriptReadInt(a); return prefabScriptCopyWord(b, out) }
	names := func() bool {
		n := prefabScriptReadInt(a)
		if n >= 4096 {
			return false
		}
		if _, err := a.Seek(int64(n), io.SeekCurrent); err != nil {
			return false
		}
		return prefabScriptCopyString(b, out)
	}
	both()
	if !names() {
		return 0
	}
	both()
	both()
	both()
	na, nb := prefabScriptReadInt(a), prefabScriptReadInt(b)
	prefabScriptWriteInt(out, na+nb)
	both()
	both()
	sizeA := prefabScriptReadInt(a)
	prefabScriptReadInt(b)
	prefabScriptWriteInt(out, sizeA)
	both()
	both()
	if !names() {
		return 0
	}
	both()
	both()
	both()
	na, nb = prefabScriptReadInt(a), prefabScriptReadInt(b)
	prefabScriptOffsets[0] = na
	prefabScriptWriteInt(out, na+nb-4)
	both()
	for i := int32(0); i < int32(na); i++ {
		prefabScriptCopyWord(a, out)
	}
	for i := 0; i < 4; i++ {
		prefabScriptReadInt(b)
	}
	for i := int32(4); i < int32(nb); i++ {
		prefabScriptCopyWord(b, out)
	}
	both()
	sizeA = prefabScriptReadInt(a)
	sizeB := prefabScriptReadInt(b)
	prefabScriptWriteInt(out, sizeA+sizeB-4)
	prefabScriptInstructions(a, out, false)
	out.Seek(-4, io.SeekCurrent)
	return prefabScriptInstructions(b, out, true)
}
func prefabScriptFunctions(a, b, out *binfile.File) uint32 {
	prefabScriptReadInt(a)
	prefabScriptCopyWord(b, out)
	na, nb := prefabScriptReadInt(a), prefabScriptReadInt(b)
	prefabScriptOffsets[2] = na
	prefabScriptWriteInt(out, na+nb-2)
	result := prefabScriptReserved(a, b, out)
	for _, v := range []struct {
		f        *binfile.File
		n        uint32
		relocate bool
	}{{a, na, false}, {b, nb, true}} {
		for i := int32(0); i < int32(v.n)-2; i++ {
			prefabScriptCopyWord(v.f, out)
			if !prefabScriptCopyString(v.f, out) {
				return 0
			}
			prefabScriptCopyWord(v.f, out)
			prefabScriptCopyWord(v.f, out)
			prefabScriptLocals(v.f, out)
			prefabScriptCopyWord(v.f, out)
			prefabScriptCopyWord(v.f, out)
			prefabScriptInstructions(v.f, out, v.relocate)
			result = v.n - 2
		}
	}
	return result
}
func prefabScriptMerge(a, b, out *binfile.File) uint32 {
	prefabScriptOffsets = [3]uint32{}
	a.Seek(8, io.SeekCurrent)
	prefabScriptCopyBytes(b, out, 8)
	prefabScriptStrings(a, b, out)
	prefabScriptFunctions(a, b, out)
	prefabScriptReadInt(a)
	var end [4]byte
	b.Read(end[:])
	n, _ := out.Write(end[:])
	return uint32(n)
}
