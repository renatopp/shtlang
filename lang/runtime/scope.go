package runtime

type Scope struct {
	Parent *Scope
	Values map[string]*Reference
}

func CreateScope(parent *Scope) *Scope {
	s := &Scope{}
	s.Parent = parent
	s.Values = make(map[string]*Reference)
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
