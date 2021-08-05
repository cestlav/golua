package state

import "golua/binary"

type luaClosure struct {
	proto *binary.ProtoType
}

func newLuaClosure(proto *binary.ProtoType) *luaClosure {
	return &luaClosure{
		proto: proto,
	}
}