package legacy

import "github.com/brutusmcforce/opennox/legacy/common/alloc"

func Sub_4AD840() {
	serverPanelsGeneralRefresh()
}
func Sub_409E70(a1 int) {
	serverConfigFlagsAdd(int32(a1))
}
func Sub_415960(a1 string) uint32 {
	return runtimeEquipmentMask(false, alloc.InternCString16(a1))
}
func Sub_415DA0(a1 string) uint32 {
	return runtimeEquipmentMask(true, alloc.InternCString16(a1))
}
