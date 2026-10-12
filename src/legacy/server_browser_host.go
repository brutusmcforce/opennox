package legacy

import (
	"encoding/binary"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func browserHostDescription() *byte {
	if onlineRetryPending != 0 && onlineRetryActive != 0 {
		NetworkLogPrint("RECON: Posting server to WOL")
	}
	browserUI.creating = 0
	browserUI.hosting = 1
	r := unsafe.Slice((*byte)(serverConfigSettings()), 184)
	copy(r[111:169], unsafe.Slice(serverConfigSlot(0), 58))
	binary.LittleEndian.PutUint16(r[163:], uint16(noxflags.GetGame()))
	for i := 135; i < 163; i++ {
		r[i] = 255
	}
	clear(r[120:135])
	copy(r[120:135], alloc.GoString(serverConfigNameGet()))
	name := Nox_xxx_mapGetMapName_409B40()
	copy(r[111:], name)
	r[111+len(name)] = 0
	if memmap.Uint32(0x5D4594, 1556160) != 0 {
		stage := uint32(1)
		if onlineRetryActive != 0 {
			stage = questRuntimeStage()
		}
		binary.LittleEndian.PutUint16(r[165:], uint16(stage))
	}
	r[104] = byte(serverConfigLimitGet())
	r[103] = byte(GetServer().S().Players.Count())
	if noxflags.HasEngine(noxflags.EngineNoRendering) {
		r[103]--
		r[104]--
	}
	video := byte(Sub_43BE50_get_video_mode_id())
	// The C translation unit always enabled its high-resolution version define.
	binary.LittleEndian.PutUint32(r[48:], 0x000f039a)
	r[102] = r[102]&0x80 | video
	binary.LittleEndian.PutUint32(r[44:], memmap.Uint32(0x5D4594, 814916))
	binary.LittleEndian.PutUint16(r[109:], uint16(GetServer().S().ServerPort()))
	ClientSetServerHost("localhost")
	slot := serverConfigSlot(1)
	s := unsafe.Slice(slot, 58)
	binary.LittleEndian.PutUint16(s[52:], binary.LittleEndian.Uint16(s[52:])&0xe90f|0x100)
	return slot
}

func sub_43AA70() *byte { return browserHostDescription() }
