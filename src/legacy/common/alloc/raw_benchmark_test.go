package alloc

import (
	"testing"
	"unsafe"
)

func BenchmarkRawAllocation(b *testing.B) {
	for _, batch := range []bool{false, true} {
		name := "serial"
		if batch {
			name = "batch256"
		}
		b.Run(name, func(b *testing.B) {
			sizes := [...]uintptr{16, 96, 256, 2048, 65536}
			var live [256]unsafe.Pointer
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				slot := n % len(live)
				if live[slot] != nil {
					RawFree(live[slot])
					live[slot] = nil
				}
				p := RawCalloc(1, sizes[n%len(sizes)])
				if p == nil {
					b.Fatal("allocation failed")
				}
				*(*byte)(p) = byte(n)
				if batch {
					live[slot] = p
				} else {
					RawFree(p)
				}
			}
			b.StopTimer()
			for _, p := range live {
				RawFree(p)
			}
		})
	}
}

func BenchmarkRawAllocationParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		sizes := [...]uintptr{16, 96, 256, 2048, 65536}
		n := 0
		for pb.Next() {
			p := RawCalloc(1, sizes[n%len(sizes)])
			if p == nil {
				panic("bounded benchmark allocation failed")
			}
			*(*byte)(p) = byte(n)
			RawFree(p)
			n++
		}
	})
}
