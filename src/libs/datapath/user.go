package datapath

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/brutusmcforce/opennox/libs/env"
)

// User data.
//
// The original game writes saves and a few other files (ban list, etc) to the game directory.
// OpenNox keeps them in a per-user directory instead, so that the game directory may stay read-only
// (e.g. installed to Program Files on Windows), and multiple users don't share their saves.

var userdir struct {
	sync.Once
	path string
}

func getUser() string {
	userdir.Do(func() {
		if userdir.path == "" {
			userdir.path = FindUser()
		}
	})
	return userdir.path
}

// User returns the directory for the user's game data. If additional args are provided, they will be joined with it.
// If no user directory was set, it is selected automatically using FindUser.
func User(path ...string) string {
	if len(path) == 0 {
		return getUser()
	}
	args := make([]string, 0, 1+len(path))
	args = append(args, getUser())
	args = append(args, path...)
	return filepath.Join(args...)
}

// SetUser sets the directory for the user's game data.
func SetUser(dir string) {
	userdir.Do(func() {})
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	userdir.path = dir
	Log.Printf("setting user data dir to: %q", dir)
}

// FindUser returns the default directory for the user's game data.
// It does not affect the directory returned by User.
//
// On Linux it's ~/.local/share/opennox (or $XDG_DATA_HOME/opennox), on Windows it's %AppData%\OpenNox.
// In tests, and if no such directory can be determined, it's the game data directory.
func FindUser() string {
	if env.IsE2E() {
		// tests expect everything in the game directory
		return Data()
	}
	if dir := env.AppUserDir(); dir != "" {
		return dir
	}
	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return Data()
		}
		return filepath.Join(dir, "OpenNox")
	}
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, "opennox")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return Data()
	}
	// Linux Snapcraft installation replaces HOME variable
	if rhome := os.Getenv("SNAP_REAL_HOME"); rhome != "" {
		home = rhome
	}
	return filepath.Join(home, ".local", "share", "opennox")
}
