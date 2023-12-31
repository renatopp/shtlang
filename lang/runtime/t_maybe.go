package runtime

import (
	"sht/lang/ast"
)

var maybeDT = &MaybeDataType{
	BaseDataType: BaseDataType{
		Name:        "Maybe",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var Maybe = &MaybeInfo{
	Type: maybeDT,
}

// ----------------------------------------------------------------------------
// MAYBE INFO
// ----------------------------------------------------------------------------
type MaybeInfo struct {
	Type         DataType
	TypeInstance *Instance
}

func (t *MaybeInfo) Create(value *Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &MaybeDataImpl{
			Value: value,
			Error: nil,
		},
	}
}

func (t *MaybeInfo) CreateError(err *Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &MaybeDataImpl{
			Value: nil,
			Error: err,
		},
	}
}

func (t *MaybeInfo) Setup() {
	t.TypeInstance = Type.Create(Maybe.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

// ----------------------------------------------------------------------------
// MAYBE DATA TYPE
// ----------------------------------------------------------------------------
type MaybeDataType struct {
	BaseDataType
}

func (d *MaybeDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self)
}

func (d *MaybeDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	maybe := self.Impl.(*MaybeDataImpl)

	if maybe.Error != nil {
		return String.Create("<Maybe:error>")
	} else {
		return String.Create("<Maybe:value>")
	}
}

// ----------------------------------------------------------------------------
// MAYBE DATA IMPL
// ----------------------------------------------------------------------------
type MaybeDataImpl struct {
	Value *Instance
	Error *Instance
}
