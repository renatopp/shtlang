package runtime

import "sht/lang/ast"

type Runtime struct {
	Global *Scope
	Stack  *Stack
}

func CreateRuntime() *Runtime {
	r := &Runtime{}
	r.Global = CreateScope(nil)
	r.Stack = NewStack(r.Global)
	return r
}

func (r *Runtime) Eval(node ast.Node, scope *Scope) string {
	switch n := node.(type) {
	case *ast.Block:
		return r.EvalBlock(n)

	case *ast.Number:
		return r.EvalNumber(n)
	}

	return ""
}

func (r *Runtime) EvalBlock(node *ast.Block) string {
	// var result obj.Object
	for _, stmt := range node.Statements {
		return r.Eval(stmt, nil)
	}

	return ""
}

func (r *Runtime) EvalNumber(node *ast.Number) string {
	inst := &Instance{
		Type: &DataType{Name: "Number"},
		Impl: NumberImpl{Value: node.Value},
	}
	return inst.Impl.Repr()
}
