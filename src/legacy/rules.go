package legacy

import (
	"bytes"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/spell"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/server"
)

// Current section persists across direct line calls and resets on every file
// attempt, including failed opens. The engine invokes this loader serially.
var ruleLoaderContext uint32

func ruleHeader(flags uint16) *int8 {
	mask := uint32(flags) & 0x17f0
	for i := 0; i < 7; i++ {
		if *memmap.PtrUint32(0x587000, uintptr(312212+8*i)) == mask {
			return (*int8)(*memmap.PtrPtr(0x587000, uintptr(312208+8*i)))
		}
	}
	return nil
}

func ruleNarrow(s []uint16) string {
	b := make([]byte, 0, len(s))
	for _, c := range s {
		if byte(c) == 0 {
			break
		}
		b = append(b, byte(c))
	}
	return string(b)
}

func ruleWideBytes(s []byte) []uint16 {
	out := make([]uint16, 0, len(s)+1)
	for _, c := range s {
		if c == 0 {
			break
		}
		out = append(out, uint16(c))
	}
	return append(out, 0)
}

// Keywords compare the original wide characters. Narrowing here would turn
// unrelated Unicode characters into ASCII commands; Unicode folding would also
// accept additional characters that the legacy ASCII keywords do not match.
func ruleEqual(a []uint16, b string) bool {
	for i, c := range a {
		if c == 0 {
			return i == len(b)
		}
		if i >= len(b) {
			return false
		}
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if c != uint16(b[i]) {
			return false
		}
	}
	return len(a) == len(b)
}

func ruleAppendRejected(list *legacyListNode, line []uint16) {
	if list == nil {
		return
	}
	p := legacyCalloc(1, 0x20c)
	if p == nil {
		return
	}
	dst := unsafe.Slice((*uint16)(unsafe.Add(p, 12)), 256)
	copy(dst[:255], line)
	listAppend((*legacyListNode)(unsafe.Pointer(list)), (*legacyListNode)(p))
}

func ruleClearSpell(st *server.Settings2, ind int) {
	st.Field24.Vals[ind/32] &^= uint32(1) << uint(ind%32)
}

// Preserve wcstok cursor behavior: quotes are recognized at the probe after
// the preceding delimiter, not by a general shell/CSV quoting grammar.
func ruleTokens(line []uint16) [][]uint16 {
	buf := append([]uint16(nil), line...)
	if len(buf) == 0 || buf[len(buf)-1] != 0 {
		buf = append(buf, 0)
	}
	next := -1
	stok := func(start int, delims string) (int, bool) {
		if start < 0 {
			start = next
		}
		if start < 0 {
			return 0, false
		}
		for buf[start] != 0 && strings.ContainsRune(delims, rune(buf[start])) {
			start++
		}
		if buf[start] == 0 {
			next = -1
			return 0, false
		}
		for i := start; buf[i] != 0; i++ {
			if strings.ContainsRune(delims, rune(buf[i])) {
				buf[i] = 0
				next = i + 1
				return start, true
			}
		}
		next = -1
		return start, true
	}
	quoted := false
	start := 0
	if buf[0] == '"' {
		start = 1
		quoted = true
	}
	delims := " \n\t\r"
	if quoted {
		delims = "\"\n\r"
	}
	pos, ok := stok(start, delims)
	var out [][]uint16
	for ok {
		end := pos
		for buf[end] != 0 {
			end++
		}
		out = append(out, append([]uint16(nil), buf[pos:end]...))
		probe := end + 1
		if quoted {
			probe++
		}
		if probe < len(buf) && buf[probe] == '"' {
			pos, ok = stok(probe+1, "\"\n\r")
			quoted = true
		} else {
			pos, ok = stok(-1, " \n\t\r")
			quoted = false
		}
	}
	return out
}

func ruleApply(tokens [][]uint16, st *server.Settings2, mode uint32) int {
	if len(tokens) == 0 {
		return 0
	}
	first := ruleNarrow(tokens[0])
	for i := 0; i < 7; i++ {
		name := *memmap.PtrPtr(0x587000, uintptr(312208+8*i))
		if first == GoStringP(unsafe.Pointer((*int8)(name))) {
			ruleLoaderContext = *memmap.PtrUint32(0x587000, uintptr(312212+8*i))
			return bool2int(mode == ruleLoaderContext)
		}
	}
	if ruleLoaderContext&mode == 0 || !ruleEqual(tokens[0], "set") || len(tokens) <= 1 {
		return 0
	}
	if len(tokens) != 4 {
		return 0
	}
	category := tokens[1]
	off := ruleEqual(tokens[3], "off")
	if ruleEqual(category, "armor") {
		if !noxflags.HasGame(noxflags.GameHost) {
			return 0
		}
		typ := GetServer().S().Sub415EC0(ruleNarrow(tokens[2]))
		if typ == nil {
			return 0
		}
		bit := GetServer().S().Armor.Sub_415D10(typ.Ind())
		if bit == 0 {
			return 0
		}
		p := &st.Field48
		if off {
			*p &^= bit
		} else {
			*p |= bit
		}
	} else if ruleEqual(category, "weapon") {
		if !noxflags.HasGame(noxflags.GameHost) {
			return 0
		}
		typ := GetServer().S().Sub415A30(ruleNarrow(tokens[2]))
		if typ == nil {
			return 0
		}
		bit := GetServer().S().Weapons.Nox_xxx_ammoCheck_415880(typ.Ind())
		if bit == 0 {
			return 0
		}
		// The legacy operation edits only the highest nonzero byte of the
		// equipment mask. Shipped entries are all single bits.
		shift := 0
		v := bit
		for v>>8 != 0 {
			v >>= 8
			shift++
		}
		p := &st.Field44[shift]
		if off {
			*p &^= byte(v)
		} else {
			*p |= byte(v)
		}
	} else if ruleEqual(category, "spell") {
		ind := int(spell.ParseID(ruleNarrow(tokens[2])))
		if ind == 0 {
			ind = int(GetServer().S().Spells.ByTitle(GoWStringSlice(tokens[2])))
		}
		if !GetServer().S().Spells.DefByInd(spell.ID(ind)).IsValid() {
			return 0
		}
		if GetServer().S().Spells.Flags(spell.ID(ind))&0x7000000 != 0 && off {
			ruleClearSpell(st, ind)
		}
	} else {
		return 0
	}
	return bool2int(mode == ruleLoaderContext)
}

func ruleParseLine(line []uint16, st *server.Settings2, list *legacyListNode, mode uint32) {
	tokens := ruleTokens(line)
	if len(tokens) != 0 && ruleApply(tokens, st, mode) == 0 {
		ruleAppendRejected(list, line)
	}
}

func ruleReadFile(path string, st *server.Settings2, list *legacyListNode, mode uint32) int {
	ruleLoaderContext = 6128
	f, err := ifs.Open(path)
	if err != nil {
		return 0
	}
	bf := binfile.NewTextFile(f)
	defer bf.Close()
	for {
		line, err := bf.ReadString()
		// The C filesystem bridge consumes a whole physical line and then
		// truncates its copy. Bytes are widened individually, never UTF-8 decoded.
		buf := make([]byte, 256)
		copy(buf[:255], line)
		if i := bytes.IndexByte(buf, '\n'); i >= 0 {
			buf[i] = 0
		}
		if buf[0] != 0 {
			wide := ruleWideBytes(buf)
			ruleParseLine(wide, st, list, mode)
		}
		// EOF may accompany a final nonempty line. Other read errors must
		// also terminate instead of spinning in the old feof-only loop.
		if err != nil {
			break
		}
	}
	return 1
}

func ruleMapName(st *server.Settings2) string {
	b := st.Field0[:8]
	if i := bytes.IndexByte(b, 0); i >= 0 {
		b = b[:i]
	}
	return string(b)
}

func ruleLoad(st *server.Settings2, user string, list *legacyListNode, selection byte, flags uint16) byte {
	if list != nil {
		runtimeRejectedClear(unsafe.Pointer(list))
	}
	for i := range st.Field24.Vals {
		st.Field24.Vals[i] = ^uint32(0)
	}
	for i := range st.Field44 {
		st.Field44[i] = 255
	}
	st.Field48 = ^uint32(0)
	mode := uint32(flags) & 0x17f0
	var custom int
	if selection&3 != 0 {
		base := "maps\\" + ruleMapName(st) + "\\"
		if selection&2 != 0 {
			custom = ruleReadFile(base+user, st, list, mode)
		}
		if selection&1 != 0 && custom == 0 {
			ruleReadFile(base+ruleMapName(st)+".rul", st, list, mode)
		}
	}
	if Get_dword_5d4594_2650652() != 0 && selection&4 != 0 {
		ruleReadFile("internet.rul", st, list, mode)
	}
	if flags&0x40 != 0 {
		ruleClearSpell(st, 132)
		return 0xef // Low byte of the original complement mask returned by 453FA0.
	}
	return byte(flags)
}

func ruleCString(s string) string {
	if i := strings.IndexByte(s, 0); i >= 0 {
		return s[:i]
	}
	return s
}
