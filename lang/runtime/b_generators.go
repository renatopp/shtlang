package runtime

var b_range = Function.CreateNative("range",
	[]*FunctionParam{
		{"min", nil, false},
		{"max", Boolean.FALSE, false},
		{"step", Number.ONE, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => min
		// args[1] => max
		// args[2] => step

		if len(args) <= 0 {
			return r.Throw(Error.Create(s, "range requires at least one argument"), s)
		}

		min := 0.0
		max := AsNumber(args[0])
		step := 1.0

		if len(args) > 1 && args[1] != Boolean.FALSE {
			min = max
			max = AsNumber(args[1])
		}

		if len(args) > 2 && args[2] != Boolean.FALSE {
			step = AsNumber(args[2])
		}

		cur := min
		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				if min < max && cur >= max {
					return Iteration.DONE
				}

				if min > max && cur <= max {
					return Iteration.DONE
				}

				val := Iteration.Create(Number.Create(cur))
				cur = cur + step
				return val
			}),
		)
	},
)
