package legacy

import (
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"image"
	"unsafe"
)

func screenParticleHead() *Nox_screenParticle {
	return (*Nox_screenParticle)(unsafe.Pointer(legacyGlobals.nox_screenParticles_head))
}
func screenParticleTail() *Nox_screenParticle {
	return (*Nox_screenParticle)(unsafe.Pointer(legacyGlobals.dword_5d4594_806052))
}
func screenParticleAdd(p *Nox_screenParticle) {
	p.Field_44, p.Field_48 = screenParticleHead(), nil
	if p.Field_44 != nil {
		p.Field_44.Field_48 = p
	} else {
		legacyGlobals.dword_5d4594_806052 = p
	}
	legacyGlobals.nox_screenParticles_head = p
}
func screenParticleUnlink(p *Nox_screenParticle) {
	if p == screenParticleTail() {
		legacyGlobals.dword_5d4594_806052 = p.Field_48
	}
	if p.Field_44 != nil {
		p.Field_44.Field_48 = p.Field_48
	}
	if p.Field_48 != nil {
		p.Field_48.Field_44 = p.Field_44
	} else {
		legacyGlobals.nox_screenParticles_head = p.Field_44
	}
}
func screenParticleDelete(p *Nox_screenParticle) {
	screenParticleUnlink(p)
	alloc.AsClassT[Nox_screenParticle](legacyGlobals.nox_alloc_screenParticles_806044).FreeObjectFirst(p)
}
func screenParticleCreate(kind, x, y, vx, vy, gravity int, size, timer, phase, mode byte) *Nox_screenParticle {
	if legacyGlobals.nox_alloc_screenParticles_806044 == nil || kind < 0 || kind > 4 {
		return nil
	}
	colors := [5][2]uintptr{{806016, 806036}, {806028, 806004}, {806032, 806040}, {806020, 806012}, {806008, 806024}}
	glow, core := *effectMapped(colors[kind][0]), *effectMapped(colors[kind][1])
	p := alloc.AsClassT[Nox_screenParticle](legacyGlobals.nox_alloc_screenParticles_806044).NewObject()
	if p == nil {
		p = screenParticleTail()
		if p == nil {
			return nil
		}
		screenParticleUnlink(p)
	}
	// Tail reuse preserves unused bytes; only the low byte of Field_32 is assigned.
	p.Field_24, p.Field_28 = uint32(x)<<16, uint32(y)<<16
	p.Field_40 = [4]byte{size, timer, phase, timer}
	p.Draw_fnc = screenParticleCallbackKey()
	p.Field_16, p.Field_20, p.Field_36 = uint32(vx)<<16, uint32(vy)<<16, uint32(gravity)<<16
	p.Field_32 = p.Field_32&0xffffff00 | uint32(mode)
	p.Field_4, p.Field_8, p.Field_12 = uint32(kind), glow, core
	screenParticleAdd(p)
	if p.Field_36 == 0 && p.Field_40[1] == 0 {
		p.Field_40[1], p.Field_40[2], p.Field_40[3] = 3, 2, 3
	}
	return p
}
func screenParticleDraw(vp *noxrender.Viewport, p *Nox_screenParticle) int {
	pos := image.Pt(int(p.Field_24>>16), int(p.Field_28>>16))
	if pos.X <= 0 || pos.Y <= 0 || pos.X >= vp.Size.X || pos.Y >= vp.Size.Y {
		screenParticleDelete(p)
		return 0
	}
	effectGlow(pos, p.Field_8, int(p.Field_40[0]), int(p.Field_40[0]))
	effectColor(p.Field_12)
	effectPoint(pos, int(p.Field_40[0]>>1))
	if !nox_client_drawIsTick {
		// particles only move on game ticks
		return 1
	}
	p.Field_20 += p.Field_36
	if p.Field_40[1] != 0 {
		p.Field_40[1]--
		if p.Field_40[1] == 0 {
			if p.Field_40[2] == 1 {
				p.Field_40[0]++
				if p.Field_40[0] >= 4 {
					p.Field_40[2] = 2
				}
			} else {
				p.Field_40[0]--
				if p.Field_40[0] == 0 {
					screenParticleDelete(p)
					return 0
				}
			}
			p.Field_40[1] = p.Field_40[3]
		}
	}
	if byte(p.Field_32) == 1 && effectRand(0, 10) >= 8 {
		timer := byte(effectRand(3, 5))
		size := byte(effectRand(2, 3))
		dx := effectRand(-2, 2)
		screenParticleCreate(int(p.Field_4), pos.X+dx, pos.Y, 0, 0, 1, size, timer, 2, 2)
	}
	p.Field_24 += p.Field_16
	p.Field_28 += p.Field_20
	return 1
}
func screenParticlesDraw(vp *noxrender.Viewport) {
	if vp == nil {
		return
	}
	Sub_430B50(vp.Screen.Min.X, vp.Screen.Min.Y, vp.Screen.Max.X, vp.Screen.Max.Y)
	for p := screenParticleHead(); p != nil; {
		GetClient().Cli().GUI.ValYYY = 1
		// The callback can delete, recycle, or create nodes, including the current one.
		next := p.Field_44
		callScreenParticleDraw(p.Draw_fnc, vp, p)
		p = next
	}
}
