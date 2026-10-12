package alloc

import (
	"fmt"
	"testing"
	"unsafe"
)

var memoryBenchmarkByte byte

func BenchmarkMemoryFill(b *testing.B) {
	for _, size := range []int{32, 256, 1024, 4096, 65536} {
		for _, v := range []byte{0, 0xac} {
			b.Run(fmt.Sprintf("bytes%d/value%d", size, v), func(b *testing.B) {
				data, free := Make([]byte(nil), size)
				defer free()
				ptr := unsafe.Pointer(&data[0])
				b.SetBytes(int64(size))
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					Memset(ptr, v, uintptr(size))
				}
				b.StopTimer()
				memoryBenchmarkByte = data[size-1]
			})
		}
	}
}

func BenchmarkMemoryCopy(b *testing.B) {
	for _, size := range []int{60, 256, 4096, 65536} {
		b.Run(fmt.Sprintf("bytes%d", size), func(b *testing.B) {
			src, freeSrc := Make([]byte(nil), size)
			defer freeSrc()
			dst, freeDst := Make([]byte(nil), size)
			defer freeDst()
			for i := range src {
				src[i] = byte(i * 37)
			}
			sp, dp := unsafe.Pointer(&src[0]), unsafe.Pointer(&dst[0])
			b.SetBytes(int64(size))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Memcpy(dp, sp, uintptr(size))
			}
			b.StopTimer()
			memoryBenchmarkByte = dst[size-1]
		})
	}
}
