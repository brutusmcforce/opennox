package opennox

import (
	"fmt"
	"image/color"
	"math"

	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/player"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/libs/types"
	ns4 "github.com/brutusmcforce/opennox/noxscript/ns/v4"
	"github.com/brutusmcforce/opennox/noxscript/ns/v4/class"
	"github.com/brutusmcforce/opennox/noxscript/ns/v4/damage"
	"github.com/brutusmcforce/opennox/noxscript/ns/v4/enchant"
	"github.com/brutusmcforce/opennox/noxscript/ns/v4/subclass"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/server"
)

func (s noxScriptNS) toObj(obj *server.Object) ns4.Obj {
	if obj == nil {
		return nil
	}
	return nsObj{s.s, asObjectS(obj)}
}

func (s noxScriptNS) asObj(obj ns4.Obj) *server.Object {
	if obj == nil {
		return nil
	}
	o, _ := obj.(nsObj)
	return o.Object.SObj()
}

func (s noxScriptNS) ObjectType(name ns4.ObjTypeName) ns4.ObjType {
	typ := s.s.Types.ByID(name)
	if typ == nil {
		return nil
	}
	return nsObjType{s.s, typ}
}

func (s noxScriptNS) ObjectTypeByInd(ind int) ns4.ObjType {
	typ := s.s.Types.ByInd(ind)
	if typ == nil {
		return nil
	}
	return nsObjType{s.s, typ}
}

func (s noxScriptNS) ObjectByHandle(h ns4.ObjHandle) ns4.Obj {
	if h == nil {
		return nil
	}
	obj := s.s.noxScript.ScriptToObject(h.ObjScriptID())
	return s.toObj(obj)
}

func (s noxScriptNS) createObject(typ ns4.ObjTypeName, p ns4.Positioner) *server.Object {
	if p == nil {
		return nil
	}
	obj := s.s.NewObjectByTypeID(string(typ))
	if obj == nil {
		return nil
	}
	s.s.CreateObjectAt(obj, nil, p.Pos())
	if noxflags.HasGame(noxflags.GameModeQuest) {
		if item := legacy.RewardFromScriptMarker(obj); item != nil {
			return item
		}
	}
	return obj
}

func (s noxScriptNS) CreateObject(typ ns4.ObjTypeName, p ns4.Positioner) ns4.Obj {
	obj := s.createObject(typ, p)
	if obj == nil {
		return nil
	}
	return s.toObj(obj)
}

func (s noxScriptNS) Object(name string) ns4.Obj {
	obj := legacy.Nox_xxx_getObjectByScrName_4DA4F0(name)
	if obj == nil {
		ScriptLog.Printf("noxscript: cannot find object: %q", name)
		return nil
	}
	return s.toObj(obj)
}

func (s noxScriptNS) ObjectGroupByHandle(h ns4.ObjGroupHandle) ns4.ObjGroup {
	if h == nil {
		return nil
	}
	g := s.s.MapGroups.GroupByInd(h.ObjGroupScriptID())
	if g == nil || s.s.MapGroups.MapGroupType(g) != server.MapGroupObjects {
		return nil
	}
	return nsObjGroup{s.s, g}
}

func (s noxScriptNS) ObjectGroup(name string) ns4.ObjGroup {
	g := s.s.MapGroups.GroupByID(name, server.MapGroupObjects)
	if g == nil {
		ScriptLog.Printf("noxscript: cannot find object group: %q", name)
		return nil
	}
	return nsObjGroup{s.s, g}
}

func (s noxScriptNS) FindObjects(iter func(obj ns4.Obj) bool, conditions ...ns4.ObjCond) int {
	// This generic iteration function will act as a fallback that scans everything.
	search := func(fnc func(obj *server.Object) bool) {
		for it := s.s.Objs.First(); it != nil; it = it.Next() {
			if !fnc(it) {
				return
			}
		}
		for it := s.s.Objs.MissileList; it != nil; it = it.Next() {
			if !fnc(it) {
				return
			}
		}
	}
	// Attempt to find a more specific iterator ad optimize it.
	for i, c := range conditions {
		found := false
		switch c := c.(type) {
		case ns4.InRectf:
			search = func(fnc func(obj *server.Object) bool) {
				s.s.Map.EachObjAndMissileInRect(types.RectFromPointsf(c.Min, c.Max), fnc)
			}
			found = true
		case ns4.InCirclef:
			search = func(fnc func(obj *server.Object) bool) {
				s.s.Map.EachObjAndMissileInCircle(c.Center.Pos(), float32(c.R), fnc)
			}
			found = true
		}
		// We only need one optimized iterator. Once found - remove it from conditions.
		if found {
			conditions = append(conditions[:i], conditions[i+1:]...)
			break
		}
	}
	// Finally start the actual iteration.
	filter := ns4.AND(conditions)
	cnt := 0
	search(func(obj *server.Object) bool {
		nobj := s.toObj(obj)
		if !filter.Matches(nobj) {
			return true // find next match
		}
		// Found, pass to user-defined handler func.
		cnt++
		if iter != nil && !iter(nobj) {
			return false
		}
		return true
	})
	return cnt
}

func (s noxScriptNS) GetTrigger() ns4.Obj {
	obj := s.s.noxScript.Trigger()
	return s.toObj(obj)
}

func (s noxScriptNS) GetCaller() ns4.Obj {
	obj := s.s.noxScript.Caller()
	return s.toObj(obj)
}

func (s noxScriptNS) DestroyEveryChat() {
	legacy.Nox_xxx_destroyEveryChatMB_528D60()
}

func (s noxScriptNS) MakeFriendly(obj ns4.Obj) {
	legacy.ScriptOwnership(s.asObj(obj), 0)
}

func (s noxScriptNS) MakeEnemy(obj ns4.Obj) {
	legacy.ScriptOwnership(s.asObj(obj), 1)
}

func (s noxScriptNS) BecomePet(obj ns4.Obj) {
	legacy.ScriptOwnership(s.asObj(obj), 2)
}

func (s noxScriptNS) BecomeEnemy(obj ns4.Obj) {
	legacy.ScriptOwnership(s.asObj(obj), 3)
}

type nsObjType struct {
	s *Server
	t *server.ObjectType
}

func (typ nsObjType) Name() ns4.ObjTypeName {
	return ns4.ObjTypeName(typ.t.ID())
}

func (typ nsObjType) Index() int {
	return typ.t.Ind()
}

func (typ nsObjType) Create(p ns4.Positioner) ns4.Obj {
	if p == nil {
		return nil
	}
	obj := typ.s.Objs.NewObject(typ.t)
	if obj == nil {
		return nil
	}
	typ.s.CreateObjectAt(obj, nil, p.Pos())
	return typ.s.noxScriptP().toObj(obj)
}

func (typ nsObjType) Class() object.Class {
	return typ.t.Class()
}

func (typ nsObjType) HasClass(class class.Class) bool {
	cls, err := object.ParseClass(string(class))
	if err != nil {
		return false
	}
	return typ.t.Class().Has(cls)
}

func (typ nsObjType) HasSubclass(subclass subclass.SubClass) bool {
	cls, err := object.ParseSubClass(string(subclass))
	if err != nil {
		return false
	}
	return typ.t.SubClass().Has(cls)
}

func (typ nsObjType) Flags() object.Flags {
	return typ.t.Flags()
}

var _ server.Obj = nsObj{}

type nsObj struct {
	s *Server
	*Object
}

func (obj nsObj) CursorPos() types.Pointf {
	if !obj.Class().Has(object.ClassPlayer) {
		return types.Pointf{}
	}
	pl := obj.Object.ControllingPlayer()
	return pl.CursorPos()
}

func (obj nsObj) CursorObj() ns4.Obj {
	if !obj.Class().Has(object.ClassPlayer) {
		return nil
	}
	ud := obj.Object.UpdateDataPlayer()
	return obj.s.noxScriptP().toObj(ud.CursorObj)
}

func (obj nsObj) ScriptID() int {
	return obj.Object.ScriptIDVal
}

func (obj nsObj) ObjScriptID() int {
	return obj.Object.ScriptIDVal
}

func (obj nsObj) Type() ns4.ObjType {
	return nsObjType{obj.s, obj.Object.ObjectTypeC()}
}

func (obj nsObj) Mass() float32 {
	return obj.Object.Mass
}

func (obj nsObj) SetMass(v float32) {
	obj.Object.Mass = v
	obj.Object.NeedSync()
}

func (obj nsObj) Team() ns4.Team {
	t := obj.Object.Team()
	if t == nil {
		return nil
	}
	return nsTeam{obj.s, t}
}

func (obj nsObj) HasTeam(t ns4.Team) bool {
	if t == nil {
		return !obj.Object.HasTeam()
	}
	tm, ok := t.(nsTeam)
	if !ok {
		return false
	}
	return obj.TeamEqual(tm.t)
}

func (obj nsObj) SetTeam(t ns4.Team) {
	if t == nil {
		return // TODO: support clearing the team
	}
	tm, ok := t.(nsTeam)
	if !ok {
		return
	}
	obj.Object.SetTeam(tm.t)
}

func (obj nsObj) HasClass(class class.Class) bool {
	cls, err := object.ParseClass(string(class))
	if err != nil {
		return false
	}
	return obj.Class().Has(cls)
}

func (obj nsObj) HasSubclass(subclass subclass.SubClass) bool {
	cls, err := object.ParseSubClass(string(subclass))
	if err != nil {
		return false
	}
	return obj.SubClass().Has(cls)
}

func (obj nsObj) HasEnchant(enc enchant.Enchant) bool {
	e, ok := server.ParseEnchant(string(enc))
	if !ok {
		return false
	}
	return obj.Object.HasEnchant(e)
}

func (obj nsObj) MonsterStatus() object.MonsterStatus {
	if !obj.Class().Has(object.ClassMonster) {
		return 0
	}
	return obj.UpdateDataMonster().StatusFlags
}

func (obj nsObj) Direction() ns4.Direction {
	return ns4.Direction(obj.Dir1())
}

func (obj nsObj) CurrentHealth() int {
	v, _ := obj.Health()
	return v
}

func (obj nsObj) MaxHealth() int {
	_, v := obj.Health()
	return v
}

func (obj nsObj) CurrentMana() int {
	v, _ := obj.Mana()
	return v
}

func (obj nsObj) MaxMana() int {
	_, v := obj.Mana()
	return v
}

func (obj nsObj) GiveXp(xp float32) {
	legacy.ScriptGiveXP(obj.SObj(), xp)
}

func (obj nsObj) Player() ns4.Player {
	if !obj.Class().Has(object.ClassPlayer) {
		return nil
	}
	pl := obj.ControllingPlayer()
	return obj.s.noxScriptP().toPlayer(pl)
}

func (obj nsObj) GetScore() int {
	pl := obj.Player()
	if pl == nil {
		return 0
	}
	return pl.GetScore()
}

func (obj nsObj) GetClass() player.Class {
	if !obj.Class().Has(object.ClassPlayer) {
		return 0
	}
	return obj.ControllingPlayer().PlayerClass()
}

func (obj nsObj) GetLevel() int {
	if !obj.Class().Has(object.ClassPlayer) {
		return 0
	}
	return int(obj.ControllingPlayer().Level)
}

func (obj nsObj) HasOwner(owner ns4.Obj) bool {
	own, _ := owner.(server.Obj)
	return obj.Object.HasOwner(server.ToObject(own))
}

func (obj nsObj) HasOwnerIn(owners ns4.ObjGroup) bool {
	if owners == nil {
		return false
	}
	ok := false
	owners.EachObject(false, func(own ns4.Obj) bool {
		ok = obj.HasOwner(own)
		return !ok // stop if has any owner
	})
	return ok
}

func (obj nsObj) SetOwner(owner ns4.Obj) {
	if owner == nil {
		obj.Object.SetOwner(nil)
		return
	}
	v := server.ToObject(owner.(server.Obj))
	obj.Object.SetOwner(v)
}

func (obj nsObj) SetOwners(owners ns4.ObjGroup) {
	if owners == nil {
		obj.Object.SetOwner(nil)
		return
	}
	owners.EachObject(false, func(it ns4.Obj) bool {
		obj.SetOwner(it)
		return true
	})
}

func (obj nsObj) CanSee(obj2 ns4.Obj) bool {
	if obj.Object == nil || obj2 == nil {
		return false
	}
	return obj.Object.CanSee(obj.s.noxScriptP().asObj(obj2))
}

func (obj nsObj) PushTo(p ns4.Positioner, force float32) {
	if p != nil {
		obj.Push(p.Pos(), force)
	}
}

func (obj nsObj) Damage(source ns4.Obj, amount int, typ damage.Type) bool {
	var src *server.Object
	if source != nil {
		src = obj.s.noxScriptP().asObj(source)
	}
	return obj.Object.DoDamage(src, amount, object.DamageType(typ))
}

// playerActions maps player states to actions, see controlActionState for all states.
var playerActions = map[server.PlayerState]ns4.PlayerAction{
	server.PlayerState13:        ns4.PlayerIdle,
	server.PlayerState0:         ns4.PlayerWalk,
	server.PlayerState5:         ns4.PlayerRun,
	server.PlayerState1:         ns4.PlayerAttack,
	server.PlayerState14:        ns4.PlayerAttack,
	server.PlayerState22:        ns4.PlayerAttack,
	server.PlayerState2:         ns4.PlayerCast,
	server.PlayerState10:        ns4.PlayerCast,
	server.PlayerState12:        ns4.PlayerJump,
	server.PlayerState15:        ns4.PlayerBlock,
	server.PlayerState16:        ns4.PlayerBlock,
	server.PlayerState17:        ns4.PlayerBlock,
	server.PlayerState3:         ns4.PlayerDying,
	server.PlayerState4:         ns4.PlayerDead,
	server.PlayerStateShakeFist: ns4.PlayerShakeFist,
	server.PlayerStateLaugh:     ns4.PlayerLaugh,
	server.PlayerStatePoint:     ns4.PlayerPoint,
}

func (obj nsObj) PlayerAction() ns4.PlayerAction {
	u := obj.SObj()
	if !u.Class().Has(object.ClassPlayer) {
		return ""
	}
	if a, ok := playerActions[u.UpdateDataPlayer().State]; ok {
		return a
	}
	return ns4.PlayerOther
}

func (obj nsObj) SetPlayerAction(a ns4.PlayerAction) bool {
	u := obj.SObj()
	if !u.Class().Has(object.ClassPlayer) {
		return false
	}
	// only actions that don't need more state (like a target or a spell)
	var st server.PlayerState
	switch a {
	case ns4.PlayerIdle:
		st = server.PlayerState13
	case ns4.PlayerShakeFist:
		st = server.PlayerStateShakeFist
	case ns4.PlayerLaugh:
		st = server.PlayerStateLaugh
	case ns4.PlayerPoint:
		st = server.PlayerStatePoint
	default:
		return false
	}
	return nox_xxx_playerSetState_4FA020(u, st)
}

func (obj nsObj) MonsterAction() ns4.MonsterAction {
	u := obj.SObj()
	if !u.Class().Has(object.ClassMonster) {
		return ""
	}
	it := u.UpdateDataMonster().AIStackHead()
	if it == nil {
		return ns4.MonsterAction(ai.ACTION_IDLE.String())
	}
	return ns4.MonsterAction(ai.ActionType(it.Action).String())
}

func (obj nsObj) IsRunning() bool {
	u := obj.SObj()
	switch {
	case u.Class().Has(object.ClassPlayer):
		return u.UpdateDataPlayer().State == server.PlayerState5
	case u.Class().Has(object.ClassMonster):
		return u.UpdateDataMonster().StatusFlags.Has(object.MonStatusRunning)
	}
	return false
}

// itemClasses are the classes of items with modifiers (weapons, armor, wands, flags).
const itemClasses = object.ClassWeapon | object.ClassArmor | object.ClassWand | object.ClassFlag

// itemModifiers returns the modifiers of an item, or nil if it has none.
func (obj nsObj) itemModifiers() *server.ModifierInitData {
	u := obj.SObj()
	if !u.Class().HasAny(itemClasses) || u.InitData == nil {
		return nil
	}
	return u.InitDataModifier()
}

func (obj nsObj) ItemModifiers() [4]string {
	var out [4]string
	if d := obj.itemModifiers(); d != nil {
		for i, m := range d.Modifiers {
			if m != nil {
				out[i] = m.Name()
			}
		}
	}
	return out
}

func (obj nsObj) SetItemModifier(slot int, name string) bool {
	d := obj.itemModifiers()
	if d == nil || slot < 0 || slot >= len(d.Modifiers) {
		return false
	}
	var m *server.ModifierEff
	if name != "" {
		if m = obj.s.Modif.Nox_xxx_modifGetDescById413330(obj.s.Modif.Nox_xxx_modifGetIdByName413290(name)); m == nil {
			return false
		}
	}
	attrs := *d
	attrs.Modifiers[slot] = m
	// this also sends the change to the clients
	legacy.Nox_xxx_modifSetItemAttrs_4E4990(obj.SObj(), &attrs)
	return true
}

// setItemColor sets a modifier slot of an item to the closest user color (UserColor1-32), like the game colors clothes.
func (obj nsObj) setItemColor(slot int, cl color.Color) {
	r, g, b, _ := cl.RGBA()
	best, bestDist := 0, -1
	for i, c := range obj.s.Modif.Colors[0] {
		dr, dg, db := int(c.R)-int(r>>8), int(c.G)-int(g>>8), int(c.B)-int(b>>8)
		if d := dr*dr + dg*dg + db*db; bestDist < 0 || d < bestDist {
			best, bestDist = i, d
		}
	}
	obj.SetItemModifier(slot, fmt.Sprintf("UserColor%d", best+1))
}

func (obj nsObj) SetColor(ind int, cl ns4.Color) {
	if cl != nil && obj.itemModifiers() != nil {
		obj.setItemColor(ind, cl)
		return
	}
	obj.Object.SetColor(ind, cl)
}

// triggerData returns the data of triggers and pressure plates, or nil for other objects.
func (obj nsObj) triggerData() *server.TriggerUpdateData {
	u := obj.SObj()
	if u.HasUpdate("TriggerUpdate") || u.HasUpdate("ToggleUpdate") {
		return u.UpdateDataTrigger()
	}
	return nil
}

func classesOf(mask uint32) []class.Class {
	var out []class.Class
	for _, c := range object.Class(mask).Split() {
		out = append(out, class.Class(c.String()))
	}
	return out
}

func classMask(list []class.Class) uint32 {
	var mask object.Class
	for _, c := range list {
		if cls, err := object.ParseClass(string(c)); err == nil {
			mask |= cls
		}
	}
	return uint32(mask)
}

func (obj nsObj) TriggerClasses() (allow, ignore []class.Class) {
	d := obj.triggerData()
	if d == nil {
		return nil, nil
	}
	return classesOf(d.AllowClass), classesOf(d.IgnoreClass)
}

func (obj nsObj) SetTriggerClasses(allow, ignore []class.Class) {
	d := obj.triggerData()
	if d == nil {
		return
	}
	d.AllowClass = classMask(allow)
	d.IgnoreClass = classMask(ignore)
	missiles := d.AllowClass&uint32(object.ClassMissile) != 0
	if missiles != obj.SObj().TriggerMissiles() {
		obj.SObj().SetExt().TriggerMissiles = missiles
	}
}

// deathBallLifetime returns the lifetime of a DeathBall (the ball of Force of Nature) in frames after creation.
func (obj nsObj) deathBallLifetime() uint32 {
	if v := obj.SObj().LifetimeOverride(); v != 0 {
		return v
	}
	return obj.s.SecToFrames(4)
}

func (obj nsObj) Lifetime() ns4.Duration {
	u := obj.SObj()
	s := obj.s
	age := s.Frame() - u.Field32 // creation frame
	left := func(life uint32) ns4.Duration {
		if life == math.MaxUint32 {
			return ns4.Infinite()
		} else if life <= age {
			return ns4.Frames(0)
		}
		return ns4.Frames(int(life - age))
	}
	switch {
	case u.HasUpdate("LifetimeUpdate"):
		return left(u.UpdateDataLifetime().Lifetime)
	case int(u.TypeInd) == s.Types.DeathBallID():
		return left(obj.deathBallLifetime())
	case u.ObjFlags&0x400000 != 0: // in the decay list, see DeleteAfter
		if u.Field34 <= s.Frame() {
			return ns4.Frames(0)
		}
		return ns4.Frames(int(u.Field34 - s.Frame()))
	}
	return ns4.Infinite()
}

func (obj nsObj) SetLifetime(dt ns4.Duration) {
	u := obj.SObj()
	s := obj.s
	age := s.Frame() - u.Field32
	life := uint32(math.MaxUint32)
	if !dt.IsInfinite() {
		life = age + uint32(max(s.AsFrames(dt), 1))
	}
	switch {
	case u.HasUpdate("LifetimeUpdate"):
		u.UpdateDataLifetime().Lifetime = life
	case int(u.TypeInd) == s.Types.DeathBallID():
		u.SetExt().Lifetime = life
	case dt.IsInfinite():
		legacy.ObjectDecayCancel(u)
	default:
		s.DeleteAfter(u, life-age)
	}
}

func (obj nsObj) DeleteAfter(dt ns4.Duration) {
	if dt.IsInfinite() {
		return // never (AsFrames gives 0 frames, which deleted it at once)
	}
	s := obj.s
	s.DeleteAfter(obj.Object.SObj(), uint32(s.AsFrames(dt)))
}

func (obj nsObj) IsAttackedBy(obj2 ns4.Obj) bool {
	if obj2 == nil {
		return false
	}
	return obj.Object.IsAttackedBy(server.ToObject(obj2.(server.Obj)))
}

func (obj nsObj) HasItem(item ns4.Obj) bool {
	if item == nil {
		return false
	}
	return obj.Object.HasItem(server.ToObject(item.(server.Obj)))
}

func (obj nsObj) HasEquipment(item ns4.Obj) bool {
	if !obj.HasItem(item) {
		return false
	}
	return item.Flags().Has(object.FlagEquipped)
}

func (obj nsObj) GetLastItem() ns4.Obj {
	it := obj.FirstItem()
	return obj.s.noxScriptP().toObj(it)
}

func (obj nsObj) GetPreviousItem() ns4.Obj {
	it := obj.NextItem()
	return obj.s.noxScriptP().toObj(it)
}

func (obj nsObj) Items(conditions ...ns4.ObjCond) []ns4.Obj {
	filter := ns4.AND(conditions)
	var out []ns4.Obj
	for it := obj.FirstItem(); it != nil; it = it.NextItem() {
		v := obj.s.noxScriptP().toObj(it)
		if filter.Matches(v) {
			out = append(out, v)
		}
	}
	return out
}

func (obj nsObj) Equipment(conditions ...ns4.ObjCond) []ns4.Obj {
	conditions = append(conditions, ns4.HasObjFlags(object.FlagEquipped))
	return obj.Items(conditions...)
}

func (obj nsObj) Weight() int {
	return int(obj.Object.Weight)
}

func (obj nsObj) InventoryWeight() int {
	total := 0
	for it := obj.FirstItem(); it != nil; it = it.NextItem() {
		total += int(it.Weight)
	}
	return total
}

func (obj nsObj) CarryMaxWeight() int {
	return int(obj.CarryCapacity) * 2
}

func (obj nsObj) SetCarryMaxWeight(v int) {
	c := v / 2
	if v%2 != 0 {
		c++
	}
	if c < 0 || c > math.MaxUint16 {
		c = math.MaxUint16
	}
	obj.CarryCapacity = uint16(c)
}

type nsObjInItems struct {
	obj     nsObj
	filters []ns4.ObjCond
}

func (s nsObjInItems) FindObjects(fnc func(it ns4.Obj) bool, conditions ...ns4.ObjCond) int {
	filter := ns4.AND(conditions)
	filter = append(filter, s.filters...)
	cnt := 0
	for it := s.obj.FirstItem(); it != nil; it = it.NextItem() {
		v := s.obj.s.noxScriptP().toObj(it)
		if !filter.Matches(v) {
			continue
		}
		cnt++
		if fnc != nil && !fnc(v) {
			break
		}
	}
	return cnt
}

func (obj nsObj) InItems() ns4.ObjSearcher {
	return nsObjInItems{obj: obj}
}

func (obj nsObj) InEquipment() ns4.ObjSearcher {
	return nsObjInItems{obj: obj, filters: []ns4.ObjCond{
		ns4.HasObjFlags(object.FlagEquipped),
	}}
}

func (obj nsObj) FindItems(fnc func(it ns4.Obj) bool, conditions ...ns4.ObjCond) int {
	return obj.InItems().FindObjects(fnc, conditions...)
}

func (obj nsObj) GetHolder() ns4.Obj {
	obj2 := obj.InvHolder
	return obj.s.noxScriptP().toObj(obj2)
}

func (obj nsObj) Pickup(item ns4.Obj) bool {
	if item == nil {
		return false
	}
	return obj.Object.DoPickup(server.ToObject(item.(server.Obj)))
}

func (obj nsObj) SetGeneratorMonsters(max int, types ...ns4.ObjTypeName) bool {
	if !obj.Class().Has(object.ClassMonsterGenerator) || max < 0 || max > 255 || len(types) == 0 || len(types) > 4 {
		return false
	}
	// every difficulty level has its own template monsters (like a generator from a map): they are freed with it
	var rows [3][4]*server.Object
	free := func() {
		for _, row := range rows {
			for _, m := range row {
				if m != nil {
					obj.s.Objs.FreeObject(m)
				}
			}
		}
	}
	for row := range rows {
		for i, t := range types {
			m := obj.s.NewObjectByTypeID(string(t))
			if m == nil || !m.Class().Has(object.ClassMonster) {
				if m != nil {
					obj.s.Objs.FreeObject(m)
				}
				free()
				return false
			}
			rows[row][i] = m
		}
	}
	ud := obj.UpdateDataMonsterGen()
	for i, old := range ud.Field0 {
		if old != nil {
			obj.s.Objs.FreeObject(old)
		}
		ud.Field0[i] = rows[i/4][i%4]
	}
	ud.Max = byte(max)
	// Field83 selects a maximum from gamedata.bin for each level, which the generator's init (when it's added
	// to the world, after CreateObject) and the start of a quest level would set instead: none of them
	ud.Field83 = [3]byte{0xff, 0xff, 0xff}
	return true
}

func (obj nsObj) Drop(item ns4.Obj) bool {
	if item == nil {
		return false
	}
	return obj.DoDrop(server.ToObject(item.(server.Obj)))
}

func (obj nsObj) Equip(item ns4.Obj) bool {
	if item == nil {
		return false
	}
	return obj.Object.Equip(server.ToObject(item.(server.Obj)))
}

func (obj nsObj) Unequip(item ns4.Obj) bool {
	if item == nil {
		return false
	}
	return obj.Object.Unequip(server.ToObject(item.(server.Obj)))
}

func (obj nsObj) CreateMover(wp ns4.WaypointObj, speed float32) ns4.Obj {
	mv := obj.Object.CreateMover(wp, speed)
	return obj.s.noxScriptP().toObj(mv)
}

func (obj nsObj) GetElevatorStatus() int {
	if !obj.Class().Has(object.ClassElevator) {
		return -1
	}
	return int(obj.UpdateDataElevator().Field_3)
}

func (obj nsObj) SetSightRange(val float32) {
	if obj.Class().Has(object.ClassMonster) {
		ud := obj.UpdateDataMonster()
		ud.SightRange = val
	}
}

func (obj nsObj) AggressionLevel(level float32) {
	obj.Object.SetAggression(level)
}

func (obj nsObj) SetRoamFlag(flags int) {
	legacy.ScriptRoamFlag(obj.SObj(), byte(flags))
}

func (obj nsObj) RetreatLevel(percent float32) {
	obj.Object.SetRetreatLevel(percent)
}

func (obj nsObj) ResumeLevel(percent float32) {
	obj.Object.SetRegroupLevel(percent)
}

func (obj nsObj) SetShopText(id ns4.StringID) {
	obj.Object.SetShopText(strman.ID(id))
}

func (obj nsObj) OnEvent(ev ns4.ObjectEvent, fnc ns4.Func) {
	if !obj.Class().Has(object.ClassMonster) {
		return
	}
	ind := int32(obj.s.noxScript.AsFuncIndex(fmt.Sprintf("OnEvent_%d", int(ev)), fnc))
	ud := obj.UpdateDataMonster()
	switch ev {
	case ns4.EventEnemySighted: // Enemy sighted
		ud.ScriptEnemySighted.Func = ind
	case ns4.EventLookingForEnemy: // Looking for enemy
		ud.ScriptLookingForEnemy.Func = ind
	case ns4.EventDeath: // Death
		ud.ScriptDeath.Func = ind
	case ns4.EventChangeFocus: // Change focus
		ud.ScriptChangeFocus.Func = ind
	case ns4.EventIsHit: // Is hit
		ud.ScriptIsHit.Func = ind
	case ns4.EventRetreat: // Retreat
		ud.ScriptRetreat.Func = ind
	case ns4.EventCollision: // Collision
		ud.ScriptCollision.Func = ind
	case ns4.EventEnemyHeard: // Enemy heard
		ud.ScriptHearEnemy.Func = ind
	case ns4.EventEndOfWaypoint: // End of waypoint
		ud.ScriptEndOfWaypoint.Func = ind
	case ns4.EventLostEnemy: // Lost sight of enemy
		ud.ScriptLostEnemy.Func = ind
	}
}

type nsObjGroup struct {
	s *Server
	g *server.MapGroup
}

func (g nsObjGroup) ScriptID() int {
	return int(g.g.Index())
}

func (g nsObjGroup) ObjGroupScriptID() int {
	return int(g.g.Index())
}

func (g nsObjGroup) Name() string {
	return g.g.ID()
}

func (g nsObjGroup) Enable(enable bool) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Enable(enable)
		return true
	})
}

func (g nsObjGroup) Toggle() bool {
	one := false
	g.EachObject(true, func(obj ns4.Obj) bool {
		if obj.Toggle() {
			one = true
		}
		return true
	})
	return one
}

func (g nsObjGroup) HasOwner(owner ns4.Obj) bool {
	ok := true
	g.EachObject(false, func(obj ns4.Obj) bool {
		ok = obj.HasOwner(owner)
		return ok
	})
	return ok
}

func (g nsObjGroup) HasOwnerIn(owners ns4.ObjGroup) bool {
	// TODO: not exactly sure if the logic was like this
	ok := true
	g.EachObject(false, func(obj ns4.Obj) bool {
		ok = obj.HasOwnerIn(owners)
		return ok
	})
	return ok
}

func (g nsObjGroup) SetOwner(owner ns4.Obj) {
	g.EachObject(false, func(obj ns4.Obj) bool {
		obj.SetOwner(owner)
		return true
	})
}

func (g nsObjGroup) SetOwners(owners ns4.ObjGroup) {
	g.EachObject(false, func(obj ns4.Obj) bool {
		obj.SetOwners(owners)
		return true
	})
}

func (g nsObjGroup) Pause(dt ns4.Duration) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Pause(dt)
		return true
	})
}

func (g nsObjGroup) Move(wp ns4.WaypointObj) {
	if wp == nil {
		return
	}
	g.EachObject(false, func(obj ns4.Obj) bool {
		obj.Move(wp)
		return true
	})
}

func (g nsObjGroup) WalkTo(p types.Pointf) {
	g.EachObject(false, func(obj ns4.Obj) bool {
		obj.WalkTo(p)
		return true
	})
}

func (g nsObjGroup) LookAtDirection(dir ns4.Direction) {
	g.EachObject(false, func(obj ns4.Obj) bool {
		obj.LookAtDirection(dir)
		return true
	})
}

func (g nsObjGroup) Damage(source ns4.Obj, amount int, typ damage.Type) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Damage(source, amount, typ)
		return true
	})
}

func (g nsObjGroup) Delete() {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Delete()
		return true
	})
}

func (g nsObjGroup) Idle() {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Idle()
		return true
	})
}

func (g nsObjGroup) Wander() {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Wander()
		return true
	})
}

func (g nsObjGroup) Hunt() {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Hunt()
		return true
	})
}

func (g nsObjGroup) Follow(targ ns4.Positioner) {
	if targ == nil {
		return
	}
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Follow(targ)
		return true
	})
}

func (g nsObjGroup) Guard(p1, p2 types.Pointf, distance float32) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Guard(p1, p2, distance)
		return true
	})
}

func (g nsObjGroup) Flee(target ns4.Positioner, dt ns4.Duration) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Flee(target, dt)
		return true
	})
}

func (g nsObjGroup) HitMelee(pos types.Pointf) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.HitMelee(pos)
		return true
	})
}

func (g nsObjGroup) HitRanged(pos types.Pointf) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.HitRanged(pos)
		return true
	})
}

func (g nsObjGroup) Attack(targ ns4.Positioner) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Attack(targ)
		return true
	})
}

func (g nsObjGroup) ZombieStayDown() {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.ZombieStayDown()
		return true
	})
}

func (g nsObjGroup) RaiseZombie() {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.RaiseZombie()
		return true
	})
}

func (g nsObjGroup) CreateMover(wp ns4.WaypointObj, speed float32) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.CreateMover(wp, speed)
		return true
	})
}

func (g nsObjGroup) AggressionLevel(level float32) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.AggressionLevel(level)
		return true
	})
}

func (g nsObjGroup) SetRoamFlag(flags int) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.SetRoamFlag(flags)
		return true
	})
}

func (g nsObjGroup) RetreatLevel(percent float32) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.RetreatLevel(percent)
		return true
	})
}

func (g nsObjGroup) ResumeLevel(percent float32) {
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.ResumeLevel(percent)
		return true
	})
}

func (g nsObjGroup) Enchant(enchant enchant.Enchant, dt ns4.Duration) {
	if _, ok := server.ParseEnchant(string(enchant)); !ok {
		return
	}
	g.EachObject(true, func(obj ns4.Obj) bool {
		obj.Enchant(enchant, dt)
		return true
	})
}

func (g nsObjGroup) EachObject(recursive bool, fnc func(obj ns4.Obj) bool) {
	s := g.s.noxScriptP()
	if recursive {
		server.EachObjectRecursive(g.s.Server, g.g, func(obj *server.Object) bool {
			return fnc(s.toObj(obj))
		})
	} else {
		server.EachObject(g.s.Server, g.g, func(obj *server.Object) bool {
			return fnc(s.toObj(obj))
		})
	}
}

func (g nsObjGroup) AllObjects() ns4.Objects {
	var out ns4.Objects
	g.EachObject(true, func(obj ns4.Obj) bool {
		out = append(out, obj)
		return true
	})
	return out
}

func (g nsObjGroup) FindObjects(fnc func(it ns4.Obj) bool, conditions ...ns4.ObjCond) int {
	filter := ns4.AND(conditions)
	cnt := 0
	g.EachObject(true, func(obj ns4.Obj) bool {
		if !filter.Matches(obj) {
			return true // continue
		}
		cnt++
		if fnc != nil && !fnc(obj) {
			return false // break
		}
		return true
	})
	return cnt
}
