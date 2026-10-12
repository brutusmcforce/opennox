package legacy

import (
	"math"
	"strconv"
	"strings"
)

func resourceSpace(c byte) bool { return c == ' ' || c >= '\t' && c <= '\r' }
func resourceScanInt(s string) (int32, int, bool) {
	i := 0
	for i < len(s) && resourceSpace(s[i]) {
		i++
	}
	start := i
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	digits := i
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == digits {
		return 0, i, false
	}
	v, _ := strconv.ParseInt(s[start:i], 10, 32) // ErrRange carries libc's saturated endpoint.
	return int32(v), i, true
}
func resourceFloatDigit(c byte, hex bool) bool {
	return c >= '0' && c <= '9' || hex && (c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F')
}
func resourceScanFloat(s string, bits int) (float64, int, bool) {
	return resourceParseFloat(s, bits, true)
}
func resourceParseFloat(s string, bits int, scan bool) (float64, int, bool) {
	i := 0
	for i < len(s) && resourceSpace(s[i]) {
		i++
	}
	start := i
	neg := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		neg = s[i] == '-'
		i++
	}
	rest := s[i:]
	if resourceASCIIPrefix(rest, "inf") {
		n := 3
		if resourceASCIIPrefix(rest, "infinity") {
			n = 8
		} else if scan && len(rest) > 3 && (rest[3] == 'i' || rest[3] == 'I') {
			return 0, i, false
		}
		v := math.Inf(1)
		if neg {
			v = -v
		}
		return v, i + n, true
	}
	if resourceASCIIPrefix(rest, "nan") {
		end := i + 3
		payload := ""
		if end < len(s) && s[end] == '(' {
			begin := end + 1
			pos := begin
			for pos < len(s) && (s[pos] >= '0' && s[pos] <= '9' || s[pos] >= 'a' && s[pos] <= 'z' || s[pos] >= 'A' && s[pos] <= 'Z' || s[pos] == '_') {
				pos++
			}
			if pos < len(s) && s[pos] == ')' {
				payload = s[begin:pos]
				end = pos + 1
			} else if scan {
				return 0, pos, false
			}
		}
		return resourceNaN(payload, neg, bits), end, true
	}
	hex := i+1 < len(s) && s[i] == '0' && (s[i+1] == 'x' || s[i+1] == 'X')
	if hex {
		i += 2
	}
	digits := 0
	for i < len(s) && resourceFloatDigit(s[i], hex) {
		i++
		digits++
	}
	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) && resourceFloatDigit(s[i], hex) {
			i++
			digits++
		}
	}
	if digits == 0 {
		return 0, i, false
	}
	mantissaEnd := i
	exponent := false
	validExponent := false
	if i < len(s) && (!hex && (s[i] == 'e' || s[i] == 'E') || hex && (s[i] == 'p' || s[i] == 'P')) {
		exponent = true
		i++
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			i++
		}
		begin := i
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		validExponent = i > begin
	}
	token := s[start:i]
	if exponent && !validExponent {
		if scan {
			return 0, i, false
		}
		token = s[start:mantissaEnd]
	}
	if hex && !validExponent {
		token += "p0"
	}
	v, err := strconv.ParseFloat(token, bits)
	if err != nil {
		if n, ok := err.(*strconv.NumError); !ok || n.Err != strconv.ErrRange {
			return 0, i, false
		}
	}
	return v, i, true
}
func resourceAtoi(s string) int32   { v, _, _ := resourceScanInt(s); return v }
func resourceAtof(s string) float64 { v, _, _ := resourceParseFloat(s, 64, false); return v }

// libc accepts decimal, octal and hexadecimal NaN payloads, but not Go's
// underscore separators or binary/octal prefix extensions.
func resourceNaN(payload string, negative bool, bits int) float64 {
	var value uint64
	if payload != "" && !strings.Contains(payload, "_") {
		base := 10
		digits := payload
		if len(payload) > 1 && payload[0] == '0' {
			base = 8
			if payload[1] == 'x' || payload[1] == 'X' {
				base = 16
				digits = payload[2:]
			}
		}
		n, err := strconv.ParseUint(digits, base, 64)
		if err == nil {
			value = n
		} else if e, ok := err.(*strconv.NumError); ok && e.Err == strconv.ErrRange {
			value = n
		}
	}
	if bits == 32 {
		word := uint32(value)&0x7fffff | 0x7fc00000
		if negative {
			word |= 0x80000000
		}
		return float64(math.Float32frombits(word))
	}
	word := value&0xfffffffffffff | 0x7ff8000000000000
	if negative {
		word |= 0x8000000000000000
	}
	return math.Float64frombits(word)
}

func resourceASCIIPrefix(s, lower string) bool {
	if len(s) < len(lower) {
		return false
	}
	for i := range len(lower) {
		b := s[i]
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		if b != lower[i] {
			return false
		}
	}
	return true
}
