// Package mp3 contains Go decoder internals translated from minimp3.
//
// The upstream minimp3 authors dedicated the software to the public domain
// worldwide, to the extent possible under law, and provide it without warranty.
// Upstream: https://github.com/lieff/minimp3
// License: http://creativecommons.org/publicdomain/zero/1.0/
// This attribution is retained from minimp3.h.
package mp3

type bitReader struct {
	buf   []byte
	pos   int
	limit int
}

var headerHalfRateKbps = [2][3][15]uint8{
	{
		{0, 4, 8, 12, 16, 20, 24, 28, 32, 40, 48, 56, 64, 72, 80},
		{0, 4, 8, 12, 16, 20, 24, 28, 32, 40, 48, 56, 64, 72, 80},
		{0, 16, 24, 28, 32, 40, 48, 56, 64, 72, 80, 88, 96, 112, 128},
	},
	{
		{0, 16, 20, 24, 28, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160},
		{0, 16, 24, 28, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192},
		{0, 16, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224},
	},
}

func bsInit(bs *bitReader, data []byte, bytes int) {
	if bytes < 0 || bytes > len(data) {
		panic("mp3: invalid bit-reader byte count")
	}
	bs.buf = data[:bytes]
	bs.pos = 0
	bs.limit = bytes * 8
}

// getBits follows minimp3's MSB-first bit extraction and position update.
// Inputs outside n in [0,32], or malformed reader state, are unsupported.
func getBits(bs *bitReader, n int) uint32 {
	if n < 0 || n > 32 || bs.pos < 0 || bs.limit < 0 || bs.limit > len(bs.buf)*8 {
		panic("mp3: unsupported bit-reader state or width")
	}
	s := uint32(bs.pos & 7)
	shl := n + int(s)
	bs.pos += n
	if bs.pos > bs.limit {
		return 0
	}
	if n == 0 {
		return 0
	}

	// The C implementation starts from the byte containing the old bit position.
	p := (bs.pos - n) >> 3
	next := uint32(bs.buf[p] & (255 >> s))
	cache := uint32(0)
	for {
		shl -= 8
		if shl <= 0 {
			break
		}
		cache |= next << uint(shl)
		p++
		next = uint32(bs.buf[p])
	}
	return cache | (next >> uint(-shl))
}

func headerValid(h [4]byte) bool {
	layer := (h[1] >> 1) & 3
	bitrate := h[2] >> 4
	sampleRate := (h[2] >> 2) & 3
	return h[0] == 0xff && ((h[1]&0xf0) == 0xf0 || (h[1]&0xfe) == 0xe2) &&
		layer != 0 && bitrate != 15 && sampleRate != 3
}

func headerCompare(h1, h2 [4]byte) bool {
	free1 := (h1[2] & 0xf0) == 0
	free2 := (h2[2] & 0xf0) == 0
	return headerValid(h2) && ((h1[1]^h2[1])&0xfe) == 0 &&
		((h1[2]^h2[2])&0x0c) == 0 && free1 == free2
}

func headerBitrateKbps(h [4]byte) uint32 {
	mpeg1 := 0
	if h[1]&0x08 != 0 {
		mpeg1 = 1
	}
	layer := int((h[1] >> 1) & 3)
	bitrate := int(h[2] >> 4)
	return 2 * uint32(headerHalfRateKbps[mpeg1][layer-1][bitrate])
}

func headerSampleRateHz(h [4]byte) uint32 {
	base := [3]uint32{44100, 48000, 32000}
	idx := (h[2] >> 2) & 3
	shift := uint32(0)
	if h[1]&0x08 == 0 {
		shift++
	}
	if h[1]&0x10 == 0 {
		shift++
	}
	return base[idx] >> shift
}

func headerFrameSamples(h [4]byte) uint32 {
	if h[1]&0x06 == 0x06 {
		return 384
	}
	n := uint32(1152)
	if h[1]&0x0e == 0x02 {
		n >>= 1
	}
	return n
}

func headerFrameBytes(h [4]byte, freeFormatBytes int32) int32 {
	frameBytes := headerFrameSamples(h) * headerBitrateKbps(h) * 125 / headerSampleRateHz(h)
	if h[1]&0x06 == 0x06 {
		frameBytes &= ^uint32(3)
	}
	if frameBytes == 0 {
		return freeFormatBytes
	}
	return int32(frameBytes)
}

func headerPadding(h [4]byte) int32 {
	if h[2]&0x02 == 0 {
		return 0
	}
	if h[1]&0x06 == 0x06 {
		return 4
	}
	return 1
}
