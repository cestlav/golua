package state

import "golua/api"

type luaStack struct {
	slots []luaValue
	top int

	prev *luaStack
	closure *closure
	varargs []luaValue
	pc int

	state *luaState
}

func newLuaStack(size int, s *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: s,
	}
}

func (s *luaStack) check(n int)  {
	free := len(s.slots) - s.top
	for i := free; i < n; i++ {
		s.slots = append(s.slots, nil)
	}
}

func (s *luaStack) push(value luaValue)  {
	if s.top == len(s.slots) {
		panic("stack overflow")
	}
	s.slots[s.top] = value
	s.top++
}

func (s *luaStack) pop() luaValue {
	if s.top < 1 {
		panic("stack underflow")
	}
	s.top--
	v := s.slots[s.top]
	s.slots[s.top] = nil
	return v
}

func (s *luaStack) absIndex(index int) int {
	if index >= 0 || index <= api.LUA_REGISTRYINDEX {
		return index
	}
	return index + s.top +1
}

func (s *luaStack) isValid(index int) bool {
	if index == api.LUA_REGISTRYINDEX {
		return true
	}
	absIndex := s.absIndex(index)
	return absIndex > 0 && absIndex <= s.top
}

func (s *luaStack) get(index int) luaValue {
	if index == api.LUA_REGISTRYINDEX {
		return s.state.registry
	}
	absindex := s.absIndex(index)
	if absindex > 0 && absindex <= s.top {
		return s.slots[absindex - 1]
	}
	return nil
}

func (s *luaStack) set(index int, v luaValue) {
	if index == api.LUA_REGISTRYINDEX {
		s.state.registry = v.(*luaTable)
		return
	}
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

func (s *luaStack) popN(n int) []luaValue {
	res := make([]luaValue, n)
	for i := n; i > 0; i-- {
		res[i - 1] = s.pop()
	}
	return res
}

func (s *luaStack) pushN(vals []luaValue, n int)  {
	nVals  := len(vals)
	if n < 0 { n = nVals}

	for i := 0; i < n; i++ {
		if i < nVals {
			s.push(vals[i])
		} else {
			s.push(nil)
		}
	}
}