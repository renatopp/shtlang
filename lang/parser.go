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
		case "??":
			return order.Unwrapping
		}

	case t.Is(tokens.Keyword):
		switch t.Literal {
		case "as", "is", "in":
			return order.Calls
		}

	case t.Is(tokens.Lparen), t.Is(tokens.Lbrace):
		return order.Calls
	case t.Is(tokens.Lbracket):
		return order.Indexing
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
	p.prefixFns[tokens.String] = p.parsePrefixString
	p.prefixFns[tokens.Bang] = p.parsePrefixOperator
	p.prefixFns[tokens.Operator] = p.parsePrefixOperator
	p.prefixFns[tokens.Lparen] = p.parsePrefixParenthesis
	p.prefixFns[tokens.Identifier] = p.parsePrefixIdentifier
	p.prefixFns[tokens.Spread] = p.parsePrefixSpread

	p.infixFns[tokens.Operator] = p.parseInfixOperator
	p.infixFns[tokens.Keyword] = p.parseInfixKeyword
	p.infixFns[tokens.Lparen] = p.parseInfixCall
	p.infixFns[tokens.Lbrace] = p.parseInfixCall
	p.infixFns[tokens.Lbracket] = p.parseInfixBracket

	p.postfixFns[tokens.Operator] = p.parsePostfixOperator
	p.postfixFns[tokens.Bang] = p.parsePostfixOperator
	p.postfixFns[tokens.Question] = p.parsePostfixOperator
	p.postfixFns[tokens.Spread] = p.parsePostfixOperator

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

	} else if cur.Is(tokens.Keyword) &&
		cur.Literal == "return" ||
		cur.Literal == "raise" ||
		cur.Literal == "yield" {
		return p.parseReturn()

	} else if cur.Is(tokens.Keyword) && cur.Literal == "for" {
		// parse for

	} else if cur.Is(tokens.Keyword) && cur.Literal == "if" {
		// parse if

	} else if cur.Is(tokens.Keyword) && (cur.Literal == "let" || cur.Literal == "const") {
		return p.parseDeclaration()
	}

	e := p.parseExpressionTuple()
	if e == nil {
		p.RegisterError(fmt.Sprintf("invalid token '%s'", cur.Literal), cur)
		return nil
	}

	cur = p.lexer.PeekToken()
	if cur.Is(tokens.Assignment) {
		return p.parseAssignment(e)
	}

	cur = p.lexer.PeekToken()
	if !isEndOfStatement(cur) {
		p.RegisterError(fmt.Sprintf("unexpected token '%s'", cur.Literal), cur)
		return nil
	}

	return e
}

func (p *Parser) parseReturn() ast.Node {
	cur := p.lexer.PeekToken()
	p.lexer.EatToken()

	exp := p.parseSingleExpression(order.Lowest)

	switch cur.Literal {
	case "return":
		return &ast.Return{
			Token:      cur,
			Expression: exp,
		}
	case "raise":
		return &ast.Raise{
			Token:      cur,
			Expression: exp,
		}
	case "yield":
		return &ast.Yield{
			Token:      cur,
			Expression: exp,
		}
	default:
		p.RegisterError(fmt.Sprintf("invalid return token '%s'", cur.Literal), cur)
		return nil
	}
}

func (p *Parser) parseDeclaration() ast.Node {
	cur := p.lexer.PeekToken()
	p.lexer.EatToken()
	constant := cur.Literal == "const"

	left := p.parseExpressionTuple()
	if !p.Expect(tokens.Assignment) {
		return nil
	}

	node := p.parseAssignment(left)
	if node == nil {
		return nil
	}

	decl := node.(*ast.Assignment)
	decl.Definition = true
	decl.Constant = constant

	return decl
}

func (p *Parser) parseAssignment(left ast.Node) ast.Node {
	if left == nil {
		p.RegisterError(fmt.Sprintf("invalid assignment target"), p.lexer.PeekToken())
		return nil
	}

	switch tp := left.(type) {
	case *ast.Tuple:
		left = tp
	default:
		left = &ast.Tuple{
			Token:  left.GetToken(),
			Values: []ast.Node{tp},
		}
	}

	inv, err := p.assertAssignmentTargets(left)
	if err != "" {
		p.RegisterError(err, inv.GetToken())
	}

	ids, _ := left.(*ast.Tuple)
	ass := p.lexer.PeekToken()
	p.lexer.EatToken()

	exp := p.parseExpressionTuple()

	if exp == nil {
		p.RegisterError(fmt.Sprintf("expected expression, got %s instead", p.lexer.PeekToken()), ass)
	}

	return &ast.Assignment{
		Token:      ass,
		Identifier: ids,
		Literal:    ass.Literal,
		Expression: exp,
		Definition: false,
	}
}

func (p *Parser) assertAssignmentTargets(t ast.Node) (ast.Node, string) {
	if t == nil {
		return t, "no left targets"
	}

	hasSpread := false
	switch t := t.(type) {
	case *ast.Identifier:
		return t, ""

	case *ast.Tuple:
		for _, v := range t.Values {
			_, ok := v.(*ast.SpreadIn)

			if ok && hasSpread {
				return t, "left-side assignments can have only one spread operator"
			}
			if ok {
				hasSpread = true
			}

			r, err := p.assertAssignmentTargets(v)
			if err != "" {
				return r, err
			}
		}

	case *ast.SpreadIn:
		return p.assertAssignmentTargets(t.Target)

	default:
		return t, fmt.Sprintf("invalid left-side assignment token '%s'", t.GetToken().Literal)
	}

	return t, ""
}

func (p *Parser) parseBoolean() ast.Node {
	t := p.lexer.EatToken()
	return &ast.Boolean{
		Token: t,
		Value: t.Literal == "true",
	}
}

func (p *Parser) parseFunctionDef() ast.Node {
	cur := p.lexer.EatToken()
	if !p.Expect(tokens.Identifier, tokens.Lparen, tokens.Lbrace, tokens.Question) {
		p.RegisterError(fmt.Sprintf("invalid function definition"), p.lexer.PeekToken())
		return nil
	}

	fn := &ast.FunctionDef{
		Token: cur,
	}

	cur = p.lexer.PeekToken()
	if cur.Is(tokens.Identifier) {
		p.lexer.EatToken()
		fn.Name = cur.Literal
	}

	cur = p.lexer.PeekToken()
	if cur.Is(tokens.Lparen) {
		fn.Params = p.parseParameters()
	}

	cur = p.lexer.PeekToken()
	if cur.Is(tokens.Question) {
		p.lexer.EatToken()
		fn.Maybe = true
	}

	cur = p.lexer.PeekToken()
	if !p.Expect(tokens.Lbrace) {
		p.RegisterError(fmt.Sprintf("invalid function definition"), p.lexer.PeekToken())
		return nil
	}

	if cur.Is(tokens.Lbrace) {
		fn.Body = p.parseBlock()
	}

	return fn
}

func (p *Parser) parseParameters() []ast.Node {
	cur := p.lexer.PeekToken()
	if cur.Is(tokens.Lparen) {
		p.lexer.EatToken()
	}

	params := []ast.Node{}

	cur = p.lexer.PeekToken()
	for !cur.Is(tokens.Rparen) {
		p.eatNewLines()

		if !p.Expect(tokens.Identifier) {
			p.RegisterError(fmt.Sprintf("invalid parameter token '%s'", cur.Literal), cur)
			return nil
		}

		p.lexer.EatToken()
		param := &ast.Parameter{
			Token: cur,
			Name:  cur.Literal,
		}

		cur = p.lexer.PeekToken()
		if cur.Is(tokens.Spread) {
			p.lexer.EatToken()
			param.Spread = true
		}

		cur = p.lexer.PeekToken()
		if cur.Is(tokens.Assignment) && cur.Literal == "=" {
			p.lexer.EatToken()
			param.Default = p.parseLiteral()
		}

		cur = p.lexer.PeekToken()
		if cur.Is(tokens.Comma) {
			p.lexer.EatToken()
		}

		params = append(params, param)

		if !p.Expect(tokens.Identifier, tokens.Rparen, tokens.Newline) {
			p.RegisterError(fmt.Sprintf("invalid parameter token '%s'", cur.Literal), cur)
			return nil
		}

		cur = p.lexer.PeekToken()
	}

	p.lexer.EatToken() // )
	return params
}

func (p *Parser) parseLiteral() ast.Node {
	cur := p.lexer.PeekToken()

	switch {
	case cur.Is(tokens.Number):
		return p.parsePrefixNumber()
	case cur.Is(tokens.String):
		return p.parsePrefixString()
	case cur.Is(tokens.Keyword) && (cur.Literal == "true" || cur.Literal == "false"):
		return p.parseBoolean()
	default:
		p.RegisterError(fmt.Sprintf("invalid literal '%s'", cur.Literal), cur)
		return nil
	}
}

func (p *Parser) parseExpressionTuple() ast.Node {
	args := make([]ast.Node, 0)

	cur := p.lexer.PeekToken()
	for {
		p.eatNewLines()

		arg := p.parseSingleExpression(order.Lowest)
		if arg == nil {
			break
		}

		args = append(args, arg)

		cur = p.lexer.PeekToken()
		if cur.Is(tokens.Comma) {
			p.lexer.EatToken()
		} else {
			break
		}
	}

	if len(args) == 0 {
		return nil
	} else if len(args) == 1 {
		return args[0]
	} else {
		return &ast.Tuple{
			Token:  p.lexer.PeekToken(),
			Values: args,
		}
	}
}

func (p *Parser) parseExpressionList() []ast.Node {
	args := make([]ast.Node, 0)

	cur := p.lexer.PeekToken()
	for !cur.Is(tokens.Rparen) {
		p.eatNewLines()

		arg := p.parseSingleExpression(order.Lowest)
		if arg == nil {
			break
		}

		args = append(args, arg)

		cur = p.lexer.PeekToken()
		if cur.Is(tokens.Comma) {
			p.lexer.EatToken()
		}
	}

	return args
}

func (p *Parser) parseSingleExpression(priority int) ast.Node {
	cur := p.lexer.PeekToken()
	// fmt.Println("parsing expression", priority, "-", cur)

	prefixFn := p.prefixFns[cur.Type]

	if prefixFn == nil {
		return nil
	}

	left := prefixFn()
	// fmt.Println("prefix", left)

	cur = p.lexer.PeekToken()

repeat_infix:
	for {
		// fmt.Println("for infix", priority, priorityOf(cur), cur)
		for !isEndOfExpression(cur) && priority < priorityOf(cur) {
			infixFn := p.infixFns[cur.Type]
			// fmt.Println("checking infix", left, infixFn, cur)
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
		// fmt.Println("for postfix", priority, priorityOf(cur), cur)
		for isPostfix(cur) {
			postfixFn := p.postfixFns[cur.Type]
			// fmt.Println("has postfix?", postfixFn, cur.Type)
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

func (p *Parser) parseInitializer() ast.Initializer {
	init := p.lexer.PeekToken()
	if !p.Expect(tokens.Lbrace) {
		p.RegisterError(fmt.Sprintf("invalid initializer"), init)
		return nil
	}

	tp := '?'
	p.lexer.EatToken()
	cur := p.lexer.PeekToken()
	var initializer ast.Initializer
	for !cur.Is(tokens.Rbrace) {
		p.eatNewLines()

		first := p.parseSingleExpression(order.Lowest)

		if first == nil {
			break
		}

		if tp == '?' {
			cur = p.lexer.PeekToken()
			if cur.Is(tokens.Colon) {
				tp = 'M'
				initializer = &ast.MapInitializer{
					Token:  init,
					Values: map[string]ast.Node{},
				}
			} else {
				tp = 'L'
				initializer = &ast.ListInitializer{
					Token:  init,
					Values: []ast.Node{},
				}
			}
		}

		if tp == 'M' {
			name, ok := first.(*ast.Identifier)
			if !ok {
				p.RegisterError(fmt.Sprintf("invalid map initializer key '%s'", cur.Literal), cur)
				return nil
			}

			if !p.Expect(tokens.Colon) {
				return nil
			}

			p.lexer.EatToken()
			p.eatNewLines()

			exp := p.parseSingleExpression(order.Lowest)
			if exp == nil {
				p.RegisterError(fmt.Sprintf("invalid map initializer value '%s'", cur.Literal), cur)
				return nil
			}
			initializer.(*ast.MapInitializer).Values[name.Value] = exp
		} else {
			initializer.(*ast.ListInitializer).Values = append(initializer.(*ast.ListInitializer).Values, first)
		}

		cur = p.lexer.PeekToken()
		if cur.Is(tokens.Comma) {
			p.lexer.EatToken()
		}
	}

	cur = p.lexer.PeekToken()
	if !p.Expect(tokens.Rbrace) {
		p.RegisterError(fmt.Sprintf("expecting '}' character to end initialized, got '%s' instead", cur.Literal), cur)
		return nil
	}
	p.lexer.EatToken()

	return initializer
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

	case "fn":
		return p.parseFunctionDef()

	default:
		p.RegisterError(fmt.Sprintf("invalid keyword '%s'", cur.Literal), cur)
		return nil
	}
}

func (p *Parser) parsePrefixNumber() ast.Node {
	cur := p.lexer.PeekToken()
	v, e := strconv.ParseFloat(cur.Literal, 64)

	p.lexer.EatToken()
	if e != nil {
		p.RegisterError(fmt.Sprintf("invalid number '%s'", cur.Literal), cur)
		return nil
	}

	return &ast.Number{
		Token: cur,
		Value: v,
	}
}

func (p *Parser) parsePrefixString() ast.Node {
	cur := p.lexer.PeekToken()

	p.lexer.EatToken()
	return &ast.String{
		Token: cur,
		Value: cur.Literal,
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
		Right:    p.parseSingleExpression(order.Unary),
	}
}

func (p *Parser) parsePrefixParenthesis() ast.Node {
	p.lexer.EatToken()
	e := p.parseSingleExpression(order.Lowest)
	p.Expect(tokens.Rparen, tokens.Comma)

	if p.lexer.PeekToken().Is(tokens.Comma) {
		p.lexer.EatToken()
		e = &ast.Tuple{
			Token:  p.lexer.PeekToken(),
			Values: append([]ast.Node{e}, p.parseExpressionList()...),
		}
	}

	p.lexer.EatToken()
	return e
}

func (p *Parser) parsePrefixIdentifier() ast.Node {
	cur := p.lexer.PeekToken()
	// fmt.Println("parsePrefixIdentifier", cur)
	p.lexer.EatToken()
	return &ast.Identifier{
		Token: cur,
		Value: cur.Literal,
	}
}

func (p *Parser) parsePrefixSpread() ast.Node {
	cur := p.lexer.PeekToken()
	p.lexer.EatToken()
	return &ast.SpreadIn{
		Token:  cur,
		Target: p.parseSingleExpression(order.Spread),
	}
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
		Right:    p.parseSingleExpression(priority),
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
		Right:    p.parseSingleExpression(priority),
	}
}

func (p *Parser) parseInfixCall(left ast.Node) ast.Node {
	node := &ast.Call{
		Target: left,
	}

	if p.lexer.PeekToken().Is(tokens.Lparen) {
		p.lexer.EatToken()
		node.Arguments = p.parseExpressionList()

		if !p.Expect(tokens.Rparen) {
			return nil
		}
		p.lexer.EatToken()
	}

	cur := p.lexer.PeekToken()
	if cur.Is(tokens.Lbrace) {
		node.Initializer = p.parseInitializer()
	}

	return node
}

func (p *Parser) parseInfixBracket(left ast.Node) ast.Node {
	p.lexer.EatToken()

	node := &ast.Indexing{
		Target: left,
		Values: p.parseExpressionList(),
	}

	p.Expect(tokens.Rbracket)
	p.lexer.EatToken()
	return node
}

// ----------------------------------------------------------------
// Postfix Functions
// ----------------------------------------------------------------
func (p *Parser) parsePostfixOperator(left ast.Node) ast.Node {
	cur := p.lexer.PeekToken()
	p.lexer.EatToken()

	if cur.Is(tokens.Question) {
		return &ast.Wrapping{
			Token:      cur,
			Expression: left,
		}

	} else if cur.Is(tokens.Bang) {
		return &ast.Unwrapping{
			Token:  cur,
			Target: left,
		}

	} else if cur.Is(tokens.Spread) {
		return &ast.SpreadOut{
			Token:  cur,
			Target: left,
		}
	}

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
	case "++", "--", "!", "?", "...":
		return true
	}

	return false
}

func isInfix(t *tokens.Token) bool {
	switch t.Literal {
	case "+", "-", "*", "/", "//", "%", "**", "==", "!=", ">", "<", ">=", "<=", "and", "nand", "or", "xor", "nor", "nxor", "..", "??":
		return true
	}

	return false
}
