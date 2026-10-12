package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func spellEffectMovable(u *server.Object) uint32 {
	if u.ObjFlags&0x8068 != 0 {
		return 0
	}
	return (^uint32(u.ObjClass) >> 22) & 1
}

type spellEffectForceContext struct {
	pos                  types.Pointf
	radius, inner, power float32
	source               *server.Object
	callback             unsafe.Pointer
	arg                  uintptr
}

var spellEffectForceCallbacks = make(map[unsafe.Pointer]func(*server.Object, uint32, uintptr))

func spellEffectForce(u *server.Object, ctx spellEffectForceContext) {
	if spellEffectMovable(u) == 0 || !spellEffectTrace(ctx.pos, u.PosVec, 0) {
		return
	}
	dx := float32(u.PosVec.X - ctx.pos.X)
	dy64 := float64(u.PosVec.Y) - float64(ctx.pos.Y)
	dy := float32(dy64)
	dist := math.Sqrt(dy64*float64(dy)+float64(dx)*float64(dx)) + 0.1
	d := float32(dist)
	if dist > float64(ctx.radius) {
		return
	}
	force := ctx.power
	if float64(d) > float64(ctx.inner) {
		force = float32((1 - (float64(d)-float64(ctx.inner))/(float64(ctx.radius)-float64(ctx.inner))) * float64(ctx.power))
	}
	acceleration := float32(float64(force) / float64(u.Mass))
	if ctx.callback != nil {
		// Distance is a raw float word, and arg is opaque integer-sized data.
		// Neither may enter a Go pointer slot visible to the garbage collector.
		fn := spellEffectForceCallbacks[ctx.callback]
		if fn == nil {
			panic("unregistered spatial force callback")
		}
		fn(u, math.Float32bits(d), ctx.arg)
	}
	u.ForceVec.X = float32(float64(acceleration)*float64(dx)/float64(d) + float64(u.ForceVec.X))
	u.ForceVec.Y = float32(float64(acceleration)*float64(dy)/float64(d) + float64(u.ForceVec.Y))
	if u.ObjClass&1 == 0 {
		collisionActivate(u)
	}
}
func spellEffectPushUnit(u *server.Object, record unsafe.Pointer) {
	ctx := spellEffectForceContext{pos: spellEffectPos(*controlPtr(record, 0), 0), radius: *temporaryFloat(record, 4), inner: *temporaryFloat(record, 8), power: *temporaryFloat(record, 12), source: spellEffectObject(record, 16), callback: *controlPtr(record, 20), arg: uintptr(*spellLifeWord(record, 24))}
	spellEffectForce(u, ctx)
}
func spellEffectPushAround(pos types.Pointf, radius, inner, power float32, source *server.Object, callback unsafe.Pointer, arg uintptr) {
	ctx := spellEffectForceContext{pos: pos, radius: radius, inner: inner, power: float32(float64(power) * 10), source: source, callback: callback, arg: arg}
	if radius < inner {
		ctx.radius = inner
	}
	rect := types.Rectf{Min: types.Ptf(pos.X-radius, pos.Y-radius), Max: types.Ptf(pos.X+radius, pos.Y+radius)}
	GetServer().S().Map.EachObjAndMissileInRect(rect, func(u *server.Object) bool { spellEffectForce(u, ctx); return true })
}
func spellEffectPull(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	power := float32(-(spellEffectScalar("PullPowerCoeff") * float64(level)))
	spellEffectPushAround(c.PosVec, 600, 10, power, nil, nil, 0)
	spellEffectAudio(id, 0, b)
	return 1
}
func spellEffectPush(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	power := float32(spellEffectScalar("PushPowerCoeff") * float64(level))
	spellEffectPushAround(c.PosVec, 600, 10, power, nil, nil, 0)
	spellEffectAudio(id, 0, b)
	return 1
}
func spellEffectInversion(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	GetServer().S().Map.EachMissileInCircle(c.PosVec, float32(spellEffectScalar("InversionRange")), func(u *server.Object) bool { Nox_xxx_changeOwner_52BE40(u, b); return true })
	spellEffectAudio(id, 0, c)
	return 1
}
func spellEffectQuakeDamage(u, source *server.Object) int16 {
	owner := source.FindOwnerChainPlayer()
	ret := memmap.PtrWord(owner.CObj())
	if u != source {
		ret = uint32(u.ObjFlags)
		if ret&0x4000 == 0 {
			hp := resourceGetHP(u)
			ret = ret&0xffff0000 | uint32(uint16(hp))
			if uint16(hp) != 0 {
				distance := float32(stateDistance(u, source))
				scale := float32(1 - float64(distance)/spellEffectScalar("EarthquakeRange"))
				level := int32(*memmap.PtrUint32(0x5d4594, 2487700))
				damage := int32(int64(spellEffectTable("EarthquakeDamage", level-1) * float64(scale)))
				ret = uint32(projectileDamage(u, owner, source, damage, 11))
			}
		}
	}
	return int16(ret)
}
func spellEffectQuake(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	*memmap.PtrUint32(0x5d4594, 2487700) = uint32(level)
	GetServer().S().Map.EachObjInCircle(c.PosVec, float32(spellEffectScalar("EarthquakeRange")), func(u *server.Object) bool { spellEffectQuakeDamage(u, c); return true })
	spellEffectAudio(id, 0, c)
	jiggle := floatToInt32(float32(spellEffectTable("EarthquakeJiggle", level-1)))
	gameplayReportEarthquake(&c.PosVec, int(jiggle))
	return 1
}
func spellEffectDoorLink(u *server.Object) {
	if u.ObjClass&0x80 != 0 {
		u.ObjOwner = spellDoorLinkSource
		u.Field34 = GetServer().S().Frame() + 60*uint32(GetServer().S().TickRate())
	}
}
func spellEffectDoorCandidate(u, source *server.Object) {
	if u.ObjClass&0x80 == 0 {
		return
	}
	dx := float64(source.PosVec.X) - float64(u.PosVec.X)
	dy := float64(source.PosVec.Y) - float64(u.PosVec.Y)
	d := float32(dy*dy + dx*dx)
	if d > 22500 || float64(d) >= float64(*memmap.PtrFloat32(0x5d4594, 2487704)) {
		return
	}
	dir := *spellLifeWord(u.UpdateData, 12)
	pos := types.Ptf(float32(float64(*memmap.PtrInt32(0x587000, 196184+uintptr(dir)*8))*0.5+float64(u.PosVec.X)), float32(float64(*memmap.PtrInt32(0x587000, 196188+uintptr(dir)*8))*0.5+float64(u.PosVec.Y)))
	if spellEffectTrace(source.PosVec, pos, 0) {
		dword_5d4594_2487708 = memmap.PtrWord(u.CObj())
		*memmap.PtrFloat32(0x5d4594, 2487704) = float32(dy*dy + dx*dx)
	}
}
func spellEffectDoorPropagate(u, source *server.Object) {
	x := int32(23 * *spellLifeWord(u.UpdateData, 16))
	y := int32(23 * *spellLifeWord(u.UpdateData, 20))
	rect := types.Rectf{Min: types.Ptf(float32(float64(x)-34), float32(float64(y)-34)), Max: types.Ptf(float32(float64(x)+34), float32(float64(y)+34))}
	spellDoorLinkSource = source
	GetServer().S().Map.EachObjInRect(rect, func(u *server.Object) bool { spellEffectDoorLink(u); return true })
	spellDoorLinkSource = nil
}
func spellEffectLock(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	rect := types.Rectf{Min: types.Ptf(b.PosVec.X-150, b.PosVec.Y-150), Max: types.Ptf(b.PosVec.X+150, b.PosVec.Y+150)}
	dword_5d4594_2487708 = 0
	*memmap.PtrUint32(0x5d4594, 2487704) = 1287568416
	GetServer().S().Map.EachObjInRect(rect, func(u *server.Object) bool { spellEffectDoorCandidate(u, c); return true })
	u := (*server.Object)(memmap.WordPtr(dword_5d4594_2487708))
	if u == nil {
		return 0
	}
	if u.ObjOwner != nil && u.ObjOwner != b {
		resourcePriority(b, "ExecSpel.c:DoorAlreadyLocked")
		return 0
	}
	u.ObjOwner = b
	u.Field34 = GetServer().S().Frame() + 60*uint32(GetServer().S().TickRate())
	spellEffectDoorPropagate(u, b)
	spellEffectAudio(id, 0, u)
	return 1
}
