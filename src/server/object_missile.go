package server

import "github.com/brutusmcforce/opennox/libs/types"

type MissileUpdateData struct {
	Owner    *Object // 0, 0
	Target   *Object // 1, 4
	Field_8  uint32  // 2, 8
	SpellID  int32   // 3, 12
	Field_16 uint32  // 4, 16
	Field_20 uint32  // 5, 20
	Field_24 uint32  // 6, 24
}

func (obj *Object) UpdateDataMissile() *MissileUpdateData {
	return updateDataAs[MissileUpdateData](obj)
}

// ChakramUpdateData is the update data of chakrams in motion (ChakramInMotionUpdate).
type ChakramUpdateData struct {
	Field0    uint32       // 0, 0
	Bounces   byte         // 1, 4
	_         [3]byte      // 1, 5
	Target    *Object      // 2, 8
	Prior     *Object      // 3, 12: the last object hit
	OwnerPos  types.Pointf // 4, 16
	Returning byte         // 6, 24
	_         [3]byte      // 6, 25
}

func (obj *Object) UpdateDataChakram() *ChakramUpdateData {
	return updateDataAs[ChakramUpdateData](obj)
}
