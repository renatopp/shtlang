package runtime

var b_map = Function.CreateNative("map",
	[]*FunctionParam{
		{"iter", nil, false},
		{"func", nil, false},
	},
	func(r *Runtime, s *Scope, args ...*Instance) *Instance {
		// args[0] => function
		// args[1] => iter
		// args[2] => func
		if len(args) != 3 {
			return r.Throw(Error.Create(s, "map does not accept additional parameters"), s)
		}
		if args[2] == Boolean.FALSE {
			return r.Throw(Error.Create(s, "map requires a function"), s)
		}

		iter := args[1]
		next := args[1].Impl.(*IteratorDataImpl).next()
		fn := args[2]

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, args ...*Instance) *Instance {
				ret := next.Type.OnCall(r, s, next, iter)
				iteration := ret.Impl.(*IterationDataImpl)

				if iteration.error() == Boolean.TRUE {
					return ret

				} else if iteration.done() == Boolean.TRUE {
					return Iteration.DONE

				} else {
					values := iteration.value().Impl.(*TupleDataImpl)
					ret := fn.Type.OnCall(r, s, append([]*Instance{fn}, values.Values...)...)
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
	func(r *Runtime, s *Scope, args ...*Instance) *Instance {
		// args[0] => function
		// args[1] => iter
		// args[2] => func
		if len(args) != 3 {
			return r.Throw(Error.Create(s, "each does not accept additional parameters"), s)
		}
		if args[2] == Boolean.FALSE {
			return r.Throw(Error.Create(s, "each requires a function"), s)
		}

		iter := args[1]
		next := args[1].Impl.(*IteratorDataImpl).next()
		fn := args[2]

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, args ...*Instance) *Instance {
				ret := next.Type.OnCall(r, s, next, iter)
				iteration := ret.Impl.(*IterationDataImpl)

				if iteration.error() == Boolean.TRUE {
					return ret

				} else if iteration.done() == Boolean.TRUE {
					return Iteration.DONE

				} else {
					values := iteration.value().Impl.(*TupleDataImpl)
					fn.Type.OnCall(r, s, append([]*Instance{fn}, values.Values...)...)
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
	func(r *Runtime, s *Scope, args ...*Instance) *Instance {
		// args[0] => function
		// args[1] => iter
		// args[2] => func
		if len(args) != 3 {
			return r.Throw(Error.Create(s, "filter does not accept additional parameters"), s)
		}
		if args[2] == Boolean.FALSE {
			return r.Throw(Error.Create(s, "filter requires a function"), s)
		}

		iter := args[1]
		next := args[1].Impl.(*IteratorDataImpl).next()
		fn := args[2]

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, args ...*Instance) *Instance {
				for {
					ret := next.Type.OnCall(r, s, next, iter)
					iteration := ret.Impl.(*IterationDataImpl)

					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						return Iteration.DONE

					} else {
						values := iteration.value().Impl.(*TupleDataImpl)
						r := fn.Type.OnCall(r, s, append([]*Instance{fn}, values.Values...)...)
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
	func(r *Runtime, s *Scope, args ...*Instance) *Instance {
		// args[0] => function
		// args[1] => iter
		// args[2] => func
		// args[3] => default
		if len(args) > 4 {
			return r.Throw(Error.Create(s, "reduce does not accept additional parameters"), s)
		}
		if args[2] == Boolean.FALSE {
			return r.Throw(Error.Create(s, "map requires a function"), s)
		}

		iter := args[1]
		next := args[1].Impl.(*IteratorDataImpl).next()
		fn := args[2]
		acc := Number.ZERO

		if len(args) >= 4 {
			acc = args[3]
		}

		finished := false

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, args ...*Instance) *Instance {
				if finished {
					return Iteration.DONE
				}

				for {
					ret := next.Type.OnCall(r, s, next, iter)
					iteration := ret.Impl.(*IterationDataImpl)

					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						finished = true
						return Iteration.Create(acc)

					} else {
						values := iteration.value().Impl.(*TupleDataImpl)
						acc = fn.Type.OnCall(r, s, append([]*Instance{fn, acc}, values.Values...)...)
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
	func(r *Runtime, s *Scope, args ...*Instance) *Instance {
		// args[0] => function
		// args[1] => iter
		// args[2] => func
		if len(args) > 3 {
			return r.Throw(Error.Create(s, "sum does not accept additional parameters"), s)
		}
		if args[2] != Boolean.FALSE {
			return r.Throw(Error.Create(s, "sum does not accepts a function"), s)
		}

		iter := args[1]
		next := args[1].Impl.(*IteratorDataImpl).next()
		total := 0.0

		finished := false

		return Iterator.Create(
			Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, args ...*Instance) *Instance {
				if finished {
					return Iteration.DONE
				}

				for {
					ret := next.Type.OnCall(r, s, next, iter)
					iteration := ret.Impl.(*IterationDataImpl)

					if iteration.error() == Boolean.TRUE {
						return ret

					} else if iteration.done() == Boolean.TRUE {
						finished = true
						return Iteration.Create(Number.Create(total))

					} else {
						values := iteration.value().Impl.(*TupleDataImpl)
						total += AsNumber(values.Values[0])
					}
				}
			}),
		)
	},
)
