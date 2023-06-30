package ast

import "sht/lang/tokens"

type FunctionDef struct {
	Token     *tokens.Token
	Scoped    bool
	Generator bool
	Name      string
	Params    []Node
	Body      Node
}

func (p *FunctionDef) GetToken() *tokens.Token {
	return p.Token
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
