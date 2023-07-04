package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type BinaryOperator struct {
	Token    *tokens.Token
	Operator string
	Left     Node
	Right    Node
}

func (p *BinaryOperator) GetToken() *tokens.Token {
	return p.Token
}

func (p *BinaryOperator) String() string {
	return fmt.Sprintf("<operator:%s>", p.Operator)
}

func (p *BinaryOperator) Children() []Node {
	return []Node{p.Left, p.Right}
}

func (p *BinaryOperator) Traverse(level int, fn tfunc) {
	fn(level, p)
	if p.Left != nil {
		p.Left.Traverse(level+1, fn)
	}
	if p.Right != nil {
		p.Right.Traverse(level+1, fn)
	}
}
