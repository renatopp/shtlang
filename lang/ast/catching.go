package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Catching struct {
	Token      *tokens.Token
	Expression Node
}

func (p *Catching) GetToken() *tokens.Token {
	return p.Token
}

func (p *Catching) String() string {
	return fmt.Sprintf("<catching>")
}

func (p *Catching) Children() []Node {
	return []Node{p.Expression}
}

func (p *Catching) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Expression.Traverse(level+1, fn)
}
