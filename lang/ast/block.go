package ast

type Block struct {
	Statements []Node
}

func (p *Block) String() string {
	return "<block>"
}

func (p *Block) Children() []Node {
	return p.Statements
}

func (p *Block) Traverse(level int, fn tfunc) {
	fn(0, p)

	for _, s := range p.Children() {
		s.Traverse(level, fn)
	}
}
