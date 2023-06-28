package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Unwrapping struct {
	Token  *tokens.Token
	Target Node
}

func (p *Unwrapping) GetToken() *tokens.Token {
	return p.Token
}

func (p *Unwrapping) String() string {
	return fmt.Sprintf("<unwrapping>")
}

func (p *Unwrapping) Children() []Node {
	return []Node{p.Target}
}

func (p *Unwrapping) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Target.Traverse(level+1, fn)
}
