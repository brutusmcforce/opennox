package legacy

import (
	"testing"
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

func TestMapNameField(t *testing.T) {
	var st server.Settings2
	copy(st.Field0[9:], "server")
	for _, name := range []string{"estate", "So_Beach", "MyLongArenaMap", "AnotherVeryLongMapName12345", "MyLongArenaMap2"} {
		mapNameSetField(&st.Field0[0], name)
		if got := mapNameField(&st.Field0[0]); got != name {
			t.Errorf("got %q, want %q", got, name)
		}
		if got := string(st.Field0[9:15]); got != "server" {
			t.Fatalf("%q overwrote the server name: %q", name, got)
		}
		// a copy of the record keeps the name
		cp := st
		if got := ruleMapName(&cp); got != name {
			t.Errorf("copy: got %q, want %q", got, name)
		}
	}
	// the same long name gets the same token
	mapNameSetField(&st.Field0[0], "MyLongArenaMap")
	a := string(st.Field0[:8])
	mapNameSetField(&st.Field0[0], "mylongarenamap")
	if b := string(st.Field0[:8]); a != b {
		t.Errorf("different tokens: %q, %q", a, b)
	}
	// unknown tokens stay as they are
	if got := mapNameResolve("~fffffff"); got != "~fffffff" {
		t.Errorf("got %q", got)
	}
	var it Nox_map_list_item
	it.SetMapName("MyLongArenaMap")
	if got := it.MapName(); got != "MyLongArenaMap" {
		t.Errorf("map list: got %q", got)
	}
	var row mapQuestRow
	row.setName("MyVeryLongQuestMap.map")
	if got := row.name(); got != "MyVeryLongQuestMap.map" {
		t.Errorf("quest row: got %q", got)
	}
	if got := mapNameField((*byte)(unsafe.Pointer(&row.Name[0]))); got == "" {
		t.Errorf("quest row empty")
	}
	row.setName("g_castle.map")
	if got := row.name(); got != "g_castle.map" {
		t.Errorf("quest row: got %q", got)
	}
}
