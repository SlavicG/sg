package Item

import (
	"bytes"
	"fmt"
	"sg_interpreter/src/sg/ast"
	"strings"
)

type ItemType string

type Item interface {
	Type() ItemType
	Output() string
}

const (
	NULL_ITEM         = "NULL"
	ERROR_ITEM        = "ERROR"
	INTEGER_ITEM      = "INTEGER"
	BOOLEAN_ITEM      = "BOOLEAN"
	FUNCTION_ITEM     = "FUNCTION"
	RETURN_VALUE_ITEM = "RETURN_VALUE"
)

type Integer struct {
	Value int64
}

func (integer *Integer) Type() ItemType {
	return INTEGER_ITEM
}
func (integer *Integer) Output() string {
	return fmt.Sprintf("%d", integer.Value)
}

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() ItemType {
	return BOOLEAN_ITEM
}
func (boolean *Boolean) Output() string {
	return fmt.Sprintf("%t", boolean.Value)
}

type Null struct{}

func (null *Null) Type() ItemType {
	return NULL_ITEM
}
func (null *Null) Output() string {
	return "null"
}

type ReturnValue struct {
	Value Item
}

func (returnValue *ReturnValue) Output() string {
	return returnValue.Value.Output()
}

func (returnValue *ReturnValue) Type() ItemType {
	return RETURN_VALUE_ITEM
}

type Error struct {
	Message string
}

func (error *Error) Type() ItemType {
	return ERROR_ITEM
}
func (error *Error) Output() string {
	return "ERROR: " + error.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Scope      *Scope
}

func (function *Function) Type() ItemType {
	return FUNCTION_ITEM
}
func (function *Function) Output() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range function.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fun")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(function.Body.String())
	out.WriteString("\n}")
	return out.String()
}
