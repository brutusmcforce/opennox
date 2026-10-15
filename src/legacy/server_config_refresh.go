package legacy

import (
	"encoding/binary"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

var serverConfigRecordDirty int32

func serverConfigRefresh() int32 {
	var spells [5]uint32
	var weapons uint32
	serverPanelsSpellSnapshot(&spells[0])
	serverPanelsWeaponSnapshot(&weapons)
	armor := serverPanelsArmorSnapshot()
	count, limit := GetServer().S().Players.Count(), int(serverConfigLimitGet())
	if noxflags.HasEngine(noxflags.EngineNoRendering) {
		count--
		limit--
	}
	config := unsafe.Slice((*byte)(serverConfigSettings()), 184)
	for _, v := range []struct{ offset, value int }{{103, count}, {104, limit}} {
		if int(config[v.offset]) != v.value {
			config[v.offset] = byte(v.value)
			serverConfigRecordDirty = 1
		}
	}
	if int(config[102]&0xef) != Sub_43BE50_get_video_mode_id() {
		config[102] = config[102]&0x80 | byte(Sub_43BE50_get_video_mode_id())
		serverConfigRecordDirty = 1
	}
	record := unsafe.Slice(serverConfigSlot(0), 58)
	flags := uint16(noxflags.GetGame())
	mode := binary.LittleEndian.Uint16(record[52:])
	if (mode^flags)&0xfff0 != 0 {
		binary.LittleEndian.PutUint16(record[52:], flags)
		mode = flags
		serverConfigRecordDirty = 1
	}
	// C promotes the signed score before comparing to unsigned 16-bit storage.
	score := serverConfigScore(int16(mode))
	if int(binary.LittleEndian.Uint16(record[54:])) != int(score) {
		binary.LittleEndian.PutUint16(record[54:], uint16(serverConfigScore(int16(mode))))
		serverConfigRecordDirty = 1
	}
	if record[56] != serverConfigMinutes(int16(mode)) {
		record[56] = serverConfigMinutes(int16(mode))
		serverConfigRecordDirty = 1
	}
	var name [15]byte
	copy(name[:], alloc.GoString(serverConfigNameGet()))
	// strncmp stops at the first shared NUL, leaving any bytes beyond it alone.
	same := true
	for i := range name {
		if record[9+i] != name[i] {
			same = false
			break
		}
		if name[i] == 0 {
			break
		}
	}
	if !same {
		copy(record[9:24], name[:])
		serverConfigRecordDirty = 1
	}
	if mapName := Nox_xxx_mapGetMapName_409B40(); mapNameField(&record[0]) != mapName {
		mapNameSetField(&record[0], mapName)
		serverConfigRecordDirty = 1
	}
	for i, v := range spells {
		off := 24 + 4*i
		if binary.LittleEndian.Uint32(record[off:]) != v {
			for j, x := range spells {
				binary.LittleEndian.PutUint32(record[24+4*j:], x)
			}
			serverConfigRecordDirty = 1
			break
		}
	}
	if !noxflags.HasGame(1) {
		return serverConfigRecordDirty
	}
	if binary.LittleEndian.Uint32(record[48:]) != armor {
		binary.LittleEndian.PutUint32(record[48:], armor)
		serverConfigRecordDirty = 1
	}
	for i := 0; i < 4; i++ {
		if int(record[44+i]) != int(int8(weapons>>uint(8*i))) {
			binary.LittleEndian.PutUint32(record[44:], weapons)
			serverConfigRecordDirty = 1
			break
		}
	}
	return serverConfigRecordDirty
}
