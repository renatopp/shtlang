package runtime

import (
	"sht/lang/ast"
)

var booleanDT = &BooleanDataType{
	BaseDataType: BaseDataType{
		Name:        "Boolean",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var Boolean = &BooleanInfo{
	Type: booleanDT,

	TRUE: &Instance{
		Type: booleanDT,
		Impl: BooleanDataImpl{
			Value: true,
		},
	},
	FALSE: &Instance{
		Type: booleanDT,
		Impl: BooleanDataImpl{
			Value: false,
		},
	},
}

// ----------------------------------------------------------------------------
// BOOLEAN INFO
// ----------------------------------------------------------------------------
type BooleanInfo struct {
	Type DataType

	TRUE  *Instance
	FALSE *Instance
}

func (t *BooleanInfo) Create(value bool) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: BooleanDataImpl{
			Value: value,
		},
	}
}

// ----------------------------------------------------------------------------
// BOOLEAN DATA TYPE
// ----------------------------------------------------------------------------
type BooleanDataType struct {
	BaseDataType
}

func (d *BooleanDataType) OnBoolean(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return args[0]
}

func (d *BooleanDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args...)
}

func (d *BooleanDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	v := AsBool(args[0])
	if v {
		return String.Create("true")
	}

	return String.Create("false")
}

func (d *BooleanDataType) OnNot(r *Runtime, s *Scope, args ...*Instance) *Instance {
	v := AsBool(args[0])
	return Boolean.Create(!v)
}

func (n *BooleanDataType) OnEq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Boolean.FALSE
	}

	this := AsBool(args[0])
	other := AsBool(args[1])

	if this == other {
		return Boolean.TRUE
	} else {
		return Boolean.FALSE
	}
}

func (n *BooleanDataType) OnNeq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Boolean.FALSE
	}

	this := AsBool(args[0])
	other := AsBool(args[1])

	if this != other {
		return Boolean.TRUE
	} else {
		return Boolean.FALSE
	}
}

// ----------------------------------------------------------------------------
// BOOLEAN DATA IMPL
// ----------------------------------------------------------------------------
type BooleanDataImpl struct {
	Value bool
}
