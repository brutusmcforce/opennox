//go:build 386

package server

// Layout assertions: these structs keep the layout of the original 32-bit binary.
// On other architectures pointers are larger and the layout differs.

import (
	"unsafe"
)

var _ = [1]struct{}{}[16-unsafe.Sizeof(SystemTime{})]
var _ = [1]struct{}{}[1316-unsafe.Sizeof(NPC{})]
var _ = ([1]struct{}{})[52-unsafe.Sizeof(Shape{})]
var _ = [1]struct{}{}[8-unsafe.Sizeof(PlayerNetData{})]
var _ = [1]struct{}{}[76-unsafe.Sizeof(PlayerJournal{})]
var _ = [1]struct{}{}[16-unsafe.Sizeof(MinimapItem{})]
var _ = [1]struct{}{}[24-unsafe.Sizeof(EquipmentData{})]
var (
	_ = [1]struct{}{}[4828-unsafe.Sizeof(Player{})]
	_ = [1]struct{}{}[2185-unsafe.Offsetof(Player{}.info)]
	_ = [1]struct{}{}[2096-unsafe.Offsetof(Player{}.Field2096Buf)]
	_ = [1]struct{}{}[2282-unsafe.Offsetof(Player{}.Field2282)]
	_ = [1]struct{}{}[3596-unsafe.Offsetof(Player{}.Frame3596)]
	_ = [1]struct{}{}[4580-unsafe.Offsetof(Player{}.Field4580)]
	_ = [1]struct{}{}[4800-unsafe.Offsetof(Player{}.data4800)]
)
var (
	_ = [1]struct{}{}[97-unsafe.Sizeof(PlayerInfo{})]
	_ = [1]struct{}{}[0-unsafe.Offsetof(PlayerInfo{}.name)]
	_ = [1]struct{}{}[66-unsafe.Offsetof(PlayerInfo{}.playerClass)]
	_ = [1]struct{}{}[89-unsafe.Offsetof(PlayerInfo{}.nameSuff)]
)
var (
	_ = [1]struct{}{}[72-unsafe.Sizeof(SavePlayerInfo{})]
	_ = [1]struct{}{}[20-unsafe.Offsetof(SavePlayerInfo{}.NameBuf)]
)
var (
	_ = [1]struct{}{}[1277-unsafe.Offsetof(SaveGameInfo{}.Stage)]
	_ = [1]struct{}{}[1280-unsafe.Sizeof(SaveGameInfo{})] // FIXME: should be 1278
)
