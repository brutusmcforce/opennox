package opennox

import (
	"unsafe"

	"github.com/spf13/viper"

	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy"
)

// configVideoSmoothFonts enables bundled outline fonts instead of the original bitmap fonts (see noxrender.SmoothFonts).
const configVideoSmoothFonts = "video.smooth_fonts"

func init() {
	viper.SetDefault(configVideoSmoothFonts, true)
}

func (c *Client) loadGameFonts() error {
	noxrender.SmoothFonts = viper.GetBool(configVideoSmoothFonts)
	if err := c.r.Fonts.Load(c.Strings().Lang()); err != nil {
		return err
	}
	c.nox_xxx_fontLoadMB_43F1C0()
	return nil
}

func (c *Client) nox_xxx_fontLoadMB_43F1C0() {
	legacy.Set_dword_5d4594_816444(0)
	v1 := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, 816464)), 5)
	for i := range v1 {
		v1[i] = 0
	}
}
