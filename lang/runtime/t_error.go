package runtime

import (
	"fmt"
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

var Error = _setupError()

type ErrorInfo struct {
	Instance *Instance
	Type     *DataType
}

type ErrorImpl struct {
	Values map[string]*Instance
}

func _setupError() *ErrorInfo {
	dataType := &DataType{
		Name:        "Error",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
		Meta:        map[meta.MetaName]Function{},
	}

	n := &ErrorInfo{
		Instance: Type.Create(dataType, true),
		Type:     dataType,
	}

	dataType.Properties["message"] = &ast.String{Value: ""}
	return n
}

// ----------------------------------------------------------------------------
// Error Implementation
// ----------------------------------------------------------------------------
func (n ErrorImpl) Repr() string {
	return fmt.Sprintf("ERR! %s", n.Values["message"].Impl.(StringImpl).Value)
}

// ----------------------------------------------------------------------------
// Error Info
// ----------------------------------------------------------------------------
func (n *ErrorInfo) Create(msg string, constant bool) *Instance {
	return &Instance{
		Type: n.Type,
		Impl: ErrorImpl{
			Values: map[string]*Instance{
				"message": String.Create(msg, false),
			},
		},
		Const: constant,
	}
}
