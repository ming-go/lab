package tsslice

type chSlice struct {
	s      []int
	ch     chan func()
	stopCh chan struct{}
}

func NewCHSlice(len int, cap int) *chSlice {
	cs := &chSlice{
		s:      make([]int, len, cap),
		ch:     make(chan func()),
		stopCh: make(chan struct{}),
	}

	go func() {
		for {
			select {
			case f := <-cs.ch:
				f()
			case <-cs.stopCh:
				break
			}
		}
	}()

	return cs
}

func (s *chSlice) Get(index int) (value int) {
	finCh := make(chan struct{})

	s.ch <- func() {
		value = s.s[index]
		finCh <- struct{}{}
	}

	<-finCh
	return
}

func (s *chSlice) Set(index int, value int) {
	finCh := make(chan struct{})

	s.ch <- func() {
		s.s[index] = value
		finCh <- struct{}{}
	}

	<-finCh
	return
}

func (s *chSlice) Len() (length int) {
	finCh := make(chan struct{})

	s.ch <- func() {
		length = len(s.s)
		finCh <- struct{}{}
	}

	<-finCh
	return
}

func (s *chSlice) Append(elements ...int) {
	finCh := make(chan struct{})

	s.ch <- func() {
		s.s = append(s.s, elements...)
		finCh <- struct{}{}
	}

	<-finCh
	return
}
