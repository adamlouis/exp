package main

var (
	_ (VisitorExpr) = (*Resolver)(nil)
	_ (VisitorStmt) = (*Resolver)(nil)
)

type Resolver struct {
	lox          *Lox
	itrp         *Interpreter
	scopes       []map[string]bool
	currentFn    FunctionType
	currentClass ClassType
}

type FunctionType string

const (
	FunctionType_NONE        FunctionType = "NONE"
	FunctionType_FUNCTION    FunctionType = "FUNCTION"
	FunctionType_METHOD      FunctionType = "METHOD"
	FunctionType_INITIALIZER FunctionType = "INITIALIZER"
)

type ClassType string

const (
	ClassType_NONE     ClassType = "NONE"
	ClassType_CLASS    ClassType = "CLASS"
	ClassType_SUBCLASS ClassType = "SUBCLASS"
)

func NewResolver(lox *Lox, itrp *Interpreter) *Resolver {
	return &Resolver{
		lox:          lox,
		itrp:         itrp,
		scopes:       []map[string]bool{},
		currentFn:    FunctionType_NONE,
		currentClass: ClassType_NONE,
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
func (r *Resolver) VisitSet(expr *Set) any {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil
}
func (r *Resolver) VisitSuper(expr *Super) any {
	if r.currentClass == ClassType_NONE {
		r.lox.error(expr.Keyword, "Can't use 'super' outside of a class.")
	} else if r.currentClass != ClassType_SUBCLASS {
		r.lox.error(expr.Keyword, "Can't use 'super' in a class with no superclass.")
	}

	r.resolveLocal(Expr{Super: expr}, expr.Keyword)
	return nil
}
func (r *Resolver) VisitThis(expr *This) any {
	if r.currentClass == ClassType_NONE {
		r.lox.error(expr.Keyword, "Can't use 'this' outside of a class.")
	}

	r.resolveLocal(Expr{This: expr}, expr.Keyword)
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
func (r *Resolver) VisitGet(stmt *Get) any {
	r.resolveExpr(stmt.Object)
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
		if r.currentFn == FunctionType_INITIALIZER {
			r.lox.error(stmt.Keyword, "Can't return a value from an initializer.")
		}

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
	enclosing := r.currentClass
	r.currentClass = ClassType_CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.SuperClass != nil && stmt.Name.lexeme == stmt.SuperClass.Name.lexeme {
		r.lox.error(stmt.SuperClass.Name, "A class can't inherit from itself.")
	}

	if stmt.SuperClass != nil {
		r.currentClass = ClassType_SUBCLASS
		r.resolveExpr(&Expr{Variable: stmt.SuperClass})
	}

	if stmt.SuperClass != nil {
		r.beginScope()
		last := r.scopes[len(r.scopes)-1]
		last["super"] = true
	}

	r.beginScope()

	last := r.scopes[len(r.scopes)-1]
	last["this"] = true

	for _, method := range stmt.Methods {
		declaration := FunctionType_METHOD
		if method.Function.Name.lexeme == "init" {
			declaration = FunctionType_INITIALIZER
		}
		r.resolveFunction(method.Function, declaration)
	}
	r.endScope()

	if stmt.SuperClass != nil {
		r.endScope()
	}

	r.currentClass = enclosing
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
