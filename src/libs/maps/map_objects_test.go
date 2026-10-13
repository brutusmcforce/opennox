package maps_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/maps"
	"github.com/brutusmcforce/opennox/libs/noxtest"
	"github.com/brutusmcforce/opennox/libs/xfer"
)

// readObjects decodes all objects of a map file.
func readObjects(t testing.TB, path string) []maps.Xfer {
	f, err := ifs.Open(path) // names differ in case
	must.NoError(t, err)
	defer func() { _ = f.Close() }()
	r, err := maps.NewReader(f)
	must.NoError(t, err)
	raw, err := r.ReadSectionsRaw()
	must.NoError(t, err)
	var (
		toc  maps.ObjectsTOC
		objs maps.Objects
	)
	for _, s := range raw {
		switch s.Name {
		case toc.MapSection():
			must.NoError(t, toc.UnmarshalBinary(s.Data))
		case objs.MapSection():
			must.NoError(t, objs.UnmarshalBinary(s.Data))
		}
	}
	list, err := objs.ReadObjects(&toc, xfer.DefaultRegistry)
	must.NoError(t, err)
	return list
}

// Every object of every map in the game data decodes, using all of its data.
func TestReadObjectsAllMaps(t *testing.T) {
	path := noxtest.DataPath(t, maps.Dir)
	dirs, err := os.ReadDir(path)
	must.NoError(t, err)
	n := 0
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		files, _ := os.ReadDir(filepath.Join(path, d.Name()))
		for _, f := range files {
			// names differ in case, e.g. con01a/Con01A.map
			if !strings.EqualFold(f.Name(), d.Name()+maps.Ext) {
				continue
			}
			n++
			t.Run(d.Name(), func(t *testing.T) {
				list := readObjects(t, filepath.Join(path, d.Name(), f.Name()))
				must.NotEq(t, 0, len(list))
			})
		}
	}
	t.Logf("%d maps", n)
}

func TestReadObjectsValues(t *testing.T) {
	path := noxtest.DataPath(t, maps.Dir)
	list := readObjects(t, filepath.Join(path, "con01a", "con01a"+maps.Ext))
	byType := make(map[string][]xfer.Xfer)
	for _, o := range list {
		byType[o.Type] = append(byType[o.Type], o.Xfer)
	}
	exit := byType["InvisibleExitArea"][0].(*xfer.Exit)
	must.Eq(t, "Con02a.map", exit.MapName)
	book := byType["SpellBook"][0].(*xfer.SpellReward)
	must.Eq(t, "SpellBook", book.Object.Type)
	must.StrHasPrefix(t, "SPELL_", book.Spell)
	maiden := byType["Maiden"][0].(*xfer.Monster)
	must.Eq(t, "CON01A:Tanya2", maiden.Object.Name)
	must.Eq(t, 0.5, maiden.RetreatLevel)
	// elevators point at each other
	for _, x := range byType["Elevator"] {
		e := x.(*xfer.Elevator)
		found := false
		for _, o := range list {
			if s, ok := o.Xfer.(*xfer.ElevatorShaft); ok && s.Object.Extent == e.PartnerExtent {
				must.Eq(t, e.Object.Extent, s.PartnerExtent)
				found = true
			}
		}
		must.True(t, found)
	}
	// NPCs carry their inventory
	npc := byType["NPC"][0].(*xfer.NPC)
	must.Eq(t, int(npc.Object.SubN), len(npc.Sub))
	for _, it := range npc.Sub {
		must.NotEq(t, "", it.(interface{ XferType() xfer.Type }).XferType())
	}
}
