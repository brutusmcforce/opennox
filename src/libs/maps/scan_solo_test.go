package maps

import (
	"log/slog"
	"slices"
	"testing"

	"github.com/brutusmcforce/opennox/libs/noxtest"
)

func TestIsSoloName(t *testing.T) {
	for name, want := range map[string]bool{
		"con01a": true, "War11A": true, "wiz07f": true,
		"Conquest": false, "warzone": false, "wizards_arena": false, "con01": false, "con1ab": false, "estate": false,
	} {
		if got := IsSoloName(name); got != want {
			t.Errorf("IsSoloName(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestScanSolo(t *testing.T) {
	path := noxtest.DataPath(t, Dir)
	has := func(list MapList, name string) bool {
		return slices.ContainsFunc(list, func(m Info) bool { return m.Filename == name })
	}
	for _, solo := range []bool{false, true} {
		list, err := Scan(slog.Default(), path, &ScanOptions{Solo: solo})
		if err != nil && len(list) == 0 {
			t.Fatal(err)
		}
		if !has(list, "Estate") && !has(list, "estate") {
			t.Errorf("solo=%v: no Estate", solo)
		}
		for _, name := range []string{"con01a", "wiz01a", "war01a"} {
			if has(list, name) != solo {
				t.Errorf("solo=%v: %s listed: %v", solo, name, has(list, name))
			}
		}
	}
}
