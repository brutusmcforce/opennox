package xfer

import (
	"fmt"
	"image"

	"github.com/brutusmcforce/opennox/libs/binenc"
	"github.com/brutusmcforce/opennox/libs/types"
)

// XFER data of world objects (the engine's readers are in legacy/object_xfer_types.go and object_xfer_world.go).

func init() {
	Register(&SpellPagePedestal{})
	Register(&Readable{})
	Register(&Exit{})
	Register(&Door{})
	Register(&Trigger{})
	Register(&Hole{})
	Register(&Transporter{})
	Register(&Elevator{})
	Register(&ElevatorShaft{})
	Register(&Mover{})
	Register(&Glyph{})
	Register(&InvisibleLight{})
	Register(&Sentry{})
}

// decodeStart reads the version and the object data that start most XFER data.
func decodeStart(reg ObjectRegistry, r *binenc.Reader, name string, max uint16, positive bool, obj *Object) (*reader, uint16) {
	rd := &reader{r: r}
	v := rd.version(name, max, positive)
	if rd.err == nil {
		rd.err = obj.DecodeXfer(reg, v, r)
	}
	return rd, v
}

type SpellPagePedestal struct {
	Vers   uint16
	Object Object
	Spell  uint32
	Sub    []Xfer
}

func (*SpellPagePedestal) XferType() Type        { return "SpellPagePedestalXfer" }
func (x *SpellPagePedestal) xferObject() *Object { return &x.Object }

func (x *SpellPagePedestal) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = SpellPagePedestal{}
	rd, v := decodeStart(reg, r, "pedestal", 60, false, &x.Object)
	x.Vers = v
	x.Spell = rd.u32()
	return rd.err
}

type Readable struct {
	Vers   uint16
	Object Object
	Text   string
	Sub    []Xfer
}

func (*Readable) XferType() Type        { return "ReadableXfer" }
func (x *Readable) xferObject() *Object { return &x.Object }

func (x *Readable) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Readable{}
	rd, v := decodeStart(reg, r, "readable", 60, false, &x.Object)
	x.Vers = v
	n := uint32(1) // before version 2: the length of an empty text
	if v >= 2 {
		n = rd.u32()
	}
	x.Text = rd.str(int(n))
	return rd.err
}

type Exit struct {
	Vers    uint16
	Object  Object
	MapName string
	ExitPos types.Pointf
	Sub     []Xfer
}

func (*Exit) XferType() Type        { return "ExitXfer" }
func (x *Exit) xferObject() *Object { return &x.Object }

func (x *Exit) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Exit{}
	rd, v := decodeStart(reg, r, "exit", 60, false, &x.Object)
	x.Vers = v
	if v >= 2 {
		x.MapName = rd.str(int(rd.u32()))
	} else {
		// zero-terminated
		var b []byte
		for rd.err == nil {
			c := rd.u8()
			if c == 0 {
				break
			}
			b = append(b, c)
		}
		x.MapName = string(b)
	}
	if v >= 31 {
		x.ExitPos = types.Pointf{X: rd.f32(), Y: rd.f32()}
	}
	return rd.err
}

type Door struct {
	Vers        uint16
	Object      Object
	Direction   uint32
	LockType    uint32
	Orientation uint32
	Sub         []Xfer
}

func (*Door) XferType() Type        { return "DoorXfer" }
func (x *Door) xferObject() *Object { return &x.Object }

func (x *Door) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Door{}
	rd, v := decodeStart(reg, r, "door", 60, false, &x.Object)
	x.Vers = v
	x.Direction = rd.u32()
	x.LockType = rd.u32()
	x.Orientation = x.Direction
	if v >= 41 {
		x.Orientation = rd.u32()
	}
	return rd.err
}

type Trigger struct {
	Vers             uint16
	Object           Object
	SizeX, SizeY     int32
	Colors           [6]byte // two RGB edge colors of a pressure plate
	Flags            uint32
	ScriptActivate   ScriptHandler
	ScriptDeactivate ScriptHandler
	ScriptCollide    ScriptHandler
	AllowClass       uint32
	IgnoreClass      uint32
	AllowTeam        byte
	IgnoreTeam       byte
	State            byte
	Field9           byte
	AnimFrame        uint32
	Sub              []Xfer
}

func (*Trigger) XferType() Type        { return "TriggerXfer" }
func (x *Trigger) xferObject() *Object { return &x.Object }

func (x *Trigger) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Trigger{}
	rd, v := decodeStart(reg, r, "trigger", 61, false, &x.Object)
	x.Vers = v
	x.SizeX, x.SizeY = rd.i32(), rd.i32()
	if v < 41 {
		rd.bytes(9)
	} else {
		copy(x.Colors[:], rd.bytes(6))
	}
	x.Flags = rd.u32()
	if v < 3 {
		if rd.err == nil {
			rd.err = fmt.Errorf("unsupported trigger xfer vers: %v", v)
		}
		return rd.err
	}
	x.ScriptActivate = rd.script()
	x.ScriptDeactivate = rd.script()
	if v >= 31 {
		x.ScriptCollide = rd.script()
	}
	if v < 31 {
		for i := 0; i < 4; i++ {
			rd.bytes(4 * int(rd.u8()))
		}
	}
	x.AllowClass = rd.u32()
	x.IgnoreClass = rd.u32()
	if v >= 21 {
		x.AllowTeam = rd.u8()
		x.IgnoreTeam = rd.u8()
	}
	if v >= 61 {
		x.State = rd.u8()
		x.Field9 = rd.u8()
		x.AnimFrame = rd.u32()
	}
	return rd.err
}

type Hole struct {
	Vers      uint16
	Object    Object
	Triggered uint32
	Script    ScriptHandler
	Dest      image.Point
	Frame     uint32
	Delay     uint16
	Sub       []Xfer
}

func (*Hole) XferType() Type        { return "HoleXfer" }
func (x *Hole) xferObject() *Object { return &x.Object }

func (x *Hole) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Hole{}
	rd, v := decodeStart(reg, r, "hole", 60, false, &x.Object)
	x.Vers = v
	if v >= 42 {
		x.Triggered = rd.u32()
	}
	if v < 41 {
		x.Dest = image.Point{X: int(rd.i32()), Y: int(rd.i32())}
	} else {
		x.Script = rd.script()
		x.Dest = image.Point{X: int(rd.i32()), Y: int(rd.i32())}
		x.Frame = rd.u32()
		x.Delay = rd.u16()
	}
	return rd.err
}

type Transporter struct {
	Vers       uint16
	Object     Object
	DestExtent uint32
	Sub        []Xfer
}

func (*Transporter) XferType() Type        { return "TransporterXfer" }
func (x *Transporter) xferObject() *Object { return &x.Object }

func (x *Transporter) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Transporter{}
	rd, v := decodeStart(reg, r, "transporter", 60, false, &x.Object)
	x.Vers = v
	x.DestExtent = rd.u32()
	return rd.err
}

type Elevator struct {
	Vers          uint16
	Object        Object
	PartnerExtent uint32
	Field4        uint32
	Field3        byte
	Sub           []Xfer
}

func (*Elevator) XferType() Type        { return "ElevatorXfer" }
func (x *Elevator) xferObject() *Object { return &x.Object }

func (x *Elevator) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Elevator{}
	rd, v := decodeStart(reg, r, "elevator", 61, false, &x.Object)
	x.Vers = v
	x.PartnerExtent = rd.u32()
	if v >= 41 {
		x.Field4 = rd.u32()
	}
	if v >= 61 {
		x.Field3 = rd.u8()
	}
	return rd.err
}

type ElevatorShaft struct {
	Vers          uint16
	Object        Object
	PartnerExtent uint32
	Sub           []Xfer
}

func (*ElevatorShaft) XferType() Type        { return "ElevatorShaftXfer" }
func (x *ElevatorShaft) xferObject() *Object { return &x.Object }

func (x *ElevatorShaft) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = ElevatorShaft{}
	rd, v := decodeStart(reg, r, "elevator shaft", 60, false, &x.Object)
	x.Vers = v
	x.PartnerExtent = rd.u32()
	return rd.err
}

type Mover struct {
	Vers         uint16
	Object       Object
	MovedExtent  uint32
	Field2       uint32
	Field8       uint32
	Field0       byte
	NextWaypoint uint32
	PrevWaypoint uint32
	SpeedBase    float32
	SpeedCur     float32
	Sub          []Xfer
}

func (*Mover) XferType() Type        { return "MoverXfer" }
func (x *Mover) xferObject() *Object { return &x.Object }

func (x *Mover) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Mover{}
	rd, v := decodeStart(reg, r, "mover", 60, false, &x.Object)
	x.Vers = v
	x.MovedExtent = rd.u32()
	x.Field2 = rd.u32()
	x.Field8 = rd.u32()
	if v >= 41 {
		x.Field0 = rd.u8()
		x.NextWaypoint = rd.u32()
		x.PrevWaypoint = rd.u32()
	}
	if v >= 42 {
		x.SpeedBase = rd.f32()
		x.SpeedCur = rd.f32()
	}
	return rd.err
}

type Glyph struct {
	Vers      uint16
	Object    Object
	Direction byte
	Target    types.Pointf
	Spells    []string
	Sub       []Xfer
}

func (*Glyph) XferType() Type        { return "GlyphXfer" }
func (x *Glyph) xferObject() *Object { return &x.Object }

func (x *Glyph) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Glyph{}
	rd, v := decodeStart(reg, r, "glyph", 60, false, &x.Object)
	x.Vers = v
	if v < 41 {
		rd.u32()
	}
	x.Direction = rd.u8()
	x.Target = types.Pointf{X: rd.f32(), Y: rd.f32()}
	n := int(rd.u8())
	if v < 31 {
		// spell IDs
		rd.bytes(20)
	} else {
		for i := 0; i < n && rd.err == nil; i++ {
			x.Spells = append(x.Spells, rd.str8())
		}
	}
	return rd.err
}

// InvisibleLight keeps the light record as stored (the layout of the game's drawable from its light fields on).
type InvisibleLight struct {
	Vers   uint16
	Object Object
	Light  []byte
	Extra  uint32
	Sub    []Xfer
}

func (*InvisibleLight) XferType() Type        { return "InvisibleLightXfer" }
func (x *InvisibleLight) xferObject() *Object { return &x.Object }

func (x *InvisibleLight) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = InvisibleLight{}
	rd, v := decodeStart(reg, r, "light", 60, false, &x.Object)
	x.Vers = v
	n := 36
	if v >= 2 {
		n += 97
	}
	x.Light = rd.bytes(n)
	if v >= 42 {
		x.Extra = rd.u32()
	} else if v > 40 {
		x.Extra = uint32(rd.u8())
	}
	return rd.err
}

type Sentry struct {
	Vers      uint16
	Object    Object
	BaseAngle float32
	Speed     float32
	Angle     float32
	Sub       []Xfer
}

func (*Sentry) XferType() Type        { return "SentryXfer" }
func (x *Sentry) xferObject() *Object { return &x.Object }

func (x *Sentry) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Sentry{}
	rd, v := decodeStart(reg, r, "sentry", 61, false, &x.Object)
	x.Vers = v
	x.BaseAngle = rd.f32()
	x.Speed = rd.f32()
	x.Angle = x.BaseAngle
	if v >= 61 {
		x.Angle = rd.f32()
	}
	return rd.err
}
