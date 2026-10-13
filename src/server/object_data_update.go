package server

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
)

// Update data blocks of simple objects. They keep the 386 layout of the original game and are pointer-free:
// object references are stored as memmap pointer words.

// LifetimeUpdateData is the update data of objects that expire after a fixed time (LifetimeUpdate).
type LifetimeUpdateData struct {
	Lifetime uint32 // 0, 0: frames after creation before the object dies
}

var _ = [1]struct{}{}[4-unsafe.Sizeof(LifetimeUpdateData{})]

// UpdateDataLifetime returns the update data of objects that expire after a fixed time (LifetimeUpdate).
func (obj *Object) UpdateDataLifetime() *LifetimeUpdateData {
	return updateDataAs[LifetimeUpdateData](obj)
}

// ToxicCloudUpdateData is the update data of toxic clouds (ToxicCloudUpdate, SmallToxicCloudUpdate).
type ToxicCloudUpdateData struct {
	Life uint32 // 0, 0: remaining frames, counts down to 0
}

var _ = [1]struct{}{}[4-unsafe.Sizeof(ToxicCloudUpdateData{})]

// UpdateDataToxicCloud returns the update data of toxic clouds (ToxicCloudUpdate, SmallToxicCloudUpdate).
func (obj *Object) UpdateDataToxicCloud() *ToxicCloudUpdateData {
	return updateDataAs[ToxicCloudUpdateData](obj)
}

// DamageUpdateData is the update data of falling spell objects that deal damage
// (FistUpdate, MeteorShowerUpdate, MeteorUpdate).
type DamageUpdateData struct {
	Damage uint32 // 0, 0: damage dealt on impact; a meteor shower passes it to its meteors
}

var _ = [1]struct{}{}[4-unsafe.Sizeof(DamageUpdateData{})]

// UpdateDataDamage returns the update data of falling spell objects (FistUpdate, MeteorShowerUpdate, MeteorUpdate).
func (obj *Object) UpdateDataDamage() *DamageUpdateData {
	return updateDataAs[DamageUpdateData](obj)
}

// HarpoonUpdateData is the update data of harpoon bolts (HarpoonUpdate).
type HarpoonUpdateData struct {
	Target uint32 // 0, 0: homing target (as a memmap word)
}

var _ = [1]struct{}{}[4-unsafe.Sizeof(HarpoonUpdateData{})]

// TargetObj decodes the Target pointer word.
func (d *HarpoonUpdateData) TargetObj() *Object {
	return (*Object)(memmap.WordPtr(d.Target))
}

// SetTargetObj stores obj as the Target pointer word.
func (d *HarpoonUpdateData) SetTargetObj(obj *Object) {
	d.Target = memmap.PtrWord(unsafe.Pointer(obj))
}

// UpdateDataHarpoon returns the update data of harpoon bolts (HarpoonUpdate).
func (obj *Object) UpdateDataHarpoon() *HarpoonUpdateData {
	return updateDataAs[HarpoonUpdateData](obj)
}

// HomingProjectileUpdateData is the update data of homing projectiles (HomingProjectileUpdate, no handler).
type HomingProjectileUpdateData struct {
	Field0 uint32 // 0, 0
}

var _ = [1]struct{}{}[4-unsafe.Sizeof(HomingProjectileUpdateData{})]

// TowerUpdateData is the update data of towers (TowerUpdate, no handler).
type TowerUpdateData struct {
	Field0 uint32 // 0, 0
	Field4 uint32 // 1, 4
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(TowerUpdateData{})]

// WeaponArmorUpdateData is the update data of weapons and armor (WeaponArmorUpdate, DamageRoundoffUpdate).
type WeaponArmorUpdateData struct {
	Roundoff float32 // 0, 0: fractional durability damage carried over between hits
	Flags    uint32  // 1, 4: 1 if the item was given on respawn (sells for nothing in quest)
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(WeaponArmorUpdateData{})]

// UpdateDataWeaponArmor returns the update data of weapons and armor (WeaponArmorUpdate, DamageRoundoffUpdate).
func (obj *Object) UpdateDataWeaponArmor() *WeaponArmorUpdateData {
	return updateDataAs[WeaponArmorUpdateData](obj)
}

// CrownUpdateData is the update data of crowns (CrownUpdate).
type CrownUpdateData struct {
	LastOwner uint32 // 0, 0: cleared when that object is destroyed (as a memmap word)
	Pending   uint32 // 1, 4: player the crown goes to next (as a memmap word)
	Field8    uint32 // 2, 8
}

var _ = [1]struct{}{}[12-unsafe.Sizeof(CrownUpdateData{})]

// LastOwnerObj decodes the LastOwner pointer word.
func (d *CrownUpdateData) LastOwnerObj() *Object {
	return (*Object)(memmap.WordPtr(d.LastOwner))
}

// PendingObj decodes the Pending pointer word.
func (d *CrownUpdateData) PendingObj() *Object {
	return (*Object)(memmap.WordPtr(d.Pending))
}

// SetPendingObj stores obj as the Pending pointer word.
func (d *CrownUpdateData) SetPendingObj(obj *Object) {
	d.Pending = memmap.PtrWord(unsafe.Pointer(obj))
}

// UpdateDataCrown returns the update data of crowns (CrownUpdate).
func (obj *Object) UpdateDataCrown() *CrownUpdateData {
	return updateDataAs[CrownUpdateData](obj)
}

// FlagUpdateData is the update data of CTF flags (FlagUpdate).
type FlagUpdateData struct {
	Home      types.Pointf // 0, 0: position the flag returns to
	DropFrame uint32       // 2, 8: frame the flag was dropped, 0 if it is home or carried
}

var _ = [1]struct{}{}[12-unsafe.Sizeof(FlagUpdateData{})]

// UpdateDataFlag returns the update data of CTF flags (FlagUpdate).
func (obj *Object) UpdateDataFlag() *FlagUpdateData {
	return updateDataAs[FlagUpdateData](obj)
}

// PushUpdateData is the update data of push fields (PushUpdate).
type PushUpdateData struct {
	Radius  float32 // 0, 0
	Radius2 float32 // 1, 4: copy of Radius
	Power   float32 // 2, 8
}

var _ = [1]struct{}{}[12-unsafe.Sizeof(PushUpdateData{})]

// UpdateDataPush returns the update data of push fields (PushUpdate).
func (obj *Object) UpdateDataPush() *PushUpdateData {
	return updateDataAs[PushUpdateData](obj)
}

// SentryGlobeUpdateData is the update data of sentry globes (SentryGlobeUpdate).
type SentryGlobeUpdateData struct {
	Angle     float32 // 0, 0: current beam angle (radians)
	BaseAngle float32 // 1, 4: angle restored while the sentry is disabled
	Speed     float32 // 2, 8: angle step per frame
}

var _ = [1]struct{}{}[12-unsafe.Sizeof(SentryGlobeUpdateData{})]

// UpdateDataSentryGlobe returns the update data of sentry globes (SentryGlobeUpdate).
func (obj *Object) UpdateDataSentryGlobe() *SentryGlobeUpdateData {
	return updateDataAs[SentryGlobeUpdateData](obj)
}

// LoopAndDamageUpdateData is the update data of periodic damage objects (LoopAndDamageUpdate).
// The three parsed values are never read.
type LoopAndDamageUpdateData struct {
	Field0  uint32 // 0, 0
	Field4  uint32 // 1, 4
	Field8  uint32 // 2, 8
	Field12 uint32 // 3, 12
}

var _ = [1]struct{}{}[16-unsafe.Sizeof(LoopAndDamageUpdateData{})]

// SparkUpdateData is the update data of sparks (SparkUpdate, initialized by SparkInit).
type SparkUpdateData struct {
	MaxLife uint32 // 0, 0: initial lifetime in frames
	Life    uint32 // 1, 4: remaining frames
	Field8  uint32 // 2, 8
	Stage   uint32 // 3, 12: spark kind (4 does not fade, 5 webbing)
}

var _ = [1]struct{}{}[16-unsafe.Sizeof(SparkUpdateData{})]

// UpdateDataSpark returns the update data of sparks (SparkUpdate).
func (obj *Object) UpdateDataSpark() *SparkUpdateData {
	return updateDataAs[SparkUpdateData](obj)
}

// BallUpdateData is the update data of flag balls (BallUpdate).
type BallUpdateData struct {
	LastOwner       uint32  // 0, 0: last player that held the ball (as a memmap word)
	LastTeam        uint32  // 1, 4: team of LastOwner
	ResetTicks      uint64  // 2, 8: platform ticks of the last possession, the ball resets after 20 seconds
	PossessFrame    uint32  // 4, 16: frame LastOwner picked up the ball
	PossessDuration uint32  // 5, 20: frames a player may keep the ball
	ResetVel        float32 // 6, 24: LastOwner is forgotten when the ball is slower than this
	Field28         uint32  // 7, 28
}

var _ = [1]struct{}{}[32-unsafe.Sizeof(BallUpdateData{})]

// LastOwnerObj decodes the LastOwner pointer word.
func (d *BallUpdateData) LastOwnerObj() *Object {
	return (*Object)(memmap.WordPtr(d.LastOwner))
}

// SetLastOwnerObj stores obj as the LastOwner pointer word.
func (d *BallUpdateData) SetLastOwnerObj(obj *Object) {
	d.LastOwner = memmap.PtrWord(unsafe.Pointer(obj))
}

// UpdateDataBall returns the update data of flag balls (BallUpdate).
func (obj *Object) UpdateDataBall() *BallUpdateData {
	return updateDataAs[BallUpdateData](obj)
}

// SkullUpdateData is the update data of skull and arrow traps (SkullUpdate).
type SkullUpdateData struct {
	ScanDelay   uint32   // 0, 0: frames until the next scan for targets
	Cooldown    uint32   // 1, 4: frames until the next shot
	State       byte     // 2, 8: 1 while a target is in sight
	Field9      [3]byte  // 2, 9
	MissileType uint32   // 3, 12: type index of MissileName (set by SkullInit)
	MissileName [32]byte // 4, 16: missile type ID
	Active      byte     // 12, 48: set while the trap is enabled
	Field49     [3]byte  // 12, 49
}

var _ = [1]struct{}{}[52-unsafe.Sizeof(SkullUpdateData{})]

// UpdateDataSkull returns the update data of skull and arrow traps (SkullUpdate).
func (obj *Object) UpdateDataSkull() *SkullUpdateData {
	return updateDataAs[SkullUpdateData](obj)
}
