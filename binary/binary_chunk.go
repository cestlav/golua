package binary

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type BinaryChunk struct {
	Header
	SizeUpvalue  byte
	MainFunction *ProtoType
}

type Header struct {
	Signature       [4]byte
	Version         byte
	Format          byte
	LuacData        [6]byte
	CintSize        byte
	SizetSize       byte
	InstructionSize byte
	LuaIntegerSize  byte
	LuaNumberSize   byte
	LuacInt         int64
	LuacNum         float64
}

type ProtoType struct {
	Source          string
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	UpValues        []UpValue
	ProtoTypes      []*ProtoType
	LineInfo        []uint32
	LocalVariables  []LocalVariable
	UpvalueNames    []string
}

type UpValue struct {
	InStack byte
	Index   byte
}

type LocalVariable struct {
	VariableName string
	StartPC      uint32
	EndPC        uint32
}

func Undump(data []byte) *ProtoType {
	reader := &BinaryReader{data}
	reader.CheckHeader()
	reader.ReadByte()
	return reader.ReadProto("")
}
