package runtime

import (
	"fmt"
	"sht/lang/ast"
	"sht/lang/tokens"
	"strings"
)

var errorDT = &ErrorDataType{
	BaseDataType: BaseDataType{
		Name:        "Error",
		Properties:  map[string]ast.Node{},
		StaticFns:   map[string]Function{},
		InstanceFns: map[string]Function{},
	},
}

var Error = &ErrorInfo{
	Type: errorDT,
}

// ----------------------------------------------------------------------------
// ERROR INFO
// ----------------------------------------------------------------------------
type ErrorInfo struct {
	Type DataType
}

func (t *ErrorInfo) Create(s *Scope, message string, a ...any) *Instance {
	msg := fmt.Sprintf(message, a...)
	return &Instance{
		Type: t.Type,
		Impl: &ErrorDataImpl{
			Properties: map[string]*Instance{
				"message": String.Create(msg),
				"trace":   String.Create(t.StackTrace(s)),
			},
		},
	}
}

func (t *ErrorInfo) StackTrace(s *Scope) string {
	stack := s.CallStack()

	trace := strings.Builder{}
	total := 0
	for i := len(stack) - 1; i >= 0; i-- {
		total++
		if total > 10 {
			trace.WriteString("     ...\n")
			break
		}

		scope := stack[i]
		trace.WriteString("     at ")

		fn, _ := scope.GetInScope(SCOPE_FN_KEY)
		if fn == nil {
			trace.WriteString("<global>")
		} else {
			switch fn := fn.Value.Impl.(type) {
			case *CustomFunctionDataImpl:
				trace.WriteString("<function " + fn.Name + ">")
				// case *BuiltinFunctionDataImpl:
				// 	trace.WriteString("<builtin " + fn.Name + ">")
			}
		}

		node := scope.CurrentNode()
		var token *tokens.Token
		if node != nil {
			token = node.GetToken()
		}
		if token != nil {
			trace.WriteString(fmt.Sprintf(" @ line %d, column %d", token.Line, token.Column))
		}

		if i > 0 {
			trace.WriteString("\n")
		}
	}

	return trace.String()
}

func (t *ErrorInfo) IncompatibleTypeOperation(s *Scope, op string, t1 *Instance, t2 *Instance) *Instance {
	return Error.Create(s, "invalid operation with incompatible types: '%s' %s '%s'", t1.Type.GetName(), op, t2.Type.GetName())
}

func (t *ErrorInfo) InvalidOperation(s *Scope, op string, t1 *Instance) *Instance {
	return Error.Create(s, "type '%s' does not implement operator '%s'", t1.Type.GetName(), op)
}

func (t *ErrorInfo) InvalidAction(s *Scope, action string, t1 *Instance) *Instance {
	return Error.Create(s, "type '%s' does not implement action '%s'", t1.Type.GetName(), action)
}

func (t *ErrorInfo) DuplicatedDefinition(s *Scope, name string) *Instance {
	return Error.Create(s, "variable '%s' is already defined", name)
}

func (t *ErrorInfo) ReassigningConstant(s *Scope, name string) *Instance {
	return Error.Create(s, "invalid constant assignment '%s'", name)
}

func (t *ErrorInfo) VariableNotDefined(s *Scope, name string) *Instance {
	return Error.Create(s, "trying to use an unidentified variable '%s'", name)
}

func (t *ErrorInfo) NoProperty(s *Scope, typeName string, name string) *Instance {
	return Error.Create(s, "instance of type '%s' does not have property '%s'", typeName, name)
}

// ----------------------------------------------------------------------------
// ERROR DATA TYPE
// ----------------------------------------------------------------------------
type ErrorDataType struct {
	BaseDataType
}

func (d *ErrorDataType) Instantiate(r *Runtime, s *Scope, init ast.Initializer) *Instance {
	return Error.Create(s, "application error")
}

func (d *ErrorDataType) OnNew(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*ErrorDataImpl)

	if len(args) > 1 {
		this.Properties["message"] = args[1]
	}

	return args[0]
}

func (d *ErrorDataType) OnSet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*ErrorDataImpl)
	name := AsString(args[1])

	_, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	this.Properties[name] = args[2]
	return args[2]
}

func (d *ErrorDataType) OnGet(r *Runtime, s *Scope, args ...*Instance) *Instance {
	this := args[0].Impl.(*ErrorDataImpl)
	name := AsString(args[1])

	value, has := this.Properties[name]
	if !has {
		return r.Throw(Error.NoProperty(s, d.Name, name), s)
	}

	return value
}

func (t *ErrorDataType) OnString(r *Runtime, s *Scope, args ...*Instance) *Instance {
	return t.OnRepr(r, s, args[0])
}

func (d *ErrorDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	msg := AsString(args[0].Impl.(*ErrorDataImpl).Properties["message"])
	trace := AsString(args[0].Impl.(*ErrorDataImpl).Properties["trace"])
	return String.Create("ERR! " + msg + "\n" + trace)
}

// ----------------------------------------------------------------------------
// ERROR DATA IMPL
// ----------------------------------------------------------------------------
type ErrorDataImpl struct {
	Properties map[string]*Instance
}
