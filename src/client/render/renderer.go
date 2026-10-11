package render

import (
	"image"
	"sync/atomic"

	"github.com/opennox/libs/client/seat"
	"github.com/opennox/libs/log"
	"github.com/opennox/libs/noximage"
)

var (
	Log = log.New("render")
)

func New(sc seat.Screen) (*Renderer, error) {
	sz := sc.ScreenSize()
	r := &Renderer{
		sc:         sc,
		fullscreen: -4, // unset
		filtering:  true,
	}
	Log.Printf("screen size: %dx%d", sz.X, sz.Y)
	sc.OnScreenResize(func(sz image.Point) {
		r.present(nil)
	})
	return r, nil
}

// Renderer manages the window and presents frames drawn by the GPU renderer.
type Renderer struct {
	ticks      uint32 // atomic
	sc         seat.Screen
	view       image.Rectangle
	fullscreen int
	borderless bool // only used for toggle
	filtering  bool
	stretch    bool
	rotate     bool
	rotated    bool
	onResize   []func(view image.Rectangle)

	presentFunc PresentFunc
	last        *noximage.Image16
}

// PresentFunc draws a render target associated with a given image to the screen.
// View is the area of the screen (top-left origin) and win is the screen size.
type PresentFunc func(img *noximage.Image16, view image.Rectangle, win image.Point)

// SetPresentFunc sets a function for presenting frames (the GPU renderer).
// Only the size of the image is used; its pixels are drawn by the GPU.
func (r *Renderer) SetPresentFunc(fnc PresentFunc) {
	r.presentFunc = fnc
}

// Ticks returns the number of present ticks since the last Reset or Init.
func (r *Renderer) Ticks() uint {
	return uint(atomic.LoadUint32(&r.ticks))
}

// Reset the renderer.
func (r *Renderer) Reset(sz image.Point) error {
	return nil
}

// OnViewResize registers a function that will be called when the render or the underlying window resizes.
func (r *Renderer) OnViewResize(fnc func(view image.Rectangle)) {
	r.onResize = append(r.onResize, fnc)
	fnc(r.view)
}

// CopyBuffer presents the render target associated with a given image.
func (r *Renderer) CopyBuffer(img *noximage.Image16) {
	r.present(img)
}

func (r *Renderer) present(img *noximage.Image16) {
	if r.presentFunc == nil {
		return
	}
	if img == nil {
		img = r.last
	}
	r.last = img
	var view image.Rectangle
	if img != nil {
		sz := img.Size()
		view = r.setViewport(float32(sz.X) / float32(sz.Y))
		if r.view != view {
			r.view = view
			for _, fnc := range r.onResize {
				fnc(view)
			}
		}
	}
	r.presentFunc(img, view, r.sc.ScreenSize())
	r.sc.Present()
	atomic.AddUint32(&r.ticks, 1)
}

func (r *Renderer) SetStretched(stretch bool) {
	r.stretch = stretch
}

func (r *Renderer) GetStretched() bool {
	return r.stretch
}

func (r *Renderer) SetFiltering(enable bool) {
	r.filtering = enable
}

func (r *Renderer) GetFiltering() bool {
	return r.filtering
}

func (r *Renderer) setViewport(ratio float32) image.Rectangle {
	wsz := r.sc.ScreenSize()
	if r.stretch {
		return image.Rect(0, 0, wsz.X, wsz.Y)
	}
	var (
		offx     = 0
		offy     = 0
		vpw, vph int
		vpx, vpy int
	)

	vpw, vph = wsz.X, wsz.Y

	// Maintain source aspect ratio
	if r.rotate && float32(vph)-ratio*float32(vpw) > float32(vpw)-ratio*float32(vph) {
		r.rotated = true
	} else {
		r.rotated = false
	}
	if r.rotated {
		ratio = 1.0 / ratio
	}

	if ratio*float32(vph) <= float32(vpw) {
		offx = int(float32(vpw)-float32(vph)*ratio) / 2
		vpw = int(float32(vph) * ratio)
	} else {
		offy = int(float32(vph)-float32(vpw)/ratio) / 2
		vph = int(float32(vpw) / ratio)
	}
	vpx = offx
	vpy = offy
	return image.Rect(vpx, vpy, vpx+vpw, vpy+vph)
}

// IsFullScreen checks if the renderer's window is set to fullscreen mode.
func (r *Renderer) IsFullScreen() bool {
	switch r.fullscreen {
	case -1, 1, -2, 2:
		return true
	}
	return false
}

// WindowMode returns the current window/fullscreen mode for the renderer's screen.
func (r *Renderer) WindowMode() int {
	return r.fullscreen
}

// SetWindowMode updates the window/fullscreen mode for the renderer's screen.
func (r *Renderer) SetWindowMode(mode int) {
	// TODO: previously it set the size of the window
	switch mode {
	case -1, 1:
		// Normal fullscreen
		r.sc.SetScreenMode(seat.Fullscreen)
	case -2, 2:
		// Borderless fullscreen
		r.sc.SetScreenMode(seat.Borderless)
	default:
		// Windowed
		r.sc.SetScreenMode(seat.Windowed)
	}
	r.fullscreen = mode
}

// ToggleWindowMode switches window/fullscreen mode for the renderer's screen.
func (r *Renderer) ToggleWindowMode() {
	switch r.fullscreen {
	case -1, 1:
		// Normal fullscreen -> Windowed
		r.borderless = false
		r.SetWindowMode(-3)
	case -2, 2:
		// Borderless fullscreen -> Windowed
		r.borderless = true
		r.SetWindowMode(-3)
	default:
		// Windowed -> last fullscreen
		if r.borderless {
			r.SetWindowMode(-2)
		} else {
			r.SetWindowMode(-1)
		}
	}
}
