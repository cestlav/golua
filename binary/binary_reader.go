package binary

import (
	"encoding/binary"
	"math"
)

type BinaryReader struct {
	data []byte
}

func (self *BinaryReader) ReadByte() byte {
	b := self.data[0]
	self.data = self.data[1:]
	return b
}

func (self *BinaryReader) ReadBytes(size uint) []byte {
	bytes := self.data[:size]
	self.data = self.data[n:]
	return bytes
}

func (self *BinaryReader) ReadUint32() uint32 {
	i := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return i
}

func (self *BinaryReader) ReadUint64() uint64 {
	i := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return i
}

func (self *BinaryReader) ReadLuaInteger() int64 {
	return int64(self.ReadUint64());
}

func (self *BinaryReader) ReadLuaNumber() float64 {
	return math.Float64frombits(self.ReadUint64());
}

func (self *BinaryReader) ReadString() string {
	size := uint(self.ReadByte())
	if size == 0 {
		return ""
	}

	if size == 0xFF {
		size = uint(self.ReadUint64())
	}

	bytes := self.ReadBytes(size - 1)
	return string(bytes)
}

func (self *BinaryReader)CheckHeader() {
	if string(self.ReadBytes(4) != LUA_SIGNATURE) {
		panic("not a valid magic number")
	} else if self.ReadByte() != LUAC_VERSION {
		panic("invalid version")
	} else if self.ReadByte() != LUAC_FORMAT {
		panic("invalid format")
	} else if string(self.ReadBytes(6)) != LUAC_DATA {
		panic("corrupted")
	} else if self.ReadByte() != CINT_SIZE {
		panic("invalid int size")
	} else if self.ReadByte() != CSIZET_SIZE {
		panic("invalid size_t size")
	} else if self.ReadByte() != INSTRUCTION_SIZE {
		panic("invalid instruction size")
	} else if self.ReadByte() != LUA_INTEGER_SIZE {
		panic("invalid lua integer size")
	} else if self.ReadByte() != LUA_NUMBER_SIZE {
		panic("invalid lua_number size")
	} else if self.ReadLuaInteger() != LUAC_INT {
		panic("invalid endianness")
	} else if self.ReadLuaNumber() != LUAC_NUM {
		panic("invalid float number")
	}
}

func (self *BinaryReader) ReadProto(parentSource string) *ProtoType {
	source := self.ReadString()
	if source == "" {
		source = parentSource
	}

	return &ProtoType {
		Source: source,
		LineDefined: self.ReadUint32(),
		LastLineDefined: self.ReadUint32(),
		NumParams: self.ReadByte(),
		IsVararg: self.ReadByte(),
		MaxStackSize: self.ReadByte(),
		Code: self.ReadCode(),
		Constants: self.ReadConstants(),
		Upvalue: self.ReadUpvalues(),
		ProtoTypes: self.ReadProtos(source),
		LineInfo: self.ReadLineInfo(),
		LocalVariables: self.ReadLocalVariables(),
		UpvalueNames: self.ReadUpvalueNames(),
	}
}

func (self *BinaryReader)ReadCode() []uint32 {
	code := make([]uint32, self.ReadUint32())
	for i := range code {
		code[i] = self.ReadUint32()
	}
	return code
}

func (self *BinaryReader)ReadConstants() []interface{} {
	constants := make([]interface{}, self.ReadUint32)
	for i := range constants {
		constants[i] = self.ReadConstant()
	}
}

func (self *BinaryReader)ReadConstant() interface{} {
	switch self.ReadByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return self.ReadByte() != 0;
	case TAG_INTEGER:
		return self.ReadLuaInteger();
	case TAG_NUMBER:
		return self.ReadLuaNumber();
	case TAG_SHORT_STR, TAG_LONG_STR:
		return self.ReadString();
	default:
		panic("invalid constant")
	}
}

func (self *BinaryReader)ReadUpvalues() []Upvalue {
	upvalues := make([]Upvalue, self.ReadUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue {
			Instack: self.ReadByte(),
			Index: self.ReadByte(),
		}
	}
	return upvalues
}

func (self *BinaryReader)ReadProtos(parentSource string) []*ProtoType {
	protos := make([]*ProtoType, self.ReadUint32())
	for i := range protos {
		protos[i] = self.ReadProto(parentSource)
	}
	return protos
}

func (self *BinaryReader)ReadLineInfo() []uint32 {
	lineInfo := make([]uint32, sefl.ReadUint32())
	for i := range lineInfo {
		lineInfo[i] = self.ReadUint32()
	}
	return lineInfo
}

func (self *BinaryReader)ReadLocalVariables() []LocalVariable {
	localVars := make([]LocalVariable, self.ReadUint32())
	for i := range localVars {
		localVars[i] = LocalVariable {
			VariableName: self.ReadString(),
			StartPC: self.ReadUint32(),
			EndPC: self.ReadUint32(),
		}
	}
	return localVars
}

func (self *BinaryReader)ReadUpvalueNames() []string {
	names := make([]string, self.ReadUint32())
	for i := range names {
		names[i] = self.ReadString()
	}
	return names
}