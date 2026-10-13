package maps

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/brutusmcforce/opennox/libs/binenc"
	"github.com/brutusmcforce/opennox/libs/xfer"
)

func init() {
	RegisterSection(&ObjectsTOC{})
	RegisterSection(&Objects{})
}

type ObjectTOC struct {
	Ind  uint16
	Type string
}

func (w *ObjectTOC) MarshalBinary() ([]byte, error) {
	if len(w.Type) > 0xff {
		return nil, fmt.Errorf("toc type name too long: %q", w.Type)
	}
	data := make([]byte, 0, 2+1+len(w.Type))
	data = binary.LittleEndian.AppendUint16(data, w.Ind)
	data = append(data, byte(len(w.Type)))
	data = append(data, w.Type...)
	return data, nil
}

func (w *ObjectTOC) Decode(r *binenc.Reader) error {
	var ok bool
	w.Ind, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Type, ok = r.ReadString8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type ObjectsTOC struct {
	Vers uint16
	TOC  []ObjectTOC
}

func (*ObjectsTOC) MapSection() string {
	return "ObjectTOC"
}

func (sect *ObjectsTOC) MarshalBinary() ([]byte, error) {
	if len(sect.TOC) > math.MaxUint16 {
		return nil, errors.New("too many map objects")
	}
	var data []byte
	data = binary.LittleEndian.AppendUint16(data, sect.Vers)
	data = binary.LittleEndian.AppendUint16(data, uint16(len(sect.TOC)))
	for _, t := range sect.TOC {
		b, err := t.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	return data, nil
}

func (sect *ObjectsTOC) Decode(r *binenc.Reader) error {
	*sect = ObjectsTOC{}
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers > 1 {
		return fmt.Errorf("unsupported object toc version: %d", vers)
	}
	sect.Vers = vers

	cnt, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sect.TOC = make([]ObjectTOC, 0, cnt)
	for i := 0; i < int(cnt); i++ {
		var t ObjectTOC
		if err := t.Decode(r); err != nil {
			return err
		}
		sect.TOC = append(sect.TOC, t)
	}
	return nil
}

func (sect *ObjectsTOC) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

type Objects struct {
	Vers uint16
	Data []byte
}

func (*Objects) MapSection() string {
	return "ObjectData"
}

func (sect *Objects) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, 2+len(sect.Data))
	data = binary.LittleEndian.AppendUint16(data, sect.Vers)
	data = append(data, sect.Data...)
	return data, nil
}

func (sect *Objects) Decode(r *binenc.Reader) error {
	*sect = Objects{}
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers > 1 {
		return fmt.Errorf("unsupported object data version: %d", vers)
	}
	sect.Vers = vers

	sect.Data = r.ReadAllBytes()
	return nil
}

func (sect *Objects) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

type Xfer struct {
	Type string
	Xfer xfer.Xfer
}

func (sect *Objects) ReadObjects(toc *ObjectsTOC, reg xfer.ObjectRegistry) ([]Xfer, error) {
	tmap := make(map[uint16]string, len(toc.TOC))
	for _, t := range toc.TOC {
		if t.Ind == 0 {
			return nil, errors.New("object TOC index must not be zero")
		}
		if t2, ok := tmap[t.Ind]; ok {
			return nil, fmt.Errorf("ambiguous object TOC entry %d: %q vs %q", t.Ind, t2, t.Type)
		}
		tmap[t.Ind] = t.Type
	}
	if reg == nil {
		reg = xfer.DefaultRegistry
	}
	// child objects (inventories) refer to the TOC too
	treg := &xfer.TOCRegistry{Base: reg, TOC: make(map[int]string, len(tmap))}
	for ind, typ := range tmap {
		treg.TOC[int(ind)] = typ
	}
	r := binenc.NewReader(sect.Data)
	var out []Xfer
	for {
		ind, ok := r.ReadU16()
		if !ok || ind == 0 {
			return out, nil
		}
		if !r.Align(+2) {
			return out, io.ErrUnexpectedEOF
		}
		sz, ok := r.ReadU64()
		if !ok {
			return out, io.ErrUnexpectedEOF
		}
		data, ok := r.ReadNext(int(sz))
		if !ok {
			return out, io.ErrUnexpectedEOF
		}
		typ, ok := tmap[ind]
		if !ok {
			return out, fmt.Errorf("object TOC entry is missing: %d", ind)
		}
		r2 := binenc.NewReader(data)
		x, err := xfer.DecodeByObjectType(treg, typ, r2)
		if err != nil {
			return out, fmt.Errorf("object %d (%s): %w", len(out), typ, err)
		}
		out = append(out, Xfer{
			Type: typ,
			Xfer: x,
		})
		if r2.Remaining() != 0 {
			return out, fmt.Errorf("partial decoding of xfer %s(%T): %d bytes left", typ, x, r2.Remaining())
		}
		// the inventory follows as more records
		if err := xfer.ReadChildren(treg, x, r, +2); err != nil {
			return out, fmt.Errorf("object %d (%s): %w", len(out)-1, typ, err)
		}
	}
}

// ReadObjects decodes the map's objects, if the map was read with its sections kept as unknown sections (see
// Reader.ReadSections). A nil registry uses the default one.
func (m *Map) ReadObjects(reg xfer.ObjectRegistry) ([]Xfer, error) {
	var (
		toc  ObjectsTOC
		objs Objects
		ok   bool
	)
	for _, s := range m.Unknown {
		switch s.Name {
		case toc.MapSection():
			if err := toc.UnmarshalBinary(s.Data); err != nil {
				return nil, err
			}
		case objs.MapSection():
			if err := objs.UnmarshalBinary(s.Data); err != nil {
				return nil, err
			}
			ok = true
		}
	}
	if !ok {
		return nil, nil
	}
	return objs.ReadObjects(&toc, reg)
}
