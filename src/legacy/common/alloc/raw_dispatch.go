package alloc

import (
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
)

// RawCalloc allocates zeroed unmanaged memory without registering it in allocs.
func RawCalloc(num, size uintptr) unsafe.Pointer {
	p := rawCalloc(num, size)
	if layout.Checking {
		layout.Track(p, num*size, nil, heapName())
	}
	return p
}

// heapName names a heap block by the first caller outside the allocators (for layout checks).
func heapName() string {
	var pcs [8]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		f, more := frames.Next()
		if !strings.Contains(f.Function, "/alloc.") && !strings.Contains(f.File, "raw_alloc") && !strings.Contains(f.Function, "legacy.legacyCalloc") && !strings.Contains(f.Function, "legacy.legacyMalloc") && !strings.Contains(f.Function, "legacy.legacyRealloc") {
			file := f.File
			if i := strings.Index(file, "/src/"); i >= 0 {
				file = file[i+len("/src/"):]
			}
			return "heap@" + file + ":" + strconv.Itoa(f.Line)
		}
		if !more {
			return "heap"
		}
	}
}

// RawFree releases unmanaged memory without consulting or changing allocs.
func RawFree(ptr unsafe.Pointer) {
	layout.Untrack(ptr)
	rawFree(ptr)
}
