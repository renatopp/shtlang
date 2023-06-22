package tokens

import (
	"encoding/json"
	"fmt"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

func CreateToken(t Type, l string, line, column int) *Token {
	return &Token{
		Type:    t,
		Literal: l,
		Line:    line,
		Column:  column,
	}
}

func (t *Token) String() string {
	return t.Literal
}

func (t *Token) Pretty() string {
	json, _ := json.Marshal(t.Literal)
	return fmt.Sprintf("<%s@%d,%d:%s>", t.Type, t.Line, t.Column, json)
}

func (t *Token) Is(tp Type) bool {
	return t.Type == tp
}
