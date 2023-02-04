package main

var _ Callable = (*LoxClass)(nil)

type LoxClass struct {
	name       string
	superClass *LoxClass
	methods    map[string]*LoxFunction
}

func NewLoxClass(name string, superClass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name, superClass, methods}
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
	v, ok := lc.methods[name]
	if ok {
		return v
	}

	if lc.superClass != nil {
		return lc.superClass.findMethod(name)
	}

	return nil
}
