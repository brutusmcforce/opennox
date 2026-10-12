package legacy

import (
	"C"
	"math"
)

// See nox_math.h.

//export nox_math_sin
func nox_math_sin(x C.double) C.double { return C.double(math.Sin(float64(x))) }

//export nox_math_cos
func nox_math_cos(x C.double) C.double { return C.double(math.Cos(float64(x))) }

//export nox_math_acos
func nox_math_acos(x C.double) C.double { return C.double(math.Acos(float64(x))) }

//export nox_math_atan2
func nox_math_atan2(y, x C.double) C.double { return C.double(math.Atan2(float64(y), float64(x))) }

//export nox_math_pow
func nox_math_pow(x, y C.double) C.double { return C.double(math.Pow(float64(x), float64(y))) }
