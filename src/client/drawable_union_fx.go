package client

import "unsafe"

// Pointer-free DrawableUnion variants of client-side effects and objects.
// Offsets are relative to the union (Drawable offset 432 in the 386 layout).

// DrawableUnionSpark is a moving spark particle (legacy effectInitSpark).
type DrawableUnionSpark struct {
	X       uint32 // 0: world X << 12
	Y       uint32 // 4: world Y << 12
	Speed   uint32 // 8
	Born    uint32 // 12: creation frame
	Expires uint32 // 16: frame the spark dies
}

// DrawableUnionTail is a tail link (magic missile, arrow), or the anchor of the next tail link.
type DrawableUnionTail struct {
	X uint32 // 0: world position of the other end
	Y uint32 // 4
}

// DrawableUnionRainOrb is a falling orb (legacy effectCreateRainOrb).
type DrawableUnionRainOrb struct {
	ToX   uint32 // 0: world position the orb flies to after landing
	ToY   uint32 // 4
	PrevZ uint16 // 8: height on the previous tick
	Size  byte   // 10: glow radius
}

// DrawableUnionOrb is an orb that flies to a target or orbits it (legacy effectCreateOrb/effectCreateOrbit).
type DrawableUnionOrb struct {
	Target [2]uint16 // 0: X, Y to fly to (orbit: center)
	Start  [2]uint16 // 4: orbit start X, Y
	Radius uint16    // 8: orbit radius
	Angle  byte      // 10: orbit start angle
	Speed  byte      // 11: move speed (orbit: direction)
	Size   byte      // 12: glow radius
	Period byte      // 13: frames per Size step
	Timer  byte      // 14
}

// DrawableUnionBubble is a bubble particle (legacy presentationBubble).
type DrawableUnionBubble struct {
	Glow       uint32 // 0: glow color
	Core       uint32 // 4: core color
	Size       byte   // 8
	Phase      byte   // 9: 1 growing, 2 shrinking, 3 dying
	Timer      byte   // 10
	Period     byte   // 11
	FlipPeriod byte   // 12
	Flip       byte   // 13: frames until Rise flips
	Rise       byte   // 14: signed height step
}

// DrawableUnionVortex is a vortex orb (legacy updateVortex).
type DrawableUnionVortex struct {
	Glow    uint32 // 0: glow color
	Core    uint32 // 4: core color
	CenterX uint32 // 8
	CenterY uint32 // 12
	Angle   byte   // 16
	Speed   byte   // 17: signed angle step
	Radius  byte   // 18
	Rise    byte   // 19: height step
}

// DrawableUnionRay is a ray between two points or two objects (lightning, plasma, harpoon, presentation rays).
// The packed words are unaligned; use the accessor methods.
type DrawableUnionRay struct {
	ByCode byte    // 0: endpoints are net codes, not positions
	Arg    [4]byte // 1: uint32 (lightning duration; low byte: plasma angle, energy)
	From   [4]byte // 5: uint32 net code, or uint16 X, Y
	To     [4]byte // 9: like From
}

func (u *DrawableUnionRay) ArgWord() *uint32    { return (*uint32)(unsafe.Pointer(&u.Arg)) }
func (u *DrawableUnionRay) FromCode() *uint32   { return (*uint32)(unsafe.Pointer(&u.From)) }
func (u *DrawableUnionRay) ToCode() *uint32     { return (*uint32)(unsafe.Pointer(&u.To)) }
func (u *DrawableUnionRay) FromPos() *[2]uint16 { return (*[2]uint16)(unsafe.Pointer(&u.From)) }
func (u *DrawableUnionRay) ToPos() *[2]uint16   { return (*[2]uint16)(unsafe.Pointer(&u.To)) }

// DrawableUnionHeight is the simulated height of a falling or bouncing object (legacy updateHeight).
type DrawableUnionHeight struct {
	Frame       uint32  // 0: last simulated frame
	Height      float32 // 4
	Velocity    float32 // 8
	Restitution float32 // 12: scaled by 0.1 on a bounce
}

// DrawableUnionBoulder is the rolling animation of a boulder (legacy spriteBoulderDraw).
type DrawableUnionBoulder struct {
	OldX  uint32 // 0: position of the last frame step
	OldY  uint32 // 4
	Frame uint32 // 8
	Bank  uint32 // 12: 0 or 16
}

// DrawableUnionGenerator is a monster generator (legacy objectGeneratorDraw).
type DrawableUnionGenerator struct {
	Remaining uint32 // 0: frames left in the spawn animation
}

// DrawableUnionShield is a spherical shield effect (legacy drawableShield).
type DrawableUnionShield struct {
	Owner uint32 // 0: net code of the shielded object
}

// DrawableUnionCloud is a cloud puff (legacy updateCloud).
type DrawableUnionCloud struct {
	Rise byte // 0: height step
}

// DrawableUnionNPC is an NPC (legacy nox_xxx_spriteNPCInfo_49A4B0).
type DrawableUnionNPC struct {
	Anim uint32 // 0: last chosen animation
}

// DrawableUnionBook is a readable item (spell book, ability book, field guide): the server sends its contents
// with MSG_INFO_BOOK_DATA.
type DrawableUnionBook struct {
	Index uint32 // 0: spell, ability or creature type index
}

// DrawableUnionColorLight is a map color light (legacy colorLightUpdate).
type DrawableUnionColorLight struct {
	Counts [3]byte // 0: signed key counts of the color, intensity and penumbra animations
}

// DrawableUnionPressurePlate is a pressure plate (legacy objectPressureDraw).
type DrawableUnionPressurePlate struct {
	Colors [6]byte // 0: two RGB edge colors
}

// All variants must fit in the union.
var (
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionSpark{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionTail{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionRainOrb{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionOrb{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionBubble{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionVortex{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionRay{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionHeight{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionBoulder{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionGenerator{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionShield{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionCloud{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionNPC{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionBook{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionColorLight{})]struct{}
	_ [unsafe.Sizeof(DrawableUnion{}) - unsafe.Sizeof(DrawableUnionPressurePlate{})]struct{}
)

// Field offsets that must match the union-relative bytes used by the legacy code.
var (
	_ [unsafe.Offsetof(DrawableUnionSpark{}.Expires) - 16]struct{}
	_ [16 - unsafe.Offsetof(DrawableUnionSpark{}.Expires)]struct{}
	_ [unsafe.Offsetof(DrawableUnionOrb{}.Timer) - 14]struct{}
	_ [14 - unsafe.Offsetof(DrawableUnionOrb{}.Timer)]struct{}
	_ [unsafe.Offsetof(DrawableUnionBubble{}.Rise) - 14]struct{}
	_ [14 - unsafe.Offsetof(DrawableUnionBubble{}.Rise)]struct{}
	_ [unsafe.Offsetof(DrawableUnionVortex{}.Rise) - 19]struct{}
	_ [19 - unsafe.Offsetof(DrawableUnionVortex{}.Rise)]struct{}
	_ [unsafe.Offsetof(DrawableUnionRay{}.To) - 9]struct{}
	_ [9 - unsafe.Offsetof(DrawableUnionRay{}.To)]struct{}
	_ [unsafe.Offsetof(DrawableUnionRainOrb{}.Size) - 10]struct{}
	_ [10 - unsafe.Offsetof(DrawableUnionRainOrb{}.Size)]struct{}
	_ [unsafe.Offsetof(DrawableUnionHeight{}.Restitution) - 12]struct{}
	_ [12 - unsafe.Offsetof(DrawableUnionHeight{}.Restitution)]struct{}
)
