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
	Type         DataType
	TypeInstance *Instance
}

func (t *TypeInfo) Create(dataType DataType) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &TypeDataImpl{
			DataType: dataType,
		},
	}
}

func (t *TypeInfo) Setup() {
	t.TypeInstance = Type.Create(Type.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
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

func (d *TypeDataType) OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*TypeDataImpl)
	return this.DataType.OnTo(r, s, args[0], args[1:]...)
}

func (d *TypeDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return String.Createf("<Type:%s>", self.Type.GetName())
}

// ----------------------------------------------------------------------------
// TYPE DATA IMPL
// ----------------------------------------------------------------------------
type TypeDataImpl struct {
	DataType     DataType
	TypeInstance *Instance
}
