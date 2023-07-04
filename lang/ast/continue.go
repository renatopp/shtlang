package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Continue struct {
	Token *tokens.Token
}

func (p *Continue) GetToken() *tokens.Token {
	return p.Token
}

func (p *Continue) String() string {
	return fmt.Sprintf("<continue>")
}

func (p *Continue) Children() []Node {
	return []Node{}
}

func (p *Continue) Traverse(level int, fn tfunc) {
	fn(level, p)
}
