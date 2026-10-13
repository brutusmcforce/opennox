package legacy

import (
	"math"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func questProgressStage() uint32              { return memmap.Uint32(0x5D4594, 2388660) }
func questProgressSetStage(v uint32) uint32   { *memmap.PtrUint32(0x5D4594, 2388660) = v; return v }
func questProgressSetMinions(v uint32) uint32 { *memmap.PtrUint32(0x5D4594, 2388656) = v; return v }
func questProgressInitMapping() {
	for off := uintptr(249896); memmap.LoadPtr(0x587000, off+8) != nil; off += 16 {
		for _, delta := range []uintptr{8, 0} {
			name := alloc.GoString((*byte)(memmap.LoadPtr(0x587000, off+delta)))
			*memmap.PtrUint32(0x587000, off+delta+4) = uint32(GetServer().S().Types.IndByID(name))
		}
	}
	*memmap.PtrUint32(0x5D4594, 2388664) = 1
}
func questProgressGeneratorType(u *server.Object) uint32 {
	if memmap.Uint32(0x5D4594, 2388664) == 0 {
		questProgressInitMapping()
	}
	if u == nil {
		return 0
	}
	for off := uintptr(249896); memmap.LoadPtr(0x587000, off+8) != nil; off += 16 {
		if memmap.Uint32(0x587000, off+4) == uint32(u.TypeInd) {
			return memmap.Uint32(0x587000, off+12)
		}
	}
	return 0
}
func questProgressSpawnBoss(pos types.Pointf, hecubah bool) {
	s := GetServer().S()
	name := "Necromancer"
	if hecubah {
		name = "Hecubah"
	}
	u := s.NewObjectByTypeID(name)
	scale := float32(questRuntimeFloat(202036))
	if u == nil {
		return
	}
	d := u.UpdateDataMonster()
	var hp int32
	if def := d.MonsterDef; def != nil {
		hp = int32(def.HealthQuest72)
	} else {
		hp = int32(s.Types.ByInd(int(u.TypeInd)).Health().Max)
	}
	if scale < 1 {
		scale = 1
	}
	value := float64(hp) * float64(scale)
	// Compiled C keeps the current-HP product wide for int64 truncation, but
	// spills a float32 copy for the independent maximum-HP conversion.
	resourceSetHP(u, uint16(int64(value)))
	u.HealthData.Max = uint16(floatToInt32(float32(value)))
	if u.HealthData.Cur == 0 {
		resourceSetHP(u, 1)
	}
	if u.HealthData.Max == 0 {
		u.HealthData.Max = 1
	}
	d.Field411 = 0x10000000
	d.Field423 = 0x10000000
	d.AIAction340 = 4
	d.Aggression = math.Float32frombits(1062501089)
	d.Field410 = 0x8000000
	d.Field444 = 0x20000000
	d.Field415 = 0x40000000
	rewards := 1
	d.Field510 = 1
	if hecubah {
		rewards = 4
		d.Field510 = 3
		d.Field388 = 0x40000000
		s.Balance.Float("HecubahQuestSkill")
		d.Field330 = math.Float32frombits(1062836634)
	}
	GetServer().CreateObjectAt(u, nil, pos)
	if marker := s.NewObjectByTypeID("RewardMarker"); marker != nil {
		for i := 0; i < rewards; i++ {
			if item := rewardMarker(marker, questRuntimeStage()+2); item != nil {
				inventoryInsert(u, item, 0)
			}
		}
		s.Objs.FreeObject(marker)
	}
}
func questProgressPrepare(level int) {
	s := GetServer().S()
	stage := int32(questRuntimeStage())
	hardcore := uint32(floatToInt32(float32(s.Balance.Float("QuestHardcoreStage"))))
	if memmap.Uint32(0x5D4594, 2388668) == 0 {
		*memmap.PtrUint32(0x5D4594, 2388668) = uint32(s.Types.IndByID("HecubahMarker"))
		*memmap.PtrUint32(0x5D4594, 2388672) = uint32(s.Types.IndByID("NecromancerMarker"))
	}
	var exits uint32
	hecubahMarkers := 0
	for u := s.Objs.First(); u != nil; {
		next := u.Next()
		if u.ObjClass&0x20 != 0 && u.ObjSubClass&1 != 0 {
			exits++
		} else if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2388668) {
			hecubahMarkers++
		}
		if u.ObjClass&0x20000 != 0 {
			d := u.UpdateDataMonsterGen()
			if template := d.Field0[4*level]; template != nil {
				rate := d.Field83[level]
				keys := [4]string{"GeneratorMaxActiveCreaturesHigh", "GeneratorMaxActiveCreaturesNormal", "GeneratorMaxActiveCreaturesLow", "GeneratorMaxActiveCreaturesSingular"}
				if rate < 4 {
					d.Max = byte(int64(s.Balance.Float(keys[rate])))
				}
				if questRuntimeStage() >= hardcore && rate != 3 {
					d.Max *= 2
				}
				if ind := questProgressGeneratorType(template); ind != 0 {
					u.TypeInd = uint16(ind)
				}
			} else {
				GetServer().DelayedDelete(u)
			}
		}
		u = next
	}
	if exits > 1 {
		chosen := s.Rand.Logic.IntClamp(0, int(int32(exits-1)))
		index := 0
		for u := s.Objs.First(); u != nil; {
			next := u.Next()
			if u.ObjClass&0x20 != 0 && u.ObjSubClass&1 != 0 {
				if index != chosen {
					GetServer().DelayedDelete(u)
				}
				index++
			}
			u = next
		}
	}
	questProgressSetMinions(0)
	if stage >= 5 {
		always := floatToInt32(float32(s.Balance.Float("MinionsAlwaysStage")))
		if stage == 5 || stage >= always || (stage&1 != 0 && s.Rand.Logic.IntClamp(1, 100) >= 50) {
			questProgressSetMinions(1)
			if hecubahMarkers != 0 {
				chosen := s.Rand.Logic.IntClamp(1, hecubahMarkers)
				index := 0
				for u := s.Objs.First(); u != nil; u = u.Next() {
					if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2388668) {
						index++
						if index == chosen {
							questProgressSpawnBoss(u.PosVec, true)
						}
					}
					if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2388672) && s.Rand.Logic.IntClamp(1, 100) >= 50 {
						questProgressSpawnBoss(u.PosVec, false)
					}
				}
			}
		}
	}
	for u := s.Objs.First(); u != nil; {
		next := u.Next()
		if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2388668) {
			GetServer().DelayedDelete(u)
		}
		if uint32(u.TypeInd) == memmap.Uint32(0x5D4594, 2388672) {
			GetServer().DelayedDelete(u)
		}
		u = next
	}
}
