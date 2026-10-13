package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
)

// Bot update data remains C-owned because the object lifecycle frees it in C.
// Morphing temporarily exchanges player and monster records without copying them.
func controlBotCreate(u *server.Object) uint32 {
	d := u.UpdateDataPlayer()
	result := controlRaw(u)
	if d.BotData == nil {
		p := legacyCalloc(1, unsafe.Sizeof(server.MonsterUpdateData{}))
		layout.TrackAs[server.MonsterUpdateData](p)
		d.BotData = (*server.MonsterUpdateData)(p)
		result = memmap.PtrWord(p)
	}
	b := d.BotData
	if b == nil {
		return result
	}
	*b = server.MonsterUpdateData{}
	b.BotPlayerData = d
	b.MonsterDef = monsterDefinitionByType(uint32(GetServer().S().Types.IndByID("NPC")))
	b.RetreatLevel = math.Float32frombits(1048576000)
	b.ResumeLevel = math.Float32frombits(1061997773)
	b.StatusFlags = 186376
	b.AIStack[0].Action = 5
	b.AIAction340 = 38
	b.Aggression2 = math.Float32frombits(1056964608)
	b.Aggression = math.Float32frombits(1062501089)
	b.SightRange = math.Float32frombits(1125515264)
	b.Field329 = math.Float32frombits(1106247680)
	b.Field330 = math.Float32frombits(1065353216)
	b.Field332 = math.Float32frombits(1056964608)
	b.Field338 = math.Float32frombits(1065353216)
	b.Field510 = 3
	b.DialogStartFunc = -1
	b.DialogEndFunc = -1
	b.Field0 = 0xdeadface
	b.Field335 = 1
	b.Field337 = 1
	b.Field333 = 255
	b.Field331 = 30
	for _, cb := range []*server.ScriptCallback{
		&b.ScriptLookingForEnemy, &b.ScriptEnemySighted, &b.ScriptChangeFocus, &b.ScriptIsHit, &b.ScriptRetreat,
		&b.ScriptDeath, &b.ScriptCollision, &b.ScriptHearEnemy, &b.ScriptEndOfWaypoint, &b.ScriptLostEnemy,
	} {
		cb.Func = -1
	}
	b.Direction94 = uint32(int32(int16(u.Direction1)))
	b.Pos95 = u.PosVec
	class := d.Player.InfoBytes[66] // PlayerClassVal
	if class == 0 {
		return 0
	}
	if class > 2 {
		return uint32(class) - 2
	}
	fps := GetServer().S().TickRate()
	f := uint32(fps)
	b.FleeRange = math.Float32frombits(1112014848)
	b.Field410 = 0x8000000
	b.StatusFlags |= 0x20
	b.Field362_2 = uint16(f >> 1)
	b.Field366_0 = uint16(3 * f)
	b.Field366_2 = uint16(30 * f)
	b.Field368_2 = uint16(2 * f)
	b.Field370_2 = uint16(6 * f)
	b.Field376 = 0x80000000
	if class == 2 {
		b.Field430 = 0x10000000
		b.Field432 = 0x20000000
		b.Field446 = 0x20000000
		for _, p := range []*uint32{&b.Field401, &b.Field424, &b.Field456, &b.Field455, &b.Field464} {
			*p = 0x40000000
		}
		b.Field364_2 = uint16(6 * f)
		return f
	}
	b.Field423 = 0x10000000
	b.Field408 = 0x10000000
	b.Field411 = 0x10000000
	b.Field384 = 0x20000000
	b.Field405 = 0x20000000
	b.Field364_2 = uint16(f)
	n := GetServer().S().Rand.Logic.IntClamp(0, 100)
	if n < 33 {
		b.Field399 = 0x40000000
	} else if n < 66 {
		b.Field388 = 0x40000000
	} else {
		b.Field415 = 0x40000000
		b.Field422 = 0x40000000
	}
	return 6 * f
}
func controlMorphFromPlayer(u *server.Object) int8 {
	if u.ObjClass&4 != 0 {
		u.ObjClass = u.ObjClass&^4 | 2
		u.UpdateData = unsafe.Pointer(u.UpdateDataPlayer().BotData)
		u.ObjSubClass = 16
	}
	return int8(u.ObjClass)
}
func controlMorphToPlayer(u *server.Object) int8 {
	if u.ObjClass&2 != 0 {
		u.ObjClass = u.ObjClass&^2 | 4
		u.UpdateData = unsafe.Pointer(u.UpdateDataMonster().BotPlayerData)
		u.ObjSubClass = 0
	}
	return int8(u.ObjClass)
}
func controlBotUpdate(u *server.Object) uint32 {
	d := u.UpdateDataPlayer()
	if d.BotData == nil {
		controlBotCreate(u)
	}
	b := d.BotData
	if b == nil {
		u.Update = updateIdentityKey(updateIDPlayer)
		return 0
	}
	result := controlRespawnBot(u)
	if result != 0 {
		return uint32(result)
	}
	b.StatusFlags |= 0x100
	controlMorphFromPlayer(u)
	Nox_xxx_unitUpdateMonster_50A5C0(u)
	controlMorphToPlayer(u)
	d.State = server.PlayerState(controlBotState(u))
	d.Field59_0 = b.Field120_1
	d.Player.Pos3632Vec = u.PosVec
	return math.Float32bits(u.PosVec.Y)
}
