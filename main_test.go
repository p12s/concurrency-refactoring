package main

import (
	"testing"
)

func BenchmarkMain(b *testing.B) {

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			main()
		}
	})
}
