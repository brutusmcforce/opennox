package maps

import (
	"os"
	"path/filepath"

	"github.com/brutusmcforce/opennox/libs/ifs"
)

// atomicFile is a downloaded file: it's written next to its destination and only replaces it when complete, so an
// aborted download never leaves a broken map. Replacing also never writes through a link at the destination.
type atomicFile struct {
	*os.File
	dst string
}

func createAtomic(path string) (*atomicFile, error) {
	path = ifs.Normalize(path)
	f, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*.part")
	if err != nil {
		return nil, err
	}
	// a temporary file is private; the map is not
	_ = f.Chmod(0644)
	return &atomicFile{File: f, dst: path}, nil
}

// Commit replaces the destination with the file.
func (f *atomicFile) Commit() error {
	if err := f.Close(); err != nil {
		_ = os.Remove(f.Name())
		return err
	}
	if err := os.Rename(f.Name(), f.dst); err != nil {
		_ = os.Remove(f.Name())
		return err
	}
	return nil
}

// Abort deletes the file and leaves the destination as it was.
func (f *atomicFile) Abort() {
	_ = f.Close()
	_ = os.Remove(f.Name())
}
