package stack


type node struct {
	v int
	next *node
}

type Stack struct {
	head *node
	n int
}

func New() *Stack { 
	return &Stack{} 
}

func (s *Stack) Push(x int) {
	s.head = &node{v: x, next: s.head}
	s.n++
}

func (s *Stack) Pop() (int, bool) {
	if s.head == nil {
		return 0, false
	}
	v := s.head.v
	next := s.head.next
	s.head.next = nil
	s.head = next
	s.n--
	return v, true
}

func (s *Stack) Peek() (int, bool) {
	if s.head == nil {
		return 0, false
	}
	return s.head.v, true
}

func (s *Stack) IsEmpty() bool {
	return s.n == 0 
}

func (s *Stack) Size() int {
	return s.n 
}

func (s *Stack) Clear() {
	for p := s.head; p != nil; {
		next := p.next
		p.next = nil
		p = next
	}
	s.head = nil
	s.n = 0
}
