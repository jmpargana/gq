package main

// TODO: refactor to use generics
type stack struct {
	el []*cmd
}

func newStack() *stack {
	return &stack{
		el: make([]*cmd, 0),
	}
}

func (s *stack) peek() *cmd {
	if len(s.el) == 0 {
		return nil
	}
	return s.el[len(s.el)-1]
}

func (s *stack) push(el *cmd) {
	s.el = append(s.el, el)
}

func (s *stack) pop() *cmd {
	if len(s.el) == 0 {
		return nil
	}
	n := len(s.el) - 1
	el := s.el[n]
	s.el[n] = nil
	s.el = s.el[:n]
	return el
}

func (s *stack) isEmpty() bool {
	return len(s.el) == 0
}
