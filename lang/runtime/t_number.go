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
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var Number = &NumberInfo{
	Type: numberDT,

	ZERO: &Instance{
		Type: numberDT,
		Impl: &NumberDataImpl{
			Value: 0,
		},
	},
	ONE: &Instance{
		Type: numberDT,
		Impl: &NumberDataImpl{
			Value: 1,
		},
	},
	TWO: &Instance{
		Type: numberDT,
		Impl: &NumberDataImpl{
			Value: 2,
		},
	},
}

// ----------------------------------------------------------------------------
// NUMBER INFO
// ----------------------------------------------------------------------------
type NumberInfo struct {
	Type         DataType
	TypeInstance *Instance

	ZERO *Instance
	ONE  *Instance
	TWO  *Instance
}

func (t *NumberInfo) Create(value float64) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &NumberDataImpl{
			Value: value,
		},
	}
}

func (t *NumberInfo) Setup() {
	t.TypeInstance = Type.Create(Number.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

// ----------------------------------------------------------------------------
// NUMBER DATA TYPE
// ----------------------------------------------------------------------------
type NumberDataType struct {
	BaseDataType
}

func (d *NumberDataType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	switch init.(type) {
	case *ast.ListInitializer, *ast.MapInitializer:
		return r.Throw(Error.Create(s, "type '%s' does not allow instantiation with initializer", d.Name), s)
	default:
		return Number.Create(0)
	}
}

func (d *NumberDataType) OnNew(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if len(args) == 0 {
		return self
	}
	return args[0].OnNumber(r, s)
}

func (d *NumberDataType) OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	iter := self.Impl.(*IteratorDataImpl)
	next := iter.next()
	tion := next.OnCall(r, s, self).Impl.(*IterationDataImpl)

	if tion.error() == Boolean.TRUE {
		tuple := tion.value().AsTuple()
		return r.Throw(tuple.Values[0], s)

	} else if tion.done() == Boolean.TRUE {
		return r.Throw(Error.Create(s, "The iteration has been finished"), s)

	} else {
		tuple := tion.value().Impl.(*TupleDataImpl)
		if tuple.Values[0].Type != Number.Type {
			return r.Throw(Error.Create(s, "Cannot convert to number"), s)
		}
		return Number.Create(AsNumber(tuple.Values[0]))
	}
}

func (d *NumberDataType) OnNumber(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return self
}

func (d *NumberDataType) OnBoolean(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return Boolean.Create(AsNumber(self) != 0)
}

func (d *NumberDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self, args...)
}

func (d *NumberDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	v := AsNumber(self)

	if math.Mod(v, 1.0) == 0 {
		return String.Create(fmt.Sprintf("%.0f", v))
	}

	return String.Create(fmt.Sprintf("%f", v))
}

func (d *NumberDataType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	cur := 0
	this := self.Impl.(*NumberDataImpl)
	return Iterator.Create(
		Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if cur >= 1 {
				return Iteration.DONE
			}

			cur++
			return Iteration.Create(Number.Create(this.Value))
		}),
	)
}

func (d *NumberDataType) OnAdd(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "+", self, args[0]), s)
	}

	return Number.Create(AsNumber(self) + AsNumber(args[0]))
}

func (d *NumberDataType) OnSub(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "-", self, args[0]), s)
	}

	return Number.Create(AsNumber(self) - AsNumber(args[0]))
}

func (d *NumberDataType) OnMul(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "*", self, args[0]), s)
	}

	return Number.Create(AsNumber(self) * AsNumber(args[0]))
}

func (d *NumberDataType) OnDiv(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "/", self, args[0]), s)
	}

	return Number.Create(AsNumber(self) / AsNumber(args[0]))
}

func (d *NumberDataType) OnIntDiv(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "//", self, args[0]), s)
	}

	return Number.Create(math.Floor(AsNumber(self) / AsNumber(args[0])))
}

func (d *NumberDataType) OnMod(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "%", self, args[0]), s)
	}

	return Number.Create(math.Mod(AsNumber(self), AsNumber(args[0])))
}

func (d *NumberDataType) OnPow(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "**", self, args[0]), s)
	}

	return Number.Create(math.Pow(AsNumber(self), AsNumber(args[0])))
}

func (d *NumberDataType) OnEq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return Boolean.FALSE
	}
	return Boolean.Create(AsNumber(self) == AsNumber(args[0]))
}

func (d *NumberDataType) OnNeq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return Boolean.TRUE
	}
	return Boolean.Create(AsNumber(self) != AsNumber(args[0]))
}

func (d *NumberDataType) OnGt(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, ">", self, args[0]), s)
	}

	return Boolean.Create(AsNumber(self) > AsNumber(args[0]))
}

func (d *NumberDataType) OnLt(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "<", self, args[0]), s)
	}

	return Boolean.Create(AsNumber(self) < AsNumber(args[0]))
}

func (d *NumberDataType) OnGte(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, ">=", self, args[0]), s)
	}

	return Boolean.Create(AsNumber(self) >= AsNumber(args[0]))
}

func (d *NumberDataType) OnLte(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "<=", self, args[0]), s)
	}

	return Boolean.Create(AsNumber(self) <= AsNumber(args[0]))
}

func (d *NumberDataType) OnNot(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return Boolean.Create(!AsBool(self))
}

func (d *NumberDataType) OnNeg(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return Number.Create(-AsNumber(self))
}

func (d *NumberDataType) OnPos(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return self
}

// ----------------------------------------------------------------------------
// NUMBER DATA IMPL
// ----------------------------------------------------------------------------
type NumberDataImpl struct {
	Value float64
}
