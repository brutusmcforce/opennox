package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/libs/spell"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func worldCollideSound(id int32, u *server.Object) {
	GetServer().S().Audio.EventObj(sound.ID(id), u, 0, 0)
}
func worldCollideMessage(u *server.Object, text string) {
	gameplayTextPrivate(u, alloc.InternCString(text), 0)
}

func worldCollideMass(a, b *server.Object) {
	if a == nil || b == nil {
		return
	}
	ma, mb := float64(a.Mass), float64(b.Mass)
	// The qualified C build keeps the sum and ab in x87 PC53 registers,
	// but spills bb to float32 before applying it to either velocity.
	sum := ma + mb
	aa := (ma - mb) / sum
	ab := (mb + mb) / sum
	bb := float64(float32((mb - ma) / sum))
	ba := (ma + ma) / sum
	ay := float64(a.VelVec.Y)*aa + ab*float64(b.VelVec.Y)
	bx := float32(bb*float64(b.VelVec.X) + ba*float64(a.VelVec.X))
	by := float32(bb*float64(b.VelVec.Y) + ba*float64(a.VelVec.Y))
	a.VelVec.X = float32(float64(a.VelVec.X)*aa + ab*float64(b.VelVec.X))
	a.VelVec.Y = float32(ay)
	b.VelVec.Y = by
	b.VelVec.X = bx
}
func worldCollidePickup(a, b *server.Object) uint32 {
	if b != nil && b.ObjClass&2 == 0 && GetServer().S().Frame()-a.Field32 >= uint32(int32(GetServer().S().TickRate())>>1) && (b.ObjClass&4 == 0 || *(*byte)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](b.UpdateData).Field60))&1 != 0) {
		return uint32(bool2int(Nox_xxx_inventoryServPlace_4F36F0(b, a, 1, 1)))
	}
	return controlRaw(b)
}
func worldCollideBarrel(a *server.Object) {
	if frame := GetServer().S().Frame(); frame > a.Field34+3 {
		a.Field34 = frame
		worldCollideSound(281, a)
	}
}
func worldCollideAudio(a, b *server.Object) {
	if b != nil && b.ObjClass&4 != 0 {
		if frame := GetServer().S().Frame(); frame > a.Field34+30 {
			a.Field34 = frame
			worldCollideSound(int32(a.CollideDataAudioEvent().Sound), a)
		}
	}
}
func worldCollidePentagram(a *server.Object) uint32 {
	layout.As[server.TransporterUpdateData](a.UpdateData).Active = 1
	return controlRaw(a)
}
func worldCollideSign(a, b *server.Object) {
	if b != nil && b.ObjClass&4 != 0 {
		a.Use.Get()(b, a)
	}
}
func worldCollideTrap(a, b *server.Object) {
	d := a.CollideDataHole()
	if b == nil || b.ObjClass&0x80 != 0 {
		return
	}
	if a.ObjFlags&0x1000000 != 0 {
		switch b.Shape.Kind {
		case server.ShapeKindBox:
			if a.Shape.Box.W < b.Shape.Box.W || a.Shape.Box.H < b.Shape.Box.H {
				return
			}
		case server.ShapeKindCircle:
			diameter := float64(b.Shape.Circle.R) + float64(b.Shape.Circle.R)
			if !(diameter <= float64(a.Shape.Box.W) && diameter <= float64(a.Shape.Box.H)) {
				return
			}
		}
		if collisionContains(&a.PosVec, (*[11]float32)(unsafe.Pointer(&a.Shape)), &b.PosVec) {
			b.ObjFlags |= 0x60000
			b.Field41 = math.Float32bits(float32(d.Dest.X))
			b.Field42 = math.Float32bits(float32(d.Dest.Y))
			b.Pos39 = a.PosVec
		}
		return
	}
	if d.Triggered == 0 && (b.ObjClass&4 == 0 || !GetServer().S().Abils.IsActive(b, server.Ability(4))) {
		if delay := d.Delay; delay != 0 {
			d.Frame = GetServer().S().Frame() + uint32(delay)
		}
		GetServer().NoxScriptC().ScriptCallback(&d.Script, b, a, server.ScriptEventType(20))
		d.Triggered = 1
	}
}
func worldCollideTeleport(a, b *server.Object) {
	d := a.CollideDataTeleport()
	if b == nil || b.ObjClass&0x80 != 0 {
		return
	}
	visibilityFXPoint(138, b.PosVec)
	worldCollideSound(147, b)
	b.Field41 = math.Float32bits(float32(d.Dest.X))
	b.Field42 = math.Float32bits(float32(d.Dest.Y))
	stateTeleport(b, (*types.Pointf)(unsafe.Pointer(&b.Field41)))
	visibilityFXPoint(137, b.PosVec)
	worldCollideSound(147, b)
}
func worldCollideSpellAward(a, b *server.Object) uint32 {
	if b == nil {
		return 0
	}
	return uint32(Nox_xxx_spellGrantToPlayer_4FB550(b, spell.ID(int32(a.CollideDataAwardSpell().Spell)), 1, 0, 0))
}
func worldCollideUndead(a, b *server.Object, normal bool) {
	if b == nil {
		if !normal {
			GetServer().DelayedDelete(a)
		}
		return
	}
	if b.ObjClass&2 == 0 || b.ObjSubClass&0x40 == 0 {
		return
	}
	pool := (*server.DurSpell)(a.CollideDataUndeadKiller().SpellPtr()) // the duration spell (see sustainedTurnUndeadStart)
	hp := uint32(uint16(resourceGetHP(b)))
	available := pool.Field72 // the remaining hit points to drain
	if available <= int32(hp) {
		if available != 0 {
			b.CallDamage(a.FindOwnerChainPlayer(), a, int(available), object.DamageType(6))
			GetServer().DelayedDelete(a)
			pool.Field72 -= available
		} else {
			GetServer().DelayedDelete(a)
		}
	} else {
		b.CallDamage(a.FindOwnerChainPlayer(), a, int(hp), object.DamageType(6))
		pool.Field72 = available - int32(hp)
	}
}
func worldCollideGenerator(a, b *server.Object) {
	if b != nil && b.ObjClass&4 != 0 {
		GetServer().NoxScriptC().ScriptCallback(a.UpdateDataMonsterGen().ScriptAt(72), b, a, server.ScriptEventType(19))
	}
}
