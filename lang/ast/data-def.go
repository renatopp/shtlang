package ast

import "sht/lang/tokens"

type DataDef struct {
	Token         *tokens.Token
	Name          string
	Properties    []Node
	Functions     []Node
	MetaFunctions []Node
}

func (p *DataDef) GetToken() *tokens.Token {
	return p.Token
}

func (p *DataDef) String() string {
	return "<datadef:" + p.Name + ">"
}

func (p *DataDef) Children() []Node {
	return append(append(append([]Node{}, p.Properties...), p.Functions...), p.MetaFunctions...)
}

func (p *DataDef) Traverse(level int, fn tfunc) {
	fn(level, p)
	for _, prop := range p.Properties {
		prop.Traverse(level+1, fn)
	}
	for _, f := range p.Functions {
		f.Traverse(level+1, fn)
	}
	for _, f := range p.MetaFunctions {
		f.Traverse(level+1, fn)
	}
}
