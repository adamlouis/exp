package main

type Expr struct {
	Binary   *Binary
	Unary    *Unary
	Literal  *Literal
	Grouping *Grouping
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}
type Unary struct {
	Operator Token
	Right    Expr
}
type Literal struct {
	Value any
}
type Grouping struct {
	Expression Expr
}
