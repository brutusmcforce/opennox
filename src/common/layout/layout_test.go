package layout

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/server"
)

func testTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[client.Drawable](), reflect.TypeFor[server.Object](), reflect.TypeFor[server.Player](),
		reflect.TypeFor[server.PlayerUpdateData](), reflect.TypeFor[server.MonsterUpdateData](), reflect.TypeFor[gui.Window](),
		reflect.TypeFor[gui.WindowData](), reflect.TypeFor[server.Team](), reflect.TypeFor[server.NPC](),
	}
}

// On 386 the 386 model must match the real layout, and the offset table must be the identity.
func TestModel386(t *testing.T) {
	if runtime.GOARCH != "386" {
		t.Skip("386 only")
	}
	forceTables = true
	defer func() { forceTables = false }()
	for _, typ := range testTypes() {
		if got := size386(typ); got != typ.Size() {
			t.Errorf("%s: size386 %d, real %d", typ, got, typ.Size())
			continue
		}
		tb := build(typ)
		for i, v := range tb.off {
			if uintptr(v) != uintptr(i) {
				t.Errorf("%s: offset %d maps to %d", typ, i, v)
				break
			}
		}
	}
}
