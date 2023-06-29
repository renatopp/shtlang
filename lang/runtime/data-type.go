package runtime

import "sht/lang/ast"

type DataType interface {
	GetName() string
	Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance
	GetProperty(name string) ast.Node
	HasProperty(name string) bool
	GetStaticFn(name string) Callable
	HasStaticFn(name string) bool
	GetInstanceFn(name string) Callable
	HasInstanceFn(name string) bool
	OnLen(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnSet(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnGet(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnSetItem(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnGetItem(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnNew(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnCall(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnBoolean(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnString(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnTo(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnIter(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnAdd(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnSub(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnMul(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnDiv(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnIntDiv(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnMod(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnPow(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnEq(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnNeq(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnGt(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnLt(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnGte(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnLte(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnPos(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnNeg(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnNot(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnPostInc(r *Runtime, s *Scope, args ...*Instance) *Instance
	OnPostDec(r *Runtime, s *Scope, args ...*Instance) *Instance
}

type BaseDataType struct {
	Name        string
	Properties  map[string]ast.Node
	StaticFns   map[string]Callable
	InstanceFns map[string]Callable
}

func (d *BaseDataType) GetName() string {
	return d.Name
}

func (d *BaseDataType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	return r.Throw(Error.Create(s, "type '%s' does not allow instantiation", d.Name), s)
}

func (d *BaseDataType) GetProperty(name string) ast.Node {
	return d.Properties[name]
}
func (d *BaseDataType) HasProperty(name string) bool {
	_, ok := d.Properties[name]
	return ok
}

func (d *BaseDataType) GetStaticFn(name string) Callable {
	return d.StaticFns[name]
}
func (d *BaseDataType) HasStaticFn(name string) bool {
	_, ok := d.StaticFns[name]
	return ok
}

func (d *BaseDataType) GetInstanceFn(name string) Callable {
	return d.InstanceFns[name]
}
func (d *BaseDataType) HasInstanceFn(name string) bool {
	_, ok := d.InstanceFns[name]
	return ok
}

func (d *BaseDataType) OnLen(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Len", args[0]), s)
}
func (d *BaseDataType) OnSet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Set", args[0]), s)
}
func (d *BaseDataType) OnGet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Get", args[0]), s)
}
func (d *BaseDataType) OnSetItem(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "SetItem", args[0]), s)
}
func (d *BaseDataType) OnGetItem(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "GetItem", args[0]), s)
}
func (d *BaseDataType) OnNew(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "New", args[0]), s)
}
func (d *BaseDataType) OnCall(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Call", args[0]), s)
}
func (d *BaseDataType) OnBoolean(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return Boolean.TRUE
}
func (d *BaseDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args...)
}
func (d *BaseDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Repr", args[0]), s)
}
func (d *BaseDataType) OnTo(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "To", args[0]), s)
}
func (d *BaseDataType) OnIter(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Iter", args[0]), s)
}
func (d *BaseDataType) OnBang(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Bang", args[0]), s)
}
func (d *BaseDataType) OnAdd(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Add", args[0]), s)
}
func (d *BaseDataType) OnSub(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Sub", args[0]), s)
}
func (d *BaseDataType) OnMul(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Mul", args[0]), s)
}
func (d *BaseDataType) OnDiv(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Div", args[0]), s)
}
func (d *BaseDataType) OnIntDiv(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "IntDiv", args[0]), s)
}
func (d *BaseDataType) OnMod(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Mod", args[0]), s)
}
func (d *BaseDataType) OnPow(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Pow", args[0]), s)
}
func (d *BaseDataType) OnEq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0] == args[1] {
		return Boolean.TRUE
	}

	return Boolean.FALSE
}
func (d *BaseDataType) OnNeq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0] != args[1] {
		return Boolean.TRUE
	}

	return Boolean.FALSE
}
func (d *BaseDataType) OnGt(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Gt", args[0]), s)
}
func (d *BaseDataType) OnLt(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Lt", args[0]), s)
}
func (d *BaseDataType) OnGte(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Gte", args[0]), s)
}
func (d *BaseDataType) OnLte(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Lte", args[0]), s)
}
func (d *BaseDataType) OnPos(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Pos", args[0]), s)
}
func (d *BaseDataType) OnNeg(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Neg", args[0]), s)
}
func (d *BaseDataType) OnNot(r *Runtime, s *Scope, args ...*Instance) *Instance {
	b := AsBool(d.OnBoolean(r, s, args...))
	if b {
		return Boolean.FALSE
	}

	return Boolean.TRUE
}
func (d *BaseDataType) OnPostInc(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "PostInc", args[0]), s)
}
func (d *BaseDataType) OnPostDec(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "PostDec", args[0]), s)
}
