package tsslice

import "testing"

func BenchmarkAVSliceGet(b *testing.B) {
	s := NewAVSlice(10, 10)

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Get(1)
		}
	})
}

func BenchmarkRWSliceGet(b *testing.B) {
	s := NewRWSlice(10, 10)

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Get(1)
		}
	})
}

func BenchmarkAVSliceSet(b *testing.B) {
	s := NewAVSlice(10, 10)

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Set(1, 1)
		}
	})
}

func BenchmarkRWSliceSet(b *testing.B) {
	s := NewRWSlice(10, 10)

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Set(1, 1)
		}
	})
}
