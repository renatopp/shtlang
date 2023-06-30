package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Access struct {
	Token *tokens.Token
	Left  Node
	Right Node
}

func (p *Access) GetToken() *tokens.Token {
	return p.Token
}

func (p *Access) String() string {
	return fmt.Sprintf("<access>")
}

func (p *Access) Children() []Node {
	return []Node{p.Left, p.Right}
}

func (p *Access) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Left.Traverse(level+1, fn)
	p.Right.Traverse(level+1, fn)
}
