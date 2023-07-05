package runtime

type Instance struct {
	Constant bool
	Type     DataType
	Impl     DataImpl
	MemberOf *Instance
}

func (i *Instance) Repr() string {
	return AsString(i.Type.OnRepr(nil, nil, i))
}

func (i *Instance) OnLen(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnLen(r, s, i, args...)
}
func (i *Instance) OnSet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnSet(r, s, i, args...)
}
func (i *Instance) OnGet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnGet(r, s, i, args...)
}
func (i *Instance) OnSetItem(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnSetItem(r, s, i, args...)
}
func (i *Instance) OnGetItem(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnGetItem(r, s, i, args...)
}
func (i *Instance) OnNew(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnNew(r, s, i, args...)
}
func (i *Instance) OnCall(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnCall(r, s, i, args...)
}
func (i *Instance) OnBoolean(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnBoolean(r, s, i, args...)
}
func (i *Instance) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnString(r, s, i, args...)
}
func (i *Instance) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnRepr(r, s, i, args...)
}
func (i *Instance) OnTo(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnTo(r, s, i, args...)
}
func (i *Instance) OnIn(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnIn(r, s, i, args...)
}
func (i *Instance) OnIs(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnIs(r, s, i, args...)
}
func (i *Instance) OnIter(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnIter(r, s, i, args...)
}
func (i *Instance) OnAdd(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnAdd(r, s, i, args...)
}
func (i *Instance) OnSub(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnSub(r, s, i, args...)
}
func (i *Instance) OnMul(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnMul(r, s, i, args...)
}
func (i *Instance) OnDiv(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnDiv(r, s, i, args...)
}
func (i *Instance) OnIntDiv(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnIntDiv(r, s, i, args...)
}
func (i *Instance) OnMod(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnMod(r, s, i, args...)
}
func (i *Instance) OnPow(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnPow(r, s, i, args...)
}
func (i *Instance) OnEq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnEq(r, s, i, args...)
}
func (i *Instance) OnNeq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnNeq(r, s, i, args...)
}
func (i *Instance) OnGt(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnGt(r, s, i, args...)
}
func (i *Instance) OnLt(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnLt(r, s, i, args...)
}
func (i *Instance) OnGte(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnGte(r, s, i, args...)
}
func (i *Instance) OnLte(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnLte(r, s, i, args...)
}
func (i *Instance) OnPos(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnPos(r, s, i, args...)
}
func (i *Instance) OnNeg(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnNeg(r, s, i, args...)
}
func (i *Instance) OnNot(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return i.Type.OnNot(r, s, i, args...)
}
