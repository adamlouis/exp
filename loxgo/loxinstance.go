package main

type LoxInstance struct {
	*LoxClass
	fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class, map[string]any{}}
}

func (li LoxInstance) String() string {
	return li.name + " instance"
}

func (li *LoxInstance) Get(name *Token) any {
	if v, ok := li.fields[name.lexeme]; ok {
		return v
	}

	method := li.LoxClass.findMethod(name.lexeme)
	if method != nil {
		return method.bind(li)
	}

	panic("Undefined property '" + name.lexeme + "'.")
}
func (li *LoxInstance) Set(name *Token, value any) {
	li.fields[name.lexeme] = value
}
