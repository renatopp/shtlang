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

func (d *NumberDataType) OnBoolean(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return Boolean.Create(AsNumber(args[0]) != 0)
}

func (d *NumberDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args...)
}

func (d *NumberDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	v := AsNumber(args[0])

	if math.Mod(v, 1.0) == 0 {
		return String.Create(fmt.Sprintf("%.0f", v))
	}

	return String.Create(fmt.Sprintf("%f", v))
}

func (d *NumberDataType) OnAdd(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, "+", args[0], args[1])
	}

	return Number.Create(AsNumber(args[0]) + AsNumber(args[1]))
}

func (d *NumberDataType) OnSub(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, "-", args[0], args[1])
	}

	return Number.Create(AsNumber(args[0]) - AsNumber(args[1]))
}

func (d *NumberDataType) OnMul(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, "*", args[0], args[1])
	}

	return Number.Create(AsNumber(args[0]) * AsNumber(args[1]))
}

func (d *NumberDataType) OnDiv(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, "/", args[0], args[1])
	}

	return Number.Create(AsNumber(args[0]) / AsNumber(args[1]))
}

func (d *NumberDataType) OnIntDiv(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, "//", args[0], args[1])
	}

	return Number.Create(math.Floor(AsNumber(args[0]) / AsNumber(args[1])))
}

func (d *NumberDataType) OnMod(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, "%", args[0], args[1])
	}

	return Number.Create(math.Mod(AsNumber(args[0]), AsNumber(args[1])))
}

func (d *NumberDataType) OnPow(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, "**", args[0], args[1])
	}

	return Number.Create(math.Pow(AsNumber(args[0]), AsNumber(args[1])))
}

func (d *NumberDataType) OnEq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Boolean.FALSE
	}
	return Boolean.Create(AsNumber(args[0]) == AsNumber(args[1]))
}

func (d *NumberDataType) OnNeq(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Boolean.TRUE
	}
	return Boolean.Create(AsNumber(args[0]) != AsNumber(args[1]))
}

func (d *NumberDataType) OnGt(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, ">", args[0], args[1])
	}

	return Boolean.Create(AsNumber(args[0]) > AsNumber(args[1]))
}

func (d *NumberDataType) OnLt(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, "<", args[0], args[1])
	}

	return Boolean.Create(AsNumber(args[0]) < AsNumber(args[1]))
}

func (d *NumberDataType) OnGte(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, ">=", args[0], args[1])
	}

	return Boolean.Create(AsNumber(args[0]) >= AsNumber(args[1]))
}

func (d *NumberDataType) OnLte(r *Runtime, s *Scope, args ...*Instance) *Instance {
	if args[0].Type != args[1].Type {
		return Error.IncompatibleTypeOperation(s, "<=", args[0], args[1])
	}

	return Boolean.Create(AsNumber(args[0]) <= AsNumber(args[1]))
}

func (d *NumberDataType) OnPostInc(r *Runtime, s *Scope, args ...*Instance) *Instance {
	impl := args[0].Impl.(NumberDataImpl)
	old := impl.Value
	impl.Value += 1
	return Number.Create(old)
}

func (d *NumberDataType) OnPostDec(r *Runtime, s *Scope, args ...*Instance) *Instance {
	impl := args[0].Impl.(NumberDataImpl)
	old := impl.Value
	impl.Value -= 1
	return Number.Create(old)
}

func (d *NumberDataType) OnNot(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return Boolean.Create(!AsBool(args[0]))
}

func (d *NumberDataType) OnNeg(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return Number.Create(-AsNumber(args[0]))
}

func (d *NumberDataType) OnPos(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return args[0]
}

// ----------------------------------------------------------------------------
// NUMBER DATA IMPL
// ----------------------------------------------------------------------------
type NumberDataImpl struct {
	Value float64
}
