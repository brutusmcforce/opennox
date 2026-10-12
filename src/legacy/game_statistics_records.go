package legacy

import (
	"bytes"
	"encoding/binary"
	"math/bits"
)

// statisticsRecord owns a report field. Records never cross the C boundary.
type statisticsRecord struct {
	tag          [4]byte
	kind, length uint16
	data         []byte
	next         *statisticsRecord
}
type statisticsRecords struct {
	kind uint16
	head *statisticsRecord
}

func (r *statisticsRecord) set(kind uint16, tag string, value int32, data []byte) {
	*r = statisticsRecord{kind: kind}
	if i := bytes.IndexByte([]byte(tag), 0); i >= 0 {
		tag = tag[:i]
	}
	copy(r.tag[:], tag)
	switch kind {
	case 2:
		r.data = []byte{byte(value)}
	case 3:
		r.data = make([]byte, 2)
		binary.LittleEndian.PutUint16(r.data, uint16(int16(int8(value))))
	case 6:
		r.data = make([]byte, 4)
		binary.LittleEndian.PutUint32(r.data, uint32(int32(int8(value))))
	case 7:
		if i := bytes.IndexByte(data, 0); i >= 0 {
			data = data[:i]
		}
		r.data = make([]byte, int(uint16(len(data)+1)))
		copy(r.data, data)
	case 20:
		r.data = append([]byte{}, data[:int(uint16(len(data)))]...)
	default:
		panic("unsupported statistics field")
	}
	r.length = uint16(len(r.data))
}
func (p *statisticsRecords) prepend(r *statisticsRecord) { r.next = p.head; p.head = r }
func (p *statisticsRecords) add(kind uint16, tag string, value int32, data []byte) *statisticsRecord {
	r := new(statisticsRecord)
	r.set(kind, tag, value, data)
	p.prepend(r)
	return r
}
func (r *statisticsRecord) swapData() uint16 {
	switch r.kind {
	case 3, 4:
		v := bits.ReverseBytes16(binary.LittleEndian.Uint16(r.data))
		binary.LittleEndian.PutUint16(r.data, v)
		return v
	case 5, 6:
		v := bits.ReverseBytes32(binary.LittleEndian.Uint32(r.data))
		binary.LittleEndian.PutUint32(r.data, v)
		return uint16(v)
	default:
		return r.kind - 3
	}
}
func (r *statisticsRecord) encode() uint16 {
	r.swapData()
	r.kind = bits.ReverseBytes16(r.kind)
	r.length = bits.ReverseBytes16(r.length)
	return r.length
}
func (r *statisticsRecord) decode() uint16 {
	r.kind = bits.ReverseBytes16(r.kind)
	r.length = bits.ReverseBytes16(r.length)
	return r.swapData()
}
func (p *statisticsRecords) serialize() []byte {
	n := 4
	for r := p.head; r != nil; r = r.next {
		n += 8 + (int(r.length)+3)&^3
	}
	out := make([]byte, n)
	binary.BigEndian.PutUint16(out, uint16(n))
	binary.BigEndian.PutUint16(out[2:], p.kind)
	pos := 4
	for r := p.head; r != nil; r = r.next {
		size := int(r.length)
		r.encode()
		copy(out[pos:], r.tag[:])
		binary.LittleEndian.PutUint16(out[pos+4:], r.kind)
		binary.LittleEndian.PutUint16(out[pos+6:], r.length)
		copy(out[pos+8:], r.data[:size])
		pos += 8 + (size+3)&^3
		r.decode()
	}
	return out
}
