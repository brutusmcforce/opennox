//go:build !server

package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

func Nox_xxx_clientDrawAll_436100_draw_A() {
	clientLoadingOverlay()
}

func Nox_xxx_drawMinimapAndLines_4738E0() {
	minimapDrawAndMessages()
}

func Nox_xxx_clientDrawAll_436100_draw_B() {
	clientWinnerOverlay()
}

func Sub_436F50() {
	clientDebugOverlay()
}

func Sub_437100() {
	sub_437100()
}

func Sub_470DE0() {
	uiMeterHeartbeat()
}

func Nox_xxx_clientEnumHover_476FA0() {
	nox_xxx_clientEnumHover_476FA0()
}

func Nox_xxx_spriteDeleteSomeList_49C4B0() {
	presentationTransientClear()
}

func Sub_49BBC0() {
	presentationChantTick()
}

func Nox_xxx_polygonDrawColor_421B80() {
	mapPolygonColor()
}

func Nox_xxx_cliToggleObsWindow_4357A0() {
	nox_xxx_cliToggleObsWindow_4357A0()
}

func Nox_xxx_motd_4467F0() {
	sessionMOTDShow()
}

func Sub_42EBA0() int {
	return int(memmap.Int32(0x5D4594, 754052))
}

func Sub_49B6E0() {
	sub_49B6E0()
}

func Get_nox_thing_glow_orb_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_glow_orb_draw)
}

func Nox_xxx_drawAllMB_475810_draw_B(vp *noxrender.Viewport) int {
	if noxflags.HasEngine(noxflags.EngineNoFloorRendering) {
		return 0
	}
	tile := tileAtPoint(types.Pointf{X: float32(vp.World.Max.X), Y: float32(vp.World.Max.Y)})
	if tile == -1 || tile == 255 {
		return 0
	}
	return 1
}
func Sub_4C5060(vp *noxrender.Viewport) {
	objectRenderBeamDraw(vp)
}
func Nox_xxx_drawWalls_473C10(vp *noxrender.Viewport, a2 *server.Wall) {
	worldWallDraw(vp, a2)
}
func Sub_4761B0(dr *client.Drawable) int {
	return int(presentationDrawableY(dr))
}
func Sub_476080(a1 unsafe.Pointer) int {
	return int(presentationWallY((*[8]byte)(a1)))
}
func Sub_459DB0(dr *client.Drawable) int {
	return drawableStatePredicate(dr)
}
func Sub_49A6A0(vp *noxrender.Viewport, dr *client.Drawable) {
	combatHealthDraw(vp, dr)
}
func Nox_xxx_sprite_4756E0_drawable(dr *client.Drawable) int {
	return bool2int(worldWallStaticPass(dr))
}
func Nox_xxx_sprite_475740_drawable(dr *client.Drawable) int {
	return bool2int(worldWallDynamicPass(dr))
}
func Nox_xxx_sprite_4757A0_drawable(dr *client.Drawable) int {
	return bool2int(worldWallSpecialPass(dr))
}
func Sub_4757D0_drawable(dr *client.Drawable) int {
	return bool2int(worldWallInactivePass(dr))
}
func Nox_xxx_tileDrawImpl_4826A0(vp *noxrender.Viewport) {
	tileCompositionFull(vp)
}
func Nox_xxx_tileDrawMB_481C20_A(vp *noxrender.Viewport, a2 int) {
	tileCompositionHorizontal(vp, a2)
}
func Nox_xxx_tileDrawMB_481C20_B(vp *noxrender.Viewport, a2 int) {
	tileCompositionVertical(vp, a2)
}
func Nox_xxx_tileCheckRedrawMB_482570(vp *noxrender.Viewport) int {
	return tileCompositionRedraw(vp)
}
