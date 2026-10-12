package legacy

import (
	"encoding/binary"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func serverTextListing(src, dst []byte) int {
	if len(src) < 12 {
		return 0
	}
	name := alloc.GoString(serverConfigNameGet())
	size := 73 + len(name)
	if len(dst) < size || browserUI.hosting == 0 || Sub_4D6F30() != 0 {
		return 0
	}
	settings := unsafe.Slice((*byte)(serverConfigSettings()), 184)
	slot := unsafe.Slice(serverConfigSlot(0), 58)
	// Assemble separately so aliased input/output retains the original token.
	var buf [72]byte
	buf[2] = 13
	buf[3] = byte(GetServer().S().Players.Count())
	buf[4] = byte(serverConfigLimitGet())
	if noxflags.HasEngine(noxflags.EngineNoRendering) {
		buf[3]--
		buf[4]--
	}
	buf[5], buf[6] = settings[101]&15, settings[101]>>4
	copy(buf[7:10], slot[44:47])
	copy(buf[10:18], Nox_xxx_mapGetMapName_409B40())
	buf[19] = settings[102] | byte(Sub_43BE50_get_video_mode_id())
	buf[20], buf[21] = settings[100], settings[100]&16
	copy(buf[24:28], settings[48:52])
	flags := uint32(noxflags.GetGame())
	if questRuntimeWord(1556160) != 0 {
		flags = flags&^0x80 | 0x1000
		binary.LittleEndian.PutUint16(buf[68:70], uint16(questRuntimeStage()))
	}
	binary.LittleEndian.PutUint32(buf[28:32], flags)
	copy(buf[32:36], slot[48:52])
	copy(buf[36:40], settings[105:109])
	copy(buf[40:44], settings[44:48])
	copy(buf[44:48], src[8:12])
	copy(buf[48:68], slot[24:44])
	copy(dst, buf[:])
	copy(dst[72:], name)
	dst[size-1] = 0
	return size
}
