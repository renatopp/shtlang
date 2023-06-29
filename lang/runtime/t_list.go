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
	Type DataType
}

func (t *ListInfo) Create(values ...*Instance) *Instance {
	return &Instance{
		Type: t.Type,
		Impl: &ListDataImpl{
			Values: values,
		},
	}
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
		return r.Throw(Error.Create(s, "invalid initializer for list"), s)
	}
}

func (d *ListDataType) OnNew(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return args[0]
}

func (d *ListDataType) OnIter(r *Runtime, s *Scope, args ...*Instance) *Instance {
	cur := 0
	this := args[0].Impl.(*ListDataImpl)
	return Iterator.Create(
		Function.CreateNative("next", []*FunctionParam{}, func(r *Runtime, s *Scope, args ...*Instance) *Instance {
			if cur >= len(this.Values) {
				return Iteration.DONE
			}

			cur++
			return Iteration.Create(this.Values[cur-1])
		}),
	)
}

func (d *ListDataType) OnLen(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*ListDataImpl)
	return Number.Create(float64(len(this.Values)))
}

func (t *ListDataType) OnGetItem(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*ListDataImpl)

	nargs := len(args)
	if nargs > 1 && !IsNumber(args[1]) {
		return r.Throw(Error.Create(s, "index of a list must be a number, '%s' provided", args[1].Type.GetName()), s)
	}

	if nargs > 2 && !IsNumber(args[2]) {
		return r.Throw(Error.Create(s, "index of a list must be a number, '%s' provided", args[2].Type.GetName()), s)
	}

	if nargs > 3 {
		return r.Throw(Error.Create(s, "list indexing accepts only 1 or 2 parameters, %d given", nargs-1), s)
	}

	if nargs == 2 {
		return this.Values[AsInteger(args[1])]
	}

	if nargs == 1 {
		return List.Create(this.Values...)
	}

	size := len(this.Values)
	idx0 := AsInteger(args[1])
	if idx0 < 0 || idx0 > size-1 {
		return r.Throw(Error.Create(s, "first index '%d' of list slicing out of bounds", idx0), s)
	}

	idx1 := AsInteger(args[2])
	if idx1 < 0 || idx1 > size {
		return r.Throw(Error.Create(s, "second index '%d' of list slicing out of bounds", idx0), s)
	}

	if idx1 <= idx0 {
		return r.Throw(Error.Create(s, "second index '%d' of list slicing must be greater than the first '%d'", idx1, idx0), s)
	}

	values := make([]*Instance, 0)
	for _, v := range this.Values[idx0:idx1] {
		values = append(values, v)
	}

	return List.Create(values...)
}

func (d *ListDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return d.OnRepr(r, s, args[0])
}

func (d *ListDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	list := args[0].Impl.(*ListDataImpl)

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
	Values []*Instance
}
