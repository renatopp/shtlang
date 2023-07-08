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

var b_sum = Function.CreateNative("sum",
	[]*FunctionParam{
		{"iter", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		if len(args) > 2 {
			return r.Throw(Error.Create(s, "sum does not accept additional parameters"), s)
		}
		if args[1] != Boolean.FALSE {
			return r.Throw(Error.Create(s, "sum does not accepts a function"), s)
		}

		iter := args[0]
		next := args[0].Impl.(*IteratorDataImpl).next()
		total := Number.Create(0.0)

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
						total = total.OnAdd(r, s, values.Values[0])
					}
				}
			}),
		)
	},
)

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

var b_min = Function.CreateNative("min",
	[]*FunctionParam{
		{"iter", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		if len(args) > 2 {
			return r.Throw(Error.Create(s, "min does not accept additional parameters"), s)
		}
		if args[1] != Boolean.FALSE {
			return r.Throw(Error.Create(s, "min does not accepts a function"), s)
		}

		iter := args[0]
		next := args[0].AsIterator().next()
		var min *Instance

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
						if min == nil {
							return Iteration.DONE
						}

						return Iteration.Create(min)

					} else {
						values := iteration.value().AsTuple()
						if min == nil || AsBool(values.Values[0].OnLt(r, s, min)) {
							min = values.Values[0]
						}
					}
				}
			}),
		)
	},
)

var b_max = Function.CreateNative("max",
	[]*FunctionParam{
		{"iter", nil, false},
	},
	func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
		// self => function
		// args[0] => iter
		// args[1] => func
		if len(args) > 2 {
			return r.Throw(Error.Create(s, "max does not accept additional parameters"), s)
		}
		if args[1] != Boolean.FALSE {
			return r.Throw(Error.Create(s, "max does not accepts a function"), s)
		}

		iter := args[0]
		next := args[0].AsIterator().next()
		var max *Instance

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
						if max == nil {
							return Iteration.DONE
						}

						return Iteration.Create(max)

					} else {
						values := iteration.value().AsTuple()
						if max == nil || AsBool(values.Values[0].OnGt(r, s, max)) {
							max = values.Values[0]
						}
					}
				}
			}),
		)
	},
)

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
