package legacy

import "math/bits"

// browserIPv4 preserves inet_addr's legacy numeric forms and 386 result byte order.
// A suffix beginning with ASCII whitespace is accepted, even if text follows it.
func browserIPv4(s string) uint32 {
	const invalid = ^uint32(0)
	var address uint32
	parts, pos := 0, 0
	for {
		if pos == len(s) || s[pos] < '0' || s[pos] > '9' {
			return invalid
		}
		base := uint32(10)
		digits := 0
		if s[pos] == '0' {
			base = 8
			pos++
			digits++
			if pos < len(s) && (s[pos] == 'x' || s[pos] == 'X') {
				base = 16
				pos++
				digits = 0
			}
		}
		var value uint32
		for pos < len(s) {
			c := s[pos]
			var digit uint32
			switch {
			case c >= '0' && c <= '9':
				digit = uint32(c - '0')
			case c >= 'a' && c <= 'f':
				digit = uint32(c-'a') + 10
			case c >= 'A' && c <= 'F':
				digit = uint32(c-'A') + 10
			default:
				digit = 16
			}
			if digit >= base {
				break
			}
			if value > (invalid-digit)/base {
				return invalid
			}
			value = value*base + digit
			pos++
			digits++
		}
		if digits == 0 {
			return invalid
		}
		if pos < len(s) && s[pos] == '.' {
			if parts == 3 || value > 255 {
				return invalid
			}
			address |= value << uint(24-8*parts)
			parts++
			pos++
			continue
		}
		if pos < len(s) && s[pos] != 0 && !(s[pos] == ' ' || s[pos] >= '\t' && s[pos] <= '\r') {
			return invalid
		}
		if value > invalid>>uint(8*parts) {
			return invalid
		}
		return bits.ReverseBytes32(address | value)
	}
}
