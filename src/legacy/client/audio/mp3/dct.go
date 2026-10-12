// Synthesis DCT translated from the scalar minimp3 implementation.
//
// The upstream minimp3 authors dedicated the software to the public domain
// worldwide, to the extent possible under law, and provide it without warranty
// (CC0 1.0 Universal): https://creativecommons.org/publicdomain/zero/1.0/
// This attribution is retained from minimp3.h.
package mp3

var dct2Section = [24]float32{
	10.19000816, 0.50060302, 0.50241929, 3.40760851, 0.50547093, 0.52249861,
	2.05778098, 0.51544732, 0.56694406, 1.48416460, 0.53104258, 0.64682180,
	1.16943991, 0.55310392, 0.78815460, 0.97256821, 0.58293498, 1.06067765,
	0.83934963, 0.62250412, 1.72244716, 0.74453628, 0.67480832, 5.10114861,
}

// dctII translates the scalar mp3d_DCT_II path. The active MINIMP3_ONLY_MP3
// synthesis caller supplies n=18; n=12 is the supplemental Layer I/II caller
// when that code is enabled. Each transform column is staged fully before its
// 32 strided outputs are written.
func dctII(grbuf []float32, n int) {
	for k := 0; k < n; k++ {
		y := grbuf[k:]
		var t [4][8]float32
		for i := 0; i < 8; i++ {
			x0 := y[i*18]
			x1 := y[(15-i)*18]
			x2 := y[(16+i)*18]
			x3 := y[(31-i)*18]
			t0 := float32(x0 + x3)
			t1 := float32(x1 + x2)
			t2 := float32(float32(x1-x2) * dct2Section[3*i])
			t3 := float32(float32(x0-x3) * dct2Section[3*i+1])
			t[0][i] = float32(t0 + t1)
			t[1][i] = float32(float32(t0-t1) * dct2Section[3*i+2])
			t[2][i] = float32(t3 + t2)
			t[3][i] = float32(float32(t3-t2) * dct2Section[3*i+2])
		}

		for group := 0; group < 4; group++ {
			x := &t[group]
			x0, x1, x2, x3 := x[0], x[1], x[2], x[3]
			x4, x5, x6, x7 := x[4], x[5], x[6], x[7]
			xt := float32(x0 - x7)
			x0 = float32(x0 + x7)
			x7 = float32(x1 - x6)
			x1 = float32(x1 + x6)
			x6 = float32(x2 - x5)
			x2 = float32(x2 + x5)
			x5 = float32(x3 - x4)
			x3 = float32(x3 + x4)
			x4 = float32(x0 - x3)
			x0 = float32(x0 + x3)
			x3 = float32(x1 - x2)
			x1 = float32(x1 + x2)
			x[0] = float32(x0 + x1)
			x[4] = float32(float32(x0-x1) * float32(0.70710677))
			x5 = float32(x5 + x6)
			x6 = float32(float32(x6+x7) * float32(0.70710677))
			x7 = float32(x7 + xt)
			x3 = float32(float32(x3+x4) * float32(0.70710677))
			x5 = float32(x5 - float32(x7*float32(0.198912367)))
			x7 = float32(x7 + float32(x5*float32(0.382683432)))
			x5 = float32(x5 - float32(x7*float32(0.198912367)))
			x0 = float32(xt - x6)
			xt = float32(xt + x6)
			x[1] = float32(float32(xt+x7) * float32(0.50979561))
			x[2] = float32(float32(x4+x3) * float32(0.54119611))
			x[3] = float32(float32(x0-x5) * float32(0.60134488))
			x[5] = float32(float32(x0+x5) * float32(0.89997619))
			x[6] = float32(float32(x4-x3) * float32(1.30656302))
			x[7] = float32(float32(xt-x7) * float32(2.56291556))
		}

		for i := 0; i < 7; i++ {
			y[0] = t[0][i]
			y[18] = float32(float32(t[2][i]+t[3][i]) + t[3][i+1])
			y[36] = float32(t[1][i] + t[1][i+1])
			y[54] = float32(float32(t[2][i+1]+t[3][i]) + t[3][i+1])
			y = y[4*18:]
		}
		y[0] = t[0][7]
		y[18] = float32(t[2][7] + t[3][7])
		y[36] = t[1][7]
		y[54] = t[3][7]
	}
}
