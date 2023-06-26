package runtime

import (
	"sht/lang/ast"
)

var stringDT = &StringDataType{
	BaseDataType: BaseDataType{
		Name:        "String",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
	},
}

var String = &StringInfo{
	Type: stringDT,

	EMPTY: &Instance{
		Type: stringDT,
		Impl: StringDataImpl{
			Value: "",
		},
	},
}

// ----------------------------------------------------------------------------
// STRING INFO
// ----------------------------------------------------------------------------
type StringInfo struct {
	Type DataType

	EMPTY *Instance
}

func (t *StringInfo) Create(value string) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: StringDataImpl{
			Value: value,
		},
	}
}

// ----------------------------------------------------------------------------
// STRING DATA TYPE
// ----------------------------------------------------------------------------
type StringDataType struct {
	BaseDataType
}

func (d *StringDataType) OnRepr(r *Runtime, args ...*Instance) *Instance {
	return args[0]
}

// ----------------------------------------------------------------------------
// STRING DATA IMPL
// ----------------------------------------------------------------------------
type StringDataImpl struct {
	Value string
}
