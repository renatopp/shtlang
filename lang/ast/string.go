package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type String struct {
	Token *tokens.Token
	Value string
}

func (p *String) String() string {
	return fmt.Sprintf("<string:%s>", p.Value)
}

func (p *String) Children() []Node {
	return []Node{}
}

func (p *String) Traverse(level int, fn tfunc) {
	fn(level, p)
}
