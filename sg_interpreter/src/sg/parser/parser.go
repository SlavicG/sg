package parser

import (
	"fmt"
	"sg_interpreter/src/sg/ast"
	"sg_interpreter/src/sg/lexer"
	"sg_interpreter/src/sg/token"
	"strconv"
)

const (
	LOWEST        = 1
	EQUALS        = 2
	LESSORGREATER = 3
	SUM           = 4
	PRODUCT       = 5
	PREFIX        = 6
	CALL          = 7
)

var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT:     LESSORGREATER,
	token.GT:     LESSORGREATER,
	token.MINUS:  SUM,
	token.PLUS:   SUM,
	token.STAR:   PRODUCT,
	token.SLASH:  PRODUCT,
	token.LP:     CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.infixParseFns = make(map[token.TokenType]infixParseFn)

	//Registering the respective functions to tokens
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.EXC, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.LP, parser.parseGroupedExpressions)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)

	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.STAR, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.LP, parser.parseCallExpression)

	parser.nextToken()
	parser.nextToken()

	return parser
}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) CurTokenIsType(tokenType token.TokenType) bool {
	return tokenType == parser.curToken.Type
}
func (parser *Parser) PeekTokenIsType(tokenType token.TokenType) bool {
	return tokenType == parser.peekToken.Type
}
func (parser *Parser) ExpectPeek(tokenType token.TokenType) bool {
	if parser.PeekTokenIsType(tokenType) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(tokenType)
		return false
	}
}
func (parser *Parser) peekPrecedence() int {
	if pp, ok := precedences[parser.peekToken.Type]; ok {
		return pp
	}
	return LOWEST
}
func (parser *Parser) curPrecedence() int {
	if pp, ok := precedences[parser.curToken.Type]; ok {
		return pp
	}

	return LOWEST
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(tokenType token.TokenType) {
	message := fmt.Sprintf("Expected next token to be %s, got %s.", tokenType, parser.peekToken.Type)
	parser.errors = append(parser.errors, message)
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
}
func (parser *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: parser.curToken}

	value, err := strconv.ParseInt(parser.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Couldn't parse %q as an Integer", parser.curToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}
func (parser *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: parser.curToken, Value: parser.CurTokenIsType(token.TRUE)}
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !parser.CurTokenIsType(token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.nextToken()
	}
	return program
}
func (parser *Parser) parseStatement() ast.Statement {
	if parser.curToken.Type == token.IDENT && parser.peekToken.Type == token.SET {
		return parser.parseSetStatement()
	}
	switch parser.curToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}
func (parser *Parser) parseSetStatement() ast.Statement {
	statement := &ast.SetStatement{Token: parser.curToken}
	statement.Id = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	if !parser.ExpectPeek(token.SET) {
		return nil
	}
	parser.nextToken()
	statement.Val = parser.parseExpression(LOWEST)

	if parser.PeekTokenIsType(token.SEMICOL) {
		parser.nextToken()
	}

	return statement
}
func (parser *Parser) parseLetStatement() ast.Statement {
	statement := &ast.LetStatement{Token: parser.curToken}
	if !parser.ExpectPeek(token.IDENT) {
		return nil
	}
	statement.Id = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	if !parser.ExpectPeek(token.SET) {
		return nil
	}
	parser.nextToken()
	statement.Val = parser.parseExpression(LOWEST)

	if parser.PeekTokenIsType(token.SEMICOL) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() ast.Statement {
	statement := &ast.ReturnStatement{Token: parser.curToken}
	parser.nextToken()
	statement.RetValue = parser.parseExpression(LOWEST)
	if parser.PeekTokenIsType(token.SEMICOL) {
		parser.nextToken()
	}
	return statement
}

func (parser *Parser) parseExpressionStatement() ast.Statement {
	statement := &ast.ExpressionStatement{Token: parser.curToken}
	statement.Expr = parser.parseExpression(LOWEST)

	if parser.PeekTokenIsType(token.SEMICOL) {
		parser.nextToken()
	}
	return statement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFns[parser.curToken.Type]
	if prefix == nil {
		parser.noPrefixParseFnError(parser.curToken.Type)

		return nil
	}
	leftExpr := prefix()

	for !parser.PeekTokenIsType(token.SEMICOL) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFns[parser.peekToken.Type]
		if infix == nil {
			return leftExpr
		}
		parser.nextToken()
		leftExpr = infix(leftExpr)
	}
	return leftExpr
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
	}
	parser.nextToken()
	expression.Right = parser.parseExpression(PREFIX)
	return expression
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
		Left:     left,
	}
	precedence := parser.curPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)
	return expression
}

func (parser *Parser) parseGroupedExpressions() ast.Expression {
	parser.nextToken()
	expression := parser.parseExpression(LOWEST)
	if !parser.ExpectPeek(token.RP) {
		return nil
	}
	return expression
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.curToken}
	if !parser.ExpectPeek(token.LP) {
		return nil
	}
	parser.nextToken()
	expression.Cond = parser.parseExpression(LOWEST)
	if !parser.ExpectPeek(token.RP) {
		return nil
	}
	if !parser.ExpectPeek(token.LB) {
		return nil
	}
	expression.Cons = parser.parseBlockStatement()
	if parser.PeekTokenIsType(token.ELSE) {
		parser.nextToken()
		if !parser.ExpectPeek(token.LB) {
			return nil
		}
		expression.Alt = parser.parseBlockStatement()
	}
	return expression
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: parser.curToken}
	block.Statements = []ast.Statement{}

	parser.nextToken()

	for !parser.CurTokenIsType(token.RB) && !parser.CurTokenIsType(token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		parser.nextToken()
	}
	return block
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: parser.curToken}

	if !parser.ExpectPeek(token.LP) {
		return nil
	}
	literal.Parameters = parser.parseFunctionParameters()

	if !parser.ExpectPeek(token.LB) {
		return nil
	}
	literal.Body = parser.parseBlockStatement()
	return literal
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if parser.PeekTokenIsType(token.RP) {
		parser.nextToken()
		return identifiers
	}
	parser.nextToken()

	ident := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
	identifiers = append(identifiers, ident)

	for parser.PeekTokenIsType(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		ident := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !parser.ExpectPeek(token.RP) {
		return nil
	}
	return identifiers
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: parser.curToken, Function: function}
	expression.Arguments = parser.parseCallArguments()
	return expression
}

func (parser *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if parser.PeekTokenIsType(token.RP) {
		parser.nextToken()
		return args
	}
	parser.nextToken()
	arg := parser.parseExpression(LOWEST)
	args = append(args, arg)

	for parser.PeekTokenIsType(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		arg := parser.parseExpression(LOWEST)
		args = append(args, arg)
	}

	if !parser.ExpectPeek(token.RP) {
		return nil
	}
	return args
}
