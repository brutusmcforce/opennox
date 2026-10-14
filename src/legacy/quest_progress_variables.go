package legacy

import (
	"math"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
)

// The list is private to quest variables; C callers only request reset/save/load.
// Preserve the 148-byte record layout for persistence and qualified snapshots.
type questProgressRecord struct {
	Name        [132]byte
	Kind, Value uint32
	Next, Prev  *questProgressRecord
}

var questProgressHead *questProgressRecord

func questProgressCString(s string) string {
	if i := strings.IndexByte(s, 0); i >= 0 {
		return s[:i]
	}
	return s
}
func questProgressNamespace(name string) {
	name = questProgressCString(name)
	if len(name) >= 132 {
		return
	}
	dst := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1570008), 132)
	copy(dst, name)
	dst[len(name)] = 0
}
func questProgressQualify(name string) (string, bool) {
	name = questProgressCString(name)
	if !strings.Contains(name, ":") {
		base := monsterDefinitionText(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1570008), 132))
		name = base + ":" + name
	}
	if len(name) >= 132 {
		return "", false
	}
	dst := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1570140), 132)
	copy(dst, name)
	dst[len(name)] = 0
	return name, true
}
func questProgressFind(name string) *questProgressRecord {
	key, ok := questProgressQualify(name)
	if !ok {
		return nil
	}
	for p := questProgressHead; p != nil; p = p.Next {
		if mapThemeLower(monsterDefinitionText(p.Name[:])) == mapThemeLower(key) {
			return p
		}
	}
	return nil
}
func questProgressSet(name string, value, kind uint32) *questProgressRecord {
	if p := questProgressFind(name); p != nil {
		p.Value = value
		return p
	}
	key, ok := questProgressQualify(name)
	if !ok {
		return nil
	}
	old := questProgressHead
	p := &questProgressRecord{Kind: kind, Value: value, Next: old}
	copy(p.Name[:], key)
	if old != nil {
		old.Prev = p
	}
	questProgressHead = p
	return old
}
func questProgressInt(name string) uint32 {
	if p := questProgressFind(name); p != nil {
		return p.Value
	}
	return 0
}
func questProgressFloat(name string) float64 {
	return float64(math.Float32frombits(questProgressInt(name)))
}
func questProgressRemove(p *questProgressRecord) {
	if p.Prev != nil {
		p.Prev.Next = p.Next
	}
	if p.Next != nil {
		p.Next.Prev = p.Prev
	}
	if p == questProgressHead {
		questProgressHead = p.Next
	}
}
func questProgressReset(name string) {
	pattern, ok := questProgressQualify(name)
	if !ok {
		return
	}
	star := strings.IndexByte(pattern, '*')
	if star < 0 {
		if p := questProgressFind(name); p != nil {
			questProgressRemove(p)
		}
		return
	}
	if pattern == "*:*" {
		questProgressHead = nil
		return
	}
	// The legacy suffix branches use the FIRST substring occurrence, not globbing.
	suffixMatch := func(value, suffix string) bool {
		i := strings.Index(value, suffix)
		return i >= 0 && len(value)-i == len(suffix)
	}
	for p := questProgressHead; p != nil; {
		next := p.Next
		value := monsterDefinitionText(p.Name[:])
		remove := false
		switch {
		case star == len(pattern)-1:
			prefix := pattern[:star]
			remove = len(value) >= len(prefix) && mapThemeLower(value[:len(prefix)]) == mapThemeLower(prefix)
		case star == 0:
			remove = suffixMatch(value, pattern[1:])
		default:
			colon := strings.IndexByte(pattern, ':')
			if colon >= 0 && colon+2 <= len(pattern) && len(value) >= colon+1 && mapThemeLower(value[:colon+1]) == mapThemeLower(pattern[:colon+1]) {
				start := colon + 2
				if start > len(value) {
					start = len(value)
				}
				remove = suffixMatch(value[start:], pattern[colon+2:])
			}
		}
		if remove {
			questProgressRemove(p)
		}
		p = next
	}
}

// QuestStatusSet sets a quest status variable; isFloat marks a float value (as bits).
func QuestStatusSet(name string, value uint32, isFloat bool) {
	kind := uint32(0)
	if isFloat {
		kind = 1
	}
	questProgressSet(name, value, kind)
}

// QuestStatus returns a quest status variable.
func QuestStatus(name string) uint32 { return questProgressInt(name) }

// QuestStatusFloat returns a float quest status variable.
func QuestStatusFloat(name string) float32 { return float32(questProgressFloat(name)) }

// QuestStatusReset removes quest status variables matching a name.
func QuestStatusReset(name string) { questProgressReset(name) }
