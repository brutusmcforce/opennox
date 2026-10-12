package legacy

import (
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
)

func mapThemeToken() *byte { return memmap.PtrUint8(0x5D4594, 2487264) }
func mapThemeText() string { return populationString(mapRoomRaw(unsafe.Pointer(mapThemeToken()))) }
func mapThemeLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}
func mapThemeByte(p uint32, off int) *byte { return (*byte)(unsafe.Add(mapRoomPointer(p), off)) }
func mapThemeCopy(p uint32, s string) {
	b := unsafe.Slice(mapThemeByte(p, 0), len(s)+1)
	copy(b, s)
	b[len(s)] = 0
}
func mapThemeRaw(f uint32, out *byte) uint32 {
	dst := out
	var previous byte
	leading := true
	for {
		c, ok := mapThemeReadByte(f)
		if !ok {
			return 0
		}
		if c == '\n' {
			*populationBlob(2487520)++
		}
		space := c == ' ' || (c >= '\t' && c <= '\r')
		if space {
			if !leading {
				*dst = 0
				return 1
			}
		} else {
			leading = false
			if c == '/' && previous == '/' {
				mapThemeSkipLine(f)
				dst = out
				leading = true
			} else {
				*dst = c
				dst = (*byte)(unsafe.Add(unsafe.Pointer(dst), 1))
			}
		}
		previous = c
	}
}
func mapThemeSkipLine(f uint32) uint32 {
	for {
		c, ok := mapThemeReadByte(f)
		if !ok {
			return ^uint32(0)
		}
		if c == '\n' {
			*populationBlob(2487520)++
			return 0
		}
	}
}
func mapThemeRead(f uint32, out *byte) uint32 {
	if mapThemeRaw(f, out) == 0 {
		return 0
	}
	return mapThemeControl(f)
}
func mapThemeNext(f uint32) bool { return mapThemeRead(f, mapThemeToken()) != 0 }
func mapThemeSkip(f uint32, allowElse bool) uint32 {
	for mapThemeRaw(f, mapThemeToken()) != 0 {
		switch mapThemeLower(mapThemeText()) {
		case "if":
			if mapThemeSkip(f, false) == 0 {
				return 0
			}
		case "endif":
			return 1
		case "else":
			if allowElse {
				return 1
			}
		}
	}
	return 0
}
func mapThemeControl(f uint32) uint32 {
	switch mapThemeLower(mapThemeText()) {
	case "if":
		var yes uint32
		if mapThemeCondition(f, &yes) == 0 {
			return 0
		}
		if yes == 0 && mapThemeSkip(f, true) == 0 {
			return 0
		}
	case "else":
		if mapThemeSkip(f, false) == 0 {
			return 0
		}
	case "endif":
	default:
		return 1
	}
	return mapThemeRead(f, mapThemeToken())
}
func mapThemeOperator(f uint32, out *uint32) uint32 {
	if !mapThemeNext(f) {
		return 0
	}
	for i, s := range []string{"<", ">", "==", "!="} {
		if mapThemeText() == s {
			*out = uint32(i)
			return 1
		}
	}
	return 0
}
func mapThemeCondition(f uint32, out *uint32) uint32 {
	if !mapThemeNext(f) {
		return 0
	}
	s := mapThemeLower(mapThemeText())
	class := -1
	switch s {
	case "warrior":
		class = 0
	case "wizard":
		class = 1
	case "conjurer":
		class = 2
	}
	if class >= 0 {
		*out = 0
		for p := GetServer().S().Players.First(); p != nil; p = GetServer().S().Players.Next(p) {
			if int(*(*byte)(unsafe.Add(p.C(), 2251))) == class {
				*out = 1
				break
			}
		}
		return 1
	}
	if s == "experience_level" || s == "numplayers" {
		count := int32(3)
		if s == "numplayers" {
			count = int32(GetServer().S().Players.Count())
		}
		var op uint32
		if mapThemeOperator(f, &op) == 0 || !mapThemeNext(f) {
			return 0
		}
		value := mapThemeInt(mapThemeText())
		yes := false
		switch op {
		case 0:
			yes = count < value
		case 1:
			yes = count > value
		case 2:
			yes = count == value
		case 3:
			yes = count != value
		}
		*out = 0
		if yes {
			*out = 1
		}
		return 1
	}
	if strings.Contains(s, "%") {
		*out = 0
		if mapRoomRandomInt(1, 100) <= mapThemeInt(mapThemeText()) {
			*out = 1
		}
		return 1
	}
	return 0
}
