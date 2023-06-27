package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Assignment struct {
	Token      *tokens.Token
	Literal    string
	Identifier Node
	Expression Node
	Definition bool
	Constant   bool
}

func (p *Assignment) String() string {
	add := ""
	if p.Definition {
		add += ";def"
	}
	if p.Constant {
		add += ";const"
	} else {
		add += ";var"
	}
	return fmt.Sprintf("<assignment%s:%s>", add, p.Literal)
}

func (p *Assignment) Children() []Node {
	return []Node{p.Identifier, p.Expression}
}

func (p *Assignment) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Identifier.Traverse(level+1, fn)
	p.Expression.Traverse(level+1, fn)
}
