package legacy

import (
	"strings"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/datapath"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/server"
)

func mapOrchestrationSetName(name string) uint32 {
	if i := strings.IndexByte(name, 0); i >= 0 {
		name = name[:i]
	}
	dst := unsafe.Slice(memmap.PtrUint8(0x587000, 197860), len(name)+1)
	copy(dst, name)
	dst[len(name)] = 0
	return uint32(len(dst))
}
func mapOrchestrationName() unsafe.Pointer { return memmap.PtrOff(0x587000, 197860) }
func mapOrchestrationStart() uint32 {
	Nox_xxx_mapSwitchLevel_4D12E0(true)
	noxflags.SetGame(0x400000)
	saved := memmap.PtrUint32(0x5D4594, 1550924)
	*saved = memmap.PtrWord(unsafe.Pointer(GetServer().S().Objs.GetAndZeroObjects()))
	clear(unsafe.Slice(memmap.PtrUint8(0x973F18, 2408), 1464))
	for tries := 100; ; {
		result := mapOrchestrationStep()
		if result == 1 {
			break
		}
		if result == 0 {
			return 0
		}
		tries--
		if tries == 0 {
			return 0
		}
	}
	// the script of the generated map is kept in memory instead of the nc.obj file
	if !mapGenUseBlendScript() {
		return 0
	}
	for obj := GetServer().S().Objs.First(); obj != nil; obj = obj.ObjNext {
		obj.ScriptIDVal = 0
	}
	Nox_xxx_mapGenMakeInfo_4D5DB0(memmap.PtrOff(0x973F18, 2408))
	name := GoStringP(mapOrchestrationName())
	directory := datapath.Data() + "\\Maps\\$" + name
	_ = ifs.Mkdir(directory)
	if !Nox_xxx_mapSaveMap_51E010(directory+"\\$"+name+".map", 1) {
		return 0
	}
	Nox_xxx_mapSwitchLevel_4D12E0(true)
	GetServer().S().Objs.SetObjects((*server.Object)(memmap.WordPtr(*saved)))
	*saved = 0
	noxflags.UnsetGame(0x400000)
	return 1
}
func mapOrchestrationStartAlt() uint32 {
	noxflags.SetGame(0x400000)
	clear(unsafe.Slice(memmap.PtrUint8(0x973F18, 2408), 1464))
	for tries := 100; ; {
		result := mapOrchestrationStep()
		if result == 1 {
			break
		}
		if result == 0 {
			return 0
		}
		tries--
		if tries == 0 {
			return 0
		}
	}
	// the script of the generated map is kept in memory instead of the nc.obj file
	if !mapGenUseBlendScript() {
		return 0
	}
	Nox_xxx_mapGenMakeInfo_4D5DB0(memmap.PtrOff(0x973F18, 2408))
	noxflags.UnsetGame(0x400000)
	return 1
}
