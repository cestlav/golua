package parser

import (
	. "golua/compiler/ast"
	. "golua/compiler/lexer"
)

func Parse(chunk, chunkName string) *Block {
	lexer := NewLexer()
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_EOF)
	return block
}