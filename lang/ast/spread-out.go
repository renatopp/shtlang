package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type SpreadOut struct {
	Token  *tokens.Token
	Target Node
}

func (p *SpreadOut) GetToken() *tokens.Token {
	return p.Token
}

func (p *SpreadOut) String() string {
	return fmt.Sprintf("<spread out>")
}

func (p *SpreadOut) Children() []Node {
	return append([]Node{}, p.Target)
}

func (p *SpreadOut) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Target.Traverse(level+1, fn)
}
