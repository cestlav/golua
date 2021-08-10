package ast

type Statement interface {}

type EmptyStatement struct {}

type BreakStatement struct {Line int}

type LabelStatement struct {Name string}

type GotoStatement struct {Name string}

type DoStatement struct {Block *Block}

type FuncCallStatement = FunctionCallExpression

type WhileStatement struct {
	Expression   Expression
	Block *Block
}

type RepeatStatement struct {
	Expression   Expression
	Block *Block
}

type IfStatement struct {
	Expression []Expression
	Blocks []*Block
}

type ForNumStatement struct {
	LineOfFor int
	LineOfDo  int
	VarName   string
	InitExpression   Expression
	LimitExpression Expression
	StepExpression  Expression
	Block     *Block
}

type ForInStatement struct {
	LineOfDo int
	NameList []string
	ExpressionList []Expression
	Block *Block
}

type LocalVariableDeclareStatement struct {
	LastLine int
	NameList []string
	ExpressionList []Expression
}

type AssignStatement struct {
	LastLine int
	VariableList []Expression
	ExpressionList []Expression
}

type LocalFunctionDefineStatement struct {
	Name string
	Expression *FunctionDefineExpression
}