package ast

import (
	"sht/lang/tokens"
	"strings"
)

type DataDef struct {
	Token         *tokens.Token
	Name          string
	Likes         []string
	Properties    []Node
	Functions     []Node
	MetaFunctions []Node
}

func (p *DataDef) GetToken() *tokens.Token {
	return p.Token
}

func (p *DataDef) String() string {
	name := "<datadef:" + p.Name + ">"
	if len(p.Likes) > 0 {
		name += " like " + strings.Join(p.Likes, ", ")
	}

	return name
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
