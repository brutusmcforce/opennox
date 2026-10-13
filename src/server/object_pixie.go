package server

// PixieUpdateData is the update data of pixies (PixieUpdate).
type PixieUpdateData struct {
	Owner     *Object // 0, 0
	Target    *Object // 1, 4
	Field2    uint32  // 2, 8
	SpellID   uint32  // 3, 12
	Field4    uint32  // 4, 16
	Deadline  uint32  // 5, 20: frame to return to the owner
	SeenOwner uint32  // 6, 24: last frame the owner was visible
}

func (obj *Object) UpdateDataPixie() *PixieUpdateData {
	return updateDataAs[PixieUpdateData](obj)
}
