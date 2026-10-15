package eval

import (
	"os"
	"path/filepath"
	"testing"
)

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
		if got != c.want || ok != c.ok {
			t.Errorf("%s: got %q %v, want %q %v", c.name, got, ok, c.want, c.ok)
		}
	}
}
