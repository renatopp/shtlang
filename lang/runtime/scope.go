package runtime

import (
	"fmt"
	"sht/lang/ast"
)

type Scope struct {
	Parent *Scope
	Caller *Scope
	Values map[string]*Reference

	InAssignment bool
	InArgument   bool

	nodeStack []ast.Node
}

func CreateScope(parent *Scope, caller *Scope) *Scope {
	s := &Scope{}
	s.Parent = parent
	s.Caller = caller
	s.Values = make(map[string]*Reference)
	s.nodeStack = make([]ast.Node, 0)
	return s
}

func (s *Scope) Get(name string) (*Reference, bool) {
	if val, ok := s.Values[name]; ok {
		return val, true
	}

	if s.Parent != nil {
		return s.Parent.Get(name)
	}

	return nil, false
}

func (s *Scope) GetInScope(name string) (*Reference, bool) {
	if val, ok := s.Values[name]; ok {
		return val, true
	}

	return nil, false
}

func (s *Scope) Set(name string, value *Reference) *Reference {
	s.Values[name] = value
	return value
}

func (s *Scope) Has(name string) bool {
	if _, ok := s.Values[name]; ok {
		return true
	}

	if s.Parent != nil {
		return s.Parent.Has(name)
	}

	return false
}

func (s *Scope) HasInScope(name string) bool {
	if _, ok := s.Values[name]; ok {
		return true
	}

	return false
}

func (s *Scope) Delete(name string) {
	delete(s.Values, name)
}

func (s *Scope) ForEach(fn func(string, *Reference)) {
	for k, v := range s.Values {
		fn(k, v)
	}
}

func (s *Scope) print(i int, stack []*Scope) {
	if i >= len(stack) {
		return
	}

	scope := stack[i]
	name, _ := scope.Get(SCOPE_NAME_KEY)
	prefix := fmt.Sprintf("%*s", (i+1)*2, "")
	prefix2 := fmt.Sprintf("%*s", (i+2)*2, "")
	fmt.Printf(prefix+"scope %s {\n", AsString(name.Value))

	scope.ForEach(func(s string, r *Reference) {
		fmt.Println(prefix2 + s + ": " + r.Value.Repr())
	})

	s.print(i+1, stack)

	fmt.Println(prefix + "}")
}

func (s *Scope) PrintSelf() {
	s.ForEach(func(s string, r *Reference) {
		fmt.Println(s, ":", r.Value.Repr())
	})
}

func (s *Scope) Print() {
	stack := s.ScopeStack()
	s.print(0, stack)
}

func (s *Scope) PushNode(node ast.Node) {
	s.nodeStack = append(s.nodeStack, node)
}

func (s *Scope) PopNode() {
	s.nodeStack = s.nodeStack[:len(s.nodeStack)-1]
}

func (s *Scope) CurrentNode() ast.Node {
	if len(s.nodeStack) <= 0 {
		return nil
	}
	return s.nodeStack[len(s.nodeStack)-1]
}

func (s *Scope) ScopeStack() []*Scope {
	stack := make([]*Scope, 0)
	scope := s
	for scope != nil {
		stack = append([]*Scope{scope}, stack...)
		scope = scope.Parent
	}

	return stack
}

func (s *Scope) CallStack() []*Scope {
	stack := make([]*Scope, 0)
	scope := s
	for scope != nil {
		stack = append([]*Scope{scope}, stack...)
		scope = scope.Caller
	}

	return stack
}
