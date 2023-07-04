package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Break struct {
	Token *tokens.Token
}

func (p *Break) GetToken() *tokens.Token {
	return p.Token
}

func (p *Break) String() string {
	return fmt.Sprintf("<continue>")
}

func (p *Break) Children() []Node {
	return []Node{}
}

func (p *Break) Traverse(level int, fn tfunc) {
	fn(level, p)
}
