package state

func (s *luaState) Len(index int)  {
	val := s.luaStack.get(index)

	if str, ok := val.(string); ok {
		s.luaStack.push(int64(len(str)))
	} else {
		panic("length error!")
	}
}

func (s *luaState) Concat(n int)  {
	if n == 0 {
		s.luaStack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if s.IsString(-1) && s.IsString(-2) {
				s2 := s.ToString(-1)
				s1 := s.ToString(-2)
				s.luaStack.pop()
				s.luaStack.pop()
				s.luaStack.push(s1 + s2)
				continue
			}

			panic("concatenation error!")
		}
	}
}