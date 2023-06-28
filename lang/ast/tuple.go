package ast

import "sht/lang/tokens"

type Tuple struct {
	Token  *tokens.Token
	Values []Node
}

func (p *Tuple) GetToken() *tokens.Token {
	return p.Token
}

func (p *Tuple) String() string {
	return "<tuple>"
}

func (p *Tuple) Children() []Node {
	return p.Values
}

func (p *Tuple) Traverse(level int, fn tfunc) {
	fn(level, p)

	for _, value := range p.Values {
		value.Traverse(level+1, fn)
	}
}
