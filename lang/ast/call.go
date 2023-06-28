package ast

import "sht/lang/tokens"

type Call struct {
	Target      Node
	Arguments   []Node
	Initializer Initializer
}

func (p *Call) GetToken() *tokens.Token {
	return p.Target.GetToken()
}

func (p *Call) String() string {
	return "<call>"
}

func (p *Call) Children() []Node {
	children := []Node{p.Target}
	if p.Initializer != nil {
		children = append(children, p.Initializer)
	}
	children = append(children, p.Arguments...)

	return children
}

func (p *Call) Traverse(level int, fn tfunc) {
	fn(level, p)
	for _, args := range p.Children() {
		args.Traverse(level+1, fn)
	}
}
