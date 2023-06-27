package ast

type Call struct {
	Target    Node
	Arguments []Node
}

func (p *Call) String() string {
	return "<call>"
}

func (p *Call) Children() []Node {
	return append(append([]Node{}, p.Target), p.Arguments...)
}

func (p *Call) Traverse(level int, fn tfunc) {
	fn(level, p)
	p.Target.Traverse(level+1, fn)
	for _, args := range p.Arguments {
		args.Traverse(level+1, fn)
	}
}
