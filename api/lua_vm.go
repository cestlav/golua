package api

type LuaVM interface {
	LuaState

	PC() int
	AddPC(n int)
	Fetch() uint32
	GetConst(int)
	GetRK(int)

	RegisterCount() int
	LoadVararg(int)
	LoadProto(int)

	CloseUpValues(int)
}