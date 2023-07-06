package runtime

import (
	"sht/lang/ast"
	"strings"
)

var listDT = &ListDataType{
	BaseDataType: BaseDataType{
		Name:        "List",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var List = &ListInfo{
	Type: listDT,
}

// ----------------------------------------------------------------------------
// LIST INFO
// ----------------------------------------------------------------------------
type ListInfo struct {
	Type         DataType
	TypeInstance *Instance
}

func (t *ListInfo) Create(values ...*Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &ListDataImpl{
			Properties: map[string]*Instance{
				"default": ThrowFn,
			},
			Values: values,
		},
	}
}

func (t *ListInfo) Setup() {
	t.TypeInstance = Type.Create(List.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

// ----------------------------------------------------------------------------
// LIST DATA TYPE
// ----------------------------------------------------------------------------
type ListDataType struct {
	BaseDataType
}

func (d *ListDataType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	switch init := init.(type) {
	case *ast.ListInitializer:
		values := make([]*Instance, 0)
		for _, value := range init.Values {
			if spread, ok := value.(*ast.SpreadOut); ok {
				var e *Instance
				target := r.Eval(spread.Target, s)
				r.ResolveIterator(target, s, func(v *Instance, err *Instance) {
					if err != nil {
						e = err
					} else if v != nil {
						t := v.Impl.(*TupleDataImpl)
						values = append(values, t.Values...)
					}
				})
				if e != nil {
					return e
				}
				continue
			}
			values = append(values, r.Eval(value, s))
		}
		return List.Create(values...)
	default:
		return List.Create()
	}
}

func (d *ListDataType) OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	iter := self.Impl.(*IteratorDataImpl)
	next := iter.next()
	values := []*Instance{}
	for {
		tion := next.OnCall(r, s, self).AsIteration()

		if AsBool(tion.error()) {
			tuple := tion.value().AsTuple()
			return r.Throw(tuple.Values[0], s)

		} else if AsBool(tion.done()) {
			return List.Create(values...)

		} else {
			tuple := tion.value().AsTuple()
			values = append(values, tuple.Values[0])
		}
	}
}

func (d *ListDataType) OnNew(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*ListDataImpl)
	if len(args) > 0 {
		this.Properties["default"] = args[0]
	}

	return self
}

func (d *ListDataType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	cur := 0
	this := self.Impl.(*ListDataImpl)
	return Iterator.Create(
		Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if cur >= len(this.Values) {
				return Iteration.DONE
			}

			cur++
			return Iteration.Create(this.Values[cur-1])
		}),
	)
}

func (d *ListDataType) OnSet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsList()
	name := AsString(args[0])

	_, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	this.Properties[name] = args[1]
	return args[1]
}

func (d *ListDataType) OnGet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsList()
	name := AsString(args[0])

	value, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	return value
}

func (d *ListDataType) OnLen(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*ListDataImpl)
	return Number.Create(float64(len(this.Values)))
}

func (t *ListDataType) OnGetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*ListDataImpl)

	nargs := len(args)
	if nargs > 0 && !IsNumber(args[0]) {
		return r.Throw(Error.Create(s, "index of a list must be a number, '%s' provided", args[0].Type.GetName()), s)
	}

	if nargs > 1 && !IsNumber(args[1]) {
		return r.Throw(Error.Create(s, "index of a list must be a number, '%s' provided", args[2].Type.GetName()), s)
	}

	if nargs == 0 {
		return List.Create(this.Values...)
	}

	if nargs == 1 {
		idx := AsInteger(args[0])
		if idx >= len(this.Values) || idx < 0 {
			fn := this.default_()
			return fn.OnCall(r, s, String.Createf("list out of bounds for item '%d'", idx))
		}
		return this.Values[idx]
	}

	if nargs > 2 {
		return r.Throw(Error.Create(s, "list indexing accepts only 1 or 2 parameters, %d given", nargs-1), s)
	}

	size := len(this.Values)
	idx0 := AsInteger(args[0])
	idx1 := AsInteger(args[1])
	if idx0 > size {
		idx0 = size
	} else if idx1 < 0 {
		idx0 = -1
	}
	if idx1 > size {
		idx1 = size
	} else if idx1 < 0 {
		idx1 = -1
	}

	values := make([]*Instance, 0)

	if idx0 > idx1 {
		for i := idx0 - 1; i >= idx1; i-- {
			if i < 0 || i >= len(this.Values) {
				continue
			}
			values = append(values, this.Values[i])
		}
	} else if idx0 < idx1 {
		for i := idx0; i < idx1; i++ {
			if i < 0 || i >= len(this.Values) {
				continue
			}
			values = append(values, this.Values[i])
		}
	}

	return List.Create(values...)
}

func (t *ListDataType) OnSetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*ListDataImpl)

	nargs := len(args)
	if nargs != 2 {
		return r.Throw(Error.Create(s, "setItem receives only one index, '%d' provided", nargs), s)
	}

	if !IsNumber(args[0]) {
		return r.Throw(Error.Create(s, "index of a list must be a number, '%s' provided", args[0].Type.GetName()), s)
	}

	idx := AsInteger(args[0])
	if idx >= len(this.Values) || idx < 0 {
		return r.Throw(Error.Create(s, "list out of bounds for item '%d'", idx), s)
	}

	this.Values[idx] = args[1]

	return args[1]
}

func (d *ListDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self)
}

func (d *ListDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	list := self.Impl.(*ListDataImpl)

	var values []string
	for _, value := range list.Values {
		values = append(values, value.Repr())
	}

	return String.Create("[" + strings.Join(values, ", ") + "]")
}

// ----------------------------------------------------------------------------
// LIST DATA IMPL
// ----------------------------------------------------------------------------
type ListDataImpl struct {
	Properties map[string]*Instance
	Values     []*Instance
}

func (impl *ListDataImpl) default_() *Instance {
	return impl.Properties["default"]
}
