//go:build !windows
// +build !windows

package ifs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

func TestNormalize(t *testing.T) {
	dir, err := os.MkdirTemp("", "nox_fs_")
	must.NoError(t, err)
	defer os.RemoveAll(dir)

	dir1 := filepath.Join(dir, "AbC", "def")
	err = os.MkdirAll(dir1, 0755)
	must.NoError(t, err)
	file1 := filepath.Join(dir1, "File.txt")

	err = os.WriteFile(file1, []byte("data"), 0644)
	must.NoError(t, err)

	got := Normalize(strings.Join([]string{dir, "abc", "Def", "FILE.TXT"}, "\\"))
	if caseInsensitiveFS(t, dir) {
		// macOS: the file system resolved the path itself, so Normalize returned it unchanged. Any case
		// opens the file there, which is all Normalize has to guarantee.
		must.EqOp(t, filepath.Join(dir, "abc", "Def", "FILE.TXT"), got)
	} else {
		must.EqOp(t, file1, got)
	}
	must.EqOp(t,
		filepath.Join(dir1, "NotExistent"),
		Normalize(strings.Join([]string{dir, "ABC", "DeF", "NotExistent"}, "\\")),
	)
}

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
