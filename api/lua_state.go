package api

type LuaType int
type ArithOp = int
type CompareOp = int

type LuaState interface {
	GetTop() int
	AbsIndex(int) int
	CheckStack(int) bool
	Pop(int)
	Copy(int, int)
	PushValue(int)
	Replace(int)
	Insert(int)
	Remove(int)
	Rotate(int, int)
	SetTop(int)

	TypeName(LuaType) string
	Type(int) LuaType
	IsNone(int) bool
	IsNil(int) bool
	IsNoneOrNil(int) bool
	IsBoolean(int) bool
	IsInteger(int) bool
	IsNumber(int) bool
	IsString(int) bool
	IsTable(int) bool
	IsThread(int) bool
	IsFunction(int) bool
	ToBoolean(int) bool
	ToInteger(int) int64
	ToIntegerX(int) (int64, bool)
	ToNumber(int) float64
	ToNumberX(int) (float64, bool)
	ToString(int) string
	ToStringX(int) (string, bool)

	PushNil()
	PushBoolean(bool)
	PushInteger(int64)
	PushNumber(float64)
	PushString(string)

	Arith(ArithOp)
	Compare(int, int, CompareOp) bool
	Len(int)
	Concat(int)

	NewTable()
	CreateTable(int, int)
	GetTable(int) LuaType
	GetField(int, string) LuaType
	GetI(int, int64) LuaType

	SetTable(int)
	SetField(int, string)
	SetI(int, int64)

	Load([]byte, string, string) int
	Call(int, int)

}