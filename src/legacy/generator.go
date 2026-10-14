package legacy

import (
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func generatorOccupied(p types.Pointf) uint32 {
	occupied := memmap.PtrUint32(0x5D4594, 2491708)
	*occupied = 0
	GetServer().S().Map.EachObjInRect(types.Rectf{Min: types.Pointf{X: p.X - 15, Y: p.Y - 15}, Max: types.Pointf{X: p.X + 15, Y: p.Y + 15}}, func(t *server.Object) bool {
		if uint32(t.ObjClass)&0x20000 == 0 && t.Field5&0x800 == 0 && collisionObjectContains(t, &p) == 1 {
			*occupied = 1
		}
		return true
	})
	return *occupied
}
func generatorRadial(radius float32, origin types.Pointf, out *types.Pointf, descriptor *server.Object) int {
	core := GetServer().S()
	angle := float32(core.Rand.Logic.FloatClamp(float64(float32(-3.1415927)), float64(float32(3.1415927))))
	flags := server.MapTraceFlags(1)
	if descriptor.ObjFlags&0x4000 != 0 {
		flags = 5
	}
	for i := 0; i < 32; i++ {
		// C keeps the double sum for cos, but passes the stored float to sin.
		next := float64(angle) + 1.8849558
		angle = float32(next)
		p := types.Pointf{X: float32(placementCos(next)*float64(radius) + float64(origin.X)), Y: float32(placementSin(float64(angle))*float64(radius) + float64(origin.Y))}
		if core.MapTraceRay(origin, p, flags) && generatorOccupied(p) == 0 && !worldTileWater(p) {
			*out = p
			return 1
		}
	}
	return 0
}
func generatorPlace(u *server.Object, out *types.Pointf, player, descriptor *server.Object) int {
	core := GetServer().S()
	if u.UpdateDataMonsterGen().Flags92&2 != 0 && player != nil {
		p := types.Pointf{X: player.PosVec.X - u.PosVec.X, Y: player.PosVec.Y - u.PosVec.Y}
		if p.X == 0 {
			p.X++
		}
		if p.Y == 0 {
			p.Y++
		}
		geometryNormalize((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&p))))
		p.X = float32(float64(p.X)*45 + float64(u.PosVec.X))
		p.Y = float32(float64(p.Y)*45 + float64(u.PosVec.Y))
		flags := server.MapTraceFlags(1)
		if descriptor.ObjFlags&0x4000 != 0 {
			flags = 5
		}
		if generatorOccupied(p) == 0 && core.MapTraceRay(u.PosVec, p, flags) && !worldTileWater(p) {
			*out = p
			d := types.Pointf{X: player.PosVec.X - p.X, Y: player.PosVec.Y - p.Y}
			direction := server.Dir16(Nox_xxx_math_509ED0(d))
			descriptor.Direction1, descriptor.Direction2 = direction, direction
			return 1
		}
	}
	if generatorRadial(45, u.PosVec, out, descriptor) == 1 {
		if player != nil {
			d := types.Pointf{X: player.PosVec.X - u.PosVec.X, Y: player.PosVec.Y - u.PosVec.Y}
			direction := server.Dir16(Nox_xxx_math_509ED0(d))
			descriptor.Direction1, descriptor.Direction2 = direction, direction
		}
		return 1
	}
	return 0
}
func generatorPick(u *server.Object, out *types.Pointf, descriptor *server.Object) int {
	core := GetServer().S()
	var choices []*server.Object
	for p := core.Players.FirstUnit(); p != nil; p = core.Players.NextUnit(p) {
		if uint32(p.ObjFlags)&0x8020 != 0 {
			continue
		}
		info := p.UpdateDataPlayer().Player
		if *(*byte)(unsafe.Pointer(&info.Field3680))&1 != 0 {
			continue
		}
		flags := u.UpdateDataMonsterGen().Flags92
		dist := float64(stateDistance(u, p))
		if flags&1 != 0 {
			if dist <= 300 {
				return generatorPlace(u, out, nil, descriptor)
			}
			continue
		}
		if flags&2 == 0 {
			return 0
		}
		w := uint32(info.Field10)
		h := uint32(info.Field12)
		radius := float32(math.Sqrt(float64(int32(w*w + h*h))))
		if dist <= float64(radius) && core.MapTraceRay(u.PosVec, p.PosVec, 69) {
			choices = append(choices, p)
		}
	}
	if len(choices) == 0 {
		return 0
	}
	return generatorPlace(u, out, choices[core.Rand.Logic.IntClamp(0, len(choices)-1)], descriptor)
}
func generatorUpdate(u *server.Object) int8 {
	core := GetServer().S()
	stage := *memmap.PtrUint32(0x587000, 202028)
	level := *memmap.PtrUint32(0x5D4594, 2388660)
	cache := func(i uintptr) *uint32 { return memmap.PtrUint32(0x5D4594, 2491716+4*i) }
	if dword_5d4594_2491716 == 0 {
		dword_5d4594_2491716 = uint32(int32(float32(core.Balance.Float("QuestHardcoreStage"))))
		*cache(1) = uint32(int32(float32(core.Balance.Float("QuestHardcoreSpawnRateIncrease"))))
		*cache(7) = math.Float32bits(float32(core.Balance.Float("QuestHardcoreSpawnCap")))
		for i, key := range []string{"SpawnRateHighValue", "SpawnRateNormalValue", "SpawnRateLowValue", "SpawnRateVeryLowValue", "SpawnRateVeryVeryLowValue"} {
			*cache(uintptr(i + 2)) = uint32(int32(float32(core.Balance.Float(key))))
		}
	}
	result := u.Field5
	if result&0x800 != 0 {
		return int8(result)
	}
	result = uint32(u.ObjFlags)
	if result&0x1000000 == 0 || result&0x8020 != 0 {
		return int8(result)
	}
	u.NeedSync()
	gen := u.UpdateDataMonsterGen()
	last := &gen.Frame88
	elapsed := core.Frame() - *last
	rate := gen.Rates[level]
	result = memmap.PtrWord(u.CObj())
	if rate <= 4 {
		result = *cache(uintptr(rate) + 2)
	}
	hardcore := uint32(dword_5d4594_2491716)
	if stage >= hardcore {
		reduction := *cache(1) * (stage - hardcore + 1)
		adjusted := uint32(0)
		if reduction <= result {
			adjusted = result - reduction
		}
		cap := float64(math.Float32frombits(*cache(7)))
		if float64(adjusted)/float64(result) < cap {
			adjusted = uint32(int32(float32(float64(result) * cap)))
		}
		result = adjusted
	}
	if elapsed > result {
		current := gen.Current
		result = result&0xffffff00 | uint32(current)
		if current < gen.Max && byte(core.Frame())&8 != 0 {
			sources := gen.Field0[4*level : 4*level+4]
			n := 0
			for _, p := range sources {
				if p != nil {
					n++
				}
			}
			if n == 0 {
				// nothing to spawn at this level (e.g. a generator created without a monster list)
				return int8(result)
			}
			descriptor := sources[core.Rand.Logic.IntClamp(0, n-1)]
			var p types.Pointf
			result = uint32(generatorPick(u, &p, descriptor))
			if result == 1 {
				result = result&0xffffff00 | uint32(byte(generatorSpawn(u, &p, descriptor)))
				*last = core.Frame()
			}
		}
	}
	return int8(result)
}
