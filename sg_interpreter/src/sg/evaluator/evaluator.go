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
		return EvalProgram(node, scope)
	case *ast.BlockStatement:
		return EvalBlockStatement(node, scope)
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
	case *ast.Boolean:
		return boolToBoolean(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, scope)
		if isError(right) {
			return right
		}
		return EvalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, scope)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, scope)
		if isError(right) {
			return right
		}
		return EvalInfixExpression(left, node.Operator, right)
	case *ast.IfExpression:
		return EvalIfExpression(node, scope)
	case *ast.Identifier:
		return EvalIdentifier(node, scope)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Item.Function{Parameters: params, Body: body, Scope: scope}
	case *ast.CallExpression:
		function := Eval(node.Function, scope)
		if isError(function) {
			return function
		}
		args := EvalExpression(node.Arguments, scope)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return ApplyFunction(function, args)
	}
	return nil
}

func EvalProgram(program *ast.Program, scope *Item.Scope) Item.Item {
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

func EvalBlockStatement(block *ast.BlockStatement, scope *Item.Scope) Item.Item {
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

func EvalPrefixExpression(operator string, expression Item.Item) Item.Item {
	switch operator {
	case "!":
		return EvalEXC(expression)
	case "-":
		return EvalMINUS(expression)
	default:
		return newError("unknown operator: %s%s", operator, expression.Type())
	}
}

func EvalInfixExpression(left Item.Item, op string, right Item.Item) Item.Item {
	switch {
	case left.Type() == Item.INTEGER_ITEM && right.Type() == Item.INTEGER_ITEM:
		return EvalIntegerInfixExpr(left, op, right)
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

func EvalEXC(expression Item.Item) Item.Item {
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

func EvalMINUS(expression Item.Item) Item.Item {
	if expression.Type() != Item.INTEGER_ITEM {
		return newError("unknown operator: -%s", expression.Type())
	}

	val := expression.(*Item.Integer).Value
	return &Item.Integer{Value: -val}
}

func EvalIntegerInfixExpr(left Item.Item, op string, right Item.Item) Item.Item {
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

func boolToBoolean(b bool) *Item.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func TrueLike(item Item.Item) bool {
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
func EvalIfExpression(is *ast.IfExpression, scope *Item.Scope) Item.Item {
	cond := Eval(is.Cond, scope)
	if isError(cond) {
		return cond
	}
	if TrueLike(cond) {
		return Eval(is.Cons, scope)
	} else if is.Alt != nil {
		return Eval(is.Alt, scope)
	} else {
		return NULL
	}
}

func EvalIdentifier(identifier *ast.Identifier, scope *Item.Scope) Item.Item {
	val, ok := scope.Get(identifier.Value)
	if !ok {
		return newError("identifier not found: " + identifier.Value)
	}
	return val
}

func EvalExpression(expressions []ast.Expression, scope *Item.Scope) []Item.Item {
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
func ApplyFunction(_function Item.Item, args []Item.Item) Item.Item {
	function, ok := _function.(*Item.Function)
	if !ok {
		return newError("not a function: %s", _function.Type())
	}

	extendedScope := extendedScope(function, args)
	eval := Eval(function.Body, extendedScope)
	return UnrwapReturnValue(eval)
}
func extendedScope(function *Item.Function, args []Item.Item) *Item.Scope {
	scope := Item.NewEnclosedScope(function.Scope)
	for paramIdx, param := range function.Parameters {
		scope.Set(param.Value, args[paramIdx])
	}
	return scope
}
func UnrwapReturnValue(item Item.Item) Item.Item {
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
