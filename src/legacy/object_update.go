package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/player"

	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_updatePlayer_4F8100         func(up *server.Object)
	Nox_xxx_updatePixie_53CD20          func(cobj *server.Object)
	Nox_xxx_updatePlayerObserver_4E62F0 func(a1p *server.Object)
	Nox_xxx_updateProjectile_53AC10     func(a1 *server.Object)
	Nox_xxx_updateDeathBall_53D080      func(a1 *server.Object)
	Nox_xxx___mkgmtime_538280           func(a1 *server.Object)
)

func init() {
	server.RegisterObjectUpdateGo("PlayerUpdate", updateIdentityKey(updateIDPlayer), func(u *server.Object) { Nox_xxx_updatePlayer_4F8100(u) }, unsafe.Sizeof(server.PlayerUpdateData{}))
	server.RegisterObjectUpdateGo("ProjectileUpdate", updateIdentityKey(updateIDProjectile), func(u *server.Object) { Nox_xxx_updateProjectile_53AC10(u) }, 0)
	server.RegisterObjectUpdateGo("SpellProjectileUpdate", updateIdentityKey(updateIDSpellProjectile), func(u *server.Object) { temporarySpellFly(u) }, unsafe.Sizeof(server.SpellProjectileUpdateData{}))
	server.RegisterObjectUpdateGo("AntiSpellProjectileUpdate", updateIdentityKey(updateIDAntiSpellProjectile), func(u *server.Object) { temporaryAntiSpell(u) }, 28)
	server.RegisterObjectUpdateGo("DoorUpdate", updateIdentityKey(updateIDDoor), func(u *server.Object) { worldDoor(u) }, 52)
	server.RegisterObjectUpdateGo("SparkUpdate", updateIdentityKey(updateIDSpark), func(u *server.Object) { temporarySparkUpdate(u) }, 16)
	server.RegisterObjectUpdateGo("ProjectileTrailUpdate", updateIdentityKey(updateIDProjectileTrail), func(u *server.Object) { temporaryTrail(u) }, 0)
	server.RegisterObjectUpdateGo("PushUpdate", updateIdentityKey(updateIDPush), func(u *server.Object) { worldPush(u) }, 12)
	server.RegisterObjectUpdateGo("TriggerUpdate", updateIdentityKey(updateIDTrigger), func(u *server.Object) { worldTrigger(u) }, 60)
	server.RegisterObjectUpdateGo("ToggleUpdate", updateIdentityKey(updateIDToggle), func(u *server.Object) { worldToggle(u) }, 60)
	server.RegisterObjectUpdateGo("MonsterUpdate", updateIdentityKey(updateIDMonster), func(u *server.Object) { Nox_xxx_unitUpdateMonster_50A5C0(u) }, unsafe.Sizeof(server.MonsterUpdateData{}))
	server.RegisterObjectUpdateGo("LoopAndDamageUpdate", updateIdentityKey(updateIDLoopAndDamage), func(u *server.Object) { worldEnabledCollision(u) }, 16)
	server.RegisterObjectUpdateGo("ElevatorUpdate", updateIdentityKey(updateIDElevator), func(u *server.Object) { worldElevator(u) }, 20)
	server.RegisterObjectUpdateGo("ElevatorShaftUpdate", updateIdentityKey(updateIDElevatorShaft), func(u *server.Object) { worldShaft(u) }, 16)
	server.RegisterObjectUpdateGo("PhantomPlayerUpdate", updateIdentityKey(updateIDPhantomPlayer), func(u *server.Object) { worldPhantom(u) }, 0)
	server.RegisterObjectUpdateGo("ObeliskUpdate", updateIdentityKey(updateIDObelisk), func(u *server.Object) { objectiveObelisk(u) }, unsafe.Sizeof(server.ObeliskUpdateData{}))
	server.RegisterObjectUpdateGo("LifetimeUpdate", updateIdentityKey(updateIDLifetime), func(u *server.Object) { temporaryLifetime(u) }, 4)
	server.RegisterObjectUpdateGo("MagicMissileUpdate", updateIdentityKey(updateIDMagicMissile), func(u *server.Object) { temporaryMagicMissile(u) }, 28)
	server.RegisterObjectUpdateGo("PixieUpdate", updateIdentityKey(updateIDPixie), func(u *server.Object) { Nox_xxx_updatePixie_53CD20(u) }, 28)
	server.RegisterObjectUpdateGo("SkullUpdate", updateIdentityKey(updateIDSkull), func(u *server.Object) { motionTrapUpdate(u) }, 52)
	server.RegisterObjectUpdateGo("PentagramUpdate", updateIdentityKey(updateIDPentagram), func(u *server.Object) { worldTeleport(u) }, 24)
	server.RegisterObjectUpdateGo("InvisiblePentagramUpdate", updateIdentityKey(updateIDInvisiblePentagram), func(u *server.Object) { worldInvisibleTeleport(u) }, 24)
	server.RegisterObjectUpdateGo("SwitchUpdate", updateIdentityKey(updateIDSwitch), func(u *server.Object) { worldSwitch(u) }, 0)
	server.RegisterObjectUpdateGo("BlowUpdate", updateIdentityKey(updateIDBlow), func(u *server.Object) { worldBlow(u) }, 0)
	server.RegisterObjectUpdateGo("MoverUpdate", updateIdentityKey(updateIDMover), func(u *server.Object) { motionMover(u) }, 36)
	server.RegisterObjectUpdateGo("BlackPowderBarrelUpdate", updateIdentityKey(updateIDBlackPowderBarrel), func(u *server.Object) { temporaryPowderBarrel(u) }, 0)
	server.RegisterObjectUpdateGo("OneSecondDieUpdate", updateIdentityKey(updateIDOneSecondDie), func(u *server.Object) { temporaryOneSecond(u) }, 0)
	server.RegisterObjectUpdateGo("WaterBarrelUpdate", updateIdentityKey(updateIDWaterBarrel), func(u *server.Object) { temporaryWaterBarrel(u) }, 0)
	server.RegisterObjectUpdateGo("SelfDestructUpdate", updateIdentityKey(updateIDSelfDestruct), func(u *server.Object) { temporarySelfDestruct(u) }, 0)
	server.RegisterObjectUpdateGo("BlackPowderBurnUpdate", updateIdentityKey(updateIDBlackPowderBurn), func(u *server.Object) { temporaryPowderBurn(u) }, 0)
	server.RegisterObjectUpdateGo("DeathBallUpdate", updateIdentityKey(updateIDDeathBall), func(u *server.Object) { Nox_xxx_updateDeathBall_53D080(u) }, 0)
	server.RegisterObjectUpdateGo("DeathBallFragmentUpdate", updateIdentityKey(updateIDDeathBallFragment), func(u *server.Object) { temporaryDeathFragment(u) }, 0)
	server.RegisterObjectUpdateGo("MoonglowUpdate", updateIdentityKey(updateIDMoonglow), func(u *server.Object) { temporaryMoonglow(u) }, 0)
	server.RegisterObjectUpdateGo("SentryGlobeUpdate", updateIdentityKey(updateIDSentryGlobe), func(u *server.Object) { motionSentryUpdate(u) }, 12)
	server.RegisterObjectUpdateGo("TelekinesisUpdate", updateIdentityKey(updateIDTelekinesis), func(u *server.Object) { temporaryTelekinesis(u) }, 0)
	server.RegisterObjectUpdateGo("FistUpdate", updateIdentityKey(updateIDFist), func(u *server.Object) { temporaryFist(u) }, 4)
	server.RegisterObjectUpdateGo("MeteorShowerUpdate", updateIdentityKey(updateIDMeteorShower), func(u *server.Object) { temporaryMeteorShower(u) }, 4)
	server.RegisterObjectUpdateGo("MeteorUpdate", updateIdentityKey(updateIDMeteor), func(u *server.Object) { temporaryMeteorExplode(u) }, 4)
	server.RegisterObjectUpdateGo("ToxicCloudUpdate", updateIdentityKey(updateIDToxicCloud), func(u *server.Object) { temporaryCloud(u, false) }, 4)
	server.RegisterObjectUpdateGo("SmallToxicCloudUpdate", updateIdentityKey(updateIDSmallToxicCloud), func(u *server.Object) { temporaryCloud(u, true) }, 4)
	server.RegisterObjectUpdateGo("ArachnaphobiaUpdate", updateIdentityKey(updateIDArachnaphobia), func(u *server.Object) { temporaryArachnaphobia(u) }, 0)
	server.RegisterObjectUpdateGo("ExpireUpdate", updateIdentityKey(updateIDExpire), func(u *server.Object) { temporaryExpire(u) }, 0)
	server.RegisterObjectUpdateGo("BreakUpdate", updateIdentityKey(updateIDBreak), func(u *server.Object) { temporaryBreak(u, false) }, 0)
	server.RegisterObjectUpdateGo("OpenUpdate", updateIdentityKey(updateIDOpen), func(u *server.Object) { temporaryBreak(u, true) }, 0)
	server.RegisterObjectUpdateGo("BreakAndRemoveUpdate", updateIdentityKey(updateIDBreakAndRemove), func(u *server.Object) { temporaryBreakRemove(u) }, 0)
	server.RegisterObjectUpdateGo("ChakramInMotionUpdate", updateIdentityKey(updateIDChakramInMotion), func(u *server.Object) { temporaryChakram(u) }, 28)
	server.RegisterObjectUpdateGo("FlagUpdate", updateIdentityKey(updateIDFlag), func(u *server.Object) { objectiveFlagUpdate(u) }, 12)
	server.RegisterObjectUpdateGo("TrapDoorUpdate", updateIdentityKey(updateIDTrapDoor), func(u *server.Object) { worldTrapDoor(u) }, 0)
	server.RegisterObjectUpdateGo("BallUpdate", updateIdentityKey(updateIDBall), func(u *server.Object) { objectiveBallUpdate(u) }, 32)
	server.RegisterObjectUpdateGo("CrownUpdate", updateIdentityKey(updateIDCrown), func(u *server.Object) { objectiveCrownUpdate(u) }, 12)
	server.RegisterObjectUpdateGo("UndeadKillerUpdate", updateIdentityKey(updateIDUndeadKiller), func(u *server.Object) { unitUndeadUpdate(u) }, 0)
	server.RegisterObjectUpdateGo("HarpoonUpdate", updateIdentityKey(updateIDHarpoon), func(u *server.Object) { Nox_xxx_updateHarpoon_54F380(u) }, 4)
	server.RegisterObjectUpdateGo("MonsterGeneratorUpdate", updateIdentityKey(updateIDMonsterGenerator), func(u *server.Object) { generatorUpdate(u) }, 164)

	server.RegisterObjectUpdateCallbackGo(updateIdentityKey(updateIDPlayerObserver), func(u *server.Object) { Nox_xxx_updatePlayerObserver_4E62F0(u) })
	server.RegisterObjectUpdateCallbackGo(updateIdentityKey(updateIDMkgmtime), func(u *server.Object) { Nox_xxx___mkgmtime_538280(u) })
	server.RegisterObjectUpdateCallbackGo(updateIdentityKey(updateIDPlayerMonsterBot), func(u *server.Object) { controlBotUpdate(u) })

	server.RegisterObjectUpdateParse("PushUpdate", resourceObjectParser("update", "push"))
	server.RegisterObjectUpdateParse("TriggerUpdate", resourceObjectParser("update", "trigger"))
	server.RegisterObjectUpdateParse("ToggleUpdate", resourceObjectParser("update", "trigger"))
	server.RegisterObjectUpdateParse("LoopAndDamageUpdate", resourceObjectParser("update", "triple"))
	server.RegisterObjectUpdateParse("LifetimeUpdate", resourceObjectParser("update", "lifetime"))
	server.RegisterObjectUpdateParse("SkullUpdate", resourceObjectParser("update", "skull"))
}

func nox_xxx_objectApplyForce_52DF80(vec *float32, obj *nox_object_t, force float32) {
	GetServer().ApplyForce(asObjectS(obj), AsPointf(unsafe.Pointer(vec)), float64(force))
}

func Get_nox_xxx___mkgmtime_538280() unsafe.Pointer {
	return updateIdentityKey(updateIDMkgmtime)
}

func Nox_server_doPlayersAutoRespawn_40A5F0() int {
	return int(int32(serverConfigRespawnGet()))
}
func Sub_4E4100() uint32 {
	return uint32(bool2int(questRuntimeRoom()))
}

func Nox_xxx_questCheckSecretArea_421C70(a1 *server.Object) {
	mapPolygonPlayer(a1)
}
func Nox_xxx_playerCanMove_4F9BC0(a1 *server.Object) int {
	return bool2int(controlCanMove(a1))
}
func Sub_4F9AB0(a1 *server.Object) int {
	return int(controlWalkWaypoint(a1))
}
func Nox_xxx_playerConfusedGetDirection_4F7A40(a1 *server.Object) server.Dir16 {
	return server.Dir16(controlConfusedDirection(a1))
}
func Nox_xxx_playerAttack_538960(a1 *server.Object) int {
	return attackPlayer(a1)
}
func Nox_xxx_playerRespawn_4F7EF0(a1 *server.Object) {
	controlRespawn(a1)
}
func Sub_4F9E10(a1 *server.Object) int {
	return int(controlFollowEnemy(a1))
}
func Sub_4F9A80(a1 *server.Object) int {
	return bool2int(controlHasWaypoint(a1))
}
func Nox_xxx_monsterTestBlockShield_533E70(a1 *server.Object) int {
	return int(uintptr(unsafe.Pointer(monsterShieldThreat(a1))))
}
func Nox_common_mapPlrActionToStateId_4FA2B0(a1 *server.Object) int {
	return int(controlActionState(a1))
}
func Nox_xxx_playerCanAttack_4F9C40(a1 *server.Object) int {
	return bool2int(controlCanAttack(a1))
}
func Nox_xxx_checkWinkFlags_4F7DF0(a1 *server.Object) int {
	return int(controlDropBall(a1))
}
func Nox_xxx_playerInputAttack_4F9C70(a1 *server.Object) {
	controlInputAttack(a1)
}
func Nox_xxx_playerSubStamina_4F7D30(a1 *server.Object, a2 int) int {
	return int(controlSubStamina(a1, int32(a2)))
}
func Nox_xxx_playerDoSchedSpell_4FB0E0(a1 *server.Object, a2 *server.Object) {
	controlScheduledSpell(a1, a2, false)
}
func Nox_xxx_playerDoSchedSpellQueue_4FB1D0(a1 *server.Object, a2 *server.Object) {
	controlScheduledSpell(a1, a2, true)
}
func Sub_4E7540(a1 *server.Object, a2 *server.Object) {
	stateRememberAttacker(a1, a2)
}
func Nox_xxx_playerCheckStrength_4F3180(a1 *server.Object, a2 *server.Object) bool {
	return equipmentCheckStrength(a1, a2)
}
func Nox_xxx_unitDamageClear_4EE5E0(a1 *server.Object, a2 int) {
	resourceDamage(a1, int32(a2))
}
func Nox_xxx_unitAdjustHP_4EE460(a1 *server.Object, a2 int) {
	resourceAdjustHP(a1, int32(a2))
}
func Sub_509CF0(a1 *byte, a2 player.Class, a3 uint32) int {
	return bool2int(matchRosterIdentityAllowed(GoStringP(unsafe.Pointer(a1)), byte(a2), a3))
}
func Sub_4D79C0(a1 *server.Object) {
	questRuntimeReconnect(a1)
}
func Sub_4D7480(a1 *server.Object) {
	questRuntimeGateReturn(a1)
}
