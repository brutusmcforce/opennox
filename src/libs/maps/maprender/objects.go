package maprender

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"

	"github.com/brutusmcforce/opennox/libs/maps"
	"github.com/brutusmcforce/opennox/libs/noximage/pcx"
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/things"
	"github.com/brutusmcforce/opennox/libs/xfer"
)

// thingInfo is what the renderer needs of an object type.
type thingInfo struct {
	*things.Thing
	flags   object.Flags // the default flags
	class   object.Class
	monster object.MonsterClass
}

func (r *Renderer) indexThings() error {
	list, err := r.tng.ReadThings()
	if err != nil {
		return err
	}
	r.things = make(map[string]*thingInfo, len(list))
	for i := range list {
		th := &list[i]
		info := &thingInfo{Thing: th}
		for _, f := range th.Flags {
			if v, err := object.ParseFlag(string(f)); err == nil {
				info.flags |= v
			}
		}
		for _, c := range th.Class {
			if v, err := object.ParseClass(string(c)); err == nil {
				info.class |= v
			}
		}
		if info.class&object.ClassMonster != 0 {
			for _, c := range th.SubClass {
				if v, err := object.ParseMonsterClass(string(c)); err == nil {
					info.monster |= v
				}
			}
		}
		// as the server sets them (see server.ObjectType)
		if th.OnCollide == nil {
			info.flags |= object.FlagNoCollide
		}
		info.flags |= object.FlagEnabled
		r.things[th.Name] = info
	}
	return nil
}

// mapObject is an object as the game's client sees it.
type mapObject struct {
	typ    *thingInfo
	x      xfer.Xfer
	obj    *xfer.Object
	pos    image.Point
	flags  object.Flags
	flags2 uint32 // the client's state flags of the object (animations)
}

// The passes the game draws objects in (see the client's nox_xxx_spriteAddQueue_475560_draw).
const (
	passNone   = iota
	passFloor  // flat objects on the floor, drawn with the floor tiles
	passBelow  // other objects below the walls and the objects
	passAbove  // airborne objects (drawn before the others in the game too)
	passSorted // sorted by Y with the walls
)

func (o *mapObject) pass() int {
	f, c := o.flags, o.typ.class
	if o.typ.Draw == nil || f&object.FlagEditVisible != 0 {
		return passNone
	}
	if f&object.FlagBelow != 0 {
		switch o.typ.Draw.(type) {
		case things.StaticDraw, things.StaticRandomDraw:
			if c&(object.ClassPickup|object.ClassVisibleEnable) == 0 &&
				(f&(object.FlagAllowOverlap|object.FlagNoCollide) != 0 || c&object.ClassImmobile != 0) &&
				f&object.FlagOwnerVisible == 0 {
				return passFloor
			}
		}
		return passBelow
	}
	if f&object.FlagAirborne != 0 && f&object.FlagNoCollide != 0 {
		return passAbove
	}
	if c&object.ClassFire == 0 || f&object.FlagEnabled != 0 {
		return passSorted
	}
	return passNone
}

// sortY is the Y the game sorts the object by, among the other objects and the walls.
func (o *mapObject) sortY() int {
	if d, ok := o.x.(*xfer.Door); ok && o.typ.class&object.ClassDoor != 0 {
		// with no player to see it from
		return o.pos.Y + doorOffset(d.Direction).Y/2
	}
	return o.pos.Y
}

// doorOffset is the offset of a door's far end for its direction (the game's table, see server.DoorSize).
func doorOffset(dir uint32) image.Point {
	if int(dir) >= len(doorTable) {
		return image.Point{}
	}
	return doorTable[dir]
}

var doorTable = []image.Point{
	{X: -23, Y: -23}, {X: -18, Y: -27}, {X: -12, Y: -30}, {X: -6, Y: -31},
	{X: 0, Y: -32}, {X: 6, Y: -31}, {X: 12, Y: -30}, {X: 18, Y: -27},
	{X: 23, Y: -23}, {X: 27, Y: -18}, {X: 30, Y: -12}, {X: 31, Y: -6},
	{X: 32, Y: 0}, {X: 31, Y: 6}, {X: 30, Y: 12}, {X: 27, Y: 18},
	{X: 23, Y: 23}, {X: 18, Y: 27}, {X: 12, Y: 30}, {X: 6, Y: 31},
	{X: 0, Y: 32}, {X: -6, Y: 31}, {X: -12, Y: 30}, {X: -18, Y: 27},
	{X: -23, Y: 23}, {X: -27, Y: 18}, {X: -30, Y: 12}, {X: -31, Y: 6},
	{X: -32, Y: 0}, {X: -31, Y: -6}, {X: -30, Y: -12}, {X: -27, Y: -18},
}

// readObjects reads the map's objects as the game's client would show them when the map starts.
func (r *Renderer) readObjects(m *maps.Map) ([]*mapObject, error) {
	list, err := m.ReadObjects(nil)
	if err != nil {
		return nil, err
	}
	out := make([]*mapObject, 0, len(list))
	for _, it := range list {
		typ := r.things[it.Type]
		ob := xfer.ObjectOf(it.Xfer)
		// objects without a drawing can still give light
		if typ == nil || ob == nil {
			continue
		}
		o := &mapObject{
			typ: typ, x: it.Xfer, obj: ob,
			pos:   image.Pt(int(math.Round(float64(ob.Pos.X))), int(math.Round(float64(ob.Pos.Y)))),
			flags: typ.flags,
		}
		if ob.Val5 != 0 {
			// the flags the map stores replace these
			o.flags = o.flags&0xEEBF7E9D | object.Flags(ob.Flags)
			o.flags2 = ob.Anim
		}
		out = append(out, o)
	}
	return out, nil
}

// objectImage returns the image the game shows for the object when the map starts, or nil for none.
func (r *Renderer) objectImage(o *mapObject) (*pcx.Image, error) {
	enabled := o.flags&object.FlagEnabled != 0
	var ref *things.ImageRef
	frame0 := func(a *things.Animation) {
		if len(a.Frames) > 0 {
			ref = &a.Frames[0]
		}
	}
	switch d := o.typ.Draw.(type) {
	case things.StaticDraw:
		if o.typ.class&object.ClassReadable == 0 || enabled {
			ref = &d.Img
		}
	case things.BaseDraw:
		ref = &d.Img
	case things.WeaponDraw:
		ref = &d.Img
	case things.ArmorDraw:
		ref = &d.Img
	case things.StaticRandomDraw:
		// the game picks one at random; this picks the same for the same object
		if len(d.Imgs) > 0 {
			ref = &d.Imgs[int(o.obj.Extent*2654435761>>16)%len(d.Imgs)]
		}
	case things.SlaveDraw:
		if len(d.Imgs) > 0 {
			ref = &d.Imgs[0]
		}
	case things.BoulderDraw:
		if len(d.Imgs) > 0 {
			ref = &d.Imgs[0]
		}
	case things.DoorDraw:
		if dr, ok := o.x.(*xfer.Door); ok && int(dr.Direction) < len(d.Imgs) {
			ref = &d.Imgs[dr.Direction]
		}
	case things.AnimateDraw:
		if d.Anim.Kind == things.AnimationLoop && o.typ.class&object.ClassFlag != 0 && !enabled {
			break
		}
		frame0(&d.Anim)
	case things.GlyphDraw:
		frame0(&d.Anim)
	case things.FlagDraw:
		frame0(&d.Anim)
	case things.WeaponAnimateDraw:
		frame0(&d.Anim)
	case things.ArmorAnimateDraw:
		frame0(&d.Anim)
	case things.SphericalShieldDraw:
		frame0(&d.Anim)
	case things.SummonEffectDraw:
		frame0(&d.Anim)
	case things.ConditionalAnimateDraw:
		state := 1
		if enabled {
			state = 0
		}
		if state < len(d.Anims) {
			frame0(&d.Anims[state])
		}
	case things.AnimateStateDraw:
		// the state the object's state flags select, and the state each animation is for
		stateOf := func(f uint32) int {
			switch {
			case f&2 != 0:
				return 0
			case f&4 != 0:
				return 1
			}
			return 2
		}
		want := 0
		if o.flags2&6 != 0 || o.flags2&8 != 0 {
			want = stateOf(o.flags2)
		}
		for i := range d.States {
			if st := &d.States[i]; stateOf(st.Flags) == want {
				frame0(&st.Anim)
				break
			}
		}
	case things.MonsterDraw:
		ref = monsterFrame(d.Anims, o.x)
	case things.MaidenDraw:
		ref = monsterFrame(d.Anims, o.x)
	}
	if ref == nil {
		return nil, nil
	}
	return r.getPCX(ref.Ind)
}

// monsterFrame returns the first frame of the monster's idle animation, in the direction it faces.
func monsterFrame(anims []things.MonsterAnimation, x xfer.Xfer) *things.ImageRef {
	var c *xfer.Creature
	switch m := x.(type) {
	case *xfer.Monster:
		c = &m.Creature
	case *xfer.NPC:
		c = &m.Creature
	default:
		return nil
	}
	var anim *things.MonsterAnimation
	for i := range anims {
		if anims[i].Type == things.MonsterAnimIdle {
			anim = &anims[i]
			break
		}
	}
	if anim == nil {
		return nil
	}
	_, dir := creatureDir(c)
	if len(anim.Frames[dir]) == 0 {
		return nil
	}
	return &anim.Frames[dir][0]
}

// drawObject draws the object as the game does: its image is offset by half its size and raised by its height.
func (r *Renderer) drawObject(out *image.RGBA, o *mapObject) error {
	if d, ok := o.typ.Draw.(things.UnknownDraw); ok {
		switch d.Type {
		case "PressurePlateDraw":
			// the game draws the outline of pressure plates; triggers are drawn translucent with what the engine
			// set before, not drawn here
			if tr, ok := o.x.(*xfer.Trigger); ok {
				drawPlate(out, o.pos, tr)
			}
		case "NPCDraw":
			if npc, ok := o.x.(*xfer.NPC); ok {
				layers, err := r.npcLayers(npc)
				saved := r.mats
				for _, l := range layers {
					r.mats = l.mats
					r.drawLit(out, o, l.img)
				}
				r.mats = saved
				if err != nil {
					return fmt.Errorf("%s: %w", o.typ.Name, err)
				}
			}
		}
		return nil
	}
	switch o.typ.Name {
	case "Ghost":
		// only visible close to the player
		return nil
	case "RewardMarker", "RewardMarkerPlus":
		// a quest replaces some with rewards and removes the others when it starts
		return nil
	}
	img, err := r.objectImage(o)
	if err != nil {
		return fmt.Errorf("%s: %w", o.typ.Name, err)
	} else if img == nil {
		return nil
	}
	// the colors the game sets for the object
	switch o.typ.Draw.(type) {
	case things.WeaponDraw, things.WeaponAnimateDraw, things.FlagDraw:
		r.setItemMaterials(o.typ.Name, false, o.x)
	case things.ArmorDraw, things.ArmorAnimateDraw:
		r.setItemMaterials(o.typ.Name, true, o.x)
	case things.MaidenDraw:
		if m, ok := o.x.(*xfer.Monster); ok {
			for i, c := range m.Colors {
				r.setMaterial5551(i+1, c)
			}
		}
	}
	r.drawLit(out, o, img)
	return nil
}

// drawLit draws an image of the object with the light at its position (see the client's objectRenderDraw).
func (r *Renderer) drawLit(out *image.RGBA, o *mapObject, img *pcx.Image) {
	l := [3]int{255, 255, 255}
	// lights themselves, and objects flagged to be, are drawn fully lit
	lit := o.typ.class&object.ClassLight != 0 && o.flags&object.FlagEnabled != 0 || o.flags&0x40000000 != 0
	if r.light != nil && !lit {
		l = r.light.sample(o.pos)
	}
	alpha := byte(0xff)
	if o.flags&object.FlagTranslucent != 0 {
		alpha = 128
	}
	drawLit(out, o.rect(img), r.shade(img), l, l, 255, alpha)
}

// rect is where the object's image is drawn.
func (o *mapObject) rect(img *pcx.Image) image.Rectangle {
	var size image.Point
	if o.typ.Size != nil {
		size = *o.typ.Size
	}
	z := 0
	if o.typ.Z != nil {
		z = *o.typ.Z
	}
	p := o.pos.Sub(image.Pt(size.X/2, size.Y/2+z)).Add(img.Point)
	b := img.Bounds()
	return b.Sub(b.Min).Add(p)
}

// objectsBounds returns the bounds of the objects' images (an empty rectangle if none).
func (r *Renderer) objectsBounds(objs []*mapObject) image.Rectangle {
	var bb image.Rectangle
	for _, o := range objs {
		if o.pass() == passNone {
			continue
		}
		if d, ok := o.typ.Draw.(things.UnknownDraw); ok {
			// pressure plates stay within the floor
			if npc, ok := o.x.(*xfer.NPC); ok && d.Type == "NPCDraw" {
				saved := r.mats
				layers, _ := r.npcLayers(npc)
				r.mats = saved
				for _, l := range layers {
					bb = bb.Union(o.rect(l.img))
				}
			}
			continue
		}
		img, err := r.objectImage(o)
		if err != nil || img == nil {
			continue
		}
		bb = bb.Union(o.rect(img))
	}
	return bb
}

// drawPlate draws the outline of a pressure plate in its two colors, as the game does.
func drawPlate(out draw.Image, pos image.Point, tr *xfer.Trigger) {
	c := tr.Colors
	if c == [6]byte{} {
		return
	}
	w, h := float64(min(tr.SizeX, 60)), float64(min(tr.SizeY, 60))
	px, py := w*0.35355, h*0.35355
	pt := func(x, y float64) image.Point { return pos.Add(image.Pt(int(math.Round(x)), int(math.Round(y)))) }
	p := [4]image.Point{pt(-px+py, -px-py), pt(px-py, px+py), pt(-px-py, -px+py), pt(px+py, px-py)}
	for i := 0; i < 2; i++ {
		cl := color.NRGBA{R: c[3*i], G: c[3*i+1], B: c[3*i+2], A: 0xff}
		s := image.Pt(i, 0)
		for _, e := range [][2]int{{0, 3}, {1, 3}, {0, 2}, {1, 2}} {
			drawLine(out, p[e[0]].Add(s), p[e[1]].Add(s), cl)
		}
	}
}

func drawLine(out draw.Image, a, b image.Point, cl color.Color) {
	dx, dy := b.X-a.X, b.Y-a.Y
	n := max(abs(dx), abs(dy))
	for i := 0; i <= n; i++ {
		x, y := a.X, a.Y
		if n > 0 {
			x += (dx*i + n/2*sign(dx)) / n
			y += (dy*i + n/2*sign(dy)) / n
		}
		out.Set(x, y, cl)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func sign(v int) int {
	switch {
	case v < 0:
		return -1
	case v > 0:
		return 1
	}
	return 0
}

// sortObjects sorts the objects by the Y the game sorts them by, keeping their order otherwise.
func sortObjects(list []*mapObject) {
	sort.SliceStable(list, func(i, j int) bool { return list[i].sortY() < list[j].sortY() })
}
