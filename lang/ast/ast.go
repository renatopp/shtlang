package ast

import "sht/lang/tokens"

type tfunc func(int, Node)

type Node interface {
	GetToken() *tokens.Token
	String() string
	Children() []Node
	Traverse(int, tfunc)
}
