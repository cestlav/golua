package state

import . "golua/api"

type luaState struct {
	registry *luaTable
	luaStack *luaStack
}

func NewLuaState() *luaState {
	registry := newLuaTable(0, 0)
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0))

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
	return ls
}

func (s *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = s.luaStack
	s.luaStack = stack
}

func (s *luaState) popLuaStack()  {
	stack := s.luaStack
	s.luaStack = stack.prev
	stack.prev = nil
}