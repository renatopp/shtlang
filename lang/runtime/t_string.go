package runtime

import (
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

var String = _setupString()

type StringInfo struct {
	Instance *Instance
	Type     *DataType

	EMPTY *Instance
}

type StringImpl struct {
	Value string
}

func _setupString() *StringInfo {
	dataType := &DataType{
		Name:        "String",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
		Meta:        map[meta.MetaName]Function{},
	}

	n := &StringInfo{
		Instance: Type.Create(dataType, true),
		Type:     dataType,
	}

	n.EMPTY = n.Create("", true)

	return n
}

// ----------------------------------------------------------------------------
// String Implementation
// ----------------------------------------------------------------------------
func (n StringImpl) Repr() string {
	return n.Value
}

// ----------------------------------------------------------------------------
// String Info
// ----------------------------------------------------------------------------
func (n *StringInfo) Create(value string, constant bool) *Instance {
	return &Instance{
		Type:  n.Type,
		Impl:  StringImpl{Value: value},
		Const: constant,
	}
}
