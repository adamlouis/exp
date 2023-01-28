package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: genast <output directory>")
		os.Exit(64)
	}
	outputDir := os.Args[1]

	if err := defineAst(outputDir, "Expr", []string{
		"Binary   : Expr left, Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal  : Object value",
		"Unary    : Token operator, Expr right",
	}); err != nil {
		fmt.Println(err.Error())
		os.Exit(64)
	}
}

func writeln(w io.Writer, s string) error {
	_, err := w.Write([]byte(s + "\n"))
	return err
}

func defineAst(outputDir string, baseName string, types []string) error {
	path := outputDir + "/" + baseName + ".java"

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = writeln(f, "package com.craftinginterpreters.lox;"); err != nil {
		return err
	}
	if err = writeln(f, ""); err != nil {
		return err
	}
	if err = writeln(f, "import java.util.List;"); err != nil {
		return err
	}
	if err = writeln(f, ""); err != nil {
		return err
	}
	if err = writeln(f, "abstract class "+baseName+" {"); err != nil {
		return err
	}

	// The AST classes.
	for _, t := range types {
		className := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		if err := defineTypeW(f, baseName, className, fields); err != nil {
			return err
		}
	}

	if err = writeln(f, "}"); err != nil {
		return err
	}

	return nil
}

func defineTypeW(w io.Writer, baseName string, className string, fieldList string) error {
	if err := writeln(w, "  static class "+className+" extends "+baseName+" {"); err != nil {
		return err
	}

	// Constructor.
	if err := writeln(w, "    "+className+"("+fieldList+") {"); err != nil {
		return err
	}

	// Store parameters in fields.
	fields := strings.Split(fieldList, ",")
	for _, field := range fields {
		name := strings.Split(field, " ")[1]
		if err := writeln(w, "      this."+name+" = "+name+";"); err != nil {
			return err
		}
	}

	if err := writeln(w, "    }"); err != nil {
		return err
	}

	// Fields.
	if err := writeln(w, ""); err != nil {
		return err
	}
	for _, field := range fields {
		if err := writeln(w, "    final "+field+";"); err != nil {
			return err
		}
	}

	if err := writeln(w, "  }"); err != nil {
		return err
	}
	return nil
}
