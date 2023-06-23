package runtime

import (
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

var ZERO = CreateNumber(0, true)
var ONE = CreateNumber(1, true)
var TRUE = CreateBoolean(true, true)
var FALSE = CreateBoolean(false, true)

const RETURN_KEY = "0_return"

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

	SetupNumber()

	r.Stack = NewStack(r.Global)
	return r
}

func (r *Runtime) Run(node ast.Node) string {
	instance := r.Eval(node, nil)
	return instance.Impl.Repr()
}

func (r *Runtime) Eval(node ast.Node, scope *Scope) *Instance {
	switch n := node.(type) {
	case *ast.Block:
		return r.EvalBlock(n)

	case *ast.Number:
		return r.EvalNumber(n)

	case *ast.Boolean:
		return r.EvalBoolean(n)

		// case *ast.UnaryOperator:
		// 	return r.EvalUnaryOperator(n, scope)

	case *ast.BinaryOperator:
		return r.EvalBinaryOperator(n)
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

// func (r *Runtime) EvalUnaryOperator(node *ast.UnaryOperator, scope *Scope) *Instance {

// }

func (r *Runtime) EvalBinaryOperator(node *ast.BinaryOperator) *Instance {

	left := r.Eval(node.Left, nil)
	right := r.Eval(node.Right, nil)

	return left.Type.Meta[meta.Add].Call(r, []*Instance{left, right})
}
