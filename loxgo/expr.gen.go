package main

type Expr struct {
	Binary   *Binary
	Grouping *Grouping
	Literal  *Literal
	Unary    *Unary
}
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}
type Grouping struct {
	Expression Expr
}
type Literal struct {
	Value any
}
type Unary struct {
	Operator Token
	Right    Expr
}
type VisitorExpr interface {
	VisitBinary(expr *Binary) any
	VisitGrouping(expr *Grouping) any
	VisitLiteral(expr *Literal) any
	VisitUnary(expr *Unary) any
}

func (e *Expr) accept(v VisitorExpr) any {
	if e.Binary != nil {
		return e.Binary.accept(v)
	}
	if e.Grouping != nil {
		return e.Grouping.accept(v)
	}
	if e.Literal != nil {
		return e.Literal.accept(v)
	}
	if e.Unary != nil {
		return e.Unary.accept(v)
	}
	return nil
}
func (e *Binary) accept(visitor VisitorExpr) any {
	return visitor.VisitBinary(e)
}
func (e *Grouping) accept(visitor VisitorExpr) any {
	return visitor.VisitGrouping(e)
}
func (e *Literal) accept(visitor VisitorExpr) any {
	return visitor.VisitLiteral(e)
}
func (e *Unary) accept(visitor VisitorExpr) any {
	return visitor.VisitUnary(e)
}
