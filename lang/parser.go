package lang

import (
	"fmt"
	"sht/lang/ast"
	"sht/lang/order"
	"sht/lang/tokens"
	"strings"

	"golang.org/x/exp/slices"
)

func priorityOf(t *Token) int {
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
		case "and":
			return order.And
		case "nand":
			return order.And
		case "or":
			return order.Or
		case "xor":
			return order.Or
		case "nor":
			return order.Or
		case "nxor":
			return order.Or
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

type Parser struct {
	lexer     *Lexer
	root      ast.Node
	prefixFns map[tokens.Type]prefixFn
	infixFns  map[tokens.Type]infixFn
	errors    []string
}

func CreateParser() *Parser {
	return &Parser{
		lexer:     nil,
		root:      nil,
		prefixFns: map[tokens.Type]prefixFn{},
		infixFns:  map[tokens.Type]infixFn{},
		errors:    []string{},
	}
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

func (p *Parser) RegisterError(e string, t *Token) {
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

func (p *Parser) parseBlock() ast.Node {
	block := &ast.Block{}

	t := p.lexer.PeekToken()
	braced := t.Is(tokens.Lbrace)
	if braced {
		p.lexer.EatToken()
	}

	cur := p.lexer.PeekToken()
	for !(cur.Is(tokens.Rbrace) || cur.Is(tokens.Eof)) {
		p.lexer.EatToken()
		cur = p.lexer.PeekToken()
	}

	if braced {
		p.Expect(tokens.Rbrace)
		p.lexer.EatToken()
	}

	return block
}
