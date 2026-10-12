package legacy

import (
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"image"
	"math"
	"unsafe"
)

func objectRenderGhost(vp *noxrender.Viewport, dr *client.Drawable) byte {
	local := LocalPlayerDrawable
	if local != nil && local.HasEnchant(21) {
		return 255
	}
	x := uint32(dr.PosVec.X) - uint32(vp.Size.X)/2 - uint32(vp.World.Min.X)
	y := uint32(dr.PosVec.Y) - uint32(vp.Size.Y)/2 - uint32(vp.World.Min.Y)
	value := int32(uint32(128) - ((x*x+y*y)<<7)/memmap.Uint32(0x587000, 185464))
	if value < 0 {
		return 0
	}
	if value > 128 {
		return 128
	}
	return byte(value)
}
func objectRenderShiny(vp *noxrender.Viewport, dr *client.Drawable) uint16 {
	if objectRenderShinyAnim == nil {
		objectRenderShinyAnim = Nox_xxx_gLoadAnim("ShinySpot")
	}
	ref := objectRenderShinyAnim
	anim := (*ImageRefAnim)(ref.Field_24)
	count := uint32(anim.ImagesSz)
	frame := gameFrame() + dr.NetCode32
	period := frame / (8 * count)
	index := (frame % (8 * count)) >> 1
	if index < count {
		pos := effectScreen(vp, dr).Sub(image.Pt(64, 64))
		r := GetClient().R2()
		r.DrawImageAt(r.GetBag().AsImage(anim.Images()[index]), pos)
	}
	return uint16(period)
}
func objectRenderTint(color uint32) {
	r := GetClient().R2()
	r.Data().SetColorize17(1)
	c := noxcolor.RGBA5551(color).ColorNRGBA()
	r.SetColorMultAndIntensityRGB(c.R, c.G, c.B)
}
func objectRenderTeam(id int) *server.ObjectTeam {
	// The shared lookup still serves remaining C code and owns host/client routing.
	return teamRuntimeObject(id)
}
func objectRenderDraw(vp *noxrender.Viewport, dr *client.Drawable, img noxrender.ImageHandle) {
	if dword_5d4594_1321520 == 0 {
		dword_5d4594_1321520 = uint32(effectType("Ghost"))
	}
	sameTeam, targetObserver := false, false
	if dr.ObjClass&4 != 0 {
		player := GetServer().S().Players.ByID(int(dr.NetCode32))
		if uint32(ClientPlayerNetCode()) == dr.NetCode32 {
			targetObserver = player != nil && player.Field3680&1 != 0
		} else if player != nil && player.Field3680&1 != 0 {
			return
		}
		if localTeam := objectRenderTeam(ClientPlayerNetCode()); localTeam != nil {
			if team := objectRenderTeam(int(dr.NetCode32)); team != nil {
				sameTeam = uint32(ClientPlayerNetCode()) == dr.NetCode32 || localTeam.SameAs(team)
			}
		}
	}
	pos := effectScreen(vp, dr).Sub(image.Pt(int(*effectByte(dr, 0)), int(*effectByte(dr, 1))))
	if memmap.Uint32(0x587000, 80808) != 0 && dr.ObjFlags&0x40000001 == 0 && dr.ObjClass&0x80 != 0 {
		a := vp.ToScreenPos(dr.PosVec)
		dir := int(*effectByte(dr, int(layoutClientDrawable.Off(299)) /* Field_74_4 */))
		b := a.Add(image.Pt(int(memmap.Int32(0x587000, 196184+8*uintptr(dir))), int(memmap.Int32(0x587000, 196188+8*uintptr(dir)))))
		mid := image.Pt((a.X+b.X)>>1, (a.Y+b.Y)>>1)
		if dir < 24 && dir != 0 && (dir < 8 || dir > 16) {
			if GetClient().Sub4C5630(mid.X-5, mid.X-5, mid.Y) != 0 {
				a.X -= 2
				b.X -= 2
			} else if GetClient().Sub4C5630(mid.X+5, mid.X+5, mid.Y) != 0 {
				a.X += 2
				b.X += 2
			}
		} else {
			if GetClient().Sub4C5630(mid.X, mid.X, mid.Y-5) != 0 {
				a.Y -= 2
				b.Y -= 2
			} else if GetClient().Sub4C5630(mid.X, mid.X, mid.Y+5) != 0 {
				a.Y += 2
				b.Y += 2
			}
		}
		// Preserve the original acceptance decision. Returned horizontal bounds were
		// not applied by C, including accepted short paths with empty scanline spans.
		minX, maxX := 0, int(nox_win_width)
		if GetClient().Sub4C42A0(a, b, &minX, &maxX) == 0 {
			return
		}
	}
	var light []uint32
	if dr.ObjClass&0x80000 != 0 && dr.ObjFlags&0x1000000 != 0 || dr.ObjFlags&0x40000000 != 0 {
		light = unsafe.Slice(memmap.PtrUint32(0x587000, 185472), 3)
	} else {
		light = GetClient().Sub469920(dr.PosVec)
	}
	crop := 0
	repeatBelow := false
	z := int(int16(dr.ZVal))
	if z < 0 {
		crop = -z
		repeatBelow = int32(*effectMapped(1321512)) < 0 && objectRenderLastBelow == dr
	}
	objectRenderLastBelow = dr
	*effectMapped(1321512) = uint32(z)
	r := GetClient().R2()
	d := r.Data()
	if dr.HasEnchant(25) || dr.ObjClass&2 != 0 && dr.ObjFlags&0x40000000 != 0 && dr.ObjFlags&0x8020 == 0 {
		objectRenderTint(uint32(nox_color_blue_2650684))
	} else {
		d.SetMultiply14(1)
		r.SetColorMultAndIntensityRGB(byte(light[0]), byte(light[1]), byte(light[2]))
	}
	ghost := dr.TypeIDVal == uint32(dword_5d4594_1321520)
	alpha := func(v byte) { d.SetAlphaEnabled(true); d.SetAlpha(v) }
	if dr.Field_120 != 0 {
		factor := 1 - float64(gameFrame()-dr.Field_85)/float64(int32(gameFPS()))
		if factor < 0 {
			factor = 0.001
		}
		value := byte(int64(factor * 255))
		if ghost {
			if g := objectRenderGhost(vp, dr); g < value {
				value = g
			}
		} else if dr.ObjFlags&0x4000000 != 0 {
			value >>= 1
		}
		alpha(value)
	} else if ghost {
		alpha(objectRenderGhost(vp, dr))
	} else if dr.ObjFlags&0x4000000 != 0 {
		alpha(128)
	}
	local := LocalPlayerDrawable
	see := local != nil && local.HasEnchant(21)
	if dr.HasEnchant(0) || targetObserver || ghost && see {
		if pl := Get_dword_8531A0_2576(); pl != nil && pl.Field3680&1 != 0 {
			alpha(128)
		} else {
			dx, dy := dr.PosVec.X-int(dr.Field_8), dr.PosVec.Y-int(dr.Field_9)
			ax, ay := dx, dy
			if ax < 0 {
				ax = -ax
			}
			if ay < 0 {
				ay = -ay
			}
			// This return deliberately preserves already-written cache/renderer state.
			if ax < 4 && ay < 4 && local != nil && !see && !sameTeam {
				return
			}
			speed := dx
			if dx != 0 || dy != 0 {
				elapsed := int(dr.Field_5 - dr.Field_10)
				if elapsed == 0 {
					elapsed = 1
				}
				speed = int(math.Sqrt(float64(dx*dx+dy*dy))) / elapsed
			}
			if local != nil && !targetObserver && !sameTeam && see {
				objectRenderTint(uint32(dword_8531A0_2572))
				alpha(255)
			} else {
				value := byte(128)
				if !(local != nil && targetObserver) && speed < 8 {
					value = byte((speed << 7) / 8)
					if value == 0 {
						value = 1
					}
				}
				if sameTeam && value <= 1 && (dr.ObjClass&2 != 0 && dr.AnimInd == 8 || dr.ObjClass&4 != 0 && dr.AnimInd == 0) {
					objectRenderTint(uint32(dword_8531A0_2572))
					value = 128
				}
				alpha(value)
			}
		}
	}
	if dr.ObjClass&4 == 0 && dr.HasEnchant(23) && !noxflags.HasGame(noxflags.GameFlag(2048)) {
		color := uint32(nox_color_blue_2650684)
		if byte(gameFrame())&1 != 0 {
			color = uint32(nox_color_white_2523948)
		}
		objectRenderTint(color)
	}
	if !dr.HasEnchant(23) && dr.HasEnchant(11) {
		objectRenderTint(memmap.Uint32(0x85B3FC, 956))
	}
	if repeatBelow {
		objectRenderSaveClip()
		r.Sub_49F7C0_def_go()
	} else {
		r.Set_dword_5d4594_3799484(crop)
	}
	r.DrawImageAt(r.GetBag().AsImage(img), pos)
	d.SetMultiply14(0)
	d.SetAlphaEnabled(false)
	d.SetColorize17(0)
	if repeatBelow {
		objectRenderRestoreClip()
	}
	*effectShort(dr, 2) = 1
	*effectShort(dr, 4) = *memmap.PtrUint16(0x973F18, 88)
	*effectShort(dr, 6) = *memmap.PtrUint16(0x973F18, 76)
	dr.Field_2 = unsafe.Pointer(img)
}
