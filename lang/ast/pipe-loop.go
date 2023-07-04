package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type PipeLoop struct {
	Token      *tokens.Token
	Iterator   Node
	Assignment Node
	Body       Node
}

func (p *PipeLoop) GetToken() *tokens.Token {
	return p.Token
}

func (p *PipeLoop) String() string {
	return fmt.Sprintf("<pipe loop>")
}

func (p *PipeLoop) Children() []Node {
	return []Node{p.Iterator, p.Assignment, p.Body}
}

func (p *PipeLoop) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Iterator.Traverse(level+1, fn)
	p.Assignment.Traverse(level+1, fn)
	p.Body.Traverse(level+1, fn)
}
