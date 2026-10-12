package gpu

import (
	"encoding/binary"
	"errors"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
)

const spriteHeaderSize = 17

// errEmptySprite is returned for sprites without any pixels. They are valid, but not drawn.
var errEmptySprite = errors.New("empty sprite")

// decodedSprite is a sprite decoded to RGBA (non-premultiplied) and a material index layer.
type decodedSprite struct {
	W, H int
	RGBA []byte // 4 bytes per pixel
	// Mat contains material index + 1 for indexed (op 4) pixels, or 0 for regular pixels.
	Mat []byte
}

func decodeSprite(typ int, data []byte) (*decodedSprite, error) {
	if len(data) < spriteHeaderSize {
		return nil, errors.New("sprite is too short")
	}
	w := int(binary.LittleEndian.Uint32(data[0:]))
	h := int(binary.LittleEndian.Uint32(data[4:]))
	if w == 0 || h == 0 {
		return nil, errEmptySprite
	}
	if w < 0 || h < 0 || w > 8192 || h > 8192 {
		return nil, errors.New("invalid sprite size")
	}
	data = data[spriteHeaderSize:]
	d := &decodedSprite{
		W: w, H: h,
		RGBA: make([]byte, 4*w*h),
		Mat:  make([]byte, w*h),
	}
	switch typ {
	case 2, 7:
		if len(data) < 2*w*h {
			return nil, errors.New("unexpected end of sprite data")
		}
		for i := 0; i < w*h; i++ {
			d.setRGB555(i, binary.LittleEndian.Uint16(data[2*i:]))
		}
		return d, nil
	case 3, 4, 5, 6, 8:
	default:
		return nil, errors.New("unsupported sprite type")
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; {
			if len(data) < 2 {
				return nil, errors.New("unexpected end of sprite data")
			}
			op, n := data[0], int(data[1])
			data = data[2:]
			if n == 0 {
				return nil, errors.New("empty sprite span")
			}
			if x+n > w {
				n = w - x
			}
			i := y*w + x
			switch op & 0xF {
			case 1: // transparent
			case 2, 7:
				if len(data) < 2*n {
					return nil, errors.New("unexpected end of sprite data")
				}
				for j := 0; j < n; j++ {
					v := binary.LittleEndian.Uint16(data[2*j:])
					if typ == 8 {
						d.setRGB565(i+j, v)
					} else {
						d.setRGB555(i+j, v)
					}
				}
				data = data[2*n:]
			case 4: // grayscale with material index
				if len(data) < n {
					return nil, errors.New("unexpected end of sprite data")
				}
				mat := op >> 4
				for j := 0; j < n; j++ {
					c := data[j]
					p := d.RGBA[4*(i+j):]
					p[0], p[1], p[2], p[3] = c, c, c, 0xff
					d.Mat[i+j] = mat + 1
				}
				data = data[n:]
			case 5: // RGBA4444
				if len(data) < 2*n {
					return nil, errors.New("unexpected end of sprite data")
				}
				for j := 0; j < n; j++ {
					c := noxcolor.RGBA4444(binary.LittleEndian.Uint16(data[2*j:])).ColorNRGBA()
					p := d.RGBA[4*(i+j):]
					p[0], p[1], p[2], p[3] = c.R, c.G, c.B, c.A
				}
				data = data[2*n:]
			case 6: // unused, no-op
				if len(data) < 2*n {
					return nil, errors.New("unexpected end of sprite data")
				}
				data = data[2*n:]
			default:
				return nil, errors.New("invalid sprite op")
			}
			x += n
		}
	}
	return d, nil
}

func (d *decodedSprite) setRGB555(i int, v uint16) {
	c := noxcolor.RGBA5551(v).ColorNRGBA()
	p := d.RGBA[4*i:]
	p[0], p[1], p[2], p[3] = c.R, c.G, c.B, 0xff
}

func (d *decodedSprite) setRGB565(i int, v uint16) {
	// additive particles: colors are added to the destination
	p := d.RGBA[4*i:]
	p[0] = byte(v>>8)&0xF8 | 7
	p[1] = byte(v>>3)&0xFC | 3
	p[2] = byte(v<<3)&0xF8 | 7
	p[3] = 0xff
}
