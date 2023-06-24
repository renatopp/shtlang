package runtime

import (
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

const RETURN_KEY = "0_return"

type Runtime struct {
	Global *Scope
	Stack  *Stack
}

func CreateRuntime() *Runtime {
	r := &Runtime{}
	r.Global = CreateScope(nil)
	r.Global.Set(Type.Type.Name, Type.Instance)
	r.Global.Set(Number.Type.Name, Number.Instance)
	r.Global.Set(Boolean.Type.Name, Boolean.Instance)

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

	case *ast.String:
		return r.EvalString(n)

	case *ast.UnaryOperator:
		return r.EvalUnaryOperator(n)

	case *ast.BinaryOperator:
		return r.EvalBinaryOperator(n)
	}

	return Boolean.FALSE
}

func (r *Runtime) EvalBlock(node *ast.Block) *Instance {
	var result *Instance
	for _, stmt := range node.Statements {
		result = r.Eval(stmt, nil)
	}

	if result == nil {
		return Boolean.FALSE
	}
	return result
}

func (r *Runtime) EvalNumber(node *ast.Number) *Instance {
	return Number.Create(node.Value, false)
}

func (r *Runtime) EvalBoolean(node *ast.Boolean) *Instance {
	return Boolean.Create(node.Value, false)
}

func (r *Runtime) EvalString(node *ast.String) *Instance {
	return String.Create(node.Value, false)
}

func (r *Runtime) EvalUnaryOperator(node *ast.UnaryOperator) *Instance {
	right := r.Eval(node.Right, nil)

	m := right.Type.Meta
	switch node.Operator {
	case "+":
		return m[meta.Pos].Call(r, right)
	case "-":
		return m[meta.Neg].Call(r, right)
	case "!":
		return m[meta.Not].Call(r, right)
	}

	return nil
}

func (r *Runtime) EvalBinaryOperator(node *ast.BinaryOperator) *Instance {
	left := r.Eval(node.Left, nil)
	right := r.Eval(node.Right, nil)

	m := left.Type.Meta
	switch node.Operator {
	case "+", "-", "*", "/", "//", "%", "**", "==", "!=", ">", "<", ">=", "<=":
		return m[meta.FromBinaryOperator(node.Operator)].Call(r, left, right)

	case "and", "or", "nand", "nor", "xor", "nxor":
		lt, rt := false, false
		if !IsBool(left) {
			lt = m[meta.Boolean].Call(r, left).Impl.(BooleanImpl).Value
		} else {
			lt = left.Impl.(BooleanImpl).Value
		}

		if !IsBool(right) {
			rt = m[meta.Boolean].Call(r, right).Impl.(BooleanImpl).Value
		} else {
			rt = right.Impl.(BooleanImpl).Value
		}

		switch node.Operator {
		case "and":
			return Boolean.Create(lt && rt, false)
		case "or":
			return Boolean.Create(lt || rt, false)
		case "nand":
			return Boolean.Create(!(lt && rt), false)
		case "nor":
			return Boolean.Create(!(lt || rt), false)
		case "xor":
			return Boolean.Create(lt != rt, false)
		case "nxor":
			return Boolean.Create(lt == rt, false)
		}

	case "..":
		lt := left.Type.Meta[meta.String].Call(r, left).Impl.(StringImpl).Value
		rt := left.Type.Meta[meta.String].Call(r, right).Impl.(StringImpl).Value
		return String.Create(lt+rt, false)
	}

	return nil
}

func IsBool(instance *Instance) bool {
	return instance.Type == Boolean.Type
}
