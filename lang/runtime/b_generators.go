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

var b_fibonacci = fn("fibonacci").
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		cur := -1
		a := 0
		b := 1
		return i(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			cur++
			if cur == 0 {
				return Iteration.Create(Number.ZERO)
			} else if cur == 1 {
				return Iteration.Create(Number.ONE)
			} else {
				c := a + b
				a = b
				b = c
				return Iteration.Create(Number.Create(float64(c)))
			}
		})
	})

var b_primes = fn("primes").
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		D := map[int]int{}
		q := 0
		return i(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if q == 0 {
				q = 1
				return Iteration.Create(Number.TWO)
			}

			for {
				q += 2
				p, ok := D[q]
				if !ok {
					delete(D, q)
					D[q*q] = q
					return Iteration.Create(Number.Create(float64(q)))
				} else {
					x := q + 2*p
					for _, ok := D[x]; ok; {
						x += 2 * p
						_, ok = D[x]
					}
					D[x] = p
				}
			}
		})
	})
