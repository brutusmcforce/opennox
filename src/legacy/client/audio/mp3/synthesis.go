// PCM synthesis translated from the scalar minimp3 implementation.
//
// The upstream minimp3 authors dedicated the software to the public domain
// worldwide, to the extent possible under law, and provide it without warranty
// (CC0 1.0 Universal): https://creativecommons.org/publicdomain/zero/1.0/
// This attribution is retained from minimp3.h.
package mp3

// scalePCM translates the active int16-output branch. Comparisons use the
// source's double constants; the addition and conversion remain float32/int.
func scalePCM(sample float32) int16 {
	if float64(sample) >= 32766.5 {
		return 32767
	}
	if float64(sample) <= -32767.5 {
		return -32768
	}
	s := int16(int32(float32(sample + float32(0.5))))
	if s < 0 {
		s--
	}
	return s
}

// synthPair translates mp3d_synth_pair. z exposes at least899 words;
// pcm addresses the channel base and exposes at least16*nch+1 samples.
func synthPair(pcm []int16, nch int, z []float32) {
	a := float32(float32(z[14*64]-z[0]) * float32(29))
	a = float32(a + float32(float32(z[1*64]+z[13*64])*float32(213)))
	a = float32(a + float32(float32(z[12*64]-z[2*64])*float32(459)))
	a = float32(a + float32(float32(z[3*64]+z[11*64])*float32(2037)))
	a = float32(a + float32(float32(z[10*64]-z[4*64])*float32(5153)))
	a = float32(a + float32(float32(z[5*64]+z[9*64])*float32(6574)))
	a = float32(a + float32(float32(z[8*64]-z[6*64])*float32(37489)))
	a = float32(a + float32(z[7*64]*float32(75038)))
	pcm[0] = scalePCM(a)

	z = z[2:]
	a = float32(z[14*64] * float32(104))
	a = float32(a + float32(z[12*64]*float32(1567)))
	a = float32(a + float32(z[10*64]*float32(9727)))
	a = float32(a + float32(z[8*64]*float32(64019)))
	a = float32(a + float32(z[6*64]*float32(-9975)))
	a = float32(a + float32(z[4*64]*float32(-45)))
	a = float32(a + float32(z[2*64]*float32(146)))
	a = float32(a + float32(z[0]*float32(-5)))
	pcm[16*nch] = scalePCM(a)
}

var synthWindow = [240]float32{
	-1, 26, -31, 208, 218, 401, -519, 2063, 2000, 4788, -5517, 7134, 5959, 35640, -39336, 74992,
	-1, 24, -35, 202, 222, 347, -581, 2080, 1952, 4425, -5879, 7640, 5288, 33791, -41176, 74856,
	-1, 21, -38, 196, 225, 294, -645, 2087, 1893, 4063, -6237, 8092, 4561, 31947, -43006, 74630,
	-1, 19, -41, 190, 227, 244, -711, 2085, 1822, 3705, -6589, 8492, 3776, 30112, -44821, 74313,
	-1, 17, -45, 183, 228, 197, -779, 2075, 1739, 3351, -6935, 8840, 2935, 28289, -46617, 73908,
	-1, 16, -49, 176, 228, 153, -848, 2057, 1644, 3004, -7271, 9139, 2037, 26482, -48390, 73415,
	-2, 14, -53, 169, 227, 111, -919, 2032, 1535, 2663, -7597, 9389, 1082, 24694, -50137, 72835,
	-2, 13, -58, 161, 224, 72, -991, 2001, 1414, 2330, -7910, 9592, 70, 22929, -51853, 72169,
	-2, 11, -63, 154, 221, 36, -1064, 1962, 1280, 2006, -8209, 9750, -998, 21189, -53534, 71420,
	-2, 10, -68, 147, 215, 2, -1137, 1919, 1131, 1692, -8491, 9863, -2122, 19478, -55178, 70590,
	-3, 9, -73, 139, 208, -29, -1210, 1870, 970, 1388, -8755, 9935, -3300, 17799, -56778, 69679,
	-3, 8, -79, 132, 200, -57, -1283, 1817, 794, 1095, -8998, 9966, -4533, 16155, -58333, 68692,
	-4, 7, -85, 125, 189, -83, -1356, 1759, 605, 814, -9219, 9959, -5818, 14548, -59838, 67629,
	-4, 7, -91, 117, 177, -106, -1428, 1698, 402, 545, -9416, 9916, -7154, 12980, -61289, 66494,
	-5, 6, -97, 111, 163, -127, -1498, 1634, 185, 288, -9585, 9838, -8540, 11455, -62684, 65290,
}

// synth translates the scalar mp3d_synth path. xl exposes at least560 plus
// 576*(nch-1) words, dstl has room for64*nch samples, and lins exposes1088
// history/work words. The granule caller advances these slices each iteration.
func synth(xl []float32, dstl []int16, nch int, lins []float32) {
	xrOff := 576 * (nch - 1)
	dstrOff := nch - 1
	zbase := 15 * 64
	lins[zbase+4*15] = xl[18*16]
	lins[zbase+4*15+1] = xl[xrOff+18*16]
	lins[zbase+4*15+2] = xl[0]
	lins[zbase+4*15+3] = xl[xrOff]
	lins[zbase+4*31] = xl[1+18*16]
	lins[zbase+4*31+1] = xl[xrOff+1+18*16]
	lins[zbase+4*31+2] = xl[1]
	lins[zbase+4*31+3] = xl[xrOff+1]

	// For mono, the first two calls write [0,16] and [32,48], which the
	// following pair calls overwrite. PCM output and lins workspace are disjoint.
	if nch != 1 {
		synthPair(dstl[dstrOff:], nch, lins[4*15+1:])
		synthPair(dstl[dstrOff+32*nch:], nch, lins[4*15+64+1:])
	}
	synthPair(dstl, nch, lins[4*15:])
	synthPair(dstl[32*nch:], nch, lins[4*15+64:])

	wpos := 0
	for i := 14; i >= 0; i-- {
		z := zbase
		lins[z+4*i] = xl[18*(31-i)]
		lins[z+4*i+1] = xl[xrOff+18*(31-i)]
		lins[z+4*i+2] = xl[1+18*(31-i)]
		lins[z+4*i+3] = xl[xrOff+1+18*(31-i)]
		lins[z+4*(i+16)] = xl[1+18*(1+i)]
		lins[z+4*(i+16)+1] = xl[xrOff+1+18*(1+i)]
		lins[z+4*(i-16)+2] = xl[18*(1+i)]
		lins[z+4*(i-16)+3] = xl[xrOff+18*(1+i)]

		var a, b [4]float32
		windowBase := wpos
		wpos += 16
		laneStep := 1
		if nch == 1 {
			// Mono stores even lanes over odd lanes at the same PCM indices.
			laneStep = 2
		}
		for j := 0; j < 4; j += laneStep {
			vz := z + 4*i
			vy := vz - 15*64
			w0 := synthWindow[windowBase]
			w1 := synthWindow[windowBase+1]
			p0 := float32(lins[vz+j] * w1)
			p1 := float32(lins[vy+j] * w0)
			q0 := float32(lins[vz+j] * w0)
			q1 := float32(lins[vy+j] * w1)
			bv := float32(p0 + p1)
			av := float32(q0 - q1)
			for step := 1; step < 8; step++ {
				windowPos := windowBase + 2*step
				w0 := synthWindow[windowPos]
				w1 := synthWindow[windowPos+1]
				k := step
				vz := z + 4*i - k*64
				vy := z + 4*i - (15-k)*64
				p0 := float32(lins[vz+j] * w1)
				p1 := float32(lins[vy+j] * w0)
				q0 := float32(lins[vz+j] * w0)
				q1 := float32(lins[vy+j] * w1)
				bv = float32(bv + float32(p0+p1))
				if step&1 != 0 {
					av = float32(av + float32(q1-q0))
				} else {
					av = float32(av + float32(q0-q1))
				}
			}
			a[j], b[j] = av, bv
		}

		pcm := func(off int, v float32) { dstl[off] = scalePCM(v) }
		if nch != 1 {
			pcm(dstrOff+(15-i)*nch, a[1])
			pcm(dstrOff+(17+i)*nch, b[1])
		}
		pcm((15-i)*nch, a[0])
		pcm((17+i)*nch, b[0])
		if nch != 1 {
			pcm(dstrOff+(47-i)*nch, a[3])
			pcm(dstrOff+(49+i)*nch, b[3])
		}
		pcm((47-i)*nch, a[2])
		pcm((49+i)*nch, b[2])
	}
}

// synthGranule translates mp3d_synth_granule for the production int16 output
// and the standard scalar mono/stereo state update.
func synthGranule(qmfState, grbuf []float32, nBands, nch int, pcm []int16, lins []float32) {
	for i := 0; i < nch; i++ {
		dctII(grbuf[576*i:], nBands)
	}
	copy(lins[:15*64], qmfState[:15*64])
	for i := 0; i < nBands; i += 2 {
		synth(grbuf[i:], pcm[32*nch*i:], nch, lins[i*64:])
	}
	if nch == 1 {
		for i := 0; i < 15*64; i += 2 {
			qmfState[i] = lins[nBands*64+i]
		}
	} else {
		copy(qmfState[:15*64], lins[nBands*64:nBands*64+15*64])
	}
}
