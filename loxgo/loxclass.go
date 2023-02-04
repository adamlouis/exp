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

	initializer := lc.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).Call(itrp, arguments)
	}

	return instance
}

func (lc *LoxClass) Arity() int {
	initializer := lc.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (lc *LoxClass) findMethod(name string) *LoxFunction {
	return lc.methods[name]
}
