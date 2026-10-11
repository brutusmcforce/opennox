package opennox

import (
	"image"
	"unsafe"

	"github.com/opennox/libs/noximage"

	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// UI scaling.
//
// When UI scale is above 100%, the game uses two pixel buffers:
//   - The main pixel buffer (noxPixBuffer) has a reduced size (window size / scale).
//     It is used for the GUI, and all the logical screen sizes (nox_win_width, mouse position, window layout)
//     refer to this buffer. Thus, the GUI is drawn exactly as before, but at a lower resolution.
//   - The world buffer has the full window size. The game world is drawn into it.
//
// While the world is being drawn, the world buffer temporarily replaces the main pixel buffer
// and the logical screen size is set to the world size (see beginWorldDraw).
//
// After the world is drawn, it is downscaled into the GUI buffer, so that translucent GUI elements
// can be blended with it. Then the GUI is drawn on top. Finally, all GUI buffer pixels that changed
// are upscaled back into the world buffer, which is presented on the screen (see composeUI).
var uiScaleState struct {
	worldSize image.Point // requested size for the world buffer; zero if UI scaling is disabled

	world      *noximage.Image16
	worldFree  func()
	worldRows  []*uint16
	worldRFree func()

	inWorld bool
	saved   struct {
		img   *noximage.Image16
		rows  []*uint16
		win   image.Point
		clip  bool
		rect1 image.Rectangle
		rect2 image.Rectangle
		rect3 image.Rectangle
	}

	snap    []uint16 // GUI buffer contents before the GUI was drawn
	hasSnap bool

	// cursor is drawn at full resolution, so that it is not scaled with the GUI
	cursorPos image.Point // in GUI coordinates
	hasCursor bool
}

// videoUIScaled checks if the world and the GUI are rendered at different resolutions.
func videoUIScaled() bool {
	return uiScaleState.world != nil
}

// videoSetWorldSize sets the size of the world buffer. It will be applied on the next pixel buffer init.
// Zero size or size equal to the main buffer disables separate world buffer.
func videoSetWorldSize(sz image.Point) {
	uiScaleState.worldSize = sz
}

// videoWorldSize returns the size of the screen area used to draw the game world.
func videoWorldSize() image.Point {
	if w := uiScaleState.world; w != nil {
		return w.Size()
	}
	return videoGetWindowSize()
}

// uiToWindow converts screen position from GUI coordinates to the full resolution screen coordinates.
func uiToWindow(p image.Point) image.Point {
	w := uiScaleState.world
	if w == nil || noxPixBuffer.img == nil {
		return p
	}
	ui := uiScaleState.saved.img
	if !uiScaleState.inWorld {
		ui = noxPixBuffer.img
	}
	usz, wsz := ui.Size(), w.Size()
	if usz.X <= 0 || usz.Y <= 0 {
		return p
	}
	return image.Point{
		X: p.X * wsz.X / usz.X,
		Y: p.Y * wsz.Y / usz.Y,
	}
}

// uiToWorldScreen converts screen position from GUI coordinates to world buffer coordinates.
func uiToWorldScreen(p image.Point) image.Point {
	return noxClient.windowToWorldScreen(uiToWindow(p))
}

func init() {
	legacy.GetDrawViewport = func() *noxrender.Viewport {
		return noxClient.drawViewport()
	}
	legacy.WorldToDrawScreen = func(p image.Point) image.Point {
		return noxClient.worldToDrawScreen(p)
	}
}

// drawingUI checks if the GUI is currently drawn at a different resolution than the world.
func drawingUI() bool {
	return uiScaleState.world != nil && !uiScaleState.inWorld && noxPixBuffer.img != nil
}

// windowToUI converts screen position from the full resolution screen coordinates to GUI coordinates.
// It must only be called when drawingUI returns true.
func windowToUI(p image.Point) image.Point {
	usz, wsz := noxPixBuffer.img.Size(), uiScaleState.world.Size()
	if wsz.X <= 0 || wsz.Y <= 0 {
		return p
	}
	return image.Point{X: p.X * usz.X / wsz.X, Y: p.Y * usz.Y / wsz.Y}
}

// worldToDrawScreen converts world position to a screen position in the coordinates of the current render target.
func (c *Client) worldToDrawScreen(p image.Point) image.Point {
	p = c.Viewport().ToScreenPos(p)
	if drawingUI() {
		p = windowToUI(c.worldScreenToWindow(p))
	}
	return p
}

// uiViewport is a copy of the client viewport, scaled to GUI coordinates.
var uiViewport *noxrender.Viewport

// drawViewport returns the viewport for the current render target.
//
// The client viewport is in world buffer coordinates. Legacy GUI code that positions things relative to it
// (e.g. centered messages) must use a scaled copy when drawing to the GUI buffer.
func (c *Client) drawViewport() *noxrender.Viewport {
	vp := c.Viewport()
	if !drawingUI() {
		return vp
	}
	if uiViewport == nil {
		uiViewport, _ = alloc.New(noxrender.Viewport{})
	}
	*uiViewport = *vp
	uiViewport.Screen = image.Rectangle{Min: windowToUI(vp.Screen.Min), Max: windowToUI(vp.Screen.Max)}
	uiViewport.Size = windowToUI(vp.Size)
	return uiViewport
}

// GetMousePos returns the mouse position in the coordinates of the current render target.
//
// This overrides client.Client.GetMousePos, which is used by the legacy code.
func (c *Client) GetMousePos() image.Point {
	p := c.Client.GetMousePos()
	if uiScaleState.inWorld {
		p = uiToWorldScreen(p)
	}
	return p
}

// sub_473970_screen converts screen position in the coordinates of the current render target to world position.
//
// It's used by the legacy code, which may call it both in GUI and world drawing code.
func sub_473970_screen(p image.Point) image.Point {
	if !uiScaleState.inWorld {
		p = uiToWorldScreen(p)
	}
	return sub_473970(p)
}

func videoInitWorldBuffer(ui image.Point) {
	videoFreeWorldBuffer()
	sz := uiScaleState.worldSize
	if sz.X <= 0 || sz.Y <= 0 || sz == ui {
		return
	}
	videoLog.Printf("initializing world buffer: %v (ui: %v)", sz, ui)
	videoAllocWorldBuffer(sz)
}

// videoNeedWorldBuffer creates or frees the world buffer when UI scaling is disabled.
// In this case, the world buffer has the same size as the GUI buffer, and it is only needed for zooming.
func videoNeedWorldBuffer(need bool) {
	st := &uiScaleState
	if st.worldSize != (image.Point{}) || st.inWorld || noxPixBuffer.img == nil || need == (st.world != nil) {
		return
	}
	if need {
		videoAllocWorldBuffer(noxPixBuffer.img.Size())
	} else {
		videoFreeWorldBuffer()
	}
}

func videoAllocWorldBuffer(sz image.Point) {
	data, free := alloc.Make([]uint16{}, sz.X*sz.Y)
	img := noximage.NewImage16WithData(data, sz)
	rows, freeRows := alloc.Make([]*uint16{}, sz.Y)
	for y := 0; y < sz.Y; y++ {
		rows[y] = &img.Row(y)[0]
	}
	uiScaleState.world = img
	if g := noxClient.gpu; g != nil && noxPixBuffer.img != nil {
		// draw the GUI at the full resolution; it allows drawing smooth text
		g.SetTargetSize(noxPixBuffer.img, sz)
	}
	uiScaleState.worldFree = free
	uiScaleState.worldRows = rows
	uiScaleState.worldRFree = freeRows
	uiScaleState.snap = make([]uint16, len(noxPixBuffer.img.Pix))
	uiScaleState.hasSnap = false
}

func videoFreeWorldBuffer() {
	if uiScaleState.inWorld {
		noxClient.endWorldDraw()
	}
	if g := noxClient.gpu; g != nil && uiScaleState.world != nil && noxPixBuffer.img != nil {
		g.SetTargetSize(noxPixBuffer.img, image.Point{})
	}
	if uiScaleState.worldFree != nil {
		uiScaleState.worldFree()
	}
	if uiScaleState.worldRFree != nil {
		uiScaleState.worldRFree()
	}
	uiScaleState.world = nil
	uiScaleState.worldFree = nil
	uiScaleState.worldRows = nil
	uiScaleState.worldRFree = nil
	uiScaleState.snap = nil
	uiScaleState.hasSnap = false
}

func (c *Client) setPixTarget(img *noximage.Image16, rows []*uint16) {
	noxPixBuffer.img = img
	noxPixBuffer.rows = rows
	legacy.Set_nox_pixbuffer_rows_3798784(unsafe.Pointer(&rows[0]))
	c.r.SetPixBuffer(img)
}

// beginWorldDraw switches the render target to the full-resolution world buffer.
// It is a no-op if UI scaling is disabled.
func (c *Client) beginWorldDraw() {
	st := &uiScaleState
	if st.world == nil || st.inWorld {
		return
	}
	st.inWorld = true
	d := c.r.Data()
	sv := &st.saved
	sv.img, sv.rows = noxPixBuffer.img, noxPixBuffer.rows
	sv.win = videoGetWindowSize()
	sv.clip = d.Clip()
	sv.rect1, sv.rect2, sv.rect3 = d.ClipRect(), d.ClipRect2(), d.Rect3()

	sz := st.world.Size()
	c.setPixTarget(st.world, st.worldRows)
	videoSetWindowSize(sz)
	d.SetClipRect(image.Rectangle{Max: sz})
	d.SetClipRect2(image.Rectangle{Max: image.Pt(sz.X-1, sz.Y-1)})
	d.SetRect3(image.Rectangle{Max: sz})
}

// endWorldDraw switches the render target back to the GUI buffer.
func (c *Client) endWorldDraw() {
	st := &uiScaleState
	if !st.inWorld {
		return
	}
	st.inWorld = false
	sv := &st.saved
	c.setPixTarget(sv.img, sv.rows)
	videoSetWindowSize(sv.win)
	d := c.r.Data()
	d.SetClip(sv.clip)
	d.SetClipRect(sv.rect1)
	d.SetClipRect2(sv.rect2)
	d.SetRect3(sv.rect3)
	sv.img, sv.rows = nil, nil
}

// worldDrawn must be called after the world is drawn into the world buffer.
// It downscales the world into the GUI buffer and remembers the result,
// so that later we can detect which pixels were drawn by the GUI.
func (c *Client) worldDrawn() {
	st := &uiScaleState
	if st.world == nil || st.inWorld {
		return
	}
	st.hasSnap = true
	c.applyZoom(st.world)
	if c.gpu != nil {
		// GUI is drawn into a separate transparent target, so we don't need the world under it
		c.gpu.ClearTransparent()
		return
	}
	ui := noxPixBuffer.img
	scaleNearest16(ui, st.world)
	copy(st.snap, ui.Pix)
}

// composeUI copies all the GUI pixels drawn since worldDrawn into the world buffer, upscaling them.
func (c *Client) composeUI() {
	st := &uiScaleState
	if st.world == nil || st.inWorld || !st.hasSnap {
		return
	}
	if c.gpu != nil {
		c.gpu.Composite(st.world, noxPixBuffer.img)
		return
	}
	composeChanged16(st.world, noxPixBuffer.img, st.snap)
}

// deferCursorDraw remembers the cursor position, if the cursor must be drawn at full resolution.
// It returns false if the cursor should be drawn right away.
func (c *Client) deferCursorDraw(pos image.Point) bool {
	st := &uiScaleState
	if st.world == nil || st.inWorld {
		return false
	}
	st.cursorPos, st.hasCursor = pos, true
	return true
}

// drawDeferredCursor draws the cursor remembered by deferCursorDraw. It must be called during the world draw.
func (c *Client) drawDeferredCursor() {
	st := &uiScaleState
	if !st.hasCursor {
		return
	}
	st.hasCursor = false
	c.nox_video_cursorDrawImpl_477A30(uiToWindow(st.cursorPos))
}

// presentScaled draws final full-screen effects on the world buffer and presents it.
// It returns false if UI scaling is disabled.
func (c *Client) presentScaled(draw func()) bool {
	st := &uiScaleState
	if st.world == nil || st.inWorld || !st.hasSnap {
		return false
	}
	c.composeUI()
	c.beginWorldDraw()
	defer c.endWorldDraw()
	draw()
	c.copyPixBuffer()
	return true
}

// scaleNearest16 scales src image into dst using nearest neighbour sampling.
func scaleNearest16(dst, src *noximage.Image16) {
	dsz, ssz := dst.Size(), src.Size()
	if dsz.X <= 0 || dsz.Y <= 0 {
		return
	}
	xs := make([]int, dsz.X)
	for x := range xs {
		xs[x] = (2*x + 1) * ssz.X / (2 * dsz.X)
	}
	for y := 0; y < dsz.Y; y++ {
		srow := src.Row((2*y + 1) * ssz.Y / (2 * dsz.Y))
		drow := dst.Row(y)
		for x, sx := range xs {
			drow[x] = srow[sx]
		}
	}
}

// composeChanged16 upscales each pixel of src that differs from prev into dst.
func composeChanged16(dst, src *noximage.Image16, prev []uint16) {
	dsz, ssz := dst.Size(), src.Size()
	if ssz.X <= 0 || ssz.Y <= 0 || len(prev) < ssz.X*ssz.Y {
		return
	}
	// pixel i in src covers [xs[i], xs[i+1]) range in dst
	xs := make([]int, ssz.X+1)
	for x := range xs {
		xs[x] = x * dsz.X / ssz.X
	}
	for y := 0; y < ssz.Y; y++ {
		srow := src.Row(y)
		prow := prev[y*ssz.X : (y+1)*ssz.X]
		y0, y1 := y*dsz.Y/ssz.Y, (y+1)*dsz.Y/ssz.Y
		for x, c := range srow {
			if c == prow[x] {
				continue
			}
			x0, x1 := xs[x], xs[x+1]
			for dy := y0; dy < y1; dy++ {
				drow := dst.Row(dy)[x0:x1]
				for i := range drow {
					drow[i] = c
				}
			}
		}
	}
}
