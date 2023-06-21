package lang

import (
	"fmt"
	"sht/lang/ast"
	"sht/lang/tokens"
	"strings"
)

type prefixFn func() ast.Node
type infixFn func(ast.Node) ast.Node

type Parser struct {
	lexer     *Lexer
	root      ast.Node
	prefixFns map[tokens.Type]prefixFn
	infixFns  map[tokens.Type]infixFn
	errors    []string
}

func CreateParser(input []byte) *Parser {
	return &Parser{
		lexer:     CreateLexer(input),
		root:      nil,
		prefixFns: map[tokens.Type]prefixFn{},
		infixFns:  map[tokens.Type]infixFn{},
		errors:    []string{},
	}
}

func Parse(input []byte) (ast.Node, error) {
	p := CreateParser(input)
	p.root = p.parseBlock()
	return p.root, p.GetError()
}

func (p *Parser) GetError() error {
	if p.HasError() {
		return fmt.Errorf("Tokenizer errors: \n- %s", strings.Join(p.errors, "\n- "))
	}

	return nil
}

func (p *Parser) HasError() bool {
	return len(p.errors) > 0
}

func (p *Parser) TooManyErrors() bool {
	return len(p.errors) >= 10
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

func (p *Parser) parseBlock() ast.Node {
	block := &ast.Block{}

	//

	return block
}
