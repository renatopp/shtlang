package runtime

type Stack struct {
	Global *Scope
	Scopes []*Scope
}

func NewStack(global *Scope) *Stack {
	return &Stack{
		Global: global,
		Scopes: []*Scope{},
	}
}

func (s *Stack) Current() *Scope {
	if len(s.Scopes) == 0 {
		return s.Global
	}

	return s.Scopes[len(s.Scopes)-1]
}

func (s *Stack) Push(scope *Scope) *Scope {
	if scope == nil {
		scope = s.Current()
	}

	s.Scopes = append(s.Scopes, CreateScope(scope))
	return s.Current()
}

func (s *Stack) Pop() *Scope {
	if len(s.Scopes) > 0 {
		s.Scopes = s.Scopes[:len(s.Scopes)-1]
	}

	return s.Current()
}

func (e *Stack) Get(name string) (*Instance, bool) {
	return e.Current().Get(name)
}

func (e *Stack) Set(name string, val *Instance) *Instance {
	return e.Current().Set(name, val)
}

func (e *Stack) Delete(name string) {
	e.Current().Delete(name)
}

func (e *Stack) ForEach(fn func(string, *Instance)) {
	e.Global.ForEach(fn)
}
