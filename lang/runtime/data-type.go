package runtime

import "sht/lang/ast"

type DataType interface {
	GetName() string
	Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance
	SetProperty(name string, node ast.Node)
	GetProperty(name string) ast.Node
	HasProperty(name string) bool
	SetStaticFn(name string, fn *Instance)
	GetStaticFn(name string) *Instance
	HasStaticFn(name string) bool
	SetInstanceFn(name string, fn *Instance)
	GetInstanceFn(name string) *Instance
	HasInstanceFn(name string) bool
	OnLen(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnSet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnGet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnSetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnGetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnNew(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnCall(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnBoolean(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnIn(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnIs(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnAdd(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnSub(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnMul(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnDiv(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnIntDiv(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnMod(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnPow(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnEq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnNeq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnGt(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnLt(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnGte(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnLte(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnPos(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnNeg(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
	OnNot(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance
}

type BaseDataType struct {
	Name        string
	Properties  map[string]ast.Node
	StaticFns   map[string]*Instance
	InstanceFns map[string]*Instance
}

func (d *BaseDataType) GetName() string {
	return d.Name
}

func (d *BaseDataType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	return r.Throw(Error.Create(s, "type '%s' does not allow instantiation", d.Name), s)
}

func (d *BaseDataType) SetProperty(name string, node ast.Node) {
	d.Properties[name] = node
}

func (d *BaseDataType) GetProperty(name string) ast.Node {
	return d.Properties[name]
}
func (d *BaseDataType) HasProperty(name string) bool {
	_, ok := d.Properties[name]
	return ok
}

func (d *BaseDataType) SetStaticFn(name string, fn *Instance) {
	d.StaticFns[name] = fn
}

func (d *BaseDataType) GetStaticFn(name string) *Instance {
	return d.StaticFns[name]
}
func (d *BaseDataType) HasStaticFn(name string) bool {
	_, ok := d.StaticFns[name]
	return ok
}

func (d *BaseDataType) SetInstanceFn(name string, fn *Instance) {
	d.InstanceFns[name] = fn
}

func (d *BaseDataType) GetInstanceFn(name string) *Instance {
	return d.InstanceFns[name]
}
func (d *BaseDataType) HasInstanceFn(name string) bool {
	_, ok := d.InstanceFns[name]
	return ok
}

func (d *BaseDataType) OnLen(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Len", self), s)
}
func (d *BaseDataType) OnSet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Set", self), s)
}
func (d *BaseDataType) OnGet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Get", self), s)
}
func (d *BaseDataType) OnSetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "SetItem", self), s)
}
func (d *BaseDataType) OnGetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "GetItem", self), s)
}
func (d *BaseDataType) OnNew(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "New", self), s)
}
func (d *BaseDataType) OnCall(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Call", self), s)
}
func (d *BaseDataType) OnBoolean(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return Boolean.TRUE
}
func (d *BaseDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self, args...)
}
func (d *BaseDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Repr", self), s)
}
func (d *BaseDataType) OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "To", self), s)
}
func (d *BaseDataType) OnIn(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "In", self), s)
}
func (d *BaseDataType) OnIs(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "In", self), s)
}
func (d *BaseDataType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Iter", self), s)
}
func (d *BaseDataType) OnBang(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidAction(s, "Bang", self), s)
}
func (d *BaseDataType) OnAdd(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Add", self), s)
}
func (d *BaseDataType) OnSub(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Sub", self), s)
}
func (d *BaseDataType) OnMul(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Mul", self), s)
}
func (d *BaseDataType) OnDiv(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Div", self), s)
}
func (d *BaseDataType) OnIntDiv(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "IntDiv", self), s)
}
func (d *BaseDataType) OnMod(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Mod", self), s)
}
func (d *BaseDataType) OnPow(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Pow", self), s)
}
func (d *BaseDataType) OnEq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self == args[0] {
		return Boolean.TRUE
	}

	return Boolean.FALSE
}
func (d *BaseDataType) OnNeq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self != args[0] {
		return Boolean.TRUE
	}

	return Boolean.FALSE
}
func (d *BaseDataType) OnGt(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Gt", self), s)
}
func (d *BaseDataType) OnLt(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Lt", self), s)
}
func (d *BaseDataType) OnGte(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Gte", self), s)
}
func (d *BaseDataType) OnLte(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Lte", self), s)
}
func (d *BaseDataType) OnPos(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Pos", self), s)
}
func (d *BaseDataType) OnNeg(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return r.Throw(Error.InvalidOperation(s, "Neg", self), s)
}
func (d *BaseDataType) OnNot(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	b := AsBool(d.OnBoolean(r, s, self, args...))
	if b {
		return Boolean.FALSE
	}

	return Boolean.TRUE
}
