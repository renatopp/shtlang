package lang

import (
	"sht/lang/tokens"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _createToken(t tokens.Type, l string) *Token {
	return &Token{
		Type:    t,
		Literal: l,
	}
}

func TestTokenizeSymbols(t *testing.T) {
	input := `;,:!?.\@%^&|+-*/><=~{}()[]`

	expected := []*Token{
		_createToken(tokens.Semicolon, ";"),
		_createToken(tokens.Comma, ","),
		_createToken(tokens.Colon, ":"),
		_createToken(tokens.Bang, "!"),
		_createToken(tokens.Question, "?"),
		_createToken(tokens.Dot, "."),
		_createToken(tokens.Backslash, "\\"),
		_createToken(tokens.At, "@"),
		_createToken(tokens.Percent, "%"),
		_createToken(tokens.Caret, "^"),
		_createToken(tokens.Ampersand, "&"),
		_createToken(tokens.Pipe, "|"),
		_createToken(tokens.Plus, "+"),
		_createToken(tokens.Minus, "-"),
		_createToken(tokens.Asterisk, "*"),
		_createToken(tokens.Slash, "/"),
		_createToken(tokens.Greater, ">"),
		_createToken(tokens.Less, "<"),
		_createToken(tokens.Equal, "="),
		_createToken(tokens.Tilde, "~"),
		_createToken(tokens.Lbrace, "{"),
		_createToken(tokens.Rbrace, "}"),
		_createToken(tokens.Lparen, "("),
		_createToken(tokens.Rparen, ")"),
		_createToken(tokens.Lbracket, "["),
		_createToken(tokens.Rbracket, "]"),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}

func TestTokenizeSpaces(t *testing.T) {
	input := `.    . 
	
	
	. 
	`

	expected := []*Token{
		_createToken(tokens.Dot, "."),
		_createToken(tokens.Space, " "),
		_createToken(tokens.Dot, "."),
		_createToken(tokens.Newline, "\n"),
		_createToken(tokens.Dot, "."),
		_createToken(tokens.Newline, "\n"),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}

func TestTokenizeIdentifier(t *testing.T) {
	input := `valid st$ing $Here on raise`

	expected := []*Token{
		_createToken(tokens.Identifier, "valid"),
		_createToken(tokens.Space, " "),
		_createToken(tokens.Identifier, "st$ing"),
		_createToken(tokens.Space, " "),
		_createToken(tokens.Identifier, "$Here"),
		_createToken(tokens.Space, " "),
		_createToken(tokens.Keyword, "on"),
		_createToken(tokens.Space, " "),
		_createToken(tokens.Keyword, "raise"),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}

func TestTokenizeStrings(t *testing.T) {
	input := "'A string with <\\'123123 ☀ ☃ ☂ ☁> \n characters'"

	expected := []*Token{
		_createToken(tokens.String, "A string with <'123123 ☀ ☃ ☂ ☁> \n characters"),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}

func TestTokenizeNumbers(t *testing.T) {
	input := `123 1e321 .12 -21e-123`

	expected := []*Token{
		_createToken(tokens.Number, "123"),
		_createToken(tokens.Space, " "),
		_createToken(tokens.Number, "1e321"),
		_createToken(tokens.Space, " "),
		_createToken(tokens.Number, ".12"),
		_createToken(tokens.Space, " "),
		_createToken(tokens.Minus, "-"),
		_createToken(tokens.Number, "21e-123"),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}

func TestTokenizeComments(t *testing.T) {
	input := `! # is a comment!`

	expected := []*Token{
		_createToken(tokens.Bang, "!"),
		_createToken(tokens.Space, " "),
		_createToken(tokens.Eof, ""),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}
