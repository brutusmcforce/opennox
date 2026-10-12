package legacy

import (
	"encoding/binary"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
)

func questProgressWrite() int {
	f := cryptfile.Global()
	if f.WriteU16(1) != nil {
		return 0
	}
	var count uint32
	if noxflags.HasGame(2048) {
		for p := questProgressHead; p != nil; p = p.Next {
			count++
		}
	}
	if f.WriteU32(count) != nil {
		return 0
	}
	if !noxflags.HasGame(2048) {
		return 1
	}
	for p := questProgressHead; p != nil; p = p.Next {
		name := monsterDefinitionText(p.Name[:])
		if f.WriteU8(byte(len(name))) != nil {
			return 0
		}
		if _, err := f.Write([]byte(name)); err != nil {
			return 0
		}
		if f.WriteU32(p.Kind) != nil {
			return 0
		}
		if p.Kind <= 1 && f.WriteU32(p.Value) != nil {
			return 0
		}
	}
	return 1
}
func questProgressRead() int {
	questProgressReset("*:*")
	f := cryptfile.Global()
	// Keep individual read boundaries: they also define the stream checksum.
	read := func(n int) ([]byte, bool) {
		b := make([]byte, n)
		got, err := f.Read(b)
		return b, err == nil && got == n
	}
	b, ok := read(2)
	if !ok || int16(binary.LittleEndian.Uint16(b)) > 1 {
		return 0
	}
	b, ok = read(4)
	if !ok {
		return 0
	}
	count := binary.LittleEndian.Uint32(b)
	for i := uint32(0); i < count; i++ {
		b, ok = read(1)
		if !ok {
			return 0
		}
		length := int(b[0])
		b, ok = read(length)
		if !ok {
			return 0
		}
		name := questProgressCString(string(b))
		b, ok = read(4)
		if !ok {
			return 0
		}
		kind := binary.LittleEndian.Uint32(b)
		if kind > 1 {
			continue
		}
		b, ok = read(4)
		if !ok {
			return 0
		}
		value := binary.LittleEndian.Uint32(b)
		if _, ok = questProgressQualify(name); !ok {
			return 0
		}
		questProgressSet(name, value, kind)
	}
	return 1
}
