package tsslice

import (
	"sync"
	"sync/atomic"
)

type avSlice struct {
	sync.Mutex
	av  atomic.Value
	len int
}

func NewAVSlice(len int, cap int) *avSlice {
	si := &avSlice{}
	si.len = len
	si.av.Store(make([]int, len, cap))
	return si
}

func (s *avSlice) Get(index int) int {
	return s.av.Load().([]int)[index]
}

func (s *avSlice) Set(index int, value int) {
	s.Lock()
	so := s.av.Load().([]int)
	sn := append(so[:0:0], so...)
	sn[index] = value
	s.av.Store(sn)
	s.Unlock()
}

func (s *avSlice) Len() int {
	// TODO: thread-safe
	return s.len
}

func (s *avSlice) Append(elements ...int) {
	s.Lock()
	so := s.av.Load().([]int)
	sn := append(so[:0:0], so...)
	s.av.Store(append(sn, elements...))
	s.Unlock()
}
