package runtime

import (
	"sht/lang/ast"
)

var moduleDT = &ModuleDataType{
	BaseDataType: BaseDataType{
		Name:        "Module",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var Module = &ModuleInfo{
	Type: moduleDT,
}

// ----------------------------------------------------------------------------
// MODULE INFO
// ----------------------------------------------------------------------------
type ModuleInfo struct {
	Type         DataType
	TypeInstance *Instance
}

func (t *ModuleInfo) Create(name string) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &ModuleDataImpl{
			Name:  name,
			Scope: CreateScope(nil, nil),
		},
	}
}

func (t *ModuleInfo) Setup() {
	t.TypeInstance = Type.Create(Module.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

func (t *ModuleInfo) Add(impl *Instance, name string, value *Instance) {
	impl.Impl.(*ModuleDataImpl).Scope.Set(name, value)
}

// ----------------------------------------------------------------------------
// MODULE DATA TYPE
// ----------------------------------------------------------------------------
type ModuleDataType struct {
	BaseDataType
}

func (d *ModuleDataType) OnGet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*ModuleDataImpl)
	name := AsString(args[0])

	value, has := this.Scope.GetInScope(name)
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	return value
}

func (d *ModuleDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self)
}

func (d *ModuleDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return String.Createf("<Module:%s>", self.Impl.(*ModuleDataImpl).Name)
}

// ----------------------------------------------------------------------------
// MODULE DATA IMPL
// ----------------------------------------------------------------------------
type ModuleDataImpl struct {
	Name  string
	Scope *Scope
}
