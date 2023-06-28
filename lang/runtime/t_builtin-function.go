package runtime

import (
	"fmt"
	"sht/lang/ast"
)

var builtinFunctionDT = &BuiltinFunctionDataType{
	BaseDataType: BaseDataType{
		Name:        "Function",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
	},
}

var BuiltinFunction = &BuiltinFunctionInfo{
	Type: builtinFunctionDT,
}

// ----------------------------------------------------------------------------
// FUNCTION INFO
// ----------------------------------------------------------------------------
type BuiltinFunctionInfo struct {
	Type DataType
}

func (t *BuiltinFunctionInfo) Create(name string, params []*FunctionParam, fn InternalFunction) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: BuiltinFunctionDataImpl{
			Name:   name,
			Params: params,
			Fn:     fn,
		},
	}
}

// ----------------------------------------------------------------------------
// FUNCTION DATA TYPE
// ----------------------------------------------------------------------------
type BuiltinFunctionDataType struct {
	BaseDataType
}

func (d *BuiltinFunctionDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	name := args[0].Impl.(BuiltinFunctionDataImpl).Name
	return String.Create(fmt.Sprintf("<function:%s>", name))
}

func (d *BuiltinFunctionDataType) OnCall(r *Runtime, s *Scope, args ...*Instance) *Instance {
	impl := args[0].Impl.(BuiltinFunctionDataImpl)
	return impl.Call(r, s, args[1:]...)
}

// ----------------------------------------------------------------------------
// FUNCTION DATA IMPL
// ----------------------------------------------------------------------------
type BuiltinFunctionDataImpl struct {
	Name   string
	Params []*FunctionParam
	Fn     InternalFunction
}

func (d *BuiltinFunctionDataImpl) Call(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(BuiltinFunctionDataImpl)
	return this.Fn(r, s, args...)
}
