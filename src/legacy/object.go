package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"

	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_unitMonsterInit_4F0040             func(obj *server.Object)
	Nox_xxx_checkSummonedCreaturesLimit_500D70 func(u *server.Object, ind int) bool
	Nox_xxx_unitDoSummonAt_5016C0              func(typID int, pos types.Pointf, owner *server.Object, dir server.Dir16) *server.Object
	Sub_4E71F0                                 func(obj *server.Object)
	Nox_bomberDead_54A150                      func(obj *server.Object) int
	Nox_xxx_dieGlyph_54DF30                    func(obj *server.Object)
	Nox_xxx_collideGlyph_4E9A00                func(obj, obj2 *server.Object)
	Nox_xxx_playerSetState_4FA020              func(a1 *server.Object, a2 server.PlayerState) bool
)

type nox_object_t = server.Object

func asObjectC(p *server.Object) *nox_object_t {
	return (*nox_object_t)(p.CObj())
}

func AsObjectP(p unsafe.Pointer) *server.Object {
	return (*server.Object)(p)
}

func asObjectS(p *nox_object_t) *server.Object {
	return (*server.Object)(unsafe.Pointer(p))
}

func ToObjS(p *nox_object_t) server.Obj {
	return objectAsInterface(asObjectS(p))
}

func objectAsInterface(p *server.Object) server.Obj {
	if p == nil {
		return nil
	}
	return p
}

func nox_server_getFirstObject_4DA790() *nox_object_t {
	return asObjectC(GetServer().S().Objs.First())
}

func nox_server_getNextObject_4DA7A0(cobj *nox_object_t) *nox_object_t {
	return asObjectC(asObjectS(cobj).Next())
}

func nox_xxx_servFinalizeDelObject_4DADE0(cobj *nox_object_t) {
	GetServer().ObjectDeleteLast(asObjectS(cobj))
}

func nox_xxx_delayedDeleteObject_4E5CC0(obj *nox_object_t) {
	GetServer().DelayedDelete(asObjectS(obj))
}

func nox_xxx_unitSetOwner_4EC290(obj1, obj2 *nox_object_t) {
	GetServer().S().ObjSetOwner(asObjectS(obj1), asObjectS(obj2))
}

func nox_xxx_unitClearOwner_4EC300(obj *nox_object_t) {
	GetServer().S().ObjClearOwner(asObjectS(obj))
}

func nox_xxx_creatureIsMonitored_500CC0(obj1, obj2 *nox_object_t) int {
	return bool2int(server.Nox_xxx_creatureIsMonitored_500CC0(asObjectS(obj1), asObjectS(obj2)))
}

func nox_xxx_netMarkMinimapObject_417190(a1 int, obj *nox_object_t, a3 uint32) {
	GetServer().S().Players.Nox_xxx_netMarkMinimapObject_417190(ntype.PlayerInd(a1), asObjectS(obj), a3)
}

func nox_xxx_netUnmarkMinimapObj_417300(a1 int, obj *nox_object_t, a3 uint32) {
	GetServer().S().Players.Nox_xxx_netUnmarkMinimapObj_417300(ntype.PlayerInd(a1), asObjectS(obj), a3)
}

func nox_xxx_monsterMarkUpdate_4E8020(obj *nox_object_t) {
	asObjectS(obj).Nox_xxx_monsterMarkUpdate_4E8020()
}

func nox_xxx_unitIsHostileMimic_4E7F90(obj1, obj2 *nox_object_t) int {
	return bool2int(GetServer().S().IsHostileMimicXxx(asObjectS(obj1), asObjectS(obj2)))
}

func AsPointf(p unsafe.Pointer) types.Pointf {
	cp := (*types.Pointf)(p)
	return types.Pointf{
		X: float32(cp.X),
		Y: float32(cp.Y),
	}
}
func AsPoint(p unsafe.Pointer) image.Point {
	cp := (*[2]int32)(p)
	return image.Point{
		X: int(cp[0]),
		Y: int(cp[1]),
	}
}

func nox_xxx_objectFreeMem_4E38A0(a1p *nox_object_t) int {
	return GetServer().S().Objs.FreeObject(asObjectS(a1p))
}

func nox_xxx_findParentChainPlayer_4EC580(obj *nox_object_t) *nox_object_t {
	return asObjectC(asObjectS(obj).FindOwnerChainPlayer())
}

func nox_xxx_unitHasThatParent_4EC4F0(obj, owner *nox_object_t) int {
	return bool2int(asObjectS(obj).HasOwner(asObjectS(owner)))
}

func nox_xxx_unitIsEnemyTo_5330C0(a, b *nox_object_t) int {
	return bool2int(GetServer().S().IsEnemyTo(asObjectS(a), asObjectS(b)))
}

func nox_xxx_checkSummonedCreaturesLimit_500D70(obj *nox_object_t, ind int) bool {
	return bool(Nox_xxx_checkSummonedCreaturesLimit_500D70(asObjectS(obj), ind))
}

func nox_xxx_unitSetXStatus_4E4800(a1 *nox_object_t, a2 uint32) {
	asObjectS(a1).SetXStatus(a2)
}

func nox_xxx_unitUnsetXStatus_4E4780(a1 *nox_object_t, a2 uint32) {
	asObjectS(a1).UnsetXStatus(a2)
}

func nox_xxx_playerSetState_4FA020(a1 *nox_object_t, a2 int) int {
	return bool2int(Nox_xxx_playerSetState_4FA020(asObjectS(a1), server.PlayerState(a2)))
}

func nox_xxx_weaponInventoryEquipFlags_415820(obj *nox_object_t) int {
	return int(GetServer().S().Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(asObjectS(obj)))
}

func nox_xxx_unitArmorInventoryEquipFlags_415C70(obj *nox_object_t) int {
	return int(GetServer().S().Armor.Nox_xxx_unitArmorInventoryEquipFlags_415C70(asObjectS(obj)))
}

func nox_xxx_ammoCheck_415880(a1 int) int {
	return int(GetServer().S().Weapons.Nox_xxx_ammoCheck_415880(a1))
}

func sub_415840(a1 int) int {
	return int(GetServer().S().Weapons.Sub_415840(uint32(a1)))
}

func Nox_server_getObjectFromNetCode_4ECCB0(a1 int) *server.Object {
	return objectLookupByNetCode(uint32(a1))
}
func Nox_xxx_monsterRemoveMonitors_4E7B60(a1 *server.Object, a2 *server.Object) {
	stateRemoveMonitors(a1, a2)
}
func Sub_4ED0C0(a1 *server.Object, a2 *server.Object) {
	inventoryRemove(a1, a2)
}
func Nox_xxx_playerCancelSpells_4FEAE0(a1 *server.Object) {
	spellLifeCancelPlayer(a1)
}
func Sub_50E210(a1 *server.Object) {
	spawnPolicyGlyphRelease(a1)
}
func Sub_506740(a1 *server.Object) {
	voteRemovePlayer(a1)
}
func Nox_xxx_unitTransferSlaves_4EC4B0(a1 *server.Object) {
	controlTransferChildren(a1)
}
func Nox_xxx_decay_5116F0(a1 *server.Object) {
	motionDecayRemove(a1)
}
func Nox_xxx_netReportDestroyObject_5289D0(a1 *server.Object) {
	visibilityDestroyReport(a1)
}
func Nox_xxx_unit_511810(a1 *server.Object) {
	motionDeactivate(a1)
}
func Nox_xxx_unitRemoveChild_4EC470(a1 *server.Object) {
	controlRemoveChildren(a1)
}
func Sub_4ECFA0(a1 *server.Object) {
	netCodeCacheInvalidate(a1)
}
func Sub_511DE0(a1 *server.Object) {
	monsterCacheRemove(a1)
}
func Sub_528990(a1 *server.Object) {
	visibilityGlobalRemove(a1)
}
func Nox_xxx_unitNewAddShadow_4DA9A0(a1 *server.Object) {
	sessionShadowAdd(a1)
}
func Nox_xxx_respawnAdd_4EC5E0(a1 *server.Object) {
	itemRespawnAdd(a1)
}
func Sub_5117F0(a1 *server.Object) {
	motionActivate(a1)
}
func Nox_xxx_action_4DA9F0(a1 *server.Object) {
	sessionShadowRemove(a1)
}
func Nox_xxx_unitPostCreateNotify_4E7F10(a1 *server.Object) {
	statePostCreate(a1)
}
func Nox_xxx_buffApplyTo_4FF380(a1 *server.Object, a2 server.EnchantID, dur int, power int) {
	spellLifeApplyBuff(a1, int32(a2), int16(dur), int8(power))
}
func Nox_xxx_spellBuffOff_4FF5B0(a1 *server.Object, a2 server.EnchantID) {
	spellLifeBuffOff(a1, int32(a2))
}
func Nox_xxx_unitRaise_4E46F0(a1 *server.Object, a2 float32) {
	stateRaise(a1, float32(a2))
}
func Nox_xxx_objectSetOff_4E7600(a1 *server.Object) {
	stateOff(a1)
}
func Nox_xxx_objectSetOn_4E75B0(a1 *server.Object) {
	stateOn(a1)
}
func Nox_xxx_drop_4ED790(a1 *server.Object, a2 *server.Object, pos types.Pointf) int {
	cpos, free := alloc.New(types.Pointf{})
	defer free()
	*cpos = pos
	return inventoryDrop(a1, a2, cpos)
}
func Nox_xxx_dropAllItems_4EDA40(a1 *server.Object) {
	inventoryDropAll(a1)
}

func Get_nox_objectDropAudEvent_4EE2F0() unsafe.Pointer {
	return itemIdentityKey(itemIDAudEventDrop)
}
func Get_nox_xxx_XFerDefault_4F49A0() unsafe.Pointer {
	return xferIdentityKey(xferIDDefault)
}
func Get_nox_xxx_updateHarpoon_54F380() unsafe.Pointer {
	return updateIdentityKey(updateIDHarpoon)
}
func Get_nox_xxx_updatePixie_53CD20() unsafe.Pointer {
	return updateIdentityKey(updateIDPixie)
}
func Nox_object_getGold_4FA6D0(obj *server.Object) int {
	return int(int32(resourceObjectGold(obj)))
}
func Nox_object_setGold_4FA620(obj *server.Object, v int) {
	resourceSetGold(obj, int32(v))
}
func Nox_xxx_script_forcedialog_548CD0(obj, obj2 *server.Object) {
	unitForceDialogue(obj, obj2)
}
func Sub_4E39F0_obj_db(obj *server.Object) string {
	return alloc.GoString16(unitNPCName(obj))
}
func Nox_xxx_scriptDialog_548D30(obj *server.Object, a2 byte) {
	unitFinishDialogue(obj, a2)
}
func Nox_xxx_mobSetFightTarg_515D30(obj, targ *server.Object) {
	monsterControlFight(obj, targ)
}
func Nox_server_scriptFleeFrom_515F70(obj, targ *server.Object, df int) {
	if targ == nil {
		return
	}
	monsterControlFlee(obj, targ.PosVec, uint32(df))
}

// FleeFromPos makes a monster flee from a position for a number of frames.
func FleeFromPos(obj *server.Object, pos types.Pointf, df int) {
	monsterControlFlee(obj, pos, uint32(df))
}
func Nox_xxx_monsterGoPatrol_515680(obj *server.Object, p1, p2 types.Pointf, dist float32) {
	monsterControlPatrol(obj, p1, p2, dist)
}
func Nox_xxx_monsterActionMelee_515A30(obj *server.Object, pos types.Pointf) {
	monsterControlMelee(obj, &pos)
}
func Nox_xxx_monsterMissileAttack_515B80(obj *server.Object, pos types.Pointf) {
	monsterControlMissile(obj, &pos)
}

func Sub_516090(obj *server.Object, df int) {
	monsterControlWait(obj, uint32(df))
}

func Nox_xxx_monsterCanCast_534300(obj *server.Object) bool {
	return monsterCanCast(obj)
}

func Nox_xxx_playerTryEquip_4F2F70(obj, item *server.Object) bool {
	return equipmentTryEquip(obj, item) != 0
}

func Nox_xxx_playerTryDequip_4F2FB0(obj, item *server.Object) bool {
	return equipmentTryDequip(obj, item) != 0
}

func Nox_xxx_inventoryPutImpl_4F3070(obj, item *server.Object, a3 int) {
	inventoryInsert(obj, item, a3)
}

func Nox_xxx_orderUnit_533900(owner, obj *server.Object, order uint32) {
	monsterOrder(owner, obj, int(order))
}

func Sub_4E9A30(a1, a2 *server.Object) bool {
	return projectileTrapEligible(a1, a2)
}

func Nox_xxx_unitsHaveSameTeam_4EC520(a1, a2 *server.Object) bool {
	return itemOwnerSameTeam(a1, a2)
}

func Nox_xxx_mapPushUnitsAround_52E040(pos types.Pointf, a2, a3, a4 float32, a5 *server.Object, a6, a7 int) {
	spellEffectPushAround(pos, a2, a3, a4, a5, unsafe.Pointer(uintptr(a6)), uintptr(a7))
}
