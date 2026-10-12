package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"strings"
	"unsafe"
)

type resourceTokenizer struct {
	b  []byte
	at int
}

func (t *resourceTokenizer) next(delim string) (string, bool) {
	for t.at < len(t.b) && t.b[t.at] != 0 && strings.IndexByte(delim, t.b[t.at]) >= 0 {
		t.at++
	}
	start := t.at
	for t.at < len(t.b) && t.b[t.at] != 0 && strings.IndexByte(delim, t.b[t.at]) < 0 {
		t.at++
	}
	if start == t.at {
		return "", false
	}
	out := string(t.b[start:t.at])
	if t.at < len(t.b) && t.b[t.at] != 0 {
		t.b[t.at] = 0
		t.at++
	}
	return out, true
}
func resourceFirstWord(s string) (string, bool) {
	i := 0
	for i < len(s) && resourceSpace(s[i]) {
		i++
	}
	start := i
	for i < len(s) && !resourceSpace(s[i]) {
		i++
	}
	return s[start:i], i > start
}
func resourceWord(data unsafe.Pointer, off uintptr) *uint32 { return (*uint32)(unsafe.Add(data, off)) }
func resourceString(data unsafe.Pointer, off uintptr, s string) {
	copy(unsafe.Slice((*byte)(unsafe.Add(data, off)), len(s)+1), s+"\x00")
}
func resourceParser(kind string, input *byte, data unsafe.Pointer) int {
	s := alloc.GoString(input)
	tokens := resourceTokenizer{b: unsafe.Slice(input, len(s)+1)}
	switch kind {
	case "lifetime", "projectile":
		if v, _, ok := resourceScanInt(s); ok {
			*resourceWord(data, 0) = uint32(v)
		}
		return 1
	case "spark":
		if v, _, ok := resourceScanInt(s); ok {
			*(*byte)(data) = byte(v)
		}
		return 1
	case "mana", "arrow":
		a, ok := tokens.next(" ")
		if !ok {
			return 0
		}
		v := resourceAtoi(a)
		if kind == "mana" {
			*(*byte)(data) = byte(v)
		} else {
			*resourceWord(data, 0) = uint32(v)
			*resourceWord(data, 4) = uint32(v)
		}
		return 1
	case "triple":
		for off := uintptr(0); off < 12; off += 4 {
			v, n, ok := resourceScanInt(s)
			if !ok {
				break
			}
			*resourceWord(data, off) = uint32(v)
			s = s[n:]
		}
		return 1
	case "push":
		if v, n, ok := resourceScanFloat(s, 32); ok {
			*resourceWord(data, 0) = math.Float32bits(float32(v))
			if v, _, ok := resourceScanFloat(s[n:], 32); ok {
				*resourceWord(data, 8) = math.Float32bits(float32(v))
			}
		}
		*resourceWord(data, 4) = *resourceWord(data, 0)
		return 1
	case "skull":
		if name, ok := resourceFirstWord(s); ok {
			resourceString(data, 16, name)
			*resourceWord(data, 12) = 0
			return 1
		}
		return 0
	case "trigger":
		for _, off := range []uintptr{36, 40} {
			if name, ok := tokens.next(" "); ok {
				*resourceWord(data, off) = uint32(sound.ByName(name))
			}
		}
		return 1
	case "audio":
		name, ok := resourceFirstWord(s)
		if !ok {
			return 0
		}
		id := sound.ByName(name)
		*resourceWord(data, 0) = uint32(id)
		return bool2int(id != 0)
	case "spawn":
		name, ok := resourceFirstWord(s)
		if !ok {
			return 0
		}
		resourceString(data, 0, name)
		i := 0
		for i < len(s) && resourceSpace(s[i]) {
			i++
		}
		i += len(name)
		if name, ok := resourceFirstWord(s[i:]); ok {
			*resourceWord(data, 128) = uint32(sound.ByName(name))
			return 1
		}
		return 0
	case "damage":
		v, ok := tokens.next(" ")
		if !ok {
			return 0
		}
		*(*byte)(data) = byte(resourceAtoi(v))
		name, ok := tokens.next(" ")
		if !ok {
			return 0
		}
		id := damageTypeByName(name)
		*resourceWord(data, 4) = uint32(id)
		return bool2int(id != 18)
	case "wand", "wandcast":
		*resourceWord(data, 0) = 0
		if kind == "wandcast" {
			*resourceWord(data, 0) = 1
		}
		charge, ok := tokens.next(" ")
		if !ok {
			return 0
		}
		*(*byte)(unsafe.Add(data, 108)) = byte(resourceAtoi(charge))
		*(*byte)(unsafe.Add(data, 109)) = byte(resourceAtoi(charge))
		*resourceWord(data, 112) = 100
		if kind == "wand" {
			name, ok := tokens.next(" ")
			if !ok {
				return 0
			}
			resourceString(data, 4, name)
			*resourceWord(data, 84) = 0
		}
		rate, ok := tokens.next(" ")
		if !ok {
			return 0
		}
		v := resourceAtof(rate)
		value := float64(GetServer().S().TickRate()) / v
		if kind == "wand" && v == 0 {
			value = 0
		}
		*resourceWord(data, 100) = uint32(effectsTruncWord(value))
		if kind == "wandcast" {
			name, ok := tokens.next(" ")
			if !ok {
				return 0
			}
			*resourceWord(data, 92) = uint32(spell.ParseID(name))
			return 1
		}
		if flag, ok := tokens.next(" "); ok && flag == "MULTI_SHOT" {
			*resourceWord(data, 96) |= 1
		}
		if name, ok := tokens.next(" "); ok {
			*resourceWord(data, 88) = uint32(sound.ByName(name))
		}
		return 1
	}
	panic(kind)
}
func resourceObjectParser(group, kind string) server.ObjectParseFunc {
	return func(typ *server.ObjectType, args []string) error {
		var dst unsafe.Pointer
		switch group {
		case "use":
			dst = typ.UseData.Ptr
		case "update":
			dst = typ.UpdateData
		case "death":
			dst = typ.DeathData
		case "collide":
			dst = typ.CollideData
		}
		raw := append([]byte(strings.Join(args, " ")), 0)
		if resourceParser(kind, &raw[0], dst) == 0 {
			return fmt.Errorf("cannot parse %s data for %q", group, typ.ID())
		}
		return nil
	}
}
