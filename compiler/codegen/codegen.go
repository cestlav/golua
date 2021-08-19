package codegen

import . "golua/binary"
import . "golua/compiler/ast"

func GenProto(chunk *Block) *ProtoType {
	fd := &FunctionDefineExpression{
		LastLine: chunk.LastLine,
		IsVararg: true,
		Block:    chunk,
	}

	fi := newFuncInfo(nil, fd)
	fi.addLocalVariable("_ENV", 0)
	cgFuncDefExp(fi, fd, 0)
	return toProtoType(fi.subFuncs[0])
}
