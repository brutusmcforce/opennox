package gpu

import (
	"image"
	"math"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Solid fills (the edges of the line-of-sight shadows are thousands of them per frame) are drawn with
// their own compact instances: the rectangle as 16-bit integers and the color, 24 bytes instead of
// a 68-byte sprite instance. The shader computes exactly what the sprite shader does for a fill
// (a white texel without a material), so the pixels are the same.

const (
	solidSize    = 6 // uint32 words: x0|y0, x1|y1 (int16 each), r, g, b, a (float32 bits)
	maxSolidQuad = 8192
)

type solidRenderer struct {
	prog    uint32
	locSize int32
	vao     uint32
	vbo     uint32
	data    []uint32
}

func (s *solidRenderer) init() error {
	var err error
	s.prog, err = newProgram(solidVert, solidFrag)
	if err != nil {
		return err
	}
	s.locSize = gl.GetUniformLocation(s.prog, gl.Str("uSize\x00"))
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)
	gl.GenBuffers(1, &s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	const stride = solidSize * 4
	// shorts are converted to floats exactly
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 4, gl.SHORT, false, stride, 0)
	gl.VertexAttribDivisor(0, 1)
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 4, gl.FLOAT, false, stride, 2*4)
	gl.VertexAttribDivisor(2, 1)
	s.data = make([]uint32, 0, maxSolidQuad*solidSize)
	return nil
}

func (s *solidRenderer) free() {
	gl.DeleteProgram(s.prog)
	gl.DeleteVertexArrays(1, &s.vao)
	gl.DeleteBuffers(1, &s.vbo)
}

// fits reports if a rectangle can be stored in a solid instance.
func solidFits(rc image.Rectangle) bool {
	return rc.Min.X >= math.MinInt16 && rc.Min.Y >= math.MinInt16 && rc.Max.X <= math.MaxInt16 && rc.Max.Y <= math.MaxInt16
}

func (s *solidRenderer) add(rc image.Rectangle, c *[4]float32) {
	s.data = append(s.data,
		uint32(uint16(int16(rc.Min.X)))|uint32(uint16(int16(rc.Min.Y)))<<16,
		uint32(uint16(int16(rc.Max.X)))|uint32(uint16(int16(rc.Max.Y)))<<16,
		math.Float32bits(c[0]), math.Float32bits(c[1]), math.Float32bits(c[2]), math.Float32bits(c[3]),
	)
}

// draw draws the batch; the program, the vertex array and the buffer must be bound.
func (s *solidRenderer) draw() {
	gl.BufferData(gl.ARRAY_BUFFER, len(s.data)*4, unsafe.Pointer(&s.data[0]), gl.STREAM_DRAW)
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(len(s.data)/solidSize))
	s.data = s.data[:0]
}

// The same as spriteVert, with a white texel and modeNormal.
const solidVert = `#version 330 core
in vec4 aDst; // x0, y0, x1, y1
in vec4 aCol;
uniform vec2 uSize;
out vec4 vCol;
const vec2 corners[6] = vec2[6](vec2(0, 0), vec2(1, 0), vec2(0, 1), vec2(1, 0), vec2(1, 1), vec2(0, 1));
void main() {
	vec2 c = corners[gl_VertexID];
	vec2 pos = mix(aDst.xy, aDst.zw, c);
	vCol = aCol;
	vec2 p = pos / uSize * 2.0 - 1.0;
	gl_Position = vec4(p.x, -p.y, 0.0, 1.0);
}
`

// The same as shade in spriteFrag for a white texel (t = 1) without a material in modeNormal.
const solidFrag = `#version 330 core
in vec4 vCol;
out vec4 oColor;
void main() {
	vec3 c = vec3(1.0);
	c *= vCol.rgb;
	float a = 1.0 * vCol.a;
	oColor = vec4(c * a, a);
}
`
