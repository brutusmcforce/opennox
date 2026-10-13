package things

import (
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/libs/noxtest"
)

func TestReadFloorsEdges(t *testing.T) {
	f, err := Open(filepath.Join(noxtest.DataPath(t), "thing.bin"))
	must.NoError(t, err)
	defer func() { _ = f.Close() }()

	floors, err := f.ReadFloors()
	must.NoError(t, err)
	must.Len(t, 145, floors)
	must.Eq(t, "Mud", floors[0].Name)
	for _, fl := range floors {
		must.Len(t, int(fl.Rows)*int(fl.Cols)*int(fl.Frames), fl.Images)
	}

	edges, err := f.ReadEdges()
	must.NoError(t, err)
	must.Len(t, 39, edges)
	must.Eq(t, "MudEdge", edges[0].Name)
	for _, e := range edges {
		must.Len(t, 2*int(e.Variants)*(int(e.Rows)+int(e.Cols)), e.Images)
	}

	data, err := f.ReadAll()
	must.NoError(t, err)
	must.Eq(t, floors, data.Floors)
	must.Eq(t, edges, data.Edges)
}

func TestReadDraws(t *testing.T) {
	f, err := Open(filepath.Join(noxtest.DataPath(t), "thing.bin"))
	must.NoError(t, err)
	defer func() { _ = f.Close() }()
	list, err := f.ReadThings()
	must.NoError(t, err)
	byName := make(map[string]*Thing)
	for i := range list {
		byName[list[i].Name] = &list[i]
	}
	lever, ok := byName["FloorLever2"].Draw.(SlaveDraw)
	must.True(t, ok)
	must.NotEq(t, 0, len(lever.Imgs))
	_, ok = byName["BoulderIndestructible"].Draw.(BoulderDraw)
	must.True(t, ok)
	pit, ok := byName["Pit"].Draw.(AnimateStateDraw)
	must.True(t, ok)
	must.NotEq(t, 0, len(pit.States))
	for _, st := range pit.States {
		must.NotEq(t, 0, len(st.Anim.Frames))
	}
}
