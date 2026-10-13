package server

import (
	"unsafe"
)

// Init data blocks. They keep the 386 layout of the original game and are pointer-free.

// DirectionInitData is the init data of objects facing a fixed direction (DirectionInit, SkullInit).
type DirectionInitData struct {
	DX int32 // 0, 0: -1 or 1
	DY int32 // 1, 4: -1 or 1
}

var _ = [1]struct{}{}[8-unsafe.Sizeof(DirectionInitData{})]

// InitDataDirection returns the init data of objects facing a fixed direction (DirectionInit, SkullInit).
func (obj *Object) InitDataDirection() *DirectionInitData {
	return initDataAs[DirectionInitData](obj)
}

// RewardMarkerInitData is the init data of quest reward markers (RewardMarkerInit).
type RewardMarkerInitData struct {
	Categories uint32    // 0, 0: allowed reward categories (bit i: category 1<<i)
	Explicit   uint32    // 1, 4: 1 spells, 2 abilities, 4 field guides are picked from the lists below
	Spells     [137]byte // 2, 8: allowed spells (1 if allowed), by spell ID
	Abilities  [6]byte   // 36, 145: allowed abilities, by ability ID
	Guides     [41]byte  // 37, 151: allowed field guides, by guide ID
	Field192   uint32    // 48, 192
	Field196   uint32    // 49, 196
	Field200   uint32    // 50, 200
	Field204   uint32    // 51, 204
	Field208   uint32    // 52, 208
	Chance     uint32    // 53, 212: 1-4 for a 25/50/75/95% chance to drop nothing; other values always drop
	Flags      byte      // 54, 216: 1 always selected, 0x80 selected to spawn a reward
	Field217   [3]byte   // 54, 217
}

var _ = [1]struct{}{}[220-unsafe.Sizeof(RewardMarkerInitData{})]

// InitDataRewardMarker returns the init data of quest reward markers (RewardMarkerInit).
func (obj *Object) InitDataRewardMarker() *RewardMarkerInitData {
	return initDataAs[RewardMarkerInitData](obj)
}

// AnkhRecord remembers a player that got an extra life from an ankh.
type AnkhRecord struct {
	Name   [25]uint16 // 0: player name (wide string)
	Class  byte       // 50: player class
	Serial [25]byte   // 51: player serial (string)
	Frame  uint32     // 76: frame of the award
}

var _ = [1]struct{}{}[80-unsafe.Sizeof(AnkhRecord{})]

// AnkhInitData is the init data of ankhs (AnkhInit).
type AnkhInitData struct {
	Records   [64]AnkhRecord // 0, 0: players that got an extra life from this ankh
	Next      byte           // 1280, 5120: next record to overwrite
	Field5121 [3]byte        // 1280, 5121
}

var _ = [1]struct{}{}[5124-unsafe.Sizeof(AnkhInitData{})]

// InitDataAnkh returns the init data of ankhs (AnkhInit).
func (obj *Object) InitDataAnkh() *AnkhInitData {
	return initDataAs[AnkhInitData](obj)
}
