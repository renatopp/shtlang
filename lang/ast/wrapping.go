package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Wrapping struct {
	Token      *tokens.Token
	Expression Node
}

func (p *Wrapping) GetToken() *tokens.Token {
	return p.Token
}

func (p *Wrapping) String() string {
	return fmt.Sprintf("<wrapping>")
}

func (p *Wrapping) Children() []Node {
	return []Node{p.Expression}
}

func (p *Wrapping) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Expression.Traverse(level+1, fn)
}
