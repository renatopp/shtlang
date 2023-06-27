package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type Indexing struct {
	Token  *tokens.Token
	Target Node
	Values []Node
}

func (p *Indexing) String() string {
	return fmt.Sprintf("<indexing>")
}

func (p *Indexing) Children() []Node {
	return append(append([]Node{}, p.Target), p.Values...)
}

func (p *Indexing) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Target.Traverse(level+1, fn)
	for _, args := range p.Values {
		args.Traverse(level+1, fn)
	}
}
