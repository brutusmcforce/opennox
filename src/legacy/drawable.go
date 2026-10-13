package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
)

func asDrawable(p *nox_drawable) *client.Drawable {
	return (*client.Drawable)(unsafe.Pointer(p))
}

func AsDrawableP(p unsafe.Pointer) *client.Drawable {
	return (*client.Drawable)(p)
}

type nox_drawable = client.Drawable

func CallDrawFunc(s *client.Drawable, vp *noxrender.Viewport) int {
	return s.CallDraw(vp)
}

func Nox_xxx_spriteGetMB_476F80() *client.Drawable {
	return asDrawable((*nox_drawable)(interactionHoverDrawable))
}

func Nox_xxx_clientGetSpriteAtCursor_476F90() *client.Drawable {
	return asDrawable((*nox_drawable)(interactionUsableCursorDrawable))
}
func Get_dword_5d4594_1096640() *client.Drawable {
	return AsDrawableP(interactionHoverDrawable)
}
func Set_dword_5d4594_1096640(dr *client.Drawable) {
	interactionHoverDrawable = dr.C()
}
func Get_nox_client_spriteUnderCursorXxx_1096644() *client.Drawable {
	return AsDrawableP(interactionUsableCursorDrawable)
}
func Set_nox_client_spriteUnderCursorXxx_1096644(dr *client.Drawable) {
	interactionUsableCursorDrawable = dr.C()
}
func Sub_495B50(fx *client.DrawableFX) {
	combatFXDetach(fx)
}
func Sub_4523D0(p unsafe.Pointer) {
	audioEventDelete((*audioEvent)(p))
}
func Sub_495FC0(p *client.DrawableFX, dr *client.Drawable) {
	combatFXAttach(p, dr)
}
func Sub_49C520(dr *client.Drawable) int {
	return bool2int(presentationRayContains(dr))
}
func Sub_45A9B0(a1, a2 *client.Drawable) {
	clientDrawableLoopAudio(a1, a2)
}
func Nox_xxx_unitSpriteCheckAlly_4951F0(id int) bool {
	return combatAllyLookup(uint32(id)) != nil
}
func Nox_xxx_draw_44C650_free_kind(ptr unsafe.Pointer, kind int) {
	spriteDataFreeKind(ptr, kind)
}
