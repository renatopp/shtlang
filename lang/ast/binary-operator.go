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

func (p *BinaryOperator) String() string {
	return fmt.Sprintf("<operator:%s>", p.Operator)
}

func (p *BinaryOperator) Children() []Node {
	return []Node{p.Left, p.Right}
}

func (p *BinaryOperator) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Left.Traverse(level+1, fn)
	p.Right.Traverse(level+1, fn)
}
