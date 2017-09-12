package utils

// Stack base structure
type Stack struct {
	top  *node
	size int
}

type node struct {
	value interface{}
	next  *node
}

// Len of stack
func (s *Stack) Len() int {
	return s.size
}

// Push to stack
func (s *Stack) Push(value interface{}) {
	s.top = &node{value, s.top}
	s.size++
}

// Pop from stack
func (s *Stack) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}

	return nil
}

// Slice from stack
func (s *Stack) Slice() []interface{} {
	ret := make([]interface{}, 0)
	c := s.top
	for {
		if c == nil {
			return ret
		}

		ret = append(ret, c.value)
		c = c.next
	}
}

// ReverseSlice to revert slice
func (s *Stack) ReverseSlice() []interface{} {
	sl := s.Slice()
	for i, j := 0, len(sl)-1; i < j; i, j = i+1, j-1 {
		sl[i], sl[j] = sl[j], sl[i]
	}
	return sl
}
