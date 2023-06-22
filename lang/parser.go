package lang

import (
	"fmt"
	"sht/lang/ast"
	"sht/lang/order"
	"sht/lang/tokens"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

func priorityOf(t *tokens.Token) int {
	switch {
	case t.Is(tokens.Operator):
		switch t.Literal {
		case "+":
			return order.Addition
		case "-":
			return order.Subtraction
		case "*":
			return order.Multiplication
		case "/":
			return order.Division
		case "//":
			return order.Division
		case "%":
			return order.Mod
		case "**":
			return order.Exponentiation
		case "==":
			return order.Comparison
		case "!=":
			return order.Comparison
		case ">":
			return order.Comparison
		case "<":
			return order.Comparison
		case ">=":
			return order.Comparison
		case "<=":
			return order.Comparison
		case "!":
			return order.Not
		case "and", "nand":
			return order.And
		case "or", "xor", "nor", "nxor":
			return order.Or
		case "++", "--":
			return order.Postfix
		case "..":
			return order.Concat
		}

	case t.Is(tokens.Keyword):
		switch t.Literal {
		case "as", "is", "in":
			return order.Calls
		}

	case t.Is(tokens.Assignment):
		return order.Assign
	case t.Is(tokens.Lparen):
		return order.Calls
	case t.Is(tokens.Lbracket):
		return order.Indexing
	case t.Is(tokens.Dot):
		return order.Chain
	case t.Is(tokens.Question):
		return order.Conditional
	default:
		return order.Lowest
	}

	return order.Lowest
}

type prefixFn func() ast.Node
type infixFn func(ast.Node) ast.Node
type postfixFn func(ast.Node) ast.Node

type Parser struct {
	lexer      *Lexer
	root       ast.Node
	prefixFns  map[tokens.Type]prefixFn
	infixFns   map[tokens.Type]infixFn
	postfixFns map[tokens.Type]postfixFn
	errors     []string
}

func CreateParser() *Parser {
	p := &Parser{
		lexer:      nil,
		root:       nil,
		prefixFns:  map[tokens.Type]prefixFn{},
		infixFns:   map[tokens.Type]infixFn{},
		postfixFns: map[tokens.Type]postfixFn{},
		errors:     []string{},
	}

	p.prefixFns[tokens.Keyword] = p.parsePrefixKeyword
	p.prefixFns[tokens.Number] = p.parsePrefixNumber
	p.prefixFns[tokens.Bang] = p.parsePrefixOperator
	p.prefixFns[tokens.Operator] = p.parsePrefixOperator
	p.prefixFns[tokens.Lparen] = p.parsePrefixParenthesis

	p.infixFns[tokens.Operator] = p.parseInfixOperator
	p.infixFns[tokens.Operator] = p.parseInfixOperator
	p.infixFns[tokens.Keyword] = p.parseInfixKeyword

	p.postfixFns[tokens.Operator] = p.parsePostfixOperator
	p.postfixFns[tokens.Bang] = p.parsePostfixOperator

	return p
}

func Parse(input []byte) (ast.Node, error) {
	p := CreateParser()
	return p.Parse(input)
}

func (p *Parser) Parse(input []byte) (ast.Node, error) {
	p.lexer = CreateLexer(input)
	p.root = p.parseBlock()
	p.Expect(tokens.Eof)

	return p.root, p.GetError()
}

func (p *Parser) TooManyErrors() bool {
	return len(p.errors) >= 10
}

func (p *Parser) HasError() bool {
	return len(p.errors) > 0
}

func (p *Parser) GetError() error {
	if p.HasError() {
		return fmt.Errorf("Parser errors: \n- %s", strings.Join(p.errors, "\n- "))
	}

	return nil
}

func (p *Parser) RegisterError(e string, t *tokens.Token) {
	if p.TooManyErrors() {
		return
	}

	p.errors = append(p.errors, fmt.Sprintf("%s at %d:%d", e, t.Line, t.Column))
	if p.TooManyErrors() {
		p.errors = append(p.errors, "too many errors, aborting")
	}
}

func (p *Parser) Expect(tps ...tokens.Type) bool {
	t := p.lexer.PeekToken()
	tp := t.Type
	if !slices.Contains(tps, tp) {
		p.RegisterError(fmt.Sprintf("expected %s, got %s", tokens.JoinTypes(tps...), tp), t)
		return false
	}

	return true
}

func (p *Parser) eatNewLines() {
	for p.lexer.PeekToken().Is(tokens.Newline) {
		p.lexer.EatToken()
	}
}

// ----------------------------------------------------------------
// Parsing Functions
// ----------------------------------------------------------------

func (p *Parser) parseBlock() ast.Node {
	block := &ast.Block{}

	t := p.lexer.PeekToken()
	braced := t.Is(tokens.Lbrace)
	if braced {
		p.lexer.EatToken()
	}

	cur := p.lexer.PeekToken()
	for !isEndOfBlock(cur) && !p.HasError() {
		s := p.parseStatement()
		if s != nil {
			block.Statements = append(block.Statements, s)
		}
		cur = p.lexer.PeekToken()
	}

	if braced {
		p.Expect(tokens.Rbrace)
		p.lexer.EatToken()
	}

	return block
}

func (p *Parser) parseStatement() ast.Node {
	p.eatNewLines()

	cur := p.lexer.PeekToken()
	if isEndOfBlock(cur) {
		return nil

	} else if cur.Is(tokens.Semicolon) {
		p.lexer.EatToken()
		return nil

	} else if cur.Is(tokens.Keyword) && cur.Literal == "return" {
		// parse return

	} else if cur.Is(tokens.Keyword) && cur.Literal == "for" {
		// parse for

	} else if cur.Is(tokens.Keyword) && cur.Literal == "if" {
		// parse if

	} else if cur.Is(tokens.Keyword) && (cur.Literal == "let" || cur.Literal == "const") {
		// parse declaration

	}

	e := p.parseExpression(order.Lowest)
	if e == nil {
		p.RegisterError(fmt.Sprintf("invalid token '%s'", cur.Literal), cur)
		return nil
	}

	cur = p.lexer.PeekToken()
	if !isEndOfStatement(cur) {
		p.RegisterError(fmt.Sprintf("unexpected token '%s'", cur.Literal), cur)
		return nil
	}

	return e
}

func (p *Parser) parseExpression(priority int) ast.Node {
	cur := p.lexer.PeekToken()

	prefixFn := p.prefixFns[cur.Type]

	if prefixFn == nil {
		return nil
	}

	left := prefixFn()
	// fmt.Println("prefix", left)

	cur = p.lexer.PeekToken()

repeat_infix:
	for {
		for !isEndOfExpression(cur) && priority < priorityOf(cur) && isInfix(cur) {
			infixFn := p.infixFns[cur.Type]
			if infixFn == nil {
				return left
			}

			newleft := infixFn(left)
			if newleft == nil {
				return left
			}
			left = newleft
			// fmt.Println("infix", left)

			cur = p.lexer.PeekToken()
		}

		cur = p.lexer.PeekToken()
		for isPostfix(cur) {

			postfixFn := p.postfixFns[cur.Type]
			if postfixFn == nil {
				return left
			}

			newleft := postfixFn(left)
			if newleft == nil {
				return left
			}
			left = newleft

			cur = p.lexer.PeekToken()
			// fmt.Println("postfix", left)

			continue repeat_infix
		}

		break
	}

	return left
}

func (p *Parser) parseBoolean() ast.Node {
	t := p.lexer.EatToken()
	return &ast.Boolean{
		Token: t,
		Value: t.Literal == "true",
	}
}

// ----------------------------------------------------------------
// Prefix Functions
// ----------------------------------------------------------------
func (p *Parser) parsePrefixKeyword() ast.Node {
	cur := p.lexer.PeekToken()

	switch cur.Literal {
	case "true", "false":
		return p.parseBoolean()

	case "!":
		return p.parsePrefixOperator()

	default:
		p.RegisterError(fmt.Sprintf("invalid keyword '%s'", cur.Literal), cur)
		return nil
	}
}

func (p *Parser) parsePrefixNumber() ast.Node {
	cur := p.lexer.PeekToken()
	v, e := strconv.ParseFloat(cur.Literal, 64)

	if e != nil {
		panic(e)
	}

	p.lexer.EatToken()
	return &ast.Number{
		Token: cur,
		Value: v,
	}
}

func (p *Parser) parsePrefixOperator() ast.Node {
	cur := p.lexer.PeekToken()

	if !isUnary(cur) {
		p.RegisterError(fmt.Sprintf("invalid unary operator '%s'", cur.Literal), cur)
		return nil
	}

	p.lexer.EatToken()
	return &ast.UnaryOperator{
		Token:    cur,
		Operator: cur.Literal,
		Right:    p.parseExpression(order.Unary),
	}
}

func (p *Parser) parsePrefixParenthesis() ast.Node {
	p.lexer.EatToken()
	e := p.parseExpression(order.Lowest)
	p.Expect(tokens.Rparen)
	p.lexer.EatToken()
	return e
}

// ----------------------------------------------------------------
// Infix Functions
// ----------------------------------------------------------------

func (p *Parser) parseInfixOperator(left ast.Node) ast.Node {
	cur := p.lexer.PeekToken()
	priority := priorityOf(cur)

	p.lexer.EatToken()
	return &ast.BinaryOperator{
		Token:    cur,
		Operator: cur.Literal,
		Left:     left,
		Right:    p.parseExpression(priority),
	}
}

func (p *Parser) parseInfixKeyword(left ast.Node) ast.Node {
	cur := p.lexer.PeekToken()
	priority := priorityOf(cur)

	p.lexer.EatToken()
	return &ast.BinaryOperator{
		Token:    cur,
		Operator: cur.Literal,
		Left:     left,
		Right:    p.parseExpression(priority),
	}
}

// ----------------------------------------------------------------
// Postfix Functions
// ----------------------------------------------------------------

func (p *Parser) parsePostfixOperator(left ast.Node) ast.Node {
	cur := p.lexer.PeekToken()
	p.lexer.EatToken()

	return &ast.PostfixOperator{
		Token:    cur,
		Operator: cur.Literal,
		Left:     left,
	}
}

// ----------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------

func isEndOfBlock(t *tokens.Token) bool {
	return t.Is(tokens.Rbrace) || t.Is(tokens.Eof)
}

func isEndOfStatement(t *tokens.Token) bool {
	return t.Is(tokens.Semicolon) || t.Is(tokens.Eof) || t.Is(tokens.Newline) || t.Is(tokens.Rbrace)
}

func isEndOfExpression(t *tokens.Token) bool {
	return t.Is(tokens.Semicolon) // t.Is(token.Newline) ||
}

func isUnary(t *tokens.Token) bool {
	switch t.Literal {
	case "+", "-", "!":
		return true
	}

	return false
}

func isPostfix(t *tokens.Token) bool {
	switch t.Literal {
	case "++", "--", "!":
		return true
	}

	return false
}

func isInfix(t *tokens.Token) bool {
	switch t.Literal {
	case "+", "-", "*", "/", "//", "%", "**", "==", "!=", ">", "<", ">=", "<=", "and", "nand", "or", "xor", "nor", "nxor", "..":
		return true
	}

	return false
}
