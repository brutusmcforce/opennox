package opennox

import (
	"github.com/spf13/viper"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy"
)

func sub_473670() int {
	return int(memmap.Uint8(0x5D4594, 1096424) & 1)
}

func nox_client_toggleMap_473610() int {
	if nox_xxx_guiCursor_477600() != 1 {
		if memmap.Uint8(0x5D4594, 1096424)&1 != 0 {
			legacy.MinimapLarge = false
			clientPlaySoundSpecial(sound.SoundMapClose, 100)
			*memmap.PtrUint32(0x5D4594, 1096424) &= 0xFFFFFFFE
		} else {
			clientPlaySoundSpecial(sound.SoundMapOpen, 100)
			*memmap.PtrUint32(0x5D4594, 1096424) |= 0x1
		}
		noxClient.GUI.ValYYY = 1
		return int(*memmap.PtrUint32(0x5D4594, 1096424))
	}
	return 0
}

// cycleMap is the map key: it cycles through a closed map, the normal map on the left and a large map
// in the middle of the screen (opennox#716).
func cycleMap() {
	if sub_473670() != 0 && !legacy.MinimapLarge && nox_xxx_guiCursor_477600() != 1 {
		legacy.MinimapLarge = true
		clientPlaySoundSpecial(sound.SoundMapOpen, 100)
		return
	}
	nox_client_toggleMap_473610() // closing also resets the large map
}

// configLargeMapOpacity is the opacity of the large map in percent.
const configLargeMapOpacity = "video.large_map_opacity"

var largeMapOpacityChoices = []int{20, 40, 60, 80, 100}

func init() {
	viper.SetDefault(configLargeMapOpacity, legacy.MinimapLargeOpacity)
	registerOnConfigRead(func() {
		setLargeMapOpacity(viper.GetInt(configLargeMapOpacity))
	})
}

func setLargeMapOpacity(v int) {
	legacy.MinimapLargeOpacity = min(max(v, 10), 100)
}
