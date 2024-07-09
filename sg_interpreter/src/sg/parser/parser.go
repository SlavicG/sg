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

	//TODO: REGISTER FUNCTIONS FOR RESPECTIVE TOKENS
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
	if parser.CurTokenIsType(token.LET) {
		return parser.parseLetStatement()
	} else if parser.CurTokenIsType(token.RETURN) {
		return parser.parseReturnStatement()
	} else {
		return parser.parseExpressionStatement()
	}
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
		//TODO: NoPrefixParseError
		return nil
	}
	leftExpr := prefix()

	for !parser.PeekTokenIsType(token.SEMICOL) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFns[parser.curToken.Type]
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

//TODO: FUNCTIONS AND IF EXPRESSIONS AND POLISHING ABOVE STUFF
