package main

import "time"

type Callable interface {
	Call(itrp *Interpreter, arguments []any) any
	Arity() int
	String() string
}

var _ Callable = (*Clock)(nil)

type Clock struct{}

func (c *Clock) Call(itrp *Interpreter, arguments []any) any {
	return float64(time.Now().Second())
}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) String() string {
	return "<native fn>"
}
