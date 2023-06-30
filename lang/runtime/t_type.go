package runtime

import (
	"sht/lang/ast"
)

var Type = &TypeInfo{
	Type: &TypeDataType{
		BaseDataType: BaseDataType{
			Name:        "Type",
			Properties:  map[string]ast.Node{},
			StaticFns:   map[string]*Instance{},
			InstanceFns: map[string]*Instance{},
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
		Impl: &TypeDataImpl{
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

// func (d *TypeDataType) OnCall(r *Runtime, s *Scope, def map[string]*Instance, args ...*Instance) *Instance {
// 	impl := args[0].Impl.(*TypeDataImpl)
// 	instance := impl.DataType.Instantiate(r, s, def)
// 	instance = impl.DataType.OnNew(r, s, append([]*Instance{instance}, args[1:]...)...)
// 	return instance
// }

func (d *TypeDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return String.Createf("<Type:%s>", args[0].Type.GetName())
}

// ----------------------------------------------------------------------------
// TYPE DATA IMPL
// ----------------------------------------------------------------------------
type TypeDataImpl struct {
	DataType DataType
}
