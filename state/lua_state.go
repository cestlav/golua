package state

type luaState struct {
	luaStack *luaStack
}

func NewLuaState() *luaState {
	return &luaState{luaStack: newLuaStack(20)}
}