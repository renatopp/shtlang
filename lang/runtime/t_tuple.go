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

// TODO GetItem

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
