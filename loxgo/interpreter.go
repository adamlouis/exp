package main

import (
	"fmt"
)

var _ = (VisitorExpr)(&Interpreter{})
var _ = (VisitorStmt)(&Interpreter{})

type Interpreter struct {
	lox     *Lox
	env     *Environment
	globals *Environment
	locals  map[Expr]int
}

func NewInterpreter(lox *Lox) *Interpreter {
	globals := NewEnvironment()

	globals.define("clock", &Clock{})

	return &Interpreter{
		lox:     lox,
		env:     globals,
		globals: globals,
		locals:  map[Expr]int{},
	}
}

func (itrp *Interpreter) interpret(stmts []*Stmt) (err error) {
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

func (itrp *Interpreter) VisitVariable(expr *Variable) any {
	return itrp.lookUpVariable(expr.Name, Expr{Variable: expr})
}

func (itrp *Interpreter) lookUpVariable(name *Token, expr Expr) any {
	distance, ok := itrp.locals[expr]
	if ok {
		return itrp.env.getAt(distance, name.lexeme)
	} else {
		return itrp.globals.get(name)
	}
}

func (itrp *Interpreter) VisitCall(expr *Call) any {
	callee := itrp.evaluate(expr.Callee)

	arguments := []any{}
	for _, argument := range expr.Arguments {
		arguments = append(arguments, itrp.evaluate(argument))
	}

	fn, ok := callee.(Callable)
	if !ok {
		lc, ok := callee.(LoxClass)
		if !ok {
			panic("can only call functions and classes")
		}
		fn = &lc
	}

	if len(arguments) != fn.Arity() {
		panic(fmt.Sprintf("expected %d arguments but got %d", fn.Arity(), len(arguments)))
	}

	return fn.Call(itrp, arguments)
}

func (itrp *Interpreter) VisitLogical(expr *Logical) any {
	left := itrp.evaluate(expr.Left)

	if expr.Operator.t == TokenType_OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}
	return itrp.evaluate(expr.Right)
}
func (itrp *Interpreter) VisitSet(expr *Set) any {
	object := itrp.evaluate(expr.Object)

	loxi, ok := object.(*LoxInstance)
	if !ok {
		panic(expr.Name.lexeme + ": only instances have fields.")
	}

	value := itrp.evaluate(expr.Value)
	loxi.Set(expr.Name, value)
	return value
}
func (itrp *Interpreter) VisitAssign(expr *Assign) any {
	value := itrp.evaluate(expr.Value)
	super := Expr{Assign: expr}
	distance, ok := itrp.locals[super]
	if ok {
		itrp.env.assignAt(distance, expr.Name, value)
	} else {
		itrp.globals.assign(expr.Name, value)
	}

	return value
}
func (itrp *Interpreter) VisitLiteral(expr *Literal) any {
	return expr.Value
}
func (itrp *Interpreter) VisitGrouping(expr *Grouping) any {
	return itrp.evaluate(expr.Expression)
}

func (itrp *Interpreter) VisitUnary(expr *Unary) any {
	right := itrp.evaluate(expr.Right)

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

func (itrp *Interpreter) VisitGet(expr *Get) any {
	object := itrp.evaluate(expr.Object)

	li, ok := object.(*LoxInstance)
	if ok {
		return li.Get(expr.Name)
	}

	panic("Only instances have properties.")
}

func (itrp *Interpreter) VisitBinary(expr *Binary) any {
	left := itrp.evaluate(expr.Left)
	right := itrp.evaluate(expr.Right)

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
	itrp.evaluate(stmt.Expression)
	return nil
}
func (itrp *Interpreter) VisitPrint(stmt *Print) any {
	v := itrp.evaluate(stmt.Expression)
	fmt.Println(stringify(v))
	return nil
}
func (itrp *Interpreter) VisitVar(stmt *Var) any {
	var value any
	if stmt.Initializer != nil {
		value = itrp.evaluate(stmt.Initializer)
	}

	itrp.env.define(stmt.Name.lexeme, value)
	return nil
}
func (itrp *Interpreter) VisitBlock(stmt *Block) any {
	itrp.executeBlock(stmt.Statements, NewEnvironmentFrom(itrp.env))
	return nil
}

func (itrp *Interpreter) executeBlock(statements []*Stmt, env *Environment) {
	previous := itrp.env
	defer func() {
		itrp.env = previous
	}()

	itrp.env = env
	for _, statement := range statements {
		itrp.execute(statement)
	}
}

func (itrp *Interpreter) VisitIf(stmt *If) any {
	if isTruthy(itrp.evaluate(stmt.Condition)) {
		itrp.execute(stmt.Then)
	} else if stmt.Else != nil {
		itrp.execute(stmt.Else)
	}
	return nil
}

func (itrp *Interpreter) VisitWhile(stmt *While) any {
	for isTruthy(itrp.evaluate(stmt.Condition)) {
		itrp.execute(stmt.Body)
	}
	return nil
}

func (itrp *Interpreter) VisitFunction(stmt *Function) any {
	f := NewLoxFunction(stmt, itrp.env)
	itrp.env.define(stmt.Name.lexeme, f)
	return nil
}

func (itrp *Interpreter) VisitReturn(stmt *Return) any {
	var value any

	if stmt.Value != nil {
		value = itrp.evaluate(stmt.Value)
	}

	panic(ReturnException{Value: value})
}

func (itrp *Interpreter) resolve(expr Expr, depth int) {
	itrp.locals[expr] = depth
}

func (itrp *Interpreter) VisitClass(stmt *Class) any {
	itrp.env.define(stmt.Name.lexeme, nil)

	methods := map[string]*LoxFunction{}
	for _, method := range stmt.Methods {
		fn := method.Function
		lfn := NewLoxFunction(fn, itrp.env)
		methods[fn.Name.lexeme] = lfn
	}

	klass := NewLoxClass(stmt.Name.lexeme, methods)
	itrp.env.assign(stmt.Name, klass)
	return nil
}

type ReturnException struct {
	Value any
}
