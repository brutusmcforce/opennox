package gpu

import (
	"image"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"

	"github.com/brutusmcforce/opennox/client/noxrender"
)

// floorRenderer draws floor tiles from the game's circular tile buffer.
type floorRenderer struct {
	prog  uint32
	loc   map[string]int32
	tiles uint32
	tsize image.Point
	light uint32
	lsize image.Point
	verts []float32
	// reused buffers for merging dirty cells (see uploadTiles)
	open, next []cellRect
}

// cellRect is a rectangle of tile buffer cells: columns [c0, c1) and rows [r0, r1).
type cellRect struct {
	c0, c1, r0, r1 int
}

func (f *floorRenderer) init() error {
	var err error
	f.prog, err = newProgram(floorVert, floorFrag)
	if err != nil {
		return err
	}
	f.loc = make(map[string]int32)
	for _, name := range []string{"uSize", "uPhys", "uTiles", "uLight", "uOff", "uPitch", "uLen", "uLightOff", "uStep", "uLightSize", "uLit"} {
		f.loc[name] = gl.GetUniformLocation(f.prog, gl.Str(name+"\x00"))
	}
	gl.GenTextures(1, &f.tiles)
	gl.BindTexture(gl.TEXTURE_2D, f.tiles)
	setTexParams(gl.NEAREST)
	gl.GenTextures(1, &f.light)
	gl.BindTexture(gl.TEXTURE_2D, f.light)
	setTexParams(gl.LINEAR)
	return nil
}

func (f *floorRenderer) free() {
	gl.DeleteProgram(f.prog)
	gl.DeleteTextures(1, &f.tiles)
	gl.DeleteTextures(1, &f.light)
}

func uploadTex(tex uint32, cur *image.Point, sz image.Point, ifmt int32, fmt, typ uint32, ptr unsafe.Pointer) {
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	if *cur != sz {
		gl.TexImage2D(gl.TEXTURE_2D, 0, ifmt, int32(sz.X), int32(sz.Y), 0, fmt, typ, ptr)
		*cur = sz
	} else {
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, int32(sz.X), int32(sz.Y), fmt, typ, ptr)
	}
}

// uploadTiles updates the tile texture: in full if its size changed or there are no dirty flags,
// otherwise only the cells that the game drew into since the last upload.
func (f *floorRenderer) uploadTiles(p *noxrender.FloorParams, rows int) {
	sz := image.Pt(p.Pitch, rows)
	const cell = noxrender.TileCell
	cols := p.DirtyCols
	if f.tsize != sz || p.Dirty == nil || cols <= 0 || len(p.Dirty) < cols*((rows+cell-1)/cell) {
		uploadTex(f.tiles, &f.tsize, sz, gl.R16UI, gl.RED_INTEGER, gl.UNSIGNED_SHORT, unsafe.Pointer(&p.Tiles[0]))
		clear(p.Dirty)
		return
	}
	gl.BindTexture(gl.TEXTURE_2D, f.tiles)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, int32(p.Pitch))
	upload := func(r cellRect) {
		x, y := r.c0*cell, r.r0*cell
		w, h := min((r.c1-r.c0)*cell, p.Pitch-x), min((r.r1-r.r0)*cell, rows-y)
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, int32(x), int32(y), int32(w), int32(h),
			gl.RED_INTEGER, gl.UNSIGNED_SHORT, unsafe.Pointer(&p.Tiles[y*p.Pitch+x]))
	}
	// Runs of dirty cells in a row are merged with the same runs in the rows above,
	// so a new column of tiles (when scrolling sideways) is a single upload. Each upload has a fixed cost.
	open, next := f.open[:0], f.next[:0]
	for cy := 0; cy*cell < rows; cy++ {
		row := p.Dirty[cy*cols : (cy+1)*cols]
		next = next[:0]
		k := 0 // open runs are sorted by column, like the runs of this row
		for c0 := 0; c0 < cols; c0++ {
			if row[c0] == 0 {
				continue
			}
			c1 := c0 + 1
			for c1 < cols && row[c1] != 0 {
				c1++
			}
			clear(row[c0:c1])
			for k < len(open) && open[k].c0 < c0 {
				upload(open[k])
				k++
			}
			if k < len(open) && open[k].c0 == c0 && open[k].c1 == c1 {
				r := open[k]
				r.r1 = cy + 1
				next = append(next, r)
				k++
			} else {
				next = append(next, cellRect{c0: c0, c1: c1, r0: cy, r1: cy + 1})
			}
			c0 = c1
		}
		for ; k < len(open); k++ {
			upload(open[k])
		}
		open, next = next, open
	}
	for _, r := range open {
		upload(r)
	}
	f.open, f.next = open[:0], next[:0]
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
}

// Floor draws floor tiles in given screen spans (rectangles).
func (r *Renderer) Floor(p *noxrender.FloorParams, spans []image.Rectangle) {
	if r.cur == nil || len(spans) == 0 || p.Pitch <= 0 || len(p.Tiles) < p.Pitch {
		return
	}
	f := &r.floor
	rows := len(p.Tiles) / p.Pitch
	r.flush()
	f.uploadTiles(p, rows)
	if p.Lit && len(p.Light) >= 3*p.LightW*p.LightH && p.LightW > 0 && p.LightH > 0 {
		uploadTex(f.light, &f.lsize, image.Pt(p.LightW, p.LightH), gl.RGB32F, gl.RGB, gl.FLOAT, unsafe.Pointer(&p.Light[0]))
	}

	f.verts = f.verts[:0]
	for _, s := range spans {
		x0, y0, x1, y1 := float32(s.Min.X), float32(s.Min.Y), float32(s.Max.X), float32(s.Max.Y)
		f.verts = append(f.verts,
			x0, y0, 0, 0, 0, 0, 0, 0, 0,
			x1, y0, 0, 0, 0, 0, 0, 0, 0,
			x0, y1, 0, 0, 0, 0, 0, 0, 0,
			x1, y0, 0, 0, 0, 0, 0, 0, 0,
			x1, y1, 0, 0, 0, 0, 0, 0, 0,
			x0, y1, 0, 0, 0, 0, 0, 0, 0,
		)
	}
	t := r.cur
	t.used = r.frame
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.Viewport(0, 0, int32(t.size.X), int32(t.size.Y))
	gl.UseProgram(f.prog)
	gl.BindVertexArray(r.vao)
	// spans are in logical coordinates, while the target may be larger or smaller (see SetTargetSize)
	gl.Uniform2f(f.loc["uSize"], float32(t.log.X), float32(t.log.Y))
	gl.Uniform2f(f.loc["uPhys"], float32(t.size.X), float32(t.size.Y))
	gl.Uniform1i(f.loc["uTiles"], 0)
	gl.Uniform1i(f.loc["uLight"], 1)
	gl.Uniform2i(f.loc["uOff"], int32(p.Off.X), int32(p.Off.Y))
	gl.Uniform1i(f.loc["uPitch"], int32(p.Pitch))
	gl.Uniform1i(f.loc["uLen"], int32(rows*p.Pitch))
	gl.Uniform2i(f.loc["uLightOff"], int32(p.LightOff.X), int32(p.LightOff.Y))
	gl.Uniform1f(f.loc["uStep"], float32(max(p.Step, 1)))
	gl.Uniform2f(f.loc["uLightSize"], float32(max(f.lsize.X, 1)), float32(max(f.lsize.Y, 1)))
	lit := int32(0)
	if p.Lit && f.lsize.X > 0 {
		lit = 1
	}
	gl.Uniform1i(f.loc["uLit"], lit)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, f.tiles)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, f.light)
	gl.Disable(gl.BLEND)
	gl.Disable(gl.SCISSOR_TEST)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(f.verts)*4, gl.Ptr(f.verts), gl.STREAM_DRAW)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(f.verts)/vertSize))
	r.restore()
}

const floorVert = `#version 330 core
in vec2 aPos;
uniform vec2 uSize;
void main() {
	vec2 p = aPos / uSize * 2.0 - 1.0;
	gl_Position = vec4(p.x, -p.y, 0.0, 1.0);
}
`

const floorFrag = `#version 330 core
uniform usampler2D uTiles;
uniform sampler2D uLight;
uniform vec2 uSize; // logical size of the target
uniform vec2 uPhys; // physical size of the target
uniform ivec2 uOff;
uniform int uPitch;
uniform int uLen;
uniform ivec2 uLightOff;
uniform float uStep;
uniform vec2 uLightSize;
uniform int uLit;
out vec4 oColor;
void main() {
	vec2 lp = vec2(gl_FragCoord.x, uPhys.y - gl_FragCoord.y) * uSize / uPhys;
	int x = int(lp.x);
	int y = int(lp.y);
	int idx = uPitch * (y + uOff.y) + x + uOff.x;
	idx = idx % uLen;
	if (idx < 0) {
		idx += uLen;
	}
	uint v = texelFetch(uTiles, ivec2(idx % uPitch, idx / uPitch), 0).r;
	vec3 c = vec3(float((v >> 10) & 31u), float((v >> 5) & 31u), float(v & 31u)) * (8.0 / 255.0);
	if (uLit != 0) {
		vec2 g = (vec2(x, y) + vec2(uLightOff)) / uStep;
		vec3 l = texture(uLight, (g + 0.5) / uLightSize).rgb;
		c = min(c * l, vec3(1.0));
	}
	oColor = vec4(c, 1.0);
}
`
