package Item

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"sg_interpreter/src/sg/ast"
	"strings"
)

type ItemType string

type Item interface {
	Type() ItemType
	Output() string
}

const (
	NULL_ITEM    = "NULL"
	ERROR_ITEM   = "ERROR"
	INTEGER_ITEM = "INTEGER"
	BOOLEAN_ITEM = "BOOLEAN"
	STRING_ITEM  = "STRING"

	FUNCTION_ITEM     = "FUNCTION"
	RETURN_VALUE_ITEM = "RETURN_VALUE"

	BUILTIN_ITEM = "BUILTIN"

	ARRAY_ITEM = "ARRAY"
	HASH_ITEM  = "HASH"
)

type HashKey struct {
	Type  ItemType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type Integer struct {
	Value int64
}

func (integer *Integer) Type() ItemType {
	return INTEGER_ITEM
}
func (integer *Integer) Output() string {
	return fmt.Sprintf("%d", integer.Value)
}
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
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
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
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

type String struct {
	Value string
}

func (s *String) Type() ItemType { return STRING_ITEM }
func (s *String) Output() string { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type BuiltinFunction func(args ...Item) Item

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ItemType { return BUILTIN_ITEM }
func (b *Builtin) Output() string { return "builtin function" }

type Array struct {
	Elements []Item
	Capacity int64
	Len      int64
}

func (ao *Array) Type() ItemType { return ARRAY_ITEM }
func (ao *Array) Output() string {
	var out bytes.Buffer

	elements := []string{}

	for i := int64(0); i < ao.Len; i++ {
		elements = append(elements, ao.Elements[i].Output())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Item
	Value Item
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ItemType { return HASH_ITEM }
func (h *Hash) Output() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Output(), pair.Value.Output()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
