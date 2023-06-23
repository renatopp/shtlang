package runtime

import "sht/lang/ast"

type FunctionParam struct {
	Name    string
	Default Instance
	Spread  bool
}

type FunctionImpl struct {
	Scope  *Scope
	Params []*FunctionParam
	Body   ast.Node
}

// func (f *FunctionImpl) Call(r *Runtime, args []Instance) Instance {
// 	scope := CreateScope(f.Scope)
// 	for i, param := range f.Params {
// 		if i < len(args) {
// 			scope.Set(param.Name, args[i])
// 		} else {
// 			scope.Set(param.Name, param.Default)
// 		}
// 	}

// 	return r.Eval(f.Body, scope)
// }

func (f *FunctionImpl) Repr() string {
	return "<function>"
}
