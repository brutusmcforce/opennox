package legacy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"io"
	"unsafe"
)

func serverConfigFileLine(f *binfile.File) bool {
	raw, err := f.ReadString()
	if err != nil && !errors.Is(err, io.EOF) {
		return false
	}
	StrCopyBytes(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1563936), 1024), string(raw))
	return !errors.Is(err, io.EOF)
}

// Match the first strtok token and its mutation of the shared scratch buffer.
func serverConfigFileToken() *byte {
	p := memmap.PtrUint8(0x5D4594, 1563936)
	buf := unsafe.Slice(p, 1024)
	i := 0
	for i < len(buf) && (buf[i] == '\r' || buf[i] == '\t' || buf[i] == '\n') {
		i++
	}
	if i == len(buf) || buf[i] == 0 {
		return nil
	}
	start := i
	for i < len(buf) && buf[i] != 0 && buf[i] != '\r' && buf[i] != '\t' && buf[i] != '\n' {
		i++
	}
	if i < len(buf) {
		buf[i] = 0
	}
	return &buf[start]
}
func serverConfigFileRead(path string) int32 {
	f, err := ifs.Open(path)
	if err != nil {
		return 0
	}
	bf := binfile.NewTextFile(f)
	defer bf.Close()
	for serverConfigFileLine(bf) {
		if serverConfigFileToken() == nil {
			continue
		}
		// The outer loop historically compares the beginning of the buffer, not the
		// token pointer, so leading delimiters do not normalize section headers.
		switch alloc.GoString(memmap.PtrUint8(0x5D4594, 1563936)) {
		case "[Banned Users]":
			if serverConfigFileReadBlocked(bf) == 0 {
				return 0
			}
		case "[Allowed Users]":
			if serverConfigFileReadAllowed(bf) == 0 {
				return 0
			}
		}
	}
	return 1
}
func serverConfigFileWide(p *byte) []uint16 {
	b := alloc.GoString(p)
	out := make([]uint16, len(b)+1)
	for i := 0; i < len(b); i++ {
		out[i] = uint16(b[i])
	}
	return out
}
func serverConfigFileReadBlocked(f *binfile.File) int32 {
	for serverConfigFileLine(f) {
		p := serverConfigFileToken()
		if p == nil {
			return 1
		}
		name := serverConfigFileWide(p)
		if !serverConfigFileLine(f) {
			return 1
		}
		address := serverConfigFileToken()
		if address == nil {
			return 0
		}
		if alloc.GoString(address) == "0" {
			address = nil
		}
		serverConfigBlockedAdd(0, &name[0], address)
	}
	return 1
}
func serverConfigFileReadAllowed(f *binfile.File) int32 {
	for serverConfigFileLine(f) {
		p := serverConfigFileToken()
		if p == nil {
			break
		}
		name := serverConfigFileWide(p)
		serverConfigAllowedAdd(&name[0])
	}
	return 1
}
func serverConfigNarrow(p *uint16) string {
	var out []byte
	for i := uintptr(0); ; i += 2 {
		c := *(*uint16)(unsafe.Add(unsafe.Pointer(p), i))
		if c == 0 {
			break
		}
		out = append(out, byte(c))
	}
	if i := bytes.IndexByte(out, 0); i >= 0 {
		out = out[:i]
	}
	return string(out)
}
func serverConfigFileWrite(path string) int32 {
	f, err := ifs.Create(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	fmt.Fprintf(f, "%s\n", alloc.GoString(memmap.PtrUint8(0x587000, 202212)))
	for p := serverConfigBlockedFirst(); p != nil; p = serverConfigBlockedNext(p) {
		if p.expires != 0 {
			continue
		}
		fmt.Fprintf(f, "%s\n", serverConfigNarrow(&p.name[0]))
		address := alloc.GoString(&p.address[0])
		if address == "" {
			address = "0"
		}
		fmt.Fprintf(f, "%s\n", address)
	}
	fmt.Fprintf(f, "\n%s\n", alloc.GoString(memmap.PtrUint8(0x587000, 202228)))
	for p := serverConfigAllowedFirst(); p != nil; p = serverConfigAllowedNext(p) {
		fmt.Fprintf(f, "%s\n", serverConfigNarrow(&p.name[0]))
	}
	return 1
}
