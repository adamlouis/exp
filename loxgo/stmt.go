package main

type Stmt struct {
	Expression *Expression
	Print      *Print
}
type Expression struct {
	Expression Expr
}
type Print struct {
	Expression Expr
}
type VisitorStmt interface {
	VisitExpression(expr *Expression) any
	VisitPrint(expr *Print) any
}

func (e *Stmt) accept(v VisitorStmt) any {
	if e.Expression != nil {
		return e.Expression.accept(v)
	}
	if e.Print != nil {
		return e.Print.accept(v)
	}
	return nil
}
func (e *Expression) accept(visitor VisitorStmt) any {
	return visitor.VisitExpression(e)
}
func (e *Print) accept(visitor VisitorStmt) any {
	return visitor.VisitPrint(e)
}
