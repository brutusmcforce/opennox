package legacy

import "unsafe"

type xferIdentityID uint8

const (
	xferIDDefault xferIdentityID = iota
	xferIDSpellPagePedestal
	xferIDSpellReward
	xferIDAbilityReward
	xferIDFieldGuide
	xferIDReadable
	xferIDExit
	xferIDDoor
	xferIDTrigger
	xferIDMonster
	xferIDHole
	xferIDTransporter
	xferIDElevator
	xferIDElevatorShaft
	xferIDMover
	xferIDGlyph
	xferIDInvisibleLight
	xferIDSentry
	xferIDWeapon
	xferIDArmor
	xferIDTeam
	xferIDGold
	xferIDAmmo
	xferIDNPC
	xferIDObelisk
	xferIDToxicCloud
	xferIDMonsterGenerator
	xferIDRewardMarker
)

type xferIdentitySlot struct {
	marker byte
}

var xferIdentitySlots [28]xferIdentitySlot

func xferIdentityKey(id xferIdentityID) unsafe.Pointer {
	return unsafe.Pointer(&xferIdentitySlots[id])
}
