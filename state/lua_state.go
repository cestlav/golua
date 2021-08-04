package state

import "golua/binary"

type luaState struct {
	luaStack *luaStack

	proto *binary.ProtoType
	pc int
}

func NewLuaState(stackSize int, proto *binary.ProtoType) *luaState {
	return &luaState{
		luaStack: newLuaStack(stackSize),

		proto: proto,
		pc: 0,
	}
}