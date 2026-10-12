package datapath

import (
	"path/filepath"
	"strings"

	"github.com/brutusmcforce/opennox/libs/common"
)

// SaveNameFromPath gives a save directory name based on a path to player save file.
func SaveNameFromPath(src string) string {
	if src == "" {
		return ""
	}
	src = strings.ReplaceAll(src, "\\", "/")
	sub := strings.Split(src, "/")
	n := len(sub)
	if n == 1 {
		return src
	}
	if last := sub[n-1]; filepath.Ext(last) == "" {
		return last
	}
	return sub[len(sub)-2]
}

// Save returns the current Nox save dir in the user data dir (see User).
// If additional args are provided, they will be joined with the save dir.
func Save(path ...string) string {
	args := make([]string, 0, 2+len(path))
	args = append(args, getUser(), common.SaveDir)
	args = append(args, path...)
	return filepath.Join(args...)
}

// GameSave returns the save dir in the game data dir, where the original game keeps saves.
// If additional args are provided, they will be joined with the save dir.
func GameSave(path ...string) string {
	args := make([]string, 0, 2+len(path))
	args = append(args, getData(), common.SaveDir)
	args = append(args, path...)
	return filepath.Join(args...)
}
