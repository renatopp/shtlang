package runtime

type Scope struct {
	Parent *Scope
	Values map[string]*Instance
}

func CreateScope(parent *Scope) *Scope {
	s := &Scope{}
	s.Parent = parent
	s.Values = make(map[string]*Instance)
	return s
}

func (s *Scope) Get(name string) (*Instance, bool) {
	if val, ok := s.Values[name]; ok {
		return val, true
	}

	if s.Parent != nil {
		return s.Parent.Get(name)
	}

	return nil, false
}

func (s *Scope) Set(name string, value *Instance) *Instance {
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

func (s *Scope) Delete(name string) {
	delete(s.Values, name)
}

func (s *Scope) ForEach(fn func(string, *Instance)) {
	for k, v := range s.Values {
		fn(k, v)
	}
}