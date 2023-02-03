// DO NOT EDIT - generated code!
package main

type Expr struct {
	Binary   *Binary
	Grouping *Grouping
	Call     *Call
	Get      *Get
	Set      *Set
	Literal  *Literal
	Unary    *Unary
	Logical  *Logical
	Variable *Variable
	Assign   *Assign
}
type Binary struct {
	Left     *Expr
	Operator *Token
	Right    *Expr
}
type Grouping struct {
	Expression *Expr
}
type Call struct {
	Callee    *Expr
	Paren     *Token
	Arguments []*Expr
}
type Get struct {
	Object *Expr
	Name   *Token
}
type Set struct {
	Object *Expr
	Name   *Token
	Value  *Expr
}
type Literal struct {
	Value any
}
type Unary struct {
	Operator *Token
	Right    *Expr
}
type Logical struct {
	Left     *Expr
	Operator *Token
	Right    *Expr
}
type Variable struct {
	Name *Token
}
type Assign struct {
	Name  *Token
	Value *Expr
}
type VisitorExpr interface {
	VisitBinary(expr *Binary) any
	VisitGrouping(expr *Grouping) any
	VisitCall(expr *Call) any
	VisitGet(expr *Get) any
	VisitSet(expr *Set) any
	VisitLiteral(expr *Literal) any
	VisitUnary(expr *Unary) any
	VisitLogical(expr *Logical) any
	VisitVariable(expr *Variable) any
	VisitAssign(expr *Assign) any
}

func (e *Expr) accept(v VisitorExpr) any {
	if e.Binary != nil {
		return e.Binary.accept(v)
	}
	if e.Grouping != nil {
		return e.Grouping.accept(v)
	}
	if e.Call != nil {
		return e.Call.accept(v)
	}
	if e.Get != nil {
		return e.Get.accept(v)
	}
	if e.Set != nil {
		return e.Set.accept(v)
	}
	if e.Literal != nil {
		return e.Literal.accept(v)
	}
	if e.Unary != nil {
		return e.Unary.accept(v)
	}
	if e.Logical != nil {
		return e.Logical.accept(v)
	}
	if e.Variable != nil {
		return e.Variable.accept(v)
	}
	if e.Assign != nil {
		return e.Assign.accept(v)
	}
	return nil
}
func (e *Binary) accept(visitor VisitorExpr) any {
	return visitor.VisitBinary(e)
}
func (e *Grouping) accept(visitor VisitorExpr) any {
	return visitor.VisitGrouping(e)
}
func (e *Call) accept(visitor VisitorExpr) any {
	return visitor.VisitCall(e)
}
func (e *Get) accept(visitor VisitorExpr) any {
	return visitor.VisitGet(e)
}
func (e *Set) accept(visitor VisitorExpr) any {
	return visitor.VisitSet(e)
}
func (e *Literal) accept(visitor VisitorExpr) any {
	return visitor.VisitLiteral(e)
}
func (e *Unary) accept(visitor VisitorExpr) any {
	return visitor.VisitUnary(e)
}
func (e *Logical) accept(visitor VisitorExpr) any {
	return visitor.VisitLogical(e)
}
func (e *Variable) accept(visitor VisitorExpr) any {
	return visitor.VisitVariable(e)
}
func (e *Assign) accept(visitor VisitorExpr) any {
	return visitor.VisitAssign(e)
}
