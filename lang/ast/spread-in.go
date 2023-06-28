package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type SpreadIn struct {
	Token  *tokens.Token
	Target Node
}

func (p *SpreadIn) GetToken() *tokens.Token {
	return p.Token
}

func (p *SpreadIn) String() string {
	return fmt.Sprintf("<spread in>")
}

func (p *SpreadIn) Children() []Node {
	return append([]Node{}, p.Target)
}

func (p *SpreadIn) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Target.Traverse(level+1, fn)
}
