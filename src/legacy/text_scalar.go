package legacy

import "unsafe"

func textFold(c uint16) uint16 {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

func textCompareWide(a, b *uint16) int {
	for i := 0; ; i++ {
		x, y := *(*uint16)(unsafe.Add(unsafe.Pointer(a), i*2)), *(*uint16)(unsafe.Add(unsafe.Pointer(b), i*2))
		fx, fy := textFold(x), textFold(y)
		if x == 0 || fx != fy {
			return int(fx) - int(fy)
		}
	}
}

func textCompareNarrow(a, b *byte) int {
	for i := 0; ; i++ {
		x := *(*byte)(unsafe.Add(unsafe.Pointer(a), i))
		y := *(*byte)(unsafe.Add(unsafe.Pointer(b), i))
		if x == 0 || y == 0 || x != y && textFold(uint16(x)) != textFold(uint16(y)) {
			d := int(textFold(uint16(x))) - int(textFold(uint16(y)))
			if d < -1 {
				return -1
			}
			if d > 1 {
				return 1
			}
			return d
		}
	}
}

func textDecimal(p *uint16) int32 {
	const maxInt32 = int64(1<<31 - 1)
	const minInt32 = -1 << 31

	i := 0
	for {
		c := *(*uint16)(unsafe.Add(unsafe.Pointer(p), i*2))
		if c != ' ' && (c < '\t' || c > '\r') {
			break
		}
		i++
	}
	neg := false
	c := *(*uint16)(unsafe.Add(unsafe.Pointer(p), i*2))
	if c == '+' || c == '-' {
		neg = c == '-'
		i++
	}
	limit := maxInt32
	if neg {
		limit = -minInt32
	}
	var n int64
	any := false
	for {
		c = *(*uint16)(unsafe.Add(unsafe.Pointer(p), i*2))
		if c < '0' || c > '9' {
			break
		}
		any = true
		d := int64(c - '0')
		if n > (limit-d)/10 {
			n = limit
			// Continue consuming digits; the saturated value remains fixed.
		} else {
			n = n*10 + d
		}
		i++
	}
	if !any {
		return 0
	}
	if neg {
		return int32(-n)
	}
	return int32(n)
}

func textCopy(dst, src *uint16) *uint16 {
	for i := 0; ; i++ {
		c := *(*uint16)(unsafe.Add(unsafe.Pointer(src), i*2))
		*(*uint16)(unsafe.Add(unsafe.Pointer(dst), i*2)) = c
		if c == 0 {
			return dst
		}
	}
}
