package xfer

import (
	"fmt"
	"io"
	"reflect"

	"github.com/brutusmcforce/opennox/libs/binenc"
)

// Type is a decoding type for XFER data.
type Type string

// Xfer is type that can be encoded or decoded to/from Nox XFER data.
type Xfer interface {
	XferType() Type
	DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error
}

// ObjectRegistry is an interface for looking up XFER Type based on object type name or ID.
type ObjectRegistry interface {
	// XferByObjectType finds XFER Type for a given object type name.
	XferByObjectType(typ string) Type
	// XferByObjectTypeID finds XFER Type for a given object type ID.
	XferByObjectTypeID(id int) Type
}

var byType = make(map[Type]reflect.Type)

// Register a new XFER data type.
func Register(x Xfer) {
	typ := x.XferType()
	if _, ok := byType[typ]; ok {
		panic("already registered")
	}
	byType[typ] = reflect.TypeOf(x).Elem()
}

// Decode XFER data with a specific Type name.
func Decode(reg ObjectRegistry, xfer Type, r *binenc.Reader) (Xfer, error) {
	return decodeTyped(reg, xfer, "", r)
}

// typedXfer is implemented by XFER types whose data depends on the object type (see ClassRegistry).
type typedXfer interface {
	decodeTyped(reg ObjectRegistry, typ string, r *binenc.Reader) error
}

// objectXfer is implemented by XFER types with the common object data.
type objectXfer interface {
	xferObject() *Object
}

func decodeTyped(reg ObjectRegistry, xfer Type, typ string, r *binenc.Reader) (Xfer, error) {
	if reg == nil {
		reg = DefaultRegistry
	}
	rt, ok := byType[xfer]
	if !ok {
		return nil, fmt.Errorf("unsupported xfer: %q", xfer)
	}
	x := reflect.New(rt).Interface().(Xfer)
	var err error
	if tx, ok := x.(typedXfer); ok {
		err = tx.decodeTyped(reg, typ, r)
	} else {
		err = x.DecodeXfer(reg, r)
	}
	if ox, ok := x.(objectXfer); ok && typ != "" {
		ox.xferObject().Type = typ
	}
	return x, err
}

// DecodeByObjectType is similar to Decode, but accepts object type name, instead of XFER Type.
func DecodeByObjectType(reg ObjectRegistry, typ string, r *binenc.Reader) (Xfer, error) {
	if reg == nil {
		reg = DefaultRegistry
	}
	xfer := reg.XferByObjectType(typ)
	if xfer == "" {
		return nil, fmt.Errorf("cannot find xfer for object type: %q", typ)
	}
	return decodeTyped(reg, xfer, typ, r)
}

// DecodeByObjectTypeID is similar to Decode, but accepts object type ID, instead of XFER Type.
func DecodeByObjectTypeID(reg ObjectRegistry, id int, r *binenc.Reader) (Xfer, error) {
	if reg == nil {
		reg = DefaultRegistry
	}
	xfer := reg.XferByObjectTypeID(id)
	if xfer == "" {
		return nil, fmt.Errorf("cannot find xfer for object type id: %d", id)
	}
	typ := ""
	if tr, ok := reg.(TypeNameRegistry); ok {
		typ = tr.ObjectTypeByID(id)
	}
	return decodeTyped(reg, xfer, typ, r)
}

// reader reads XFER data, keeping the first error.
type reader struct {
	r   *binenc.Reader
	err error
}

func (x *reader) fail() {
	if x.err == nil {
		x.err = io.ErrUnexpectedEOF
	}
}

func (x *reader) u8() byte {
	if x.err != nil {
		return 0
	}
	v, ok := x.r.ReadU8()
	if !ok {
		x.fail()
	}
	return v
}

func (x *reader) u16() uint16 {
	if x.err != nil {
		return 0
	}
	v, ok := x.r.ReadU16()
	if !ok {
		x.fail()
	}
	return v
}

func (x *reader) u32() uint32 {
	if x.err != nil {
		return 0
	}
	v, ok := x.r.ReadU32()
	if !ok {
		x.fail()
	}
	return v
}

func (x *reader) i32() int32 { return int32(x.u32()) }

func (x *reader) f32() float32 {
	if x.err != nil {
		return 0
	}
	v, ok := x.r.ReadF32()
	if !ok {
		x.fail()
	}
	return v
}

func (x *reader) bytes(n int) []byte {
	if x.err != nil {
		return nil
	}
	v, ok := x.r.ReadBytes(n)
	if !ok {
		x.fail()
	}
	return v
}

// str8 reads a string with a byte length.
func (x *reader) str8() string {
	if x.err != nil {
		return ""
	}
	v, ok := x.r.ReadString8()
	if !ok {
		x.fail()
	}
	return v
}

// str reads a string of n bytes (up to a zero byte).
func (x *reader) str(n int) string {
	if x.err != nil {
		return ""
	}
	v, ok := x.r.ReadString(n)
	if !ok {
		x.fail()
	}
	return v
}

func (x *reader) script() ScriptHandler {
	var h ScriptHandler
	if x.err != nil {
		return h
	}
	x.err = h.DecodeXfer(nil, x.r)
	return h
}

// version reads an XFER version that must be in 0..max (or 1..max with positive).
func (x *reader) version(name string, max uint16, positive bool) uint16 {
	v := x.u16()
	if x.err == nil && (int16(v) > int16(max) || positive && int16(v) <= 0) {
		x.err = fmt.Errorf("unsupported %s xfer vers: %v", name, int16(v))
	}
	return v
}
