package ast

import (
	"testing"
	"github.com/BentleyOph/monke/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token : token.Token{Type: token.IDENT, Literal: "myVar"},
					Value : "myVar",
				},
				Value: &Identifier{
					Token : token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},

			},

		},
	}
	if program.String() != "let myVar = anotherVar;"{
		t.Errorf("program.String() wrong. got=%q",program.String())
	}
}

// The test above creates a new Program node with a single LetStatement node inside it. The LetStatement node has a Name and a Value, both of which are Identifier nodes. The test then calls the String method on the Program node and checks the output. The expected output is let myVar = anotherVar;.

