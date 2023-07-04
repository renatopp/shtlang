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
	"else",
	"for",
	"pipe",
	"break",
	"continue",

	"return",
	"raise",
	"yield",

	"on",
	"fn",
	"data",

	"module",
	"use",
	"async",
	"await",

	"as",
	"is",
	"in",
	"to",
}

var operators = []string{
	"and",
	"or",
	"xor",
	"nand",
	"nor",
	"nxor",
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
	tokenQueue []*tokens.Token
	charQueue  []*char
	errors     []string
	line       int
	column     int
	cursor     int
	eof        *tokens.Token
	builder    *strings.Builder
}

func CreateLexer(input []byte) *Lexer {
	return &Lexer{
		input:      input,
		tokenQueue: []*tokens.Token{},
		charQueue:  []*char{},
		errors:     []string{},
		line:       1,
		column:     1,
		cursor:     0,
		builder:    &strings.Builder{},
	}
}

func Tokenize(input []byte) ([]*tokens.Token, error) {
	lexer := CreateLexer(input)

	r := []*tokens.Token{}
	for {
		token := lexer.EatToken()
		r = append(r, token)

		if token.Is(tokens.Eof) {
			break
		}
	}

	return r, lexer.GetError()
}

func (l *Lexer) EatToken() *tokens.Token {
	token := l.PeekToken()
	l.tokenQueue = l.tokenQueue[1:]
	return token
}

func (l *Lexer) PeekToken() *tokens.Token {
	return l.PeekTokenN(0)
}

func (l *Lexer) PeekTokenN(i int) *tokens.Token {
	if len(l.tokenQueue) <= i {
		t := l.parseNextToken()

		if t.Literal == "??" {
			l.tokenQueue = append(l.tokenQueue, tokens.CreateToken(tokens.Question, "?", t.Line, t.Column))
		}

		l.tokenQueue = append(l.tokenQueue, t)
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

func (l *Lexer) RegisterError(e string, c *char) {
	if l.TooManyErrors() {
		return
	}

	l.errors = append(l.errors, fmt.Sprintf("%s at %d:%d", e, c.Line, c.Column))

	if l.TooManyErrors() {
		l.errors = append(l.errors, "too many errors, aborting")
	}
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

func (l *Lexer) isOperatorKeyword(lit string) bool {
	return slices.Contains(operators, lit)
}

func (l *Lexer) isDoubleOperator(a rune, b rune) bool {
	switch {
	case a == '+' && b == '+',
		a == '-' && b == '-',
		a == '*' && b == '*',
		a == '/' && b == '/',
		a == '<' && b == '=',
		a == '>' && b == '=',
		a == '=' && b == '=',
		a == '!' && b == '=',
		a == '.' && b == '.',
		a == '?' && b == '?':
		return true
	default:
		return false
	}
}

func (l *Lexer) isOperator(r rune) bool {
	switch {
	case r == '+',
		r == '-',
		r == '*',
		r == '/',
		r == '%',
		r == '<',
		r == '>':
		return true
	default:
		return false
	}
}

func (l *Lexer) isCompositeAssignment(a rune, b rune) bool {
	switch {
	case a == '+' && b == '=',
		a == '-' && b == '=',
		a == '*' && b == '=',
		a == '/' && b == '=',
		a == ':' && b == '=':
		return true
	default:
		return false
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

func (l *Lexer) parseWhitespaces() *tokens.Token {
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
		if l.PeekChar().Is('|') {
			return nil
		}
		return tokens.CreateToken(tokens.Newline, "\n", first.Line, first.Column)
	}

	return nil
}

func (l *Lexer) parseIdentifier() *tokens.Token {
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
	tp := tokens.Identifier
	if l.isOperatorKeyword(literal) {
		tp = tokens.Operator
	} else if l.isKeyword(literal) {
		tp = tokens.Keyword
	}
	return tokens.CreateToken(tokens.Type(tp), literal, first.Line, first.Column)
}

func (l *Lexer) parseNumber() *tokens.Token {
	l.builder.Reset()
	dot := false
	exp := false

	first := l.PeekChar()
	for {
		c := l.PeekChar()

		if c.Is('.') {
			if dot {
				l.RegisterError("unexpected '.'", c)
				l.EatChar()
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

	return tokens.CreateToken(tokens.Number, literal, first.Line, first.Column)
}

func (l *Lexer) parseString() *tokens.Token {
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
	return tokens.CreateToken(tokens.String, l.builder.String(), first.Line, first.Column)
}

func (l *Lexer) parseBacklash() {
	nl := false

	l.EatChar()
	for {
		c := l.PeekChar()

		if !l.isWhitespace(c.Rune) {
			break
		}

		if c.Is('\n') {
			if nl {
				break
			}

			nl = true
		}

		l.EatChar()
	}
}

// Parse the next token given the parts queue
func (l *Lexer) parseNextToken() *tokens.Token {
	if l.eof != nil {
		return l.eof
	}

	var token *tokens.Token

	for {
		if l.TooManyErrors() {
			l.eof = tokens.CreateToken(tokens.Eof, "", l.line, l.column)
			return l.eof
		}

		c := l.PeekChar()
		if c.Rune == 0 {
			l.eof = tokens.CreateToken(tokens.Eof, "", c.Line, c.Column)
			return l.eof
		}

		cr := l.PeekCharN(0).Rune
		nr := l.PeekCharN(1).Rune
		nnr := l.PeekCharN(2).Rune

		switch {
		case c.Is('#'):
			l.parseComment()
			continue

		case c.Is('\\'):
			l.parseBacklash()
			continue

		case cr == '.' && nr == '.' && nnr == '.':
			token = tokens.CreateToken(tokens.Spread, "...", c.Line, c.Column)
			l.EatChar()
			l.EatChar()
			l.EatChar()

		case c.Is('=') && nr == '>':
			token = tokens.CreateToken(tokens.Arrow, "=>", c.Line, c.Column)
			l.EatChar()
			l.EatChar()

		case cr == '/' && nr == '/' && nnr == '=',
			cr == '.' && nr == '.' && nnr == '=':
			token = tokens.CreateToken(tokens.Assignment, string(cr)+string(nr)+string(nnr), c.Line, c.Column)
			l.EatChar()
			l.EatChar()
			l.EatChar()

		case l.isCompositeAssignment(cr, nr):
			token = tokens.CreateToken(tokens.Assignment, string(cr)+string(nr), c.Line, c.Column)
			l.EatChar()
			l.EatChar()

		case l.isDoubleOperator(cr, nr):
			token = tokens.CreateToken(tokens.Operator, string(cr)+string(nr), c.Line, c.Column)
			l.EatChar()
			l.EatChar()

		case l.isOperator(cr):
			token = tokens.CreateToken(tokens.Operator, string(cr), c.Line, c.Column)
			l.EatChar()

		case c.Is('='):
			token = tokens.CreateToken(tokens.Assignment, "=", c.Line, c.Column)
			l.EatChar()

		case c.Is('!'):
			token = tokens.CreateToken(tokens.Bang, "!", c.Line, c.Column)
			l.EatChar()
		case c.Is(';'):
			token = tokens.CreateToken(tokens.Semicolon, ";", c.Line, c.Column)
			l.EatChar()
		case c.Is(','):
			token = tokens.CreateToken(tokens.Comma, ",", c.Line, c.Column)
			l.EatChar()
		case c.Is(':'):
			token = tokens.CreateToken(tokens.Colon, ":", c.Line, c.Column)
			l.EatChar()
		case c.Is('?'):
			token = tokens.CreateToken(tokens.Question, "?", c.Line, c.Column)
			l.EatChar()
		case c.Is('@'):
			token = tokens.CreateToken(tokens.At, "@", c.Line, c.Column)
			l.EatChar()
		case c.Is('|'):
			token = tokens.CreateToken(tokens.Pipe, "|", c.Line, c.Column)
			l.EatChar()
		case c.Is('{'):
			token = tokens.CreateToken(tokens.Lbrace, "{", c.Line, c.Column)
			l.EatChar()
		case c.Is('}'):
			token = tokens.CreateToken(tokens.Rbrace, "}", c.Line, c.Column)
			l.EatChar()
		case c.Is('('):
			token = tokens.CreateToken(tokens.Lparen, "(", c.Line, c.Column)
			l.EatChar()
		case c.Is(')'):
			token = tokens.CreateToken(tokens.Rparen, ")", c.Line, c.Column)
			l.EatChar()
		case c.Is('['):
			token = tokens.CreateToken(tokens.Lbracket, "[", c.Line, c.Column)
			l.EatChar()
		case c.Is(']'):
			token = tokens.CreateToken(tokens.Rbracket, "]", c.Line, c.Column)
			l.EatChar()

		case l.isWhitespace(c.Rune):
			token = l.parseWhitespaces()
			if token == nil {
				continue
			}

		case c.Is('.') && !l.isDigit(nr):
			token = tokens.CreateToken(tokens.Dot, ".", c.Line, c.Column)
			l.EatChar()

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
