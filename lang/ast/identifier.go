package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Identifier struct {
	Token *tokens.Token
	Value string
}

func (p *Identifier) String() string {
	return fmt.Sprintf("<identifier:%s>", p.Value)
}

func (p *Identifier) Children() []Node {
	return []Node{}
}

func (p *Identifier) Traverse(level int, fn tfunc) {
	fn(level, p)
}
