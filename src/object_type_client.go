package opennox

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/strman"

	"github.com/brutusmcforce/opennox/client"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func init() {
	// TODO: pass strman
	client.RegisterThingParse("PRETTYNAME", func(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error {
		typ.PrettyName = alloc.InternCString16(noxClient.Strings().GetStringInFile(strman.ID(str), "drawdb.c"))
		return nil
	})
	client.RegisterThingParse("DESCRIPTION", func(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error {
		typ.Desc = alloc.InternCString16(noxClient.Strings().GetStringInFile(strman.ID(str), "drawdb.c"))
		return nil
	})
}

func (c *Client) Nox_things_free_44C580() {
	c.Things.Each(func(it *client.ObjectType) {
		if it.DrawData != nil {
			freeDrawableData(it.DrawData, it.DrawFunc)
		}
	})
	c.Things.Nox_things_free_44C580_B()
	if !noxflags.HasGame(noxflags.GameHost) {
		nox_xxx_free_42BF80()
	}
}

func freeDrawableData(ptr unsafe.Pointer, fnc unsafe.Pointer) {
	legacy.Nox_xxx_draw_44C650_free_kind(ptr, client.DrawableDataKind(fnc))
}
