package ast

import "sht/lang/tokens"

type Boolean struct {
	Token *tokens.Token
	Value bool
}

func (p *Boolean) GetToken() *tokens.Token {
	return p.Token
}

func (p *Boolean) String() string {
	return "<boolean:" + p.Token.Literal + ">"
}

func (p *Boolean) Children() []Node {
	return []Node{}
}

func (p *Boolean) Traverse(level int, fn tfunc) {
	fn(level, p)
}
