package main

import "fmt"

type Parser struct {
	lox     *Lox
	tokens  []Token
	current int
}

func NewParser(l *Lox, tokens []Token) *Parser {
	return &Parser{
		lox:     l,
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) parse() []*Stmt {
	ret := []*Stmt{}
	defer func() {
		_ = recover()
	}()

	for !p.isAtEnd() {
		ret = append(ret, p.declaration())
	}

	return ret
}

func (p *Parser) declaration() *Stmt {
	var ret *Stmt

	defer func() {
		if r := recover(); r != nil {
			_ = recover()
		}
		p.synchronize()
		ret = nil
	}()

	if p.match(TokenType_VAR) {
		ret = p.varDeclaration()
	} else {
		ret = p.statement()
	}
	return ret
}

func (p *Parser) varDeclaration() *Stmt {
	name := p.consume(TokenType_IDENTIFIER, "Expect variable name.")

	var initializer *Expr
	if p.match(TokenType_EQUAL) {
		i := p.expression()
		initializer = &i
	}

	p.consume(TokenType_SEMICOLON, "Expect ';' after variable declaration.")
	return &Stmt{
		Var: &Var{name, initializer},
	}
}

func (p *Parser) statement() *Stmt {
	if p.match(TokenType_PRINT) {
		return p.printStatement()
	}
	if p.match(TokenType_LEFT_BRACE) {
		return &Stmt{
			Block: &Block{p.block()},
		}
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() *Stmt {
	value := p.expression()
	p.consume(TokenType_SEMICOLON, "Expect ';' after value.")
	return &Stmt{
		Print: &Print{value},
	}
}

func (p *Parser) expressionStatement() *Stmt {
	expr := p.expression()
	p.consume(TokenType_SEMICOLON, "Expect ';' after expression.")
	return &Stmt{
		Expression: &Expression{expr},
	}
}

func (p *Parser) block() []*Stmt {
	statements := []*Stmt{}

	fmt.Println("BLOCK START")
	for !p.check(TokenType_RIGHT_BRACE) && !p.isAtEnd() {
		fmt.Println("block loop")
		statements = append(statements, p.declaration())
	}
	fmt.Println("BLOCK END")

	p.consume(TokenType_RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) expression() Expr {
	return p.assignment()
}
func (p *Parser) assignment() Expr {
	expr := p.equality()

	if p.match(TokenType_EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if expr.Variable != nil {
			return Expr{
				Assign: &Assign{expr.Variable.Name, value},
			}
		}

		p.error(equals, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(TokenType_BANG_EQUAL, TokenType_EQUAL_EQUAL) {
		operator := p.previous() // Token
		right := p.comparison()  // Expr
		expr = Expr{
			Binary: &Binary{expr, operator, right},
		}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(TokenType_GREATER, TokenType_GREATER_EQUAL, TokenType_LESS, TokenType_LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Expr{
			Binary: &Binary{expr, operator, right},
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(TokenType_MINUS, TokenType_PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Expr{
			Binary: &Binary{expr, operator, right},
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(TokenType_SLASH, TokenType_STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Expr{
			Binary: &Binary{expr, operator, right},
		}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(TokenType_BANG, TokenType_MINUS) {
		operator := p.previous()
		right := p.unary()
		return Expr{
			Unary: &Unary{operator, right},
		}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(TokenType_FALSE) {
		return Expr{Literal: &Literal{Value: false}}
	}
	if p.match(TokenType_TRUE) {
		return Expr{Literal: &Literal{Value: true}}
	}
	if p.match(TokenType_NIL) {
		return Expr{Literal: &Literal{Value: nil}}
	}

	if p.match(TokenType_NUMBER, TokenType_STRING) {
		return Expr{Literal: &Literal{Value: p.previous().literal}}
	}

	if p.match(TokenType_IDENTIFIER) {
		return Expr{
			Variable: &Variable{p.previous()},
		}
	}

	if p.match(TokenType_LEFT_PAREN) {
		expr := p.expression()
		p.consume(TokenType_RIGHT_PAREN, "Expect ')' after expression.")
		return Expr{
			Grouping: &Grouping{expr},
		}
	}

	panic(p.error(p.peek(), "Expect expression."))
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

func (p *Parser) consume(t TokenType, message string) Token {
	if p.check(t) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

func (p *Parser) error(token Token, message string) *ParseError {
	p.lox.error(token, message)
	return &ParseError{}
}

func (l *Lox) error(token Token, message string) {
	if token.t == TokenType_EOF {
		l.report(token.line, " at end", message)
	} else {
		l.report(token.line, " at '"+token.lexeme+"'", message)
	}
}

type ParseError struct{}

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

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().t == TokenType_SEMICOLON {
			return
		}

		switch p.peek().t {
		case TokenType_CLASS:
		case TokenType_FUN:
		case TokenType_VAR:
		case TokenType_FOR:
		case TokenType_IF:
		case TokenType_WHILE:
		case TokenType_PRINT:
		case TokenType_RETURN:
			return
		}
		p.advance()
	}
}
