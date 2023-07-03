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
	Type         DataType
	TypeInstance *Instance

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

func (t *BooleanInfo) Setup() {
	t.TypeInstance = Type.Create(Boolean.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

// ----------------------------------------------------------------------------
// BOOLEAN DATA TYPE
// ----------------------------------------------------------------------------
type BooleanDataType struct {
	BaseDataType
}

func (d *BooleanDataType) OnTo(r *Runtime, s *Scope, args ...*Instance) *Instance {
	iter := args[0].Impl.(*IteratorDataImpl)
	next := iter.next()
	tion := next.Type.OnCall(r, s, next, args[0]).Impl.(*IterationDataImpl)

	if tion.error() == Boolean.TRUE {
		return Boolean.FALSE
	} else if tion.done() == Boolean.TRUE {
		return r.Throw(Error.Create(s, "The iteration has been finished"), s)
	} else {
		tuple := tion.value().Impl.(*TupleDataImpl)
		return Boolean.Create(AsBool(tuple.Values[0]))
	}
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
