package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type MatchCase struct {
	Token     *tokens.Token
	Condition Node
	Body      Node
}

func (p *MatchCase) GetToken() *tokens.Token {
	return p.Token
}

func (p *MatchCase) String() string {
	return fmt.Sprintf("<case>")
}

func (p *MatchCase) Children() []Node {
	return []Node{p.Condition, p.Body}
}

func (p *MatchCase) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Condition.Traverse(level+1, fn)
	p.Body.Traverse(level+1, fn)
}
