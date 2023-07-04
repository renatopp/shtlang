package ast

import "sht/lang/tokens"

type Block struct {
	Unscoped   bool
	Statements []Node
}

func (p *Block) GetToken() *tokens.Token {
	return nil
}

func (p *Block) String() string {
	return "<block>"
}

func (p *Block) Children() []Node {
	return p.Statements
}

func (p *Block) Traverse(level int, fn tfunc) {
	fn(level, p)

	for _, s := range p.Children() {
		s.Traverse(level+1, fn)
	}
}
