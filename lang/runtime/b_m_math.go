package runtime

import (
	"math"
)

var b_m_math = createMathModule()

func createMathModule() *Instance {
	module := Module.Create("math")

	Module.Add(module, "pi", Number.Create(3.141592653589793))
	Module.Add(module, "e", Number.Create(2.718281828459045))

	Module.Add(module, "abs", fn("abs", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Abs(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "acos", fn("acos", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Acos(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "acosh", fn("acosh", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Acosh(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "asin", fn("asin", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Asin(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "asinh", fn("asinh", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Asinh(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "atan", fn("atan", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Atan(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "atan2", fn("atan2", p("num1"), p("num2")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Atan2(AsNumber(args[0]), AsNumber(args[1])))
		}),
	)

	Module.Add(module, "atanh", fn("atanh", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Atanh(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "cbrt", fn("cbrt", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			n := AsNumber(args[0])
			if n < 0 {
				throw(r, s, "number must be positive")
			}

			return Number.Create(math.Cbrt(n))
		}),
	)

	Module.Add(module, "ceil", fn("ceil", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Ceil(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "copysign", fn("copysign", p("num"), p("sign")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Copysign(AsNumber(args[0]), AsNumber(args[1])))
		}),
	)

	Module.Add(module, "cos", fn("cos", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Cos(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "cosh", fn("cosh", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Cosh(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "deg2rad", fn("deg2rad", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(AsNumber(args[0]) * math.Pi / 180)
		}),
	)

	Module.Add(module, "erf", fn("erf", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Erf(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "erfc", fn("erfc", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Erfc(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "exp", fn("exp", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Exp(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "exp2", fn("exp2", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Exp2(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "expm1", fn("expm1", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Expm1(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "factorial", fn("factorial", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			n := AsNumber(args[0])
			if n < 0 {
				throw(r, s, "number must be positive")
			}

			return Number.Create(_factorial(n))
		}),
	)

	Module.Add(module, "floor", fn("floor", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Floor(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "gamma", fn("gamma", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			n := AsNumber(args[0])
			if n < 0 {
				throw(r, s, "number must be positive")
			}

			return Number.Create(math.Gamma(n))
		}),
	)

	Module.Add(module, "gcd", fn("gcd", p("num1"), p("num2")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			a := AsInteger(args[0])
			b := AsInteger(args[1])

			for b != 0 {
				t := b
				b = a % b
				a = t

				if b == 0 {
					return Number.Create(float64(a))
				}
			}

			return Number.Create(float64(a))
		}),
	)

	Module.Add(module, "hypot", fn("hypot", p("num1"), p("num2")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Hypot(AsNumber(args[0]), AsNumber(args[1])))
		}),
	)

	Module.Add(module, "log", fn("log", p("num"), p("base", Number.Create(2.718281828459045))).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			n := AsNumber(args[0])
			if len(args) == 1 {
				return Number.Create(math.Log(n))
			}

			base := AsNumber(args[1])
			if base == 10 {
				return Number.Create(math.Log10(n))
			} else if base == 2 {
				return Number.Create(math.Log2(n))
			}

			return Number.Create(math.Log(n) / math.Log(base))
		}),
	)

	Module.Add(module, "max", fn("max", p("num1"), p("num2")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Max(AsNumber(args[0]), AsNumber(args[1])))
		}),
	)

	Module.Add(module, "min", fn("min", p("num1"), p("num2")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Min(AsNumber(args[0]), AsNumber(args[1])))
		}),
	)

	Module.Add(module, "pow", fn("pow", p("num"), p("exp")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Pow(AsNumber(args[0]), AsNumber(args[1])))
		}),
	)

	Module.Add(module, "rad2deg", fn("rad2deg", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(AsNumber(args[0]) * 180 / math.Pi)
		}),
	)

	Module.Add(module, "round", fn("round", p("num"), p("precision", Number.Create(0))).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Round(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "sin", fn("sin", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Sin(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "sincos", fn("sincos", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			sin, cos := math.Sincos(AsNumber(args[0]))
			return Tuple.Create(Number.Create(sin), Number.Create(cos))
		}),
	)

	Module.Add(module, "sinh", fn("sinh", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Sinh(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "sqrt", fn("sqrt", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			n := AsNumber(args[0])
			if n < 0 {
				throw(r, s, "number must be positive")
			}

			return Number.Create(math.Sqrt(n))
		}),
	)

	Module.Add(module, "tan", fn("tan", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Tan(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "tanh", fn("tanh", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Tanh(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "trunc", fn("trunc", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Number.Create(math.Trunc(AsNumber(args[0])))
		}),
	)

	Module.Add(module, "even", fn("even", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Boolean.Create(AsInteger(args[0])%2 == 0)
		}),
	)

	Module.Add(module, "odd", fn("odd", p("num")).
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			return Boolean.Create(AsInteger(args[0])%2 == 1)
		}),
	)

	Module.Add(module, "primes", fn("primes").
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
		}),
	)

	Module.Add(module, "fibonacci", fn("fibonacci").
		as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			cur := -1
			a := 1
			b := 1
			return i(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				cur++
				if cur == 0 {
					return Iteration.Create(Number.ONE)
				} else if cur == 1 {
					return Iteration.Create(Number.ONE)
				} else {
					c := a + b
					a = b
					b = c
					return Iteration.Create(Number.Create(float64(c)))
				}
			})
		}),
	)

	return module
}

func _factorial(n float64) float64 {
	if n == 0 {
		return 1
	}
	return n * _factorial(n-1)
}
