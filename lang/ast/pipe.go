package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Pipe struct {
	Token  *tokens.Token
	Left   Node
	PipeFn Node
	ArgFn  Node
}

func (p *Pipe) GetToken() *tokens.Token {
	return p.Token
}

func (p *Pipe) String() string {
	return fmt.Sprintf("<pipe>")
}

func (p *Pipe) Children() []Node {
	if p.ArgFn != nil {
		return []Node{p.Left, p.PipeFn, p.ArgFn}
	} else {
		return []Node{p.Left, p.PipeFn}
	}
}

func (p *Pipe) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Left.Traverse(level+1, fn)
	p.PipeFn.Traverse(level+1, fn)
	if p.ArgFn != nil {
		p.ArgFn.Traverse(level+1, fn)
	}
}
