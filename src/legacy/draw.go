package legacy

import (
	"image"
	"unsafe"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
)

var (
	Nox_xxx_cliUpdateCameraPos_435600    func(x, y int)
	Sub_437260                           func()
	Get_nox_client_texturedFloors_154956 func() bool
	Sub_473970                           func(a1 image.Point) image.Point
	Nox_client_isConnected               func() bool
	Nox_video_inFadeTransition_44E0D0    func() int
)

func nox_xxx_cliUpdateCameraPos_435600(x, y int) {
	Nox_xxx_cliUpdateCameraPos_435600(x, y)
}

func nox_draw_setMaterial_4340A0(ind, r, g, b int) {
	GetClient().R2().Data().SetMaterialRGB(ind, r, g, b)
}

func nox_xxx_drawSetTextColor_434390(a1 int) {
	GetClient().R2().Data().SetTextColor(noxcolor.RGBA5551(a1))
}

func nox_client_drawSetColor_434460(a1 int) {
	GetClient().R2().Data().SetColor2(noxcolor.RGBA5551(a1))
}

func nox_client_drawEnableAlpha_434560(a1 int) {
	GetClient().R2().Data().SetAlphaEnabled(a1 != 0)
}

func nox_client_drawAddPoint_49F500(x, y int) {
	GetClient().R2().AddPoint(image.Pt(x, y))
}

func nox_xxx_rasterPointRel_49F570(x, y int) {
	GetClient().R2().AddPointRel(image.Pt(x, y))
}

func nox_client_drawLineFromPoints_49E4B0() int {
	r := GetClient().R2()
	return bool2int(r.DrawLineFromPoints(r.Data().Color2()))
}

func nox_draw_set54RGB32_434040(cl int) {
	c := noxrender.SplitColor(noxcolor.RGBA5551(cl))
	GetClient().R2().Data().SetColorInt54(noxrender.RGB{
		R: int(c.R),
		G: int(c.G),
		B: int(c.B),
	})
}

func nox_client_drawRectFilledOpaque_49CE30(a1, a2, a3, a4 int) {
	r := GetClient().R2()
	r.DrawRectFilledOpaque(a1, a2, a3, a4, r.Data().Color2())
}

func nox_client_drawRectFilledAlpha_49CF10(a1, a2, a3, a4 int) {
	GetClient().R2().DrawRectFilledAlpha(a1, a2, a3, a4)
}

func nox_client_drawBorderLines_49CC70(a1, a2, a3, a4 int) {
	r := GetClient().R2()
	r.DrawBorder(a1, a2, a3, a4, r.Data().Color2())
}

func nox_client_drawPixel_49EFA0(a1, a2 int) {
	r := GetClient().R2()
	r.DrawPixel(image.Pt(a1, a2), r.Data().Color2())
}

func nox_xxx_drawPointMB_499B70(a1, a2, a3 int) {
	r := GetClient().R2()
	r.DrawPoint(image.Pt(a1, a2), a3, r.Data().Color2())
}

func nox_xxx_guiFontHeightMB_43F320(fnt unsafe.Pointer) int {
	r := GetClient().R2()
	return r.FontHeight(r.GetFonts().AsFont(fnt))
}

func nox_xxx_drawStringWrap_43FAF0(font unsafe.Pointer, sp *wchar2_t, x, y, maxW, maxH int) int {
	r := GetClient().R2()
	return r.DrawStringWrapped(r.GetFonts().AsFont(font), GoWString(sp), image.Rect(x, y, x+maxW, y+maxH))
}

func nox_video_drawCircleColored_4C3270(a1, a2, a3, a4 int) {
	GetClient().R2().DrawCircle(a1, a2, a3, noxcolor.RGBA5551(a4))
}

// GetDrawViewport returns the viewport for the current drawing pass. It may differ from the client viewport
// when the GUI is drawn at a different resolution than the world (UI scaling).
var GetDrawViewport func() *noxrender.Viewport

func drawViewport() *noxrender.Viewport {
	if GetDrawViewport != nil {
		return GetDrawViewport()
	}
	return GetClient().Viewport()
}

// WorldToDrawScreen converts world position to a screen position in the coordinates of the current drawing pass.
var WorldToDrawScreen func(p image.Point) image.Point

func worldToDrawScreen(p image.Point) image.Point {
	if WorldToDrawScreen != nil {
		return WorldToDrawScreen(p)
	}
	return GetClient().Viewport().ToScreenPos(p)
}

// Nox_client_particleFxDraw draws a particle effect associated with the drawable (see nox_thing_particle_draw).
var Nox_client_particleFxDraw func(dr *client.Drawable)

// Nox_client_gpu_edgeDraw draws a lit wall sprite on the GPU. It returns false if GPU rendering is disabled.
var Nox_client_gpu_edgeDraw func(img *noxrender.Image, pos image.Point, left, right [3]int, clip image.Rectangle, cut int) bool

func nox_client_drawImageAt_47D2C0(img noxrender.ImageHandle, x, y int) {
	GetClient().R2().DrawImageAt(asImage(img), image.Point{X: x, Y: y})
}

func sub_4AE6F0(cx, cy, rad, ang, ccl int) {
	GetClient().R2().DrawCircleSegment(cx, cy, rad, ang, noxcolor.RGBA5551(ccl))
}

func nox_client_isConnected_43C700() int {
	return bool2int(Nox_client_isConnected())
}

func Sub_437180() {
	chatBubbleDraw(drawViewport())
}

func Sub_476AE0(vp *noxrender.Viewport, dr *client.Drawable) {
	presentationBake(vp, dr)
}

func Nox_xxx_drawShinySpot_4C4F40(vp *noxrender.Viewport, dr *client.Drawable) {
	objectRenderShiny(vp, dr)
}

func Sub_499F60(p int, pos image.Point, a4 int, a5, a6, a7, a8, a9 int, a10 int) {
	presentationBubble(p, pos, int16(a4), byte(a5), byte(a6), byte(a7), byte(a8), byte(a9), a10)
}

func Sub_435120(a1 unsafe.Pointer, a2 unsafe.Pointer) {
	clientPaletteExpand(a1, a2)
}
func Sub_435040() {
	clientPaletteSort()
}
func Sub_435150(a1 unsafe.Pointer, a2 unsafe.Pointer) {
	clientPaletteCompact(a1, a2)
}
func Nox_xxx_wndDraw_49F7F0() {
	objectRenderSaveClip()
}
func Sub_49F780(a1 int, a2 int) {
	uiRenderNarrowClip(a1, a2)
}
func Sub_49F860() {
	objectRenderRestoreClip()
}
func Nox_xxx_drawEnergyBolt_499710(a1 int, a2 int, a3 int, a4 int) {
	effectCreateEnergySparks(a1, a2, int16(a3), a4)
}
func Nox_xxx_drawShield_499810(vp *noxrender.Viewport, dr *client.Drawable) {
	presentationShieldDraw(vp, dr)
}
func Sub_474B40(dr *client.Drawable) int {
	return bool2int(worldWallPlayerVisible(dr))
}
func Sub_495BB0(dr *client.Drawable, vp *noxrender.Viewport) {
	combatFXDraw(vp, dr)
}
