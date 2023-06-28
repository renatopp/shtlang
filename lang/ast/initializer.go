package ast

import (
	"fmt"
	"sht/lang/tokens"
)

type InitializerType string

var INITIALIZER_MAP InitializerType = "map"
var INITIALIZER_LIST InitializerType = "list"

type Initializer interface {
	Node

	GetType() InitializerType
}

type ListInitializer struct {
	Token  *tokens.Token
	Values []Node
}

func (p *ListInitializer) GetType() InitializerType {
	return INITIALIZER_LIST
}

func (p *ListInitializer) GetToken() *tokens.Token {
	return p.Token
}

func (p *ListInitializer) String() string {
	return fmt.Sprintf("<list initializer>")
}

func (p *ListInitializer) Children() []Node {
	return p.Values
}

func (p *ListInitializer) Traverse(level int, fn tfunc) {
	fn(level, p)
	for _, args := range p.Values {
		args.Traverse(level+1, fn)
	}
}

type MapInitializer struct {
	Token  *tokens.Token
	Values map[string]Node
}

func (p *MapInitializer) GetType() InitializerType {
	return INITIALIZER_MAP
}

func (p *MapInitializer) GetToken() *tokens.Token {
	return p.Token
}

func (p *MapInitializer) String() string {
	return fmt.Sprintf("<map initializer>")
}

func (p *MapInitializer) Children() []Node {
	values := []Node{}
	for _, value := range p.Values {
		values = append(values, value)
	}
	return values
}

func (p *MapInitializer) Traverse(level int, fn tfunc) {
	fn(level, p)
	for _, args := range p.Values {
		args.Traverse(level+1, fn)
	}
}
