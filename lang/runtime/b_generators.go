package runtime

var b_range = Function.CreateNative("range",
	[]*FunctionParam{
		{"min", nil, false},
		{"max", Boolean.FALSE, false},
	},
	func(r *Runtime, s *Scope, args ...*Instance) *Instance {
		// args[0] => function
		// args[1] => min
		// args[2] => max
		// args[3] => step

		if len(args) <= 1 {
			return r.Throw(Error.Create(s, "range requires at least one argument"), s)
		}

		min := 0.0
		max := AsNumber(args[1])
		step := 1.0

		if len(args) > 2 && args[2] != Boolean.FALSE {
			min = max
			max = AsNumber(args[2])
		}

		if len(args) > 3 && args[3] != Boolean.FALSE {
			step = AsNumber(args[3])
		}

		cur := min
		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, args ...*Instance) *Instance {
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
