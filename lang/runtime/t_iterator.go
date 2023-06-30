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
	Type DataType
}

func (t *IteratorInfo) Setup() {
	t.Type.SetInstanceFn("next", Iterator_Next)
}

func (t *IteratorInfo) Create(nextFn *Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &IteratorDataImpl{
			Properties: map[string]*Instance{
				"finished": Boolean.FALSE,
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

func (d *IteratorDataType) OnIter(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return args[0]
}

func (d *IteratorDataType) OnGet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*IteratorDataImpl)
	name := AsString(args[1])

	value, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	return value
}

func (d *IteratorDataType) OnNew(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if len(args) == 1 {
		return args[0]
	}

	_, ok := args[1].Impl.(*FunctionDataImpl)
	if !ok {
		return r.Throw(Error.Create(s, "Expected function, %s given", args[0].Type.GetName()), s)
	}

	this := args[0].Impl.(*IteratorDataImpl)
	this.Next = args[1]

	return args[0]
}

func (d *IteratorDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args[0])
}

func (d *IteratorDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*IteratorDataImpl)
	if AsBool(this.finished()) {
		return String.Create("<Iterator:finished>")
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

func (impl *IteratorDataImpl) finished() *Instance {
	return impl.Properties["finished"]
}

var Iterator_Next = Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*IteratorDataImpl)
	nextFn := this.Next.Impl.(*FunctionDataImpl)
	ret := nextFn.Call(r, s)

	if ret.Type != Iteration.Type {
		return r.Throw(Error.Create(s, "Expected Iteration, %s given", ret.Type.GetName()), s)
	}

	if ret == Iteration.DONE {
		this.Properties["finished"] = Boolean.TRUE
	}

	return ret
})
