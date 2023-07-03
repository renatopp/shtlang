package runtime

import (
	"fmt"
	"strings"
)

var b_print = Function.CreateNative("print",
	[]*FunctionParam{
		{"msgs", nil, true},
	},
	func(r *Runtime, s *Scope, args ...*Instance) *Instance {
		if len(args) == 1 {
			fmt.Println()
			return String.EMPTY
		}

		msgs := []string{}
		for _, arg := range args[1:] {
			msgs = append(msgs, AsString(arg))
		}

		final := strings.Join(msgs, " ")
		fmt.Println(final)
		return String.Create(final)
	},
)
