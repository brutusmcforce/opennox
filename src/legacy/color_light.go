package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
)

// The animation stores the interpolation step in a byte, even for longer periods.
func colorLightInterpolation(d *client.Drawable, kind int) (first, next int, step, period float64, ok bool) {
	count := int(int8((*client.DrawableUnionColorLight)(unsafe.Pointer(&d.Union)).Counts[kind]))
	p := int(*effectShort(d, 258+2*kind))
	if count <= 1 || p == 0 {
		return
	}
	period = float64(p)
	cycle := float64(GetServer().S().Frame()) / float64(p*count)
	whole := int64(cycle)
	pos := (cycle - float64(int32(whole))) * float64(count)
	idx := int64(pos)
	first = int(int8(idx))
	next = first + 1
	if *(*byte)(unsafe.Pointer(&d.Field_44))&(1<<(2*kind)) != 0 {
		if next >= count {
			next = 0
		}
	} else if whole&1 != 0 {
		first = count - int(idx) - 1
		next = first - 1
		if next < 0 {
			next = 0
		}
	} else if next >= count {
		next = count - 1
	}
	step = float64(byte(int64((pos - float64(int32(idx))) * period)))
	ok = true
	return
}
func colorLightColor(d *client.Drawable) {
	a, b, step, period, ok := colorLightInterpolation(d, 0)
	if !ok {
		return
	}
	var rgb [3]int
	for i := range rgb {
		start := int(*effectByte(d, 178+3*a+i))
		end := int(*effectByte(d, 178+3*b+i))
		rgb[i] = int(int64(float64(end-start)/period*step + float64(start)))
	}
	particleLightColor(d, rgb[0], rgb[1], rgb[2])
}
func colorLightIntensity(d *client.Drawable) {
	a, b, step, period, ok := colorLightInterpolation(d, 1)
	if !ok {
		return
	}
	start := int(*effectByte(d, 226+a))
	end := int(*effectByte(d, 226+b))
	value := float32(step*(float64(end-start)/period) + float64(start))
	particleLightIntensity(d, value, true)
}
func colorLightPenumbra(d *client.Drawable) int {
	mode := d.Field_42
	if mode != 0 {
		return int(mode)
	}
	a, b, step, period, ok := colorLightInterpolation(d, 2)
	if !ok {
		return 0
	}
	start := int(*effectByte(d, 242+a))
	end := int(*effectByte(d, 242+b))
	value := int(int64(step*(float64(end-start)/period) + float64(start)))
	return int(particleLightAngle(d, value, true))
}
func colorLightTarget(d *client.Drawable) {
	if *(*byte)(unsafe.Pointer(&d.Field_44))&0x40 == 0 {
		return
	}
	target := GetClient().Cli().Objs.ByNetCodeStatic(int(d.Field_66))
	if target == nil {
		return
	}
	dx, dy := int32(target.PosVec.X-d.PosVec.X), int32(target.PosVec.Y-d.PosVec.Y)
	if dx == 0 && dy == 0 {
		return
	}
	angle := float64(float32(math.Acos(float64(dx)/math.Sqrt(float64(dx*dx+dy*dy))))) * 57.295776
	if dy < 0 {
		angle = 360 - angle
	}
	particleLightAngle(d, int(effectsTruncWord(angle)), false)
}
func colorLightRotate(d *client.Drawable) {
	if d.Field_42 == 0 {
		return
	}
	flags := *(*uint16)(unsafe.Pointer(&d.Field_44))
	speed := int16(*(*uint16)(unsafe.Add(unsafe.Pointer(&d.Field_67), 2)))
	if flags&0x80 == 0 || speed == 0 {
		return
	}
	start := float64(*(*uint16)(unsafe.Pointer(&d.Field_67)))
	span := 360.0
	if flags&0x100 != 0 {
		span = float64(int(*(*uint16)(unsafe.Pointer(&d.Field_68))) - int(*(*uint16)(unsafe.Pointer(&d.Field_67))))
	}
	if span == 0 {
		return
	}
	fps := float64(GetServer().S().TickRate())
	v := float64(speed)
	period := span / v * fps
	q := float64(GetServer().S().Frame()) / period
	value := (q-float64(effectsTruncWord(q)))*period*(v/fps) + start
	angle := int32(value)
	if math.IsNaN(value) || value >= 2147483648 || value < -2147483648 {
		angle = -2147483648
	}
	if angle >= 360 {
		angle -= 360
	} else if angle < 0 {
		angle += 360
	}
	particleLightAngle(d, int(angle), false)
}
func colorLightUpdate(vp *noxrender.Viewport, d *client.Drawable) int {
	if (*client.DrawableUnionColorLight)(unsafe.Pointer(&d.Union)).Counts == [3]byte{} {
		GetClient().Cli().Objs.List5Delete(d)
		return 1
	}
	x, y := d.PosVec.X, d.PosVec.Y
	if x >= vp.World.Min.X-100 && x <= vp.World.Min.X+vp.Size.X+100 && y >= vp.World.Min.Y-100 && y <= vp.World.Min.Y+vp.Size.Y+100 && !noxflags.HasGame(noxflags.GameFlag(0x200000)) {
		colorLightColor(d)
		colorLightIntensity(d)
		colorLightPenumbra(d)
		colorLightRotate(d)
		colorLightTarget(d)
	}
	return 1
}
