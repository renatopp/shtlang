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

var b_printf = Function.CreateNative("printf",
	[]*FunctionParam{
		{"msg", nil, false},
		{"values", nil, true},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		if len(args) == 0 {
			fmt.Println()
			return String.EMPTY
		}

		if len(args) == 1 {
			str := AsString(args[1])
			fmt.Println(str)
			return String.Create(str)
		}

		msgs := []any{}
		for _, arg := range args[1:] {
			msgs = append(msgs, AsString(arg))
		}

		str := AsString(args[1])
		v := fmt.Sprintf(str, msgs...)
		fmt.Println(v)
		return String.Create(v)
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

var b_even = Function.CreateNative("even",
	[]*FunctionParam{
		{"num", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		if len(args) == 0 {
			return r.Throw(Error.Create(s, "Expecting a number, got none"), s)
		}

		if !args[0].IsNumber() {
			return r.Throw(Error.Create(s, "Expecting a number, got "+args[0].Type.GetName()), s)
		}

		return Boolean.Create(AsInteger(args[0])%2 == 0)
	},
)

var b_odd = Function.CreateNative("odd",
	[]*FunctionParam{
		{"num", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		if len(args) == 0 {
			return r.Throw(Error.Create(s, "Expecting a number, got none"), s)
		}

		if !args[0].IsNumber() {
			return r.Throw(Error.Create(s, "Expecting a number, got "+args[0].Type.GetName()), s)
		}

		return Boolean.Create(AsInteger(args[0])%2 == 1)
	},
)

var b_palindrome = Function.CreateNative("palindrome",
	[]*FunctionParam{
		{"str", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		if len(args) == 0 {
			return r.Throw(Error.Create(s, "Expecting an obj, got none"), s)
		}

		str := AsString(args[0].OnString(r, s))
		for i := 0; i < len(str)/2; i++ {
			if str[i] != str[len(str)-1-i] {
				return Boolean.FALSE
			}
		}
		return Boolean.TRUE
	},
)
