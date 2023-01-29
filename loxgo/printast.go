package main

import "fmt"

// TODO: use functional style, drop the use of `any`

type ASTPrinter struct {
}

func (p *ASTPrinter) print(expr Expr) string {
	return fmt.Sprintf("%v", expr.accept(p))
}

func (p *ASTPrinter) parenthesize(name string, exprs ...Expr) string {
	ret := ""

	ret += "(" + name
	for _, expr := range exprs {
		ret += " "
		v := expr.accept(p)
		ret += fmt.Sprintf("%v", v)
	}
	ret += ")"

	return ret
}

func (p *ASTPrinter) VisitBinary(expr *Binary) any {
	return p.parenthesize(
		expr.Operator.lexeme,
		expr.Left,
		expr.Right,
	)
}

func (p *ASTPrinter) VisitGrouping(expr *Grouping) any {
	return p.parenthesize(
		"group",
		expr.Expression,
	)
}

func (p *ASTPrinter) VisitLiteral(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p *ASTPrinter) VisitUnary(expr *Unary) any {
	return p.parenthesize(expr.Operator.lexeme, expr.Right)
}
