package runtime

import (
	"sht/lang/ast"
	"strings"
)

var tupleDT = &TupleDataType{
	BaseDataType: BaseDataType{
		Name:        "Tuple",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var Tuple = &TupleInfo{
	Type: tupleDT,
}

// ----------------------------------------------------------------------------
// TUPLE INFO
// ----------------------------------------------------------------------------
type TupleInfo struct {
	Type         DataType
	TypeInstance *Instance
}

func (t *TupleInfo) Create(values ...*Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &TupleDataImpl{
			Values: values,
		},
	}
}

func (t *TupleInfo) Setup() {
	t.TypeInstance = Type.Create(Tuple.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

// ----------------------------------------------------------------------------
// TUPLE DATA TYPE
// ----------------------------------------------------------------------------
type TupleDataType struct {
	BaseDataType
}

func (d *TupleDataType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	switch init := init.(type) {
	case *ast.ListInitializer:
		values := make([]*Instance, 0)
		for _, value := range init.Values {

			if spread, ok := value.(*ast.SpreadOut); ok {
				var e *Instance
				target := r.Eval(spread.Target, s)
				r.ResolveIterator(target, s, func(v *Instance, err *Instance) {
					if err != nil {
						e = err
					} else if v != nil {
						t := v.Impl.(*TupleDataImpl)
						values = append(values, t.Values...)
					}
				})
				if e != nil {
					return e
				}
				continue
			}

			values = append(values, r.Eval(value, s))
		}
		return Tuple.Create(values...)
	default:
		return r.Throw(Error.Create(s, "invalid initializer for tuple"), s)
	}
}

func (d *TupleDataType) OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	iter := self.Impl.(*IteratorDataImpl)
	next := iter.next()
	values := []*Instance{}
	for {

		tion := next.OnCall(r, s, self).Impl.(*IterationDataImpl)

		if tion.error() == Boolean.TRUE {
			tuple := tion.value().AsTuple()
			return r.Throw(tuple.Values[0], s)

		} else if tion.done() == Boolean.TRUE {
			return Tuple.Create(values...)

		} else {
			tuple := tion.value().Impl.(*TupleDataImpl)
			values = append(values, tuple.Values[0])
		}
	}
}

func (d *TupleDataType) OnNew(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return self
}

func (d *TupleDataType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	cur := 0
	this := self.Impl.(*TupleDataImpl)
	return Iterator.Create(
		Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if cur >= len(this.Values) {
				return Iteration.DONE
			}

			cur++
			return Iteration.Create(this.Values[cur-1])
		}),
	)
}

func (d *TupleDataType) OnLen(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*TupleDataImpl)
	return Number.Create(float64(len(this.Values)))
}

func (t *TupleDataType) OnGetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*TupleDataImpl)

	nargs := len(args)
	if nargs > 0 && !IsNumber(args[0]) {
		return r.Throw(Error.Create(s, "index of a tuple must be a number, '%s' provided", args[1].Type.GetName()), s)
	}

	if nargs > 1 && !IsNumber(args[1]) {
		return r.Throw(Error.Create(s, "index of a tuple must be a number, '%s' provided", args[2].Type.GetName()), s)
	}

	if nargs > 2 {
		return r.Throw(Error.Create(s, "tuple indexing accepts only 1 or 2 parameters, %d given", nargs-1), s)
	}

	if nargs == 1 {
		return this.Values[AsInteger(args[0])]
	}

	if nargs == 0 {
		return Tuple.Create(this.Values...)
	}

	size := len(this.Values)
	idx0 := AsInteger(args[1])
	if idx0 < 0 || idx0 > size-1 {
		return r.Throw(Error.Create(s, "first index '%d' of tuple slicing out of bounds", idx0), s)
	}

	idx1 := AsInteger(args[2])
	if idx1 < 0 || idx1 > size {
		return r.Throw(Error.Create(s, "second index '%d' of tuple slicing out of bounds", idx0), s)
	}

	if idx1 <= idx0 {
		return r.Throw(Error.Create(s, "second index '%d' of tuple slicing must be greater than the first '%d'", idx1, idx0), s)
	}

	values := make([]*Instance, 0)
	for _, v := range this.Values[idx0:idx1] {
		values = append(values, v)
	}

	return Tuple.Create(values...)
}

func (d *TupleDataType) OnEq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return Boolean.FALSE
	}

	this := self.Impl.(*TupleDataImpl)
	other := args[0].Impl.(*TupleDataImpl)

	if len(this.Values) != len(other.Values) {
		return Boolean.FALSE
	}

	for i, value := range this.Values {
		t := value
		o := other.Values[i]

		if t.Type != o.Type {
			return Boolean.FALSE
		}

		if !AsBool(t.OnEq(r, s, o)) {
			return Boolean.FALSE
		}
	}

	return Boolean.TRUE
}

func (d *TupleDataType) OnNeq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	v := AsBool(d.OnEq(r, s, self, args...))
	if v {
		return Boolean.FALSE
	}
	return Boolean.TRUE
}

func (d *TupleDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self)
}

func (d *TupleDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	tuple := self.Impl.(*TupleDataImpl)

	var values []string
	for _, value := range tuple.Values {
		values = append(values, value.Repr())
	}

	if len(values) == 1 {
		return String.Create("(" + values[0] + ",)")
	} else {
		return String.Create("(" + strings.Join(values, ", ") + ")")
	}

}

// ----------------------------------------------------------------------------
// TUPLE DATA IMPL
// ----------------------------------------------------------------------------
type TupleDataImpl struct {
	Values []*Instance
}
