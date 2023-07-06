package runtime

import (
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

func CreateCustomType(
	name string,
	properties map[string]ast.Node,
	staticFns map[string]*Instance,
	instanceFns map[string]*Instance,
	meta map[string]*Instance,
) *Instance {
	// println("CreateCustomType", name)
	// println("... properties")
	// for k, v := range properties {
	// 	println("......", k, v.String())
	// }
	// println("... staticFns")
	// for k, v := range staticFns {
	// 	println("......", k, v)
	// }
	// println("... instanceFns")
	// for k, v := range instanceFns {
	// 	println("......", k, v)
	// }
	// println("... meta")
	// for k, v := range meta {
	// 	println("......", k, v)
	// }

	return Type.Create(
		&CustomType{
			BaseDataType: BaseDataType{
				Name:        name,
				Properties:  properties,
				StaticFns:   staticFns,
				InstanceFns: instanceFns,
			},
			MetaFunctions: meta,
		},
	)
}

type CustomType struct {
	BaseDataType
	MetaFunctions map[string]*Instance
}

type CustomImpl struct {
	Properties map[string]*Instance
}

func (d *CustomType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	properties := map[string]*Instance{}

	switch init := init.(type) {
	case *ast.ListInitializer:
		return r.Throw(Error.Create(s, "Cannot instantiate custom type with list initializer"), s)

	case *ast.MapInitializer:
		for name, node := range init.Values {
			properties[name] = r.Eval(node, s)
		}
	}

	for name, node := range d.Properties {
		if _, ok := properties[name]; !ok {
			properties[name] = r.Eval(node, s)
		}
	}

	return &Instance{
		Type: d,
		Impl: &CustomImpl{
			Properties: properties,
		},
	}
}

func (d *CustomType) OnNew(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.New)]; fn != nil {
		fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return self
}
func (d *CustomType) OnLen(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Len)]; fn != nil {
		ret := fn.OnCall(r, s, append([]*Instance{self}, args...)...)

		if !ret.IsNumber() {
			return r.Throw(Error.Create(s, "Expected number on meta function '%s', got %s", string(meta.Len), ret.Type.GetName()), s)
		}

		return ret
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Len), self), s)
}
func (d *CustomType) OnSet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsCustom()
	name := AsString(args[0])
	old, hasOld := this.Properties[name]
	new := args[1]

	if fn := d.MetaFunctions[string(meta.SetProperty)]; fn != nil {
		old = Error.NoProperty(s, d.Name, name)
		new = fn.OnCall(r, s, self, args[0], old, new)

	} else if !hasOld {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)

	}

	this.Properties[name] = new
	return new
}
func (d *CustomType) OnGet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsCustom()
	name := AsString(args[0])

	value, has := this.Properties[name]
	if d.InstanceFns[name] != nil {
		value = d.InstanceFns[name]
		has = true
	}

	if fn := d.MetaFunctions[string(meta.GetProperty)]; fn != nil {
		value = Error.NoProperty(s, d.Name, name)
		ret := fn.OnCall(r, s, self, args[0], value)
		return ret

	} else if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	return value
}
func (d *CustomType) OnSetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.SetItem)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.SetItem), self), s)
}
func (d *CustomType) OnGetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.GetItem)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.GetItem), self), s)
}
func (d *CustomType) OnCall(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Call)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Call), self), s)
}
func (d *CustomType) OnBoolean(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Boolean)]; fn != nil {
		ret := fn.OnCall(r, s, append([]*Instance{self}, args...)...)
		return Boolean.Create(AsBool(ret))
	}
	return Boolean.TRUE
}
func (d *CustomType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.String)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return d.OnRepr(r, s, self, args...)
}
func (d *CustomType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Repr)]; fn != nil {
		ret := fn.OnCall(r, s, append([]*Instance{self}, args...)...)

		if !ret.IsString() {
			return ret.OnString(r, s, args...)
		}

		return ret
	}
	return String.Createf("<CustomType %s>", d.Name)
}
func (d *CustomType) OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.To)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.To), self), s)
}
func (d *CustomType) OnIn(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.In)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.In), self), s)
}
func (d *CustomType) OnIs(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Is)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Is), self), s)
}
func (d *CustomType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Iter)]; fn != nil {
		ret := fn.OnCall(r, s, append([]*Instance{self}, args...)...)

		if !ret.IsIterator() {
			return r.Throw(Error.Create(s, "Expected iterator on meta function '%s', got %s", string(meta.Iter), ret.Type.GetName()), s)
		}

		return ret
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Iter), self), s)
}
func (d *CustomType) OnAdd(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Add)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Add), self), s)
}
func (d *CustomType) OnSub(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Sub)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Sub), self), s)
}
func (d *CustomType) OnMul(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Mul)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Mul), self), s)
}
func (d *CustomType) OnDiv(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Div)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Div), self), s)
}
func (d *CustomType) OnIntDiv(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.IntDiv)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.IntDiv), self), s)
}
func (d *CustomType) OnMod(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Mod)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Mod), self), s)
}
func (d *CustomType) OnPow(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Pow)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidAction(s, string(meta.Pow), self), s)
}
func (d *CustomType) OnEq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Eq)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	if self == args[0] {
		return Boolean.TRUE
	}
	return Boolean.FALSE
}
func (d *CustomType) OnNeq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Neq)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	if self == args[0] {
		return Boolean.FALSE
	}
	return Boolean.TRUE
}
func (d *CustomType) OnGt(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Gt)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidOperation(s, string(meta.Gt), self), s)
}
func (d *CustomType) OnLt(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Lt)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidOperation(s, string(meta.Lt), self), s)
}
func (d *CustomType) OnGte(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Gte)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidOperation(s, string(meta.Gte), self), s)
}
func (d *CustomType) OnLte(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Lte)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidOperation(s, string(meta.Lte), self), s)
}
func (d *CustomType) OnPos(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Pos)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidOperation(s, string(meta.Pos), self), s)
}
func (d *CustomType) OnNeg(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Neg)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	return r.Throw(Error.InvalidOperation(s, string(meta.Neg), self), s)
}
func (d *CustomType) OnNot(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if fn := d.MetaFunctions[string(meta.Not)]; fn != nil {
		return fn.OnCall(r, s, append([]*Instance{self}, args...)...)
	}
	b := AsBool(d.OnBoolean(r, s, self, args...))
	if b {
		return Boolean.FALSE
	}

	return Boolean.TRUE
}
