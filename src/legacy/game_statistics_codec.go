package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

var statisticsRandomA, statisticsRandomB uint32

func statisticsRuns(data []byte) []byte {
	var counts [256]int
	for _, c := range data {
		counts[c]++
	}
	marker, least := byte(0), counts[0]
	for i, n := range counts {
		if n == 0 {
			marker = byte(i)
			break
		}
		if n < least {
			marker = byte(i)
			least = n
		}
	}
	out := []byte{marker}
	if len(data) == 0 {
		return out
	}
	last, count := data[0], byte(1)
	for i := 1; i <= len(data); i++ {
		// The original compares an unsigned input byte with a signed char.
		if i == len(data) || int(data[i]) != int(int8(last)) {
			if count <= 3 {
				for j := byte(0); j < count; j++ {
					if last == marker {
						out = append(out, last)
					}
					out = append(out, last)
				}
			} else {
				out = append(out, marker, count, last)
			}
			if i < len(data) {
				last = data[i]
			}
			count = 1
		} else {
			count++
		}
	}
	return out
}
func statisticsRandom(seed *int32) float64 {
	table := memmap.PtrInt32(0x5D4594, 741384)
	at := func(i uint32) *int32 { return statisticsI32(unsafe.Pointer(table), int(i)*4) }
	initialized := memmap.PtrUint32(0x5D4594, 741656)
	if *seed < 0 || *initialized == 0 {
		*initialized = 1
		v := *seed
		if v < 0 {
			v = -v
		}
		current := (int32(161803398) - v) % 1000000000
		*at(55) = current
		next := int32(1)
		for j := uint32(21); j <= 1134; j += 21 {
			*at(j % 55) = next
			next = current - next
			if next < 0 {
				next += 1000000000
			}
			current = *at(j % 55)
		}
		for pass := 0; pass < 4; pass++ {
			for i := uint32(1); i <= 55; i++ {
				v := *at(i) - *at(1 + (i+30)%55)
				if v < 0 {
					v += 1000000000
				}
				*at(i) = v
			}
		}
		statisticsRandomA = 0
		statisticsRandomB = 31
		*seed = 1
	}
	statisticsRandomA++
	if statisticsRandomA == 56 {
		statisticsRandomA = 1
	}
	statisticsRandomB++
	if statisticsRandomB == 56 {
		statisticsRandomB = 1
	}
	v := *at(statisticsRandomA) - *at(statisticsRandomB)
	if v < 0 {
		v += 1000000000
	}
	*at(statisticsRandomA) = v
	return float64(v) * memmap.Float64(0x581450, 8368)
}
func statisticsRandomBytes(seed int32, count int) []byte {
	if seed > 0 {
		seed = -seed
	}
	v := statisticsRandom(&seed)
	out := make([]byte, count)
	for i := range out {
		n := int32(v * 255)
		if n < 0 {
			n = -n
		}
		out[i] = byte(n)
		v = statisticsRandom(&seed)
	}
	return out
}
func statisticsFrame(data []byte) ([]byte, int32) {
	if len(data) < 15 {
		return nil, -2
	}
	seed := int32(statisticsTime())
	if seed > 0 {
		seed = -seed
	}
	stream := statisticsRandomBytes(seed, len(data))
	next := int32(statisticsTime())
	if next > 0 {
		next = -next
	}
	span := len(data) - 14
	if len(data) >= 241 {
		span = 241
	}
	position := byte(int64(statisticsRandom(&next)*float64(span)) + 10)
	out := make([]byte, len(data)+5)
	j := 0
	for i, v := range data {
		if j == 5 {
			out[j] = position
			j++
		}
		if j == int(position) {
			binary.BigEndian.PutUint32(out[j:], uint32(seed))
			j += 4
		}
		out[j] = v ^ stream[i]
		j++
	}
	return out, int32(len(out))
}
func statisticsEncode(data []byte) ([]byte, int32) {
	frame, n := statisticsFrame(statisticsRuns(data))
	if frame == nil {
		return nil, n
	}
	records := statisticsRecords{}
	records.add(20, statisticsTag(71480), 0, frame)
	out := records.serialize()
	return out, int32(len(out))
}
