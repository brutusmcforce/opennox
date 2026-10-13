package xfer

import (
	"fmt"
	"io"
	"reflect"

	"github.com/brutusmcforce/opennox/libs/binenc"
	"github.com/brutusmcforce/opennox/libs/types"
)

const DefaultType = Type("DefaultXfer")

func init() {
	Register(&Default{})
}

type ScriptHandler struct {
	Vers uint16
	Func string
	Val2 uint32
}

func (x *ScriptHandler) DecodeXfer(_ ObjectRegistry, r *binenc.Reader) error {
	*x = ScriptHandler{}
	var ok bool
	x.Vers, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if x.Vers > 1 {
		return fmt.Errorf("unsupported script handler xfer vers: %v", x.Vers)
	}
	x.Func, ok = r.ReadString32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Val2, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// Object is the data common to all objects.
type Object struct {
	// Type is the object type name, if known (it is not a part of the data).
	Type      string
	Vers      uint16
	Extent    uint32
	ID        uint32 // script ID
	Pos       types.Pointf
	Val5      byte // 0 if only the extent, ID and position are stored
	Flags     uint32
	Name      string
	Team      byte
	SubN      byte     // number of inventory items
	Owned     []uint32 // script IDs of owned objects
	Anim      uint32   // extra status flags
	Handler11 *ScriptHandler
	DeadFrame int32
}

func (x *Object) xferObject() *Object { return x }

// ObjectOf returns the data common to all objects of x, or nil if it has none.
func ObjectOf(x Xfer) *Object {
	if ox, ok := x.(interface{ xferObject() *Object }); ok {
		return ox.xferObject()
	}
	return nil
}

func (x *Object) DecodeXfer(reg ObjectRegistry, gvers uint16, r *binenc.Reader) error {
	*x = Object{}
	var ok bool
	if gvers >= 40 {
		x.Vers, ok = r.ReadU16()
		if !ok {
			return io.ErrUnexpectedEOF
		} else if x.Vers > 64 {
			return fmt.Errorf("unsupported object xfer vers: %v", x.Vers)
		}
	}
	if gvers < 40 || x.Vers < 61 {
		return x.decodeOld(gvers, r)
	}
	x.Extent, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.ID, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Pos.X, ok = r.ReadF32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Pos.Y, ok = r.ReadF32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Val5, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	if x.Val5 == 0 {
		return nil
	}
	x.Flags, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Name, ok = r.ReadString8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Team, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.SubN, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sz, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Owned = make([]uint32, sz)
	for i := range x.Owned {
		x.Owned[i], ok = r.ReadU32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	}
	x.Anim, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	if x.Vers >= 63 {
		x.Handler11 = new(ScriptHandler)
		err := x.Handler11.DecodeXfer(reg, r)
		if err != nil {
			return err
		}
		if x.Vers >= 64 {
			x.DeadFrame, ok = r.ReadI32()
			if !ok {
				return io.ErrUnexpectedEOF
			}
		}
	}
	return nil
}

// decodeOld decodes the object data of the versions before 61.
func (x *Object) decodeOld(gvers uint16, r *binenc.Reader) error {
	rd := &reader{r: r}
	x.Val5 = 1
	x.Extent = rd.u32()
	x.Flags = rd.u32()
	if gvers < 40 || x.Vers < 4 {
		px, py := rd.i32(), rd.i32()
		x.Pos = types.Pointf{X: float32(px), Y: float32(py)}
	} else {
		x.Pos = types.Pointf{X: rd.f32(), Y: rd.f32()}
	}
	if gvers >= 10 {
		x.Name = rd.str8()
	}
	if gvers >= 20 {
		x.Team = rd.u8()
	}
	if gvers >= 30 {
		x.SubN = rd.u8()
	}
	if gvers < 40 {
		return rd.err
	}
	x.ID = rd.u32()
	if x.Vers >= 2 {
		var n int
		if x.Vers < 5 {
			n = int(uint16(rd.u32()))
		} else {
			n = int(rd.u16())
		}
		if rd.err == nil {
			x.Owned = make([]uint32, n)
			for i := range x.Owned {
				x.Owned[i] = rd.u32()
			}
		}
	}
	if x.Vers >= 3 {
		x.Anim = rd.u32()
	}
	return rd.err
}

// ReadChildren reads the child objects of x (its inventory), which follow its data in a stream of object records,
// into its Sub field, with their children. Each is a type (a name before version 60 of x, from it an index, see
// TOCRegistry), then its data, with the length at the next 8-byte boundary. align is the offset of r's start from
// an 8-byte boundary of the stream.
func ReadChildren(reg ObjectRegistry, x Xfer, r *binenc.Reader, align int) error {
	ox, ok := x.(objectXfer)
	if !ok {
		return nil
	}
	n := int(ox.xferObject().SubN)
	if n == 0 {
		return nil
	}
	rv := reflect.ValueOf(x).Elem()
	vers := uint16(rv.FieldByName("Vers").Uint())
	sub := rv.FieldByName("Sub")
	if !sub.IsValid() {
		return fmt.Errorf("xfer %T has no inventory", x)
	}
	for i := 0; i < n; i++ {
		var (
			typ string
			id  = -1
		)
		if vers < 60 {
			name, ok := r.ReadString8()
			if !ok {
				return io.ErrUnexpectedEOF
			}
			typ = name
		} else {
			v, ok := r.ReadU16()
			if !ok {
				return io.ErrUnexpectedEOF
			}
			id = int(v)
		}
		data, err := readSized(r, align)
		if err != nil {
			return err
		}
		r2 := binenc.NewReader(data)
		var c Xfer
		if id >= 0 {
			c, err = DecodeByObjectTypeID(reg, id, r2)
		} else {
			c, err = DecodeByObjectType(reg, typ, r2)
		}
		if err == nil && r2.Remaining() != 0 {
			err = errPartial(c, r2)
		}
		if err != nil {
			return fmt.Errorf("inventory item %d: %w", i, err)
		}
		sub.Set(reflect.Append(sub, reflect.ValueOf(c)))
		if err := ReadChildren(reg, c, r, align); err != nil {
			return err
		}
	}
	return nil
}

// readSized reads the length at the next 8-byte boundary and data of that length.
func readSized(r *binenc.Reader, align int) ([]byte, error) {
	if !r.Align(align) {
		return nil, io.ErrUnexpectedEOF
	}
	sz, ok := r.ReadU64()
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	data, ok := r.ReadNext(int(sz))
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	return data, nil
}

func errPartial(x Xfer, r *binenc.Reader) error {
	typ := ""
	if ox, ok := x.(objectXfer); ok {
		typ = ox.xferObject().Type
	}
	return fmt.Errorf("partial decoding of xfer %s(%T): %d bytes left", typ, x, r.Remaining())
}

type Default struct {
	Vers   uint16
	Object Object
	Sub    []Xfer
}

func (x *Default) xferObject() *Object { return &x.Object }

func (*Default) XferType() Type {
	return DefaultType
}

func (x *Default) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Default{}
	var ok bool
	x.Vers, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if x.Vers > 60 {
		return fmt.Errorf("unsupported default xfer vers: %v", x.Vers)
	}
	err := x.Object.DecodeXfer(reg, x.Vers, r)
	if err != nil {
		return err
	}
	return err
}
