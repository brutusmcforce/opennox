package maps

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/shoenig/test/must"
)

// Map files linked from elsewhere are sent; links to other files and to directories are not.
func TestCompressMapSymlinks(t *testing.T) {
	src := t.TempDir()
	secret := filepath.Join(src, "secret.txt")
	must.NoError(t, os.WriteFile(secret, []byte("secret"), 0o644))
	mapFile := filepath.Join(src, "real.map")
	must.NoError(t, os.WriteFile(mapFile, []byte("map data"), 0o644))
	must.NoError(t, os.Mkdir(filepath.Join(src, "sub"), 0o755))

	dir := filepath.Join(t.TempDir(), "test")
	must.NoError(t, os.Mkdir(dir, 0o755))
	must.NoError(t, os.Symlink(mapFile, filepath.Join(dir, "test.map")))
	must.NoError(t, os.WriteFile(filepath.Join(dir, "test.rul"), []byte("rules"), 0o644))
	must.NoError(t, os.Symlink(secret, filepath.Join(dir, "notes.txt")))
	must.NoError(t, os.Symlink(filepath.Join(src, "sub"), filepath.Join(dir, "linked")))
	must.NoError(t, os.Symlink(filepath.Join(src, "missing.map"), filepath.Join(dir, "broken.map")))

	var buf bytes.Buffer
	must.NoError(t, CompressMap(&buf, nil, dir))
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	must.NoError(t, err)
	var names []string
	for _, f := range zr.File {
		names = append(names, f.Name)
	}
	slices.Sort(names)
	must.Eq(t, []string{"test.map", "test.rul"}, names)
}
