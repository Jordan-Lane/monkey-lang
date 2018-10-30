package ast

import (
	"monkeylang/token"
	"testing"
)

func TestString(t *testing.T) {
	// TODO: remove this once the AST is built properly
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "y"},
					Value: "y",
				},
			},
		},
	}

	expectedString := "let x = y;"
	if program.String() != expectedString {
		t.Errorf("program.String() invalid. Expected: %q, Got:%q", expectedString, program.String())
	}
}
