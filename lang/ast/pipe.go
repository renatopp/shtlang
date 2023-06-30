package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Pipe struct {
	Token *tokens.Token
	Left  Node
	Right Node
}

func (p *Pipe) GetToken() *tokens.Token {
	return p.Token
}

func (p *Pipe) String() string {
	return fmt.Sprintf("<pipe>")
}

func (p *Pipe) Children() []Node {
	return []Node{p.Left, p.Right}
}

func (p *Pipe) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Left.Traverse(level+1, fn)
	p.Right.Traverse(level+1, fn)
}
