package legacy

import "github.com/brutusmcforce/opennox/client"

func worldWallPlayerVisible(dr *client.Drawable) bool {
	localTeam := objectRenderTeam(ClientPlayerNetCode())
	if localTeam != nil {
		team := objectRenderTeam(int(dr.NetCode32))
		if team != nil && (ClientPlayerNetCode() == int(dr.NetCode32) || localTeam.SameAs(team)) {
			return true
		}
	}
	local := minimapLocalDrawable()
	if dr == local {
		return true
	}
	if local != nil && local.Buffs&(1<<21) != 0 {
		return true
	}
	// Preserve the existing player lookup even though its result is unused.
	if dr.ObjClass&4 != 0 && dr != local {
		GetServer().S().Players.ByID(int(dr.NetCode32))
	}
	return false
}
func worldWallStaticPass(dr *client.Drawable) bool {
	flags, class, fn := dr.ObjFlags, dr.ObjClass, dr.DrawFuncPtr
	return fn != nil && flags&0x1000 == 0 && flags&1 != 0 && (fn == drawableDrawKey(drawKey_nox_thing_static_draw) || fn == drawableDrawKey(drawKey_nox_thing_static_random_draw)) && class&0x80800000 == 0 && (flags&0x48 != 0 || class&0x400000 != 0) && flags&0x800 == 0
}
func worldWallDynamicPass(dr *client.Drawable) bool {
	return dr.DrawFuncPtr != nil && dr.ObjFlags&0x1000 == 0 && dr.ObjFlags&1 != 0 && !worldWallStaticPass(dr)
}
func worldWallSpecialPass(dr *client.Drawable) bool {
	return dr.DrawFuncPtr != nil && dr.ObjFlags&0x1000 == 0 && dr.ObjFlags&0x4000 != 0 && dr.ObjFlags&0x40 != 0
}
func worldWallInactivePass(dr *client.Drawable) bool {
	return dr.DrawFuncPtr != nil && dr.ObjFlags&1 == 0 && (dr.ObjClass&0x2000 == 0 || dr.ObjFlags&0x1000000 != 0) && dr.ObjFlags&0x1000 == 0
}
