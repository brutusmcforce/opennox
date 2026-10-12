package legacy

import (
	"errors"
	"unsafe"

	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_savegame_rm                      func(name string, rmDir bool) error
	Nox_client_countPlayerFiles04_4DC7D0 func() int
	Nox_xxx_gameGet_4DB1B0               func() bool
	Sub_4DCC90                           func() int
	Sub_4DB1C0                           func() unsafe.Pointer
	Sub_4DCBF0                           func(a1 int)
	Nox_xxx_serverIsClosing_446180       func() int
	Sub_4DCC10                           func(a1p *server.Object) int
	Sub_4DCFB0                           func(a1p *server.Object)
	Sub_4DD0B0                           func(a1p *server.Object)
	Nox_setSaveFileName_4DB130           func(s string)
)

func nox_xxx_gameGet_4DB1B0() int { return bool2int(Nox_xxx_gameGet_4DB1B0()) }

func sub_4DD0B0(a1p *nox_object_t) { Sub_4DD0B0(asObjectS(a1p)) }
func Nox_xxx_destroyEveryChatMB_528D60() {
	visibilityClearChats()
}
func Nox_xxx_quickBarClose_4606B0() {
	quickbarCloseExpanded()
}
func Nox_xxx_monstersAllBelongToHost_4DB6A0() {
	runtimeHostOwnership()
}
func Nox_xxx_mapSaveMap_51E010(a1 string, a2 int) bool {
	return worldMapSaveService(a1, a2)
}

func Sub_41A590(cf *cryptfile.CryptFile, u *server.Object, pinfo *server.PlayerInfo) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileAttributes(u, pinfo.C()) == 0 {
		return errors.New("failed")
	}
	return nil
}

func Sub_41AA30(cf *cryptfile.CryptFile, u *server.Object, pinfo *server.PlayerInfo) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileStatus(u) == 0 {
		return errors.New("failed")
	}
	return nil
}

func Sub_41AC30(cf *cryptfile.CryptFile, u *server.Object, pinfo *server.PlayerInfo) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileInventory(u) == 0 {
		return errors.New("failed")
	}
	return nil
}

func Nox_xxx_guiFieldbook_41B420(cf *cryptfile.CryptFile, u *server.Object, pinfo *server.PlayerInfo) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileGuides(u) == 0 {
		return errors.New("failed")
	}
	return nil
}

func Nox_xxx_guiSpellbook_41B660(cf *cryptfile.CryptFile, u *server.Object, pinfo *server.PlayerInfo) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileSpells(u) == 0 {
		return errors.New("failed")
	}
	return nil
}

func Nox_xxx_guiEnchantment_41B9C0(cf *cryptfile.CryptFile, u *server.Object, pinfo *server.PlayerInfo) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileEnchantment(u) == 0 {
		return errors.New("failed")
	}
	return nil
}

func Sub_41BEC0(cf *cryptfile.CryptFile, u *server.Object, pinfo *server.PlayerInfo) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileJournal(u) == 0 {
		return errors.New("failed")
	}
	return nil
}

func Sub_41C080(cf *cryptfile.CryptFile, u *server.Object, pinfo *server.PlayerInfo) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileGame(u) == 0 {
		return errors.New("failed")
	}
	return nil
}

func Sub_41C280(cf *cryptfile.CryptFile) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileGUI() == 0 {
		return errors.New("failed")
	}
	return nil
}

func Nox_xxx_parseFileInfoData_41C3B0(cf *cryptfile.CryptFile) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileMetadata() == 0 {
		return errors.New("failed")
	}
	return nil
}

func Sub_41C780(cf *cryptfile.CryptFile) error {
	old := cryptfile.Global()
	cryptfile.SetGlobal(cf)
	defer cryptfile.SetGlobal(old)
	if playerFileMusic() == 0 {
		return errors.New("failed")
	}
	return nil
}
