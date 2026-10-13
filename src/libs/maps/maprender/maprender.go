package maprender

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"path/filepath"
	"sort"

	"github.com/brutusmcforce/opennox/libs/bag"
	"github.com/brutusmcforce/opennox/libs/common"
	"github.com/brutusmcforce/opennox/libs/maps"
	"github.com/brutusmcforce/opennox/libs/modifiers"
	"github.com/brutusmcforce/opennox/libs/noximage/pcx"
	"github.com/brutusmcforce/opennox/libs/things"
)

// DrawMap is a helper for drawing a single map. For multiple maps it's more efficient to use Renderer.
func DrawMap(datadir string, m *maps.Map, opts *Options) (*image.RGBA, error) {
	r, err := NewRenderer(datadir)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return r.DrawMap(m, opts)
}

// NewRenderer creates a Nox map renderer based on a given Nox data dir.
// The renderer can be sued to draw multiple maps.
func NewRenderer(datadir string) (*Renderer, error) {
	tng, err := things.Open(filepath.Join(datadir, "thing.bin"))
	if err != nil {
		return nil, err
	}
	bg, err := bag.Open(filepath.Join(datadir, "video.bag"))
	if err != nil {
		_ = tng.Close()
		return nil, err
	}
	r := &Renderer{tng: tng, bag: bg}
	if err = r.indexThingWalls(); err != nil {
		_ = r.Close()
		return nil, err
	}
	if err = r.indexBagImages(); err != nil {
		_ = r.Close()
		return nil, err
	}
	// the colors of items; drawn white without it
	if mods, err := modifiers.ReadFile(filepath.Join(datadir, "modifier.bin")); err == nil {
		r.mods = mods
	}
	if err = r.indexThings(); err != nil {
		_ = r.Close()
		return nil, err
	}
	if r.floors, err = tng.ReadFloors(); err != nil {
		_ = r.Close()
		return nil, err
	}
	if r.edges, err = tng.ReadEdges(); err != nil {
		_ = r.Close()
		return nil, err
	}
	return r, nil
}

// Renderer is a Nox map renderer.
type Renderer struct {
	tng        *things.Reader
	bag        *bag.File
	wallByMat  map[int]*things.Wall
	floors     []things.Floor
	things     map[string]*thingInfo
	mods       *modifiers.File // nil if missing
	mats       [16]color.RGBA  // the current material colors
	light      *lightGrid      // the lights of the map drawn; nil to draw everything fully lit
	edges      []things.Edge
	images     []*bag.ImageRec
	imageByInd map[int]*pcx.Image
}

// Close the renderer,
func (r *Renderer) Close() error {
	_ = r.tng.Close()
	_ = r.bag.Close()
	return nil
}

func (r *Renderer) indexThingWalls() error {
	wallvar, err := r.tng.ReadWalls()
	if err != nil {
		return err
	}
	r.wallByMat = make(map[int]*things.Wall)
	for i := range wallvar {
		r.wallByMat[i] = &wallvar[i]
	}
	return nil
}

func (r *Renderer) indexBagImages() error {
	imgs, err := r.bag.Images()
	if err != nil {
		return err
	}
	r.images = imgs
	r.imageByInd = make(map[int]*pcx.Image)
	return nil
}

func (r *Renderer) getPCX(ind int) (*pcx.Image, error) {
	if ind < 0 || ind >= len(r.images) {
		return nil, fmt.Errorf("invalid image index: %d", ind)
	}
	img, ok := r.imageByInd[ind]
	if !ok {
		var err error
		img, err = r.images[ind].Decode()
		if err != nil {
			return nil, err
		}
		r.imageByInd[ind] = img
	}
	return img, nil
}

func (r *Renderer) getImage(ind int) (image.Image, image.Point, error) {
	img, err := r.getPCX(ind)
	if err != nil {
		return nil, image.Point{}, err
	}
	return img.Image, img.Point, nil
}

// getFloorImage returns the image of a floor tile variant.
func (r *Renderer) getFloorImage(tile byte, variant uint16) (*pcx.Image, error) {
	if int(tile) >= len(r.floors) {
		return nil, fmt.Errorf("unknown floor tile: %d", tile)
	}
	fl := &r.floors[tile]
	if int(variant) >= len(fl.Images) {
		return nil, fmt.Errorf("floor tile variant out of bounds: %s %d", fl.Name, variant)
	}
	return r.getPCX(fl.Images[variant].Ind)
}

// getTileImage returns the image of a floor tile with its edges, which blend the neighboring tiles into it.
func (r *Renderer) getTileImage(t *maps.Tile) (*image.NRGBA, error) {
	base, err := r.getFloorImage(t.Image, t.Variant)
	if err != nil {
		return nil, err
	}
	out := image.NewNRGBA(base.Bounds())
	draw.Draw(out, out.Rect, base, base.Bounds().Min, draw.Src)
	for _, e := range t.Edges {
		cover, err := r.getFloorImage(e.Image, e.Variant)
		if err != nil {
			return out, err
		}
		if int(e.Edge) >= len(r.edges) {
			return out, fmt.Errorf("unknown edge: %d", e.Edge)
		}
		ed := &r.edges[e.Edge]
		if int(e.Dir) >= len(ed.Images) {
			return out, fmt.Errorf("edge image out of bounds: %s %d", ed.Name, e.Dir)
		}
		mask, err := r.getPCX(ed.Images[e.Dir].Ind)
		if err != nil {
			return out, err
		}
		if mask.Cover == nil {
			return out, fmt.Errorf("not an edge image: %s %d", ed.Name, e.Dir)
		}
		// the edge's own pixels, then the neighboring tile's
		draw.Draw(out, out.Rect, mask.Image, image.Point{}, draw.Over)
		draw.DrawMask(out, out.Rect, cover.Image, image.Point{}, mask.Cover, image.Point{}, draw.Over)
	}
	return out, nil
}

// wallFlags are the wall flags that change its image; the map stores them in sections of their own.
type wallFlags struct {
	window    bool
	breakable bool
}

// getWallImage returns the wall's image and the offset of its top-left corner from the wall's grid position in
// pixels, or a nil image if the game draws none: like the game, for a variant the wall's direction does not have.
func (r *Renderer) getWallImage(w *maps.Wall, f wallFlags) (image.Image, image.Point, error) {
	wt := r.wallByMat[int(w.Material)]
	if wt == nil {
		return nil, image.Point{}, fmt.Errorf("unknown wall material: %d", w.Material)
	}
	dir := int(w.Dir)
	if dir >= len(wt.Directions) {
		return nil, image.Point{}, fmt.Errorf("invalid wall direction: %d", w.Dir)
	}
	// windows have images of their own in the two straight directions
	if f.window {
		switch dir {
		case 0:
			dir = 11
		case 1:
			dir = 12
		}
	}
	wd := wt.Directions[dir]
	if int(w.Variant) >= len(wd.Variants) {
		return nil, image.Point{}, nil
	}
	// the images are: normal, in front of the player, and both again for breakable walls
	layer := 0
	if f.breakable {
		layer = 2
	}
	wi := wd.Variants[w.Variant].Images[layer]
	ind := wi.Img.Ind
	if ind < 0 || ind >= len(r.images) {
		return nil, image.Point{}, fmt.Errorf("unsupported wall image: %+v", wi.Img)
	}
	img, pt, err := r.getImage(ind)
	if err != nil {
		return nil, image.Point{}, err
	}
	// as the game draws it
	pt = pt.Add(image.Pt(wi.Pt.X-50, -wi.Pt.Y-72))
	return img, pt, nil
}

type Options struct {
	FailFast bool // return on the first error
	NoFloor  bool // do not draw floor tiles
	NoWalls  bool // do not draw walls
	// NoObjects disables drawing the objects
	NoObjects bool
	// NoLight draws everything fully lit
	NoLight bool
}

// DrawMapFile reads and renders the map file. See DrawMap for details.
// As opposed to DrawMap, the returned image may be nil in case map decoding fails.
func (r *Renderer) DrawMapFile(dir string, opts *Options) (*image.RGBA, error) {
	m, err := maps.ReadMap(dir)
	if err != nil {
		return nil, err
	}
	return r.DrawMap(m, opts)
}

// DrawMap renders the map. It will keep processing the map and return a partial image in case of an error.
// If FailFast option is set, it will fail on the first error instead of returning the last one.
// Passing nil options will use defaults.
func (r *Renderer) DrawMap(m *maps.Map, opts *Options) (*image.RGBA, error) {
	if opts == nil {
		opts = &Options{}
	}
	bb := m.GridBoundingBox()

	// give some space for the wall images, which reach up to 3 cells left of and 4 above their position (and 2 right
	// of and 2 below it), and for the floor tiles, which reach 3 cells below theirs
	bb.Min = bb.Min.Add(image.Pt(-3, -4))
	bb.Max = bb.Max.Add(image.Pt(3, 3))

	// convert grid coordinates to pixels
	bb.Min.X *= common.GridStep
	bb.Min.Y *= common.GridStep
	bb.Max.X *= common.GridStep
	bb.Max.Y *= common.GridStep

	var objs []*mapObject
	if !opts.NoObjects {
		var err error
		objs, err = r.readObjects(m)
		if err != nil {
			return nil, err
		}
		// objects can be far from the floor and the walls
		bb = bb.Union(r.objectsBounds(objs))
	}

	r.resetMaterials()
	r.light = nil
	if !opts.NoLight {
		r.light = lightMap(m, objs, bb)
	}
	out := image.NewRGBA(bb)
	// black background
	draw.Draw(out, out.Rect, image.NewUniform(color.Black), image.Pt(0, 0), draw.Src)

	var last error
	failed := func(err error) bool {
		if err == nil {
			return false
		}
		last = err
		return opts.FailFast
	}
	if !opts.NoFloor {
		if failed(r.drawTiles(out, m, opts)) {
			return out, last
		}
	}
	// the game draws the objects in passes, and the walls sorted by Y among the other objects
	var sorted []*mapObject
	if !opts.NoObjects {
		for _, pass := range []int{passFloor, passBelow, passAbove} {
			for _, o := range objs {
				if o.pass() == pass && failed(r.drawObject(out, o)) {
					return out, last
				}
			}
			if pass == passFloor && r.light != nil {
				// the floor, with the objects drawn as a part of it
				lightFloor(out, r.light)
			}
		}
		for _, o := range objs {
			if o.pass() == passSorted {
				sorted = append(sorted, o)
			}
		}
		sortObjects(sorted)
	}
	var walls []maps.Wall
	if !opts.NoWalls && m.Walls != nil {
		walls = append(walls, m.Walls.Walls...)
		// with no player to see them from
		sort.SliceStable(walls, func(i, j int) bool { return walls[i].Pos.Y < walls[j].Pos.Y })
	}
	if opts.NoObjects && r.light != nil {
		lightFloor(out, r.light)
	}
	flags := wallFlagsOf(m)
	for len(walls) > 0 || len(sorted) > 0 {
		if len(walls) > 0 && (len(sorted) == 0 || sorted[0].sortY() >= int(walls[0].Pos.Y)*common.GridStep+11) {
			if failed(r.drawWall(out, &walls[0], flags)) {
				return out, last
			}
			walls = walls[1:]
			continue
		}
		if failed(r.drawObject(out, sorted[0])) {
			return out, last
		}
		sorted = sorted[1:]
	}
	return out, last
}

func (r *Renderer) drawTiles(out draw.Image, m *maps.Map, opts *Options) error {
	if m.Floor == nil {
		return nil
	}
	var last error
	for _, p := range m.Floor.Tiles {
		for _, it := range []struct {
			pos  maps.FloorPos
			tile *maps.Tile
		}{{p.LeftPos(), p.L}, {p.RightPos(), p.R}} {
			if it.tile == nil {
				continue
			}
			img, err := r.getTileImage(it.tile)
			if err != nil {
				last = err
				if opts.FailFast {
					return last
				}
				continue
			}
			// as the game draws it (a right tile of the first row is at -1)
			pt := image.Pt(int(int16(it.pos.X))*common.GridStep-11, int(int16(it.pos.Y))*common.GridStep+12)
			draw.Draw(out, img.Rect.Add(pt), img, image.Point{}, draw.Over)
		}
	}
	return last
}

func wallFlagsOf(m *maps.Map) map[image.Point]wallFlags {
	flags := make(map[image.Point]wallFlags)
	if m.WindowWalls != nil {
		for _, w := range m.WindowWalls.Walls {
			f := flags[w.Pos]
			f.window = true
			flags[w.Pos] = f
		}
	}
	if m.DestructableWalls != nil {
		for _, w := range m.DestructableWalls.Walls {
			f := flags[w.Pos]
			f.breakable = true
			flags[w.Pos] = f
		}
	}
	return flags
}

func (r *Renderer) drawWall(out *image.RGBA, w *maps.Wall, flags map[image.Point]wallFlags) error {
	img, pt, err := r.getWallImage(w, flags[image.Pt(int(w.Pos.X), int(w.Pos.Y))])
	if err != nil || img == nil {
		return err
	}
	p := image.Pt(
		int(w.Pos.X)*common.GridStep+pt.X,
		int(w.Pos.Y)*common.GridStep+pt.Y,
	)
	rect := img.Bounds()
	// wall lights are 8.8 fixed point multipliers
	left, right := [3]int{256, 256, 256}, [3]int{256, 256, 256}
	if r.light != nil {
		left, right = r.light.wallLight(w)
	}
	drawLit(out, image.Rect(p.X, p.Y, p.X+rect.Dx(), p.Y+rect.Dy()), img, left, right, 256, 0xff)
	return nil
}
