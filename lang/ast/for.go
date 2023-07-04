package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type For struct {
	Token     *tokens.Token
	Condition Node
	Body      Node
}

func (p *For) GetToken() *tokens.Token {
	return p.Token
}

func (p *For) String() string {
	return fmt.Sprintf("<for>")
}

func (p *For) Children() []Node {
	values := []Node{p.Condition}
	if p.Condition != nil {
		values = append(values, p.Condition, p.Body)
	}
	return values
}

func (p *For) Traverse(level int, fn tfunc) {
	fn(level, p)
	if p.Condition != nil {
		p.Condition.Traverse(level+1, fn)
	}
	p.Body.Traverse(level+1, fn)
}
