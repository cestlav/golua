package ast

type Block struct {
	LastLine int
	Statements []Statement
	RetExpressions []Expression
}