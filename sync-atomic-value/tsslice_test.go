package tsslice

import (
	"sync"
	"testing"
)

func BenchmarkAVSliceGetSet(b *testing.B) {
	s := NewAVSlice(10, 10)

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				s.Get(1)
				wg.Done()
			}()

			go func() {
				s.Set(1, 1)
				wg.Done()
			}()
			wg.Wait()
		}
	})
}

func BenchmarkRWSliceGetSet(b *testing.B) {
	s := NewRWSlice(10, 10)

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				s.Get(1)
				wg.Done()
			}()

			go func() {
				s.Set(1, 1)
				wg.Done()
			}()
			wg.Wait()
		}
	})
}

func BenchmarkCHSliceGetSet(b *testing.B) {
	s := NewCHSlice(10, 10)

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				s.Get(1)
				wg.Done()
			}()

			go func() {
				s.Set(1, 1)
				wg.Done()
			}()
			wg.Wait()
		}
	})
}

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

func BenchmarkCHSliceGet(b *testing.B) {
	s := NewCHSlice(10, 10)

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

func BenchmarkCHSliceSet(b *testing.B) {
	s := NewCHSlice(10, 10)

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Set(1, 1)
		}
	})
}
