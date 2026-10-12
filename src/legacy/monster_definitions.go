package legacy

import (
	"errors"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"strings"
	"unsafe"
)

var monsterDefinitions *server.MonsterDef
var errMonsterDefinitionToken = errors.New("monster definition token exceeds its field")

func monsterDefinitionReadByte(f *binfile.Binfile) (byte, bool) {
	var b [1]byte
	_, err := f.Read(b[:])
	if err != nil {
		f.File.Err = err
	}
	return b[0], f.File.Err == nil
}
func monsterDefinitionSkipLine(f *binfile.Binfile) int {
	for {
		c, ok := monsterDefinitionReadByte(f)
		if !ok {
			return -1
		}
		if c == '\n' {
			return 0
		}
	}
}
func monsterDefinitionToken(f *binfile.Binfile, out []byte) bool {
	n := 0
	leading := true
	var previous byte
	for {
		c, ok := monsterDefinitionReadByte(f)
		if !ok {
			return false
		}
		space := c == ' ' || (c >= '\t' && c <= '\r')
		if space {
			if !leading {
				out[n] = 0
				return true
			}
		} else {
			leading = false
			if c == '/' && previous == '/' {
				monsterDefinitionSkipLine(f)
				n = 0
				leading = true
			} else {
				if n >= len(out)-1 {
					f.File.Err = errMonsterDefinitionToken
					return false
				}
				out[n] = c
				n++
			}
		}
		previous = c
	}
}
func monsterDefinitionText(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
func monsterDefinitionStatus(name string) uint16 {
	if strings.HasPrefix(name, "NULL") {
		return 0
	}
	var bits uint32
	for _, part := range strings.Split(name, "+") {
		if part == "" {
			continue
		}
		for i := uintptr(0); ; i++ {
			p := memmap.LoadPtr(0x587000, 247536+4*i)
			if p == nil {
				break
			}
			if mapThemeLower(alloc.GoString((*byte)(p))) == mapThemeLower(part) {
				bits |= uint32(1) << uint(i)
				break
			}
		}
	}
	return uint16(bits)
}
func monsterDefinitionParse(f *binfile.Binfile, name string) bool {
	if len(name) >= 64 {
		return false
	}
	d, _ := alloc.New(server.MonsterDef{})
	*d = server.MonsterDef{}
	copy(d.Name0[:], name)
	accepted := false
	defer func() {
		if !accepted {
			alloc.Free(d)
		}
	}()
	var token [256]byte
	for {
		ok := monsterDefinitionToken(f, token[:])
		key := monsterDefinitionText(token[:])
		if errors.Is(f.File.Err, errMonsterDefinitionToken) {
			return false
		}
		if !ok || mapThemeLower(key) == "end" {
			d.Next244 = monsterDefinitions
			monsterDefinitions = d
			accepted = true
			return true
		}
		if noxflags.HasGame(2048 | 0x200000) {
			switch mapThemeLower(key) {
			case "arena":
				monsterDefinitionSkipLine(f)
				continue
			case "solo":
				continue
			}
		}
		if noxflags.HasGame(0x2000) {
			switch mapThemeLower(key) {
			case "solo":
				monsterDefinitionSkipLine(f)
				continue
			case "arena":
				continue
			}
		}
		field := uintptr(248192)
		for {
			p := memmap.LoadPtr(0x587000, field)
			if p == nil {
				return false
			}
			if mapThemeLower(alloc.GoString((*byte)(p))) == mapThemeLower(key) {
				break
			}
			field += 12
		}
		kind := memmap.Uint32(0x587000, field+4)
		if kind > 8 {
			continue
		}
		offset := memmap.Uint32(0x587000, field+8)
		dst := unsafe.Add(unsafe.Pointer(d), uintptr(offset))
		if kind == 7 {
			if !monsterDefinitionToken(f, unsafe.Slice((*byte)(dst), 64)) && errors.Is(f.File.Err, errMonsterDefinitionToken) {
				return false
			}
			if alloc.GoString((*byte)(dst)) == "NULL" {
				*(*byte)(dst) = 0
			}
			continue
		}
		// Preserve the original token buffer across a short/EOF value read.
		monsterDefinitionToken(f, token[:])
		if errors.Is(f.File.Err, errMonsterDefinitionToken) {
			return false
		}
		value := monsterDefinitionText(token[:])
		switch kind {
		case 0:
			*(*uint32)(dst) = uint32(mapThemeInt(value))
		case 1:
			*(*float32)(dst) = float32(mapThemeFloat(value))
		case 2:
			*(*uint32)(dst) = uint32(sound.ByName(value))
		case 3, 4, 5:
			table := map[uint32]uintptr{3: 287096, 4: 287280, 5: 287192}[kind]
			if !monsterLoadCallback(unsafe.Pointer(d), value, table, uintptr(offset)) {
				return false
			}
		case 6:
			*(*uint16)(dst) = monsterDefinitionStatus(value)
		case 8:
			d.MeleeAttackDamageType124 = 18
			for i := uintptr(0); i < 18; i++ {
				name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, 247464+4*i)))
				if mapThemeLower(value) == mapThemeLower(name[7:]) {
					d.MeleeAttackDamageType124 = uint32(i)
					break
				}
			}
			if d.MeleeAttackDamageType124 == 18 {
				return false
			}
		}
	}
}
func monsterDefinitionLoad() int {
	monsterDefinitions = nil
	f, err := binfile.BinfileOpen("monster.bin", binfile.ReadOnly)
	if err != nil {
		return 0
	}
	defer f.Close()
	if err = f.SetKey(23); err != nil {
		return 0
	}
	var name [256]byte
	for monsterDefinitionToken(f, name[:]) && monsterDefinitionParse(f, monsterDefinitionText(name[:])) {
	}
	return 1
}
func monsterDefinitionFree() uint32 {
	for p := monsterDefinitions; p != nil; {
		next := p.Next244
		alloc.Free(p)
		p = next
	}
	monsterDefinitions = nil
	return 0
}
func monsterDefinitionBind() int {
	for p := monsterDefinitions; p != nil; p = p.Next244 {
		p.TypeInd240 = uint32(GetServer().S().Types.IndByID(monsterDefinitionText(p.Name0[:])))
		if p.TypeInd240 == 0 {
			monsterDefinitionFree()
			return 0
		}
	}
	return 1
}
func monsterDefinitionByType(id uint32) *server.MonsterDef {
	for p := monsterDefinitions; p != nil; p = p.Next244 {
		if p.TypeInd240 == id {
			return p
		}
	}
	return nil
}
