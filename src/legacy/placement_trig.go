// Copyright (C) 2001-2026 Free Software Foundation, Inc.
// Adapted to Go from the IBM Accurate Mathematical Library in glibc.
// SPDX-License-Identifier: LGPL-2.1-or-later
// See docs/licenses/GLIBC-MATH-LGPL-2.1.txt.
package legacy

import "math"

const (
	trigS1    float64 = -0x1.5555555555555p-3
	trigS2    float64 = 0x1.1111111110ECEp-7
	trigS3    float64 = -0x1.A01A019DB08B8p-13
	trigS4    float64 = 0x1.71DE27B9A7ED9p-19
	trigS5    float64 = -0x1.ADDFFC2FCDF59p-26
	trigBig   float64 = 0x1.8000000000000p45
	trigHP0   float64 = 0x1.921FB54442D18p0
	trigHP1   float64 = 0x1.1A62633145C07p-54
	trigMP1   float64 = 0x1.921FB58000000p0
	trigMP2   float64 = -0x1.DDE973C000000p-27
	trigPP3   float64 = -0x1.CB3B398000000p-55
	trigPP4   float64 = -0x1.d747f23e32ed7p-83
	trigHPInv float64 = 0x1.45F306DC9C883p-1
	trigToInt float64 = 0x1.8000000000000p52
	trigSN3   float64 = -1.66666666666664880952546298448555e-01
	trigSN5   float64 = 8.33333214285722277379541354343671e-03
	trigCS2   float64 = 4.99999999999999999999950396842453e-01
	trigCS4   float64 = -4.16666666666664434524222570944589e-02
	trigCS6   float64 = 1.38888874007937613028114285595617e-03
)

func trigLookup(u float64) (sn, ssn, cs, ccs float64) {
	i := int(uint32(math.Float64bits(u))) << 2
	return trigTable[i], trigTable[i+1], trigTable[i+2], trigTable[i+3]
}
func trigDoCos(x, dx float64) float64 {
	if x < 0 {
		dx = -dx
	}
	u := trigBig + math.Abs(x)
	x = math.Abs(x) - (u - trigBig) + dx
	xx := x * x
	s := x + x*xx*(trigSN3+xx*trigSN5)
	c := xx * (trigCS2 + xx*(trigCS4+xx*trigCS6))
	sn, ssn, cs, ccs := trigLookup(u)
	cor := (ccs - s*ssn - cs*c) - sn*s
	return cs + cor
}
func trigDoSin(x, dx float64) float64 {
	old := x
	if math.Abs(x) < 0.126 {
		xx := x * x
		polynomial := ((((trigS5*xx+trigS4)*xx+trigS3)*xx + trigS2) * xx) + trigS1
		t := (polynomial*x-0.5*dx)*xx + dx
		return x + t
	}
	if x <= 0 {
		dx = -dx
	}
	u := trigBig + math.Abs(x)
	x = math.Abs(x) - (u - trigBig)
	xx := x * x
	s := x + (dx + x*xx*(trigSN3+xx*trigSN5))
	c := x*dx + xx*(trigCS2+xx*(trigCS4+xx*trigCS6))
	sn, ssn, cs, ccs := trigLookup(u)
	cor := (ssn + s*ccs - sn*c) + cs*s
	return math.Copysign(sn+cor, old)
}
func trigReduce(x float64) (a, da float64, n uint32) {
	t := x*trigHPInv + trigToInt
	xn := t - trigToInt
	y := (x - xn*trigMP1) - xn*trigMP2
	n = uint32(math.Float64bits(t)) & 3
	t1 := xn * trigPP3
	t2 := y - t1
	db := (y - t2) - t1
	t1 = xn * trigPP4
	b := t2 - t1
	db += (t2 - b) - t1
	return b, db, n
}
func trigDoSinCos(a, da float64, n uint32) float64 {
	var v float64
	if n&1 != 0 {
		v = trigDoCos(a, da)
	} else {
		v = trigDoSin(a, da)
	}
	if n&2 != 0 {
		return -v
	}
	return v
}

// placementSin and placementCos preserve the qualified libc rounding of the
// two placement owners. Their RNG starts in [-pi, pi] and advances at most 64
// times by 1.8849558, so every argument has magnitude below 128. Keep each
// binary64 evaluation boundary; this is not a general libm replacement.
func placementSin(x float64) float64 {
	k := uint32(math.Float64bits(x)>>32) & 0x7fffffff
	switch {
	case k < 0x3e500000:
		return x
	case k < 0x3feb6000:
		return trigDoSin(x, 0)
	case k < 0x400368fd:
		return math.Copysign(trigDoCos(trigHP0-math.Abs(x), trigHP1), x)
	case k < 0x419921fb:
		a, da, n := trigReduce(x)
		return trigDoSinCos(a, da, n)
	default:
		return math.Sin(x) // Outside the engine placement angle range.
	}
}
func placementCos(x float64) float64 {
	k := uint32(math.Float64bits(x)>>32) & 0x7fffffff
	switch {
	case k < 0x3e400000:
		return 1
	case k < 0x3feb6000:
		return trigDoCos(x, 0)
	case k < 0x400368fd:
		y := trigHP0 - math.Abs(x)
		a := y + trigHP1
		da := (y - a) + trigHP1
		return trigDoSin(a, da)
	case k < 0x419921fb:
		a, da, n := trigReduce(x)
		return trigDoSinCos(a, da, n+1)
	default:
		return math.Cos(x) // Outside the engine placement angle range.
	}
}
