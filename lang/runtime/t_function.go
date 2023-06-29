package runtime

import (
	"fmt"
	"sht/lang/ast"
)

var functionDT = &FunctionDataType{
	BaseDataType: BaseDataType{
		Name:        "Function",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Callable{},
		InstanceFns: map[string]Callable{},
	},
}

var Function = &FunctionInfo{
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
			ParentScope: scope,
			Name:        name,
			Params:      params,
			Body:        body,
		},
	}
}

func (t *FunctionInfo) CreateNative(name string, params []*FunctionParam, fn MetaFunction) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: FunctionDataImpl{
			Name:     name,
			Params:   params,
			NativeFn: fn,
		},
	}
}

// ----------------------------------------------------------------------------
// FUNCTION DATA TYPE
// ----------------------------------------------------------------------------
type FunctionDataType struct {
	BaseDataType
}

func (d *FunctionDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	name := args[0].Impl.(FunctionDataImpl).Name
	return String.Create(fmt.Sprintf("<function:%s>", name))
}

func (d *FunctionDataType) OnCall(r *Runtime, s *Scope, args ...*Instance) *Instance {
	impl := args[0].Impl.(FunctionDataImpl)
	return impl.Call(r, s, args[1:]...)
}

// ----------------------------------------------------------------------------
// FUNCTION DATA IMPL
// ----------------------------------------------------------------------------
type FunctionDataImpl struct {
	ParentScope *Scope
	Name        string
	Params      []*FunctionParam
	Body        ast.Node
	NativeFn    MetaFunction
}

type FunctionParam struct {
	Name    string
	Default *Instance
	Spread  bool
}

func (d *FunctionDataImpl) Call(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if d.NativeFn != nil {
		return d.NativeFn(r, s, args...)
	}

	parentScope := d.ParentScope
	if parentScope == nil {
		parentScope = s
	}

	depth, _ := parentScope.GetInScope(SCOPE_DEPTH_KEY)

	scope := CreateScope(parentScope, s)
	scope.Set(SCOPE_NAME_KEY, Constant(String.Create(d.Name)))
	scope.Set(SCOPE_DEPTH_KEY, Constant(Number.Create(AsNumber(depth.Value)+1)))
	scope.Set(SCOPE_ID_KEY, Constant(String.Create(Id())))
	scope.Set(SCOPE_FN_KEY, Constant(&Instance{
		Type: Function.Type,
		Impl: d,
	}))

	tArgs := len(args)
	// tParams := len(d.Params)
	g := 0
	for i, v := range d.Params {
		if !v.Spread {
			var value *Instance
			if g >= tArgs {
				if v.Default == nil {
					return r.Throw(Error.Create(scope, "missing argument: '%s'", v.Name), scope)
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
			return r.Throw(Error.Create(scope, "spread arguments are not supported yet: '%s'", v.Name), scope)

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

	res := r.Eval(d.Body, scope)

	if scope.HasInScope(RAISE_KEY) {
		err, _ := scope.GetInScope(RAISE_KEY)
		s.Set(RAISE_KEY, err)
	}

	return res
}
