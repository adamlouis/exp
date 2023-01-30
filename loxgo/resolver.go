package main

var (
	_ (VisitorExpr) = (*Resolver)(nil)
	_ (VisitorStmt) = (*Resolver)(nil)
)

type Resolver struct{}
