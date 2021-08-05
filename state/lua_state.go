package state

type luaState struct {
	luaStack *luaStack
}

func NewLuaState(stackSize int) *luaState {
	return &luaState{
		luaStack: newLuaStack(stackSize),
	}
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