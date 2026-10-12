// Package mp3 contains Go decoder internals translated from minimp3.
//
// The upstream minimp3 authors dedicated the software to the public domain
// worldwide, to the extent possible under law, and provide it without warranty.
// Upstream: https://github.com/lieff/minimp3
// License: http://creativecommons.org/publicdomain/zero/1.0/
// This attribution is retained from minimp3.h.
// This file reuses bitReader/getBits and header helpers from bits_headers.go.
package mp3

const shortBlockType = 2

// grInfo mirrors minimp3's L3_gr_info_t. sfbTable borrows the exact static row
// selected by C's sfbtab pointer and remains unchanged if parsing returns early.
type grInfo struct {
	sfbTable         []uint8
	part23Length     uint16
	bigValues        uint16
	scalefacCompress uint16
	globalGain       uint8
	blockType        uint8
	mixedBlockFlag   uint8
	nLongSFB         uint8
	nShortSFB        uint8
	tableSelect      [3]uint8
	regionCount      [3]uint8
	subblockGain     [3]uint8
	preflag          uint8
	scalefacScale    uint8
	count1Table      uint8
	scfsi            uint8
}

// The following tables are the static C g_scf_long/short/mixed contents;
// treat them as immutable. The row index is the source's srIdx calculation.
var scfLong = [8][23]uint8{
	{6, 6, 6, 6, 6, 6, 8, 10, 12, 14, 16, 20, 24, 28, 32, 38, 46, 52, 60, 68, 58, 54, 0},
	{12, 12, 12, 12, 12, 12, 16, 20, 24, 28, 32, 40, 48, 56, 64, 76, 90, 2, 2, 2, 2, 2, 0},
	{6, 6, 6, 6, 6, 6, 8, 10, 12, 14, 16, 20, 24, 28, 32, 38, 46, 52, 60, 68, 58, 54, 0},
	{6, 6, 6, 6, 6, 6, 8, 10, 12, 14, 16, 18, 22, 26, 32, 38, 46, 54, 62, 70, 76, 36, 0},
	{6, 6, 6, 6, 6, 6, 8, 10, 12, 14, 16, 20, 24, 28, 32, 38, 46, 52, 60, 68, 58, 54, 0},
	{4, 4, 4, 4, 4, 4, 6, 6, 8, 8, 10, 12, 16, 20, 24, 28, 34, 42, 50, 54, 76, 158, 0},
	{4, 4, 4, 4, 4, 4, 6, 6, 6, 8, 10, 12, 16, 18, 22, 28, 34, 40, 46, 54, 54, 192, 0},
	{4, 4, 4, 4, 4, 4, 6, 6, 8, 10, 12, 16, 20, 24, 30, 38, 46, 56, 68, 84, 102, 26, 0},
}

var scfShort = [8][40]uint8{
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 6, 6, 6, 8, 8, 8, 10, 10, 10, 12, 12, 12, 14, 14, 14, 18, 18, 18, 24, 24, 24, 30, 30, 30, 40, 40, 40, 18, 18, 18, 0},
	{8, 8, 8, 8, 8, 8, 8, 8, 8, 12, 12, 12, 16, 16, 16, 20, 20, 20, 24, 24, 24, 28, 28, 28, 36, 36, 36, 2, 2, 2, 2, 2, 2, 2, 2, 2, 26, 26, 26, 0},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 6, 6, 6, 6, 6, 6, 8, 8, 8, 10, 10, 10, 14, 14, 14, 18, 18, 18, 26, 26, 26, 32, 32, 32, 42, 42, 42, 18, 18, 18, 0},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 6, 6, 6, 8, 8, 8, 10, 10, 10, 12, 12, 12, 14, 14, 14, 18, 18, 18, 24, 24, 24, 32, 32, 32, 44, 44, 44, 12, 12, 12, 0},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 6, 6, 6, 8, 8, 8, 10, 10, 10, 12, 12, 12, 14, 14, 14, 18, 18, 18, 24, 24, 24, 30, 30, 30, 40, 40, 40, 18, 18, 18, 0},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 6, 6, 6, 8, 8, 8, 10, 10, 10, 12, 12, 12, 14, 14, 14, 18, 18, 18, 22, 22, 22, 30, 30, 30, 56, 56, 56, 0},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 6, 6, 6, 6, 6, 6, 10, 10, 10, 12, 12, 12, 14, 14, 14, 16, 16, 16, 20, 20, 20, 26, 26, 26, 66, 66, 66, 0},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 6, 6, 6, 8, 8, 8, 12, 12, 12, 16, 16, 16, 20, 20, 20, 26, 26, 26, 34, 34, 34, 42, 42, 42, 12, 12, 12, 0},
}

var scfMixed = [8][40]uint8{
	{6, 6, 6, 6, 6, 6, 6, 6, 6, 8, 8, 8, 10, 10, 10, 12, 12, 12, 14, 14, 14, 18, 18, 18, 24, 24, 24, 30, 30, 30, 40, 40, 40, 18, 18, 18, 0, 0, 0, 0},
	{12, 12, 12, 4, 4, 4, 8, 8, 8, 12, 12, 12, 16, 16, 16, 20, 20, 20, 24, 24, 24, 28, 28, 28, 36, 36, 36, 2, 2, 2, 2, 2, 2, 2, 2, 2, 26, 26, 26, 0},
	{6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 8, 8, 8, 10, 10, 10, 14, 14, 14, 18, 18, 18, 26, 26, 26, 32, 32, 32, 42, 42, 42, 18, 18, 18, 0, 0, 0, 0},
	{6, 6, 6, 6, 6, 6, 6, 6, 6, 8, 8, 8, 10, 10, 10, 12, 12, 12, 14, 14, 14, 18, 18, 18, 24, 24, 24, 32, 32, 32, 44, 44, 44, 12, 12, 12, 0, 0, 0, 0},
	{6, 6, 6, 6, 6, 6, 6, 6, 6, 8, 8, 8, 10, 10, 10, 12, 12, 12, 14, 14, 14, 18, 18, 18, 24, 24, 24, 30, 30, 30, 40, 40, 40, 18, 18, 18, 0, 0, 0, 0},
	{4, 4, 4, 4, 4, 4, 6, 6, 4, 4, 4, 6, 6, 6, 8, 8, 8, 10, 10, 10, 12, 12, 12, 14, 14, 14, 18, 18, 18, 22, 22, 22, 30, 30, 30, 56, 56, 56, 0, 0},
	{4, 4, 4, 4, 4, 4, 6, 6, 4, 4, 4, 6, 6, 6, 6, 6, 6, 10, 10, 10, 12, 12, 12, 14, 14, 14, 16, 16, 16, 20, 20, 20, 26, 26, 26, 66, 66, 66, 0, 0},
	{4, 4, 4, 4, 4, 4, 6, 6, 4, 4, 4, 6, 6, 6, 8, 8, 8, 12, 12, 12, 16, 16, 16, 20, 20, 20, 26, 26, 26, 34, 34, 34, 42, 42, 42, 12, 12, 12, 0, 0},
}

// readSideInfo is a direct translation of minimp3's Layer III side-info parser.
// hdr must be a valid Layer III header; gr supplies four existing output slots.
// It intentionally does not clear gr, because C leaves some fields untouched.
func readSideInfo(bs *bitReader, gr *[4]grInfo, hdr [4]byte) int32 {
	srIdx := int((hdr[2] >> 2) & 3)
	srIdx += (int((hdr[1]>>3)&1) + int((hdr[1]>>4)&1)) * 3
	if srIdx != 0 {
		srIdx--
	}
	grCount := 2
	if hdr[3]&0xc0 == 0xc0 {
		grCount = 1
	}
	if hdr[1]&0x08 != 0 {
		grCount *= 2
	}
	mainDataBegin := int32(0)
	scfsi := uint32(0)
	part23Sum := int32(0)
	if hdr[1]&0x08 != 0 {
		mainDataBegin = int32(getBits(bs, 9))
		scfsi = getBits(bs, 7+grCount)
	} else {
		mainDataBegin = int32(getBits(bs, 8+grCount) >> uint(grCount))
	}
	for i := 0; i < grCount; i++ {
		g := &gr[i]
		if hdr[3]&0xc0 == 0xc0 {
			scfsi <<= 4
		}
		g.part23Length = uint16(getBits(bs, 12))
		part23Sum += int32(g.part23Length)
		g.bigValues = uint16(getBits(bs, 9))
		if g.bigValues > 288 {
			return -1
		}
		g.globalGain = uint8(getBits(bs, 8))
		if hdr[1]&0x08 != 0 {
			g.scalefacCompress = uint16(getBits(bs, 4))
		} else {
			g.scalefacCompress = uint16(getBits(bs, 9))
		}
		g.sfbTable = scfLong[srIdx][:]
		g.nLongSFB = 22
		g.nShortSFB = 0
		if getBits(bs, 1) != 0 {
			g.blockType = uint8(getBits(bs, 2))
			if g.blockType == 0 {
				return -1
			}
			g.mixedBlockFlag = uint8(getBits(bs, 1))
			g.regionCount[0] = 7
			g.regionCount[1] = 255
			if g.blockType == shortBlockType {
				scfsi &= 0x0f0f
				if g.mixedBlockFlag == 0 {
					g.regionCount[0] = 8
					g.sfbTable = scfShort[srIdx][:]
					g.nLongSFB = 0
					g.nShortSFB = 39
				} else {
					g.sfbTable = scfMixed[srIdx][:]
					if hdr[1]&0x08 != 0 {
						g.nLongSFB = 8
					} else {
						g.nLongSFB = 6
					}
					g.nShortSFB = 30
				}
			}
			tables := getBits(bs, 10) << 5
			g.subblockGain[0] = uint8(getBits(bs, 3))
			g.subblockGain[1] = uint8(getBits(bs, 3))
			g.subblockGain[2] = uint8(getBits(bs, 3))
			g.tableSelect[0] = uint8(tables >> 10)
			g.tableSelect[1] = uint8((tables >> 5) & 31)
			g.tableSelect[2] = uint8(tables & 31)
		} else {
			g.blockType = 0
			g.mixedBlockFlag = 0
			tables := getBits(bs, 15)
			g.regionCount[0] = uint8(getBits(bs, 4))
			g.regionCount[1] = uint8(getBits(bs, 3))
			g.regionCount[2] = 255
			g.tableSelect[0] = uint8(tables >> 10)
			g.tableSelect[1] = uint8((tables >> 5) & 31)
			g.tableSelect[2] = uint8(tables & 31)
		}
		if hdr[1]&0x08 != 0 {
			g.preflag = uint8(getBits(bs, 1))
		} else if g.scalefacCompress >= 500 {
			g.preflag = 1
		} else {
			g.preflag = 0
		}
		g.scalefacScale = uint8(getBits(bs, 1))
		g.count1Table = uint8(getBits(bs, 1))
		g.scfsi = uint8((scfsi >> 12) & 15)
		scfsi <<= 4
	}
	if part23Sum+int32(bs.pos) > int32(bs.limit)+mainDataBegin*8 {
		return -1
	}
	return mainDataBegin
}
