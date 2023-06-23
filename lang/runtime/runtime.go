package runtime

import "sht/lang/ast"

type Runtime struct {
	// Context *Context
}

func CreateRuntime() *Runtime {
	r := &Runtime{}
	// r.Context = CreateContext()
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
		return r.Eval(stmt)
	}

	return ""
}

func (r *Runtime) EvalNumber(node *ast.Number) string {
	inst := &Instance[NumberImpl]{
		Type: DataType{Name: "Number"},
		Impl: NumberImpl{Value: node.Value},
	}
	return inst.Impl.Repr()
}
