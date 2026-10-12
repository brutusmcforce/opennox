package legacy

import (
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"unsafe"
)

var drawableUpdateIdentitySlots [28]byte

const (
	updateID_colorlight    = 0
	updateID_cloud         = 1
	updateID_sub_4CE360    = 2
	updateID_dball         = 3
	updateID_sub_4CE0A0    = 4
	updateID_dball_charge  = 5
	updateID_magic         = 6
	updateID_sparkle       = 7
	updateID_magic_missile = 8
	updateID_teleport      = 9
	updateID_sub_4CD690    = 10
	updateID_sub_4CD450    = 11
	updateID_sub_4CD400    = 12
	updateID_sub_4CCE70    = 13
	updateID_sub_4CD090    = 14
	updateID_sub_4CD0C0    = 15
	updateID_sub_4CD0F0    = 16
	updateID_sub_4CD120    = 17
	updateID_fist          = 18
	updateID_sub_4CCD00    = 19
	updateID_undead        = 20
	updateID_manabomb      = 21
	updateID_vortex        = 22
	updateID_sub_4CA650    = 23
	updateID_monstergen    = 24
	updateID_sub_4CE340    = 25
	updateID_sub_4CA720    = 26
	updateID_sprite        = 27
)

func drawableUpdateIdentity(id int) unsafe.Pointer {
	return unsafe.Pointer(&drawableUpdateIdentitySlots[id])
}

func init() {
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_colorlight), func(vp *noxrender.Viewport, dr *client.Drawable) int32 { return int32(colorLightUpdate(vp, dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_cloud), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateCloudFrame(dr, 75)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CE360), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateCloudFrame(dr, 35)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_dball), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { updateDeathBallSparks(dr, 3); return 1 })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CE0A0), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { updateDeathBallSparks(dr, 1); return 1 })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_dball_charge), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateDeathBallCharge(dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_magic), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { updateMagicTrail(dr); return 0 })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sparkle), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateTrailSparks(dr, false)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_magic_missile), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateMagicMissile(dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_teleport), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateTeleportWake(dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CD690), func(vp *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateHealDrain(vp, dr, false)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CD450), func(vp *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateHealDrain(vp, dr, true)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CD400), func(vp *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateCharm(vp, dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CCE70), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateFireballFrame(dr, 5)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CD090), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateFireballFrame(dr, 4)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CD0C0), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateFireballFrame(dr, 3)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CD0F0), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateFireballFrame(dr, 2)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CD120), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateFireballFrame(dr, 1)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_fist), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateHeight(dr, true)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CCD00), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateHeight(dr, false)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_undead), func(_ *noxrender.Viewport, _ *client.Drawable) int32 { return 1 })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_manabomb), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateManaBomb(dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_vortex), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateVortex(dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CA650), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(drawableMotionTarget(dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_monstergen), func(_ *noxrender.Viewport, _ *client.Drawable) int32 { return 1 })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CE340), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(updateCloudRise(dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sub_4CA720), func(_ *noxrender.Viewport, dr *client.Drawable) int32 { return int32(effectOrbitUpdate(dr)) })
	client.RegisterDrawableUpdateCallbackGo(drawableUpdateIdentity(updateID_sprite), func(vp *noxrender.Viewport, dr *client.Drawable) int32 { return int32(drawableMotionDamped(vp, dr)) })
}
