package opennox

import (
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy"
)

func init() {
	client.RegisterDraw("MaidenDraw", legacy.Get_nox_thing_maiden_draw(), 7, parseMaidenDraw)
}

func parseMaidenDraw(obj *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error {
	err := parseMonsterDraw(obj, f, str, buf)
	obj.DrawFunc = legacy.Get_nox_thing_maiden_draw()
	return err
}
