package state

import (
	"golua/api"
	"golua/binary"
)

type closure struct {
	proto *binary.ProtoType
	goFunc api.GoFunction
}

func newLuaClosure(proto *binary.ProtoType) *closure {
	return &closure{
		proto: proto,
	}
}

func newGoClosure(f api.GoFunction) *closure {
	return &closure{
		goFunc: f,
	}
}