package legacy

import (
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func generatorDeath(u *server.Object) {
	core := GetServer().S()
	// C calls both quest services before any script/audio side effect.
	Sub_4D71E0(int(core.Frame()))
	Sub_4D7520(0)
	GetServer().NoxScriptC().ScriptCallback(u.UpdateDataMonsterGen().ScriptAt(56), u.Obj130, u, server.NoxEventGeneratorDead)
	core.Audio.EventObj(1000, u, 0, 0)
	visibilityFXGeneratorBreak(u.PosVec, 200)

	if noxflags.HasGame(4096) && u.Obj130 != nil {
		// The retained owner-chain service returns the terminal object.
		p := nox_xxx_findParentChainPlayer_4EC580(asObjectC(u.Obj130))
		if asObjectS(p).Class().Has(4) {
			questRuntimeIncrement(asObjectS(p), 4668, 8)
		}
	}
	if t := core.NewObjectByTypeID("DestroyedGenerator"); t != nil {
		GetServer().CreateObjectAt(t, nil, u.PosVec)
	}
	GetServer().DelayedDelete(u)
}

// generatorCopy is the private owner of the old 54F2B0 policy. The retained
// inventory/modifier/NPC-equip engines stay on their current C boundaries.
func generatorCopy(src, dst *server.Object) {
	*layout.As[server.MonsterUpdateData](dst.UpdateData) = *layout.As[server.MonsterUpdateData](src.UpdateData)
	if src.ObjSubClass&0x10 != 0 {
		for it := src.InvFirstItem; it != nil; it = it.InvNextItem {
			clone := GetServer().S().NewObjectByTypeInd(int(it.TypeInd))
			if clone == nil {
				continue
			}
			if clone.ObjClass&0x13001000 != 0 {
				stateAttributes(clone, it.InitDataModifier())
			}
			inventoryInsert(dst, clone, 0)
			if it.ObjFlags&0x100 != 0 {
				if clone.ObjClass&0x1001000 != 0 {
					equipmentNPCEquipWeapon(dst, clone)
				} else if clone.ObjClass&0x2000000 != 0 {
					equipmentNPCEquipArmor(dst, clone)
				}
			}
		}
	}
	dst.Direction1, dst.Direction2 = src.Direction1, src.Direction1
}

func generatorHealth(scale float32, n int32) uint16 {
	return uint16(floatToInt32(float32(float64(scale) * float64(n))))
}

// generatorSpawn returns the original EAX bits. They are normally child object
// pointer bits; admission or allocation failure returns 0, and the
// 50E030 rejection path returns objectFreeMem's integer result.
func generatorSpawn(gen *server.Object, point *types.Pointf, src *server.Object) uint32 {
	pos := *point
	core := GetServer().S()
	beholder := memmap.PtrUint32(0x5D4594, 2491712)
	if *beholder == 0 {
		*beholder = uint32(core.Types.IndByID("Beholder"))
	}
	result := uint32(spawnPolicyAdmission(gen, pos))
	if result == 0 {
		return 0
	}
	child := core.NewObjectByTypeInd(int(src.TypeInd))
	if child == nil {
		return 0
	}
	generatorCopy(src, child)

	scale := *memmap.PtrFloat32(0x587000, 202036)
	if def := child.UpdateDataMonster().MonsterDef; def != nil {
		child.HealthData.Cur = generatorHealth(scale, int32(def.HealthQuest72))
		child.HealthData.Max = generatorHealth(scale, int32(def.HealthQuest72))
	} else {
		health := core.Types.ByInd(int(child.TypeInd)).Health()
		child.HealthData.Cur = generatorHealth(scale, int32(health.Cur))
		child.HealthData.Max = generatorHealth(scale, int32(health.Max))
	}
	if child.HealthData.Cur == 0 {
		child.HealthData.Cur = 1
	}
	if child.HealthData.Max == 0 {
		child.HealthData.Max = 1
	}
	if uint32(child.TypeInd) == *beholder {
		child.UpdateDataMonster().Field376 = 0
	}
	if spawnPolicyRegister(gen, child) == 0 {
		return uint32(nox_xxx_objectFreeMem_4E38A0(asObjectC(child)))
	}

	GetServer().CreateObjectAt(child, nil, pos)
	GetServer().NoxScriptC().ScriptCallback(gen.UpdateDataMonsterGen().ScriptAt(64), child, gen, server.NoxEventGeneratorSpawn)
	dir, freeDir := alloc.New(types.Pointf{})
	*dir = pos.Sub(gen.PosVec)
	geometryNormalize((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(dir))))
	fx := [4]int32{
		int32(floatToInt32(gen.PosVec.X)),
		int32(floatToInt32(gen.PosVec.Y)) - 50,
		int32(floatToInt32(float32(float64(dir.X)*30 + float64(pos.X)))),
		int32(floatToInt32(float32(float64(dir.Y)*30 + float64(pos.Y)))),
	}
	freeDir()
	visibilityFXGeneratorSpawn([4]int32{int32(fx[0]), int32(fx[1]), int32(fx[2]), int32(fx[3])}, 10)
	core.Audio.EventObj(1002, child, 0, 0)
	return memmap.PtrWord(child.CObj())
}
