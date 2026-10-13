// Command fieldat prints the fields at 386 offsets of a struct, and where they are on amd64.
//
// Usage (from src/):
//
//	go run ./internal/fieldat server.MonsterUpdateData 1336 1344 2180
package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/server"
)

var typesByName = map[string]reflect.Type{}

func add[T any]() {
	t := reflect.TypeFor[T]()
	typesByName[t.String()] = t
}

func init() {
	add[server.Object]()
	add[server.Player]()
	add[server.PlayerUpdateData]()
	add[server.MonsterUpdateData]()
	add[server.MonsterDef]()
	add[server.SpellProjectileUpdateData]()
	add[server.DurSpell]()
	add[server.Team]()
	add[server.ModifierEff]()
	add[server.Modifier]()
	add[server.HealthData]()
	add[client.Drawable]()
	add[gui.Window]()
	add[gui.WindowData]()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: fieldat <type> <offset>...")
		os.Exit(2)
	}
	t, ok := typesByName[os.Args[1]]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown type %q\n", os.Args[1])
		os.Exit(2)
	}
	for _, a := range os.Args[2:] {
		off, err := strconv.ParseUint(a, 0, 32)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		leaf, off64, moved, resized := layout.Moves(t, uintptr(off), 1)
		flags := ""
		if moved {
			flags += " moved"
		}
		if resized {
			flags += " resized"
		}
		in := ""
		if uintptr(off) != leaf.Off {
			in = fmt.Sprintf(" (+%d)", uintptr(off)-leaf.Off)
		}
		fmt.Printf("%6d %s%s %v -> %d%s\n", off, leaf.Path, in, leaf.Type, off64, flags)
	}
}
