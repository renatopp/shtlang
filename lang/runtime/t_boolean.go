package runtime

import (
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

var Boolean = _setupBoolean()

type BooleanInfo struct {
	Instance *Instance
	Type     *DataType

	TRUE  *Instance
	FALSE *Instance
}

type BooleanImpl struct {
	Value bool
}

func _setupBoolean() *BooleanInfo {
	dataType := &DataType{
		Name:        "Boolean",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
		Meta:        map[meta.MetaName]Function{},
	}

	n := &BooleanInfo{
		Instance: Type.Create(dataType, true),
		Type:     dataType,
	}

	n.TRUE = n.Create(true, true)
	n.FALSE = n.Create(false, true)

	dataType.Meta[meta.SetProperty] = n.notImplemented(string(meta.SetProperty))
	dataType.Meta[meta.GetProperty] = n.notImplemented(string(meta.GetProperty))
	dataType.Meta[meta.SetItem] = n.notImplemented(string(meta.SetItem))
	dataType.Meta[meta.GetItem] = n.notImplemented(string(meta.GetItem))
	dataType.Meta[meta.Call] = n.notImplemented(string(meta.Call))

	dataType.Meta[meta.Boolean] = CreateNativeFunction(n.MetaBoolean)
	dataType.Meta[meta.String] = CreateNativeFunction(n.MetaString)
	dataType.Meta[meta.Repr] = CreateNativeFunction(n.MetaRepr)

	dataType.Meta[meta.Add] = n.notImplemented(string(meta.Add))
	dataType.Meta[meta.Sub] = n.notImplemented(string(meta.Sub))
	dataType.Meta[meta.Mul] = n.notImplemented(string(meta.Mul))
	dataType.Meta[meta.Div] = n.notImplemented(string(meta.Div))
	dataType.Meta[meta.IntDiv] = n.notImplemented(string(meta.IntDiv))
	dataType.Meta[meta.Mod] = n.notImplemented(string(meta.Mod))
	dataType.Meta[meta.Pow] = n.notImplemented(string(meta.Pow))
	dataType.Meta[meta.Eq] = CreateNativeFunction(n.MetaEq)
	dataType.Meta[meta.Neq] = CreateNativeFunction(n.MetaNeq)
	dataType.Meta[meta.Gt] = n.notImplemented(string(meta.Gt))
	dataType.Meta[meta.Lt] = n.notImplemented(string(meta.Lt))
	dataType.Meta[meta.Gte] = n.notImplemented(string(meta.Gte))
	dataType.Meta[meta.Lte] = n.notImplemented(string(meta.Lte))
	dataType.Meta[meta.Pos] = n.notImplemented(string(meta.Pos))
	dataType.Meta[meta.Neg] = n.notImplemented(string(meta.Neg))
	dataType.Meta[meta.Not] = CreateNativeFunction(n.MetaNot)
	dataType.Meta[meta.PostInc] = n.notImplemented(string(meta.PostInc))
	dataType.Meta[meta.PostDec] = n.notImplemented(string(meta.PostDec))

	return n
}

// ----------------------------------------------------------------------------
// Boolean Implementation
// ----------------------------------------------------------------------------
func (n BooleanImpl) Repr() string {
	if n.Value {
		return "true"
	}

	return "false"
}

// ----------------------------------------------------------------------------
// Boolean Info
// ----------------------------------------------------------------------------
func (n *BooleanInfo) Create(value bool, constant bool) *Instance {
	return &Instance{
		Type:  n.Type,
		Impl:  BooleanImpl{Value: value},
		Const: constant,
	}
}

func (n *BooleanInfo) val(instance *Instance) bool {
	return instance.Impl.(BooleanImpl).Value
}

func (n *BooleanInfo) MetaBoolean(r *Runtime, args ...*Instance) *Instance {
	return args[0]
}

func (n *BooleanInfo) MetaString(r *Runtime, args ...*Instance) *Instance {
	return String.Create(args[0].Impl.Repr(), false)
}

func (n *BooleanInfo) MetaRepr(r *Runtime, args ...*Instance) *Instance {
	return String.Create(args[0].Impl.Repr(), false)
}

func (n *BooleanInfo) MetaNot(r *Runtime, args ...*Instance) *Instance {
	this := n.val(args[0])
	if this {
		return n.FALSE
	} else {
		return n.TRUE
	}
}

func (n *BooleanInfo) MetaEq(r *Runtime, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return n.FALSE
	}

	this := n.val(args[0])
	other := n.val(args[1])

	if this == other {
		return n.TRUE
	} else {
		return n.FALSE
	}
}

func (n *BooleanInfo) MetaNeq(r *Runtime, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return n.FALSE
	}

	this := n.val(args[0])
	other := n.val(args[1])

	if this != other {
		return n.TRUE
	} else {
		return n.FALSE
	}
}

func (n *BooleanInfo) notImplemented(name string) Function {
	return CreateNativeFunction(func(r *Runtime, args ...*Instance) *Instance {
		return NotImplemented(name, n.Instance)
	})
}
