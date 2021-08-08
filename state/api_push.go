package state

import . "golua/api"

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

func (s *luaState) PushGoFunction(f GoFunction)  {
	s.luaStack.push(newGoClosure(f))
}

func (s *luaState) PushGlobalTable()  {
	global := s.registry.get(LUA_RIDX_GLOBALS)
	s.luaStack.push(global)
}
