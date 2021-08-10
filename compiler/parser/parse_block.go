package parser

import (
	. "golua/compiler/ast"
	. "golua/compiler/lexer"
)

func parseBlock(lexer *Lexer) *Block {
	return &Block{
		LastLine: lexer.CurrentLine(),
		Statements:    parseStatements(lexer),
		RetExpressions:  parseRetExpressions(lexer),
	}
}

func parseStatements(lexer *Lexer) []Statement {
	statements := make([]Statement, 0, 8)
	for !isReturnOrBlockEnd(lexer.LookAhead()) {
		statement := parseStatement(lexer)
		if _, ok := statement.(*EmptyStatement); !ok {
			statements = append(statements, statement)
		}
	}
	return statements
}

func isReturnOrBlockEnd(tokenKind int) bool {
	switch tokenKind {
	case TOKEN_KW_RETURN, TOKEN_EOF, TOKEN_KW_END, TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
		return true
	}
	return false
}

func parseRetExpressions(lexer *Lexer) []Expression {
	if lexer.LookAhead() != TOKEN_KW_RETURN {
		return nil
	}

	lexer.NextToken()
	switch lexer.LookAhead() {
	case TOKEN_EOF, TOKEN_KW_END, TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
		return []Expression{}
	case TOKEN_SEP_SEMI:
		lexer.NextToken()
		return []Expression{}
	default:
		exps := parseExpressionList(lexer)
		if lexer.LookAhead() == TOKEN_SEP_SEMI {
			lexer.NextToken()
		}
		return exps
	}
}
