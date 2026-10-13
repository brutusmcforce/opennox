package maprender

import (
	"image"
	"image/color"
	"image/draw"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/modifiers"
	"github.com/brutusmcforce/opennox/libs/noximage/pcx"
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/player"
	"github.com/brutusmcforce/opennox/libs/things"
	"github.com/brutusmcforce/opennox/libs/xfer"
)

// Sprites have material pixels: a gray level, drawn multiplied by the color the game sets for the material. Like
// the game's renderer, the renderer keeps the colors from one object to the next.

var white = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

func (r *Renderer) resetMaterials() {
	for i := range r.mats {
		r.mats[i] = white
	}
}

// shade returns the image with the current material colors.
func (r *Renderer) shade(img *pcx.Image) image.Image {
	mat := img.Material
	if mat == nil {
		return img.Image
	}
	b := img.Bounds()
	out := image.NewNRGBA(b)
	draw.Draw(out, b, img.Image, b.Min, draw.Src)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			m := int(mat.ColorIndexAt(x, y))
			if m == 0 {
				continue
			}
			c := r.mats[(m-1)%len(r.mats)]
			p := out.NRGBAAt(x, y)
			p.R = byte(int(p.R) * int(c.R) / 255)
			p.G = byte(int(p.G) * int(c.G) / 255)
			p.B = byte(int(p.B) * int(c.B) / 255)
			out.SetNRGBA(x, y, p)
		}
	}
	return out
}

func (r *Renderer) setMaterial(slot int32, c modifiers.Color) {
	if slot < 0 || int(slot) >= len(r.mats) {
		return
	}
	r.mats[slot] = color.RGBA{R: byte(c.R), G: byte(c.G), B: byte(c.B), A: 0xff}
}

// setMaterial5551 sets a material color the way the game stores some: in 15 bits.
func (r *Renderer) setMaterial5551(slot int, c [3]byte) {
	n := noxcolor.RGB5551Color(c[0], c[1], c[2]).ColorNRGBA()
	r.mats[slot] = color.RGBA{R: n.R, G: n.G, B: n.B, A: 0xff}
}

// baseColor is the color of a slot of an item definition: white if it has none (slot 0 has none).
func baseColor(def *modifiers.WeaponOrArmor, slot int) modifiers.Color {
	if slot >= 1 && slot <= len(def.Colors) {
		return def.Colors[slot-1].Color
	}
	return modifiers.Color{R: 0xff, G: 0xff, B: 0xff}
}

// slots are the material slots the four kinds of an item's modifiers color (effectiveness, material, primary and
// secondary enchantment), or -1.
func slots(def *modifiers.WeaponOrArmor) [4]int32 {
	var out [4]int32
	for i, v := range []modifiers.ColorIndex{def.Effectiveness, def.Material, def.PriEnchant, def.SecondEnchant} {
		out[i] = -1
		if v >= 0 {
			out[i] = int32(v) + 1
		}
	}
	return out
}

// modifierColor returns the color of an item modifier (an enchantment, a material or an effectiveness).
func (r *Renderer) modifierColor(name string) (modifiers.Color, bool) {
	if r.mods == nil || name == "" {
		return modifiers.Color{}, false
	}
	for _, list := range [][]modifiers.Effect{r.mods.Effectiveness, r.mods.Materials, r.mods.Enchants} {
		for i := range list {
			if list[i].Name == name {
				return list[i].Color, true
			}
		}
	}
	return modifiers.Color{}, false
}

// itemModifiers returns the names of an item's modifiers.
func itemModifiers(x xfer.Xfer) []*xfer.Modifier {
	switch v := x.(type) {
	case *xfer.Weapon:
		return v.Modifiers
	case *xfer.Armor:
		return v.Modifiers
	case *xfer.Ammo:
		return v.Modifiers
	case *xfer.Team:
		return v.Modifiers
	}
	return nil
}

// setItemMaterials sets the colors of an item as the game does for one on the ground: its definition's, then its
// modifiers'.
func (r *Renderer) setItemMaterials(typ string, armor bool, x xfer.Xfer) {
	if r.mods == nil {
		return
	}
	var def *modifiers.WeaponOrArmor
	if armor {
		def = r.mods.ArmorByName(typ)
	} else {
		def = r.mods.WeaponByName(typ)
	}
	if def == nil {
		return
	}
	for i := 1; i <= 6; i++ {
		r.setMaterial(int32(i), baseColor(def, i))
	}
	r.setModifierMaterials(def, itemModifiers(x))
}

func (r *Renderer) setModifierMaterials(def *modifiers.WeaponOrArmor, mods []*xfer.Modifier) {
	sl := slots(def)
	for i, m := range mods {
		if i >= len(sl) || m == nil {
			continue
		}
		if c, ok := r.modifierColor(m.Name); ok {
			r.setMaterial(sl[i], c)
		}
	}
}

// The armor and the weapons a character can wear, with the bit each has in the game, which also selects its
// animation (see server.PlayerArmorNames and server.PlayerWeaponNames). The first item of a bit is the one the
// game takes the colors from.
var armorBits = []struct {
	name string
	bit  uint
}{
	{"LeatherArmor", 15}, {"ChainTunic", 16}, {"Breastplate", 17}, {"LeatherHelm", 21}, {"SteelHelm", 22},
	{"SteelShield", 25}, {"WoodenShield", 24}, {"StreetShirt", 10}, {"StreetPants", 2}, {"StreetSneakers", 0},
	{"LeatherBoots", 6}, {"LeatherArmoredBoots", 7}, {"PlateBoots", 8}, {"MedievalPants", 3},
	{"LeatherLeggings", 4}, {"ChainLeggings", 5}, {"PlateLeggings", 9}, {"MedievalShirt", 11}, {"WizardRobe", 14},
	{"LeatherArmbands", 12}, {"PlateArms", 13}, {"MedievalCloak", 1}, {"ChainCoif", 18}, {"WizardHelm", 19},
	{"ConjurerHelm", 20}, {"OrnateHelm", 23},
}

var armorParts = []player.AnimPart{
	"STREET_SNEAKERS", "MEDIEVAL_CLOAK", "STREET_PANTS", "MEDIEVAL_PANTS", "LEATHER_LEGGINGS", "CHAIN_LEGGINGS",
	"LEATHER_BOOTS", "LEATHER_ARMORED_BOOTS", "PLATE_BOOTS", "PLATE_LEGGINGS", "STREET_SHIRT", "MEDIEVAL_SHIRT",
	"LEATHER_ARMBANDS", "PLATE_ARMS", "WIZARD_ROBE", "LEATHER_TUNIC", "CHAIN_TUNIC", "PLATE_BREAST", "CHAIN_COIF",
	"WIZARD_HELM", "CONJURER_HELM", "LEATHER_HELM", "PLATE_HELM", "ORNATE_HELM", "ROUND_SHIELD", "KITE_SHIELD",
}

var weaponBits = []struct {
	name string
	bit  uint
}{
	{"Bow", 2}, {"CrossBow", 3}, {"Longsword", 9}, {"GreatSword", 10}, {"MorningStar", 11}, {"BattleAxe", 12},
	{"OgreAxe", 13}, {"Sword", 8}, {"WarHammer", 14}, {"RoundChakram", 6}, {"FanChakram", 7}, {"Quiver", 1},
	{"StaffWooden", 15}, {"SulphorousFlareWand", 16}, {"SulphorousShowerWand", 17}, {"LesserFireballWand", 19},
	{"DemonsBreathWand", 19}, {"DeathRayWand", 22}, {"InfinitePainWand", 21}, {"FireStormWand", 20},
	{"ForceWand", 18}, {"OblivionHalberd", 23}, {"OblivionHeart", 24}, {"OblivionWierdling", 25},
	{"OblivionOrb", 26}, {"Flag", 0},
}

var weaponParts = []player.AnimPart{
	"FLAG", "QUIVER", "BOW", "CROSSBOW", "ARROW", "BOLT", "CHAKRAM", "SHURIKEN", "SWORD", "LONG_SWORD",
	"GREAT_SWORD", "MACE", "AXE", "OGRE_AXE", "HAMMER", "STAFF", "STAFF_SULPHOROUS_FLARE", "STAFF_SULPHOROUS_SHOWER",
	"STAFF_LIGHTNING", "STAFF_FIREBALL", "STAFF_TRIPLE_FIREBALL", "STAFF_FORCE_OF_NATURE", "STAFF_DEATH_RAY",
	"STAFF_OBLIVION_HALBERD", "STAFF_OBLIVION_HEART", "STAFF_OBLIVION_WIERDLING", "STAFF_OBLIVION_ORB",
}

func bitOf(list []struct {
	name string
	bit  uint
}, name string) (uint, bool) {
	for _, it := range list {
		if it.name == name {
			return it.bit, true
		}
	}
	return 0, false
}

func firstOfBit(list []struct {
	name string
	bit  uint
}, bit uint) string {
	for _, it := range list {
		if it.bit == bit {
			return it.name
		}
	}
	return ""
}

// equipment is what a character wears: the modifiers of the item of each bit.
type equipment struct {
	armor, weapons uint32
	armorMods      map[uint][]*xfer.Modifier
	weaponMods     map[uint][]*xfer.Modifier
}

// equipmentOf returns what the items of an inventory the game equips (those flagged as equipped).
func equipmentOf(inv []xfer.Xfer) *equipment {
	e := &equipment{armorMods: make(map[uint][]*xfer.Modifier), weaponMods: make(map[uint][]*xfer.Modifier)}
	for _, it := range inv {
		ob := xfer.ObjectOf(it)
		if ob == nil || ob.Val5 == 0 || object.Flags(ob.Flags)&object.FlagEquipped == 0 {
			continue
		}
		if bit, ok := bitOf(armorBits, ob.Type); ok {
			e.armor |= 1 << bit
			e.armorMods[bit] = itemModifiers(it)
		} else if bit, ok := bitOf(weaponBits, ob.Type); ok {
			e.weapons |= 1 << bit
			e.weaponMods[bit] = itemModifiers(it)
		}
	}
	return e
}

// setEquipMaterials sets the colors of a worn item as the game does: slots 0 to n-1 from its definition (by the
// item that comes first for its bit), then its modifiers'.
func (r *Renderer) setEquipMaterials(armor bool, bit uint, mods []*xfer.Modifier) {
	if r.mods == nil {
		return
	}
	var (
		def *modifiers.WeaponOrArmor
		n   int
	)
	if armor {
		def, n = r.mods.ArmorByName(firstOfBit(armorBits, bit)), 6
	} else {
		def, n = r.mods.WeaponByName(firstOfBit(weaponBits, bit)), 5
	}
	if def == nil {
		return
	}
	for i := 0; i < n; i++ {
		r.setMaterial(int32(i), baseColor(def, i))
	}
	r.setModifierMaterials(def, mods)
}

// npcLayer is an image of an NPC (its body or what it wears) with the material colors to draw it with.
type npcLayer struct {
	img  *pcx.Image
	mats [16]color.RGBA
}

// npcLayers returns the images the game draws for an NPC standing idle: its body and what it wears (see the
// client's DrawNPC and DrawArmorAndWeapon).
func (r *Renderer) npcLayers(npc *xfer.NPC) ([]npcLayer, error) {
	pt := r.things["NewPlayer"]
	if pt == nil {
		return nil, nil
	}
	pd, ok := pt.Draw.(things.PlayerDraw)
	if !ok {
		return nil, nil
	}
	anim := pd.Anims[player.AnimIdle]
	if anim == nil {
		return nil, nil
	}
	gridDir, dir := creatureDir(&npc.Creature)
	var (
		out  []npcLayer
		last error
	)
	add := func(part player.AnimPart) {
		frames := anim.Parts[part][dir]
		if len(frames) == 0 {
			return
		}
		img, err := r.getPCX(frames[0].Ind)
		if err != nil {
			last = err
			return
		}
		out = append(out, npcLayer{img: img, mats: r.mats})
	}
	for i, c := range npc.Colors {
		r.setMaterial5551(i+1, c)
	}
	add("NAKED")
	eq := equipmentOf(npc.Sub)
	weapons := func(mask uint32) {
		for bit := uint(1); bit < uint(len(weaponParts)); bit++ {
			if eq.weapons&mask&(1<<bit) == 0 {
				continue
			}
			if len(anim.Parts[weaponParts[bit]][dir]) == 0 {
				continue
			}
			r.setEquipMaterials(false, bit, eq.weaponMods[bit])
			add(weaponParts[bit])
		}
	}
	armor := func() {
		wear := func(bit uint) {
			if len(anim.Parts[armorParts[bit]][dir]) == 0 {
				return
			}
			r.setEquipMaterials(true, bit, eq.armorMods[bit])
			add(armorParts[bit])
		}
		const cloak, shields = 1, 1<<24 | 1<<25
		drawCloak := false
		if eq.armor&(1<<cloak) != 0 && (gridDir == 6 || gridDir == 7 || gridDir == 8) {
			if len(anim.Parts[armorParts[cloak]][dir]) != 0 {
				wear(cloak)
				drawCloak = true
			}
		}
		drawShield := false
		for bit := uint(0); bit < uint(len(armorParts)); bit++ {
			b := uint32(1) << bit
			if bit == cloak || len(anim.Parts[armorParts[bit]][dir]) == 0 {
				continue
			}
			if b&shields != 0 && (gridDir == 3 || gridDir == 0 || gridDir == 6 || gridDir == 7) {
				continue
			}
			// the shirts are under a robe
			if eq.armor&b != 0 && (b&0xC00 == 0 || eq.armor&0x4000 == 0) {
				wear(bit)
				if b&shields != 0 {
					drawShield = true
				}
			}
		}
		if eq.armor&(1<<cloak) != 0 && !drawCloak {
			wear(cloak)
		}
		if !drawShield && eq.armor&shields != 0 {
			bit := uint(25)
			if eq.armor&(1<<24) != 0 {
				bit = 24
			}
			wear(bit)
		}
	}
	// the quiver is on the back
	if gridDir != 0 && gridDir != 1 && gridDir != 2 && gridDir != 3 && gridDir != 6 {
		weapons(2)
		armor()
		weapons(^uint32(2))
	} else {
		weapons(^uint32(2))
		armor()
		weapons(2)
	}
	for i := 1; i <= 6; i++ {
		r.mats[i] = white
	}
	return out, last
}

// creatureDir returns the direction a creature faces on the game's 3x3 grid (0 to 8, without the center 4), and the
// index of its animations for it (which have the 8 directions around the center).
func creatureDir(c *xfer.Creature) (grid, anim int) {
	dx, dy := int(int32(c.Direction[0])), int(int32(c.Direction[1]))
	grid = 3*dy + dx + 4
	if grid < 0 || grid > 8 || grid == 4 {
		grid = 7 // south
	}
	anim = grid
	if anim > 4 {
		anim--
	}
	return grid, anim
}
