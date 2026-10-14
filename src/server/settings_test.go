package server

import (
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/shoenig/test/must"
)

// The password is at byte 78 of the settings ((wchar_t*)settings + 39 in the original), opennox#574.
func TestSettingsPassword(t *testing.T) {
	must.EqOp(t, 78, unsafe.Offsetof(Settings{}.ServerPass78))
	var s Settings
	must.EqOp(t, "", s.Password())
	copy(s.ServerPass78[:], utf16.Encode([]rune("secret")))
	must.EqOp(t, "secret", s.Password())
	copy(s.ServerPass78[:], utf16.Encode([]rune("ninechars")))
	must.EqOp(t, "ninechars", s.Password()) // no terminator: the whole field
}
