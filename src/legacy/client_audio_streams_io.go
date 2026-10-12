package legacy

import (
	"encoding/binary"
	"fmt"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"io"
	"unsafe"
)

func audioStreamClose(p *audioStreamCatalog) {
	p.Active = nil
	if p.Override != nil {
		Nox_fs_close(p.Override)
		p.Override = nil
	}
}
func audioStreamRead(p *audioStreamCatalog, dst unsafe.Pointer, requested int32) int32 {
	if p.Active == nil {
		return 0
	}
	n := requested
	if n > int32(p.Remaining) {
		n = int32(p.Remaining)
	}
	got := int32(0)
	if n > 0 {
		got = int32(nox_fs_fread(p.Active, dst, int(n)))
		if got < 0 {
			got = 0
		}
	}
	p.Remaining -= uint32(got)
	return got
}
func audioStreamOpen(p *audioStreamCatalog, index int32) int32 {
	entry := audioStreamCatalogEntry(p, index)
	audioStreamClose(p)
	p.Active = p.Bag
	p.Remaining = entry.Length
	result := int32(1)
	if _, err := fileByHandle(p.Bag).Seek(int64(int32(entry.Offset)), io.SeekStart); err != nil {
		result = 0
	}
	if entry.Length == 0 {
		result = 0
	}
	if p.UseOverride == 0 {
		return result
	}
	path := alloc.GoString(&p.Directory[0]) + alloc.GoString(&entry.Name[0]) + ".wav"
	if len(path) >= 280 {
		return result
	}
	file, err := ifs.Open(path)
	if err != nil {
		return result
	}
	p.Override = NewFileHandle(binfile.NewFile(file))
	stream := fileByHandle(p.Override)
	invalid := func() int32 { Nox_fs_close(p.Override); p.Override = nil; return result }
	var header [12]byte
	n, _ := stream.Read(header[:])
	if n != 12 || string(header[:4]) != "RIFF" || string(header[8:]) != "WAVE" {
		fmt.Printf("error: '%s' is bad - cannot read\n", path)
		return invalid()
	}
	var format [16]byte
	var chunk [8]byte
	haveFormat := false
	length := uint32(0)
	for {
		n, _ = stream.Read(chunk[:])
		if n != 8 {
			break
		}
		size := binary.LittleEndian.Uint32(chunk[4:])
		switch string(chunk[:4]) {
		case "fmt ":
			if size < 16 {
				return invalid()
			}
			n, _ = stream.Read(format[:])
			if n != 16 || binary.LittleEndian.Uint16(format[2:]) == 0 {
				return invalid()
			}
			haveFormat = true
			stream.Seek(int64(int32(size-16)), io.SeekCurrent)
		case "data":
			length = size
			goto complete
		default:
			stream.Seek(int64(int32(size)), io.SeekCurrent)
		}
	}
complete:
	if !haveFormat {
		return invalid()
	}
	channels := binary.LittleEndian.Uint16(format[2:])
	align := binary.LittleEndian.Uint16(format[12:])
	entry.Flags = 2
	if align/channels == 2 {
		entry.Flags = 6
	}
	if channels == 2 {
		entry.Flags |= 1
	}
	entry.Rate = binary.LittleEndian.Uint32(format[4:])
	p.Remaining = length
	p.Active = p.Override
	return 1
}
