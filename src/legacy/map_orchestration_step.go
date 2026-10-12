package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func mapOrchestrationStep() uint32 {
	*populationGlobal(0) = 0
	Sub_57C490_2("theme")
	mapPaintAfterWalls(0)
	*memmap.PtrUint32(0x973F18, 35976) = 0
	cfg := memmap.PtrWord(mapGrowthConfig())
	if mapThemeFile(cfg, memmap.PtrWord(memmap.PtrOff(0x587000, 197860))) == 0 {
		return 0
	}
	*populationWord(cfg, 68) = uint32(int64(float64(*populationFloat(cfg, 64)) * 0.030743772))
	mapRoomSeed(*populationWord(cfg, 76))
	mapPopulationMetadataInit()
	if mapRoomScratchAlloc() == 0 {
		return 0
	}
	if mapRoomGridInit(mapGrowthConfig()) == 0 {
		return 0
	}
	mapGrowthInitial(cfg)
	result := uint32(2)
	if mapPopulationSelectPrefabs(cfg) != 0 {
		mapGrowthFrontiers()
		if mapPopulationConnectPrefabs(cfg) != 0 {
			if mapPopulationApplyPrefabs(cfg) == 0 {
				result = 0
			} else {
				mapPopulationDistance(*populationGlobal(0), 0, 0)
				mapPopulationThemes(*populationGlobal(0))
				if *populationWord(cfg, 184) != 0 {
					radius := int32(int64(float64(*populationFloat(cfg, 64)) * 0.030743772))
					r := mapRoomNew(2*radius+1, 2*radius+1)
					x := float32(float64(-radius) * 32.526913)
					p := types.Ptf(x, x)
					mapRoomSetPos(r, &p)
					for n := mapRoomHead(); n != nil; n = n.Next {
						mapRoomAddExclusion(r, &n.Pos, n.Size.X, n.Size.Y)
					}
					mapRoomAssignDecoration(mapGrowthConfig(), r)
					mapPaintRoomWalls(mapGrowthConfig(), r)
					mapPopulationRoom(cfg, memmap.PtrWord(unsafe.Pointer(r)))
					mapRoomFree(r)
				}
				if mapRoomAssignRequiredDecorations(mapGrowthConfig()) != 0 {
					for r := mapRoomHead(); r != nil; r = r.Next {
						if mapRoomIsHall(r) != 0 {
							mapPopulationProgress(156)
							mapPaintRoomWalls(mapGrowthConfig(), r)
						}
					}
					for r := mapRoomHead(); r != nil; r = r.Next {
						if mapRoomIsHall(r) == 0 {
							mapPopulationProgress(156)
							mapPaintRoomWalls(mapGrowthConfig(), r)
						}
					}
					mapPopulationHallwayWaypoints(cfg)
					mapGrowthDoors()
					mapPopulationRoomWaypoints(cfg)
					mapPopulationFinish(cfg)
					result = 1
				}
			}
		}
	}
	mapRoomFreeAll()
	mapRoomScratchFree()
	mapRoomGridFree()
	mapThemeCleanup(cfg)
	return result
}
