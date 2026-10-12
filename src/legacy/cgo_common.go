package legacy

/*
#cgo CFLAGS: -fshort-wchar -fno-strict-aliasing -fno-strict-overflow
#cgo linux CFLAGS: -fstack-protector-strong
// Use SSE2 instead of the x87 FPU for floating point: x87 keeps intermediate results in 80-bit registers and only
// rounds them when they are spilled to memory, which depends on register allocation. The Linux (PIE) and Windows
// builds allocate registers differently, so the game logic gave different results on each. Go requires SSE2 on 386 too.
#cgo 386 CFLAGS: -msse2 -mfpmath=sse

#cgo CFLAGS: -Werror=return-type
#cgo CFLAGS: -Werror=implicit-function-declaration
#cgo CFLAGS: -Werror=pointer-arith
#cgo CFLAGS: -Werror=implicit-int
#cgo CFLAGS: -Werror=unused-label
#cgo CFLAGS: -Werror=address
//#cgo CFLAGS: -Werror=unused-variable
//#cgo CFLAGS: -Werror=cast-function-type
#cgo CFLAGS: -Werror=stringop-overflow
//#cgo linux CFLAGS: -Werror=attribute-warning
#cgo linux CFLAGS: -Werror=uninitialized
//#cgo linux CFLAGS: -Werror=parentheses // FIXME: resolve issues related to this flag
//#cgo linux CFLAGS: -Werror=switch-outside-range

#cgo CFLAGS: -Wno-pointer-to-int-cast
#cgo CFLAGS: -Wno-int-to-pointer-cast
#cgo CFLAGS: -Wno-incompatible-pointer-types
#cgo CFLAGS: -Wno-int-conversion
#cgo CFLAGS: -Wno-format
#cgo CFLAGS: -Wno-shift-count-overflow
#cgo CFLAGS: -Wno-pedantic
#cgo CFLAGS: -Wno-bad-function-cast
#cgo CFLAGS: -Wno-strict-prototypes
#cgo CFLAGS: -Wno-discarded-qualifiers
#cgo CFLAGS: -Wno-return-local-addr
#cgo CFLAGS: -Wno-unused-result

#cgo windows LDFLAGS: -lws2_32
#cgo linux LDFLAGS: -lm
#cgo CFLAGS: -DNOX_CGO
*/
import "C"
