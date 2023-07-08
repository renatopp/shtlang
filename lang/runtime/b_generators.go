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

var b_fibonacci = Function.CreateNative("fibonacci",
	[]*FunctionParam{},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function

		cur := -1
		a := 0
		b := 1
		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
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
			}),
		)
	},
)

var b_primes = Function.CreateNative("primes",
	[]*FunctionParam{},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function

		D := map[int]int{}
		q := 0
		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
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
			}),
		)
	},
)
