package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type PostfixOperator struct {
	Token    *tokens.Token
	Operator string
	Left     Node
}

func (p *PostfixOperator) String() string {
	return fmt.Sprintf("<postfix:%s>", p.Operator)
}

func (p *PostfixOperator) Children() []Node {
	return []Node{p.Left}
}

func (p *PostfixOperator) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Left.Traverse(level+1, fn)
}
