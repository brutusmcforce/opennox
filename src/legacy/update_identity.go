package legacy

import "unsafe"

type updateIdentitySlot struct{ marker byte }

var updateIdentitySlots [56]updateIdentitySlot

type updateIdentityID uint8

const (
	updateIDPlayer              updateIdentityID = 0
	updateIDProjectile          updateIdentityID = 1
	updateIDSpellProjectile     updateIdentityID = 2
	updateIDAntiSpellProjectile updateIdentityID = 3
	updateIDDoor                updateIdentityID = 4
	updateIDSpark               updateIdentityID = 5
	updateIDProjectileTrail     updateIdentityID = 6
	updateIDPush                updateIdentityID = 7
	updateIDTrigger             updateIdentityID = 8
	updateIDToggle              updateIdentityID = 9
	updateIDMonster             updateIdentityID = 10
	updateIDLoopAndDamage       updateIdentityID = 11
	updateIDElevator            updateIdentityID = 12
	updateIDElevatorShaft       updateIdentityID = 13
	updateIDPhantomPlayer       updateIdentityID = 14
	updateIDObelisk             updateIdentityID = 15
	updateIDLifetime            updateIdentityID = 16
	updateIDMagicMissile        updateIdentityID = 17
	updateIDPixie               updateIdentityID = 18
	updateIDSkull               updateIdentityID = 19
	updateIDPentagram           updateIdentityID = 20
	updateIDInvisiblePentagram  updateIdentityID = 21
	updateIDSwitch              updateIdentityID = 22
	updateIDBlow                updateIdentityID = 23
	updateIDMover               updateIdentityID = 24
	updateIDBlackPowderBarrel   updateIdentityID = 25
	updateIDOneSecondDie        updateIdentityID = 26
	updateIDWaterBarrel         updateIdentityID = 27
	updateIDSelfDestruct        updateIdentityID = 28
	updateIDBlackPowderBurn     updateIdentityID = 29
	updateIDDeathBall           updateIdentityID = 30
	updateIDDeathBallFragment   updateIdentityID = 31
	updateIDMoonglow            updateIdentityID = 32
	updateIDSentryGlobe         updateIdentityID = 33
	updateIDTelekinesis         updateIdentityID = 34
	updateIDFist                updateIdentityID = 35
	updateIDMeteorShower        updateIdentityID = 36
	updateIDMeteor              updateIdentityID = 37
	updateIDToxicCloud          updateIdentityID = 38
	updateIDSmallToxicCloud     updateIdentityID = 39
	updateIDArachnaphobia       updateIdentityID = 40
	updateIDExpire              updateIdentityID = 41
	updateIDBreak               updateIdentityID = 42
	updateIDOpen                updateIdentityID = 43
	updateIDBreakAndRemove      updateIdentityID = 44
	updateIDChakramInMotion     updateIdentityID = 45
	updateIDFlag                updateIdentityID = 46
	updateIDTrapDoor            updateIdentityID = 47
	updateIDBall                updateIdentityID = 48
	updateIDCrown               updateIdentityID = 49
	updateIDUndeadKiller        updateIdentityID = 50
	updateIDHarpoon             updateIdentityID = 51
	updateIDMonsterGenerator    updateIdentityID = 52
	updateIDMkgmtime            updateIdentityID = 53
	updateIDPlayerMonsterBot    updateIdentityID = 54
	updateIDPlayerObserver      updateIdentityID = 55
)

func updateIdentityKey(id updateIdentityID) unsafe.Pointer {
	return unsafe.Pointer(&updateIdentitySlots[id])
}
