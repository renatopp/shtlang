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

func (n *BooleanInfo) invalid(name string) Function {
	return CreateNativeFunction(func(r *Runtime, args ...*Instance) *Instance {
		return NotImplemented(name, n.Instance)
	})
}
