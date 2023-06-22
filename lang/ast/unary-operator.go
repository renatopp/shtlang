package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type UnaryOperator struct {
	Token    *tokens.Token
	Operator string
	Right    Node
}

func (p *UnaryOperator) String() string {
	return fmt.Sprintf("<unary:%s>", p.Operator)
}

func (p *UnaryOperator) Children() []Node {
	return []Node{p.Right}
}

func (p *UnaryOperator) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Right.Traverse(level+1, fn)
}
