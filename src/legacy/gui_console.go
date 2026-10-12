package legacy

import (
	"github.com/brutusmcforce/opennox/libs/console"
	"unicode/utf16"
)

var (
	GetConsole                     func() *console.Console
	Nox_gui_console_flagXxx_451410 func() int
	Nox_gui_console_Hide_4512B0    func() int
)

func Nox_xxx_consoleEsc_49B7A0() {
	nox_xxx_consoleEsc_49B7A0()
}

func textFormatConsole(color byte, format *uint16, args ...textFormatArgument) int {
	GetConsole().Print(console.Color(color), string(utf16.Decode(textFormatUnits(format, args...))))
	return 1
}
