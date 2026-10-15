package legacy

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

// Long map names.
//
// The original game limits map names to 8 characters (all its maps follow the 8.3 file names), and keeps them in
// 9-byte fields of records that are copied around as raw memory: the game settings (map and server name), the map
// list, the server list entries. OpenNox allows longer names: a field gets a short token for a long name instead,
// which mapNameResolve turns back into the name. Copies of a record keep the token, so they keep the name too.
//
// The tokens are only valid in this process: a record sent to another game must not carry one (see
// clientSessionSettings).

// MaxMapNameLen is the longest map name OpenNox supports. MsgUseMap, which tells clients which map to load, has
// 32 bytes for the file name with the ".map" extension.
const MaxMapNameLen = 28

// mapNameFieldLen is the size of the original map name fields, with the terminating zero.
const mapNameFieldLen = 9

const mapNameTokenPrefix = '~'

var longMapNames struct {
	list  []string
	index map[string]int
}

// mapNameToken returns the name to store in a short field: the name itself, or a token for a long name.
func mapNameToken(name string) string {
	if len(name) < mapNameFieldLen {
		return name
	}
	key := strings.ToLower(name)
	i, ok := longMapNames.index[key]
	if !ok {
		if longMapNames.index == nil {
			longMapNames.index = make(map[string]int)
		}
		i = len(longMapNames.list)
		longMapNames.list = append(longMapNames.list, name)
		longMapNames.index[key] = i
	}
	return fmt.Sprintf("%c%07x", mapNameTokenPrefix, i)
}

// mapNameResolve returns the map name for a value read from a short field (see mapNameToken).
func mapNameResolve(v string) string {
	if len(v) != mapNameFieldLen-1 || v[0] != mapNameTokenPrefix {
		return v
	}
	i, err := strconv.ParseUint(v[1:], 16, 32)
	if err != nil || int(i) >= len(longMapNames.list) {
		return v
	}
	return longMapNames.list[i]
}

// mapNameIsToken reports if a short field holds a token (see mapNameToken).
func mapNameIsToken(v string) bool {
	return len(v) == mapNameFieldLen-1 && v[0] == mapNameTokenPrefix
}

// mapNameField returns the map name stored in a 9-byte field at p.
func mapNameField(p *byte) string {
	return mapNameResolve(alloc.GoStringS(unsafe.Slice(p, mapNameFieldLen)))
}

// mapNameSetField stores a map name in a 9-byte field at p.
func mapNameSetField(p *byte, name string) {
	field := unsafe.Slice(p, mapNameFieldLen)
	clear(field)
	copy(field[:mapNameFieldLen-1], mapNameToken(name))
}

// MapNameField returns the map name stored in the 9-byte field at p, e.g. the start of the game settings.
func MapNameField(p unsafe.Pointer) string {
	return mapNameField((*byte)(p))
}
