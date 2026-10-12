package legacy

import (
	"github.com/brutusmcforce/opennox/libs/player"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/memmap"
)

var glyphClientType, glyphItemType uint32

func glyphType(cache *uint32) uint32 {
	if *cache == 0 {
		*cache = uint32(GetClient().Cli().Things.IndByID("Glyph"))
	}
	return *cache
}

func glyphSelectionAllowed(dr *client.Drawable) int {
	typ := glyphType(&glyphClientType)
	pl := Get_dword_8531A0_2576()
	if pl == nil {
		return 0
	}
	return bool2int(dr.TypeIDVal != typ || pl.PlayerClass() == player.Wizard)
}

func glyphItemAllowed(dr *client.Drawable) int {
	typ := glyphType(&glyphItemType)
	if dr == nil || memmap.Uint32(0x852978, 8) == 0 {
		return 0
	}
	pl := Get_dword_8531A0_2576()
	if pl == nil || (dr.TypeIDVal == typ && pl.PlayerClass() != player.Wizard) {
		return 0
	}
	if nox_cheat_allowall != 0 {
		return 1
	}
	// Preserve the observed 386 shift semantics and evaluate before the callback.
	mask := byte(uint32(1) << (uint32(pl.PlayerClass()) & 31))
	allowed := Sub_57B370(dr.Class(), dr.SubClass(), int(int32(dr.TypeIDVal)))
	return bool2int(mask&allowed != 0)
}
