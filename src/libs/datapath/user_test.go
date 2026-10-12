package datapath

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/brutusmcforce/opennox/libs/common"
)

func TestFindUser(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses Linux paths")
	}
	t.Setenv("NOX_E2E", "")
	t.Setenv("NOX_E2E_RECORD", "")
	t.Setenv("SNAP_NAME", "")
	t.Setenv("SNAP_REAL_HOME", "")
	t.Setenv("XDG_DATA_HOME", "/xdg/data")
	if got, want := FindUser(), filepath.Join("/xdg/data", "opennox"); got != want {
		t.Errorf("with XDG_DATA_HOME: got %q, want %q", got, want)
	}
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("HOME", "/home/user")
	if got, want := FindUser(), filepath.Join("/home/user", ".local", "share", "opennox"); got != want {
		t.Errorf("with HOME: got %q, want %q", got, want)
	}
	t.Setenv("SNAP_REAL_HOME", "/home/real")
	if got, want := FindUser(), filepath.Join("/home/real", ".local", "share", "opennox"); got != want {
		t.Errorf("with SNAP_REAL_HOME: got %q, want %q", got, want)
	}
}

func TestSaveDirs(t *testing.T) {
	user, data := t.TempDir(), t.TempDir()
	SetUser(user)
	SetData(data)
	if got, want := Save("a", "b.dat"), filepath.Join(user, common.SaveDir, "a", "b.dat"); got != want {
		t.Errorf("Save: got %q, want %q", got, want)
	}
	if got, want := GameSave("a"), filepath.Join(data, common.SaveDir, "a"); got != want {
		t.Errorf("GameSave: got %q, want %q", got, want)
	}
	if got, want := User("ban.txt"), filepath.Join(user, "ban.txt"); got != want {
		t.Errorf("User: got %q, want %q", got, want)
	}
}
