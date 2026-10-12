// Package mp3 contains Go decoder internals translated from minimp3.
//
// The upstream minimp3 authors dedicated the software to the public domain
// worldwide, to the extent possible under law, and provide it without warranty
// (CC0 1.0 Universal). This attribution is retained from minimp3.h.
// Upstream: https://github.com/lieff/minimp3
// License: http://creativecommons.org/publicdomain/zero/1.0/
package mp3

const maxScalefacIndex = 44

var scfPartitions = [3][28]uint8{
	{6, 5, 5, 5, 6, 5, 5, 5, 6, 5, 7, 3, 11, 10, 0, 0, 7, 7, 7, 0, 6, 6, 6, 3, 8, 8, 5, 0},
	{8, 9, 6, 12, 6, 9, 9, 9, 6, 9, 12, 6, 15, 18, 0, 0, 6, 15, 12, 0, 6, 12, 9, 6, 6, 18, 9, 0},
	{9, 9, 6, 12, 9, 9, 9, 9, 9, 9, 12, 6, 18, 18, 0, 0, 12, 12, 12, 0, 12, 9, 9, 6, 15, 12, 9, 0},
}

var scfCompressDecode = [16]uint8{0, 1, 2, 3, 12, 5, 6, 7, 9, 10, 11, 13, 14, 15, 18, 19}

var scfMod = [24]uint8{
	5, 5, 4, 4, 5, 5, 4, 1, 4, 3, 1, 1,
	5, 6, 6, 1, 4, 4, 4, 1, 4, 3, 1, 1,
}

var scfPreamp = [10]uint8{1, 1, 1, 1, 2, 2, 3, 3, 3, 2}

var scfExpFrac = [4]float32{
	9.31322575e-10,
	7.83145814e-10,
	6.58544508e-10,
	5.53767716e-10,
}

// readScalefactors translates L3_read_scalefactors. It retains the source's
// byte wrapping for sentinel values and leaves output tails untouched.
func readScalefactors(scf *[40]uint8, istPos *[39]uint8, scfSize [4]uint8, scfCount [4]uint8, bitbuf *bitReader, scfsi int32) {
	scfOffset, istOffset := 0, 0
	for i := 0; i < 4 && scfCount[i] != 0; i++ {
		count := int(scfCount[i])
		if scfsi&8 != 0 {
			copy(scf[scfOffset:scfOffset+count], istPos[istOffset:istOffset+count])
		} else {
			bits := int(scfSize[i])
			if bits == 0 {
				clear(scf[scfOffset : scfOffset+count])
				clear(istPos[istOffset : istOffset+count])
			} else {
				maxScf := int32(-1)
				if scfsi < 0 {
					maxScf = (1 << uint(bits)) - 1
				}
				for k := 0; k < count; k++ {
					s := int32(getBits(bitbuf, bits))
					if s == maxScf {
						istPos[istOffset+k] = 0xff
					} else {
						istPos[istOffset+k] = uint8(s)
					}
					scf[scfOffset+k] = uint8(s)
				}
			}
		}
		istOffset += count
		scfOffset += count
		scfsi *= 2
	}
	scf[scfOffset] = 0
	scf[scfOffset+1] = 0
	scf[scfOffset+2] = 0
}

// ldexpQ2 mirrors L3_ldexp_q2's repeated float32 scaling. The caller domain
// supplies nonnegative exponents; the operation order and intermediate
// float32 rounding are kept explicit for the 386 SSE2 target.
func ldexpQ2(y float32, expQ2 int32) float32 {
	for {
		e := int32(120)
		if expQ2 < e {
			e = expQ2
		}
		frac := scfExpFrac[uint32(e)&3]
		scaleInt := int32(1<<30) >> uint(e>>2)
		scale := float32(float32(frac) * float32(scaleInt))
		y = float32(y * scale)
		expQ2 -= e
		if expQ2 <= 0 {
			return y
		}
	}
}

// decodeScalefactors translates L3_decode_scalefactors for the caller's
// MPEG-1/MPEG-2 Layer III domain. Existing istPos and scf tails are preserved.
func decodeScalefactors(hdr [4]byte, istPos *[39]uint8, bs *bitReader, gr *grInfo, scf *[40]float32, ch int) {
	partitionIndex := 0
	if gr.nShortSFB != 0 {
		partitionIndex++
	}
	if gr.nLongSFB == 0 {
		partitionIndex++
	}
	scfPartition := &scfPartitions[partitionIndex]
	partitionOffset := 0
	var scfSize [4]uint8
	var iscf [40]uint8
	scfShift := int32(gr.scalefacScale) + 1
	scfsI := int32(gr.scfsi)

	if hdr[1]&0x08 != 0 {
		part := scfCompressDecode[gr.scalefacCompress]
		scfSize[0] = part >> 2
		scfSize[1] = part >> 2
		scfSize[2] = part & 3
		scfSize[3] = part & 3
	} else {
		ist := int32(0)
		if hdr[3]&0x10 != 0 && ch != 0 {
			ist = 1
		}
		sfc := int32(gr.scalefacCompress) >> uint(ist)
		k := int(ist * 12)
		for sfc >= 0 {
			modprod := int32(1)
			for i := 3; i >= 0; i-- {
				mod := int32(scfMod[k+i])
				scfSize[i] = uint8((sfc / modprod) % mod)
				modprod *= mod
			}
			sfc -= modprod
			k += 4
		}
		partitionOffset = k
		scfsI = -16
	}
	scfCount := [4]uint8(scfPartition[partitionOffset:])
	readScalefactors(&iscf, istPos, scfSize, scfCount, bs, scfsI)

	if gr.nShortSFB != 0 {
		sh := int32(3) - scfShift
		for i := 0; i < int(gr.nShortSFB); i += 3 {
			for win := 0; win < 3; win++ {
				idx := int(gr.nLongSFB) + i + win
				iscf[idx] += uint8(uint32(gr.subblockGain[win]) << uint(sh))
			}
		}
	} else if gr.preflag != 0 {
		for i, value := range scfPreamp {
			iscf[11+i] += value
		}
	}

	gainExp := int32(gr.globalGain) + (-1)*4 - 210
	if hdr[3]&0xe0 == 0x60 {
		gainExp -= 2
	}
	gain := ldexpQ2(1<<uint(maxScalefacIndex/4), maxScalefacIndex-gainExp)
	count := int(gr.nLongSFB) + int(gr.nShortSFB)
	for i := 0; i < count; i++ {
		exp := int32(iscf[i]) << uint(scfShift)
		scf[i] = ldexpQ2(gain, exp)
	}
}
