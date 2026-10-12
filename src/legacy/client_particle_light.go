package legacy

import (
	"math"
	"unsafe"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client"
)

func particleRGB(r, g, b int) uint32 {
	return noxcolor.RGB5551Color(byte(r), byte(g), byte(b)).Color32()
}
func particleLightWord(p unsafe.Pointer, off int) *uint32 { return (*uint32)(unsafe.Add(p, off)) }
func particleLightColor(p unsafe.Pointer, r, g, b int) unsafe.Pointer {
	*particleLightWord(p, 16) = uint32(r)
	*particleLightWord(p, 0) = 2
	*particleLightWord(p, 20) = uint32(g)
	*particleLightWord(p, 24) = uint32(b)
	return p
}
func particleLightAngle(p unsafe.Pointer, angle int, penumbra bool) int64 {
	value := int64(float64(angle)*0.0027777778*math.Float64frombits(uint64(qword_581450_9552)) + math.Float64frombits(uint64(qword_581450_9544)))
	off := 28
	if penumbra {
		off = 30
	} else {
		*particleLightWord(p, 32) = 0
	}
	*(*uint16)(unsafe.Add(p, off)) = uint16(value)
	return value
}
func particleLightIntensity(p unsafe.Pointer, intensity float32, fixed bool) int {
	if intensity > 63 {
		intensity = 63
	}
	*(*float32)(unsafe.Add(p, 4)) = intensity
	if fixed {
		value := int64(float64(intensity)*math.Float64frombits(uint64(qword_581450_9552)) + math.Float64frombits(uint64(qword_581450_9544)))
		*particleLightWord(p, 12) = uint32(value)
	}
	radius := client.LightRadius(intensity)
	*particleLightWord(p, 8) = uint32(radius)
	return radius
}
func initParticlePalettes() int {
	for i := 0; i < 64; i++ {
		// Preserve the source's signed-byte first palette, including its final step.
		v := int(byte(int8(-i) / 63))
		*effectMapped(1312500 + uintptr(4*i)) = particleRGB(v/3, 3*v/5, v)
	}
	for i := 0; i < 32; i++ {
		v := i * 255 / 32
		*effectMapped(1312756 + uintptr(4*i)) = particleRGB(v, v, 0)
		*effectMapped(1312884 + uintptr(4*i)) = particleRGB(255, -1-v, 0)
	}
	for i := 0; i < 64; i++ {
		v := i * 255 / 63
		*effectMapped(1313012 + uintptr(4*i)) = particleRGB(v, v, v)
		*effectMapped(1313268 + uintptr(4*i)) = particleRGB(v, 50, 50)
	}
	return int(*effectMapped(1313520))
}
func initParticleColors() int {
	for _, v := range [][4]int{
		{1313524, 255, 255, 0}, {1313528, 255, 100, 0}, {1313544, 0, 200, 200},
		{1313548, 50, 255, 255}, {1313552, 255, 0, 255}, {1313556, 255, 200, 255},
		{1313560, 255, 200, 0}, {1313568, 100, 255, 50}, {1313572, 150, 255, 150},
		{1313576, 255, 255, 0}, {1313580, 0, 220, 0}, {1313584, 150, 255, 150},
		{1313588, 200, 200, 200}, {1313592, 255, 255, 255},
	} {
		*effectMapped(uintptr(v[0])) = particleRGB(v[1], v[2], v[3])
	}
	dword_5d4594_1313532 = uint32(particleRGB(255, 255, 0))
	dword_5d4594_1313536 = uint32(particleRGB(0, 0, 255))
	dword_5d4594_1313540 = uint32(particleRGB(0, 200, 255))
	dword_5d4594_1313564 = uint32(particleRGB(255, 255, 100))
	red, green, blue := 255, 255, 255
	for i := 0; i < 16; i++ {
		if i <= 3 {
			blue = (765 - 255*i) / 3
		} else if i <= 7 {
			green = (600-200*i)/4 - 1
		} else {
			red = (1085-155*i)/9 - 1
		}
		*effectMapped(1313656 - uintptr(4*i)) = particleRGB(red, green, blue)
	}
	return initParticlePalettes()
}
