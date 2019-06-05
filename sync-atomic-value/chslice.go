package tsslice

type chSlice struct {
	s      []int
	ch     chan func()
	stopCh chan struct{}
}

func NewCHSlice(len int, cap int) *chSlice {
	cs := &chSlice{
		s: make([]int, len, cap),
	}

	go func() {
		for {
			select {
			case (<-cs.ch)():
			case stppCh:
			}
		}
	}()
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
	}

	<-finCh
	return
}

func (s *chSlice) Len() (len int) {
	finCh := make(chan struct{})
	s.ch <- func() {
		len = len(s.s[index])
	}
	<-finCh
	return
}

func (s *chSlice) Append(elements ...int) {
	finCh := make(chan struct{})
	s.ch <- func() {
		s.s = append(s.s, elements...)
	}
	<-finCh
	return
}
