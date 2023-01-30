// DO NOT EDIT - generated code!
package main

type Stmt struct {
	Expression *Expression
	Print      *Print
	Var        *Var
	Block      *Block
}
type Expression struct {
	Expression Expr
}
type Print struct {
	Expression Expr
}
type Var struct {
	Name        Token
	Initializer *Expr
}
type Block struct {
	Statements []*Stmt
}
type VisitorStmt interface {
	VisitExpression(expr *Expression) any
	VisitPrint(expr *Print) any
	VisitVar(expr *Var) any
	VisitBlock(expr *Block) any
}

func (e *Stmt) accept(v VisitorStmt) any {
	if e.Expression != nil {
		return e.Expression.accept(v)
	}
	if e.Print != nil {
		return e.Print.accept(v)
	}
	if e.Var != nil {
		return e.Var.accept(v)
	}
	if e.Block != nil {
		return e.Block.accept(v)
	}
	return nil
}
func (e *Expression) accept(visitor VisitorStmt) any {
	return visitor.VisitExpression(e)
}
func (e *Print) accept(visitor VisitorStmt) any {
	return visitor.VisitPrint(e)
}
func (e *Var) accept(visitor VisitorStmt) any {
	return visitor.VisitVar(e)
}
func (e *Block) accept(visitor VisitorStmt) any {
	return visitor.VisitBlock(e)
}
