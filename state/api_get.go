package state

import . "golua/api"

func (s *luaState) NewTable() {
	s.CreateTable(0, 0)
}

func (s *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	s.luaStack.push(t)
}

func (s *luaState) GetTable(idx int) LuaType {
	t := s.luaStack.get(idx)
	k := s.luaStack.pop()
	return s.getTable(t, k)
}

func (s *luaState) GetField(idx int, k string) LuaType {
	t := s.luaStack.get(idx)
	return s.getTable(t, k)
}

func (s *luaState) GetI(idx int, i int64) LuaType {
	t := s.luaStack.get(idx)
	return s.getTable(t, i)
}

func (s *luaState) getTable(t, k luaValue) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		s.luaStack.push(v)
		return typeOf(v)
	}

	panic("not a table!") // todo
}

func (s *luaState) GetGlobal(name string) LuaType {
	t := s.registry.get(LUA_RIDX_GLOBALS)
	return s.getTable(t, name)
}