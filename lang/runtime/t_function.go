package runtime

import (
	"sht/lang/ast"
)

var functionDT = &FunctionDataType{
	BaseDataType: BaseDataType{
		Name:        "Function",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
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
		Impl: &FunctionDataImpl{
			ParentScope: scope,
			Name:        name,
			Params:      params,
			Body:        body,
			Generator:   false,
		},
	}
}

func (t *FunctionInfo) CreateNative(name string, params []*FunctionParam, fn MetaFunction) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &FunctionDataImpl{
			Name:      name,
			Params:    params,
			NativeFn:  fn,
			Generator: false,
		},
	}
}

// ----------------------------------------------------------------------------
// FUNCTION DATA TYPE
// ----------------------------------------------------------------------------
type FunctionDataType struct {
	BaseDataType
}

func (d *FunctionDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args[0])
}

func (d *FunctionDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	name := args[0].Impl.(*FunctionDataImpl).Name
	return String.Createf("<Function:%s>", name)
}

func (d *FunctionDataType) OnCall(r *Runtime, s *Scope, args ...*Instance) *Instance {
	impl := args[0].Impl.(*FunctionDataImpl)
	return impl.Call(r, s, args...)
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
	Generator   bool
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

	scope := CreateScope(parentScope, s)
	scope.Set(SCOPE_NAME_KEY, Constant(String.Create(d.Name)))
	scope.Set(SCOPE_FN_KEY, Constant(&Instance{
		Type: Function.Type,
		Impl: d,
	}))

	arguments := []*Instance{}
	paramsLength := len(d.Params)
	argsLength := len(args)

	j := 0
	for i, pv := range d.Params {
		if pv.Spread {
			spreadAmount := (argsLength - i) - (paramsLength - i - 1)
			spreadItems := []*Instance{}

			for k := 0; k < spreadAmount; k++ {
				spreadItems = append(spreadItems, args[j])
				j++
			}

			rv := List.Create(spreadItems...)
			arguments = append(arguments, rv)

		} else {
			if j >= argsLength {
				if pv.Default == nil {
					return r.Throw(Error.Create(scope, "missing arguments for parameter '%s'", pv.Name), scope)
				} else {
					arguments = append(arguments, pv.Default)
				}
			} else {
				arguments = append(arguments, args[j])
				j++
			}
		}
	}

	for i, pv := range d.Params {
		if pv.Name != "_" {
			scope.Set(pv.Name, Variable(arguments[i]))
		}
	}

	if d.Generator {
		iter := Iterator.Create(Function.CreateNative("generator", []*FunctionParam{}, func(r *Runtime, s *Scope, args ...*Instance) *Instance {
			res := r.Eval(d.Body, scope)

			if err, ok := scope.GetInScope(RAISE_KEY); ok {
				s.Set(RAISE_KEY, err)
				return Iteration.Error(err.Value)

			} else if _, ok := scope.GetInScope(YIELD_KEY); ok {
				scope.Delete(YIELD_KEY)
				scope.Delete(JUST_YIELDED_KEY)
				return Iteration.Create(res)

			} else {
				return Iteration.DONE
			}
		}))

		return iter

	} else {
		res := r.Eval(d.Body, scope)

		if scope.HasInScope(RAISE_KEY) {
			err, _ := scope.GetInScope(RAISE_KEY)
			s.Set(RAISE_KEY, err)
		}

		return res
	}
}
