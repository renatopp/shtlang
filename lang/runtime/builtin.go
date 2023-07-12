package runtime

import (
	"fmt"
	"strings"
)

type BuiltinFunction struct {
	name   string
	params []*FunctionParam
}

func (b BuiltinFunction) as(impl MetaFunction) *Instance {
	return Function.CreateNative(b.name, b.params, impl)
}

// Creates a new builtin function with the given name and parameters.
func fn(name string, params ...*FunctionParam) *BuiltinFunction {
	return &BuiltinFunction{
		name:   name,
		params: params,
	}
}

// Creates a function parameter
func p(name string, args ...any) *FunctionParam {
	if len(args) == 0 {
		return &FunctionParam{Name: name}
	}

	var def *Instance
	if args[0] != nil {
		def = args[0].(*Instance)
	}

	if len(args) == 1 {
		return &FunctionParam{Name: name, Default: def}
	}

	return &FunctionParam{Name: name, Default: def, Spread: args[1].(bool)}
}

// Creates an iterator given the next function
func i(fn MetaFunction) *Instance {
	return Iterator.Create(
		Function.CreateNative("next", []*FunctionParam{}, fn),
	)
}

func throw(r *Runtime, s *Scope, msg string, args ...any) *Instance {
	return r.Throw(Error.Create(s, msg, args...), s)
}

// ----------------------------------------------------------------------------

type BuiltinArg struct {
	args     []*Instance
	index    int
	optional bool
	default_ *Instance
	types    []string
}

func arg(args []*Instance, index int) *BuiltinArg {
	return &BuiltinArg{
		args:  args,
		index: index,
	}
}

func (b *BuiltinArg) Optional(def ...*Instance) *BuiltinArg {
	b.optional = true
	if len(def) > 0 {
		b.default_ = def[0]
	}
	return b
}

func (b *BuiltinArg) IsString() *BuiltinArg {
	b.types = []string{"String"}
	return b
}

func (b *BuiltinArg) IsNumber() *BuiltinArg {
	b.types = []string{"Number"}
	return b
}

func (b *BuiltinArg) IsBoolean() *BuiltinArg {
	b.types = []string{"Boolean"}
	return b
}

func (b *BuiltinArg) IsFunction() *BuiltinArg {
	b.types = []string{"Function"}
	return b
}

func (b *BuiltinArg) IsIterator() *BuiltinArg {
	b.types = []string{"Iterator"}
	return b
}

func (b *BuiltinArg) IsList() *BuiltinArg {
	b.types = []string{"List"}
	return b
}

func (b *BuiltinArg) IsTuple() *BuiltinArg {
	b.types = []string{"Tuple"}
	return b
}

func (b *BuiltinArg) OrString() *BuiltinArg {
	b.types = append(b.types, "String")
	return b
}

func (b *BuiltinArg) OrNumber() *BuiltinArg {
	b.types = append(b.types, "Number")
	return b
}

func (b *BuiltinArg) OrBoolean() *BuiltinArg {
	b.types = append(b.types, "Boolean")
	return b
}

func (b *BuiltinArg) OrFunction() *BuiltinArg {
	b.types = append(b.types, "Function")
	return b
}

func (b *BuiltinArg) OrIterator() *BuiltinArg {
	b.types = append(b.types, "Iterator")
	return b
}

func (b *BuiltinArg) OrList() *BuiltinArg {
	b.types = append(b.types, "List")
	return b
}

func (b *BuiltinArg) OrTuple() *BuiltinArg {
	b.types = append(b.types, "Tuple")
	return b
}

func (b *BuiltinArg) Validate() (*Instance, error) {
	if b.index >= len(b.args) || b.args[b.index] == Boolean.FALSE {
		if b.optional {
			return b.default_, nil
		}
		return nil, fmt.Errorf("Expecting argument at index %d, got none", b.index)
	}

	arg := b.args[b.index]

	ok := false
	for _, t := range b.types {
		switch t {
		case "String":
			if arg.IsString() {
				ok = true
				break
			}
		case "Number":
			if arg.IsNumber() {
				ok = true
				break
			}
		case "Boolean":
			if arg.IsBoolean() {
				ok = true
				break
			}
		case "Function":
			if arg.IsFunction() {
				ok = true
				break
			}
		case "Iterator":
			if arg.IsIterator() {
				ok = true
				break
			}
		case "List":
			if arg.IsList() {
				ok = true
				break
			}
		case "Tuple":
			if arg.IsTuple() {
				ok = true
				break
			}
		}
	}

	if !ok {
		return nil, fmt.Errorf("Expecting argument at index '%d' to be a '%s', got '%s'", b.index, strings.Join(b.types, "', or '"), arg.Type.GetName())
	}

	return arg, nil
}
