package legacy

func nox_xxx_packetDynamicUnitCode_578B40(value int32) int32 {
	return int32(networkDynamicUnitCode(uint32(value)))
}

func networkDynamicUnitCode(code uint32) uint32 {
	if code&0x8000 == 0 {
		return code
	}
	obj := GetServer().S().Objs.GetObjectByInd(int(code &^ 0x8000))
	if obj == nil {
		return 0
	}
	return obj.NetCode
}
