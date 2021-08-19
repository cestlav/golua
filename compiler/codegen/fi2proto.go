package codegen

import . "golua/binary"

func toProtoType(fi *funcInfo) *ProtoType {
	proto := &ProtoType{
		LineDefined:     uint32(fi.line),
		LastLineDefined: uint32(fi.lastLine),
		NumParams:       byte(fi.numParams),
		MaxStackSize:    byte(fi.maxRegs),
		Code:            fi.insts,
		Constants:       getConstants(fi),
		UpValues:        getUpValues(fi),
		ProtoTypes:      toProtoTypes(fi.subFuncs),
		LineInfo:        fi.lineNums,
		LocalVariables:  getLocalVariables(fi),
		UpvalueNames:    getUpValueNames(fi),
	}

	if fi.line == 0 {
		proto.LastLineDefined = 0
	}
	if proto.MaxStackSize < 2 {
		proto.MaxStackSize = 2 // todo
	}
	if fi.isVararg {
		proto.IsVararg = 1 // todo
	}

	return proto
}

func toProtoTypes(fis []*funcInfo) []*ProtoType {
	protos := make([]*ProtoType, len(fis))
	for i, fi := range fis {
		protos[i] = toProtoType(fi)
	}
	return protos
}

func getConstants(fi *funcInfo) []interface{} {
	consts := make([]interface{}, len(fi.constants))
	for k, idx := range fi.constants {
		consts[idx] = k
	}
	return consts
}

func getLocalVariables(fi *funcInfo) []LocalVariable {
	locVars := make([]LocalVariable, len(fi.locVars))
	for i, locVar := range fi.locVars {
		locVars[i] = LocalVariable{
			VariableName: locVar.name,
			StartPC: uint32(locVar.startPC),
			EndPC:   uint32(locVar.endPC),
		}
	}
	return locVars
}

func getUpValues(fi *funcInfo) []UpValue {
	upvals := make([]UpValue, len(fi.upvalues))
	for _, uv := range fi.upvalues {
		if uv.locVarSlot >= 0 { // instack
			upvals[uv.index] = UpValue{1, byte(uv.locVarSlot)}
		} else {
			upvals[uv.index] = UpValue{0, byte(uv.upvalIndex)}
		}
	}
	return upvals
}

func getUpValueNames(fi *funcInfo) []string {
	names := make([]string, len(fi.upvalues))
	for name, uv := range fi.upvalues {
		names[uv.index] = name
	}
	return names
}
