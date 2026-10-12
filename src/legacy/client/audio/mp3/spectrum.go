// Package mp3 contains Go decoder internals translated from minimp3.
//
// The upstream minimp3 authors dedicated the software to the public domain
// worldwide, to the extent possible under law, and provide it without warranty
// (CC0 1.0 Universal). This attribution is retained from minimp3.h.
// Upstream: https://github.com/lieff/minimp3
// License: http://creativecommons.org/publicdomain/zero/1.0/
package mp3

var stereoPan = [14]float32{
	0, 1,
	0.21132487, 0.78867513,
	0.36602540, 0.63397460,
	0.5, 0.5,
	0.63397460, 0.36602540,
	0.78867513, 0.21132487,
	1, 0,
}

var antialiasCoefficients = [2][8]float32{
	{0.85749293, 0.88174200, 0.94962865, 0.98331459, 0.99551782, 0.99916056, 0.99989920, 0.99999316},
	{0.51449576, 0.47173197, 0.31337745, 0.18191320, 0.09457419, 0.04096558, 0.01419856, 0.00369997},
}

// midsideStereo translates the scalar L3_midside_stereo branch. left contains
// both contiguous 576-sample channels; only the first n samples are changed.
func midsideStereo(left []float32, n int) {
	right := 576
	for i := 0; i < n; i++ {
		a, b := left[i], left[right+i]
		left[i] = float32(a + b)
		left[right+i] = float32(a - b)
	}
}

// intensityStereoBand translates L3_intensity_stereo_band, retaining each
// source sample before either channel slot is overwritten.
func intensityStereoBand(left []float32, n int, kl, kr float32) {
	for i := 0; i < n; i++ {
		x := left[i]
		left[i+576] = float32(x * kr)
		left[i] = float32(x * kl)
	}
}

// stereoTopBand translates L3_stereo_top_band and returns its three
// modulo-three last-nonzero band indices, initialized to -1.
func stereoTopBand(right []float32, sfb []uint8, nBands int) [3]int {
	maxBand := [3]int{-1, -1, -1}
	off := 0
	for i := 0; i < nBands; i++ {
		width := int(sfb[i])
		for k := 0; k < width; k += 2 {
			if right[off+k] != 0 || right[off+k+1] != 0 {
				maxBand[i%3] = i
				break
			}
		}
		off += width
	}
	return maxBand
}

// stereoProcess translates scalar L3_stereo_process. sfb includes the source
// zero terminator; left holds left and right channel data in adjacent 576-word
// spans. msStereo intentionally checks the raw mode-extension bit 0x20.
func stereoProcess(left []float32, istPos, sfb []uint8, hdr [4]byte, maxBand [3]int, mpeg2Shift int) {
	mpeg1 := hdr[1]&0x08 != 0
	msStereo := hdr[3]&0x20 != 0
	maxPos := uint8(64)
	if mpeg1 {
		maxPos = 7
	}
	off := 0
	for i := 0; sfb[i] != 0; i++ {
		width := int(sfb[i])
		ipos := istPos[i]
		if i > maxBand[i%3] && ipos < maxPos {
			kl, kr := float32(1), float32(1)
			s := float32(1)
			if msStereo {
				s = 1.41421356
			}
			if mpeg1 {
				kl = stereoPan[2*int(ipos)]
				kr = stereoPan[2*int(ipos)+1]
			} else {
				exp := int32((uint32(ipos)+1)>>1) << uint(mpeg2Shift)
				kr = ldexpQ2(1, exp)
				if ipos&1 != 0 {
					kl = kr
					kr = 1
				}
			}
			kl = float32(kl * s)
			kr = float32(kr * s)
			intensityStereoBand(left[off:], width, kl, kr)
		} else if msStereo {
			midsideStereo(left[off:], width)
		}
		off += width
	}
}

// intensityStereo translates L3_intensity_stereo, including its updates to
// the final intensity-position entries and its gr[1] scalefactor lookup.
func intensityStereo(left []float32, istPos []uint8, gr *[2]grInfo, hdr [4]byte) {
	nSFB := int(gr[0].nLongSFB) + int(gr[0].nShortSFB)
	maxBlocks := 1
	if gr[0].nShortSFB != 0 {
		maxBlocks = 3
	}
	maxBand := stereoTopBand(left[576:], gr[0].sfbTable, nSFB)
	if gr[0].nLongSFB != 0 {
		v := maxBand[0]
		if maxBand[1] > v {
			v = maxBand[1]
		}
		if maxBand[2] > v {
			v = maxBand[2]
		}
		maxBand = [3]int{v, v, v}
	}
	defaultPos := uint8(0)
	if hdr[1]&0x08 != 0 {
		defaultPos = 3
	}
	for i := 0; i < maxBlocks; i++ {
		iTop := nSFB - maxBlocks + i
		prev := iTop - maxBlocks
		if maxBand[i] >= prev {
			istPos[iTop] = defaultPos
		} else {
			istPos[iTop] = istPos[prev]
		}
	}
	stereoProcess(left, istPos, gr[0].sfbTable, hdr, maxBand, int(gr[1].scalefacCompress&1))
}

// reorder translates L3_reorder. The input and scratch spans must be disjoint,
// as required by the source memcpy; sfb ends at its zero terminator. grbuf is
// sized for the sum of the selected band widths. For the MPEG-2.5 mixed-block
// caller, the suffix sums to 528 (within one 576-float row), but the caller's
// starting offset 72 makes the written channel-relative extent end at 600; see
// README.md for the cross-field consequence.
func reorder(grbuf, scratch []float32, sfb []uint8) {
	srcOff, dstOff := 0, 0
	for band := 0; sfb[band] != 0; band += 3 {
		length := int(sfb[band])
		for i := 0; i < length; i++ {
			scratch[dstOff] = grbuf[srcOff+i]
			scratch[dstOff+1] = grbuf[srcOff+length+i]
			scratch[dstOff+2] = grbuf[srcOff+2*length+i]
			dstOff += 3
		}
		srcOff += 3 * length
	}
	copy(grbuf[:dstOff], scratch[:dstOff])
}

// antialias translates the scalar L3_antialias branch. grbuf has the
// contiguous subband rows required by nbands, including the following row.
func antialias(grbuf []float32, nBands int) {
	for band := 0; band < nBands; band++ {
		row := 18 * band
		for i := 0; i < 8; i++ {
			u, d := grbuf[row+18+i], grbuf[row+17-i]
			upper0 := float32(u * antialiasCoefficients[0][i])
			lower1 := float32(d * antialiasCoefficients[1][i])
			upper1 := float32(u * antialiasCoefficients[1][i])
			lower0 := float32(d * antialiasCoefficients[0][i])
			grbuf[row+18+i] = float32(upper0 - lower1)
			grbuf[row+17-i] = float32(upper1 + lower0)
		}
	}
}
