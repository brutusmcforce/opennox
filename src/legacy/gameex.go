package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/server"
)

func Sub_4BDFD0() { serverPanelsAdvancedServerOpen() }
func Mix_MouseKeyboardWeaponRoll(u *server.Object, direction int8) int {
	return extensionWeaponRoll(u, direction)
}
func PlayerInfoStructParser_0(record unsafe.Pointer) int {
	return extensionPlayerName(record, nil)
}
func PlayerInfoStructParser_1(record unsafe.Pointer, team *int32) int {
	if team == nil {
		return 0
	}
	return extensionPlayerName(record, team)
}
func PlayerDropATrap(u *server.Object) { extensionDropTrap(u) }
func GetFlagValueFromFlagIndex(index int) uint32 {
	// The settings UI uses 1..5. Explicitly handle the wider Go input domain,
	// including values that made the old C helper divide by zero.
	if index < 0 || index >= 32 {
		return 0
	}
	return uint32(1) << uint(index)
}
