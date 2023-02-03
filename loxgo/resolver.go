package main

var (
	_ (VisitorExpr) = (*Resolver)(nil)
	_ (VisitorStmt) = (*Resolver)(nil)
)

type Resolver struct {
	lox       *Lox
	itrp      *Interpreter
	scopes    []map[string]bool
	currentFn FunctionType
}

type FunctionType string

const (
	FunctionType_NONE     = "NONE"
	FunctionType_FUNCTION = "FUNCTION"
)

func NewResolver(lox *Lox, itrp *Interpreter) *Resolver {
	return &Resolver{
		lox:       lox,
		itrp:      itrp,
		scopes:    []map[string]bool{},
		currentFn: FunctionType_NONE,
	}
}

func (r *Resolver) VisitBinary(expr *Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}
func (r *Resolver) VisitGrouping(expr *Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}
func (r *Resolver) VisitCall(expr *Call) any {
	r.resolveExpr(expr.Callee)
	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}
	return nil
}
func (r *Resolver) VisitLiteral(expr *Literal) any {
	return nil
}
func (r *Resolver) VisitUnary(expr *Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}
func (r *Resolver) VisitLogical(expr *Logical) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}
func (r *Resolver) VisitVariable(expr *Variable) any {
	if len(r.scopes) > 0 {
		last := r.scopes[len(r.scopes)-1]
		v, ok := last[expr.Name.lexeme]
		if ok && !v {
			r.lox.error(expr.Name,
				"Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(Expr{Variable: expr}, expr.Name)
	return nil
}
func (r *Resolver) VisitAssign(expr *Assign) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(Expr{Assign: expr}, expr.Name)
	return nil
}
func (r *Resolver) VisitExpression(stmt *Expression) any {
	r.resolveExpr(stmt.Expression)
	return nil
}
func (r *Resolver) VisitIf(stmt *If) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Then)
	if stmt.Else != nil {
		r.resolveStmt(stmt.Else)
	}
	return nil
}
func (r *Resolver) VisitFunction(stmt *Function) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FunctionType_FUNCTION)
	return nil
}
func (r *Resolver) VisitReturn(stmt *Return) any {
	if r.currentFn == FunctionType_NONE {
		r.lox.error(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
	return nil
}
func (r *Resolver) VisitPrint(stmt *Print) any {
	r.resolveExpr(stmt.Expression)
	return nil
}
func (r *Resolver) VisitVar(stmt *Var) any {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}
func (r *Resolver) VisitWhile(stmt *While) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}
func (r *Resolver) VisitBlock(stmt *Block) any {
	r.beginScope()
	r.resolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitClass(stmt *Class) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) resolveStmts(statements []*Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(stmt *Stmt) {
	stmt.accept(r)
}

func (r *Resolver) resolveExpr(expr *Expr) {
	expr.accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}
func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}
func (r *Resolver) declare(name *Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]

	if _, ok := scope[name.lexeme]; ok {
		r.lox.error(name,
			"Already a variable with this name in this scope.")
	}

	scope[name.lexeme] = false
}
func (r *Resolver) define(name *Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.lexeme]; ok {
			r.itrp.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}
func (r *Resolver) resolveFunction(fn *Function, ft FunctionType) {
	enclosingFn := r.currentFn
	r.currentFn = ft

	r.beginScope()
	for _, param := range fn.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(fn.Body)
	r.endScope()
	r.currentFn = enclosingFn
}
