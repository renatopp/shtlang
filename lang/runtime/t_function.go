package runtime

import (
	"fmt"
	"sht/lang/ast"
)

var functionDT = &FunctionDataType{
	BaseDataType: BaseDataType{
		Name:        "Function",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
	},
}

var CustomFunction = &FunctionInfo{
	Type: functionDT,
}

// ----------------------------------------------------------------------------
// FUNCTION INFO
// ----------------------------------------------------------------------------
type FunctionInfo struct {
	Type DataType
}

func (t *FunctionInfo) Create(name string, params []*FunctionParam, body ast.Node, scope *Scope) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: FunctionDataImpl{
			Scope:  scope,
			Name:   name,
			Params: params,
			Body:   body,
		},
	}
}

// ----------------------------------------------------------------------------
// FUNCTION DATA TYPE
// ----------------------------------------------------------------------------
type FunctionDataType struct {
	BaseDataType
}

func (d *FunctionDataType) OnRepr(r *Runtime, args ...*Instance) *Instance {
	name := args[0].Impl.(FunctionDataImpl).Name
	return String.Create(fmt.Sprintf("<function:%s>", name))
}

// ----------------------------------------------------------------------------
// FUNCTION DATA IMPL
// ----------------------------------------------------------------------------
type FunctionDataImpl struct {
	Scope  *Scope
	Name   string
	Params []*FunctionParam
	Body   ast.Node
}

type FunctionParam struct {
	Name    string
	Default *Instance
	Spread  bool
}
