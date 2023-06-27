package runtime

import (
	"sht/lang/ast"
)

const RETURN_KEY = "0_return"

type Runtime struct {
	Global *Scope
	Stack  *Stack
}

func CreateRuntime() *Runtime {
	r := &Runtime{}
	r.Global = CreateScope(nil)
	// r.Global.Set(Type.Type.Name, Type.Instance)
	// r.Global.Set(Number.Type.Name, Number.Instance)
	// r.Global.Set(Boolean.Type.Name, Boolean.Instance)

	r.Stack = NewStack(r.Global)
	return r
}

func (r *Runtime) Run(node ast.Node) string {
	instance := r.Eval(node, nil)
	return instance.Repr()
}

func (r *Runtime) Eval(node ast.Node, scope *Scope) *Instance {
	if scope == nil {
		scope = r.Stack.Current()
	}

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

	case *ast.Assignment:
		scope.InAssignment = true
		v := r.EvalAssignment(n, scope)
		scope.InAssignment = false
		return v

	case *ast.Identifier:
		return r.EvalIdentifier(n, scope)

	case *ast.FunctionDef:
		return r.EvalFunctionDef(n, scope)
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
	return Number.Create(node.Value)
}

func (r *Runtime) EvalBoolean(node *ast.Boolean) *Instance {
	return Boolean.Create(node.Value)
}

func (r *Runtime) EvalString(node *ast.String) *Instance {
	return String.Create(node.Value)
}

func (r *Runtime) EvalUnaryOperator(node *ast.UnaryOperator) *Instance {
	right := r.Eval(node.Right, nil)

	switch node.Operator {
	case "+":
		return right.Type.OnPos(r, right)
	case "-":
		return right.Type.OnNeg(r, right)
	case "!":
		return right.Type.OnNot(r, right)
	}

	return nil
}

func (r *Runtime) EvalBinaryOperator(node *ast.BinaryOperator) *Instance {
	left := r.Eval(node.Left, nil)
	right := r.Eval(node.Right, nil)

	switch node.Operator {
	case "+":
		return left.Type.OnAdd(r, left, right)
	case "-":
		return left.Type.OnSub(r, left, right)
	case "*":
		return left.Type.OnMul(r, left, right)
	case "/":
		return left.Type.OnDiv(r, left, right)
	case "//":
		return left.Type.OnIntDiv(r, left, right)
	case "%":
		return left.Type.OnMod(r, left, right)
	case "**":
		return left.Type.OnPow(r, left, right)
	case "==":
		return left.Type.OnEq(r, left, right)
	case "!=":
		return left.Type.OnNeq(r, left, right)
	case ">":
		return left.Type.OnGt(r, left, right)
	case "<":
		return left.Type.OnLt(r, left, right)
	case ">=":
		return left.Type.OnGte(r, left, right)
	case "<=":
		return left.Type.OnLte(r, left, right)

	case "and", "or", "nand", "nor", "xor", "nxor":
		lt := AsBool(left)
		rt := AsBool(right)

		switch node.Operator {
		case "and":
			return Boolean.Create(lt && rt)
		case "or":
			return Boolean.Create(lt || rt)
		case "nand":
			return Boolean.Create(!(lt && rt))
		case "nor":
			return Boolean.Create(!(lt || rt))
		case "xor":
			return Boolean.Create(lt != rt)
		case "nxor":
			return Boolean.Create(lt == rt)
		}

	case "..":
		lt := AsString(left)
		rt := AsString(right)
		return String.Create(lt + rt)
	}

	return nil
}

func (r *Runtime) EvalAssignment(node *ast.Assignment, scope *Scope) *Instance {
	name := node.Identifier.(*ast.Identifier).Value
	if node.Definition && scope.HasInScope(name) {
		return Error.DuplicatedDefinition(name)
	}

	ref, ok := scope.Get(name)
	if !node.Definition && !ok {
		return Error.VariableNotDefined(name)
	}

	if !node.Definition && ok && ref.Constant {
		return Error.ReassigningConstant(name)
	}

	exp := r.Eval(node.Expression, scope)
	if exp.Type == Error.Type {
		// TODO: Convert to maybe
	}

	if ok {
		ref.Value = exp
	} else {
		scope.Set(name, &Reference{
			Value:    exp,
			Constant: node.Constant,
		})
	}

	return exp
}

func (r *Runtime) EvalIdentifier(node *ast.Identifier, scope *Scope) *Instance {
	name := node.Value
	ref, ok := scope.Get(name)
	if !ok {
		return Error.VariableNotDefined(name)
	}

	return ref.Value
}

func (r *Runtime) EvalFunctionDef(node *ast.FunctionDef, scope *Scope) *Instance {
	name := node.Name
	params := make([]*FunctionParam, len(node.Params))

	hasSpread := false
	hasDefault := false
	for i, param := range node.Params {
		param := param.(*ast.Parameter)
		p := &FunctionParam{
			Name:    param.Name,
			Spread:  param.Spread,
			Default: r.Eval(param.Default, nil),
		}

		if p.Spread && p.Default != nil {
			return Error.Create("spread arguments cannot have default values: '%s'", p.Name)
		}

		if p.Spread {
			if hasSpread {
				return Error.Create("only one spread argument is allowed: '%s'", p.Name)
			}

			hasSpread = true
		}

		if p.Default != nil {
			if hasSpread {
				return Error.Create("default arguments cannot proceed spread arguments: '%s'", p.Name)
			}
			hasDefault = true
		} else if hasDefault {
			return Error.Create("default arguments must be at the end: '%s'", p.Name)
		}

		params[i] = p
	}

	fn := CustomFunction.Create(name, params, node.Body, scope)

	if !scope.InAssignment && !scope.InArgument && name != "" {
		scope.Set(name, &Reference{
			Value:    fn,
			Constant: false,
		})
	}

	return fn
}
