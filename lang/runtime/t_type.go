package runtime

import (
	"sht/lang/ast"
)

var Type = &TypeInfo{
	Type: &TypeDataType{
		BaseDataType: BaseDataType{
			Name:        "Type",
			Properties:  map[string]ast.Node{},
			StaticFns:   map[string]Function{},
			InstanceFns: map[string]Function{},
		},
	},
}

// ----------------------------------------------------------------------------
// TYPE INFO
// ----------------------------------------------------------------------------
type TypeInfo struct {
	Type DataType
}

func (t *TypeInfo) Create(dataType DataType) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: TypeDataImpl{
			DataType: dataType,
		},
	}
}

// ----------------------------------------------------------------------------
// TYPE DATA TYPE
// ----------------------------------------------------------------------------
type TypeDataType struct {
	BaseDataType
}

func (d *TypeDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return String.Create(args[0].Type.GetName())
}

// ----------------------------------------------------------------------------
// TYPE DATA IMPL
// ----------------------------------------------------------------------------
type TypeDataImpl struct {
	DataType DataType
}
