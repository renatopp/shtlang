package runtime

import (
	"fmt"
	"math"
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

var Number = _setupNumber()

type NumberInfo struct {
	Instance *Instance
	Type     *DataType

	ZERO *Instance
	ONE  *Instance
}

type NumberImpl struct {
	Value float64
}

func _setupNumber() *NumberInfo {
	dataType := &DataType{
		Name:        "Number",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
		Meta:        map[meta.MetaName]Function{},
	}

	n := &NumberInfo{
		Instance: Type.Create(dataType, true),
		Type:     dataType,
	}

	n.ZERO = n.Create(0, true)
	n.ONE = n.Create(1, true)
	dataType.Meta[meta.Add] = CreateNativeFunction(n.MetaAdd)
	dataType.Meta[meta.Sub] = CreateNativeFunction(n.MetaSub)

	return n
}

// ----------------------------------------------------------------------------
// Number Implementation
// ----------------------------------------------------------------------------
func (n NumberImpl) Repr() string {
	if math.Mod(n.Value, 1.0) == 0 {
		return fmt.Sprintf("%.0f", n.Value)
	} else {
		return fmt.Sprintf("%f", n.Value)
	}
}

// ----------------------------------------------------------------------------
// Number Info
// ----------------------------------------------------------------------------
func (n *NumberInfo) Create(value float64, constant bool) *Instance {
	return &Instance{
		Type:  n.Type,
		Impl:  NumberImpl{Value: value},
		Const: constant,
	}
}

func (n *NumberInfo) val(instance *Instance) float64 {
	return instance.Impl.(NumberImpl).Value
}

func (n *NumberInfo) MetaAdd(r *Runtime, args ...*Instance) *Instance {
	if len(args) <= 1 {
		return Error.Create("not enough arguments", false)
	}

	if args[0].Type != args[1].Type {
		return Error.Create("invalid types", false)
	}

	this := n.val(args[0])
	other := n.val(args[1])

	return n.Create(this+other, false)
}

func (n *NumberInfo) MetaSub(r *Runtime, args ...*Instance) *Instance {
	if len(args) <= 1 {
		return Error.Create("not enough arguments", false)
	}

	if args[0].Type != args[1].Type {
		return Error.Create("invalid types", false)
	}

	this := n.val(args[0])
	other := n.val(args[1])

	return n.Create(this-other, false)
}
