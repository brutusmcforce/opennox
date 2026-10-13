package server

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/types"
)

// Collide and death data blocks. They keep the 386 layout of the original game and are pointer-free:
// object and spell references are stored as memmap pointer words.

func collideDataAs[T any](obj *Object) *T {
	if alloc.IsDead(obj.CollideData) {
		panic("object already deleted")
	}
	// TODO: verify this conversion by checking ObjectType
	return (*T)(obj.CollideData)
}

// ProjectileCollideData is the collide data of projectiles (ProjectileCollide, ProjectileSparkCollide,
// WallReflectCollide, WallReflectSparkCollide, PixieCollide, SparkCollide, YellowStarShotCollide,
// ArrowCollide, HarpoonCollide).
type ProjectileCollideData struct {
	Damage uint32 // 0, 0: damage on contact
	Owner  uint32 // 1, 4: the shooter of arrows and harpoon bolts (as a memmap word)
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(ProjectileCollideData{})]

// OwnerObj decodes the Owner pointer word.
func (d *ProjectileCollideData) OwnerObj() *Object {
	return (*Object)(memmap.WordPtr(d.Owner))
}

// SetOwnerObj stores obj as the Owner pointer word.
func (d *ProjectileCollideData) SetOwnerObj(obj *Object) {
	d.Owner = memmap.PtrWord(unsafe.Pointer(obj))
}

// CollideDataProjectile returns the collide data of projectiles (see ProjectileCollideData).
func (obj *Object) CollideDataProjectile() *ProjectileCollideData {
	return collideDataAs[ProjectileCollideData](obj)
}

// UnknownCollideData is the collide data of collide handlers that never access it (BombCollide, ElevatorCollide).
type UnknownCollideData struct {
	Field0 uint32 // 0, 0
	Field4 uint32 // 1, 4
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(UnknownCollideData{})]

// DamageCollideData is the collide data of damage fields (DamageCollide).
type DamageCollideData struct {
	Damage     byte    // 0, 0
	Field1     [3]byte // 0, 1
	DamageType uint32  // 1, 4
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(DamageCollideData{})]

// CollideDataDamage returns the collide data of damage fields (DamageCollide).
func (obj *Object) CollideDataDamage() *DamageCollideData {
	return collideDataAs[DamageCollideData](obj)
}

// ManaDrainCollideData is the collide data of mana drains (ManaDrainCollide).
type ManaDrainCollideData struct {
	Mana   byte    // 0, 0: mana drained per collision
	Field1 [3]byte // 0, 1
	Field4 uint32  // 1, 4
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(ManaDrainCollideData{})]

// CollideDataManaDrain returns the collide data of mana drains (ManaDrainCollide).
func (obj *Object) CollideDataManaDrain() *ManaDrainCollideData {
	return collideDataAs[ManaDrainCollideData](obj)
}

// SparkExplosionCollideData is the collide data of fireballs (SparkExplosionCollide).
type SparkExplosionCollideData struct {
	Damage byte // 0, 0
}

var _ = [1]struct{}{}[1-unsafe.Sizeof(SparkExplosionCollideData{})]

// CollideDataSparkExplosion returns the collide data of fireballs (SparkExplosionCollide).
func (obj *Object) CollideDataSparkExplosion() *SparkExplosionCollideData {
	return collideDataAs[SparkExplosionCollideData](obj)
}

// AudioEventCollideData is the collide data of audio event triggers (AudioEventCollide).
type AudioEventCollideData struct {
	Sound uint32 // 0, 0: sound ID
}

var _ = [1]struct{}{}[4-unsafe.Sizeof(AudioEventCollideData{})]

// CollideDataAudioEvent returns the collide data of audio event triggers (AudioEventCollide).
func (obj *Object) CollideDataAudioEvent() *AudioEventCollideData {
	return collideDataAs[AudioEventCollideData](obj)
}

// TeleportCollideData is the collide data of teleports (TeleportCollide).
type TeleportCollideData struct {
	Dest ntype.Point32 // 0, 0: destination
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(TeleportCollideData{})]

// CollideDataTeleport returns the collide data of teleports (TeleportCollide).
func (obj *Object) CollideDataTeleport() *TeleportCollideData {
	return collideDataAs[TeleportCollideData](obj)
}

// TeleportWakeCollideData is the collide data of teleport wakes (TeleportWakeCollide).
type TeleportWakeCollideData struct {
	Dest types.Pointf // 0, 0: destination
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(TeleportWakeCollideData{})]

// CollideDataTeleportWake returns the collide data of teleport wakes (TeleportWakeCollide).
func (obj *Object) CollideDataTeleportWake() *TeleportWakeCollideData {
	return collideDataAs[TeleportWakeCollideData](obj)
}

// AwardSpellCollideData is the collide data of spell page pedestals (AwardSpellCollide).
type AwardSpellCollideData struct {
	Spell uint32 // 0, 0: spell ID
}

var _ = [1]struct{}{}[4-unsafe.Sizeof(AwardSpellCollideData{})]

// CollideDataAwardSpell returns the collide data of spell page pedestals (AwardSpellCollide).
func (obj *Object) CollideDataAwardSpell() *AwardSpellCollideData {
	return collideDataAs[AwardSpellCollideData](obj)
}

// UndeadKillerCollideData is the collide data of turn undead waves (UndeadKillerCollide).
type UndeadKillerCollideData struct {
	Spell uint32 // 0, 0: the turn undead duration spell (as a memmap word)
}

var _ = [1]struct{}{}[4-unsafe.Sizeof(UndeadKillerCollideData{})]

// SpellPtr decodes the Spell pointer word.
func (d *UndeadKillerCollideData) SpellPtr() unsafe.Pointer {
	return memmap.WordPtr(d.Spell)
}

// CollideDataUndeadKiller returns the collide data of turn undead waves (UndeadKillerCollide).
func (obj *Object) CollideDataUndeadKiller() *UndeadKillerCollideData {
	return collideDataAs[UndeadKillerCollideData](obj)
}

// MonsterArrowCollideData is the collide data of monster arrows (MonsterArrowCollide).
type MonsterArrowCollideData struct {
	DamageCoop uint32 // 0, 0: damage in coop games
	Damage     uint32 // 1, 4: damage in other games
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(MonsterArrowCollideData{})]

// CollideDataMonsterArrow returns the collide data of monster arrows (MonsterArrowCollide).
func (obj *Object) CollideDataMonsterArrow() *MonsterArrowCollideData {
	return collideDataAs[MonsterArrowCollideData](obj)
}

// HoleCollideData is the collide data of trap doors and holes (TrapDoorCollide).
type HoleCollideData struct {
	Script    ScriptCallback // 0, 0: script called on collision
	Dest      ntype.Point32  // 2, 8: fall destination (holes)
	Frame     uint32         // 4, 16: frame when the trap door rearms
	Delay     uint16         // 5, 20: rearm delay in frames
	Field22   uint16         // 5, 22
	Triggered uint32         // 6, 24
}

var _ = [1]struct{}{}[28-unsafe.Sizeof(HoleCollideData{})]

// CollideDataHole returns the collide data of trap doors and holes (TrapDoorCollide).
func (obj *Object) CollideDataHole() *HoleCollideData {
	return collideDataAs[HoleCollideData](obj)
}

// CollideDataHole returns the collide data of trap door and hole types (TrapDoorCollide).
func (t *ObjectType) CollideDataHole() *HoleCollideData {
	// TODO: verify this conversion by checking the collide function
	return (*HoleCollideData)(t.CollideData)
}

// ExitCollideData is the collide data of exits (ExitCollide).
type ExitCollideData struct {
	MapName [80]byte     // 0, 0: the map to switch to
	Dest    types.Pointf // 20, 80: position to return to (quest warp gates, saves)
}

var _ = [1]struct{}{}[88-unsafe.Sizeof(ExitCollideData{})]

// CollideDataExit returns the collide data of exits (ExitCollide).
func (obj *Object) CollideDataExit() *ExitCollideData {
	return collideDataAs[ExitCollideData](obj)
}

// SoulGateCollideData is the collide data of soul gates (SoulGateCollide).
type SoulGateCollideData struct {
	Frame uint32 // 0, 0: frame of the last collision
}

var _ = [1]struct{}{}[4-unsafe.Sizeof(SoulGateCollideData{})]

// CollideDataSoulGate returns the collide data of soul gates (SoulGateCollide).
func (obj *Object) CollideDataSoulGate() *SoulGateCollideData {
	return collideDataAs[SoulGateCollideData](obj)
}

// CreateObjectDeathData is the death data of objects that create another object on death
// (CreateObjectDie, SpawnObjectDie).
type CreateObjectDeathData struct {
	Name  [128]byte // 0, 0: type ID of the object to create
	Sound uint32    // 32, 128: sound ID
}

var _ = [1]struct{}{}[132-unsafe.Sizeof(CreateObjectDeathData{})]

// DeathDataCreateObject returns the death data of CreateObjectDie and SpawnObjectDie.
func (obj *Object) DeathDataCreateObject() *CreateObjectDeathData {
	// TODO: verify this conversion by checking ObjectType
	return (*CreateObjectDeathData)(obj.DeathData)
}
