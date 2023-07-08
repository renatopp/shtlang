package runtime

import (
	"sht/lang/ast"
	"strings"
)

var dictDT = &DictDataType{
	BaseDataType: BaseDataType{
		Name:        "Dict",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]*Instance{},
		InstanceFns: map[string]*Instance{},
	},
}

var Dict = &DictInfo{
	Type: dictDT,
}

// ----------------------------------------------------------------------------
// DICT INFO
// ----------------------------------------------------------------------------
type DictInfo struct {
	Type         DataType
	TypeInstance *Instance
}

func (t *DictInfo) Create(values map[string]*Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &DictDataImpl{
			Properties: map[string]*Instance{
				"default": ThrowFn,
			},
			Values: values,
		},
	}
}

func (t *DictInfo) Setup() {
	t.TypeInstance = Type.Create(Dict.Type)
	t.TypeInstance.Impl.(*TypeDataImpl).TypeInstance = t.TypeInstance
}

// ----------------------------------------------------------------------------
// DICT DATA TYPE
// ----------------------------------------------------------------------------
type DictDataType struct {
	BaseDataType
}

func (d *DictDataType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	switch init := init.(type) {
	case *ast.MapInitializer:
		values := map[string]*Instance{}

		for k, v := range init.Values {
			values[k] = r.Eval(v, s)
		}

		return Dict.Create(values)
	case *ast.ListInitializer:
		return r.Throw(Error.Create(s, "type '%s' does not allow instantiation with list initializer", d.Name), s)
	default:
		return Dict.Create(map[string]*Instance{})
	}
}

func (d *DictDataType) OnTo(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	iter := self.AsIterator()
	next := iter.next()
	values := map[string]*Instance{}
	for {
		tion := next.OnCall(r, s, self).AsIteration()

		if AsBool(tion.error()) {
			tuple := tion.value().AsTuple()
			return r.Throw(tuple.Values[0], s)

		} else if AsBool(tion.done()) {
			return Dict.Create(values)

		} else {
			tuple := tion.value().AsTuple()
			if len(tuple.Values) < 2 {
				return r.Throw(Error.Create(s, "invalid tuple for dict, dict requires two elements as (key, value)"), s)
			}

			values[AsString(tuple.Values[0])] = tuple.Values[1]
		}
	}
}

func (d *DictDataType) OnNew(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsDict()
	if len(args) > 0 {
		this.Properties["default"] = args[0]
	}

	return self
}

func (d *DictDataType) OnIter(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsDict()

	keys := []string{}
	for k := range this.Values {
		keys = append(keys, k)
	}

	cur := 0
	return Iterator.Create(
		Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
			if cur >= len(keys) {
				return Iteration.DONE
			}

			cur++
			k := keys[cur-1]
			return Iteration.Create(
				String.Create(k),
				this.Values[k],
			)
		}),
	)
}

func (d *DictDataType) OnSet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsDict()
	name := AsString(args[0])

	_, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	this.Properties[name] = args[1]
	return args[1]
}

func (d *DictDataType) OnGet(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsDict()
	name := AsString(args[0])

	value, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	return value
}

func (d *DictDataType) OnLen(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.Impl.(*DictDataImpl)
	return Number.Create(float64(len(this.Values)))
}

func (t *DictDataType) OnGetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsDict()

	if len(args) != 1 {
		return r.Throw(Error.Create(s, "dict getItem receives only one index, '%d' provided", len(args)), s)
	}

	key := AsString(args[0])

	if _, has := this.Values[key]; !has {
		val := this.default_().OnCall(r, s, self)
		if s.IsInterruptedAs(FlowRaise) {
			return val
		}

		this.Values[key] = val
	}

	return this.Values[key]
}

func (t *DictDataType) OnSetItem(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	this := self.AsDict()

	nargs := len(args)
	if nargs != 2 {
		return r.Throw(Error.Create(s, "dict setItem receives only one index, '%d' provided", nargs), s)
	}

	name := AsString(args[0])
	this.Values[name] = args[1]
	return args[1]
}

func (d *DictDataType) OnString(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	return d.OnRepr(r, s, self)
}

func (d *DictDataType) OnRepr(r *Runtime, s *Scope, self *Instance, args ...*Instance) *Instance {
	dict := self.AsDict()

	var values []string
	for key, value := range dict.Values {
		values = append(values, key+": "+value.Repr())
	}

	return String.Create("{" + strings.Join(values, ", ") + "}")
}

// ----------------------------------------------------------------------------
// DICT DATA IMPL
// ----------------------------------------------------------------------------
type DictDataImpl struct {
	Properties map[string]*Instance
	Values     map[string]*Instance
}

func (impl *DictDataImpl) default_() *Instance {
	return impl.Properties["default"]
}
