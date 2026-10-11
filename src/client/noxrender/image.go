package noxrender

import (
	"image"
)

type Image16 interface {
	Type() int
	Pixdata() []byte
}

type rawImage struct {
	typ  int
	data []byte
}

func (img *rawImage) Type() int {
	return img.typ
}

func (img *rawImage) Pixdata() []byte {
	return img.data
}

func NewRawImage16(typ int, data []byte) Image16 {
	return &rawImage{typ: typ, data: data}
}

func (r *NoxRender) DrawImage16(img Image16, pos image.Point) {
	if img == nil || r.gpu == nil {
		return
	}
	if p, ok := img.(*Image); ok && p == nil {
		return
	}
	r.gpuDrawImage(img, pos)
}
