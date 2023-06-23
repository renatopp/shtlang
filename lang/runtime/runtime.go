package runtime

import "sht/lang/ast"

var ZERO = CreateNumber(0, true)
var ONE = CreateNumber(1, true)
var TRUE = CreateBoolean(true, true)
var FALSE = CreateBoolean(false, true)

type Runtime struct {
	Global *Scope
	Stack  *Stack
}

func CreateRuntime() *Runtime {
	r := &Runtime{}
	r.Global = CreateScope(nil)
	r.Global.Set(Type.Type.Name, Type)
	r.Global.Set(Number.Type.Name, Number)
	r.Global.Set(Boolean.Type.Name, Boolean)

	r.Stack = NewStack(r.Global)
	return r
}

func (r *Runtime) Eval(node ast.Node, scope *Scope) *Instance {
	switch n := node.(type) {
	case *ast.Block:
		return r.EvalBlock(n)

	case *ast.Number:
		return r.EvalNumber(n)

	case *ast.Boolean:
		return r.EvalBoolean(n)
	}

	return FALSE
}

func (r *Runtime) EvalBlock(node *ast.Block) *Instance {
	var result *Instance
	for _, stmt := range node.Statements {
		result = r.Eval(stmt, nil)
	}

	if result == nil {
		return FALSE
	}
	return result
}

func (r *Runtime) EvalNumber(node *ast.Number) *Instance {
	return CreateNumber(node.Value, false)
}

func (r *Runtime) EvalBoolean(node *ast.Boolean) *Instance {
	return CreateBoolean(node.Value, false)
}
