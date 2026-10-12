// Layer III orchestration translated from scalar minimp3.
// The upstream minimp3 authors dedicated this software to the public domain
// worldwide, to the extent possible under law, and provide it without warranty
// (CC0 1.0 Universal): https://creativecommons.org/publicdomain/zero/1.0/
// This attribution is retained from minimp3.h.
package mp3

// decodeScratch adds the C decoder's Layer III working regions. spectral keeps
// both 576-float channels immediately before the 40-float scalefactor area;
// synthesis retains the existing C scratch extent and supplies reorder scratch.
type decodeScratch struct {
	scratchState
	grInfo    [4]grInfo
	istPos    [2][39]uint8
	spectral  [1192]float32
	synthesis [2112]float32
}

// decodeLayer3 follows L3_decode for the selected granule's gr_info slice.
// gr must begin at the corresponding C gr_info entry; in stereo it contains
// both channels, and retains the backing entries needed by intensity stereo.
// nch is the source's one- or two-channel value.
func decodeLayer3(dec *decoderState, s *decodeScratch, gr []grInfo, nch int) {
	scf := (*[40]float32)(s.spectral[1152:])
	for ch := 0; ch < nch; ch++ {
		layer3Limit := s.bs.pos + int(gr[ch].part23Length)
		decodeScalefactors(dec.header, &s.istPos[ch], &s.bs, &gr[ch], scf, ch)
		huffman(s.spectral[ch*576:ch*576+576], &s.bs, &gr[ch], scf[:], layer3Limit)
	}

	if dec.header[3]&0x10 != 0 {
		intensityStereo(s.spectral[:1152], s.istPos[1][:], (*[2]grInfo)(gr), dec.header)
	} else if dec.header[3]&0xe0 == 0x60 {
		midsideStereo(s.spectral[:1152], 576)
	}

	sampleRateIndex := int((dec.header[2] >> 2) & 3)
	sampleRateIndex += (int((dec.header[1]>>3)&1) + int((dec.header[1]>>4)&1)) * 3
	for ch := 0; ch < nch; ch++ {
		aaBands := 31
		nLongBands := 0
		if gr[ch].mixedBlockFlag != 0 {
			nLongBands = 2
			if sampleRateIndex == 2 {
				nLongBands <<= 1
			}
		}
		grbuf := s.spectral[ch*576:]
		if gr[ch].nShortSFB != 0 {
			aaBands = nLongBands - 1
			reorder(grbuf[nLongBands*18:], s.synthesis[:576], gr[ch].sfbTable[int(gr[ch].nLongSFB):])
		}
		antialias(grbuf, aaBands)
		imdctGranule(grbuf, dec.mdctOverlap[ch][:], int(gr[ch].blockType), nLongBands)
		changeSign(grbuf)
	}
}
