package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Match struct {
	Token      *tokens.Token
	Expression Node
	Cases      []Node
}

func (p *Match) GetToken() *tokens.Token {
	return p.Token
}

func (p *Match) String() string {
	return fmt.Sprintf("<match>")
}

func (p *Match) Children() []Node {
	return append([]Node{p.Expression}, p.Cases...)
}

func (p *Match) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Expression.Traverse(level+1, fn)
	for _, c := range p.Cases {
		c.Traverse(level+1, fn)
	}
}
