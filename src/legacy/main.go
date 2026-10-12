package legacy

var (
	Nox_exit                                              func(exitCode int)
	Nox_xxx_gameGetScreenBoundaries_43BEB0_get_video_mode func(w, h, d *int)
	Sub_4AA9C0                                            func() int
)

func Nox_xxx_loadLook_415D50() {
	runtimeEquipmentLoad(true)
}

func Nox_xxx_loadModifyers_4158C0() {
	runtimeEquipmentLoad(false)
}

func Sub_4D11A0() {
	sessionRosterInit()
}

func Sub_431370() int {
	uiEntryContext = true
	return 1
}

func Nox_xxx_tileAlloc_410F60_init() int {
	return int(worldGridAllocate())
}

func Nox_xxx_initSinCosTables_414C90() {
	worldInitTrigTables()
}

func Nox_xxx_loadMapCycle_4D0A30() {
	mapCycleLoad()
}

func Nox_xxx_mapSelectFirst_4D0E00() {
	mapQuestBuild()
}

func Sub_413920() {
	Set_dword_5d4594_251744(0)
}

func Sub_431380() {
	uiEntryContext = false
}

func Nox_xxx_tileFree_410FC0_free() {
	worldGridFreeRows()
}

func Sub_4D11D0() {
	sessionRosterClear()
}

func Sub_4D0DA0() {
	mapCycleReset()
}

func Nox_common_maplist_free_4D0970() {
	mapCatalogFree()
}

func Sub_451970() {
	audioEventFree()
}

func Sub_431270() {
	clientAudioContextDestroy()
}

func Sub_4875F0() {
	audioStreamContextDestroyAll()
}

func Sub_4870A0() {
	audioStreamDeviceDestroyAll()
}

func Sub_431290() {
	clientAudioContextStop()
}

func Nox_xxx_servSetPlrLimit_409F80(v int) {
	serverConfigLimitSet(int32(v))
}

func Nox_xxx_guiChatShowHide_445730(v bool) {
	nox_xxx_guiChatShowHide_445730(int32(bool2int(v)))
}
