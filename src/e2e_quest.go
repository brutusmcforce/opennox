package opennox

import (
	"fmt"
	"os"

	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/libs/ifs"
	crypt "github.com/brutusmcforce/opennox/noxcrypt"

	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/server"
)

// Quest test actions (e2e/quest.yaml, run with -autoquest): they check the items of a quest game through the server
// state, and save and load the host's character with the code that a quest uses between maps.

var e2eQuestVerbose = os.Getenv("NOX_E2E_QUEST_VERBOSE") == "true"

// e2eQuestPlayer returns the host's player unit.
func e2eQuestPlayer() *Object {
	host := noxServer.Players.Host()
	if host == nil || host.PlayerUnit == nil {
		e2eError(fmt.Errorf("quest: no host player"))
		return nil
	}
	return asObjectS(host.PlayerUnit)
}

func e2eItemName(it *server.Object) string {
	if t := it.ObjectTypeC(); t != nil {
		return t.ID()
	}
	return fmt.Sprintf("type#%d", it.TypeInd)
}

// e2eEachItem calls fnc for every object in the world and every item in their inventories (and in those, e.g. chests
// in containers), with the object holding it (nil in the world).
func e2eEachItem(fnc func(it, holder *server.Object)) {
	var inv func(holder *server.Object)
	inv = func(holder *server.Object) {
		for it := holder.FirstItem(); it != nil; it = it.NextItem() {
			fnc(it, holder)
			inv(it)
		}
	}
	for obj := noxServer.Objs.First(); obj != nil; obj = obj.Next() {
		fnc(obj, nil)
		inv(obj)
	}
}

// QuestItems checks that every weapon and armor (in the world, in containers, held by monsters and players) has full
// durability, scaled for quests (opennox#740: the map's items had the base durability; a saved character
// loaded them as damaged). It logs how many it checked, and fails if there are none.
func (sc *e2eScenario) QuestItems(name string) {
	sc.add(0, name, func() {
		n := 0
		e2eEachItem(func(it, holder *server.Object) {
			want, _, ok := legacy.E2EItemDurability(it)
			if !ok || it.Flags().Has(object.FlagDestroyed) {
				return
			}
			n++
			if e2eQuestVerbose {
				e2eLog.Printf("QUEST ITEM %s holder=%v durability %d/%d want %d flags=%#x", e2eItemName(it), holder != nil, it.HealthData.Cur, it.HealthData.Max, want, uint32(it.Flags()))
			}
			if it.HealthData.Cur != want || it.HealthData.Max != want {
				where := "world"
				if holder != nil {
					where = e2eItemName(holder)
				}
				e2eError(fmt.Errorf("%s: %s (in %s) has durability %d/%d, expected %d/%d",
					name, e2eItemName(it), where, it.HealthData.Cur, it.HealthData.Max, want, want))
			}
		})
		e2eLog.Printf("QUEST ITEMS: %d weapons and armor checked", n)
		if n == 0 {
			e2eError(fmt.Errorf("%s: no weapons or armor found", name))
		}
	})
}

// QuestTake makes the host pick up the first weapon or armor (of the item type, if set) that lies on the map or in a
// container, so that it's saved with the character.
func (sc *e2eScenario) QuestTake(item, name string) {
	sc.add(0, name, func() {
		u := e2eQuestPlayer()
		if u == nil {
			return
		}
		var found, from *server.Object
		e2eEachItem(func(it, holder *server.Object) {
			if found != nil || it.Flags().Has(object.FlagDestroyed) || !it.Class().HasAny(object.ClassWeapon|object.ClassArmor) {
				return
			}
			if holder != nil && holder.Class().HasAny(object.MaskUnits) {
				return // carried by a monster or a player
			}
			if item != "" && e2eItemName(it) != item {
				return
			}
			if _, _, ok := legacy.E2EItemDurability(it); ok {
				found, from = it, holder
			}
		})
		if found == nil {
			e2eError(fmt.Errorf("%s: no weapon or armor on the map", name))
			return
		}
		if from != nil && !asObjectS(from).DoDrop(found) {
			e2eError(fmt.Errorf("%s: cannot take %s out of %s", name, e2eItemName(found), e2eItemName(from)))
			return
		}
		if !u.DoPickup(found) {
			e2eError(fmt.Errorf("%s: cannot pick up %s", name, e2eItemName(found)))
			return
		}
		e2eLog.Printf("QUEST TAKE: %s, durability %d/%d", e2eItemName(found), found.HealthData.Cur, found.HealthData.Max)
	})
}

// QuestMapItem places a weapon or armor next to the host the way a map file has it: an item placed in the map editor
// has the base durability, and the map loader reads it (with the object's xfer function, outside a character load).
// Quest maps place their items with scripts, so this is how custom quest maps get them (opennox#740).
func (sc *e2eScenario) QuestMapItem(item, name string) {
	sc.add(0, name, func() {
		u := e2eQuestPlayer()
		if u == nil {
			return
		}
		src := noxServer.NewObjectByTypeID(item)
		if src == nil {
			e2eError(fmt.Errorf("%s: unknown item type %q", name, item))
			return
		}
		_, base, ok := legacy.E2EItemDurability(src)
		if !ok {
			e2eError(fmt.Errorf("%s: %s is not a weapon or armor", name, item))
			return
		}
		src.HealthData.Cur, src.HealthData.Max = base, base
		src.PosVec = u.Pos()
		src.PosVec.X += 30
		path := savePath("_e2e_mapitem_.obj")
		defer ifs.Remove(path)
		if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, crypt.MapKey); err != nil {
			e2eError(fmt.Errorf("%s: %w", name, err))
			return
		}
		err := src.CallXfer(nil)
		_ = cryptfile.Close()
		noxServer.Objs.FreeObject(src) // never placed
		if err != nil {
			e2eError(fmt.Errorf("%s: write: %w", name, err))
			return
		}
		obj := noxServer.NewObjectByTypeID(item)
		if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, crypt.MapKey); err != nil {
			e2eError(fmt.Errorf("%s: %w", name, err))
			return
		}
		err = obj.CallXfer(nil)
		_ = cryptfile.Close()
		if err != nil {
			e2eError(fmt.Errorf("%s: read: %w", name, err))
			return
		}
		legacy.Nox_xxx_servMapLoadPlaceObj_4F3F50(obj, 0, nil)
		noxServer.ObjectsAddPending()
		e2eLog.Printf("QUEST MAP ITEM: %s saved with durability %d, loaded as %d/%d", item, base, obj.HealthData.Cur, obj.HealthData.Max)
	})
}

// QuestGive creates an item by its type and gives it to the host.
func (sc *e2eScenario) QuestGive(item, name string) {
	sc.add(0, name, func() {
		u := e2eQuestPlayer()
		if u == nil {
			return
		}
		it := noxServer.NewObjectByTypeID(item)
		if it == nil {
			e2eError(fmt.Errorf("%s: unknown item type %q", name, item))
			return
		}
		noxServer.CreateObjectAt(it, nil, u.Pos())
		noxServer.ObjectsAddPending()
		if !u.DoPickup(it) {
			e2eError(fmt.Errorf("%s: cannot pick up %s", name, item))
			return
		}
		e2eLog.Printf("QUEST GIVE: %s (allowed in quests: %v)", item, legacy.E2EQuestItemAllowed(it))
	})
}

// QuestHas checks whether the host carries an item of the type (want) or not.
func (sc *e2eScenario) QuestHas(item string, want bool, name string) {
	sc.add(0, name, func() {
		u := e2eQuestPlayer()
		if u == nil {
			return
		}
		has := false
		for it := u.FirstItem(); it != nil; it = it.NextItem() {
			if !it.Flags().Has(object.FlagDestroyed) && e2eItemName(it) == item {
				has = true
			}
		}
		e2eLog.Printf("QUEST HAS %s: %v", item, has)
		if has != want {
			e2eError(fmt.Errorf("%s: carrying %s: %v, expected %v", name, item, has, want))
		}
	})
}

func e2eQuestSavePath() string {
	return savePath("_e2e_quest_.plr")
}

// QuestSave saves the host's character like a quest does before the next map.
func (sc *e2eScenario) QuestSave(name string) {
	sc.add(0, name, func() {
		u := e2eQuestPlayer()
		if u == nil {
			return
		}
		if !savePlayerServerData(e2eQuestSavePath(), u.ControllingPlayer().PlayerIndex()) {
			e2eError(fmt.Errorf("%s: cannot save the character", name))
			return
		}
		e2eLog.Printf("QUEST SAVE: %s", e2eQuestSavePath())
	})
}

// QuestLoad gives the host the character saved by QuestSave, the way a quest does it for a host that continues with
// a saved character: the file is queued (sub4DCEE0), and the server loads it into the host's player on the next
// frames (sub_4DCF20), which replaces the inventory. The game deletes the queued file.
func (sc *e2eScenario) QuestLoad(name string) {
	sc.add(0, name, func() {
		path := savePath("_e2e_quest_load_.plr")
		if err := ifs.Copy(e2eQuestSavePath(), path); err != nil {
			e2eError(fmt.Errorf("%s: %w", name, err))
			return
		}
		sub4DCEE0(path)
		e2eLog.Printf("QUEST LOAD: queued %s", path)
	})
}

// QuestExit moves the host onto the map's exit, which finishes the quest stage and loads the next map.
func (sc *e2eScenario) QuestExit(name string) {
	sc.add(0, name, func() {
		u := e2eQuestPlayer()
		if u == nil {
			return
		}
		// the warp gates (exit subclass 2) only open after some stages
		exit := server.CollideFuncByName("ExitCollide")
		for obj := noxServer.Objs.First(); obj != nil; obj = obj.Next() {
			if exit != nil && obj.Collide == exit && obj.ObjSubClass&2 == 0 {
				e2eLog.Printf("QUEST EXIT: %s at (%.0f,%.0f)", e2eItemName(obj), obj.PosVec.X, obj.PosVec.Y)
				u.SetPos(obj.Pos())
				return
			}
		}
		e2eError(fmt.Errorf("%s: no exit on the map", name))
	})
}

// Drawables logs how many drawables the client uses (opennox#715: a long quest ran out of them).
func (sc *e2eScenario) Drawables(name string) {
	sc.add(0, name, func() {
		n, max := noxServer.Objs.Active()
		e2eLog.Printf("DRAWABLES: %d active (preallocated %d), server objects %d of %d",
			noxClient.Objs.Alloc.Active(), noxClient.Objs.Alloc.Cap(), n, max)
	})
}
