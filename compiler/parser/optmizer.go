package parser

import (
	. "golua/compiler/ast"
	. "golua/compiler/lexer"
	"golua/number"
	"math"
)

func optimizeLogicalOr(exp *BinaryOperatorExpression) Expression {
	if isTrue(exp.Expression1) {
		return exp.Expression1 // true or x => true
	}
	if isFalse(exp.Expression1) && !isVarargOrFuncCall(exp.Expression2) {
		return exp.Expression2 // false or x => x
	}
	return exp
}

func optimizeLogicalAnd(exp *BinaryOperatorExpression) Expression {
	if isFalse(exp.Expression1) {
		return exp.Expression1 // false and x => false
	}
	if isTrue(exp.Expression1) && !isVarargOrFuncCall(exp.Expression2) {
		return exp.Expression2 // true and x => x
	}
	return exp
}

func optimizeBitwiseBinaryOp(exp *BinaryOperatorExpression) Expression {
	if i, ok := castToInt(exp.Expression1); ok {
		if j, ok := castToInt(exp.Expression2); ok {
			switch exp.Operator {
			case TOKEN_OP_BAND:
				return &IntegerExpression{exp.Line, i & j}
			case TOKEN_OP_BOR:
				return &IntegerExpression{exp.Line, i | j}
			case TOKEN_OP_BXOR:
				return &IntegerExpression{exp.Line, i ^ j}
			case TOKEN_OP_SHL:
				return &IntegerExpression{exp.Line, number.ShiftLeft(i, j)}
			case TOKEN_OP_SHR:
				return &IntegerExpression{exp.Line, number.ShiftRight(i, j)}
			}
		}
	}
	return exp
}

func optimizeArithBinaryOp(exp *BinaryOperatorExpression) Expression {
	if x, ok := exp.Expression1.(*IntegerExpression); ok {
		if y, ok := exp.Expression2.(*IntegerExpression); ok {
			switch exp.Operator {
			case TOKEN_OP_ADD:
				return &IntegerExpression{exp.Line, x.Value + y.Value}
			case TOKEN_OP_SUB:
				return &IntegerExpression{exp.Line, x.Value - y.Value}
			case TOKEN_OP_MUL:
				return &IntegerExpression{exp.Line, x.Value * y.Value}
			case TOKEN_OP_IDIV:
				if y.Value != 0 {
					return &IntegerExpression{exp.Line, number.IFloorDiv(x.Value, y.Value)}
				}
			case TOKEN_OP_MOD:
				if y.Value != 0 {
					return &IntegerExpression{exp.Line, number.IMod(x.Value, y.Value)}
				}
			}
		}
	}
	if f, ok := castToFloat(exp.Expression1); ok {
		if g, ok := castToFloat(exp.Expression2); ok {
			switch exp.Operator {
			case TOKEN_OP_ADD:
				return &FloatExpression{exp.Line, f + g}
			case TOKEN_OP_SUB:
				return &FloatExpression{exp.Line, f - g}
			case TOKEN_OP_MUL:
				return &FloatExpression{exp.Line, f * g}
			case TOKEN_OP_DIV:
				if g != 0 {
					return &FloatExpression{exp.Line, f / g}
				}
			case TOKEN_OP_IDIV:
				if g != 0 {
					return &FloatExpression{exp.Line, number.FFloorDiv(f, g)}
				}
			case TOKEN_OP_MOD:
				if g != 0 {
					return &FloatExpression{exp.Line, number.FMod(f, g)}
				}
			case TOKEN_OP_POW:
				return &FloatExpression{exp.Line, math.Pow(f, g)}
			}
		}
	}
	return exp
}

func optimizePow(exp Expression) Expression {
	if binop, ok := exp.(*BinaryOperatorExpression); ok {
		if binop.Operator == TOKEN_OP_POW {
			binop.Expression2 = optimizePow(binop.Expression2)
		}
		return optimizeArithBinaryOp(binop)
	}
	return exp
}

func optimizeUnaryOp(exp *UnaryOperatorExpression) Expression {
	switch exp.Operator {
	case TOKEN_OP_UNM:
		return optimizeUnm(exp)
	case TOKEN_OP_NOT:
		return optimizeNot(exp)
	case TOKEN_OP_BNOT:
		return optimizeBnot(exp)
	default:
		return exp
	}
}

func optimizeUnm(exp *UnaryOperatorExpression) Expression {
	switch x := exp.Expression.(type) { // number?
	case *IntegerExpression:
		x.Value = -x.Value
		return x
	case *FloatExpression:
		if x.Value != 0 {
			x.Value = -x.Value
			return x
		}
	}
	return exp
}

func optimizeNot(exp *UnaryOperatorExpression) Expression {
	switch exp.Expression.(type) {
	case *NilExpression, *FalseExpression: // false
		return &TrueExpression{exp.Line}
	case *TrueExpression, *IntegerExpression, *FloatExpression, *StringExpression: // true
		return &FalseExpression{exp.Line}
	default:
		return exp
	}
}

func optimizeBnot(exp *UnaryOperatorExpression) Expression {
	switch x := exp.Expression.(type) { // number?
	case *IntegerExpression:
		x.Value = ^x.Value
		return x
	case *FloatExpression:
		if i, ok := number.FloatToInteger(x.Value); ok {
			return &IntegerExpression{x.Line, ^i}
		}
	}
	return exp
}

func isFalse(exp Expression) bool {
	switch exp.(type) {
	case *FalseExpression, *NilExpression:
		return true
	default:
		return false
	}
}

func isTrue(exp Expression) bool {
	switch exp.(type) {
	case *TrueExpression, *IntegerExpression, *FloatExpression, *StringExpression:
		return true
	default:
		return false
	}
}

// todo
func isVarargOrFuncCall(exp Expression) bool {
	switch exp.(type) {
	case *VarargExpression, *FunctionCallExpression:
		return true
	}
	return false
}

func castToInt(exp Expression) (int64, bool) {
	switch x := exp.(type) {
	case *IntegerExpression:
		return x.Value, true
	case *FloatExpression:
		return number.FloatToInteger(x.Value)
	default:
		return 0, false
	}
}

func castToFloat(exp Expression) (float64, bool) {
	switch x := exp.(type) {
	case *IntegerExpression:
		return float64(x.Value), true
	case *FloatExpression:
		return x.Value, true
	default:
		return 0, false
	}
}
