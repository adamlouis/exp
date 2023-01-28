package main

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Expr() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(TokenType_BANG_EQUAL, TokenType_EQUAL_EQUAL) {
		operator := p.previous() // Token
		right := p.comparison()  // Expr
		// 		  expr = new Expr.Binary(expr, operator, right);
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(TokenType_GREATER, TokenType_GREATER_EQUAL, TokenType_LESS, TokenType_LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		//   expr = new Expr.Binary(expr, operator, right);
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(TokenType_MINUS, TokenType_PLUS) {
		operator := p.previous()
		right := p.factor()
		//   expr = new Expr.Binary(expr, operator, right);
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(TokenType_SLASH, TokenType_STAR) {
		operator := p.previous()
		right := p.unary()
		//   expr = new Expr.Binary(expr, operator, right);
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(TokenType_BANG, TokenType_MINUS) {
		operator := p.previous()
		right := p.unary()
		//   return new Expr.Unary(operator, right);
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(TokenType_FALSE) {
		// return new Expr.Literal(false)
	}
	if p.match(TokenType_TRUE) {
		// return new Expr.Literal(true)
	}
	if p.match(TokenType_NIL) {
		// return new Expr.Literal(null)
	}

	if p.match(TokenType_NUMBER, TokenType_STRING) {
		//   return new Expr.Literal(previous().literal);
	}

	if p.match(TokenType_LEFT_PAREN) {
		expr := p.expression()
		p.consume(TokenType_RIGHT_PAREN, "Expect ')' after expression.")
		//   return new Expr.Grouping(expr);
	}

	// TODO:
	// - finishing parsing expressions: https://craftinginterpreters.com/parsing-expressions.html#syntax-errors
	// - generate AST for Go ... using Java reference & requirements above
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().t == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().t == TokenType_EOF
}
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
