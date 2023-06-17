package lang

import (
	"fmt"
	"sht/lang/tokens"
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/slices"
)

var keywords = []string{
	"true",
	"false",
	"if",
	"for",
	"while",
	"break",
	"continue",
	"return",
	"fn",
	"let",
	"const",
	"as",
	"data",
	"raise",
	"on",
}

func CreateLexer(input []byte) *Lexer {
	return &Lexer{
		input:      input,
		tokenQueue: []*Token{},
		charQueue:  []*char{},
		errors:     []string{},
		line:       1,
		column:     1,
		cursor:     0,
		builder:    &strings.Builder{},
	}
}

func Tokenize(input []byte) ([]*Token, error) {
	lexer := CreateLexer(input)

	r := []*Token{}
	for {
		token := lexer.EatToken()
		r = append(r, token)

		if token.Is(tokens.Eof) {
			break
		}
	}

	return r, lexer.GetError()
}

type char struct {
	Rune   rune
	Size   int
	Line   int
	Column int
}

func (p *char) Is(r rune) bool {
	return p.Rune == r
}

type Lexer struct {
	input      []byte
	tokenQueue []*Token
	charQueue  []*char
	errors     []string
	line       int
	column     int
	cursor     int
	eof        *Token
	builder    *strings.Builder
}

func (l *Lexer) EatToken() *Token {
	token := l.PeekToken()
	l.tokenQueue = l.tokenQueue[1:]
	return token
}

func (l *Lexer) PeekToken() *Token {
	return l.PeekTokenN(0)
}

func (l *Lexer) PeekTokenN(i int) *Token {
	if len(l.tokenQueue) <= i {
		l.tokenQueue = append(l.tokenQueue, l.parseNextToken())
	}

	return l.tokenQueue[i]
}

func (l *Lexer) EatChar() *char {
	c := l.PeekChar()
	l.charQueue = l.charQueue[1:]
	return c
}

func (l *Lexer) PeekChar() *char {
	return l.PeekCharN(0)
}

func (l *Lexer) PeekCharN(n int) *char {
	if len(l.charQueue) <= n {
		l.charQueue = append(l.charQueue, l.parseNextChar())
	}

	return l.charQueue[n]
}

func (l *Lexer) TooManyErrors() bool {
	return len(l.errors) >= 10
}

func (l *Lexer) HasError() bool {
	return len(l.errors) > 0
}

func (l *Lexer) GetError() error {
	if l.HasError() {
		return fmt.Errorf("Tokenizer errors: \n- %s", strings.Join(l.errors, "\n- "))
	}

	return nil
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

func (l *Lexer) isKeyword(lit string) bool {
	return slices.Contains(keywords, lit)
}

func (l *Lexer) RegisterError(e string, c *char) {
	if l.TooManyErrors() {
		return
	}

	l.errors = append(l.errors, fmt.Sprintf("%s at %d:%d", e, c.Line, c.Column))

	if l.TooManyErrors() {
		l.errors = append(l.errors, "too many errors, aborting")
	}
}

// Parse the next character given the cursor position
func (l *Lexer) parseNextChar() *char {
	for l.cursor < len(l.input) {
		if l.TooManyErrors() {
			return &char{
				Rune:   0,
				Size:   0,
				Line:   l.line,
				Column: l.column,
			}
		}

		r, size := utf8.DecodeRune(l.input[l.cursor:])
		c := &char{
			Rune:   r,
			Size:   size,
			Line:   l.line,
			Column: l.column,
		}

		l.column++

		if r == utf8.RuneError {
			l.RegisterError("invalid UTF-8 character", c)
			l.cursor++
			continue
		}

		if r == '\n' {
			l.line++
			l.column = 1
		}

		l.cursor += size
		return c
	}

	return &char{
		Rune:   0,
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

	var token *Token

	for {
		if l.TooManyErrors() {
			l.eof = CreateToken(tokens.Eof, "", l.line, l.column)
			return l.eof
		}

		c := l.PeekChar()
		if c.Rune == 0 {
			l.eof = CreateToken(tokens.Eof, "", c.Line, c.Column)
			return l.eof
		}

		switch {
		case c.Is(';'):
			token = CreateToken(tokens.Semicolon, ";", c.Line, c.Column)
			l.EatChar()
		case c.Is(','):
			token = CreateToken(tokens.Comma, ",", c.Line, c.Column)
			l.EatChar()
		case c.Is(':'):
			token = CreateToken(tokens.Colon, ":", c.Line, c.Column)
			l.EatChar()
		case c.Is('!'):
			token = CreateToken(tokens.Bang, "!", c.Line, c.Column)
			l.EatChar()
		case c.Is('?'):
			token = CreateToken(tokens.Question, "?", c.Line, c.Column)
			l.EatChar()
		case c.Is('.') && !l.isDigit(l.PeekCharN(1).Rune):
			token = CreateToken(tokens.Dot, ".", c.Line, c.Column)
			l.EatChar()
		case c.Is('\\'):
			token = CreateToken(tokens.Backslash, "\\", c.Line, c.Column)
			l.EatChar()
		case c.Is('@'):
			token = CreateToken(tokens.At, "@", c.Line, c.Column)
			l.EatChar()
		case c.Is('%'):
			token = CreateToken(tokens.Percent, "%", c.Line, c.Column)
			l.EatChar()
		case c.Is('^'):
			token = CreateToken(tokens.Caret, "^", c.Line, c.Column)
			l.EatChar()
		case c.Is('&'):
			token = CreateToken(tokens.Ampersand, "&", c.Line, c.Column)
			l.EatChar()
		case c.Is('|'):
			token = CreateToken(tokens.Pipe, "|", c.Line, c.Column)
			l.EatChar()
		case c.Is('+'):
			token = CreateToken(tokens.Plus, "+", c.Line, c.Column)
			l.EatChar()
		case c.Is('-'):
			token = CreateToken(tokens.Minus, "-", c.Line, c.Column)
			l.EatChar()
		case c.Is('*'):
			token = CreateToken(tokens.Asterisk, "*", c.Line, c.Column)
			l.EatChar()
		case c.Is('/'):
			token = CreateToken(tokens.Slash, "/", c.Line, c.Column)
			l.EatChar()
		case c.Is('>'):
			token = CreateToken(tokens.Greater, ">", c.Line, c.Column)
			l.EatChar()
		case c.Is('<'):
			token = CreateToken(tokens.Less, "<", c.Line, c.Column)
			l.EatChar()
		case c.Is('='):
			token = CreateToken(tokens.Equal, "=", c.Line, c.Column)
			l.EatChar()
		case c.Is('~'):
			token = CreateToken(tokens.Tilde, "~", c.Line, c.Column)
			l.EatChar()
		case c.Is('{'):
			token = CreateToken(tokens.Lbrace, "{", c.Line, c.Column)
			l.EatChar()
		case c.Is('}'):
			token = CreateToken(tokens.Rbrace, "}", c.Line, c.Column)
			l.EatChar()
		case c.Is('('):
			token = CreateToken(tokens.Lparen, "(", c.Line, c.Column)
			l.EatChar()
		case c.Is(')'):
			token = CreateToken(tokens.Rparen, ")", c.Line, c.Column)
			l.EatChar()
		case c.Is('['):
			token = CreateToken(tokens.Lbracket, "[", c.Line, c.Column)
			l.EatChar()
		case c.Is(']'):
			token = CreateToken(tokens.Rbracket, "]", c.Line, c.Column)
			l.EatChar()

		case c.Is('#'):
			l.parseComment()
			continue

		case l.isWhitespace(c.Rune):
			token = l.parseWhitespaces()

		case l.isLetter(c.Rune):
			token = l.parseIdentifier()

		case l.isDigit(c.Rune) || c.Is('.') && l.isDigit(l.PeekCharN(1).Rune):
			token = l.parseNumber()

		case c.Is('\''):
			token = l.parseString()

		default:
			l.RegisterError(fmt.Sprintf("invalid character '%c'", c.Rune), c)
			l.EatChar()
			continue
		}

		break
	}

	return token
}

func (l *Lexer) parseComment() {
	l.EatChar()

	for {
		c := l.PeekChar()

		if c.Is('\n') || l.isEOF(c.Rune) {
			break
		}

		l.EatChar()
	}
}

func (l *Lexer) parseWhitespaces() *Token {
	nl := false

	first := l.PeekChar()
	for {
		c := l.PeekChar()

		if !l.isWhitespace(c.Rune) {
			break
		}

		if c.Is('\n') {
			nl = true
		}

		l.EatChar()
	}

	if nl {
		return CreateToken(tokens.Newline, "\n", first.Line, first.Column)
	}
	return CreateToken(tokens.Space, " ", first.Line, first.Column)
}

func (l *Lexer) parseIdentifier() *Token {
	l.builder.Reset()

	first := l.EatChar()
	l.builder.WriteRune(first.Rune)
	for {
		c := l.PeekChar()

		if l.isLetter(c.Rune) || l.isDigit(c.Rune) {
			l.builder.WriteRune(c.Rune)
		} else {
			break
		}

		l.EatChar()
	}

	literal := l.builder.String()
	if l.isKeyword(literal) {
		return CreateToken(tokens.Keyword, literal, first.Line, first.Column)
	} else {
		return CreateToken(tokens.Identifier, literal, first.Line, first.Column)
	}
}

func (l *Lexer) parseNumber() *Token {
	l.builder.Reset()
	dot := false
	exp := false

	first := l.PeekChar()
	for {
		c := l.PeekChar()

		if c.Is('.') {
			if dot {
				l.RegisterError("unexpected '.'", c)
				continue
			}

			dot = true
			l.builder.WriteRune(c.Rune)

		} else if c.Is('e') || c.Is('E') {
			if exp {
				l.RegisterError("unexpected 'e'", c)
				l.EatChar()
				continue
			}

			exp = true
			l.builder.WriteRune(c.Rune)

			if l.PeekCharN(1).Is('+') || l.PeekCharN(1).Is('-') {
				l.EatChar()
				l.builder.WriteRune(l.PeekChar().Rune)
			}

		} else if l.isDigit(c.Rune) {
			l.builder.WriteRune(c.Rune)

		} else {
			break
		}

		l.EatChar()
	}

	literal := l.builder.String()

	return CreateToken(tokens.Number, literal, first.Line, first.Column)
}

func (l *Lexer) parseString() *Token {
	l.builder.Reset()
	esc := false

	first := l.EatChar()
	for {
		c := l.PeekChar()

		if !esc && c.Is('\\') {
			esc = !esc
			l.EatChar()
			continue

		} else if l.isEOF(c.Rune) || !esc && c.Is('\'') {
			break

		} else {
			esc = false
		}

		l.builder.WriteRune(c.Rune)
		l.EatChar()
	}

	l.EatChar()
	return CreateToken(tokens.String, l.builder.String(), first.Line, first.Column)
}
