package ast

import "sht/lang/tokens"

type FunctionDef struct {
	Token  *tokens.Token
	Name   string
	Params []Node
	Body   Node
	Maybe  bool
}

func (p *FunctionDef) String() string {
	return "<funcdef:" + p.Name + ">"
}

func (p *FunctionDef) Children() []Node {
	return append(append([]Node{}, p.Params...), p.Body)
}

func (p *FunctionDef) Traverse(level int, fn tfunc) {
	fn(level, p)
	for _, param := range p.Params {
		param.Traverse(level+1, fn)
	}
	p.Body.Traverse(level+1, fn)
}
