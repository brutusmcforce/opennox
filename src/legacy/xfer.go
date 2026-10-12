package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/server"
)

var (
	Nox_xxx_xfer_saveObj51DF90        func(cf *cryptfile.CryptFile, a1p *server.Object) int
	Nox_xxx_XFerDefault4F49A0         func(cf *cryptfile.CryptFile, v1 *server.Object, a2 unsafe.Pointer) error
	Nox_xxx_XFer_ReadShopItem_52A840  func(a1 unsafe.Pointer, a2 int)
	Nox_xxx_XFer_WriteShopItem_52A5F0 func(a1 unsafe.Pointer)
)

func init() {
	server.RegisterObjectXferGo("DefaultXfer", xferIdentityKey(xferIDDefault), func(u *server.Object, arg unsafe.Pointer) int {
		if err := Nox_xxx_XFerDefault4F49A0(cryptfile.Global(), u, arg); err != nil {
			mapLog.Println("nox_xxx_XFerDefault_4F49A0:", err)
			return 0
		}
		return 1
	})
	server.RegisterObjectXferGo("SpellPagePedestalXfer", xferIdentityKey(xferIDSpellPagePedestal), func(u *server.Object, arg unsafe.Pointer) int { return objectXferPedestal(u) })
	server.RegisterObjectXferGo("SpellRewardXfer", xferIdentityKey(xferIDSpellReward), func(u *server.Object, arg unsafe.Pointer) int { return itemXferSpellReward(u) })
	server.RegisterObjectXferGo("AbilityRewardXfer", xferIdentityKey(xferIDAbilityReward), func(u *server.Object, arg unsafe.Pointer) int { return itemXferAbilityReward(u) })
	server.RegisterObjectXferGo("FieldGuideXfer", xferIdentityKey(xferIDFieldGuide), func(u *server.Object, arg unsafe.Pointer) int { return itemXferFieldGuide(u) })
	server.RegisterObjectXferGo("ReadableXfer", xferIdentityKey(xferIDReadable), func(u *server.Object, arg unsafe.Pointer) int { return objectXferReadable(u) })
	server.RegisterObjectXferGo("ExitXfer", xferIdentityKey(xferIDExit), func(u *server.Object, arg unsafe.Pointer) int { return objectXferExit(u) })
	server.RegisterObjectXferGo("DoorXfer", xferIdentityKey(xferIDDoor), func(u *server.Object, arg unsafe.Pointer) int { return objectXferDoor(u) })
	server.RegisterObjectXferGo("TriggerXfer", xferIdentityKey(xferIDTrigger), func(u *server.Object, arg unsafe.Pointer) int { return objectXferTrigger(u) })
	server.RegisterObjectXferGo("MonsterXfer", xferIdentityKey(xferIDMonster), func(u *server.Object, arg unsafe.Pointer) int { return creatureXferMonster(u) })
	server.RegisterObjectXferGo("HoleXfer", xferIdentityKey(xferIDHole), func(u *server.Object, arg unsafe.Pointer) int { return objectXferHole(u) })
	server.RegisterObjectXferGo("TransporterXfer", xferIdentityKey(xferIDTransporter), func(u *server.Object, arg unsafe.Pointer) int { return objectXferTransporter(u) })
	server.RegisterObjectXferGo("ElevatorXfer", xferIdentityKey(xferIDElevator), func(u *server.Object, arg unsafe.Pointer) int { return objectXferElevator(u) })
	server.RegisterObjectXferGo("ElevatorShaftXfer", xferIdentityKey(xferIDElevatorShaft), func(u *server.Object, arg unsafe.Pointer) int { return objectXferShaft(u) })
	server.RegisterObjectXferGo("MoverXfer", xferIdentityKey(xferIDMover), func(u *server.Object, arg unsafe.Pointer) int { return objectXferMover(u) })
	server.RegisterObjectXferGo("GlyphXfer", xferIdentityKey(xferIDGlyph), func(u *server.Object, arg unsafe.Pointer) int { return objectXferGlyph(u) })
	server.RegisterObjectXferGo("InvisibleLightXfer", xferIdentityKey(xferIDInvisibleLight), func(u *server.Object, arg unsafe.Pointer) int { return objectXferLight(u) })
	server.RegisterObjectXferGo("SentryXfer", xferIdentityKey(xferIDSentry), func(u *server.Object, arg unsafe.Pointer) int { return objectXferSentry(u) })
	server.RegisterObjectXferGo("WeaponXfer", xferIdentityKey(xferIDWeapon), func(u *server.Object, arg unsafe.Pointer) int { return itemXferWeapon(u) })
	server.RegisterObjectXferGo("ArmorXfer", xferIdentityKey(xferIDArmor), func(u *server.Object, arg unsafe.Pointer) int { return itemXferArmor(u) })
	server.RegisterObjectXferGo("TeamXfer", xferIdentityKey(xferIDTeam), func(u *server.Object, arg unsafe.Pointer) int { return itemXferTeam(u) })
	server.RegisterObjectXferGo("GoldXfer", xferIdentityKey(xferIDGold), func(u *server.Object, arg unsafe.Pointer) int { return itemXferGold(u) })
	server.RegisterObjectXferGo("AmmoXfer", xferIdentityKey(xferIDAmmo), func(u *server.Object, arg unsafe.Pointer) int { return itemXferAmmo(u) })
	server.RegisterObjectXferGo("NPCXfer", xferIdentityKey(xferIDNPC), func(u *server.Object, arg unsafe.Pointer) int { return creatureXferNPC(u) })
	server.RegisterObjectXferGo("ObeliskXfer", xferIdentityKey(xferIDObelisk), func(u *server.Object, arg unsafe.Pointer) int { return itemXferObelisk(u) })
	server.RegisterObjectXferGo("ToxicCloudXfer", xferIdentityKey(xferIDToxicCloud), func(u *server.Object, arg unsafe.Pointer) int { return itemXferToxicCloud(u) })
	server.RegisterObjectXferGo("MonsterGeneratorXfer", xferIdentityKey(xferIDMonsterGenerator), func(u *server.Object, arg unsafe.Pointer) int { return itemXferGenerator(u) })
	server.RegisterObjectXferGo("RewardMarkerXfer", xferIdentityKey(xferIDRewardMarker), func(u *server.Object, arg unsafe.Pointer) int { return itemXferRewardMarker(u) })
}

func Get_nox_xxx_XFerFieldGuide_4F6390() unsafe.Pointer {
	return xferIdentityKey(xferIDFieldGuide)
}

func Get_nox_xxx_XFerAbilityReward_4F6240() unsafe.Pointer {
	return xferIdentityKey(xferIDAbilityReward)
}
func Nox_xxx_mapReadWriteObjData_4F4530(a1 *server.Object, a2 int) int {
	return objectXferCommon(a1, a2)
}
func Nox_xxx_xfer_4F3E30(a1 int, a2 *server.Object, a3 uint32) int {
	return objectXferInventory(uint16(a1), a2, int32(a3))
}
