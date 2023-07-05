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
				"value": Tuple.Create(Boolean.FALSE),
				"done":  Boolean.TRUE,
				"error": Boolean.FALSE,
			},
		},
	},
}

// ----------------------------------------------------------------------------
// ITERATION INFO
// ----------------------------------------------------------------------------
type IterationInfo struct {
	Type         DataType
	TypeInstance *Instance

	DONE *Instance
}

func (t *IterationInfo) Create(values ...*Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &IterationDataImpl{
			Properties: map[string]*Instance{
				"value": Tuple.Create(values...),
				"done":  Boolean.FALSE,
				"error": Boolean.FALSE,
			},
		},
	}
}

func (t *IterationInfo) CreateAsTuple(tuple *Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &IterationDataImpl{
			Properties: map[string]*Instance{
				"value": tuple,
				"done":  Boolean.FALSE,
				"error": Boolean.FALSE,
			},
		},
	}
}

func (t *IterationInfo) Error(values ...*Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &IterationDataImpl{
			Properties: map[string]*Instance{
				"value": Tuple.Create(values...),
				"done":  Boolean.TRUE,
				"error": Boolean.TRUE,
			},
		},
	}
}

func (t *IterationInfo) Setup() {
	t.TypeInstance = Type.Create(Iteration.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
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
				"error": Boolean.FALSE,
			},
		},
	}
}

func (d *IterationDataType) OnNew(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*IterationDataImpl)

	this.Properties["done"] = Boolean.FALSE
	if len(args) == 0 {
		this.Properties["value"] = Tuple.Create(Boolean.FALSE)
	}

	if len(args) > 0 {
		this.Properties["value"] = Tuple.Create(args...)
	}

	return self
}

func (d *IterationDataType) OnSet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*IterationDataImpl)
	name := AsString(args[0])

	_, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	this.Properties[name] = args[1]
	return args[1]
}

func (d *IterationDataType) OnGet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*IterationDataImpl)
	name := AsString(args[0])

	value, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	return value
}

func (d *IterationDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self)
}

func (d *IterationDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*IterationDataImpl)
	if AsBool(this.error()) {
		return String.Create("<Iteration:error>")
	}

	if AsBool(this.done()) {
		return String.Create("<Iteration:done>")
	}

	return String.Create("<Iteration>")
}

func (d *IterationDataType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return self
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

func (impl *IterationDataImpl) error() *Instance {
	return impl.Properties["error"]
}
