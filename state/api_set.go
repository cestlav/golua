package state

import . "golua/api"

func (s *luaState) SetTable(idx int) {
	t := s.luaStack.get(idx)
	v := s.luaStack.pop()
	k := s.luaStack.pop()
	s.setTable(t, k, v)
}

func (s *luaState) SetField(idx int, k string) {
	t := s.luaStack.get(idx)
	v := s.luaStack.pop()
	s.setTable(t, k, v)
}

func (s *luaState) SetI(idx int, i int64) {
	t := s.luaStack.get(idx)
	v := s.luaStack.pop()
	s.setTable(t, i, v)
}

func (s *luaState) setTable(t, k, v luaValue) {
	if tbl, ok := t.(*luaTable); ok {
		tbl.put(k, v)
		return
	}

	panic("not a table!")
}

func (s *luaState) SetGlobal(name string) {
	t := s.registry.get(LUA_RIDX_GLOBALS)
	v := s.luaStack.pop()
	s.setTable(t, name, v)
}

func (s *luaState) Register(name string, f GoFunction) {
	s.PushGoFunction(f)
	s.SetGlobal(name)
}