package ast

type tfunc func(int, Node)

type Node interface {
	String() string
	Children() []Node
	Traverse(int, tfunc)
}
