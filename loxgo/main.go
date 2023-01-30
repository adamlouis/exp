package main

import (
	"fmt"
	"os"
)

func main() {
	// took shortcuts to get java patterns into go
	// take a pass at end to write idomatic go
	l := &Lox{
		interpreter: NewInterpreter(nil),
	}
	l.interpreter.lox = l

	switch len(os.Args) {
	case 2:
		if err := l.runFile(os.Args[1]); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	case 1:
		if err := l.runPrompt(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	default:
		panic("usage: loxgo [script]")
	}
}
