package state

func (s *luaState) PC() int {
	return s.luaStack.pc
}

func (s * luaState) AddPC(n int) {
	s.luaStack.pc += n
}

func (s *luaState) Fetch() uint32 {
	i := s.luaStack.closure.proto.Code[s.luaStack.pc]
	s.luaStack.pc++
	return i
}

func (s *luaState) GetConst(index int)  {
	c := s.luaStack.closure.proto.Constants[index]
	s.luaStack.push(c)
}

func (s *luaState) GetRK(rk int)  {
	if rk > 0xFF {
		s.GetConst(rk & 0xFF)
	} else {
		s.PushValue(rk + 1)
	}
}

func (s *luaState) RegisterCount() int {
	return int(s.luaStack.closure.proto.MaxStackSize)
}

func (s *luaState) LoadVararg(n int)  {
	if n < 0 {
		n = len(s.luaStack.varargs)
	}
	s.luaStack.check(n)
	s.luaStack.pushN(s.luaStack.varargs, n)
}

func (s *luaState) LoadProto(index int)  {
	proto := s.luaStack.closure.proto.ProtoTypes[index]
	closure := newLuaClosure(proto)
	s.luaStack.push(closure)

}