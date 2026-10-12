package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"math"
)

func questRuntimeFloat(off uintptr) float64 { return float64(memmap.Float32(0x587000, off)) }
func questRuntimeSetDifficulty(value float32) uint32 {
	*memmap.PtrFloat32(0x587000, 202024) = value
	return math.Float32bits(value)
}
func questRuntimeCap(off uintptr, key string, value float32) {
	cap := GetServer().S().Balance.Float(key)
	if float64(value) <= cap {
		*memmap.PtrFloat32(0x587000, off) = value
	} else {
		*memmap.PtrFloat32(0x587000, off) = float32(cap)
	}
}
func questRuntimeDifficulty() uint32 {
	if questRuntimeWord(1563928) == 0 {
		*memmap.PtrFloat32(0x5D4594, 1563912) = float32(GetServer().S().Balance.Float("PlayerDifficultyDelta"))
		questRuntimeSetWord(1563928, 1)
	}
	value := float64(questRuntimeStage()) * (float64(questRuntimeCount()-1)*float64(memmap.Float32(0x5D4594, 1563912)) + 1)
	return questRuntimeSetDifficulty(float32(value))
}
func questRuntimeAbsHealth(value float32) uint16 {
	return uint16(floatToInt32(math.Float32frombits(math.Float32bits(value) & 0x7fffffff)))
}
func questRuntimeScaleHealth() {
	core := GetServer().S()
	limit := uint16(floatToInt32(float32(core.Balance.Float("GeneratorMaxHealth"))))
	if questRuntimeWord(1563932) == 0 {
		for _, item := range []struct {
			off uintptr
			key string
		}{{1563908, "PlayerDamageDiffInit"}, {1563916, "SystemHealthDiffInit"}, {1563920, "PlayerDamageDiffCoeff"}, {1563924, "SystemHealthDiffCoeff"}} {
			*memmap.PtrFloat32(0x5D4594, item.off) = float32(core.Balance.Float(item.key))
		}
		questRuntimeSetWord(1563932, 1)
	}
	questRuntimeCap(202032, "PlayerDamageCap", float32((questRuntimeFloat(202024)-1)*float64(memmap.Float32(0x5D4594, 1563920))+float64(memmap.Float32(0x5D4594, 1563908))))
	questRuntimeCap(202036, "SystemHealthCap", float32((questRuntimeFloat(202024)-1)*float64(memmap.Float32(0x5D4594, 1563924))+float64(memmap.Float32(0x5D4594, 1563916))))
	for u := core.Objs.First(); u != nil; {
		next := u.Next()
		if u.ObjFlags&0x8000 == 0 && (u.ObjClass&0x20000 != 0 || u.ObjClass&2 != 0) {
			h := u.HealthData
			// The original promotes its maximum through signed short before comparing.
			if h.Cur != 0 && h.Max != 0 && int(h.Cur) == int(int16(h.Max)) {
				base := core.Types.ByInd(int(u.TypeInd)).Health()
				if u.ObjClass&0x20000 != 0 {
					max := questRuntimeAbsHealth(float32(questRuntimeFloat(202036) * float64(base.Max)))
					cur := questRuntimeAbsHealth(float32(questRuntimeFloat(202036) * float64(base.Cur)))
					if cur == 0 {
						cur = 1
					}
					if max == 0 {
						max = 1
					}
					if cur > limit {
						cur = limit
					}
					if max > limit {
						max = limit
					}
					resourceSetHP(u, cur)
					u.HealthData.Max = max
				} else {
					d := u.UpdateData
					hp := base.Max
					if p := *controlPtr(d, 484); p != nil {
						hp = *controlHalf(p, 72)
					}
					if int8(*controlByte(d, 1440)) >= 0 {
						value := questRuntimeAbsHealth(float32(questRuntimeFloat(202036) * float64(hp)))
						if value == 0 {
							value = 1
						}
						resourceSetHP(u, value)
						u.HealthData.Max = value
						for i := 0; i < 32; i++ {
							*controlHalf(d, 412+2*i) = u.HealthData.Cur
						}
					}
				}
			}
		}
		u = next
	}
}
