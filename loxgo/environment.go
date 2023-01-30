package main

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    map[string]any{},
	}
}

func NewEnvironmentFrom(env *Environment) *Environment {
	return &Environment{
		enclosing: env,
		values:    map[string]any{},
	}
}

func (e *Environment) define(name string, v any) {
	e.values[name] = v
}
func (e *Environment) assign(name Token, v any) {
	_, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			e.enclosing.assign(name, v)
			return
		}
		panic("Undefined variable '" + name.lexeme + "'.")
	}
	e.values[name.lexeme] = v
}
func (e *Environment) get(name Token) any {
	v, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}
		panic("Undefined variable '" + name.lexeme + "'.")
	}
	return v
}
