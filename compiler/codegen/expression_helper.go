package codegen

import . "golua/compiler/ast"

func isVarargOrFuncCall(exp Expression) bool {
	switch exp.(type) {
	case *VarargExpression, *FunctionCallExpression:
		return true
	}
	return false
}

func removeTailNils(exps []Expression) []Expression {
	for n := len(exps) - 1; n >= 0; n-- {
		if _, ok := exps[n].(*NilExpression); !ok {
			return exps[0 : n+1]
		}
	}
	return nil
}

func lineOf(exp Expression) int {
	switch x := exp.(type) {
	case *NilExpression:
		return x.Line
	case *TrueExpression:
		return x.Line
	case *FalseExpression:
		return x.Line
	case *IntegerExpression:
		return x.Line
	case *FloatExpression:
		return x.Line
	case *StringExpression:
		return x.Line
	case *VarargExpression:
		return x.Line
	case *NameExpression:
		return x.Line
	case *FunctionDefineExpression:
		return x.Line
	case *FunctionCallExpression:
		return x.Line
	case *TableConstructorExpression:
		return x.Line
	case *UnaryOperatorExpression:
		return x.Line
	case *TableAccessExpression:
		return lineOf(x.PrefixExpression)
	case *ConcatExpression:
		return lineOf(x.Expressions[0])
	case *BinaryOperatorExpression:
		return lineOf(x.Expression1)
	default:
		panic("unreachable!")
	}
}

func lastLineOf(exp Expression) int {
	switch x := exp.(type) {
	case *NilExpression:
		return x.Line
	case *TrueExpression:
		return x.Line
	case *FalseExpression:
		return x.Line
	case *IntegerExpression:
		return x.Line
	case *FloatExpression:
		return x.Line
	case *StringExpression:
		return x.Line
	case *VarargExpression:
		return x.Line
	case *NameExpression:
		return x.Line
	case *FunctionDefineExpression:
		return x.LastLine
	case *FunctionCallExpression:
		return x.LastLine
	case *TableConstructorExpression:
		return x.LastLine
	case *TableAccessExpression:
		return x.LastLine
	case *ConcatExpression:
		return lastLineOf(x.Expressions[len(x.Expressions)-1])
	case *BinaryOperatorExpression:
		return lastLineOf(x.Expression2)
	case *UnaryOperatorExpression:
		return lastLineOf(x.Expression)
	default:
		panic("unreachable!")
	}
}
