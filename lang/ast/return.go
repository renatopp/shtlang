package ast

import "sht/lang/tokens"

type Return struct {
	Token      *tokens.Token
	Expression Node
}

func (p *Return) GetToken() *tokens.Token {
	return p.Token
}

func (p *Return) String() string {
	return "<return>"
}

func (p *Return) Children() []Node {
	return []Node{}
}

func (p *Return) Traverse(level int, fn tfunc) {
	fn(level, p)

	if p.Expression != nil {
		p.Expression.Traverse(level+1, fn)
	}
}
