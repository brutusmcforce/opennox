//go:build !server

package opennox

import (
	"path/filepath"

	"github.com/brutusmcforce/opennox/libs/datapath"
	libsdl "github.com/brutusmcforce/opennox/libs/sdl"
)

// askDataDirPrompt is the title of the folder dialog shown when the Nox game files are not found.
const askDataDirPrompt = "OpenNox needs the files of the original Nox (for example from GoG). Select the Nox folder: it contains thing.bin and gamedata.bin."

// askDataDir asks the player for the Nox game folder with the system's folder dialog (the desktop's own through
// xdg-desktop-portal on Linux), until a folder with the game files is picked. It returns "" if the player cancels,
// or if there's no dialog.
func askDataDir() string {
	for {
		path, err := libsdl.ShowOpenFolderDialog(askDataDirPrompt)
		if err != nil {
			datapath.Log.Println("cannot show the folder dialog:", err)
			return ""
		}
		if path == "" {
			return "" // cancelled
		}
		path = filepath.Clean(path)
		if datapath.CheckData(path) {
			return path
		}
		errorMessage("This folder doesn't contain the Nox game files (thing.bin, gamedata.bin). Please select the folder where Nox is installed.")
	}
}
