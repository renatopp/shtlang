package runtime

import "sht/lang/ast"

type DataType interface {
	GetName() string
	GetProperty(name string) ast.Node
	HasProperty(name string) bool
	GetStaticFn(name string) Function
	HasStaticFn(name string) bool
	GetInstanceFn(name string) Function
	HasInstanceFn(name string) bool
	OnSet(r *Runtime, args ...*Instance) *Instance
	OnGet(r *Runtime, args ...*Instance) *Instance
	OnSetItem(r *Runtime, args ...*Instance) *Instance
	OnGetItem(r *Runtime, args ...*Instance) *Instance
	OnNew(r *Runtime, args ...*Instance) *Instance
	OnCall(r *Runtime, args ...*Instance) *Instance
	OnBoolean(r *Runtime, args ...*Instance) *Instance
	OnString(r *Runtime, args ...*Instance) *Instance
	OnRepr(r *Runtime, args ...*Instance) *Instance
	OnTo(r *Runtime, args ...*Instance) *Instance
	OnIter(r *Runtime, args ...*Instance) *Instance
	OnBang(r *Runtime, args ...*Instance) *Instance
	OnAdd(r *Runtime, args ...*Instance) *Instance
	OnSub(r *Runtime, args ...*Instance) *Instance
	OnMul(r *Runtime, args ...*Instance) *Instance
	OnDiv(r *Runtime, args ...*Instance) *Instance
	OnIntDiv(r *Runtime, args ...*Instance) *Instance
	OnMod(r *Runtime, args ...*Instance) *Instance
	OnPow(r *Runtime, args ...*Instance) *Instance
	OnEq(r *Runtime, args ...*Instance) *Instance
	OnNeq(r *Runtime, args ...*Instance) *Instance
	OnGt(r *Runtime, args ...*Instance) *Instance
	OnLt(r *Runtime, args ...*Instance) *Instance
	OnGte(r *Runtime, args ...*Instance) *Instance
	OnLte(r *Runtime, args ...*Instance) *Instance
	OnPos(r *Runtime, args ...*Instance) *Instance
	OnNeg(r *Runtime, args ...*Instance) *Instance
	OnNot(r *Runtime, args ...*Instance) *Instance
	OnPostInc(r *Runtime, args ...*Instance) *Instance
	OnPostDec(r *Runtime, args ...*Instance) *Instance
}

type BaseDataType struct {
	Name        string
	Properties  map[string]ast.Node
	StaticFns   map[string]Function
	InstanceFns map[string]Function
}

func (d *BaseDataType) GetName() string {
	return d.Name
}

func (d *BaseDataType) GetProperty(name string) ast.Node {
	return d.Properties[name]
}
func (d *BaseDataType) HasProperty(name string) bool {
	_, ok := d.Properties[name]
	return ok
}

func (d *BaseDataType) GetStaticFn(name string) Function {
	return d.StaticFns[name]
}
func (d *BaseDataType) HasStaticFn(name string) bool {
	_, ok := d.StaticFns[name]
	return ok
}

func (d *BaseDataType) GetInstanceFn(name string) Function {
	return d.InstanceFns[name]
}
func (d *BaseDataType) HasInstanceFn(name string) bool {
	_, ok := d.InstanceFns[name]
	return ok
}

func (d *BaseDataType) OnSet(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("Set", args[0])
}
func (d *BaseDataType) OnGet(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("Get", args[0])
}
func (d *BaseDataType) OnSetItem(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("SetItem", args[0])
}
func (d *BaseDataType) OnGetItem(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("GetItem", args[0])
}
func (d *BaseDataType) OnNew(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("New", args[0])
}
func (d *BaseDataType) OnCall(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("Call", args[0])
}
func (d *BaseDataType) OnBoolean(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("Boolean", args[0])
}
func (d *BaseDataType) OnString(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("String", args[0])
}
func (d *BaseDataType) OnRepr(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("Repr", args[0])
}
func (d *BaseDataType) OnTo(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("To", args[0])
}
func (d *BaseDataType) OnIter(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("Iter", args[0])
}
func (d *BaseDataType) OnBang(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidAction("Bang", args[0])
}
func (d *BaseDataType) OnAdd(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Add", args[0])
}
func (d *BaseDataType) OnSub(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Sub", args[0])
}
func (d *BaseDataType) OnMul(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Mul", args[0])
}
func (d *BaseDataType) OnDiv(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Div", args[0])
}
func (d *BaseDataType) OnIntDiv(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("IntDiv", args[0])
}
func (d *BaseDataType) OnMod(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Mod", args[0])
}
func (d *BaseDataType) OnPow(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Pow", args[0])
}
func (d *BaseDataType) OnEq(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Eq", args[0])
}
func (d *BaseDataType) OnNeq(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Neq", args[0])
}
func (d *BaseDataType) OnGt(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Gt", args[0])
}
func (d *BaseDataType) OnLt(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Lt", args[0])
}
func (d *BaseDataType) OnGte(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Gte", args[0])
}
func (d *BaseDataType) OnLte(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Lte", args[0])
}
func (d *BaseDataType) OnPos(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Pos", args[0])
}
func (d *BaseDataType) OnNeg(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Neg", args[0])
}
func (d *BaseDataType) OnNot(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("Not", args[0])
}
func (d *BaseDataType) OnPostInc(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("PostInc", args[0])
}
func (d *BaseDataType) OnPostDec(r *Runtime, args ...*Instance) *Instance {
	return Error.InvalidOperation("PostDec", args[0])
}
