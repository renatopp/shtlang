package runtime

import (
	"sht/lang/ast"
)

type Runtime struct {
	Global *Scope
}

func CreateRuntime() *Runtime {
	r := &Runtime{}
	r.Global = CreateScope(nil)
	r.Global.Set(SCOPE_NAME_KEY, Constant(String.Create("Global")))
	r.Global.Set(SCOPE_DEPTH_KEY, Constant(Number.ZERO))
	r.Global.Set(SCOPE_ID_KEY, Constant(String.Create(Id())))

	// r.Global.Set(Type.Type.Name, Type.Instance)
	// r.Global.Set(Number.Type.Name, Number.Instance)
	// r.Global.Set(Boolean.Type.Name, Boolean.Instance)

	return r
}

func (r *Runtime) Run(node ast.Node) string {
	instance := r.Eval(node, nil)
	return instance.Repr()
}

func (r *Runtime) Eval(node ast.Node, scope *Scope) *Instance {
	if scope == nil {
		scope = r.Global
	}

	if node == nil {
		return Boolean.FALSE
	}

	switch n := node.(type) {
	case *ast.Block:
		return r.EvalBlock(n, scope)

	case *ast.Number:
		return r.EvalNumber(n, scope)

	case *ast.Boolean:
		return r.EvalBoolean(n, scope)

	case *ast.String:
		return r.EvalString(n, scope)

	case *ast.UnaryOperator:
		return r.EvalUnaryOperator(n, scope)

	case *ast.BinaryOperator:
		return r.EvalBinaryOperator(n, scope)

	case *ast.Assignment:
		scope.InAssignment = true
		v := r.EvalAssignment(n, scope)
		scope.InAssignment = false
		return v

	case *ast.Identifier:
		return r.EvalIdentifier(n, scope)

	case *ast.FunctionDef:
		return r.EvalFunctionDef(n, scope)

	case *ast.Call:
		return r.EvalCall(n, scope)

	case *ast.Return:
		return r.EvalReturn(n, scope)

	}
	return Boolean.FALSE
}

func (r *Runtime) EvalBlock(node *ast.Block, scope *Scope) *Instance {
	var result *Instance
	for _, stmt := range node.Statements {
		result = r.Eval(stmt, scope)
		if _, ok := scope.Get(RETURN_KEY); ok {
			break
		}
	}

	if result == nil {
		return Boolean.FALSE
	}

	return result
}

func (r *Runtime) EvalNumber(node *ast.Number, scope *Scope) *Instance {
	return Number.Create(node.Value)
}

func (r *Runtime) EvalBoolean(node *ast.Boolean, scope *Scope) *Instance {
	return Boolean.Create(node.Value)
}

func (r *Runtime) EvalString(node *ast.String, scope *Scope) *Instance {
	return String.Create(node.Value)
}

func (r *Runtime) EvalUnaryOperator(node *ast.UnaryOperator, scope *Scope) *Instance {
	right := r.Eval(node.Right, scope)

	switch node.Operator {
	case "+":
		return right.Type.OnPos(r, scope, right)
	case "-":
		return right.Type.OnNeg(r, scope, right)
	case "!":
		return right.Type.OnNot(r, scope, right)
	}

	return nil
}

func (r *Runtime) EvalBinaryOperator(node *ast.BinaryOperator, scope *Scope) *Instance {
	left := r.Eval(node.Left, scope)
	right := r.Eval(node.Right, scope)

	switch node.Operator {
	case "+":
		return left.Type.OnAdd(r, scope, left, right)
	case "-":
		return left.Type.OnSub(r, scope, left, right)
	case "*":
		return left.Type.OnMul(r, scope, left, right)
	case "/":
		return left.Type.OnDiv(r, scope, left, right)
	case "//":
		return left.Type.OnIntDiv(r, scope, left, right)
	case "%":
		return left.Type.OnMod(r, scope, left, right)
	case "**":
		return left.Type.OnPow(r, scope, left, right)
	case "==":
		return left.Type.OnEq(r, scope, left, right)
	case "!=":
		return left.Type.OnNeq(r, scope, left, right)
	case ">":
		return left.Type.OnGt(r, scope, left, right)
	case "<":
		return left.Type.OnLt(r, scope, left, right)
	case ">=":
		return left.Type.OnGte(r, scope, left, right)
	case "<=":
		return left.Type.OnLte(r, scope, left, right)

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

	globalRef, _ := scope.Get(name)
	localRef, _ := scope.GetInScope(name)
	ref := localRef
	if localRef == nil && !node.Definition {
		ref = globalRef
	}

	if !node.Definition && ref == nil {
		return Error.VariableNotDefined(name)
	}

	if !node.Definition && ref != nil && ref.Constant {
		return Error.ReassigningConstant(name)
	}

	exp := r.Eval(node.Expression, scope)
	if exp.Type == Error.Type {
		// TODO: Convert to maybe
	}

	if ref != nil {
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

func (r *Runtime) EvalCall(node *ast.Call, scope *Scope) *Instance {
	target := r.Eval(node.Target, scope)

	args := make([]*Instance, len(node.Arguments)+1)
	args[0] = target
	for i, v := range node.Arguments {
		args[i+1] = r.Eval(v, scope)
	}

	return target.Type.OnCall(r, scope, args...)
}

func (r *Runtime) EvalReturn(node *ast.Return, scope *Scope) *Instance {
	exp := r.Eval(node.Expression, scope)
	scope.Set(RETURN_KEY, &Reference{
		Value:    exp,
		Constant: true,
	})

	return exp
}
