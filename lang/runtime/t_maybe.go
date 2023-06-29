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
	Type DataType
}

func (t *MaybeInfo) Create() *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &MaybeDataImpl{
			Value: nil,
			Error: nil,
		},
	}
}

// ----------------------------------------------------------------------------
// MAYBE DATA TYPE
// ----------------------------------------------------------------------------
type MaybeDataType struct {
	BaseDataType
}

func (d *MaybeDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args[0])
}

func (d *MaybeDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	maybe := args[0].Impl.(*MaybeDataImpl)

	if maybe.Error != nil {
		return String.Create("Maybe<error>")
	} else {
		return String.Create("Maybe<value>")
	}
}

// ----------------------------------------------------------------------------
// MAYBE DATA IMPL
// ----------------------------------------------------------------------------
type MaybeDataImpl struct {
	Value *Instance
	Error *Instance
}
