package ast

type Expression interface {}

type NilExpression struct {Line int}

type TrueExpression struct {Line int}

type FalseExpression struct {Lint int}

type VarargExpression struct {Line int}

type IntegerExpression struct {Line int; Value int64}

type FloatExpression struct {Line int; Value float64}

type StringExpression struct {Line int; Value string}

type NameExpression struct {Line int; Name string}

type UnaryOperatorExpression struct {
	Line int
	Operator int
	Expression Expression
}

type BinaryOperatorExpression struct {
	Line int
	Operator int
	Expression1 Expression
	Expression2 Expression
}

type ConcatExpression struct {
	Line int
	Expressions []Expression
}

type TableConstructionExpression struct {
	Line int
	LastLine int
	KeyExpressions []Expression
	ValueExpressions []Expression
}

type TableAccessExpression struct {
	LastLine int
	PrefixExpression Expression
	KeyExpression Expression
}

type FunctionDefineExpression struct {
	Line int
	LastLine int
	ParamList []string
	IsVararg bool
	Block *Block
}

type ParenthesesExpression struct {
	Expression Expression
}

type FunctionCallExpression struct {
	Line int
	LastLine int
	PrefixExpression Expression
	NameExpression *StringExpression
	Args []Expression
}