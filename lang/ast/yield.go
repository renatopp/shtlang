package ast

import "sht/lang/tokens"

type Yield struct {
	Token      *tokens.Token
	Expression Node
}

func (p *Yield) GetToken() *tokens.Token {
	return p.Token
}

func (p *Yield) String() string {
	return "<yield>"
}

func (p *Yield) Children() []Node {
	return []Node{}
}

func (p *Yield) Traverse(level int, fn tfunc) {
	fn(level, p)

	if p.Expression != nil {
		p.Expression.Traverse(level+1, fn)
	}
}
