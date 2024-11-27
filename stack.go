package xml2go

type stack[T any] struct {
	elements []T
}

func (s *stack[T]) push(elem T) {
	s.elements = append(s.elements, elem)
}

func (s *stack[T]) pop() T {
	elem := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return elem
}

func (s *stack[T]) top() T {
	return s.elements[len(s.elements)-1]
}

func (s *stack[T]) empty() bool {
	return len(s.elements) == 0
}
