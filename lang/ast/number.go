package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Number struct {
	Token *tokens.Token
	Value float64
}

func (p *Number) GetToken() *tokens.Token {
	return p.Token
}

func (p *Number) String() string {
	return fmt.Sprintf("<number:%f>", p.Value)
}

func (p *Number) Children() []Node {
	return []Node{}
}

func (p *Number) Traverse(level int, fn tfunc) {
	fn(level, p)
}
