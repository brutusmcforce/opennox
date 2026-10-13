package layout_test

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	. "github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/server"
)

func testTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[client.Drawable](), reflect.TypeFor[server.Object](), reflect.TypeFor[server.Player](),
		reflect.TypeFor[server.PlayerUpdateData](), reflect.TypeFor[server.MonsterUpdateData](), reflect.TypeFor[gui.Window](),
		reflect.TypeFor[gui.WindowData](), reflect.TypeFor[server.Team](), reflect.TypeFor[server.NPC](),
	}
}

// The 386 model must match the real layout of the original 32-bit binary, and the offset table must be the identity
// there. The game no longer runs on 386; run with GOARCH=386 CGO_ENABLED=0 to check the model.
func TestModel386(t *testing.T) {
	if runtime.GOARCH != "386" {
		t.Skip("386 only")
	}
	for _, typ := range testTypes() {
		if got := Size386(typ); got != typ.Size() {
			t.Errorf("%s: size386 %d, real %d", typ, got, typ.Size())
			continue
		}
		tb := Build(typ)
		for i, v := range tb.Offsets() {
			if v != ^uint32(0) && uintptr(v) != uintptr(i) {
				t.Errorf("%s: offset %d maps to %d", typ, i, v)
				break
			}
		}
	}
}

// The amd64 model must match the real layout on amd64.
func TestModelNative(t *testing.T) {
	m := Model386
	if runtime.GOARCH == "amd64" {
		m = ModelAmd64
	} else if runtime.GOARCH != "386" {
		t.Skip("386 and amd64 only")
	}
	for _, typ := range testTypes() {
		if got := m.Size(typ); got != typ.Size() {
			t.Errorf("%s: model size %d, real %d", typ, got, typ.Size())
		}
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if got := m.OffsetOf(typ, FieldName(typ, i)); got != f.Offset {
				t.Errorf("%s.%s: model offset %d, real %d", typ, f.Name, got, f.Offset)
			}
		}
	}
}

// On amd64, Moves must agree with the offset tables.
func TestMovesMatchesOff(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("amd64 only")
	}
	for _, typ := range testTypes() {
		tb := Build(typ)
		for off := uintptr(0); off < Size386(typ); off++ {
			leaf, off64, _, _ := Moves(typ, off, 1)
			if leaf.Path == "<padding>" {
				continue
			}
			if want := tb.Off(off); off64 != want {
				t.Errorf("%s+%d (%s): Moves %d, Off %d", typ, off, leaf.Path, off64, want)
				break
			}
		}
	}
}

func TestRuns(t *testing.T) {
	type S struct {
		A uint32
		P *int
		B [3]uint16
	}
	tb := Build(reflect.TypeFor[S]())
	var got [][2]uintptr
	tb.Runs(0, Size386(reflect.TypeFor[S]()), func(_, off, n uintptr) { got = append(got, [2]uintptr{off, n}) })
	var want [][2]uintptr
	if runtime.GOARCH == "386" {
		want = [][2]uintptr{{0, 14}} // identity on 386, padding skipped
	} else {
		want = [][2]uintptr{{0, 4}, {8, 4}, {16, 6}} // A, first half of P, B
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
