package runtime

import (
	"fmt"
	"strings"
)

var b_print = Function.CreateNative("print",
	[]*FunctionParam{
		{"msgs", nil, true},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		if len(args) == 0 {
			fmt.Println()
			return String.EMPTY
		}

		msgs := []string{}
		for _, arg := range args {
			msgs = append(msgs, AsString(arg))
		}

		final := strings.Join(msgs, " ")
		fmt.Println(final)
		return String.Create(final)
	},
)

var b_len = Function.CreateNative("len",
	[]*FunctionParam{
		{"obj", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		if len(args) == 0 {
			return r.Throw(Error.Create(s, "Expecting an object on len function, got none"), s)
		}

		return args[0].OnLen(r, s)
	},
)
