package runtime

import (
	"sht/lang/ast"
)

var iteratorDT = &IteratorDataType{
	BaseDataType: BaseDataType{
		Name:        "Iterator",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var Iterator = &IteratorInfo{
	Type: iteratorDT,
}

// ----------------------------------------------------------------------------
// ITERATOR INFO
// ----------------------------------------------------------------------------
type IteratorInfo struct {
	Type         DataType
	TypeInstance *Instance
}

func (t *IteratorInfo) Setup() {
	t.TypeInstance = Type.Create(Iterator.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
	t.Type.SetInstanceFn("next", Iterator_Next)
}

func (t *IteratorInfo) Create(nextFn *Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &IteratorDataImpl{
			Properties: map[string]*Instance{
				"done": Boolean.FALSE,
			},
			Next: nextFn,
		},
	}
}

// ----------------------------------------------------------------------------
// ITERATOR DATA TYPE
// ----------------------------------------------------------------------------
type IteratorDataType struct {
	BaseDataType
}

func (d *IteratorDataType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	return Iterator.Create(DoneFn)
}

func (d *IteratorDataType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return self
}

func (d *IteratorDataType) OnGet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*IteratorDataImpl)
	name := AsString(args[0])

	value, has := d.InstanceFns[name]
	if has {
		return value
	}

	value, has = this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	return value
}

func (d *IteratorDataType) OnNew(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if len(args) == 0 {
		return self
	}

	_, ok := args[0].Impl.(*FunctionDataImpl)
	if !ok {
		return r.Throw(Error.Create(s, "Expected function, %s given", self.Type.GetName()), s)
	}

	this := self.Impl.(*IteratorDataImpl)
	this.Next = args[0]

	return self
}

func (d *IteratorDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self)
}

func (d *IteratorDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*IteratorDataImpl)
	if AsBool(this.done()) {
		return String.Create("<Iterator:done>")
	} else {
		return String.Create("<Iterator>")
	}
}

// ----------------------------------------------------------------------------
// ITERATOR DATA IMPL
// ----------------------------------------------------------------------------
type IteratorDataImpl struct {
	Properties map[string]*Instance
	Next       *Instance
}

func (impl *IteratorDataImpl) done() *Instance {
	return impl.Properties["done"]
}

func (impl *IteratorDataImpl) next() *Instance {
	return Iterator.Type.GetInstanceFn("next")
}

var Iterator_Next = Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	// self => function
	// args[0] => this (the iterator object)
	this := args[0].AsIterator()
	if this.Properties["done"] == Boolean.TRUE {
		return Iteration.DONE
	}

	ret := this.Next.OnCall(r, s, args[0])

	if ret.Type != Iteration.Type {
		this.Properties["done"] = Boolean.TRUE
		return Iteration.Error(Error.Create(s, "Expected iteration, %s given", ret.Type.GetName()))
	}

	if s.IsInterruptedAs(FlowRaise) {
		this.Properties["done"] = Boolean.TRUE
		return Iteration.Error(s.Interruption.Value)
	}

	if AsBool(ret.AsIteration().error()) {
		this.Properties["done"] = Boolean.TRUE
	}

	if ret == Iteration.DONE {
		this.Properties["done"] = Boolean.TRUE
	}

	return ret
})
