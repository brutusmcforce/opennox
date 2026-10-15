package eval

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// caseInsensitiveFS reports whether names in dir resolve regardless of case, as they do on macOS.
func caseInsensitiveFS(t *testing.T, dir string) bool {
	t.Helper()
	probe := filepath.Join(dir, "CaseProbe")
	if err := os.WriteFile(probe, nil, 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(probe)
	_, err := os.Stat(filepath.Join(dir, "caseprobe"))
	return err == nil
}

func TestMapScriptDir(t *testing.T) {
	maps := t.TempDir()
	for _, d := range []string{"BluDeath", "con01a"} {
		if err := os.Mkdir(filepath.Join(maps, d), 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(maps, "estate"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	ci := caseInsensitiveFS(t, maps)
	for _, c := range []struct {
		name, want string
		ok         bool
	}{
		{"bludeath", "BluDeath", true},
		{"BluDeath", "BluDeath", true},
		{"con01a", "con01a", true},
		{"estate", "", false}, // a file, not a folder
		{"missing", "", false},
	} {
		got, ok := mapScriptDir(maps, c.name)
		// macOS: the file system resolves the name itself, so mapScriptDir returns it as it was given;
		// it still points at the right folder, which is what the caller needs.
		same := got == c.want
		if ci {
			same = strings.EqualFold(got, c.want)
		}
		if !same || ok != c.ok {
			t.Errorf("%s: got %q %v, want %q %v", c.name, got, ok, c.want, c.ok)
		}
	}
}
