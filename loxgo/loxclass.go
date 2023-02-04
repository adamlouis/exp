package main

var _ Callable = (*LoxClass)(nil)

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name, methods}
}

func (lc LoxClass) String() string {
	return lc.name
}

func (lc *LoxClass) Call(itrp *Interpreter, arguments []any) any {
	instance := NewLoxInstance(lc)
	return instance
}

func (lc *LoxClass) Arity() int {
	return 0
}

func (lc *LoxClass) findMethod(name string) *LoxFunction {
	return lc.methods[name]
}
