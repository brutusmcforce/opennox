package legacy

import (
	"sync"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/datapath"

	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc/handles"
)

var files struct {
	sync.RWMutex
	byHandle map[unsafe.Pointer]*binfile.File
}

// FILE identifies a file in the protected handle region. It is never dereferenced
// and does not contain libc FILE storage. Keep it non-zero-sized for pointer identity.
type FILE struct{ _ byte }

// SaveRoot returns the directory that contains the save directory. It may differ from the game directory.
var SaveRoot func() string

func saveRoot() string {
	if SaveRoot != nil {
		return SaveRoot()
	}
	return datapath.Data()
}

// UserFilePath returns the path of a file written by the game (e.g. the ban list), which is not stored
// in the game directory.
var UserFilePath func(name string) string

func userFilePath(name string) string {
	if UserFilePath != nil {
		return UserFilePath(name)
	}
	return name
}

func nox_fs_fread(f *FILE, dst unsafe.Pointer, sz int) int {
	fp := fileByHandle(f)
	n, _ := fp.Read(unsafe.Slice((*byte)(dst), sz))
	return n
}

func fileByHandle(f *FILE) *binfile.File {
	h := unsafe.Pointer(f)
	handles.AssertValidPtr(h)
	files.RLock()
	fp := files.byHandle[h]
	files.RUnlock()
	return fp
}

func nox_fs_close(f *FILE) {
	if f == nil {
		return
	}
	h := unsafe.Pointer(f)
	handles.AssertValidPtr(h)
	files.Lock()
	defer files.Unlock()
	fp := files.byHandle[h]
	if fp != nil {
		_ = fp.Close()
		delete(files.byHandle, h)
	}
}

func Nox_fs_close(f *FILE) {
	nox_fs_close(f)
}

func NewFileHandle(f *binfile.File) *FILE {
	if f.Handle != nil {
		return (*FILE)(f.Handle)
	}
	f.Handle = handles.NewPtr()
	files.Lock()
	defer files.Unlock()
	if files.byHandle == nil {
		files.byHandle = make(map[unsafe.Pointer]*binfile.File)
	}
	files.byHandle[f.Handle] = f
	return (*FILE)(f.Handle)
}
