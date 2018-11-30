package evaluator

import "monkeylang/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Invalid number of argument to `len` function. Expected: 1, Got: %d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("Invalid argument to `len` function. Got: %s", args[0].Type())
			}
		},
	},
}
