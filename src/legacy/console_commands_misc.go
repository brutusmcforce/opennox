package legacy

import (
	"github.com/brutusmcforce/opennox/libs/console"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func Nox_cmd_allow_user(_ int, _ []string) bool {
	consoleCommandPrint("notyetimplemented")
	return true
}
func Nox_cmd_allow_ip(i int, t []string) bool     { return Nox_cmd_allow_user(i, t) }
func Nox_cmd_set_spellpts(i int, t []string) bool { return Nox_cmd_allow_user(i, t) }
func Nox_cmd_offonly1(_ int, t []string) bool {
	if len(t) != 2 {
		return false
	}
	teamRuntimeCreateMap()
	return true
}
func Nox_cmd_set_net_debug(_ int, t []string) bool {
	if len(t) != 2 {
		return false
	}
	noxflags.SetEngine(noxflags.EngineNetDebug)
	return true
}
func Nox_cmd_unset_net_debug(_ int, t []string) bool {
	if len(t) != 2 {
		return false
	}
	noxflags.UnsetEngine(noxflags.EngineNetDebug)
	return true
}
func Nox_cmd_show_info(_ int, _ []string) bool {
	Nox_xxx_clientPlaySoundSpecial_452D80(921, 100)
	sub_435F60()
	return true
}
func Nox_cmd_show_mem(_ int, _ []string) bool {
	sessionMapFilename()
	Nox_xxx_gameLoopMemDump_413E30()
	return true
}
func Nox_cmd_show_rank(_ int, _ []string) bool {
	if noxflags.HasGame(8192) {
		Sub_4703F0()
	}
	return true
}
func Nox_cmd_show_motd(_ int, t []string) bool {
	if len(t) != 2 {
		return false
	}
	sessionMOTDShow()
	return true
}
func Nox_cmd_show_seq(_ int, t []string) bool {
	if len(t) != 2 {
		return false
	}
	// The former queue walk produced no output or state changes.
	return true
}
func Nox_cmd_list_maps(_ int, _ []string) bool {
	row := ""
	n := 0
	for p := mapCatalogFirst(); p != nil; p = mapCatalogNext(p) {
		row += consoleCommandWiden(alloc.GoString(&p.Name[0])+alloc.GoString(memmap.PtrUint8(0x587000, 103292))) + "\t\t"
		n++
		if n%4 == 0 {
			GetConsole().Print(console.ColorRed, consoleCommandFormat(alloc.GoString16(memmap.PtrUint16(0x587000, 103276)), row))
			row = ""
		}
	}
	if n%4 != 0 {
		GetConsole().Print(console.ColorRed, consoleCommandFormat(alloc.GoString16(memmap.PtrUint16(0x587000, 103284)), row))
	}
	return true
}
func Nox_cmd_window(i int, t []string) bool {
	if len(t) > 1 {
		n := int(consoleCommandNumber(t[i]))
		if len(t[i]) > 0 && (t[i][0] == '+' || t[i][0] == '-') {
			Nox_draw_setCutSize_476700(0, n)
		} else {
			Nox_draw_setCutSize_476700(n, 0)
		}
	}
	return true
}
func Nox_cmd_menu_options(_ int, _ []string) bool {
	if !noxflags.HasGame(8) && noxflags.HasGame(8192) {
		serverOptionsConstruct()
	}
	return true
}
func Nox_cmd_menu_vidopt(_ int, _ []string) bool { optionsShow(); return true }
func Nox_cmd_reenter(_ int, _ []string) bool {
	if !noxflags.HasGame(8192) {
		playerStateReentry(1)
	}
	return true
}
