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
	return s.getTable(t, k, false)
}

func (s *luaState) GetField(idx int, k string) LuaType {
	t := s.luaStack.get(idx)
	return s.getTable(t, k, false)
}

func (s *luaState) GetI(idx int, i int64) LuaType {
	t := s.luaStack.get(idx)
	return s.getTable(t, i, false)
}

func (s *luaState) getTable(t, k luaValue, raw bool) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		if raw || v != nil || !tbl.hasMetaField("__index") {
			s.luaStack.push(v)
			return typeOf(v)
		}
	}

	if !raw {
		if mf := getMetaField(t, "__index", s); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return s.getTable(x, k, false)
			case *closure:
				s.luaStack.push(mf)
				s.luaStack.push(t)
				s.luaStack.push(k)
				s.Call(2, 1)
				v := s.luaStack.get(-1)
				return typeOf(v)
			}
		}
	}

	panic("index error!")
}

func (s *luaState) GetGlobal(name string) LuaType {
	t := s.registry.get(LUA_RIDX_GLOBALS)
	return s.getTable(t, name, false)
}

func (s *luaState) GetMetaTable(index int) bool {
	val := s.luaStack.get(index)

	if mt := getMetatable(val, s); mt != nil {
		s.luaStack.push(mt)
		return true
	} else {
		return false
	}
}

func (s *luaState) RawGet(idx int) LuaType {
	t := s.luaStack.get(idx)
	k := s.luaStack.pop()
	return s.getTable(t, k, true)
}

func (s *luaState) RawGetI(idx int, i int64) LuaType {
	t := s.luaStack.get(idx)
	return s.getTable(t, i, true)
}