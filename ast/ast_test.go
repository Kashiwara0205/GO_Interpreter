package ast

import(
	"../token"
	"testing"
)

// let myvar = anoutherVar;
func Testing(t *testing.T){
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myvar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anoutherVar"},
					Value: "anoutherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anoutherVar;"{
		t.Errorf("prgoram String() wrong. got = %q", program.String())
	}
}