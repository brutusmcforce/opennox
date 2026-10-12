// Inverse transforms translated from the scalar minimp3 implementation.
//
// The upstream minimp3 authors dedicated the software to the public domain
// worldwide, to the extent possible under law, and provide it without warranty
// (CC0 1.0 Universal): https://creativecommons.org/publicdomain/zero/1.0/
// This attribution is retained from minimp3.h.
package mp3

var mdctTwiddle9 = [18]float32{
	0.73727734, 0.79335334, 0.84339145, 0.88701083, 0.92387953, 0.95371695,
	0.97629601, 0.99144486, 0.99904822, 0.67559021, 0.60876143, 0.53729961,
	0.46174861, 0.38268343, 0.30070580, 0.21643961, 0.13052619, 0.04361938,
}

var mdctTwiddle3 = [6]float32{0.79335334, 0.92387953, 0.99144486, 0.60876143, 0.38268343, 0.13052619}

var mdctWindow = [2][18]float32{
	{
		0.99904822, 0.99144486, 0.97629601, 0.95371695, 0.92387953, 0.88701083,
		0.84339145, 0.79335334, 0.73727734, 0.04361938, 0.13052619, 0.21643961,
		0.30070580, 0.38268343, 0.46174861, 0.53729961, 0.60876143, 0.67559021,
	},
	{
		1, 1, 1, 1, 1, 1, 0.99144486, 0.92387953, 0.79335334,
		0, 0, 0, 0, 0, 0, 0.13052619, 0.38268343, 0.60876143,
	},
}

// dct3Nine translates the scalar L3_dct3_9 butterfly in source order.
func dct3Nine(y *[9]float32) {
	s0, s2, s4, s6, s8 := y[0], y[2], y[4], y[6], y[8]
	t0 := float32(s6 * float32(0.5))
	t0 = float32(s0 + t0)
	s0 = float32(s0 - s6)
	t4 := float32(s4 + s2)
	t4 = float32(t4 * float32(0.93969262))
	t2 := float32(s8 + s2)
	t2 = float32(t2 * float32(0.76604444))
	s6 = float32(s4 - s8)
	s6 = float32(s6 * float32(0.17364818))
	s4 = float32(s4 + float32(s8-s2))

	s2 = float32(s0 - float32(s4*float32(0.5)))
	y[4] = float32(s4 + s0)
	s8 = float32(float32(t0-t2) + s6)
	s0 = float32(float32(t0-t4) + t2)
	s4 = float32(float32(t0+t4) - s6)

	s1, s3, s5, s7 := y[1], y[3], y[5], y[7]
	s3 = float32(s3 * float32(0.86602540))
	t0 = float32(s5 + s1)
	t0 = float32(t0 * float32(0.98480775))
	t4 = float32(s5 - s7)
	t4 = float32(t4 * float32(0.34202014))
	t2 = float32(s1 + s7)
	t2 = float32(t2 * float32(0.64278761))
	s1 = float32(float32(float32(s1-s5)-s7) * float32(0.86602540))

	s5 = float32(float32(t0-s3) - t2)
	s7 = float32(float32(t4-s3) - t0)
	s3 = float32(float32(t4+s3) - t2)
	y[0] = float32(s4 - s7)
	y[1] = float32(s2 + s1)
	y[2] = float32(s0 - s3)
	y[3] = float32(s8 + s5)
	y[5] = float32(s8 - s5)
	y[6] = float32(s0 + s3)
	y[7] = float32(s2 - s1)
	y[8] = float32(s4 + s7)
}

// imdct36 translates the scalar L3_imdct36 path. grbuf and overlap must have
// 18*nbands and 9*nbands elements respectively; window has 18 elements.
func imdct36(grbuf, overlap []float32, window []float32, nBands int) {
	for band := 0; band < nBands; band++ {
		g := grbuf[18*band : 18*(band+1)]
		o := overlap[9*band : 9*(band+1)]
		var co, si [9]float32
		co[0] = -g[0]
		si[0] = g[17]
		for i := 0; i < 4; i++ {
			si[8-2*i] = float32(g[4*i+1] - g[4*i+2])
			co[1+2*i] = float32(g[4*i+1] + g[4*i+2])
			si[7-2*i] = float32(g[4*i+4] - g[4*i+3])
			co[2+2*i] = -float32(g[4*i+3] + g[4*i+4])
		}
		dct3Nine(&co)
		dct3Nine(&si)
		si[1] = -si[1]
		si[3] = -si[3]
		si[5] = -si[5]
		si[7] = -si[7]

		for i := 0; i < 9; i++ {
			ovl := o[i]
			sum0 := float32(co[i] * mdctTwiddle9[9+i])
			sum1 := float32(si[i] * mdctTwiddle9[i])
			sum := float32(sum0 + sum1)
			o0 := float32(co[i] * mdctTwiddle9[i])
			o1 := float32(si[i] * mdctTwiddle9[9+i])
			o[i] = float32(o0 - o1)
			g[i] = float32(float32(ovl*window[i]) - float32(sum*window[9+i]))
			g[17-i] = float32(float32(ovl*window[9+i]) + float32(sum*window[i]))
		}
	}
}

// idct3 translates the scalar L3_idct3 operation order into three outputs.
func idct3(x0, x1, x2 float32, dst *[3]float32) {
	m1 := float32(x1 * float32(0.86602540))
	a1 := float32(x0 - float32(x2*float32(0.5)))
	dst[1] = float32(x0 + x2)
	dst[0] = float32(a1 + m1)
	dst[2] = float32(a1 - m1)
}

// imdct12 translates L3_imdct12. x provides indices 0..15; dst and overlap
// provide six and three writable elements. The destination spans may be
// adjacent but dst[0:6] and overlap[0:3] must be disjoint as in the C caller.
func imdct12(x, dst, overlap []float32) {
	var co, si [3]float32
	x0 := -x[0]
	x1 := float32(x[6] + x[3])
	x2 := float32(x[12] + x[9])
	idct3(x0, x1, x2, &co)
	x0 = x[15]
	x1 = float32(x[12] - x[9])
	x2 = float32(x[6] - x[3])
	idct3(x0, x1, x2, &si)
	si[1] = -si[1]
	for i := 0; i < 3; i++ {
		ovl := overlap[i]
		sum0 := float32(co[i] * mdctTwiddle3[3+i])
		sum1 := float32(si[i] * mdctTwiddle3[i])
		sum := float32(sum0 + sum1)
		o0 := float32(co[i] * mdctTwiddle3[i])
		o1 := float32(si[i] * mdctTwiddle3[3+i])
		overlap[i] = float32(o0 - o1)
		dst[i] = float32(float32(ovl*mdctTwiddle3[2-i]) - float32(sum*mdctTwiddle3[5-i]))
		dst[5-i] = float32(float32(ovl*mdctTwiddle3[5-i]) + float32(sum*mdctTwiddle3[2-i]))
	}
}

// imdctShort translates L3_imdct_short and preserves its three overlapping
// transform windows by copying the original 18 values before each update.
func imdctShort(grbuf, overlap []float32, nBands int) {
	for band := 0; band < nBands; band++ {
		g := grbuf[18*band : 18*(band+1)]
		o := overlap[9*band : 9*(band+1)]
		var tmp [18]float32
		copy(tmp[:], g[:18])
		copy(g[:6], o[:6])
		imdct12(tmp[:], g[6:], o[6:])
		imdct12(tmp[1:], g[12:], o[6:])
		imdct12(tmp[2:], o[:], o[6:])
	}
}

// changeSign translates the observable writes of L3_change_sign: odd indices
// in odd-numbered 18-float subbands only, ending at sample 575.
func changeSign(grbuf []float32) {
	for band := 1; band < 32; band += 2 {
		base := band * 18
		for i := 1; i < 18; i += 2 {
			grbuf[base+i] = -grbuf[base+i]
		}
	}
}

// imdctGranule translates L3_imdct_gr for nLongBands in 0..32. The caller
// provides the full 576-sample granule and 288 overlap values.
func imdctGranule(grbuf, overlap []float32, blockType, nLongBands int) {
	if nLongBands != 0 {
		imdct36(grbuf, overlap, mdctWindow[0][:], nLongBands)
		grbuf = grbuf[18*nLongBands:]
		overlap = overlap[9*nLongBands:]
	}
	if blockType == int(shortBlockType) {
		imdctShort(grbuf, overlap, 32-nLongBands)
	} else {
		window := 0
		if blockType == 3 {
			window = 1
		}
		imdct36(grbuf, overlap, mdctWindow[window][:], 32-nLongBands)
	}
}
