package runtime

import (
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

type DataType struct {
	Name        string
	Properties  map[string]*ast.Node
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
	Call(r *Runtime, args []*Instance) *Instance
}
