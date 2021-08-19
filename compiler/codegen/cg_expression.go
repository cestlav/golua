package codegen

import . "golua/compiler/ast"
import . "golua/compiler/lexer"
import . "golua/vm"

// kind of operands
const (
	ARG_CONST = 1 // const index
	ARG_REG   = 2 // register index
	ARG_UPVAL = 4 // upvalue index
	ARG_RK    = ARG_REG | ARG_CONST
	ARG_RU    = ARG_REG | ARG_UPVAL
	ARG_RUK   = ARG_REG | ARG_UPVAL | ARG_CONST
)

// todo: rename to evalExp()?
func cgExp(fi *funcInfo, node Expression, a, n int) {
	switch exp := node.(type) {
	case *NilExpression:
		fi.emitLoadNil(exp.Line, a, n)
	case *FalseExpression:
		fi.emitLoadBool(exp.Line, a, 0, 0)
	case *TrueExpression:
		fi.emitLoadBool(exp.Line, a, 1, 0)
	case *IntegerExpression:
		fi.emitLoadK(exp.Line, a, exp.Value)
	case *FloatExpression:
		fi.emitLoadK(exp.Line, a, exp.Value)
	case *StringExpression:
		fi.emitLoadK(exp.Line, a, exp.Str)
	case *ParenthesesExpression:
		cgExp(fi, exp.Expression, a, 1)
	case *VarargExpression:
		cgVarargExp(fi, exp, a, n)
	case *FunctionDefineExpression:
		cgFuncDefExp(fi, exp, a)
	case *TableConstructorExpression:
		cgTableConstructorExp(fi, exp, a)
	case *UnaryOperatorExpression:
		cgUnopExp(fi, exp, a)
	case *BinaryOperatorExpression:
		cgBinopExp(fi, exp, a)
	case *ConcatExpression:
		cgConcatExp(fi, exp, a)
	case *NameExpression:
		cgNameExp(fi, exp, a)
	case *TableAccessExpression:
		cgTableAccessExp(fi, exp, a)
	case *FunctionCallExpression:
		cgFuncCallExp(fi, exp, a, n)
	}
}

func cgVarargExp(fi *funcInfo, node *VarargExpression, a, n int) {
	if !fi.isVararg {
		panic("cannot use '...' outside a vararg function")
	}
	fi.emitVararg(node.Line, a, n)
}

// f[a] := function(args) body end
func cgFuncDefExp(fi *funcInfo, node *FunctionDefineExpression, a int) {
	subFI := newFuncInfo(fi, node)
	fi.subFuncs = append(fi.subFuncs, subFI)

	for _, param := range node.ParamList {
		subFI.addLocalVariable(param, 0)
	}

	cgBlock(subFI, node.Block)
	subFI.exitScope(subFI.pc() + 2)
	subFI.emitReturn(node.LastLine, 0, 0)

	bx := len(fi.subFuncs) - 1
	fi.emitClosure(node.LastLine, a, bx)
}

func cgTableConstructorExp(fi *funcInfo, node *TableConstructorExpression, a int) {
	nArr := 0
	for _, keyExp := range node.KeyExpressions {
		if keyExp == nil {
			nArr++
		}
	}
	nExps := len(node.KeyExpressions)
	multRet := nExps > 0 &&
		isVarargOrFuncCall(node.ValueExpressions[nExps-1])

	fi.emitNewTable(node.Line, a, nArr, nExps-nArr)

	arrIdx := 0
	for i, keyExp := range node.KeyExpressions {
		valExp := node.ValueExpressions[i]

		if keyExp == nil {
			arrIdx++
			tmp := fi.allocReg()
			if i == nExps-1 && multRet {
				cgExp(fi, valExp, tmp, -1)
			} else {
				cgExp(fi, valExp, tmp, 1)
			}

			if arrIdx%50 == 0 || arrIdx == nArr { // LFIELDS_PER_FLUSH
				n := arrIdx % 50
				if n == 0 {
					n = 50
				}
				fi.freeRegs(n)
				line := lastLineOf(valExp)
				c := (arrIdx-1)/50 + 1 // todo: c > 0xFF
				if i == nExps-1 && multRet {
					fi.emitSetList(line, a, 0, c)
				} else {
					fi.emitSetList(line, a, n, c)
				}
			}

			continue
		}

		b := fi.allocReg()
		cgExp(fi, keyExp, b, 1)
		c := fi.allocReg()
		cgExp(fi, valExp, c, 1)
		fi.freeRegs(2)

		line := lastLineOf(valExp)
		fi.emitSetTable(line, a, b, c)
	}
}

// r[a] := op exp
func cgUnopExp(fi *funcInfo, node *UnaryOperatorExpression, a int) {
	oldRegs := fi.usedRegs
	b, _ := expToOpArg(fi, node.Expression, ARG_REG)
	fi.emitUnaryOp(node.Line, node.Operator, a, b)
	fi.usedRegs = oldRegs
}

// r[a] := exp1 op exp2
func cgBinopExp(fi *funcInfo, node *BinaryOperatorExpression, a int) {
	switch node.Operator {
	case TOKEN_OP_AND, TOKEN_OP_OR:
		oldRegs := fi.usedRegs

		b, _ := expToOpArg(fi, node.Expression1, ARG_REG)
		fi.usedRegs = oldRegs
		if node.Operator == TOKEN_OP_AND {
			fi.emitTestSet(node.Line, a, b, 0)
		} else {
			fi.emitTestSet(node.Line, a, b, 1)
		}
		pcOfJmp := fi.emitJmp(node.Line, 0, 0)

		b, _ = expToOpArg(fi, node.Expression2, ARG_REG)
		fi.usedRegs = oldRegs
		fi.emitMove(node.Line, a, b)
		fi.fixSbx(pcOfJmp, fi.pc()-pcOfJmp)
	default:
		oldRegs := fi.usedRegs
		b, _ := expToOpArg(fi, node.Expression1, ARG_RK)
		c, _ := expToOpArg(fi, node.Expression2, ARG_RK)
		fi.emitBinaryOp(node.Line, node.Operator, a, b, c)
		fi.usedRegs = oldRegs
	}
}

// r[a] := exp1 .. exp2
func cgConcatExp(fi *funcInfo, node *ConcatExpression, a int) {
	for _, subExp := range node.Expressions {
		a := fi.allocReg()
		cgExp(fi, subExp, a, 1)
	}

	c := fi.usedRegs - 1
	b := c - len(node.Expressions) + 1
	fi.freeRegs(c - b + 1)
	fi.emitABC(node.Line, OP_CONCAT, a, b, c)
}

// r[a] := name
func cgNameExp(fi *funcInfo, node *NameExpression, a int) {
	if r := fi.slotOfLocalVariable(node.Name); r >= 0 {
		fi.emitMove(node.Line, a, r)
	} else if idx := fi.indexOfUpValue(node.Name); idx >= 0 {
		fi.emitGetUpval(node.Line, a, idx)
	} else { // x => _ENV['x']
		taExp := &TableAccessExpression{
			LastLine:  node.Line,
			PrefixExpression: &NameExpression{node.Line, "_ENV"},
			KeyExpression:    &StringExpression{node.Line, node.Name},
		}
		cgTableAccessExp(fi, taExp, a)
	}
}

// r[a] := prefix[key]
func cgTableAccessExp(fi *funcInfo, node *TableAccessExpression, a int) {
	oldRegs := fi.usedRegs
	b, kindB := expToOpArg(fi, node.PrefixExpression, ARG_RU)
	c, _ := expToOpArg(fi, node.KeyExpression, ARG_RK)
	fi.usedRegs = oldRegs

	if kindB == ARG_UPVAL {
		fi.emitGetTabUp(node.LastLine, a, b, c)
	} else {
		fi.emitGetTable(node.LastLine, a, b, c)
	}
}

// r[a] := f(args)
func cgFuncCallExp(fi *funcInfo, node *FunctionCallExpression, a, n int) {
	nArgs := prepFuncCall(fi, node, a)
	fi.emitCall(node.Line, a, nArgs, n)
}

// return f(args)
func cgTailCallExp(fi *funcInfo, node *FunctionCallExpression, a int) {
	nArgs := prepFuncCall(fi, node, a)
	fi.emitTailCall(node.Line, a, nArgs)
}

func prepFuncCall(fi *funcInfo, node *FunctionCallExpression, a int) int {
	nArgs := len(node.Args)
	lastArgIsVarargOrFuncCall := false

	cgExp(fi, node.PrefixExpression, a, 1)
	if node.NameExpression != nil {
		fi.allocReg()
		c, k := expToOpArg(fi, node.NameExpression, ARG_RK)
		fi.emitSelf(node.Line, a, a, c)
		if k == ARG_REG {
			fi.freeRegs(1)
		}
	}
	for i, arg := range node.Args {
		tmp := fi.allocReg()
		if i == nArgs-1 && isVarargOrFuncCall(arg) {
			lastArgIsVarargOrFuncCall = true
			cgExp(fi, arg, tmp, -1)
		} else {
			cgExp(fi, arg, tmp, 1)
		}
	}
	fi.freeRegs(nArgs)

	if node.NameExpression != nil {
		fi.freeReg()
		nArgs++
	}
	if lastArgIsVarargOrFuncCall {
		nArgs = -1
	}

	return nArgs
}

func expToOpArg(fi *funcInfo, node Expression, argKinds int) (arg, argKind int) {
	if argKinds&ARG_CONST > 0 {
		idx := -1
		switch x := node.(type) {
		case *NilExpression:
			idx = fi.indexOfConstant(nil)
		case *FalseExpression:
			idx = fi.indexOfConstant(false)
		case *TrueExpression:
			idx = fi.indexOfConstant(true)
		case *IntegerExpression:
			idx = fi.indexOfConstant(x.Value)
		case *FloatExpression:
			idx = fi.indexOfConstant(x.Value)
		case *StringExpression:
			idx = fi.indexOfConstant(x.Str)
		}
		if idx >= 0 && idx <= 0xFF {
			return 0x100 + idx, ARG_CONST
		}
	}

	if nameExp, ok := node.(*NameExpression); ok {
		if argKinds&ARG_REG > 0 {
			if r := fi.slotOfLocalVariable(nameExp.Name); r >= 0 {
				return r, ARG_REG
			}
		}
		if argKinds&ARG_UPVAL > 0 {
			if idx := fi.indexOfUpValue(nameExp.Name); idx >= 0 {
				return idx, ARG_UPVAL
			}
		}
	}

	a := fi.allocReg()
	cgExp(fi, node, a, 1)
	return a, ARG_REG
}
