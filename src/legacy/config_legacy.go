package legacy

import "unsafe"

func Sub_453070() int {
	return int(int32(*audioEventPlayback))
}
func Sub_44D990() int {
	return int(audioEventDialogEnabled())
}
func Sub_43DC30() int {
	return int(audioEventMusicEnabled())
}
func Nox_xxx_sysopGetPass_40A630() string {
	return GoWStringP(unsafe.Pointer(serverConfigPasswordGet()))
}
func Sub_4D0D70() int {
	return mapCycleEnabled()
}
func Nox_xxx_getServerSubFlags_409E60() uint32 {
	return uint32(serverConfigFlagsGet())
}
func Nox_xxx_gameSetServername_40A440(a1 string) {
	serverConfigNameSet((*byte)(unsafe.Pointer(internCStr(a1))))
}
func Nox_xxx_sysopSetPass_40A610(a1 string) {
	serverConfigPasswordSet((*uint16)(unsafe.Pointer(internWStr(a1))))
}
func Nox_xxx_rateUpdate_40A6D0(a1 int) {
	serverConfigRateSet(int32(a1))
}
func Sub_4D0D90(a1 int) {
	mapCycleSetEnabled(uint32(a1))
}
func Sub_409FB0_settings(a1 uint16, a2 uint16) {
	serverConfigScoreSet(int16(a1), a2)
}
func Sub_409EC0(a1 int) {
	serverConfigFlagsRemove(int32(a1))
}
func Sub_4D0DC0(a1 uint32, a2 int) {
	mapCycleSetIndex(a1, uint32(a2))
}
func Sub_489FF0(a1 int, a2 int, a3 unsafe.Pointer) {
	sessionFilterConfig(a1, a2, a3)
}
func Sub_4D0DE0(a1 uint32) int {
	return int(mapCycleGetIndex(a1))
}
