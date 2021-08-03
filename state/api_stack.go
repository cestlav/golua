package state

func (s *luaState) GetTop() int {
	return s.luaStack.top
}

func (s *luaState) AbsIndex(index int) int {
	return s.luaStack.absIndex(index)
}

func (s *luaState) CheckStack(n int) bool {
	s.luaStack.check(n)
	return true
}

func (s *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		s.luaStack.pop()
	}
}

func (s *luaState) Copy(from, to int) {
	v := s.luaStack.get(from)
	s.luaStack.set(to, v)
}

func (s *luaState) PushValue(index int) {
	v := s.luaStack.get(index)
	s.luaStack.push(v)
}

func (s *luaState) Replace(index int) {
	v := s.luaStack.pop()
	s.luaStack.set(index, v)
}

func (s *luaState) Insert(index int) {
	s.Rotate(index, 1)
}

func (s *luaState) Remove(index int) {
	s.Rotate(index, -1)
	s.Pop(1)
}

func (s *luaState) Rotate(index, n int)  {
	t := s.luaStack.top - 1
	p := s.luaStack.absIndex(index)

	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	s.luaStack.reverse(p, m)
	s.luaStack.reverse(m + 1, t)
	s.luaStack.reverse(p, t)
}

func (s *luaState) SetTop(index int)  {
	newTop := s.luaStack.absIndex(index)
	if newTop < 0 {
		panic("stack underflow")
	}

	n := s.luaStack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			s.luaStack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			s.luaStack.push(nil)
		}
	}
}