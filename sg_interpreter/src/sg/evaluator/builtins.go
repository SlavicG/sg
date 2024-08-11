package evaluator

import (
	"fmt"
	"math/rand"
	"sg_interpreter/src/sg/Item"
	"time"
)

var builtins = map[string]*Item.Builtin{
	"len": {Fn: func(args ...Item.Item) Item.Item {
		if len(args) != 1 {
			return newError("Wrong number of arguments! Expected = 1. Received=%d",
				len(args))
		}
		switch arg := args[0].(type) {
		case *Item.Array:
			return &Item.Integer{Value: arg.Len}
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
				return newError("Wrong number of arguments! Expected=1. Received=%d", len(args))
			}
			if args[0].Type() != Item.ARRAY_ITEM {
				return newError("Argument to `last` must be ARRAY. Received %s", args[0].Type())
			}
			arr := args[0].(*Item.Array)
			if arr.Len > 0 {
				return arr.Elements[arr.Len-1]
			} else {
				return NULL
			}
		},
	},
	"push": {
		Fn: func(args ...Item.Item) Item.Item {
			if len(args) != 2 {
				return newError("Wrong number of arguments! Expected=2. Received=%d.", len(args))
			}
			if args[0].Type() != Item.ARRAY_ITEM {
				return newError("Argument to `push` must be ARRAY. Expected %s", args[0].Type())
			}
			arr := args[0].(*Item.Array)
			if arr.Len < arr.Capacity {
				arr.Elements[arr.Len] = args[1]
				arr.Len++
			} else {
				if arr.Capacity == 0 {
					arr.Capacity = 1
				} else {
					arr.Capacity *= 2
				}
				newElements := make([]Item.Item, arr.Capacity)
				copy(newElements, arr.Elements)
				newElements[arr.Len] = args[1]
				arr.Len++
				arr.Elements = newElements
			}
			return arr
		},
	},
	"set": {
		Fn: func(args ...Item.Item) Item.Item {
			if len(args) != 3 {
				return newError("Wrong number of arguments! Expected=3. Received=%d.", len(args))
			}
			if args[0].Type() != Item.ARRAY_ITEM {
				return newError("Argument to `set` must be ARRAY. Expected %s", args[0].Type())
			}
			arr := args[0].(*Item.Array)
			index, ok := args[1].(*Item.Integer)

			if !ok {
				return newError("Argument to `set` must be an INTEGER. Received %s", args[1].Type())
			}
			idx := index.Value
			if idx < 0 || idx >= arr.Len {
				return newError("Index Argument is out of bounds!")
			} else {
				arr.Elements[idx] = args[2]
			}
			return arr
		},
	},
	"get": {
		Fn: func(args ...Item.Item) Item.Item {
			if len(args) != 2 {
				return newError("Wrong number of arguments! Expected=2. Received=%d.", len(args))
			}
			if args[0].Type() != Item.STRING_ITEM {
				return newError("Argument to `get` must be STRING. Expected %s", args[0].Type())
			}
			s := args[0].(*Item.String)
			index, ok := args[1].(*Item.Integer)

			if !ok {
				return newError("Argument to `get` must be an INTEGER. Received %s", args[1].Type())
			}
			idx := index.Value
			if idx < 0 || idx >= int64(len(s.Value)) {
				return newError("Index Argument is out of bounds!")
			} else {
				return &Item.String{Value: string(s.Value[idx])} // Correctly slicing the string
			}
		},
	},
	"shuffle": {Fn: func(args ...Item.Item) Item.Item {
		if len(args) != 1 {
			return newError("Wrong number of arguments! Expected = 1. Received=%d",
				len(args))
		}
		arr := args[0].(*Item.Array)
		rand.Seed(time.Now().UnixNano())
		for i := arr.Len - 1; i > 0; i-- {
			j := rand.Intn(int(i + 1))
			arr.Elements[i], arr.Elements[j] = arr.Elements[j], arr.Elements[i]
		}
		return arr
	},
	},
	"reverse": {Fn: func(args ...Item.Item) Item.Item {
		if len(args) != 1 {
			return newError("Wrong number of arguments! Expected = 1. Received=%d",
				len(args))
		}
		arr := args[0].(*Item.Array)
		for i := 0; i < int(arr.Len/2); i++ {
			arr.Elements[i], arr.Elements[int(arr.Len)-1-i] = arr.Elements[int(arr.Len)-1-i], arr.Elements[i]
		}
		return arr
	},
	},

	"sort": {Fn: func(args ...Item.Item) Item.Item {
		if len(args) != 1 {
			return newError("Wrong number of arguments! Expected = 1. Received=%d",
				len(args))
		}
		arr := args[0].(*Item.Array)
		quicksort(arr, 0, int(arr.Len-1))
		return arr
	},
	},
}

func quicksort(arr *Item.Array, l int, r int) *Item.Error {
	if r-l+1 <= 1 {
		return nil
	}
	if arr.Elements[0].Type() != Item.INTEGER_ITEM {
		return newError("We can only sort Integer Arrays!")
	}

	pivot := r
	pos := l

	for i := l; i <= r; i++ {
		if arr.Elements[i].(*Item.Integer).Value < arr.Elements[pivot].(*Item.Integer).Value {
			swap(arr.Elements, i, pos)
			pos++
		}
	}

	swap(arr.Elements, pos, pivot)

	if err := quicksort(arr, l, pos-1); err != nil {
		return err
	}
	if err := quicksort(arr, pos+1, r); err != nil {
		return err
	}
	return nil
}

func swap(elements []Item.Item, i, j int) {
	elements[i], elements[j] = elements[j], elements[i]
}
