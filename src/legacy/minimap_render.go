package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/object"
	"github.com/brutusmcforce/opennox/server"
)

// Current object in the optional debug overlay traversal; cleared at its end.
var minimapDebugIterator uint32

func minimapLocalDrawable() *client.Drawable { return LocalPlayerDrawable }
func minimapLevel(dr *client.Drawable) int {
	var p image.Point
	if dr == minimapLocalDrawable() {
		p = GetClient().Viewport().World.Max
	} else {
		p = dr.PosVec
	}
	index := int(mapPolygonIndex(&[2]int32{int32(p.X), int32(p.Y)}, memmap.Uint32(0x5D4594, 1096312)))
	if index != 0 {
		*memmap.PtrUint32(0x5D4594, 1096312) = uint32(index)
	} else {
		index = int(memmap.Uint32(0x5D4594, 1096312))
	}
	if index == 0 {
		return 1
	}
	return int(mapPolygonGet(uint32(index)).Level)
}
func minimapPolygonLevel(p image.Point) (int, bool) {
	poly := mapPolygonFind(&[2]int32{int32(p.X), int32(p.Y)}, 0, false)
	if poly == nil {
		return 0, false
	}
	return int(poly.Level), true
}
func minimapDrawSprite(dr *client.Drawable) {
	if local := minimapLocalDrawable(); local != nil && local.Buffs&(1<<2) != 0 {
		return
	}
	Sub_437260()
	level := minimapLevel(dr)
	*memmap.PtrUint32(0x5D4594, 1096316) = uint32(level)
	minimapDraw(dr, level)
	GetClient().R2().SetRectFullScreen()
}
func minimapDrawAndMessages() int {
	if nox_client_gui_flag_1556112 == 1 {
		return 1
	}
	if memmap.Uint8(0x5D4594, 1096424)&1 != 0 {
		minimapDrawSprite(GetClient().Cli().Objs.ByNetCodeDynamic(ClientPlayerNetCode()))
	}
	return int(interactionMessagesDraw())
}
func minimapDraw(_ *client.Drawable, level int) int {
	c, s := GetClient(), GetServer().S()
	r := c.R2()
	d := r.Data()
	if memmap.Uint8(0x5D4594, 1096300) == 0 {
		*memmap.PtrUint8(0x5D4594, 1096300) = byte(s.Walls.DefIndByName("InvisibleWallSet"))
		*memmap.PtrUint8(0x5D4594, 1096301) = byte(s.Walls.DefIndByName("InvisibleBlockingWallSet"))
	}
	d.SetAlphaEnabled(false)
	objectRenderSaveClip()
	width, height := int(nox_win_width), int(nox_win_height)
	size := width / 6
	top := (height - size) / 2
	uiRenderCopyRect(0, 0, width, height)
	left := drawViewport().Screen.Min.X
	if left == 0 {
		r.DrawRectFilledAlpha(0, top, size, size)
	} else {
		d.SetColor2(noxcolor.RGBA5551(nox_color_black_2650656))
		if left >= size {
			r.DrawRectFilledOpaque(0, top, size, size, d.Color2())
		} else {
			r.DrawRectFilledOpaque(0, top, left, size, d.Color2())
			r.DrawRectFilledAlpha(left, top, size-left, size)
		}
	}
	d.SetAlphaEnabled(true)
	d.SetColor2(noxcolor.RGBA5551(nox_color_black_2650656))
	for i, a := range []byte{90, 60, 40} {
		n := i + 1
		d.SetAlpha(a)
		minimapOutline(-n, top-n, size+2*n, size+2*n)
	}
	d.SetAlphaEnabled(false)
	uiRenderCopyRect(0, top, size, size)
	zoom := int(minimapZoom)
	span := size * zoom / 100
	origin := drawViewport().World.Max.Sub(image.Pt(span/2, span/2))
	project := func(p image.Point) image.Point { return image.Pt(100*(p.X-origin.X)/zoom, top+100*(p.Y-origin.Y)/zoom) }
	xMin, xMax := origin.X/23, (origin.X+span)/23
	yMin, yMax := origin.Y/23, (origin.Y+span)/23
	for y := yMin; y <= yMax; y++ {
		for w := s.Walls.IndexByY(y); w != nil; w = w.NextByY24 {
			if w.Tile1 == memmap.Uint8(0x5D4594, 1096300) || w.Tile1 == memmap.Uint8(0x5D4594, 1096301) {
				continue
			}
			if byte(w.Field8) != 0 && int(byte(w.Field8)) != level {
				continue
			}
			x := int(w.X5)
			if x < xMin {
				continue
			}
			if x > xMax {
				break
			}
			if w.Flags4&0x10 != 0 {
				if w.Field32 == 0 {
					continue
				}
				revealed := false
				for i := uintptr(0); i < 4; i++ {
					dx, dy := int(memmap.Int32(0x587000, 149240+8*i)), int(memmap.Int32(0x587000, 149244+8*i))
					adj := s.Walls.GetWallAtGrid(image.Pt(int(w.X5)+dx, int(w.Y6)+dy))
					if adj != nil {
						revealed = adj.Field12 != 0
					} else {
						a := s.Walls.GetWallAtGrid(image.Pt(int(w.X5)+dx, int(w.Y6)))
						b := s.Walls.GetWallAtGrid(image.Pt(int(w.X5), int(w.Y6)+dy))
						revealed = (a != nil && a.Field12 != 0) || (b != nil && b.Field12 != 0)
					}
					if revealed {
						break
					}
				}
				if revealed {
					door := (*client.Drawable)(memmap.WordPtr(w.Field32))
					p := project(door.PosVec)
					off := uintptr(8 * int(door.Field_74_4))
					dx, dy := 100*int(memmap.Int32(0x587000, 196184+off))/zoom, 100*int(memmap.Int32(0x587000, 196188+off))/zoom
					d.SetColor2(noxcolor.RGBA5551(memmap.Uint32(0x85B3FC, 940)))
					r.AddPoint(p)
					r.AddPointRel(image.Pt(dx, dy))
					r.DrawLineFromPoints(d.Color2())
				}
				continue
			}
			if !noxflags.HasGame(0x10000) && w.Field12 == 0 {
				continue
			}
			if w.Flags4&4 != 0 {
				v := w.SecretData().State
				if v == server.SecretWallClosing || v == server.SecretWallOpen {
					continue
				}
			}
			if w.Flags4&0x20 == 0 {
				minimapWall(project(image.Pt(23*x, 23*y)), w.Dir0, 2300/zoom)
			}
		}
	}
	if noxflags.HasEngine(noxflags.EngineShowAI) {
		points := s.AI.Paths.Points()
		if len(points) >= 2 {
			d.SetColor2(noxcolor.RGBA5551(dword_8531A0_2572))
			for i := 1; i < len(points); i++ {
				a, b := points[i-1], points[i]
				minimapLine(project(image.Pt(int(int64(a.X)), int(int64(a.Y)))), project(image.Pt(int(int64(b.X)), int(int64(b.Y)))))
			}
		}
		for u := s.Objs.First(); u != nil; u = u.ObjNext {
			minimapDebugIterator = memmap.PtrWord(unsafe.Pointer(u))
			if u.ObjClass&object.ClassMonster == 0 {
				continue
			}
			d.SetColor2(noxcolor.RGBA5551(memmap.Uint32(0x85B3FC, 940)))
			p := project(image.Pt(int(int64(u.PosVec.X)), int(int64(u.PosVec.Y))))
			minimapPoint(p.X, p.Y, false)
		}
		minimapDebugIterator = 0
	}
	minimapDrawObjects(level, project)
	return objectRenderRestoreClip()
}

func minimapDrawObjects(level int, project func(image.Point) image.Point) {
	c, s := GetClient(), GetServer().S()
	d := c.R2().Data()
	if memmap.Uint32(0x5D4594, 1096304) == 0 {
		*memmap.PtrUint32(0x5D4594, 1096304) = uint32(c.Cli().Things.IndByID("Crown"))
		*memmap.PtrUint32(0x5D4594, 1096308) = uint32(c.Cli().Things.IndByID("GameBall"))
	}
	localTeam := objectRenderTeam(ClientPlayerNetCode())
	local := minimapLocalDrawable()
	teamColor := func(t *server.Team) { d.SetColor2(noxcolor.ToRGBA5551Color(s.Teams.GetTeamColor(t))) }
	for dr := c.Cli().Objs.FirstMinimapList(); dr != nil; dr = dr.Field_102 {
		floor, ok := minimapPolygonLevel(dr.PosVec)
		if (ok && floor != level) || (!ok && level != 1) {
			continue
		}
		p := project(dr.PosVec)
		color := memmap.Uint32(0x85B3FC, 940)
		if dr.ObjClass&0x400000 != 0 && dr.ObjSubClass&8 != 0 {
			color = uint32(nox_color_blue_2650684)
		}
		d.SetColor2(noxcolor.RGBA5551(color))
		if dr.TypeIDVal == memmap.Uint32(0x5D4594, 1096304) {
			held := false
			if s.Teams.First() == nil {
				for pl := c.Cli().Objs.FirstPlayerList(); pl != nil; pl = pl.Field_104 {
					if pl.Buffs&(1<<30) != 0 {
						held = true
						break
					}
				}
			}
			if held {
				continue
			}
			d.SetColor2(noxcolor.RGBA5551(dword_8531A0_2572))
			if team := objectRenderTeam(int(dr.NetCode32)); team != nil {
				if t := s.Teams.ByID(team.ID); t != nil {
					teamColor(t)
				}
			}
			minimapCrown(p)
			continue
		}
		if dr.TypeIDVal == memmap.Uint32(0x5D4594, 1096308) {
			if team := objectRenderTeam(int(dr.NetCode32)); team.Has() {
				if t := s.Teams.ByID(team.ID); t != nil {
					teamColor(t)
				}
			} else {
				d.SetColor2(noxcolor.RGBA5551(nox_color_white_2523948))
			}
			minimapCircle(p)
			continue
		}
		if dr.ObjClass&0x10000000 != 0 {
			if dr.ObjFlags&0x1000000 != 0 {
				d.SetColor2(noxcolor.RGBA5551(nox_color_white_2523948))
				if t := s.Teams.ByID(server.TeamID(objectDrawableTeamColor(dr))); t != nil {
					teamColor(t)
				}
				minimapFlag(p)
			}
		} else if dr.ObjClass&4 == 0 {
			minimapPoint(p.X, p.Y, false)
		} else if !noxflags.HasGame(32) {
			minimapPoint(p.X, p.Y, dr == local)
		} else if pl := s.Players.ByID(int(dr.NetCode32)); pl != nil && pl.WeaponEquip&1 != 0 {
			d.SetColor2(noxcolor.RGBA5551(nox_color_white_2523948))
			if team := objectRenderTeam(int(dr.NetCode32)); team != nil {
				id := server.TeamID(1)
				if team.ID == 1 {
					id = 2
				}
				if t := s.Teams.ByID(id); t != nil {
					teamColor(t)
				}
			}
			minimapFlag(p)
		}
	}
	for dr := c.Cli().Objs.FirstPlayerList(); dr != nil; dr = dr.Field_104 {
		crown := dr.Buffs&(1<<30) != 0
		team := objectRenderTeam(int(dr.NetCode32))
		pl := s.Players.ByID(int(dr.NetCode32))
		same := localTeam.Has() && localTeam.SameAs(team)
		if !same && !crown && dr != local && (*server.Player)(uiMeterPlayer()).Field3680&1 == 0 {
			continue
		}
		floor, ok := minimapPolygonLevel(dr.PosVec)
		if (ok && floor != level) || pl == nil || pl.Field3680&1 != 0 {
			continue
		}
		p := project(dr.PosVec)
		color := memmap.Uint32(0x85B3FC, 940)
		if dr == local || same {
			color = uint32(dword_8531A0_2572)
		}
		d.SetColor2(noxcolor.RGBA5551(color))
		if crown {
			if team != nil {
				if t := s.Teams.ByID(team.ID); t != nil {
					teamColor(t)
				}
			}
			minimapCrown(p)
		} else {
			minimapPoint(p.X, p.Y, false)
		}
	}
}
