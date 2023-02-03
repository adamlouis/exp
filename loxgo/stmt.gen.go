// DO NOT EDIT - generated code!
package main

type Stmt struct {
	Expression *Expression
	If         *If
	Function   *Function
	Return     *Return
	Print      *Print
	Var        *Var
	While      *While
	Block      *Block
	Class      *Class
}
type Expression struct {
	Expression *Expr
}
type If struct {
	Condition *Expr
	Then      *Stmt
	Else      *Stmt
}
type Function struct {
	Name   *Token
	Params []*Token
	Body   []*Stmt
}
type Return struct {
	Keyword *Token
	Value   *Expr
}
type Print struct {
	Expression *Expr
}
type Var struct {
	Name        *Token
	Initializer *Expr
}
type While struct {
	Condition *Expr
	Body      *Stmt
}
type Block struct {
	Statements []*Stmt
}
type Class struct {
	Name    *Token
	Methods []*Stmt
}
type VisitorStmt interface {
	VisitExpression(expr *Expression) any
	VisitIf(expr *If) any
	VisitFunction(expr *Function) any
	VisitReturn(expr *Return) any
	VisitPrint(expr *Print) any
	VisitVar(expr *Var) any
	VisitWhile(expr *While) any
	VisitBlock(expr *Block) any
	VisitClass(expr *Class) any
}

func (e *Stmt) accept(v VisitorStmt) any {
	if e.Expression != nil {
		return e.Expression.accept(v)
	}
	if e.If != nil {
		return e.If.accept(v)
	}
	if e.Function != nil {
		return e.Function.accept(v)
	}
	if e.Return != nil {
		return e.Return.accept(v)
	}
	if e.Print != nil {
		return e.Print.accept(v)
	}
	if e.Var != nil {
		return e.Var.accept(v)
	}
	if e.While != nil {
		return e.While.accept(v)
	}
	if e.Block != nil {
		return e.Block.accept(v)
	}
	if e.Class != nil {
		return e.Class.accept(v)
	}
	return nil
}
func (e *Expression) accept(visitor VisitorStmt) any {
	return visitor.VisitExpression(e)
}
func (e *If) accept(visitor VisitorStmt) any {
	return visitor.VisitIf(e)
}
func (e *Function) accept(visitor VisitorStmt) any {
	return visitor.VisitFunction(e)
}
func (e *Return) accept(visitor VisitorStmt) any {
	return visitor.VisitReturn(e)
}
func (e *Print) accept(visitor VisitorStmt) any {
	return visitor.VisitPrint(e)
}
func (e *Var) accept(visitor VisitorStmt) any {
	return visitor.VisitVar(e)
}
func (e *While) accept(visitor VisitorStmt) any {
	return visitor.VisitWhile(e)
}
func (e *Block) accept(visitor VisitorStmt) any {
	return visitor.VisitBlock(e)
}
func (e *Class) accept(visitor VisitorStmt) any {
	return visitor.VisitClass(e)
}
