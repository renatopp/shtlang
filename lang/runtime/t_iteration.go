package runtime

import (
	"sht/lang/ast"
)

var iterationDT = &IterationDataType{
	BaseDataType: BaseDataType{
		Name:        "Iteration",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var Iteration = &IterationInfo{
	Type: iterationDT,

	DONE: &Instance{
		Type: iterationDT,
		Impl: &IterationDataImpl{
			Properties: map[string]*Instance{
				"values": Tuple.Create(Boolean.FALSE),
				"done":   Boolean.TRUE,
			},
		},
	},
}

// ----------------------------------------------------------------------------
// ITERATION INFO
// ----------------------------------------------------------------------------
type IterationInfo struct {
	Type DataType

	DONE *Instance
}

func (t *IterationInfo) Create(values ...*Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &IterationDataImpl{
			Properties: map[string]*Instance{
				"value": Tuple.Create(values...),
				"done":  Boolean.FALSE,
			},
		},
	}
}

// ----------------------------------------------------------------------------
// ITERATION DATA TYPE
// ----------------------------------------------------------------------------
type IterationDataType struct {
	BaseDataType
}

func (d *IterationDataType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	return &Instance{
		Type: d,
		Impl: &IterationDataImpl{
			Properties: map[string]*Instance{
				"value": Boolean.FALSE,
				"done":  Boolean.TRUE,
			},
		},
	}
}

func (d *IterationDataType) OnNew(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*IterationDataImpl)

	this.Properties["done"] = Boolean.FALSE
	if len(args) == 1 {
		this.Properties["value"] = Tuple.Create(Boolean.FALSE)
	}

	if len(args) > 1 {
		this.Properties["value"] = Tuple.Create(args[1:]...)
	}

	return args[0]
}

func (d *IterationDataType) OnSet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*IterationDataImpl)
	name := AsString(args[1])

	_, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	this.Properties[name] = args[2]
	return args[2]
}

func (d *IterationDataType) OnGet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*IterationDataImpl)
	name := AsString(args[1])

	value, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	return value
}

func (d *IterationDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args[0])
}

func (d *IterationDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*IterationDataImpl)
	if AsBool(this.Properties["done"]) {
		return String.Create("Iteration<done>")
	}

	return String.Create("Iteration")
}

// ----------------------------------------------------------------------------
// ITERATION DATA IMPL
// ----------------------------------------------------------------------------
type IterationDataImpl struct {
	Properties map[string]*Instance
}

func (impl *IterationDataImpl) value() *Instance {
	return impl.Properties["value"]
}

func (impl *IterationDataImpl) done() *Instance {
	return impl.Properties["done"]
}
