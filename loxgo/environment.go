package main

type Environment struct {
	values map[string]any
}

func (e *Environment) define(name string, v any) {
	e.values[name] = v
}
func (e *Environment) assign(name Token, v any) {
	_, ok := e.values[name.lexeme]
	if !ok {
		panic("Undefined variable '" + name.lexeme + "'.")
	}
	e.values[name.lexeme] = v
}
func (e *Environment) get(name Token) any {
	v, ok := e.values[name.lexeme]
	if !ok {
		panic("Undefined variable '" + name.lexeme + "'.")
	}
	return v
}
