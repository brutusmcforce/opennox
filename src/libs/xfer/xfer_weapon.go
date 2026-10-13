package xfer

import (
	"github.com/brutusmcforce/opennox/libs/binenc"
)

func init() {
	Register(&Weapon{})
}

type Weapon struct {
	Vers      uint16
	Object    Object
	Modifiers []*Modifier
	Val4      byte   // wand charges
	Val5      byte   // wand max charges
	Percent   uint32 // wand recharge progress
	Health    uint16
	Val13     byte
	Val14     uint32 // flags
	Sub       []Xfer
}

func (*Weapon) XferType() Type {
	return "WeaponXfer"
}

func (x *Weapon) xferObject() *Object { return &x.Object }

// DecodeXfer decodes weapon data. Without the object type, wands are decoded as other weapons, which fails for
// those that store charges (see DecodeByObjectType).
func (x *Weapon) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	return x.decodeTyped(reg, "", r)
}

// weaponHasCharges checks if a weapon type stores wand charges: wands of the staff subclasses that cast spells.
func weaponHasCharges(c ObjectClass, vers uint16) bool {
	const (
		classWand  = 0x1000
		chargedSub = 0x47f0000
		orb        = 0x4000000 // stored from version 62
	)
	return c.Class&classWand != 0 && c.SubClass&chargedSub != 0 && (vers >= 62 || c.SubClass&orb == 0)
}

func (x *Weapon) decodeTyped(reg ObjectRegistry, typ string, r *binenc.Reader) error {
	*x = Weapon{}
	rd := &reader{r: r}
	x.Vers = rd.version("weapon", 64, false)
	if rd.err != nil {
		return rd.err
	}
	if err := x.Object.DecodeXfer(reg, x.Vers, r); err != nil {
		return err
	}
	if x.Vers < 11 {
		return nil
	}
	for i := 0; i < 4; i++ {
		x.Modifiers = append(x.Modifiers, &Modifier{Name: rd.str8()})
	}
	if x.Vers >= 41 && weaponHasCharges(objectClass(reg, typ), x.Vers) {
		x.Val4 = rd.u8()
		x.Val5 = rd.u8()
		if x.Vers >= 61 {
			x.Percent = rd.u32()
		}
	}
	if x.Vers >= 42 {
		x.Health = rd.u16()
	}
	if x.Vers == 63 {
		x.Val13 = rd.u8()
	}
	if x.Vers >= 64 {
		x.Val14 = rd.u32()
	}
	return rd.err
}
