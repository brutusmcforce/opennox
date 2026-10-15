package maps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNativeDownloaderAtomic(t *testing.T) {
	dir := t.TempDir()
	// the map is a link to a file elsewhere, e.g. a shared game install
	shared := filepath.Join(dir, "shared.nxz")
	require.NoError(t, os.WriteFile(shared, []byte("original"), 0644))
	mdir := filepath.Join(dir, "maps", "test")
	require.NoError(t, os.MkdirAll(mdir, 0755))
	dst := filepath.Join(mdir, "test.nxz")
	require.NoError(t, os.Symlink(shared, dst))

	d := NewNativeDownloader(filepath.Join(dir, "maps"))
	data := []byte("downloaded map")

	// an aborted download keeps the old map
	require.NoError(t, d.Start(dst, uint(len(data))))
	d.WritePart(1, data[:4])
	d.CancelAndDelete()
	got, err := os.ReadFile(dst)
	require.NoError(t, err)
	require.Equal(t, "original", string(got))

	// a complete one replaces the link, not the file it points to
	require.NoError(t, d.Start(dst, uint(len(data))))
	d.WritePart(2, data[4:]) // out of order
	d.WritePart(1, data[:4])
	require.True(t, d.Complete())
	require.NoError(t, d.Finish())
	got, err = os.ReadFile(dst)
	require.NoError(t, err)
	require.Equal(t, string(data), string(got))
	got, err = os.ReadFile(shared)
	require.NoError(t, err)
	require.Equal(t, "original", string(got))
	fi, err := os.Lstat(dst)
	require.NoError(t, err)
	require.True(t, fi.Mode().IsRegular())

	// no temporary files are left
	list, err := os.ReadDir(mdir)
	require.NoError(t, err)
	require.Len(t, list, 1)
}
