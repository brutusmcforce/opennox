package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

var orchestrationRestoreCleanup, orchestrationRewardMarker uint32

func orchestrationRoundFlags() {
	s := GetServer().S()
	count := s.Players.Count()
	for team := s.Teams.First(); team != nil; team = s.Teams.Next(team) {
		flag := (*server.Object)(team.Field_76)
		if flag.InvHolder != nil || teamRuntimeCount(team) == 0 {
			continue
		}
		for {
			index := s.Rand.Logic.IntClamp(0, count-1)
			pl := s.Players.First()
			for ; pl != nil && index != 0; pl = s.Players.Next(pl) {
				index--
			}
			if u := pl.PlayerUnit; u != nil && teamRuntimeContains(&u.TeamVal, team.ID()) {
				inventoryCrownPickup(u, flag, 1)
				break
			}
		}
	}
}

func orchestrationTransitionPlayers() {
	s := GetServer().S()
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		u := pl.PlayerUnit
		if u == nil {
			continue
		}
		if session := layout.As[server.PlayerUpdateData](u.UpdateData).Trade70; session != nil {
			shopCancel(session)
		}
		*(*unsafe.Pointer)(unsafe.Pointer(&layout.As[server.PlayerUpdateData](u.UpdateData).Trade70)) = nil
		s.Players.Control.Player(int(pl.PlayerInd)).Reset()
		if noxflags.HasGame(4096) {
			pl.Field4680 = pl.Field4676
			pl.Field4676 = 0
		}
		if !Nox_xxx_playerSetState_4FA020(u, 13) || !noxflags.HasGame(512) {
			pl.Field4700 = 0
			controlDefaultItems(u, int32(bool2int(!Nox_xxx_gameIsSwitchToSolo_4DB240())), int32(uint32(sessionMapState())>>1&1))
		}
		if pl.Field3680&0x20 != 0 {
			controlLeaveObserver((*server.Player)(pl.C()))
			Nox_xxx_playerCameraUnlock_4E6040(u)
		}
		u.Field34 = s.Frame()
	}
}

func orchestrationRestore(arg int32) {
	s := GetServer().S()
	pl := s.Players.ByInd(31)
	if pl == nil || pl.PlayerUnit == nil {
		return
	}
	host := pl.PlayerUnit
	if memmap.Uint32(0x5D4594, 1563128) == 0 {
		*memmap.PtrUint32(0x5D4594, 1563128) = uint32(s.Types.IndByID("SaveGameLocation"))
		*memmap.PtrUint32(0x5D4594, 1563132) = uint32(s.Types.IndByID("Glyph"))
	}
	if arg == 1 {
		for u := s.Objs.List; u != nil; {
			next := u.ObjNext
			if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 1563128) {
				Nox_xxx_unitMove_4E7010(pl.PlayerUnit, u.PosVec)
				pl.PlayerUnit.ScriptIDVal = u.ScriptIDVal
				u.ScriptIDVal = 0
				GetServer().DelayedDelete(u)
			} else {
				data := u.UpdateData
				switch {
				case u.ObjClass&2 != 0:
					// The loaded data holds script IDs in place of object pointers (see creatureXferAction).
					md := layout.As[server.MonsterUpdateData](data)
					scriptID := func(p **server.Object) uint32 { return *(*uint32)(unsafe.Pointer(p)) }
					for i := 0; i < int(md.Field282_1); i++ {
						ref := objectLookupByScriptID(scriptID(&md.Field283[i]))
						md.Field283[i] = ref // the seen slots hold object pointers
						if ref == nil {
							md.Field282_1 = 0
						}
					}
					md.CurrentEnemy = objectLookupByScriptID(scriptID(&md.CurrentEnemy))
					if int32(u.ObjFlags) >= 0 {
						creatureXferPostload(u)
					}
					for _, p := range []*uint32{&md.Field98, &md.Field300} {
						ref := objectLookupByScriptID(*p)
						*p = 0
						if ref != nil {
							*p = ref.NetCode
						}
					}
					md.Field304 = objectLookupByScriptID(scriptID(&md.Field304))
					for i := 0; i < int(md.Field543_0); i++ {
						p := &md.Field535[i]
						ref := objectLookupByScriptID(*p)
						if ref != nil {
							*p = ref.NetCode
						} else {
							md.Field543_0 = 0
						}
					}
				case u.ObjClass&0x4000 != 0:
					if layout.As[server.ElevatorUpdateData](data).Field_4 != 0 {
						u.NeedSync()
					}
				case u.ObjClass&0x8000 != 0:
					if target := layout.As[server.ElevatorShaftUpdateData](data).Partner; target != nil && layout.As[server.ElevatorUpdateData](target.UpdateData).Field_4 != 0 {
						u.NeedSync()
					}
				case u.ObjClass&0x80 != 0:
					if dd := u.UpdateDataDoor(); dd.Dir != dd.BaseDir {
						s.Objs.AddToUpdatable(u)
					}
				}
			}
			u = next
		}
		GetServer().NoxScriptC().ActResolveObjs()
		monsterPendingResolve()
		if orchestrationRestoreCleanup != 0 {
			for u := s.Objs.List; u != nil; {
				next := u.ObjNext
				if int32(u.ObjFlags) >= 0 && stateIsUnit(u) {
					if u.ObjClass&2 != 0 && u.ObjSubClass&0x2000 != 0 {
						for it := u.InvFirstItem; it != nil; {
							following := it.InvNextItem
							if uint32(it.TypeInd) == memmap.Uint32(0x5D4594, 1563132) {
								GetServer().DelayedDelete(it)
							}
							it = following
						}
					}
					GetServer().DelayedDelete(u)
				}
				u = next
			}
			for u := s.Objs.MissileList; u != nil; {
				next := u.ObjNext
				if int32(u.ObjFlags) >= 0 && stateIsPixie(u) {
					GetServer().DelayedDelete(u)
				}
				u = next
			}
		}
	}
	for u := host.Field129; u != nil; u = u.Field128 {
		if u.ObjClass&2 == 0 {
			continue
		}
		if *(*byte)(unsafe.Pointer(&layout.As[server.MonsterUpdateData](u.UpdateData).StatusFlags))&0x80 != 0 {
			gameplayReportAcquireCreature(int(pl.PlayerInd), u)
		} else if u.ObjSubClass&0x80 != 0 {
			gameplayReportMonitor(int(pl.PlayerInd), u)
		} else {
			continue
		}
		s.Players.Nox_xxx_netMarkMinimapObject_417190(pl.PlayerIndex(), u, 1)
	}
	Nox_xxx_gameSetSwitchSolo_4DB220(0)
	orchestrationRestoreCleanup = 0
	Nox_ticks_reset_416D40()
}

func orchestrationLoadState() int32 { return int32(memmap.Uint32(0x5D4594, 1563072)) }
func orchestrationDifficulty() {
	s := GetServer().S()
	if s.Frame()%(5*uint32(s.TickRate())) != 0 {
		return
	}
	before := questRuntimeFloat(202024)
	questRuntimeDifficulty()
	if questRuntimeFloat(202024) != before {
		questRuntimeScaleHealth()
	}
}
func orchestrationDropFlags() {
	s := GetServer().S()
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		for it := u.InvFirstItem; it != nil; {
			next := it.InvNextItem
			if it.ObjClass&0x10000000 != 0 {
				var pos types.Pointf
				inventoryRandomPlacement(50, &u.PosVec, &pos)
				inventoryDrop(u, it, &pos)
			}
			it = next
		}
	}
}
func orchestrationResetPlayer(u *server.Object) uint32 {
	data := u.UpdateDataPlayer()
	pl := data.Player
	data.Field29 = [4]*server.Object{}
	data.Field77 = nil
	pl.Ankhs = [5]*server.Object{}
	if !noxflags.HasGame(2048) {
		*(*byte)(unsafe.Pointer(&data.CurTraps)) = byte(controlGlyphCount(u))
	}
	data.Field66 = 0
	u.Obj130 = nil
	return memmap.PtrWord(mapPolygonActorInit(pl.C()))
}

func orchestrationRewards() {
	s := GetServer().S()
	stage := memmap.Uint32(0x587000, 202028)
	if orchestrationRewardMarker == 0 {
		orchestrationRewardMarker = uint32(s.Types.IndByID("RewardMarker"))
		*memmap.PtrUint32(0x5D4594, 1568304) = uint32(s.Types.IndByID("RewardMarkerPlus"))
	}
	rewardPlaceAnkh()
	rewardSelectMarkers()
	chestInit := lifecycleInitKey(initIDChest)
	marker := func(u *server.Object) bool {
		return uint32(u.TypeInd) == orchestrationRewardMarker || uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 1568304)
	}
	bow := func(u *server.Object) bool { return u.ObjClass&0x1000000 != 0 && u.ObjSubClass&0xc != 0 }
	for u := s.Objs.List; u != nil; {
		next := u.ObjNext
		if marker(u) {
			if u.InitDataRewardMarker().Flags&0x80 != 0 {
				if item := rewardMarker(u, stage); item != nil {
					GetServer().CreateObjectAt(item, nil, u.PosVec)
					if bow(item) {
						if quiver := s.NewObjectByTypeID("Quiver"); quiver != nil {
							pos := item.PosVec
							inventoryRandomPlacement(30, &item.PosVec, &pos)
							GetServer().CreateObjectAt(quiver, nil, pos)
						}
					}
				}
			}
			GetServer().DelayedDelete(u)
		} else if u.Init == chestInit {
			for it := u.InvFirstItem; it != nil; {
				following := it.InvNextItem
				if marker(it) {
					item := rewardMarker(it, stage+1)
					inventoryRemove(u, it)
					GetServer().DelayedDelete(it)
					if item != nil {
						inventoryInsert(u, item, 0)
						if bow(item) {
							if quiver := s.NewObjectByTypeID("Quiver"); quiver != nil {
								inventoryInsert(u, quiver, 0)
							}
						}
					}
				}
				it = following
			}
		}
		u = next
	}
}

func orchestrationWalls() {
	s := GetServer().S()
	scan := false
	for node := Get_dword_5d4594_251560(); node != nil; node = worldSecretNext(node) {
		sw := server.AsSecretWallData(node)
		phase, progress, flags := &sw.State, &sw.Progress, sw.Flags
		timer := &sw.Timer
		audio := func(open bool) {
			wall := worldSecretWall(node)
			def := s.Walls.DefByInd(int(wall.Tile1))
			name := def.CloseSound()
			if open {
				name = def.OpenSound()
			}
			x, y := int32(sw.GridX)*23+11, int32(sw.GridY)*23+11
			s.Audio.EventPos(sound.ByName(name), types.Pointf{X: float32(x), Y: float32(y)}, 0, 0)
		}
		switch *phase {
		case 1:
			scan = false
			if flags&12 == 12 {
				*timer--
				if *timer == 0 {
					*phase = 4
					audio(true)
				}
			}
		case 2:
			*progress--
			if *progress == 0 {
				*phase = 1
				*timer = uint32(s.TickRate()) * sw.OpenDelay
			}
			scan = true
		case 3:
			scan = false
			if flags&4 != 0 && flags&8 == 0 {
				*timer--
				if *timer == 0 {
					*phase = 2
					audio(false)
				}
			}
		case 4:
			*progress++
			if *progress == 23 {
				*phase = 3
				*timer = uint32(s.TickRate()) * sw.OpenDelay
			}
			scan = true
		}
		// Unknown phases deliberately preserve the previous wall's scan state.
		if scan {
			x, y := float64(int32(sw.GridX*23))+11.5, float64(int32(sw.GridY*23))+11.5
			rect := types.Rectf{Min: types.Pointf{X: float32(x - 42.5), Y: float32(y - 42.5)}, Max: types.Pointf{X: float32(x + 42.5), Y: float32(y + 42.5)}}
			s.Map.EachObjInRect(rect, func(u *server.Object) bool { collisionActivate(u); return true })
		}
	}
}
func orchestrationTimeout() bool {
	return memmap.Uint64(0x5D4594, 2523796) != 0 && memmap.Uint64(0x5D4594, 2523788)+5000 < PlatformTicks()
}
