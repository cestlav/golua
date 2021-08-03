package state

type luaStack struct {
	slots []LuaValue
	top int
}

func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]LuaValue, size),
		top:   0,
	}
}

func (s *luaStack) check(n int)  {
	free := len(s.slots) - s.top
	for i := free; i < n; i++ {
		s.slots = append(s.slots, nil)
	}
}

func (s *luaStack) push(value LuaValue)  {
	if s.top == len(s.slots) {
		panic("stack overflow")
	}
	s.slots[s.top] = value
	s.top++
}

func (s *luaStack) pop() LuaValue {
	if s.top < 1 {
		panic("stack underflow")
	}
	s.top--
	v := s.slots[s.top]
	s.slots[s.top] = nil
	return v
}

func (s *luaStack) absIndex(index int) int {
	if index >= 0 {
		return index
	}
	return index + s.top +1
}

func (s *luaStack) isValid(index int) bool {
	absIndex := s.absIndex(index)
	return absIndex > 0 && absIndex <= s.top
}

func (s *luaStack) get(index int) LuaValue {
	absindex := s.absIndex(index)
	if absindex > 0 && absindex <= s.top {
		return s.slots[absindex - 1]
	}
	return nil
}

func (s *luaStack) set(index int, v LuaValue) {
	absIndex := s.absIndex(index)
	if absIndex > 0 && absIndex <= s.top {
		s.slots[absIndex - 1] = v
		return
	}
	panic("invalid index")
}

func (s *luaStack) reverse(from, to int) {
	slots := s.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}