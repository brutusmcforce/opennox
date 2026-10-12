package cnxz

import (
	"fmt"
	"io"
	"sort"
)

// NXZ retains its dictionary, symbol order and frequencies across byte-aligned
// blocks. The file header supplies the total uncompressed size.
type mapDecoder struct {
	source  []byte
	pos     int
	bits    byte
	current byte
	symbols [274]int
	counts  [274]int16
	groups  [16]struct{ bits, offset int }
	window  [65536]byte
	written int
}

func newMapDecoder(src []byte) *mapDecoder {
	d := &mapDecoder{source: src}
	sizes := [16]int{2, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 5, 5, 5}
	off := 0
	for i, n := range sizes {
		d.groups[i].bits = n
		d.groups[i].offset = off
		off += 1 << n
	}
	pos := 0
	for i := 256; i < 272; i++ {
		d.symbols[pos] = i
		pos++
	}
	for _, i := range []int{0, 32, 48, 255} {
		d.symbols[pos] = i
		pos++
	}
	for i := 1; i < 274; i++ {
		if i == 32 || i == 48 || i >= 255 && i < 272 {
			continue
		}
		d.symbols[pos] = i
		pos++
	}
	return d
}
func (d *mapDecoder) read(n int) (uint32, error) {
	if n < 0 || n > 24 {
		return 0, fmt.Errorf("nxz: invalid bit count %d", n)
	}
	var v uint32
	for n > 0 {
		if d.bits == 0 {
			if d.pos == len(d.source) {
				return 0, io.ErrUnexpectedEOF
			}
			d.current = d.source[d.pos]
			d.pos++
			d.bits = 8
		}
		take := min(n, int(d.bits))
		v = v<<take | uint32(d.current>>(8-take))
		d.current <<= take
		d.bits -= byte(take)
		n -= take
	}
	return v, nil
}
func (d *mapDecoder) rebuild() error {
	for i := range d.symbols {
		d.symbols[i] = i
	}
	sort.Slice(d.symbols[:], func(i, j int) bool {
		a, b := d.symbols[i], d.symbols[j]
		if d.counts[a] != d.counts[b] {
			return d.counts[a] > d.counts[b]
		}
		return a > b
	})
	for i := range d.counts {
		d.counts[i] >>= 1
	}
	bits, off := 0, 0
	for i := range d.groups {
		for {
			b, err := d.read(1)
			if err != nil {
				return err
			}
			if b != 0 {
				break
			}
			bits++
			if bits > 24 {
				return fmt.Errorf("nxz: invalid symbol table")
			}
		}
		d.groups[i].bits = bits
		d.groups[i].offset = off
		off += 1 << bits
	}
	return nil
}
func (d *mapDecoder) decode(dst []byte) error {
	put := func(v byte) { dst[d.written] = v; d.window[d.written&65535] = v; d.written++ }
	for d.written < len(dst) {
		prefix, err := d.read(4)
		if err != nil {
			return err
		}
		group := d.groups[prefix]
		tail, err := d.read(group.bits)
		if err != nil {
			return err
		}
		index := group.offset + int(tail)
		if index >= len(d.symbols) {
			return fmt.Errorf("nxz: invalid symbol index %d", index)
		}
		sym := d.symbols[index]
		d.counts[sym]++
		switch {
		case sym < 256:
			put(byte(sym))
		case sym == 272:
			if err := d.rebuild(); err != nil {
				return err
			}
		case sym == 273:
			d.bits = 0
			d.current = 0 // Next compressed block starts at the next byte.
		default:
			length := sym - 256 + 4
			if sym >= 264 {
				n := sym - 264 + 1
				extra, err := d.read(n)
				if err != nil {
					return err
				}
				length = 4 + 6 + (1 << n) + int(extra)
			}
			prefix, err := d.read(3)
			if err != nil {
				return err
			}
			n, base := 0, int(prefix)
			if prefix > 1 {
				n = int(prefix) - 1
				base = 1 << n
			}
			extra, err := d.read(n + 9)
			if err != nil {
				return err
			}
			distance := (base << 9) + int(extra)
			if distance == 0 {
				return fmt.Errorf("nxz: zero copy distance")
			}
			if length > len(dst)-d.written {
				return fmt.Errorf("nxz: copy exceeds declared output size")
			}
			for i := 0; i < length; i++ {
				put(d.window[(d.written-distance)&65535])
			}
		}
	}
	return nil
}
