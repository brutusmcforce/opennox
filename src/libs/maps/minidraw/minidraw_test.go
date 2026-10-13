package minidraw

import (
	"image/color"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/libs/maps"
	"github.com/brutusmcforce/opennox/libs/maps/maprender"
	"github.com/brutusmcforce/opennox/libs/noxtest"
)

var casesMapDraw = []struct {
	Name string
	Hash string
}{
	{Name: "con01a", Hash: "311087a7b3a58dd60048586643f6a869"},
	{Name: "estate", Hash: "747bb24e891dc4ba76d1baa0fe75ded6"},
	{Name: "g_castle", Hash: "9c76e17c0ff224ea92b1ada92262210b"},
	{Name: "g_mines", Hash: "2be98df6635a0932badc1a377754aaff"},
	{Name: "so_brin", Hash: "14df1f4edf3305ab97a99d165739cf96"},
}

func TestDraw(t *testing.T) {
	r, err := maprender.NewRenderer(noxtest.DataPath(t))
	must.NoError(t, err)
	defer r.Close()
	path := noxtest.DataPath(t, maps.Dir)
	for _, m := range casesMapDraw {
		t.Run(m.Name, func(t *testing.T) {
			mp, err := maps.ReadMap(filepath.Join(path, m.Name))
			must.NoError(t, err)
			img := MinimapRGBA(mp, color.Black)
			noxtest.WritePNG(t, m.Name+".png", img, m.Hash)
		})
	}
}
