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

	for i, uvInfo := range proto.UpValues {
		uvIndex := int(uvInfo.Index)
		if uvInfo.InStack == 1 {
			if s.luaStack.openuvs == nil {
				s.luaStack.openuvs = make(map[int]*upValue)
			}

			if openuv, found := s.luaStack.openuvs[uvIndex]; found {
				closure.upValues[i] = openuv
			} else {
				closure.upValues[i] = &upValue{&s.luaStack.slots[uvIndex]}
				s.luaStack.openuvs[uvIndex] = closure.upValues[i]
			}
		} else {
			closure.upValues[i] = s.luaStack.closure.upValues[uvIndex]
		}
	}
}

func (s *luaState) CloseUpValues(n int)  {
	for i, openuv := range s.luaStack.openuvs {
		if i >= n - 1 {
			v := *openuv.value
			openuv.value = &v
			delete(s.luaStack.openuvs, i)
		}
	}
}