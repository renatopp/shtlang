package ast

import "sht/lang/tokens"

type Parameter struct {
	Token   *tokens.Token
	Name    string
	Spread  bool
	Default Node
}

func (p *Parameter) String() string {
	prefix := ""
	suffix := ""
	if p.Spread {
		prefix = "..."
	}

	if p.Default != nil {
		suffix = " = " + p.Default.String()
	}

	return "<param:" + prefix + p.Name + suffix + ">"
}

func (p *Parameter) Children() []Node {
	r := []Node{}
	if p.Default != nil {
		r = append(r, p.Default)
	}
	return r
}

func (p *Parameter) Traverse(level int, fn tfunc) {
	fn(level, p)
	if p.Default != nil {
		p.Default.Traverse(level+1, fn)
	}
}
