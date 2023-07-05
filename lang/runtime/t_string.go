package runtime

import (
	"fmt"
	"sht/lang/ast"
	"strings"
)

var stringDT = &StringDataType{
	BaseDataType: BaseDataType{
		Name:        "String",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var String = &StringInfo{
	Type: stringDT,

	EMPTY: &Instance{
		Type: stringDT,
		Impl: StringDataImpl{
			Value: "",
		},
	},
}

// ----------------------------------------------------------------------------
// STRING INFO
// ----------------------------------------------------------------------------
type StringInfo struct {
	Type         DataType
	TypeInstance *Instance

	EMPTY *Instance
}

func (t *StringInfo) Create(value string) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: StringDataImpl{
			Value: value,
		},
	}
}

func (t *StringInfo) Createf(value string, v ...any) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: StringDataImpl{
			Value: fmt.Sprintf(value, v...),
		},
	}
}

func (t *StringInfo) Setup() {
	t.TypeInstance = Type.Create(String.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

// ----------------------------------------------------------------------------
// STRING DATA TYPE
// ----------------------------------------------------------------------------
type StringDataType struct {
	BaseDataType
}

func (d *StringDataType) OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	iter := self.Impl.(*IteratorDataImpl)
	next := iter.next()
	values := []*Instance{}
	for {

		tion := next.OnCall(r, s, self).Impl.(*IterationDataImpl)

		if tion.error() == Boolean.TRUE {
			return List.Create()

		} else if tion.done() == Boolean.TRUE {
			builder := strings.Builder{}
			for _, value := range values {
				builder.WriteString(AsString(value))
			}
			return String.Create(builder.String())

		} else {
			tuple := tion.value().Impl.(*TupleDataImpl)
			values = append(values, tuple.Values[0])
		}
	}
}

func (d *StringDataType) OnLen(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := AsString(self)
	return Number.Create(float64(len(this)))
}

func (d *StringDataType) OnAdd(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return r.Throw(Error.IncompatibleTypeOperation(s, "+", self, args[0]), s)
	}

	return String.Create(AsString(self) + AsString(args[0]))
}

func (d *StringDataType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	cur := 0
	this := self.Impl.(StringDataImpl)
	return Iterator.Create(
		Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if cur >= len(this.Value) {
				return Iteration.DONE
			}

			val := Iteration.Create(String.Create(this.Value[cur : cur+1]))
			cur++
			return val
		}),
	)
}

func (d *StringDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return self
}

func (d *StringDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return self
}

func (d *StringDataType) OnBoolean(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := AsString(self)

	if this == "" {
		return Boolean.FALSE
	}

	return Boolean.TRUE
}

func (d *StringDataType) OnNot(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := AsBool(d.OnBoolean(r, s, self, args...))

	if this {
		return Boolean.FALSE
	}

	return Boolean.TRUE
}

func (n *StringDataType) OnEq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return Boolean.FALSE
	}

	this := AsString(self)
	other := AsString(args[0])

	if this == other {
		return Boolean.TRUE
	} else {
		return Boolean.FALSE
	}
}

func (n *StringDataType) OnNeq(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	if self.Type != args[0].Type {
		return Boolean.FALSE
	}

	this := AsString(self)
	other := AsString(args[0])

	if this != other {
		return Boolean.TRUE
	} else {
		return Boolean.FALSE
	}
}

func (n *StringDataType) OnGetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := AsString(self)

	nargs := len(args)
	if nargs > 0 && !IsNumber(args[0]) {
		return r.Throw(Error.Create(s, "index of a string must be a number, '%s' provided", args[0].Type.GetName()), s)
	}

	if nargs > 1 && !IsNumber(args[1]) {
		return r.Throw(Error.Create(s, "index of a string must be a number, '%s' provided", args[2].Type.GetName()), s)
	}

	if nargs > 2 {
		return r.Throw(Error.Create(s, "string indexing accepts only 0, 1 or 2 parameters, %d given", nargs-1), s)
	}

	idx0 := 0
	if nargs >= 1 {
		idx0 = int(AsNumber(args[0]))
		if idx0 < 0 {
			idx0 = 0
		}
		if idx0 >= len(this) {
			idx0 = len(this) - 1
		}
	}

	idx1 := idx0 + 1
	if nargs >= 2 {
		idx1 = int(AsNumber(args[1]))
		if idx1 < 1 {
			idx1 = 1
		}
		if idx1 > len(this) {
			idx1 = len(this)
		}
	}

	if nargs == 0 {
		idx0 = 0
		idx1 = len(this)
	}

	if idx1 <= idx0 {
		return r.Throw(Error.Create(s, "second index '%d' of string slicing must be greater than the first '%d'", idx1, idx0), s)
	}

	return String.Create(this[idx0:idx1])
}

// ----------------------------------------------------------------------------
// STRING DATA IMPL
// ----------------------------------------------------------------------------
type StringDataImpl struct {
	Value string
}
