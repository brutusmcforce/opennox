package legacy

func sub_56F250() int32 {
	var result int32
	for i := 0; i < 7; i++ {
		result = int32(createProtectionRecord(uint32(dword_5d4594_2516356), 0))
		// Reserved slots consume an ID even if record allocation fails.
		dword_5d4594_2516356++
	}
	return result
}

func nox_xxx_protectionCreateInt_56F400(value int32) int32 {
	if createProtectionRecord(uint32(dword_5d4594_2516356), uint32(value)) == 0 {
		return 0
	}
	id := dword_5d4594_2516356
	dword_5d4594_2516356++
	return int32(id)
}
