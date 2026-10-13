package server

// MoverUpdateData is the update data of movers (MoverUpdate), which move another object along waypoints.
type MoverUpdateData struct {
	Field_0 uint8     // 0, 0: state
	_       [3]byte   // 0, 1
	Field_1 float32   // 1, 4: speed (read as the integer bits by the update, see legacy motionMover)
	Field_2 int32     // 2, 8: first waypoint index
	Next    *Waypoint // 3, 12
	NextInd uint32    // 4, 16: index of Next, until it is resolved
	Prev    *Waypoint // 5, 20
	PrevInd uint32    // 6, 24: index of Prev, until it is resolved
	Field_7 *Object   // 7, 28: the moved object
	Field_8 uint32    // 8, 32: extent of the moved object, until it is resolved
}

func (obj *Object) UpdateDataMover() *MoverUpdateData {
	return updateDataAs[MoverUpdateData](obj)
}
