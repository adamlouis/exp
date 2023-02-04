package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLox(t *testing.T) {

	prog := `class Bacon {
  eat() {
    print "Crunch crunch crunch!";
  }
}

print "After";
Bacon().eat(); // Prints "Crunch crunch crunch!".
print "Before";
`
	l := &Lox{
		interpreter: NewInterpreter(nil),
	}
	l.interpreter.lox = l

	err := l.run(prog)
	require.Nil(t, err)
}
