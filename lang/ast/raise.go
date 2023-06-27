package ast

import "sht/lang/tokens"

type Raise struct {
	Token      *tokens.Token
	Expression Node
}

func (p *Raise) GetToken() *tokens.Token {
	return p.Token
}

func (p *Raise) String() string {
	return "<raise>"
}

func (p *Raise) Children() []Node {
	return []Node{}
}

func (p *Raise) Traverse(level int, fn tfunc) {
	fn(level, p)

	if p.Expression != nil {
		p.Expression.Traverse(level+1, fn)
	}
}
