package state

import . "golua/api"

func (s *luaState) SetTable(idx int) {
	t := s.luaStack.get(idx)
	v := s.luaStack.pop()
	k := s.luaStack.pop()
	s.setTable(t, k, v, false)
}

func (s *luaState) SetField(idx int, k string) {
	t := s.luaStack.get(idx)
	v := s.luaStack.pop()
	s.setTable(t, k, v, false)
}

func (s *luaState) SetI(idx int, i int64) {
	t := s.luaStack.get(idx)
	v := s.luaStack.pop()
	s.setTable(t, i, v, false)
}

func (s *luaState) setTable(t, k, v luaValue, raw bool) {
	if tbl, ok := t.(*luaTable); ok {
		if raw || tbl.get(k) != nil || !tbl.hasMetaField("__newindex") {
			tbl.put(k, v)
			return
		}
	}

	if !raw {
		if mf := getMetaField(t, "__newindex", s); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				s.setTable(x, k, v, false)
				return
			case *closure:
				s.luaStack.push(mf)
				s.luaStack.push(t)
				s.luaStack.push(k)
				s.luaStack.push(v)
				s.Call(3, 0)
				return
			}
		}
	}

	panic("index error!")
}

func (s *luaState) SetGlobal(name string) {
	t := s.registry.get(LUA_RIDX_GLOBALS)
	v := s.luaStack.pop()
	s.setTable(t, name, v, false)
}

func (s *luaState) Register(name string, f GoFunction) {
	s.PushGoFunction(f)
	s.SetGlobal(name)
}

func (s *luaState) SetMetaTable(idx int) {
	val := s.luaStack.get(idx)
	mtVal := s.luaStack.pop()

	if mtVal == nil {
		setMetatable(val, nil, s)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, s)
	} else {
		panic("table expected!") // todo
	}
}

func (s *luaState) RawSet(idx int) {
	t := s.luaStack.get(idx)
	v := s.luaStack.pop()
	k := s.luaStack.pop()
	s.setTable(t, k, v, true)
}

func (s *luaState) RawSetI(idx int, i int64) {
	t := s.luaStack.get(idx)
	v := s.luaStack.pop()
	s.setTable(t, i, v, true)
}