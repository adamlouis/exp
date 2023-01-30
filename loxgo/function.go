package main

var _ Callable = (*LoxFunction)(nil)

type LoxFunction struct {
	decl    *Function
	closure *Environment
}

func NewLoxFunction(decl *Function, closure *Environment) *LoxFunction {
	return &LoxFunction{
		decl:    decl,
		closure: closure,
	}
}

func (f *LoxFunction) Call(itrp *Interpreter, arguments []any) (ret any) {
	env := NewEnvironmentFrom(f.closure)

	defer func() {
		if r := recover(); r != nil {
			retex, ok := r.(ReturnException)
			if ok {
				ret = retex.Value
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
	return ret
}

func (c *LoxFunction) Arity() int {
	return len(c.decl.Params)
}

func (c *LoxFunction) String() string {
	return "<fn " + c.decl.Name.lexeme + ">"
}
