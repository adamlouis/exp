package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type ExprType struct {
	Name   string
	Fields []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: genast <output directory>")
		os.Exit(64)
	}
	outputDir := os.Args[1]

	// TODO(adam): fix bad dupe names .. Var / Exp

	if err := genAST(outputDir, "Expr", []ExprType{
		{"Binary", []string{
			"Left *Expr",
			"Operator *Token",
			"Right *Expr",
		}},
		{"Grouping", []string{
			"Expression *Expr",
		}},
		{"Call", []string{
			"Callee *Expr ",
			"Paren *Token ",
			"Arguments []*Expr",
		}},
		{"Get", []string{
			"Object *Expr",
			"Name *Token",
		}},
		{"Set", []string{
			"Object *Expr",
			"Name *Token",
			"Value *Expr",
		}},
		{"Literal", []string{
			"Value any",
		}},
		{"Unary", []string{
			"Operator *Token",
			"Right *Expr",
		}},
		{"This", []string{
			"Keyword *Token",
		}},
		{"Logical", []string{
			"Left *Expr",
			"Operator *Token",
			"Right *Expr",
		}},
		{"Variable", []string{
			"Name *Token",
		}},
		{"Assign", []string{
			"Name *Token",
			"Value *Expr",
		}},
	}); err != nil {
		fmt.Println(err.Error())
		os.Exit(64)
	}

	if err := genAST(outputDir, "Stmt", []ExprType{
		{"Expression", []string{
			"Expression *Expr",
		}},
		{"If", []string{
			"Condition *Expr",
			"Then *Stmt",
			"Else *Stmt",
		}},
		{"Function", []string{
			"Name *Token",
			"Params []*Token",
			"Body []*Stmt",
		}},
		{"Return", []string{
			"Keyword *Token",
			"Value *Expr",
		}},
		{"Print", []string{
			"Expression *Expr",
		}},
		{"Var", []string{
			"Name *Token",
			"Initializer *Expr",
		}},
		{"While", []string{
			"Condition *Expr",
			"Body *Stmt",
		}},
		{"Block", []string{
			"Statements []*Stmt",
		}},
		{"Class", []string{
			"Name *Token",
			"Methods []*Stmt",
		}},
	}); err != nil {
		fmt.Println(err.Error())
		os.Exit(64)
	}
}

func writeln(w io.Writer, s string) error {
	_, err := w.Write([]byte(s + "\n"))
	return err
}

func genAST(outputDir string, baseName string, types []ExprType) error {
	path := outputDir + "/" + strings.ToLower(baseName) + ".gen.go"

	_ = os.RemoveAll(path)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = writeln(f, "// DO NOT EDIT - generated code!"); err != nil {
		return err
	}
	if err = writeln(f, "package main"); err != nil {
		return err
	}

	if err = writeln(f, "type "+baseName+" struct {"); err != nil {
		return err
	}

	for _, t := range types {
		if err = writeln(f, fmt.Sprintf("	%s *%s", t.Name, t.Name)); err != nil {
			return err
		}
	}
	if err = writeln(f, "}"); err != nil {
		return err
	}

	for _, t := range types {
		if err = writeln(f, "type "+t.Name+" struct {"); err != nil {
			return err
		}
		for _, field := range t.Fields {
			if err = writeln(f, field); err != nil {
				return err
			}
		}
		if err = writeln(f, "}"); err != nil {
			return err
		}
	}

	return genVisitor(f, baseName, types)
}

func genVisitor(w io.Writer, baseName string, types []ExprType) error {
	iname := "Visitor" + baseName

	if err := writeln(w, `type `+iname+` interface {`); err != nil {
		return err
	}
	for _, t := range types {
		if err := writeln(w, "	Visit"+t.Name+"(expr *"+t.Name+") any"); err != nil {
			return err
		}
	}
	if err := writeln(w, `}`); err != nil {
		return err
	}

	if err := writeln(w, "func (e *"+baseName+") accept(v "+iname+") any {"); err != nil {
		return err
	}
	for _, t := range types {
		if err := writeln(w, "if e."+t.Name+" != nil {"); err != nil {
			return err
		}
		if err := writeln(w, "	return e."+t.Name+".accept(v)"); err != nil {
			return err
		}
		if err := writeln(w, "}"); err != nil {
			return err
		}
	}
	if err := writeln(w, "	return nil"); err != nil {
		return err
	}
	if err := writeln(w, "}"); err != nil {
		return err
	}

	for _, t := range types {
		if err := writeln(w, "func (e *"+t.Name+") accept(visitor "+iname+") any {"); err != nil {
			return err
		}

		if err := writeln(w, "	return visitor.Visit"+t.Name+"(e)"); err != nil {
			return err
		}

		if err := writeln(w, "}"); err != nil {
			return err
		}
	}
	return nil
}
