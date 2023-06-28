package lang

import (
	"sht/lang/tokens"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _createToken(t tokens.Type, l string) *tokens.Token {
	return &tokens.Token{
		Type:    t,
		Literal: l,
	}
}

func TestTokenizeSymbols(t *testing.T) {
	input := `; , : ! ? . @ { } ( ) [ ] => ...`

	expected := []*tokens.Token{
		_createToken(tokens.Semicolon, ";"),
		_createToken(tokens.Comma, ","),
		_createToken(tokens.Colon, ":"),
		_createToken(tokens.Bang, "!"),
		_createToken(tokens.Question, "?"),
		_createToken(tokens.Dot, "."),
		_createToken(tokens.At, "@"),
		_createToken(tokens.Lbrace, "{"),
		_createToken(tokens.Rbrace, "}"),
		_createToken(tokens.Lparen, "("),
		_createToken(tokens.Rparen, ")"),
		_createToken(tokens.Lbracket, "["),
		_createToken(tokens.Rbracket, "]"),
		_createToken(tokens.Arrow, "=>"),
		_createToken(tokens.Spread, "..."),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}

func TestTokenizeOperators(t *testing.T) {
	input := `+ - * / // ?? % ** ++ -- < <= > >= == != ..`

	expected := []*tokens.Token{
		_createToken(tokens.Operator, "+"),
		_createToken(tokens.Operator, "-"),
		_createToken(tokens.Operator, "*"),
		_createToken(tokens.Operator, "/"),
		_createToken(tokens.Operator, "//"),
		_createToken(tokens.Operator, "??"),
		_createToken(tokens.Operator, "%"),
		_createToken(tokens.Operator, "**"),
		_createToken(tokens.Operator, "++"),
		_createToken(tokens.Operator, "--"),
		_createToken(tokens.Operator, "<"),
		_createToken(tokens.Operator, "<="),
		_createToken(tokens.Operator, ">"),
		_createToken(tokens.Operator, ">="),
		_createToken(tokens.Operator, "=="),
		_createToken(tokens.Operator, "!="),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}

func TestTokenizeAssignment(t *testing.T) {
	input := `= += -= *= /= //= ..=`

	expected := []*tokens.Token{
		_createToken(tokens.Assignment, "="),
		_createToken(tokens.Assignment, "+="),
		_createToken(tokens.Assignment, "-="),
		_createToken(tokens.Assignment, "*="),
		_createToken(tokens.Assignment, "/="),
		_createToken(tokens.Assignment, "//="),
		_createToken(tokens.Assignment, "..="),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}

func TestTokenizeSpaces(t *testing.T) {
	input := `. \
	 .           .



	.
	`

	expected := []*tokens.Token{
		_createToken(tokens.Dot, "."),
		_createToken(tokens.Dot, "."),
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

	expected := []*tokens.Token{
		_createToken(tokens.Identifier, "valid"),
		_createToken(tokens.Identifier, "st$ing"),
		_createToken(tokens.Identifier, "$Here"),
		_createToken(tokens.Keyword, "on"),
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

	expected := []*tokens.Token{
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

	expected := []*tokens.Token{
		_createToken(tokens.Number, "123"),
		_createToken(tokens.Number, "1e321"),
		_createToken(tokens.Number, ".12"),
		_createToken(tokens.Operator, "-"),
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

	expected := []*tokens.Token{
		_createToken(tokens.Bang, "!"),
		_createToken(tokens.Eof, ""),
	}

	result, err := Tokenize([]byte(input))

	assert.Equal(t, err, nil)
	for i, token := range expected {
		assert.Equal(t, token.Type, result[i].Type)
		assert.Equal(t, token.Literal, result[i].Literal)
	}
}

func TestInvalidToken(t *testing.T) {
	input := `�������`
	_, err := Tokenize([]byte(input))
	assert.NotEqual(t, err, nil)
}

func TestInvalidCharacter(t *testing.T) {
	input := `☂`
	_, err := Tokenize([]byte(input))
	assert.NotEqual(t, err, nil)
}
