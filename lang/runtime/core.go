package runtime

import (
	"fmt"
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

type DataType struct {
	Name        string
	Properties  map[string]ast.Node
	StaticFns   map[string]Function
	InstanceFns map[string]Function
	Meta        map[meta.MetaName]Function
}

type DataImpl interface {
	Repr() string
}

type Instance struct {
	Type  *DataType
	Impl  DataImpl
	Const bool
}

type Function interface {
	Call(r *Runtime, args ...*Instance) *Instance
}

func InvalidOperationType(op string, t1 *Instance, t2 *Instance) *Instance {
	msg := fmt.Sprintf("invalid operation with incompatible types: %s %s %s", t1.Type.Name, op, t2.Type.Name)
	return Error.Create(msg, false)
}

func InvalidOperation(op string, t1 *Instance) *Instance {
	msg := fmt.Sprintf("type %s does not implement operator %s", t1.Type.Name, op)
	return Error.Create(msg, false)
}

func NotImplemented(action string, t1 *Instance) *Instance {
	msg := fmt.Sprintf("type %s does not implement action %s", t1.Type.Name, action)
	return Error.Create(msg, false)
}
