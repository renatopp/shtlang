package runtime

import (
	"fmt"
	"sht/lang/ast"
)

var functionDT = &CustomFunctionDataType{
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
		Impl: CustomFunctionDataImpl{
			ParentScope: scope,
			Name:        name,
			Params:      params,
			Body:        body,
		},
	}
}

// ----------------------------------------------------------------------------
// FUNCTION DATA TYPE
// ----------------------------------------------------------------------------
type CustomFunctionDataType struct {
	BaseDataType
}

func (d *CustomFunctionDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	name := args[0].Impl.(CustomFunctionDataImpl).Name
	return String.Create(fmt.Sprintf("<function:%s>", name))
}

func (d *CustomFunctionDataType) OnCall(r *Runtime, s *Scope, args ...*Instance) *Instance {
	impl := args[0].Impl.(CustomFunctionDataImpl)
	return impl.Call(r, s, args[1:]...)
}

// ----------------------------------------------------------------------------
// FUNCTION DATA IMPL
// ----------------------------------------------------------------------------
type CustomFunctionDataImpl struct {
	ParentScope *Scope
	Name        string
	Params      []*FunctionParam
	Body        ast.Node
}

type FunctionParam struct {
	Name    string
	Default *Instance
	Spread  bool
}

func (d *CustomFunctionDataImpl) Call(r *Runtime, s *Scope, args ...*Instance) *Instance {
	parentScope := d.ParentScope
	if parentScope == nil {
		parentScope = s
	}

	depth, _ := parentScope.GetInScope(SCOPE_DEPTH_KEY)

	scope := CreateScope(parentScope)
	scope.Set(SCOPE_NAME_KEY, Constant(String.Create(d.Name)))
	scope.Set(SCOPE_DEPTH_KEY, Constant(Number.Create(AsNumber(depth.Value)+1)))
	scope.Set(SCOPE_ID_KEY, Constant(String.Create(Id())))

	tArgs := len(args)
	// tParams := len(d.Params)
	g := 0
	for i, v := range d.Params {
		if !v.Spread {
			var value *Instance
			if g >= tArgs {
				if v.Default == nil {
					return Error.Create("missing argument: '%s'", v.Name)
				}
				value = v.Default
			} else {
				value = args[g]
			}

			g++
			scope.Set(v.Name, &Reference{
				Value:    value,
				Constant: false,
			})
		} else {
			// TODO: Handle spread arguments
			return Error.Create("spread arguments are not supported yet: '%s'", v.Name)

			if i == 0 {
			} // TODO: REMOVE

			// missing := tParams - i - 1

			// total := 0
			// sv := make([]*Instance, 0)
			// for j := i; j < (tArgs - missing); j++ {
			// 	t := args[j]
			// 	if t.Type() == obj.TList {
			// 		sv = append(sv, t.(*obj.List).Values...)
			// 	} else {
			// 		sv = append(sv, t)
			// 	}
			// 	total++
			// }

			// g += total
			// scope.Set(v.Name, &obj.List{Values: sv})
		}
	}

	return r.Eval(d.Body, scope)
}
