package runtime

import (
	"fmt"
	"math"
	"sht/lang/ast"
)

var numberDT = &NumberDataType{
	BaseDataType: BaseDataType{
		Name:        "Number",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
	},
}

var Number = &NumberInfo{
	Type: numberDT,

	ZERO: &Instance{
		Type: numberDT,
		Impl: NumberDataImpl{
			Value: 0,
		},
	},
	ONE: &Instance{
		Type: numberDT,
		Impl: NumberDataImpl{
			Value: 1,
		},
	},
}

// ----------------------------------------------------------------------------
// NUMBER INFO
// ----------------------------------------------------------------------------
type NumberInfo struct {
	Type DataType

	ZERO *Instance
	ONE  *Instance
}

func (t *NumberInfo) Create(value float64) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: NumberDataImpl{
			Value: value,
		},
	}
}

// ----------------------------------------------------------------------------
// NUMBER DATA TYPE
// ----------------------------------------------------------------------------
type NumberDataType struct {
	BaseDataType
}

func (d *NumberDataType) OnBoolean(r *Runtime, args ...*Instance) *Instance {
	return Boolean.Create(AsBool(args[0]))
}

func (d *NumberDataType) OnString(r *Runtime, args ...*Instance) *Instance {
	return d.OnRepr(r, args...)
}

func (d *NumberDataType) OnRepr(r *Runtime, args ...*Instance) *Instance {
	v := AsNumber(args[0])

	if math.Mod(v, 1.0) == 0 {
		return String.Create(fmt.Sprintf("%.0f", v))
	}

	return String.Create(fmt.Sprintf("%f", v))
}

func (d *NumberDataType) OnNeg(r *Runtime, args ...*Instance) *Instance {
	return Number.Create(-AsNumber(args[0]))
}

// ----------------------------------------------------------------------------
// NUMBER DATA IMPL
// ----------------------------------------------------------------------------
type NumberDataImpl struct {
	Value float64
}
