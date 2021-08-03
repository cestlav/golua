package binary

import (
	"encoding/binary"
	"math"
)

type BinaryReader struct {
	data []byte
}

func (r *BinaryReader) ReadByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

func (r *BinaryReader) ReadBytes(size uint) []byte {
	bytes := r.data[:size]
	r.data = r.data[size:]
	return bytes
}

func (r *BinaryReader) ReadUint32() uint32 {
	i := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i
}

func (r *BinaryReader) ReadUint64() uint64 {
	i := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i
}

func (r *BinaryReader) ReadLuaInteger() int64 {
	return int64(r.ReadUint64());
}

func (r *BinaryReader) ReadLuaNumber() float64 {
	return math.Float64frombits(r.ReadUint64());
}

func (r *BinaryReader) ReadString() string {
	size := uint(r.ReadByte())
	if size == 0 {
		return ""
	}

	if size == 0xFF {
		size = uint(r.ReadUint64())
	}

	bytes := r.ReadBytes(size - 1)
	return string(bytes)
}

func (r *BinaryReader) CheckHeader() {
	if string(r.ReadBytes(4)) != LUA_SIGNATURE {
		panic("not a valid magic number")
	} else if r.ReadByte() != LUAC_VERSION {
		panic("invalid version")
	} else if r.ReadByte() != LUAC_FORMAT {
		panic("invalid format")
	} else if string(r.ReadBytes(6)) != LUAC_DATA {
		panic("corrupted")
	} else if r.ReadByte() != CINT_SIZE {
		panic("invalid int size")
	} else if r.ReadByte() != CSIZET_SIZE {
		panic("invalid size_t size")
	} else if r.ReadByte() != INSTRUCTION_SIZE {
		panic("invalid instruction size")
	} else if r.ReadByte() != LUA_INTEGER_SIZE {
		panic("invalid lua integer size")
	} else if r.ReadByte() != LUA_NUMBER_SIZE {
		panic("invalid lua_number size")
	} else if r.ReadLuaInteger() != LUAC_INT {
		panic("invalid endianness")
	} else if r.ReadLuaNumber() != LUAC_NUM {
		panic("invalid float number")
	}
}

func (r *BinaryReader) ReadProto(parentSource string) *ProtoType {
	source := r.ReadString()
	if source == "" {
		source = parentSource
	}

	return &ProtoType {
		Source: source,
		LineDefined: r.ReadUint32(),
		LastLineDefined: r.ReadUint32(),
		NumParams: r.ReadByte(),
		IsVararg: r.ReadByte(),
		MaxStackSize: r.ReadByte(),
		Code: r.ReadCode(),
		Constants: r.ReadConstants(),
		Upvalues: r.ReadUpvalues(),
		ProtoTypes: r.ReadProtos(source),
		LineInfo: r.ReadLineInfo(),
		LocalVariables: r.ReadLocalVariables(),
		UpvalueNames: r.ReadUpvalueNames(),
	}
}

func (r *BinaryReader) ReadCode() []uint32 {
	code := make([]uint32, r.ReadUint32())
	for i := range code {
		code[i] = r.ReadUint32()
	}
	return code
}

func (r *BinaryReader) ReadConstants() []interface{} {
	constants := make([]interface{}, r.ReadUint32())
	for i := range constants {
		constants[i] = r.ReadConstant()
	}
	return constants
}

func (r *BinaryReader) ReadConstant() interface{} {
	switch r.ReadByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return r.ReadByte() != 0;
	case TAG_INTEGER:
		return r.ReadLuaInteger();
	case TAG_NUMBER:
		return r.ReadLuaNumber();
	case TAG_SHORT_STR, TAG_LONG_STR:
		return r.ReadString();
	default:
		panic("invalid constant")
	}
}

func (r *BinaryReader) ReadUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.ReadUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue {
			Instack: r.ReadByte(),
			Index: r.ReadByte(),
		}
	}
	return upvalues
}

func (r *BinaryReader) ReadProtos(parentSource string) []*ProtoType {
	protos := make([]*ProtoType, r.ReadUint32())
	for i := range protos {
		protos[i] = r.ReadProto(parentSource)
	}
	return protos
}

func (r *BinaryReader) ReadLineInfo() []uint32 {
	lineInfo := make([]uint32, r.ReadUint32())
	for i := range lineInfo {
		lineInfo[i] = r.ReadUint32()
	}
	return lineInfo
}

func (r *BinaryReader) ReadLocalVariables() []LocalVariable {
	localVars := make([]LocalVariable, r.ReadUint32())
	for i := range localVars {
		localVars[i] = LocalVariable {
			VariableName: r.ReadString(),
			StartPC: r.ReadUint32(),
			EndPC: r.ReadUint32(),
		}
	}
	return localVars
}

func (r *BinaryReader) ReadUpvalueNames() []string {
	names := make([]string, r.ReadUint32())
	for i := range names {
		names[i] = r.ReadString()
	}
	return names
}