package runtime

var b_range = fn("range", p("min"), p("max", Boolean.FALSE), p("step", Number.ONE)).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		i_min, err := arg(args, 0).IsNumber().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		i_max, err := arg(args, 1).Optional(Boolean.FALSE).IsNumber().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		i_step, err := arg(args, 2).Optional(Number.ONE).IsNumber().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		min := 0.0
		max := AsNumber(i_min)
		step := AsNumber(i_step)

		if i_max != Boolean.FALSE {
			min = max
			max = AsNumber(i_max)
		}

		cur := min
		return i(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if min < max && cur >= max {
				return Iteration.DONE
			}

			if min > max && cur <= max {
				return Iteration.DONE
			}

			val := Iteration.Create(Number.Create(cur))
			cur = cur + step
			return val
		})
	})
