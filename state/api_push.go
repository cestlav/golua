package state

func (s *luaState) PushNil()  {
	s.luaStack.push(nil)
}

func (s *luaState) PushBoolean(b bool) {
	s.luaStack.push(b)
}

func (s *luaState) PushInteger(i int64)  {
	s.luaStack.push(i)
}

func (s *luaState) PushNumber(f float64)  {
	s.luaStack.push(f)
}

func (s *luaState) PushString(str string) {
	s.luaStack.push(str)
}