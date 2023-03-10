package main

var _ Callable = (*LoxFunction)(nil)

type LoxFunction struct {
	decl          *Function
	closure       *Environment
	isInitializer bool
}

func NewLoxFunction(decl *Function, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		decl:          decl,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *LoxFunction) Call(itrp *Interpreter, arguments []any) (ret any) {
	env := NewEnvironmentFrom(f.closure)

	defer func() {
		if r := recover(); r != nil {
			retex, ok := r.(ReturnException)
			if ok {

				if f.isInitializer {
					f.closure.getAt(0, "this")
				} else {
					ret = retex.Value
				}
			}
		}
	}()

	for i := 0; i < len(f.decl.Params); i++ {
		env.define(
			f.decl.Params[i].lexeme,
			arguments[i],
		)
	}

	itrp.executeBlock(f.decl.Body, env)

	if f.isInitializer {
		ret = f.closure.getAt(0, "this")
	}

	return ret
}

func (c *LoxFunction) Arity() int {
	return len(c.decl.Params)
}

func (c *LoxFunction) String() string {
	return "<fn " + c.decl.Name.lexeme + ">"
}

func (c *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	env := NewEnvironmentFrom(c.closure)
	env.define("this", instance)
	return NewLoxFunction(c.decl, env, c.isInitializer)
}
