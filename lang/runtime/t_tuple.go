package runtime

import (
	"sht/lang/ast"
	"strings"
)

var tupleDT = &TupleDataType{
	BaseDataType: BaseDataType{
		Name:        "Tuple",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
	},
}

var Tuple = &TupleInfo{
	Type: tupleDT,
}

// ----------------------------------------------------------------------------
// TUPLE INFO
// ----------------------------------------------------------------------------
type TupleInfo struct {
	Type DataType
}

func (t *TupleInfo) Create(values ...*Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &TupleDataImpl{
			Values: values,
		},
	}
}

// ----------------------------------------------------------------------------
// TUPLE DATA TYPE
// ----------------------------------------------------------------------------
type TupleDataType struct {
	BaseDataType
}

func (t *TupleDataType) OnGetItem(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*TupleDataImpl)

	nargs := len(args)
	if nargs > 1 && !IsNumber(args[1]) {
		return r.Throw(Error.Create(s, "index of a tuple must be a number, '%s' provided", args[1].Type.GetName()), s)
	}

	if nargs > 2 && !IsNumber(args[2]) {
		return r.Throw(Error.Create(s, "index of a tuple must be a number, '%s' provided", args[2].Type.GetName()), s)
	}

	if nargs > 3 {
		return r.Throw(Error.Create(s, "tuple indexing accepts only 1 or 2 parameters, %d given", nargs-1), s)
	}

	if nargs == 2 {
		return this.Values[AsInteger(args[1])]
	}

	if nargs == 1 {
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

func (d *TupleDataType) OnEq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Boolean.FALSE
	}

	this := args[0].Impl.(*TupleDataImpl)
	other := args[1].Impl.(*TupleDataImpl)

	if len(this.Values) != len(other.Values) {
		return Boolean.FALSE
	}

	for i, value := range this.Values {
		t := value
		o := other.Values[i]

		if t.Type != o.Type {
			return Boolean.FALSE
		}

		if !AsBool(t.Type.OnEq(r, s, t, o)) {
			return Boolean.FALSE
		}
	}

	return Boolean.TRUE
}

func (d *TupleDataType) OnNeq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	v := AsBool(d.OnEq(r, s, args...))
	if v {
		return Boolean.FALSE
	}
	return Boolean.TRUE
}

func (d *TupleDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args[0])
}

func (d *TupleDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	tuple := args[0].Impl.(*TupleDataImpl)

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
