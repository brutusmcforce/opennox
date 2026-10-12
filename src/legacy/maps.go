package legacy

import (
	"fmt"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/log"
	"github.com/brutusmcforce/opennox/libs/types"

	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/server"
)

var (
	gameLog                           = log.New("game")
	mapLog                            = log.New("map")
	Nox_common_checkMapFile           func(name string) error
	Nox_xxx_mapWriteSectionsMB_426E20 func(a1 unsafe.Pointer) int
	// Nox_xxx_mapGenUseBlendScript makes the script of the generated map (blend.obj) the current map script.
	Nox_xxx_mapGenUseBlendScript func() bool
	Nox_xxx_mapReadSection       func(*cryptfile.CryptFile, unsafe.Pointer, string) (bool, error)
)

func mapGenUseBlendScript() bool {
	return Nox_xxx_mapGenUseBlendScript != nil && Nox_xxx_mapGenUseBlendScript()
}

func Nox_server_mapRWMapInfo_42A6E0(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if mapMetadataInfo() == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWWallMap_429B20(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if mapSectionWalls((*[8]uint32)(a1)) == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWFloorMap_422230(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if mapSectionFloor((*[8]uint32)(a1)) == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWSecretWalls_4297C0(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if mapSectionMetadata(2, (*[8]uint32)(a1)) == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWDestructableWalls_429530(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if mapSectionMetadata(1, (*[8]uint32)(a1)) == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWWaypoints_506260(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if prefabWaypointSection((*[8]uint32)(a1)) == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWWindowWalls_4292C0(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if mapSectionMetadata(0, (*[8]uint32)(a1)) == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWGroupData_505C30(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if prefabGroupSection() == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWAmbientData_429200(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if mapMetadataAmbient() == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWPolygons_428CD0(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if mapPolygonSection(int(uintptr(a1))) == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWMapIntro_505080(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if prefabIntroSection() == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}
func Nox_server_mapRWObjectTOC_428B30(_ *cryptfile.CryptFile, a1 unsafe.Pointer) error {
	if mapMetadataTOC() == 0 {
		return fmt.Errorf("%s failed", caller(0))
	}
	return nil
}

func Nox_xxx_servMapLoadPlaceObj_4F3F50(obj *server.Object, a2 int, a3 unsafe.Pointer) int {
	return objectXferPlace(obj, unsafe.Pointer(uintptr(uint32(a2))), a3)
}

func Nox_xxx_prepareLightningEffects_4BAB30() {
	effectPrepareLightning()
}

func Sub_4B64C0() {
	initParticleColors()
}

func Nox_xxx_bookSetColor_45AC40() {
	bookSetColor()
}

func Nox_xxx_colorInit_4C4FD0() {
	objectRenderBeamColors()
}

func Sub_445FF0() {
	sessionQuitColors()
}

func Sub_470680() {
	scoreboardClearFlags()
}

func Sub_461520() {
	clientAllPlayerColors()
}

func Nox_xxx_tile_486060() {
	tileRasterWrapSetup()
}

func Sub_461400() {
	quickbarRestoreSlots()
}

func Sub_461450() int {
	return int(*quickbarWord(1049688))
}

func Nox_xxx_cliShowHideTubes_470AA0(v int) {
	nox_xxx_cliShowHideTubes_470AA0(v)
}

func Sub_51DED0() int {
	return worldSaveDrawables()
}

func Sub_428170(a1, a2 unsafe.Pointer) {
	geometryWallBounds((*[8]uint32)(unsafe.Pointer(a1)), (*[4]int32)(unsafe.Pointer(a2)))
}

func Nox_xxx_tileNFromPoint_411160(p types.Pointf) int {
	return int(tileAtPoint(p))
}

func Nox_xxx_unitSetDecayTime_511660(obj *server.Object, a2 int) {
	motionDecaySet(obj, int32(a2))
}

func Nox_xxx_tileFreeTileOne_4221E0(p unsafe.Pointer) {
	mapPaintSubtileFree((*[5]uint32)(p))
}

func Sub_4DE410(pli ntype.PlayerInd) {
	matchRosterReportMask(int(pli))
}
