package parser

import (
	. "golua/compiler/ast"
	. "golua/compiler/lexer"
	"golua/number"
)

func parseExpressionList(lexer *Lexer) []Expression {
	exps := make([]Expression, 0, 4)
	exps = append(exps, parseExpression(lexer))
	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()
		exps = append(exps, parseExpression(lexer))
	}
	return exps
}

func parseExpression(lexer *Lexer) Expression {
	return parseExpression12(lexer)
}

func parseExpression12(lexer *Lexer) Expression {
	expression := parseExpression11(lexer)
	for lexer.LookAhead() == TOKEN_OP_OR {
		line, op, _ := lexer.NextToken()
		lor := &BinaryOperatorExpression{line , op, expression, parseExpression11(lexer)}
		expression = optimizeLogicalOr(lor)
	}
	return expression
}

func parseExpression11(lexer *Lexer) Expression {
	expression := parseExpression10(lexer)
	for lexer.LookAhead() == TOKEN_OP_AND {
		line, op, _ := lexer.NextToken()
		land := &BinaryOperatorExpression{line, op, expression, parseExpression10(lexer)}
		expression = optimizeLogicalAnd(land)
	}
	return expression
}

func parseExpression10(lexer *Lexer) Expression {
	expression := parseExpression9(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_LT, TOKEN_OP_GT, TOKEN_OP_NE, TOKEN_OP_LE, TOKEN_OP_GE, TOKEN_OP_EQ:
			line, op, _ := lexer.NextToken()
			expression = &BinaryOperatorExpression{line, op, expression, parseExpression9(lexer)}
		default:
			return expression
		}
	}
	return expression
}

func parseExpression9(lexer *Lexer) Expression {
	expression := parseExpression8(lexer)
	for lexer.LookAhead() == TOKEN_OP_BOR {
		line, op, _ := lexer.NextToken()
		bor := &BinaryOperatorExpression{line, op, expression, parseExpression8(lexer)}
		expression = optimizeBitwiseBinaryOp(bor)
	}
	return expression
}

func parseExpression8(lexer *Lexer) Expression {
	exp := parseExpression7(lexer)
	for lexer.LookAhead() == TOKEN_OP_BXOR {
		line, op, _ := lexer.NextToken()
		bxor := &BinaryOperatorExpression{line, op, exp, parseExpression7(lexer)}
		exp = optimizeBitwiseBinaryOp(bxor)
	}
	return exp
}

// x & y
func parseExpression7(lexer *Lexer) Expression {
	exp := parseExpression6(lexer)
	for lexer.LookAhead() == TOKEN_OP_BAND {
		line, op, _ := lexer.NextToken()
		band := &BinaryOperatorExpression{line, op, exp, parseExpression6(lexer)}
		exp = optimizeBitwiseBinaryOp(band)
	}
	return exp
}

// shift
func parseExpression6(lexer *Lexer) Expression {
	exp := parseExpression5(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_SHL, TOKEN_OP_SHR:
			line, op, _ := lexer.NextToken()
			shx := &BinaryOperatorExpression{line, op, exp, parseExpression5(lexer)}
			exp = optimizeBitwiseBinaryOp(shx)
		default:
			return exp
		}
	}
	return exp
}

// a .. b
func parseExpression5(lexer *Lexer) Expression {
	exp := parseExpression4(lexer)
	if lexer.LookAhead() != TOKEN_OP_CONCAT {
		return exp
	}

	line := 0
	exps := []Expression{exp}
	for lexer.LookAhead() == TOKEN_OP_CONCAT {
		line, _, _ = lexer.NextToken()
		exps = append(exps, parseExpression4(lexer))
	}
	return &ConcatExpression{line, exps}
}

// x +/- y
func parseExpression4(lexer *Lexer) Expression {
	exp := parseExpression3(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_ADD, TOKEN_OP_SUB:
			line, op, _ := lexer.NextToken()
			arith := &BinaryOperatorExpression{line, op, exp, parseExpression3(lexer)}
			exp = optimizeArithBinaryOp(arith)
		default:
			return exp
		}
	}
	return exp
}

// *, %, /, //
func parseExpression3(lexer *Lexer) Expression {
	exp := parseExpression2(lexer)
	for {
		switch lexer.LookAhead() {
		case TOKEN_OP_MUL, TOKEN_OP_MOD, TOKEN_OP_DIV, TOKEN_OP_IDIV:
			line, op, _ := lexer.NextToken()
			arith := &BinaryOperatorExpression{line, op, exp, parseExpression2(lexer)}
			exp = optimizeArithBinaryOp(arith)
		default:
			return exp
		}
	}
	return exp
}

// unary
func parseExpression2(lexer *Lexer) Expression {
	switch lexer.LookAhead() {
	case TOKEN_OP_UNM, TOKEN_OP_BNOT, TOKEN_OP_LEN, TOKEN_OP_NOT:
		line, op, _ := lexer.NextToken()
		exp := &UnaryOperatorExpression{line, op, parseExpression2(lexer)}
		return optimizeUnaryOp(exp)
	}
	return parseExpression1(lexer)
}

// x ^ y
func parseExpression1(lexer *Lexer) Expression { // pow is right associative
	exp := parseExp0(lexer)
	if lexer.LookAhead() == TOKEN_OP_POW {
		line, op, _ := lexer.NextToken()
		exp = &BinaryOperatorExpression{line, op, exp, parseExpression1(lexer)}
	}
	return optimizePow(exp)
}

func parseExp0(lexer *Lexer) Expression {
	switch lexer.LookAhead() {
	case TOKEN_VARARG: // ...
		line, _, _ := lexer.NextToken()
		return &VarargExpression{line}
	case TOKEN_KW_NIL: // nil
		line, _, _ := lexer.NextToken()
		return &NilExpression{line}
	case TOKEN_KW_TRUE: // true
		line, _, _ := lexer.NextToken()
		return &TrueExpression{line}
	case TOKEN_KW_FALSE: // false
		line, _, _ := lexer.NextToken()
		return &FalseExpression{line}
	case TOKEN_STRING: // LiteralString
		line, _, token := lexer.NextToken()
		return &StringExpression{line, token}
	case TOKEN_NUMBER: // Numeral
		return parseNumberExpression(lexer)
	case TOKEN_SEP_LCURLY: // tableconstructor
		return parseTableConstructorExpression(lexer)
	case TOKEN_KW_FUNCTION: // functiondef
		lexer.NextToken()
		return parseFunctionDefineExpression(lexer)
	default: // prefixexp
		return parsePrefixExp(lexer)
	}
}

func parseNumberExpression(lexer *Lexer) Expression {
	line, _, token := lexer.NextToken()
	if i, ok := number.ParseInteger(token); ok {
		return &IntegerExpression{line, i}
	} else if f, ok := number.ParseFloat(token); ok {
		return &FloatExpression{line, f}
	} else { // todo
		panic("not a number: " + token)
	}
}

func parseFunctionDefineExpression(lexer *Lexer) *FunctionDefineExpression {
	line := lexer.CurrentLine()                               // function
	lexer.NextTokenOfKind(TOKEN_SEP_LPAREN)            // (
	parList, isVararg := _parseParamList(lexer)          // [parlist]
	lexer.NextTokenOfKind(TOKEN_SEP_RPAREN)            // )
	block := parseBlock(lexer)                         // block
	lastLine, _ := lexer.NextTokenOfKind(TOKEN_KW_END) // end
	return &FunctionDefineExpression{line, lastLine, parList, isVararg, block}
}

// [parlist]
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
func _parseParamList(lexer *Lexer) (names []string, isVararg bool) {
	switch lexer.LookAhead() {
	case TOKEN_SEP_RPAREN:
		return nil, false
	case TOKEN_VARARG:
		lexer.NextToken()
		return nil, true
	}

	_, name := lexer.NextIdentifier()
	names = append(names, name)
	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()
		if lexer.LookAhead() == TOKEN_IDENTIFIER {
			_, name := lexer.NextIdentifier()
			names = append(names, name)
		} else {
			lexer.NextTokenOfKind(TOKEN_VARARG)
			isVararg = true
			break
		}
	}
	return
}

// tableconstructor ::= ‘{’ [fieldlist] ‘}’
func parseTableConstructorExpression(lexer *Lexer) *TableConstructorExpression {
	line := lexer.CurrentLine()
	lexer.NextTokenOfKind(TOKEN_SEP_LCURLY)    // {
	keyExps, valExps := _parseFieldList(lexer) // [fieldlist]
	lexer.NextTokenOfKind(TOKEN_SEP_RCURLY)    // }
	lastLine := lexer.CurrentLine()
	return &TableConstructorExpression{line, lastLine, keyExps, valExps}
}

// fieldlist ::= field {fieldsep field} [fieldsep]
func _parseFieldList(lexer *Lexer) (ks, vs []Expression) {
	if lexer.LookAhead() != TOKEN_SEP_RCURLY {
		k, v := _parseField(lexer)
		ks = append(ks, k)
		vs = append(vs, v)

		for _isFieldSep(lexer.LookAhead()) {
			lexer.NextToken()
			if lexer.LookAhead() != TOKEN_SEP_RCURLY {
				k, v := _parseField(lexer)
				ks = append(ks, k)
				vs = append(vs, v)
			} else {
				break
			}
		}
	}
	return
}

// fieldsep ::= ‘,’ | ‘;’
func _isFieldSep(tokenKind int) bool {
	return tokenKind == TOKEN_SEP_COMMA || tokenKind == TOKEN_SEP_SEMI
}

// field ::= ‘[’ exp ‘]’ ‘=’ exp | Name ‘=’ exp | exp
func _parseField(lexer *Lexer) (k, v Expression) {
	if lexer.LookAhead() == TOKEN_SEP_LBRACK {
		lexer.NextToken()                       // [
		k = parseExpression(lexer)                     // exp
		lexer.NextTokenOfKind(TOKEN_SEP_RBRACK) // ]
		lexer.NextTokenOfKind(TOKEN_OP_ASSIGN)  // =
		v = parseExpression(lexer)                     // exp
		return
	}

	exp := parseExpression(lexer)
	if nameExp, ok := exp.(*NameExpression); ok {
		if lexer.LookAhead() == TOKEN_OP_ASSIGN {
			// Name ‘=’ exp => ‘[’ LiteralString ‘]’ = exp
			lexer.NextToken()
			k = &StringExpression{nameExp.Line, nameExp.Name}
			v = parseExpression(lexer)
			return
		}
	}

	return nil, exp
}
