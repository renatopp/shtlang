package lang

import (
	"fmt"
	"sht/lang/tokens"
	"unicode/utf8"
)

func CreateLexer(input []byte) *Lexer {
	return &Lexer{
		input:      input,
		tokenQueue: []*Token{},
		partQueue:  []*lexerPart{},
		errors:     []string{},
		line:       1,
		column:     1,
		cursor:     0,
	}
}

func Tokenize(input []byte) ([]*Token, error) {
	lexer := CreateLexer(input)

	r := []*Token{}
	for {
		token := lexer.Next()
		r = append(r, token)

		if token.Is(tokens.Eof) {
			break
		}
	}

	return r, nil
}

type lexerPart struct {
	Char   rune
	Size   int
	Line   int
	Column int
}

func (p *lexerPart) Is(r rune) bool {
	return p.Char == r
}

type Lexer struct {
	input      []byte
	tokenQueue []*Token
	partQueue  []*lexerPart
	errors     []string
	line       int
	column     int
	cursor     int
	eof        *Token
}

func (l *Lexer) Next() *Token {
	token := l.Peek()
	l.tokenQueue = l.tokenQueue[1:]
	return token
}

func (l *Lexer) Peek() *Token {
	return l.PeekN(0)
}

func (l *Lexer) PeekN(i int) *Token {
	if len(l.tokenQueue) <= i {
		l.tokenQueue = append(l.tokenQueue, l.parseNextToken())
	}

	return l.tokenQueue[i]
}

func (l *Lexer) nextPart() *lexerPart {
	part := l.peekPart()
	l.partQueue = l.partQueue[1:]
	return part
}

func (l *Lexer) peekPart() *lexerPart {
	return l.peekPartN(0)
}

func (l *Lexer) peekPartN(n int) *lexerPart {
	if len(l.partQueue) <= n {
		l.partQueue = append(l.partQueue, l.parseNextPart())
	}

	return l.partQueue[n]
}

func (l *Lexer) isWhitespace(r rune) bool {
	return r == '\n' || r == '\r' || r == '\t' || r == ' '
}

func (l *Lexer) isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '$'
}

func (l *Lexer) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l *Lexer) isEOF(r rune) bool {
	return r == 0
}

func (l *Lexer) registerError(e string, part *lexerPart) {
	l.errors = append(l.errors, fmt.Sprintf("%s at %d:%d", e, part.Line, part.Column))
}

// Parse the next character given the cursor position
func (l *Lexer) parseNextPart() *lexerPart {
	for l.cursor < len(l.input) {
		r, size := utf8.DecodeRune(l.input[l.cursor:])
		part := &lexerPart{
			Char:   r,
			Size:   size,
			Line:   l.line,
			Column: l.column,
		}

		l.column++

		if r == utf8.RuneError {
			l.registerError("invalid UTF-8 character", part)
			l.cursor++
			continue
		}

		if r == '\n' {
			l.line++
			l.column = 1
		}

		l.cursor += size
		return part
	}

	return &lexerPart{
		Char:   0,
		Size:   0,
		Line:   l.line,
		Column: l.column,
	}
}

// Parse the next token given the parts queue
func (l *Lexer) parseNextToken() *Token {
	if l.eof != nil {
		return l.eof
	}

	part := l.peekPart()
	if part.Char == 0 {
		l.eof = CreateToken(tokens.Eof, "", part.Line, part.Column)
		return l.eof
	}

	var token *Token
	switch {
	case part.Is('!'):
		token = CreateToken(tokens.Bang, "!", part.Line, part.Column)
		l.nextPart()
	}
	// CREATE TOKENS HERE

	return token
}
