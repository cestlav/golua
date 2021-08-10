package ast

type Block struct {
	LastLine int
	Stats []Statement
	RetExps []Expression
}