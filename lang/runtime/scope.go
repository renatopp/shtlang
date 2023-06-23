package runtime

type Scope struct {
	Parent *Scope
	Values map[string]interface{}
}

func CreateScope(parent *Scope) *Scope {
	s := &Scope{}
	s.Parent = parent
	s.Values = make(map[string]interface{})
	return s
}

func (s *Scope) Get(name string) interface{} {
	if val, ok := s.Values[name]; ok {
		return val
	}

	if s.Parent != nil {
		return s.Parent.Get(name)
	}

	return nil
}

func (s *Scope) Set(name string, value interface{}) {
	s.Values[name] = value
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
