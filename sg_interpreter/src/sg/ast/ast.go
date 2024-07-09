package ast

import (
	"bytes"
	"sg_interpreter/src/sg/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

type Expression interface {
	Node
	expressionNode()
}
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

func (program *Program) TokenLiteral() string {
	if program.Statements != nil && len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (program *Program) String() string {
	var output bytes.Buffer
	for _, s := range program.Statements {
		output.WriteString(s.String())
	}
	return output.String()
}

type LetStatement struct {
	Token token.Token
	Id    *Identifier
	Val   Expression
}

func (letStatement *LetStatement) statementNode()       {}
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }
func (letStatement *LetStatement) String() string {
	var output bytes.Buffer
	output.WriteString(letStatement.TokenLiteral() + " ")
	output.WriteString(letStatement.Id.String())
	output.WriteString(" = ")
	if letStatement.Val != nil {
		output.WriteString(letStatement.Val.String())
	}
	output.WriteString(";")
	return output.String()
}

type ReturnStatement struct {
	Token    token.Token
	RetValue Expression
}

func (returnStatement *ReturnStatement) statementNode()       {}
func (returnStatement *ReturnStatement) TokenLiteral() string { return returnStatement.Token.Literal }
func (returnStatement *ReturnStatement) String() string {
	var output bytes.Buffer
	output.WriteString(returnStatement.TokenLiteral() + " ")
	if returnStatement.RetValue != nil {
		output.WriteString(returnStatement.RetValue.String())
	}
	output.WriteString(";")
	return output.String()
}

type ExpressionStatement struct {
	Token token.Token
	Expr  Expression
}

func (expressionStatement *ExpressionStatement) statementNode() {}
func (expressionStatement *ExpressionStatement) TokenLiteral() string {
	return expressionStatement.Token.Literal
}
func (expressionStatement *ExpressionStatement) String() string {
	if expressionStatement.Expr != nil {
		return expressionStatement.Expr.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (blockStatement *BlockStatement) statementNode()       {}
func (blockStatement *BlockStatement) TokenLiteral() string { return blockStatement.Token.Literal }
func (blockStatement *BlockStatement) String() string {
	var output bytes.Buffer
	for _, s := range blockStatement.Statements {
		output.WriteString(s.String())
	}
	return output.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (boolean *Boolean) expressionNode()      {}
func (boolean *Boolean) TokenLiteral() string { return boolean.Token.Literal }
func (boolean *Boolean) String() string       { return boolean.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (integerLiteral *IntegerLiteral) expressionNode()      {}
func (integerLiteral *IntegerLiteral) TokenLiteral() string { return integerLiteral.Token.Literal }
func (integerLiteral *IntegerLiteral) String() string       { return integerLiteral.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (expression *PrefixExpression) expressionNode()      {}
func (expression *PrefixExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *PrefixExpression) String() string {
	var output bytes.Buffer

	output.WriteString("(")
	output.WriteString(expression.Operator)
	output.WriteString(expression.Right.String())
	output.WriteString(")")

	return output.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (expression *InfixExpression) expressionNode()      {}
func (expression *InfixExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *InfixExpression) String() string {
	var output bytes.Buffer
	output.WriteString("(")
	output.WriteString(expression.Left.String())
	output.WriteString(" " + expression.Operator + " ")
	output.WriteString(expression.Right.String())
	output.WriteString(")")
	return output.String()
}

type IfExpression struct {
	Token token.Token
	Cond  Expression
	Cons  *BlockStatement
	Alt   *BlockStatement
}

func (ifExpression *IfExpression) expressionNode()      {}
func (ifExpression *IfExpression) TokenLiteral() string { return ifExpression.Token.Literal }
func (ifExpression *IfExpression) String() string {
	var output bytes.Buffer
	output.WriteString("if")
	output.WriteString(ifExpression.Cons.String())
	output.WriteString(" ")
	output.WriteString(ifExpression.Cons.String())
	output.WriteString(" ")
	output.WriteString(ifExpression.Alt.String())
	return output.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (functionLiteral *FunctionLiteral) expressionNode()      {}
func (functionLiteral *FunctionLiteral) TokenLiteral() string { return functionLiteral.Token.Literal }
func (functionLiteral *FunctionLiteral) String() string {
	var output bytes.Buffer
	params := []string{}
	for _, p := range functionLiteral.Parameters {
		params = append(params, p.String())
	}
	output.WriteString(functionLiteral.TokenLiteral())
	output.WriteString("(")
	output.WriteString(strings.Join(params, ", "))
	output.WriteString(") ")
	output.WriteString(functionLiteral.Body.String())
	return output.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (callExpression *CallExpression) expressionNode()      {}
func (callExpression *CallExpression) TokenLiteral() string { return callExpression.Token.Literal }
func (callExpression *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range callExpression.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(callExpression.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}