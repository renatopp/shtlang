package runtime

import (
	"fmt"
	"sht/lang/ast"
	"sht/lang/tokens"
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

func (t *ErrorInfo) Create(message string, a ...any) *Instance {
	msg := fmt.Sprintf(message, a...)
	return &Instance{
		Type: t.Type,
		Impl: ErrorDataImpl{
			Values: map[string]*Instance{
				"message": String.Create(msg),
			},
		},
	}
}

func (t *ErrorInfo) StackTrace(s *Scope) {
	stack := s.Stack()

	for _, scope := range stack {
		fn, _ := scope.GetInScope(SCOPE_FN_KEY)
		if fn == nil {
			fmt.Print("global")
		} else {
			fn := fn.Value.Impl.(*CustomFunctionDataImpl)
			fmt.Print(fn.Name)
		}

		node := scope.CurrentNode()
		var token *tokens.Token
		if node != nil {
			token = node.GetToken()
		}
		if token != nil {
			fmt.Printf(" at %d, %d\n", token.Line, token.Column)
		} else {
			fmt.Println()
		}
	}
}

func (t *ErrorInfo) IncompatibleTypeOperation(op string, t1 *Instance, t2 *Instance) *Instance {
	return Error.Create("invalid operation with incompatible types: '%s' %s '%s'", t1.Type.GetName(), op, t2.Type.GetName())
}

func (t *ErrorInfo) InvalidOperation(op string, t1 *Instance) *Instance {
	return Error.Create("type '%s' does not implement operator '%s'", t1.Type.GetName(), op)
}

func (t *ErrorInfo) InvalidAction(action string, t1 *Instance) *Instance {
	return Error.Create("type '%s' does not implement action '%s'", t1.Type.GetName(), action)
}

func (t *ErrorInfo) DuplicatedDefinition(name string) *Instance {
	return Error.Create("variable '%s' is already defined", name)
}

func (t *ErrorInfo) ReassigningConstant(name string) *Instance {
	return Error.Create("invalid constant assignment '%s'", name)
}

func (t *ErrorInfo) VariableNotDefined(name string) *Instance {
	return Error.Create("trying to use an unidentified variable '%s'", name)
}

// ----------------------------------------------------------------------------
// ERROR DATA TYPE
// ----------------------------------------------------------------------------
type ErrorDataType struct {
	BaseDataType
}

func (d *ErrorDataType) OnRepr(r *Runtime, s *Scope, args ...*Instance) *Instance {
	msg := AsString(args[0].Impl.(ErrorDataImpl).Values["message"])
	return String.Create("ERR! " + msg)
}

// ----------------------------------------------------------------------------
// ERROR DATA IMPL
// ----------------------------------------------------------------------------
type ErrorDataImpl struct {
	Values map[string]*Instance
}
