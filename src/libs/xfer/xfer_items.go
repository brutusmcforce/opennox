package xfer

import (
	"github.com/brutusmcforce/opennox/libs/binenc"
)

// XFER data of items and other objects (the engine's readers are in legacy/item_xfer_*.go).

func init() {
	Register(&Ammo{})
	Register(&Team{})
	Register(&SpellReward{})
	Register(&AbilityReward{})
	Register(&FieldGuide{})
	Register(&Gold{})
	Register(&ToxicCloud{})
	Register(&Obelisk{})
	Register(&MonsterGenerator{})
	Register(&RewardMarker{})
}

func (x *reader) modifiers() []*Modifier {
	var out []*Modifier
	for i := 0; i < 4; i++ {
		out = append(out, &Modifier{Name: x.str8()})
	}
	return out
}

type Ammo struct {
	Vers       uint16
	Object     Object
	Modifiers  []*Modifier
	Charges    byte
	MaxCharges byte
	Sub        []Xfer
}

func (*Ammo) XferType() Type        { return "AmmoXfer" }
func (x *Ammo) xferObject() *Object { return &x.Object }

func (x *Ammo) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Ammo{}
	rd, v := decodeStart(reg, r, "ammo", 60, false, &x.Object)
	x.Vers = v
	x.Modifiers = rd.modifiers()
	x.Charges = rd.u8()
	x.MaxCharges = rd.u8()
	return rd.err
}

// Team is a team item (a flag or a crown).
type Team struct {
	Vers      uint16
	Object    Object
	Modifiers []*Modifier
	Sub       []Xfer
}

func (*Team) XferType() Type        { return "TeamXfer" }
func (x *Team) xferObject() *Object { return &x.Object }

func (x *Team) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Team{}
	rd, v := decodeStart(reg, r, "team", 60, false, &x.Object)
	x.Vers = v
	x.Modifiers = rd.modifiers()
	return rd.err
}

// SpellReward is a spell book. Old versions stored spell IDs, kept in SpellID.
type SpellReward struct {
	Vers    uint16
	Object  Object
	Spell   string
	SpellID byte
	Sub     []Xfer
}

func (*SpellReward) XferType() Type        { return "SpellRewardXfer" }
func (x *SpellReward) xferObject() *Object { return &x.Object }

func (x *SpellReward) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = SpellReward{}
	rd, v := decodeStart(reg, r, "spell reward", 60, false, &x.Object)
	x.Vers = v
	switch {
	case v < 31:
		rd.u8()
		a, b := rd.u8(), rd.u8()
		x.SpellID = a
		if b != 0 {
			x.SpellID = b
		}
		if v == 10 {
			rd.u8()
		}
	case v >= 41:
		x.Spell = rd.str8()
	default:
		rd.str8()
		a, b := rd.str8(), rd.str8()
		x.Spell = a
		if b != "" {
			x.Spell = b
		}
	}
	return rd.err
}

type AbilityReward struct {
	Vers    uint16
	Object  Object
	Ability string
	Sub     []Xfer
}

func (*AbilityReward) XferType() Type        { return "AbilityRewardXfer" }
func (x *AbilityReward) xferObject() *Object { return &x.Object }

func (x *AbilityReward) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = AbilityReward{}
	rd, v := decodeStart(reg, r, "ability reward", 61, false, &x.Object)
	x.Vers = v
	x.Ability = rd.str8()
	return rd.err
}

type FieldGuide struct {
	Vers     uint16
	Object   Object
	Creature string
	Sub      []Xfer
}

func (*FieldGuide) XferType() Type        { return "FieldGuideXfer" }
func (x *FieldGuide) xferObject() *Object { return &x.Object }

func (x *FieldGuide) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = FieldGuide{}
	rd, v := decodeStart(reg, r, "field guide", 60, false, &x.Object)
	x.Vers = v
	x.Creature = rd.str8()
	return rd.err
}

type Gold struct {
	Vers   uint16
	Object Object
	Amount uint32
	Sub    []Xfer
}

func (*Gold) XferType() Type        { return "GoldXfer" }
func (x *Gold) xferObject() *Object { return &x.Object }

func (x *Gold) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Gold{}
	rd, v := decodeStart(reg, r, "gold", 60, false, &x.Object)
	x.Vers = v
	x.Amount = rd.u32()
	return rd.err
}

type ToxicCloud struct {
	Vers   uint16
	Object Object
	Life   uint32
	Sub    []Xfer
}

func (*ToxicCloud) XferType() Type        { return "ToxicCloudXfer" }
func (x *ToxicCloud) xferObject() *Object { return &x.Object }

func (x *ToxicCloud) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = ToxicCloud{}
	rd, v := decodeStart(reg, r, "toxic cloud", 61, true, &x.Object)
	x.Vers = v
	x.Life = rd.u32()
	return rd.err
}

type Obelisk struct {
	Vers    uint16
	Object  Object
	Mana    uint32
	Minimap byte
	Sub     []Xfer
}

func (*Obelisk) XferType() Type        { return "ObeliskXfer" }
func (x *Obelisk) xferObject() *Object { return &x.Object }

func (x *Obelisk) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Obelisk{}
	rd, v := decodeStart(reg, r, "obelisk", 61, false, &x.Object)
	x.Vers = v
	if v >= 61 {
		x.Mana = rd.u32()
		x.Minimap = rd.u8()
	}
	return rd.err
}

// MonsterGenerator holds the monsters it makes (up to 4 for each of 3 difficulty levels), with their XFER data.
type MonsterGenerator struct {
	Vers         uint16
	Object       Object
	Rates        []byte
	Current      byte
	Max          byte
	Frame        uint32
	Scripts      [4]ScriptHandler
	Monsters     [3][]Xfer
	MonsterTypes [3][]string
	Field83      []byte
	Flags        uint32
	Sub          []Xfer
}

func (*MonsterGenerator) XferType() Type        { return "MonsterGeneratorXfer" }
func (x *MonsterGenerator) xferObject() *Object { return &x.Object }

func (x *MonsterGenerator) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = MonsterGenerator{}
	rd, v := decodeStart(reg, r, "monster generator", 63, true, &x.Object)
	x.Vers = v
	x.Rates = rd.bytes(int(rd.u8()))
	x.Current = rd.u8()
	x.Max = rd.u8()
	x.Frame = rd.u32()
	for i := range x.Scripts {
		x.Scripts[i] = rd.script()
	}
	rows := int(rd.u8())
	for row := 0; row < rows && rd.err == nil; row++ {
		n := int(rd.u8())
		for i := 0; i < n && rd.err == nil; i++ {
			typ := rd.str8()
			rd.u16()
			if rd.err != nil {
				break
			}
			data, err := readSized(r, 0)
			if err != nil {
				return err
			}
			r2 := binenc.NewReader(data)
			m, err := DecodeByObjectType(reg, typ, r2)
			if err == nil && r2.Remaining() != 0 {
				err = errPartial(m, r2)
			}
			if err != nil {
				return err
			}
			if row < len(x.Monsters) {
				x.Monsters[row] = append(x.Monsters[row], m)
				x.MonsterTypes[row] = append(x.MonsterTypes[row], typ)
			}
		}
	}
	if v >= 62 {
		x.Field83 = rd.bytes(int(rd.u8()))
	}
	if v >= 63 {
		x.Flags = rd.u32()
	}
	return rd.err
}

// RewardMarker marks where random rewards appear.
type RewardMarker struct {
	Vers       uint16
	Object     Object
	Categories uint32
	Explicit   uint32
	Spells     []string
	Abilities  []string
	Guides     []string
	Fields     [5]uint32
	Chance     uint32
	Flags      byte
	Sub        []Xfer
}

func (*RewardMarker) XferType() Type        { return "RewardMarkerXfer" }
func (x *RewardMarker) xferObject() *Object { return &x.Object }

func (x *RewardMarker) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = RewardMarker{}
	rd, v := decodeStart(reg, r, "reward marker", 63, true, &x.Object)
	x.Vers = v
	x.Categories = rd.u32()
	x.Explicit = rd.u32()
	for _, list := range []*[]string{&x.Spells, &x.Abilities, &x.Guides} {
		n := int(rd.u16())
		for i := 0; i < n && rd.err == nil; i++ {
			*list = append(*list, rd.str8())
		}
	}
	for i := range x.Fields {
		x.Fields[i] = rd.u32()
	}
	if v >= 62 {
		x.Chance = rd.u32()
	}
	if v >= 63 {
		x.Flags = rd.u8()
	}
	return rd.err
}
