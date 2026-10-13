package legacy

import (
	"fmt"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/server"
)

func ruleWrite(name string, st *server.Settings2, list *legacyListNode) byte {
	if st.Field52&0x80 != 0 {
		return byte(st.Field52)
	}
	f, err := ifs.Create("maps\\" + ruleMapName(st) + "\\" + name)
	if err != nil {
		return 0
	}
	defer f.Close()
	online := Get_dword_5d4594_2650652() != 0
	// The decompiled C split each of these into separate 24/36-byte arrays.
	// GCC reordered those arrays, overlapping the loads and leaving a mask
	// uninitialized. Each load now owns a complete, independent settings value.
	var internet, local server.Settings2
	if online {
		// Only the name is input to ruleLoad; the masks are reset there. Avoid
		// copying caller padding or trailing fields from legacy-sized buffers.
		internet.Field0 = st.Field0
		local.Field0 = st.Field0
		ruleLoad(&internet, "user.rul", nil, 4, st.Field52)
		ruleLoad(&local, "user.rul", nil, 3, st.Field52)
	}
	if list != nil {
		for p := list.next; p != list; p = p.next {
			var narrow []byte
			for _, c := range (*ruleRejected)(unsafe.Pointer(p)).Line {
				if c == 0 {
					break
				}
				narrow = append(narrow, byte(c))
			}
			narrow = append(narrow, '\n')
			// A NUL introduced by narrowing terminates fputs before even the LF.
			_, _ = f.WriteString(ruleCString(string(narrow)))
		}
	}
	_, _ = f.WriteString(GoStringP(unsafe.Pointer(ruleHeader(st.Field52))) + "\n")
	s := GetServer().S()
	for i := 1; i <= 136; i++ {
		id := spell.ID(i)
		mask := uint32(1) << uint(i%32)
		if !s.Spells.DefByInd(id).IsValid() || st.Field24.Vals[i/32]&mask != 0 || s.Spells.Flags(id)&0x7000000 == 0 {
			continue
		}
		if online && internet.Field24.Vals[i/32]&mask == 0 && local.Field24.Vals[i/32]&mask != 0 {
			continue
		}
		_, _ = fmt.Fprintf(f, "set spell \"%s\" off\n", id.String())
	}
	for i := 0; i < 26; i++ {
		bit := uint32(1) << uint(i)
		if st.Field48&bit == 0 {
			_, _ = fmt.Fprintf(f, "set Armor \"%s\" off\n", s.Armor.Sub_415E40(bit))
		}
	}
	for i := 0; i < 27; i++ {
		if st.Field44[i/8]&(1<<uint(i%8)) == 0 {
			_, _ = fmt.Fprintf(f, "set weapon \"%s\" off\n", s.Weapons.Sub_4159B0(1<<uint(i)))
		}
	}
	// The legacy ABI reports zero for both successful writing and failed creation.
	return 0
}
