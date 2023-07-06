package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Property struct {
	Token *tokens.Token
	Name  string
	Value Node
}

func (p *Property) GetToken() *tokens.Token {
	return p.Token
}

func (p *Property) String() string {
	return fmt.Sprintf("<property:%s>", p.Name)
}

func (p *Property) Children() []Node {
	return []Node{p.Value}
}

func (p *Property) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Value.Traverse(level+1, fn)
}
