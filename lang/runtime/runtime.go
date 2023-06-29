package runtime

import (
	"sht/lang/ast"
)

type Runtime struct {
	Global *Scope
}

func CreateRuntime() *Runtime {
	r := &Runtime{}

	Iterator.Setup()

	r.Global = CreateScope(nil, nil)
	r.Global.Set(SCOPE_NAME_KEY, Constant(String.Create("Global")))
	r.Global.Set(SCOPE_DEPTH_KEY, Constant(Number.ZERO))
	r.Global.Set(SCOPE_ID_KEY, Constant(String.Create(Id())))

	r.Global.Set(Boolean.Type.GetName(), Constant(Type.Create(Boolean.Type)))
	r.Global.Set(Error.Type.GetName(), Constant(Type.Create(Error.Type)))
	r.Global.Set(Iteration.Type.GetName(), Constant(Type.Create(Iteration.Type)))
	r.Global.Set(Iterator.Type.GetName(), Constant(Type.Create(Iterator.Type)))
	r.Global.Set(Function.Type.GetName(), Constant(Type.Create(Function.Type)))
	r.Global.Set(List.Type.GetName(), Constant(Type.Create(List.Type)))
	r.Global.Set(Maybe.Type.GetName(), Constant(Type.Create(Maybe.Type)))
	r.Global.Set(Number.Type.GetName(), Constant(Type.Create(Number.Type)))
	r.Global.Set(String.Type.GetName(), Constant(Type.Create(String.Type)))
	r.Global.Set(Tuple.Type.GetName(), Constant(Type.Create(Tuple.Type)))
	r.Global.Set(Type.Type.GetName(), Constant(Type.Create(Type.Type)))

	r.Global.Set("Done", Constant(Iteration.DONE))

	return r
}

func (r *Runtime) Run(node ast.Node) string {
	instance := r.Eval(node, nil)

	delete(r.Global.Values, RAISE_KEY)
	return instance.Repr()
}

func (r *Runtime) Eval(node ast.Node, scope *Scope) *Instance {
	if scope == nil {
		scope = r.Global
	}

	if node == nil {
		return Boolean.FALSE
	}

	scope.PushNode(node)
	var result *Instance
	switch n := node.(type) {
	case *ast.Block:
		result = r.EvalBlock(n, scope)

	case *ast.Number:
		result = r.EvalNumber(n, scope)

	case *ast.Boolean:
		result = r.EvalBoolean(n, scope)

	case *ast.String:
		result = r.EvalString(n, scope)

	case *ast.Tuple:
		result = r.EvalTuple(n, scope)

	case *ast.UnaryOperator:
		result = r.EvalUnaryOperator(n, scope)

	case *ast.BinaryOperator:
		result = r.EvalBinaryOperator(n, scope)

	case *ast.Assignment:
		scope.InAssignment = true
		result = r.EvalAssignment(n, scope)
		scope.InAssignment = false

	case *ast.Identifier:
		result = r.EvalIdentifier(n, scope)

	case *ast.FunctionDef:
		result = r.EvalFunctionDef(n, scope)

	case *ast.Call:
		result = r.EvalCall(n, scope)

	case *ast.Return:
		result = r.EvalReturn(n, scope)

	case *ast.Raise:
		result = r.EvalRaise(n, scope)

	case *ast.Indexing:
		result = r.EvalIndexing(n, scope)

	case *ast.Wrapping:
		result = r.EvalWrapping(n, scope)

	case *ast.Unwrapping:
		result = r.EvalUnwrap(n, scope)

	case *ast.If:
		result = r.EvalIf(n, scope)

	case *ast.SpreadOut:
		result = r.EvalSpreadOut(n, scope)

	}

	scope.PopNode()

	if result != nil {
		return result
	}

	return Boolean.FALSE
}

func (r *Runtime) Throw(err *Instance, scope *Scope) *Instance {
	e, ok := scope.GetInScope(RAISE_KEY)
	if ok {
		return e.Value
	}

	scope.Set(RAISE_KEY, &Reference{
		Value:    err,
		Constant: true,
	})

	return err
}

func (r *Runtime) EvalBlock(node *ast.Block, scope *Scope) *Instance {
	var result *Instance
	for _, stmt := range node.Statements {
		result = r.Eval(stmt, scope)
		if err, ok := scope.Get(RAISE_KEY); ok {
			result = err.Value
			break
		}

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

func (r *Runtime) EvalTuple(node *ast.Tuple, scope *Scope) *Instance {
	values := make([]*Instance, 0)

	for _, v := range node.Values {
		s, isSpread := v.(*ast.SpreadOut)

		if isSpread {
			spread := r.Eval(s.Target, scope)
			var e *Instance
			r.ResolveIterator(spread, scope, func(v *Instance, err *Instance) {
				if err != nil {
					e = err
				} else if v != nil {
					t := v.Impl.(*TupleDataImpl)
					values = append(values, t.Values...)
				}
			})

			if e != nil {
				return e
			}

		} else {
			values = append(values, r.Eval(v, scope))
		}
	}

	return Tuple.Create(values...)
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

	switch node.Operator {
	case "+":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnAdd(r, scope, left, right)
	case "-":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnSub(r, scope, left, right)
	case "*":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnMul(r, scope, left, right)
	case "/":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnDiv(r, scope, left, right)
	case "//":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnIntDiv(r, scope, left, right)
	case "%":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnMod(r, scope, left, right)
	case "**":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnPow(r, scope, left, right)
	case "==":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnEq(r, scope, left, right)
	case "!=":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnNeq(r, scope, left, right)
	case ">":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnGt(r, scope, left, right)
	case "<":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnLt(r, scope, left, right)
	case ">=":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnGte(r, scope, left, right)
	case "<=":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		return left.Type.OnLte(r, scope, left, right)

	case "and", "or", "nand", "nor", "xor", "nxor":
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		lt := AsBool(left)
		rt := AsBool(right)

		// TODO: only evaluating if necessary
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
		left := r.Eval(node.Left, scope)
		right := r.Eval(node.Right, scope)
		lt := AsString(left)
		rt := AsString(right)
		return String.Create(lt + rt)

	case "??":
		left := r.Eval(node.Left, scope)
		if left.Type != Maybe.Type && left.Type != Error.Type {
			return left
		}

		if left.Type == Error.Type {
			return r.Eval(node.Right, scope)
		}

		maybe := left.Impl.(*MaybeDataImpl)
		if maybe.Error != nil {
			return r.Eval(node.Right, scope)
		}

		r.SolveMaybe(left, scope)
		return maybe.Value
	}

	return nil
}

func (r *Runtime) EvalAssignment(node *ast.Assignment, scope *Scope) *Instance {
	// idpool := []ast.Node{}
	// expool := []*Instance{}

	identifiers := node.Identifier.(*ast.Tuple)
	values := r.Eval(node.Expression, scope)

	if len(identifiers.Values) > 1 {
		if values.Type != Tuple.Type {
			return r.Throw(Error.Create(scope, "cannot unpack non-tuple type"), scope)
		}
		tupledValues := values.Impl.(*TupleDataImpl)
		if len(identifiers.Values) != len(tupledValues.Values) {
			return r.Throw(Error.Create(scope, "cannot unpack tuple of length %d into %d variables",
				len(tupledValues.Values), len(identifiers.Values)), scope)
		}

		for i, _ := range identifiers.Values {
			id := identifiers.Values[i].(*ast.Identifier).Value
			exp := tupledValues.Values[i]
			r.Assign(id, exp, node.Definition, node.Constant, scope)
		}

		return values
	}

	name := identifiers.Values[0].(*ast.Identifier).Value
	return r.Assign(name, values, node.Definition, node.Constant, scope)
}

func (r *Runtime) Assign(name string, exp *Instance, def, cnst bool, scope *Scope) *Instance {
	if def && scope.HasInScope(name) {
		return r.Throw(Error.DuplicatedDefinition(scope, name), scope)
	}

	globalRef, _ := scope.Get(name)
	localRef, _ := scope.GetInScope(name)
	ref := localRef
	if localRef == nil && !def {
		ref = globalRef
	}

	if !def && ref == nil {
		return r.Throw(Error.VariableNotDefined(scope, name), scope)
	}

	if !def && ref != nil && ref.Constant {
		return r.Throw(Error.ReassigningConstant(scope, name), scope)
	}

	if ref != nil {
		ref.Value = exp
	} else {
		scope.Set(name, &Reference{
			Value:    exp,
			Constant: cnst,
		})
	}

	return exp
}

func (r *Runtime) EvalIdentifier(node *ast.Identifier, scope *Scope) *Instance {
	name := node.Value
	ref, ok := scope.Get(name)
	if !ok {
		return r.Throw(Error.VariableNotDefined(scope, name), scope)
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
			return r.Throw(
				Error.Create(scope, "spread arguments cannot have default values: '%s'", p.Name),
				scope,
			)
		}

		if p.Spread {
			if hasSpread {
				return r.Throw(
					Error.Create(scope, "only one spread argument is allowed: '%s'", p.Name),
					scope,
				)
			}

			hasSpread = true
		}

		if p.Default != nil {
			if hasSpread {
				return r.Throw(
					Error.Create(scope, "default arguments cannot proceed spread arguments: '%s'", p.Name),
					scope,
				)
			}
			hasDefault = true
		} else if hasDefault {
			return r.Throw(
				Error.Create(scope, "default arguments must be at the end: '%s'", p.Name),
				scope,
			)
		}

		params[i] = p
	}

	fn := Function.Create(name, params, node.Body, scope)

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

	isType := target.Type == Type.Type
	if !isType && node.Initializer != nil {
		return r.Throw(Error.Create(scope, "cannot initialize non-type"), scope)
	}

	args := []*Instance{target}
	for _, v := range node.Arguments {
		if spread, ok := v.(*ast.SpreadOut); ok {
			var e *Instance
			target := r.Eval(spread.Target, scope)
			r.ResolveIterator(target, scope, func(v *Instance, err *Instance) {
				if err != nil {
					e = err
				} else if v != nil {
					t := v.Impl.(*TupleDataImpl)
					args = append(args, t.Values...)
				}
			})
			if e != nil {
				return e
			}
			continue
		}

		args = append(args, r.Eval(v, scope))
	}

	if isType {
		impl := target.Impl.(*TypeDataImpl)
		value := impl.DataType.Instantiate(r, scope, node.Initializer)
		return value.Type.OnNew(r, scope, append([]*Instance{value}, args[1:]...)...)

	} else {
		return target.Type.OnCall(r, scope, args...)
	}
}

func (r *Runtime) EvalReturn(node *ast.Return, scope *Scope) *Instance {
	exp := r.Eval(node.Expression, scope)
	if exp == nil {
		exp = Boolean.FALSE
	}

	scope.Set(RETURN_KEY, &Reference{
		Value:    exp,
		Constant: true,
	})

	return exp
}

func (r *Runtime) EvalRaise(node *ast.Raise, scope *Scope) *Instance {
	exp := r.Eval(node.Expression, scope)
	if exp == nil {
		exp = Boolean.FALSE
	}

	if exp.Type == Error.Type {
		return r.Throw(exp, scope)
	} else {
		return r.Throw(Error.Create(scope,
			AsString(exp.Type.OnString(r, scope, exp)),
		), scope)
	}
}

func (r *Runtime) EvalIndexing(node *ast.Indexing, scope *Scope) *Instance {
	target := r.Eval(node.Target, scope)

	args := make([]*Instance, len(node.Values)+1)
	args[0] = target
	for i, v := range node.Values {
		args[i+1] = r.Eval(v, scope)
	}

	return target.Type.OnGetItem(r, scope, args...)
}

func (r *Runtime) EvalWrapping(node *ast.Wrapping, scope *Scope) *Instance {
	exp := r.Eval(node.Expression, scope)

	if exp.Type == Maybe.Type {
		return exp
	}

	maybe := Maybe.Create()
	impl := maybe.Impl.(*MaybeDataImpl)
	err, ok := scope.GetInScope(RAISE_KEY)

	if ok {
		scope.Delete(RAISE_KEY)
		impl.Error = err.Value
	} else {
		impl.Value = exp
	}

	return maybe
}

func (r *Runtime) EvalUnwrap(node *ast.Unwrapping, scope *Scope) *Instance {
	target := r.Eval(node.Target, scope)
	return r.SolveMaybe(target, scope)
}

func (r *Runtime) SolveMaybe(target *Instance, scope *Scope) *Instance {
	if target.Type != Maybe.Type {
		return r.Throw(Error.Create(scope, "cannot unwrap non-maybe type"), scope)
	}

	maybe := target.Impl.(*MaybeDataImpl)

	if maybe.Error != nil {
		target.Type = maybe.Error.Type
		target.Impl = maybe.Error.Impl
		return maybe.Error
	} else {
		target.Type = maybe.Value.Type
		target.Impl = maybe.Value.Impl
		return Boolean.FALSE
	}
}

func (r *Runtime) EvalIf(node *ast.If, scope *Scope) *Instance {
	condition := r.Eval(node.Condition, scope)
	if condition == nil {
		return r.Throw(Error.Create(scope, "invalid condition"), scope)
	}

	if AsBool(condition) {
		return r.Eval(node.TrueBody, scope)

	} else if node.FalseBody != nil {
		return r.Eval(node.FalseBody, scope)
	}
	return Boolean.FALSE
}

func (r *Runtime) EvalSpreadOut(node *ast.SpreadOut, scope *Scope) *Instance {
	target := r.Eval(node.Target, scope)

	values := []*Instance{}
	var e *Instance
	r.ResolveIterator(target, scope, func(v *Instance, err *Instance) {
		if err != nil {
			e = err
		} else if v != nil {
			t := v.Impl.(*TupleDataImpl)
			values = append(values, t.Values...)
		}
	})

	if e != nil {
		return e
	}

	return Tuple.Create(values...)
}

func (r *Runtime) ResolveIterator(target *Instance, scope *Scope, up func(*Instance, *Instance)) {
	iter := target.Type.OnIter(r, scope, target)
	if iter.Type != Iterator.Type {
		up(nil, r.Throw(Error.Create(scope, "cannot iterate non-iterable type"), scope))
	}

	// impl := iter.Impl.(*IteratorDataImpl)
	fn := AsFunction(iter.Type.GetInstanceFn("next"))

	v := fn(r, scope, iter)
	for v != Iteration.DONE {
		ret := v.Impl.(*IterationDataImpl)
		up(ret.value(), nil)
		v = fn(r, scope, iter)
	}
	up(nil, nil)
}
