package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
)

func init() {
	client.RegisterThingParse("LIGHTDIRECTION", resourceClientField("direction"))
	client.RegisterThingParse("LIGHTPENUMBRA", resourceClientField("penumbra"))
	client.RegisterThingParse("CLIENTUPDATE", resourceClientField("update"))
	client.RegisterThingParse("PRETTYIMAGE", resourceClientField("image"))
	client.ThingDrawDefault = drawableDrawKey(drawKey_nox_thing_debug_draw)
}

func nox_get_thing_name(i int) *int8 {
	t := GetClient().Cli().Things.TypeByInd(i)
	if t == nil {
		return nil
	}
	return (*int8)(unsafe.Pointer(t.Name))
}

func nox_get_thing_pretty_name(i int) *wchar2_t {
	t := GetClient().Cli().Things.TypeByInd(i)
	if t == nil {
		return nil
	}
	return (*wchar2_t)(unsafe.Pointer(t.PrettyName))
}
