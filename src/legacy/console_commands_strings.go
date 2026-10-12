package legacy

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/brutusmcforce/opennox/libs/console"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

var consoleCommandServer bool
var consoleCommandSender *server.Player

// ConsoleCommand retains the registered legacy token-index contract without a C call.
func ConsoleCommand(fn func(int, []string) bool) console.CommandLegacyFunc {
	return func(ctx context.Context, _ *console.Console, index int, tokens []string) bool {
		consoleCommandServer = !console.IsClient(ctx)
		normalized := make([]string, len(tokens))
		for i, s := range tokens {
			normalized[i] = consoleCommandString(s)
		}
		return fn(index, normalized)
	}
}
func consoleCommandString(s string) string {
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	// Match conversion to UTF16 and back, including invalid UTF8 input.
	return string([]rune(s))
}
func consoleCommandNarrow(s string) string {
	var b []byte
	for _, v := range utf16.Encode([]rune(s)) {
		if byte(v) == 0 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}
func consoleCommandWiden(s string) string {
	r := make([]rune, 0, len(s))
	for i := 0; i < len(s) && s[i] != 0; i++ {
		r = append(r, rune(s[i]))
	}
	return string(r)
}

// C strtol on this target accepts an ASCII decimal prefix and saturates to int32.
func consoleCommandNumber(s string) int32 {
	s = strings.TrimLeft(s, " \t\n\r\v\f")
	neg := false
	if len(s) > 0 && (s[0] == '+' || s[0] == '-') {
		neg = s[0] == '-'
		s = s[1:]
	}
	limit := uint32(2147483647)
	if neg {
		limit++
	}
	var n uint32
	for i := 0; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
		d := uint32(s[i] - '0')
		if n > (limit-d)/10 {
			if neg {
				return -2147483648
			}
			return 2147483647
		}
		n = n*10 + d
	}
	if neg {
		return -int32(n)
	}
	return int32(n)
}
func consoleCommandText(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), "parsecmd.c")
}

// The C formatter accepts two literal escapes and ignores width/precision for
// strings. Scan conversions so an escaped %%S stays uppercase literal text.
var consoleCommandEscapes = strings.NewReplacer("%%", "%", "%!", "!")

func consoleCommandFormat(format string, args ...any) string {
	if len(args) == 0 {
		return consoleCommandEscapes.Replace(format)
	}
	var out strings.Builder
	for i := 0; i < len(format); {
		if format[i] != '%' {
			out.WriteByte(format[i])
			i++
			continue
		}
		start := i
		i++
		if i == len(format) {
			out.WriteByte('%')
			break
		}
		if format[i] == '%' {
			out.WriteString("%%")
			i++
			continue
		}
		if format[i] == '!' {
			out.WriteByte('!')
			i++
			continue
		}
		for i < len(format) && strings.ContainsRune("-0123456789.", rune(format[i])) {
			i++
		}
		if i == len(format) {
			out.WriteString(format[start:])
			break
		}
		switch format[i] {
		case 's', 'S':
			out.WriteString("%s")
		case 'i', 'u':
			out.WriteString(format[start:i])
			out.WriteByte('d')
		default:
			out.WriteString(format[start : i+1])
		}
		i++
	}
	return fmt.Sprintf(out.String(), args...)
}
func consoleCommandPrint(id string, args ...any) {
	GetConsole().Print(console.ColorRed, consoleCommandFormat(consoleCommandText(id), args...))
}
func consoleCommandCentered(id string, args ...any) {
	Nox_xxx_printCentered_445490(consoleCommandFormat(consoleCommandText(id), args...))
}
func consoleCommandToken(i uintptr) string {
	return alloc.GoString16((*uint16)(memmap.LoadPtr(0x587000, 94468+4*i)))
}
