package maprender

import (
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/libs/maps"
	"github.com/brutusmcforce/opennox/libs/noxtest"
)

var casesMapDraw = []struct {
	Name string
	Hash string
}{
	{Name: "con01a", Hash: "53a042d4d82e00a92f94939606b7fedb"},
	{Name: "estate", Hash: "76d6b9937b5f16d08383df402ce4614b"},
	{Name: "g_castle", Hash: "11090440f62f350c791f06bcaa031042"},
	{Name: "g_mines", Hash: "abdfd297c8caca310eb3e6a573f2e038"},
	{Name: "so_brin", Hash: "43ceb85e9ca3575f5a8bb315b4e03c5c"},
}

func TestDraw(t *testing.T) {
	r, err := NewRenderer(noxtest.DataPath(t))
	must.NoError(t, err)
	defer r.Close()
	path := noxtest.DataPath(t, maps.Dir)
	for _, m := range casesMapDraw {
		t.Run(m.Name, func(t *testing.T) {
			mp, err := maps.ReadMap(filepath.Join(path, m.Name))
			must.NoError(t, err)
			img, err := r.DrawMap(mp, nil)
			noxtest.WritePNG(t, m.Name+".png", img, m.Hash)
			must.NoError(t, err)
		})
	}
}

func BenchmarkDraw(b *testing.B) {
	b.ReportAllocs()
	r, err := NewRenderer(noxtest.DataPath(b))
	must.NoError(b, err)
	defer r.Close()
	path := noxtest.DataPath(b, maps.Dir)
	b.StopTimer()
	for _, m := range casesMapDraw {
		b.Run(m.Name, func(b *testing.B) {
			mp, err := maps.ReadMap(filepath.Join(path, m.Name))
			must.NoError(b, err)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = r.DrawMap(mp, nil)
			}
		})
	}
}
