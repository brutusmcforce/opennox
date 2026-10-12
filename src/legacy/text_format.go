package legacy

import (
	"github.com/brutusmcforce/opennox/server"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

// Varargs are gone; each caller supplies the actual word, real or string value.
// Pointer arguments remain raw UTF-16/narrow text until the conversion is read.
type textFormatArgument struct {
	word uint32
	real float64
	ptr  unsafe.Pointer
}

func textFormatWord(v uint32) textFormatArgument            { return textFormatArgument{word: v} }
func textFormatReal(v float64) textFormatArgument           { return textFormatArgument{real: v} }
func textFormatPointer(v unsafe.Pointer) textFormatArgument { return textFormatArgument{ptr: v} }

// Preserve the historical Nox formatter, including its deliberately non-printf
// padding and signed word formatting. Return the complete UTF-16 length even if
// the destination is short; append a terminator only when it fits.
func textFormatInto(dst *uint16, count int, format *uint16, args ...textFormatArgument) int {
	var buffer []uint16
	if dst != nil && count > 0 {
		buffer = unsafe.Slice(dst, count)
	}
	out := 0
	emit := func(ch uint16) {
		if out < len(buffer) {
			buffer[out] = ch
		}
		out++
	}
	emitBytes := func(s string) {
		for i := 0; i < len(s); i++ {
			emit(uint16(s[i]))
		}
	}
	repeat := func(ch uint16, n int32) {
		for ; n > 0; n-- {
			emit(ch)
		}
	}
	pos := uintptr(0)
	read := func() uint16 { v := *(*uint16)(unsafe.Add(unsafe.Pointer(format), pos)); pos += 2; return v }
	ai := 0
	argument := func() textFormatArgument {
		if ai >= len(args) {
			panic("missing text formatting argument")
		}
		v := args[ai]
		ai++
		return v
	}
	for ch := read(); ch != 0; ch = read() {
		if ch != '%' {
			emit(ch)
			continue
		}
		ch = read()
		flag := uint16(0)
		if ch == '+' || ch == '-' || ch == '0' {
			flag = ch
			ch = read()
		}
		width, precision := int32(-1), int32(-1)
		for ch >= '0' && ch <= '9' {
			if width < 0 {
				width = 0
			}
			width = width*10 + int32(ch-'0')
			ch = read()
		}
		if ch == '.' {
			ch = read()
			precision = 0
			for ch >= '0' && ch <= '9' {
				precision = precision*10 + int32(ch-'0')
				ch = read()
			}
		}
		switch ch {
		case 'c':
			emit(uint16(argument().word))
		case 's', 'S':
			a := argument()
			if a.ptr == nil {
				emitBytes("(null)")
				continue
			}
			if ch == 's' {
				for off := uintptr(0); ; off += 2 {
					v := *(*uint16)(unsafe.Add(a.ptr, off))
					if v == 0 {
						break
					}
					emit(v)
				}
			} else {
				for off := uintptr(0); ; off++ {
					v := *(*byte)(unsafe.Add(a.ptr, off))
					if v == 0 {
						break
					}
					emit(uint16(v))
				}
			}
		case 'd', 'i', 'u', 'x', 'X', 'o':
			base := 10
			if ch == 'x' || ch == 'X' {
				base = 16
			}
			if ch == 'o' {
				base = 8
			}
			digits := strconv.FormatInt(int64(int32(argument().word)), base)
			if ch == 'X' {
				digits = strings.ToUpper(digits)
			}
			n := int32(len(digits))
			used := n
			if precision > 0 {
				used = precision
			}
			pad := uint16(' ')
			if flag == '0' {
				pad = '0'
			}
			repeat(pad, width-used)
			repeat('0', precision-n)
			emitBytes(digits)
		case 'f':
			value := argument().real
			digits := int(precision)
			if digits <= 0 {
				digits = 5
			}
			// Every finite binary64 fraction has an exact decimal expansion within
			// 1074 fractional digits; greater precision only adds zeros beyond the
			// old 31-byte result. Avoid allocating arbitrary format-specified sizes.
			digits = min(digits, 1074)
			s := strconv.FormatFloat(value, 'f', digits, 64)
			if math.IsInf(value, 1) {
				s = "inf"
			} else if math.IsInf(value, -1) {
				s = "-inf"
			} else if math.IsNaN(value) {
				s = "nan"
				if math.Signbit(value) {
					s = "-nan"
				}
			}
			if len(s) > 31 {
				s = s[:31]
			}
			emitBytes(s)
		case '%', '!':
			emit(ch)
		default:
			panic("unsupported text formatting conversion")
		}
	}
	emit(0)
	return out - 1
}

// textFormatBuffer is for terminated-string consumers; the low-level bounded
// formatter above deliberately does not terminate a full destination.
func textFormatBuffer(dst []uint16, format *uint16, args ...textFormatArgument) int {
	var p *uint16
	if len(dst) > 0 {
		p = &dst[0]
	}
	n := textFormatInto(p, len(dst), format, args...)
	if len(dst) > 0 && n >= len(dst) {
		dst[len(dst)-1] = 0
	}
	return n
}
func textFormatUnits(format *uint16, args ...textFormatArgument) []uint16 {
	n := textFormatInto(nil, 0, format, args...)
	out := make([]uint16, n+1)
	textFormatInto(&out[0], len(out), format, args...)
	// Embedded %c NUL ends the old downstream C-string reader.
	for i, c := range out {
		if c == 0 {
			return out[:i]
		}
	}
	return out[:n]
}
func textFormatLine(u *server.Object, format *uint16, args ...textFormatArgument) int {
	if u == nil || uint32(u.ObjClass)&4 == 0 {
		return gameplayReportPtr(u.CObj())
	}
	return gameplayTextLine(u, textFormatUnits(format, args...))
}
func textFormatAll(flags byte, format *uint16, args ...textFormatArgument) int {
	return gameplayTextAll(flags, textFormatUnits(format, args...))
}
