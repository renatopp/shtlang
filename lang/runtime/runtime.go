package runtime

import (
	"sht/lang/ast"
)

type Runtime struct {
	Global *Scope
}

func CreateRuntime() *Runtime {
	r := &Runtime{}

	Boolean.Setup()
	Error.Setup()
	Iteration.Setup()
	Iterator.Setup()
	Function.Setup()
	List.Setup()
	Maybe.Setup()
	Number.Setup()
	String.Setup()
	Tuple.Setup()
	Type.Setup()

	r.Global = CreateScope(nil, nil)
	r.Global.Set(SCOPE_NAME_KEY, Constant(String.Create("Global")))
	r.Global.Set(SCOPE_DEPTH_KEY, Constant(Number.ZERO))
	r.Global.Set(SCOPE_ID_KEY, Constant(String.Create(Id())))

	r.Global.Set(Boolean.Type.GetName(), Constant(Boolean.TypeInstance))
	r.Global.Set(Error.Type.GetName(), Constant(Error.TypeInstance))
	r.Global.Set(Iteration.Type.GetName(), Constant(Iteration.TypeInstance))
	r.Global.Set(Iterator.Type.GetName(), Constant(Iterator.TypeInstance))
	r.Global.Set(Function.Type.GetName(), Constant(Function.TypeInstance))
	r.Global.Set(List.Type.GetName(), Constant(List.TypeInstance))
	r.Global.Set(Maybe.Type.GetName(), Constant(Maybe.TypeInstance))
	r.Global.Set(Number.Type.GetName(), Constant(Number.TypeInstance))
	r.Global.Set(String.Type.GetName(), Constant(String.TypeInstance))
	r.Global.Set(Tuple.Type.GetName(), Constant(Tuple.TypeInstance))
	r.Global.Set(Type.Type.GetName(), Constant(Type.TypeInstance))

	r.Global.Set("Done", Constant(Iteration.DONE))
	r.Global.Set("map", Constant(b_map))
	r.Global.Set("each", Constant(b_each))
	r.Global.Set("filter", Constant(b_filter))
	r.Global.Set("reduce", Constant(b_reduce))
	r.Global.Set("sum", Constant(b_sum))

	r.Global.Set("range", Constant(b_range))

	r.Global.Set("print", Constant(b_print))

	return r
}

func (r *Runtime) Run(node ast.Node) string {
	instance := r.Eval(node, r.Global)

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

	case *ast.PostfixOperator:
		result = r.EvalPostfixOperator(n, scope)

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

	case *ast.Continue:
		result = r.EvalContinue(n, scope)

	case *ast.Break:
		result = r.EvalBreak(n, scope)

	case *ast.Return:
		result = r.EvalReturn(n, scope)

	case *ast.Raise:
		result = r.EvalRaise(n, scope)

	case *ast.Yield:
		result = r.EvalYield(n, scope)

	case *ast.Indexing:
		result = r.EvalIndexing(n, scope)

	case *ast.Wrapping:
		result = r.EvalWrapping(n, scope)

	case *ast.Unwrapping:
		result = r.EvalUnwrap(n, scope)

	case *ast.If:
		result = r.EvalIf(n, scope)

	case *ast.For:
		result = r.EvalFor(n, scope)

	case *ast.SpreadOut:
		result = r.EvalSpreadOut(n, scope)

	case *ast.Access:
		result = r.EvalAccess(n, scope)

	case *ast.Pipe:
		result = r.EvalPipe(n, scope)

	case *ast.PipeLoop:
		result = r.EvalPipeLoop(n, scope)

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

	scope.Set(RAISE_KEY, Constant(err))

	return err
}

func (r *Runtime) EvalBlock(node *ast.Block, scope *Scope) *Instance {
	var newScope *Scope
	var currentStatement int
	if scope.State != nil {
		state := scope.State.(*BlockState)
		currentStatement = state.Current
		newScope = state.Scope
	} else if node.Unscoped {
		newScope = scope
		currentStatement = 0
	} else {
		newScope = CreateScope(scope, scope.Caller)
		currentStatement = 0
	}

	var result *Instance
	// fmt.Println("BLOCK STATE", currentStatement)
	for i := currentStatement; i < len(node.Statements); i++ {
		stmt := node.Statements[i]
		result = r.Eval(stmt, newScope)
		if err, ok := newScope.Get(RAISE_KEY); ok {
			scope.Set(RAISE_KEY, err)
			result = err.Value
			break

		} else if v, ok := newScope.Get(CONTINUE_KEY); ok {
			scope.Set(CONTINUE_KEY, v)
			break

		} else if v, ok := newScope.Get(BREAK_KEY); ok {
			scope.Set(BREAK_KEY, v)
			break

		} else if v, ok := newScope.Get(RETURN_KEY); ok {
			scope.Set(RETURN_KEY, v)
			break

		} else if v, ok := newScope.Get(YIELD_KEY); ok {

			cur := i
			if newScope.HasInScope(JUST_YIELDED_KEY) {
				cur += 1
			}
			scope.State = &BlockState{
				Scope:   newScope,
				Current: cur,
			}
			newScope.Delete(YIELD_KEY)
			newScope.Delete(JUST_YIELDED_KEY)
			scope.Set(YIELD_KEY, v)
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
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnAdd(r, scope, left, right)

	case "-":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnSub(r, scope, left, right)

	case "*":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnMul(r, scope, left, right)

	case "/":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnDiv(r, scope, left, right)

	case "//":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnIntDiv(r, scope, left, right)

	case "%":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnMod(r, scope, left, right)

	case "**":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnPow(r, scope, left, right)

	case "==":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnEq(r, scope, left, right)

	case "!=":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnNeq(r, scope, left, right)

	case ">":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnGt(r, scope, left, right)

	case "<":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnLt(r, scope, left, right)

	case ">=":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnGte(r, scope, left, right)

	case "<=":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return left.Type.OnLte(r, scope, left, right)

	case "and", "or", "nand", "nor", "xor", "nxor":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

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
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		lt := AsString(left)
		rt := AsString(right)
		return String.Create(lt + rt)

	case "??":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

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

	case "as":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		id, ok := node.Right.(*ast.Identifier)
		if !ok {
			return r.Throw(Error.Create(scope, "'as' expression requires an identifier on the right side"), scope)
		}

		scope.Set(id.Value, Variable(left))
		return left

	case "is":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return right.Type.OnIs(r, scope, right, left)

	case "in":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return right.Type.OnIn(r, scope, right, left)

	case "to":
		left := r.Eval(node.Left, scope)
		if scope.HasInScope(RAISE_KEY) {
			return left
		}
		iter := left.Type.OnIter(r, scope, left)

		right := r.Eval(node.Right, scope)
		if scope.HasInScope(RAISE_KEY) {
			return right
		}

		return right.Type.OnTo(r, scope, right, iter)
	}

	return nil
}

func (r *Runtime) EvalPostfixOperator(node *ast.PostfixOperator, scope *Scope) *Instance {
	left := r.Eval(node.Left, scope)

	switch node.Operator {
	case "++":
		return left.Type.OnPostInc(r, scope, left)
	case "--":
		return left.Type.OnPostDec(r, scope, left)
	}

	return nil
}

func (r *Runtime) ResolveAssignment(left ast.Node, right *Instance, assignment *ast.Assignment, scope *Scope) *Instance {
	switch id := left.(type) {
	case *ast.Tuple:
		if len(id.Values) == 1 {
			return r.ResolveAssignment(id.Values[0], right, assignment, scope)
		}

		leftLength := len(id.Values)
		rightLength := AsNumber(right.Type.OnLen(r, scope, right))

		j := 0
		for i, lv := range id.Values {
			lvspread, isSpread := lv.(*ast.SpreadIn)
			if isSpread {
				//              remainer of right side - remainer of left side
				spreadAmount := (int(rightLength) - i) - (int(leftLength) - i - 1)
				spreadItems := []*Instance{}

				for k := 0; k < spreadAmount; k++ {
					rv := right.Type.OnGetItem(r, scope, right, Number.Create(float64(j)))
					spreadItems = append(spreadItems, rv)
					j++
				}

				rv := List.Create(spreadItems...)
				r.ResolveAssignment(lvspread.Target, rv, assignment, scope)

			} else {
				if j >= int(rightLength) {
					return r.Throw(Error.Create(scope, "assignment right side has less elements than left side"), scope)
				}

				rv := right.Type.OnGetItem(r, scope, right, Number.Create(float64(j)))
				r.ResolveAssignment(lv, rv, assignment, scope)
				j++
			}
		}

		if j < int(rightLength) {
			return r.Throw(Error.Create(scope, "assignment right side has more elements than left side"), scope)
		}

		return right

	case *ast.Identifier:
		return r.Assign(id.Value, right, assignment.Definition, assignment.Constant, scope)

	case *ast.Indexing:
		target := r.Eval(id.Target, scope)
		if target == nil {
			return r.Throw(Error.Create(scope, "invalid assignment target"), scope)
		}

		if id.Values == nil || len(id.Values) == 0 {
			return r.Throw(Error.Create(scope, "invalid assignment target"), scope)
		}

		idx := r.Eval(id.Values[0], scope)
		return target.Type.OnSetItem(r, scope, target, idx, right)

	case *ast.Access:
		target := r.Eval(id.Left, scope)
		if target == nil {
			return r.Throw(Error.Create(scope, "invalid assignment target"), scope)
		}

		name := id.Right.(*ast.Identifier).Value
		return target.Type.OnSet(r, scope, target, String.Create(name), right)

	default:
		return r.Throw(Error.Create(scope, "cannot assign to non-identifier"), scope)
	}
}

func (r *Runtime) EvalAssignment(node *ast.Assignment, scope *Scope) *Instance {
	right := r.Eval(node.Expression, scope)
	return r.ResolveAssignment(node.Identifier, right, node, scope)
}

func (r *Runtime) Assign(name string, exp *Instance, def, cnst bool, scope *Scope) *Instance {
	if name == "_" {
		return exp
	}

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
			Default: nil,
		}

		if param.Default != nil {
			p.Default = r.Eval(param.Default, nil)
		}

		if p.Spread {
			if p.Default != nil {
				return r.Throw(
					Error.Create(scope, "spread arguments cannot have default values: '%s'", p.Name),
					scope,
				)
			}

			if hasSpread {
				return r.Throw(
					Error.Create(scope, "arguments can only have one spread argument: '%s'", p.Name),
					scope,
				)
			}

			hasSpread = true
		}

		if p.Default != nil {
			hasDefault = true

		} else if hasDefault && !p.Spread {
			return r.Throw(
				Error.Create(scope, "default arguments must be at the end: '%s'", p.Name),
				scope,
			)
		}

		params[i] = p
	}

	fn := Function.Create(name, params, node.Body, scope)
	impl := fn.Impl.(*FunctionDataImpl)
	impl.Generator = node.Generator

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
	if target.MemberOf != nil {
		args = append(args, target.MemberOf)
	}
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

func (r *Runtime) EvalContinue(node *ast.Continue, scope *Scope) *Instance {
	scope.Set(CONTINUE_KEY, Constant(Boolean.TRUE))
	return Boolean.TRUE
}

func (r *Runtime) EvalBreak(node *ast.Break, scope *Scope) *Instance {
	scope.Set(BREAK_KEY, Constant(Boolean.TRUE))
	return Boolean.TRUE
}

func (r *Runtime) EvalReturn(node *ast.Return, scope *Scope) *Instance {
	exp := r.Eval(node.Expression, scope)
	if exp == nil {
		exp = Boolean.FALSE
	}

	scope.Set(RETURN_KEY, Constant(exp))
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

func (r *Runtime) EvalYield(node *ast.Yield, scope *Scope) *Instance {
	exp := r.Eval(node.Expression, scope)
	if exp == nil {
		exp = Boolean.FALSE
	}

	scope.Set(YIELD_KEY, Constant(exp))
	scope.Set(JUST_YIELDED_KEY, Constant(Boolean.TRUE))
	return exp
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

	err, ok := scope.GetInScope(RAISE_KEY)

	if ok {
		scope.Delete(RAISE_KEY)
		return Maybe.CreateError(err.Value)
	} else {
		return Maybe.Create(exp)
	}
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
	var newScope *Scope
	var condition *bool
	if scope.State != nil {
		// fmt.Println("IF STATE")
		state := scope.State.(*IfState)
		condition = &state.Condition
		newScope = state.Scope
	} else {
		// fmt.Println("IF NO STATE")
		newScope = CreateScope(scope, scope.Caller)
		condition = nil
	}

	if condition == nil {
		c := r.Eval(node.Condition, newScope)
		if c == nil {
			return r.Throw(Error.Create(scope, "invalid condition"), scope)
		}

		t := true
		f := false
		if AsBool(c) {
			condition = &t
		} else {
			condition = &f
		}
	}

	var ret *Instance
	if *condition {
		ret = r.Eval(node.TrueBody, newScope)

	} else if node.FalseBody != nil {
		ret = r.Eval(node.FalseBody, newScope)
	}

	if newScope.HasInScope(CONTINUE_KEY) {
		err, _ := newScope.GetInScope(CONTINUE_KEY)
		scope.Set(CONTINUE_KEY, err)

	} else if newScope.HasInScope(BREAK_KEY) {
		err, _ := newScope.GetInScope(BREAK_KEY)
		scope.Set(BREAK_KEY, err)

	} else if newScope.HasInScope(RAISE_KEY) {
		err, _ := newScope.GetInScope(RAISE_KEY)
		scope.Set(RAISE_KEY, err)

	} else if newScope.HasInScope(RETURN_KEY) {
		err, _ := newScope.GetInScope(RETURN_KEY)
		scope.Set(RETURN_KEY, err)

	} else if newScope.HasInScope(YIELD_KEY) {
		v, _ := newScope.GetInScope(YIELD_KEY)

		scope.State = &IfState{
			Scope:     newScope,
			Condition: *condition,
		}
		newScope.Delete(YIELD_KEY)
		newScope.Delete(JUST_YIELDED_KEY)
		scope.Set(YIELD_KEY, v)
	}

	return ret
}

func (r *Runtime) EvalFor(node *ast.For, scope *Scope) *Instance {
	var newScope *Scope
	var evalCondition bool
	if scope.State != nil {
		state := scope.State.(*ForState)
		newScope = state.Scope
		evalCondition = false
	} else {
		newScope = CreateScope(scope, scope.Caller)
		evalCondition = true
	}

	for {
		if evalCondition {
			newScope.Clear()

			c := r.Eval(node.Condition, newScope)
			if c == nil {
				return r.Throw(Error.Create(scope, "invalid condition"), scope)
			}

			if !AsBool(c) {
				break
			}
		}

		evalCondition = true
		r.Eval(node.Body, newScope)

		// execute block, if return evalCondition = true
		if newScope.HasInScope(BREAK_KEY) {
			newScope.Delete(BREAK_KEY)
			break

		} else if newScope.HasInScope(CONTINUE_KEY) {
			newScope.Delete(CONTINUE_KEY)
			continue

		} else if newScope.HasInScope(RAISE_KEY) {
			err, _ := newScope.GetInScope(RAISE_KEY)
			scope.Set(RAISE_KEY, err)
			break

		} else if newScope.HasInScope(RETURN_KEY) {
			err, _ := newScope.GetInScope(RETURN_KEY)
			scope.Set(RETURN_KEY, err)
			break

		} else if newScope.HasInScope(YIELD_KEY) {
			v, _ := newScope.GetInScope(YIELD_KEY)

			scope.State = &ForState{
				Scope: newScope,
			}
			newScope.Delete(YIELD_KEY)
			newScope.Delete(JUST_YIELDED_KEY)
			scope.Set(YIELD_KEY, v)
			break
		}
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

	if scope.HasInScope(RAISE_KEY) {
		up(nil, Boolean.TRUE)
		return
	}

	impl := iter.Type.GetInstanceFn("next")
	fn := AsFunction(impl)

	v := fn(r, scope, impl, iter)
	if scope.HasInScope(RAISE_KEY) {
		up(nil, Boolean.TRUE)
		return
	}

	it := v.Impl.(*IterationDataImpl)
	for v != Iteration.DONE {
		if scope.HasInScope(RAISE_KEY) {
			up(nil, it.error())
			return
		}

		up(it.value(), nil)
		v = fn(r, scope, impl, iter)
		it = v.Impl.(*IterationDataImpl)

	}

	up(nil, nil)
}

func (r *Runtime) EvalAccess(node *ast.Access, scope *Scope) *Instance {
	left := r.Eval(node.Left, scope)
	right := node.Right.(*ast.Identifier).Value

	res := left.Type.OnGet(r, scope, left, String.Create(right))
	res.MemberOf = left
	return res
}

func (r *Runtime) EvalPipe(node *ast.Pipe, scope *Scope) *Instance {
	if scope.PipeCounter != 0 && node.To != nil {
		return r.Throw(Error.Create(scope, "'to' expression can only be used at the end of a pipe"), scope)
	}

	scope.PipeCounter += 1

	left := r.Eval(node.Left, scope)
	if left.Type != Iterator.Type {
		left = left.Type.OnIter(r, scope, left)
		if left.Type != Iterator.Type {
			return r.Throw(Error.Create(scope, "cannot iterate non-iterable type"), scope)
		}
	}

	if node.To != nil {
		to := r.Eval(node.To, scope)
		scope.PipeCounter -= 1
		return to.Type.OnTo(r, scope, to, left)
	}

	var pipeFn *Instance
	pipeArgs := []*Instance{}
	addArgs := []*Instance{}
	switch t := node.PipeFn.(type) {
	case *ast.Identifier:
		pipeFn = r.Eval(t, scope)
		pipeArgs = append(pipeArgs, pipeFn)
		pipeArgs = append(pipeArgs, left)

	case *ast.Call:
		pipeFn = r.Eval(t.Target, scope)
		pipeArgs = append(pipeArgs, pipeFn)
		pipeArgs = append(pipeArgs, left)
		for _, v := range t.Arguments {
			addArgs = append(addArgs, r.Eval(v, scope))
		}
	}

	if pipeFn.Type == Error.Type {
		return r.Throw(Error.Create(scope, "invalid pipe function"), scope)
	}

	var argFn *Instance
	if node.ArgFn != nil {
		argFn = r.Eval(node.ArgFn, scope)
		if argFn.Type != Function.Type {
			return r.Throw(Error.Create(scope, "cannot use non-function as argument function"), scope)
		}
		argFn.Impl.(*FunctionDataImpl).Piped = true
	} else {
		argFn = Boolean.FALSE
	}
	pipeArgs = append(pipeArgs, argFn)

	for _, v := range addArgs {
		pipeArgs = append(pipeArgs, v)
	}

	pipe := pipeFn.Type.OnPipe(r, scope, pipeArgs...)

	scope.PipeCounter -= 1
	if scope.PipeCounter == 0 {
		values := []*Instance{}
		r.ResolveIterator(pipe, scope, func(v *Instance, err *Instance) {
			if err != nil {
				values = append(values, err)
			} else if v != nil {
				t := v.Impl.(*TupleDataImpl)
				values = append(values, t.Values...)
			}
		})

		return List.Create(values...)
	}

	return pipe
}

func (r *Runtime) EvalPipeLoop(node *ast.PipeLoop, scope *Scope) *Instance {
	var newScope *Scope
	var evalCondition bool
	var i_iterator *Instance
	if scope.State != nil {
		state := scope.State.(*PipeLoopState)
		newScope = state.Scope
		i_iterator = state.Iterator
		evalCondition = false
	} else {
		newScope = CreateScope(scope, scope.Caller)
		i_eval := r.Eval(node.Iterator, newScope)
		if i_eval == nil {
			return r.Throw(Error.Create(scope, "invalid iterator"), scope)
		}

		i_iterator = i_eval.Type.OnIter(r, scope, i_eval)
		evalCondition = true
	}

	iterator := i_iterator.Impl.(*IteratorDataImpl)
	for {
		if evalCondition {
			newScope.Clear()

			i_next := iterator.next()
			i_iteration := i_next.Type.OnCall(r, scope, i_next, i_iterator)
			iteration := i_iteration.Impl.(*IterationDataImpl)

			if iteration.done() == Boolean.TRUE || iteration.error() == Boolean.TRUE {
				break
			}

			i_value := iteration.value()
			value := i_value.Impl.(*TupleDataImpl)

			asTuple, isTuple := node.Assignment.(*ast.Tuple)
			if !isTuple {
				asTuple = &ast.Tuple{
					Values: []ast.Node{node.Assignment},
				}
			} else {
				r.Throw(Error.Create(scope, "tuple assignment not implement yet"), scope)
			}

			r.ResolveAssignment(asTuple, value.Values[0], &ast.Assignment{
				Definition: true,
				Constant:   false,
			}, newScope)
		}

		r.Eval(node.Body, newScope)

		// execute block, if return evalCondition = true
		if newScope.HasInScope(BREAK_KEY) {
			newScope.Delete(BREAK_KEY)
			break

		} else if newScope.HasInScope(CONTINUE_KEY) {
			newScope.Delete(CONTINUE_KEY)
			continue

		} else if newScope.HasInScope(RAISE_KEY) {
			err, _ := newScope.GetInScope(RAISE_KEY)
			scope.Set(RAISE_KEY, err)
			break

		} else if newScope.HasInScope(RETURN_KEY) {
			err, _ := newScope.GetInScope(RETURN_KEY)
			scope.Set(RETURN_KEY, err)
			break

		} else if newScope.HasInScope(YIELD_KEY) {
			v, _ := newScope.GetInScope(YIELD_KEY)

			scope.State = &ForState{
				Scope: newScope,
			}
			newScope.Delete(YIELD_KEY)
			newScope.Delete(JUST_YIELDED_KEY)
			scope.Set(YIELD_KEY, v)
			break
		}
	}

	return Boolean.FALSE
}
