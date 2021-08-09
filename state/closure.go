package state

import (
	"golua/api"
	"golua/binary"
)

type closure struct {
	proto *binary.ProtoType
	goFunc api.GoFunction
	upValues []*upValue
}

type upValue struct {
	value *luaValue
}

func newLuaClosure(proto *binary.ProtoType) *closure {
	c := &closure{
		proto: proto,
	}

	if nUpValues := len(proto.UpValues); nUpValues > 0 {
		c.upValues = make([]*upValue, nUpValues)
	}
	return c
}

func newGoClosure(f api.GoFunction, nUpValues int) *closure {
	c:=  &closure {
		goFunc: f,
	}

	if nUpValues > 0 {
		c.upValues = make([]*upValue, nUpValues)
	}
	return c
}