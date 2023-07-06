package runtime

import (
	"sht/lang/ast"
)

var wildcardDT = &WildCardDataType{
	BaseDataType: BaseDataType{
		Name:        "WildCard",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var WildCard = &WildCardInfo{
	Type: wildcardDT,

	UNDERSCORE: &Instance{
		Type: wildcardDT,
		Impl: &WildCardDataImpl{},
	},
}

// ----------------------------------------------------------------------------
// WILDCARD INFO
// ----------------------------------------------------------------------------
type WildCardInfo struct {
	Type         DataType
	TypeInstance *Instance

	UNDERSCORE *Instance
}

func (t *WildCardInfo) Create() *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &WildCardDataImpl{},
	}
}

func (t *WildCardInfo) Setup() {
	t.TypeInstance = Type.Create(WildCard.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

// ----------------------------------------------------------------------------
// WILDCARD DATA TYPE
// ----------------------------------------------------------------------------
type WildCardDataType struct {
	BaseDataType
}

func (d *WildCardDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return String.Create("_")
}

// ----------------------------------------------------------------------------
// WILDCARD DATA IMPL
// ----------------------------------------------------------------------------
type WildCardDataImpl struct {
}
