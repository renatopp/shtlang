package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type If struct {
	Token     *tokens.Token
	Condition Node
	TrueBody  Node
	FalseBody Node
}

func (p *If) GetToken() *tokens.Token {
	return p.Token
}

func (p *If) String() string {
	return fmt.Sprintf("<if>")
}

func (p *If) Children() []Node {
	values := []Node{p.Condition, p.TrueBody}
	if p.FalseBody != nil {
		values = append(values, p.FalseBody)
	}
	return values
}

func (p *If) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Condition.Traverse(level+1, fn)
	p.TrueBody.Traverse(level+1, fn)
	if p.FalseBody != nil {
		p.FalseBody.Traverse(level+1, fn)
	}
}
