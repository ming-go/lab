package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type slice struct {
	sync.Mutex
	av  atomic.Value
	len int
}

func NewSlice(len int) *slice {
	si := &slice{}
	si.len = len
	si.av.Store(make([]int, len, 1000))
	return si
}

func (s *slice) Get(index int) int {
	return s.av.Load().([]int)[index]
}

func (s *slice) Set(index int, value int) {
	s.Lock()
	si := s.av.Load().([]int)
	si[index] = value
	s.av.Store(si)
	s.Unlock()
}

func (s *slice) Len() int {
	return s.len
}

func (s *slice) Append(elem int) {
	s.Lock()
	ins := s.av.Load().([]int)
	s.av.Store(append(ins, elem))
	s.Unlock()
}

func main() {
	s := NewSlice(1000)

	//go func() {
	//	for i := 0 {
	//		s.Set(3, 1)
	//	}
	//}()

	go func() {
		for i := 0; i < 5000; i++ {
			go func(i int) {
				x := s.Get(0)
				s.Set(0, x+1)
			}(i)
		}
	}()

	<-time.After(10 * time.Second)

	fmt.Println(s.Get(0))

	/*
		for i := 0; i < 5000; i++ {
			fmt.Println(s.Get(i))
		}
	*/

	//fmt.Println(s.Len())

	<-time.After(10 * time.Minute)
}
