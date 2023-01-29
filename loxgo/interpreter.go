package main

import (
	"fmt"
)

var _ = (VisitorExpr)(&Interpreter{})
var _ = (VisitorStmt)(&Interpreter{})

type Interpreter struct {
	lox *Lox
}

func (itrp *Interpreter) interpret(stmts []*Stmt) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			itrp.lox.runtimeError(r, Token{})
		}
	}()

	for _, stmt := range stmts {
		itrp.execute(stmt)
	}
	return err
}

func (itrp *Interpreter) execute(stmt *Stmt) {
	stmt.accept(itrp)
}

func (itrp *Interpreter) evaluate(expr *Expr) any {
	return expr.accept(itrp)
}

func (itrp *Interpreter) VisitLiteral(expr *Literal) any {
	return expr.Value
}
func (itrp *Interpreter) VisitGrouping(expr *Grouping) any {
	return itrp.evaluate(&expr.Expression)
}

func (itrp *Interpreter) VisitUnary(expr *Unary) any {
	right := itrp.evaluate(&expr.Right)

	switch expr.Operator.t {
	case TokenType_BANG:
		return !isTruthy(right)
	case TokenType_MINUS:
		rf, ok := right.(float64)
		if !ok {
			panic("failed to cast float in unary interpreter")
		}
		return -1.0 * rf
	}

	panic("unreachable")
}

func (itrp *Interpreter) VisitBinary(expr *Binary) any {
	left := itrp.evaluate(&expr.Left)
	right := itrp.evaluate(&expr.Right)

	switch expr.Operator.t {
	case TokenType_GREATER:
		fs, err := toFloats([]any{left, right})
		if err != nil {
			panic(err)
		}
		return fs[0] > fs[1]
	case TokenType_GREATER_EQUAL:
		fs, err := toFloats([]any{left, right})
		if err != nil {
			panic(err)
		}
		return fs[0] >= fs[1]
	case TokenType_LESS:
		fs, err := toFloats([]any{left, right})
		if err != nil {
			panic(err)
		}
		return fs[0] < fs[1]
	case TokenType_LESS_EQUAL:
		fs, err := toFloats([]any{left, right})
		if err != nil {
			panic(err)
		}
		return fs[0] <= fs[1]
	case TokenType_BANG_EQUAL:
		return !isEqual(left, right)
	case TokenType_EQUAL_EQUAL:
		return isEqual(left, right)
	case TokenType_MINUS:
		fs, err := toFloats([]any{left, right})
		if err != nil {
			panic(err)
		}
		return fs[0] - fs[1]
	case TokenType_PLUS:
		lr := []any{left, right}
		fs, err := toFloats(lr)
		if err == nil {
			return fs[0] + fs[1]
		}
		ss, err := toStrs(lr)
		if err == nil {
			return ss[0] + ss[1]
		}
	case TokenType_SLASH:
		fs, err := toFloats([]any{left, right})
		if err != nil {
			panic(err)
		}
		return fs[0] / fs[1]
	case TokenType_STAR:
		fs, err := toFloats([]any{left, right})
		if err != nil {
			panic(err)
		}
		return fs[0] * fs[1]
	}

	panic("unreachable")
}

func stringify(v any) string {
	return fmt.Sprintf("%v", v)
}

func isEqual(a, b any) bool {
	return a == b
}

func toFloats(vs []any) ([]float64, error) {
	return toTs(vs, make([]float64, len(vs)))
}
func toStrs(vs []any) ([]string, error) {
	return toTs(vs, make([]string, len(vs)))
}
func toTs[T any](src []any, dst []T) ([]T, error) {
	if len(src) != len(dst) {
		return nil, fmt.Errorf("src and dst must be equal in len")
	}
	for i := range src {
		f, ok := src[i].(T)
		if !ok {
			return nil, fmt.Errorf("failed to cast")
		}
		dst[i] = f
	}
	return dst, nil
}
func isTruthy(v any) bool {
	if isNil(v) {
		return true
	}
	bool, ok := v.(bool)
	if ok {
		return bool
	}
	return false
}
func isNil(v any) bool {
	return v == nil
}

func (itrp *Interpreter) VisitExpression(stmt *Expression) any {
	itrp.evaluate(&stmt.Expression)
	return nil
}
func (itrp *Interpreter) VisitPrint(stmt *Print) any {
	v := itrp.evaluate(&stmt.Expression)
	fmt.Println(stringify(v))
	return nil
}
