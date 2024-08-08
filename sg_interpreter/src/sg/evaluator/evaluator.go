package evaluator

import (
	"fmt"
	"sg_interpreter/src/sg/Item"
	"sg_interpreter/src/sg/ast"
)

var (
	NULL  = &Item.Null{}
	TRUE  = &Item.Boolean{Value: true}
	FALSE = &Item.Boolean{Value: false}
)

func Eval(node ast.Node, scope *Item.Scope) Item.Item {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, scope)
	case *ast.BlockStatement:
		return evalBlockStatement(node, scope)
	case *ast.ExpressionStatement:
		return Eval(node.Expr, scope)
	case *ast.ReturnStatement:
		val := Eval(node.RetValue, scope)
		if isError(val) {
			return val
		}
		return &Item.ReturnValue{Value: val}
	case *ast.SetStatement:
		val := Eval(node.Val, scope)
		if isError(val) {
			return val
		}
		_, b := scope.Get(node.Id.Value)
		if !b {
			return newError("Variable %s not defined in current scope!", node.Id.Value)
		}
		scope.Set(node.Id.Value, val)
	case *ast.LetStatement:
		val := Eval(node.Val, scope)
		if isError(val) {
			return val
		}
		_, b := scope.Mp[node.Id.Value]
		if b {
			return newError("Variable %s already is defined in this function's scope!", node.Id.Value)
		}
		scope.Set(node.Id.Value, val)
	case *ast.IntegerLiteral:
		return &Item.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &Item.String{Value: node.Value}
	case *ast.Boolean:
		return boolToBoolean(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, scope)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, scope)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, scope)
		if isError(right) {
			return right
		}
		return evalInfixExpression(left, node.Operator, right)
	case *ast.IfExpression:
		return evalIfExpression(node, scope)
	case *ast.Identifier:
		return evalIdentifier(node, scope)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Item.Function{Parameters: params, Body: body, Scope: scope}
	case *ast.CallExpression:
		function := Eval(node.Function, scope)
		if isError(function) {
			return function
		}
		args := evalExpression(node.Arguments, scope)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.ArrayLiteral:
		elements := evalExpression(node.Elements, scope)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Item.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, scope)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, scope)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.MapLiteral:
		return evalMapLiteral(node, scope)

	}
	return nil
}

func evalProgram(program *ast.Program, scope *Item.Scope) Item.Item {
	var res Item.Item
	for _, statement := range program.Statements {
		res = Eval(statement, scope)
		switch result := res.(type) {
		case *Item.ReturnValue:
			return result.Value
		case *Item.Error:
			return result
		}
	}
	return res
}

func evalBlockStatement(block *ast.BlockStatement, scope *Item.Scope) Item.Item {
	var res Item.Item
	for _, statement := range block.Statements {
		res = Eval(statement, scope)
		if res != nil {
			if res.Type() == Item.RETURN_VALUE_ITEM || res.Type() == Item.ERROR_ITEM {
				return res
			}
		}
	}
	return res
}

func evalPrefixExpression(operator string, expression Item.Item) Item.Item {
	switch operator {
	case "!":
		return evalEXC(expression)
	case "-":
		return evalMINUS(expression)
	default:
		return newError("unknown operator: %s%s", operator, expression.Type())
	}
}

func evalInfixExpression(left Item.Item, op string, right Item.Item) Item.Item {
	switch {
	case left.Type() == Item.INTEGER_ITEM && right.Type() == Item.INTEGER_ITEM:
		return evalIntegerInfixExpr(left, op, right)
	case left.Type() == Item.STRING_ITEM && right.Type() == Item.STRING_ITEM:
		return evalStringInfixExpression(op, left, right)
	case op == "==":
		return boolToBoolean(left == right)
	case op == "!=":
		return boolToBoolean(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), op, right.Type())
	}
}

func evalEXC(expression Item.Item) Item.Item {
	switch expression {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMINUS(expression Item.Item) Item.Item {
	if expression.Type() != Item.INTEGER_ITEM {
		return newError("unknown operator: -%s", expression.Type())
	}

	val := expression.(*Item.Integer).Value
	return &Item.Integer{Value: -val}
}

func evalIntegerInfixExpr(left Item.Item, op string, right Item.Item) Item.Item {
	leftVal := left.(*Item.Integer).Value
	rightVal := right.(*Item.Integer).Value
	switch op {
	case "-":
		return &Item.Integer{Value: leftVal - rightVal}
	case "+":
		return &Item.Integer{Value: leftVal + rightVal}
	case "/":
		return &Item.Integer{Value: leftVal / rightVal}
	case "*":
		return &Item.Integer{Value: leftVal * rightVal}
	case "==":
		return boolToBoolean(leftVal == rightVal)
	case "!=":
		return boolToBoolean(leftVal != rightVal)
	case "<":
		return boolToBoolean(leftVal < rightVal)
	case ">":
		return boolToBoolean(leftVal > rightVal)
	}
	return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
}
func evalStringInfixExpression(
	operator string,
	left, right Item.Item,
) Item.Item {
	if operator != "+" {
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

	leftVal := left.(*Item.String).Value
	rightVal := right.(*Item.String).Value
	return &Item.String{Value: leftVal + rightVal}
}
func boolToBoolean(b bool) *Item.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func trueLike(item Item.Item) bool {
	switch item {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
func evalIfExpression(is *ast.IfExpression, scope *Item.Scope) Item.Item {
	cond := Eval(is.Cond, scope)
	if isError(cond) {
		return cond
	}
	if trueLike(cond) {
		return Eval(is.Cons, scope)
	} else if is.Alt != nil {
		return Eval(is.Alt, scope)
	} else {
		return NULL
	}
}

func evalIdentifier(
	node *ast.Identifier,
	env *Item.Scope,
) Item.Item {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpression(expressions []ast.Expression, scope *Item.Scope) []Item.Item {
	var res []Item.Item
	for _, expression := range expressions {
		eval := Eval(expression, scope)
		if isError(eval) {
			return []Item.Item{eval}
		}
		res = append(res, eval)
	}
	return res
}
func applyFunction(fn Item.Item, args []Item.Item) Item.Item {
	switch fn := fn.(type) {

	case *Item.Function:
		extendedEnv := extendedScope(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *Item.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}
func extendedScope(function *Item.Function, args []Item.Item) *Item.Scope {
	scope := Item.NewEnclosedScope(function.Scope)
	for paramIdx, param := range function.Parameters {
		scope.Set(param.Value, args[paramIdx])
	}
	return scope
}
func unwrapReturnValue(item Item.Item) Item.Item {
	if returnValue, ok := item.(*Item.ReturnValue); ok {
		return returnValue.Value
	}
	return item
}

func newError(format string, a ...interface{}) *Item.Error {
	return &Item.Error{Message: fmt.Sprintf(format, a...)}
}
func isError(item Item.Item) bool {
	if item != nil {
		return item.Type() == Item.ERROR_ITEM
	}
	return false
}

func evalArrayIndexExpression(array, index Item.Item) Item.Item {
	arrayObject := array.(*Item.Array)
	idx := index.(*Item.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalMapLiteral(
	node *ast.MapLiteral,
	env *Item.Scope,
) Item.Item {
	pairs := make(map[Item.HashKey]Item.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(Item.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = Item.HashPair{Key: key, Value: value}
	}

	return &Item.Hash{Pairs: pairs}
}

func evalMapIndexExpression(hash, index Item.Item) Item.Item {
	hashObject := hash.(*Item.Hash)

	key, ok := index.(Item.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
func evalIndexExpression(left, index Item.Item) Item.Item {
	switch {
	case left.Type() == Item.ARRAY_ITEM && index.Type() == Item.INTEGER_ITEM:
		return evalArrayIndexExpression(left, index)
	case left.Type() == Item.HASH_ITEM:
		return evalMapIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}
