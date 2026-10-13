package legacy

import (
	"encoding/binary"
	"image"
	"strings"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func browserMapPolygons() {
	if browserUI.polygonsReady != 0 {
		return
	}
	for off := uintptr(165744); off < 166016; off += 16 {
		count := int(int8(memmap.Uint8(0x587000, off)))
		for i := 0; i < count; i++ {
			index := int32(int8(memmap.Uint8(0x587000, off+1+uintptr(i))))
			p := uintptr(uint32(165104 + 8*index))
			mapPolygonSelect(memmap.Float32(0x587000, p), memmap.Float32(0x587000, p+4))
		}
		var polygon mapPolygon
		name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, off+12)))
		copy(polygon.Name[:], name)
		mapPolygonDefaultsWrite(&polygon)
		mapPolygonConstruct()
	}
	browserUI.polygonsReady = 1
}
func browserMapPolygonsClear() unsafe.Pointer {
	if browserUI.polygonsReady == 0 {
		return nil
	}
	p := mapPolygonReset()
	browserUI.polygonsReady = 0
	return p
}
func browserSwitchChatMap() {
	Nox_client_gui_set_flag_815132(0)
	noxflags.SetGame(5)
	GetClient().SetMouseBounds(image.Rect(0, 0, int(nox_win_width)-1, int(nox_win_height)-1))
	name := Nox_client_getChatMap_49FF40()
	if strings.Contains(name, ".") {
		return
	}
	p, free := alloc.CString(name)
	sessionSetSelectedMap(p)
	free()
	suffix := unsafe.Slice(memmap.PtrUint8(0x587000, 90856), 4)
	name += string(suffix)
	p, free = alloc.CString(name)
	defer free()
	if !noxflags.HasEngine(noxflags.EngineReplayRead) {
		sessionSetMapPath(p)
	}
	noxflags.UnsetGame(55280)
	noxflags.SetGame(128)
	slot := unsafe.Slice(serverConfigSlotSelect(0), 58)
	binary.LittleEndian.PutUint16(slot[52:], binary.LittleEndian.Uint16(slot[52:])&0x280f|0x80)
}

func sub_49FDB0(_ int32) { browserMapPolygons() }

func sub_49FF20() *uint32 { return (*uint32)(browserMapPolygonsClear()) }

func nox_client_xxx_switchChatMap_43B510() { browserSwitchChatMap() }
