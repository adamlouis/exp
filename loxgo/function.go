package main

var _ Callable = (*LoxFunction)(nil)

type LoxFunction struct {
	decl *Function
}

func NewLoxFunction(decl *Function) *LoxFunction {
	return &LoxFunction{
		decl: decl,
	}
}

func (f *LoxFunction) Call(itrp *Interpreter, arguments []any) any {
	env := NewEnvironmentFrom(itrp.globals)

	for i := 0; i < len(f.decl.Params); i++ {
		env.define(
			f.decl.Params[i].lexeme,
			arguments[i],
		)
	}

	itrp.executeBlock(f.decl.Body, env)
	return nil
}

func (c *LoxFunction) Arity() int {
	return len(c.decl.Params)
}

func (c *LoxFunction) String() string {
	return "<fn " + c.decl.Name.lexeme + ">"
}
