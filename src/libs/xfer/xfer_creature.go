package xfer

import (
	"fmt"

	"github.com/brutusmcforce/opennox/libs/binenc"
)

// XFER data of monsters and NPCs (the engine's readers are in legacy/creature_xfer_*.go). Field names follow the
// engine where it names them; the numbers are offsets in the game's monster data.

func init() {
	Register(&Monster{})
	Register(&NPC{})
}

// Creature is the data monsters and NPCs share.
type Creature struct {
	Direction [2]uint32 // a direction vector
	// Scripts: enemy sighted, death, looking for enemy, change focus, is hit, retreat, collision, enemy heard,
	// end of waypoint, lost enemy.
	Scripts    [10]*ScriptHandler
	Field305   uint16
	HeaderWord uint32

	// stats (names and offsets of the engine's MonsterUpdateData)
	Field333      byte   // 1332, roam flags?
	StatusFlags   uint32 // 1440
	Field338      float32
	RetreatLevel  float32
	ResumeLevel   float32
	SightRange    float32
	Health        uint16 // current health (monsters: from version 49)
	Aggression    float32
	AIAction      uint32 // 1360, before named actions
	EscortName    string // 1364
	SpellTable    []byte // before version 34
	Spells        map[string]uint32
	Timings       [5][2]uint16 // 1448
	Field329      float32
	Field510      uint32 // 2040
	Field331      byte   // 1324, NPCs
	Field332      float32
	Field330      float32
	DefaultStats  bool // before version 42: no stats stored
	SpellSlots    [3]string
	SpellSlotsOld [3]uint32
	ActionName    string

	ActionVers  uint16
	ActionState bool // the AI state of a saved game was stored (not decoded)

	Byte1445 byte   // 1445: the creature is dead
	Word0    uint32 // 0
	SubClass uint32 // subclass bits stored with the object
	Colors   [6][3]byte
	Voice    string
	Buffs    []Buff
	Poison   byte
}

// Buff is an enchantment on a creature.
type Buff struct {
	Name  string
	Power byte
	Timer uint32
	Extra uint32 // ENCHANT_SHIELD from version 2: the shield's strength
}

// ShopItem is an item a shopkeeper sells.
type ShopItem struct {
	Count     byte
	Type      string
	Reward    string // the spell, ability or creature of a book
	Modifiers [4]string
}

type Monster struct {
	Vers   uint16
	Object Object
	Creature
	Field337, Field335, Field361, Field509 byte // 1348, 1340, 1444, 2036
	ShopField1716                          uint32
	ShopField1720                          uint32
	ShopText                               string
	Shop                                   []ShopItem
	Sub                                    []Xfer
}

func (*Monster) XferType() Type        { return "MonsterXfer" }
func (x *Monster) xferObject() *Object { return &x.Object }

// DecodeXfer decodes monster data. Without the object type (see DecodeByObjectType), the parts that depend on its
// subclass (shopkeepers, colors, voices) cannot be decoded.
func (x *Monster) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	return x.decodeTyped(reg, "", r)
}

const (
	subclassShopkeeper = 0x8
	subclassFemaleNPC  = 0x20 // stores colors and a voice
	subclassWoundedNPC = 0x80000
)

func (x *Monster) decodeTyped(reg ObjectRegistry, typ string, r *binenc.Reader) error {
	*x = Monster{}
	rd, v := decodeStart(reg, r, "monster", 64, false, &x.Object)
	x.Vers = v
	c := &x.Creature
	sub := objectClass(reg, typ).SubClass
	rd.creatureHeader(c, v, false)
	rd.creatureStats(c, v, false)
	if v >= 41 {
		rd.creatureAction(c)
	}
	if v >= 42 {
		c.Byte1445 = rd.u8()
	}
	if v >= 43 && sub&subclassShopkeeper != 0 {
		if v >= 50 {
			x.ShopField1716 = rd.u32()
		}
		if v >= 61 {
			x.ShopField1720 = rd.u32()
		}
		if v >= 48 {
			x.ShopText = rd.str8()
		}
		n := int(rd.u8())
		for i := 0; i < n && rd.err == nil; i++ {
			var it ShopItem
			if v < 50 {
				rd.u32()
			}
			it.Count = rd.u8()
			it.Type = rd.str8()
			if v >= 47 {
				it.Reward = rd.str8()
			}
			for k := range it.Modifiers {
				it.Modifiers[k] = rd.str8()
			}
			x.Shop = append(x.Shop, it)
		}
	}
	if v >= 44 {
		c.Word0 = rd.u32()
	}
	if v >= 45 {
		c.SubClass = rd.u32()
		sub |= c.SubClass
	}
	if v >= 49 {
		c.Health = rd.u16()
	}
	if v >= 51 {
		x.Field337, x.Field335, x.Field361, x.Field509 = rd.u8(), rd.u8(), rd.u8(), rd.u8()
	}
	if sub&subclassFemaleNPC != 0 && v >= 54 {
		rd.colors(c)
	}
	if sub&subclassFemaleNPC != 0 && v >= 55 {
		c.Voice = rd.str8()
	}
	if v >= 62 {
		rd.buffs(c)
	}
	if v >= 63 && sub&subclassWoundedNPC != 0 {
		c.Voice = rd.str8()
	}
	if v >= 64 {
		c.Poison = rd.u8()
	}
	return rd.err
}

type NPC struct {
	Vers   uint16
	Object Object
	Creature
	HealthMax  uint32
	Experience float32
	Sub        []Xfer
}

func (*NPC) XferType() Type        { return "NPCXfer" }
func (x *NPC) xferObject() *Object { return &x.Object }

func (x *NPC) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = NPC{}
	rd, v := decodeStart(reg, r, "NPC", 62, false, &x.Object)
	x.Vers = v
	c := &x.Creature
	rd.creatureHeader(c, v, true)
	rd.colors(c)
	if v == 31 {
		rd.u16()
		n := int(rd.u8())
		for i := 0; i < n && rd.err == nil; i++ {
			rd.str8()
			rd.u8()
			rd.u8()
		}
	}
	rd.creatureStats(c, v, true)
	if v >= 41 {
		rd.creatureAction(c)
	}
	if v >= 42 {
		c.Byte1445 = rd.u8()
	}
	if v >= 44 {
		c.Word0 = rd.u32()
	}
	if v >= 45 {
		x.HealthMax = rd.u32()
	}
	if v >= 46 {
		c.SubClass = rd.u32()
	}
	if v >= 48 {
		c.Health = rd.u16()
	}
	if v >= 51 {
		x.Experience = rd.f32()
	}
	if v >= 52 {
		c.Voice = rd.str8()
	}
	if v >= 61 {
		rd.buffs(c)
	}
	if v >= 62 {
		c.Poison = rd.u8()
	}
	return rd.err
}

func (rd *reader) colors(c *Creature) {
	for i := range c.Colors {
		copy(c.Colors[i][:], rd.bytes(3))
	}
}

func (rd *reader) creatureHeader(c *Creature, v uint16, npc bool) {
	c.Direction = [2]uint32{rd.u32(), rd.u32()}
	script := func(i int) {
		if !npc && v < 3 {
			if rd.err == nil {
				rd.err = fmt.Errorf("unsupported monster xfer vers: %v", v)
			}
			return
		}
		h := rd.script()
		c.Scripts[i] = &h
	}
	script(0) // enemy sighted
	script(1) // death
	c.Field305 = rd.u16()
	script(2) // looking for enemy
	extra, last := uint16(31), uint16(52)
	if npc {
		extra, last = 32, 50
	}
	if v >= extra {
		for i := 3; i <= 8; i++ {
			script(i)
		}
		if v >= last {
			script(9)
		}
	}
	if npc || v >= 11 {
		c.HeaderWord = rd.u32()
	}
}

func (rd *reader) creatureStats(c *Creature, v uint16, npc bool) {
	start, statusWide, timingWide, actionNamed, more := uint16(31), uint16(51), uint16(46), uint16(34), uint16(33)
	if npc {
		start, statusWide, timingWide, actionNamed, more = 32, 49, 47, 35, 34
	}
	if v < start {
		return
	}
	c.Field333 = rd.u8()
	if v < statusWide {
		c.StatusFlags = uint32(rd.u16())
	} else {
		c.StatusFlags = rd.u32()
	}
	c.Field338 = rd.f32()
	c.RetreatLevel = rd.f32()
	c.ResumeLevel = rd.f32()
	c.SightRange = rd.f32()
	if npc {
		c.Health = rd.u16()
	} else if v < 33 {
		rd.u16()
	}
	c.Aggression = rd.f32()
	if v < actionNamed {
		c.AIAction = rd.u32()
	}
	c.EscortName = rd.str8()
	if v < 34 {
		c.SpellTable = rd.bytes(548)
	} else {
		n := int(rd.u32())
		c.Spells = make(map[string]uint32)
		for i := 0; i < n && rd.err == nil; i++ {
			name := rd.str8()
			c.Spells[name] = rd.u32()
		}
	}
	for i := range c.Timings {
		if v < timingWide {
			c.Timings[i] = [2]uint16{uint16(rd.u8()), uint16(rd.u8())}
		} else {
			c.Timings[i] = [2]uint16{rd.u16(), rd.u16()}
		}
		if v < more {
			rd.u32()
		}
	}
	if v >= more-1 {
		c.Field329 = rd.f32()
	}
	if v >= more {
		c.Field510 = rd.u32()
		if npc {
			c.Field331 = rd.u8()
			c.Field332 = rd.f32()
		}
		c.Field330 = rd.f32()
		if v < 42 && rd.u16() == 0 {
			c.DefaultStats = true
		}
		for i := range c.SpellSlots {
			if npc || v >= 53 {
				c.SpellSlots[i] = rd.str8()
			} else {
				c.SpellSlotsOld[i] = rd.u32()
			}
		}
	}
	if v >= actionNamed {
		c.ActionName = rd.str8()
	}
}

// creatureAction reads the AI state, which maps do not store (only saved games do).
func (rd *reader) creatureAction(c *Creature) {
	c.ActionVers = rd.u16()
	if rd.err != nil {
		return
	}
	if int16(c.ActionVers) > 4 {
		rd.err = fmt.Errorf("unsupported creature action xfer vers: %v", c.ActionVers)
		return
	}
	if c.ActionVers >= 2 && rd.u8() == 0 {
		return
	}
	c.ActionState = true
	if rd.err == nil {
		rd.err = fmt.Errorf("creature AI state (of a saved game) is not supported")
	}
}

func (rd *reader) buffs(c *Creature) {
	v := rd.u16()
	if rd.err != nil {
		return
	}
	if int16(v) <= 0 || v > 2 {
		rd.err = fmt.Errorf("unsupported creature buffs xfer vers: %v", v)
		return
	}
	n := int(rd.u8())
	for i := 0; i < n && rd.err == nil; i++ {
		b := Buff{Name: rd.str8(), Power: rd.u8(), Timer: rd.u32()}
		if b.Name == "ENCHANT_SHIELD" && v >= 2 {
			b.Extra = rd.u32()
		}
		c.Buffs = append(c.Buffs, b)
	}
}
