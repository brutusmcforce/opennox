package things

import (
	"encoding/binary"
	"io"
)

func (f *Reader) skipFLOR() error {
	if err := f.skip(4); err != nil {
		return err
	}
	if err := f.skipBytes8(); err != nil {
		return err
	}
	if err := f.skip(3*1 + 2*4 + 1); err != nil {
		return err
	}
	n1, err := f.readU8()
	if err != nil {
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
	if err := f.skip(1); err != nil {
		return err
	}
	n := int(n1) * int(n2) * int(n3)
	for i := 0; i < n; i++ {
		if err := f.skipImageRef(); err != nil {
			return err
		}
	}
	return f.checkEND()
}

// Floor is a floor tile definition: its images are all variants of the tile, as many as Rows*Cols*Frames.
type Floor struct {
	Name   string     `json:"name"`
	Color  [3]byte    `json:"color"` // the color on the minimap
	Unk1   int32      `json:"unk1,omitempty"`
	Unk2   int32      `json:"unk2,omitempty"`
	Unk3   byte       `json:"unk3,omitempty"`
	Rows   byte       `json:"rows"`
	Cols   byte       `json:"cols"`
	Frames byte       `json:"frames"`
	Unk4   byte       `json:"unk4,omitempty"`
	Images []ImageRef `json:"images,omitempty"`
}

// ReadFloors reads the floor tile definitions, in the order the maps refer to them.
func (f *Reader) ReadFloors() ([]Floor, error) {
	if err := f.seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	var out []Floor
	for {
		ok, err := f.skipUntil("FLOR")
		if !ok {
			return out, err
		}
		fl, err := f.readFLOR()
		if err != nil {
			return out, err
		}
		out = append(out, *fl)
	}
}

func (f *Reader) readFLOR() (*Floor, error) {
	if err := f.skip(4); err != nil {
		return nil, err
	}
	var (
		fl  Floor
		err error
	)
	if fl.Name, err = f.readString8(); err != nil {
		return nil, err
	}
	var b [3 + 2*4 + 5]byte
	if _, err = f.read(b[:]); err != nil {
		return nil, err
	}
	copy(fl.Color[:], b[0:3])
	fl.Unk1 = int32(binary.LittleEndian.Uint32(b[3:]))
	fl.Unk2 = int32(binary.LittleEndian.Uint32(b[7:]))
	fl.Unk3 = b[11]
	fl.Rows, fl.Cols, fl.Frames, fl.Unk4 = b[12], b[13], b[14], b[15]
	n := int(fl.Rows) * int(fl.Cols) * int(fl.Frames)
	fl.Images = make([]ImageRef, 0, n)
	for i := 0; i < n; i++ {
		ref, err := f.readImageRef()
		if err != nil {
			return nil, err
		}
		fl.Images = append(fl.Images, *ref)
	}
	return &fl, f.checkEND()
}
