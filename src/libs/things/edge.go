package things

import (
	"encoding/binary"
	"errors"
	"io"
)

func (f *Reader) skipEDGE() error {
	if err := f.skip(4); err != nil {
		return err
	}
	if err := f.skipBytes8(); err != nil {
		return err
	}
	if err := f.skip(2*4 + 1); err != nil {
		return err
	}
	n1, err := f.readU8()
	if err != nil {
		return err
	}
	if err := f.skip(2 * 1); err != nil {
		return err
	}
	n2, err := f.readU8()
	if err != nil {
		return err
	}
	n3, err := f.readU8()
	if err != nil {
		return err
	}

	n := 2 * int(n1) * (int(n2) + int(n3))
	for i := 0; i < n; i++ {
		if err := f.skipImageRef(); err != nil {
			return err
		}
	}
	return f.checkEND()
}

// Edge is an edge definition: masks that blend a neighboring floor tile into a tile. Its images are as many as
// 2*Variants*(Rows+Cols).
type Edge struct {
	Name     string     `json:"name"`
	Unk1     int32      `json:"unk1,omitempty"`
	Unk2     int32      `json:"unk2,omitempty"`
	Unk3     byte       `json:"unk3,omitempty"`
	Variants byte       `json:"variants"`
	Unk4     byte       `json:"unk4,omitempty"`
	Rows     byte       `json:"rows"`
	Cols     byte       `json:"cols"`
	Images   []ImageRef `json:"images,omitempty"`
}

// ReadEdges reads the edge definitions, in the order the maps refer to them.
func (f *Reader) ReadEdges() ([]Edge, error) {
	if err := f.seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	var out []Edge
	for {
		ok, err := f.skipUntil("EDGE")
		if !ok {
			return out, err
		}
		e, err := f.readEDGE()
		if err != nil {
			return out, err
		}
		out = append(out, *e)
	}
}

func (f *Reader) readEDGE() (*Edge, error) {
	if err := f.skip(4); err != nil {
		return nil, err
	}
	var (
		e   Edge
		err error
	)
	if e.Name, err = f.readString8(); err != nil {
		return nil, err
	}
	var b [2*4 + 6]byte
	if _, err = f.read(b[:]); err != nil {
		return nil, err
	}
	e.Unk1 = int32(binary.LittleEndian.Uint32(b[0:]))
	e.Unk2 = int32(binary.LittleEndian.Uint32(b[4:]))
	e.Unk3, e.Variants, e.Unk4 = b[8], b[9], b[10]
	if b[11] == 1 {
		return nil, errors.New("unsupported edge definition")
	}
	e.Rows, e.Cols = b[12], b[13]
	n := 2 * int(e.Variants) * (int(e.Rows) + int(e.Cols))
	e.Images = make([]ImageRef, 0, n)
	for i := 0; i < n; i++ {
		ref, err := f.readImageRef()
		if err != nil {
			return nil, err
		}
		e.Images = append(e.Images, *ref)
	}
	return &e, f.checkEND()
}
