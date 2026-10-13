package maprender

import (
	"encoding/binary"
	"image"
	"image/draw"
	"math"

	"github.com/brutusmcforce/opennox/libs/common"
	"github.com/brutusmcforce/opennox/libs/maps"
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/libs/xfer"
)

// Lighting as the game's client computes it (see the client's sub_468F80 and nox_xxx_cliLight16_469140): a grid of
// 8.16 fixed point light values, 23 pixels apart, starts with the ambient light, and each light adds to it. The floor
// is lit by the grid interpolated at each pixel, an object by its value at the object's position, and a wall by its
// values at the wall's ends.
//
// The game uses the ambient light of the polygon (or the map) the player is in for all of the screen; here each
// point gets that of the polygon it is in, as if the player stood there.

const (
	lightMinIntensity = 0x10000 // 1.0
	lightMaxIntensity = 0x1f0000
	lightBase         = 0x40000
	lightFalloff      = 961
	lightFull         = 0xff0000
)

type lightGrid struct {
	org  image.Point // the world position of the first point
	w, h int
	v    [][3]int // 8.16 fixed point, [x*h+y]
}

// newLightGrid returns a light grid covering the rectangle, with its points where the game has them: at the corners
// of the floor tiles.
func newLightGrid(r image.Rectangle) *lightGrid {
	const s = common.GridStep
	x0 := floorDiv(r.Min.X+11, s) - 1
	y0 := floorDiv(r.Min.Y-12, s) - 1
	x1 := floorDiv(r.Max.X+11, s) + 2
	y1 := floorDiv(r.Max.Y-12, s) + 2
	g := &lightGrid{org: image.Pt(s*x0-11, s*y0+12), w: x1 - x0, h: y1 - y0}
	g.v = make([][3]int, g.w*g.h)
	return g
}

func floorDiv(a, b int) int {
	q := a / b
	if a%b != 0 && a < 0 {
		q--
	}
	return q
}

func (g *lightGrid) at(x, y int) *[3]int {
	return &g.v[x*g.h+y]
}

// sample returns the light at a world position, 0-255 for each component (the game's Sub469920).
func (g *lightGrid) sample(p image.Point) [3]int {
	x, y := p.X-g.org.X, p.Y-g.org.Y
	const s = common.GridStep
	xd, yd := x/s, y/s
	if x < 0 || y < 0 || xd+1 >= g.w || yd+1 >= g.h {
		return [3]int{}
	}
	xr, yr := x%s, y%s
	c00, c10, c01, c11 := g.at(xd, yd), g.at(xd+1, yd), g.at(xd, yd+1), g.at(xd+1, yd+1)
	var out [3]int
	for i := range out {
		a := c00[i] + xr*(c10[i]-c00[i])/s
		b := c01[i] + xr*(c11[i]-c01[i])/s
		out[i] = (a + yr*(b-a)/s) >> 16
	}
	return out
}

// sampleF returns the light at a world position as multipliers (the floor shader's bilinear interpolation).
func (g *lightGrid) sampleF(x, y float64) [3]float64 {
	const s = common.GridStep
	gx, gy := (x-float64(g.org.X))/s, (y-float64(g.org.Y))/s
	xd, yd := math.Floor(gx), math.Floor(gy)
	ix, iy := int(xd), int(yd)
	if ix < 0 || iy < 0 || ix+1 >= g.w || iy+1 >= g.h {
		return [3]float64{}
	}
	fx, fy := gx-xd, gy-yd
	c00, c10, c01, c11 := g.at(ix, iy), g.at(ix+1, iy), g.at(ix, iy+1), g.at(ix+1, iy+1)
	var out [3]float64
	for i := range out {
		a := float64(c00[i])*(1-fx) + float64(c10[i])*fx
		b := float64(c01[i])*(1-fx) + float64(c11[i])*fx
		out[i] = (a*(1-fy) + b*fy) / (1 << 24)
	}
	return out
}

// ambient sets each point to the ambient light of the polygon it is in, or of the map.
func (g *lightGrid) ambient(m *maps.Map) {
	amb := types.RGB{R: 25, G: 25, B: 25}
	if m.Ambient != nil {
		amb = m.Ambient.AmbientColor
	}
	polys := mapPolygons(m)
	for x := 0; x < g.w; x++ {
		for y := 0; y < g.h; y++ {
			c := amb
			p := image.Pt(g.org.X+common.GridStep*x, g.org.Y+common.GridStep*y)
			for i := range polys {
				if polys[i].contains(p) {
					c = polys[i].color
					break
				}
			}
			*g.at(x, y) = [3]int{int(c.R) << 16, int(c.G) << 16, int(c.B) << 16}
		}
	}
}

type polygon struct {
	pts    []image.Point
	bounds image.Rectangle
	color  types.RGB
}

func (p *polygon) contains(pt image.Point) bool {
	if !pt.In(p.bounds.Inset(-1)) {
		return false
	}
	in := false
	for i, j := 0, len(p.pts)-1; i < len(p.pts); j, i = i, i+1 {
		a, b := p.pts[i], p.pts[j]
		if (a.Y > pt.Y) != (b.Y > pt.Y) && pt.X < (b.X-a.X)*(pt.Y-a.Y)/(b.Y-a.Y)+a.X {
			in = !in
		}
	}
	return in
}

func mapPolygons(m *maps.Map) []polygon {
	var sect maps.Polygons
	found := false
	for _, s := range m.Unknown {
		if s.Name == sect.MapSection() {
			if err := sect.UnmarshalBinary(s.Data); err == nil {
				found = true
			}
		}
	}
	if !found {
		return nil
	}
	byID := make(map[uint32]image.Point, len(sect.Points))
	for _, p := range sect.Points {
		byID[p.ID] = image.Pt(int(math.Round(float64(p.Pos.X))), int(math.Round(float64(p.Pos.Y))))
	}
	var out []polygon
	for _, p := range sect.Polygons {
		pl := polygon{color: p.AmbientLight}
		for _, id := range p.Points {
			if pt, ok := byID[id]; ok {
				pl.pts = append(pl.pts, pt)
			}
		}
		if len(pl.pts) < 3 {
			continue
		}
		pl.bounds = image.Rectangle{Min: pl.pts[0], Max: pl.pts[0]}
		for _, pt := range pl.pts {
			pl.bounds = pl.bounds.Union(image.Rectangle{Min: pt, Max: pt.Add(image.Pt(1, 1))})
		}
		out = append(out, pl)
	}
	return out
}

// light is a light source.
type light struct {
	pos       image.Point
	intensity int // 16.16
	radius    int
	color     [3]int
	subtract  bool
}

// lightRadius is the reach of a light of an intensity (the client's LightRadius).
func lightRadius(intens float32) int {
	if intens <= 1 {
		return 0
	}
	i2 := min(intens, 31)
	return int(math.Sqrt(float64((float32(4)+i2)/(float32(4)+1)+1) * float64(intens*intens/(float32(961)*float32(0.001)))))
}

func lightOfIntensity(v float32) (intens, rad int) {
	v = min(v, 63)
	return int(v*0x10000 + 0.5), lightRadius(v)
}

// objectLight returns the light the object gives, if any.
func objectLight(o *mapObject) (light, bool) {
	typ := o.typ
	l := light{pos: o.pos, color: [3]int{255, 255, 255}}
	var v float32
	if typ.LightIntensity != 0 {
		v = float32(typ.LightIntensity)
		if v < 0 {
			v, l.subtract = -v, true
		}
		if typ.LightColor != nil {
			l.color = [3]int{int(typ.LightColor.R), int(typ.LightColor.G), int(typ.LightColor.B)}
		}
	}
	l.intensity, l.radius = lightOfIntensity(v)
	// lights placed on the map store their own
	if il, ok := o.x.(*xfer.InvisibleLight); ok && len(il.Light) >= 36 {
		d := il.Light
		l.radius = int(binary.LittleEndian.Uint32(d[8:]))
		l.intensity = int(binary.LittleEndian.Uint32(d[12:]))
		for i := range l.color {
			l.color[i] = int(int32(binary.LittleEndian.Uint32(d[16+4*i:])))
		}
		l.subtract = il.Extra != 0
	}
	switch {
	case typ.class&object.ClassMonster != 0 && typ.class&object.ClassLight == 0:
		// creatures glow a little
		l.color = [3]int{255, 255, 255}
		switch {
		case typ.monster&object.MonsterSmall != 0:
			l.intensity, l.radius = lightOfIntensity(25)
		case typ.monster&object.MonsterMedium != 0:
			l.intensity, l.radius = lightOfIntensity(35)
		case typ.monster&object.MonsterLarge != 0:
			l.intensity, l.radius = lightOfIntensity(45)
		}
	case typ.class&object.ClassLight != 0 && o.flags&object.FlagEnabled != 0 && l.radius > 0:
	default:
		return light{}, false
	}
	if l.intensity <= lightMinIntensity {
		return light{}, false
	}
	return l, true
}

// add adds a round light (the client's nox_xxx_cliLight16_469140).
func (g *lightGrid) add(l light) {
	intens2 := min(l.intensity, lightMaxIntensity)
	xx, yy := l.pos.X-g.org.X, l.pos.Y-g.org.Y
	rad := l.radius
	const s = common.GridStep
	xmin, xmax := max((xx-rad)/s, 0), min((xx+rad)/s, g.w-1)
	ymin, ymax := max((yy-rad)/s, 0), min((yy+rad)/s, g.h-1)
	limit := rad * rad
	v36 := int32((l.intensity >> 16) * (l.intensity >> 16))
	if v36 == 0 {
		return
	}
	for y := ymin; y <= ymax; y++ {
		dy := yy - s*y
		for x := xmin; x <= xmax; x++ {
			dx := xx - s*x
			dist := dx*dx + dy*dy
			if dist > limit {
				continue
			}
			// 32-bit like the original
			v := fixDiv(intens2+lightBase, int(66*int32(dist)*lightFalloff/v36+0x10000))
			if v > lightBase {
				g.addAt(x, y, l.color, 8*(v-lightBase), l.subtract)
			}
		}
	}
}

func fixDiv(a, b int) int {
	return int(int32((uint64(a) << 16) / uint64(b)))
}

func fixMul(a, b int) int {
	return int(int32((int64(a) * int64(b)) >> 16))
}

func (g *lightGrid) addAt(x, y int, c [3]int, v int, subtract bool) {
	if subtract {
		v = -v
	}
	p := g.at(x, y)
	for i := range p {
		p[i] = max(0, min(p[i]+fixMul(v, c[i])<<8, lightFull))
	}
}

// lightMap computes the lights of the map and its objects.
func lightMap(m *maps.Map, objs []*mapObject, bounds image.Rectangle) *lightGrid {
	g := newLightGrid(bounds)
	g.ambient(m)
	for _, o := range objs {
		if o.flags&object.FlagEditVisible != 0 {
			continue
		}
		l, ok := objectLight(o)
		if !ok {
			continue
		}
		// lights placed on the map can be cones
		if il, ok := o.x.(*xfer.InvisibleLight); ok && len(il.Light) >= 36 && binary.LittleEndian.Uint32(il.Light[32:]) != 0xffff {
			g.addCone(l, binary.LittleEndian.Uint16(il.Light[28:]), binary.LittleEndian.Uint16(il.Light[30:]))
			continue
		}
		g.add(l)
	}
	return g
}

// lightFloor multiplies the image by the light, as the game does for the floor.
func lightFloor(out *image.RGBA, g *lightGrid) {
	b := out.Rect
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			i := out.PixOffset(x, y)
			p := out.Pix[i : i+3 : i+3]
			if p[0]|p[1]|p[2] == 0 {
				continue
			}
			l := g.sampleF(float64(x)+0.5, float64(y)+0.5)
			for k := range p {
				p[k] = byte(min(float64(p[k])*l[k], 255))
			}
		}
	}
}

// drawLit draws the image multiplied by a light for each column (from the left one to the right one, which are
// 0-255 for objects, as the game's color multiplier, and 0-256 for walls), with an alpha.
func drawLit(out *image.RGBA, r image.Rectangle, img image.Image, left, right [3]int, scale int, alpha byte) {
	src, ok := img.(*image.NRGBA)
	if !ok {
		src = image.NewNRGBA(img.Bounds())
		draw.Draw(src, src.Rect, img, img.Bounds().Min, draw.Src)
	}
	b := src.Rect
	w := max(b.Dx()-1, 1)
	for y := 0; y < b.Dy(); y++ {
		py := r.Min.Y + y
		if py < out.Rect.Min.Y || py >= out.Rect.Max.Y {
			continue
		}
		for x := 0; x < b.Dx(); x++ {
			px := r.Min.X + x
			if px < out.Rect.Min.X || px >= out.Rect.Max.X {
				continue
			}
			si := src.PixOffset(b.Min.X+x, b.Min.Y+y)
			c := src.Pix[si : si+4 : si+4]
			a := int(c[3]) * int(alpha) / 255
			if a == 0 {
				continue
			}
			di := out.PixOffset(px, py)
			d := out.Pix[di : di+4 : di+4]
			for i := 0; i < 3; i++ {
				m := left[i] + (right[i]-left[i])*x/w
				v := min(int(c[i])*m/scale, 255)
				d[i] = byte((v*a + int(d[i])*(255-a)) / 255)
			}
			d[3] = 0xff
		}
	}
}

// wallLight returns the light at the two ends of a wall, as the game samples it for the wall's direction.
func (g *lightGrid) wallLight(w *maps.Wall) (left, right [3]int) {
	x, y := common.GridStep*int(w.Pos.X), common.GridStep*int(w.Pos.Y)
	var a, b image.Point
	switch w.Dir {
	case 0, 3:
		a, b = image.Pt(x, y+23), image.Pt(x+23, y)
	case 1, 4:
		a, b = image.Pt(x, y), image.Pt(x+23, y+23)
	case 7:
		a, b = image.Pt(x, y), image.Pt(x+23, y)
	case 8:
		a, b = image.Pt(x+11, y+11), image.Pt(x-23, y-23)
	case 10:
		a, b = image.Pt(x, y+23), image.Pt(x+11, y+12)
	default:
		a, b = image.Pt(x, y+23), image.Pt(x+23, y+23)
	}
	return g.sample(a), g.sample(b)
}

// sinTable is the game's table of sines for 4096 angles, 4096 for 1.
var sinTable = func() (t [4096]int32) {
	for i := range t {
		t[i] = int32(int64(math.Sin(float64(i)*0.0015339808) * 4096))
	}
	return t
}()

// addCone adds a cone of light: the light's direction and penumbra are angles of 65536 for a turn. Like the game, it
// lights the triangle between the light and the points at its radius in its direction and at the penumbra from it
// (twice), with the math of the original game on a 32-bit system.
func (g *lightGrid) addCone(l light, dir, penumbra uint16) {
	const s = common.GridStep
	xx, yy := l.pos.X-g.org.X, l.pos.Y-g.org.Y
	a1 := image.Pt((xx<<16)/s, (yy<<16)/s)
	r := (l.radius << 16) / s
	at := func(ang uint16) image.Point {
		return a1.Add(image.Pt(
			fixMul(r, 16*int(sinTable[uint16(ang+0x4000)>>4])),
			fixMul(r, 16*int(sinTable[ang>>4])),
		))
	}
	a3 := at(dir)
	a2 := at(dir + penumbra)
	center := image.Pt(xx, yy)
	g.triangle(a1, a2, a3, center, l.intensity, l.color)
	g.triangle(a1, a3, a2, center, l.intensity, l.color)
}

// triangle lights the grid points in a triangle of 16.16 grid coordinates (the client's sub_4696B0).
func (g *lightGrid) triangle(a1, a2, a3, center image.Point, intens int, c [3]int) {
	var top, mid, bot image.Point
	if a1.Y > a2.Y {
		if a2.Y <= a3.Y {
			top, mid, bot = a2, a3, a1
		} else {
			top, mid, bot = a3, a1, a2
		}
	} else {
		if a1.Y <= a3.Y {
			top, mid, bot = a1, a2, a3
		} else {
			top, mid, bot = a3, a1, a2
		}
	}
	e1, e2, e3 := g.edge(top, mid), g.edge(top, bot), []int(nil)
	g.span(e1, e2, center, intens, c)
	if mid.Y < bot.Y {
		e3 = g.edge(mid, bot)
		g.span(e3, e2, center, intens, c)
	} else if mid.Y > bot.Y {
		e3 = g.edge(bot, mid)
		g.span(e1, e3, center, intens, c)
	}
}

// edge returns the first row of an edge, its number of rows, and the column of each row (the client's sub_484DC0).
func (g *lightGrid) edge(a, b image.Point) []int {
	out := make([]int, 2+g.h)
	out[0] = (a.Y + 0x8000) >> 16
	out[1] = ((b.Y + 0x8000) >> 16) - out[0] + 1
	if out[1] <= 0 {
		return out
	}
	dn := (b.X - a.X) / out[1]
	end := out[0] + out[1]
	if out[0] < 0 {
		out[0] = 0
		out[1] = end
	}
	if end > g.h {
		out[1] += g.h - end
		end = g.h
	}
	x := a.X + 0x8000
	for row := out[0]; row < end; row++ {
		out[row+2] = min(max(x>>16, 0), g.w-1)
		x += dn
	}
	return out
}

// span lights the grid points between two edges (the client's sub_4697C0).
func (g *lightGrid) span(left, right []int, center image.Point, intens int, c [3]int) {
	from := max(left[0], right[0])
	to := min(left[0]+left[1], right[0]+right[1])
	if to-from <= 1 {
		return
	}
	intens2 := min(intens, lightMaxIntensity)
	v21 := int32((intens >> 16) * (intens >> 16))
	if v21 == 0 {
		return
	}
	const s = common.GridStep
	dy := s*from - center.Y
	for row := from; row < to; row++ {
		dx := s*left[row+2] - center.X
		for col := left[row+2]; col < right[row+2]; col++ {
			v := fixDiv(lightBase+intens2, int(66*lightFalloff*int32(dy*dy+dx*dx)/v21+0x10000))
			if v > lightBase {
				g.addAt(col, row, c, 8*(v-lightBase), false)
			}
			dx += s
		}
		dy += s
	}
}
