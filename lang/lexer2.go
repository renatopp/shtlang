package lang

// import (
// 	"fmt"
// 	"sht/lang/tokens"
// 	"unicode/utf8"
// )

// type Lexer struct {
// 	input  []byte
// 	tokens []*Token
// 	errors []string
// 	line   int
// 	column int
// 	cursor int
// }

// func Tokenize(input []byte) ([]*Token, error) {
// 	lexer := &Lexer{
// 		input:  input,
// 		tokens: []*Token{},
// 		errors: []string{},
// 		line:   1,
// 		column: 1,
// 		cursor: 0,
// 	}

// 	return lexer.tokenize()
// }

// // func (l *Lexer) tokenize() ([]*Token, error) {
// // 	for {
// // 		c := l.peek()
// // 		line := l.line
// // 		column := l.column

// // 		if c == 0 {
// // 			l.add(tokens.Eof, "", line, column)
// // 			break
// // 		}

// // 		switch{
// // 		case ';':
// // 			l.add(tokens.Semicolon, string(c), line, column)
// // 			l.column++
// // 		case ',':
// // 			l.add(tokens.Comma, string(c), line, column)
// // 			l.column++
// // 		case ':':
// // 			l.add(tokens.Colon, string(c), line, column)
// // 			l.column++
// // 		case '!':
// // 			l.add(tokens.Bang, string(c), line, column)
// // 			l.column++
// // 		case '?':
// // 			l.add(tokens.Question, string(c), line, column)
// // 			l.column++
// // 		case '.':
// // 			if l.isDigit(l.peek()) {
// // 				// parse number
// // 				break
// // 			}
// // 			l.add(tokens.Dot, string(c), line, column)
// // 			l.column++
// // 		case '\\':
// // 			l.add(tokens.Backslash, string(c), line, column)
// // 			l.column++
// // 		case '@':
// // 			l.add(tokens.At, string(c), line, column)
// // 			l.column++
// // 		case '#':
// // 			l.add(tokens.Hash, string(c), line, column)
// // 			l.column++
// // 		case '%':
// // 			l.add(tokens.Percent, string(c), line, column)
// // 			l.column++
// // 		case '^':
// // 			l.add(tokens.Caret, string(c), line, column)
// // 			l.column++
// // 		case '&':
// // 			l.add(tokens.Ampersand, string(c), line, column)
// // 			l.column++
// // 		case '|':
// // 			l.add(tokens.Pipe, string(c), line, column)
// // 			l.column++
// // 		case '+':
// // 			l.add(tokens.Plus, string(c), line, column)
// // 			l.column++
// // 		case '-':
// // 			l.add(tokens.Minus, string(c), line, column)
// // 			l.column++
// // 		case '*':
// // 			l.add(tokens.Asterisk, string(c), line, column)
// // 			l.column++
// // 		case '/':
// // 			l.add(tokens.Slash, string(c), line, column)
// // 			l.column++
// // 		case '>':
// // 			l.add(tokens.Greater, string(c), line, column)
// // 			l.column++
// // 		case '<':
// // 			l.add(tokens.Less, string(c), line, column)
// // 			l.column++
// // 		case '=':
// // 			l.add(tokens.Equal, string(c), line, column)
// // 			l.column++
// // 		case '~':
// // 			l.add(tokens.Tilde, string(c), line, column)
// // 			l.column++
// // 		case '{':
// // 			l.add(tokens.Lbrace, string(c), line, column)
// // 			l.column++
// // 		case '}':
// // 			l.add(tokens.Rbrace, string(c), line, column)
// // 			l.column++
// // 		case '(':
// // 			l.add(tokens.Lparen, string(c), line, column)
// // 			l.column++
// // 		case ')':
// // 			l.add(tokens.Rparen, string(c), line, column)
// // 			l.column++
// // 		case '[':
// // 			l.add(tokens.Lbracket, string(c), line, column)
// // 			l.column++
// // 		case ']':
// // 			l.add(tokens.Rbracket, string(c), line, column)
// // 			l.column++
// // 		case '\n', '\r', '\t', ' ':
// // 			space := l.parseWhitespaces(c)
// // 			switch space {
// // 			case '\n':
// // 				l.add(tokens.Newline, string(space), line, column)
// // 			case ' ':
// // 				l.add(tokens.Space, string(space), line, column)
// // 			}
// // 		case l.isLetter(c):
// // 			parseIdentifier(c)
// // 		// case '\'':
// // 		// parse string
// // 		// case isDigit(c):
// // 		// parse number
// // 		default:
// // 			l.registerError(fmt.Sprintf("invalid character '%c'", c), line, column)
// // 			l.column++
// // 		}
// // 	}

// // 	return l.tokens, nil
// // }

// func (l *Lexer) next() rune {
// 	for l.cursor < len(l.input) {
// 		r, size := utf8.DecodeRune(l.input[l.cursor:])
// 		if r == utf8.RuneError {
// 			l.registerError("invalid UTF-8 character", l.line, l.column)
// 			l.cursor++
// 			l.column++
// 			continue
// 		}

// 		l.cursor += size
// 		return r
// 	}

// 	return 0
// }

// func (l *Lexer) peek() rune {
// 	for l.cursor < len(l.input) {
// 		r, _ := utf8.DecodeRune(l.input[l.cursor:])
// 		if r == utf8.RuneError {
// 			l.registerError("invalid UTF-8 character", l.line, l.column)
// 			l.cursor++
// 			l.column++
// 			continue
// 		}

// 		return r
// 	}

// 	return 0
// }

// func (l *Lexer) add(tp tokens.Type, lit string, line, column int) {
// 	l.tokens = append(l.tokens, CreateToken(tp, lit, line, column))
// }

// func (l *Lexer) registerError(e string, line, column int) {
// 	l.errors = append(l.errors, fmt.Sprintf("%s at %d:%d", e, line, column))
// }

// func (l *Lexer) parseWhitespaces(cur rune) rune {
// 	r := cur
// 	for {
// 		if cur == '\n' {
// 			l.line++
// 			l.column = 1
// 		} else {
// 			l.column++
// 		}

// 		nxt := l.peek()
// 		if !l.isWhitespace(nxt) {
// 			break
// 		}

// 		if nxt == '\n' {
// 			r = '\n'
// 		} else if r != '\n' {
// 			r = ' '
// 		}
// 		cur = l.next()
// 	}

// 	return r
// }

// func (l *Lexer) isWhitespace(r rune) bool {
// 	return r == '\n' || r == '\r' || r == '\t' || r == ' '
// }

// func (l *Lexer) isLetter(r rune) bool {
// 	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '$'
// }

// func (l *Lexer) isDigit(r rune) bool {
// 	return r >= '0' && r <= '9'
// }
