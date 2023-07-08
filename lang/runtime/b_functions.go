package runtime

import (
	"fmt"
	"strings"
)

var b_print = fn("print", p("msgs", String.EMPTY, true)).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
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
	})

var b_printf = fn("printf", p("msg"), p("values", nil, true)).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
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
	})

var b_len = fn("len", p("obj")).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		obj, err := arg(args, 0).Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		return obj.OnLen(r, s)
	})

var b_even = fn("even", p("num")).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		num, err := arg(args, 0).IsNumber().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		return Boolean.Create(AsInteger(num)%2 == 0)
	})

var b_odd = fn("odd", p("num")).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		num, err := arg(args, 0).IsNumber().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		return Boolean.Create(AsInteger(num)%2 == 1)
	})

var b_palindrome = fn("palindrome", p("str")).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		str, err := arg(args, 0).Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		value := AsString(str.OnString(r, s))
		for i := 0; i < len(value)/2; i++ {
			if value[i] != value[len(value)-1-i] {
				return Boolean.FALSE
			}
		}
		return Boolean.TRUE
	})
