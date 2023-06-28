package runtime

import (
	"sht/lang/ast"
)

var iteratorDT = &IteratorDataType{
	BaseDataType: BaseDataType{
		Name:        "Iterator",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
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

func (t *IteratorInfo) Create(nextFn *Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &IteratorDataImpl{
			Properties: map[string]*Instance{
				"next": nextFn,
			},
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

func (d *IteratorDataType) OnSet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*IteratorDataImpl)
	name := AsString(args[1])

	_, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	this.Properties[name] = args[2]
	return args[2]
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

	_, ok := args[1].Impl.(Function)
	if !ok {
		return r.Throw(Error.Create(s, "Expected function, %s given", args[0].Type.GetName()), s)
	}

	this := args[0].Impl.(*IteratorDataImpl)
	this.Properties["next"] = args[1]

	return args[0]
}

func (d *IteratorDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args[0])
}

func (d *IteratorDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return String.Create("Iterator")
}

// ----------------------------------------------------------------------------
// ITERATOR DATA IMPL
// ----------------------------------------------------------------------------
type IteratorDataImpl struct {
	Properties map[string]*Instance
}
