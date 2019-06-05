package tsslice

import "sync"

type rwSlice struct {
	sync.RWMutex
	s []int
}

func NewRWSlice(len int, cap int) *rwSlice {
	si := &rwSlice{}
	si.s = make([]int, len, cap)
	return si
}

func (s *rwSlice) Get(index int) int {
	s.RLock()
	defer s.RUnlock()
	return s.s[index]
}

func (s *rwSlice) Set(index int, value int) {
	s.Lock()
	s.s[index] = value
	s.Unlock()
}

func (s *rwSlice) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.s)
}

func (s *rwSlice) Append(elements ...int) {
	s.Lock()
	s.s = append(s.s, elements...)
	s.Unlock()
}
