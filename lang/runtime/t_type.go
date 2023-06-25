package runtime

import (
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

var Type = _setupType()

type TypeInfo struct {
	Instance *Instance
	Type     *DataType
}

type TypeImpl struct {
	DataType *DataType
}

func _setupType() *TypeInfo {
	dataType := &DataType{
		Name:        "Type",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
		Meta:        map[meta.MetaName]Function{},
	}

	n := &TypeInfo{
		Instance: &Instance{
			Type:  dataType,
			Impl:  TypeImpl{DataType: dataType},
			Const: true,
		},
		Type: dataType,
	}

	return n
}

// ----------------------------------------------------------------------------
// Type Implementation
// ----------------------------------------------------------------------------
func (n TypeImpl) Repr() string {
	return n.DataType.Name
}

// ----------------------------------------------------------------------------
// Type Info
// ----------------------------------------------------------------------------
func (t *TypeInfo) Create(dataType *DataType, constant bool) *Instance {
	return &Instance{
		Type:  t.Type,
		Impl:  TypeImpl{DataType: dataType},
		Const: constant,
	}
}
