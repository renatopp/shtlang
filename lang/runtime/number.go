package runtime

import (
	"fmt"
	"math"
	"sht/lang/ast"
	"sht/lang/runtime/meta"
)

var NumberType = &DataType{
	Name:        "Number",
	Properties:  map[string]*ast.Node{},
	StaticFns:   map[string]Function{},
	InstanceFns: map[string]Function{},
	Meta:        map[meta.MetaName]Function{},
}
var Number = CreateType(NumberType)

func CreateNumber(value float64, constant bool) *Instance {
	return &Instance{
		Type:  NumberType,
		Impl:  NumberImpl{Value: value},
		Const: constant,
	}
}

type NumberImpl struct {
	Value float64
}

func (n NumberImpl) Repr() string {
	if math.Mod(n.Value, 1.0) == 0 {
		return fmt.Sprintf("%.0f", n.Value)
	} else {
		return fmt.Sprintf("%f", n.Value)
	}
}

func SetupNumber() {
	NumberType.Meta[meta.Add] = CreateNativeFunction(number_add)
}

func numberValue(instance *Instance) float64 {
	return instance.Impl.(NumberImpl).Value
}

func number_add(r *Runtime, args []*Instance) *Instance {
	if len(args) <= 1 {
		// return ERROR
	}

	if args[0].Type != args[1].Type {
		// return ERROR
	}

	this := numberValue(args[0])
	other := numberValue(args[1])

	return CreateNumber(this+other, false)
}
