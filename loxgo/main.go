package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    TokenType_AND,
	"class":  TokenType_CLASS,
	"else":   TokenType_ELSE,
	"false":  TokenType_FALSE,
	"for":    TokenType_FOR,
	"fun":    TokenType_FUN,
	"if":     TokenType_IF,
	"nil":    TokenType_NIL,
	"or":     TokenType_OR,
	"print":  TokenType_PRINT,
	"return": TokenType_RETURN,
	"super":  TokenType_SUPER,
	"this":   TokenType_THIS,
	"true":   TokenType_TRUE,
	"var":    TokenType_VAR,
	"while":  TokenType_WHILE,
}

func main() {
	// take a pass at end to make java patterns idomatic go
	l := &Lox{
		interpreter: NewInterpreter(nil),
	}
	l.interpreter.lox = l

	switch len(os.Args) {
	case 2:
		if err := l.runFile(os.Args[1]); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	case 1:
		if err := l.runPrompt(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	default:
		panic("usage: loxgo [script]")
	}
}

type Lox struct {
	interpreter     *Interpreter
	hadError        bool
	hadRuntimeError bool
}

func (l *Lox) runFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if err := l.run(string(b)); err != nil {
		return err
	}

	if l.hadError {
		os.Exit(65)
	}
	if l.hadRuntimeError {
		os.Exit(65)
	}
	return nil
}

func (l *Lox) runPrompt() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return err
		}
		if isPrefix {
			panic(fmt.Errorf("TODO: handle prefix at line: %s", line))
		}
		if err := l.run(string(line)); err != nil {
			return err
		}
		l.hadError = false
	}
}

func (l *Lox) reportErr(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}

func (l *Lox) runtimeError(v any, token Token) {
	fmt.Printf("%v [line %d]\n", v, token.line)
	l.hadRuntimeError = true
}

func (l *Lox) run(source string) error {
	scanner := NewScanner(l, source)
	tokens, err := scanner.scanTokens()
	if err != nil {
		return err
	}

	parser := NewParser(l, tokens)
	statements := parser.parse()

	// Stop if there was a syntax error.
	if l.hadError {
		return nil
	}

	return l.interpreter.interpret(statements)
}

type TokenType string

const (
	// Single-character tokens.
	TokenType_LEFT_PAREN  TokenType = "LEFT_PAREN"
	TokenType_RIGHT_PAREN TokenType = "RIGHT_PAREN"
	TokenType_LEFT_BRACE  TokenType = "LEFT_BRACE"
	TokenType_RIGHT_BRACE TokenType = "RIGHT_BRACE"
	TokenType_COMMA       TokenType = "COMMA"
	TokenType_DOT         TokenType = "DOT"
	TokenType_MINUS       TokenType = "MINUS"
	TokenType_PLUS        TokenType = "PLUS"
	TokenType_SEMICOLON   TokenType = "SEMICOLON"
	TokenType_SLASH       TokenType = "SLASH"
	TokenType_STAR        TokenType = "STAR"
	// One or two character tokens.
	TokenType_BANG          TokenType = "BANG"
	TokenType_BANG_EQUAL    TokenType = "BANG_EQUAL"
	TokenType_EQUAL         TokenType = "EQUAL"
	TokenType_EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	TokenType_GREATER       TokenType = "GREATER"
	TokenType_GREATER_EQUAL TokenType = "GREATER_EQUAL"
	TokenType_LESS          TokenType = "LESS"
	TokenType_LESS_EQUAL    TokenType = "LESS_EQUAL"
	// Literals.
	TokenType_IDENTIFIER TokenType = "IDENTIFIER"
	TokenType_STRING     TokenType = "STRING"
	TokenType_NUMBER     TokenType = "NUMBER"
	// Keywords.
	TokenType_AND    TokenType = "AND"
	TokenType_CLASS  TokenType = "CLASS"
	TokenType_ELSE   TokenType = "ELSE"
	TokenType_FALSE  TokenType = "FALSE"
	TokenType_FUN    TokenType = "FUN"
	TokenType_FOR    TokenType = "FOR"
	TokenType_IF     TokenType = "IF"
	TokenType_NIL    TokenType = "NIL"
	TokenType_OR     TokenType = "OR"
	TokenType_PRINT  TokenType = "PRINT"
	TokenType_RETURN TokenType = "RETURN"
	TokenType_SUPER  TokenType = "SUPER"
	TokenType_THIS   TokenType = "THIS"
	TokenType_TRUE   TokenType = "TRUE"
	TokenType_VAR    TokenType = "VAR"
	TokenType_WHILE  TokenType = "WHILE"
	TokenType_EOF    TokenType = "EOF"
)

type Token struct {
	t       TokenType
	lexeme  string
	literal any
	line    int
}

func NewToken(t TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		t:       t,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %s", t.t, t.lexeme, t.literal)
}

type Scanner struct {
	lox     *Lox
	source  string
	tokens  []*Token
	start   int
	current int
	line    int
}

func NewScanner(lox *Lox, source string) *Scanner {
	return &Scanner{
		lox:     lox,
		source:  source,
		tokens:  []*Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) scanTokens() ([]*Token, error) {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		if err := s.scanToken(); err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, NewToken(TokenType_EOF, "", nil, s.line))
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(TokenType_LEFT_PAREN)
	case ')':
		s.addToken(TokenType_RIGHT_PAREN)
	case '{':
		s.addToken(TokenType_LEFT_BRACE)
	case '}':
		s.addToken(TokenType_RIGHT_BRACE)
	case ',':
		s.addToken(TokenType_COMMA)
	case '.':
		s.addToken(TokenType_DOT)
	case '-':
		s.addToken(TokenType_MINUS)
	case '+':
		s.addToken(TokenType_PLUS)
	case ';':
		s.addToken(TokenType_SEMICOLON)
	case '*':
		s.addToken(TokenType_STAR)
	case '!':
		s.addToken(tern(s.match('='), TokenType_BANG_EQUAL, TokenType_BANG))
	case '=':
		s.addToken(tern(s.match('='), TokenType_EQUAL_EQUAL, TokenType_EQUAL))
	case '<':
		s.addToken(tern(s.match('='), TokenType_LESS_EQUAL, TokenType_LESS))
	case '>':
		s.addToken(tern(s.match('='), TokenType_GREATER_EQUAL, TokenType_GREATER))
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TokenType_SLASH)
		}
	case ' ':
		// Ignore whitespace.
	case '\r':
		// Ignore whitespace.
	case '\t':
		// Ignore whitespace.
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(c) {
			if err := s.number(); err != nil {
				return err
			}
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.lox.reportErr(s.line, "Unexpected character.")
		}
	}
	return nil
}

func tern[T any](c bool, t T, f T) T {
	if c {
		return t
	}
	return f
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) advance() byte {
	ret := s.source[s.current]
	s.current++
	return ret
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenL(t, nil)
}

func (s *Scanner) addTokenL(t TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(t, text, literal, s.line))
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.lox.reportErr(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addTokenL(TokenType_STRING, value)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() error {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	f, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return err
	}
	s.addTokenL(TokenType_NUMBER, f)
	return nil
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	t, ok := keywords[text]
	if !ok {
		t = TokenType_IDENTIFIER
	}
	s.addToken(t)
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
