package runtime

import "fmt"

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
	type_    string
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
	b.type_ = "string"
	return b
}

func (b *BuiltinArg) IsNumber() *BuiltinArg {
	b.type_ = "number"
	return b
}

func (b *BuiltinArg) IsBoolean() *BuiltinArg {
	b.type_ = "boolean"
	return b
}

func (b *BuiltinArg) IsFunction() *BuiltinArg {
	b.type_ = "function"
	return b
}

func (b *BuiltinArg) IsIterator() *BuiltinArg {
	b.type_ = "iterator"
	return b
}

func (b *BuiltinArg) IsList() *BuiltinArg {
	b.type_ = "list"
	return b
}

func (b *BuiltinArg) IsTuple() *BuiltinArg {
	b.type_ = "tuple"
	return b
}

func (b *BuiltinArg) Validate() (*Instance, error) {
	if b.index >= len(b.args) {
		if b.optional {
			return b.default_, nil
		}
		return nil, fmt.Errorf("Expecting argument at index %d, got none", b.index)
	}

	arg := b.args[b.index]
	if b.type_ != "" {
		switch b.type_ {
		case "string":
			if !arg.IsString() {
				return nil, fmt.Errorf("Expecting argument at index %d to be a string, got %s", b.index, arg.Type.GetName())
			}
		case "number":
			if !arg.IsNumber() {
				return nil, fmt.Errorf("Expecting argument at index %d to be a number, got %s", b.index, arg.Type.GetName())
			}
		case "boolean":
			if !arg.IsBoolean() {
				return nil, fmt.Errorf("Expecting argument at index %d to be a boolean, got %s", b.index, arg.Type.GetName())
			}
		case "function":
			if !arg.IsFunction() {
				return nil, fmt.Errorf("Expecting argument at index %d to be a function, got %s", b.index, arg.Type.GetName())
			}
		case "iterator":
			if !arg.IsIterator() {
				return nil, fmt.Errorf("Expecting argument at index %d to be an iterator, got %s", b.index, arg.Type.GetName())
			}
		case "list":
			if !arg.IsList() {
				return nil, fmt.Errorf("Expecting argument at index %d to be a list, got %s", b.index, arg.Type.GetName())
			}
		case "tuple":
			if !arg.IsTuple() {
				return nil, fmt.Errorf("Expecting argument at index %d to be a tuple, got %s", b.index, arg.Type.GetName())
			}
		}
	}

	return arg, nil
}
