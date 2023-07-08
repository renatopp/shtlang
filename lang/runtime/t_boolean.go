package runtime

import (
	"sht/lang/ast"
)

var booleanDT = &BooleanDataType{
	BaseDataType: BaseDataType{
		Name:        "Boolean",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var Boolean = &BooleanInfo{
	Type: booleanDT,

	TRUE: &Instance{
		Type: booleanDT,
		Impl: &BooleanDataImpl{
			Value: true,
		},
	},
	FALSE: &Instance{
		Type: booleanDT,
		Impl: &BooleanDataImpl{
			Value: false,
		},
	},
}

// ----------------------------------------------------------------------------
// BOOLEAN INFO
// ----------------------------------------------------------------------------
type BooleanInfo struct {
	Type         DataType
	TypeInstance *Instance

	TRUE  *Instance
	FALSE *Instance
}

func (t *BooleanInfo) Create(value bool) *Instance {
	if value {
		return t.TRUE
	} else {
		return t.FALSE
	}
}

func (t *BooleanInfo) Setup() {
	t.TypeInstance = Type.Create(Boolean.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

// ----------------------------------------------------------------------------
// BOOLEAN DATA TYPE
// ----------------------------------------------------------------------------
type BooleanDataType struct {
	BaseDataType
}

func (d *BooleanDataType) OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	iter := self.Impl.(*IteratorDataImpl)
	next := iter.next()
	tion := next.OnCall(r, s, self).Impl.(*IterationDataImpl)

	if tion.error() == Boolean.TRUE {
		tuple := tion.value().AsTuple()
		return r.Throw(tuple.Values[0], s)

	} else if tion.done() == Boolean.TRUE {
		return r.Throw(Error.Create(s, "The iteration has been finished"), s)

	} else {
		tuple := tion.value().AsTuple()
		return Boolean.Create(AsBool(tuple.Values[0]))
	}
}

func (d *BooleanDataType) OnBoolean(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return self
}

func (d *BooleanDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self, args...)
}

func (d *BooleanDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	v := AsBool(self)
	if v {
		return String.Create("true")
	}

	return String.Create("false")
}

func (d *BooleanDataType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	cur := 0
	this := self.Impl.(*BooleanDataImpl)
	return Iterator.Create(
		Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if cur >= 1 {
				return Iteration.DONE
			}

			cur++
			return Iteration.Create(Boolean.Create(this.Value))
		}),
	)
}

func (d *BooleanDataType) OnNot(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	v := AsBool(self)
	return Boolean.Create(!v)
}

func (n *BooleanDataType) OnEq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return Boolean.FALSE
	}

	this := AsBool(self)
	other := AsBool(args[0])

	if this == other {
		return Boolean.TRUE
	} else {
		return Boolean.FALSE
	}
}

func (n *BooleanDataType) OnNeq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return Boolean.FALSE
	}

	this := AsBool(self)
	other := AsBool(args[0])

	if this != other {
		return Boolean.TRUE
	} else {
		return Boolean.FALSE
	}
}

// ----------------------------------------------------------------------------
// BOOLEAN DATA IMPL
// ----------------------------------------------------------------------------
type BooleanDataImpl struct {
	Value bool
}
