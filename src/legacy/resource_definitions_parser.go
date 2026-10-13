package legacy

import (
	"fmt"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
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
func resourceWord(data unsafe.Pointer, off uintptr) *uint32 { return rawAt[uint32](data, int(off)) }
func resourceString(data unsafe.Pointer, off uintptr, s string) {
	copy(unsafe.Slice((*byte)(unsafe.Add(data, off)), len(s)+1), s+"\x00")
}
func resourceParser(kind string, input *byte, data unsafe.Pointer) int {
	s := alloc.GoString(input)
	tokens := resourceTokenizer{b: unsafe.Slice(input, len(s)+1)}
	switch kind {
	case "lifetime":
		if v, _, ok := resourceScanInt(s); ok {
			(*server.LifetimeUpdateData)(data).Lifetime = uint32(v)
		}
		return 1
	case "projectile":
		if v, _, ok := resourceScanInt(s); ok {
			(*server.ProjectileCollideData)(data).Damage = uint32(v)
		}
		return 1
	case "spark":
		if v, _, ok := resourceScanInt(s); ok {
			(*server.SparkExplosionCollideData)(data).Damage = byte(v)
		}
		return 1
	case "mana", "arrow":
		a, ok := tokens.next(" ")
		if !ok {
			return 0
		}
		v := resourceAtoi(a)
		if kind == "mana" {
			(*server.ManaDrainCollideData)(data).Mana = byte(v)
		} else {
			d := (*server.MonsterArrowCollideData)(data)
			d.DamageCoop = uint32(v)
			d.Damage = uint32(v)
		}
		return 1
	case "triple":
		d := (*server.LoopAndDamageUpdateData)(data)
		for _, dst := range []*uint32{&d.Field0, &d.Field4, &d.Field8} {
			v, n, ok := resourceScanInt(s)
			if !ok {
				break
			}
			*dst = uint32(v)
			s = s[n:]
		}
		return 1
	case "push":
		d := (*server.PushUpdateData)(data)
		if v, n, ok := resourceScanFloat(s, 32); ok {
			d.Radius = float32(v)
			if v, _, ok := resourceScanFloat(s[n:], 32); ok {
				d.Power = float32(v)
			}
		}
		d.Radius2 = d.Radius
		return 1
	case "skull":
		if name, ok := resourceFirstWord(s); ok {
			d := (*server.SkullUpdateData)(data)
			resourceString(unsafe.Pointer(&d.MissileName[0]), 0, name)
			d.MissileType = 0
			return 1
		}
		return 0
	case "trigger":
		d := (*server.TriggerUpdateData)(data)
		for _, dst := range []*uint32{&d.SoundOn, &d.SoundOff} {
			if name, ok := tokens.next(" "); ok {
				*dst = uint32(sound.ByName(name))
			}
		}
		return 1
	case "audio":
		name, ok := resourceFirstWord(s)
		if !ok {
			return 0
		}
		id := sound.ByName(name)
		(*server.AudioEventCollideData)(data).Sound = uint32(id)
		return bool2int(id != 0)
	case "spawn":
		name, ok := resourceFirstWord(s)
		if !ok {
			return 0
		}
		d := (*server.CreateObjectDeathData)(data)
		resourceString(unsafe.Pointer(&d.Name[0]), 0, name)
		i := 0
		for i < len(s) && resourceSpace(s[i]) {
			i++
		}
		i += len(name)
		if name, ok := resourceFirstWord(s[i:]); ok {
			d.Sound = uint32(sound.ByName(name))
			return 1
		}
		return 0
	case "damage":
		v, ok := tokens.next(" ")
		if !ok {
			return 0
		}
		d := (*server.DamageCollideData)(data)
		d.Damage = byte(resourceAtoi(v))
		name, ok := tokens.next(" ")
		if !ok {
			return 0
		}
		id := damageTypeByName(name)
		d.DamageType = uint32(id)
		return bool2int(id != 18)
	case "wand", "wandcast":
		d := (*server.WandUseData)(data)
		d.Kind = 0
		if kind == "wandcast" {
			d.Kind = 1
		}
		charge, ok := tokens.next(" ")
		if !ok {
			return 0
		}
		d.Charges = byte(resourceAtoi(charge))
		d.MaxCharges = byte(resourceAtoi(charge))
		d.Percent = 100
		if kind == "wand" {
			name, ok := tokens.next(" ")
			if !ok {
				return 0
			}
			resourceString(unsafe.Pointer(&d.MissileBuf[0]), 0, name)
			d.MissileInd = 0
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
		d.Delay = uint32(effectsTruncWord(value))
		if kind == "wandcast" {
			name, ok := tokens.next(" ")
			if !ok {
				return 0
			}
			d.Spell = uint32(spell.ParseID(name))
			return 1
		}
		if flag, ok := tokens.next(" "); ok && flag == "MULTI_SHOT" {
			d.Flags |= 1
		}
		if name, ok := tokens.next(" "); ok {
			d.Sound = uint32(sound.ByName(name))
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
