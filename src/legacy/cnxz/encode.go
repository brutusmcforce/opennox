package cnxz

import (
	"math/bits"
	"sort"
)

type mapCode struct {
	bits  int
	value uint32
}
type mapEncoder struct {
	window    [65536]byte
	positions [32768]uint16
	position  uint32
	counts    [274]int16
	codes     [274]mapCode
	remaining int
	output    []byte
	word      uint32
	nbits     int
}

func newMapEncoder() *mapEncoder {
	e := &mapEncoder{remaining: 4096}
	for i := range e.positions {
		e.positions[i] = 65535
	}
	// Invert the initial wire-format symbol order, also used by the decoder.
	d := newMapDecoder(nil)
	for group, g := range d.groups {
		for j := 0; j < 1<<g.bits && g.offset+j < len(d.symbols); j++ {
			e.codes[d.symbols[g.offset+j]] = mapCode{4 + g.bits, uint32(group<<g.bits | j)}
		}
	}
	return e
}
func (e *mapEncoder) put(n int, v uint32) {
	// Each original write is at most 16 bits, with a 16-bit flush threshold.
	for n > 16 {
		e.put(n-16, v>>16)
		n = 16
		v &= 65535
	}
	e.word |= v << (32 - e.nbits - n)
	e.nbits += n
	if e.nbits >= 16 {
		e.output = append(e.output, byte(e.word>>24), byte(e.word>>16))
		e.word <<= 16
		e.nbits -= 16
	}
}
func (e *mapEncoder) symbol(s int) {
	e.remaining--
	if e.remaining <= 0 {
		e.rebuild()
	}
	e.counts[s]++
	c := e.codes[s]
	e.put(c.bits, c.value)
}
func (e *mapEncoder) rebuild() {
	// The rebuild marker itself uses the old table and contributes to the counts.
	e.remaining = 2
	e.symbol(272)
	var order [274]int
	total := 0
	for i, n := range e.counts {
		order[i] = i
		total += int(n)
	}
	sort.Slice(order[:], func(i, j int) bool {
		a, b := order[i], order[j]
		if e.counts[a] != e.counts[b] {
			return e.counts[a] > e.counts[b]
		}
		return a > b
	})
	var frequency [274]int
	for i, s := range order {
		frequency[i] = int(e.counts[s])
	}
	for i := range e.counts {
		e.counts[i] >>= 1
	}
	e.remaining = 4096
	var groups [16]int
	used, weight := 0, 0
	abs := func(v int) int {
		if v < 0 {
			return -v
		}
		return v
	}
	for group := 0; group < 14; group++ {
		target := (total - weight) / (16 - group)
		n, count, sum, previous := 0, 0, 0, 0
		for {
			take := 1 << n
			clipped := used+take+count > 274
			if clipped {
				take = 274 - used
			}
			for count < take {
				sum += frequency[used+count]
				count++
			}
			if clipped || n >= 8 || sum > target {
				break
			}
			previous = sum
			n++
		}
		if n != 0 && abs(previous-target) <= abs(sum-target) {
			sum = previous
			n--
		}
		groups[group] = n
		sort.Ints(groups[:group+1])
		weight += sum
		used += 1 << n
	}
	rest := 0
	for i := used; i < 274; i++ {
		rest += frequency[i]
	}
	count, sum, best := 0, 0, int(0x7fffffff)
	left, right := 0, 0
	for n := 0; used+(1<<n)+count <= 274; n++ {
		for count < 1<<n {
			sum += frequency[used+count]
			count++
		}
		r := 0
		for 1<<r < 274-count-used {
			r++
		}
		if n <= 8 && r <= 8 {
			cost := sum*n + r*(rest-sum)
			if cost >= best {
				break
			}
			best = cost
			left = n
			right = r
		}
	}
	groups[14] = left
	sort.Ints(groups[:15])
	groups[15] = right
	sort.Ints(groups[:])
	offset, prev := 0, 0
	for g, n := range groups {
		for j := 0; j < 1<<n && offset+j < 274; j++ {
			e.codes[order[offset+j]] = mapCode{n + 4, uint32(g<<n | j)}
		}
		offset += 1 << n
		e.put(n-prev+1, 1)
		prev = n
	}
}
func (e *mapEncoder) literals(src []byte) {
	for _, v := range src {
		e.symbol(int(v))
	}
}
func (e *mapEncoder) match(prefix []byte, length, distance int) {
	e.literals(prefix)
	if length < 12 {
		e.symbol(256 + length - 4)
	} else {
		n := bits.Len(uint(length-10)) - 1
		e.symbol(263 + n)
		e.put(n, uint32(length-10-(1<<n)))
	}
	high := distance >> 9
	n, prefixCode := 0, high
	if high > 1 {
		n = bits.Len(uint(high)) - 1
		prefixCode = n + 1
	}
	e.put(3, uint32(prefixCode))
	e.put(n+9, uint32(distance-((high>>n)<<n<<9)))
}
func mapHash(src []byte) uint32 {
	var h uint32
	for _, v := range src[:5] {
		h = bits.RotateLeft32(h^uint32(v), 5)
	}
	return h
}
func mapHashNext(h uint32, out, in byte) uint32 {
	return bits.RotateLeft32(h^uint32(in)^bits.RotateLeft32(uint32(out), 25), 5)
}
func (e *mapEncoder) remember(h uint32, pos uint32) uint16 {
	index := (214013*h + 2531011) >> 17
	old := e.positions[index]
	e.positions[index] = uint16(pos)
	return old
}
func (e *mapEncoder) copyWindow(src []byte) {
	for _, v := range src {
		e.window[uint16(e.position)] = v
		e.position++
	}
}

// source includes the following file bytes and five zero padding bytes so that
// rolling hashes at the end of a block can retain the original lookahead.
func (e *mapEncoder) block(source []byte, size int) []byte {
	e.output = nil
	e.word = 0
	e.nbits = 0
	at := 0
	var hash uint32
	if size >= 5 {
		hash = mapHash(source)
	}
	for size-at >= 5 {
		remaining := size - at
		limit := min(remaining-5, 64)
		i, bestLen, bestAt := 0, 0, 0
		var bestPos uint32
		var bestOld uint16
		var bestHash uint32
		consumed := false
		for {
			old := e.remember(hash, e.position)
			pos := uint16(e.position)
			length := 0
			backed := false
			if old != 65535 && old != pos {
				distance := int(uint16(pos - old))
				bound := min(distance, remaining-i, 521)
				for length < bound && e.window[uint16(int(old)+length)] == source[at+i+length] {
					length++
				}
				if uint16(int(old)+length) == pos {
					bound = min(521, remaining-i)
					for length < bound && source[at+i+length-distance] == source[at+i+length] {
						length++
					}
				}
				if length >= 3 {
					backLimit := min(521-length, i, distance-length, 65536-distance)
					back := 0
					for back < backLimit && e.window[uint16(int(old)-back-1)] == source[at+i-back-1] {
						back++
					}
					if back > 0 {
						i -= back
						e.position -= uint32(back)
						old -= uint16(back)
						length += back
						pos = uint16(e.position)
						hash = mapHash(source[at+i:])
						backed = true
					}
				}
			}
			if bestLen >= 4 || length >= 4 && backed {
				startUpdate := 1
				if bestLen >= 4 && length <= bestLen {
					i = bestAt
					length = bestLen
					e.position = bestPos
					old = bestOld
					hash = bestHash
					startUpdate = 2
				}
				e.match(source[at:at+i], length, int(uint16(uint16(e.position)-old)))
				h := hash
				for j := 0; j < startUpdate; j++ {
					h = mapHashNext(h, source[at+i+j], source[at+i+j+5])
				}
				for j := startUpdate; j < length; j++ {
					e.remember(h, e.position+uint32(j))
					h = mapHashNext(h, source[at+i+j], source[at+i+j+5])
				}
				hash = h
				e.copyWindow(source[at+i : at+i+length])
				at += i + length
				consumed = true
				break
			}
			if length >= 4 {
				bestLen = length
				bestAt = i
				bestPos = e.position
				bestOld = old
				bestHash = hash
			}
			if i+1 > limit {
				break
			}
			hash = mapHashNext(hash, source[at+i], source[at+i+5])
			e.copyWindow(source[at+i : at+i+1])
			i++
		}
		if consumed {
			continue
		}
		if bestLen >= 4 {
			e.match(source[at:at+i], bestLen, int(uint16(uint16(e.position)-bestOld)))
			h := mapHashNext(hash, source[at+i], source[at+i+5])
			for j := 1; j < bestLen; j++ {
				e.remember(h, e.position+uint32(j))
				h = mapHashNext(h, source[at+i+j], source[at+i+j+5])
			}
			hash = h
			e.copyWindow(source[at+i : at+i+bestLen])
			at += i + bestLen
		} else {
			if i+5 >= remaining && remaining <= 64 {
				e.copyWindow(source[at+i : at+remaining])
				i = remaining
			}
			e.literals(source[at : at+i])
			at += i
		}
	}
	if at < size {
		e.copyWindow(source[at:size])
		e.literals(source[at:size])
	}
	e.symbol(273)
	for e.nbits > 0 {
		e.output = append(e.output, byte(e.word>>24))
		e.word <<= 8
		e.nbits -= 8
	}
	return e.output
}
