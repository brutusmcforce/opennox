package legacy

import (
	"image"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
)

// Particle initialization consumes randomness only after allocation succeeds.
// Creation callers then set height/vertical velocity before joining list 3/4.
func effectInitSpark(dr *client.Drawable, pos image.Point, speed, ttlMin, ttlMax int) {
	u := (*client.DrawableUnionSpark)(unsafe.Pointer(&dr.Union))
	u.X = uint32(pos.X) << 12
	u.Y = uint32(pos.Y) << 12
	dr.Field_74_4 = byte(effectRand(0, 255))
	u.Speed = uint32(effectRand(1, speed))
	u.Expires = gameFrame() + uint32(effectRand(ttlMin, ttlMax))
	u.Born = gameFrame()
}
func effectCreateOrb(typ int, coords *[4]uint16, x, y int, speed, period byte) {
	dr := effectSpawn(typ, image.Pt(x+int(coords[2]), y+int(coords[3])))
	if dr == nil {
		return
	}
	u := (*client.DrawableUnionOrb)(unsafe.Pointer(&dr.Union))
	u.Target[0] = coords[0]
	u.Target[1] = coords[1]
	u.Speed = speed
	u.Size = byte(effectRand(3, 10))
	u.Timer = period
	u.Period = period
	effectLink(dr)
}
func effectCreateOrbit(typ int, coords *[4]int16, angle int16, direction, period byte) {
	pos := image.Pt(int(coords[2]), int(coords[3]))
	dr := effectSpawn(typ, pos)
	if dr == nil {
		return
	}
	u := (*client.DrawableUnionOrb)(unsafe.Pointer(&dr.Union))
	u.Target = [2]uint16{uint16(coords[0]), uint16(coords[1])}
	u.Start = [2]uint16{uint16(coords[2]), uint16(coords[3])}
	u.Angle = byte(angle)
	dx, dy := pos.X-int(coords[0]), pos.Y-int(coords[1])
	u.Radius = uint16(int64(math.Sqrt(float64(dx*dx + dy*dy))))
	u.Speed = direction
	u.Size = byte(effectRand(3, 10))
	u.Timer = period
	u.Period = period
	dr.ClientUpdateFuncPtr = drawableUpdateIdentity(updateID_sub_4CA720)
	*(*uint16)(unsafe.Pointer(&dr.Field_127)) = uint16(angle)
	effectLink(dr)
}
func effectCreatePointSparks(typ, count, speed, ttl, x, y int) int {
	for count > 0 {
		if dr := effectSpawn(typ, image.Pt(x, y)); dr != nil {
			effectInitSpark(dr, dr.PosVec, speed, ttl, 64)
			dr.ZVal = 0
			dr.VelZ = int8(effectRand(2, 10))
			effectLink(dr)
		}
		count--
	}
	return count
}
func effectCreateEnergySparks(x, y int, z int16, typ int) int {
	for i := 0; i < 2; i++ {
		if dr := effectSpawn(typ, image.Pt(x, y)); dr != nil {
			effectInitSpark(dr, image.Pt(x, y), 3000, 5, 20)
			dr.ZVal = uint16(int(z) + effectRand(0, 20))
			dr.VelZ = int8(effectRand(0, 4))
			effectLink(dr)
		}
	}
	return 0
}
func effectCreateRainOrb(typ int, from, to image.Point, z uint16, velocity int8) *client.Drawable {
	dr := effectSpawn(typ, from)
	if dr == nil {
		return nil
	}
	dr.ZVal = z
	dr.ZVal2 = 0
	dr.VelZ = velocity
	u := (*client.DrawableUnionRainOrb)(unsafe.Pointer(&dr.Union))
	u.PrevZ = z
	u.Size = byte(effectRand(3, 10))
	u.ToX = uint32(to.X)
	u.ToY = uint32(to.Y)
	effectLink(dr)
	return dr
}
func effectLightningParticles(typ int, from, to image.Point) int {
	dx, dy := to.X-from.X, to.Y-from.Y
	distance := int(int64(math.Sqrt(float64(dx*dx + dy*dy))))
	result := distance
	if distance > 0 {
		result = effectRand(0, distance)
		for offset := result; offset <= distance; {
			pos := image.Pt(from.X+dx*offset/distance, from.Y+dy*offset/distance)
			if dr := effectSpawn(typ, pos); dr != nil {
				effectInitSpark(dr, pos, 3000, 5, 20)
				dr.ZVal = uint16(effectRand(15, 30))
				dr.VelZ = int8(effectRand(-4, 4))
				effectLink(dr)
			}
			result = effectRand(8, 100)
			offset += result
		}
	}
	// The historical return is the final spacing sample, not the emission count.
	return result
}
func effectScreenParticles(kind, x, y, width, height, axis, direction int) int {
	for i := 0; i < 100; i++ {
		px, py := x, y
		if axis == 2 {
			px += effectRand(0, width)
		} else {
			py += effectRand(0, height)
		}
		radius := effectRand(2, 5)
		vy := effectRand(-40, -20)
		var vx int
		if direction == 1 {
			vx = effectRand(-20, 0)
		} else {
			vx = effectRand(0, 20)
		}
		screenParticleCreate(kind, px, py, vx, vy, 1, byte(radius), 0, 0, 1)
	}
	return 0
}
func effectSparkBurst(pos image.Point, typ int, amount byte) int {
	speed := 2400*int(amount)/255 + 200
	ttl := 10*int(amount)/255 + 5
	for count := 180*int(amount)/255 + 10; count > 0; count-- {
		if dr := effectSpawn(typ, pos); dr != nil {
			effectInitSpark(dr, dr.PosVec, speed, ttl, 96)
			dr.ZVal = uint16(effectRand(5, 15))
			dr.ZVal2 = 0
			dr.VelZ = int8(effectRand(0, 8))
			effectLink(dr)
		}
	}
	return 0
}
