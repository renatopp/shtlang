package runtime

var b_map = Function.CreateNative("map",
	[]*FunctionParam{
		{"iter", nil, false},
		{"func", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		if len(args) != 2 {
			return r.Throw(Error.Create(s, "map does not accept additional parameters"), s)
		}
		if args[1] == Boolean.FALSE {
			return r.Throw(Error.Create(s, "map requires a function"), s)
		}

		iter := args[0]
		next := args[0].Impl.(*IteratorDataImpl).next()
		fn := args[1]

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				ret := next.OnCall(r, s, iter)
				if s.IsInterruptedAs(FlowRaise) {
					return ret
				}
				iteration := ret.AsIteration()

				if AsBool(iteration.error()) {
					return ret

				} else if AsBool(iteration.done()) {
					return Iteration.DONE

				} else {
					values := iteration.value().AsTuple()
					ret := fn.OnCall(r, s, values.Values...)
					return Iteration.Create(ret)
				}
			}),
		)
	},
)

var b_each = Function.CreateNative("each",
	[]*FunctionParam{
		{"iter", nil, false},
		{"func", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		if len(args) != 2 {
			return r.Throw(Error.Create(s, "each does not accept additional parameters"), s)
		}
		if args[1] == Boolean.FALSE {
			return r.Throw(Error.Create(s, "each requires a function"), s)
		}

		iter := args[0]
		next := args[0].Impl.(*IteratorDataImpl).next()
		fn := args[1]

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				ret := next.OnCall(r, s, iter)
				if s.IsInterruptedAs(FlowRaise) {
					return ret
				}
				iteration := ret.Impl.(*IterationDataImpl)

				if iteration.error() == Boolean.TRUE {
					return ret

				} else if iteration.done() == Boolean.TRUE {
					return Iteration.DONE

				} else {
					values := iteration.value().Impl.(*TupleDataImpl)
					fn.OnCall(r, s, values.Values...)
					return ret
				}
			}),
		)
	},
)

var b_filter = Function.CreateNative("filter",
	[]*FunctionParam{
		{"iter", nil, false},
		{"func", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		if len(args) != 2 {
			return r.Throw(Error.Create(s, "filter does not accept additional parameters"), s)
		}
		if args[1] == Boolean.FALSE {
			return r.Throw(Error.Create(s, "filter requires a function"), s)
		}

		iter := args[0]
		next := args[0].Impl.(*IteratorDataImpl).next()
		fn := args[1]

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				for {
					ret := next.OnCall(r, s, iter)
					if s.IsInterruptedAs(FlowRaise) {
						return ret
					}

					iteration := ret.Impl.(*IterationDataImpl)

					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						return Iteration.DONE

					} else {
						values := iteration.value().Impl.(*TupleDataImpl)
						r := fn.OnCall(r, s, values.Values...)
						if AsBool(r) {
							return ret
						}
					}
				}
			}),
		)
	},
)

var b_reduce = Function.CreateNative("reduce",
	[]*FunctionParam{
		{"iter", nil, false},
		{"func", nil, false},
		{"default", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		// args[2] => default
		if len(args) > 3 {
			return r.Throw(Error.Create(s, "reduce does not accept additional parameters"), s)
		}
		if args[1] == Boolean.FALSE {
			return r.Throw(Error.Create(s, "map requires a function"), s)
		}

		iter := args[0]
		next := args[0].Impl.(*IteratorDataImpl).next()
		fn := args[1]
		acc := Number.ZERO

		if len(args) >= 3 {
			acc = args[2]
		}

		finished := false

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				if finished {
					return Iteration.DONE
				}

				for {
					ret := next.OnCall(r, s, iter)
					if s.IsInterruptedAs(FlowRaise) {
						return ret
					}

					iteration := ret.Impl.(*IterationDataImpl)

					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						finished = true
						return Iteration.Create(acc)

					} else {
						values := iteration.value().Impl.(*TupleDataImpl)
						acc = fn.OnCall(r, s, append([]*Instance{acc}, values.Values...)...)
					}
				}
			}),
		)
	},
)

var b_sum = fn("sum", p("iter"), p("first", GetFirstFn)).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		i_iter, err := arg(args, 0).IsIterator().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		i_fn, err := arg(args, 1).Optional(GetFirstFn).IsFunction().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		next := i_iter.AsIterator().next()
		total := Number.Create(0.0)

		finished := false

		return i(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if finished {
				return Iteration.DONE
			}

			for {
				ret := next.OnCall(r, s, i_iter)
				if s.IsInterruptedAs(FlowRaise) {
					return ret
				}
				iteration := ret.AsIteration()

				if iteration.error() == Boolean.TRUE {
					return ret

				} else if iteration.done() == Boolean.TRUE {
					finished = true
					return Iteration.Create(total)

				} else {
					val := i_fn.OnCall(r, s, iteration.value().AsTuple().Values...)
					total = total.OnAdd(r, s, val)
				}
			}
		})
	})

var b_takeWhile = Function.CreateNative("takeWhile",
	[]*FunctionParam{
		{"iter", nil, false},
		{"func", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		if len(args) != 2 {
			return r.Throw(Error.Create(s, "takeWhile does not accept additional parameters"), s)
		}
		if args[1] == Boolean.FALSE {
			return r.Throw(Error.Create(s, "takeWhile requires a function"), s)
		}

		iter := args[0]
		next := args[0].AsIterator().next()
		fn := args[1]

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				for {
					ret := next.OnCall(r, s, iter)
					if s.IsInterruptedAs(FlowRaise) {
						return ret
					}

					iteration := ret.Impl.(*IterationDataImpl)

					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						return Iteration.DONE

					} else {
						values := iteration.value().Impl.(*TupleDataImpl)
						r := fn.OnCall(r, s, values.Values...)
						if AsBool(r) {
							return ret
						} else {
							return Iteration.DONE
						}
					}
				}
			}),
		)
	},
)

var b_take = Function.CreateNative("take",
	[]*FunctionParam{
		{"iter", nil, false},
		{"func", nil, false},
		{"amount", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		// args[2] => amount
		if len(args) > 3 {
			return r.Throw(Error.Create(s, "take does not accept additional parameters"), s)
		}
		if len(args) < 3 {
			return r.Throw(Error.Create(s, "take requires an amount"), s)
		}
		if !args[2].IsNumber() {
			return r.Throw(Error.Create(s, "take requires a number as parameter"), s)
		}

		if args[1] != Boolean.FALSE {
			return r.Throw(Error.Create(s, "take does not support a function"), s)
		}

		iter := args[0]
		next := args[0].AsIterator().next()
		amount := AsInteger(args[2])
		total := 0
		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				for total < amount {
					ret := next.OnCall(r, s, iter)
					if s.IsInterruptedAs(FlowRaise) {
						return ret
					}

					iteration := ret.AsIteration()

					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						return Iteration.DONE

					} else {
						total++
						return ret
					}
				}
				return Iteration.DONE
			}),
		)
	},
)

var b_min = fn("min", p("iter"), p("func", GetFirstFn)).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		i_iter, err := arg(args, 0).IsIterator().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		i_fn, err := arg(args, 1).Optional(GetFirstFn).IsFunction().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		next := i_iter.AsIterator().next()
		var min *Instance
		var minValue *Instance

		finished := false
		return i(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if finished {
				return Iteration.DONE
			}

			for {
				ret := next.OnCall(r, s, i_iter)
				if s.IsInterruptedAs(FlowRaise) {
					return ret
				}

				iteration := ret.AsIteration()
				if iteration.error() == Boolean.TRUE {
					return ret

				} else if iteration.done() == Boolean.TRUE {
					finished = true
					if min == nil {
						return Iteration.DONE
					}

					return Iteration.Create(min)

				} else {
					values := iteration.value().AsTuple().Values[0]
					val := i_fn.OnCall(r, s, values)
					if min == nil || AsBool(val.OnLt(r, s, minValue)) {
						min = values
						minValue = val
					}
				}
			}
		})
	})

var b_max = fn("max", p("iter"), p("func", GetFirstFn)).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		i_iter, err := arg(args, 0).IsIterator().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		i_fn, err := arg(args, 1).Optional(GetFirstFn).IsFunction().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		next := i_iter.AsIterator().next()
		var max *Instance
		var maxValue *Instance

		finished := false
		return i(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if finished {
				return Iteration.DONE
			}

			for {
				ret := next.OnCall(r, s, i_iter)
				if s.IsInterruptedAs(FlowRaise) {
					return ret
				}

				iteration := ret.AsIteration()
				if iteration.error() == Boolean.TRUE {
					return ret

				} else if iteration.done() == Boolean.TRUE {
					finished = true
					if max == nil {
						return Iteration.DONE
					}

					return Iteration.Create(max)

				} else {
					values := iteration.value().AsTuple().Values[0]
					val := i_fn.OnCall(r, s, values)
					if max == nil || AsBool(val.OnGt(r, s, maxValue)) {
						max = values
						maxValue = val
					}
				}
			}
		})
	})

var b_first = Function.CreateNative("first",
	[]*FunctionParam{
		{"iter", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		if len(args) > 2 {
			return r.Throw(Error.Create(s, "first does not accept additional parameters"), s)
		}
		if args[1] != Boolean.FALSE {
			return r.Throw(Error.Create(s, "first does not accepts a function"), s)
		}

		iter := args[0]
		next := args[0].AsIterator().next()

		finished := false
		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				if finished {
					return Iteration.DONE
				}

				for {
					ret := next.OnCall(r, s, iter)
					if s.IsInterruptedAs(FlowRaise) {
						return ret
					}

					iteration := ret.AsIteration()
					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						return ret

					} else {
						finished = true
						return ret
					}
				}
			}),
		)
	},
)

var b_last = Function.CreateNative("last",
	[]*FunctionParam{
		{"iter", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		if len(args) > 2 {
			return r.Throw(Error.Create(s, "last does not accept additional parameters"), s)
		}
		if args[1] != Boolean.FALSE {
			return r.Throw(Error.Create(s, "last does not accepts a function"), s)
		}

		iter := args[0]
		next := args[0].AsIterator().next()

		var last *Instance
		finished := false
		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				if finished {
					return Iteration.DONE
				}

				for {
					ret := next.OnCall(r, s, iter)
					if s.IsInterruptedAs(FlowRaise) {
						return ret
					}

					iteration := ret.AsIteration()
					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						finished = true
						if last == nil {
							return Iteration.DONE
						}
						return last

					} else {
						last = ret
					}
				}
			}),
		)
	},
)

var b_window = fn("window", p("iter"), p("size")).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		iter, err := arg(args, 0).IsIterator().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		i_size, err := arg(args, 2).IsNumber().Validate()
		if err != nil {
			return throw(r, s, err.Error())
		}

		size := AsInteger(i_size)
		if size < 1 {
			return throw(r, s, "window size must be greater than 0")
		}

		next := iter.AsIterator().next()
		values := make([]*Instance, size)
		total := 0
		return i(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			for total < size {
				ret := next.OnCall(r, s, iter)
				if s.IsInterruptedAs(FlowRaise) {
					return ret
				}

				iteration := ret.AsIteration()
				if iteration.error() == Boolean.TRUE {
					return ret

				} else if iteration.done() == Boolean.TRUE {
					return Iteration.DONE

				} else {
					values[total] = iteration.value()
					total++
				}

				if total == size {
					v := []*Instance{}
					for _, value := range values {
						v = append(v, value.AsTuple().Values...)
					}
					return Iteration.Create(Tuple.Create(v...))
				}
			}

			ret := next.OnCall(r, s, iter)
			if s.IsInterruptedAs(FlowRaise) {
				return ret
			}

			iteration := ret.AsIteration()
			if iteration.error() == Boolean.TRUE {
				return ret

			} else if iteration.done() == Boolean.TRUE {
				return Iteration.DONE

			} else {
				for i := 0; i < size-1; i++ {
					values[i] = values[i+1]
				}
				values[size-1] = iteration.value()
				v := []*Instance{}
				for _, value := range values {
					v = append(v, value.AsTuple().Values...)
				}
				return Iteration.Create(Tuple.Create(v...))
			}
		})
	})

var b_multiply = fn("multiply", p("next")).
	as(func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		if len(args) > 2 {
			return r.Throw(Error.Create(s, "multiply does not accept additional parameters"), s)
		}
		if args[1] != Boolean.FALSE {
			return r.Throw(Error.Create(s, "multiply does not accepts a function"), s)
		}

		iter := args[0]
		next := args[0].Impl.(*IteratorDataImpl).next()
		total := Number.Create(1.0)

		finished := false

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
				if finished {
					return Iteration.DONE
				}

				for {
					ret := next.OnCall(r, s, iter)
					if s.IsInterruptedAs(FlowRaise) {
						return ret
					}
					iteration := ret.Impl.(*IterationDataImpl)

					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						finished = true
						return Iteration.Create(total)

					} else {
						values := iteration.value().Impl.(*TupleDataImpl)
						total = total.OnMul(r, s, values.Values[0])
					}
				}
			}),
		)
	},
	)
