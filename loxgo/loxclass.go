package main

var _ Callable = (*LoxClass)(nil)

type LoxClass struct {
	name string
}

func (lc LoxClass) String() string {
	return lc.name
}

func (lc *LoxClass) Call(itrp *Interpreter, arguments []any) any {
	instance := LoxInstance{lc}
	return instance
}

func (lc *LoxClass) Arity() int {
	return 0
}
