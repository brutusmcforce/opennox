package opennox

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateScreenshotFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "screenshots")
	now := time.Date(2026, 10, 13, 22, 55, 1, 0, time.Local)
	var names []string
	for i := 0; i < 3; i++ {
		f, err := createScreenshotFile(dir, now)
		require.NoError(t, err)
		require.NoError(t, f.Close())
		names = append(names, filepath.Base(f.Name()))
	}
	require.Equal(t, []string{"opennox-20261013-225501.png", "opennox-20261013-225501-2.png", "opennox-20261013-225501-3.png"}, names)
}
