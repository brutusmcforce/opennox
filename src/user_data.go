package opennox

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/opennox/libs/common"
	"github.com/opennox/libs/datapath"
	"github.com/opennox/libs/env"
	"github.com/opennox/libs/ifs"

	"github.com/opennox/opennox/v1/legacy"
)

// User data.
//
// The original game writes saves and a few other files (ban list, etc) to the game directory.
// OpenNox keeps them in a per-user directory instead, so that the game directory may stay read-only
// (e.g. installed to Program Files on Windows), and multiple users don't share their saves.
// Existing files are copied from the game directory on the first start (see migrateUserData).

var userData struct {
	once sync.Once
	dir  string
}

// userDataDir returns the directory for the user's game data.
//
// On Linux it's ~/.local/share/opennox (or $XDG_DATA_HOME/opennox), on Windows it's %AppData%\OpenNox.
func userDataDir() string {
	userData.once.Do(func() {
		userData.dir = findUserDataDir()
	})
	return userData.dir
}

func findUserDataDir() string {
	if env.IsE2E() {
		// tests expect everything in the game directory
		return datapath.Data()
	}
	if sdir := env.AppUserDir(); sdir != "" {
		return sdir
	}
	if runtime.GOOS == "windows" {
		if dir := configUserDir(); dir != "" {
			return dir
		}
		return datapath.Data()
	}
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, "opennox")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return datapath.Data()
	}
	// Linux Snapcraft installation replaces HOME variable
	if rhome := os.Getenv("SNAP_REAL_HOME"); rhome != "" {
		home = rhome
	}
	return filepath.Join(home, ".local", "share", "opennox")
}

// savePath returns a path in the save directory.
func savePath(path ...string) string {
	args := make([]string, 0, 2+len(path))
	args = append(args, userDataDir(), common.SaveDir)
	args = append(args, path...)
	return filepath.Join(args...)
}

// userFiles are files written by the game, other than saves. They are kept in the user data directory.
var userFiles = []string{
	"ban.txt",
	"server.pem",
}

// userFilePath returns the path of a file written by the game, in the user data directory.
func userFilePath(name string) string {
	return filepath.Join(userDataDir(), name)
}

func init() {
	legacy.SaveRoot = userDataDir
	legacy.UserFilePath = userFilePath
}

// migrateUserData copies saves and other user files from the game directory to the user data directory,
// if they are not there yet. Files in the game directory are not modified.
func migrateUserData() {
	dst := userDataDir()
	src := datapath.Data()
	if dst == "" || src == "" || filepath.Clean(dst) == filepath.Clean(src) {
		return
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		configLog.Printf("cannot create user data dir: %v", err)
		return
	}
	saves := savePath()
	if _, err := os.Stat(saves); os.IsNotExist(err) {
		if old := datapath.Save(); dirExists(old) {
			configLog.Printf("copying saves from %q to %q", old, saves)
			if err := copyDirTree(old, saves); err != nil {
				configLog.Printf("cannot copy saves: %v", err)
			}
		}
	}
	for _, name := range userFiles {
		path := userFilePath(name)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}
		old := datapath.Data(name)
		if _, err := ifs.Stat(old); err != nil {
			continue
		}
		configLog.Printf("copying %q to %q", old, path)
		if err := ifs.Copy(old, path); err != nil {
			configLog.Printf("cannot copy %q: %v", name, err)
		}
	}
}

func dirExists(path string) bool {
	fi, err := ifs.Stat(path)
	return err == nil && fi.IsDir()
}

// copyDirTree copies a directory recursively. The destination must not exist.
func copyDirTree(src, dst string) error {
	src = ifs.Normalize(src)
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}
