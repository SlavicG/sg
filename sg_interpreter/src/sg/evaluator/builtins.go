package evaluator

import (
	"fmt"
	"sg_interpreter/src/sg/Item"
)

var builtins = map[string]*Item.Builtin{
	"len": {Fn: func(args ...Item.Item) Item.Item {
		if len(args) != 1 {
			return newError("Wrong number of arguments! Expected = 1. Received=%d",
				len(args))
		}

		switch arg := args[0].(type) {
		case *Item.Array:
			return &Item.Integer{Value: int64(len(arg.Elements))}
		case *Item.String:
			return &Item.Integer{Value: int64(len(arg.Value))}
		default:
			return newError("Argument `len` not supported. Received %s",
				args[0].Type())
		}
	},
	},
	"puts": {
		Fn: func(args ...Item.Item) Item.Item {
			res := ""
			first := true
			for _, arg := range args {
				if first == false {
					res += " "
				}
				res += arg.Output()
				first = false
			}
			fmt.Println(res)
			return NULL
		},
	},
	"first": {
		Fn: func(args ...Item.Item) Item.Item {
			if len(args) != 1 {
				return newError("Wrong number of arguments! Expected=1. Received=%d",
					len(args))
			}
			if args[0].Type() != Item.ARRAY_ITEM {
				return newError("Argument to `first` must be ARRAY. Received %s",
					args[0].Type())
			}

			arr := args[0].(*Item.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": {
		Fn: func(args ...Item.Item) Item.Item {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != Item.ARRAY_ITEM {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*Item.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			} else {
				return NULL
			}
		},
	},
	"push": {
		Fn: func(args ...Item.Item) Item.Item {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != Item.ARRAY_ITEM {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*Item.Array)
			length := len(arr.Elements)
			newElements := make([]Item.Item, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &Item.Array{Elements: newElements}
		},
	},
}
