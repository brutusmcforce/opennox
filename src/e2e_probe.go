package opennox

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/libs/wall"
	ns4 "github.com/brutusmcforce/opennox/noxscript/ns/v4"
	"github.com/brutusmcforce/opennox/server"
)

// Probes for investigating bugs in e2e scripts: `action: probe` with `args: [name, ...]`. They act on the server
// state directly (load a map, teleport, list or damage objects), so a script can get to a situation quickly.

var e2eProbes = map[string]func(args []string) error{
	// load <map>: switch to a map, ignoring the game mode checks (like "set maps allow.all" + "load")
	"load": func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: load <map>")
		}
		return noxLoadMap(args[0], mapLoadOptions{Force: true})
	},
	// tp <x> <y>: move the host's player
	"tp": func(args []string) error {
		pos, err := e2eProbePos(args)
		if err != nil {
			return err
		}
		u := e2eQuestPlayer()
		if u == nil {
			return fmt.Errorf("no player")
		}
		u.SetPos(pos)
		return nil
	},
	// list <type prefix>: log the objects whose type ID starts with the prefix
	"list": func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: list <type prefix>")
		}
		n := 0
		for _, obj := range append(noxServer.Objs.AllNonMissiles(), noxServer.Objs.AllMissiles()...) {
			if name := e2eItemName(obj); strings.HasPrefix(strings.ToLower(name), strings.ToLower(args[0])) {
				hp, max := 0, 0
				if obj.HealthData != nil {
					hp, max = int(obj.HealthData.Cur), int(obj.HealthData.Max)
				}
				e2eLog.Printf("PROBE OBJ %s #%d at (%.0f,%.0f) hp %d/%d flags %#x", name, obj.NetCode, obj.PosVec.X, obj.PosVec.Y, hp, max, uint32(obj.Flags()))
				n++
			}
		}
		e2eLog.Printf("PROBE LIST %s: %d", args[0], n)
		return nil
	},
	// damageid <net code> <amount> [kind]: damage one object
	"damageid": func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: damageid <net code> <amount> [kind]")
		}
		code, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		kind := 0
		if len(args) > 2 {
			if kind, err = strconv.Atoi(args[2]); err != nil {
				return err
			}
		}
		for obj := noxServer.Objs.First(); obj != nil; obj = obj.Next() {
			if int(obj.NetCode) == code {
				e2eLog.Printf("PROBE DAMAGE %s #%d by %d (kind %d)", e2eItemName(obj), obj.NetCode, amount, kind)
				asObjectS(obj).CallDamage(e2eQuestPlayer(), nil, amount, object.DamageType(kind))
				return nil
			}
		}
		return fmt.Errorf("no object #%d", code)
	},
	// find <script name>: look up an object by its script name, like the map scripts
	"find": func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: find <name>")
		}
		obj := legacy.Nox_xxx_getObjectByScrName_4DA4F0(args[0])
		if obj == nil {
			e2eLog.Printf("PROBE FIND %q: not found", args[0])
			return nil
		}
		e2eLog.Printf("PROBE FIND %q: %s #%d at (%.0f,%.0f)", args[0], e2eItemName(obj), obj.NetCode, obj.PosVec.X, obj.PosVec.Y)
		return nil
	},
	// names [substring]: log the objects' script names (containing the substring)
	"names": func(args []string) error {
		n, named := 0, 0
		for obj := noxServer.Objs.First(); obj != nil; obj = obj.Next() {
			n++
			id := obj.ID()
			if id == "" {
				continue
			}
			named++
			if len(args) == 0 || strings.Contains(strings.ToLower(id), strings.ToLower(args[0])) {
				e2eLog.Printf("PROBE NAME %q %s #%d", id, e2eItemName(obj), obj.NetCode)
			}
		}
		e2eLog.Printf("PROBE NAMES: %d of %d objects have names", named, n)
		return nil
	},
	// group <name>: look up an object group, like the map scripts
	"group": func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: group <name>")
		}
		g := noxServer.MapGroups.GroupByID(args[0], server.MapGroupObjects)
		e2eLog.Printf("PROBE GROUP %q: found %v", args[0], g != nil)
		return nil
	},
	// spawn <type> <x> <y> [vx vy]: create an object owned by the host's player, moving with the velocity;
	// x and y are relative to the player if they start with "~"
	"spawn": func(args []string) error {
		if len(args) != 3 && len(args) != 5 {
			return fmt.Errorf("usage: spawn <type> <x> <y> [vx vy]")
		}
		u := e2eQuestPlayer()
		if u == nil {
			return fmt.Errorf("no player")
		}
		var vals [4]float64
		rel := [2]bool{}
		for i, a := range args[1:] {
			if i < 2 && strings.HasPrefix(a, "~") {
				rel[i] = true
				a = a[1:]
				if a == "" {
					a = "0"
				}
			}
			v, err := strconv.ParseFloat(a, 64)
			if err != nil {
				return err
			}
			vals[i] = v
		}
		pos := types.Pointf{X: float32(vals[0]), Y: float32(vals[1])}
		if rel[0] {
			pos.X += u.PosVec.X
		}
		if rel[1] {
			pos.Y += u.PosVec.Y
		}
		obj := noxServer.NewObjectByTypeID(args[0])
		if obj == nil {
			return fmt.Errorf("unknown type %q", args[0])
		}
		noxServer.CreateObjectAt(obj, u, pos)
		obj.VelVec = types.Pointf{X: float32(vals[2]), Y: float32(vals[3])}
		if vals[2] != 0 || vals[3] != 0 {
			d := server.DirFromVec(obj.VelVec)
			obj.Direction1, obj.Direction2 = d, d
		}
		e2eLog.Printf("PROBE SPAWN %s #%d at (%.0f,%.0f) vel (%.1f,%.1f)", args[0], obj.NetCode, pos.X, pos.Y, vals[2], vals[3])
		return nil
	},
	// gen <monster type> <max>: create a monster generator next to the host with the script API (opennox#666)
	"gen": func(args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("usage: gen <monster type> <max>")
		}
		max, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		u := e2eQuestPlayer()
		if u == nil {
			return fmt.Errorf("no player")
		}
		ns := noxServer.noxScriptP()
		pos := u.Pos()
		pos.X += 150
		g := ns.CreateObject("MonsterGenerator", ns4.Ptf(pos.X, pos.Y))
		if g == nil {
			return fmt.Errorf("cannot create a generator")
		}
		if !g.SetGeneratorMonsters(max, ns4.ObjTypeName(args[0])) {
			return fmt.Errorf("SetGeneratorMonsters failed")
		}
		e2eLog.Printf("PROBE GEN %s max %d at (%.0f,%.0f)", args[0], max, pos.X, pos.Y)
		return nil
	},
	// genstat: log the monster generators' spawn counts
	"genstat": func(args []string) error {
		for obj := noxServer.Objs.First(); obj != nil; obj = obj.Next() {
			if obj.Class().Has(object.ClassMonsterGenerator) {
				d := obj.UpdateDataMonsterGen()
				e2eLog.Printf("PROBE GENSTAT %s #%d current %d max %d rates %v flags %#x", e2eItemName(obj), obj.NetCode, d.Current, d.Max, d.Rates, d.Flags92)
			}
		}
		return nil
	},
	// loadmap <map> [hide]: load a map with the script API (LoadMap), "hide" hides the title screen and the scores
	"loadmap": func(args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("usage: loadmap <map> [hide]")
		}
		opts := &ns4.LoadMapOptions{}
		if len(args) > 1 && args[1] == "hide" {
			opts.HideTitleScreen, opts.HideScores = true, true
		}
		return noxServer.noxScriptP().loadMap(args[0], opts)
	},
	// hidequest <title> <scores>: HideQuestScreens of the script API (true/false)
	"hidequest": func(args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("usage: hidequest <title> <scores>")
		}
		noxServer.noxScriptP().HideQuestScreens(args[0] == "true", args[1] == "true")
		return nil
	},
	// chat <string id>: the host's player says a line like the script Chat (which also plays its voice file)
	"chat": func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: chat <string id>")
		}
		u := e2eQuestPlayer()
		if u == nil {
			return fmt.Errorf("no player")
		}
		u.Chat(ns4.StringID(args[0]))
		return nil
	},
	// dialog: log the state of the dialog voice player
	"dialog": func(args []string) error {
		e2eLog.Printf("PROBE DIALOG %s", legacy.Dialogs.GoString())
		return nil
	},
	// savedata set <key> <value> | get <key> <expected>: the scripts' SetSaveData and SaveData
	"savedata": func(args []string) error {
		ns := noxServer.noxScriptP()
		switch {
		case len(args) == 3 && args[0] == "set":
			ns.SetSaveData(args[1], []byte(args[2]))
			return nil
		case len(args) == 3 && args[0] == "get":
			got := string(ns.SaveData(args[1]))
			e2eLog.Printf("PROBE SAVEDATA %s = %q", args[1], got)
			if got != args[2] {
				return fmt.Errorf("save data %q: %q, expected %q", args[1], got, args[2])
			}
			return nil
		}
		return fmt.Errorf("usage: savedata set <key> <value> | get <key> <expected>")
	},
	// save <name>: save the solo game (e.g. SAVE0009), like the game menu
	"save": func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: save <name>")
		}
		if !saveCoopGame(args[0]) {
			return fmt.Errorf("cannot save %s", args[0])
		}
		return nil
	},
	// loadsave <name>: load a solo save, like the game menu
	"loadsave": func(args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: loadsave <name>")
		}
		if !clientLoadCoopGame(args[0]) {
			return fmt.Errorf("cannot load %s", args[0])
		}
		return nil
	},
	// types: log how many objects of each type there are
	"types": func(args []string) error {
		cnt := make(map[string]int)
		for obj := noxServer.Objs.First(); obj != nil; obj = obj.Next() {
			cnt[e2eItemName(obj)]++
		}
		for name, n := range cnt {
			e2eLog.Printf("PROBE TYPE %s %d", name, n)
		}
		return nil
	},
	// walls [x y r]: log the breakable walls (near a point, if given)
	"walls": func(args []string) error {
		near := len(args) == 3
		var p types.Pointf
		var r float64
		if near {
			var err error
			if p, err = e2eProbePos(args[:2]); err != nil {
				return err
			}
			if r, err = strconv.ParseFloat(args[2], 64); err != nil {
				return err
			}
		}
		n := 0
		for it := noxServer.Walls.FirstBreakable(); it != nil; it = it.Next() {
			wl := it.Wall
			x, y := float64(wl.X5)*23+11.5, float64(wl.Y6)*23+11.5
			if near && (x-float64(p.X))*(x-float64(p.X))+(y-float64(p.Y))*(y-float64(p.Y)) > r*r {
				continue
			}
			e2eLog.Printf("PROBE WALL (%d,%d) at (%.0f,%.0f) health %d broken %v", wl.X5, wl.Y6, x, y, wl.Health7, wl.Flags4.Has(wall.FlagBroken))
			n++
		}
		e2eLog.Printf("PROBE WALLS: %d breakable", n)
		return nil
	},
	// damage <type prefix> <amount> [kind]: damage the objects of the type (kind: object.DamageType, default 0)
	"damage": func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: damage <type prefix> <amount> [kind]")
		}
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		kind := 0
		if len(args) > 2 {
			if kind, err = strconv.Atoi(args[2]); err != nil {
				return err
			}
		}
		var list []*server.Object
		for obj := noxServer.Objs.First(); obj != nil; obj = obj.Next() {
			if strings.HasPrefix(strings.ToLower(e2eItemName(obj)), strings.ToLower(args[0])) && !obj.Flags().Has(object.FlagDestroyed) {
				list = append(list, obj)
			}
		}
		host := e2eQuestPlayer()
		for _, obj := range list {
			e2eLog.Printf("PROBE DAMAGE %s #%d by %d", e2eItemName(obj), obj.NetCode, amount)
			asObjectS(obj).CallDamage(host, nil, amount, object.DamageType(kind))
		}
		return nil
	},
}

func e2eProbePos(args []string) (types.Pointf, error) {
	if len(args) != 2 {
		return types.Pointf{}, fmt.Errorf("expected x y")
	}
	x, err := strconv.ParseFloat(args[0], 32)
	if err != nil {
		return types.Pointf{}, err
	}
	y, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return types.Pointf{}, err
	}
	return types.Pointf{X: float32(x), Y: float32(y)}, nil
}

// Probe runs a probe from e2eProbes.
func (sc *e2eScenario) Probe(args []string, name string) {
	sc.add(0, name, func() {
		if len(args) == 0 || e2eProbes[args[0]] == nil {
			e2eError(fmt.Errorf("%s: unknown probe %q", name, args))
			return
		}
		e2eLog.Printf("PROBE %v", args)
		if err := e2eProbes[args[0]](args[1:]); err != nil {
			e2eError(fmt.Errorf("%s: probe %v: %w", name, args, err))
		}
	})
}
