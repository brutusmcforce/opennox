package opennox

import (
	"image"

	"github.com/brutusmcforce/opennox/libs/noximage"

	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
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
// After the world is drawn, the GUI is drawn into a separate transparent render target at the full resolution
// (coordinates are still in GUI pixels, see gpu.Renderer.SetTargetSize), which is then composed over the world
// (see composeUI). The world buffer is presented on the screen.
var uiScaleState struct {
	worldSize image.Point // requested size for the world buffer; zero if UI scaling is disabled

	world *noximage.Image16

	inWorld bool
	saved   struct {
		img   *noximage.Image16
		win   image.Point
		clip  bool
		rect1 image.Rectangle
		rect2 image.Rectangle
		rect3 image.Rectangle
	}

	worldReady bool // the world was drawn in this frame

	viewSize image.Point // the world view size that the client state was set up for (see resizeWorldView)

	overlay *noximage.Image16 // full resolution layer for the cursor, when the world buffer is larger than the screen

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

func videoAllocWorldBuffer(sz image.Point) {
	uiScaleState.world = noximage.NewImage16(image.Rectangle{Max: sz})
	if g := noxClient.gpu; g != nil && noxPixBuffer.img != nil {
		// draw the GUI at the full resolution; it allows drawing smooth text
		g.SetTargetSize(noxPixBuffer.img, worldPhysSize())
	}
	uiScaleState.worldReady = false
}

func videoFreeWorldBuffer() {
	if uiScaleState.inWorld {
		noxClient.endWorldDraw()
	}
	if g := noxClient.gpu; g != nil && uiScaleState.world != nil && noxPixBuffer.img != nil {
		g.SetTargetSize(noxPixBuffer.img, image.Point{})
	}
	uiScaleState.world = nil
	uiScaleState.worldReady = false
}

func (c *Client) setPixTarget(img *noximage.Image16) {
	noxPixBuffer.img = img
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
	sv.img = noxPixBuffer.img
	sv.win = videoGetWindowSize()
	sv.clip = d.Clip()
	sv.rect1, sv.rect2, sv.rect3 = d.ClipRect(), d.ClipRect2(), d.Rect3()

	sz := st.world.Size()
	c.setPixTarget(st.world)
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
	c.setPixTarget(sv.img)
	videoSetWindowSize(sv.win)
	d := c.r.Data()
	d.SetClip(sv.clip)
	d.SetClipRect(sv.rect1)
	d.SetClipRect2(sv.rect2)
	d.SetRect3(sv.rect3)
	sv.img = nil
}

// worldDrawn must be called after the world is drawn into the world buffer.
// It applies the zoom and prepares the GUI target.
func (c *Client) worldDrawn() {
	st := &uiScaleState
	if st.world == nil || st.inWorld || c.gpu == nil {
		return
	}
	st.worldReady = true
	c.applyZoom(st.world)
	// GUI is drawn into a separate transparent target, which is composed over the world later
	c.gpu.ClearTransparent()
}

// composeUI draws the GUI over the world buffer.
func (c *Client) composeUI() {
	st := &uiScaleState
	if st.world == nil || st.inWorld || !st.worldReady || c.gpu == nil {
		return
	}
	c.gpu.Composite(st.world, noxPixBuffer.img)
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
	if st.world == nil || st.inWorld || !st.worldReady {
		return false
	}
	c.composeUI()
	c.drawCursorOverlay()
	c.beginWorldDraw()
	defer c.endWorldDraw()
	draw()
	c.copyPixBuffer()
	return true
}

// drawCursorOverlay draws the deferred cursor at the screen resolution, if the world buffer is larger than the screen
// (when zoomed out). Otherwise, the cursor is drawn directly into the world buffer (see drawDeferredCursor).
func (c *Client) drawCursorOverlay() {
	st := &uiScaleState
	phys := worldPhysSize()
	if !st.hasCursor || st.world == nil || st.world.Size() == phys || c.gpu == nil || phys.X <= 0 || phys.Y <= 0 {
		return
	}
	st.hasCursor = false
	if st.overlay == nil || st.overlay.Size() != phys {
		st.overlay = noximage.NewImage16(image.Rectangle{Max: phys})
	}
	ui := noxPixBuffer.img.Size()
	pos := image.Pt(st.cursorPos.X*phys.X/ui.X, st.cursorPos.Y*phys.Y/ui.Y)

	d := c.r.Data()
	saved := noxPixBuffer.img
	clip, r1, r2, r3 := d.Clip(), d.ClipRect(), d.ClipRect2(), d.Rect3()
	win := videoGetWindowSize()
	c.setPixTarget(st.overlay)
	videoSetWindowSize(phys)
	d.SetClip(false)
	d.SetClipRect(image.Rectangle{Max: phys})
	d.SetClipRect2(image.Rectangle{Max: phys.Sub(image.Pt(1, 1))})
	d.SetRect3(image.Rectangle{Max: phys})
	c.gpu.ClearTransparent()
	c.nox_video_cursorDrawImpl_477A30(pos)
	c.setPixTarget(saved)
	videoSetWindowSize(win)
	d.SetClip(clip)
	d.SetClipRect(r1)
	d.SetClipRect2(r2)
	d.SetRect3(r3)
	c.gpu.Composite(st.world, st.overlay)
}
