package parser

import (
	. "golua/compiler/ast"
	. "golua/compiler/lexer"
)

var _statEmpty = &EmptyStatement{}

func parseStatement(lexer *Lexer) Statement {
	switch lexer.LookAhead() {
	case TOKEN_SEP_SEMI:
		return parseEmptyStat(lexer)
	case TOKEN_KW_BREAK:
		return parseBreakStat(lexer)
	case TOKEN_SEP_LABEL:
		return parseLabelStat(lexer)
	case TOKEN_KW_GOTO:
		return parseGotoStat(lexer)
	case TOKEN_KW_DO:
		return parseDoStat(lexer)
	case TOKEN_KW_WHILE:
		return parseWhileStat(lexer)
	case TOKEN_KW_REPEAT:
		return parseRepeatStat(lexer)
	case TOKEN_KW_IF:
		return parseIfStat(lexer)
	case TOKEN_KW_FOR:
		return parseForStat(lexer)
	case TOKEN_KW_FUNCTION:
		return parseFuncDefStat(lexer)
	case TOKEN_KW_LOCAL:
		return parseLocalAssignOrFuncDefStat(lexer)
	default:
		return parseAssignOrFuncCallStat(lexer)
	}
}

func parseEmptyStat(lexer *Lexer) *EmptyStatement {
	lexer.NextTokenOfKind(TOKEN_SEP_SEMI)
	return _statEmpty
}

// break
func parseBreakStat(lexer *Lexer) *BreakStatement {
	lexer.NextTokenOfKind(TOKEN_KW_BREAK)
	return &BreakStatement{lexer.CurrentLine()}
}

// ‘::’ Name ‘::’
func parseLabelStat(lexer *Lexer) *LabelStatement {
	lexer.NextTokenOfKind(TOKEN_SEP_LABEL) // ::
	_, name := lexer.NextIdentifier()      // name
	lexer.NextTokenOfKind(TOKEN_SEP_LABEL) // ::
	return &LabelStatement{name}
}

// goto Name
func parseGotoStat(lexer *Lexer) *GotoStatement {
	lexer.NextTokenOfKind(TOKEN_KW_GOTO) // goto
	_, name := lexer.NextIdentifier()    // name
	return &GotoStatement{name}
}

// do block end
func parseDoStat(lexer *Lexer) *DoStatement {
	lexer.NextTokenOfKind(TOKEN_KW_DO)  // do
	block := parseBlock(lexer)          // block
	lexer.NextTokenOfKind(TOKEN_KW_END) // end
	return &DoStatement{block}
}

// while exp do block end
func parseWhileStat(lexer *Lexer) *WhileStatement {
	lexer.NextTokenOfKind(TOKEN_KW_WHILE) // while
	exp := parseExpression(lexer)                // exp
	lexer.NextTokenOfKind(TOKEN_KW_DO)    // do
	block := parseBlock(lexer)            // block
	lexer.NextTokenOfKind(TOKEN_KW_END)   // end
	return &WhileStatement{exp, block}
}

// repeat block until exp
func parseRepeatStat(lexer *Lexer) *RepeatStatement {
	lexer.NextTokenOfKind(TOKEN_KW_REPEAT) // repeat
	block := parseBlock(lexer)             // block
	lexer.NextTokenOfKind(TOKEN_KW_UNTIL)  // until
	exp := parseExpression(lexer)                 // exp
	return &RepeatStatement{exp, block}
}

// if exp then block {elseif exp then block} [else block] end
func parseIfStat(lexer *Lexer) *IfStatement {
	exps := make([]Expression, 0, 4)
	blocks := make([]*Block, 0, 4)

	lexer.NextTokenOfKind(TOKEN_KW_IF)         // if
	exps = append(exps, parseExpression(lexer))       // exp
	lexer.NextTokenOfKind(TOKEN_KW_THEN)       // then
	blocks = append(blocks, parseBlock(lexer)) // block

	for lexer.LookAhead() == TOKEN_KW_ELSEIF {
		lexer.NextToken()                          // elseif
		exps = append(exps, parseExpression(lexer))       // exp
		lexer.NextTokenOfKind(TOKEN_KW_THEN)       // then
		blocks = append(blocks, parseBlock(lexer)) // block
	}

	// else block => elseif true then block
	if lexer.LookAhead() == TOKEN_KW_ELSE {
		lexer.NextToken()                           // else
		exps = append(exps, &TrueExpression{lexer.CurrentLine()}) //
		blocks = append(blocks, parseBlock(lexer))  // block
	}

	lexer.NextTokenOfKind(TOKEN_KW_END) // end
	return &IfStatement{exps, blocks}
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
// for namelist in explist do block end
func parseForStat(lexer *Lexer) Statement {
	lineOfFor, _ := lexer.NextTokenOfKind(TOKEN_KW_FOR)
	_, name := lexer.NextIdentifier()
	if lexer.LookAhead() == TOKEN_OP_ASSIGN {
		return _finishForNumStat(lexer, lineOfFor, name)
	} else {
		return _finishForInStat(lexer, name)
	}
}

func _finishForNumStat(lexer *Lexer, lineOfFor int, varName string) *ForNumStatement {
	lexer.NextTokenOfKind(TOKEN_OP_ASSIGN) // for name =
	initExp := parseExpression(lexer)             // exp
	lexer.NextTokenOfKind(TOKEN_SEP_COMMA) // ,
	limitExp := parseExpression(lexer)            // exp

	var stepExp Expression
	if lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()         // ,
		stepExp = parseExpression(lexer) // exp
	} else {
		stepExp = &IntegerExpression{lexer.CurrentLine(), 1}
	}

	lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO) // do
	block := parseBlock(lexer)                        // block
	lexer.NextTokenOfKind(TOKEN_KW_END)               // end

	return &ForNumStatement{lineOfFor, lineOfDo,
		varName, initExp, limitExp, stepExp, block}
}

func _finishForInStat(lexer *Lexer, name0 string) *ForInStatement {
	nameList := _finishNameList(lexer, name0)         // for namelist
	lexer.NextTokenOfKind(TOKEN_KW_IN)                // in
	expList := parseExpressionList(lexer)                    // explist
	lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO) // do
	block := parseBlock(lexer)                        // block
	lexer.NextTokenOfKind(TOKEN_KW_END)               // end
	return &ForInStatement{lineOfDo, nameList, expList, block}
}

func _finishNameList(lexer *Lexer, name0 string) []string {
	names := []string{name0}
	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()                 // ,
		_, name := lexer.NextIdentifier() // Name
		names = append(names, name)
	}
	return names
}

func parseLocalAssignOrFuncDefStat(lexer *Lexer) Statement {
	lexer.NextTokenOfKind(TOKEN_KW_LOCAL)
	if lexer.LookAhead() == TOKEN_KW_FUNCTION {
		return _finishLocalFuncDefStat(lexer)
	} else {
		return _finishLocalVarDeclStat(lexer)
	}
}

func _finishLocalFuncDefStat(lexer *Lexer) *LocalFunctionDefineStatement {
	lexer.NextTokenOfKind(TOKEN_KW_FUNCTION) // local function
	_, name := lexer.NextIdentifier()        // name
	fdExp := parseFuncDefExp(lexer)          // funcbody
	return &LocalFunctionDefineStatement{Name: name, Expression: fdExp}
}

// local namelist [‘=’ explist]
func _finishLocalVarDeclStat(lexer *Lexer) *LocalVariableDeclareStatement {
	_, name0 := lexer.NextIdentifier()        // local Name
	nameList := _finishNameList(lexer, name0) // { , Name }
	var expList []Expression = nil
	if lexer.LookAhead() == TOKEN_OP_ASSIGN {
		lexer.NextToken()             // ==
		expList = parseExpressionList(lexer) // explist
	}
	lastLine := lexer.CurrentLine()
	return &LocalVariableDeclareStatement{lastLine, nameList, expList}
}

func parseAssignOrFuncCallStat(lexer *Lexer) Statement{
	prefixExp := parsePrefixExp(lexer)
	if fc, ok := prefixExp.(*FunctionCallExpression); ok {
		return fc
	} else {
		return parseAssignStat(lexer, prefixExp)
	}
}

func parseFuncDefStat(lexer *Lexer) *AssignStatement {
	lexer.NextTokenOfKind(TOKEN_KW_FUNCTION) // function
	fnExp, hasColon := _parseFuncName(lexer) // funcname
	fdExp := parseFuncDefExp(lexer)          // funcbody
	if hasColon {                            // insert self
		fdExp.ParamList = append(fdExp.ParamList, "")
		copy(fdExp.ParamList[1:], fdExp.ParamList)
		fdExp.ParamList[0] = "self"
	}

	return &AssignStatement{
		LastLine: fdExp.Line,
		VariableList:  []Expression{fnExp},
		ExpressionList:  []Expression{fdExp},
	}
}

func parseAssignStat(lexer *Lexer, var0 Expression) *AssignStatement {
	varList := _finishVarList(lexer, var0) // varlist
	lexer.NextTokenOfKind(TOKEN_OP_ASSIGN) // =
	expList := parseExpressionList(lexer)         // explist
	lastLine := lexer.CurrentLine()
	return &AssignStatement{lastLine, varList, expList}
}

func _finishVarList(lexer *Lexer, var0 Expression) []Expression {
	vars := []Expression{_checkVar(lexer, var0)}      // var
	for lexer.LookAhead() == TOKEN_SEP_COMMA { // {
		lexer.NextToken()                          // ,
		exp := parsePrefixExp(lexer)               // var
		vars = append(vars, _checkVar(lexer, exp)) //
	} // }
	return vars
}

func _checkVar(lexer *Lexer, exp Expression) Expression {
	switch exp.(type) {
	case *NameExpression, *TableAccessExpression:
		return exp
	}
	lexer.NextTokenOfKind(-1) // trigger error
	panic("unreachable!")
}

func _parseFuncName(lexer *Lexer) (exp Expression, hasColon bool) {
	line, name := lexer.NextIdentifier()
	exp = &NameExpression{line, name}

	for lexer.LookAhead() == TOKEN_SEP_DOT {
		lexer.NextToken()
		line, name := lexer.NextIdentifier()
		idx := &StringExpression{line, name}
		exp = &TableAccessExpression{line, exp, idx}
	}
	if lexer.LookAhead() == TOKEN_SEP_COLON {
		lexer.NextToken()
		line, name := lexer.NextIdentifier()
		idx := &StringExpression{line, name}
		exp = &TableAccessExpression{line, exp, idx}
		hasColon = true
	}

	return
}