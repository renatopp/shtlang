package runtime

import (
	"fmt"
	"sht/lang/ast"
)

var errorDT = &ErrorDataType{
	BaseDataType: BaseDataType{
		Name:        "Error",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
	},
}

var Error = &ErrorInfo{
	Type: errorDT,
}

// ----------------------------------------------------------------------------
// ERROR INFO
// ----------------------------------------------------------------------------
type ErrorInfo struct {
	Type DataType
}

func (t *ErrorInfo) Create(message string, a ...any) *Instance {
	msg := fmt.Sprintf(message, a...)
	return &Instance{
		Type: t.Type,
		Impl: ErrorDataImpl{
			Values: map[string]*Instance{
				"message": String.Create(msg),
			},
		},
	}
}

func (t *ErrorInfo) IncompatibleTypeOperation(op string, t1 *Instance, t2 *Instance) *Instance {
	return Error.Create("invalid operation with incompatible types: %s %s %s", t1.Type.GetName(), op, t2.Type.GetName())
}

func (t *ErrorInfo) InvalidOperation(op string, t1 *Instance) *Instance {
	return Error.Create("type %s does not implement operator %s", t1.Type.GetName(), op)
}

func (t *ErrorInfo) InvalidAction(action string, t1 *Instance) *Instance {
	return Error.Create("type %s does not implement action %s", t1.Type.GetName(), action)
}

// ----------------------------------------------------------------------------
// ERROR DATA TYPE
// ----------------------------------------------------------------------------
type ErrorDataType struct {
	BaseDataType
}

func (d *ErrorDataType) OnRepr(r *Runtime, args ...*Instance) *Instance {
	msg := AsString(args[0].Impl.(ErrorDataImpl).Values["message"])
	return String.Create("ERR! " + msg)
}

// ----------------------------------------------------------------------------
// ERROR DATA IMPL
// ----------------------------------------------------------------------------
type ErrorDataImpl struct {
	Values map[string]*Instance
}
