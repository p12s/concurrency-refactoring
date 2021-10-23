package main

import (
	"testing"
)

func BenchmarkNew(b *testing.B) {

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			main()
		}
	})
}
