package legacy

func nox_xxx_inventoryGetFirst_4E7980(a int32) int32 {
	return int32(inventoryInt(objectFromInt(a).InvFirstItem))
}

func nox_xxx_inventoryGetNext_4E7990(a int32) int32 {
	if a == 0 {
		return 0
	}
	return int32(inventoryInt(objectFromInt(a).InvNextItem))
}
